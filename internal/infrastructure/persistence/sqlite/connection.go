package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Connection manages SQLite database connection
type Connection struct {
	DB *sql.DB
}

// NewConnection creates a new SQLite connection with optimizations
func NewConnection(dbPath string) (*Connection, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys and WAL mode for better performance
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if _, err := conn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	connection := &Connection{DB: conn}

	if err := connection.initSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("SQLite database initialized successfully")
	return connection, nil
}

// initSchema creates the necessary tables and indexes
func (c *Connection) initSchema() error {
	// Main memories table with biological fields
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS memories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		chat_id INTEGER NOT NULL,
		text_content TEXT NOT NULL,
		search_content TEXT,
		tags TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_reviewed DATETIME,
		review_count INTEGER DEFAULT 0,
		last_consolidated DATETIME DEFAULT CURRENT_TIMESTAMP,
		priority_score REAL DEFAULT 0.0,
		emotional_weight REAL DEFAULT 0.0,
		time_of_day TEXT DEFAULT '',
		day_of_week TEXT DEFAULT '',
		chat_source TEXT DEFAULT 'Telegram'
	);`

	if _, err := c.DB.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create memories table: %w", err)
	}

	// Create composite index for fast user-based queries
	createIndexSQL := `
	CREATE INDEX IF NOT EXISTS idx_user_time 
	ON memories (user_id, created_at DESC);`

	if _, err := c.DB.Exec(createIndexSQL); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Create biological memory system indexes
	biologicalIndexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_memories_time_of_day ON memories(time_of_day);`,
		`CREATE INDEX IF NOT EXISTS idx_memories_day_of_week ON memories(day_of_week);`,
		`CREATE INDEX IF NOT EXISTS idx_memories_emotional_weight ON memories(emotional_weight DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_memories_priority_score ON memories(priority_score DESC);`,
	}

	for _, indexSQL := range biologicalIndexes {
		if _, err := c.DB.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create biological index: %w", err)
		}
	}

	// Create FTS5 virtual table for full-text search
	// Uses search_content which contains plain text (not encrypted)
	createFTSSQL := `
	CREATE VIRTUAL TABLE IF NOT EXISTS memories_fts USING fts5(
		text_content,
		tags,
		content='memories',
		content_rowid='id',
		tokenize='porter unicode61'
	);`

	if _, err := c.DB.Exec(createFTSSQL); err != nil {
		return fmt.Errorf("failed to create FTS5 table: %w", err)
	}

	// Create triggers to keep FTS5 table in sync
	if err := c.createTriggers(); err != nil {
		return err
	}

	return nil
}

// createTriggers creates database triggers for FTS5 synchronization
func (c *Connection) createTriggers() error {
	triggers := []string{
		`CREATE TRIGGER IF NOT EXISTS memories_ai AFTER INSERT ON memories BEGIN
			INSERT INTO memories_fts(rowid, text_content, tags)
			VALUES (new.id, COALESCE(new.search_content, new.text_content), new.tags);
		END;`,
		`CREATE TRIGGER IF NOT EXISTS memories_au AFTER UPDATE ON memories BEGIN
			UPDATE memories_fts 
			SET text_content = COALESCE(new.search_content, new.text_content), tags = new.tags
			WHERE rowid = new.id;
		END;`,
		`CREATE TRIGGER IF NOT EXISTS memories_ad AFTER DELETE ON memories BEGIN
			DELETE FROM memories_fts WHERE rowid = old.id;
		END;`,
	}

	for _, trigger := range triggers {
		if _, err := c.DB.Exec(trigger); err != nil {
			return fmt.Errorf("failed to create trigger: %w", err)
		}
	}

	return nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
