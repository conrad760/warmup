package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

// TopicTag represents a topic/category tag from the leetgo database.
type TopicTag struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// CodeSnippet represents a code template for a specific language.
type CodeSnippet struct {
	LangSlug string `json:"langSlug"`
	Lang     string `json:"lang"`
	Code     string `json:"code"`
}

// findLeetgoWorkspace searches for a directory containing leetgo.yaml.
func findLeetgoWorkspace() string {
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		return ""
	}

	dir, err := os.Getwd()
	if err == nil {
		for {
			if _, err := os.Stat(filepath.Join(dir, "leetgo.yaml")); err == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		parent := filepath.Dir(cwd)
		if entries, err := os.ReadDir(parent); err == nil {
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				sibling := filepath.Join(parent, e.Name())
				if _, err := os.Stat(filepath.Join(sibling, "leetgo.yaml")); err == nil {
					return sibling
				}
				if subs, err := os.ReadDir(sibling); err == nil {
					for _, sub := range subs {
						if !sub.IsDir() {
							continue
						}
						nested := filepath.Join(sibling, sub.Name())
						if _, err := os.Stat(filepath.Join(nested, "leetgo.yaml")); err == nil {
							return nested
						}
					}
				}
			}
		}
	}

	bases := []string{
		filepath.Join(homeDir, "leetcode"),
		filepath.Join(homeDir, "code", "leetcode"),
		filepath.Join(homeDir, "projects", "leetcode"),
	}
	for _, base := range bases {
		if _, err := os.Stat(filepath.Join(base, "leetgo.yaml")); err == nil {
			return base
		}
		ws := filepath.Join(base, "workspace")
		if _, err := os.Stat(filepath.Join(ws, "leetgo.yaml")); err == nil {
			return ws
		}
	}

	return ""
}

// findLeetgoDatabase locates the leetgo cache database on disk.
func findLeetgoDatabase() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	candidates := []string{
		filepath.Join(homeDir, ".config", "leetgo", "cache", "leetcode-questions-full.db"),
		filepath.Join(homeDir, ".leetgo", "cache", "leetcode-questions-full.db"),
		filepath.Join(homeDir, "Library", "Application Support", "leetgo", "cache", "leetcode-questions-full.db"),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("leetgo database not found — run 'leetgo cache update' first")
}

// loadQuestions loads curated questions by pulling descriptions from the leetgo database.
func loadQuestions(dbPath string, curated []CuratedQuestion) ([]Question, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var questions []Question
	var failed []string

	for _, cq := range curated {
		var titleSlug, title, difficulty, topicTags, content string
		err := db.QueryRow(`
			SELECT titleSlug, title, difficulty, topicTags, content
			FROM questions
			WHERE titleSlug = ?
			LIMIT 1
		`, cq.Slug).Scan(&titleSlug, &title, &difficulty, &topicTags, &content)
		if err != nil {
			failed = append(failed, cq.Slug)
			continue
		}

		var diff Difficulty
		switch strings.ToLower(difficulty) {
		case "easy":
			diff = Easy
		case "medium":
			diff = Medium
		case "hard":
			diff = Hard
		}

		desc, example := parseContent(content)

		cat := cq.Category
		if cat == "" {
			cat = parsePrimaryCategory(topicTags)
		}

		q := Question{
			Title:        title,
			Difficulty:   diff,
			Category:     cat,
			Description:  desc,
			Example:      example,
			Options:      cq.Options,
			Solution:     cq.Solution,
			LeetcodeSlug: titleSlug,
		}
		questions = append(questions, q)
	}

	if len(failed) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: %d questions not found in leetgo DB: %s\n",
			len(failed), strings.Join(failed, ", "))
	}

	if len(questions) == 0 {
		return nil, fmt.Errorf("no questions could be loaded from leetgo database")
	}

	return questions, nil
}

// loadQuestionsFromJSON loads additional curated questions from a JSON file.
func loadQuestionsFromJSON(dbPath, jsonPath string) ([]Question, error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var curated []CuratedQuestion
	if err := json.Unmarshal(data, &curated); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	for i, cq := range curated {
		if cq.Slug == "" {
			return nil, fmt.Errorf("question %d: missing slug", i)
		}
		if len(cq.Options) < 2 {
			return nil, fmt.Errorf("question %d (%s): needs at least 2 options", i, cq.Slug)
		}
		hasOptimal := false
		for _, opt := range cq.Options {
			if opt.Rating == Optimal {
				hasOptimal = true
				break
			}
		}
		if !hasOptimal {
			return nil, fmt.Errorf("question %d (%s): needs at least one OPTIMAL option", i, cq.Slug)
		}
	}

	return loadQuestions(dbPath, curated)
}

var (
	reBlockClose = regexp.MustCompile(`</(p|div|li|pre|h\d)>`)
	reBlockOpen  = regexp.MustCompile(`<(p|div|li|h\d)[^>]*>`)
	rePreOpen    = regexp.MustCompile(`<pre[^>]*>`)
	reTags       = regexp.MustCompile(`<[^>]+>`)
	reSup        = regexp.MustCompile(`<sup>([^<]+)</sup>`)
)

func parseContent(html string) (description, example string) {
	html = reSup.ReplaceAllString(html, "^$1")

	parts := regexp.MustCompile(`(?i)<strong[^>]*>\s*Example\s*\d*\s*:?\s*</strong>`).Split(html, 3)

	descHTML := parts[0]
	description = htmlToText(descHTML)

	if len(parts) >= 2 {
		exHTML := parts[1]
		// Trim at next Example or Constraints
		if idx := regexp.MustCompile(`(?i)<strong[^>]*>\s*(Example|Constraint)`).FindStringIndex(exHTML); idx != nil {
			exHTML = exHTML[:idx[0]]
		}
		example = htmlToText(exHTML)
	}

	if idx := strings.Index(strings.ToLower(description), "constraints"); idx > 0 {
		description = strings.TrimSpace(description[:idx])
	}

	if idx := strings.Index(strings.ToLower(description), "follow-up"); idx > 0 {
		description = strings.TrimSpace(description[:idx])
	}

	return description, example
}

func htmlToText(rawHTML string) string {
	text := rawHTML

	text = rePreOpen.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, "</pre>", "\n")

	text = reBlockClose.ReplaceAllString(text, "\n")
	text = reBlockOpen.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "<ul>", "\n")
	text = strings.ReplaceAll(text, "</ul>", "")
	text = strings.ReplaceAll(text, "<ol>", "\n")
	text = strings.ReplaceAll(text, "</ol>", "")

	text = reTags.ReplaceAllString(text, "")

	text = html.UnescapeString(text)
	text = strings.ReplaceAll(text, "\u00a0", " ")

	lines := strings.Split(text, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}

// parsePrimaryCategory extracts the first topic tag as the category.
// This is used as a fallback when the curated question doesn't specify its own category.
func parsePrimaryCategory(topicTagsJSON string) string {
	var tags []TopicTag
	if err := json.Unmarshal([]byte(topicTagsJSON), &tags); err != nil {
		return "Unknown"
	}
	if len(tags) > 0 {
		return tags[0].Name
	}
	return "Unknown"
}
