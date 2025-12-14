package command

import (
	"context"
	"fmt"
	"log"

	"memory-bot/internal/application/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// RecentCommand handles the /recent command
type RecentCommand struct {
	useCase *usecase.GetRecentMemoriesUseCase
}

// NewRecentCommand creates a new recent command
func NewRecentCommand(useCase *usecase.GetRecentMemoriesUseCase) *RecentCommand {
	return &RecentCommand{
		useCase: useCase,
	}
}

// Name returns the command name
func (c *RecentCommand) Name() string {
	return "recent"
}

// Description returns the command description
func (c *RecentCommand) Description() string {
	return "Recent memories"
}

// Execute executes the recent command
func (c *RecentCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
	input := usecase.GetRecentMemoriesInput{
		UserID: message.From.ID,
		Limit:  10,
	}

	output, err := c.useCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Error getting recent memories: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to retrieve memories.")
		bot.Send(msg)
		return err
	}

	if len(output.Memories) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You don't have any memories yet. Start saving some!")
		_, err := bot.Send(msg)
		return err
	}

	response := "📋 *Your Recent Memories:*\n\n"

	for i, mem := range output.Memories {
		content := mem.Content
		if len(content) > 100 {
			content = content[:100] + "..."
		}

		response += fmt.Sprintf("%d. %s\n\n", i+1, content)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"
	_, err = bot.Send(msg)
	return err
}
