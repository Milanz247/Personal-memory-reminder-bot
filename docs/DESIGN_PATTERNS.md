# Design Patterns Documentation

This document provides a detailed explanation of the design patterns used in the Personal Memory Reminder Bot refactoring.

---

## Table of Contents

1. [Repository Pattern](#1-repository-pattern)
2. [Strategy Pattern](#2-strategy-pattern)
3. [Command Pattern](#3-command-pattern)
4. [Observer Pattern](#4-observer-pattern)
5. [Dependency Injection](#5-dependency-injection)
6. [Factory Pattern](#6-factory-pattern)

---

## 1. Repository Pattern

### Purpose
Abstract the data access layer from business logic, providing a collection-like interface for accessing domain objects.

### Problem Solved
- Business logic was tightly coupled to database implementation
- Difficult to test without a real database
- Hard to switch database technologies

### Implementation

**Interface (Domain Layer)**:
```go
// File: internal/domain/repository/memory_repository.go
package repository

type MemoryRepository interface {
    Save(ctx context.Context, memory *entity.Memory) (int64, error)
    FindByID(ctx context.Context, id int) (*entity.Memory, error)
    Search(ctx context.Context, userID int64, query string, opts SearchOptions) ([]*entity.Memory, error)
    GetRecent(ctx context.Context, userID int64, limit int) ([]*entity.Memory, error)
    GetForReview(ctx context.Context, intervals []int) ([]*entity.Memory, error)
    Update(ctx context.Context, memory *entity.Memory) error
    Delete(ctx context.Context, id int, userID int64) error
    Count(ctx context.Context, userID int64) (int, error)
}
```

**Implementation (Infrastructure Layer)**:
```go
// File: internal/infrastructure/persistence/sqlite/memory_repository.go
package sqlite

type MemoryRepository struct {
    conn *Connection
}

func (r *MemoryRepository) Save(ctx context.Context, memory *entity.Memory) (int64, error) {
    // SQLite-specific implementation
}

// ... other methods
```

**Usage (Application Layer)**:
```go
// File: internal/application/usecase/save_memory.go
package usecase

type SaveMemoryUseCase struct {
    repo repository.MemoryRepository  // Depends on interface, not concrete implementation
}

func (uc *SaveMemoryUseCase) Execute(ctx context.Context, input SaveMemoryInput) (*SaveMemoryOutput, error) {
    memory := entity.NewMemory(input.UserID, input.ChatID, input.Content)
    id, err := uc.repo.Save(ctx, memory)  // Uses interface method
    // ...
}
```

### Benefits
✅ **Testability**: Easy to create mock repositories for testing  
✅ **Flexibility**: Can swap SQLite for PostgreSQL without changing business logic  
✅ **Separation**: Database concerns separated from business logic  
✅ **Single Responsibility**: Repository only handles data access  

### Testing Example
```go
// Mock repository for testing
type MockMemoryRepository struct {
    SaveFunc func(ctx context.Context, memory *entity.Memory) (int64, error)
}

func (m *MockMemoryRepository) Save(ctx context.Context, memory *entity.Memory) (int64, error) {
    return m.SaveFunc(ctx, memory)
}

// Test
func TestSaveMemoryUseCase(t *testing.T) {
    mockRepo := &MockMemoryRepository{
        SaveFunc: func(ctx context.Context, memory *entity.Memory) (int64, error) {
            return 1, nil
        },
    }
    
    useCase := usecase.NewSaveMemoryUseCase(mockRepo)
    // Test without real database
}
```

---

## 2. Strategy Pattern

### Purpose
Define a family of algorithms, encapsulate each one, and make them interchangeable.

### Problem Solved
- Multiple search algorithms hardcoded in queries
- Difficult to add new search strategies
- No way to select algorithm at runtime

### Implementation

**Strategy Interface**:
```go
// File: internal/infrastructure/search/strategy/search_strategy.go
package strategy

type SearchStrategy interface {
    Search(ctx context.Context, query SearchQuery) ([]*entity.Memory, error)
    Name() string
}

type SearchQuery struct {
    UserID  int64
    Keyword string
    Limit   int
    Offset  int
}
```

**Concrete Strategies**:
```go
// File: internal/infrastructure/search/strategy/smart_strategy.go
package strategy

type SmartSearchStrategy struct {
    repo repository.MemoryRepository
}

func (s *SmartSearchStrategy) Search(ctx context.Context, query SearchQuery) ([]*entity.Memory, error) {
    // 1. Try primary FTS5 search
    // 2. Fallback to AND search
    // 3. Fallback to OR search
}

func (s *SmartSearchStrategy) Name() string {
    return "SmartSearch"
}
```

**Other Strategies (Can be added)**:
```go
type FTS5Strategy struct { ... }        // Basic FTS5 search
type ProximityStrategy struct { ... }   // NEAR searches
type SemanticStrategy struct { ... }    // Embedding-based search
type FuzzyStrategy struct { ... }       // Fuzzy matching
```

**Usage**:
```go
// File: cmd/bot/main.go

// Choose strategy at startup
searchStrategy := strategy.NewSmartSearchStrategy(memoryRepo)
searchMemoryUC := usecase.NewSearchMemoryUseCase(searchStrategy)

// Can easily switch strategies:
// searchStrategy := strategy.NewFTS5Strategy(memoryRepo)
// searchStrategy := strategy.NewSemanticStrategy(memoryRepo)
```

### Benefits
✅ **Open/Closed Principle**: Open for extension, closed for modification  
✅ **Runtime Selection**: Can choose strategy based on user preferences  
✅ **Independent Algorithms**: Each strategy is self-contained  
✅ **Easy Testing**: Test each strategy independently  

---

## 3. Command Pattern

### Purpose
Encapsulate a request as an object, allowing parameterization and queuing of requests.

### Problem Solved
- Bot commands mixed with handler logic
- Difficult to add new commands
- No command history or undo capability

### Implementation

**Command Interface**:
```go
// File: internal/presentation/handler/command/command.go
package command

type Command interface {
    Name() string
    Description() string
    Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error
}

type BotAPI interface {
    Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
```

**Command Registry**:
```go
type CommandRegistry struct {
    commands map[string]Command
}

func (r *CommandRegistry) Register(cmd Command) {
    r.commands[cmd.Name()] = cmd
}

func (r *CommandRegistry) Execute(ctx context.Context, name string, bot BotAPI, msg *tgbotapi.Message) error {
    cmd, exists := r.Get(name)
    if !exists {
        return ErrCommandNotFound
    }
    return cmd.Execute(ctx, bot, msg)
}
```

**Concrete Commands**:
```go
// File: internal/presentation/handler/command/save_command.go
package command

type SaveCommand struct {
    useCase *usecase.SaveMemoryUseCase
}

func (c *SaveCommand) Name() string {
    return "save"
}

func (c *SaveCommand) Description() string {
    return "Save a memory"
}

func (c *SaveCommand) Execute(ctx context.Context, bot BotAPI, message *tgbotapi.Message) error {
    input := usecase.SaveMemoryInput{
        UserID:  message.From.ID,
        ChatID:  message.Chat.ID,
        Content: message.CommandArguments(),
    }
    
    output, err := c.useCase.Execute(ctx, input)
    // Handle response...
}
```

### Benefits
✅ **Decoupling**: Invoker doesn't know about concrete commands  
✅ **Extensibility**: Easy to add new commands  
✅ **Single Responsibility**: Each command has one job  

---

## 4. Observer Pattern

### Purpose
Define a one-to-many dependency so that when one object changes state, all dependents are notified.

### Implementation

**Subject (Scheduler)**:
```go
// File: internal/infrastructure/scheduler/spaced_repetition.go
package scheduler

type SpacedRepetitionScheduler struct {
    api       *tgbotapi.BotAPI
    useCase   *usecase.ReviewMemoryUseCase
    intervals []int
    ticker    *time.Ticker
    stopChan  chan bool
}

func (s *SpacedRepetitionScheduler) checkAndSendReviews() {
    memories, _ := s.useCase.Execute(ctx, input)
    
    // Notify users (observers)
    for userID, mems := range groupByUser(memories) {
        s.sendReviewToUser(ctx, userID, mems)
    }
}
```

### Benefits
✅ **Decoupling**: Subject doesn't know about concrete observers  
✅ **Multiple Observers**: Can notify via Telegram, Email, SMS  
✅ **Event-Driven**: Clean event-driven architecture  

---

## 5. Dependency Injection

### Purpose
Provide dependencies from the outside rather than creating them internally.

### Implementation

**Dependency Wiring (main.go)**:
```go
// File: cmd/bot/main.go

func main() {
    // 1. Infrastructure
    dbConn := sqlite.NewConnection(cfg.DBPath)
    memoryRepo := sqlite.NewMemoryRepository(dbConn)
    
    // 2. Application
    saveMemoryUC := usecase.NewSaveMemoryUseCase(memoryRepo)
    searchStrategy := strategy.NewSmartSearchStrategy(memoryRepo)
    searchMemoryUC := usecase.NewSearchMemoryUseCase(searchStrategy)
    
    // 3. Presentation
    registry := command.NewCommandRegistry()
    registry.Register(command.NewSaveCommand(saveMemoryUC))
    
    // 4. Start application
    bot := telegram.NewBot(token, registry, saveMemoryUC)
    bot.Start()
}
```

### Benefits
✅ **Testability**: Easy to inject mocks for testing  
✅ **Flexibility**: Can swap implementations  
✅ **Clear Dependencies**: Explicit dependency graph  

---

## 6. Factory Pattern

### Purpose
Define an interface for creating objects.

### Implementation

**Simple Factory**:
```go
func NewStartCommand() *StartCommand {
    return &StartCommand{}
}

func NewSaveCommand(useCase *usecase.SaveMemoryUseCase) *SaveCommand {
    return &SaveCommand{
        useCase: useCase,
    }
}
```

### Benefits
✅ **Centralized Creation**: Object creation in one place  
✅ **Encapsulation**: Creation logic hidden  
✅ **Extensibility**: Easy to add new types  

---

## Summary

| Pattern | Purpose | Location | Benefit |
|---------|---------|----------|---------|
| **Repository** | Data access abstraction | `domain/repository/` | Database-agnostic business logic |
| **Strategy** | Interchangeable algorithms | `infrastructure/search/` | Easy to add new search methods |
| **Command** | Encapsulate requests | `presentation/handler/` | Easy to add new bot commands |
| **Observer** | Event notification | `infrastructure/scheduler/` | Decoupled review system |
| **Dependency Injection** | Provide dependencies | `cmd/bot/main.go` | Testable, flexible |
| **Factory** | Object creation | Constructor functions | Centralized creation logic |

---

**These patterns work together to create a maintainable, testable, and scalable architecture.**
