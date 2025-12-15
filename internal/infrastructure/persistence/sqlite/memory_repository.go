package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"memory-bot/internal/domain/entity"
	"memory-bot/internal/domain/repository"
	"memory-bot/pkg/encryption"
)

// MemoryRepository is the SQLite implementation of repository.MemoryRepository
type MemoryRepository struct {
	conn      *Connection
	encryptor *encryption.Encryptor
}

// NewMemoryRepository creates a new SQLite memory repository
func NewMemoryRepository(conn *Connection, encryptor *encryption.Encryptor) *MemoryRepository {
	return &MemoryRepository{
		conn:      conn,
		encryptor: encryptor,
	}
}

// Save stores a new memory in the database
func (r *MemoryRepository) Save(ctx context.Context, memory *entity.Memory) (int64, error) {
	if err := memory.Validate(); err != nil {
		return 0, err
	}

	// Encrypt content before storing
	encryptedContent, err := encryption.EncryptIfEnabled(r.encryptor, memory.Content)
	if err != nil {
		return 0, fmt.Errorf("failed to encrypt content: %w", err)
	}

	stmt, err := r.conn.DB.PrepareContext(ctx, `
		INSERT INTO memories (
			user_id, chat_id, text_content, search_content, tags, created_at,
			last_consolidated, priority_score, emotional_weight,
			time_of_day, day_of_week, chat_source, parent_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		memory.UserID,
		memory.ChatID,
		encryptedContent, // Encrypted version
		memory.Content,   // Plain text for FTS5 searching
		memory.GetTagsString(),
		memory.CreatedAt,
		memory.LastConsolidated,
		memory.PriorityScore,
		memory.EmotionalWeight,
		memory.TimeOfDay,
		memory.DayOfWeek,
		memory.ChatSource,
		memory.ParentID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to save memory: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	log.Printf("Memory saved: ID=%d, UserID=%d, EmotionalWeight=%.2f, Context=%s %s",
		id, memory.UserID, memory.EmotionalWeight, memory.DayOfWeek, memory.TimeOfDay)
	return id, nil
}

// FindByID retrieves a memory by its ID
func (r *MemoryRepository) FindByID(ctx context.Context, id int) (*entity.Memory, error) {
	query := `
		SELECT id, user_id, chat_id, text_content, tags, 
		       created_at, last_reviewed, review_count,
		       last_consolidated, priority_score, emotional_weight,
		       time_of_day, day_of_week, chat_source, parent_id
		FROM memories
		WHERE id = ?
	`

	var m entity.Memory
	var tags string
	var lastReviewed sql.NullTime
	var parentID sql.NullInt64

	err := r.conn.DB.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.UserID,
		&m.ChatID,
		&m.Content,
		&tags,
		&m.CreatedAt,
		&lastReviewed,
		&m.ReviewCount,
		&m.LastConsolidated,
		&m.PriorityScore,
		&m.EmotionalWeight,
		&m.TimeOfDay,
		&m.DayOfWeek,
		&m.ChatSource,
		&parentID,
	)

	if err == sql.ErrNoRows {
		return nil, entity.ErrMemoryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find memory: %w", err)
	}

	m.Tags = strings.Fields(tags)
	if lastReviewed.Valid {
		m.LastReviewed = &lastReviewed.Time
	}

	// Decrypt content after reading
	decryptedContent, err := encryption.DecryptIfEnabled(r.encryptor, m.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt content: %w", err)
	}
	m.Content = decryptedContent

	return &m, nil
}

// Search performs FTS5 search with ranking and optional contextual filtering
func (r *MemoryRepository) Search(ctx context.Context, userID int64, query string, opts repository.SearchOptions) ([]*entity.Memory, error) {
	searchTerm := prepareFTS5SearchTerm(query)

	// Build dynamic SQL query with advanced ranking
	// Ranking factors: BM25 score + emotional weight + priority score + recency
	sqlQuery := `
		SELECT 
			m.id,
			m.user_id,
			m.chat_id,
			m.text_content,
			m.tags,
			m.created_at,
			m.last_reviewed,
			m.review_count,
			m.emotional_weight,
			m.priority_score,
			memories_fts.rank as rank,
			(
				memories_fts.rank + 
				(m.emotional_weight * 2.0) + 
				(m.priority_score * 1.5) +
				(CASE 
					WHEN julianday('now') - julianday(m.created_at) < 7 THEN 1.0
					WHEN julianday('now') - julianday(m.created_at) < 30 THEN 0.5
					ELSE 0.0
				END)
			) as combined_rank
		FROM 
			memories AS m
		JOIN 
			memories_fts ON m.id = memories_fts.rowid
		WHERE 
			m.user_id = ? AND 
			memories_fts MATCH ?`

	args := []interface{}{userID, searchTerm}

	// Apply contextual filtering directly in SQL for better performance
	if opts.ContextFilter != nil {
		if opts.ContextFilter.TimeOfDay != "" {
			sqlQuery += " AND m.time_of_day = ?"
			args = append(args, opts.ContextFilter.TimeOfDay)
			log.Printf("Search: Applying TimeOfDay filter: %s", opts.ContextFilter.TimeOfDay)
		}
		if opts.ContextFilter.DayOfWeek != "" {
			sqlQuery += " AND m.day_of_week = ?"
			args = append(args, opts.ContextFilter.DayOfWeek)
			log.Printf("Search: Applying DayOfWeek filter: %s", opts.ContextFilter.DayOfWeek)
		}
	}

	// Order by combined ranking (BM25 + emotional + priority + recency)
	sqlQuery += ` ORDER BY combined_rank DESC LIMIT ? OFFSET ?`
	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.conn.DB.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %w", err)
	}
	defer rows.Close()

	memories := []*entity.Memory{}
	for rows.Next() {
		m, err := scanMemoryRow(rows)
		if err != nil {
			return nil, err
		}

		// Decrypt content after reading
		decryptedContent, err := encryption.DecryptIfEnabled(r.encryptor, m.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt content: %w", err)
		}
		m.Content = decryptedContent

		memories = append(memories, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Found %d memories for user %d with keyword '%s'", len(memories), userID, query)
	return memories, nil
}

// GetRecent retrieves the most recent memories for a user
func (r *MemoryRepository) GetRecent(ctx context.Context, userID int64, limit int) ([]*entity.Memory, error) {
	query := `
		SELECT 
			id, user_id, chat_id, text_content, tags, 
			created_at, last_reviewed, review_count, parent_id
		FROM memories
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.conn.DB.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent memories: %w", err)
	}
	defer rows.Close()

	memories, err := scanMemories(rows)
	if err != nil {
		return nil, err
	}

	// Decrypt content for each memory
	for _, m := range memories {
		decryptedContent, err := encryption.DecryptIfEnabled(r.encryptor, m.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt content: %w", err)
		}
		m.Content = decryptedContent
	}

	return memories, nil
}

// GetForReview retrieves memories that need review based on intervals
func (r *MemoryRepository) GetForReview(ctx context.Context, intervals []int) ([]*entity.Memory, error) {
	conditions := make([]string, 0, len(intervals))
	args := make([]interface{}, 0, len(intervals))

	for _, days := range intervals {
		conditions = append(conditions, "(julianday('now') - julianday(COALESCE(last_reviewed, created_at)) >= ?)")
		args = append(args, days)
	}

	query := fmt.Sprintf(`
		SELECT 
			id, user_id, chat_id, text_content, tags,
			created_at, last_reviewed, review_count, parent_id
		FROM memories
		WHERE %s
		ORDER BY COALESCE(last_reviewed, created_at) ASC
		LIMIT 50
	`, strings.Join(conditions, " OR "))

	rows, err := r.conn.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories for review: %w", err)
	}
	defer rows.Close()

	memories, err := scanMemories(rows)
	if err != nil {
		return nil, err
	}

	// Decrypt content for each memory
	for _, m := range memories {
		decryptedContent, err := encryption.DecryptIfEnabled(r.encryptor, m.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt content: %w", err)
		}
		m.Content = decryptedContent
	}

	return memories, nil
}

// Update updates an existing memory
func (r *MemoryRepository) Update(ctx context.Context, memory *entity.Memory) error {
	if err := memory.Validate(); err != nil {
		return err
	}

	stmt, err := r.conn.DB.PrepareContext(ctx, `
		UPDATE memories 
		SET last_reviewed = ?, review_count = ?
		WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, memory.LastReviewed, memory.ReviewCount, memory.ID)
	if err != nil {
		return fmt.Errorf("failed to update memory: %w", err)
	}

	return nil
}

// Delete removes a memory with authorization check
func (r *MemoryRepository) Delete(ctx context.Context, id int, userID int64) error {
	stmt, err := r.conn.DB.PrepareContext(ctx, `
		DELETE FROM memories 
		WHERE id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if affected == 0 {
		return entity.ErrUnauthorized
	}

	log.Printf("Memory deleted: ID=%d, UserID=%d", id, userID)
	return nil
}

// Count returns the total number of memories for a user
func (r *MemoryRepository) Count(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.conn.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM memories WHERE user_id = ?",
		userID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get memory count: %w", err)
	}
	return count, nil
}

// Helper functions

// scanMemoryRow scans a single memory row with rank
func scanMemoryRow(rows *sql.Rows) (*entity.Memory, error) {
	var m entity.Memory
	var tags string
	var lastReviewed sql.NullTime
	var emotionalWeight float64
	var priorityScore float64
	var rank float64
	var combinedRank float64

	err := rows.Scan(
		&m.ID,
		&m.UserID,
		&m.ChatID,
		&m.Content,
		&tags,
		&m.CreatedAt,
		&lastReviewed,
		&m.ReviewCount,
		&emotionalWeight,
		&priorityScore,
		&rank,
		&combinedRank,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	m.Tags = strings.Fields(tags)
	if lastReviewed.Valid {
		m.LastReviewed = &lastReviewed.Time
	}
	m.EmotionalWeight = emotionalWeight
	m.PriorityScore = priorityScore
	m.Rank = combinedRank // Use combined rank for display

	return &m, nil
}

// scanMemories scans multiple memory rows
func scanMemories(rows *sql.Rows) ([]*entity.Memory, error) {
	var memories []*entity.Memory

	for rows.Next() {
		var m entity.Memory
		var tags string
		var lastReviewed sql.NullTime
		var parentID sql.NullInt64

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.ChatID,
			&m.Content,
			&tags,
			&m.CreatedAt,
			&lastReviewed,
			&m.ReviewCount,
			&parentID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		m.Tags = strings.Fields(tags)
		if lastReviewed.Valid {
			m.LastReviewed = &lastReviewed.Time
		}
		if parentID.Valid {
			m.ParentID = &parentID.Int64
		}

		memories = append(memories, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return memories, nil
}

// prepareFTS5SearchTerm prepares a search term for FTS5 MATCH with wildcard support
func prepareFTS5SearchTerm(keyword string) string {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return keyword
	}

	// If query contains FTS5 operators, return as-is (already prepared by strategy)
	if strings.Contains(keyword, " OR ") || strings.Contains(keyword, "NEAR(") {
		return keyword
	}

	// Handle hashtag search (exact matching)
	if strings.HasPrefix(keyword, "#") {
		return keyword // Return hashtag as-is for exact tag matching
	}

	words := strings.Fields(keyword)
	if len(words) == 0 {
		return keyword
	}

	var ftsTerms []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		// Escape special FTS5 characters
		word = escapeFTS5SpecialChars(word)

		// Only add wildcard if not already present
		if !strings.HasSuffix(word, "*") {
			word = word + "*"
		}

		ftsTerms = append(ftsTerms, word)
	}

	if len(ftsTerms) > 1 {
		return strings.Join(ftsTerms, " ")
	}

	if len(ftsTerms) == 0 {
		return keyword
	}

	return ftsTerms[0]
}

// escapeFTS5SpecialChars escapes special characters in FTS5 queries
func escapeFTS5SpecialChars(word string) string {
	// Skip if word is already quoted or has wildcards
	if strings.HasPrefix(word, "\"") || strings.HasSuffix(word, "*") {
		return word
	}

	// Escape double quotes in the word
	word = strings.ReplaceAll(word, "\"", "\"\"")

	// If word contains special chars, wrap in quotes
	specialChars := []string{"(", ")", "^", "-", "+"}
	for _, char := range specialChars {
		if strings.Contains(word, char) {
			return "\"" + word + "\""
		}
	}

	return word
}

// GetFragileMemories retrieves recently created memories that need consolidation
// Biological principle: Memories created within the last 7 days with low review counts
func (r *MemoryRepository) GetFragileMemories(ctx context.Context) ([]*entity.Memory, error) {
	query := `
		SELECT 
			id, user_id, chat_id, text_content, tags,
			created_at, last_reviewed, review_count,
			last_consolidated, priority_score, emotional_weight,
			time_of_day, day_of_week, chat_source, parent_id
		FROM memories
		WHERE 
			julianday('now') - julianday(created_at) <= 7
			AND review_count < 2
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get fragile memories: %w", err)
	}
	defer rows.Close()

	memories, err := scanMemoriesWithBiologicalFields(rows)
	if err != nil {
		return nil, err
	}

	// Decrypt content for each memory
	for _, m := range memories {
		decryptedContent, err := encryption.DecryptIfEnabled(r.encryptor, m.Content)
		if err != nil {
			log.Printf("Warning: failed to decrypt memory %d: %v", m.ID, err)
			continue
		}
		m.Content = decryptedContent
	}

	log.Printf("Found %d fragile memories for consolidation", len(memories))
	return memories, nil
}

// UpdateConsolidation updates consolidation-related fields
func (r *MemoryRepository) UpdateConsolidation(ctx context.Context, memory *entity.Memory) error {
	stmt, err := r.conn.DB.PrepareContext(ctx, `
		UPDATE memories
		SET last_consolidated = ?, priority_score = ?
		WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare consolidation update: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, memory.LastConsolidated, memory.PriorityScore, memory.ID)
	if err != nil {
		return fmt.Errorf("failed to update consolidation: %w", err)
	}

	log.Printf("Updated consolidation for memory %d: PriorityScore=%.2f", memory.ID, memory.PriorityScore)
	return nil
}

// scanMemoriesWithBiologicalFields scans memory rows including biological fields
func scanMemoriesWithBiologicalFields(rows *sql.Rows) ([]*entity.Memory, error) {
	memories := []*entity.Memory{}

	for rows.Next() {
		var m entity.Memory
		var tags string
		var lastReviewed sql.NullTime
		var parentID sql.NullInt64

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.ChatID,
			&m.Content,
			&tags,
			&m.CreatedAt,
			&lastReviewed,
			&m.ReviewCount,
			&m.LastConsolidated,
			&m.PriorityScore,
			&m.EmotionalWeight,
			&m.TimeOfDay,
			&m.DayOfWeek,
			&m.ChatSource,
			&parentID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		m.Tags = strings.Fields(tags)
		if lastReviewed.Valid {
			m.LastReviewed = &lastReviewed.Time
		}
		if parentID.Valid {
			m.ParentID = &parentID.Int64
		}

		memories = append(memories, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return memories, nil
}
