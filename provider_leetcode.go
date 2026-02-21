package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	leetcodeBaseURL   = "https://leetcode.com"
	leetcodeUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36"

	defaultPollInterval = 1 * time.Second
	defaultPollTimeout  = 30 * time.Second
)

// LeetCodeProvider fetches problem data from LeetCode's public GraphQL API.
// Authentication (LEETCODE_SESSION + LEETCODE_CSRFTOKEN env vars) is only
// needed for test/submit — fetching problems is unauthenticated.
type LeetCodeProvider struct {
	client       *http.Client
	baseURL      string        // defaults to "https://leetcode.com"
	session      string        // LEETCODE_SESSION cookie value
	csrfToken    string        // LEETCODE_CSRFTOKEN (also sent as X-CSRFToken header)
	pollInterval time.Duration // time between check polls (default 1s)
	pollTimeout  time.Duration // max time to wait for result (default 30s)

	// questionIDs caches slug → numeric question ID to avoid redundant GraphQL lookups.
	// Populated during FetchProblem and used by RunTests/Submit.
	questionIDs map[string]string
}

func init() {
	RegisterProvider("leetcode", func() Provider {
		return NewLeetCodeProvider()
	})
}

// NewLeetCodeProvider creates a provider with production defaults.
func NewLeetCodeProvider() *LeetCodeProvider {
	return &LeetCodeProvider{
		client:       &http.Client{Timeout: 30 * time.Second},
		baseURL:      leetcodeBaseURL,
		pollInterval: defaultPollInterval,
		pollTimeout:  defaultPollTimeout,
		questionIDs:  make(map[string]string),
	}
}

// URL builders — all endpoints derived from baseURL.
func (lc *LeetCodeProvider) graphqlURL() string {
	return lc.baseURL + "/graphql"
}
func (lc *LeetCodeProvider) interpretURL(slug string) string {
	return lc.baseURL + "/problems/" + slug + "/interpret_solution/"
}
func (lc *LeetCodeProvider) submitURL(slug string) string {
	return lc.baseURL + "/problems/" + slug + "/submit/"
}
func (lc *LeetCodeProvider) checkURL(id string) string {
	return lc.baseURL + "/submissions/detail/" + id + "/check/"
}

func (lc *LeetCodeProvider) Name() string { return "leetcode" }

// PollInterval implements PollConfig.
func (lc *LeetCodeProvider) PollInterval() time.Duration { return lc.pollInterval }

// PollTimeout implements PollConfig.
func (lc *LeetCodeProvider) PollTimeout() time.Duration { return lc.pollTimeout }

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

	req, err := http.NewRequest("POST", lc.graphqlURL(), bytes.NewReader(body))
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

	// Cache the numeric question ID for later test/submit calls.
	if q.QuestionID != "" {
		if lc.questionIDs == nil {
			lc.questionIDs = make(map[string]string)
		}
		lc.questionIDs[q.TitleSlug] = q.QuestionID
	}

	return &ProblemData{
		ID:          q.TitleSlug,
		QuestionID:  q.QuestionID,
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

// --- Authenticator interface ---

const (
	envLeetCodeSession   = "LEETCODE_SESSION"
	envLeetCodeCSRFToken = "LEETCODE_CSRFTOKEN"
)

// Authenticate loads credentials from environment variables.
// On failure, both fields are cleared to avoid a half-authenticated state
// where session is set but csrfToken is empty (or vice versa).
func (lc *LeetCodeProvider) Authenticate() error {
	session := os.Getenv(envLeetCodeSession)
	csrfToken := os.Getenv(envLeetCodeCSRFToken)
	if session == "" || csrfToken == "" {
		lc.session = ""
		lc.csrfToken = ""
		return fmt.Errorf("missing credentials: set %s and %s env vars", envLeetCodeSession, envLeetCodeCSRFToken)
	}
	lc.session = session
	lc.csrfToken = csrfToken
	return nil
}

// IsAuthenticated returns whether valid credentials are available.
func (lc *LeetCodeProvider) IsAuthenticated() bool {
	return lc.session != "" && lc.csrfToken != ""
}

// AuthHelp returns instructions for setting up LeetCode credentials.
func (lc *LeetCodeProvider) AuthHelp() string {
	return `LeetCode authentication required for test/submit.

1. Log in to leetcode.com in your browser
2. Open DevTools > Application > Cookies > leetcode.com
3. Copy these two cookie values:
   - LEETCODE_SESSION  (the long one)
   - csrftoken

4. Set them as environment variables:
   export LEETCODE_SESSION="<value>"
   export LEETCODE_CSRFTOKEN="<value>"

   Or add them to your shell profile / .envrc file.
`
}

// setAuthHeaders adds authentication cookies and CSRF header to a request.
func (lc *LeetCodeProvider) setAuthHeaders(req *http.Request) {
	req.Header.Set("Cookie", fmt.Sprintf("LEETCODE_SESSION=%s; csrftoken=%s", lc.session, lc.csrfToken))
	req.Header.Set("X-CSRFToken", lc.csrfToken)
	req.Header.Set("Referer", lc.baseURL)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", leetcodeUserAgent)
}

// --- Tester interface ---

// RunTests submits code for testing against provided test cases.
// Returns an interpret_id for polling.
func (lc *LeetCodeProvider) RunTests(id string, lang string, code string, input string) (string, error) {
	if !lc.IsAuthenticated() {
		return "", fmt.Errorf("not authenticated: %s", lc.AuthHelp())
	}

	// We need the numeric question_id. Look it up via a quick GraphQL query
	// if we don't have it cached. For now, fetch it from the problem data.
	qID, err := lc.fetchQuestionID(id)
	if err != nil {
		return "", fmt.Errorf("fetch question ID for %q: %w", id, err)
	}

	payload := map[string]string{
		"lang":        normalizeLangSlug(lang),
		"question_id": qID,
		"typed_code":  code,
		"data_input":  input,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal test payload: %w", err)
	}

	req, err := http.NewRequest("POST", lc.interpretURL(id), bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create test request: %w", err)
	}
	lc.setAuthHeaders(req)

	resp, err := lc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("test request for %q: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("test %q: HTTP %d: %s", id, resp.StatusCode, string(respBody))
	}

	var result struct {
		InterpretID string `json:"interpret_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode test response: %w", err)
	}

	if result.InterpretID == "" {
		return "", fmt.Errorf("test %q: empty interpret_id (auth may have expired)", id)
	}

	return result.InterpretID, nil
}

// CheckTestResult polls for the result of a test run.
// Returns (result, done, error). When done is false, the caller should poll again.
func (lc *LeetCodeProvider) CheckTestResult(runID string) (*TestResult, bool, error) {
	req, err := http.NewRequest("GET", lc.checkURL(runID), nil)
	if err != nil {
		return nil, false, fmt.Errorf("create check request: %w", err)
	}
	lc.setAuthHeaders(req)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("check: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var raw lcCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, false, fmt.Errorf("decode check response: %w", err)
	}

	if raw.State == "PENDING" || raw.State == "STARTED" {
		return nil, false, nil
	}

	return raw.toTestResult(), true, nil
}

// --- Submitter interface ---

// Submit submits code for grading. Returns a submission ID for polling.
func (lc *LeetCodeProvider) Submit(id string, lang string, code string) (string, error) {
	if !lc.IsAuthenticated() {
		return "", fmt.Errorf("not authenticated: %s", lc.AuthHelp())
	}

	qID, err := lc.fetchQuestionID(id)
	if err != nil {
		return "", fmt.Errorf("fetch question ID for %q: %w", id, err)
	}

	payload := map[string]string{
		"lang":        normalizeLangSlug(lang),
		"question_id": qID,
		"typed_code":  code,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal submit payload: %w", err)
	}

	req, err := http.NewRequest("POST", lc.submitURL(id), bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create submit request: %w", err)
	}
	lc.setAuthHeaders(req)

	resp, err := lc.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("submit request for %q: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("submit %q: HTTP %d: %s", id, resp.StatusCode, string(respBody))
	}

	var result struct {
		SubmissionID json.Number `json:"submission_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode submit response: %w", err)
	}

	subID := result.SubmissionID.String()
	if subID == "" || subID == "0" {
		return "", fmt.Errorf("submit %q: empty submission_id (auth may have expired)", id)
	}

	return subID, nil
}

// CheckSubmission polls for the result of a submission.
func (lc *LeetCodeProvider) CheckSubmission(subID string) (*SubmitResult, bool, error) {
	req, err := http.NewRequest("GET", lc.checkURL(subID), nil)
	if err != nil {
		return nil, false, fmt.Errorf("create check request: %w", err)
	}
	lc.setAuthHeaders(req)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("check request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("check: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var raw lcCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, false, fmt.Errorf("decode check response: %w", err)
	}

	if raw.State == "PENDING" || raw.State == "STARTED" {
		return nil, false, nil
	}

	return raw.toSubmitResult(), true, nil
}

// --- Shared polling + response types ---

// flexStrings handles LeetCode's inconsistent JSON where a field may be
// either a string or an array of strings depending on the response.
type flexStrings []string

func (f *flexStrings) UnmarshalJSON(data []byte) error {
	// Try array first.
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*f = arr
		return nil
	}
	// Try single string.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if s != "" {
			*f = []string{s}
		} else {
			*f = nil
		}
		return nil
	}
	// Graceful fallback — don't fail the whole response.
	*f = nil
	return nil
}

// lcCheckResponse is the unified response from LeetCode's /check/ endpoint.
// Used for both test and submit results.
type lcCheckResponse struct {
	State              string      `json:"state"`
	StatusCode         int         `json:"status_code"`        // 10=Accepted, 11=Wrong, 12=MLE, 13=OLE, 14=TLE, 15=RE, 20=CE
	StatusMsg          string      `json:"status_msg"`         // "Accepted", "Wrong Answer", etc.
	StatusRuntime      string      `json:"status_runtime"`     // e.g. "5 ms"
	StatusMemory       string      `json:"status_memory"`      // e.g. "3.2 MB"
	RuntimePercentile  json.Number `json:"runtime_percentile"` // e.g. 95.2
	MemoryPercentile   json.Number `json:"memory_percentile"`
	CodeAnswer         flexStrings `json:"code_answer"`          // test: actual outputs
	ExpectedCodeAnswer flexStrings `json:"expected_code_answer"` // test: expected outputs
	CodeOutput         flexStrings `json:"code_output"`          // stdout lines
	StdOutput          string      `json:"std_output"`           // raw stdout
	TotalCorrect       *int        `json:"total_correct"`        // submit: passed count
	TotalTestcases     *int        `json:"total_testcases"`      // submit: total count
	CompileError       string      `json:"compile_error"`
	FullCompileError   string      `json:"full_compile_error"`
	RuntimeError       string      `json:"runtime_error"`
	FullRuntimeError   string      `json:"full_runtime_error"`
	LastTestcase       string      `json:"last_testcase"` // failing test input
	ExpectedOutput     string      `json:"expected_output"`
	InputFormatted     string      `json:"input_formatted"`
	Input              string      `json:"input"` // raw test input
}

func (r *lcCheckResponse) toTestResult() *TestResult {
	tr := &TestResult{
		RuntimeMs:    parseRuntimeMs(r.StatusRuntime),
		CompileError: firstNonEmpty(r.FullCompileError, r.CompileError),
		RuntimeError: firstNonEmpty(r.FullRuntimeError, r.RuntimeError),
	}

	// Build raw output for display.
	var out strings.Builder
	if tr.CompileError != "" {
		out.WriteString("Compile Error:\n" + tr.CompileError)
	} else if tr.RuntimeError != "" {
		out.WriteString("Runtime Error:\n" + tr.RuntimeError)
	} else if r.StatusCode != 10 {
		// Non-success status (TLE, MLE, etc.) — not a correctness check.
		tr.Passed = false
		out.WriteString(fmt.Sprintf("%s\n", r.StatusMsg))
		if r.LastTestcase != "" {
			tr.Input = r.LastTestcase
			out.WriteString(fmt.Sprintf("Input:    %s\n", tr.Input))
		}
	} else {
		// Status 10 means "executed without error" for interpret_solution.
		// Actual correctness must be determined by comparing outputs.
		// LeetCode sometimes returns trailing empty strings — strip them.
		answers := trimTrailingEmpty(r.CodeAnswer)
		expected := trimTrailingEmpty(r.ExpectedCodeAnswer)

		if len(answers) > 0 {
			tr.Actual = strings.Join(answers, "\n")
		}
		if len(expected) > 0 {
			tr.Expected = strings.Join(expected, "\n")
		}
		if r.LastTestcase != "" {
			tr.Input = r.LastTestcase
		}

		tr.Passed = flexStringsEqual(answers, expected)

		if tr.Passed {
			out.WriteString(fmt.Sprintf("Accepted (%s)\n", r.StatusRuntime))
		} else {
			out.WriteString("Wrong Answer\n")
		}

		// Show per-test-case results when there are multiple.
		if len(answers) > 1 || len(expected) > 1 {
			n := len(answers)
			if len(expected) > n {
				n = len(expected)
			}
			for i := 0; i < n; i++ {
				actual := ""
				exp := ""
				if i < len(answers) {
					actual = answers[i]
				}
				if i < len(expected) {
					exp = expected[i]
				}
				mark := "PASS"
				if actual != exp {
					mark = "FAIL"
				}
				out.WriteString(fmt.Sprintf("  Case %d: %s  output=%s  expected=%s\n", i+1, mark, actual, exp))
			}
		} else {
			out.WriteString(fmt.Sprintf("Output:   %s\n", tr.Actual))
			out.WriteString(fmt.Sprintf("Expected: %s\n", tr.Expected))
		}
	}

	if r.StdOutput != "" {
		out.WriteString(fmt.Sprintf("\nStdout:\n%s", r.StdOutput))
	}
	tr.RawOutput = out.String()
	return tr
}

// flexStringsEqual compares two flexStrings for element-wise equality.
func flexStringsEqual(a, b flexStrings) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// trimTrailingEmpty removes trailing empty strings from a slice.
// LeetCode's API sometimes returns ["4", "-1", ""] — the trailing empty
// element produces a phantom "Case 3" with blank output/expected.
func trimTrailingEmpty(ss []string) []string {
	i := len(ss)
	for i > 0 && ss[i-1] == "" {
		i--
	}
	return ss[:i]
}

func (r *lcCheckResponse) toSubmitResult() *SubmitResult {
	sr := &SubmitResult{
		Accepted:     r.StatusCode == 10,
		StatusMsg:    r.StatusMsg,
		RuntimeMs:    parseRuntimeMs(r.StatusRuntime),
		CompileError: firstNonEmpty(r.FullCompileError, r.CompileError),
		RuntimeError: firstNonEmpty(r.FullRuntimeError, r.RuntimeError),
	}

	if r.TotalCorrect != nil {
		sr.PassedCases = *r.TotalCorrect
	}
	if r.TotalTestcases != nil {
		sr.TotalCases = *r.TotalTestcases
	}

	rPct, _ := r.RuntimePercentile.Float64()
	if rPct > 0 {
		sr.RuntimePct = fmt.Sprintf("faster than %.1f%%", rPct)
	}
	mPct, _ := r.MemoryPercentile.Float64()
	if mPct > 0 {
		sr.MemoryPct = fmt.Sprintf("less than %.1f%%", mPct)
	}

	// Build raw output for display.
	var out strings.Builder
	if sr.CompileError != "" {
		out.WriteString("Compile Error:\n" + sr.CompileError)
	} else if sr.RuntimeError != "" {
		out.WriteString("Runtime Error:\n" + sr.RuntimeError)
	} else if sr.Accepted {
		out.WriteString("Accepted\n")
		if sr.RuntimePct != "" {
			out.WriteString(fmt.Sprintf("Runtime: %s (%s)\n", r.StatusRuntime, sr.RuntimePct))
		} else {
			out.WriteString(fmt.Sprintf("Runtime: %s\n", r.StatusRuntime))
		}
		if sr.MemoryPct != "" {
			out.WriteString(fmt.Sprintf("Memory:  %s (%s)\n", r.StatusMemory, sr.MemoryPct))
		} else {
			out.WriteString(fmt.Sprintf("Memory:  %s\n", r.StatusMemory))
		}
		out.WriteString(fmt.Sprintf("Cases:   %d/%d passed\n", sr.PassedCases, sr.TotalCases))
	} else {
		out.WriteString(fmt.Sprintf("%s\n", sr.StatusMsg))
		out.WriteString(fmt.Sprintf("Cases: %d/%d passed\n", sr.PassedCases, sr.TotalCases))
		if r.LastTestcase != "" {
			out.WriteString(fmt.Sprintf("Failing input:\n%s\n", r.LastTestcase))
		}
		if r.ExpectedOutput != "" {
			out.WriteString(fmt.Sprintf("Expected: %s\n", r.ExpectedOutput))
		}
	}
	sr.RawOutput = out.String()
	return sr
}

// fetchQuestionID returns the numeric question ID for a given slug.
// Checks the in-memory cache first (populated by FetchProblem), falling back
// to a lightweight GraphQL query. This avoids redundant network calls when
// the problem was already fetched during startup.
func (lc *LeetCodeProvider) fetchQuestionID(slug string) (string, error) {
	// Check cache first.
	if lc.questionIDs != nil {
		if qID, ok := lc.questionIDs[slug]; ok {
			return qID, nil
		}
	}

	query := `query questionId($titleSlug: String!) {
		question(titleSlug: $titleSlug) {
			questionId
		}
	}`

	payload := graphqlPayload{
		Query:     query,
		Variables: map[string]string{"titleSlug": slug},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", lc.graphqlURL(), bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", leetcodeUserAgent)

	resp, err := lc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Question struct {
				QuestionID string `json:"questionId"`
			} `json:"question"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Data.Question.QuestionID == "" {
		return "", fmt.Errorf("question %q not found", slug)
	}

	// Cache for future calls.
	if lc.questionIDs == nil {
		lc.questionIDs = make(map[string]string)
	}
	lc.questionIDs[slug] = result.Data.Question.QuestionID

	return result.Data.Question.QuestionID, nil
}

// parseRuntimeMs extracts milliseconds from LeetCode's "X ms" string.
func parseRuntimeMs(s string) int {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, " ms")
	var ms int
	fmt.Sscanf(s, "%d", &ms)
	return ms
}

// firstNonEmpty returns the first non-empty string.
func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
