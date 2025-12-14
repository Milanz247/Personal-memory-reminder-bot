#!/bin/bash

# Stop any running instances of memory-bot
echo "🛑 Stopping any running memory-bot instances..."
pkill -f memory-bot
sleep 1

# Check if any instances are still running
if pgrep -f memory-bot > /dev/null; then
    echo "⚠️  Warning: Some instances may still be running. Force killing..."
    pkill -9 -f memory-bot
    sleep 1
fi

# Verify all instances are stopped
if pgrep -f memory-bot > /dev/null; then
    echo "❌ Error: Could not stop all instances"
    exit 1
else
    echo "✅ All memory-bot instances stopped"
fi
