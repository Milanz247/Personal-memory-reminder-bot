package usecase

import (
	"context"
	"memory-bot/internal/domain/repository"
)

// GetStatsInput represents the input for getting statistics
type GetStatsInput struct {
	UserID int64
}

// GetStatsOutput represents the statistics output
type GetStatsOutput struct {
	TotalMemories int
}

// GetStatsUseCase handles retrieving user statistics
type GetStatsUseCase struct {
	repo repository.MemoryRepository
}

// NewGetStatsUseCase creates a new use case
func NewGetStatsUseCase(repo repository.MemoryRepository) *GetStatsUseCase {
	return &GetStatsUseCase{
		repo: repo,
	}
}

// Execute retrieves user statistics
func (uc *GetStatsUseCase) Execute(ctx context.Context, input GetStatsInput) (*GetStatsOutput, error) {
	count, err := uc.repo.Count(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &GetStatsOutput{
		TotalMemories: count,
	}, nil
}
