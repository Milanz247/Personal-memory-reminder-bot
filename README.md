<div align="center">

# 🧠 Personal Memory Reminder Bot

**Telegram Bot for Intelligent Memory Management with Biological Features**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)
[![Telegram](https://img.shields.io/badge/Platform-Telegram-blue?style=flat-square&logo=telegram)](https://telegram.org/)

**A neuroscience-inspired personal memory assistant powered by biological memory principles, advanced full-text search, spaced repetition algorithms, and military-grade encryption.**

[English](#english) • [සිංහල](#sinhala)

</div>

---

<a name="english"></a>

## 📋 English Guide

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

<a name="sinhala"></a>

## 🇱🇰 සිංහල මාර්ගෝපදේශය

### මෙය කුමක්ද?

මෙය Telegram හරහා ක්‍රියාත්මක වන **බුද්ධිමත් පුද්ගලික මතක සහායකයෙකි**. මෙය ඔබට උදව් කරයි:
- 📝 වැදගත් මතකයන්, සටහන් සුරකින්න
- 🔍 ඔබේ මතකයන් ඉක්මනින් සොයා ගන්න
- 🧠 ස්නායු විද්‍යාව මත පදනම් වූ ක්‍රම භාවිතයෙන් දේවල් මතක තබා ගන්න
- 📊 ඔබේ මතක රටා නිරීක්ෂණය කරන්න
- 🔒 ඔබේ දත්ත සංකේතනය කර ආරක්ෂිතව තබා ගන්න

### 🌟 ප්‍රධාන ලක්ෂණ

#### 🧠 ජීව විද්‍යාත්මක මතක පද්ධතිය
- **Amygdala චිත්තවේගීය ටැග් කිරීම** - චිත්තවේගීය බර ස්වයංක්‍රීයව හඳුනා ගනී (0-100%)
- **Hippocampus සන්දර්භ සංකේතනය** - වේලාව, දිනය ස්වයංක්‍රීයව සටහන් කරයි
- **නිදි සම්පීඩනය** - වැදගත් මතකයන් "නිදි" කාලය තුළ ශක්තිමත් කරයි
- **LTP පරතරය සමඟ නැවත සමාලෝචනය** - ප්‍රශස්ත කාල පරතරවලදී සමාලෝචනය කරයි (1, 3, 7, 14, 30 දින)
- **අමතක වීමේ වක්‍රය වැළැක්වීම** - අමතක වීමට පෙර මතක් කරයි

#### 🔍 උසස් සෙවීම (සෙවුම් උපාය 8ක්)
1. **Hashtag සෙවීම** - ටැග් මගින් සොයන්න: `#work`, `#health`
2. **සන්දර්භීය සෙවීම** - කාලය අනුව සොයන්න: "ඊයේ", "අද උදෑසන"
3. **FTS5 සෙවීම** - wildcards සහිත ශක්තිමත් සෙවීම
4. **Fuzzy සෙවීම** - සමාන වචන සොයා ගනී (අකුරු වැරදි හසුරුවයි)
5. **AND සෙවීම** - සියලු වචන ගැලපිය යුතුය
6. **අර්ධ ගැලපුම** - වචනවල කොටසක් ගැලපේ
7. **OR සෙවීම** - ඕනෑම වචනයක් ගැලපේ
8. **NEAR සෙවීම** - එකිනෙකට ආසන්න වචන සොයා ගනී

---

### 🚀 ස්ථාපන මාර්ගෝපදේශය (පියවරෙන් පියවර)

#### පියවර 1: අවශ්‍ය මෘදුකාංග ස්ථාපනය කරන්න

**Ubuntu/Debian සඳහා:**
```bash
sudo apt update
sudo apt install golang-go sqlite3 git
```

**Fedora/RHEL සඳහා:**
```bash
sudo dnf install golang sqlite git
```

#### පියවර 2: Telegram Bot එකක් සාදන්න

1. Telegram විවෘත කර `@BotFather` සොයන්න
2. `/newbot` විධානය යවන්න
3. ඔබේ bot එක සෑදීමට උපදෙස් අනුගමනය කරන්න
4. **Bot Token** එක පිටපත් කරන්න (උදා: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`)
5. මෙම token එක සුරකින්න - පියවර 4 හිදී අවශ්‍ය වේ

#### පියවර 3: ව්‍යාපෘතිය Clone කර Setup කරන්න

```bash
# Repository එක clone කරන්න
git clone https://github.com/Milanz247/Personal-memory-reminder-bot.git

# ව්‍යාපෘති ෆෝල්ඩරයට යන්න
cd Personal-memory-reminder-bot

# උදාහරණ පරිසර ගොනුව පිටපත් කරන්න
cp .env.example .env
```

#### පියවර 4: පරිසර විචල්‍යයන් සකසන්න

`.env` ගොනුව text editor එකකින් විවෘත කරන්න:

```bash
nano .env
```

පහත දෑ යාවත්කාලීන කරන්න:

```env
# අවශ්‍යයි: පියවර 2 සිට ඔබේ bot token එක paste කරන්න
TELEGRAM_BOT_TOKEN=your_bot_token_here

# Database ගොනු ස්ථානය (පෙරනිමි හොඳයි)
DB_PATH=./memories.db

# විකල්පය: ආරක්ෂාව සඳහා encryption key එකක් ජනනය කරන්න
ENCRYPTION_KEY=your-32-character-key

# විකල්පය: සමාලෝචන කාල පරතර දින වලින් (පෙරනිමි හොඳයි)
REVIEW_INTERVAL_1=1
REVIEW_INTERVAL_2=3
REVIEW_INTERVAL_3=7
REVIEW_INTERVAL_4=14
REVIEW_INTERVAL_5=30
```

**ආරක්ෂිත encryption key එකක් ජනනය කිරීමට:**
```bash
openssl rand -base64 32
```

#### පියවර 5: Bot එක Build කරන්න

```bash
# Scripts executable කරන්න
chmod +x build.sh run.sh stop.sh

# Bot එක build කරන්න
./build.sh
```

ඔබට පෙනේවි:
```
🔨 Building Biological Memory Bot...
✅ Build successful!
Binary: memory-bot
Size: 14M
```

#### පියවර 6: Bot එක ධාවනය කරන්න

```bash
./run.sh
```

ඔබට පෙනේවි:
```
🤖 Starting Biological Memory Bot...
✅ Environment validated
📊 Active Biological Features:
   • Amygdala Emotional Tagging
   • Hippocampus Context Encoding
   • Sleep Consolidation
   • LTP Spaced Repetition

Bot is running!
```

#### පියවර 7: Telegram හි භාවිතා කිරීම ආරම්භ කරන්න

1. Telegram විවෘත කරන්න
2. ඔබේ bot එක සොයන්න (පියවර 2 හි ඔබ සෑදූ username)
3. "Start" ක්ලික් කරන්න හෝ `/start` යවන්න
4. මතකයන් සුරැකීම ආරම්භ කරන්න!

---

### 📱 Bot එක භාවිතා කරන්නේ කෙසේද

#### මූලික විධාන

| විධානය | විස්තරය | උදාහරණය |
|---------|-------------|---------|
| `/start` | පිළිගැනීමේ පණිවිඩය පෙන්වන්න | `/start` |
| `/save` | මතකයක් සුරකින්න | `/save හෙට රැස්වීම 3 ට #work` |
| `/search` | මතකයන් සොයන්න | `/search රැස්වීම` |
| `/recent` | මෑත මතකයන් පෙන්වන්න | `/recent` |
| `/stats` | සංඛ්‍යාන බලන්න | `/stats` |
| `/help` | උදව් ලබා ගන්න | `/help` |

#### මතකයන් සුරකින්න

**සරල save:**
```
/save හෙට ග්‍රාහකයා සමඟ රැස්වීම පස්වරු 3 ට
```

**Hashtags සමඟ (සංවිධානය සඳහා):**
```
/save ව්‍යාපෘති සන්ධිස්ථානය සම්පූර්ණයි! ඉතා සතුටුයි! #work #achievement
```

**චිත්තවේගීය අන්තර්ගතය සමඟ (ඉහළ ශ්‍රේණිගත කිරීමක් ලැබේ):**
```
/save අද විශිෂ්ට පෙරළියක්! Bug එක විසඳා ගත්තා! #coding
```

#### මතකයන් සොයන්න

**සරල සෙවීම:**
```
/search රැස්වීම
```

**Tag සෙවීම:**
```
/search #work
```

**සන්දර්භීය සෙවීම:**
```
/search ඊයේ
/search අද උදෑසන
/search පසුගිය සඳුදා
```

---

### 🛑 Bot කළමනාකරණය

#### Bot එක නවත්වන්න
```bash
./stop.sh
```

#### Bot එක නැවත ආරම්භ කරන්න
```bash
./stop.sh
./run.sh
```

#### Bot එක ධාවනය වේදැයි පරීක්ෂා කරන්න
```bash
ps aux | grep memory-bot
```

---

### 🔧 ගැටළු නිරාකරණය

#### Bot එක start නොවේ නම්
- `.env` ගොනුව ඇත්ද සහ නිවැරදි `TELEGRAM_BOT_TOKEN` ඇතිද පරීක්ෂා කරන්න
- Go version පරීක්ෂා කරන්න: `go version` (1.21+ අවශ්‍යයි)

#### සෙවීමේදී "මතකයන් හමු නොවුනි" යැයි පෙන්වන්නේ නම්
- පළමුව `/save` භාවිතයෙන් මතකයන් සුරකින්න
- සරල සෙවුම් පද භාවිතා කරන්න

#### Build fail වේ නම්
```bash
# FTS5 support ඇතිද පරීක්ෂා කරන්න
go build -tags "fts5" -o memory-bot cmd/bot/main.go
```

#### අවසරය ප්‍රතික්ෂේප වේ නම්
```bash
chmod +x build.sh run.sh stop.sh memory-bot
```

---

### 💡 හොඳම ප්‍රතිඵල සඳහා උපදෙස්

#### 1. චිත්තවේගීය වචන භාවිතා කරන්න
"විශිෂ්ට", "භයානක", "උද්යෝගිමත්", "කනස්සල්ල" වැනි වචන චිත්තවේගීය බර වැඩි කරයි:
```
/save නව ව්‍යාපෘතිය ගැන මම ඉතා උද්යෝගිමත්! #work
```

#### 2. සන්දර්භය එක් කරන්න
වේලාව, ස්ථානය, හෝ පුද්ගලයන් ඇතුළත් කරන්න:
```
/save කෝපි කඩයේදී සාරා මුණගැසී ව්‍යාපෘති කාල සටහන සාකච්ඡා කළා #meeting
```

#### 3. Hashtags භාවිතා කරන්න
මතකයන් tags සමඟ සංවිධානය කරන්න:
```
/save මොඩියුල 3 සම්පූර්ණයි #project #milestone #coding
```

#### 4. නිශ්චිත වන්න
වැඩි නිශ්චිත = වඩා හොඳ මතක කැඳවීමක්:
```
❌ /save රැස්වීම
✅ /save Q4 අරමුණු පිළිබඳ කණ්ඩායම සමඟ කාර්තුමය සමාලෝචන රැස්වීම #work
```

---

### 🔐 ආරක්ෂාව සහ පෞද්ගලිකත්වය

#### සංකේතනය (Encryption)
- ඔබේ මතකයන් AES-256-GCM සමඟ සංකේතනය කළ හැක
- ඔබ පමණක් සතුව encryption key ඇත
- Key නොමැතිව මතකයන් විකේතනය කළ නොහැක

#### දත්ත ගබඩාව
- සියලුම දත්ත `memories.db` ගොනුවේ දේශීයව ගබඩා වේ
- Cloud storage නැත
- ඔබ ඔබේ දත්ත පාලනය කරයි

#### ඔබේ දත්ත උපස්ථය කරන්න
```bash
# Database උපස්ථය
cp memories.db memories-backup-$(date +%Y%m%d).db

# Encryption key උපස්ථය
cp .env .env.backup
```

⚠️ **වැදගත්:** Encryption භාවිතා කරන්නේ නම් `.env` ගොනුව කිසිවිටෙක අහිමි නොකරන්න!

---

## 📊 පද්ධති ප්‍රමිතීන්

- **සෙවුම් වේගය:** <100ms
- **මතක භාවිතය:** ~15 MB
- **Binary ප්‍රමාණය:** ~14 MB
- **Database:** SQLite WAL ප්‍රකාරය

---

## 📁 ව්‍යාපෘති ව්‍යුහය

```
├── cmd/bot/main.go              # ප්‍රධාන ආරම්භක ස්ථානය
├── internal/
│   ├── domain/                  # ප්‍රධාන business logic
│   ├── application/             # Use cases
│   ├── infrastructure/          # බාහිර සේවා
│   └── presentation/            # විධාන handlers
├── pkg/
│   ├── config/                  # වින්‍යාසය
│   └── encryption/              # AES සංකේතනය
├── build.sh                     # Build script
├── run.sh                       # Run script
└── stop.sh                      # Stop script
```

---

## 🤝 දායකත්වය (Contributing)

දායකත්වය පිළිගනිමු! Pull requests යවන්න.

---

## 📄 බලපත්‍රය (License)

MIT License

---

## 👨‍💻 කතුවරයා (Author)

**Milan Madusanka**
- GitHub: [@Milanz247](https://github.com/Milanz247)
- Repository: [Personal-memory-reminder-bot](https://github.com/Milanz247/Personal-memory-reminder-bot)

---

<div align="center">

**Clean Architecture සහ ස්නායු විද්‍යා මූලධර්ම සමඟ සාදා ඇත**

**Built with Clean Architecture & Neuroscience Principles**

[⬆ Back to Top](#-personal-memory-reminder-bot)

</div>
