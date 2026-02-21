package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// loadQuestions loads questions from the curated bank using the provider system.
// It groups questions by provider, fetches problem data (via cache), and merges
// with the curated approach options and solutions.
func loadQuestionsFromProviders(curated []CuratedQuestion, cache *QuestionCache, lang string) ([]Question, error) {
	// Group curated questions by provider.
	byProvider := make(map[string][]CuratedQuestion)
	for _, cq := range curated {
		pName := cq.Provider
		if pName == "" {
			pName = DefaultProviderName
		}
		byProvider[pName] = append(byProvider[pName], cq)
	}

	// Initialize each needed provider.
	providerInstances := make(map[string]Provider)
	for pName := range byProvider {
		p, err := GetProvider(pName)
		if err != nil {
			return nil, fmt.Errorf("provider %q: %w", pName, err)
		}
		providerInstances[pName] = p
	}

	// Fetch problem data and build questions.
	var questions []Question
	var failed []string

	for pName, cqs := range byProvider {
		p := providerInstances[pName]
		for _, cq := range cqs {
			problem, err := FetchWithCache(cache, p, cq.ProblemID, lang)
			if err != nil {
				failed = append(failed, fmt.Sprintf("%s/%s", pName, cq.ProblemID))
				continue
			}

			diff := parseDifficulty(problem.Difficulty)

			// Use curated category if set, otherwise pick from provider tags.
			cat := cq.Category
			if cat == "" && len(problem.Tags) > 0 {
				cat = problem.Tags[0]
			}
			if cat == "" {
				cat = "Unknown"
			}

			q := Question{
				Title:       problem.Title,
				Difficulty:  diff,
				Category:    cat,
				Description: problem.Description,
				Example:     problem.Examples,
				Options:     cq.Options,
				Solution:    cq.Solution,
				Provider:    pName,
				ProblemID:   cq.ProblemID,
				CodeSnippet: problem.CodeSnippet,
				TestInput:   problem.TestInput,
				Meta:        problem.Meta,
			}
			questions = append(questions, q)
		}
	}

	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: %d questions could not be loaded: %s\n",
			len(failed), strings.Join(failed, ", "))
	}

	if len(questions) == 0 {
		return nil, fmt.Errorf("no questions could be loaded from any provider")
	}

	return questions, nil
}

// loadQuestionsFromJSON loads additional curated questions from a JSON file,
// then fetches their problem data via providers.
func loadQuestionsFromJSONFile(jsonPath string, cache *QuestionCache, lang string) ([]Question, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var curated []CuratedQuestion
	if err := json.Unmarshal(data, &curated); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Validate.
	for i, cq := range curated {
		if cq.ProblemID == "" {
			return nil, fmt.Errorf("question %d: missing ProblemID", i)
		}
		if len(cq.Options) < 2 {
			return nil, fmt.Errorf("question %d (%s): needs at least 2 options", i, cq.ProblemID)
		}
		hasOptimal := false
		for _, opt := range cq.Options {
			if opt.Rating == Optimal {
				hasOptimal = true
				break
			}
		}
		if !hasOptimal {
			return nil, fmt.Errorf("question %d (%s): needs at least one OPTIMAL option", i, cq.ProblemID)
		}
	}

	return loadQuestionsFromProviders(curated, cache, lang)
}

// parseDifficulty converts a string difficulty to the Difficulty type.
func parseDifficulty(s string) Difficulty {
	switch strings.ToLower(s) {
	case "easy":
		return Easy
	case "medium":
		return Medium
	case "hard":
		return Hard
	default:
		return Easy
	}
}
