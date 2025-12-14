package usecase

import (
	"context"
	"memory-bot/internal/domain/entity"
	"memory-bot/internal/infrastructure/search/strategy"
)

// SearchMemoryInput represents the input for searching memories
type SearchMemoryInput struct {
	UserID  int64
	Keyword string
	Limit   int
	Offset  int
}

// SearchMemoryOutput represents the output after searching memories
type SearchMemoryOutput struct {
	Memories []*entity.Memory
	Total    int
	HasMore  bool
}

// SearchMemoryUseCase handles the business logic for searching memories
type SearchMemoryUseCase struct {
	strategy strategy.SearchStrategy
}

// NewSearchMemoryUseCase creates a new search memory use case
func NewSearchMemoryUseCase(searchStrategy strategy.SearchStrategy) *SearchMemoryUseCase {
	return &SearchMemoryUseCase{
		strategy: searchStrategy,
	}
}

// Execute searches for memories using the configured strategy
func (uc *SearchMemoryUseCase) Execute(ctx context.Context, input SearchMemoryInput) (*SearchMemoryOutput, error) {
	// Create search query
	query := strategy.SearchQuery{
		UserID:  input.UserID,
		Keyword: input.Keyword,
		Limit:   input.Limit + 1, // Fetch one extra to check if there are more
		Offset:  input.Offset,
	}

	// Execute search
	memories, err := uc.strategy.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	// Check if there are more results
	hasMore := len(memories) > input.Limit
	if hasMore {
		memories = memories[:input.Limit]
	}

	return &SearchMemoryOutput{
		Memories: memories,
		Total:    len(memories),
		HasMore:  hasMore,
	}, nil
}
