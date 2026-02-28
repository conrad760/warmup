# Day 17 — Dynamic Programming: 2D & Grids

---

## 1. Curated Learning Resources

| # | Resource | Format | Why This One |
|---|----------|--------|-------------|
| 1 | [NeetCode — Unique Paths](https://www.youtube.com/watch?v=IlEsdxuD4lY) | Video (10 min) | The cleanest introduction to grid DP. Draws the grid, fills every cell, and makes the recurrence obvious: each cell is the sum of the cell above and to the left. Shows the 1D space optimization at the end. Watch this first — it's the simplest grid DP and sets the pattern for everything else today. |
| 2 | [NeetCode — 0/1 Knapsack](https://www.youtube.com/watch?v=cjWnW0hdF1Y) | Video (18 min) | Builds the knapsack recurrence from scratch with a clear 2D table. The key moment: he shows why the 1D optimization requires iterating capacity right-to-left — if you go left-to-right, you accidentally use the same item twice. This single insight is the entire difference between 0/1 and unbounded knapsack. |
| 3 | [Abdul Bari — 0/1 Knapsack (Dynamic Programming)](https://www.youtube.com/watch?v=nLmhmB6NzcM) | Video (21 min) | Whiteboard-style, slower pace. Excellent at explaining the "take or skip" decision at each item. He fills the full 2D table cell by cell, which makes the logic impossible to miss. Good if the NeetCode version moves too fast. |
| 4 | [Back To Back SWE — Partition Equal Subset Sum](https://www.youtube.com/watch?v=s6FhG--P7z0) | Video (15 min) | The best walkthrough of reducing partition to subset sum to boolean knapsack. Shows the chain of reductions clearly: "can we partition?" becomes "does a subset sum to total/2?" becomes "boolean knapsack with capacity = total/2." This connection is the key insight for interview day. |
| 5 | [VisuAlgo — Knapsack Visualization](https://visualgo.net/en/knapsack) | Interactive | Step through the knapsack DP table cell by cell. Toggle between 0/1 and unbounded variants. Watch how the cell dependencies change between the two — this makes the iteration direction difference visceral. Use during the review block and again during knapsack implementation. |
| 6 | [Aditya Verma — 0/1 Knapsack Space Optimization](https://www.youtube.com/watch?v=ntCGbPMeqgg) | Video (14 min) | Dedicated walkthrough of collapsing the 2D knapsack table to 1D. He draws both the 2D table and the 1D array side-by-side and traces through exactly why right-to-left iteration preserves the "previous row" values. Watch after you implement the full 2D version. |
| 7 | [Techdose — Minimum Path Sum](https://www.youtube.com/watch?v=lBRtnuxg-gU) | Video (12 min) | Clean grid DP walkthrough for minimum path sum. Shows both the 2D approach and the in-place modification trick (use the input grid as the DP table). Good for reinforcing the pattern: grid DP is always about "which cells can I come FROM?" |
| 8 | [MIT 6.006 — DP IV: Guitar Fingering, Tetris, Knapsack](https://www.youtube.com/watch?v=Tw1k46ywN6E) | Video (50 min) | Erik Demaine covers knapsack as a pseudopolynomial algorithm — O(nW) where W is the capacity, not the input size. This theoretical perspective explains why knapsack is NP-hard despite having a "polynomial-looking" DP. Watch the knapsack section (starts ~25:00) if you want the deeper understanding. |

**Suggested order:** Resource 1 (NeetCode unique paths) first to establish the grid DP pattern. Resource 7 (min path sum) during grid implementation. Then pivot: Resource 2 (NeetCode knapsack) for the 0/1 formulation, Resource 6 (Aditya Verma) after your 2D knapsack works, and Resource 4 (Back To Back SWE) before implementing subset sum / partition. Keep Resource 5 (VisuAlgo) open throughout.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 - 12:20 — Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:06 | Read Section 3 (Core Concepts) end to end. Focus on the two paradigms today: grid DP (the grid IS the DP table) and knapsack DP (items x capacity). Internalize the core question for grid DP: "which cells can I arrive FROM?" and for knapsack: "do I take this item or skip it?" |
| 12:06 - 12:12 | **Hand-trace the unique paths grid for a 4x3 grid on paper.** Fill every cell. Notice: the first row is all 1s, the first column is all 1s, and every other cell is the sum of the cell above and the cell to the left. Then trace a 4x3 min path sum grid with values `[[1,3,1],[1,5,1],[4,2,1]]`. Write the min-cost at each cell. |
| 12:12 - 12:17 | **Hand-trace the 0/1 knapsack table** for items `[(weight=1, value=1), (weight=3, value=4), (weight=4, value=5), (weight=5, value=7)]` with capacity W=7. Fill the 2D table `dp[i][w]`. At each cell, ask: "take or skip?" Circle the cells where "take" wins. |
| 12:17 - 12:20 | **Read the Knapsack Family section (Section 5).** Map the relationships in your head: 0/1 knapsack -> subset sum (value = weight) -> partition (target = total/2). Unbounded knapsack -> coin change. The recurrences are nearly identical — only the iteration direction changes. |

### 12:20 - 1:20 — Implement

| Time | Activity |
|------|----------|
| 12:20 - 12:30 | **Unique Paths — full 2D table, then 1D optimized.** Write `UniquePaths(m, n int) int` with a full `m x n` table. Then write `UniquePathsOptimized(m, n int) int` using a single 1D row of length `n`. Test with: `(3,7)` -> 28, `(3,2)` -> 3, `(1,1)` -> 1, `(1,5)` -> 1, `(5,1)` -> 1. Verify both versions agree. |
| 12:30 - 12:40 | **Min Path Sum.** Write `MinPathSum(grid [][]int) int`. Use the recurrence `dp[r][c] = grid[r][c] + min(dp[r-1][c], dp[r][c-1])`. Handle first row and first column separately (only one direction to come from). Test with: `[[1,3,1],[1,5,1],[4,2,1]]` -> 7, `[[1,2,3]]` -> 6 (single row), `[[1],[2],[3]]` -> 6 (single column), `[[5]]` -> 5 (single cell). |
| 12:40 - 12:48 | **Unique Paths with Obstacles.** Write `UniquePathsWithObstacles(grid [][]int) int`. Obstacle cells (`grid[r][c] == 1`) have `dp[r][c] = 0`. Check: if start or end is blocked, return 0 immediately. In the first row/column, once you hit an obstacle, everything after it is also 0. Test with: `[[0,0,0],[0,1,0],[0,0,0]]` -> 2, `[[1,0]]` -> 0 (start blocked), `[[0,1]]` -> 0 (end blocked), `[[0]]` -> 1. |
| 12:48 - 13:00 | **0/1 Knapsack — full 2D table.** Write `Knapsack01(weights, values []int, W int) int`. Build the `(n+1) x (W+1)` table. For each item `i` and capacity `w`: skip if `weights[i-1] > w`, otherwise `dp[i][w] = max(dp[i-1][w], dp[i-1][w-weights[i-1]] + values[i-1])`. Test with: items `[(1,1),(3,4),(4,5),(5,7)]`, W=7 -> 9 (take items 2 and 3, or 1 and 4). |
| 13:00 - 13:08 | **0/1 Knapsack — 1D space optimized.** Write `Knapsack01Optimized(weights, values []int, W int) int`. Single array of size `W+1`. For each item, iterate `w` from `W` DOWN to `weights[i]`. This right-to-left traversal is the critical detail. Verify results match the 2D version on all test cases. |
| 13:08 - 13:15 | **Subset Sum.** Write `SubsetSum(nums []int, target int) bool`. This is boolean 0/1 knapsack: `dp[w]` is `true/false`. Initialize `dp[0] = true`. For each number, iterate right-to-left: `dp[w] = dp[w] || dp[w-nums[i]]`. Test with: `([3,34,4,12,5,2], 9)` -> true (4+5), `([3,34,4,12,5,2], 30)` -> false, `([1,2,3], 6)` -> true, `([1], 2)` -> false. |
| 13:15 - 13:20 | **Partition Equal Subset Sum.** Write `CanPartition(nums []int) bool`. Sum all elements. If sum is odd, return false. Otherwise, call `SubsetSum(nums, sum/2)`. Test with: `[1,5,11,5]` -> true (1+5+5=11), `[1,2,3,5]` -> false, `[1,1]` -> true, `[1]` -> false, `[2,2,1,1]` -> true. |

### 1:20 - 1:50 — Solidify

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | **Edge cases for all problems.** Single-cell grids. Single-row and single-column grids. Grid entirely filled with obstacles (except start/end). Knapsack with zero capacity. Knapsack with items heavier than capacity. Subset sum with target 0 (always true — empty subset). Empty nums array. Partition with all-zero elements. |
| 1:30 - 1:40 | **Space-optimize min path sum.** Modify your MinPathSum to use a single 1D array of size `cols`. Process row by row: for the first cell of each row, add from above. For other cells, `dp[c] = grid[r][c] + min(dp[c], dp[c-1])`. Alternatively, modify the input grid in-place (if allowed). Verify results match. |
| 1:40 - 1:50 | **Think about unbounded knapsack.** Take your `Knapsack01Optimized` and change the inner loop to iterate LEFT-TO-RIGHT instead of right-to-left. This gives unbounded knapsack (each item usable multiple times). Test: with items `[(1,1),(3,4),(4,5),(5,7)]`, W=7, unbounded should give 10 (two items of weight 3 and one of weight 1). Recognize: coin change is unbounded knapsack where value = 1 and you minimize. |

### 1:50 - 2:00 — Recap (From Memory)

Write down without looking:

1. The unique paths recurrence and base cases.
2. Why the 1D optimization works for grid DP (what does each cell in the row represent before and after updating?).
3. The 0/1 knapsack recurrence (take-or-skip).
4. Why the 1D knapsack iterates right-to-left (0/1) vs. left-to-right (unbounded).
5. How subset sum is a boolean knapsack. How partition reduces to subset sum.
6. One gotcha you hit during implementation.

---

## 3. Core Concepts Deep Dive

### Grid DP: The Cell Is the State

In grid DP, the grid coordinates **are** the DP state. `dp[r][c]` represents the answer (number of paths, minimum cost, maximum value) for the subproblem "from the start to cell (r, c)."

The pattern is always:

> **Look at cell (r, c). Ask: "which cells can I arrive FROM?" The recurrence combines those cells.**

For a grid where you can only move right or down:
- You can arrive at `(r, c)` from `(r-1, c)` (came from above) or `(r, c-1)` (came from the left).
- The recurrence combines these two sources — addition for counting paths, min for shortest path, etc.

**Movement constraints define the recurrence.** If you could also move diagonally, you'd add `dp[r-1][c-1]` to the recurrence. If you could move in all four directions, it's no longer simple DP — you'd need BFS/DFS instead (because of cycles).

**Boundary conditions are the base cases.** The first row can only be reached from the left. The first column can only be reached from above. These are your base cases, and they often have trivial values (all 1s for unique paths, cumulative sums for min path sum).

This is one of the most visual DP patterns. You can literally watch the table fill in, cell by cell, left-to-right, top-to-bottom. There's nothing abstract — every cell is a physical location on the grid.

---

### Why Grid DP Is So Visual

Most DP problems require you to imagine an abstract state space. "What does `dp[i]` mean?" is often the hardest question. With grid DP, you don't have to imagine anything — the state space is the grid sitting right in front of you.

Consider unique paths on a 3x4 grid:

```
  ┌───┬───┬───┬───┐
  │ 1 │ 1 │ 1 │ 1 │   First row: only one way (go right)
  ├───┼───┼───┼───┤
  │ 1 │ 2 │ 3 │ 4 │   Each cell = above + left
  ├───┼───┼───┼───┤
  │ 1 │ 3 │ 6 │10 │   Answer is bottom-right: 10
  └───┴───┴───┴───┘
```

You can SEE the table being filled. You can POINT to the two cells that each value came from. This visual immediacy is why grid DP is typically the first 2D DP pattern taught — it builds intuition that transfers to more abstract 2D problems like knapsack.

---

### 0/1 Knapsack: The Classic Formulation

You have `n` items, each with a `weight` and a `value`. You have a knapsack with capacity `W`. Maximize total value without exceeding capacity. Each item can be used at most once.

**State:** `dp[i][w]` = maximum value achievable using the first `i` items with capacity `w`.

**Recurrence — the take-or-skip decision:**

```
if weights[i-1] > w:
    dp[i][w] = dp[i-1][w]                              // can't take: too heavy
else:
    dp[i][w] = max(
        dp[i-1][w],                                     // SKIP item i
        dp[i-1][w - weights[i-1]] + values[i-1]        // TAKE item i
    )
```

**Base case:** `dp[0][w] = 0` for all `w` (no items, no value). `dp[i][0] = 0` for all `i` (no capacity, no value).

**Why O(n * W):** We fill an `(n+1) x (W+1)` table, each cell in O(1). This is pseudopolynomial — polynomial in `n` and `W`, but `W` is a *value*, not a count. If `W` is encoded in `log(W)` bits, the algorithm is actually exponential in input size. This is why knapsack is NP-hard despite having a "polynomial" DP.

**The key insight:** When deciding about item `i`, the "skip" option looks at `dp[i-1][w]` (same capacity, one fewer item), and the "take" option looks at `dp[i-1][w - weight]` (reduced capacity, one fewer item). Both reference row `i-1`. This is what makes the 1D space optimization possible — and what determines its iteration direction.

---

### Unbounded Knapsack: The One-Line Change

Same setup, but each item can be used unlimited times.

The recurrence changes subtly:

```
dp[i][w] = max(
    dp[i-1][w],                                     // SKIP item i entirely
    dp[i][w - weights[i-1]] + values[i-1]           // TAKE item i (may take again!)
)
```

Notice: the "take" option references `dp[i][...]` (same row), not `dp[i-1][...]` (previous row). After taking item `i`, you're still in the "considering item `i`" row — because you can take it again.

**Coin change is unbounded knapsack.** "Minimum coins to make amount A" is unbounded knapsack where:
- Capacity = target amount
- Item weight = coin denomination
- Item value = 1 (each coin costs 1 to use)
- Objective = minimize instead of maximize

Same recurrence, different bookkeeping.

---

### Subset Sum: Boolean Knapsack

Subset sum asks: given a set of positive integers and a target, does any subset sum to exactly the target?

This is 0/1 knapsack where value = weight. But you don't maximize — you just need a boolean: is the target achievable?

**State:** `dp[w]` = `true` if some subset of the numbers seen so far sums to exactly `w`.

**Recurrence:** `dp[w] = dp[w] || dp[w - nums[i]]`

**Base case:** `dp[0] = true` (the empty subset sums to 0).

Iterate right-to-left (it's 0/1 — each number used at most once).

---

### Partition Equal Subset Sum: One Reduction Away

"Can you partition an array into two subsets with equal sum?"

**Reduction:**
1. Compute `total = sum(nums)`.
2. If `total` is odd, return `false` (can't split an odd number into two equal integers).
3. Return `SubsetSum(nums, total/2)`.

Why? If one subset sums to `total/2`, the remaining elements must also sum to `total/2`. Finding one subset is sufficient.

This chain of reductions — partition -> subset sum -> boolean knapsack — is exactly the kind of connection interviewers love. Each step is trivial once you see it, but spotting the chain is the skill.

---

### Space Optimization: The Direction Rule

**Grid DP — 1D row:**

Since each cell depends only on the current row (left neighbor) and the previous row (above neighbor), you can use a single 1D array. When you process left-to-right:
- `dp[c]` before updating: holds the value from the previous row (the "above" value).
- `dp[c-1]` after updating: holds the value from the current row (the "left" value).
- After `dp[c] += dp[c-1]`, it now holds the correct value for the current row.

This works because the two dependencies (above and left) are conveniently available: "above" is the old value (not yet overwritten), "left" is the new value (just computed).

**0/1 Knapsack — 1D with REVERSE (right-to-left) iteration:**

```go
dp := make([]int, W+1)
for i := 0; i < n; i++ {
    for w := W; w >= weights[i]; w-- {   // RIGHT TO LEFT
        dp[w] = max(dp[w], dp[w-weights[i]] + values[i])
    }
}
```

Why right-to-left? When computing `dp[w]`, you need `dp[w - weight[i]]` from the PREVIOUS item's row. If you iterate left-to-right, `dp[w - weight[i]]` might already be updated (it has a smaller index, so you'd have processed it first). That updated value reflects "took item `i`" — so using it again means taking item `i` twice. Right-to-left ensures `dp[w - weight[i]]` still holds the previous row's value.

**Unbounded Knapsack — 1D with FORWARD (left-to-right) iteration:**

```go
dp := make([]int, W+1)
for i := 0; i < n; i++ {
    for w := weights[i]; w <= W; w++ {   // LEFT TO RIGHT
        dp[w] = max(dp[w], dp[w-weights[i]] + values[i])
    }
}
```

Why left-to-right? You WANT `dp[w - weight[i]]` to reflect the current row (including the possibility of having already taken item `i`). Left-to-right gives you exactly that — the "double counting" that was a bug in 0/1 is a feature in unbounded.

**This is the single most important insight for the knapsack family.** The only code difference between 0/1 and unbounded is the loop direction. Everything else — the recurrence, the base case, the initialization — is identical.

---

## 4. Implementation Checklist

### Function Signatures

```go
package dp

// --- Grid DP ---

// UniquePaths returns the number of unique paths from the top-left to the
// bottom-right of an m x n grid, moving only right or down.
// Time: O(m*n), Space: O(m*n)
func UniquePaths(m, n int) int

// UniquePathsOptimized returns the same result using a single 1D row.
// Time: O(m*n), Space: O(n)
func UniquePathsOptimized(m, n int) int

// MinPathSum returns the minimum sum path from top-left to bottom-right
// in a grid of non-negative integers, moving only right or down.
// Time: O(m*n), Space: O(m*n) or O(n) with 1D optimization
func MinPathSum(grid [][]int) int

// UniquePathsWithObstacles returns the number of unique paths in a grid
// where 1 = obstacle and 0 = free. Returns 0 if start or end is blocked.
// Time: O(m*n), Space: O(m*n) or O(n) with 1D optimization
func UniquePathsWithObstacles(grid [][]int) int

// --- Knapsack Family ---

// Knapsack01 returns the maximum value achievable with the given items
// (each used at most once) and capacity W. Full 2D table.
// Time: O(n*W), Space: O(n*W)
func Knapsack01(weights, values []int, W int) int

// Knapsack01Optimized returns the same result using a single 1D array
// with right-to-left iteration.
// Time: O(n*W), Space: O(W)
func Knapsack01Optimized(weights, values []int, W int) int

// SubsetSum returns true if any subset of nums sums to exactly target.
// Uses boolean 0/1 knapsack with right-to-left iteration.
// Time: O(n*target), Space: O(target)
func SubsetSum(nums []int, target int) bool

// CanPartition returns true if nums can be partitioned into two subsets
// with equal sum. Reduces to SubsetSum(nums, total/2).
// Time: O(n*sum/2), Space: O(sum/2)
func CanPartition(nums []int) bool
```

### Test Cases & Edge Cases

**UniquePaths / UniquePathsOptimized:**

| Input | Expected | Notes |
|-------|----------|-------|
| `(3, 7)` | `28` | Standard case |
| `(3, 2)` | `3` | Small grid |
| `(1, 1)` | `1` | Single cell |
| `(1, 5)` | `1` | Single row — only one path (all right) |
| `(5, 1)` | `1` | Single column — only one path (all down) |
| `(2, 2)` | `2` | Right-Down or Down-Right |
| `(10, 10)` | `48620` | Larger grid, verify no overflow |

**MinPathSum:**

| Input | Expected | Notes |
|-------|----------|-------|
| `[[1,3,1],[1,5,1],[4,2,1]]` | `7` | Path: 1->3->1->1->1 or 1->1->5->1->1? Nope: 1->1->2->1->1 = ... wait. 1->3->1->1->1=7, also 1->1->5->1->1=9. Optimal: right→right→down→down = 1+3+1+1+1=7 |
| `[[1,2,3]]` | `6` | Single row |
| `[[1],[2],[3]]` | `6` | Single column |
| `[[5]]` | `5` | Single cell |
| `[[1,2],[1,1]]` | `3` | 1->1->1 (down then right) |

**UniquePathsWithObstacles:**

| Input | Expected | Notes |
|-------|----------|-------|
| `[[0,0,0],[0,1,0],[0,0,0]]` | `2` | Classic — obstacle in the middle |
| `[[1,0]]` | `0` | Start blocked |
| `[[0,1]]` | `0` | End blocked |
| `[[0]]` | `1` | Single free cell |
| `[[1]]` | `0` | Single obstacle |
| `[[0,0],[0,0]]` | `2` | No obstacles |
| `[[0,1,0],[0,1,0],[0,0,0]]` | `1` | Wall in the middle, only one path |
| `[[0,0,0],[1,1,0],[0,0,0]]` | `1` | Horizontal wall, only one path |

**Knapsack01 / Knapsack01Optimized:**

| Weights | Values | W | Expected | Notes |
|---------|--------|---|----------|-------|
| `[1,3,4,5]` | `[1,4,5,7]` | `7` | `9` | Take items with w=3,v=4 and w=4,v=5 |
| `[2,3,4,5]` | `[3,4,5,6]` | `5` | `7` | Take w=2,v=3 and w=3,v=4 |
| `[10]` | `[100]` | `5` | `0` | Item too heavy |
| `[]` | `[]` | `10` | `0` | No items |
| `[1,2,3]` | `[6,10,12]` | `5` | `22` | Take all (total weight = 6 > 5), take w=2+w=3, v=10+12=22 |
| `[5]` | `[10]` | `5` | `10` | Exactly fits |

**SubsetSum:**

| Input | Target | Expected | Notes |
|-------|--------|----------|-------|
| `[3,34,4,12,5,2]` | `9` | `true` | 4+5=9 |
| `[3,34,4,12,5,2]` | `30` | `false` | No subset sums to 30 |
| `[1,2,3]` | `6` | `true` | All elements: 1+2+3=6 |
| `[1]` | `2` | `false` | Single element, too small |
| `[1]` | `1` | `true` | Single element, exact match |
| `[]` | `0` | `true` | Empty subset sums to 0 |
| `[1,2,3]` | `0` | `true` | Empty subset |

**CanPartition:**

| Input | Expected | Notes |
|-------|----------|-------|
| `[1,5,11,5]` | `true` | {1,5,5} and {11}, both sum to 11 |
| `[1,2,3,5]` | `false` | Total=11, odd |
| `[1,1]` | `true` | {1} and {1} |
| `[1]` | `false` | Can't split single element |
| `[2,2,1,1]` | `true` | {2,1} and {2,1} |
| `[1,2,5]` | `false` | Total=8, target=4, but no subset sums to 4 |
| `[0,0,0,0]` | `true` | All zeros, total=0, target=0 |

---

## 5. The Knapsack Family

### The Family Tree

All knapsack variants share the same skeleton: "for each item, decide how to use it, subject to a capacity constraint." The differences are in how many times each item can be used and what you're optimizing.

```
                        ┌─────────────────────┐
                        │    KNAPSACK FAMILY   │
                        │                      │
                        │  State: dp[i][w]     │
                        │  Choice: take/skip   │
                        │  Constraint: capacity│
                        └──────────┬───────────┘
                                   │
                 ┌─────────────────┼─────────────────┐
                 │                 │                  │
          ┌──────▼──────┐  ┌──────▼──────┐  ┌───────▼───────┐
          │ 0/1 Knapsack│  │  Unbounded  │  │    Bounded    │
          │             │  │  Knapsack   │  │   Knapsack    │
          │ Each item   │  │             │  │               │
          │ at most once│  │ Each item   │  │ Item i has    │
          │             │  │ unlimited   │  │ quantity q[i] │
          │ 1D: R-to-L  │  │             │  │               │
          └──────┬──────┘  │ 1D: L-to-R  │  │ Binary decomp │
                 │         │             │  │ into 0/1      │
                 │         │ Ex: Coin    │  └───────────────┘
                 │         │ Change      │
                 │         └─────────────┘
                 │
        ┌────────┼────────┐
        │                 │
  ┌─────▼──────┐  ┌──────▼───────────────┐
  │ Subset Sum │  │ Target Sum           │
  │            │  │ (assign +/- to each  │
  │ value =    │  │  number to reach     │
  │ weight     │  │  target)             │
  │            │  │                      │
  │ Boolean:   │  │ Reduces to subset    │
  │ dp[w]=T/F  │  │ sum with shifted     │
  └─────┬──────┘  │ target               │
        │         └──────────────────────┘
        │
  ┌─────▼───────────────┐
  │ Partition Equal      │
  │ Subset Sum           │
  │                      │
  │ target = total / 2   │
  │ If total is odd:     │
  │   impossible         │
  └──────────────────────┘
```

### Side-by-Side Recurrences

The power of the knapsack family is that every variant uses nearly the same recurrence. Here they are, aligned to make the differences pop:

**0/1 Knapsack (maximize value, each item once):**
```
dp[w] = max(dp[w], dp[w - weight[i]] + value[i])
iterate w: W down to weight[i]       ← RIGHT-TO-LEFT
```

**Unbounded Knapsack (maximize value, each item unlimited):**
```
dp[w] = max(dp[w], dp[w - weight[i]] + value[i])
iterate w: weight[i] up to W         ← LEFT-TO-RIGHT
```

**Coin Change — Minimum coins (unbounded, minimize):**
```
dp[w] = min(dp[w], dp[w - coin] + 1)
iterate w: coin up to amount          ← LEFT-TO-RIGHT
init: dp[0] = 0, dp[1..] = infinity
```

**Coin Change II — Count ways (unbounded, count):**
```
dp[w] += dp[w - coin]
iterate w: coin up to amount          ← LEFT-TO-RIGHT
init: dp[0] = 1
```

**Subset Sum (0/1, boolean):**
```
dp[w] = dp[w] || dp[w - num]
iterate w: target down to num         ← RIGHT-TO-LEFT
init: dp[0] = true
```

Five problems. One skeleton. The differences: loop direction (0/1 vs. unbounded), operation (max vs. min vs. sum vs. or), and what you're tracking (value, count, boolean).

---

## 6. Visual Diagrams

### Grid DP: Unique Paths (4 rows x 5 cols)

Each cell shows the number of unique paths from the top-left to that cell. Arrows show which cells contribute to each value.

```
  Moving only → and ↓

  ┌────┬────┬────┬────┬────┐
  │  1 │  1 │  1 │  1 │  1 │   Row 0: all 1s (only way: go right)
  │    │ ←  │ ←  │ ←  │ ←  │
  ├────┼────┼────┼────┼────┤
  │  1 │  2 │  3 │  4 │  5 │   dp[1][1] = dp[0][1] + dp[1][0] = 1+1 = 2
  │ ↑  │↑←  │↑←  │↑←  │↑←  │   dp[1][2] = dp[0][2] + dp[1][1] = 1+2 = 3
  ├────┼────┼────┼────┼────┤
  │  1 │  3 │  6 │ 10 │ 15 │   dp[2][1] = dp[1][1] + dp[2][0] = 2+1 = 3
  │ ↑  │↑←  │↑←  │↑←  │↑←  │   dp[2][2] = dp[1][2] + dp[2][1] = 3+3 = 6
  ├────┼────┼────┼────┼────┤
  │  1 │  4 │ 10 │ 20 │ 35 │   dp[3][4] = dp[2][4] + dp[3][3] = 15+20 = 35
  │ ↑  │↑←  │↑←  │↑←  │↑←  │
  └────┴────┴────┴────┴────┘
                         ^^^
                    Answer: 35

  1D optimization (single row, updated per row):

  After row 0:  [1, 1, 1, 1, 1]
  After row 1:  [1, 2, 3, 4, 5]     dp[c] += dp[c-1]
  After row 2:  [1, 3, 6, 10, 15]   dp[c] += dp[c-1]
  After row 3:  [1, 4, 10, 20, 35]  dp[c] += dp[c-1]
```

---

### Grid DP: Min Path Sum

Grid values:
```
  ┌───┬───┬───┬───┐
  │ 1 │ 3 │ 1 │ 2 │
  ├───┼───┼───┼───┤
  │ 1 │ 5 │ 1 │ 3 │
  ├───┼───┼───┼───┤
  │ 4 │ 2 │ 1 │ 1 │
  └───┴───┴───┴───┘
```

DP table (minimum cost to reach each cell):
```
  ┌────┬────┬────┬────┐
  │  1 │  4 │  5 │  7 │   Row 0: cumulative sum (only from left)
  ├────┼────┼────┼────┤      dp[0][1] = 1+3 = 4
  │  2 │  7 │  6 │  9 │   dp[1][1] = 5 + min(dp[0][1], dp[1][0])
  ├────┼────┼────┼────┤            = 5 + min(4, 2) = 7
  │  6 │  8 │  7 │  8 │   dp[2][2] = 1 + min(dp[1][2], dp[2][1])
  └────┴────┴────┴────┘            = 1 + min(6, 8) = 7
                    ^^^
               Answer: 8

  Optimal path: 1 → 3 → 1 → 1 → 1 → 1 = 8
                (right, right, down, down, right)
    or:         1 → 1 → 5 → 1 → 1 → 1 = 10  (worse)
    or:         1 → 1 → 2 → 1 → 1      = path is 1→1→(need to get right)
    Best:       1 → 3 → 1 → 1 → 1 → 1
                Start at (0,0)=1, go right to (0,1)=3, right to (0,2)=1,
                down to (1,2)=1, down to (2,2)=1, right to (2,3)=1
                Total = 1+3+1+1+1+1 = 8 ✓
```

---

### 0/1 Knapsack: 2D Table

Items: `[(w=1,v=1), (w=3,v=4), (w=4,v=5), (w=5,v=7)]`, Capacity W=7

```
  dp[i][w] = max value using first i items with capacity w

         w=  0    1    2    3    4    5    6    7
        ┌────┬────┬────┬────┬────┬────┬────┬────┐
  i=0   │  0 │  0 │  0 │  0 │  0 │  0 │  0 │  0 │  no items
  (none)├────┼────┼────┼────┼────┼────┼────┼────┤
  i=1   │  0 │  1 │  1 │  1 │  1 │  1 │  1 │  1 │  item 1: w=1, v=1
  (w=1) │    │ T  │ T  │ T  │ T  │ T  │ T  │ T  │  T=take, S=skip
        ├────┼────┼────┼────┼────┼────┼────┼────┤
  i=2   │  0 │  1 │  1 │  4 │  5 │  5 │  5 │  5 │  item 2: w=3, v=4
  (w=3) │    │ S  │ S  │ T  │ T  │ T  │ T  │ T  │  dp[2][4]=max(1, dp[1][1]+4)
        ├────┼────┼────┼────┼────┼────┼────┼────┤            =max(1, 5)=5
  i=3   │  0 │  1 │  1 │  4 │  5 │  6 │  6 │  9 │  item 3: w=4, v=5
  (w=4) │    │ S  │ S  │ S  │ S  │ T  │ T  │ T  │  dp[3][7]=max(5, dp[2][3]+5)
        ├────┼────┼────┼────┼────┼────┼────┼────┤            =max(5, 4+5)=9
  i=4   │  0 │  1 │  1 │  4 │  5 │  7 │  8 │  9 │  item 4: w=5, v=7
  (w=5) │    │ S  │ S  │ S  │ S  │ T  │ T  │ S  │  dp[4][7]=max(9, dp[3][2]+7)
        └────┴────┴────┴────┴────┴────┴────┴────┘            =max(9, 1+7)=8→9 wins

  Answer: dp[4][7] = 9 (take items 2 and 3: w=3+4=7, v=4+5=9)
```

---

### Knapsack 1D Space Optimization: Why Right-to-Left

This is the single trickiest concept. Let's trace it step by step.

**Setup:** Items: `[(w=2,v=3), (w=3,v=4)]`, Capacity W=5

**0/1 Knapsack — RIGHT-TO-LEFT (correct):**

```
  Initial dp:  [0, 0, 0, 0, 0, 0]
                w=0  1  2  3  4  5

  Processing item 1 (w=2, v=3), iterate w from 5 down to 2:

    w=5: dp[5] = max(dp[5], dp[5-2]+3) = max(0, dp[3]+3) = max(0, 0+3) = 3
    w=4: dp[4] = max(dp[4], dp[4-2]+3) = max(0, dp[2]+3) = max(0, 0+3) = 3
    w=3: dp[3] = max(dp[3], dp[3-2]+3) = max(0, dp[1]+3) = max(0, 0+3) = 3
    w=2: dp[2] = max(dp[2], dp[2-2]+3) = max(0, dp[0]+3) = max(0, 0+3) = 3

  After item 1: [0, 0, 3, 3, 3, 3]
                              ↑ ↑ ↑ ↑  each capacity ≥ 2 can hold one of item 1

  Processing item 2 (w=3, v=4), iterate w from 5 down to 3:

    w=5: dp[5] = max(dp[5], dp[5-3]+4) = max(3, dp[2]+4) = max(3, 3+4) = 7 ✓
    w=4: dp[4] = max(dp[4], dp[4-3]+4) = max(3, dp[1]+4) = max(3, 0+4) = 4
    w=3: dp[3] = max(dp[3], dp[3-3]+4) = max(3, dp[0]+4) = max(3, 0+4) = 4

  After item 2: [0, 0, 3, 4, 4, 7]

  Answer: dp[5] = 7 (both items, w=2+3=5, v=3+4=7) ✓
```

**WRONG: LEFT-TO-RIGHT (what happens if you iterate forward):**

```
  Initial dp:  [0, 0, 0, 0, 0, 0]

  Processing item 1 (w=2, v=3), iterate w from 2 to 5:

    w=2: dp[2] = max(0, dp[0]+3) = 3       dp is now [0,0,3,0,0,0]
    w=3: dp[3] = max(0, dp[1]+3) = 3       dp is now [0,0,3,3,0,0]
    w=4: dp[4] = max(0, dp[2]+3) = 3+3 = 6  ← BUG! dp[2] was already updated!
                                                This uses item 1 TWICE
    w=5: dp[5] = max(0, dp[3]+3) = 3+3 = 6  ← BUG! item 1 used TWICE again

  After item 1: [0, 0, 3, 3, 6, 6]
                                ↑ ↑  WRONG: can't use item 1 twice in 0/1

  With 0/1 knapsack, dp[4] should be 3 (one item), not 6 (two items).
```

**The fix is simple: iterate right-to-left so that when you read `dp[w-weight]`, it still holds the PREVIOUS item's value (hasn't been updated yet in this pass).**

For unbounded knapsack, the left-to-right "bug" is exactly the behavior you want — using an item multiple times is legal.

---

## 7. Self-Assessment

Answer these without looking at the notes. If you can't, revisit the relevant section.

### Question 1
**Why must the inner loop in 0/1 knapsack iterate right-to-left for the 1D space optimization, but left-to-right for unbounded knapsack?**

<details>
<summary>Answer</summary>

In the 1D array, `dp[w]` stores the best value for capacity `w` using items processed so far. When you compute `dp[w] = max(dp[w], dp[w - weight[i]] + value[i])`, the value `dp[w - weight[i]]` is a smaller index.

**Right-to-left (0/1):** You process larger indices first. When you read `dp[w - weight[i]]`, that smaller index hasn't been updated yet for item `i` — it still reflects the state from the previous item's iteration. This means item `i` is only counted once.

**Left-to-right (unbounded):** You process smaller indices first. When you read `dp[w - weight[i]]`, that smaller index has already been updated to include item `i`. So `dp[w]` can build on a solution that already used item `i` — allowing multiple uses of the same item.

The only code difference is `for w := W; w >= weight[i]; w--` vs. `for w := weight[i]; w <= W; w++`. Same recurrence, different semantics.
</details>

### Question 2
**How do you detect if unique paths with obstacles should return 0 immediately, without filling the table?**

<details>
<summary>Answer</summary>

Check two things immediately:
1. **Start cell is an obstacle:** `grid[0][0] == 1` → return 0. No path can begin.
2. **End cell is an obstacle:** `grid[m-1][n-1] == 1` → return 0. No path can end there.

These are O(1) checks that save you from filling an entire table for an impossible case.

During table filling, there are additional early-exit opportunities: if the entire first row or first column is blocked before reaching any cell, everything beyond the blockage is unreachable. But the start/end checks are the ones that matter most for clean, defensive code.
</details>

### Question 3
**In the 1D space optimization for grid DP (unique paths), what does `dp[c]` represent BEFORE you update it in the current row, and what does it represent AFTER?**

<details>
<summary>Answer</summary>

**Before updating `dp[c]` in row `r`:** It holds the value from the previous row — `dp[r-1][c]` — the number of paths to reach cell `(r-1, c)`. This is the "from above" contribution.

**After updating `dp[c]` with `dp[c] += dp[c-1]`:** It holds the value for the current row — `dp[r][c]` — the total paths to reach cell `(r, c)`.

The trick works because:
- `dp[c]` (before update) = value from above (previous row, not yet overwritten)
- `dp[c-1]` (after update) = value from the left (current row, just computed)

Adding them gives `dp[r-1][c] + dp[r][c-1]`, which is exactly the recurrence. The 1D array simultaneously holds "old" values (not yet updated, to its right) and "new" values (already updated, to its left). This dual role is what makes it work.
</details>

### Question 4
**You're given a partition equal subset sum problem. Walk through the complete chain of reductions to arrive at the 1D DP solution.**

<details>
<summary>Answer</summary>

**Step 1: Partition → Subset Sum.**
Compute `total = sum(nums)`. If `total` is odd, return false (can't split evenly). Otherwise, the problem becomes: "does any subset of `nums` sum to `total/2`?"

**Step 2: Subset Sum → Boolean 0/1 Knapsack.**
Create `dp[0..total/2]` of booleans. Set `dp[0] = true` (empty subset). For each number `num` in `nums`, iterate `w` from `total/2` down to `num` (right-to-left because each number is used at most once): `dp[w] = dp[w] || dp[w - num]`.

**Step 3: Return `dp[total/2]`.**

Three reductions, each trivial, composing into a clean O(n * sum/2) solution with O(sum/2) space.

The chain: partition → "does a subset sum to half?" → "boolean knapsack with capacity = half" → "1D array, right-to-left iteration."
</details>

### Question 5
**Coin change asks for the minimum number of coins to make an amount. Is this 0/1 or unbounded knapsack? How does the recurrence differ from the standard knapsack?**

<details>
<summary>Answer</summary>

Coin change is **unbounded knapsack** — each coin denomination can be used unlimited times.

The differences from standard knapsack:
1. **Objective:** Minimize instead of maximize. Use `min` instead of `max`.
2. **Value per item:** Every coin has value 1 (it costs one coin to use it). So the recurrence is `dp[w] = min(dp[w], dp[w - coin] + 1)`.
3. **Initialization:** `dp[0] = 0` (zero coins to make amount 0). `dp[1..amount] = infinity` (not yet achievable). Use `math.MaxInt32` (not `MaxInt`) to avoid overflow when adding 1.
4. **Loop direction:** Left-to-right (unbounded — same coin usable multiple times).
5. **Return value:** If `dp[amount]` is still infinity, return -1 (impossible).

The structure is identical to unbounded knapsack — only the operator (min vs. max), the "value" (1 vs. item-specific), and the initialization change.
</details>
