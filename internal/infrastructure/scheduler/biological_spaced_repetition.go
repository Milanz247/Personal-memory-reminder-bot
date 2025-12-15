package scheduler

import (
	"fmt"
	"math"
	"memory-bot/internal/domain/entity"
	"time"
)

// BiologicalSpacedRepetition implements spaced repetition with biological principles
// Simulates LTP (Long-Term Potentiation) and the Forgetting Curve
type BiologicalSpacedRepetition struct {
	baseIntervals []int // Base intervals in days: [1, 3, 7, 14, 30]
}

// NewBiologicalSpacedRepetition creates a new biological spaced repetition calculator
func NewBiologicalSpacedRepetition(baseIntervals []int) *BiologicalSpacedRepetition {
	if len(baseIntervals) == 0 {
		baseIntervals = []int{1, 3, 7, 14, 30} // Default intervals
	}
	return &BiologicalSpacedRepetition{
		baseIntervals: baseIntervals,
	}
}

// CalculateNextReviewInterval calculates when a memory should be reviewed next
// Biological Principles Applied:
// 1. LTP (Long-Term Potentiation): Each successful review strengthens the synapse
// 2. Emotional Modulation: High emotional weight = slower decay (Amygdala effect)
// 3. Priority Boost: New fragile memories get temporary protection
func (b *BiologicalSpacedRepetition) CalculateNextReviewInterval(memory *entity.Memory) time.Duration {
	// 1. Get base factor from review count (LTP simulation)
	baseFactor := b.getBaseFactor(memory.ReviewCount)

	// 2. Emotional Modulation (Amygdala's influence)
	// High emotional weight means stronger encoding and slower forgetting
	// EmotionalWeight ranges from 0.0 to 1.0
	// We apply a boost of up to 50% for highly emotional memories
	emotionalBoost := 1.0 + (memory.EmotionalWeight * 0.5)

	// 3. Priority Score (from consolidation job)
	// New memories get temporary boost to prevent early decay
	priorityBoost := 1.0
	if memory.PriorityScore > 0 {
		priorityBoost = 1.0 + memory.PriorityScore
	}

	// 4. Apply all modulations
	finalFactor := baseFactor * emotionalBoost * priorityBoost

	// 5. Calculate final interval in hours
	hoursUntilReview := finalFactor * 24.0

	return time.Duration(hoursUntilReview) * time.Hour
}

// getBaseFactor returns the multiplication factor based on review count
// Simulates synaptic strengthening (LTP)
func (b *BiologicalSpacedRepetition) getBaseFactor(reviewCount int) float64 {
	// If we have predefined intervals, use them
	if reviewCount < len(b.baseIntervals) {
		return float64(b.baseIntervals[reviewCount])
	}

	// After all intervals are used, apply exponential growth
	// This simulates very strong long-term memories
	lastInterval := float64(b.baseIntervals[len(b.baseIntervals)-1])
	additionalReviews := reviewCount - len(b.baseIntervals)

	// Each additional review doubles the interval
	return lastInterval * math.Pow(2, float64(additionalReviews))
}

// GetNextReviewTime calculates the exact time for next review
func (b *BiologicalSpacedRepetition) GetNextReviewTime(memory *entity.Memory) time.Time {
	interval := b.CalculateNextReviewInterval(memory)

	// Base time is either last review or creation time
	var baseTime time.Time
	if memory.LastReviewed != nil {
		baseTime = *memory.LastReviewed
	} else {
		baseTime = memory.CreatedAt
	}

	return baseTime.Add(interval)
}

// ShouldReviewNow checks if a memory needs review right now
func (b *BiologicalSpacedRepetition) ShouldReviewNow(memory *entity.Memory) bool {
	nextReview := b.GetNextReviewTime(memory)
	return time.Now().After(nextReview)
}

// GetReviewDaysDescription returns a human-readable description
func (b *BiologicalSpacedRepetition) GetReviewDaysDescription(memory *entity.Memory) string {
	interval := b.CalculateNextReviewInterval(memory)
	days := int(interval.Hours() / 24)

	if days == 1 {
		return "tomorrow"
	} else if days < 7 {
		return fmt.Sprintf("in %d days", days)
	} else if days < 30 {
		weeks := days / 7
		return fmt.Sprintf("in %d week(s)", weeks)
	} else {
		months := days / 30
		return fmt.Sprintf("in %d month(s)", months)
	}
}

// CalculateForgettingCurve estimates memory retention probability
// Based on Ebbinghaus forgetting curve: R(t) = e^(-t/S)
// Where R is retention, t is time, and S is memory strength
func (b *BiologicalSpacedRepetition) CalculateForgettingCurve(memory *entity.Memory, daysSinceReview int) float64 {
	// Memory strength increases with review count and emotional weight
	memoryStrength := float64(memory.ReviewCount+1) * (1.0 + memory.EmotionalWeight)

	// Apply forgetting curve formula
	retention := math.Exp(-float64(daysSinceReview) / memoryStrength)

	return retention
}

// NeedsUrgentReview checks if retention has dropped below threshold
func (b *BiologicalSpacedRepetition) NeedsUrgentReview(memory *entity.Memory) bool {
	daysSince := memory.DaysSinceLastReview()
	retention := b.CalculateForgettingCurve(memory, daysSince)

	// If retention drops below 30%, it's urgent
	return retention < 0.3
}
