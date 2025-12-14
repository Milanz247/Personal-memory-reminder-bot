package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration
type Config struct {
	TelegramBotToken string
	DBPath           string
	ReviewIntervals  []int // in days
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	token := getEnv("TELEGRAM_BOT_TOKEN", "")
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	dbPath := getEnv("DB_PATH", "./memories.db")

	// Load review intervals
	intervals := []int{1, 3, 7, 14, 30} // default intervals
	if intervalsStr := os.Getenv("REVIEW_INTERVALS"); intervalsStr != "" {
		parsedIntervals, err := parseIntervals(intervalsStr)
		if err == nil && len(parsedIntervals) > 0 {
			intervals = parsedIntervals
		}
	}

	return &Config{
		TelegramBotToken: token,
		DBPath:           dbPath,
		ReviewIntervals:  intervals,
	}, nil
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseIntervals parses comma-separated interval string
func parseIntervals(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	intervals := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid interval value: %s", part)
		}
		intervals = append(intervals, val)
	}

	return intervals, nil
}
