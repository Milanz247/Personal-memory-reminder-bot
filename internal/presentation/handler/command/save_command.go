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
		// Interactive template with structured guidance
		response := "📝 *Save a New Memory: Use the Template Below!*\n" +
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n" +
			"To maximize recall, structure your memory like this:\n\n" +
			"*🔑 Template:*\n" +
			"`[What I did or learned] [How I felt] #tag1 #tag2`\n\n" +
			"*📚 Optimal Examples:*\n\n" +
			"1️⃣ `I felt great excitement when I finally finished the complex database migration at work. #project #tech`\n\n" +
			"2️⃣ `Amazing breakthrough in my research today! Discovered a solution to the optimization problem. #research #achievement`\n\n" +
			"3️⃣ `Had a wonderful conversation with mom about childhood memories. Felt nostalgic and happy. #family #personal`\n\n" +
			"*🧠 Biological Analysis - What Happens:*\n" +
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
			"• *Amygdala Tagging:* Words like 'excitement', 'amazing', 'wonderful' boost emotional weight (0-100%)\n" +
			"• *Hippocampus Context:* Time/Day/Location automatically captured\n" +
			"• *Priority Calculation:* Higher emotion = Higher priority = Better recall\n" +
			"• *LTP Scheduling:* Smart review intervals (1→3→7→14→30 days)\n\n" +
			"*💡 Pro Tips for Better Memories:*\n" +
			"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
			"✅ Use emotional words → Better Amygdala tagging\n" +
			"✅ Add #hashtags → Easy organization & search\n" +
			"✅ Be specific → More context = Better retrieval\n" +
			"✅ Include feelings → Emotions strengthen memory\n\n" +
			"*📤 Ready? Send your memory text now!*"
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
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ *Failed to save memory*\n\nPlease try again or contact support if the problem persists.")
		msg.ParseMode = "Markdown"
		bot.Send(msg)
		return err
	}

	// Professional success message with biological features
	response := "✅ *Memory Saved Successfully!*\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━\n\n" +
		"📊 *Biological Analysis:*\n\n"

	// Emotional Weight with visual representation
	emotionalCategory := getEmotionalCategory(output.EmotionalWeight)
	emotionalBar := getEmotionalBar(output.EmotionalWeight)
	response += fmt.Sprintf("😊 *Emotional Weight:* %.0f%% %s\n", output.EmotionalWeight*100, emotionalBar)
	response += fmt.Sprintf("   Category: `%s`\n\n", emotionalCategory)

	// Context information
	if output.Context != "" {
		response += fmt.Sprintf("📍 *Context:* %s\n\n", output.Context)
	}

	// Tags if any
	if len(output.Tags) > 0 {
		response += fmt.Sprintf("🏷️ *Tags:* %s\n\n", strings.Join(output.Tags, " "))
	}

	// Memory ID for reference
	response += fmt.Sprintf("🆔 *Memory ID:* `%d`\n\n", output.MemoryID)

	// Review schedule info
	response += "🔄 *Next Steps:*\n" +
		"• Sleep consolidation will strengthen this memory tonight\n" +
		"• First review scheduled based on emotional weight\n" +
		"• Use /recent to see your latest memories"

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"

	// Add action buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📝 Save Another", "cmd_save"),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search", "cmd_search"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 My Stats", "cmd_stats"),
			tgbotapi.NewInlineKeyboardButtonData("📚 Recent", "cmd_recent"),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	return err
}

// getEmotionalCategory returns the emotional category name
func getEmotionalCategory(weight float64) string {
	if weight < 0.3 {
		return "Neutral 😐"
	} else if weight < 0.6 {
		return "Moderate 🙂"
	} else if weight < 0.8 {
		return "Strong 😊"
	}
	return "Intense 🤩"
}

// getEmotionalBar returns a visual bar representation
func getEmotionalBar(weight float64) string {
	filled := int(weight * 10)
	bar := ""
	for i := 0; i < 10; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
