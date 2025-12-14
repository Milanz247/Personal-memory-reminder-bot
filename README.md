<div align="center">

# 🧠 Personal Memory Reminder Bot

**A sophisticated Telegram bot for intelligent memory management with AI-powered search**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Architecture](https://img.shields.io/badge/Architecture-Clean-blue?style=for-the-badge)](docs/ARCHITECTURE.md)
[![Telegram](https://img.shields.io/badge/Telegram-Bot-26A5E4?style=for-the-badge&logo=telegram)](https://telegram.org/)

[Features](#-features) •
[Quick Start](#-quick-start) •
[Documentation](#-documentation) •
[Architecture](#-architecture) •
[Security](#-security)

</div>

---

## 🌟 Overview

A production-ready Telegram bot that helps you store, search, and review your personal memories using neuroscience-inspired spaced repetition techniques. Built with **Clean Architecture** principles and **6 design patterns** for maximum maintainability and scalability.

### Why This Bot?

- 🔍 **Instant Recall**: Find any memory in milliseconds with AI-powered FTS5 search
- 🧠 **Spaced Repetition**: Scientifically-backed memory retention system
- 🔒 **Secure**: Optional AES-256-GCM encryption for sensitive data
- 🏗️ **Professional Codebase**: Clean Architecture with SOLID principles
- ⚡ **Blazing Fast**: Optimized SQLite with WAL mode and composite indexes
- 📦 **Zero Dependencies**: Single binary deployment

---

## ✨ Features

### Core Functionality

| Feature | Description |
|---------|-------------|
| 💾 **Smart Storage** | Automatic tag extraction and context-aware memory storage |
| 🔎 **Intelligent Search** | Multi-strategy search with fallback mechanisms |
| 📅 **Spaced Repetition** | Automated review reminders at optimal intervals (1, 3, 7, 14, 30 days) |
| 🏷️ **Tag Organization** | Hashtag-based categorization and filtering |
| 📊 **Statistics** | Track your memory collection growth |
| 📱 **Pagination** | Browse search results with inline keyboards |

### Technical Features

| Feature | Description |
|---------|-------------|
| 🔒 **Encryption** | Optional AES-256-GCM encryption with searchable index |
| 🔍 **FTS5 Search** | SQLite Full-Text Search with Okapi BM25 ranking |
| 🎯 **Design Patterns** | Repository, Strategy, Command, Observer, DI, Factory |
| 🏗️ **Clean Architecture** | 4-layer separation: Domain, Application, Infrastructure, Presentation |
| ⚡ **Performance** | WAL mode, composite indexes, <100ms search queries |
| 🧪 **Testable** | Interface-based design with easy mocking |

---

## 🚀 Quick Start

### Prerequisites

```bash
# Required
- Go 1.21 or higher
- SQLite3
- Telegram Bot Token (get from @BotFather)

# Optional
- OpenSSL (for encryption key generation)
```

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/Personal-memory-reminder-bot.git
cd Personal-memory-reminder-bot

# 2. Install dependencies
go mod download

# 3. Configure environment
cp .env.example .env
nano .env  # Add your TELEGRAM_BOT_TOKEN

# 4. Build and run
./build.sh
./memory-bot
```

### Quick Commands

```bash
./build.sh    # Build the binary
./run.sh      # Build and run
./stop.sh     # Stop all instances
```

---

## 📖 Usage

### Available Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/start` | Initialize bot and show welcome | `/start` |
| `/help` | Display help with action buttons | `/help` |
| `/save <text>` | Save a new memory | `/save Meeting with John tomorrow #work` |
| `/search <keyword>` | Search memories | `/search meeting` |
| `/recent [limit]` | View recent memories | `/recent 10` |
| `/stats` | Show statistics | `/stats` |

### Usage Examples

**Saving Memories:**
```
/save Remember to buy milk tomorrow #shopping

Project deadline is Friday #work #important

Sister's birthday is on March 15th #family #reminder
```

**Searching:**
```
/search meeting          # Find all meeting-related memories
/search #work           # Find all work memories
/search John project    # Multi-word search
```

**Smart Features:**
- Send any text without command → Bot asks if you want to save or search
- Automatic hashtag extraction and indexing
- Partial word matching (`tele` finds `telegram`)
- Multiple search strategies with auto-fallback

---

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────┐
│         Presentation Layer                      │
│  • Command Handlers                             │
│  • Telegram Bot Adapter                         │
│  • Input Validation                             │
└────────────────┬────────────────────────────────┘
                 │ depends on ↓
┌────────────────▼────────────────────────────────┐
│         Application Layer                        │
│  • Use Cases (Business Logic)                   │
│  • DTOs & Input/Output Models                   │
└────────────────┬────────────────────────────────┘
                 │ depends on ↓
┌────────────────▼────────────────────────────────┐
│         Domain Layer (Core)                      │
│  • Entities (Memory)                            │
│  • Repository Interfaces                        │
│  • Business Rules & Validation                  │
└────────────────▲────────────────────────────────┘
                 │ implements ↑
┌────────────────┴────────────────────────────────┐
│         Infrastructure Layer                     │
│  • SQLite Repository                            │
│  • Search Strategies                            │
│  • Telegram API                                 │
│  • Encryption Service                           │
│  • Spaced Repetition Scheduler                 │
└─────────────────────────────────────────────────┘
```

### Design Patterns

1. **Repository Pattern** - Data access abstraction
2. **Strategy Pattern** - Pluggable search algorithms
3. **Command Pattern** - Encapsulated bot commands
4. **Observer Pattern** - Event-driven notifications
5. **Dependency Injection** - Loose coupling
6. **Factory Pattern** - Object creation

📚 **[Read Full Architecture Guide →](docs/ARCHITECTURE.md)**  
📚 **[Explore Design Patterns →](docs/DESIGN_PATTERNS.md)**

---

## 🔐 Security

### Encryption Setup

The bot supports **optional AES-256-GCM encryption** for sensitive memories while maintaining full search functionality through a hybrid approach.

**Quick Setup:**

```bash
# 1. Generate encryption key
openssl rand -base64 32

# 2. Add to .env
echo "ENCRYPTION_KEY=your-generated-key-here" >> .env

# 3. Migrate existing database (if any)
./migrate_db.sh

# 4. Restart bot
./stop.sh && ./memory-bot
```

**How It Works:**
- `text_content`: Encrypted with AES-256-GCM (stored securely)
- `search_content`: Plain text (indexed by FTS5 for searching)
- Backward compatible with existing unencrypted data

### Security Best Practices

```bash
# Secure your environment file
chmod 600 .env

# Never commit .env to git (already in .gitignore)
# Keep regular backups of your database
cp memories.db memories.db.backup

# Use strong encryption keys (minimum 32 characters)
```

---

## 📊 Project Structure

```
Personal-memory-reminder-bot/
├── cmd/bot/                    # Application entry point
│   └── main.go                 # Dependency injection & startup
│
├── internal/                   # Private application code
│   ├── domain/                 # 🔵 Domain Layer
│   │   ├── entity/             # Business entities
│   │   └── repository/         # Repository interfaces
│   │
│   ├── application/            # 🟢 Application Layer
│   │   └── usecase/            # Business use cases
│   │
│   ├── infrastructure/         # 🟡 Infrastructure Layer
│   │   ├── persistence/sqlite/ # SQLite implementation
│   │   ├── search/strategy/    # Search strategies
│   │   ├── messaging/telegram/ # Telegram bot adapter
│   │   └── scheduler/          # Spaced repetition
│   │
│   └── presentation/           # 🔴 Presentation Layer
│       └── handler/command/    # Command handlers
│
├── pkg/                        # Public reusable packages
│   ├── config/                 # Configuration
│   └── encryption/             # AES-256 encryption
│
├── docs/                       # Documentation
│   ├── ARCHITECTURE.md
│   └── DESIGN_PATTERNS.md
│
├── build.sh                    # Build script
├── run.sh                      # Run script
├── stop.sh                     # Stop script
├── migrate_db.sh               # Database migration
└── README.md                   # This file
```

---

## 🛠️ Development

### Building

```bash
# Simple build
go build -tags "fts5" -o memory-bot cmd/bot/main.go

# Or use build script
./build.sh
```

### Running

```bash
# Direct run
go run -tags "fts5" cmd/bot/main.go

# Or use run script (recommended)
./run.sh
```

### Adding Features

Thanks to Clean Architecture, extending functionality is straightforward:

**Example: Adding a `/delete` command**

```go
// 1. Create use case (application/usecase/delete_memory.go)
type DeleteMemoryUseCase struct {
    repo repository.MemoryRepository
}

// 2. Create command handler (presentation/handler/command/delete_command.go)
type DeleteCommand struct {
    useCase *usecase.DeleteMemoryUseCase
}

// 3. Register in main.go
registry.Register(command.NewDeleteCommand(deleteUC))
```

**No changes to existing code!** ✨

---

## 📚 Documentation

| Document | Description |
|----------|-------------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Complete architecture guide with diagrams |
| [DESIGN_PATTERNS.md](docs/DESIGN_PATTERNS.md) | Design patterns explained with examples |
| [README.md](README.md) | This file - project overview |

---

## 🧪 Testing

### Manual Testing

```bash
# Build and run
./build.sh
./memory-bot

# Test commands
# In Telegram:
/start
/save Test memory #test
/search test
/recent
/stats
```

### Unit Tests (Coming Soon)

```bash
go test ./...
```

---

## 📈 Performance

| Metric | Value |
|--------|-------|
| **Search Speed** | <100ms (typical) |
| **Binary Size** | ~13.7 MB |
| **Memory Usage** | ~15 MB |
| **Startup Time** | ~350ms |
| **Database** | SQLite WAL mode |

**Optimizations:**
- Composite indexes on (user_id, created_at DESC)
- FTS5 with Porter stemming
- Prepared statements
- Connection pooling

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **Inspired by** neuroscience research on memory formation and spaced repetition
- **Built with** [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
- **Powered by** SQLite FTS5 full-text search
- **Architecture** based on Clean Architecture by Robert C. Martin

---

## 📞 Support

- 📧 **Creator**: Milan Madusanka
- 🐛 **GitHub**: [https://github.com/Milanz247](https://github.com/Milanz247)
- 💬 **Issues**: [Report Issues](https://github.com/Milanz247/Personal-memory-reminder-bot/issues)

---

<div align="center">

**Made with ❤️ and 🧠**

**[⬆ back to top](#-personal-memory-reminder-bot)**

</div>
