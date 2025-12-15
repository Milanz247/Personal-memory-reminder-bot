#!/bin/bash

# Script to rebuild FTS5 table with new tokenizer
# This is needed when updating the tokenizer configuration

echo "🔄 Rebuilding FTS5 Index with Enhanced Tokenizer..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if database exists
if [ ! -f "memories.db" ]; then
    echo "❌ Database file 'memories.db' not found!"
    echo "   The database will be created automatically when you run the bot."
    exit 1
fi

# Backup the database
BACKUP_FILE="memories-backup-$(date +%Y%m%d-%H%M%S).db"
echo "📦 Creating backup: $BACKUP_FILE"
cp memories.db "$BACKUP_FILE"

# Rebuild FTS5 table
echo "🔨 Rebuilding FTS5 virtual table..."
sqlite3 memories.db << 'EOF'
-- Drop existing FTS5 table and triggers
DROP TRIGGER IF EXISTS memories_ai;
DROP TRIGGER IF EXISTS memories_au;
DROP TRIGGER IF EXISTS memories_ad;
DROP TABLE IF EXISTS memories_fts;

-- Recreate FTS5 table with enhanced tokenizer
CREATE VIRTUAL TABLE memories_fts USING fts5(
    text_content,
    tags,
    content='memories',
    content_rowid='id',
    tokenize="unicode61 tokenchars '.'"
);

-- Recreate triggers
CREATE TRIGGER memories_ai AFTER INSERT ON memories BEGIN
    INSERT INTO memories_fts(rowid, text_content, tags)
    VALUES (new.id, COALESCE(new.search_content, new.text_content), new.tags);
END;

CREATE TRIGGER memories_au AFTER UPDATE ON memories BEGIN
    UPDATE memories_fts 
    SET text_content = COALESCE(new.search_content, new.text_content), tags = new.tags
    WHERE rowid = new.id;
END;

CREATE TRIGGER memories_ad AFTER DELETE ON memories BEGIN
    DELETE FROM memories_fts WHERE rowid = old.id;
END;

-- Rebuild FTS5 index from existing data
INSERT INTO memories_fts(rowid, text_content, tags)
SELECT id, COALESCE(search_content, text_content), tags FROM memories;

-- Show statistics
SELECT '✅ FTS5 rebuild complete!';
SELECT 'Total memories indexed: ' || COUNT(*) FROM memories_fts;
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "✅ FTS5 index rebuilt successfully!"
    echo ""
    echo "🎯 Enhanced Search Features:"
    echo "   • Version numbers: v1.24.27, 24.27"
    echo "   • IP addresses: 192.168.1.1"
    echo "   • Dates: 2024-12-15"
    echo "   • Port numbers: port8889, admin8889"
    echo "   • Mixed alphanumeric: abc123, test2024"
    echo ""
    echo "📦 Backup saved: $BACKUP_FILE"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
else
    echo ""
    echo "❌ Error rebuilding FTS5 index!"
    echo "🔄 Restoring backup..."
    cp "$BACKUP_FILE" memories.db
    echo "✅ Database restored from backup"
    exit 1
fi
