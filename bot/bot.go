package bot

import (
	"fmt"
	"log"
	"memory-bot/database"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	PageSize = 5 // Number of results per page
)

// Bot represents the Telegram bot
type Bot struct {
	api          *tgbotapi.BotAPI
	db           *database.Database
	userStates   map[int64]string            // Track pending user actions ("save" or "search")
	userMessages map[int64]*tgbotapi.Message // Cache messages for save/search decision
}

// NewBot creates a new Telegram bot instance
func NewBot(token string, db *database.Database) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	bot := &Bot{
		api:          api,
		db:           db,
		userStates:   make(map[int64]string),
		userMessages: make(map[int64]*tgbotapi.Message),
	}

	// Set bot commands menu
	if err := bot.setCommands(); err != nil {
		log.Printf("Warning: Failed to set bot commands: %v", err)
	}

	return bot, nil
}

// setCommands sets the bot command menu
func (b *Bot) setCommands() error {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Welcome message"},
		{Command: "save", Description: "Save a memory"},
		{Command: "search", Description: "Search memories"},
		{Command: "recent", Description: "Recent memories"},
		{Command: "stats", Description: "Statistics"},
		{Command: "help", Description: "Show help"},
	}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	_, err := b.api.Request(cfg)
	if err != nil {
		return fmt.Errorf("failed to set commands: %w", err)
	}

	log.Println("Bot commands menu configured successfully")
	return nil
}

// Start starts the bot and handles incoming messages
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	log.Println("Bot started. Listening for messages...")

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.handleCallbackQuery(update.CallbackQuery)
		}
	}
}

// handleMessage handles incoming text messages
func (b *Bot) handleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		b.handleCommand(message)
		return
	}

	userID := message.From.ID

	// Check if user has a pending action
	if state, exists := b.userStates[userID]; exists {
		switch state {
		case "search":
			// User clicked search button, now searching with this text
			delete(b.userStates, userID)
			b.performSearch(message.Chat.ID, userID, message.Text, 0, 0)
			return
		case "save":
			// User clicked save button, now saving this text
			delete(b.userStates, userID)
			b.saveMemory(message)
			return
		}
	}

	// If no pending state, ask user what to do
	// Cache the message for later use
	b.userMessages[userID] = message

	response := fmt.Sprintf("What do you want to do with:\n\n`%s`", message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"

	// Add buttons for save or search
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💾 Save", fmt.Sprintf("dosave:%d", userID)),
			tgbotapi.NewInlineKeyboardButtonData("🔍 Search", fmt.Sprintf("dosearch:%d", userID)),
		),
	)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

// handleCommand handles bot commands
func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.handleStart(message)
	case "save":
		b.handleSaveCommand(message)
	case "search":
		b.handleSearchCommand(message, 0)
	case "recent":
		b.handleRecentCommand(message)
	case "stats":
		b.handleStatsCommand(message)
	case "help":
		b.handleHelp(message)
	default:
		b.sendMessage(message.Chat.ID, "Unknown command. Use /help to see available commands.")
	}
}

// handleStart handles the /start command
func (b *Bot) handleStart(message *tgbotapi.Message) {
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

	b.api.Send(msg)
}

// handleHelp handles the /help command
func (b *Bot) handleHelp(message *tgbotapi.Message) {
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

	b.api.Send(msg)
}

// saveMemory saves a memory from a message
func (b *Bot) saveMemory(message *tgbotapi.Message) {
	content := message.Text
	if content == "" {
		return
	}

	// Extract tags from content
	tags := extractTags(content)

	userID := message.From.ID
	chatID := message.Chat.ID

	_, err := b.db.SaveMemory(userID, chatID, content, strings.Join(tags, " "))
	if err != nil {
		log.Printf("Error saving memory: %v", err)
		b.sendMessage(chatID, "❌ Failed to save memory. Please try again.")
		return
	}

	response := "✅ Memory saved!"
	if len(tags) > 0 {
		response += fmt.Sprintf(" Tags: %s", strings.Join(tags, ", "))
	}

	b.sendMessage(chatID, response)
}

// handleSaveCommand handles the /save command
func (b *Bot) handleSaveCommand(message *tgbotapi.Message) {
	args := message.CommandArguments()
	if args != "" {
		// If text provided with command, save it directly
		msg := &tgbotapi.Message{
			From: message.From,
			Chat: message.Chat,
			Text: args,
		}
		b.saveMemory(msg)
	} else {
		// Set state to wait for next message
		b.userStates[message.From.ID] = "save"
		response := "💾 *Save a Memory*\n\nSend your memory text now:\n\n*Example:*\n`Doctor appointment tomorrow 3 PM #health`\n\n*Tips:*\n• Use #hashtags to organize\n• Be specific"
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
	}
}

// handleSearchCommand handles the /search command
func (b *Bot) handleSearchCommand(message *tgbotapi.Message, page int) {
	keyword := message.CommandArguments()
	if keyword != "" {
		// If keyword provided with command, search directly
		b.performSearch(message.Chat.ID, message.From.ID, keyword, page, 0)
	} else {
		// Set state to wait for next message
		b.userStates[message.From.ID] = "search"
		response := "🔍 *Search Memories*\n\nSend your search keywords now:\n\n*Examples:*\n• `Milan` - find memories with \"Milan\"\n• `doctor health` - both words\n• `#work` - all work memories\n\n*Tips:*\n• Partial words work (\"tele\" finds \"telegram\")"
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
	}
}

// performSearch performs a search and sends results with pagination
func (b *Bot) performSearch(chatID, userID int64, keyword string, page int, editMessageID int) {
	offset := page * PageSize

	// Use smart search for better fragment matching
	var memories []database.Memory
	var err error

	if page == 0 {
		// First page: use SmartSearch with fallback strategies
		memories, err = b.db.SmartSearch(userID, keyword, PageSize+1)
	} else {
		// Subsequent pages: use standard search with offset
		memories, err = b.db.SearchAndRankMemories(userID, keyword, PageSize+1, offset)
	}

	if err != nil {
		log.Printf("Error searching memories: %v", err)
		b.sendMessage(chatID, "❌ Search failed. Please try again.")
		return
	}

	if len(memories) == 0 {
		response := fmt.Sprintf("🔍 No memories found for: `%s`\n\n💡 *Tips:*\n• Try partial words (e.g., \"tele\" finds \"telegram\")\n• Use fewer words\n• Check spelling", keyword)
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
		return
	}

	// Check if there are more results
	hasMore := len(memories) > PageSize
	if hasMore {
		memories = memories[:PageSize]
	}

	// Format results - simple and clean
	response := fmt.Sprintf("🔍 *Search Results* – %s\n\n━━━━━━━━━━━━━━━\n", keyword)

	for i, mem := range memories {
		// Truncate long memories
		content := mem.TextContent
		if len(content) > 200 {
			content = content[:200] + "..."
		}

		// Format number with emoji
		numEmoji := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}
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

	response += fmt.Sprintf("━━━━━━━━━━━━━━━\n\n📌 *Total Results:* %d", len(memories))

	log.Printf("Sending search results to chat %d: %d memories found", chatID, len(memories))

	if editMessageID != 0 {
		// Edit existing message
		msg := tgbotapi.NewEditMessageText(chatID, editMessageID, response)
		msg.ParseMode = "Markdown"

		// Only add keyboard if there are buttons
		if page > 0 || hasMore {
			var buttons []tgbotapi.InlineKeyboardButton
			if page > 0 {
				prevCallback := fmt.Sprintf("search:%s:%d", keyword, page-1)
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("⏪ Previous", prevCallback))
			}
			if hasMore {
				nextCallback := fmt.Sprintf("search:%s:%d", keyword, page+1)
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Next ⏩", nextCallback))
			}
			keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
			msg.ReplyMarkup = &keyboard
		}

		_, err := b.api.Send(msg)
		if err != nil {
			log.Printf("Error editing message: %v", err)
		}
	} else {
		// Send new message
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown"

		// Only add keyboard if there are buttons
		if page > 0 || hasMore {
			var buttons []tgbotapi.InlineKeyboardButton
			if page > 0 {
				prevCallback := fmt.Sprintf("search:%s:%d", keyword, page-1)
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("⏪ Previous", prevCallback))
			}
			if hasMore {
				nextCallback := fmt.Sprintf("search:%s:%d", keyword, page+1)
				buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Next ⏩", nextCallback))
			}
			keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
			msg.ReplyMarkup = keyboard
		}

		_, err := b.api.Send(msg)
		if err != nil {
			log.Printf("Error sending search results: %v", err)
		}
	}
}

// handleRecentCommand handles the /recent command
func (b *Bot) handleRecentCommand(message *tgbotapi.Message) {
	memories, err := b.db.GetRecentMemories(message.From.ID, 10)
	if err != nil {
		log.Printf("Error getting recent memories: %v", err)
		b.sendMessage(message.Chat.ID, "❌ Failed to retrieve memories.")
		return
	}

	if len(memories) == 0 {
		b.sendMessage(message.Chat.ID, "You don't have any memories yet. Start saving some!")
		return
	}

	response := "📋 *Your Recent Memories:*\n\n"

	for i, mem := range memories {
		content := mem.TextContent
		if len(content) > 100 {
			content = content[:100] + "..."
		}

		response += fmt.Sprintf("%d. %s\n\n",
			i+1,
			content)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

// handleStatsCommand handles the /stats command
func (b *Bot) handleStatsCommand(message *tgbotapi.Message) {
	count, err := b.db.GetMemoryCount(message.From.ID)
	if err != nil {
		log.Printf("Error getting memory count: %v", err)
		b.sendMessage(message.Chat.ID, "❌ Failed to retrieve statistics.")
		return
	}

	response := fmt.Sprintf(`📊 *Your Memory Statistics:*

Total Memories: %d
Bot Active Since: %s

Keep building your memory collection! 🧠`, count, time.Now().Format("2006-01-02"))

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

// handleCallbackQuery handles callback queries from inline keyboards
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	// Parse callback data
	parts := strings.Split(query.Data, ":")
	if len(parts) < 2 {
		b.api.Send(tgbotapi.NewCallback(query.ID, "Invalid request"))
		return
	}

	action := parts[0]

	switch action {
	case "dosave":
		// User chose to save the message
		userID := query.From.ID
		if msg, exists := b.userMessages[userID]; exists {
			b.saveMemory(msg)
			delete(b.userMessages, userID)
			// Delete the "What do you want to do" message
			deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
			b.api.Send(deleteMsg)
		}
		b.api.Send(tgbotapi.NewCallback(query.ID, "Saved!"))

	case "dosearch":
		// User chose to search with the message
		userID := query.From.ID
		if msg, exists := b.userMessages[userID]; exists {
			b.performSearch(msg.Chat.ID, userID, msg.Text, 0, 0)
			delete(b.userMessages, userID)
			// Delete the "What do you want to do" message
			deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
			b.api.Send(deleteMsg)
		}
		b.api.Send(tgbotapi.NewCallback(query.ID, "Searching..."))

	case "action":
		// Handle action buttons (save, search, recent, stats)
		b.handleActionButton(query, parts[1])
	case "search":
		// Handle search pagination
		if len(parts) != 3 {
			b.api.Send(tgbotapi.NewCallback(query.ID, "Invalid request"))
			return
		}
		keyword := parts[1]
		page, err := strconv.Atoi(parts[2])
		if err != nil {
			b.api.Send(tgbotapi.NewCallback(query.ID, "Invalid page number"))
			return
		}
		b.performSearch(query.Message.Chat.ID, query.From.ID, keyword, page, query.Message.MessageID)
		b.api.Send(tgbotapi.NewCallback(query.ID, ""))
	default:
		b.api.Send(tgbotapi.NewCallback(query.ID, "Unknown action"))
	}
}

// handleActionButton handles action button clicks
func (b *Bot) handleActionButton(query *tgbotapi.CallbackQuery, actionType string) {
	chatID := query.Message.Chat.ID

	switch actionType {
	case "save":
		// Set user state to expect save input
		b.userStates[query.From.ID] = "save"
		// Show save help and wait for next message
		response := "💾 *Save a Memory*\n\nJust type your memory and send it. I'll save it automatically!\n\n*Example:*\n`Remember to call doctor tomorrow #health`\n\n*Tips:*\n• Use #hashtags to organize\n• Be specific and descriptive\n• No need to use /save command"
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
		b.api.Send(tgbotapi.NewCallback(query.ID, "Send your memory as next message"))

	case "search":
		// Set user state to expect search input
		b.userStates[query.From.ID] = "search"
		// Show search help and wait for next message
		response := "🔍 *Search Memories*\n\nJust type your search keywords and send. I'll find matching memories!\n\n*Examples:*\n• `Milan` - find memories with \"Milan\"\n• `doctor health` - find memories with both words\n• `#work` - find all work memories\n\n*Tips:*\n• Use partial words (\"tele\" finds \"telegram\")\n• Multiple words search together\n• No need to use /search command"
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
		b.api.Send(tgbotapi.NewCallback(query.ID, "Send search keywords as next message"))

	case "recent":
		// Show recent memories
		b.handleRecentCommand(&tgbotapi.Message{
			Chat: query.Message.Chat,
			From: query.From,
		})
		b.api.Send(tgbotapi.NewCallback(query.ID, ""))

	case "stats":
		// Show statistics
		b.handleStatsCommand(&tgbotapi.Message{
			Chat: query.Message.Chat,
			From: query.From,
		})
		b.api.Send(tgbotapi.NewCallback(query.ID, ""))

	case "help":
		// Show help
		b.handleHelp(&tgbotapi.Message{
			Chat: query.Message.Chat,
			From: query.From,
		})
		b.api.Send(tgbotapi.NewCallback(query.ID, ""))

	default:
		b.api.Send(tgbotapi.NewCallback(query.ID, "Unknown action"))
	}
}

// sendMessage is a helper function to send text messages
func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}

// extractTags extracts hashtags from text
func extractTags(text string) []string {
	var tags []string
	words := strings.Fields(text)

	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	return tags
}
