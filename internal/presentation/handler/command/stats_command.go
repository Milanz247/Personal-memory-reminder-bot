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
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Failed to retrieve statistics.")
		bot.Send(msg)
		return err
	}

	response := fmt.Sprintf(`📊 *Your Memory Statistics:*

Total Memories: %d
Bot Active Since: %s

Keep building your memory collection! 🧠`, output.TotalMemories, time.Now().Format("2006-01-02"))

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"
	_, err = bot.Send(msg)
	return err
}
