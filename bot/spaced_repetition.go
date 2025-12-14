package bot

import (
	"fmt"
	"log"
	"memory-bot/database"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SpacedRepetition handles automatic memory review reminders
type SpacedRepetition struct {
	bot       *Bot
	db        *database.Database
	intervals []int // Review intervals in days
	ticker    *time.Ticker
	stopChan  chan bool
}

// NewSpacedRepetition creates a new spaced repetition handler
func NewSpacedRepetition(bot *Bot, db *database.Database, intervals []int) *SpacedRepetition {
	return &SpacedRepetition{
		bot:       bot,
		db:        db,
		intervals: intervals,
		stopChan:  make(chan bool),
	}
}

// Start starts the spaced repetition scheduler
// It checks for memories to review every hour
func (sr *SpacedRepetition) Start() {
	log.Println("Spaced repetition scheduler started")

	// Run immediately on start
	sr.checkAndSendReviews()

	// Then run every hour
	sr.ticker = time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-sr.ticker.C:
				sr.checkAndSendReviews()
			case <-sr.stopChan:
				sr.ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the spaced repetition scheduler
func (sr *SpacedRepetition) Stop() {
	log.Println("Stopping spaced repetition scheduler")
	sr.stopChan <- true
}

// checkAndSendReviews checks for memories that need review and sends them to users
func (sr *SpacedRepetition) checkAndSendReviews() {
	log.Println("Checking for memories to review...")

	memories, err := sr.db.GetMemoriesForReview(sr.intervals)
	if err != nil {
		log.Printf("Error getting memories for review: %v", err)
		return
	}

	if len(memories) == 0 {
		log.Println("No memories need review at this time")
		return
	}

	log.Printf("Found %d memories for review", len(memories))

	// Group memories by user
	userMemories := make(map[int64][]database.Memory)
	for _, mem := range memories {
		userMemories[mem.UserID] = append(userMemories[mem.UserID], mem)
	}

	// Send review reminders to each user
	for userID, mems := range userMemories {
		sr.sendReviewToUser(userID, mems)
	}
}

// sendReviewToUser sends memory review reminders to a specific user
func (sr *SpacedRepetition) sendReviewToUser(userID int64, memories []database.Memory) {
	if len(memories) == 0 {
		return
	}

	// Get the first memory's chat ID (assuming user uses same chat)
	chatID := memories[0].ChatID

	// Send header message
	headerText := fmt.Sprintf("🔔 *Memory Review Time!*\n\nYou have %d memories to review:\n", len(memories))
	msg := tgbotapi.NewMessage(chatID, headerText)
	msg.ParseMode = "Markdown"
	sr.bot.api.Send(msg)

	// Send each memory with review buttons
	for i, mem := range memories {
		if i >= 5 { // Limit to 5 reviews per session
			remainingText := fmt.Sprintf("\n📚 +%d more memories waiting for review. They'll appear in the next session.", len(memories)-5)
			sr.bot.sendMessage(chatID, remainingText)
			break
		}

		sr.sendMemoryForReview(chatID, mem)

		// Mark as reviewed
		if err := sr.db.UpdateLastReviewed(mem.ID); err != nil {
			log.Printf("Error updating last_reviewed for memory %d: %v", mem.ID, err)
		}

		// Small delay between messages
		time.Sleep(500 * time.Millisecond)
	}

	// Send completion message
	completionText := "✅ Review session complete! Great job maintaining your memories. 🧠"
	sr.bot.sendMessage(chatID, completionText)

	log.Printf("Sent %d review reminders to user %d", min(len(memories), 5), userID)
}

// sendMemoryForReview sends a single memory for review
func (sr *SpacedRepetition) sendMemoryForReview(chatID int64, mem database.Memory) {
	// Calculate days since last review
	var daysSince int
	if mem.LastReviewed != nil {
		daysSince = int(time.Since(*mem.LastReviewed).Hours() / 24)
	} else {
		daysSince = int(time.Since(mem.CreatedAt).Hours() / 24)
	}

	reviewText := fmt.Sprintf(
		"💭 *Memory #%d*\n\n%s\n\n"+
			"📅 Created: %s\n"+
			"🔄 Last reviewed: %d days ago\n"+
			"📊 Review count: %d\n\n"+
			"_Take a moment to recall this memory..._",
		mem.ID,
		mem.TextContent,
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

	sr.bot.api.Send(msg)
}

// HandleReviewCallback handles callback queries for review feedback
func (sr *SpacedRepetition) HandleReviewCallback(query *tgbotapi.CallbackQuery, action string, memoryID int) {
	var response string

	switch action {
	case "remember":
		response = "Great! 🎉 Your memory is strong."
	case "forgot":
		response = "No worries! You'll see it again soon. 💪"
	default:
		response = "Thanks for the feedback!"
	}

	// Send callback response
	callback := tgbotapi.NewCallback(query.ID, response)
	sr.bot.api.Send(callback)

	// Update the message to remove buttons
	editText := query.Message.Text + "\n\n✓ _Reviewed_"
	edit := tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		editText,
	)
	edit.ParseMode = "Markdown"
	sr.bot.api.Send(edit)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
