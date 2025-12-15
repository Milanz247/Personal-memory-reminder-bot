package usecase

import (
	"context"
	"memory-bot/internal/domain/entity"
	"memory-bot/internal/domain/repository"
	"memory-bot/internal/domain/service"
	"time"
)

// SaveMemoryInput represents the input for saving a memory
type SaveMemoryInput struct {
	UserID  int64
	ChatID  int64
	Content string
}

// SaveMemoryOutput represents the output after saving a memory
type SaveMemoryOutput struct {
	MemoryID        int64
	Tags            []string
	EmotionalWeight float64
	Context         string
}

// SaveMemoryUseCase handles the business logic for saving memories
// Simulates the Hippocampus encoding new memories with emotional and contextual tags
type SaveMemoryUseCase struct {
	repo              repository.MemoryRepository
	sentimentAnalyzer *service.SentimentAnalyzer
	contextService    *service.ContextualMetadataService
}

// NewSaveMemoryUseCase creates a new save memory use case
func NewSaveMemoryUseCase(repo repository.MemoryRepository) *SaveMemoryUseCase {
	return &SaveMemoryUseCase{
		repo:              repo,
		sentimentAnalyzer: service.NewSentimentAnalyzer(),
		contextService:    service.NewContextualMetadataService(),
	}
}

// Execute saves a new memory with biological encoding
func (uc *SaveMemoryUseCase) Execute(ctx context.Context, input SaveMemoryInput) (*SaveMemoryOutput, error) {
	// 1. Create new memory entity (Sensory Input Processing)
	memory := entity.NewMemory(input.UserID, input.ChatID, input.Content)

	// 2. Emotional Encoding (The Amygdala's role)
	// Analyze emotional content and tag the memory
	memory.EmotionalWeight = uc.sentimentAnalyzer.Analyze(input.Content)

	// 3. Contextual Encoding (The Hippocampus's role)
	// Capture when and where the memory was created
	contextData := uc.contextService.GetCurrentContext(time.Now(), "Telegram")
	memory.TimeOfDay = contextData.TimeOfDay
	memory.DayOfWeek = contextData.DayOfWeek
	memory.ChatSource = contextData.ChatSource

	// 4. Validate
	if err := memory.Validate(); err != nil {
		return nil, err
	}

	// 5. Consolidation (Save to Cortex - long-term storage)
	id, err := uc.repo.Save(ctx, memory)
	if err != nil {
		return nil, err
	}

	return &SaveMemoryOutput{
		MemoryID:        id,
		Tags:            memory.Tags,
		EmotionalWeight: memory.EmotionalWeight,
		Context:         uc.contextService.GetContextDescription(contextData),
	}, nil
}
