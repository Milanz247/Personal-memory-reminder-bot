package command

import (
	"context"
	"fmt"
	"log"

	"memory-bot/internal/application/usecase"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const PageSize = 5

// SearchCommand handles the /search command
type SearchCommand struct {
	useCase *usecase.SearchMemoryUseCase
}

// NewSearchCommand creates a new search command
func NewSearchCommand(useCase *usecase.SearchMemoryUseCase) *SearchCommand {
	return &SearchCommand{
		useCase: useCase,
	}
}

// Name returns the command name
func (c *SearchCommand) Name() string {
	return "search"
}

// Description returns the command description
func (c *SearchCommand) Description() string {
	return "Search memories"
}

// Execute executes the search command
func (c *SearchCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
	keyword := message.CommandArguments()

	if keyword == "" {
		// Prompt user for input
		response := "🔍 *Search Memories*\n\nSend your search keywords now:\n\n*Examples:*\n• `Milan` - find memories with \"Milan\"\n• `doctor health` - both words\n• `#work` - all work memories\n\n*Tips:*\n• Partial words work (\"tele\" finds \"telegram\")"
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		return err
	}

	// Perform search
	input := usecase.SearchMemoryInput{
		UserID:  message.From.ID,
		Keyword: keyword,
		Limit:   PageSize,
		Offset:  0,
	}

	output, err := c.useCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Error searching memories: %v", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Search failed. Please try again.")
		bot.Send(msg)
		return err
	}

	if len(output.Memories) == 0 {
		response := fmt.Sprintf("🔍 No memories found for: `%s`\n\n💡 *Tips:*\n• Try partial words (e.g., \"tele\" finds \"telegram\")\n• Use fewer words\n• Check spelling", keyword)
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		msg.ParseMode = "Markdown"
		_, err := bot.Send(msg)
		return err
	}

	// Format results (no image, just text)
	response := fmt.Sprintf("🔍 *Search:* `%s`\n*Found:* %d\n\n━━━━━━━━━━━━━━━\n", keyword, len(output.Memories))

	numEmoji := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}
	for i, mem := range output.Memories {
		content := mem.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}

		numberDisplay := numEmoji[i]
		if i >= len(numEmoji) {
			numberDisplay = fmt.Sprintf("%d.", i+1)
		}

		response += fmt.Sprintf("%s %s\n🕒 %s – %s\n\n",
			numberDisplay,
			content,
			mem.CreatedAt.Format("2006-01-02"),
			mem.CreatedAt.Format("03:04 PM"))
	}

	response += fmt.Sprintf("━━━━━━━━━━━━━━━\n\n📌 *Total Results:* %d", len(output.Memories))

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"

	// Add pagination if needed
	if output.HasMore {
		nextCallback := fmt.Sprintf("search:%s:%d", keyword, 1)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next ⏩", nextCallback),
			),
		)
		msg.ReplyMarkup = keyboard
	}

	_, err = bot.Send(msg)
	return err
}
