package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"memory-bot/internal/application/usecase"
	"memory-bot/internal/domain/entity"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SpacedRepetitionScheduler handles automatic memory review reminders
type SpacedRepetitionScheduler struct {
	api       *tgbotapi.BotAPI
	useCase   *usecase.ReviewMemoryUseCase
	intervals []int
	ticker    *time.Ticker
	stopChan  chan bool
}

// NewSpacedRepetitionScheduler creates a new scheduler
func NewSpacedRepetitionScheduler(
	api *tgbotapi.BotAPI,
	useCase *usecase.ReviewMemoryUseCase,
	intervals []int,
) *SpacedRepetitionScheduler {
	return &SpacedRepetitionScheduler{
		api:       api,
		useCase:   useCase,
		intervals: intervals,
		stopChan:  make(chan bool),
	}
}

// Start starts the spaced repetition scheduler
func (s *SpacedRepetitionScheduler) Start() {
	log.Println("Spaced repetition scheduler started")

	// Run immediately on start
	s.checkAndSendReviews()

	// Then run every hour
	s.ticker = time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.checkAndSendReviews()
			case <-s.stopChan:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the scheduler
func (s *SpacedRepetitionScheduler) Stop() {
	log.Println("Stopping spaced repetition scheduler")
	s.stopChan <- true
}

// checkAndSendReviews checks for memories that need review and sends them
func (s *SpacedRepetitionScheduler) checkAndSendReviews() {
	ctx := context.Background()
	log.Println("Checking for memories to review...")

	input := usecase.ReviewMemoryInput{
		Intervals: s.intervals,
	}

	output, err := s.useCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Error getting memories for review: %v", err)
		return
	}

	if len(output.Memories) == 0 {
		log.Println("No memories need review at this time")
		return
	}

	log.Printf("Found %d memories for review", len(output.Memories))

	// Group memories by user
	userMemories := make(map[int64][]*entity.Memory)
	for _, mem := range output.Memories {
		userMemories[mem.UserID] = append(userMemories[mem.UserID], mem)
	}

	// Send review reminders to each user
	for userID, mems := range userMemories {
		s.sendReviewToUser(ctx, userID, mems)
	}
}

// sendReviewToUser sends memory review reminders to a specific user
func (s *SpacedRepetitionScheduler) sendReviewToUser(ctx context.Context, userID int64, memories []*entity.Memory) {
	if len(memories) == 0 {
		return
	}

	chatID := memories[0].ChatID

	// Send header message
	headerText := fmt.Sprintf("🔔 *Memory Review Time!*\n\nYou have %d memories to review:\n", len(memories))
	msg := tgbotapi.NewMessage(chatID, headerText)
	msg.ParseMode = "Markdown"
	s.api.Send(msg)

	// Send each memory with review buttons
	for i, mem := range memories {
		if i >= 5 { // Limit to 5 reviews per session
			remainingText := fmt.Sprintf("\n📚 +%d more memories waiting for review. They'll appear in the next session.", len(memories)-5)
			s.api.Send(tgbotapi.NewMessage(chatID, remainingText))
			break
		}

		s.sendMemoryForReview(chatID, mem)

		// Mark as reviewed
		if err := s.useCase.MarkAsReviewed(ctx, mem.ID); err != nil {
			log.Printf("Error updating last_reviewed for memory %d: %v", mem.ID, err)
		}

		// Small delay between messages
		time.Sleep(500 * time.Millisecond)
	}

	// Send completion message
	completionText := "✅ Review session complete! Great job maintaining your memories. 🧠"
	s.api.Send(tgbotapi.NewMessage(chatID, completionText))

	log.Printf("Sent %d review reminders to user %d", min(len(memories), 5), userID)
}

// sendMemoryForReview sends a single memory for review
func (s *SpacedRepetitionScheduler) sendMemoryForReview(chatID int64, mem *entity.Memory) {
	daysSince := mem.DaysSinceLastReview()

	reviewText := fmt.Sprintf(
		"💭 *Memory #%d*\n\n%s\n\n"+
			"📅 Created: %s\n"+
			"🔄 Last reviewed: %d days ago\n"+
			"📊 Review count: %d\n\n"+
			"_Take a moment to recall this memory..._",
		mem.ID,
		mem.Content,
		mem.CreatedAt.Format("2006-01-02"),
		daysSince,
		mem.ReviewCount,
	)

	msg := tgbotapi.NewMessage(chatID, reviewText)
	msg.ParseMode = "Markdown"

	// Add inline keyboard with feedback options
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Remembered", fmt.Sprintf("review:remember:%d", mem.ID)),
			tgbotapi.NewInlineKeyboardButtonData("❓ Forgot", fmt.Sprintf("review:forgot:%d", mem.ID)),
		),
	)
	msg.ReplyMarkup = keyboard

	s.api.Send(msg)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
