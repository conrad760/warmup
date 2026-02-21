package main

import (
	"strings"
	"testing"
)

func TestGetProvider_Registered(t *testing.T) {
	// "mock" and "leetcode" are registered via init() in their respective files.
	for _, name := range []string{"mock", "leetcode"} {
		p, err := GetProvider(name)
		if err != nil {
			t.Fatalf("GetProvider(%q) returned error: %v", name, err)
		}
		if p.Name() != name {
			t.Errorf("GetProvider(%q).Name() = %q, want %q", name, p.Name(), name)
		}
	}
}

func TestGetProvider_Unknown(t *testing.T) {
	_, err := GetProvider("nonexistent-provider")
	if err == nil {
		t.Fatal("GetProvider(\"nonexistent-provider\") should return error")
	}
	if !strings.Contains(err.Error(), "unknown provider") {
		t.Errorf("error should mention 'unknown provider', got: %v", err)
	}
}

func TestAvailableProviders(t *testing.T) {
	avail := AvailableProviders()
	if !strings.Contains(avail, "mock") {
		t.Errorf("AvailableProviders() should include 'mock', got: %s", avail)
	}
	if !strings.Contains(avail, "leetcode") {
		t.Errorf("AvailableProviders() should include 'leetcode', got: %s", avail)
	}
}
