# Day 15 — Dynamic Programming: Foundations

---

## 1. Curated Learning Resources

| # | Resource | Format | Why This One |
|---|----------|--------|-------------|
| 1 | [MIT 6.006 — Dynamic Programming I (Fibonacci, Shortest Paths)](https://www.youtube.com/watch?v=OQ5jsbhAv_M) | Video (52 min) | Erik Demaine builds DP from first principles: memoization on Fibonacci, then generalizes to the 5-step recipe. Watch at 1.5x during review, re-watch the recipe segment after implementing. This is the canonical "DP mental model" lecture. |
| 2 | [MIT 6.006 — Dynamic Programming II (Text Justification, Knapsack)](https://www.youtube.com/watch?v=ENyox7kNKeY) | Video (51 min) | Continues the recipe from Lecture I, applying it to new problems. Reinforces that the hard part is defining the subproblem. Watch sections on subproblem definition. |
| 3 | [Errichto — Dynamic Programming lecture #1](https://www.youtube.com/watch?v=YBSt1jYwVfU) | Video (36 min) | Competitive programmer walks through Fibonacci, climbing stairs, and frog problems. Excellent at building intuition for *why* memoization works — draws recursion trees and counts repeated calls. |
| 4 | [Errichto — Dynamic Programming lecture #2](https://www.youtube.com/watch?v=1mtvm2ubHCY) | Video (34 min) | Covers coin change, LIS, and grid DP. Focuses on transitioning from recursive thinking to tabulation. Good complement to lecture #1. |
| 5 | [Back To Back SWE — The Change Making Problem](https://www.youtube.com/watch?v=jgiZlGzXMBw) | Video (17 min) | Best single visual explanation of coin change DP. Draws the entire DP table, fills it cell by cell, and explains why each cell takes the value it does. |
| 6 | [NeetCode — House Robber](https://www.youtube.com/watch?v=73r3KWiEvyk) | Video (10 min) | Clean derivation of the recurrence from the decision tree. Shows the "rob or skip" choice and how it maps to `dp[i] = max(dp[i-1], dp[i-2] + nums[i])`. |
| 7 | [Reducible — The essence of Dynamic Programming](https://www.youtube.com/watch?v=aPQY__2H3tE) | Video (30 min) | Animated visualization of memoization vs tabulation on the same problems. Best visual resource for understanding the difference between top-down and bottom-up. Shows the DAG of subproblems explicitly. |
| 8 | [VisuAlgo — Recursion tree / DP visualization](https://visualgo.net/en/recursion) | Interactive | Step through recursion trees for Fibonacci, coin change, etc. Toggle memoization on/off and watch repeated subproblems disappear. Hands-on — use this during the review block. |

**Suggested order:** Resource 7 (Reducible) for the mental model, then 1 (MIT lecture) for the recipe during the review block. Keep 8 (VisuAlgo) open to play with. Resources 3-6 during implementation as needed. Resource 2 as follow-up reading.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 - 12:20 — Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:05 | Read Section 3 (Core Concepts) below end to end. Focus on the DP recipe — internalize the 5 steps. Do not skim; this is the framework you'll apply to every problem today. |
| 12:05 - 12:10 | **Draw the Fibonacci recursion tree by hand** for `fib(6)`. Circle every repeated subproblem. Count how many times `fib(3)` gets called. This is *why* DP works — burn this image into your brain. (Reference the ASCII diagram in Section 6 to check your work.) |
| 12:10 - 12:15 | **Apply the recipe to climbing stairs on paper.** Write down: (1) state: `dp[i]` = number of ways to reach step `i`. (2) recurrence: `dp[i] = dp[i-1] + dp[i-2]`. (3) base cases: `dp[0] = 1, dp[1] = 1`. (4) order: left to right. (5) space: collapse to two variables. |
| 12:15 - 12:20 | **Apply the recipe to coin change on paper.** Write down all 5 steps. Reference Section 5 to verify. Understand why the traversal order is left to right and why we initialize with `math.MaxInt32`. |

### 12:20 - 1:20 — Implement

| Time | Activity |
|------|----------|
| 12:20 - 12:30 | **Climbing Stairs — memoized (top-down).** Write `ClimbStairsMemo(n int) int` with a recursive helper and a `map[int]int` cache. Verify with `n=0,1,2,5,10`. |
| 12:30 - 12:38 | **Climbing Stairs — tabulated (bottom-up).** Write `ClimbStairsTabulated(n int) int` with a `[]int` table, then optimize to two variables. Compare the two versions side by side and note: same answers, different mechanics. |
| 12:38 - 12:52 | **House Robber.** Write `HouseRobber(nums []int) int`. Start with the recurrence: rob this house + best from two back, or skip this house + best from one back. Implement with full array first, then optimize to two variables (`prev2, prev1`). Test: `[1,2,3,1]` -> 4, `[2,7,9,3,1]` -> 12, empty, single element. |
| 12:52 - 1:08 | **Coin Change.** Write `CoinChange(coins []int, amount int) int`. Initialize `dp[0] = 0`, all others `math.MaxInt32`. For each amount `i` from 1 to `amount`, try each coin. Test: `coins=[1,5,11], amount=15` -> 3 (5+5+5), `coins=[2], amount=3` -> -1, `amount=0` -> 0. |
| 1:08 - 1:20 | **Longest Increasing Subsequence.** Write `LengthOfLIS_N2(nums []int) int` with the O(n^2) DP. Then write `LengthOfLIS_NLogN(nums []int) int` using patience sorting (binary search on tails array). Test: `[10,9,2,5,3,7,101,18]` -> 4, `[0,1,0,3,2,3]` -> 4, single element -> 1. |

### 1:20 - 1:50 — Solidify

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | **Edge cases for all problems.** Climbing stairs: `n=0`, `n=1`, `n=45` (large). House robber: all same values, two houses, `[100]`. Coin change: impossible amounts, single coin `[1]`, large amount. LIS: strictly decreasing input (answer=1), all same values, already sorted. |
| 1:30 - 1:40 | **Trace coin change table by hand** for `coins=[1,3,4], amount=6`. Fill out the full DP table on paper. At each cell, write which coin was used. Verify against your implementation. (Reference ASCII diagram in Section 6.) |
| 1:40 - 1:50 | **Compare top-down vs bottom-up.** Re-implement coin change as a memoized recursive function. Note: (1) the code is almost identical to the math, (2) but deep recursion on `amount=100000` would overflow the stack. Write down when each approach is better (see Section 3). |

### 1:50 - 2:00 — Recap (From Memory)

Write down without looking:

1. The 5-step DP recipe.
2. The recurrence for house robber and coin change.
3. The difference between memoization and tabulation (one sentence each).
4. Why the O(n log n) LIS works (patience sorting in one sentence).
5. One gotcha you hit during implementation.

---

## 3. Core Concepts Deep Dive

### The DP Recipe

Every DP problem follows the same 5-step process. The steps are not equally hard — step 1 is where you win or lose.

**Step 1: Define the state — what does `dp[i]` (or `dp[i][j]`) represent?**

This is the hardest and most important step. Everything else flows from it. If you get this wrong, no amount of clever recurrence-writing will save you. Spend 60-70% of your thinking time here.

Good states have these properties:
- The final answer is easily extracted from the table (often `dp[n]` or `max(dp[...])`).
- The state captures *enough* information to make the optimal decision without knowing how you got here.
- The number of distinct states is polynomial (that's what makes DP efficient).

Common state definitions:
- `dp[i]` = answer considering the first `i` elements.
- `dp[i]` = answer where element `i` is the *last* element used (e.g., LIS).
- `dp[i]` = answer for target value `i` (e.g., coin change where `i` is the amount).
- `dp[i][j]` = answer for subproblem on `s[i..j]` or prefixes `s[0..i]` and `t[0..j]`.

**Step 2: Write the recurrence — how does `dp[i]` relate to smaller subproblems?**

Express `dp[i]` in terms of `dp[j]` where `j < i` (or a smaller subproblem in general). This is the "optimal substructure" at work. The recurrence usually involves a `min`, `max`, or `sum` over a set of choices.

**Step 3: Identify base cases — what are the trivial subproblems?**

These are the subproblems small enough to solve directly: `dp[0]`, empty string, amount of 0, etc. Getting base cases wrong is the #1 source of off-by-one bugs in DP.

**Step 4: Determine the traversal order — in what order do you fill the table?**

When you compute `dp[i]`, every `dp[j]` it depends on must already be computed. For 1D problems going left-to-right, this is usually obvious. For 2D problems (especially on intervals like `dp[i][j]`), the order can be subtle — you might need to iterate by length of the interval.

**Step 5: Optimize space (optional) — can you reduce the DP table?**

If `dp[i]` depends only on `dp[i-1]` (and maybe `dp[i-2]`), you don't need the full array — just keep two variables. This reduces space from O(n) to O(1). If it depends on a full previous row in a 2D table, reduce from O(m*n) to O(n).

---

### Top-Down (Memoization) vs Bottom-Up (Tabulation)

These are two ways to implement the same recurrence. The *math* is identical; the *mechanics* differ.

**Side-by-side comparison on Climbing Stairs:**

```
Top-Down (Memoized)                     Bottom-Up (Tabulated)
─────────────────────────────────────   ─────────────────────────────────────
func climbMemo(n int,                   func climbTab(n int) int {
    memo map[int]int) int {                 if n <= 1 { return 1 }
    if n <= 1 { return 1 }                 dp := make([]int, n+1)
    if v, ok := memo[n]; ok {              dp[0], dp[1] = 1, 1
        return v                           for i := 2; i <= n; i++ {
    }                                          dp[i] = dp[i-1] + dp[i-2]
    memo[n] = climbMemo(n-1, memo) +       }
              climbMemo(n-2, memo)          return dp[n]
    return memo[n]                     }
}

How it works:                          How it works:
- Starts at the full problem (n)       - Starts at the base cases (0, 1)
- Recurses down to base cases          - Iterates up to the full problem (n)
- Caches results on the way back up    - Fills table left to right
- Only computes reachable subproblems  - Computes ALL subproblems 0..n

Pros:                                  Pros:
- Easier to derive from recursion      - No recursion overhead
- Only solves needed subproblems       - No stack overflow risk
- Natural for problems with sparse     - Easier to optimize space
  state space                          - Generally faster (no hash lookups)

Cons:                                  Cons:
- Stack overflow on deep recursion     - May compute unused subproblems
  (Go default stack ~1GB, but still)   - Harder to derive traversal order
- Hash map overhead per lookup           for complex state spaces
- Harder to optimize space
```

**When to use which:**

| Situation | Prefer |
|-----------|--------|
| Deriving the solution for the first time | Top-down (think recursively, add cache) |
| The state space is sparse (many states never needed) | Top-down (only computes what's reachable) |
| The problem has a clear iterative order | Bottom-up (simpler, faster) |
| You need space optimization | Bottom-up (rolling array is straightforward) |
| Deep recursion possible (`n > 10000`) | Bottom-up (avoids stack overflow) |
| Interview setting, running low on time | Top-down (faster to code, less bug-prone) |

---

### Overlapping Subproblems: WHY DP Works

DP is not just "recursion with caching." The caching matters because the same subproblems are solved **exponentially many times** in the naive recursive approach.

Consider `fib(6)`. Without memoization, the recursion tree has **25 nodes** but only **7 unique values** (fib(0) through fib(6)). `fib(3)` is computed 3 times, `fib(2)` is computed 5 times. For `fib(n)`, the naive approach takes O(2^n) time. With memoization, each unique subproblem is solved exactly once: O(n) time.

This exponential-to-polynomial collapse is the entire point of DP. If a recursive problem does NOT have overlapping subproblems (e.g., merge sort — each subarray is unique), memoization doesn't help, and you don't need DP.

See the recursion tree in Section 6 for a visual.

---

### Optimal Substructure: WHAT DP Exploits

A problem has **optimal substructure** if the optimal solution to the whole problem contains optimal solutions to its subproblems.

**Example (coin change):** If the optimal way to make 15 cents uses a 5-cent coin, then the remaining 10 cents must also be made optimally. If there were a better way to make 10 cents, you could substitute it and get a better solution for 15 cents — contradiction.

**Counterexample (longest simple path in a general graph):** The longest simple path from A to C through B does NOT necessarily contain the longest simple path from A to B. Why? The longest A-to-B path might visit vertices that the B-to-C path also needs, making the combined path visit a vertex twice (not simple). Cycles in the graph break the independence of subproblems.

This is why longest path is NP-hard in general graphs but solvable by DP on DAGs (no cycles, so subproblems are independent).

Both overlapping subproblems and optimal substructure are needed for DP:
- Optimal substructure alone: greedy or divide-and-conquer might work.
- Overlapping subproblems alone: not useful — caching doesn't help if subproblems don't compose optimally.
- Both together: DP.

---

### State Reduction / Space Optimization

Once you have a working DP, look at the dependency pattern:

| Dependency Pattern | Space Optimization |
|--------------------|--------------------|
| `dp[i]` depends on `dp[i-1]` and `dp[i-2]` | Two variables: `prev2, prev1` → O(1) |
| `dp[i]` depends on `dp[i-1]` only | One variable → O(1) |
| `dp[i][j]` depends on `dp[i-1][...]` | Two rows (prev, curr) → O(n) |
| `dp[i]` depends on all `dp[j]` for `j < i` | **Cannot easily optimize** — need full array → O(n) |

**Climbing stairs / House robber:** `dp[i]` depends on `dp[i-1]` and `dp[i-2]`. Collapse to:
```go
prev2, prev1 := 1, 1
for i := 2; i <= n; i++ {
    prev2, prev1 = prev1, prev2+prev1
}
return prev1
```

**LIS (O(n^2) version):** `dp[i] = max(dp[j]+1)` for all `j < i` where `nums[j] < nums[i]`. Each `dp[i]` depends on *all* previous entries. You need the full array — no space optimization is possible for the standard O(n^2) approach.

**Coin change:** `dp[i]` depends on `dp[i-coin]` for each coin. The coin values can be large, so `dp[i]` might depend on entries far back in the array. You need the full array.

---

### How to Recognize a DP Problem

Look for these signals:

| Signal | Example |
|--------|---------|
| "How many ways..." | Climbing stairs, unique paths |
| "What is the minimum/maximum cost..." | Coin change, minimum path sum |
| "Is it possible to..." | Subset sum, word break |
| "What is the longest/shortest..." | LIS, LCS, edit distance |

Combined with:
- **Choices at each step** — at each index, you decide to take/skip, use this coin or not, extend or start fresh.
- **Overlapping subproblems** — different sequences of choices lead to the same intermediate state.

If the problem has choices but NO overlapping subproblems, it's likely greedy or divide-and-conquer. If it asks to enumerate ALL solutions (not just count/optimize), it's likely backtracking, not DP.

---

## 4. Implementation Checklist

### Function Signatures

```go
package dp

import "math"

// --- Climbing Stairs ---

// ClimbStairsMemo returns the number of distinct ways to climb n stairs
// (1 or 2 steps at a time) using top-down memoization.
func ClimbStairsMemo(n int) int

// ClimbStairsTabulated returns the number of distinct ways to climb n stairs
// using bottom-up tabulation with O(1) space.
func ClimbStairsTabulated(n int) int

// --- House Robber ---

// HouseRobber returns the maximum amount of money you can rob
// without robbing two adjacent houses.
func HouseRobber(nums []int) int

// --- Coin Change ---

// CoinChange returns the fewest number of coins needed to make up `amount`.
// Returns -1 if the amount cannot be made up by any combination.
func CoinChange(coins []int, amount int) int

// --- Longest Increasing Subsequence ---

// LengthOfLIS_N2 returns the length of the longest strictly increasing
// subsequence using O(n^2) DP.
func LengthOfLIS_N2(nums []int) int

// LengthOfLIS_NLogN returns the length of the longest strictly increasing
// subsequence using patience sorting with binary search in O(n log n).
func LengthOfLIS_NLogN(nums []int) int
```

### Test Cases & Edge Cases

**Climbing Stairs:**

| Input | Expected | Notes |
|-------|----------|-------|
| `n = 0` | `1` | 1 way to stay at ground (base case) |
| `n = 1` | `1` | Only one step |
| `n = 2` | `2` | (1+1) or (2) |
| `n = 5` | `8` | |
| `n = 10` | `89` | Fibonacci-like growth |
| `n = 45` | `1836311903` | Large input, ensure no overflow (fits int32) |

**House Robber:**

| Input | Expected | Notes |
|-------|----------|-------|
| `[]` | `0` | Empty input |
| `[5]` | `5` | Single house |
| `[1, 2]` | `2` | Pick the larger |
| `[1, 2, 3, 1]` | `4` | Rob house 0 and 2 |
| `[2, 7, 9, 3, 1]` | `12` | Rob house 0, 2, 4 |
| `[5, 5, 5, 5, 5]` | `15` | Rob alternating |
| `[100, 1, 1, 100]` | `200` | Rob first and last |

**Coin Change:**

| Input | Expected | Notes |
|-------|----------|-------|
| `coins=[1,5,11], amount=15` | `3` | 5+5+5 (greedy would pick 11+1+1+1+1 = 5 coins!) |
| `coins=[1,2,5], amount=11` | `3` | 5+5+1 |
| `coins=[2], amount=3` | `-1` | Impossible |
| `coins=[1], amount=0` | `0` | Zero amount |
| `coins=[1], amount=1` | `1` | |
| `coins=[1], amount=10000` | `10000` | Large amount, stress test |
| `coins=[186,419,83,408], amount=6249` | `20` | Non-trivial denomination set |

**LengthOfLIS:**

| Input | Expected | Notes |
|-------|----------|-------|
| `[10,9,2,5,3,7,101,18]` | `4` | [2,3,7,101] or [2,3,7,18] |
| `[0,1,0,3,2,3]` | `4` | [0,1,2,3] |
| `[7,7,7,7,7]` | `1` | Strictly increasing, not non-decreasing |
| `[1]` | `1` | Single element |
| `[5,4,3,2,1]` | `1` | Strictly decreasing |
| `[1,2,3,4,5]` | `5` | Already sorted |
| `[]` | `0` | Empty input |

---

## 5. The DP Recipe Applied Step by Step

Use this section as a **template**. For every new DP problem, walk through these exact 5 steps on paper before writing code.

---

### Problem A: Coin Change

**Problem statement:** Given coins of different denominations and a total amount, find the fewest number of coins needed to make up that amount. If it's impossible, return -1.

**Step 1 — Define the state:**

> `dp[i]` = the minimum number of coins needed to make amount `i`.

Why this state? Because the answer we want is `dp[amount]`, and we can build up from smaller amounts. The key insight is that the state is the *remaining amount*, not the index of coins used. We don't need to track which coins we've used because coins are unlimited (unbounded).

**Step 2 — Write the recurrence:**

> `dp[i] = min(dp[i - coin] + 1)` for each `coin` in `coins`, where `i - coin >= 0`.

At amount `i`, we try using each coin. If we use a coin of value `c`, we need 1 coin plus the optimal solution for amount `i - c`. We take the minimum over all valid coin choices.

**Step 3 — Base cases:**

> `dp[0] = 0` — zero coins needed to make amount 0.
> `dp[i] = math.MaxInt32` for all `i > 0` — initialized to "impossible" (infinity).

We use `MaxInt32` (not `MaxInt64`) so that `dp[i-coin] + 1` doesn't overflow.

**Step 4 — Traversal order:**

> Left to right: `i` from `1` to `amount`.

When computing `dp[i]`, we need `dp[i - coin]`. Since all coin values are positive, `i - coin < i`, so `dp[i - coin]` has already been computed. Simple.

**Step 5 — Space optimization:**

> **Not applicable.** `dp[i]` can depend on `dp[i - 1]`, `dp[i - 5]`, `dp[i - 11]`, etc. The lookback distance depends on the coin values, and we might need any previous entry. The full array of size `amount + 1` is required.

**Final implementation:**

```go
func CoinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := 1; i <= amount; i++ {
        dp[i] = math.MaxInt32
    }
    for i := 1; i <= amount; i++ {
        for _, coin := range coins {
            if coin <= i && dp[i-coin] != math.MaxInt32 {
                if dp[i-coin]+1 < dp[i] {
                    dp[i] = dp[i-coin] + 1
                }
            }
        }
    }
    if dp[amount] == math.MaxInt32 {
        return -1
    }
    return dp[amount]
}
```

---

### Problem B: House Robber

**Problem statement:** Given an array where each element represents the money in a house, find the maximum money you can rob without robbing two adjacent houses.

**Step 1 — Define the state:**

> `dp[i]` = the maximum money we can rob from the first `i` houses (houses `0` through `i-1`).

Alternative (equally valid): `dp[i]` = max money robbing from houses `i..n-1`. Both work. We'll use the first definition because left-to-right traversal feels more natural.

**Step 2 — Write the recurrence:**

> `dp[i] = max(dp[i-1], dp[i-2] + nums[i-1])`

At house `i`, we have two choices:
- **Skip** house `i-1`: take the best from the first `i-1` houses → `dp[i-1]`.
- **Rob** house `i-1`: take its money `nums[i-1]` plus the best from the first `i-2` houses → `dp[i-2] + nums[i-1]`.

(Note the index shift: `dp[i]` considers `i` houses, so the last house is `nums[i-1]`.)

**Step 3 — Base cases:**

> `dp[0] = 0` — no houses, no money.
> `dp[1] = nums[0]` — one house, rob it.

**Step 4 — Traversal order:**

> Left to right: `i` from `2` to `n`.

`dp[i]` depends on `dp[i-1]` and `dp[i-2]`, both already computed.

**Step 5 — Space optimization:**

> YES. `dp[i]` depends only on `dp[i-1]` and `dp[i-2]`. Collapse to two variables.

```go
func HouseRobber(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
    if len(nums) == 1 {
        return nums[0]
    }
    prev2, prev1 := 0, nums[0]
    for i := 2; i <= len(nums); i++ {
        curr := prev1 // skip house i-1
        if prev2+nums[i-1] > curr {
            curr = prev2 + nums[i-1] // rob house i-1
        }
        prev2, prev1 = prev1, curr
    }
    return prev1
}
```

---

### Template for Any New DP Problem

Copy this checklist and fill it in on paper before coding:

```
1. STATE:    dp[___] = ________________________________________________
             (What does each entry represent? What is the "question" it answers?)

2. RECURRENCE: dp[___] = ______________________________________________
               (How does it combine smaller subproblems? What choices exist?)

3. BASE CASES: dp[___] = ___   dp[___] = ___
               (What are the trivial/smallest subproblems?)

4. ORDER:    Fill table from ___ to ___ because ________________________
             (Why are dependencies satisfied in this order?)

5. SPACE:    dp[i] depends on dp[___], so we can reduce to _____________
             (Two variables? One row? Full table required?)
```

---

## 6. Visual Diagrams

### Fibonacci Recursion Tree — Overlapping Subproblems

This is `fib(6)`. Nodes marked with `*` are **repeated computations** — the same subproblem solved multiple times. This is why naive recursion is O(2^n) and why memoization collapses it to O(n).

```
                              fib(6)
                            /        \
                       fib(5)          fib(4)*
                      /      \          /     \
                 fib(4)    fib(3)*   fib(3)*  fib(2)*
                /     \     /    \    /    \
           fib(3)  fib(2)* f(2)* f(1) f(2)* f(1)
           /    \    / \    / \
        fib(2) f(1) f(1) f(0) f(1) f(0)
        /   \
     f(1)  f(0)

  Unique subproblems:  fib(0), fib(1), fib(2), fib(3), fib(4), fib(5), fib(6)  → 7
  Total recursive calls without memo: 25
  Total with memoization: 7 (each computed once)

  * = repeated subproblem (would be a cache hit with memoization)
```

With memoization, the tree collapses to a linear chain:

```
  fib(6) → fib(5) → fib(4) → fib(3) → fib(2) → fib(1), fib(0)
                                                   ↑ base cases

  Each node computed exactly once. O(n) time, O(n) space.
```

---

### Coin Change DP Table

`coins = [1, 3, 4]`, `amount = 6`

Each cell `dp[i]` = minimum coins to make amount `i`. We fill left to right. For each cell, we check each coin and take `min(dp[i], dp[i-coin] + 1)`.

```
  Amount:    0     1     2     3     4     5     6
           ┌─────┬─────┬─────┬─────┬─────┬─────┬─────┐
  dp[i]:   │  0  │  1  │  2  │  1  │  1  │  2  │  2  │
           └─────┴─────┴─────┴─────┴─────┴─────┴─────┘

  How each cell was computed:
  ─────────────────────────────────────────────────────
  dp[0] = 0                          (base case)
  dp[1] = dp[1-1]+1 = dp[0]+1 = 1   (used coin 1)
  dp[2] = dp[2-1]+1 = dp[1]+1 = 2   (used coin 1)
  dp[3] = min(dp[3-1]+1, dp[3-3]+1)
        = min(dp[2]+1, dp[0]+1)
        = min(3, 1) = 1              (used coin 3)
  dp[4] = min(dp[4-1]+1, dp[4-3]+1, dp[4-4]+1)
        = min(dp[3]+1, dp[1]+1, dp[0]+1)
        = min(2, 2, 1) = 1           (used coin 4)
  dp[5] = min(dp[4]+1, dp[2]+1, dp[1]+1)
        = min(2, 3, 2) = 2           (used coin 1 or 4)
  dp[6] = min(dp[5]+1, dp[3]+1, dp[2]+1)
        = min(3, 2, 3) = 2           (used coin 3)

  Answer: dp[6] = 2 → coins used: 3 + 3

  Note: A greedy approach (always pick largest coin) would give
  4 + 1 + 1 = 3 coins. DP finds the true optimum.
```

---

### LIS DP Array — Step by Step

`nums = [10, 9, 2, 5, 3, 7, 101, 18]`

`dp[i]` = length of the longest increasing subsequence **ending at index `i`**.

```
  Index:     0     1     2     3     4     5     6     7
  nums:     10     9     2     5     3     7   101    18
           ┌─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┐
  dp[i]:   │  1  │  1  │  1  │  2  │  2  │  3  │  4  │  4  │
           └─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┘

  Step-by-step (each dp[i] starts at 1, then checks all j < i):
  ─────────────────────────────────────────────────────────────
  dp[0] = 1  (just [10])

  dp[1] = 1  (9 < 10? No. Just [9])

  dp[2] = 1  (2 < 10? No. 2 < 9? No. Just [2])

  dp[3] = ?  Check j=0: nums[0]=10 > 5 → skip
             Check j=1: nums[1]=9  > 5 → skip
             Check j=2: nums[2]=2  < 5 → dp[3] = max(1, dp[2]+1) = 2
             dp[3] = 2  → subsequence: [2, 5]

  dp[4] = ?  Check j=0: 10 > 3 → skip
             Check j=1: 9  > 3 → skip
             Check j=2: 2  < 3 → dp[4] = max(1, dp[2]+1) = 2
             Check j=3: 5  > 3 → skip
             dp[4] = 2  → subsequence: [2, 3]

  dp[5] = ?  Check j=0: 10 > 7 → skip
             Check j=1: 9  > 7 → skip
             Check j=2: 2  < 7 → dp[5] = max(1, dp[2]+1) = 2
             Check j=3: 5  < 7 → dp[5] = max(2, dp[3]+1) = 3
             Check j=4: 3  < 7 → dp[5] = max(3, dp[4]+1) = 3
             dp[5] = 3  → subsequence: [2, 5, 7] or [2, 3, 7]

  dp[6] = ?  Check j=0..5: 10<101, 9<101, 2<101, 5<101, 3<101, 7<101
             Best: dp[5]+1 = 4
             dp[6] = 4  → subsequence: [2, 5, 7, 101] or [2, 3, 7, 101]

  dp[7] = ?  Check j=0..6: best valid predecessor is j=5 (7 < 18, dp[5]=3)
             dp[7] = dp[5]+1 = 4
             dp[7] = 4  → subsequence: [2, 5, 7, 18] or [2, 3, 7, 18]

  Answer: max(dp[0..7]) = 4
```

**O(n log n) — Patience Sorting visualization on the same input:**

```
  Process each number, maintain a "tails" array where tails[k] is the
  smallest tail of any increasing subsequence of length k+1.

  nums:  10    9    2    5    3    7   101   18
         ──    ──   ──   ──   ──   ──  ───   ──
  tails after each step:

  10  → [10]                   (new pile)
   9  → [9]                    (replace 10: 9 < 10)
   2  → [2]                    (replace 9:  2 < 9)
   5  → [2, 5]                 (extend: 5 > 2)
   3  → [2, 3]                 (replace 5:  3 < 5, bisect to pos 1)
   7  → [2, 3, 7]             (extend: 7 > 3)
  101 → [2, 3, 7, 101]        (extend: 101 > 7)
   18 → [2, 3, 7, 18]         (replace 101: 18 < 101, bisect to pos 3)

  Answer: len(tails) = 4

  Note: tails is NOT the actual LIS — it's a bookkeeping array.
  The length of tails is always the LIS length.
```

---

## 7. Self-Assessment

Answer these without looking at the notes. If you can't, revisit the relevant section.

### Question 1
**What's the difference between overlapping subproblems and optimal substructure? Can you have one without the other?**

<details>
<summary>Answer</summary>

**Overlapping subproblems:** Different branches of the recursion solve the *same* subproblem multiple times. This is what makes caching (memoization) effective.

**Optimal substructure:** The optimal solution to the full problem contains optimal solutions to its subproblems. This is what makes the recurrence *correct* — you can build the optimal answer from optimal sub-answers.

**Can you have one without the other?** Yes.
- *Optimal substructure without overlapping subproblems:* Merge sort. The optimal sort of the whole array uses optimal sorts of the two halves, but the subproblems (different subarrays) don't overlap. Use divide-and-conquer, not DP.
- *Overlapping subproblems without optimal substructure:* Longest simple path in a graph with cycles. Subproblems repeat, but the optimal path from A to C through B doesn't necessarily use the optimal path from A to B (the paths may conflict on shared vertices). DP doesn't apply.

You need **both** for DP to work.
</details>

### Question 2
**When is top-down (memoization) better than bottom-up (tabulation)? When is bottom-up better?**

<details>
<summary>Answer</summary>

**Top-down is better when:**
- You're deriving the solution for the first time (think recursively, add cache).
- The state space is large but sparse — many subproblems never get reached.
- You're in an interview and need to code quickly with fewer bugs.

**Bottom-up is better when:**
- The input is large and deep recursion risks stack overflow.
- You want to optimize space (rolling arrays are straightforward with iteration).
- The state space is dense (you'd compute most states anyway).
- You want maximum performance (no hash map lookups or function call overhead).
</details>

### Question 3
**You're solving a new DP problem. You've spent 15 minutes and can't figure out the recurrence. What should you do?**

<details>
<summary>Answer</summary>

Go back to Step 1: **redefine the state.** If the recurrence isn't coming together, the state definition is almost certainly wrong or incomplete. Common fixes:

- Add more information to the state. Maybe `dp[i]` isn't enough and you need `dp[i][j]` or `dp[i][flag]`.
- Change *what* `dp[i]` represents. Instead of "best answer for first `i` elements," try "best answer *ending at* element `i`" (this is the LIS trick).
- Think about what information you'd need to make the optimal decision at step `i`. If `dp[i]` doesn't give you that information, your state is missing something.

The recurrence is just the math that connects subproblems. If the subproblems (states) are well-defined, the recurrence usually writes itself.
</details>

### Question 4
**For coin change with `coins = [1, 5, 11]` and `amount = 15`, a greedy algorithm (always pick the largest coin) gives 5 coins (11+1+1+1+1). Why does greedy fail, and what does DP give?**

<details>
<summary>Answer</summary>

Greedy fails because the locally optimal choice (pick the largest coin that fits) doesn't lead to the globally optimal solution. Picking the 11-coin forces you into a suboptimal path for the remaining 4 (four 1-coins).

DP gives 3 coins: 5+5+5. The DP table considers *all* ways to make each sub-amount and finds the true minimum. Specifically:
- `dp[15] = min(dp[14]+1, dp[10]+1, dp[4]+1) = min(4+1, 2+1, 4+1) = 3`

This is a great example to remember: **coin change is NOT a greedy problem** in general. It's only greedy for specific coin systems (like US coins: 1, 5, 10, 25).
</details>

### Question 5
**You have a working O(n) space DP solution. `dp[i]` depends on `dp[i-1]`, `dp[i-2]`, and `dp[i-5]`. Can you optimize to O(1) space?**

<details>
<summary>Answer</summary>

No, not easily. You'd need to keep 5 previous values (a sliding window of size 5), which is still O(1) technically — but the general principle is: you can optimize to O(1) when `dp[i]` depends on a *fixed, small* number of prior entries. With a dependency on `dp[i-5]`, you'd need to maintain the last 5 values in a circular buffer. It works here because 5 is a constant, but the approach doesn't scale when the lookback distance is variable (like coin change where it depends on the coin values).

Key distinction: fixed lookback → O(1) possible with rolling variables. Variable lookback → full array needed.
</details>
