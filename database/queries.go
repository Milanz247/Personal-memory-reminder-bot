package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// SaveMemory stores a new memory with context in the database
func (db *Database) SaveMemory(userID, chatID int64, content, tags string) (int64, error) {
	stmt, err := db.conn.Prepare(`
		INSERT INTO memories (user_id, chat_id, text_content, tags, created_at)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID, chatID, content, tags, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to save memory: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	log.Printf("Memory saved: ID=%d, UserID=%d", id, userID)
	return id, nil
}

// SearchAndRankMemories performs FTS5 search with ranking and pagination
func (db *Database) SearchAndRankMemories(userID int64, keyword string, limit, offset int) ([]Memory, error) {
	// FTS5 MATCH query with ranking
	query := `
		SELECT 
			m.id,
			m.user_id,
			m.chat_id,
			m.text_content,
			m.tags,
			m.created_at,
			m.last_reviewed,
			m.review_count,
			memories_fts.rank as rank
		FROM 
			memories AS m
		JOIN 
			memories_fts ON m.id = memories_fts.rowid
		WHERE 
			m.user_id = ? AND 
			memories_fts MATCH ?
		ORDER BY 
			rank
		LIMIT ? OFFSET ?
	`

	// Prepare the search term for FTS5 with wildcard and proximity matching
	searchTerm := prepareFTS5SearchTerm(keyword)

	rows, err := db.conn.Query(query, userID, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search memories: %w", err)
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var m Memory
		var lastReviewed sql.NullTime

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.ChatID,
			&m.TextContent,
			&m.Tags,
			&m.CreatedAt,
			&lastReviewed,
			&m.ReviewCount,
			&m.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if lastReviewed.Valid {
			m.LastReviewed = &lastReviewed.Time
		}

		memories = append(memories, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	log.Printf("Found %d memories for user %d with keyword '%s' (FTS term: '%s')", len(memories), userID, keyword, searchTerm)
	return memories, nil
}

// SmartSearch performs intelligent search with fallback strategies
// 1. Try proximity search first (NEAR)
// 2. If no results, try wildcard AND search
// 3. If still no results, try OR search
func (db *Database) SmartSearch(userID int64, keyword string, limit int) ([]Memory, error) {
	log.Printf("SmartSearch: Starting search for userID=%d, keyword='%s'", userID, keyword)

	// Try primary search with wildcards
	memories, err := db.SearchAndRankMemories(userID, keyword, limit, 0)
	if err != nil {
		log.Printf("SmartSearch: Primary search error: %v", err)
	}
	if err == nil && len(memories) > 0 {
		log.Printf("SmartSearch: Found %d results with primary search", len(memories))
		return memories, nil
	}

	// Fallback 1: Try simple wildcard AND search
	words := strings.Fields(strings.TrimSpace(keyword))
	if len(words) > 1 {
		andTerms := make([]string, len(words))
		for i, word := range words {
			andTerms[i] = word + "*"
		}
		fallbackTerm := strings.Join(andTerms, " ")

		log.Printf("SmartSearch: Trying AND fallback with term: %s", fallbackTerm)
		memories, err = db.searchWithCustomTerm(userID, fallbackTerm, limit)
		if err != nil {
			log.Printf("SmartSearch: AND fallback error: %v", err)
		}
		if err == nil && len(memories) > 0 {
			log.Printf("SmartSearch: Found %d results with AND fallback", len(memories))
			return memories, nil
		}
	}

	// Fallback 2: Try OR search (broader results)
	if len(words) > 1 {
		orTerms := make([]string, len(words))
		for i, word := range words {
			orTerms[i] = word + "*"
		}
		fallbackTerm := strings.Join(orTerms, " OR ")

		log.Printf("SmartSearch: Trying OR fallback with term: %s", fallbackTerm)
		memories, err = db.searchWithCustomTerm(userID, fallbackTerm, limit)
		if err != nil {
			log.Printf("SmartSearch: OR fallback error: %v", err)
		}
		if err == nil && len(memories) > 0 {
			log.Printf("SmartSearch: Found %d results with OR fallback", len(memories))
			return memories, nil
		}
	}

	// No results found with any strategy
	log.Printf("SmartSearch: No results found for keyword '%s'", keyword)
	return []Memory{}, nil
}

// searchWithCustomTerm performs search with a custom FTS5 term
func (db *Database) searchWithCustomTerm(userID int64, ftsTerm string, limit int) ([]Memory, error) {
	query := `
		SELECT 
			m.id,
			m.user_id,
			m.chat_id,
			m.text_content,
			m.tags,
			m.created_at,
			m.last_reviewed,
			m.review_count,
			memories_fts.rank as rank
		FROM 
			memories AS m
		JOIN 
			memories_fts ON m.id = memories_fts.rowid
		WHERE 
			m.user_id = ? AND 
			memories_fts MATCH ?
		ORDER BY 
			rank
		LIMIT ?
	`

	rows, err := db.conn.Query(query, userID, ftsTerm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search with custom term: %w", err)
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var m Memory
		var lastReviewed sql.NullTime

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.ChatID,
			&m.TextContent,
			&m.Tags,
			&m.CreatedAt,
			&lastReviewed,
			&m.ReviewCount,
			&m.Rank,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if lastReviewed.Valid {
			m.LastReviewed = &lastReviewed.Time
		}

		memories = append(memories, m)
	}

	return memories, rows.Err()
}

// GetRecentMemories retrieves the most recent memories for a user
func (db *Database) GetRecentMemories(userID int64, limit int) ([]Memory, error) {
	query := `
		SELECT 
			id, user_id, chat_id, text_content, tags, 
			created_at, last_reviewed, review_count
		FROM memories
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := db.conn.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent memories: %w", err)
	}
	defer rows.Close()

	return scanMemories(rows)
}

// GetMemoriesForReview retrieves memories that need to be reviewed
// Based on spaced repetition intervals
func (db *Database) GetMemoriesForReview(reviewIntervals []int) ([]Memory, error) {
	// Build the query dynamically based on review intervals
	conditions := make([]string, 0, len(reviewIntervals))
	args := make([]interface{}, 0, len(reviewIntervals))

	for _, days := range reviewIntervals {
		conditions = append(conditions, "(julianday('now') - julianday(COALESCE(last_reviewed, created_at)) >= ?)")
		args = append(args, days)
	}

	query := fmt.Sprintf(`
		SELECT 
			id, user_id, chat_id, text_content, tags,
			created_at, last_reviewed, review_count
		FROM memories
		WHERE %s
		ORDER BY COALESCE(last_reviewed, created_at) ASC
		LIMIT 50
	`, strings.Join(conditions, " OR "))

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories for review: %w", err)
	}
	defer rows.Close()

	return scanMemories(rows)
}

// UpdateLastReviewed updates the last_reviewed timestamp and increments review_count
func (db *Database) UpdateLastReviewed(memoryID int) error {
	stmt, err := db.conn.Prepare(`
		UPDATE memories 
		SET last_reviewed = ?, review_count = review_count + 1
		WHERE id = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), memoryID)
	if err != nil {
		return fmt.Errorf("failed to update last_reviewed: %w", err)
	}

	return nil
}

// DeleteMemory deletes a memory by ID
func (db *Database) DeleteMemory(userID int64, memoryID int) error {
	stmt, err := db.conn.Prepare(`
		DELETE FROM memories 
		WHERE id = ? AND user_id = ?
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(memoryID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete memory: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("memory not found or access denied")
	}

	log.Printf("Memory deleted: ID=%d, UserID=%d", memoryID, userID)
	return nil
}

// GetMemoryCount returns the total number of memories for a user
func (db *Database) GetMemoryCount(userID int64) (int, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM memories WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get memory count: %w", err)
	}
	return count, nil
}

// Helper functions

// prepareFTS5SearchTerm prepares a search term for FTS5 MATCH with wildcard support
// Implements fragment matching - allows partial word matches using wildcards
// Example: "Tele bot" -> "Tele* AND bot*" (matches Telegram, Telegraph, etc.)
func prepareFTS5SearchTerm(keyword string) string {
	// Remove special characters and trim spaces
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return keyword
	}

	// Split into words
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

		// Check if word already has wildcard
		if !strings.HasSuffix(word, "*") {
			word = word + "*"
		}

		ftsTerms = append(ftsTerms, word)
	}

	// For multiple words, use simple AND with wildcards
	// FTS5 AND gives good relevance ranking with BM25
	if len(ftsTerms) > 1 {
		return strings.Join(ftsTerms, " ")
	}

	return ftsTerms[0]
}

// prepareFTS5SearchTermExact prepares exact phrase search (for advanced users)
func prepareFTS5SearchTermExact(keyword string) string {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return keyword
	}
	// Wrap in quotes for exact phrase matching
	return `"` + keyword + `"`
}

// scanMemories is a helper function to scan multiple memory rows
func scanMemories(rows *sql.Rows) ([]Memory, error) {
	var memories []Memory

	for rows.Next() {
		var m Memory
		var lastReviewed sql.NullTime

		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.ChatID,
			&m.TextContent,
			&m.Tags,
			&m.CreatedAt,
			&lastReviewed,
			&m.ReviewCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if lastReviewed.Valid {
			m.LastReviewed = &lastReviewed.Time
		}

		memories = append(memories, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return memories, nil
}
