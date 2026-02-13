# DSA Warmup

A CLI-based DSA interview prep tool with spaced repetition. Pulls problem descriptions from leetgo's local database, pairs them with curated approach-based multiple choice, and lets you code, test, and submit solutions without leaving the terminal.

## Quick Start

```bash
# 1. Install leetgo
go install github.com/j178/leetgo@latest

# 2. Create a workspace and download the question database
mkdir ~/leetcode && cd ~/leetcode
leetgo init
leetgo cache update

# 3. Build and run
cd /path/to/warmup
go build -o dsa-warmup .
./dsa-warmup -workspace ~/leetcode
```

If you skip any step, the app will tell you exactly what's missing and how to fix it.

## Prerequisites

- **Go 1.24+** (required by bubbletea dependency)
- **[leetgo](https://github.com/j178/leetgo)** — manages LeetCode problems locally

## Usage

```bash
# Start a session (auto-detects workspace if leetgo.yaml is in cwd or parent)
./dsa-warmup

# Specify workspace explicitly
./dsa-warmup -workspace ~/leetcode

# Load additional custom questions
./dsa-warmup -questions my-questions.json

# View lifetime stats without starting a session
./dsa-warmup --stats
```

## Controls

| Key | When | Action |
|-----|------|--------|
| `j/k` or `Up/Down` | Always | Navigate options / scroll content |
| `Enter` or `Space` | Before answering | Select answer |
| `t` | After optimal/plausible | Open problem in leetgo (scaffolds first time, edits after) |
| `T` | After `t` | Run `leetgo test` — results shown inline |
| `S` | After `t` | Run `leetgo submit` — results shown inline |
| `s` | After answering | Toggle Go solution view |
| `n` | After answering | Next question |
| `p` | Before answering | Pause/resume timer |
| `r` | Before answering | Reset timer to 5:00 |
| `q` or `Ctrl+C` | Always | Quit (shows session report) |

## How It Works

1. **83 curated problems** covering the Blind 75 and more, across 17 categories
2. **Problem descriptions** pulled from leetgo's local SQLite cache — always current
3. **Approach options** are hand-written: each describes a strategy + complexity, not code
4. After answering, all options are **color-coded** (green/yellow/orange/red)
5. Press `s` for a **syntax-highlighted Go solution** with pattern name and complexity
6. Press `t` to **open in leetgo** and solve it — `T` to test, `S` to submit
7. Test and submit results appear **inline** in the TUI — no lost output
8. **Spaced repetition (SM-2)** tracks your performance and schedules reviews

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

### Lifetime Stats

```bash
./dsa-warmup --stats
```

Shows: total/reviewed/mastered counts, category breakdown with accuracy, weakest problems by ease factor, and due/upcoming counts.

### Session Report

On quit, a session summary is printed showing each problem reviewed, your approach rating, and submit result.

## Try It Flow

When you answer **optimal** or **plausible**:

1. Press `t` — scaffolds the problem (first time) or reopens your existing code (subsequent times)
2. Write your solution in nvim (or your configured editor)
3. Quit nvim to return to the TUI
4. Press `T` to test — results appear inline with a spinner while running
5. If wrong, press `t` to edit your code again (your work is preserved)
6. Press `S` to submit — verdict appears inline

The app uses `leetgo edit` (not `leetgo pick`) after the first scaffold, so your solution is never overwritten.

## Adding Custom Questions

Create a JSON file with LeetCode slugs and curated approaches:

```json
[
  {
    "Slug": "valid-anagram",
    "Options": [
      {"Text": "Frequency array for 26 chars — O(n) time, O(1) space", "Rating": "OPTIMAL"},
      {"Text": "Sort both strings and compare — O(n log n) time, O(1) space", "Rating": "PLAUSIBLE"},
      {"Text": "Check every permutation — O(n!) time", "Rating": "WRONG"}
    ],
    "Solution": "// Pattern: Frequency Count\nfunc isAnagram(s, t string) bool {\n    // ...\n}"
  }
]
```

The slug must match a problem in your leetgo database. Title, difficulty, description, and examples are pulled from the DB automatically.

**Validation rules:**
- Each entry needs a `Slug` matching a leetgo problem
- At least 2 options per question
- At least one option rated `OPTIMAL`
- Rating values: `OPTIMAL`, `PLAUSIBLE`, `SUBOPTIMAL`, `WRONG`

## Built-in Questions (83)

| Category | Count | Examples |
|----------|-------|---------|
| Arrays & Hashing | 8 | Two Sum, Group Anagrams, Product of Array Except Self, Longest Consecutive Sequence |
| Two Pointers | 4 | Valid Palindrome, 3Sum, Container With Most Water, Two Sum II |
| Sliding Window | 4 | Best Time to Buy/Sell Stock, Longest Substring No Repeat, Min Window Substring |
| Stack | 3 | Valid Parentheses, Min Stack, Evaluate RPN |
| Binary Search | 5 | Binary Search, Search Rotated Array, Koko Eating Bananas, Time Based Key-Value Store |
| Linked List | 6 | Reverse Linked List, Merge Two Lists, Reorder List, Merge K Sorted Lists |
| Trees | 11 | Invert Binary Tree, Validate BST, LCA of BST, Binary Tree Max Path Sum, Serialize/Deserialize |
| Heap | 1 | Find Median from Data Stream |
| Graphs | 7 | Number of Islands, Clone Graph, Course Schedule I/II, Pacific Atlantic Water Flow |
| Dynamic Programming | 10 | Climbing Stairs, House Robber I/II, Coin Change, Word Break, Longest Common Subsequence |
| Backtracking | 3 | Subsets, Combination Sum, Word Search |
| Greedy | 4 | Jump Game I/II, Maximum Subarray, Task Scheduler |
| Intervals | 5 | Insert Interval, Merge Intervals, Meeting Rooms I/II |
| Design | 3 | LRU Cache, Implement Trie, Design Twitter |
| Math & Geometry | 3 | Rotate Image, Spiral Matrix, Set Matrix Zeroes |
| Bit Manipulation | 5 | Number of 1 Bits, Counting Bits, Reverse Bits, Missing Number, Sum of Two Integers |
| Advanced Graphs | 1 | Alien Dictionary |

## Architecture

| File | Purpose |
|------|---------|
| `main.go` | TUI (Bubble Tea), model, view, key handling, styles, syntax highlighting |
| `questions.go` | 36 curated questions (slug + approach options + Go solution) |
| `questions_extended.go` | 47 more curated questions (same format) |
| `leetgo_loader.go` | DB discovery, SQLite queries, HTML-to-text parsing |
| `review.go` | SM-2 spaced repetition, review persistence, stats/reports |

**Data split:** The Go source files provide slugs, approach options, and solutions. The leetgo SQLite DB provides titles, difficulty, categories, and problem descriptions at runtime.
