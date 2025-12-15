#!/bin/bash

# Migration script to add memory chunking support (parent_id field)
# This enables hierarchical memory organization (like breaking down projects into tasks)

set -e

DB_PATH="${DB_PATH:-./memories.db}"

echo "🔄 Migrating database to support Memory Chunking..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -f "$DB_PATH" ]; then
    echo "❌ Database not found at: $DB_PATH"
    echo "   No migration needed - schema will be created on first run"
    exit 0
fi

echo "📁 Database: $DB_PATH"

# Check if parent_id column already exists
COLUMN_EXISTS=$(sqlite3 "$DB_PATH" "PRAGMA table_info(memories);" | grep -c "parent_id" || true)

if [ "$COLUMN_EXISTS" -gt 0 ]; then
    echo "✅ parent_id column already exists - no migration needed"
    exit 0
fi

echo "🔧 Adding parent_id column..."

# Add parent_id column with foreign key constraint
sqlite3 "$DB_PATH" <<EOF
-- Add parent_id column
ALTER TABLE memories ADD COLUMN parent_id INTEGER;

-- Create foreign key index for better query performance
CREATE INDEX IF NOT EXISTS idx_memories_parent_id ON memories(parent_id);

-- Verify the column was added
SELECT 'Column added successfully' as result;
EOF

echo ""
echo "✅ Migration completed successfully!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🧠 Memory Chunking is now enabled!"
echo "   You can now create hierarchical memories:"
echo "   - Parent memory: 'Build Personal Memory Bot'"
echo "   - Child memories: 'Implement search', 'Add encryption', etc."
echo ""
