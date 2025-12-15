-- Migration: Add Biologically-Inspired Memory Fields
-- Date: 2025-12-15
-- Description: Enhances the memory system with neuroscience-inspired features:
--   1. Emotional tagging (Amygdala function)
--   2. Contextual encoding (Hippocampus function)
--   3. Memory consolidation tracking (Sleep-based strengthening)

-- Add emotional weight field (0.0 to 1.0)
-- Simulates the Amygdala's role in emotional memory encoding
ALTER TABLE memories ADD COLUMN emotional_weight REAL DEFAULT 0.0;

-- Add consolidation tracking fields
-- Simulates sleep-based memory consolidation
ALTER TABLE memories ADD COLUMN last_consolidated DATETIME DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE memories ADD COLUMN priority_score REAL DEFAULT 0.0;

-- Add contextual metadata fields
-- Simulates Hippocampus encoding of time and context
ALTER TABLE memories ADD COLUMN time_of_day TEXT DEFAULT '';
ALTER TABLE memories ADD COLUMN day_of_week TEXT DEFAULT '';
ALTER TABLE memories ADD COLUMN chat_source TEXT DEFAULT 'Telegram';

-- Create indexes for efficient contextual searches
CREATE INDEX IF NOT EXISTS idx_memories_time_of_day ON memories(time_of_day);
CREATE INDEX IF NOT EXISTS idx_memories_day_of_week ON memories(day_of_week);
CREATE INDEX IF NOT EXISTS idx_memories_emotional_weight ON memories(emotional_weight DESC);
CREATE INDEX IF NOT EXISTS idx_memories_priority_score ON memories(priority_score DESC);

-- Create index for consolidation job queries
CREATE INDEX IF NOT EXISTS idx_memories_fragile ON memories(created_at, review_count)
WHERE julianday('now') - julianday(created_at) <= 7 AND review_count < 2;

-- Update existing memories with default contextual data
-- This ensures backward compatibility
UPDATE memories 
SET 
    last_consolidated = created_at,
    priority_score = 0.0,
    emotional_weight = 0.1,
    time_of_day = CASE 
        WHEN CAST(strftime('%H', created_at) AS INTEGER) BETWEEN 5 AND 11 THEN 'Morning'
        WHEN CAST(strftime('%H', created_at) AS INTEGER) BETWEEN 12 AND 16 THEN 'Afternoon'
        WHEN CAST(strftime('%H', created_at) AS INTEGER) BETWEEN 17 AND 20 THEN 'Evening'
        ELSE 'Night'
    END,
    day_of_week = CASE CAST(strftime('%w', created_at) AS INTEGER)
        WHEN 0 THEN 'Sunday'
        WHEN 1 THEN 'Monday'
        WHEN 2 THEN 'Tuesday'
        WHEN 3 THEN 'Wednesday'
        WHEN 4 THEN 'Thursday'
        WHEN 5 THEN 'Friday'
        WHEN 6 THEN 'Saturday'
    END,
    chat_source = 'Telegram'
WHERE emotional_weight IS NULL OR emotional_weight = 0.0;

-- Verify the migration
SELECT 
    'Migration completed successfully' AS status,
    COUNT(*) AS total_memories,
    COUNT(CASE WHEN emotional_weight > 0 THEN 1 END) AS emotional_memories,
    COUNT(CASE WHEN time_of_day != '' THEN 1 END) AS contextualized_memories
FROM memories;
