package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"memory-bot/internal/application/usecase"
	"memory-bot/internal/infrastructure/messaging/telegram"
	"memory-bot/internal/infrastructure/persistence/sqlite"
	"memory-bot/internal/infrastructure/scheduler"
	"memory-bot/internal/infrastructure/search/strategy"
	"memory-bot/internal/presentation/handler/command"
	"memory-bot/pkg/config"
	"memory-bot/pkg/encryption"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Memory Storage Bot with Clean Architecture...")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dbConn, err := sqlite.NewConnection(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	// Initialize encryptor if encryption key is provided
	var encryptor *encryption.Encryptor
	if cfg.EncryptionKey != "" {
		encryptor = encryption.NewEncryptor(cfg.EncryptionKey)
		log.Println("🔒 Encryption enabled for sensitive memory data")
	} else {
		log.Println("⚠️  Warning: Encryption is disabled. Set ENCRYPTION_KEY environment variable to enable encryption.")
	}

	// Initialize repository
	memoryRepo := sqlite.NewMemoryRepository(dbConn, encryptor)

	// Initialize use cases
	saveMemoryUC := usecase.NewSaveMemoryUseCase(memoryRepo)
	getRecentUC := usecase.NewGetRecentMemoriesUseCase(memoryRepo)
	getStatsUC := usecase.NewGetStatsUseCase(memoryRepo)
	reviewMemoryUC := usecase.NewReviewMemoryUseCase(memoryRepo)

	// Initialize search strategy (Smart Search)
	searchStrategy := strategy.NewSmartSearchStrategy(memoryRepo)
	searchMemoryUC := usecase.NewSearchMemoryUseCase(searchStrategy)

	// Initialize command registry
	registry := command.NewCommandRegistry()

	// Register commands
	registry.Register(command.NewStartCommand())
	registry.Register(command.NewHelpCommand())
	registry.Register(command.NewSaveCommand(saveMemoryUC))
	registry.Register(command.NewSearchCommand(searchMemoryUC))
	registry.Register(command.NewRecentCommand(getRecentUC))
	registry.Register(command.NewStatsCommand(getStatsUC))

	// Create Telegram bot
	bot, err := telegram.NewBot(cfg.TelegramBotToken, registry, saveMemoryUC)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Create bot API for scheduler
	botAPI, err := createTelegramBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("Failed to create bot API for scheduler: %v", err)
	}

	// Initialize spaced repetition scheduler
	sr := scheduler.NewSpacedRepetitionScheduler(botAPI, reviewMemoryUC, cfg.ReviewIntervals)
	sr.Start()
	defer sr.Stop()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start bot in a goroutine
	go bot.Start()

	log.Println("✅ Bot is running with Clean Architecture!")
	log.Println("📦 Design Patterns: Repository, Strategy, Command, Observer, Dependency Injection")

	// Wait for interrupt signal
	<-sigChan
	log.Println("\n🛑 Shutting down gracefully...")
}

// createTelegramBotAPI creates a Telegram bot API instance for the scheduler
func createTelegramBotAPI(token string) (*tgbotapi.BotAPI, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return api, nil
}
