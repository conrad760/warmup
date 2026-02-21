package main

import (
	"fmt"
	"sync"
)

// ProblemData is the normalized problem representation returned by any provider.
type ProblemData struct {
	ID          string    // provider-specific identifier (e.g. "two-sum" for LeetCode)
	Title       string    // human-readable title
	Description string    // cleaned plain text, not HTML
	Examples    string    // first example, plain text
	Difficulty  string    // "Easy", "Medium", "Hard" — normalized by provider
	Tags        []string  // topic/category tags
	CodeSnippet string    // language-specific function signature / template
	TestInput   string    // default test case input
	Meta        *FuncMeta // function signature metadata for test harness generation
}

// FuncMeta describes the function signature of a problem for generating test harnesses.
type FuncMeta struct {
	Name         string      // function name (e.g. "twoSum")
	Params       []ParamMeta // function parameters in order
	Return       *ParamMeta  // return type (nil for void/in-place problems)
	SystemDesign bool        // true for Constructor-based design problems
}

// ParamMeta describes a single function parameter or return type.
type ParamMeta struct {
	Name string // parameter name (e.g. "nums", "target")
	Type string // LeetCode type string (e.g. "integer[]", "string", "TreeNode")
}

// TestResult holds the outcome of running code against test cases.
type TestResult struct {
	Passed       bool
	Input        string
	Expected     string
	Actual       string
	RuntimeMs    int
	MemoryMB     float64
	CompileError string
	RuntimeError string
	RawOutput    string // full provider output for display
}

// SubmitResult holds the outcome of a graded submission.
type SubmitResult struct {
	Accepted     bool
	StatusMsg    string // e.g. "Accepted", "Wrong Answer", "Time Limit Exceeded"
	RuntimeMs    int
	RuntimePct   string // e.g. "faster than 95%"
	MemoryMB     float64
	MemoryPct    string
	TotalCases   int
	PassedCases  int
	CompileError string
	RuntimeError string
	RawOutput    string // full provider output for display
}

// Provider fetches problem data from an external platform.
// Every provider must implement at least this interface.
type Provider interface {
	// Name returns the provider identifier ("leetcode", "mock", "codewars", etc.).
	Name() string

	// FetchProblem fetches a single problem by its platform-specific ID.
	// lang is the target programming language (e.g. "go", "python").
	FetchProblem(id string, lang string) (*ProblemData, error)
}

// Tester is implemented by providers that support remote code execution against test cases.
type Tester interface {
	// RunTests submits code for testing. Returns a run ID for polling.
	RunTests(id string, lang string, code string, input string) (runID string, err error)

	// CheckTestResult polls for the result of a test run.
	CheckTestResult(runID string) (*TestResult, bool, error) // result, done, error
}

// Submitter is implemented by providers that support graded solution submission.
type Submitter interface {
	// Submit submits code for grading. Returns a submission ID for polling.
	Submit(id string, lang string, code string) (subID string, err error)

	// CheckSubmission polls for the result of a submission.
	CheckSubmission(subID string) (*SubmitResult, bool, error) // result, done, error
}

// Authenticator is implemented by providers that require credentials for some operations.
type Authenticator interface {
	// Authenticate loads and validates credentials.
	Authenticate() error

	// IsAuthenticated returns whether valid credentials are available.
	IsAuthenticated() bool

	// AuthHelp returns human-readable instructions for setting up credentials.
	AuthHelp() string
}

// --- Provider Registry ---

var (
	providersMu sync.RWMutex
	providers   = make(map[string]func() Provider)
)

// RegisterProvider registers a provider factory under the given name.
// Typically called from init() in each provider_*.go file.
func RegisterProvider(name string, factory func() Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	providers[name] = factory
}

// GetProvider returns a new instance of the named provider.
func GetProvider(name string) (Provider, error) {
	providersMu.RLock()
	defer providersMu.RUnlock()
	factory, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %q (available: %s)", name, AvailableProviders())
	}
	return factory(), nil
}

// AvailableProviders returns a comma-separated list of registered provider names.
func AvailableProviders() string {
	providersMu.RLock()
	defer providersMu.RUnlock()
	names := make([]string, 0, len(providers))
	for name := range providers {
		names = append(names, name)
	}
	return fmt.Sprintf("%v", names)
}

// DefaultProviderName is used when a CuratedQuestion doesn't specify a provider.
const DefaultProviderName = "leetcode"
