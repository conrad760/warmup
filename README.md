# DSA Warmup

A CLI-based DSA interview prep tool with spaced repetition. Fetches problem descriptions from LeetCode (or other providers), pairs them with curated approach-based multiple choice, and lets you code solutions without leaving the terminal.

## Quick Start

```bash
# Install and run
go install github.com/conrad760/warmup@latest
warmup
```

That's it. No external tools, no database setup, no credentials needed for studying.

## Prerequisites

- **Go 1.24+**

## Usage

```bash
# Start a session
warmup

# Focus on a category
warmup -category "dynamic programming"

# Choose a different language for code snippets
warmup -lang python

# Load additional custom questions
warmup -questions my-questions.json

# List available categories
warmup -categories

# View lifetime stats without starting a session
warmup --stats
```

**Supported languages:** `go` (default), `python`, `java`, `cpp`, `javascript`, `typescript`, `rust`, `c`

## Controls

| Key | When | Action |
|-----|------|--------|
| `j/k` or `Up/Down` | Always | Navigate options / scroll content |
| `Enter` or `Space` | Before answering | Select answer |
| `t` | After optimal/plausible | Open problem in your editor |
| `T` | After `t` | Run test (coming in Phase 2) |
| `S` | After `t` | Submit solution (coming in Phase 2) |
| `s` | After answering | Toggle Go solution view |
| `n` | After answering | Next question |
| `p` | Before answering | Pause/resume timer |
| `r` | Before answering | Reset timer to 5:00 |
| `q` or `Ctrl+C` | Always | Quit (shows session report) |

## Editor Configuration

When you press `t`, warmup scaffolds a workspace and opens your editor. By default it opens `$EDITOR` with the solution file. For a multi-pane layout like leetgo, set `WARMUP_EDITOR`:

```bash
# nvim with 3 panes: solution, description, test cases
export WARMUP_EDITOR="nvim {{.CodeFile}} -c 'split {{.DescriptionFile}}' -c 'vsplit {{.TestCasesFile}}'"

# vscode: open the problem folder
export WARMUP_EDITOR="code {{.Folder}}"

# vim with tabs
export WARMUP_EDITOR="vim -p {{.CodeFile}} {{.DescriptionFile}} {{.TestCasesFile}}"
```

**Available templates:**

| Template | Resolves to |
|----------|-------------|
| `{{.CodeFile}}` | `solution.go` (or `.py`, `.java`, etc.) |
| `{{.DescriptionFile}}` | `question.md` |
| `{{.TestCasesFile}}` | `testcases.txt` |
| `{{.Folder}}` | Problem directory |

**Fallback chain:** `WARMUP_EDITOR` > `$EDITOR` (solution file only) > `nvim` > `vim` > `vi`

## Workspace Layout

When you press `t`, warmup creates a leetgo-style workspace per problem:

```
~/.config/warmup/workspace/
└── leetcode/
    └── two-sum/
        ├── solution.go      # your code between @lc markers + test harness
        ├── question.md      # problem description + examples
        └── testcases.txt    # input:/output: test case pairs
```

The solution file includes `// @lc code=begin` / `// @lc code=end` markers (matching leetgo's convention) and a generated `func main()` test harness that reads from stdin using `encoding/json` -- no external dependencies.

If you've already edited a solution, it is never overwritten.

## How It Works

1. **85 curated problems** covering the Blind 75 and more, across 17 categories
2. **Problem descriptions** fetched from LeetCode's public API and cached locally
3. **Approach options** are hand-written: each describes a strategy + complexity, not code
4. After answering, all options are **color-coded** (green/yellow/orange/red)
5. Press `s` for a **syntax-highlighted Go solution** with pattern name and complexity
6. Press `t` to **open in your editor** and solve it
7. **Spaced repetition (SM-2)** tracks your performance and schedules reviews

## Architecture

Built on a pluggable **provider** system. LeetCode is the default, but the architecture supports adding Codewars, HackerRank, or custom servers.

```
warmup core (TUI, SRS, scaffold)
         │
    ┌────┼────┐
    │    │    │
 LeetCode  Mock  (future providers)
```

| File | Purpose |
|------|---------|
| `main.go` | TUI (Bubble Tea), model, view, key handling, styles, syntax highlighting |
| `provider.go` | Provider interfaces + registry + FuncMeta types |
| `provider_leetcode.go` | LeetCode GraphQL API client + metadata parsing |
| `provider_mock.go` | Mock provider for testing |
| `loader.go` | Provider dispatch, question loading |
| `cache.go` | Local question data cache (7-day TTL) |
| `scaffold.go` | Workspace + file generation + editor integration |
| `questions.go` | 36 curated questions (approach options + Go solutions) |
| `questions_extended.go` | 49 more curated questions |
| `review.go` | SM-2 spaced repetition, review persistence, stats/reports |

## Provider System

Questions specify which provider they come from. The default is LeetCode:

```go
CuratedQuestion{
    Provider:  "leetcode",  // or "mock", "codewars", etc.
    ProblemID: "two-sum",   // provider-specific identifier
    Category:  "Arrays & Hashing",
    Options:   []Option{...},
    Solution:  "...",
}
```

Adding a new provider is a single file implementing the `Provider` interface.

## Spaced Repetition

The app uses an SM-2-based algorithm to track how well you know each problem. Your review quality is derived from two signals:

| Approach | Submit | Quality Score | Effect |
|----------|--------|---------------|--------|
| Optimal | Accepted | 5 | Interval grows fast |
| Optimal | No submit | 4 | Interval grows |
| Plausible | Accepted | 3 | Interval grows slowly |
| Plausible | No submit | 2 | Interval resets |
| Suboptimal | Any | 1 | Interval resets |
| Wrong | Any | 0 | Interval resets |

Problems you struggle with come back sooner. Problems you nail get pushed further out. A problem is considered "mastered" after 3+ consecutive correct reviews with a 21+ day interval.

**Question selection priority:**
1. Overdue problems (most overdue first)
2. New problems (never reviewed)
3. Not-yet-due problems (as fallback)

Review data is stored in `~/.config/warmup/reviews.json`.

## Adding Custom Questions

Create a JSON file with problem IDs and curated approaches:

```json
[
  {
    "ProblemID": "valid-anagram",
    "Options": [
      {"Text": "Frequency array for 26 chars -- O(n) time, O(1) space", "Rating": "OPTIMAL"},
      {"Text": "Sort both strings and compare -- O(n log n) time, O(1) space", "Rating": "PLAUSIBLE"},
      {"Text": "Check every permutation -- O(n!) time", "Rating": "WRONG"}
    ],
    "Solution": "// Pattern: Frequency Count\nfunc isAnagram(s, t string) bool {\n    // ...\n}"
  }
]
```

The `ProblemID` must match a problem on the provider. Title, difficulty, description, and examples are fetched automatically.

**Validation rules:**
- Each entry needs a `ProblemID`
- At least 2 options per question
- At least one option rated `OPTIMAL`
- Rating values: `OPTIMAL`, `PLAUSIBLE`, `SUBOPTIMAL`, `WRONG`

## Data Storage

All data is stored under `~/.config/warmup/`:

```
~/.config/warmup/
├── reviews.json                    # SRS review history
├── cache/
│   └── leetcode/
│       ├── two-sum.json           # cached problem data (7-day TTL)
│       └── ...
└── workspace/
    └── leetcode/
        └── two-sum/
            ├── solution.go        # @lc markers + test harness
            ├── question.md        # problem description
            └── testcases.txt      # input:/output: pairs
```
