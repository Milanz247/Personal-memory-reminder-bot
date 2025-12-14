package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Build tag to enable FTS5 support
// To build: go build -tags "fts5"

// Memory represents a stored memory with context
type Memory struct {
	ID           int
	UserID       int64
	ChatID       int64
	TextContent  string
	Tags         string
	CreatedAt    time.Time
	LastReviewed *time.Time
	ReviewCount  int
	Rank         float64 // For FTS5 ranking results
}

// Database handles all database operations
type Database struct {
	conn *sql.DB
}

// NewDatabase creates a new database connection and initializes tables
func NewDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better performance
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if _, err := conn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	db := &Database{conn: conn}

	if err := db.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// initTables creates the necessary tables and indexes
func (db *Database) initTables() error {
	// Main memories table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS memories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		chat_id INTEGER NOT NULL,
		text_content TEXT NOT NULL,
		tags TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_reviewed DATETIME,
		review_count INTEGER DEFAULT 0
	);`

	if _, err := db.conn.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create memories table: %w", err)
	}

	// Create composite index for fast user-based queries
	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_user_time 
	ON memories (user_id, created_at DESC);`

	if _, err := db.conn.Exec(createIndexSQL); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Create FTS5 virtual table for full-text search
	createFTSSQL := `
	CREATE VIRTUAL TABLE IF NOT EXISTS memories_fts USING fts5(
		text_content,
		tags,
		content='memories',
		content_rowid='id',
		tokenize='porter unicode61'
	);`

	if _, err := db.conn.Exec(createFTSSQL); err != nil {
		return fmt.Errorf("failed to create FTS5 table: %w", err)
	}

	// Create triggers to keep FTS5 table in sync
	createTriggerInsertSQL := `
	CREATE TRIGGER IF NOT EXISTS memories_ai AFTER INSERT ON memories BEGIN
		INSERT INTO memories_fts(rowid, text_content, tags)
		VALUES (new.id, new.text_content, new.tags);
	END;`

	if _, err := db.conn.Exec(createTriggerInsertSQL); err != nil {
		return fmt.Errorf("failed to create insert trigger: %w", err)
	}

	createTriggerUpdateSQL := `
	CREATE TRIGGER IF NOT EXISTS memories_au AFTER UPDATE ON memories BEGIN
		UPDATE memories_fts 
		SET text_content = new.text_content, tags = new.tags
		WHERE rowid = new.id;
	END;`

	if _, err := db.conn.Exec(createTriggerUpdateSQL); err != nil {
		return fmt.Errorf("failed to create update trigger: %w", err)
	}

	createTriggerDeleteSQL := `
	CREATE TRIGGER IF NOT EXISTS memories_ad AFTER DELETE ON memories BEGIN
		DELETE FROM memories_fts WHERE rowid = old.id;
	END;`

	if _, err := db.conn.Exec(createTriggerDeleteSQL); err != nil {
		return fmt.Errorf("failed to create delete trigger: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.conn.Close()
}
