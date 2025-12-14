# Changelog

## [2.0.0] - 2025-12-14

### 🏗️ Architecture Refactoring
- **Complete Clean Architecture Implementation**
  - Separated into 4 layers: Domain, Application, Infrastructure, Presentation
  - 24 new Go files implementing clean separation of concerns
  - Removed old monolithic structure (1,479 lines of legacy code)

### 🎯 Design Patterns Implemented
1. **Repository Pattern** - Abstract data access from business logic
2. **Strategy Pattern** - Interchangeable search algorithms (Smart Search)
3. **Command Pattern** - Bot command encapsulation
4. **Observer Pattern** - Event-driven spaced repetition
5. **Dependency Injection** - Clean dependency management
6. **Factory Pattern** - Centralized object creation

### 🔒 Security Features
- **AES-256-GCM Encryption** support for sensitive memory data
- **Searchable Encryption** - Hybrid approach maintaining full search while encrypting stored data
- **Backward Compatible Decryption** - Handles both encrypted and plain text data
- Optional encryption via environment variable

### 🔍 Search Improvements
- **Fixed FTS5 Search Syntax** - Corrected OR query handling
- **Smart Search Strategy** with multiple fallbacks:
  - Primary: FTS5 with wildcards
  - Fallback 1: AND search
  - Fallback 2: OR search (without wildcards)
- **Searchable Encryption** - Separate search index maintains full-text search on encrypted data

### 📚 Documentation
- **ARCHITECTURE.md** - Complete clean architecture guide (500+ lines)
- **DESIGN_PATTERNS.md** - Detailed pattern explanations (400+ lines)
- **README.md** - Updated with new structure and encryption setup
- **TROUBLESHOOTING.md** - Common issues and solutions

### 🛠️ Developer Tools
- `stop.sh` - Stop all running bot instances
- `build.sh` - Build with automatic instance cleanup
- `run.sh` - Build and run with automatic cleanup
- `migrate_db.sh` - Database migration for searchable encryption

### 🔧 Technical Improvements
- Better error handling throughout the codebase
- Proper context handling in all operations
- Type-safe operations with clear interfaces
- Comprehensive logging

### 🐛 Bug Fixes
- Fixed multiple bot instance conflicts
- Fixed FTS5 search syntax errors with OR operator
- Fixed encryption/decryption type mismatches
- Fixed search compatibility with encrypted data

### ⚠️ Breaking Changes
- New entry point: `cmd/bot/main.go` (old `main.go` removed)
- New folder structure with `internal/` and `pkg/` directories
- Database schema updated with `search_content` column

### 📦 Dependencies
- No new external dependencies
- Same core stack: Go 1.21+, SQLite FTS5, telegram-bot-api/v5

### 🔄 Migration Notes
- Run `./migrate_db.sh` to update existing database for searchable encryption
- Old database file remains compatible
- Backup created automatically during migration

---

## [1.0.0] - Initial Release
- Basic memory storage and retrieval
- FTS5 full-text search
- Spaced repetition system
- Tag organization
- Simple monolithic structure
