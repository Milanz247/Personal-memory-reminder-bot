package scheduler_test

import (
	"math"
	"memory-bot/internal/domain/entity"
	"memory-bot/internal/infrastructure/scheduler"
	"testing"
	"time"
)

func TestBiologicalSpacedRepetition_CalculateNextReviewInterval(t *testing.T) {
	baseIntervals := []int{1, 3, 7, 14, 30}
	bio := scheduler.NewBiologicalSpacedRepetition(baseIntervals)

	tests := []struct {
		name            string
		memory          *entity.Memory
		expectedMinDays float64
		expectedMaxDays float64
	}{
		{
			name: "First review - neutral memory",
			memory: &entity.Memory{
				ReviewCount:     0,
				EmotionalWeight: 0.0,
				PriorityScore:   0.0,
			},
			expectedMinDays: 0.9,
			expectedMaxDays: 1.1,
		},
		{
			name: "First review - highly emotional",
			memory: &entity.Memory{
				ReviewCount:     0,
				EmotionalWeight: 0.8,
				PriorityScore:   0.0,
			},
			expectedMinDays: 1.3, // 1 * 1.4 (emotional boost)
			expectedMaxDays: 1.5,
		},
		{
			name: "First review - with priority boost",
			memory: &entity.Memory{
				ReviewCount:     0,
				EmotionalWeight: 0.5,
				PriorityScore:   0.5,
			},
			expectedMinDays: 1.8, // 1 * 1.25 * 1.5
			expectedMaxDays: 2.0,
		},
		{
			name: "First review - max emotional and priority",
			memory: &entity.Memory{
				ReviewCount:     0,
				EmotionalWeight: 1.0,
				PriorityScore:   0.9,
			},
			expectedMinDays: 2.7, // 1 * 1.5 * 1.9
			expectedMaxDays: 3.0,
		},
		{
			name: "Second review - neutral",
			memory: &entity.Memory{
				ReviewCount:     1,
				EmotionalWeight: 0.0,
				PriorityScore:   0.0,
			},
			expectedMinDays: 2.9,
			expectedMaxDays: 3.1,
		},
		{
			name: "Third review - emotional",
			memory: &entity.Memory{
				ReviewCount:     2,
				EmotionalWeight: 0.7,
				PriorityScore:   0.0,
			},
			expectedMinDays: 9.0, // 7 * 1.35
			expectedMaxDays: 10.0,
		},
		{
			name: "Fifth review - max emotional",
			memory: &entity.Memory{
				ReviewCount:     4,
				EmotionalWeight: 1.0,
				PriorityScore:   0.0,
			},
			expectedMinDays: 44.0, // 30 * 1.5
			expectedMaxDays: 46.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interval := bio.CalculateNextReviewInterval(tt.memory)
			days := interval.Hours() / 24

			if days < tt.expectedMinDays || days > tt.expectedMaxDays {
				t.Errorf("CalculateNextReviewInterval() = %.2f days, want between %.2f and %.2f",
					days, tt.expectedMinDays, tt.expectedMaxDays)
			}

			t.Logf("✅ ReviewCount=%d, Emotional=%.1f, Priority=%.1f → %.2f days",
				tt.memory.ReviewCount, tt.memory.EmotionalWeight, tt.memory.PriorityScore, days)
		})
	}
}

func TestBiologicalSpacedRepetition_GetNextReviewTime(t *testing.T) {
	baseIntervals := []int{1, 3, 7, 14, 30}
	bio := scheduler.NewBiologicalSpacedRepetition(baseIntervals)

	now := time.Now()
	lastWeek := now.AddDate(0, 0, -7)

	tests := []struct {
		name          string
		memory        *entity.Memory
		expectedAfter time.Time
	}{
		{
			name: "New memory never reviewed",
			memory: &entity.Memory{
				CreatedAt:       now,
				LastReviewed:    nil,
				ReviewCount:     0,
				EmotionalWeight: 0.0,
				PriorityScore:   0.0,
			},
			expectedAfter: now,
		},
		{
			name: "Memory reviewed last week",
			memory: &entity.Memory{
				CreatedAt:       lastWeek,
				LastReviewed:    &lastWeek,
				ReviewCount:     1,
				EmotionalWeight: 0.0,
				PriorityScore:   0.0,
			},
			expectedAfter: lastWeek,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextReview := bio.GetNextReviewTime(tt.memory)

			if nextReview.Before(tt.expectedAfter) {
				t.Errorf("GetNextReviewTime() = %v, should be after %v", nextReview, tt.expectedAfter)
			}

			t.Logf("✅ Next review at: %v", nextReview.Format("2006-01-02 15:04"))
		})
	}
}

func TestBiologicalSpacedRepetition_ShouldReviewNow(t *testing.T) {
	baseIntervals := []int{1, 3, 7, 14, 30}
	bio := scheduler.NewBiologicalSpacedRepetition(baseIntervals)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastWeek := now.AddDate(0, 0, -7)

	tests := []struct {
		name         string
		memory       *entity.Memory
		shouldReview bool
	}{
		{
			name: "Memory due yesterday",
			memory: &entity.Memory{
				CreatedAt:       lastWeek,
				LastReviewed:    &yesterday,
				ReviewCount:     1,
				EmotionalWeight: 0.0,
				PriorityScore:   0.0,
			},
			shouldReview: true,
		},
		{
			name: "Memory due tomorrow",
			memory: &entity.Memory{
				CreatedAt:       now,
				LastReviewed:    &now,
				ReviewCount:     0,
				EmotionalWeight: 0.0,
				PriorityScore:   0.0,
			},
			shouldReview: false,
		},
		{
			name: "New memory created just now",
			memory: &entity.Memory{
				CreatedAt:       now,
				LastReviewed:    nil,
				ReviewCount:     0,
				EmotionalWeight: 0.0,
				PriorityScore:   0.0,
			},
			shouldReview: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bio.ShouldReviewNow(tt.memory)

			if result != tt.shouldReview {
				t.Errorf("ShouldReviewNow() = %v, want %v", result, tt.shouldReview)
			}

			t.Logf("✅ Should review: %v", result)
		})
	}
}

func TestBiologicalSpacedRepetition_CalculateForgettingCurve(t *testing.T) {
	baseIntervals := []int{1, 3, 7, 14, 30}
	bio := scheduler.NewBiologicalSpacedRepetition(baseIntervals)

	tests := []struct {
		name            string
		memory          *entity.Memory
		daysSinceReview int
		minRetention    float64
		maxRetention    float64
	}{
		{
			name: "New memory - 1 day later",
			memory: &entity.Memory{
				ReviewCount:     0,
				EmotionalWeight: 0.0,
			},
			daysSinceReview: 1,
			minRetention:    0.3,
			maxRetention:    0.4,
		},
		{
			name: "Strong memory - 1 day later",
			memory: &entity.Memory{
				ReviewCount:     5,
				EmotionalWeight: 0.8,
			},
			daysSinceReview: 1,
			minRetention:    0.8,
			maxRetention:    1.0,
		},
		{
			name: "New emotional memory - 3 days later",
			memory: &entity.Memory{
				ReviewCount:     0,
				EmotionalWeight: 0.9,
			},
			daysSinceReview: 3,
			minRetention:    0.3,
			maxRetention:    0.6,
		},
		{
			name: "Weak memory - 7 days later",
			memory: &entity.Memory{
				ReviewCount:     0,
				EmotionalWeight: 0.0,
			},
			daysSinceReview: 7,
			minRetention:    0.0,
			maxRetention:    0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retention := bio.CalculateForgettingCurve(tt.memory, tt.daysSinceReview)

			if retention < tt.minRetention || retention > tt.maxRetention {
				t.Errorf("CalculateForgettingCurve() = %.3f, want between %.3f and %.3f",
					retention, tt.minRetention, tt.maxRetention)
			}

			t.Logf("✅ ReviewCount=%d, Emotional=%.1f, Days=%d → Retention=%.2f%%",
				tt.memory.ReviewCount, tt.memory.EmotionalWeight, tt.daysSinceReview, retention*100)
		})
	}
}

func TestBiologicalSpacedRepetition_NeedsUrgentReview(t *testing.T) {
	baseIntervals := []int{1, 3, 7, 14, 30}
	bio := scheduler.NewBiologicalSpacedRepetition(baseIntervals)

	now := time.Now()
	monthAgo := now.AddDate(0, 0, -30)

	tests := []struct {
		name           string
		memory         *entity.Memory
		expectedUrgent bool
	}{
		{
			name: "New weak memory - long overdue",
			memory: &entity.Memory{
				CreatedAt:       monthAgo,
				LastReviewed:    &monthAgo,
				ReviewCount:     0,
				EmotionalWeight: 0.0,
			},
			expectedUrgent: true,
		},
		{
			name: "Strong emotional memory - same duration",
			memory: &entity.Memory{
				CreatedAt:       monthAgo,
				LastReviewed:    &monthAgo,
				ReviewCount:     5,
				EmotionalWeight: 0.9,
			},
			expectedUrgent: false,
		},
		{
			name: "Recently reviewed",
			memory: &entity.Memory{
				CreatedAt:       now,
				LastReviewed:    &now,
				ReviewCount:     0,
				EmotionalWeight: 0.0,
			},
			expectedUrgent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bio.NeedsUrgentReview(tt.memory)

			if result != tt.expectedUrgent {
				t.Errorf("NeedsUrgentReview() = %v, want %v", result, tt.expectedUrgent)
			}

			t.Logf("✅ Urgent review needed: %v", result)
		})
	}
}

func TestBiologicalSpacedRepetition_ExponentialGrowth(t *testing.T) {
	baseIntervals := []int{1, 3, 7, 14, 30}
	bio := scheduler.NewBiologicalSpacedRepetition(baseIntervals)

	// Test that intervals grow exponentially after base intervals
	memory := &entity.Memory{
		ReviewCount:     len(baseIntervals) + 3, // Beyond base intervals
		EmotionalWeight: 0.0,
		PriorityScore:   0.0,
	}

	interval := bio.CalculateNextReviewInterval(memory)
	days := interval.Hours() / 24

	// After 5 base intervals, should be: 30 * 2^3 = 240 days
	expectedDays := 240.0

	if math.Abs(days-expectedDays) > 10 {
		t.Errorf("Exponential growth: got %.0f days, want approximately %.0f days", days, expectedDays)
	}

	t.Logf("✅ Review %d → %.0f days (exponential growth working)", memory.ReviewCount, days)
}
