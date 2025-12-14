package strategy

import (
	"context"
	"memory-bot/internal/domain/entity"
)

// SearchQuery encapsulates search parameters
type SearchQuery struct {
	UserID  int64
	Keyword string
	Limit   int
	Offset  int
}

// SearchStrategy defines the interface for different search algorithms
// This implements the Strategy Pattern for flexible search algorithms
type SearchStrategy interface {
	// Search executes the search with the specific strategy
	Search(ctx context.Context, query SearchQuery) ([]*entity.Memory, error)

	// Name returns the strategy name for logging/debugging
	Name() string
}

// SearchStrategyFactory creates search strategies
type SearchStrategyFactory interface {
	// CreateStrategy creates a strategy based on the search query
	CreateStrategy(query SearchQuery) SearchStrategy
}
