# 🧠 Personal Memory Reminder Bot

**Telegram Bot for Intelligent Memory Management with Biological Features**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)

A personal memory assistant with neuroscience-based features: emotional tagging, context encoding, spaced repetition, and full-text search.

---

## ✨ Features

### 🧠 Biological Memory System
- **Amygdala Emotional Tagging** - Analyzes emotional weight (0-100%)
- **Hippocampus Context Encoding** - Captures time, day, location
- **Sleep Consolidation** - Priority boost during rest periods
- **LTP Spaced Repetition** - Smart review scheduling (1,3,7,14,30 days)
- **Forgetting Curve Algorithm** - Prevents memory decay

### 🔍 Smart Search
- SQLite FTS5 full-text search with Porter stemming
- BM25 relevance ranking
- Multi-strategy fallback (FTS5 → AND → OR)
- Tag-based filtering (`#work`, `#health`)

### 🔒 Security
- AES-256-GCM encryption for sensitive data
- Searchable plaintext index (hybrid architecture)
- Optional encryption key

### 💬 Interactive Bot Interface
- Professional formatted messages
- Emotional analysis display with visual bars
- Context information (time of day, day of week)
- Interactive buttons for quick actions
- Memory statistics and insights

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- SQLite with FTS5
- Telegram Bot Token

### Installation

```bash
# Clone repository
git clone https://github.com/Milanz247/Personal-memory-reminder-bot.git
cd Personal-memory-reminder-bot

# Configure environment
cp .env.example .env
# Edit .env and add your TELEGRAM_BOT_TOKEN

# Build and run
./build.sh
./run.sh
```

### Environment Variables

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
DB_PATH=./memories.db
ENCRYPTION_KEY=your-32-character-key  # Optional
REVIEW_INTERVALS=1,3,7,14,30
```

---

## 📱 Bot Commands

| Command | Description |
|---------|-------------|
| `/start` | Welcome message with biological features overview |
| `/save <text>` | Save memory with emotion & context analysis |
| `/search <query>` | Search memories with smart ranking |
| `/recent` | View latest memories |
| `/stats` | Memory statistics with biological insights |
| `/help` | Command reference |

### Example Usage

```
/save Amazing breakthrough in my project! #work
```

**Bot Response:**
```
✅ Memory Saved Successfully!
━━━━━━━━━━━━━━━━━━━━━━━━

📊 Biological Analysis:

😊 Emotional Weight: 87% ████████▓▓
   Category: Intense 🤩

📍 Context: Monday Afternoon

🏷️ Tags: #work

🆔 Memory ID: 15

[📝 Save Another] [🔍 Search] [📊 Stats]
```

---

## 🏗️ Architecture

**Clean Architecture (4 Layers)**

```
Presentation → Application → Domain ← Infrastructure
   (UI)         (Use Cases)   (Core)   (External)
```

**Design Patterns:**
- Repository Pattern (data abstraction)
- Strategy Pattern (search algorithms)
- Command Pattern (bot commands)
- Observer Pattern (spaced repetition)
- Dependency Injection (loose coupling)

---

## 📁 Project Structure

```
├── cmd/bot/main.go                  # Entry point
├── internal/
│   ├── domain/                      # Core business logic
│   │   ├── entity/                  # Memory entity
│   │   ├── repository/              # Repository interface
│   │   └── service/                 # Sentiment, context services
│   ├── application/usecase/         # Use cases
│   ├── infrastructure/              # External services
│   │   ├── persistence/sqlite/      # Database
│   │   ├── messaging/telegram/      # Bot
│   │   ├── scheduler/               # Spaced repetition
│   │   └── search/strategy/         # Search algorithms
│   └── presentation/handler/        # Command handlers
├── pkg/
│   ├── config/                      # Configuration
│   └── encryption/                  # AES encryption
├── build.sh                         # Build script
├── run.sh                          # Run script
└── stop.sh                         # Stop script
```

---

## 🛠️ Development

### Build
```bash
go build -tags "fts5" -o memory-bot cmd/bot/main.go
# or
./build.sh
```

### Run
```bash
./memory-bot
# or
./run.sh
```

---

## 📊 Performance

- Search Speed: <100ms
- Memory Usage: ~15 MB
- Binary Size: ~14 MB
- Database: SQLite WAL mode

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file

---

## 👨‍💻 Author

**Milan Madusanka**
- GitHub: [@Milanz247](https://github.com/Milanz247)
- Repository: [Personal-memory-reminder-bot](https://github.com/Milanz247/Personal-memory-reminder-bot)

---

**Built with Clean Architecture & Neuroscience Principles**
