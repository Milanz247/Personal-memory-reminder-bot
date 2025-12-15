package service

import (
	"strings"
	"time"
)

// ContextualData represents the context in which a memory was created
// Mimics the Hippocampus's role in encoding contextual information
type ContextualData struct {
	TimeOfDay  string // "Morning", "Afternoon", "Evening", "Night"
	DayOfWeek  string // "Monday", "Tuesday", etc.
	ChatSource string // "Telegram", "WhatsApp", etc.
}

// ContextualMetadataService captures contextual information for memories
// Biological principle: The brain encodes when and where memories are formed
type ContextualMetadataService struct{}

// NewContextualMetadataService creates a new contextual metadata service
func NewContextualMetadataService() *ContextualMetadataService {
	return &ContextualMetadataService{}
}

// GetCurrentContext captures the current contextual information
func (s *ContextualMetadataService) GetCurrentContext(timestamp time.Time, source string) ContextualData {
	return ContextualData{
		TimeOfDay:  s.getTimeOfDay(timestamp),
		DayOfWeek:  timestamp.Weekday().String(),
		ChatSource: source,
	}
}

// getTimeOfDay categorizes the hour into time periods
func (s *ContextualMetadataService) getTimeOfDay(t time.Time) string {
	hour := t.Hour()

	switch {
	case hour >= 5 && hour < 12:
		return "Morning"
	case hour >= 12 && hour < 17:
		return "Afternoon"
	case hour >= 17 && hour < 21:
		return "Evening"
	default:
		return "Night"
	}
}

// ExtractContextCue parses user search queries for contextual hints
// Examples: "last night", "yesterday morning", "Tuesday's meeting"
func (s *ContextualMetadataService) ExtractContextCue(query string) (ContextualData, bool) {
	lowerQuery := strings.ToLower(query)
	context := ContextualData{}
	hasContext := false

	// Time of day patterns
	if strings.Contains(lowerQuery, "morning") {
		context.TimeOfDay = "Morning"
		hasContext = true
	} else if strings.Contains(lowerQuery, "afternoon") {
		context.TimeOfDay = "Afternoon"
		hasContext = true
	} else if strings.Contains(lowerQuery, "evening") {
		context.TimeOfDay = "Evening"
		hasContext = true
	} else if strings.Contains(lowerQuery, "night") {
		context.TimeOfDay = "Night"
		hasContext = true
	}

	// Day of week patterns
	daysOfWeek := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	for _, day := range daysOfWeek {
		if strings.Contains(lowerQuery, day) {
			context.DayOfWeek = strings.Title(day)
			hasContext = true
			break
		}
	}

	// Relative time patterns
	if strings.Contains(lowerQuery, "yesterday") {
		yesterday := time.Now().AddDate(0, 0, -1)
		context.DayOfWeek = yesterday.Weekday().String()
		hasContext = true
	} else if strings.Contains(lowerQuery, "last week") {
		lastWeek := time.Now().AddDate(0, 0, -7)
		context.DayOfWeek = lastWeek.Weekday().String()
		hasContext = true
	}

	return context, hasContext
}

// MatchesContext checks if a memory matches the given contextual criteria
func (s *ContextualMetadataService) MatchesContext(memoryTimeOfDay, memoryDayOfWeek string, searchContext ContextualData) bool {
	timeMatches := searchContext.TimeOfDay == "" || memoryTimeOfDay == searchContext.TimeOfDay
	dayMatches := searchContext.DayOfWeek == "" || memoryDayOfWeek == searchContext.DayOfWeek

	return timeMatches && dayMatches
}

// GetContextDescription returns a human-readable description of the context
func (s *ContextualMetadataService) GetContextDescription(context ContextualData) string {
	parts := []string{}

	if context.DayOfWeek != "" {
		parts = append(parts, context.DayOfWeek)
	}

	if context.TimeOfDay != "" {
		parts = append(parts, context.TimeOfDay)
	}

	if len(parts) == 0 {
		return "Unknown context"
	}

	return strings.Join(parts, " ")
}
