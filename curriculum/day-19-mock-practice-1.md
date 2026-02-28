# Day 19: Mixed Pattern Practice — Session 1

> **Time:** 2 hours | **Format:** Timed mock interview | **Language:** Go
>
> You've covered every topic. Today is not about learning — it's about
> **pattern recognition under pressure.** The skill gap between "I know
> the pattern" and "I see the pattern in 60 seconds on a new problem" is
> the gap between studying and passing. This session closes that gap.

---

## Session Format

**4 problems. 25 minutes each. 100 minutes total. 10 minutes for reflection.**

For each problem, follow this exact sequence:

```
0:00 - 2:00   Read the problem. Clarify edge cases out loud.
2:00 - 5:00   Classify the pattern. Decide on your approach.
5:00 - 8:00   Explain the approach out loud (as if to an interviewer).
8:00 - 20:00  Code the solution in Go.
20:00 - 25:00 Test with examples. Fix bugs. State time/space complexity.
```

**Rules — simulate real interview conditions:**

1. **No notes, no references.** Close every browser tab. Close this guide
   after reading each problem statement. Reopen only to check hints or
   the solution after your attempt.
2. **Time yourself.** Use a phone timer. When 25 minutes is up, stop. A
   partial solution is a valid data point — it tells you what to drill.
3. **Talk out loud.** Every thought goes through your mouth. If you're
   silent for more than 30 seconds, you're not practicing the right skill.
   In a real interview, silence is a red flag.
4. **Write on a single file or blank editor.** No autocomplete, no LSP,
   no running the code until you've finished writing. You won't have
   these in a whiteboard interview.

**Score each problem immediately after attempting it** using the rubric at
the bottom. Be brutally honest — inflated scores help no one.

---

## Problem 1: Longest Consecutive Sequence

**Difficulty:** Medium
**Time limit:** 25 minutes
**Topic area:** Arrays & Hashing (Day 1)

### Problem Statement

Given an unsorted array of integers `nums`, find the length of the longest
consecutive elements sequence.

A consecutive sequence is a set of numbers where each number is exactly 1
greater than the previous (e.g., `[100, 101, 102, 103]` has length 4).

You must write an algorithm that runs in **O(n)** time.

**Examples:**

```
Input:  nums = [100, 4, 200, 1, 3, 2]
Output: 4
Explanation: The longest consecutive sequence is [1, 2, 3, 4].

Input:  nums = [0, 3, 7, 2, 5, 8, 4, 6, 0, 1]
Output: 9
Explanation: The longest consecutive sequence is [0, 1, 2, 3, 4, 5, 6, 7, 8].

Input:  nums = []
Output: 0

Input:  nums = [1, 2, 0, 1]
Output: 3
Explanation: [0, 1, 2]. Note the duplicate 1 — it doesn't extend the sequence.
```

**Constraints:**
- `0 <= len(nums) <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This is a **hash set** problem. The O(n) constraint rules out sorting
(O(n log n)). Think about how a hash set lets you check membership in O(1),
and how you can avoid redundant work by only starting a count from the
**beginning** of a sequence.

Maps to: **Day 1 — Arrays & Hashing, Pattern 1/2** (frequency/set-based thinking).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

1. Put all numbers into a hash set for O(1) lookup.
2. For each number `n` in the set, check if `n-1` exists in the set.
   - If `n-1` exists, skip — this number is not the start of a sequence.
   - If `n-1` does NOT exist, `n` is the start of a sequence. Count upward:
     check `n+1`, `n+2`, ... until you find a gap.
3. Track the maximum length across all sequence starts.

The key insight: you only count forward from sequence **starts**. Every
number is visited at most twice (once in the outer loop, once in a counting
chain), so total work is O(n).

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

```go
func longestConsecutive(nums []int) int {
    set := make(map[int]bool, len(nums))
    for _, n := range nums {
        set[n] = true
    }

    best := 0

    for n := range set {
        // Only start counting from the beginning of a sequence
        if set[n-1] {
            continue
        }

        // n is the start of a sequence — count how long it goes
        length := 1
        curr := n
        for set[curr+1] {
            curr++
            length++
        }

        if length > best {
            best = length
        }
    }

    return best
}
```

**Complexity:**
- Time: O(n). Each number is inserted into the set once (O(n)) and visited
  at most twice in the counting phase — once when checking if it's a start,
  and once when it's counted as part of a chain. Total: O(n).
- Space: O(n) for the hash set.

**Edge cases handled:**
- Empty array → `best` stays 0.
- Duplicates → the set deduplicates them automatically.
- Negative numbers → no issue, the logic is value-agnostic.
- Single element → correctly returns 1.

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Pattern** | Hash set for O(1) membership checks |
| **Day** | Day 1 — Arrays & Hashing |
| **Trigger** | "Unsorted array" + "O(n) time" + "consecutive" = can't sort, need set-based lookup |
| **Key insight** | Only start counting from sequence beginnings (no left neighbor in set). This turns a seemingly O(n²) problem into O(n). |
| **Common mistake** | Iterating over `nums` instead of `set` — with duplicates, you do redundant work. Iterating over the set guarantees each unique value is processed once. |

---

## Problem 2: Binary Tree Right Side View

**Difficulty:** Medium
**Time limit:** 25 minutes
**Topic area:** Binary Trees (Day 6)

### Problem Statement

Given the `root` of a binary tree, imagine yourself standing on the **right
side** of it. Return the values of the nodes you can see, ordered from top
to bottom.

In other words, for each level of the tree, return the rightmost node's value.

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}
```

**Examples:**

```
Input:
        1
       / \
      2   3
       \   \
        5   4

Output: [1, 3, 4]
Explanation:
  Level 0: nodes [1]         → rightmost is 1
  Level 1: nodes [2, 3]      → rightmost is 3
  Level 2: nodes [5, 4]      → rightmost is 4


Input:
    1
     \
      3

Output: [1, 3]


Input:
    1
   /
  2

Output: [1, 2]
Explanation: Node 2 is the rightmost (and only) node at level 1.


Input: nil (empty tree)
Output: []
```

**Constraints:**
- The number of nodes is in the range `[0, 100]`.
- `-100 <= Node.Val <= 100`

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This is a **level-order BFS** problem. The phrase "for each level" is the
trigger. You process the tree layer by layer, and at each layer you care
about a specific position (the rightmost node).

Maps to: **Day 6 — Binary Trees, Pattern 3** (Level-Order BFS).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

**BFS approach:**
1. Standard BFS with level-size snapshot (`size := len(queue)` before the
   inner loop).
2. In the inner loop, when `i == size-1`, that node is the rightmost at
   this level. Append its value to the result.

**Alternative DFS approach:**
1. DFS with a "depth" parameter, visiting right child before left.
2. If `depth == len(result)`, this is the first node seen at this depth
   (which is the rightmost, because you visited right first). Append it.

Both are O(n) time. BFS is more intuitive for level-based problems.

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

**BFS solution (recommended):**

```go
func rightSideView(root *TreeNode) []int {
    if root == nil {
        return nil
    }

    var result []int
    queue := []*TreeNode{root}

    for len(queue) > 0 {
        size := len(queue) // snapshot level size
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]

            // Last node in this level is the rightmost
            if i == size-1 {
                result = append(result, node.Val)
            }

            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
    }

    return result
}
```

**DFS solution (alternative):**

```go
func rightSideView(root *TreeNode) []int {
    var result []int
    var dfs func(node *TreeNode, depth int)
    dfs = func(node *TreeNode, depth int) {
        if node == nil {
            return
        }
        // First time reaching this depth → this is the rightmost node
        // (because we visit right before left)
        if depth == len(result) {
            result = append(result, node.Val)
        }
        dfs(node.Right, depth+1) // right first
        dfs(node.Left, depth+1)
    }
    dfs(root, 0)
    return result
}
```

**Complexity (both approaches):**
- Time: O(n) — visit every node exactly once.
- Space: O(w) for BFS where w = max tree width (up to n/2). O(h) for DFS
  where h = tree height.

**Edge cases handled:**
- Empty tree → nil check at the top.
- Left-skewed tree → the leftmost node is also the rightmost at each level.
- Single node → returns `[root.Val]`.

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Pattern** | Level-order BFS |
| **Day** | Day 6 — Binary Trees, Pattern 3 |
| **Trigger** | "Each level" / "right side" / "layer by layer" = BFS |
| **Key insight** | Snapshot `len(queue)` before the inner loop. The last node processed in each level (`i == size-1`) is the rightmost. |
| **Common mistake** | Forgetting the level-size snapshot — processing nodes from mixed levels. Also: assuming you only need to traverse the right subtree (fails when a left child at a deeper level has no right sibling). |

---

## Problem 3: Coin Change

**Difficulty:** Medium
**Time limit:** 25 minutes
**Topic area:** Dynamic Programming (Day 11)

### Problem Statement

You are given an integer array `coins` representing coin denominations and
an integer `amount` representing a total amount of money.

Return the **fewest number of coins** that you need to make up that amount.
If that amount of money cannot be made up by any combination of the coins,
return `-1`.

You may assume that you have an **infinite number** of each coin denomination.

**Examples:**

```
Input:  coins = [1, 5, 11], amount = 15
Output: 3
Explanation: 5 + 5 + 5 = 15 (3 coins).
             Note: greedy would pick 11 + 1 + 1 + 1 + 1 = 5 coins. Greedy fails.

Input:  coins = [2], amount = 3
Output: -1
Explanation: Cannot make 3 with only 2-cent coins.

Input:  coins = [1], amount = 0
Output: 0
Explanation: 0 coins needed for amount 0.

Input:  coins = [1, 2, 5], amount = 11
Output: 3
Explanation: 5 + 5 + 1 = 11 (3 coins).

Input:  coins = [186, 419, 83, 408], amount = 6249
Output: 20
```

**Constraints:**
- `1 <= len(coins) <= 12`
- `1 <= coins[i] <= 2^31 - 1`
- `0 <= amount <= 10^4`

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This is a **1D DP — unbounded choices** problem. The signals: "minimum
number of X to reach a target" + "unlimited use of each option." This is
the exact coin change archetype.

Maps to: **Day 11 — DP 1D, Pattern 3** (Unbounded Choices).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

Apply the 5-step DP recipe:

1. **State:** `dp[i]` = minimum number of coins to make amount `i`.
2. **Recurrence:** `dp[i] = min(dp[i-c] + 1)` for each coin `c` where
   `c <= i` and `dp[i-c]` is reachable (not infinity).
3. **Base case:** `dp[0] = 0` (zero coins for zero amount).
4. **Iteration order:** Left to right (small amounts before large).
5. **Space:** O(amount) — already 1D, can't reduce.

Critical: initialize `dp[1..amount]` to `math.MaxInt32` (not 0). This is a
minimization problem — unreachable states must be infinity.

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

```go
func coinChange(coins []int, amount int) int {
    // Step 1: dp[i] = minimum coins to make amount i
    dp := make([]int, amount+1)

    // Step 3: base case dp[0] = 0; all others = infinity
    for i := 1; i <= amount; i++ {
        dp[i] = math.MaxInt32
    }

    // Step 4: iterate left to right
    for i := 1; i <= amount; i++ {
        // Step 2: try each coin
        for _, c := range coins {
            if c <= i && dp[i-c] != math.MaxInt32 {
                if dp[i-c]+1 < dp[i] {
                    dp[i] = dp[i-c] + 1
                }
            }
        }
    }

    // Return
    if dp[amount] == math.MaxInt32 {
        return -1
    }
    return dp[amount]
}
```

**Requires:** `import "math"`

**Complexity:**
- Time: O(amount × len(coins)). For each of the `amount` states, we check
  each coin. With amount=10^4 and coins up to 12, that's ~120,000 operations.
- Space: O(amount) for the dp array.

**Trace through example 1:** `coins = [1, 5, 11], amount = 15`
```
dp[0]  = 0
dp[1]  = dp[0]+1  = 1
dp[2]  = dp[1]+1  = 2
dp[3]  = dp[2]+1  = 3
dp[4]  = dp[3]+1  = 4
dp[5]  = min(dp[4]+1, dp[0]+1) = 1
dp[10] = min(dp[9]+1, dp[5]+1) = 2
dp[11] = min(dp[10]+1, dp[6]+1, dp[0]+1) = 1
dp[15] = min(dp[14]+1, dp[10]+1, dp[4]+1)
       = min(5, 3, 5) = 3   ← three coins of 5
```

**Edge cases handled:**
- `amount = 0` → returns 0 (base case).
- Impossible amount → dp stays at MaxInt32, returns -1.
- Greedy trap (coins [1,5,11], amount 15) → DP correctly finds 5+5+5, not
  the greedy 11+1+1+1+1.

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Pattern** | 1D DP — Unbounded Choices |
| **Day** | Day 11 — DP 1D, Pattern 3 |
| **Trigger** | "Minimum number of X to reach a target" + "unlimited use of each option" |
| **Key insight** | Greedy fails (coins [1,5,11] amount 15). The DP recurrence tries every coin at each amount and takes the minimum. |
| **Common mistakes** | (1) Initializing dp with 0 instead of infinity — min() returns 0 for everything. (2) Not guarding `dp[i-c] != MaxInt32` before adding 1 — integer overflow. (3) dp array size is `amount+1`, not `amount`. |

---

## Problem 4: Word Search in a Grid

**Difficulty:** Hard (combines two patterns)
**Time limit:** 25 minutes
**Topic area:** Grid Traversal + Backtracking (Days 8 & 13)

### Problem Statement

Given an `m x n` grid of characters `board` and a string `word`, return
`true` if `word` exists in the grid.

The word can be constructed from letters of sequentially adjacent cells,
where adjacent cells are horizontally or vertically neighboring. **The same
cell may not be used more than once** in a single word.

**Examples:**

```
Input:
board = [
  ['A','B','C','E'],
  ['S','F','C','S'],
  ['A','D','E','E'],
]
word = "ABCCED"

Output: true
Explanation: Path: (0,0)→(0,1)→(0,2)→(1,2)→(2,2)→(2,1)

Input:
board = [
  ['A','B','C','E'],
  ['S','F','C','S'],
  ['A','D','E','E'],
]
word = "SEE"

Output: true
Explanation: Path: (1,3)→(2,3)→(2,2)

Input:
board = [
  ['A','B','C','E'],
  ['S','F','C','S'],
  ['A','D','E','E'],
]
word = "ABCB"

Output: false
Explanation: After matching A→B→C, the B at (0,1) is already used.
             No other adjacent B is available from C at (0,2).

Input:
board = [['A']]
word = "A"

Output: true

Input:
board = [['A','B'],['C','D']]
word = "ABDC"

Output: true
Explanation: Path: (0,0)→(0,1)→(1,1)→(1,0)
```

**Constraints:**
- `m == len(board)`, `n == len(board[0])`
- `1 <= m, n <= 6`
- `1 <= len(word) <= 15`
- `board` and `word` consist of only lowercase and uppercase English letters.

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This combines **grid DFS** (Day 8, Pattern 1) with **backtracking** (Day 13).
You need to explore all four directions from each cell, track visited cells
to avoid reuse, and backtrack (unmark visited) when a path doesn't work out.

Maps to: **Day 8 — Grid DFS** + **Day 13 — Backtracking** (constraint-based
search with undo).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

1. For each cell `(r, c)` in the grid, if `board[r][c] == word[0]`, start
   a DFS/backtracking search from that cell.
2. The DFS function takes `(r, c, index)` where `index` is the current
   position in the word we're trying to match.
3. Base case: if `index == len(word)`, we've matched everything → return true.
4. Boundary/validity check: out of bounds, already visited, or
   `board[r][c] != word[index]` → return false.
5. Mark the cell as visited (temporarily modify `board[r][c]` to a sentinel
   like `'#'`).
6. Recurse in all 4 directions with `index+1`.
7. **Backtrack:** restore `board[r][c]` to its original value.
8. Return true if any recursive call returned true.

The "mark and restore" step is the backtracking piece. Without it, you'd
incorrectly block valid paths that reuse a cell from a different search branch.

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

```go
func exist(board [][]byte, word string) bool {
    rows, cols := len(board), len(board[0])

    var dfs func(r, c, idx int) bool
    dfs = func(r, c, idx int) bool {
        // All characters matched
        if idx == len(word) {
            return true
        }

        // Boundary check + character match
        if r < 0 || r >= rows || c < 0 || c >= cols {
            return false
        }
        if board[r][c] != word[idx] {
            return false
        }

        // Mark visited (in-place, no extra space)
        original := board[r][c]
        board[r][c] = '#'

        // Explore all 4 directions
        found := dfs(r+1, c, idx+1) ||
            dfs(r-1, c, idx+1) ||
            dfs(r, c+1, idx+1) ||
            dfs(r, c-1, idx+1)

        // Backtrack: restore the cell
        board[r][c] = original

        return found
    }

    // Try starting from every cell
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if dfs(r, c, 0) {
                return true
            }
        }
    }
    return false
}
```

**Complexity:**
- Time: O(m × n × 3^L) where L = len(word). From each cell, we branch into
  at most 3 directions (not 4 — we came from one direction). At each step we
  go up to L levels deep. In the worst case we start from every cell.
- Space: O(L) for the recursion stack. We modify the board in place for
  visited tracking, so no extra space for a visited set.

**Why 3^L, not 4^L:** After the first step, the cell we came from is marked
visited, so we only branch into 3 neighbors, not 4. This is a meaningful
constant factor reduction.

**Trace through example 1:** `board`, `word = "ABCCED"`
```
Start at (0,0), board[0][0]='A' == word[0]. Mark '#'.
  → (0,1), 'B' == word[1]. Mark '#'.
    → (0,2), 'C' == word[2]. Mark '#'.
      → (1,2), 'C' == word[3]. Mark '#'.
        → (2,2), 'E' == word[4]. Mark '#'.
          → (2,1), 'D' == word[5]. Mark '#'.
            → idx=6 == len(word). Return true!
```

**Edge cases handled:**
- Single cell board matching single char word → base case triggers immediately.
- Word longer than total cells → DFS naturally fails (can't revisit cells).
- Short-circuit with `||` → stops exploring once a valid path is found.

**Optional optimization (mention to interviewer):**
Before starting the search, check that the board contains all characters
in the word with sufficient frequency. This prunes impossible cases early:

```go
// Quick frequency check (optional pruning)
freq := make(map[byte]int)
for _, row := range board {
    for _, ch := range row {
        freq[ch]++
    }
}
for i := 0; i < len(word); i++ {
    freq[word[i]]--
    if freq[word[i]] < 0 {
        return false
    }
}
```

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Patterns** | Grid DFS + Backtracking (two patterns combined) |
| **Days** | Day 8 — Grid DFS (Pattern 1) + Day 13 — Backtracking |
| **Trigger** | "Grid" + "path of adjacent cells" + "same cell may not be used twice" = grid DFS with visited tracking. "Find if a path exists" + constraint on reuse = backtracking (mark → explore → unmark). |
| **Key insight** | In-place modification of the board (`'#'` sentinel) serves as both the visited set and the backtracking mechanism. Restoring the cell after recursion is the "undo" step that makes backtracking work. |
| **Common mistakes** | (1) Forgetting to backtrack — leaving cells marked as visited, which blocks valid paths in other branches. (2) Checking `board[r][c] != word[idx]` after marking visited — the order matters (check match, then mark). (3) Not short-circuiting with `||` — exploring all 4 directions even after finding a valid path wastes time. |

---

## Scoring Rubric

**Score each problem immediately after your 25 minutes are up.** Be honest.

| Criterion | Points | How to Score |
|-----------|--------|-------------|
| **Pattern identified within 5 min** | +2 | Did you name the correct pattern/data structure before minute 5? |
| **Correct approach within 10 min** | +2 | Could you explain the full algorithm (recurrence, traversal type, etc.) before minute 10? |
| **Working solution within 25 min** | +3 | Does your code produce correct output for all given examples? Partial credit: +1 if the logic is right but has 1-2 bugs. |
| **Clean code (no bugs on first test)** | +1 | When you traced through your first example, did it work without edits? |
| **Optimal time complexity** | +1 | Is your solution the intended complexity? (O(n) for P1, O(n) for P2, O(amount×coins) for P3, O(mn×3^L) for P4) |
| **Edge cases handled** | +1 | Did you explicitly handle or mention: empty input, single element, impossible cases, duplicates? |

**Maximum: 10 points per problem. 40 points total.**

### Score Interpretation

| Score | Assessment | Action |
|-------|-----------|--------|
| 36-40 | Interview ready | You're sharp. Focus Day 20-21 on speed and communication polish. |
| 28-35 | Strong foundation, gaps in execution | Review the problems you dropped points on. The pattern recognition is there; drill the implementation. |
| 20-27 | Pattern recognition needs work | Go back to the specific days for topics you missed. Redo those days' drills. |
| Below 20 | Core gaps remain | Focus Day 20 on re-drilling the weakest 2-3 topics rather than doing another mixed session. |

---

## Scorecard

Copy this and fill it in:

```
Problem 1 — Longest Consecutive Sequence
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:     ___ / 2
  Working solution:       ___ / 3
  Clean first test:       ___ / 1
  Optimal complexity:     ___ / 1
  Edge cases:             ___ / 1
  Subtotal:               ___ / 10
  Time used:              ___ min

Problem 2 — Binary Tree Right Side View
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:     ___ / 2
  Working solution:       ___ / 3
  Clean first test:       ___ / 1
  Optimal complexity:     ___ / 1
  Edge cases:             ___ / 1
  Subtotal:               ___ / 10
  Time used:              ___ min

Problem 3 — Coin Change
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:     ___ / 2
  Working solution:       ___ / 3
  Clean first test:       ___ / 1
  Optimal complexity:     ___ / 1
  Edge cases:             ___ / 1
  Subtotal:               ___ / 10
  Time used:              ___ min

Problem 4 — Word Search
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:     ___ / 2
  Working solution:       ___ / 3
  Clean first test:       ___ / 1
  Optimal complexity:     ___ / 1
  Edge cases:             ___ / 1
  Subtotal:               ___ / 10
  Time used:              ___ min

TOTAL:                    ___ / 40
```

---

## Reflection Template

**Complete this after the session. Spend the full 10 minutes. Writing forces
clarity — don't skip this.**

### 1. Which problem took the longest? Why?

```
Problem #___. Time used: ___ min.
Root cause (circle one):
  [ ] Didn't recognize the pattern
  [ ] Recognized pattern but couldn't translate to code
  [ ] Had the code but got stuck on bugs/edge cases
  [ ] Ran out of time during implementation

Specific detail:


```

### 2. Which pattern did you recognize fastest?

```
Problem #___. Recognized in ~___ seconds.
What was the trigger word/phrase in the problem statement?


```

### 3. Where did you get stuck — recognition, implementation, or edge cases?

```
Breakdown across all 4 problems:
  Pattern recognition:  ___ problems with no hesitation
  Implementation:       ___ problems where I knew what to do but struggled to code it
  Edge cases:           ___ problems where my first code missed a case

Weakest area:


```

### 4. What topics need more drilling before Day 21?

```
Based on today's performance, I need to revisit:
  1.
  2.
  3.

For each: which day's material should I re-read?


```

### 5. If this were a real interview, would I have passed?

```
Realistically (be honest):
  Problem 1: [ ] Passed  [ ] Struggled  [ ] Failed
  Problem 2: [ ] Passed  [ ] Struggled  [ ] Failed
  Problem 3: [ ] Passed  [ ] Struggled  [ ] Failed
  Problem 4: [ ] Passed  [ ] Struggled  [ ] Failed

A typical bar: solve 2 mediums cleanly in 45 min, or 1 medium + meaningful
progress on a hard. Did you meet that bar?


```

---

## Session Schedule (2 Hours)

| Time | Activity |
|------|----------|
| 0:00 - 0:02 | Read the session format and rules above. Set up your timer and blank editor. |
| 0:02 - 0:27 | **Problem 1** — Longest Consecutive Sequence (25 min) |
| 0:27 - 0:29 | Score Problem 1. Briefly review the solution if needed. Reset your editor. |
| 0:29 - 0:54 | **Problem 2** — Binary Tree Right Side View (25 min) |
| 0:54 - 0:56 | Score Problem 2. Brief review. Reset. |
| 0:56 - 1:21 | **Problem 3** — Coin Change (25 min) |
| 1:21 - 1:23 | Score Problem 3. Brief review. Reset. |
| 1:23 - 1:48 | **Problem 4** — Word Search (25 min) |
| 1:48 - 1:50 | Score Problem 4. Brief review. |
| 1:50 - 2:00 | **Reflection.** Fill out the entire reflection template. Identify your Day 20 focus areas. |

---

*Tomorrow (Day 20): Mixed Pattern Practice — Session 2. Different problems,
same format. Your Day 20 performance minus your Day 19 performance is your
rate of improvement.*
