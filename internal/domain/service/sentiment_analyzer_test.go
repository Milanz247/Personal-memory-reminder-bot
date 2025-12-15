package service_test

import (
	"memory-bot/internal/domain/service"
	"testing"
)

func TestSentimentAnalyzer_Analyze(t *testing.T) {
	analyzer := service.NewSentimentAnalyzer()

	tests := []struct {
		name           string
		content        string
		minWeight      float64
		maxWeight      float64
		expectedLevel  string
	}{
		{
			name:          "Highly positive emotional content",
			content:       "This is an amazing and wonderful breakthrough! I'm so excited and proud!",
			minWeight:     0.7,
			maxWeight:     1.0,
			expectedLevel: "Strong",
		},
		{
			name:          "Multiple positive words",
			content:       "Excellent work today! Fantastic results, absolutely brilliant and perfect!",
			minWeight:     0.8,
			maxWeight:     1.0,
			expectedLevel: "Intense",
		},
		{
			name:          "Highly negative emotional content",
			content:       "Terrible day. Feeling anxious, frustrated and stressed about the crisis.",
			minWeight:     0.7,
			maxWeight:     1.0,
			expectedLevel: "Strong",
		},
		{
			name:          "Neutral content - short",
			content:       "Meeting at 3pm tomorrow.",
			minWeight:     0.0,
			maxWeight:     0.3,
			expectedLevel: "Neutral",
		},
		{
			name:          "Neutral content - medium",
			content:       "Discussed project timeline and deliverables in the conference room.",
			minWeight:     0.0,
			maxWeight:     0.3,
			expectedLevel: "Neutral",
		},
		{
			name:          "Long detailed content",
			content:       "Today we had a comprehensive discussion about the project architecture, database design, API endpoints, authentication mechanisms, and deployment strategies. We covered all aspects thoroughly.",
			minWeight:     0.2,
			maxWeight:     0.5,
			expectedLevel: "Neutral",
		},
		{
			name:          "Content with exclamation marks",
			content:       "Important update!! Need to review this urgently!!",
			minWeight:     0.4,
			maxWeight:     0.7,
			expectedLevel: "Moderate",
		},
		{
			name:          "Mixed emotions",
			content:       "Bad news about the project delay but great progress on the new feature!",
			minWeight:     0.5,
			maxWeight:     0.9,
			expectedLevel: "Moderate",
		},
		{
			name:          "Empty content",
			content:       "",
			minWeight:     0.0,
			maxWeight:     0.0,
			expectedLevel: "Neutral",
		},
		{
			name:          "Single positive word",
			content:       "Amazing!",
			minWeight:     0.7,
			maxWeight:     1.0,
			expectedLevel: "Strong",
		},
		{
			name:          "Single negative word",
			content:       "Disaster.",
			minWeight:     0.8,
			maxWeight:     1.0,
			expectedLevel: "Intense",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weight := analyzer.Analyze(tt.content)

			if weight < tt.minWeight || weight > tt.maxWeight {
				t.Errorf("Analyze() weight = %v, want between %v and %v", weight, tt.minWeight, tt.maxWeight)
			}

			category := analyzer.GetEmotionalCategory(weight)
			if category != tt.expectedLevel {
				t.Errorf("GetEmotionalCategory() = %v, want %v (weight: %v)", category, tt.expectedLevel, weight)
			}

			t.Logf("✅ Content: %q → Weight: %.2f, Category: %s", tt.content[:min(len(tt.content), 50)], weight, category)
		})
	}
}

func TestSentimentAnalyzer_GetEmotionalCategory(t *testing.T) {
	analyzer := service.NewSentimentAnalyzer()

	tests := []struct {
		weight   float64
		expected string
	}{
		{0.0, "Neutral"},
		{0.1, "Neutral"},
		{0.29, "Neutral"},
		{0.3, "Moderate"},
		{0.5, "Moderate"},
		{0.59, "Moderate"},
		{0.6, "Strong"},
		{0.7, "Strong"},
		{0.79, "Strong"},
		{0.8, "Intense"},
		{0.9, "Intense"},
		{1.0, "Intense"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := analyzer.GetEmotionalCategory(tt.weight)
			if result != tt.expected {
				t.Errorf("GetEmotionalCategory(%v) = %v, want %v", tt.weight, result, tt.expected)
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
