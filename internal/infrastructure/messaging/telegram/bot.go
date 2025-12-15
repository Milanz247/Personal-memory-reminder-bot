package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	"memory-bot/internal/application/usecase"
	"memory-bot/internal/presentation/handler/command"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents the Telegram bot adapter
type Bot struct {
	api          *tgbotapi.BotAPI
	registry     *command.CommandRegistry
	saveUseCase  *usecase.SaveMemoryUseCase
	userStates   map[int64]string
	userMessages map[int64]*tgbotapi.Message
}

// NewBot creates a new Telegram bot instance
func NewBot(
	token string,
	registry *command.CommandRegistry,
	saveUseCase *usecase.SaveMemoryUseCase,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	bot := &Bot{
		api:          api,
		registry:     registry,
		saveUseCase:  saveUseCase,
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
	commands := []tgbotapi.BotCommand{}

	for _, cmd := range b.registry.GetAll() {
		commands = append(commands, tgbotapi.BotCommand{
			Command:     cmd.Name(),
			Description: cmd.Description(),
		})
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
	ctx := context.Background()

	if message.IsCommand() {
		b.handleCommand(ctx, message)
		return
	}

	userID := message.From.ID

	// Check if user has a pending action
	if state, exists := b.userStates[userID]; exists {
		switch state {
		case "search":
			delete(b.userStates, userID)
			// Execute search command
			if searchCmd, ok := b.registry.Get("search"); ok {
				searchMsg := &tgbotapi.Message{
					From: message.From,
					Chat: message.Chat,
					Text: "/search " + message.Text,
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Offset: 0, Length: 7},
					},
				}
				searchCmd.Execute(ctx, b.api, searchMsg)
			}
			return
		case "save":
			delete(b.userStates, userID)
			b.saveMemoryFromText(ctx, message)
			return
		}
	}

	// If no pending state, ask user what to do
	b.userMessages[userID] = message

	response := fmt.Sprintf("What do you want to do with:\n\n`%s`", message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "Markdown"

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
func (b *Bot) handleCommand(ctx context.Context, message *tgbotapi.Message) {
	cmdName := message.Command()

	err := b.registry.Execute(ctx, cmdName, b.api, message)
	if err != nil {
		if err == command.ErrCommandNotFound {
			b.sendMessage(message.Chat.ID, "Unknown command. Use /help to see available commands.")
		} else {
			log.Printf("Error executing command %s: %v", cmdName, err)
		}
	}
}

// handleCallbackQuery handles callback queries from inline keyboards
func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	ctx := context.Background()

	parts := strings.Split(query.Data, "_")
	if len(parts) < 2 {
		// Try old format with colon
		parts = strings.Split(query.Data, ":")
		if len(parts) < 2 {
			b.api.Send(tgbotapi.NewCallback(query.ID, "Invalid request"))
			return
		}
	}

	action := parts[0]

	switch action {
	case "cmd":
		// Handle command buttons (from start, save, stats)
		cmdName := parts[1]
		if cmd, ok := b.registry.Get(cmdName); ok {
			msg := &tgbotapi.Message{
				Chat: query.Message.Chat,
				From: query.From,
			}
			cmd.Execute(ctx, b.api, msg)
			b.api.Send(tgbotapi.NewCallback(query.ID, ""))
		} else {
			b.api.Send(tgbotapi.NewCallback(query.ID, "Command not found"))
		}

	case "dosave":
		userID := query.From.ID
		if msg, exists := b.userMessages[userID]; exists {
			b.saveMemoryFromText(ctx, msg)
			delete(b.userMessages, userID)
			deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
			b.api.Send(deleteMsg)
		}
		b.api.Send(tgbotapi.NewCallback(query.ID, "Saved!"))

	case "dosearch":
		userID := query.From.ID
		if msg, exists := b.userMessages[userID]; exists {
			if searchCmd, ok := b.registry.Get("search"); ok {
				searchMsg := &tgbotapi.Message{
					From: msg.From,
					Chat: msg.Chat,
					Text: "/search " + msg.Text,
					Entities: []tgbotapi.MessageEntity{
						{Type: "bot_command", Offset: 0, Length: 7},
					},
				}
				searchCmd.Execute(ctx, b.api, searchMsg)
			}
			delete(b.userMessages, userID)
			deleteMsg := tgbotapi.NewDeleteMessage(query.Message.Chat.ID, query.Message.MessageID)
			b.api.Send(deleteMsg)
		}
		b.api.Send(tgbotapi.NewCallback(query.ID, "Searching..."))

	case "action":
		b.handleActionButton(ctx, query, parts[1])

	default:
		b.api.Send(tgbotapi.NewCallback(query.ID, "Unknown action"))
	}
}

// handleActionButton handles action button clicks
func (b *Bot) handleActionButton(ctx context.Context, query *tgbotapi.CallbackQuery, actionType string) {
	chatID := query.Message.Chat.ID

	switch actionType {
	case "save":
		b.userStates[query.From.ID] = "save"
		response := "💾 *Save a Memory*\n\nJust type your memory and send it. I'll save it automatically!\n\n*Example:*\n`Remember to call doctor tomorrow #health`\n\n*Tips:*\n• Use #hashtags to organize\n• Be specific and descriptive\n• No need to use /save command"
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
		b.api.Send(tgbotapi.NewCallback(query.ID, "Send your memory as next message"))

	case "search":
		b.userStates[query.From.ID] = "search"
		response := "🔍 *Search Memories*\n\nJust type your search keywords and send. I'll find matching memories!\n\n*Examples:*\n• `Milan` - find memories with \"Milan\"\n• `doctor health` - find memories with both words\n• `#work` - find all work memories\n\n*Tips:*\n• Use partial words (\"tele\" finds \"telegram\")\n• Multiple words search together\n• No need to use /search command"
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown"
		b.api.Send(msg)
		b.api.Send(tgbotapi.NewCallback(query.ID, "Send search keywords as next message"))

	case "recent":
		if recentCmd, ok := b.registry.Get("recent"); ok {
			msg := &tgbotapi.Message{
				Chat: query.Message.Chat,
				From: query.From,
			}
			recentCmd.Execute(ctx, b.api, msg)
		}
		b.api.Send(tgbotapi.NewCallback(query.ID, ""))

	case "stats":
		if statsCmd, ok := b.registry.Get("stats"); ok {
			msg := &tgbotapi.Message{
				Chat: query.Message.Chat,
				From: query.From,
			}
			statsCmd.Execute(ctx, b.api, msg)
		}
		b.api.Send(tgbotapi.NewCallback(query.ID, ""))

	default:
		b.api.Send(tgbotapi.NewCallback(query.ID, "Unknown action"))
	}
}

// saveMemoryFromText saves a memory from a text message
func (b *Bot) saveMemoryFromText(ctx context.Context, message *tgbotapi.Message) {
	input := usecase.SaveMemoryInput{
		UserID:  message.From.ID,
		ChatID:  message.Chat.ID,
		Content: message.Text,
	}

	output, err := b.saveUseCase.Execute(ctx, input)
	if err != nil {
		log.Printf("Error saving memory: %v", err)
		b.sendMessage(message.Chat.ID, "❌ Failed to save memory. Please try again.")
		return
	}

	response := "✅ Memory saved!"
	if len(output.Tags) > 0 {
		response += fmt.Sprintf(" Tags: %s", strings.Join(output.Tags, ", "))
	}

	b.sendMessage(message.Chat.ID, response)
}

// sendMessage is a helper function to send text messages
func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.api.Send(msg)
}
