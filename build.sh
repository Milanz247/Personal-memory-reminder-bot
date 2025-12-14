#!/bin/bash

# Build Memory Bot with FTS5 support

echo "🔨 Building Memory Bot with FTS5 support..."
echo ""

# Enable CGO and build with fts5 tag
CGO_ENABLED=1 go build -tags "fts5" -o memory-bot

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Build successful!"
    echo ""
    echo "Run with: ./memory-bot"
else
    echo ""
    echo "❌ Build failed"
    exit 1
fi
