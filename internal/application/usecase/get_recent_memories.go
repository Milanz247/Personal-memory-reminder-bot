package usecase

import (
	"context"
	"memory-bot/internal/domain/entity"
	"memory-bot/internal/domain/repository"
)

// GetRecentMemoriesInput represents the input for getting recent memories
type GetRecentMemoriesInput struct {
	UserID int64
	Limit  int
}

// GetRecentMemoriesOutput represents the output
type GetRecentMemoriesOutput struct {
	Memories []*entity.Memory
}

// GetRecentMemoriesUseCase handles retrieving recent memories
type GetRecentMemoriesUseCase struct {
	repo repository.MemoryRepository
}

// NewGetRecentMemoriesUseCase creates a new use case
func NewGetRecentMemoriesUseCase(repo repository.MemoryRepository) *GetRecentMemoriesUseCase {
	return &GetRecentMemoriesUseCase{
		repo: repo,
	}
}

// Execute retrieves recent memories
func (uc *GetRecentMemoriesUseCase) Execute(ctx context.Context, input GetRecentMemoriesInput) (*GetRecentMemoriesOutput, error) {
	memories, err := uc.repo.GetRecent(ctx, input.UserID, input.Limit)
	if err != nil {
		return nil, err
	}

	return &GetRecentMemoriesOutput{
		Memories: memories,
	}, nil
}
