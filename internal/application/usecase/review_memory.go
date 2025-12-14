package usecase

import (
	"context"
	"memory-bot/internal/domain/entity"
	"memory-bot/internal/domain/repository"
)

// ReviewMemoryInput represents the input for reviewing memories
type ReviewMemoryInput struct {
	Intervals []int
}

// ReviewMemoryOutput represents memories that need review
type ReviewMemoryOutput struct {
	Memories []*entity.Memory
}

// ReviewMemoryUseCase handles the spaced repetition review logic
type ReviewMemoryUseCase struct {
	repo repository.MemoryRepository
}

// NewReviewMemoryUseCase creates a new review memory use case
func NewReviewMemoryUseCase(repo repository.MemoryRepository) *ReviewMemoryUseCase {
	return &ReviewMemoryUseCase{
		repo: repo,
	}
}

// Execute retrieves memories that need review
func (uc *ReviewMemoryUseCase) Execute(ctx context.Context, input ReviewMemoryInput) (*ReviewMemoryOutput, error) {
	memories, err := uc.repo.GetForReview(ctx, input.Intervals)
	if err != nil {
		return nil, err
	}

	return &ReviewMemoryOutput{
		Memories: memories,
	}, nil
}

// MarkAsReviewed marks a memory as reviewed
func (uc *ReviewMemoryUseCase) MarkAsReviewed(ctx context.Context, memoryID int) error {
	memory, err := uc.repo.FindByID(ctx, memoryID)
	if err != nil {
		return err
	}

	memory.MarkAsReviewed()

	return uc.repo.Update(ctx, memory)
}
