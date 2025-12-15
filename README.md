<div align="center">

# 🧠 Personal Memory Reminder Bot

**Telegram Bot for Intelligent Memory Management with Biological Features**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)
[![Telegram](https://img.shields.io/badge/Platform-Telegram-blue?style=flat-square&logo=telegram)](https://telegram.org/)

**A neuroscience-inspired personal memory assistant powered by biological memory principles, advanced full-text search, spaced repetition algorithms, and military-grade encryption.**

</div>

---

## 📋 Overview

### What is This Bot?

This is a **smart personal memory assistant** that works through Telegram. It helps you:
- 📝 Save important memories, notes, and information
- 🔍 Search and find your memories quickly
- 🧠 Remember things better using neuroscience-based techniques
- 📊 Track your memory patterns and insights
- 🔒 Keep your data secure with encryption

### 🌟 Key Features

#### 🧠 Biological Memory System (Inspired by Brain Science)
- **Amygdala Emotional Tagging** - Automatically detects emotional content (0-100%)
- **Hippocampus Context Encoding** - Records time, day, and location automatically
- **Sleep Consolidation** - Boosts important memories during "sleep" periods
- **LTP Spaced Repetition** - Reviews memories at optimal intervals (1, 3, 7, 14, 30 days)
- **Forgetting Curve Prevention** - Reminds you before you forget

#### 🔍 Advanced Search (8 Search Strategies)
1. **Hashtag Search** - Find by tags: `#work`, `#health`, `#family`
2. **Contextual Search** - Search by time: "yesterday", "this morning"
3. **FTS5 Search** - Powerful full-text search with wildcards
4. **Fuzzy Search** - Finds similar words (handles typos)
5. **AND Search** - All words must match
6. **Partial Match** - Matches part of words
7. **OR Search** - Any word matches
8. **NEAR Search** - Finds words close to each other

#### 🎯 Smart Ranking
- **Emotional memories** ranked 2× higher
- **Recently consolidated** memories boosted
- **Recent memories** get recency advantage
- **BM25 algorithm** for relevance scoring

#### 🔒 Security Features
- **AES-256-GCM encryption** for sensitive data
- **Hybrid architecture** - encrypted storage + searchable index
- **Optional encryption** - you choose when to enable it

---

### 🚀 Installation Guide (Step by Step)

#### Step 1: Install Prerequisites

**For Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go sqlite3 git
```

**For Fedora/RHEL:**
```bash
sudo dnf install golang sqlite git
```

**For macOS:**
```bash
brew install go sqlite3 git
```

**For Windows:**
- Download and install Go from: https://golang.org/dl/
- Download and install Git from: https://git-scm.com/
- SQLite is included with Go

#### Step 2: Create Telegram Bot

1. Open Telegram and search for `@BotFather`
2. Send `/newbot` command
3. Follow the prompts to create your bot
4. Copy the **Bot Token** (looks like: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`)
5. Save this token - you'll need it in Step 4

#### Step 3: Clone and Setup

```bash
# Clone the repository
git clone https://github.com/Milanz247/Personal-memory-reminder-bot.git

# Navigate to project directory
cd Personal-memory-reminder-bot

# Copy example environment file
cp .env.example .env
```

#### Step 4: Configure Environment

Open `.env` file with any text editor:

```bash
nano .env
# or
vim .env
# or use any text editor
```

Update the following:

```env
# REQUIRED: Paste your bot token from Step 2
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Database file location (default is fine)
DB_PATH=./memories.db

# OPTIONAL: Generate encryption key for security
# Run: openssl rand -base64 32
ENCRYPTION_KEY=your-32-character-key

# OPTIONAL: Review intervals in days (default is fine)
REVIEW_INTERVAL_1=1
REVIEW_INTERVAL_2=3
REVIEW_INTERVAL_3=7
REVIEW_INTERVAL_4=14
REVIEW_INTERVAL_5=30
```

**To generate a secure encryption key:**
```bash
openssl rand -base64 32
```
Copy the output and paste it as your `ENCRYPTION_KEY`.

#### Step 5: Build the Bot

```bash
# Make scripts executable
chmod +x build.sh run.sh stop.sh

# Build the bot
./build.sh
```

You should see:
```
🔨 Building Biological Memory Bot...
✅ Build successful!
Binary: memory-bot
Size: 14M
```

#### Step 6: Run the Bot

```bash
./run.sh
```

You should see:
```
🤖 Starting Biological Memory Bot...
✅ Environment validated
📊 Active Biological Features:
   • Amygdala Emotional Tagging
   • Hippocampus Context Encoding
   • Sleep Consolidation
   • LTP Spaced Repetition

2025/12/15 14:00:00 Starting Memory Storage Bot...
2025/12/15 14:00:00 ✅ Bot is running!
```

#### Step 7: Start Using in Telegram

1. Open Telegram
2. Search for your bot (the username you created in Step 2)
3. Click "Start" or send `/start`
4. Start saving memories!

---

### 📱 How to Use the Bot

#### Basic Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/start` | Show welcome message | `/start` |
| `/save` | Save a memory | `/save Meeting tomorrow 3 PM #work` |
| `/search` | Search memories | `/search meeting` |
| `/recent` | Show recent memories | `/recent` |
| `/stats` | View statistics | `/stats` |
| `/help` | Get help | `/help` |

#### Saving Memories

**Simple save:**
```
/save Meeting with client tomorrow at 3 PM
```

**With hashtags (for organization):**
```
/save Completed project milestone! Very happy! #work #achievement
```

**With emotional content (gets higher ranking):**
```
/save Amazing breakthrough today! Solved the bug! #coding
```

**Interactive save (shows template):**
```
/save
```
Then follow the template shown.

#### Searching Memories

**Simple search:**
```
/search meeting
```

**Tag search:**
```
/search #work
```

**Contextual search:**
```
/search yesterday
/search this morning
/search last Monday
```

**Multi-word search:**
```
/search project meeting client
```

#### Understanding Search Results

Results are ranked by:
- **Relevance** (how well it matches)
- **Emotional weight** (emotional memories ranked higher)
- **Recency** (recent memories boosted)
- **Priority** (consolidated memories)

---

### 🛑 Managing the Bot

#### Stop the Bot
```bash
./stop.sh
```

#### Restart the Bot
```bash
./stop.sh
./run.sh
```

#### Check if Bot is Running
```bash
ps aux | grep memory-bot
```

#### View Logs
The bot outputs logs to the terminal. To save logs:
```bash
./run.sh > bot.log 2>&1 &
tail -f bot.log
```

---

### 🔧 Troubleshooting

#### Bot doesn't start
- Check if `.env` file exists and has correct `TELEGRAM_BOT_TOKEN`
- Make sure port is not already in use
- Check Go version: `go version` (need 1.21+)

#### "No memories found" when searching
- Make sure you've saved some memories first using `/save`
- Try simpler search terms
- Check if database file `memories.db` exists

#### Build fails
```bash
# Make sure you have FTS5 support
go build -tags "fts5" -o memory-bot cmd/bot/main.go
```

#### Permission denied
```bash
chmod +x build.sh run.sh stop.sh memory-bot
```

---

### 💡 Tips for Best Results

#### 1. Use Emotional Words
Words like "amazing", "terrible", "excited", "worried" increase emotional weight:
```
/save I'm so excited about the new project! #work
```

#### 2. Add Context
Include time, place, or people:
```
/save Met Sarah at coffee shop, discussed project timeline #meeting
```

#### 3. Use Hashtags
Organize memories with tags:
```
/save Completed module 3 #project #milestone #coding
```

#### 4. Be Specific
More specific = better recall:
```
❌ /save meeting
✅ /save Quarterly review meeting with team about Q4 goals #work
```

#### 5. Regular Reviews
Let the bot remind you - don't skip review notifications!

---

### 📊 Understanding Statistics

When you use `/stats`, you'll see:
- **Total Memories** - How many memories you've saved
- **Average Emotional Weight** - How emotional your memories are
- **Most Used Tags** - Your common categories
- **Review Completion Rate** - How well you're maintaining memories
- **Biological Features Status** - What's active in your system

---

### 🔐 Security & Privacy

#### Encryption
- Your memories can be encrypted with AES-256-GCM
- Only you have the encryption key
- Without the key, memories cannot be decrypted

#### Data Storage
- All data stored locally in `memories.db` file
- No cloud storage
- You control your data

#### Backup Your Data
```bash
# Backup database
cp memories.db memories-backup-$(date +%Y%m%d).db

# Backup encryption key
cp .env .env.backup
```

⚠️ **Important:** Never lose your `.env` file if you're using encryption!

---

## 📊 System Metrics

- **Search Speed:** <100ms
- **Memory Usage:** ~15 MB
- **Binary Size:** ~14 MB
- **Database:** SQLite WAL mode

---

## 📁 Project Structure

```
├── cmd/bot/main.go              # Entry point
├── internal/
│   ├── domain/                  # Core business logic
│   ├── application/             # Use cases
│   ├── infrastructure/          # External services
│   └── presentation/            # Command handlers
├── pkg/
│   ├── config/                  # Configuration
│   └── encryption/              # AES encryption
├── build.sh                     # Build script
├── run.sh                       # Run script
└── stop.sh                      # Stop script
```

---

## 🤝 Contributing

Contributions are welcome! Feel free to submit pull requests.

---

## 📄 License

MIT License

---

## 👨‍💻 Author

**Milan Madusanka**
- GitHub: [@Milanz247](https://github.com/Milanz247)
- Repository: [Personal-memory-reminder-bot](https://github.com/Milanz247/Personal-memory-reminder-bot)

---

<div align="center">

**Built with Clean Architecture & Neuroscience Principles**

[⬆ Back to Top](#-personal-memory-reminder-bot)

</div>
