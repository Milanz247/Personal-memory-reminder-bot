#!/bin/bash

# Run script for Biological Memory Bot
# This script starts the bot with proper environment setup

set -e  # Exit on error

echo "🤖 Starting Biological Memory Bot..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Stop any running instances
echo "${YELLOW}🛑 Stopping any running instances...${NC}"
pkill -f memory-bot 2>/dev/null || true
sleep 1

# Check if binary exists, rebuild if needed
if [ ! -f "memory-bot" ]; then
    echo "${YELLOW}⚠️  Binary not found. Building...${NC}"
    go build -tags "fts5" -o memory-bot cmd/bot/main.go
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
    echo "${GREEN}✅ Build successful!${NC}"
fi

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "${YELLOW}⚠️  .env file not found!${NC}"
    echo "Creating .env from .env.example..."
    cp .env.example .env
    echo ""
    echo "${YELLOW}Please edit .env file and add your Telegram Bot Token${NC}"
    echo "Then run this script again."
    exit 1
fi

# Source environment variables
export $(cat .env | grep -v '^#' | xargs)

# Check if TELEGRAM_BOT_TOKEN is set
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "❌ TELEGRAM_BOT_TOKEN not set in .env file"
    echo "Please add your bot token to .env file"
    exit 1
fi

echo "${GREEN}✅ Environment loaded${NC}"
echo ""
echo "🧠 Biological Features Active:"
echo "  • Amygdala Emotional Tagging ✓"
echo "  • Hippocampus Context Encoding ✓"
echo "  • Sleep Consolidation (2 AM daily) ✓"
echo "  • LTP Spaced Repetition ✓"
echo "  • Ebbinghaus Forgetting Curve ✓"
echo ""
echo "${BLUE}🚀 Starting bot...${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Run the bot
./memory-bot
