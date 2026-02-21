package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 7 * 24 * time.Hour

// cacheEntry wraps a ProblemData with a timestamp for freshness checks.
type cacheEntry struct {
	FetchedAt time.Time    `json:"fetched_at"`
	Data      *ProblemData `json:"data"`
}

// QuestionCache provides a provider-aware local cache for fetched problem data.
// Cache layout: ~/.config/warmup/cache/<provider>/<problem-id>.json
type QuestionCache struct {
	baseDir string
}

// NewQuestionCache creates a cache rooted at the given directory.
// If baseDir is empty, defaults to ~/.config/warmup/cache.
func NewQuestionCache(baseDir string) (*QuestionCache, error) {
	if baseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home dir: %w", err)
		}
		baseDir = filepath.Join(home, ".config", "warmup", "cache")
	}
	return &QuestionCache{baseDir: baseDir}, nil
}

// Get retrieves a cached problem, returning nil if not cached or stale.
func (c *QuestionCache) Get(provider, id string) *ProblemData {
	path := c.path(provider, id)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil
	}

	if time.Since(entry.FetchedAt) > cacheTTL {
		return nil // stale
	}

	return entry.Data
}

// Put stores a problem in the cache.
func (c *QuestionCache) Put(provider string, problem *ProblemData) error {
	dir := filepath.Join(c.baseDir, provider)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	entry := cacheEntry{
		FetchedAt: time.Now(),
		Data:      problem,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal cache entry: %w", err)
	}

	path := c.path(provider, problem.ID)
	return os.WriteFile(path, data, 0o644)
}

// path returns the file path for a cached problem.
func (c *QuestionCache) path(provider, id string) string {
	return filepath.Join(c.baseDir, provider, id+".json")
}

// FetchWithCache attempts to load from cache first, falling back to the provider.
func FetchWithCache(cache *QuestionCache, p Provider, id string, lang string) (*ProblemData, error) {
	// Try cache first.
	if cached := cache.Get(p.Name(), id); cached != nil {
		return cached, nil
	}

	// Fetch from provider.
	problem, err := p.FetchProblem(id, lang)
	if err != nil {
		return nil, err
	}

	// Store in cache (best-effort, don't fail on cache write errors).
	if cacheErr := cache.Put(p.Name(), problem); cacheErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to cache %s/%s: %v\n", p.Name(), id, cacheErr)
	}

	return problem, nil
}
