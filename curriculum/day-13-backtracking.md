# Day 13: Backtracking

> **Time:** 2 hours | **Level:** Refresher | **Language:** Go

Backtracking = DFS + choose/explore/unchoose. Every problem follows the same skeleton.
The only things that change are: what choices you make, how you explore, and when you prune.

---

## Pattern Catalog

---

### Pattern 1: Subsets

**Trigger:** "Generate all subsets," "power set," "all combinations of any length."

**Template:**

```go
func subsets(nums []int) [][]int {
    var result [][]int
    var path []int

    var backtrack func(start int)
    backtrack = func(start int) {
        // Every node in the recursion tree is a valid subset — record it
        tmp := make([]int, len(path))
        copy(tmp, path)               // CRITICAL: copy the slice before saving
        result = append(result, tmp)

        for i := start; i < len(nums); i++ {
            path = append(path, nums[i])  // --- CHOOSE ---
            backtrack(i + 1)              // --- EXPLORE --- (i+1: no reuse)
            path = path[:len(path)-1]     // --- UNCHOOSE ---
        }
    }

    backtrack(0)
    return result
}
```

**Complexity:** O(n * 2^n) — 2^n subsets, each up to length n to copy.

**Watch out:**
- `copy` before appending to `result`. Without it, later recursive calls mutate saved slices.
- `start` parameter prevents generating duplicate subsets like {1,2} and {2,1}.

---

### Pattern 2: Permutations

**Trigger:** "Generate all permutations," "all orderings," "arrange n elements."

**Template A — `used[]` boolean array (clearer, preferred in interviews):**

```go
func permute(nums []int) [][]int {
    var result [][]int
    path := make([]int, 0, len(nums))
    used := make([]bool, len(nums))

    var backtrack func()
    backtrack = func() {
        if len(path) == len(nums) {
            tmp := make([]int, len(path))
            copy(tmp, path)
            result = append(result, tmp)
            return
        }

        for i := 0; i < len(nums); i++ {
            if used[i] {
                continue              // skip already-chosen elements
            }
            used[i] = true            // --- CHOOSE ---
            path = append(path, nums[i])

            backtrack()               // --- EXPLORE ---

            path = path[:len(path)-1] // --- UNCHOOSE ---
            used[i] = false
        }
    }

    backtrack()
    return result
}
```

**Template B — swap-based (in-place, saves memory):**

```go
func permuteSwap(nums []int) [][]int {
    var result [][]int

    var backtrack func(first int)
    backtrack = func(first int) {
        if first == len(nums) {
            tmp := make([]int, len(nums))
            copy(tmp, nums)
            result = append(result, tmp)
            return
        }

        for i := first; i < len(nums); i++ {
            nums[first], nums[i] = nums[i], nums[first] // --- CHOOSE ---
            backtrack(first + 1)                         // --- EXPLORE ---
            nums[first], nums[i] = nums[i], nums[first] // --- UNCHOOSE ---
        }
    }

    backtrack(0)
    return result
}
```

**Complexity:** O(n * n!) — n! permutations, each costs O(n) to copy.

**Watch out:**
- Swap-based does NOT handle duplicates cleanly. Use the `used[]` approach for permutations with duplicates.
- The loop starts at `0` every time (not `start`) — that is the key difference from subsets.

---

### Pattern 3: Combinations (choose k from n)

**Trigger:** "All combinations of size k," "choose k elements from n."

**Template:**

```go
func combine(n, k int) [][]int {
    var result [][]int
    var path []int

    var backtrack func(start int)
    backtrack = func(start int) {
        if len(path) == k {
            tmp := make([]int, k)
            copy(tmp, path)
            result = append(result, tmp)
            return                    // no need to explore further
        }

        // PRUNING: need (k - len(path)) more elements, so stop early
        // if there aren't enough elements remaining
        for i := start; i <= n-(k-len(path))+1; i++ {
            path = append(path, i)    // --- CHOOSE ---
            backtrack(i + 1)          // --- EXPLORE ---
            path = path[:len(path)-1] // --- UNCHOOSE ---
        }
    }

    backtrack(1)
    return result
}
```

**Complexity:** O(k * C(n,k)).

**Watch out:**
- The pruning bound `i <= n-(k-len(path))+1` is a common follow-up optimization interviewers ask about. Know how to derive it: you need `k - len(path)` more picks, and there are `n - i + 1` candidates left.

---

### Pattern 4: Subsets / Combinations with Duplicates

**Trigger:** "Input may contain duplicates," "unique subsets/combinations."

**Template (Subsets II):**

```go
func subsetsWithDup(nums []int) [][]int {
    sort.Ints(nums)                   // STEP 1: sort so duplicates are adjacent
    var result [][]int
    var path []int

    var backtrack func(start int)
    backtrack = func(start int) {
        tmp := make([]int, len(path))
        copy(tmp, path)
        result = append(result, tmp)

        for i := start; i < len(nums); i++ {
            // STEP 2: skip duplicate at same recursion level
            //   i > start  (NOT i > 0)
            //   This allows picking the same value deeper in the tree,
            //   but prevents picking it again at the same branching level.
            if i > start && nums[i] == nums[i-1] {
                continue
            }

            path = append(path, nums[i])  // --- CHOOSE ---
            backtrack(i + 1)              // --- EXPLORE ---
            path = path[:len(path)-1]     // --- UNCHOOSE ---
        }
    }

    backtrack(0)
    return result
}
```

**Why `i > start` and not `i > 0`?**

Consider `nums = [1, 2, 2]`. At `start=1`, both `nums[1]` and `nums[2]` are `2`.
- `i=1` (first `2`): we MUST allow this — it's the first time we see `2` at this level.
- `i=2` (second `2`): `i > start` is true AND `nums[2] == nums[1]` — SKIP.

If we used `i > 0`, we'd skip `nums[1]` when `start=1`, which is wrong — that would
prevent `{1, 2}` from ever being generated.

**Complexity:** O(n * 2^n) worst case.

**Watch out:**
- You MUST sort first. The skip condition only works on sorted input.
- Same skip pattern works for Combination Sum II (choose elements with duplicates, target sum).

---

### Pattern 5: Constraint Satisfaction (N-Queens, Sudoku)

**Trigger:** "Place N items with constraints," "fill a grid following rules," "is there a valid arrangement."

**Template (N-Queens):**

```go
func solveNQueens(n int) [][]string {
    var result [][]string
    board := make([][]byte, n)
    for i := range board {
        board[i] = make([]byte, n)
        for j := range board[i] {
            board[i][j] = '.'
        }
    }

    // Track which columns and diagonals are under attack
    cols := make([]bool, n)
    diag1 := make([]bool, 2*n)       // row - col + n (anti-diagonal)
    diag2 := make([]bool, 2*n)       // row + col (main diagonal)

    var backtrack func(row int)
    backtrack = func(row int) {
        if row == n {
            // Record solution
            solution := make([]string, n)
            for i := range board {
                solution[i] = string(board[i])
            }
            result = append(result, solution)
            return
        }

        for col := 0; col < n; col++ {
            // --- PRUNE: check constraints before choosing ---
            if cols[col] || diag1[row-col+n] || diag2[row+col] {
                continue
            }

            board[row][col] = 'Q'         // --- CHOOSE ---
            cols[col] = true
            diag1[row-col+n] = true
            diag2[row+col] = true

            backtrack(row + 1)            // --- EXPLORE ---

            board[row][col] = '.'         // --- UNCHOOSE ---
            cols[col] = false
            diag1[row-col+n] = false
            diag2[row+col] = false
        }
    }

    backtrack(0)
    return result
}
```

**Complexity:** O(n!) upper bound (much less in practice due to pruning).

**Watch out:**
- The key insight is encoding diagonals as `row-col` and `row+col`. Practice deriving this.
- Check constraints BEFORE recursing (prune early), not after placing.
- Constraint satisfaction problems often ask for ONE solution. Add a `found` boolean and return early.

---

### Pattern 6: Word Search / Grid Backtracking

**Trigger:** "Find a word in a grid," "find a path in a matrix," "DFS on grid with backtracking."

**Template (LC 79 — Word Search):**

```go
func exist(board [][]byte, word string) bool {
    rows, cols := len(board), len(board[0])
    dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

    var backtrack func(r, c, idx int) bool
    backtrack = func(r, c, idx int) bool {
        if idx == len(word) {
            return true               // matched entire word
        }
        // Bounds check + character match
        if r < 0 || r >= rows || c < 0 || c >= cols {
            return false
        }
        if board[r][c] != word[idx] {
            return false
        }

        saved := board[r][c]          // --- CHOOSE: mark visited ---
        board[r][c] = '#'

        for _, d := range dirs {      // --- EXPLORE: try all 4 directions ---
            if backtrack(r+d[0], c+d[1], idx+1) {
                return true           // early exit on first match
            }
        }

        board[r][c] = saved           // --- UNCHOOSE: restore cell ---
        return false
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if backtrack(r, c, 0) {
                return true
            }
        }
    }
    return false
}
```

**Complexity:** O(m * n * 4^L) where L = word length, m*n = grid size.

**Watch out:**
- Mutate the board in-place to mark visited (set to `'#'`), then restore. Faster and simpler than a separate `visited` matrix.
- Don't forget to restore the cell even on the path that returns `true` if you need to continue searching for other matches. Here we return immediately so it's fine.
- Start the search from every cell — don't just start from `(0,0)`.

---

## Decision Framework

```
What does the problem ask for?
│
├─ "All subsets / power set"
│   └─ Pattern 1: Subsets
│       └─ Has duplicates? → Pattern 4: sort + skip
│
├─ "All permutations / orderings"
│   └─ Pattern 2: Permutations
│       └─ Has duplicates? → sort + skip with used[] array
│
├─ "All combinations of size k"
│   └─ Pattern 3: Combinations
│       └─ Has duplicates? → Pattern 4: sort + skip
│       └─ Elements reusable? → recurse with i (not i+1)
│
├─ "Place items with rules" / "Fill grid"
│   └─ Pattern 5: Constraint Satisfaction
│
└─ "Find path/word in grid"
    └─ Pattern 6: Grid Backtracking
```

**Key differences between patterns:**

| Aspect          | Subsets        | Permutations     | Combinations     |
|-----------------|----------------|------------------|------------------|
| Loop start      | `i = start`    | `i = 0`          | `i = start`      |
| Recurse with    | `i + 1`        | any unused       | `i + 1`          |
| Base case       | every node     | `len == n`       | `len == k`       |
| Result count    | 2^n            | n!               | C(n,k)           |

---

## Common Interview Traps

### 1. The Go Slice Append Gotcha

This is the #1 backtracking bug in Go. Understand it cold.

```go
// WRONG — result entries get silently corrupted
result = append(result, path)

// WHY: append may return a slice sharing the same underlying array.
// Later mutations to path overwrite what you already saved.

// RIGHT — always copy
tmp := make([]int, len(path))
copy(tmp, path)
result = append(result, tmp)
```

If an interviewer asks why your subsets solution returns wrong answers, this is almost certainly it.

### 2. Forgetting to Undo the Choice

Every CHOOSE must have a matching UNCHOOSE. Checklist:
- `path = path[:len(path)-1]` — undo append
- `used[i] = false` — undo marking
- `board[r][c] = saved` — undo grid mutation
- `cols[col] = false` — undo constraint flag

If your solution produces too few results, you probably forgot to undo somewhere.

### 3. Duplicate Skipping: `i > start`, Not `i > 0`

```go
// WRONG — skips valid first picks
if i > 0 && nums[i] == nums[i-1] { continue }

// RIGHT — only skips duplicates at the same branching level
if i > start && nums[i] == nums[i-1] { continue }
```

### 4. Permutations with Duplicates

Use `used[]` approach (not swap), sort first, then:

```go
// Skip if: same value as previous, AND previous was not used
// (meaning previous was un-chosen at this level, so this is a duplicate branch)
if i > 0 && nums[i] == nums[i-1] && !used[i-1] {
    continue
}
```

### 5. Complexity Anxiety

Backtracking is inherently exponential. That's expected. Don't waste time trying to
make it polynomial. The interviewer wants to see:
- Clean structure (choose/explore/unchoose)
- Correct pruning (cut bad branches early)
- Proper handling of duplicates

---

## Thought Process Walkthrough

### Walkthrough 1: Subsets (LC 78)

> Given `nums = [1, 2, 3]`, return all subsets.

**Step 1 — Identify the pattern.**
"All subsets" maps directly to Pattern 1.

**Step 2 — State the recursion tree in words.**
"At each position `i`, I branch: either include `nums[i]` or skip to `nums[i+1]`.
Every node is a valid subset, so I record at the top of each call."

**Step 3 — Write the skeleton, then fill in.**

```
backtrack(start):
    record current path           ← every node counts
    for i = start to n-1:
        CHOOSE:   add nums[i]
        EXPLORE:  backtrack(i+1)
        UNCHOOSE: remove last
```

**Step 4 — Trace through the example.**

```
backtrack(0), path=[]
  record []
  i=0: choose 1 → backtrack(1), path=[1]
    record [1]
    i=1: choose 2 → backtrack(2), path=[1,2]
      record [1,2]
      i=2: choose 3 → backtrack(3), path=[1,2,3]
        record [1,2,3]
        (loop ends)
      unchoose 3
      (loop ends)
    unchoose 2
    i=2: choose 3 → backtrack(3), path=[1,3]
      record [1,3]
      (loop ends)
    unchoose 3
  unchoose 1
  i=1: choose 2 → backtrack(2), path=[2]
    record [2]
    i=2: choose 3 → backtrack(3), path=[2,3]
      record [2,3]
    unchoose 3
  unchoose 2
  i=2: choose 3 → backtrack(3), path=[3]
    record [3]
  unchoose 3

Result: [], [1], [1,2], [1,2,3], [1,3], [2], [2,3], [3]   ✓ (8 = 2^3)
```

**Step 5 — Verify complexity.**
2^3 = 8 subsets. Each copy costs O(n). Total: O(n * 2^n). Correct.

---

### Walkthrough 2: Word Search (LC 79)

> Given board and word = "ABCCED", determine if the word exists in the grid.

**Step 1 — Identify the pattern.**
"Find a word in a grid by moving to adjacent cells" → Pattern 6: grid backtracking.

**Step 2 — State the approach.**
"For each cell matching `word[0]`, start a DFS. At each step, check if the current cell
matches `word[idx]`. If so, mark visited and recurse in 4 directions. Backtrack by
restoring the cell."

**Step 3 — Handle the tricky parts.**

*Visited tracking:* Overwrite `board[r][c]` with `'#'` instead of a separate matrix.
Saves space and is idiomatic in interviews.

*Early termination:* Return `true` as soon as any path matches. Don't explore all paths.

*Starting point:* Must try every cell as a starting point, not just `(0,0)`.

**Step 4 — Trace a short example.**

```
Board:       Word: "AB"
A B
C D

Start at (0,0): board[0][0]='A' == word[0] ✓
  Mark (0,0) as '#'
  Try right → (0,1): board[0][1]='B' == word[1] ✓
    idx+1 == len(word) → return true
  (would restore (0,0) but we returned early)
```

**Step 5 — Edge cases to mention.**
- Empty word → return true
- Word longer than total cells → return false (quick check)
- Single cell board
- Word uses same cell twice (the visited marking prevents this)

---

## Time Targets

| Problem                     | LC # | Target | Notes                            |
|-----------------------------|------|--------|----------------------------------|
| Subsets                     | 78   | 5 min  | Base template, must be instant   |
| Subsets II                  | 90   | 7 min  | Sort + skip                      |
| Permutations                | 46   | 5 min  | used[] approach                  |
| Permutations II             | 47   | 8 min  | Sort + skip + used[]             |
| Combinations                | 77   | 5 min  | Subset variant with len check    |
| Combination Sum             | 39   | 8 min  | Recurse with `i` (reuse allowed) |
| Combination Sum II          | 40   | 8 min  | Sort + skip, recurse with `i+1`  |
| Word Search                 | 79   | 10 min | Grid backtracking                |
| N-Queens                    | 51   | 15 min | Constraint satisfaction          |
| Letter Combinations (Phone) | 17   | 5 min  | Simple backtrack over mapping    |

---

## Quick Drill (30 minutes)

Do these without looking at templates. Time yourself.

1. **Subsets** (LC 78) — Write from scratch. Target: 5 minutes.
2. **Permutations** (LC 46) — Both `used[]` and swap approaches. Target: 5 min each.
3. **Subsets II** (LC 90) — Focus on getting the skip condition right on the first try.
4. **Word Search** (LC 79) — Focus on the visited marking and restoration.

After each one, check:
- Did I copy the slice before saving?
- Did every CHOOSE have an UNCHOOSE?
- Is my skip condition `i > start` (not `i > 0`)?

---

## Self-Assessment

### Can I explain these from memory?

| Question                                                        | Confident? |
|-----------------------------------------------------------------|------------|
| Why must I `copy` a slice before appending to results in Go?    |            |
| What is the difference between subsets and permutations loops?   |            |
| Why is the skip condition `i > start` and not `i > 0`?          |            |
| How do I handle permutations with duplicates?                   |            |
| How do I encode diagonal attacks for N-Queens?                  |            |
| What is the time complexity of generating all subsets?           |            |
| Why is backtracking inherently exponential?                     |            |
| How does grid backtracking differ from tree backtracking?        |            |

### Red flags that you need more practice:
- You wrote a backtracking solution and got wrong answers → likely forgot copy or unchoose.
- You can't write Subsets from memory in under 5 minutes.
- You confuse when to loop from `0` (permutations) vs `start` (subsets/combinations).
- You can't articulate WHY the duplicate skip works.

### Green lights — you're ready:
- You can write all 6 templates from memory with correct choose/explore/unchoose.
- You can explain the slice copy gotcha to someone else.
- You can modify a base template for variants (with duplicates, with reuse, with constraints) without hesitation.
- You hit the time targets above consistently.
