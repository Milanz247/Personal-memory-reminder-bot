package job

import (
	"context"
	"log"
	"time"

	"memory-bot/internal/domain/entity"
	"memory-bot/internal/domain/repository"
)

// DailyConsolidationJob simulates sleep-based memory consolidation
// Biological principle: During sleep, the brain strengthens new memories
// and transfers them from short-term (Hippocampus) to long-term storage (Cortex)
type DailyConsolidationJob struct {
	repo repository.MemoryRepository
}

// NewDailyConsolidationJob creates a new consolidation job
func NewDailyConsolidationJob(repo repository.MemoryRepository) *DailyConsolidationJob {
	return &DailyConsolidationJob{
		repo: repo,
	}
}

// Execute runs the consolidation process
// This should be scheduled to run daily, preferably at night (e.g., 2:00 AM)
func (j *DailyConsolidationJob) Execute() error {
	ctx := context.Background()
	log.Println("🌙 Starting daily memory consolidation (simulating sleep)...")

	// 1. Identify fragile memories (recently created, not yet reviewed)
	fragileMemories, err := j.getFragileMemories(ctx)
	if err != nil {
		log.Printf("Error getting fragile memories: %v", err)
		return err
	}

	if len(fragileMemories) == 0 {
		log.Println("No fragile memories to consolidate")
		return nil
	}

	log.Printf("Found %d fragile memories for consolidation", len(fragileMemories))

	// 2. Apply consolidation to each memory
	consolidated := 0
	for _, memory := range fragileMemories {
		if err := j.consolidateMemory(ctx, memory); err != nil {
			log.Printf("Error consolidating memory %d: %v", memory.ID, err)
			continue
		}
		consolidated++
	}

	log.Printf("✅ Consolidation complete: %d/%d memories strengthened", consolidated, len(fragileMemories))
	return nil
}

// getFragileMemories identifies memories that need consolidation
// Fragile = Created within last 7 days AND review count < 2
func (j *DailyConsolidationJob) getFragileMemories(ctx context.Context) ([]*entity.Memory, error) {
	// This would be a new repository method
	// For now, we'll get recent memories and filter
	return j.repo.GetFragileMemories(ctx)
}

// consolidateMemory applies consolidation to a single memory
// Biological principle: New memories receive temporary strengthening
func (j *DailyConsolidationJob) consolidateMemory(ctx context.Context, memory *entity.Memory) error {
	// Calculate days since creation
	daysSinceCreation := int(time.Since(memory.CreatedAt).Hours() / 24)

	// Apply priority boost based on age and emotional weight
	// Newer memories get higher boost
	if daysSinceCreation <= 1 {
		// First day: strong boost
		memory.PriorityScore = 0.5 * (1.0 + memory.EmotionalWeight)
	} else if daysSinceCreation <= 3 {
		// Days 2-3: moderate boost
		memory.PriorityScore = 0.3 * (1.0 + memory.EmotionalWeight)
	} else if daysSinceCreation <= 7 {
		// Days 4-7: small boost
		memory.PriorityScore = 0.15 * (1.0 + memory.EmotionalWeight)
	} else {
		// After 7 days, memory is considered consolidated
		memory.PriorityScore = 0.0
	}

	// Mark as consolidated
	memory.LastConsolidated = time.Now()

	// Update in repository
	return j.repo.UpdateConsolidation(ctx, memory)
}

// ScheduleDaily schedules the job to run at a specific time each day
func (j *DailyConsolidationJob) ScheduleDaily(hour, minute int) {
	log.Printf("Scheduling daily consolidation at %02d:%02d", hour, minute)

	go func() {
		for {
			now := time.Now()

			// Calculate next run time
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
			if next.Before(now) {
				// If the time has passed today, schedule for tomorrow
				next = next.Add(24 * time.Hour)
			}

			// Wait until next run time
			duration := next.Sub(now)
			log.Printf("Next consolidation job in %v", duration)

			time.Sleep(duration)

			// Execute the job
			if err := j.Execute(); err != nil {
				log.Printf("Consolidation job failed: %v", err)
			}
		}
	}()
}

// RunNow executes the job immediately (useful for testing)
func (j *DailyConsolidationJob) RunNow() error {
	return j.Execute()
}
