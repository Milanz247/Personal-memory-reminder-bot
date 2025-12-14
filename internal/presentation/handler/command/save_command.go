package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"memory-bot/internal/application/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SaveCommand handles the /save command
type SaveCommand struct {
	useCase *usecase.SaveMemoryUseCase
}

// NewSaveCommand creates a new save command
func NewSaveCommand(useCase *usecase.SaveMemoryUseCase) *SaveCommand {
	return &SaveCommand{
		useCase: useCase,
	}
}

// Name returns the command name
func (c *SaveCommand) Name() string {
	return "save"
}

// Description returns the command description
func (c *SaveCommand) Description() string {
	return "Save a memory"
}

// Execute executes the save command
func (c *SaveCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
	args := message.CommandArguments()

	if args == "" {
		// Prompt user for input
		response := "💾 *Save a Memory*\n\nSend your memory text now:\n\n*Example:*\n`Doctor appointment tomorrow 3 PM #health`\n\n*Tips:*\n• Use #hashtags to organize\n• Be specific"
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		return err
	}

	// Save the memory
	input := usecase.SaveMemoryInput{
		UserID:  message.From.ID,
		ChatID:  message.Chat.ID,
		Content: args,
	}

	output, err := c.useCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Error saving memory: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to save memory. Please try again.")
		bot.Send(msg)
		return err
	}

	response := "✅ Memory saved!"
	if len(output.Tags) > 0 {
		response += fmt.Sprintf(" Tags: %s", strings.Join(output.Tags, ", "))
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	_, err = bot.Send(msg)
	return err
}
