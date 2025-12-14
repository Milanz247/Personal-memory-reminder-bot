# Personal Memory Reminder Bot - Architecture Documentation

## Overview

This document describes the clean architecture implementation of the Personal Memory Reminder Bot, a Telegram bot for storing and retrieving memories using AI-powered full-text search.

## Architecture Principles

The project follows **Clean Architecture** (also known as Hexagonal Architecture or Ports and Adapters), which provides:

- **Separation of Concerns**: Each layer has a specific responsibility
- **Dependency Rule**: Dependencies point inward (outer layers depend on inner layers)
- **Testability**: Easy to test each layer independently
- **Maintainability**: Changes in one layer don't affect others
- **Scalability**: Easy to add new features or swap implementations

---

## Architecture Layers

```
┌─────────────────────────────────────────────────┐
│         Presentation Layer                      │
│  (Command Handlers, Telegram Adapter)          │
└────────────────┬────────────────────────────────┘
                 │ depends on
┌────────────────▼────────────────────────────────┐
│         Application Layer                        │
│  (Use Cases, Business Logic)                    │
└────────────────┬────────────────────────────────┘
                 │ depends on
┌────────────────▼────────────────────────────────┐
│         Domain Layer                             │
│  (Entities, Repository Interfaces, Rules)       │
└────────────────▲────────────────────────────────┘
                 │ implements
┌────────────────┴────────────────────────────────┐
│         Infrastructure Layer                     │
│  (Database, External APIs, Frameworks)          │
└─────────────────────────────────────────────────┘
```

### 1. Domain Layer (Core Business Logic)

**Location**: `internal/domain/`

**Purpose**: Contains the core business entities and rules. This is the heart of the application and is completely independent of external frameworks.

**Components**:
- **Entities** (`entity/`): Core business objects (Memory)
- **Repository Interfaces** (`repository/`): Data access contracts
- **Business Rules**: Validation, domain logic

**Files**:
- `entity/memory.go` - Memory entity with business methods
- `entity/errors.go` - Domain-level errors
- `repository/memory_repository.go` - Repository interface

**Key Principles**:
- No dependencies on external libraries
- Contains only business logic
- Defines interfaces that outer layers must implement

---

### 2. Application Layer (Use Cases)

**Location**: `internal/application/`

**Purpose**: Contains application-specific business rules and orchestrates the flow of data between layers.

**Components**:
- **Use Cases** (`usecase/`): Application-specific business logic
- **DTOs** (`dto/`): Data transfer objects (if needed)

**Files**:
- `usecase/save_memory.go` - Save memory use case
- `usecase/search_memory.go` - Search memory use case
- `usecase/get_recent_memories.go` - Get recent memories
- `usecase/get_stats.go` - Get statistics
- `usecase/review_memory.go` - Review memory (spaced repetition)

**Key Principles**:
- Depends only on domain layer
- No knowledge of UI or database implementations
- Orchestrates domain entities and repository interfaces

---

### 3. Infrastructure Layer (External Concerns)

**Location**: `internal/infrastructure/`

**Purpose**: Implements the interfaces defined by inner layers. Deals with external concerns like databases, APIs, and frameworks.

**Components**:
- **Persistence** (`persistence/sqlite/`): Database implementation
- **Search** (`search/strategy/`): Search algorithm implementations
- **Messaging** (`messaging/telegram/`): Telegram bot adapter
- **Scheduler** (`scheduler/`): Background job scheduler

**Files**:
```
infrastructure/
├── persistence/
│   └── sqlite/
│       ├── connection.go          # Database connection
│       └── memory_repository.go   # Repository implementation
├── search/
│   └── strategy/
│       ├── search_strategy.go     # Strategy interface
│       └── smart_strategy.go      # Smart search implementation
├── messaging/
│   └── telegram/
│       └── bot.go                 # Telegram bot adapter
└── scheduler/
    └── spaced_repetition.go       # Spaced repetition scheduler
```

**Key Principles**:
- Implements domain interfaces
- Can be easily swapped (e.g., SQLite → PostgreSQL)
- Contains framework-specific code

---

### 4. Presentation Layer (User Interface)

**Location**: `internal/presentation/`

**Purpose**: Handles user interaction and translates user input into use case calls.

**Components**:
- **Command Handlers** (`handler/command/`): Bot command implementations
- **Validators** (`validator/`): Input validation

**Files**:
```
presentation/
└── handler/
    └── command/
        ├── command.go         # Command interface & registry
        ├── errors.go          # Command errors
        ├── start_command.go   # /start command
        ├── help_command.go    # /help command
        ├── save_command.go    # /save command
        ├── search_command.go  # /search command
        ├── recent_command.go  # /recent command
        └── stats_command.go   # /stats command
```

**Key Principles**:
- Depends on application layer
- Handles user input/output
- No business logic (delegates to use cases)

---

## Design Patterns Implemented

### 1. Repository Pattern

**Purpose**: Abstract data access logic from business logic

**Implementation**:
- Interface: `internal/domain/repository/memory_repository.go`
- Implementation: `internal/infrastructure/persistence/sqlite/memory_repository.go`

**Benefits**:
- Business logic doesn't know about database
- Easy to swap SQLite for PostgreSQL/MongoDB
- Easier to test with mock repositories

**Example**:
```go
// Domain defines the interface
type MemoryRepository interface {
    Save(ctx context.Context, memory *entity.Memory) (int64, error)
    FindByID(ctx context.Context, id int) (*entity.Memory, error)
    // ... other methods
}

// Infrastructure implements it
type SQLiteMemoryRepository struct {
    conn *Connection
}
```

---

### 2. Strategy Pattern

**Purpose**: Define a family of algorithms and make them interchangeable

**Implementation**:
- Interface: `internal/infrastructure/search/strategy/search_strategy.go`
- Implementation: `internal/infrastructure/search/strategy/smart_strategy.go`

**Benefits**:
- Easy to add new search algorithms
- Can switch strategies at runtime
- Algorithms are independent and reusable

**Example**:
```go
// Strategy interface
type SearchStrategy interface {
    Search(ctx context.Context, query SearchQuery) ([]*entity.Memory, error)
    Name() string
}

// Concrete strategies
type SmartSearchStrategy struct { ... }
type FTS5Strategy struct { ... }
type ProximityStrategy struct { ... }
```

---

### 3. Command Pattern

**Purpose**: Encapsulate bot commands as objects

**Implementation**:
- Interface: `internal/presentation/handler/command/command.go`
- Implementations: `start_command.go`, `save_command.go`, etc.

**Benefits**:
- Easy to add new commands
- Commands are self-contained
- Can implement command history/undo
- Clean separation of command logic

**Example**:
```go
// Command interface
type Command interface {
    Name() string
    Description() string
    Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error
}

// Command registry manages all commands
type CommandRegistry struct {
    commands map[string]Command
}
```

---

### 4. Observer Pattern (Pub/Sub)

**Purpose**: Implement event-driven spaced repetition

**Implementation**:
- `internal/infrastructure/scheduler/spaced_repetition.go`

**Benefits**:
- Decoupled review notification system
- Can have multiple subscribers
- Event-driven architecture

**Example**:
```go
// Scheduler observes time and notifies users
type SpacedRepetitionScheduler struct {
    api       *tgbotapi.BotAPI
    useCase   *usecase.ReviewMemoryUseCase
    intervals []int
    ticker    *time.Ticker
}
```

---

### 5. Dependency Injection

**Purpose**: Manage dependencies and improve testability

**Implementation**:
- `cmd/bot/main.go` - Dependency wiring

**Benefits**:
- Loose coupling between components
- Easy to test with mocks
- Clear dependency graph
- No hard-coded dependencies

**Example**:
```go
// Dependencies are injected, not created inside
saveMemoryUC := usecase.NewSaveMemoryUseCase(memoryRepo)
searchStrategy := strategy.NewSmartSearchStrategy(memoryRepo)
searchMemoryUC := usecase.NewSearchMemoryUseCase(searchStrategy)
```

---

### 6. Factory Pattern (Implicit)

**Purpose**: Create command objects

**Implementation**:
- Command constructors (`NewStartCommand()`, `NewSaveCommand()`, etc.)

**Benefits**:
- Centralized object creation
- Encapsulates creation logic
- Easy to extend

---

## Complete Folder Structure

```
Personal-memory-reminder-bot/
├── cmd/
│   └── bot/
│       └── main.go                        # Application entry point
│
├── internal/
│   ├── domain/                            # DOMAIN LAYER
│   │   ├── entity/
│   │   │   ├── memory.go                  # Memory entity
│   │   │   └── errors.go                  # Domain errors
│   │   └── repository/
│   │       └── memory_repository.go       # Repository interface
│   │
│   ├── application/                       # APPLICATION LAYER
│   │   └── usecase/
│   │       ├── save_memory.go             # Save memory use case
│   │       ├── search_memory.go           # Search use case
│   │       ├── get_recent_memories.go     # Get recent use case
│   │       ├── get_stats.go               # Statistics use case
│   │       └── review_memory.go           # Review use case
│   │
│   ├── infrastructure/                    # INFRASTRUCTURE LAYER
│   │   ├── persistence/
│   │   │   └── sqlite/
│   │   │       ├── connection.go          # DB connection
│   │   │       └── memory_repository.go   # SQLite implementation
│   │   ├── search/
│   │   │   └── strategy/
│   │   │       ├── search_strategy.go     # Strategy interface
│   │   │       └── smart_strategy.go      # Smart search
│   │   ├── messaging/
│   │   │   └── telegram/
│   │   │       └── bot.go                 # Telegram adapter
│   │   └── scheduler/
│   │       └── spaced_repetition.go       # Review scheduler
│   │
│   └── presentation/                      # PRESENTATION LAYER
│       └── handler/
│           └── command/
│               ├── command.go             # Command pattern
│               ├── errors.go              # Command errors
│               ├── start_command.go       # /start handler
│               ├── help_command.go        # /help handler
│               ├── save_command.go        # /save handler
│               ├── search_command.go      # /search handler
│               ├── recent_command.go      # /recent handler
│               └── stats_command.go       # /stats handler
│
├── pkg/
│   └── config/
│       └── config.go                      # Configuration
│
├── bot/                                   # OLD STRUCTURE (deprecated)
│   ├── bot.go                             # Old bot implementation
│   └── spaced_repetition.go              # Old scheduler
│
├── database/                              # OLD STRUCTURE (deprecated)
│   ├── database.go                        # Old DB code
│   └── queries.go                         # Old queries
│
├── config/                                # OLD STRUCTURE (deprecated)
│   └── config.go                          # Old config
│
├── main.go                                # OLD ENTRY POINT (deprecated)
├── go.mod
├── go.sum
├── .env
├── .env.example
├── Makefile
└── README.md
```

---

## Dependency Flow

```
Presentation Layer (Commands)
        ↓ uses
Application Layer (Use Cases)
        ↓ uses
Domain Layer (Entities & Interfaces)
        ↑ implements
Infrastructure Layer (SQLite, Telegram, etc.)
```

**Key Point**: The dependency arrow always points **inward**. Inner layers never depend on outer layers.

---

## Example: Save Memory Flow

```
User sends message
        ↓
Telegram Bot (Infrastructure)
        ↓
SaveCommand (Presentation)
        ↓
SaveMemoryUseCase (Application)
        ↓
Memory Entity (Domain) ← validates business rules
        ↓
MemoryRepository Interface (Domain)
        ↑
SQLiteMemoryRepository (Infrastructure) ← saves to database
```

---

## Running the Application

### With New Architecture

```bash
# Build and run from cmd/bot
cd /home/milanmadusanka/Pictures/Personal-memory-reminder-bot
go run -tags "fts5" cmd/bot/main.go

# Or build binary
go build -tags "fts5" -o memory-bot cmd/bot/main.go
./memory-bot
```

### With Old Entry Point (Deprecated)

```bash
# Old way (still works but deprecated)
go run -tags "fts5" main.go
```

---

## Benefits of This Architecture

### 1. **Maintainability**
- Clear separation of concerns
- Easy to locate functionality
- Changes in one layer don't affect others

### 2. **Testability**
- Each layer can be tested independently
- Easy to mock dependencies
- High test coverage achievable

### 3. **Flexibility**
- Easy to swap implementations (SQLite → PostgreSQL)
- Add new commands without modifying existing code
- New search strategies without touching business logic

### 4. **Scalability**
- Can add new features without refactoring
- Clear boundaries for team collaboration
- Easy to split into microservices if needed

### 5. **Code Quality**
- SOLID principles applied throughout
- Design patterns for common problems
- Self-documenting structure

---

## Migration Guide

### From Old to New

**Old Structure:**
```go
// main.go
db := database.NewDatabase()
bot := bot.NewBot(db)
bot.Start()
```

**New Structure:**
```go
// cmd/bot/main.go
// 1. Infrastructure Setup
dbConn := sqlite.NewConnection(dbPath)
memoryRepo := sqlite.NewMemoryRepository(dbConn)

// 2. Use Cases
saveUC := usecase.NewSaveMemoryUseCase(memoryRepo)

// 3. Presentation
registry := command.NewCommandRegistry()
registry.Register(command.NewSaveCommand(saveUC))

// 4. Start
bot := telegram.NewBot(token, registry, saveUC)
bot.Start()
```

### Backward Compatibility

- Old `main.go` still works (uses old structure)
- New `cmd/bot/main.go` uses clean architecture
- Both use the same database file
- Gradual migration supported

---

## Adding New Features

### Example: Adding a /delete Command

**1. Create Command Handler** (`presentation/handler/command/delete_command.go`):
```go
type DeleteCommand struct {
    useCase *usecase.DeleteMemoryUseCase
}

func (c *DeleteCommand) Execute(ctx context.Context, bot BotAPI, msg *tgbotapi.Message) error {
    // Implementation
}
```

**2. Create Use Case** (`application/usecase/delete_memory.go`):
```go
type DeleteMemoryUseCase struct {
    repo repository.MemoryRepository
}

func (uc *DeleteMemoryUseCase) Execute(ctx context.Context, input DeleteInput) error {
    return uc.repo.Delete(ctx, input.ID, input.UserID)
}
```

**3. Register Command** (`cmd/bot/main.go`):
```go
deleteUC := usecase.NewDeleteMemoryUseCase(memoryRepo)
registry.Register(command.NewDeleteCommand(deleteUC))
```

**No changes needed in**: Domain layer, Infrastructure layer (already has Delete method)

---

## Testing Strategy

### Unit Tests

**Domain Layer**:
```go
func TestMemory_NeedsReview(t *testing.T) {
    memory := entity.NewMemory(123, 456, "Test")
    intervals := []int{1, 3, 7}
    
    // Test business logic
    assert.False(t, memory.NeedsReview(intervals))
}
```

**Application Layer**:
```go
func TestSaveMemoryUseCase(t *testing.T) {
    mockRepo := &MockMemoryRepository{}
    useCase := usecase.NewSaveMemoryUseCase(mockRepo)
    
    // Test use case logic
    output, err := useCase.Execute(ctx, input)
    assert.NoError(t, err)
}
```

### Integration Tests

```go
func TestMemoryRepository_Save(t *testing.T) {
    conn := sqlite.NewConnection(":memory:")
    repo := sqlite.NewMemoryRepository(conn)
    
    memory := entity.NewMemory(123, 456, "Test")
    id, err := repo.Save(ctx, memory)
    
    assert.NoError(t, err)
    assert.Greater(t, id, int64(0))
}
```

---

## Performance Considerations

1. **Database**: Uses SQLite with WAL mode for better concurrency
2. **Indexes**: Composite index on (user_id, created_at DESC)
3. **FTS5**: Full-text search with Porter stemming
4. **Connection Pooling**: Single connection per instance (SQLite limitation)
5. **Caching**: Can add repository caching layer if needed

---

## Future Enhancements

1. **Add More Strategies**:
   - Semantic search using embeddings
   - Fuzzy search
   - Tag-based search strategy

2. **Add Middleware**:
   - Logging middleware
   - Authentication middleware
   - Rate limiting

3. **Add Events**:
   - Event bus for domain events
   - Publish events when memories are saved/deleted
   - Subscribe to events for analytics

4. **Add Services**:
   - Email notification service
   - Export service (PDF, CSV)
   - Backup service

5. **Add Tests**:
   - Unit tests for all layers
   - Integration tests
   - End-to-end tests

---

## Conclusion

This clean architecture implementation provides:

✅ **Clear separation of concerns**  
✅ **6 design patterns** (Repository, Strategy, Command, Observer, DI, Factory)  
✅ **4-layer architecture** (Domain, Application, Infrastructure, Presentation)  
✅ **High maintainability** and testability  
✅ **Easy to extend** with new features  
✅ **Professional code organization**  

All existing functionality is preserved while dramatically improving code quality and maintainability.

---

**Built with 🧠 and ❤️ using Clean Architecture principles**
