package service

import (
	"strings"
)

// SentimentAnalyzer mimics the Amygdala's role in emotional tagging
// It assigns emotional weight to memories based on their content
type SentimentAnalyzer struct {
	positiveWords map[string]float64
	negativeWords map[string]float64
}

// NewSentimentAnalyzer creates a new sentiment analyzer
func NewSentimentAnalyzer() *SentimentAnalyzer {
	return &SentimentAnalyzer{
		positiveWords: map[string]float64{
			"amazing":     0.9,
			"excellent":   0.8,
			"wonderful":   0.8,
			"great":       0.7,
			"love":        0.8,
			"happy":       0.7,
			"excited":     0.8,
			"success":     0.7,
			"achievement": 0.8,
			"proud":       0.7,
			"fantastic":   0.9,
			"brilliant":   0.8,
			"perfect":     0.8,
			"awesome":     0.8,
			"incredible":  0.9,
			"beautiful":   0.7,
			"joy":         0.8,
			"celebrate":   0.7,
			"win":         0.7,
			"victory":     0.8,
		},
		negativeWords: map[string]float64{
			"terrible":     0.9,
			"horrible":     0.9,
			"awful":        0.8,
			"bad":          0.6,
			"hate":         0.8,
			"angry":        0.7,
			"sad":          0.7,
			"failure":      0.8,
			"disappointed": 0.7,
			"frustrated":   0.7,
			"crisis":       0.9,
			"disaster":     0.9,
			"worried":      0.7,
			"anxious":      0.7,
			"stress":       0.7,
			"fear":         0.8,
			"panic":        0.8,
			"upset":        0.7,
			"miserable":    0.8,
			"devastating":  0.9,
		},
	}
}

// Analyze determines the emotional weight of a memory
// Returns a value between 0.0 (neutral) and 1.0 (highly emotional)
// Biological principle: The Amygdala strengthens memories with strong emotional content
func (s *SentimentAnalyzer) Analyze(content string) float64 {
	if content == "" {
		return 0.0
	}

	content = strings.ToLower(content)
	words := strings.Fields(content)

	var totalEmotionalScore float64
	var emotionalWordCount int

	// Scan for emotional keywords
	for _, word := range words {
		// Clean punctuation
		word = strings.Trim(word, ".,!?;:")

		if weight, exists := s.positiveWords[word]; exists {
			totalEmotionalScore += weight
			emotionalWordCount++
		}

		if weight, exists := s.negativeWords[word]; exists {
			totalEmotionalScore += weight
			emotionalWordCount++
		}
	}

	// If no emotional words found, check content length and punctuation
	if emotionalWordCount == 0 {
		// Long, detailed memories might be important
		if len(words) > 50 {
			return 0.3
		}

		// Multiple exclamation marks indicate emotion
		exclamationCount := strings.Count(content, "!")
		if exclamationCount >= 2 {
			return 0.5
		}

		return 0.1 // Default neutral
	}

	// Calculate average emotional weight
	avgEmotionalWeight := totalEmotionalScore / float64(emotionalWordCount)

	// Boost if multiple emotional words (indicates strong emotion)
	if emotionalWordCount > 2 {
		avgEmotionalWeight = avgEmotionalWeight * 1.2
	}

	// Cap at 1.0
	if avgEmotionalWeight > 1.0 {
		avgEmotionalWeight = 1.0
	}

	return avgEmotionalWeight
}

// GetEmotionalCategory returns a human-readable category
func (s *SentimentAnalyzer) GetEmotionalCategory(weight float64) string {
	if weight < 0.3 {
		return "Neutral"
	} else if weight < 0.6 {
		return "Moderate"
	} else if weight < 0.8 {
		return "Strong"
	}
	return "Intense"
}
