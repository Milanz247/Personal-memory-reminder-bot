package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartCommand handles the /start command
type StartCommand struct{}

// NewStartCommand creates a new start command
func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

// Name returns the command name
func (c *StartCommand) Name() string {
	return "start"
}

// Description returns the command description
func (c *StartCommand) Description() string {
	return "Welcome message"
}

// Execute executes the start command
func (c *StartCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
	imagePath := "assets/images/welcome_banner.png"

	// Professional welcome with biological features
	caption := `🧠 *Biological Memory System*
━━━━━━━━━━━━━━━━━━━━━━━

Welcome to your intelligent personal memory assistant powered by neuroscience principles!

🔬 *Based on Brain Science:*
• 😊 Amygdala - Emotional tagging
• 🧮 Hippocampus - Context encoding
• 💤 Sleep consolidation
• 🔄 LTP spaced repetition
• 📉 Ebbinghaus forgetting curve

✨ *Smart Features:*
• Emotional weight analysis (0-100%)
• Time & day context capture
• Priority score calculation
• Intelligent search with FTS5
• Encrypted storage
• Automatic review scheduling

🚀 *Quick Start:*
/save - Save memories with emotion
/search - Smart contextual search
/recent - View latest memories
/stats - Memory statistics
/help - Detailed help

━━━━━━━━━━━━━━━━━━━━━━━
👨‍💻 *Created by:* Milan Madusanka
🔗 [GitHub](https://github.com/Milanz247)`

	if fileExists(imagePath) {
		photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FilePath(imagePath))
		photo.Caption = caption
		photo.ParseMode = "Markdown"

		// Add quick action buttons
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💾 Save Memory", "cmd_save"),
				tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "cmd_search"),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📊 Stats", "cmd_stats"),
				tgbotapi.NewInlineKeyboardButtonData("❓ Help", "cmd_help"),
			),
		)
		photo.ReplyMarkup = keyboard

		_, err := bot.Send(photo)
		return err
	}

	// Fallback if image not found
	msg := tgbotapi.NewMessage(message.Chat.ID, caption)
	msg.ParseMode = "Markdown"

	// Add quick action buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 Save Memory", "cmd_save"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "cmd_search"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Stats", "cmd_stats"),
			tgbotapi.NewInlineKeyboardButtonData("❓ Help", "cmd_help"),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	return err
}
