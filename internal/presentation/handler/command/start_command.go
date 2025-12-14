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
	welcomeText := `🧠 *Welcome to Memory Storage Bot!*

This bot helps you store and retrieve memories using AI-powered search.

*Features:*
• Fast full-text search with relevance ranking
• Automatic spaced repetition for better retention
• Tag-based organization
• Context-aware storage

*Quick Start:*
Use the menu button (☰) at the bottom to access all commands:
• /save - Save a new memory
• /search - Search your memories
• /recent - View recent memories
• /stats - View statistics
• /help - Get detailed help

Or just type any text to save it as a memory!`

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ParseMode = "Markdown"

	_, err := bot.Send(msg)
	return err
}
