# 🧠 Personal Memory Reminder Bot

A Telegram bot that helps you store, search, and review your personal memories using spaced repetition techniques.

## 🔒 Security Features

- **AES-256-GCM Encryption**: All memory content is encrypted before storage
- **Optional Encryption**: Works with or without encryption enabled
- **Secure Key Management**: Uses environment variables for key storage

[📖 Read the Encryption Setup Guide](docs/ENCRYPTION_SETUP.md)

## ✨ Features

- 💾 **Save Memories**: Store your thoughts, ideas, and important information
- 🔍 **Smart Search**: Find memories using intelligent full-text search
- 📅 **Spaced Repetition**: Automatic review reminders based on proven learning intervals
- 🏷️ **Auto-tagging**: Automatically extracts and indexes hashtags
- 📊 **Statistics**: Track your memory collection growth
- 🔒 **Encryption**: Optional AES-256 encryption for sensitive data

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Telegram Bot Token ([Get one from @BotFather](https://t.me/botfather))
- SQLite3

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Personal-memory-reminder-bot
```

2. Copy and configure environment variables:
```bash
cp .env.example .env
```

3. Edit `.env` and add your credentials:
```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
DB_PATH=./memories.db
ENCRYPTION_KEY=your-secret-encryption-key  # Optional but recommended
```

4. Build and run:
```bash
./build.sh
./run.sh
```

## 🎯 Usage

### Available Commands

- `/start` - Initialize the bot and get welcome message
- `/help` - Show available commands
- `/save <text>` - Save a new memory
- `/search <keyword>` - Search for memories
- `/recent [limit]` - View recent memories (default: 5)
- `/stats` - View memory statistics

### Examples

```
/save Remember to call mom on Sunday #family
/search family
/recent 10
/stats
```

## 🔐 Encryption Setup

To enable encryption for your memories:

1. Generate a strong encryption key:
```bash
openssl rand -base64 32
```

2. Add it to your `.env` file:
```env
ENCRYPTION_KEY=your-generated-key-here
```

3. Restart the bot:
```bash
./stop.sh
./run.sh
```

For detailed encryption setup instructions, see [ENCRYPTION_SETUP.md](docs/ENCRYPTION_SETUP.md).

## 📚 Documentation

- [Architecture Overview](docs/ARCHITECTURE.md)
- [Design Patterns](docs/DESIGN_PATTERNS.md)
- [Encryption Setup Guide](docs/ENCRYPTION_SETUP.md)

## 🏗️ Architecture

This project follows Clean Architecture principles with clear separation of concerns:

```
cmd/bot/          - Application entry point
internal/
  ├── domain/     - Business entities and interfaces
  ├── application/ - Use cases (business logic)
  ├── infrastructure/ - External implementations
  └── presentation/   - Command handlers
pkg/
  ├── config/     - Configuration management
  └── encryption/ - AES-256 encryption utilities
```

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Database**: SQLite with FTS5 (Full-Text Search)
- **Bot Framework**: telegram-bot-api/v5
- **Architecture**: Clean Architecture
- **Design Patterns**: Repository, Strategy, Command, Observer
- **Encryption**: AES-256-GCM

## 🧪 Development

### Build
```bash
make build
# or
./build.sh
```

### Run
```bash
make run
# or
./run.sh
```

### Stop
```bash
./stop.sh
```

## 🔒 Security Best Practices

1. **Never commit** `.env` file (already in `.gitignore`)
2. **Use strong encryption keys** (minimum 32 characters)
3. **Keep your encryption key safe** - losing it means losing access to encrypted memories
4. **Secure your `.env` file**: `chmod 600 .env`
5. **Regularly backup** your `memories.db` file

## 📝 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📞 Support

For issues or questions, please open an issue on GitHub.

---

**Made with ❤️ for better memory management**
