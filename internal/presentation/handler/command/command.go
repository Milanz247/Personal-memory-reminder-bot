package command

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Command represents a bot command interface (Command Pattern)
type Command interface {
	// Name returns the command name
	Name() string

	// Description returns the command description
	Description() string

	// Execute executes the command
	Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error
}

// BotAPI defines the interface for bot operations
type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

// CommandRegistry manages and executes commands
type CommandRegistry struct {
	commands map[string]Command
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

// Register registers a new command
func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Get retrieves a command by name
func (r *CommandRegistry) Get(name string) (Command, bool) {
	cmd, exists := r.commands[name]
	return cmd, exists
}

// Execute executes a command by name
func (r *CommandRegistry) Execute(ctx context.Context, name string, bot BotAPI, message *tgbotapi.Message) error {
	cmd, exists := r.Get(name)
	if !exists {
		return ErrCommandNotFound
	}

	return cmd.Execute(ctx, bot, message)
}

// GetAll returns all registered commands
func (r *CommandRegistry) GetAll() []Command {
	commands := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, cmd)
	}
	return commands
}
