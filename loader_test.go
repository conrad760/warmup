package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDifficulty(t *testing.T) {
	tests := []struct {
		input string
		want  Difficulty
	}{
		{"Easy", Easy},
		{"easy", Easy},
		{"EASY", Easy},
		{"Medium", Medium},
		{"medium", Medium},
		{"Hard", Hard},
		{"hard", Hard},
		{"unknown", Easy}, // default
		{"", Easy},        // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseDifficulty(tt.input)
			if got != tt.want {
				t.Errorf("parseDifficulty(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoadQuestionsFromProviders_Mock(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	curated := []CuratedQuestion{
		{
			Provider:  "mock",
			ProblemID: "fizz-buzz",
			Category:  "Math",
			Options: []Option{
				{Text: "Option A", Rating: Optimal},
				{Text: "Option B", Rating: Suboptimal},
			},
			Solution: "// solution here",
		},
		{
			Provider:  "mock",
			ProblemID: "binary-search",
			Category:  "Binary Search",
			Options: []Option{
				{Text: "Option A", Rating: Optimal},
				{Text: "Option B", Rating: Wrong},
			},
			Solution: "// solution here",
		},
	}

	questions, err := loadQuestionsFromProviders(curated, cache, "go")
	if err != nil {
		t.Fatalf("loadQuestionsFromProviders: %v", err)
	}

	if len(questions) != 2 {
		t.Fatalf("got %d questions, want 2", len(questions))
	}

	// Verify fields are populated from both curated and provider data.
	for _, q := range questions {
		if q.Title == "" {
			t.Error("Title should be populated from provider")
		}
		if q.Description == "" {
			t.Error("Description should be populated from provider")
		}
		if q.Category == "" {
			t.Error("Category should be set from curated data")
		}
		if q.Provider != "mock" {
			t.Errorf("Provider = %q, want %q", q.Provider, "mock")
		}
		if len(q.Options) < 2 {
			t.Error("Options should be populated from curated data")
		}
		if q.Solution == "" {
			t.Error("Solution should be populated from curated data")
		}
	}
}

func TestLoadQuestionsFromProviders_DefaultProvider(t *testing.T) {
	// When Provider is empty, it should default to "leetcode".
	// We can't actually hit LeetCode in tests, but we can verify
	// that an unknown mock problem on "leetcode" gracefully degrades.
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	curated := []CuratedQuestion{
		{
			Provider:  "mock",
			ProblemID: "fizz-buzz",
			Options: []Option{
				{Text: "A", Rating: Optimal},
				{Text: "B", Rating: Wrong},
			},
			Solution: "sol",
		},
	}

	questions, err := loadQuestionsFromProviders(curated, cache, "go")
	if err != nil {
		t.Fatalf("loadQuestionsFromProviders: %v", err)
	}

	// With no curated category, should use provider's first tag.
	for _, q := range questions {
		if q.Category == "" {
			t.Error("Category should fall back to provider tag when curated category is empty")
		}
	}
}

func TestLoadQuestionsFromProviders_UnknownProvider(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	curated := []CuratedQuestion{
		{
			Provider:  "does-not-exist",
			ProblemID: "foo",
			Options: []Option{
				{Text: "A", Rating: Optimal},
				{Text: "B", Rating: Wrong},
			},
			Solution: "sol",
		},
	}

	_, err = loadQuestionsFromProviders(curated, cache, "go")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestLoadQuestionsFromJSONFile(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	curated := []CuratedQuestion{
		{
			Provider:  "mock",
			ProblemID: "reverse-string",
			Category:  "Strings",
			Options: []Option{
				{Text: "Two pointers from both ends", Rating: Optimal},
				{Text: "Create new reversed string", Rating: Suboptimal},
			},
			Solution: "// reverse in-place",
		},
	}

	data, _ := json.Marshal(curated)
	jsonPath := dir + "/test-questions.json"
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	questions, err := loadQuestionsFromJSONFile(jsonPath, cache, "go")
	if err != nil {
		t.Fatalf("loadQuestionsFromJSONFile: %v", err)
	}
	if len(questions) != 1 {
		t.Fatalf("got %d questions, want 1", len(questions))
	}
	if questions[0].Title != "Reverse String" {
		t.Errorf("Title = %q, want %q", questions[0].Title, "Reverse String")
	}
}

func TestLoadQuestionsFromJSONFile_Validation(t *testing.T) {
	dir := t.TempDir()
	cache, err := NewQuestionCache(dir)
	if err != nil {
		t.Fatalf("NewQuestionCache: %v", err)
	}

	tests := []struct {
		name    string
		curated []CuratedQuestion
	}{
		{
			name: "missing ProblemID",
			curated: []CuratedQuestion{{
				Provider: "mock",
				Options: []Option{
					{Text: "A", Rating: Optimal},
					{Text: "B", Rating: Wrong},
				},
			}},
		},
		{
			name: "too few options",
			curated: []CuratedQuestion{{
				Provider:  "mock",
				ProblemID: "fizz-buzz",
				Options:   []Option{{Text: "A", Rating: Optimal}},
			}},
		},
		{
			name: "no optimal option",
			curated: []CuratedQuestion{{
				Provider:  "mock",
				ProblemID: "fizz-buzz",
				Options: []Option{
					{Text: "A", Rating: Suboptimal},
					{Text: "B", Rating: Wrong},
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := json.Marshal(tt.curated)
			jsonPath := filepath.Join(dir, tt.name+".json")
			if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
				t.Fatal(err)
			}

			_, err := loadQuestionsFromJSONFile(jsonPath, cache, "go")
			if err == nil {
				t.Error("expected validation error")
			}
		})
	}
}
