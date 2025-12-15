# 🚀 Quick Implementation Guide

This guide shows how to integrate the biological memory system into your running bot.

## Step 1: Run Database Migration

```bash
# Navigate to project directory
cd /home/milanmadusanka/Projects/Personal-memory-reminder-bot

# Run migration script
./migrate_biological.sh
```

Expected output:
```
🧠 Starting Biological Memory System Migration...
📁 Database: ./memory_bot.db
💾 Creating backup: ./memory_bot.db.backup.20251215_143022
🔄 Applying migration...
✅ Migration completed successfully!
```

## Step 2: Update Main Application

Add these imports to `cmd/bot/main.go`:

```go
import (
    "memory-bot/internal/infrastructure/job"
    "memory-bot/internal/infrastructure/scheduler"
)
```

## Step 3: Initialize Consolidation Job

Add after repository initialization in `main()`:

```go
// Initialize daily consolidation job (simulates sleep-based memory strengthening)
log.Println("🌙 Initializing memory consolidation system...")
consolidationJob := job.NewDailyConsolidationJob(memoryRepo)

// Schedule to run at 2:00 AM every day
consolidationJob.ScheduleDaily(2, 0)

// Optional: Run once on startup for immediate consolidation
go func() {
    time.Sleep(5 * time.Second)  // Wait for bot to fully initialize
    if err := consolidationJob.RunNow(); err != nil {
        log.Printf("Initial consolidation run failed: %v", err)
    }
}()
```

## Step 4: Update Spaced Repetition Scheduler (Optional)

If you want to use the enhanced biological scheduler:

Replace the existing scheduler initialization with:

```go
// Create biological spaced repetition scheduler
bioScheduler := scheduler.NewBiologicalSpacedRepetition(config.ReviewIntervals)

// Use in your review logic
spacedRepScheduler := scheduler.NewSpacedRepetitionScheduler(
    botAPI,
    reviewMemoryUseCase,
    config.ReviewIntervals,
)
spacedRepScheduler.Start()
```

## Step 5: Test the System

### Test Emotional Tagging

Send these messages to your bot:

```
/save This is an amazing and wonderful breakthrough! I'm so excited!
```

Expected: High emotional weight (0.7-0.9)

```
/save Meeting scheduled for 3pm tomorrow.
```

Expected: Low emotional weight (0.1-0.2)

### Test Contextual Search

```
/save Had a great brainstorming session this morning

# Later or next day:
/search yesterday morning brainstorm
```

Should find your morning memory with contextual priority.

### Test Consolidation

Check logs for consolidation activity:
```bash
tail -f logs/bot.log | grep -i consolidation
```

You should see entries like:
```
2025-12-15 02:00:00 🌙 Starting daily memory consolidation...
2025-12-15 02:00:01 Found 5 fragile memories for consolidation
2025-12-15 02:00:02 Updated consolidation for memory 123: PriorityScore=0.50
```

## Step 6: Monitor Performance

### Check Memory Statistics

Add a new stats display to show biological features:

```go
// In your stats command handler
func handleStatsCommand(ctx context.Context, userID int64) string {
    stats, err := getStatsUseCase.Execute(ctx, userID)
    if err != nil {
        return "Error getting stats"
    }

    return fmt.Sprintf(`
📊 *Memory Statistics*

Total Memories: %d
Reviewed: %d (%.1f%%)
Average Review Count: %.1f

🧠 *Biological Features*
High Emotion Memories: %d
Consolidated Today: %d
Context-Tagged: %d
    `,
        stats.TotalMemories,
        stats.ReviewedMemories,
        stats.ReviewPercentage,
        stats.AvgReviewCount,
        stats.EmotionalMemories,
        stats.ConsolidatedToday,
        stats.ContextTagged,
    )
}
```

## Common Issues

### Issue 1: Migration fails with "no such table"

**Solution**: Ensure the database file exists and is accessible.
```bash
ls -la memory_bot.db
sqlite3 memory_bot.db ".tables"
```

### Issue 2: Consolidation job doesn't run

**Solution**: Check scheduler logs and ensure time is configured correctly.
```go
// Add debug logging
log.Printf("Next consolidation at %v", nextRunTime)
```

### Issue 3: Context not being captured

**Solution**: Verify SaveMemoryUseCase is using the new services.
```bash
# Check if new fields are being saved
sqlite3 memory_bot.db "SELECT id, emotional_weight, time_of_day FROM memories LIMIT 5;"
```

## Verification Checklist

- [ ] Migration completed successfully
- [ ] Database backup created
- [ ] New fields visible in database schema
- [ ] Consolidation job scheduled and running
- [ ] Emotional weights being calculated (check DB)
- [ ] Context fields being populated (check DB)
- [ ] Search working with contextual queries
- [ ] Logs showing consolidation activity

## Performance Tips

1. **Index Usage**: Verify indexes are created:
   ```bash
   sqlite3 memory_bot.db ".indexes memories"
   ```

2. **Consolidation Timing**: Schedule during low-usage hours (2-4 AM).

3. **Batch Size**: If you have many memories, consider batching consolidation:
   ```go
   // Process in batches of 100
   const batchSize = 100
   ```

4. **Monitor DB Size**: Check database growth:
   ```bash
   du -h memory_bot.db
   ```

## Next Steps

1. **Monitor** the system for 7 days to see consolidation patterns
2. **Tune** emotional weights based on user feedback
3. **Analyze** which contextual searches are most common
4. **Optimize** review intervals based on retention data

## Need Help?

Check the full documentation: `docs/BIOLOGICAL_MEMORY_SYSTEM.md`

Report issues: https://github.com/Milanz247/Personal-memory-reminder-bot/issues
