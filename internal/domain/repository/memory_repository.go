package repository

import (
	"context"
	"memory-bot/internal/domain/entity"
)

// SearchOptions defines options for memory search
type SearchOptions struct {
	Limit  int
	Offset int
}

// MemoryRepository defines the interface for memory data access
// This follows the Repository Pattern for data abstraction
type MemoryRepository interface {
	// Save stores a new memory
	Save(ctx context.Context, memory *entity.Memory) (int64, error)

	// FindByID retrieves a memory by its ID
	FindByID(ctx context.Context, id int) (*entity.Memory, error)

	// Search performs a search query with the given options
	Search(ctx context.Context, userID int64, query string, opts SearchOptions) ([]*entity.Memory, error)

	// GetRecent retrieves the most recent memories for a user
	GetRecent(ctx context.Context, userID int64, limit int) ([]*entity.Memory, error)

	// GetForReview retrieves memories that need review based on intervals
	GetForReview(ctx context.Context, intervals []int) ([]*entity.Memory, error)

	// Update updates an existing memory
	Update(ctx context.Context, memory *entity.Memory) error

	// Delete removes a memory (with authorization check)
	Delete(ctx context.Context, id int, userID int64) error

	// Count returns the total number of memories for a user
	Count(ctx context.Context, userID int64) (int, error)

	// Biological memory system methods

	// GetFragileMemories retrieves recently created memories that need consolidation
	GetFragileMemories(ctx context.Context) ([]*entity.Memory, error)

	// UpdateConsolidation updates consolidation-related fields
	UpdateConsolidation(ctx context.Context, memory *entity.Memory) error
}
