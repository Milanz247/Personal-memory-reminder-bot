#!/bin/bash

# Run Memory Bot with FTS5 support

echo "🚀 Starting Memory Bot with FTS5 support..."
echo ""

# Run with fts5 tag
CGO_ENABLED=1 go run -tags "fts5" main.go
