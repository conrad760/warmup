package main

import "testing"

func TestMockProvider_FetchKnown(t *testing.T) {
	p := &MockProvider{}

	tests := []struct {
		id        string
		wantTitle string
		wantDiff  string
	}{
		{"fizz-buzz", "Fizz Buzz", "Easy"},
		{"reverse-string", "Reverse String", "Easy"},
		{"valid-palindrome", "Valid Palindrome", "Easy"},
		{"binary-search", "Binary Search", "Easy"},
		{"merge-sort", "Merge Sort Implementation", "Medium"},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			pd, err := p.FetchProblem(tt.id, "go")
			if err != nil {
				t.Fatalf("FetchProblem(%q) error: %v", tt.id, err)
			}
			if pd.ID != tt.id {
				t.Errorf("ID = %q, want %q", pd.ID, tt.id)
			}
			if pd.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", pd.Title, tt.wantTitle)
			}
			if pd.Difficulty != tt.wantDiff {
				t.Errorf("Difficulty = %q, want %q", pd.Difficulty, tt.wantDiff)
			}
			if pd.Description == "" {
				t.Error("Description should not be empty")
			}
			if pd.CodeSnippet == "" {
				t.Error("CodeSnippet should not be empty")
			}
		})
	}
}

func TestMockProvider_FetchUnknown(t *testing.T) {
	p := &MockProvider{}
	_, err := p.FetchProblem("nonexistent-problem", "go")
	if err == nil {
		t.Fatal("expected error for unknown problem")
	}
}

func TestMockProvider_ReturnsCopy(t *testing.T) {
	p := &MockProvider{}
	pd1, _ := p.FetchProblem("fizz-buzz", "go")
	pd2, _ := p.FetchProblem("fizz-buzz", "go")

	// Mutating one should not affect the other.
	pd1.Title = "MUTATED"
	if pd2.Title == "MUTATED" {
		t.Error("FetchProblem should return a copy, not a reference to shared data")
	}
}
