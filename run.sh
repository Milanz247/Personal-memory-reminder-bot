#!/bin/bash

# Run Memory Bot with FTS5 support#!/bin/bash

echo "🛑 Stopping any running instances..."
pkill -f memory-bot 2>/dev/null
sleep 1

echo "🔨 Rebuilding..."
go build -tags "fts5" -o memory-bot cmd/bot/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "🚀 Starting Memory Bot..."
    echo ""
    ./memory-bot
else
    echo "❌ Build failed!"
    exit 1
fi
