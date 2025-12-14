package main

import (
	"log"
	"memory-bot/bot"
	"memory-bot/config"
	"memory-bot/database"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Memory Storage Bot...")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create bot
	memoryBot, err := bot.NewBot(cfg.TelegramBotToken, db)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Initialize spaced repetition
	sr := bot.NewSpacedRepetition(memoryBot, db, cfg.ReviewIntervals)
	sr.Start()
	defer sr.Stop()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start bot in a goroutine
	go memoryBot.Start()

	// Wait for interrupt signal
	<-sigChan
	log.Println("\nShutting down gracefully...")
}
