package usecase

import (
	"context"
	"memory-bot/internal/domain/entity"
	"memory-bot/internal/domain/repository"
)

// SaveMemoryInput represents the input for saving a memory
type SaveMemoryInput struct {
	UserID  int64
	ChatID  int64
	Content string
}

// SaveMemoryOutput represents the output after saving a memory
type SaveMemoryOutput struct {
	MemoryID int64
	Tags     []string
}

// SaveMemoryUseCase handles the business logic for saving memories
type SaveMemoryUseCase struct {
	repo repository.MemoryRepository
}

// NewSaveMemoryUseCase creates a new save memory use case
func NewSaveMemoryUseCase(repo repository.MemoryRepository) *SaveMemoryUseCase {
	return &SaveMemoryUseCase{
		repo: repo,
	}
}

// Execute saves a new memory
func (uc *SaveMemoryUseCase) Execute(ctx context.Context, input SaveMemoryInput) (*SaveMemoryOutput, error) {
	// Create new memory entity
	memory := entity.NewMemory(input.UserID, input.ChatID, input.Content)

	// Validate
	if err := memory.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	id, err := uc.repo.Save(ctx, memory)
	if err != nil {
		return nil, err
	}

	return &SaveMemoryOutput{
		MemoryID: id,
		Tags:     memory.Tags,
	}, nil
}
