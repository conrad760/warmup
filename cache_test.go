package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCache_PutAndGet(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	problem := &ProblemData{
		ID:          "two-sum",
		Title:       "Two Sum",
		Description: "Given an array...",
		Difficulty:  "Easy",
		Tags:        []string{"Array", "Hash Table"},
		CodeSnippet: "func twoSum(nums []int, target int) []int {\n}",
	}

	if err := cache.Put("leetcode", problem); err != nil {
		t.Fatalf("Put: %v", err)
	}

	got := cache.Get("leetcode", "two-sum")
	if got == nil {
		t.Fatal("Get returned nil for cached problem")
	}
	if got.Title != "Two Sum" {
		t.Errorf("Title = %q, want %q", got.Title, "Two Sum")
	}
	if got.Description != "Given an array..." {
		t.Errorf("Description = %q, want %q", got.Description, "Given an array...")
	}
	if len(got.Tags) != 2 {
		t.Errorf("Tags len = %d, want 2", len(got.Tags))
	}
}

func TestCache_GetMissing(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	got := cache.Get("leetcode", "nonexistent")
	if got != nil {
		t.Error("Get should return nil for missing entry")
	}
}

func TestCache_ProviderIsolation(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	problem := &ProblemData{ID: "test-problem", Title: "Test"}
	if err := cache.Put("providerA", problem); err != nil {
		t.Fatalf("Put: %v", err)
	}

	// Same ID under different provider should be a miss.
	got := cache.Get("providerB", "test-problem")
	if got != nil {
		t.Error("Get should return nil for different provider")
	}

	// Same provider should hit.
	got = cache.Get("providerA", "test-problem")
	if got == nil {
		t.Error("Get should return cached data for same provider")
	}
}

func TestCache_TTLExpiry(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	// Write a cache entry with a timestamp in the past.
	entry := cacheEntry{
		FetchedAt: time.Now().Add(-(cacheTTL + time.Hour)), // expired
		Data:      &ProblemData{ID: "old-problem", Title: "Old"},
	}
	providerDir := filepath.Join(dir, "leetcode")
	if err := os.MkdirAll(providerDir, 0o755); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(entry)
	if err := os.WriteFile(filepath.Join(providerDir, "old-problem.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}

	got := cache.Get("leetcode", "old-problem")
	if got != nil {
		t.Error("Get should return nil for expired entry")
	}
}

func TestCache_CorruptedJSON(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	providerDir := filepath.Join(dir, "leetcode")
	if err := os.MkdirAll(providerDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(providerDir, "bad.json"), []byte("not json{{{"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := cache.Get("leetcode", "bad")
	if got != nil {
		t.Error("Get should return nil for corrupted cache entry")
	}
}

func TestFetchWithCache_UsesCacheThenProvider(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	mock := &MockProvider{}

	// First call should hit the provider and populate cache.
	pd, err := FetchWithCache(cache, mock, "fizz-buzz", "go")
	if err != nil {
		t.Fatalf("FetchWithCache: %v", err)
	}
	if pd.Title != "Fizz Buzz" {
		t.Errorf("Title = %q, want %q", pd.Title, "Fizz Buzz")
	}

	// Verify it's in cache now.
	cached := cache.Get("mock", "fizz-buzz")
	if cached == nil {
		t.Fatal("problem should be cached after FetchWithCache")
	}

	// Second call should use cache (even if provider were unavailable).
	pd2, err := FetchWithCache(cache, mock, "fizz-buzz", "go")
	if err != nil {
		t.Fatalf("FetchWithCache second call: %v", err)
	}
	if pd2.Title != "Fizz Buzz" {
		t.Errorf("Title = %q, want %q", pd2.Title, "Fizz Buzz")
	}
}
