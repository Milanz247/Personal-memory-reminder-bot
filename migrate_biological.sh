#!/bin/bash

# Biologically-Inspired Memory System Migration
# This script adds neuroscience-based features to the memory bot database

set -e  # Exit on error

echo "🧠 Starting Biological Memory System Migration..."

# Get database path from environment or use default
DB_PATH="${DB_PATH:-./memories.db}"

if [ ! -f "$DB_PATH" ]; then
    echo "⚠️  Database not found at $DB_PATH"
    echo "   This is normal for new installations."
    echo "   The biological fields will be created automatically when you first run the bot."
    echo ""
    echo "   If you have an existing database, set DB_PATH:"
    echo "   export DB_PATH=/path/to/your/memories.db"
    echo "   ./migrate_biological.sh"
    exit 0
fi

echo "📁 Database: $DB_PATH"

# Backup database before migration
BACKUP_PATH="${DB_PATH}.backup.$(date +%Y%m%d_%H%M%S)"
echo "💾 Creating backup: $BACKUP_PATH"
cp "$DB_PATH" "$BACKUP_PATH"

# Check if columns already exist
EMOTIONAL_EXISTS=$(sqlite3 "$DB_PATH" "PRAGMA table_info(memories);" | grep -c "emotional_weight" || true)

if [ "$EMOTIONAL_EXISTS" -gt 0 ]; then
    echo "✅ Biological fields already exist! No migration needed."
    echo "   Backup saved at: $BACKUP_PATH"
    exit 0
fi

# Run migration
echo "🔄 Applying migration..."
sqlite3 "$DB_PATH" < migrations/add_biological_fields.sql

if [ $? -eq 0 ]; then
    echo "✅ Migration completed successfully!"
    echo ""
    echo "🎉 New biologically-inspired features enabled:"
    echo "  • 🧠 Emotional tagging (Amygdala function)"
    echo "  • 📍 Contextual encoding (Hippocampus function)"
    echo "  • 😴 Sleep-based consolidation tracking"
    echo ""
    echo "Backup saved at: $BACKUP_PATH"
else
    echo "❌ Migration failed!"
    echo "Restoring from backup..."
    cp "$BACKUP_PATH" "$DB_PATH"
    echo "Database restored"
    exit 1
fi
