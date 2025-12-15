#!/bin/bash

# Build script for Biological Memory Bot
# This script compiles the Go application with required tags

set -e  # Exit on error

echo "🔨 Building Biological Memory Bot..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Stop any running instances
echo "${YELLOW}🛑 Stopping any running instances...${NC}"
pkill -f memory-bot 2>/dev/null || true
sleep 1

# Clean old binary
if [ -f "memory-bot" ]; then
    echo "${BLUE}🗑️  Removing old binary...${NC}"
    rm memory-bot
fi

# Build with FTS5 support
echo "${BLUE}📦 Compiling with SQLite FTS5 support...${NC}"
go build -tags "fts5" -o memory-bot cmd/bot/main.go

# Check if build was successful
if [ $? -eq 0 ]; then
    echo ""
    echo "${GREEN}✅ Build successful!${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Binary: memory-bot"
    echo "Size: $(du -h memory-bot | cut -f1)"
    echo ""
    echo "🧠 Biological Features:"
    echo "  • Amygdala Emotional Tagging"
    echo "  • Hippocampus Context Encoding"
    echo "  • Sleep Consolidation"
    echo "  • LTP Spaced Repetition"
    echo ""
    echo "To run the bot:"
    echo "  ${BLUE}./run.sh${NC}"
else
    echo ""
    echo "❌ Build failed!"
    exit 1
fi
