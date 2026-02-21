# Provider Abstraction — Implementation Plan

## Goal

Replace the hard dependency on `leetgo` with a pluggable provider system so that:
1. `go install github.com/conrad760/warmup@latest` just works — no external tools needed
2. LeetCode is the first provider, but the architecture supports Codewars, HackerRank, custom servers, etc.
3. Multiple providers can be active in a single session (curated bank specifies provider per question)

## Architecture

```
┌─────────────────────────────────────────────┐
│                  warmup core                 │
│                                              │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│  │ TUI      │  │ SRS      │  │ Scaffold  │  │
│  │ (bubble  │  │ (review  │  │ (files,   │  │
│  │  tea)    │  │  .go)    │  │  editor)  │  │
│  └────┬─────┘  └────┬─────┘  └─────┬─────┘  │
│       └──────┬───────┘──────────────┘        │
│              │                               │
│       ┌──────┴──────┐                        │
│       │   loader    │  curated bank +        │
│       │   + cache   │  provider dispatch     │
│       └──────┬──────┘                        │
└──────────────┼───────────────────────────────┘
               │
    ┌──────────┼──────────┐
    │          │          │
┌───┴───┐ ┌───┴───┐ ┌────┴────┐
│Leet   │ │Mock   │ │Future:  │
│Code   │ │       │ │Codewars │
│       │ │       │ │Hacker-  │
│GraphQL│ │Hard-  │ │Rank,    │
│API    │ │coded  │ │Custom   │
└───────┘ └───────┘ └─────────┘
```

## Provider Interfaces

```go
// Provider fetches problem data from an external platform.
type Provider interface {
    Name() string
    FetchProblem(id string, lang string) (*ProblemData, error)
}

// Tester runs code against test cases remotely.
type Tester interface {
    RunTests(id string, lang string, code string, input string) (*TestResult, error)
    CheckResult(runID string) (*TestResult, error)
}

// Submitter grades solutions remotely.
type Submitter interface {
    Submit(id string, lang string, code string) (*SubmitResult, error)
    CheckSubmission(subID string) (*SubmitResult, error)
}

// Authenticator handles credentials for providers that need them.
type Authenticator interface {
    Authenticate() error
    IsAuthenticated() bool
    AuthHelp() string
}
```

Providers implement the interfaces they support. The core uses type assertions:
- `if t, ok := provider.(Tester); ok { ... }` — test/submit only available if provider supports it
- `if a, ok := provider.(Authenticator); ok && !a.IsAuthenticated() { ... }` — show auth help

## Data Model Changes

### CuratedQuestion (questions.go, questions_extended.go)

```go
type CuratedQuestion struct {
    Provider  string   // "leetcode" (default if empty), "mock", "codewars", etc.
    ProblemID string   // provider-specific identifier (was "Slug")
    Category  string
    Options   []Option
    Solution  string
}
```

### Question (runtime, main.go)

```go
type Question struct {
    Title       string
    Difficulty  Difficulty
    Category    string
    Description string
    Example     string
    Options     []Option
    Solution    string
    Provider    string   // which provider this came from
    ProblemID   string   // replaces LeetcodeSlug
    CodeSnippet string   // for scaffolding
}
```

## File Layout (after Phase 1)

```
warmup/
├── main.go               # TUI, model, view, key handling (modified)
├── questions.go           # curated bank (Slug→ProblemID rename)
├── questions_extended.go  # curated bank ext (same rename)
├── review.go              # SRS (unchanged)
├── provider.go            # interfaces + registry + normalized types
├── provider_leetcode.go   # LeetCode GraphQL implementation
├── provider_mock.go       # mock provider for testing
├── loader.go              # replaces leetgo_loader.go, uses providers
├── cache.go               # local question cache (~/.config/warmup/)
├── scaffold.go            # workspace + file management + editor
├── setup.go               # --setup guided flow
├── go.mod                 # updated module path, no sqlite
├── Makefile
└── README.md
```

## Phases

### Phase 1: Provider abstraction + fetch + scaffold (this branch)

Removes the leetgo dependency for startup + studying. Users can `go install` and
immediately use warmup. Test/submit stubbed with links to the website.

| #    | Task                                    | Size | Notes |
|------|-----------------------------------------|------|-------|
| 1.1  | Update go.mod module path               | XS   | `module github.com/conrad760/warmup`, drop sqlite |
| 1.2  | Create provider.go                      | S    | Interfaces, registry, ProblemData/TestResult/SubmitResult types |
| 1.3  | Create provider_mock.go                 | S    | Hardcoded problems for testing without network |
| 1.4  | Create provider_leetcode.go             | M    | GraphQL fetch, HTML-to-text, no auth needed |
| 1.5  | Create cache.go                         | S    | Provider-aware, ~/.config/warmup/cache/, 7-day TTL |
| 1.6  | Create loader.go                        | M    | Groups curated bank by provider, fetches via cache, produces []Question |
| 1.7  | Create scaffold.go                      | M    | ~/.config/warmup/workspace/, solution.go from CodeSnippet, $EDITOR |
| 1.8  | Create setup.go                         | S    | --setup flag, test connectivity, populate cache, auto-trigger |
| 1.9  | Update curated banks                    | S    | Slug→ProblemID rename, add Provider field |
| 1.10 | Update main.go                          | L    | Remove leetgo plumbing, wire providers + scaffold |
| 1.11 | Delete leetgo_loader.go                 | XS   | Fully replaced |
| 1.12 | Update README, Makefile                 | S    | New install instructions, remove leetgo refs |
| 1.13 | Verify build, tag v0.1.0                | XS   | `go install` works end-to-end |

**Dependency changes:**
- Remove: modernc.org/sqlite (and ~5 transitive deps)
- Add: none (stdlib only)

### Phase 2: Native test + submit (future branch)

| #    | Task                                    | Size | Notes |
|------|-----------------------------------------|------|-------|
| 2.1  | Create auth.go                          | M    | Env vars, .env file, provider-namespaced |
| 2.2  | Extend provider_leetcode.go             | M    | Implement Tester + Submitter + Authenticator |
| 2.3  | Wire test/submit in main.go             | M    | Type-assert, run as background cmds |
| 2.4  | Update setup.go for credentials         | S    | Test auth, show AuthHelp() |
| 2.5  | Browser cookie reading (optional)       | M    | kooky library, defer if needed |

### Phase 3: Additional providers (future)

Each provider is a single file (~100-150 lines):
- provider_codewars.go
- provider_hackerrank.go
- provider_custom.go (generic REST)

## Key Design Decisions

1. **Flags + env vars only** — no config file. Provider is specified per curated question. Auth via env vars.
2. **Multi-provider sessions** — curated bank specifies provider per question. Loader initializes all needed providers.
3. **Compile-time registry** — providers register in init(). Adding a provider = adding a Go file.
4. **Separate interfaces** — Provider, Tester, Submitter, Authenticator are distinct. Compose what you need.
5. **Cache is provider-aware** — `~/.config/warmup/cache/<provider>/<problem-id>.json`.
6. **Workspace is provider-aware** — `~/.config/warmup/workspace/<provider>/<problem-id>/solution.go`.
