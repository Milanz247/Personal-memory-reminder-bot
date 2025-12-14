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

	// Short, decorated caption with founder info
	caption := `🧠 *Memory Storage Bot*

Your intelligent personal memory assistant powered by AI.

✨ *Features*
• Smart search with FTS5
• Spaced repetition reminders
• Encrypted storage
• Tag organization

🚀 *Quick Start*
/save - Save memories
/search - Find anything
/help - Get help

━━━━━━━━━━━━━━━
👨‍💻 *Created by:* Milan Madusanka
🔗 [GitHub](https://github.com/Milanz247)
━━━━━━━━━━━━━━━`

	if fileExists(imagePath) {
		photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FilePath(imagePath))
		photo.Caption = caption
		photo.ParseMode = "Markdown"
		_, err := bot.Send(photo)
		return err
	}

	// Fallback if image not found
	msg := tgbotapi.NewMessage(message.Chat.ID, caption)
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	return err
}
