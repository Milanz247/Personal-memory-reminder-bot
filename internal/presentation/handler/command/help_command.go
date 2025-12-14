package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HelpCommand handles the /help command
type HelpCommand struct{}

// NewHelpCommand creates a new help command
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

// Name returns the command name
func (c *HelpCommand) Name() string {
	return "help"
}

// Description returns the command description
func (c *HelpCommand) Description() string {
	return "Show help"
}

// Execute executes the help command
func (c *HelpCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
	helpText := `📚 *Available Commands:*

*Saving Memories:*
Tap the button below to save a memory

*Searching:*
Tap the button below to search memories

*Other Commands:*
/recent - Show recent memories
/stats - Show statistics

*Tips:*
• Just type text to save it as a memory
• Use #tags to organize memories
• Search with partial words (e.g., "tele" finds "telegram")`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"

	// Add action buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 Save Memory", "action:save"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "action:search"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📋 Recent", "action:recent"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Stats", "action:stats"),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	return err
}
