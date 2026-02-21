package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	leetcodeGraphQLURL = "https://leetcode.com/graphql"
	leetcodeUserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"
)

// LeetCodeProvider fetches problem data from LeetCode's public GraphQL API.
// No authentication is needed for fetching problem descriptions — only for test/submit (Phase 2).
type LeetCodeProvider struct {
	client *http.Client
}

func init() {
	RegisterProvider("leetcode", func() Provider {
		return &LeetCodeProvider{
			client: &http.Client{Timeout: 30 * time.Second},
		}
	})
}

func (lc *LeetCodeProvider) Name() string { return "leetcode" }

// FetchProblem fetches a single problem from LeetCode's public GraphQL API.
// The lang parameter selects which code snippet to return (e.g. "go", "python3").
func (lc *LeetCodeProvider) FetchProblem(id string, lang string) (*ProblemData, error) {
	query := `query questionData($titleSlug: String!) {
		question(titleSlug: $titleSlug) {
			questionId
			questionFrontendId
			title
			titleSlug
			content
			difficulty
			topicTags { name slug }
			codeSnippets { lang langSlug code }
			sampleTestCase
			exampleTestcases
			metaData
		}
	}`

	payload := graphqlPayload{
		Query:         query,
		OperationName: "questionData",
		Variables:     map[string]string{"titleSlug": id},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal graphql request: %w", err)
	}

	req, err := http.NewRequest("POST", leetcodeGraphQLURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", leetcodeUserAgent)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch %q: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("fetch %q: HTTP %d: %s", id, resp.StatusCode, string(respBody))
	}

	var result graphqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response for %q: %w", id, err)
	}

	if result.Data.Question.TitleSlug == "" {
		return nil, fmt.Errorf("problem %q not found on LeetCode", id)
	}

	q := result.Data.Question

	// Find the code snippet for the requested language.
	langSlug := normalizeLangSlug(lang)
	var codeSnippet string
	for _, cs := range q.CodeSnippets {
		if cs.LangSlug == langSlug {
			codeSnippet = cs.Code
			break
		}
	}

	// Parse HTML content into plain text description + example.
	desc, example := lcParseContent(q.Content)

	// Extract tag names.
	tags := make([]string, len(q.TopicTags))
	for i, t := range q.TopicTags {
		tags[i] = t.Name
	}

	// Parse function metadata for test harness generation.
	meta := parseLCMetaData(q.MetaData)

	return &ProblemData{
		ID:          q.TitleSlug,
		Title:       q.Title,
		Description: desc,
		Examples:    example,
		Difficulty:  q.Difficulty,
		Tags:        tags,
		CodeSnippet: codeSnippet,
		TestInput:   q.ExampleTestcases,
		Meta:        meta,
	}, nil
}

// --- MetaData parsing ---

// lcRawMeta matches LeetCode's metaData JSON field structure.
type lcRawMeta struct {
	Name         string       `json:"name"`
	Params       []lcRawParam `json:"params"`
	Return       *lcRawParam  `json:"return,omitempty"`
	SystemDesign bool         `json:"systemdesign,omitempty"`
}

type lcRawParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// parseLCMetaData parses the metaData JSON string from the LeetCode API.
// Returns nil if the metadata is missing or unparseable (graceful degradation).
func parseLCMetaData(raw string) *FuncMeta {
	if raw == "" {
		return nil
	}

	var rm lcRawMeta
	if err := json.Unmarshal([]byte(raw), &rm); err != nil {
		return nil
	}

	if rm.Name == "" {
		return nil
	}

	meta := &FuncMeta{
		Name:         rm.Name,
		SystemDesign: rm.SystemDesign,
	}

	for _, p := range rm.Params {
		meta.Params = append(meta.Params, ParamMeta{Name: p.Name, Type: p.Type})
	}

	if rm.Return != nil && rm.Return.Type != "" {
		meta.Return = &ParamMeta{Name: rm.Return.Name, Type: rm.Return.Type}
	}

	return meta
}

// --- GraphQL types ---

type graphqlPayload struct {
	Query         string            `json:"query"`
	OperationName string            `json:"operationName,omitempty"`
	Variables     map[string]string `json:"variables"`
}

type graphqlResponse struct {
	Data struct {
		Question lcQuestion `json:"question"`
	} `json:"data"`
}

type lcQuestion struct {
	QuestionID         string          `json:"questionId"`
	QuestionFrontendID string          `json:"questionFrontendId"`
	Title              string          `json:"title"`
	TitleSlug          string          `json:"titleSlug"`
	Content            string          `json:"content"`
	Difficulty         string          `json:"difficulty"`
	TopicTags          []lcTopicTag    `json:"topicTags"`
	CodeSnippets       []lcCodeSnippet `json:"codeSnippets"`
	SampleTestCase     string          `json:"sampleTestCase"`
	ExampleTestcases   string          `json:"exampleTestcases"`
	MetaData           string          `json:"metaData"`
}

type lcTopicTag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type lcCodeSnippet struct {
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
}

// normalizeLangSlug maps common language names to LeetCode's langSlug values.
func normalizeLangSlug(lang string) string {
	switch strings.ToLower(lang) {
	case "go", "golang":
		return "golang"
	case "python", "python3", "py":
		return "python3"
	case "java":
		return "java"
	case "cpp", "c++":
		return "cpp"
	case "c":
		return "c"
	case "js", "javascript":
		return "javascript"
	case "ts", "typescript":
		return "typescript"
	case "rust", "rs":
		return "rust"
	default:
		return strings.ToLower(lang)
	}
}

// --- HTML to text conversion ---

var (
	lcReBlockClose = regexp.MustCompile(`</(p|div|li|pre|h\d)>`)
	lcReBlockOpen  = regexp.MustCompile(`<(p|div|li|h\d)[^>]*>`)
	lcRePreOpen    = regexp.MustCompile(`<pre[^>]*>`)
	lcReTags       = regexp.MustCompile(`<[^>]+>`)
	lcReSup        = regexp.MustCompile(`<sup>([^<]+)</sup>`)
	lcReExample    = regexp.MustCompile(`(?i)<strong[^>]*>\s*Example\s*\d*\s*:?\s*</strong>`)
	lcReNextSect   = regexp.MustCompile(`(?i)<strong[^>]*>\s*(Example|Constraint)`)
)

func lcParseContent(rawHTML string) (description, example string) {
	rawHTML = lcReSup.ReplaceAllString(rawHTML, "^$1")

	parts := lcReExample.Split(rawHTML, 3)

	descHTML := parts[0]
	description = lcHTMLToText(descHTML)

	if len(parts) >= 2 {
		exHTML := parts[1]
		if idx := lcReNextSect.FindStringIndex(exHTML); idx != nil {
			exHTML = exHTML[:idx[0]]
		}
		example = lcHTMLToText(exHTML)
	}

	if idx := strings.Index(strings.ToLower(description), "constraints"); idx > 0 {
		description = strings.TrimSpace(description[:idx])
	}
	if idx := strings.Index(strings.ToLower(description), "follow-up"); idx > 0 {
		description = strings.TrimSpace(description[:idx])
	}

	return description, example
}

func lcHTMLToText(rawHTML string) string {
	text := rawHTML

	text = lcRePreOpen.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, "</pre>", "\n")

	text = lcReBlockClose.ReplaceAllString(text, "\n")
	text = lcReBlockOpen.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "<ul>", "\n")
	text = strings.ReplaceAll(text, "</ul>", "")
	text = strings.ReplaceAll(text, "<ol>", "\n")
	text = strings.ReplaceAll(text, "</ol>", "")

	text = lcReTags.ReplaceAllString(text, "")

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
