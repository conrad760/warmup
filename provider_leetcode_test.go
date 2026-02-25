package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNormalizeLangSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"go", "golang"},
		{"golang", "golang"},
		{"Go", "golang"},
		{"python", "python3"},
		{"python3", "python3"},
		{"py", "python3"},
		{"java", "java"},
		{"cpp", "cpp"},
		{"c++", "cpp"},
		{"c", "c"},
		{"js", "javascript"},
		{"javascript", "javascript"},
		{"ts", "typescript"},
		{"typescript", "typescript"},
		{"rust", "rust"},
		{"rs", "rust"},
		{"ruby", "ruby"}, // unmapped, lowercased as-is
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeLangSlug(tt.input)
			if got != tt.want {
				t.Errorf("normalizeLangSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLcParseContent_BasicHTML(t *testing.T) {
	html := `<p>Given an array of integers <code>nums</code> and an integer <code>target</code>, return indices of the two numbers.</p>
<p><strong>Example 1:</strong></p>
<pre>
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
</pre>
<p><strong>Constraints:</strong></p>
<ul><li>2 &lt;= nums.length</li></ul>`

	desc, example, constraints := lcParseContent(html)

	if !strings.Contains(desc, "array of integers") {
		t.Errorf("description should contain problem text, got: %s", desc)
	}
	if strings.Contains(desc, "<") {
		t.Errorf("description should not contain HTML tags, got: %s", desc)
	}
	// Constraints should be stripped from description.
	if strings.Contains(strings.ToLower(desc), "constraint") {
		t.Errorf("description should not contain constraints, got: %s", desc)
	}
	// Constraints should be extracted.
	if !strings.Contains(constraints, "nums.length") {
		t.Errorf("constraints should contain constraint text, got: %s", constraints)
	}
	// Example is in the second split part — may or may not be populated
	// depending on how the <strong>Example tags are structured.
	_ = example
}

func TestLcParseContent_WithExampleSections(t *testing.T) {
	html := `<p>Description text here.</p>
<p><strong>Example 1:</strong></p>
<pre>Input: x = 5
Output: 25</pre>
<p><strong>Example 2:</strong></p>
<pre>Input: x = -3
Output: 9</pre>
<p><strong>Constraints:</strong></p>
<ul><li>-100 &lt;= x &lt;= 100</li></ul>`

	desc, example, _ := lcParseContent(html)

	if !strings.Contains(desc, "Description text here") {
		t.Errorf("description missing main text, got: %s", desc)
	}
	if example == "" {
		t.Error("example should not be empty when Example section exists")
	}
	if !strings.Contains(example, "x = 5") {
		t.Errorf("example should contain first example, got: %s", example)
	}
	if !strings.Contains(example, "x = -3") {
		t.Errorf("example should contain second example, got: %s", example)
	}
	if strings.Contains(example, "Constraint") {
		t.Error("example should not include constraints section")
	}
}

func TestLcHTMLToText_UnescapesEntities(t *testing.T) {
	got := lcHTMLToText(`<p>a &amp; b &lt; c &gt; d</p>`)
	if !strings.Contains(got, "a & b < c > d") {
		t.Errorf("should unescape HTML entities, got: %s", got)
	}
}

func TestLcHTMLToText_SuperscriptConversion(t *testing.T) {
	html := `<p>2<sup>31</sup> - 1</p>`
	// Superscripts are converted in lcParseContent (before lcHTMLToText),
	// so test via lcParseContent.
	desc, _, _ := lcParseContent(html)
	if !strings.Contains(desc, "^31") {
		t.Errorf("superscripts should be converted to ^N, got: %s", desc)
	}
}

func TestLcHTMLToText_StripsAllTags(t *testing.T) {
	html := `<div class="foo"><span>hello</span> <a href="bar">world</a></div>`
	got := lcHTMLToText(html)
	if strings.Contains(got, "<") || strings.Contains(got, ">") {
		t.Errorf("should strip all HTML tags, got: %s", got)
	}
	if !strings.Contains(got, "hello") || !strings.Contains(got, "world") {
		t.Errorf("should preserve text content, got: %s", got)
	}
}

// --- Authenticator tests ---

func TestLeetCodeProvider_Authenticate_EnvVars(t *testing.T) {
	lc := &LeetCodeProvider{client: http.DefaultClient}

	// Clear env vars first.
	origSession := os.Getenv(envLeetCodeSession)
	origCSRF := os.Getenv(envLeetCodeCSRFToken)
	defer func() {
		os.Setenv(envLeetCodeSession, origSession)
		os.Setenv(envLeetCodeCSRFToken, origCSRF)
	}()

	// Missing both → error.
	os.Unsetenv(envLeetCodeSession)
	os.Unsetenv(envLeetCodeCSRFToken)
	if err := lc.Authenticate(); err == nil {
		t.Error("Authenticate should fail with missing env vars")
	}
	if lc.IsAuthenticated() {
		t.Error("should not be authenticated with missing env vars")
	}

	// Set only session → error, and both fields should be cleared.
	os.Setenv(envLeetCodeSession, "test-session")
	os.Unsetenv(envLeetCodeCSRFToken)
	if err := lc.Authenticate(); err == nil {
		t.Error("Authenticate should fail with only session set")
	}
	if lc.session != "" || lc.csrfToken != "" {
		t.Error("Authenticate failure should clear both session and csrfToken (no half-state)")
	}

	// Set both → success.
	os.Setenv(envLeetCodeSession, "test-session")
	os.Setenv(envLeetCodeCSRFToken, "test-csrf")
	if err := lc.Authenticate(); err != nil {
		t.Errorf("Authenticate should succeed with both env vars: %v", err)
	}
	if !lc.IsAuthenticated() {
		t.Error("should be authenticated after setting both env vars")
	}
}

func TestLeetCodeProvider_AuthHelp(t *testing.T) {
	lc := &LeetCodeProvider{client: http.DefaultClient}
	help := lc.AuthHelp()
	if !strings.Contains(help, "LEETCODE_SESSION") {
		t.Error("AuthHelp should mention LEETCODE_SESSION")
	}
	if !strings.Contains(help, "LEETCODE_CSRFTOKEN") {
		t.Error("AuthHelp should mention LEETCODE_CSRFTOKEN")
	}
}

func TestLeetCodeProvider_RunTests_NotAuthenticated(t *testing.T) {
	lc := &LeetCodeProvider{client: http.DefaultClient}
	_, err := lc.RunTests("two-sum", "golang", "code", "input")
	if err == nil {
		t.Error("RunTests should fail when not authenticated")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("error should mention auth, got: %v", err)
	}
}

func TestLeetCodeProvider_Submit_NotAuthenticated(t *testing.T) {
	lc := &LeetCodeProvider{client: http.DefaultClient}
	_, err := lc.Submit("two-sum", "golang", "code")
	if err == nil {
		t.Error("Submit should fail when not authenticated")
	}
	if !strings.Contains(err.Error(), "not authenticated") {
		t.Errorf("error should mention auth, got: %v", err)
	}
}

// --- lcCheckResponse conversion tests ---

func TestLcCheckResponse_ToTestResult_Accepted(t *testing.T) {
	raw := lcCheckResponse{
		State:              "SUCCESS",
		StatusCode:         10,
		StatusMsg:          "Accepted",
		StatusRuntime:      "5 ms",
		StatusMemory:       "3.2 MB",
		CodeAnswer:         []string{"[0,1]"},
		ExpectedCodeAnswer: []string{"[0,1]"},
	}

	tr := raw.toTestResult()
	if !tr.Passed {
		t.Error("Accepted test should be marked as passed")
	}
	if tr.Actual != "[0,1]" {
		t.Errorf("Actual = %q, want %q", tr.Actual, "[0,1]")
	}
	if tr.Expected != "[0,1]" {
		t.Errorf("Expected = %q, want %q", tr.Expected, "[0,1]")
	}
	if tr.RuntimeMs != 5 {
		t.Errorf("RuntimeMs = %d, want 5", tr.RuntimeMs)
	}
	if !strings.Contains(tr.RawOutput, "Accepted") {
		t.Errorf("RawOutput should contain 'Accepted', got: %s", tr.RawOutput)
	}
}

func TestLcCheckResponse_ToTestResult_WrongAnswer(t *testing.T) {
	raw := lcCheckResponse{
		State:              "SUCCESS",
		StatusCode:         11,
		StatusMsg:          "Wrong Answer",
		CodeAnswer:         []string{"[1,0]"},
		ExpectedCodeAnswer: []string{"[0,1]"},
		LastTestcase:       "[2,7,11,15]\n9",
	}

	tr := raw.toTestResult()
	if tr.Passed {
		t.Error("Wrong Answer test should not be marked as passed")
	}
	if !strings.Contains(tr.RawOutput, "Wrong Answer") {
		t.Errorf("RawOutput should contain 'Wrong Answer', got: %s", tr.RawOutput)
	}
	if tr.Input != "[2,7,11,15]\n9" {
		t.Error("Input should be populated from LastTestcase")
	}
}

func TestLcCheckResponse_ToTestResult_Status10_MismatchedOutput(t *testing.T) {
	// This is the key bug fix: LeetCode's interpret_solution returns status_code 10
	// even when outputs don't match. Passed must be determined by comparing outputs.
	raw := lcCheckResponse{
		State:              "SUCCESS",
		StatusCode:         10,
		StatusMsg:          "Accepted",
		StatusRuntime:      "0 ms",
		CodeAnswer:         []string{"-1", "-1"},
		ExpectedCodeAnswer: []string{"4", "-1"},
	}

	tr := raw.toTestResult()
	if tr.Passed {
		t.Error("test with mismatched outputs should NOT be marked as passed, even with status_code 10")
	}
	if !strings.Contains(tr.RawOutput, "Wrong Answer") {
		t.Errorf("RawOutput should say 'Wrong Answer', got: %s", tr.RawOutput)
	}
	if !strings.Contains(tr.RawOutput, "FAIL") {
		t.Errorf("RawOutput should show per-case FAIL marker, got: %s", tr.RawOutput)
	}
	if !strings.Contains(tr.RawOutput, "PASS") {
		t.Errorf("RawOutput should show per-case PASS for the correct case, got: %s", tr.RawOutput)
	}
}

func TestLcCheckResponse_ToTestResult_MultiCase_AllPass(t *testing.T) {
	raw := lcCheckResponse{
		State:              "SUCCESS",
		StatusCode:         10,
		StatusMsg:          "Accepted",
		StatusRuntime:      "3 ms",
		CodeAnswer:         []string{"4", "-1"},
		ExpectedCodeAnswer: []string{"4", "-1"},
	}

	tr := raw.toTestResult()
	if !tr.Passed {
		t.Error("all matching outputs should be marked as passed")
	}
	if !strings.Contains(tr.RawOutput, "Accepted") {
		t.Errorf("RawOutput should say 'Accepted', got: %s", tr.RawOutput)
	}
}

func TestLcCheckResponse_ToTestResult_CompileError(t *testing.T) {
	raw := lcCheckResponse{
		State:            "SUCCESS",
		StatusCode:       20,
		StatusMsg:        "Compile Error",
		FullCompileError: "Line 1: undefined: foo",
	}

	tr := raw.toTestResult()
	if tr.Passed {
		t.Error("Compile error should not be marked as passed")
	}
	if tr.CompileError != "Line 1: undefined: foo" {
		t.Errorf("CompileError = %q, want full error", tr.CompileError)
	}
	if !strings.Contains(tr.RawOutput, "Compile Error") {
		t.Errorf("RawOutput should show compile error, got: %s", tr.RawOutput)
	}
}

func TestLcCheckResponse_ToTestResult_RuntimeError(t *testing.T) {
	raw := lcCheckResponse{
		State:            "SUCCESS",
		StatusCode:       15,
		StatusMsg:        "Runtime Error",
		FullRuntimeError: "index out of range [5] with length 3",
	}

	tr := raw.toTestResult()
	if tr.Passed {
		t.Error("Runtime error should not be marked as passed")
	}
	if tr.RuntimeError != "index out of range [5] with length 3" {
		t.Error("RuntimeError should use full error when available")
	}
}

func TestLcCheckResponse_ToSubmitResult_Accepted(t *testing.T) {
	totalCorrect := 50
	totalTestcases := 50
	raw := lcCheckResponse{
		State:             "SUCCESS",
		StatusCode:        10,
		StatusMsg:         "Accepted",
		StatusRuntime:     "3 ms",
		StatusMemory:      "2.8 MB",
		RuntimePercentile: json.Number("95.2"),
		MemoryPercentile:  json.Number("80.1"),
		TotalCorrect:      &totalCorrect,
		TotalTestcases:    &totalTestcases,
	}

	sr := raw.toSubmitResult()
	if !sr.Accepted {
		t.Error("Accepted submit should be marked as accepted")
	}
	if sr.StatusMsg != "Accepted" {
		t.Errorf("StatusMsg = %q, want 'Accepted'", sr.StatusMsg)
	}
	if sr.PassedCases != 50 || sr.TotalCases != 50 {
		t.Errorf("Cases = %d/%d, want 50/50", sr.PassedCases, sr.TotalCases)
	}
	if sr.RuntimePct != "faster than 95.2%" {
		t.Errorf("RuntimePct = %q", sr.RuntimePct)
	}
	if sr.MemoryPct != "less than 80.1%" {
		t.Errorf("MemoryPct = %q", sr.MemoryPct)
	}
	if !strings.Contains(sr.RawOutput, "Accepted") {
		t.Error("RawOutput should contain Accepted")
	}
}

func TestLcCheckResponse_ToSubmitResult_Accepted_NoPercentile(t *testing.T) {
	totalCorrect := 50
	totalTestcases := 50
	raw := lcCheckResponse{
		State:          "SUCCESS",
		StatusCode:     10,
		StatusMsg:      "Accepted",
		StatusRuntime:  "3 ms",
		StatusMemory:   "2.8 MB",
		TotalCorrect:   &totalCorrect,
		TotalTestcases: &totalTestcases,
		// No RuntimePercentile or MemoryPercentile — should not render "()"
	}

	sr := raw.toSubmitResult()
	if !sr.Accepted {
		t.Error("should be accepted")
	}
	if strings.Contains(sr.RawOutput, "()") {
		t.Errorf("should not contain empty parens, got:\n%s", sr.RawOutput)
	}
	if !strings.Contains(sr.RawOutput, "Runtime: 3 ms") {
		t.Errorf("should still show runtime, got:\n%s", sr.RawOutput)
	}
}

func TestLcCheckResponse_ToSubmitResult_WrongAnswer(t *testing.T) {
	totalCorrect := 45
	totalTestcases := 50
	raw := lcCheckResponse{
		State:          "SUCCESS",
		StatusCode:     11,
		StatusMsg:      "Wrong Answer",
		TotalCorrect:   &totalCorrect,
		TotalTestcases: &totalTestcases,
		LastTestcase:   "[1,2,3]",
		ExpectedOutput: "[3,2,1]",
	}

	sr := raw.toSubmitResult()
	if sr.Accepted {
		t.Error("Wrong Answer should not be accepted")
	}
	if sr.PassedCases != 45 || sr.TotalCases != 50 {
		t.Errorf("Cases = %d/%d, want 45/50", sr.PassedCases, sr.TotalCases)
	}
	if !strings.Contains(sr.RawOutput, "Wrong Answer") {
		t.Error("RawOutput should contain Wrong Answer")
	}
	if !strings.Contains(sr.RawOutput, "[1,2,3]") {
		t.Error("RawOutput should include failing input")
	}
}

func TestLcCheckResponse_ToSubmitResult_TLE(t *testing.T) {
	totalCorrect := 40
	totalTestcases := 50
	raw := lcCheckResponse{
		State:          "SUCCESS",
		StatusCode:     14,
		StatusMsg:      "Time Limit Exceeded",
		TotalCorrect:   &totalCorrect,
		TotalTestcases: &totalTestcases,
	}

	sr := raw.toSubmitResult()
	if sr.Accepted {
		t.Error("TLE should not be accepted")
	}
	if sr.StatusMsg != "Time Limit Exceeded" {
		t.Error("StatusMsg should be TLE")
	}
}

func TestParseRuntimeMs(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"5 ms", 5},
		{"150 ms", 150},
		{"0 ms", 0},
		{"N/A", 0},
		{"", 0},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseRuntimeMs(tt.input); got != tt.want {
				t.Errorf("parseRuntimeMs(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestFirstNonEmpty(t *testing.T) {
	if got := firstNonEmpty("", "", "third"); got != "third" {
		t.Errorf("firstNonEmpty = %q, want %q", got, "third")
	}
	if got := firstNonEmpty("first", "second"); got != "first" {
		t.Errorf("firstNonEmpty = %q, want %q", got, "first")
	}
	if got := firstNonEmpty("", ""); got != "" {
		t.Errorf("firstNonEmpty = %q, want %q", got, "")
	}
}

// --- HTTP integration tests with httptest ---

// newTestLeetCode creates a LeetCodeProvider pointed at the given httptest server.
func newTestLeetCode(serverURL string) *LeetCodeProvider {
	return &LeetCodeProvider{
		client:       &http.Client{Timeout: 5 * time.Second},
		baseURL:      serverURL,
		session:      "test-session",
		csrfToken:    "test-csrf",
		pollInterval: 10 * time.Millisecond,
		pollTimeout:  2 * time.Second,
		questionIDs:  make(map[string]string),
	}
}

func TestLeetCodeProvider_RunTests_FullFlow(t *testing.T) {
	checkCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":"runcode_abc123"}`)
		case strings.Contains(r.URL.Path, "check"):
			checkCalls++
			if checkCalls <= 2 {
				fmt.Fprint(w, `{"state":"PENDING"}`)
			} else {
				fmt.Fprint(w, `{"state":"SUCCESS","status_code":10,"status_msg":"Accepted","status_runtime":"5 ms","code_answer":["[0,1]"],"expected_code_answer":["[0,1]"],"code_output":"","std_output":""}`)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)

	// Step 1: RunTests should return the interpret_id.
	runID, err := lc.RunTests("two-sum", "golang", "func twoSum() {}", "[2,7,11,15]\n9")
	if err != nil {
		t.Fatalf("RunTests: %v", err)
	}
	if runID != "runcode_abc123" {
		t.Errorf("runID = %q, want %q", runID, "runcode_abc123")
	}

	// Step 2: Poll until done.
	var result *TestResult
	for i := 0; i < 10; i++ {
		time.Sleep(lc.pollInterval)
		r, done, err := lc.CheckTestResult(runID)
		if err != nil {
			t.Fatalf("CheckTestResult: %v", err)
		}
		if done {
			result = r
			break
		}
	}
	if result == nil {
		t.Fatal("never got a result from polling")
	}
	if !result.Passed {
		t.Error("result should be Passed")
	}
	if result.Actual != "[0,1]" {
		t.Errorf("Actual = %q, want %q", result.Actual, "[0,1]")
	}
	if !strings.Contains(result.RawOutput, "Accepted") {
		t.Errorf("RawOutput should contain 'Accepted', got: %s", result.RawOutput)
	}
}

func TestLeetCodeProvider_RunTests_WrongAnswer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":"runcode_wrong"}`)
		case strings.Contains(r.URL.Path, "check"):
			fmt.Fprint(w, `{"state":"SUCCESS","status_code":11,"status_msg":"Wrong Answer","code_answer":["[1,0]"],"expected_code_answer":["[0,1]"],"last_testcase":"[2,7,11,15]\n9","code_output":"","std_output":""}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	runID, err := lc.RunTests("two-sum", "golang", "func twoSum() {}", "[2,7,11,15]\n9")
	if err != nil {
		t.Fatalf("RunTests: %v", err)
	}

	result, done, err := lc.CheckTestResult(runID)
	if err != nil {
		t.Fatalf("CheckTestResult: %v", err)
	}
	if !done {
		t.Fatal("should be done")
	}
	if result.Passed {
		t.Error("wrong answer should not be Passed")
	}
	if !strings.Contains(result.RawOutput, "Wrong Answer") {
		t.Errorf("RawOutput should say Wrong Answer, got: %s", result.RawOutput)
	}
}

func TestLeetCodeProvider_RunTests_Status10_OutputMismatch(t *testing.T) {
	// Regression: LeetCode returns status_code 10 for interpret_solution even when
	// outputs don't match. The old code showed "Accepted" — should show "Wrong Answer".
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":"runcode_mismatch"}`)
		case strings.Contains(r.URL.Path, "check"):
			fmt.Fprint(w, `{"state":"SUCCESS","status_code":10,"status_msg":"Accepted","status_runtime":"0 ms","code_answer":["-1","-1"],"expected_code_answer":["4","-1"],"code_output":"","std_output":""}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	runID, err := lc.RunTests("binary-search", "golang", "func search() {}", "[-1,0,3,5,9,12]\n9\n[-1,0,3,5,9,12]\n2")
	if err != nil {
		t.Fatalf("RunTests: %v", err)
	}

	result, done, err := lc.CheckTestResult(runID)
	if err != nil {
		t.Fatalf("CheckTestResult: %v", err)
	}
	if !done {
		t.Fatal("should be done")
	}
	if result.Passed {
		t.Error("mismatched outputs should NOT be marked as Passed")
	}
	if !strings.Contains(result.RawOutput, "Wrong Answer") {
		t.Errorf("should show 'Wrong Answer', got: %s", result.RawOutput)
	}
	if !strings.Contains(result.RawOutput, "FAIL") {
		t.Errorf("should show per-case FAIL, got: %s", result.RawOutput)
	}
}

func TestLeetCodeProvider_RunTests_CompileError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":"runcode_ce"}`)
		case strings.Contains(r.URL.Path, "check"):
			fmt.Fprint(w, `{"state":"SUCCESS","status_code":20,"status_msg":"Compile Error","compile_error":"short","full_compile_error":"Line 1: undefined: foo","code_output":"","std_output":""}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	runID, _ := lc.RunTests("two-sum", "golang", "bad code", "")
	result, done, _ := lc.CheckTestResult(runID)
	if !done {
		t.Fatal("should be done")
	}
	if result.Passed {
		t.Error("compile error should not be Passed")
	}
	if result.CompileError != "Line 1: undefined: foo" {
		t.Errorf("CompileError = %q, want full error", result.CompileError)
	}
}

func TestLeetCodeProvider_RunTests_RuntimeError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":"runcode_re"}`)
		case strings.Contains(r.URL.Path, "check"):
			fmt.Fprint(w, `{"state":"SUCCESS","status_code":15,"status_msg":"Runtime Error","runtime_error":"short","full_runtime_error":"index out of range [5] with length 3","code_output":"","std_output":""}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	runID, _ := lc.RunTests("two-sum", "golang", "code", "")
	result, done, _ := lc.CheckTestResult(runID)
	if !done || result.Passed {
		t.Error("runtime error should be done and not Passed")
	}
	if result.RuntimeError != "index out of range [5] with length 3" {
		t.Errorf("RuntimeError = %q", result.RuntimeError)
	}
}

func TestLeetCodeProvider_Submit_FullFlow(t *testing.T) {
	checkCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"238"}}}`)
		case strings.HasSuffix(r.URL.Path, "/submit/"):
			fmt.Fprint(w, `{"submission_id":123456}`)
		case strings.Contains(r.URL.Path, "check"):
			checkCalls++
			if checkCalls <= 1 {
				fmt.Fprint(w, `{"state":"STARTED"}`)
			} else {
				fmt.Fprint(w, `{"state":"SUCCESS","status_code":10,"status_msg":"Accepted","status_runtime":"3 ms","status_memory":"2.8 MB","runtime_percentile":"95.2","memory_percentile":"80.1","total_correct":50,"total_testcases":50,"code_output":"","std_output":""}`)
			}
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)

	subID, err := lc.Submit("product-of-array-except-self", "golang", "func productExceptSelf() {}")
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}
	if subID != "123456" {
		t.Errorf("subID = %q, want %q", subID, "123456")
	}

	var result *SubmitResult
	for i := 0; i < 10; i++ {
		time.Sleep(lc.pollInterval)
		r, done, err := lc.CheckSubmission(subID)
		if err != nil {
			t.Fatalf("CheckSubmission: %v", err)
		}
		if done {
			result = r
			break
		}
	}
	if result == nil {
		t.Fatal("never got a submission result")
	}
	if !result.Accepted {
		t.Error("should be Accepted")
	}
	if result.PassedCases != 50 || result.TotalCases != 50 {
		t.Errorf("cases = %d/%d", result.PassedCases, result.TotalCases)
	}
	if !strings.Contains(result.RawOutput, "Accepted") {
		t.Error("RawOutput should contain Accepted")
	}
	if !strings.Contains(result.RawOutput, "3 ms") {
		t.Error("RawOutput should show runtime")
	}
}

func TestLeetCodeProvider_Submit_WrongAnswer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.HasSuffix(r.URL.Path, "/submit/"):
			fmt.Fprint(w, `{"submission_id":999}`)
		case strings.Contains(r.URL.Path, "check"):
			fmt.Fprint(w, `{"state":"SUCCESS","status_code":11,"status_msg":"Wrong Answer","total_correct":45,"total_testcases":50,"last_testcase":"[1,2,3]","expected_output":"[3,2,1]","code_output":"","std_output":""}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	subID, _ := lc.Submit("two-sum", "golang", "code")
	result, done, _ := lc.CheckSubmission(subID)
	if !done {
		t.Fatal("should be done")
	}
	if result.Accepted {
		t.Error("Wrong Answer should not be Accepted")
	}
	if result.PassedCases != 45 || result.TotalCases != 50 {
		t.Errorf("cases = %d/%d, want 45/50", result.PassedCases, result.TotalCases)
	}
	if !strings.Contains(result.RawOutput, "Wrong Answer") {
		t.Error("RawOutput should say Wrong Answer")
	}
	if !strings.Contains(result.RawOutput, "[1,2,3]") {
		t.Error("RawOutput should include failing input")
	}
}

func TestLeetCodeProvider_Submit_TLE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.HasSuffix(r.URL.Path, "/submit/"):
			fmt.Fprint(w, `{"submission_id":888}`)
		case strings.Contains(r.URL.Path, "check"):
			fmt.Fprint(w, `{"state":"SUCCESS","status_code":14,"status_msg":"Time Limit Exceeded","total_correct":40,"total_testcases":50,"code_output":"","std_output":""}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	subID, _ := lc.Submit("two-sum", "golang", "code")
	result, _, _ := lc.CheckSubmission(subID)
	if result.Accepted {
		t.Error("TLE should not be Accepted")
	}
	if result.StatusMsg != "Time Limit Exceeded" {
		t.Errorf("StatusMsg = %q", result.StatusMsg)
	}
}

// --- Error scenario tests ---

func TestLeetCodeProvider_RunTests_AuthExpired403(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"error":"auth expired"}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	_, err := lc.RunTests("two-sum", "golang", "code", "input")
	if err == nil {
		t.Fatal("should fail on 403")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Errorf("error should mention 403, got: %v", err)
	}
}

func TestLeetCodeProvider_Submit_EmptySubmissionID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.HasSuffix(r.URL.Path, "/submit/"):
			fmt.Fprint(w, `{"submission_id":0}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	_, err := lc.Submit("two-sum", "golang", "code")
	if err == nil {
		t.Fatal("should fail on empty submission_id")
	}
	if !strings.Contains(err.Error(), "empty submission_id") {
		t.Errorf("error = %v", err)
	}
}

func TestLeetCodeProvider_RunTests_EmptyInterpretID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":""}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	_, err := lc.RunTests("two-sum", "golang", "code", "input")
	if err == nil {
		t.Fatal("should fail on empty interpret_id")
	}
	if !strings.Contains(err.Error(), "empty interpret_id") {
		t.Errorf("error = %v", err)
	}
}

func TestLeetCodeProvider_CheckTestResult_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{invalid json`)
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	_, _, err := lc.CheckTestResult("some-id")
	if err == nil {
		t.Fatal("should fail on malformed JSON")
	}
}

func TestLeetCodeProvider_FetchQuestionID_CacheHit(t *testing.T) {
	// Pre-populate the cache so no GraphQL call is needed.
	graphqlCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			graphqlCalled = true
			fmt.Fprint(w, `{"data":{"question":{"questionId":"99"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":"runcode_cached"}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	lc.questionIDs["two-sum"] = "1" // pre-cached

	runID, err := lc.RunTests("two-sum", "golang", "code", "input")
	if err != nil {
		t.Fatalf("RunTests: %v", err)
	}
	if runID != "runcode_cached" {
		t.Errorf("runID = %q, want %q", runID, "runcode_cached")
	}
	if graphqlCalled {
		t.Error("GraphQL should NOT have been called when questionID is cached")
	}
}

func TestLeetCodeProvider_FetchQuestionID_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data":{"question":{"questionId":""}}}`)
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	_, err := lc.RunTests("nonexistent-problem", "golang", "code", "input")
	if err == nil {
		t.Fatal("should fail for unknown problem")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %v", err)
	}
}

func TestLeetCodeProvider_FlexStrings_InRealResponse(t *testing.T) {
	// The exact scenario that caused the original crash: code_output as a string.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			fmt.Fprint(w, `{"interpret_id":"runcode_flex"}`)
		case strings.Contains(r.URL.Path, "check"):
			// code_output as string, not array — the bug that crashed.
			fmt.Fprint(w, `{"state":"SUCCESS","status_code":10,"status_msg":"Accepted","status_runtime":"3 ms","code_answer":["[0,1]"],"expected_code_answer":["[0,1]"],"code_output":"debug line here","std_output":"stdout stuff"}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	runID, err := lc.RunTests("two-sum", "golang", "code", "input")
	if err != nil {
		t.Fatalf("RunTests: %v", err)
	}

	result, done, err := lc.CheckTestResult(runID)
	if err != nil {
		t.Fatalf("CheckTestResult should not crash on string code_output: %v", err)
	}
	if !done || !result.Passed {
		t.Error("should be done and passed")
	}
}

func TestLeetCodeProvider_Submit_RequestPayload(t *testing.T) {
	// Verify the payload sent to LeetCode is correct.
	var gotPayload map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"42"}}}`)
		case strings.HasSuffix(r.URL.Path, "/submit/"):
			json.NewDecoder(r.Body).Decode(&gotPayload)
			fmt.Fprint(w, `{"submission_id":777}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	_, err := lc.Submit("two-sum", "go", "func twoSum() {}")
	if err != nil {
		t.Fatalf("Submit: %v", err)
	}

	if gotPayload["lang"] != "golang" {
		t.Errorf("lang = %q, want 'golang' (normalized)", gotPayload["lang"])
	}
	if gotPayload["question_id"] != "42" {
		t.Errorf("question_id = %q, want '42'", gotPayload["question_id"])
	}
	if gotPayload["typed_code"] != "func twoSum() {}" {
		t.Errorf("typed_code = %q", gotPayload["typed_code"])
	}
}

func TestLeetCodeProvider_RunTests_RequestPayload(t *testing.T) {
	var gotPayload map[string]string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			json.NewDecoder(r.Body).Decode(&gotPayload)
			fmt.Fprint(w, `{"interpret_id":"runcode_payload"}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	_, err := lc.RunTests("two-sum", "python", "def twoSum(): pass", "[2,7]\n9")
	if err != nil {
		t.Fatalf("RunTests: %v", err)
	}

	if gotPayload["lang"] != "python3" {
		t.Errorf("lang = %q, want 'python3' (normalized)", gotPayload["lang"])
	}
	if gotPayload["data_input"] != "[2,7]\n9" {
		t.Errorf("data_input = %q", gotPayload["data_input"])
	}
}

func TestLeetCodeProvider_AuthHeaders_SentCorrectly(t *testing.T) {
	var gotCookie, gotCSRF, gotReferer string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "graphql"):
			fmt.Fprint(w, `{"data":{"question":{"questionId":"1"}}}`)
		case strings.Contains(r.URL.Path, "interpret_solution"):
			gotCookie = r.Header.Get("Cookie")
			gotCSRF = r.Header.Get("X-CSRFToken")
			gotReferer = r.Header.Get("Referer")
			fmt.Fprint(w, `{"interpret_id":"runcode_auth"}`)
		}
	}))
	defer server.Close()

	lc := newTestLeetCode(server.URL)
	lc.session = "my-session-123"
	lc.csrfToken = "my-csrf-456"

	lc.RunTests("two-sum", "golang", "code", "input")

	if !strings.Contains(gotCookie, "LEETCODE_SESSION=my-session-123") {
		t.Errorf("Cookie = %q, should contain session", gotCookie)
	}
	if !strings.Contains(gotCookie, "csrftoken=my-csrf-456") {
		t.Errorf("Cookie = %q, should contain csrftoken", gotCookie)
	}
	if gotCSRF != "my-csrf-456" {
		t.Errorf("X-CSRFToken = %q", gotCSRF)
	}
	if gotReferer != server.URL {
		t.Errorf("Referer = %q, want %q (baseURL)", gotReferer, server.URL)
	}
}

// --- flexStrings tests ---

func TestFlexStrings_Array(t *testing.T) {
	input := `{"code_output":["line1","line2"]}`
	var result struct {
		CodeOutput flexStrings `json:"code_output"`
	}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.CodeOutput) != 2 || result.CodeOutput[0] != "line1" {
		t.Errorf("got %v, want [line1 line2]", result.CodeOutput)
	}
}

func TestFlexStrings_String(t *testing.T) {
	input := `{"code_output":"single line"}`
	var result struct {
		CodeOutput flexStrings `json:"code_output"`
	}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(result.CodeOutput) != 1 || result.CodeOutput[0] != "single line" {
		t.Errorf("got %v, want [single line]", result.CodeOutput)
	}
}

func TestFlexStrings_EmptyString(t *testing.T) {
	input := `{"code_output":""}`
	var result struct {
		CodeOutput flexStrings `json:"code_output"`
	}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.CodeOutput != nil {
		t.Errorf("got %v, want nil", result.CodeOutput)
	}
}

func TestFlexStrings_Null(t *testing.T) {
	input := `{"code_output":null}`
	var result struct {
		CodeOutput flexStrings `json:"code_output"`
	}
	if err := json.Unmarshal([]byte(input), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.CodeOutput != nil {
		t.Errorf("got %v, want nil", result.CodeOutput)
	}
}

func TestFlexStrings_FullCheckResponse(t *testing.T) {
	// Simulate a real LeetCode response where code_output is a string
	// (the exact case that caused the original crash).
	input := `{
		"state": "SUCCESS",
		"status_code": 10,
		"status_msg": "Accepted",
		"status_runtime": "3 ms",
		"code_answer": ["[0,1]"],
		"expected_code_answer": ["[0,1]"],
		"code_output": "some debug output",
		"std_output": ""
	}`
	var resp lcCheckResponse
	if err := json.Unmarshal([]byte(input), &resp); err != nil {
		t.Fatalf("unmarshal full response: %v", err)
	}
	if resp.StatusCode != 10 {
		t.Error("StatusCode should be 10")
	}
	if len(resp.CodeOutput) != 1 || resp.CodeOutput[0] != "some debug output" {
		t.Errorf("CodeOutput = %v, want [some debug output]", resp.CodeOutput)
	}
	// Verify conversion still works.
	tr := resp.toTestResult()
	if !tr.Passed {
		t.Error("should be passed")
	}
}

func TestLeetCodeProvider_SetAuthHeaders(t *testing.T) {
	lc := &LeetCodeProvider{
		client:    http.DefaultClient,
		baseURL:   "https://leetcode.com",
		session:   "my-session-value",
		csrfToken: "my-csrf-value",
	}

	req, _ := http.NewRequest("POST", "https://example.com", nil)
	lc.setAuthHeaders(req)

	cookie := req.Header.Get("Cookie")
	if !strings.Contains(cookie, "LEETCODE_SESSION=my-session-value") {
		t.Errorf("Cookie should contain session, got: %s", cookie)
	}
	if !strings.Contains(cookie, "csrftoken=my-csrf-value") {
		t.Errorf("Cookie should contain csrftoken, got: %s", cookie)
	}
	if req.Header.Get("X-CSRFToken") != "my-csrf-value" {
		t.Errorf("X-CSRFToken = %q, want %q", req.Header.Get("X-CSRFToken"), "my-csrf-value")
	}
	if req.Header.Get("Referer") != "https://leetcode.com" {
		t.Error("Referer should be set to baseURL")
	}
}

func TestTrimTrailingEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want int // expected length after trim
	}{
		{"no trailing", []string{"a", "b"}, 2},
		{"one trailing", []string{"a", "b", ""}, 2},
		{"two trailing", []string{"a", "", ""}, 1},
		{"all empty", []string{"", "", ""}, 0},
		{"nil", nil, 0},
		{"empty slice", []string{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimTrailingEmpty(tt.in)
			if len(got) != tt.want {
				t.Errorf("trimTrailingEmpty(%v) len = %d, want %d", tt.in, len(got), tt.want)
			}
		})
	}
}

func TestLcCheckResponse_ToTestResult_TrailingEmpty(t *testing.T) {
	// Regression: LeetCode returns ["4", "-1", ""] — trailing empty string
	// should not produce a phantom "Case 3: PASS  output=  expected=" line.
	raw := lcCheckResponse{
		State:              "SUCCESS",
		StatusCode:         10,
		StatusMsg:          "Accepted",
		StatusRuntime:      "3 ms",
		CodeAnswer:         []string{"4", "-1", ""},
		ExpectedCodeAnswer: []string{"4", "-1", ""},
	}

	tr := raw.toTestResult()
	if !tr.Passed {
		t.Error("matching outputs (after trim) should be Passed")
	}
	if strings.Contains(tr.RawOutput, "Case 3") {
		t.Errorf("should not show phantom Case 3, got:\n%s", tr.RawOutput)
	}
	if !strings.Contains(tr.RawOutput, "Accepted") {
		t.Errorf("should show Accepted, got:\n%s", tr.RawOutput)
	}
}

func TestFlexStringsEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b flexStrings
		want bool
	}{
		{"both nil", nil, nil, true},
		{"both empty", flexStrings{}, flexStrings{}, true},
		{"equal single", flexStrings{"a"}, flexStrings{"a"}, true},
		{"equal multi", flexStrings{"a", "b"}, flexStrings{"a", "b"}, true},
		{"different values", flexStrings{"a"}, flexStrings{"b"}, false},
		{"different lengths", flexStrings{"a", "b"}, flexStrings{"a"}, false},
		{"nil vs empty", nil, flexStrings{}, true}, // both len 0
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flexStringsEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("flexStringsEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestLeetCodeProvider_SetAuthHeaders_CustomBaseURL(t *testing.T) {
	lc := &LeetCodeProvider{
		client:    http.DefaultClient,
		baseURL:   "http://localhost:9999",
		session:   "s",
		csrfToken: "c",
	}

	req, _ := http.NewRequest("POST", "http://localhost:9999/graphql", nil)
	lc.setAuthHeaders(req)

	if req.Header.Get("Referer") != "http://localhost:9999" {
		t.Errorf("Referer should use baseURL, got: %s", req.Header.Get("Referer"))
	}
}
