package entity

import (
	"strings"
	"time"
)

// Memory represents a stored memory with context
// This is the core domain entity containing business logic
type Memory struct {
	ID           int
	UserID       int64
	ChatID       int64
	Content      string
	Tags         []string
	CreatedAt    time.Time
	LastReviewed *time.Time
	ReviewCount  int
	Rank         float64 // For search result ranking
}

// NewMemory creates a new Memory entity with validation
func NewMemory(userID, chatID int64, content string) *Memory {
	memory := &Memory{
		UserID:      userID,
		ChatID:      chatID,
		Content:     strings.TrimSpace(content),
		CreatedAt:   time.Now(),
		ReviewCount: 0,
	}
	
	// Automatically extract tags
	memory.Tags = memory.extractTags()
	
	return memory
}

// NeedsReview determines if a memory needs review based on spaced repetition intervals
func (m *Memory) NeedsReview(intervals []int) bool {
	if m.ReviewCount >= len(intervals) {
		return false
	}
	
	var lastTime time.Time
	if m.LastReviewed != nil {
		lastTime = *m.LastReviewed
	} else {
		lastTime = m.CreatedAt
	}
	
	daysSince := int(time.Since(lastTime).Hours() / 24)
	
	// Get the interval for current review count
	interval := intervals[m.ReviewCount]
	
	return daysSince >= interval
}

// MarkAsReviewed updates the review metadata
func (m *Memory) MarkAsReviewed() {
	now := time.Now()
	m.LastReviewed = &now
	m.ReviewCount++
}

// extractTags extracts hashtags from the memory content
func (m *Memory) extractTags() []string {
	var tags []string
	words := strings.Fields(m.Content)
	
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}
	
	return tags
}

// GetTagsString returns tags as a space-separated string for storage
func (m *Memory) GetTagsString() string {
	return strings.Join(m.Tags, " ")
}

// DaysSinceLastReview returns the number of days since last review or creation
func (m *Memory) DaysSinceLastReview() int {
	var lastTime time.Time
	if m.LastReviewed != nil {
		lastTime = *m.LastReviewed
	} else {
		lastTime = m.CreatedAt
	}
	
	return int(time.Since(lastTime).Hours() / 24)
}

// Validate checks if the memory entity is valid
func (m *Memory) Validate() error {
	if m.UserID == 0 {
		return ErrInvalidUserID
	}
	if m.ChatID == 0 {
		return ErrInvalidChatID
	}
	if strings.TrimSpace(m.Content) == "" {
		return ErrEmptyContent
	}
	return nil
}
