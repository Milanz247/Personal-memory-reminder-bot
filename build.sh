#!/bin/bash

# Build Memory Bot with FTS5 support

echo "🛑 Stopping any running instances..."
pkill -f memory-bot 2>/dev/null
sleep 1

echo "🔨 Building Memory Bot with FTS5 support..."
echo ""

# Enable CGO and build with fts5 tag
go build -tags "fts5" -o memory-bot cmd/bot/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "📦 Binary created: memory-bot"
    echo ""
    echo "To run: ./memory-bot"
else
    echo "❌ Build failed!"
    exit 1
fi
