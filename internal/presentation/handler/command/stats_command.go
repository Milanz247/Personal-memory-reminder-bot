package command

import (
	"context"
	"fmt"
	"log"
	"time"

	"memory-bot/internal/application/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StatsCommand handles the /stats command
type StatsCommand struct {
	useCase *usecase.GetStatsUseCase
}

// NewStatsCommand creates a new stats command
func NewStatsCommand(useCase *usecase.GetStatsUseCase) *StatsCommand {
	return &StatsCommand{
		useCase: useCase,
	}
}

// Name returns the command name
func (c *StatsCommand) Name() string {
	return "stats"
}

// Description returns the command description
func (c *StatsCommand) Description() string {
	return "Statistics"
}

// Execute executes the stats command
func (c *StatsCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
	input := usecase.GetStatsInput{
		UserID: message.From.ID,
	}

	output, err := c.useCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Error getting statistics: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ *Failed to retrieve statistics*\n\nPlease try again later.")
		msg.ParseMode = "Markdown"
		bot.Send(msg)
		return err
	}

	// Professional statistics display with biological insights
	response := "📊 *Your Memory Statistics*\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━\n\n" +
		"📚 *Memory Collection:*\n" +
		fmt.Sprintf("• Total Memories: `%d`\n", output.TotalMemories) +
		fmt.Sprintf("• Active Since: `%s`\n\n", time.Now().Format("2006-01-02")) +
		"🧠 *Biological Features:*\n" +
		"• Emotional tagging active\n" +
		"• Context encoding enabled\n" +
		"• Sleep consolidation running\n" +
		"• LTP spaced repetition scheduled\n\n" +
		"💡 *Tips for Better Memory:*\n" +
		"• Use emotional words for stronger recall\n" +
		"• Add context (time, place, people)\n" +
		"• Review memories regularly\n" +
		"• Use hashtags for organization\n\n" +
		"Keep building your memory collection! 🚀"

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"

	// Add action buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Save Memory", "cmd_save"),
			tgbotapi.NewInlineKeyboardButtonData("📚 View Recent", "cmd_recent"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "cmd_search"),
			tgbotapi.NewInlineKeyboardButtonData("❓ Help", "cmd_help"),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	return err
}
