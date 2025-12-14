#!/bin/bash

echo "🔄 Migrating database to support searchable encryption..."

DB_PATH="${1:-./memories.db}"

if [ ! -f "$DB_PATH" ]; then
    echo "❌ Database not found: $DB_PATH"
    exit 1
fi

echo "📁 Database: $DB_PATH"

# Backup first
cp "$DB_PATH" "${DB_PATH}.backup-$(date +%Y%m%d-%H%M%S)"
echo "✅ Backup created"

# Add search_content column
sqlite3 "$DB_PATH" "ALTER TABLE memories ADD COLUMN search_content TEXT;" 2>/dev/null
echo "✅ Added search_content column"

# Populate search_content with existing data
sqlite3 "$DB_PATH" "UPDATE memories SET search_content = text_content WHERE search_content IS NULL;"
echo "✅ Populated search_content from text_content"

# Drop old triggers
sqlite3 "$DB_PATH" "DROP TRIGGER IF EXISTS memories_ai;"
sqlite3 "$DB_PATH" "DROP TRIGGER IF EXISTS memories_au;"
sqlite3 "$DB_PATH" "DROP TRIGGER IF EXISTS memories_ad;"
echo "✅ Dropped old triggers"

# Rebuild FTS5 index with plain text
sqlite3 "$DB_PATH" "DELETE FROM memories_fts;"
sqlite3 "$DB_PATH" "INSERT INTO memories_fts(rowid, text_content, tags) SELECT id, search_content, tags FROM memories;"
echo "✅ Rebuilt FTS5 index with searchable content"

echo ""
echo "🎉 Migration complete!"
echo ""
echo "Now run: ./memory-bot"
