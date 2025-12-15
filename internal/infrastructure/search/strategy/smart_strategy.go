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

	// Step 1: Check for hashtag search (exact tag matching)
	if strings.HasPrefix(strings.TrimSpace(query.Keyword), "#") {
		log.Printf("SmartSearch: Detected hashtag search")
		// Tag search is handled by FTS5 with high precision
		memories, err := s.repo.Search(ctx, query.UserID, query.Keyword, opts)
		if err == nil && len(memories) > 0 {
			log.Printf("SmartSearch: Found %d results with tag search", len(memories))
			return memories, nil
		}
	}

	// Step 2: Check for contextual cues (Biological principle: Associative recall)
	contextData, hasContext := s.contextService.ExtractContextCue(query.Keyword)
	if hasContext {
		log.Printf("SmartSearch: Detected contextual cue - %s", s.contextService.GetContextDescription(contextData))
		// Apply context filter directly at SQL level for better performance
		opts.ContextFilter = &contextData
	}

	// Step 3: Try primary FTS5 search with wildcard (with context filter if applicable)
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

	// Step 4: Try fuzzy search with shorter substrings (3+ chars)
	words := strings.Fields(strings.TrimSpace(query.Keyword))
	if len(words) == 1 && len(words[0]) >= 3 {
		word := words[0]
		// Try searching with first 3 characters for fuzzy matching
		fuzzyQuery := word[:3] + "*"
		log.Printf("SmartSearch: Trying fuzzy search with term: %s", fuzzyQuery)
		memories, err = s.repo.Search(ctx, query.UserID, fuzzyQuery, opts)
		if err == nil && len(memories) > 0 {
			log.Printf("SmartSearch: Found %d results with fuzzy search", len(memories))
			return memories, nil
		}
	}

	// Step 5: Try AND search (all words must match)
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

	// Step 6: Try partial match (any word matches)
	if len(words) > 1 {
		// Try each word individually with wildcards
		for _, word := range words {
			if len(word) < 3 {
				continue // Skip very short words
			}
			partialQuery := word + "*"
			log.Printf("SmartSearch: Trying partial match with term: %s", partialQuery)
			memories, err = s.repo.Search(ctx, query.UserID, partialQuery, opts)
			if err == nil && len(memories) > 0 {
				log.Printf("SmartSearch: Found %d results with partial match", len(memories))
				return memories, nil
			}
		}
	}

	// Step 7: Try OR search (broader search - any word matches)
	if len(words) > 1 {
		// FTS5 OR search without wildcards for maximum recall
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

	// Step 8: Last resort - search with NEAR operator for proximity matching
	if len(words) > 1 {
		nearQuery := "NEAR(" + strings.Join(words, " ") + ", 10)"
		log.Printf("SmartSearch: Trying NEAR proximity search with term: %s", nearQuery)
		memories, err = s.repo.Search(ctx, query.UserID, nearQuery, opts)
		if err == nil && len(memories) > 0 {
			log.Printf("SmartSearch: Found %d results with NEAR search", len(memories))
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
