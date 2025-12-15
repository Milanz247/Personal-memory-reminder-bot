package strategy

import (
	"context"
	"log"
	"strings"

	"memory-bot/internal/domain/entity"
	"memory-bot/internal/domain/repository"
	"memory-bot/internal/domain/service"
)

// SmartSearchStrategy implements intelligent search with fallback strategies
// Enhanced with biological contextual recall (Hippocampus function)
// 1. Try contextual search if user provides time/day cues
// 2. Try primary FTS5 search
// 3. If no results, try AND search with wildcards
// 4. If still no results, try OR search for broader results
type SmartSearchStrategy struct {
	repo           repository.MemoryRepository
	contextService *service.ContextualMetadataService
}

// NewSmartSearchStrategy creates a new smart search strategy
func NewSmartSearchStrategy(repo repository.MemoryRepository) *SmartSearchStrategy {
	return &SmartSearchStrategy{
		repo:           repo,
		contextService: service.NewContextualMetadataService(),
	}
}

// Search executes smart search with contextual awareness and fallback strategies
func (s *SmartSearchStrategy) Search(ctx context.Context, query SearchQuery) ([]*entity.Memory, error) {
	log.Printf("SmartSearch: Starting search for userID=%d, keyword='%s'", query.UserID, query.Keyword)

	opts := repository.SearchOptions{
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	// Step 1: Check for contextual cues (Biological principle: Associative recall)
	contextData, hasContext := s.contextService.ExtractContextCue(query.Keyword)
	if hasContext {
		log.Printf("SmartSearch: Detected contextual cue - %s", s.contextService.GetContextDescription(contextData))
		// Apply context filter directly at SQL level for better performance
		opts.ContextFilter = &contextData
	}

	// Step 2: Try primary search (with context filter if applicable)
	memories, err := s.repo.Search(ctx, query.UserID, query.Keyword, opts)
	if err != nil {
		log.Printf("SmartSearch: Primary search error: %v", err)
	}
	if err == nil && len(memories) > 0 {
		log.Printf("SmartSearch: Found %d results with primary search", len(memories))
		return memories, nil
	}

	// For fallbacks, reset context filter to broaden search
	opts.ContextFilter = nil

	// Fallback 1: Try AND search
	words := strings.Fields(strings.TrimSpace(query.Keyword))
	if len(words) > 1 {
		andTerms := make([]string, len(words))
		for i, word := range words {
			andTerms[i] = word + "*"
		}
		fallbackQuery := strings.Join(andTerms, " ")

		log.Printf("SmartSearch: Trying AND fallback with term: %s", fallbackQuery)
		memories, err = s.repo.Search(ctx, query.UserID, fallbackQuery, opts)
		if err != nil {
			log.Printf("SmartSearch: AND fallback error: %v", err)
		}
		if err == nil && len(memories) > 0 {
			log.Printf("SmartSearch: Found %d results with AND fallback", len(memories))
			return memories, nil
		}
	}

	// Fallback 2: Try OR search (without wildcards in OR to avoid syntax errors)
	if len(words) > 1 {
		// FTS5 doesn't support wildcards with OR operator properly
		// Use individual terms without wildcards for OR search
		orTerms := make([]string, len(words))
		copy(orTerms, words)
		fallbackQuery := strings.Join(orTerms, " OR ")

		log.Printf("SmartSearch: Trying OR fallback with term: %s", fallbackQuery)
		memories, err = s.repo.Search(ctx, query.UserID, fallbackQuery, opts)
		if err != nil {
			log.Printf("SmartSearch: OR fallback error: %v", err)
		}
		if err == nil && len(memories) > 0 {
			log.Printf("SmartSearch: Found %d results with OR fallback", len(memories))
			return memories, nil
		}
	}

	// No results found
	log.Printf("SmartSearch: No results found for keyword '%s'", query.Keyword)
	return []*entity.Memory{}, nil
}

// Name returns the strategy name
func (s *SmartSearchStrategy) Name() string {
	return "SmartSearch"
}
