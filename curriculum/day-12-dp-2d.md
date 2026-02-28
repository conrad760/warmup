# Day 12: Dynamic Programming — 2D, Sequences & Knapsack

> 2 hours | Refresher | The hardest DP day — patterns that appear in real interviews

This is the day where DP stops being about a single array and starts being about
**two dimensions of state**. The make-or-break skill is defining what `dp[i][j]`
means. Get the state definition right and the recurrence writes itself. Get it
wrong and you'll stare at the board for 20 minutes.

---

## Pattern Catalog

### Pattern 1: Two-Sequence Comparison

**Trigger:** Two strings or sequences. Asked about similarity, distance,
alignment, longest common something.

**State definition:** `dp[i][j]` = answer for `s1[0..i-1]` and `s2[0..j-1]`
(the first `i` characters of s1, first `j` of s2).

**Why 1-indexed offset:** Row 0 and column 0 represent the empty prefix. This
gives you clean base cases (`dp[0][j]` = answer when s1 is empty).

**Recurrence (LCS):**
```
if s1[i-1] == s2[j-1]:
    dp[i][j] = dp[i-1][j-1] + 1
else:
    dp[i][j] = max(dp[i-1][j], dp[i][j-1])
```

**Recurrence (Edit Distance):**
```
if s1[i-1] == s2[j-1]:
    dp[i][j] = dp[i-1][j-1]          // no op needed
else:
    dp[i][j] = 1 + min(
        dp[i-1][j-1],                 // replace
        dp[i-1][j],                   // delete from s1
        dp[i][j-1],                   // insert into s1
    )
```

**Go template — Edit Distance:**
```go
func minDistance(word1, word2 string) int {
    m, n := len(word1), len(word2)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }
    // base cases: transforming to/from empty string
    for i := 0; i <= m; i++ {
        dp[i][0] = i
    }
    for j := 0; j <= n; j++ {
        dp[0][j] = j
    }
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if word1[i-1] == word2[j-1] {
                dp[i][j] = dp[i-1][j-1]
            } else {
                dp[i][j] = 1 + min(dp[i-1][j-1], min(dp[i-1][j], dp[i][j-1]))
            }
        }
    }
    return dp[m][n]
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**Complexity:** O(m*n) time, O(m*n) space. Can optimize to O(min(m,n)) space
with two rows.

**Watch out:**
- `dp[i][j]` uses `s1[i-1]` and `s2[j-1]` — the off-by-one is the #1 bug.
- For space optimization, you must save `dp[i-1][j-1]` (the diagonal) before
  overwriting it. Store it in a `prev` variable.

---

### Pattern 2: Grid DP

**Trigger:** 2D grid. Asked about number of paths, minimum cost path, or
reachability from top-left to bottom-right.

**State definition:** `dp[r][c]` = answer (path count or min cost) to reach
cell `(r, c)`.

**Recurrence (Unique Paths):**
```
dp[r][c] = dp[r-1][c] + dp[r][c-1]
```

**Recurrence (Min Path Sum):**
```
dp[r][c] = grid[r][c] + min(dp[r-1][c], dp[r][c-1])
```

**Go template — Unique Paths with Obstacles:**
```go
func uniquePathsWithObstacles(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    if grid[0][0] == 1 || grid[m-1][n-1] == 1 {
        return 0 // blocked start or end
    }
    dp := make([]int, n)
    dp[0] = 1
    for r := 0; r < m; r++ {
        for c := 0; c < n; c++ {
            if grid[r][c] == 1 {
                dp[c] = 0
            } else if c > 0 {
                dp[c] += dp[c-1]
            }
            // dp[c] already has the value from the row above (from previous iteration)
        }
    }
    return dp[n-1]
}
```

**Complexity:** O(m*n) time, O(n) space with single-row optimization.

**Watch out:**
- Check if start or end is blocked before doing any work.
- First row and first column are special: only one direction feeds them. If any
  cell in the first row is blocked, all cells to its right are also unreachable.
- Space optimization works because you only look at current row and row above.

---

### Pattern 3: 0/1 Knapsack

**Trigger:** A set of items, each used at most once. A capacity/target
constraint. Maximize value, count subsets, or check feasibility.

**State definition:** `dp[i][w]` = best answer considering items `0..i-1` with
capacity `w` remaining.

**Recurrence:**
```
dp[i][w] = dp[i-1][w]                           // skip item i-1
if w >= weight[i-1]:
    dp[i][w] = max(dp[i][w], dp[i-1][w-weight[i-1]] + value[i-1])  // take it
```

The key insight: when you take item `i-1`, you look at `dp[i-1][...]` (previous
row), so the item can't be reused.

**Go template — 0/1 Knapsack (space optimized):**
```go
func knapsack01(weights, values []int, capacity int) int {
    dp := make([]int, capacity+1)
    for i := 0; i < len(weights); i++ {
        // RIGHT TO LEFT — prevents using the same item twice
        for w := capacity; w >= weights[i]; w-- {
            dp[w] = max(dp[w], dp[w-weights[i]]+values[i])
        }
    }
    return dp[capacity]
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

**Why right to left:** When you process `w` in decreasing order, `dp[w-weight]`
still holds the value from the *previous* item round (row i-1). If you go left
to right, `dp[w-weight]` may already be updated in the *current* round, meaning
you'd use the item again.

**Complexity:** O(n * W) time, O(W) space.

**Watch out:**
- Direction of inner loop is the most common interview mistake. RIGHT TO LEFT
  for 0/1. Interviewers specifically check for this.
- Subset Sum is knapsack where weight = value = the number itself, and you
  check `dp[target] == true` instead of maximizing.

---

### Pattern 4: Unbounded Knapsack

**Trigger:** Same as knapsack, but items can be reused. Coin change (count
ways), coin change (min coins), rod cutting.

**State definition:** `dp[w]` = answer for capacity/amount `w`, considering
all items.

**Recurrence (Coin Change — min coins):**
```
dp[w] = min over all coins c where c <= w: dp[w-c] + 1
```

**Recurrence (Coin Change II — count ways):**
```
for each coin c:
    for w from c to amount:
        dp[w] += dp[w-c]
```

**Go template — Coin Change II (count combinations):**
```go
func change(amount int, coins []int) int {
    dp := make([]int, amount+1)
    dp[0] = 1 // one way to make amount 0: use nothing
    for _, c := range coins {
        // LEFT TO RIGHT — allows reusing the same coin
        for w := c; w <= amount; w++ {
            dp[w] += dp[w-c]
        }
    }
    return dp[amount]
}
```

**Why left to right:** You *want* `dp[w-c]` to include combinations that
already used coin `c` in this round — that's what "unbounded" means.

**Why outer loop is coins:** Looping coins outside and amounts inside counts
*combinations* (order doesn't matter). Reversing the loop order counts
*permutations* (order matters). Interviewers love asking about this distinction.

**Complexity:** O(n * W) time, O(W) space.

**Watch out:**
- Combinations vs permutations depends on loop order. State this explicitly
  in your interview.
- For "minimum coins" variant, initialize dp to `math.MaxInt32` (not 0), with
  `dp[0] = 0`, and use `min` instead of `+=`.

---

### Pattern 5: Interval DP

**Trigger:** Answer depends on a contiguous subarray or subsequence `[i..j]`.
Merging, splitting, or evaluating ranges.

**State definition:** `dp[i][j]` = answer for the subproblem on elements
`i` through `j`.

**Recurrence (Longest Palindromic Subsequence):**
```
if s[i] == s[j]:
    dp[i][j] = dp[i+1][j-1] + 2
else:
    dp[i][j] = max(dp[i+1][j], dp[i][j-1])
```

**Go template — Longest Palindromic Subsequence:**
```go
func longestPalinSubseq(s string) int {
    n := len(s)
    dp := make([][]int, n)
    for i := range dp {
        dp[i] = make([]int, n)
        dp[i][i] = 1 // single char is a palindrome of length 1
    }
    // iterate by LENGTH of interval, not by starting index
    for length := 2; length <= n; length++ {
        for i := 0; i <= n-length; i++ {
            j := i + length - 1
            if s[i] == s[j] {
                dp[i][j] = dp[i+1][j-1] + 2
            } else {
                dp[i][j] = max(dp[i+1][j], max(dp[i][j-1], 0))
            }
        }
    }
    return dp[0][n-1]
}
```

**Complexity:** O(n^2) time, O(n^2) space.

**Watch out:**
- Loop order matters critically. You must process shorter intervals before
  longer ones. Always loop by `length` from 2 to n, then by starting index `i`.
  If you loop `i` from 0 to n and `j` from i to n, dependencies may not be
  ready.
- Base case: intervals of length 1 (`dp[i][i]`). Sometimes length 0 too.
- Burst Balloons and Matrix Chain Multiplication split at a midpoint `k`:
  `for k := i; k <= j; k++` — try every split.

---

## Decision Framework

```
What are you given?
│
├─ Two strings/sequences?
│  └─ "Similarity / distance / alignment / common subsequence"
│     → Two-Sequence DP: dp[i][j] on prefixes
│
├─ A 2D grid?
│  └─ "Paths / cost / reachability"
│     → Grid DP: dp[r][c] at each cell
│
├─ A set of items + capacity constraint?
│  ├─ Each item used at most once?
│  │  └─ 0/1 Knapsack (inner loop RIGHT TO LEFT)
│  └─ Items reusable?
│     └─ Unbounded Knapsack (inner loop LEFT TO RIGHT)
│
├─ "Partition into two equal subsets"?
│  └─ Subset Sum: knapsack where target = totalSum / 2
│     (if totalSum is odd, answer is immediately false)
│
└─ Answer depends on a contiguous range [i..j]?
   └─ Interval DP: dp[i][j], loop by length
```

**The hidden knapsack:** Many problems don't mention "knapsack" explicitly.
If you see "can you reach a target sum by selecting from a set?" or "count
ways to assign +/- to reach a target" — that's knapsack in disguise.

---

## Common Interview Traps

### Trap 1: 2D Index Alignment
`dp[i][j]` refers to `s1[0..i-1]` and `s2[0..j-1]`. When you access the
characters, use `s1[i-1]` and `s2[j-1]`. Off-by-one here silently produces
wrong answers.

**Fix:** Always write a comment: `// dp[i][j] = answer for s1[:i], s2[:j]`.

### Trap 2: Knapsack Inner Loop Direction
| Variant    | Inner loop     | Why                                       |
|------------|----------------|-------------------------------------------|
| 0/1        | Right to left  | Read from previous row's values           |
| Unbounded  | Left to right  | Read from current row (reuse allowed)     |

Getting this backwards is the single most common knapsack bug. Interviewers
watch for it.

### Trap 3: Space-Optimized Edit Distance
When reducing from 2 rows to 1 row, you overwrite `dp[j-1]` before you need
it as the diagonal. Save it first:
```go
prev := dp[0]  // will be our "diagonal"
dp[0] = i      // base case for this row
for j := 1; j <= n; j++ {
    temp := dp[j]       // save before overwrite
    if word1[i-1] == word2[j-1] {
        dp[j] = prev
    } else {
        dp[j] = 1 + min(prev, min(dp[j], dp[j-1]))
    }
    prev = temp          // old dp[j] becomes next diagonal
}
```

### Trap 4: Grid Obstacles
If `grid[0][0]` or `grid[m-1][n-1]` is blocked, return 0 immediately.
Don't forget to propagate 0s along the first row/column after a blocked cell.

### Trap 5: Interval DP Loop Order
```go
// WRONG: dependencies not ready
for i := 0; i < n; i++ {
    for j := i; j < n; j++ { ... }
}

// CORRECT: shorter intervals computed first
for length := 2; length <= n; length++ {
    for i := 0; i <= n-length; i++ {
        j := i + length - 1
        ...
    }
}
```

---

## Thought Process Walkthroughs

### Walkthrough 1: Edit Distance (LC 72)

> Given two strings word1 and word2, return the minimum number of operations
> (insert, delete, replace) to convert word1 into word2.

**Step 1 — Recognize the pattern.**
Two strings. "Convert one to the other." This is two-sequence DP.

**Step 2 — Define the state.**
`dp[i][j]` = minimum operations to convert `word1[:i]` into `word2[:j]`.

Say this out loud in the interview. It's the most important sentence.

**Step 3 — Base cases.**
- `dp[0][j] = j` — converting empty string to `word2[:j]` needs `j` inserts.
- `dp[i][0] = i` — converting `word1[:i]` to empty string needs `i` deletes.

**Step 4 — Recurrence.**
At each `(i, j)`, look at characters `word1[i-1]` and `word2[j-1]`:
- If they match: `dp[i][j] = dp[i-1][j-1]` (no operation needed).
- If they don't match, take the best of three operations:
  - Replace: `dp[i-1][j-1] + 1`
  - Delete from word1: `dp[i-1][j] + 1`
  - Insert into word1: `dp[i][j-1] + 1`

**Step 5 — Verify with a small example.**
word1 = "ab", word2 = "a"

```
     ""  "a"
""  [ 0,  1 ]
"a" [ 1,  0 ]
"b" [ 2,  1 ]
```

dp[2][1] = 1 (delete 'b'). Correct.

**Step 6 — Code it.** (See template above.)

**Step 7 — Complexity.** O(m*n) time, O(m*n) space. Mention that you can
optimize to O(n) space if needed.

---

### Walkthrough 2: Partition Equal Subset Sum (LC 416)

> Given an integer array nums, return true if you can partition it into two
> subsets with equal sum.

**Step 1 — Recognize the pattern.**
"Partition into two equal subsets" — this is subset sum, which is 0/1 knapsack.

**Step 2 — Reduce to knapsack.**
Total sum must be even (otherwise impossible). Target = totalSum / 2. Now the
question is: can we select a subset that sums to exactly `target`?

**Step 3 — Define the state.**
`dp[w]` = true if we can form sum `w` from items considered so far.

**Step 4 — Base case.**
`dp[0] = true` — sum of 0 is always achievable (empty subset).

**Step 5 — Recurrence (space-optimized 0/1 knapsack).**
For each number `num`, iterate `w` from `target` down to `num`:
```
dp[w] = dp[w] || dp[w-num]
```
Right to left because each number can be used at most once.

**Step 6 — Code it.**
```go
func canPartition(nums []int) bool {
    total := 0
    for _, v := range nums {
        total += v
    }
    if total%2 != 0 {
        return false
    }
    target := total / 2
    dp := make([]bool, target+1)
    dp[0] = true
    for _, num := range nums {
        for w := target; w >= num; w-- {
            dp[w] = dp[w] || dp[w-num]
        }
    }
    return dp[target]
}
```

**Step 7 — Complexity.** O(n * target) time, O(target) space.

**Step 8 — Edge cases.**
- Single element: can't partition, return false.
- Total is odd: return false immediately.
- Any single element > target: skip it (the loop condition handles this).

---

## Practice Problems with Time Targets

### Tier 1: Must-Solve (core patterns, high interview frequency)

| # | Problem | Pattern | Target | LC |
|---|---------|---------|--------|----|
| 1 | Edit Distance | Two-sequence | 20 min | 72 |
| 2 | Longest Common Subsequence | Two-sequence | 15 min | 1143 |
| 3 | Unique Paths II | Grid | 12 min | 63 |
| 4 | Partition Equal Subset Sum | 0/1 Knapsack | 15 min | 416 |
| 5 | Coin Change | Unbounded Knapsack | 12 min | 322 |
| 6 | Coin Change II | Unbounded Knapsack | 15 min | 518 |

### Tier 2: Reinforce (important variants)

| # | Problem | Pattern | Target | LC |
|---|---------|---------|--------|----|
| 7 | Minimum Path Sum | Grid | 10 min | 64 |
| 8 | Target Sum | 0/1 Knapsack | 20 min | 494 |
| 9 | Longest Palindromic Subsequence | Interval | 15 min | 516 |
| 10 | Interleaving String | Two-sequence | 20 min | 97 |

### Tier 3: Stretch (harder, less common)

| # | Problem | Pattern | Target | LC |
|---|---------|---------|--------|----|
| 11 | Burst Balloons | Interval | 30 min | 312 |
| 12 | Distinct Subsequences | Two-sequence | 25 min | 115 |

---

## Quick Drill: State Definition Practice

For each problem below, write the state definition in one sentence before
looking at the answer. This is the single most important skill for 2D DP.

**Problem → State Definition:**

1. **LCS of s1, s2**
   `dp[i][j]` = length of LCS of `s1[:i]` and `s2[:j]`

2. **Edit Distance**
   `dp[i][j]` = min operations to convert `s1[:i]` to `s2[:j]`

3. **Unique Paths in m x n grid**
   `dp[r][c]` = number of paths from `(0,0)` to `(r,c)`

4. **0/1 Knapsack**
   `dp[i][w]` = max value using items `0..i-1` with capacity `w`

5. **Subset Sum (target T)**
   `dp[i][w]` = can we achieve sum `w` using items `0..i-1`? (boolean)

6. **Coin Change (min coins for amount)**
   `dp[w]` = minimum coins needed to make amount `w`

7. **Longest Palindromic Subsequence**
   `dp[i][j]` = length of longest palindromic subsequence in `s[i..j]`

8. **Interleaving String (s1, s2, s3)**
   `dp[i][j]` = can `s1[:i]` and `s2[:j]` interleave to form `s3[:i+j]`?

If you can produce these instantly, you're ready.

---

## Self-Assessment Checklist

### Conceptual (answer before coding)
- [ ] Given a new problem, can you identify which 2D DP pattern applies
      within 2 minutes?
- [ ] Can you write the state definition for any of the 5 patterns without
      looking at notes?
- [ ] Can you explain why 0/1 knapsack iterates right-to-left and unbounded
      iterates left-to-right?
- [ ] Can you explain the loop order for interval DP and why it matters?

### Implementation
- [ ] Edit Distance: full implementation under 20 minutes, correct on first
      run.
- [ ] Partition Equal Subset Sum: recognize as knapsack, implement under 15
      minutes.
- [ ] Can you space-optimize a 2D DP to 1D without bugs? (diagonal save trick)
- [ ] Coin Change II: can you explain what changes if you swap the loop order
      (combinations vs permutations)?

### Interview Readiness
- [ ] Can you talk through your state definition and recurrence out loud while
      coding? (Practice this — silent coding is a red flag for interviewers.)
- [ ] On a problem you haven't seen, can you derive the recurrence from the
      state definition in under 5 minutes?
- [ ] Can you handle follow-ups like "optimize space" or "reconstruct the
      actual solution" (backtrack through the DP table)?

### Scoring
- 10+ checks: Ready for interviews on 2D DP topics.
- 7-9 checks: Solid foundation, practice the gaps.
- <7 checks: Revisit the patterns above and solve Tier 1 problems again.

---

## Session Plan (2 hours)

| Block | Time | Activity |
|-------|------|----------|
| 1 | 0:00-0:15 | Read pattern catalog. Do the state definition drill on paper. |
| 2 | 0:15-0:35 | Solve Edit Distance (LC 72). Talk out loud. |
| 3 | 0:35-0:50 | Solve Partition Equal Subset Sum (LC 416). Spot the knapsack. |
| 4 | 0:50-1:05 | Solve Coin Change II (LC 518). Explain loop order. |
| 5 | 1:05-1:20 | Solve Longest Palindromic Subsequence (LC 516). Nail the loop. |
| 6 | 1:20-1:40 | Pick one from Tier 2 you're weakest on. |
| 7 | 1:40-2:00 | Review mistakes. Redo any problem where you had a bug. |

If a problem takes more than 5 minutes past its target, stop, study the
solution, then re-solve from scratch. Burning time on a stuck attempt teaches
less than understanding and re-implementing.
