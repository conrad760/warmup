# Day 13 — DFS & Backtracking: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Focus | Time |
|---|----------|-------|------|
| 1 | [DFS vs BFS Side-by-Side — VisuAlgo](https://visualgo.net/en/dfsbfs) | Interactive visualization where you step through DFS and BFS on the same graph. Toggle between them to see the stack vs queue, the traversal order differences, and how DFS dives deep while BFS goes wide. Essential for building intuition about when each shines. | 10 min |
| 2 | [Graph Algorithms for Technical Interviews — freeCodeCamp](https://www.youtube.com/watch?v=tWVWeAqZ0WU) | Watch from 20:00 onward for the DFS walkthrough. Shows recursive and iterative DFS with clear stack traces. Covers connected components and path existence using DFS. Good contrast with the BFS coverage from the first 20 minutes (Day 12). | 15 min (20:00–35:00) |
| 3 | [Backtracking — Abdul Bari](https://www.youtube.com/watch?v=DKCbsiDBN6c) | Whiteboard explanation of the backtracking decision tree. Walks through the "choose-explore-unchoose" framework step by step. Draws the full state-space tree for subsets and permutations, showing where branches get pruned. Best conceptual foundation. | 18 min |
| 4 | [Subsets / Permutations / Combinations — NeetCode](https://www.youtube.com/watch?v=REOH22Xwdkk) | All three backtracking patterns in one video with decision tree diagrams for each. Shows the "include or exclude" model for subsets vs the "pick from remaining" model for permutations. Go-friendly pseudocode that maps directly to recursive implementations. | 15 min |
| 5 | [Word Search Backtracking — NeetCode](https://www.youtube.com/watch?v=pfiQ_PS1g8E) | Grid-based backtracking with visited cell marking and restoration. Shows how DFS on a 2D grid becomes backtracking when you need to undo visited markers. The in-place `board[r][c] = '#'` trick and why you must restore it. | 10 min |
| 6 | [Cycle Detection in Directed Graphs — William Fiset](https://www.youtube.com/watch?v=rKQaZuoUR4M) | Animated three-color DFS (white/gray/black) on a directed graph. Visualizes why a back edge to a gray node proves a cycle exists. Contrasts with undirected cycle detection. Clear step-by-step with the color state at each moment. | 12 min |
| 7 | [Go Closures and Recursive Patterns — Go Blog](https://go.dev/blog/func) | Understanding closures in Go is critical for backtracking. This covers how anonymous functions capture variables by reference, which is exactly how backtracking helpers typically accumulate results via a `result` slice in the enclosing scope. Short and practical. | 8 min |

---

## 2. Detailed 2-Hour Session Plan

### 12:00 — 12:20 | Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:06 | Read Section 3.1 (DFS: recursive vs iterative) and Section 3.2 (three-color marking). On paper, trace the DFS in Section 6.1 — write down the stack state, discovery order, and finish order for each node. Verify your trace matches the diagram. |
| 12:06 - 12:12 | Read Section 3.3 (DFS vs BFS comparison table). For each row, say out loud why that property holds. Then read Section 3.4 (backtracking framework). Study the "choose → explore → unchoose" template. On paper, draw the first 3 levels of the subset decision tree for `[1,2,3]` (Section 6.2). |
| 12:12 - 12:18 | Read Section 3.5 (pruning strategies) and Section 3.6 (subsets vs permutations) and Section 3.7 (handling duplicates). For duplicates, trace through the example: `[1,2,2]` — which branches get skipped and why. Study the diagram in Section 6.4. |
| 12:18 - 12:20 | Review the complexity table. Say out loud: "DFS is O(V+E). Subsets is O(n * 2^n). Permutations is O(n * n!). Backtracking is exponential — pruning is the only way to make it practical." |

### 12:20 — 1:20 | Implement (From Scratch)

| Time | Problem | Notes |
|------|---------|-------|
| 12:20 - 12:32 | `DFS` (graph traversal) | Implement recursive DFS returning traversal order. Use a `visited []bool`. Then implement iterative DFS with an explicit `[]int` stack. Test: linear chain, graph with cycle, disconnected component (only reachable nodes returned), single node. Compare output order with your BFS from Day 12 on the same graph. |
| 12:32 - 12:44 | `HasCycleDirected` | Three-color DFS. Use `color []int` with 0=white, 1=gray, 2=black. If you encounter a gray neighbor, return true. Test: DAG (no cycle), simple cycle A→B→C→A, self-loop, diamond graph (A→B, A→C, B→D, C→D — no cycle), disconnected graph with a cycle in one component. |
| 12:44 - 12:56 | `Subsets` | Generate all subsets of `[]int`. Use the include/exclude pattern. Copy the `current` slice before appending to `result`. Test: empty input → `[[]]`, `[1]` → `[[], [1]]`, `[1,2,3]` → 8 subsets. Verify count is 2^n. |
| 12:56 - 13:08 | `Permutations` | Generate all permutations. First implement with a `used []bool` array. Then implement the swap-based approach. Test: `[1]` → `[[1]]`, `[1,2,3]` → 6 permutations. Verify count is n!. Try both implementations and confirm identical output (order may differ). |
| 13:08 - 13:20 | `WordSearch` | Backtracking DFS on a 2D grid. Mark visited by setting `board[r][c] = '#'`, restore on backtrack. Try all 4 directions. Test: word exists, word doesn't exist, word longer than grid, single-character word, word snakes back near start but can't reuse a cell. |

### 1:20 — 1:50 | Solidify (Edge Cases & Variants)

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | `SubsetsWithDuplicates`: Modify `Subsets` to handle `[1,2,2]`. Sort first, skip when `nums[i] == nums[i-1]` and `i > start`. Test: `[1,2,2]` → 6 subsets (not 8). Verify no duplicate subsets. Then implement `CombinationSum`: given candidates and target, find all unique combinations summing to target. Candidates may be reused. Test: `[2,3,6,7]` target 7 → `[[2,2,3], [7]]`. |
| 1:30 - 1:40 | Stress test `HasCycleDirected` on a larger graph. Build a 6-node graph with multiple components, one containing a cycle. Verify cycle detection works. Then test on a complete DAG (topological ordering exists) — should return false. Try adding a single back edge and verify it flips to true. |
| 1:40 - 1:50 | Review all implementations. Extract the common backtracking skeleton across `Subsets`, `Permutations`, `CombinationSum`, and `WordSearch`. They all follow "choose → recurse → unchoose." Write a comment in each function labeling the three steps. Verify every backtracking function properly copies slices before saving to results and restores state on every path (including early returns in `WordSearch`). |

### 1:50 — 2:00 | Recap (From Memory)

Write down without looking:
1. The DFS algorithm (recursive) in pseudocode. Then the iterative version with explicit stack.
2. The three-color marking scheme: what does each color mean, and what constitutes a cycle detection (back edge to gray).
3. The backtracking template: choose → explore → unchoose. Why the undo step is essential.
4. The difference between the subsets decision tree (include/exclude at each index) and the permutations decision tree (pick from remaining at each level).
5. How to handle duplicates in subsets: sort + skip condition.

---

## 3. Core Concepts Deep Dive

### 3.1 DFS: Recursive vs Iterative

DFS explores as deep as possible along each branch before backtracking. It uses a stack — either the call stack (recursive) or an explicit stack (iterative).

**Recursive DFS:**

```go
func DFSRecursive(graph [][]int, start int) []int {
    visited := make([]bool, len(graph))
    order := []int{}

    var dfs func(node int)
    dfs = func(node int) {
        visited[node] = true
        order = append(order, node)
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                dfs(neighbor)
            }
        }
    }

    dfs(start)
    return order
}
```

**Iterative DFS:**

```go
func DFSIterative(graph [][]int, start int) []int {
    visited := make([]bool, len(graph))
    stack := []int{start}
    order := []int{}

    for len(stack) > 0 {
        // pop from stack
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if visited[node] {
            continue
        }
        visited[node] = true
        order = append(order, node)

        // push neighbors in reverse order to visit in original order
        for i := len(graph[node]) - 1; i >= 0; i-- {
            neighbor := graph[node][i]
            if !visited[neighbor] {
                stack = append(stack, neighbor)
            }
        }
    }
    return order
}
```

**When to use each:**

| | Recursive | Iterative |
|---|-----------|-----------|
| **Simplicity** | Simpler, more natural for tree/graph DFS | More boilerplate |
| **Stack overflow** | Risk on deep graphs (Go default goroutine stack starts at 8KB, grows to 1GB, but very deep recursion is still risky) | No stack overflow risk — heap-allocated stack |
| **State management** | Natural for backtracking (function scope preserves state) | Manual state management on the stack |
| **When to prefer** | Trees, backtracking, typical interview problems | Very deep graphs (100K+ depth), when you need explicit control over the stack |

**Important difference in traversal order:** The iterative version may visit neighbors in a different order than the recursive version depending on how you push to the stack. In the recursive version, you iterate neighbors left-to-right, visiting the first neighbor's subtree completely before the second. In the iterative version with a simple stack, you push all neighbors at once. To get the same order, push neighbors in reverse order (last neighbor pushed first so it's popped first).

---

### 3.2 Three-Color Marking for Directed Cycle Detection

In directed graphs, a simple `visited` boolean is insufficient for cycle detection. Consider:

```
A → B → C
A → C
```

When DFS from A visits C (via A→B→C), C is marked visited. Later, when DFS explores A→C, it sees C is visited — but this is NOT a cycle. It's just two paths to the same node. A simple visited check would incorrectly report a cycle.

**Three-color marking** fixes this by distinguishing between nodes that are on the current DFS path (in-progress) and nodes that have been fully explored (completed):

| Color | Meaning | Value |
|-------|---------|-------|
| **White** (unvisited) | Not yet discovered | 0 |
| **Gray** (in-progress) | Currently on the DFS recursion stack; we're exploring its descendants | 1 |
| **Black** (completed) | Fully processed; all descendants have been explored | 2 |

**The rule:** If DFS encounters a **gray** node, it means we've found a back edge — a path from a node back to one of its ancestors in the current DFS tree. This is a cycle.

Encountering a **black** node is fine — it's a cross-edge or forward-edge, not a cycle.

```go
func HasCycleDirected(graph [][]int, n int) bool {
    color := make([]int, n) // 0=white, 1=gray, 2=black

    var dfs func(node int) bool
    dfs = func(node int) bool {
        color[node] = 1 // mark gray (entering)

        for _, neighbor := range graph[node] {
            if color[neighbor] == 1 {
                return true // back edge to gray node = CYCLE
            }
            if color[neighbor] == 0 {
                if dfs(neighbor) {
                    return true
                }
            }
            // color[neighbor] == 2 (black) → already fully explored, skip
        }

        color[node] = 2 // mark black (leaving)
        return false
    }

    for i := 0; i < n; i++ {
        if color[i] == 0 {
            if dfs(i) {
                return true
            }
        }
    }
    return false
}
```

**Why two colors (simple visited) fail for directed graphs:** With a simple boolean `visited`, once a node is visited, you can't distinguish between "this node is an ancestor on my current path" (cycle!) and "this node was visited by a completely different DFS path" (no cycle). The gray/black distinction encodes exactly this.

**For undirected graphs**, simple visited + parent tracking works because any visited neighbor that isn't your direct parent forms a cycle. There are no "cross edges" in undirected DFS.

---

### 3.3 DFS vs BFS: Comparison Table

| Property | DFS | BFS |
|----------|-----|-----|
| **Data structure** | Stack (call stack or explicit) | Queue |
| **Exploration pattern** | Goes deep first, backtracks | Goes wide first, layer by layer |
| **Finds shortest path (unweighted)?** | No — finds *a* path, not necessarily shortest | Yes — first visit is shortest path |
| **Space complexity** | O(V) worst case, O(h) for trees (h = height) | O(V) worst case, O(w) for trees (w = max width) |
| **Better for deep, narrow graphs** | Yes — uses O(h) space on the stack | No — BFS would still use O(V) for the queue |
| **Better for wide, shallow graphs** | No — DFS may go unnecessarily deep | Yes — BFS stays at shallow depths |
| **Cycle detection (directed)** | Three-color marking | Kahn's algorithm (in-degree method) |
| **Topological sort** | DFS + reverse post-order | Kahn's BFS |
| **Path finding** | Natural — recursion builds the path | Requires backtracking through `parent[]` array |
| **Connected components** | DFS from each unvisited node | BFS from each unvisited node (both work) |
| **Use for backtracking** | Yes — DFS IS the recursion engine for backtracking | No — BFS doesn't backtrack |

**Rules of thumb:**
- Need **shortest path** (unweighted)? → BFS
- Need to **explore all paths** or **generate combinatorial objects**? → DFS/Backtracking
- Need **cycle detection** in a directed graph? → DFS (three-color)
- Need **topological sort**? → Either works (Kahn's BFS or DFS post-order)
- Working with **trees**? → DFS is usually simpler (pre/in/post-order are DFS)

---

### 3.4 Backtracking Framework: Choose → Explore → Unchoose

Backtracking is DFS on a **decision tree** (also called a state-space tree). At each node in the tree, you:

1. **Choose** — Make a decision that modifies the current state.
2. **Explore** — Recurse to explore consequences of that decision.
3. **Unchoose** — Undo the decision, restoring state to what it was before.

The "unchoose" step is what makes it backtracking rather than plain DFS. Without it, the state accumulates and every path sees modifications from all previous paths.

**Generic template:**

```go
func backtrack(result *[][]int, current []int, choices []int, start int) {
    if isGoalState(current) {
        // save a COPY of current to result
        tmp := make([]int, len(current))
        copy(tmp, current)
        *result = append(*result, tmp)
        return
    }

    for i := start; i < len(choices); i++ {
        if !isValid(choices[i]) {
            continue // prune
        }

        // CHOOSE: modify state
        current = append(current, choices[i])

        // EXPLORE: recurse
        backtrack(result, current, choices, i+1) // or i for reuse, or 0 for permutations

        // UNCHOOSE: restore state
        current = current[:len(current)-1]
    }
}
```

**Why undoing the choice is essential:** Backtracking explores multiple branches of the decision tree. Each branch should start from the same state. If you add element X to your current path and recurse, when you return from that recursion, the current path still contains X. If you don't remove X before trying the next branch (element Y), you'll end up with both X and Y in the path — giving wrong results.

**Go closure pattern:** In Go, it's idiomatic to use a closure that captures `result` from the enclosing scope, avoiding the need to pass `*[][]int` explicitly:

```go
func Subsets(nums []int) [][]int {
    result := [][]int{}
    current := []int{}

    var backtrack func(start int)
    backtrack = func(start int) {
        // save copy
        tmp := make([]int, len(current))
        copy(tmp, current)
        result = append(result, tmp)

        for i := start; i < len(nums); i++ {
            current = append(current, nums[i])     // choose
            backtrack(i + 1)                         // explore
            current = current[:len(current)-1]       // unchoose
        }
    }

    backtrack(0)
    return result
}
```

---

### 3.5 Pruning Strategies

Backtracking is inherently exponential. Pruning cuts branches of the decision tree early, avoiding exploration of states that cannot lead to valid solutions.

**Strategy 1: Skip invalid branches immediately**

Before making a choice, check if it can possibly lead to a valid solution. If not, `continue`.

```go
// Combination Sum: skip if adding this candidate exceeds target
if current_sum + candidates[i] > target {
    continue // or break if sorted (all subsequent are larger)
}
```

**Strategy 2: Sort to enable pruning**

Sorting the input allows you to break out of loops early when values exceed a bound.

```go
sort.Ints(candidates)
for i := start; i < len(candidates); i++ {
    if candidates[i] > remaining {
        break // all subsequent candidates are also too large
    }
    // ...
}
```

**Strategy 3: Use constraints to bound the search**

In N-Queens, you don't try all N^N placements. You place one queen per row, and for each placement you check column conflicts and diagonal conflicts. This reduces the search space dramatically.

```go
// Instead of trying all positions:
for col := 0; col < n; col++ {
    if cols[col] || diag1[row-col+n-1] || diag2[row+col] {
        continue // prune: column or diagonal already attacked
    }
    // ... place queen and recurse
}
```

**Strategy 4: Skip duplicates**

When the input contains duplicates and you want unique results, sort first and skip consecutive equal elements (see Section 3.7).

---

### 3.6 Subsets vs Permutations: Two Different Decision Trees

These are the two fundamental backtracking patterns. Understanding their decision trees is the key to recognizing which pattern to apply.

**Subsets — "Include or exclude at each position"**

At each index, you decide: include this element or skip it. The tree has depth n (one decision per element) and 2^n leaves (each leaf is a subset).

- The `start` parameter advances forward: once you've decided about element `i`, you only consider elements `i+1, i+2, ...`
- Order doesn't matter: `{1,2}` and `{2,1}` are the same subset
- Every node in the tree is a valid result (not just leaves)

**Permutations — "Choose an unused element at each position"**

At each level, you choose which element to place at this position. The tree has depth n (one element per position) and n! leaves (each leaf is a complete permutation).

- At each level, you can pick any unused element: iterate from 0 to n-1, skip those already used
- Order matters: `[1,2,3]` and `[3,2,1]` are different permutations
- Only leaves (complete permutations of length n) are valid results

| | Subsets | Permutations |
|---|---------|-------------|
| **Decision** | Include or skip element at index i | Pick which unused element goes at position i |
| **Tree depth** | n | n |
| **Tree leaves** | 2^n | n! |
| **Valid results** | Every node (partial subsets are valid) | Only complete leaves |
| **Loop start** | `i = start` (move forward) | `i = 0` (consider all, skip used) |
| **Time complexity** | O(n * 2^n) | O(n * n!) |

---

### 3.7 Handling Duplicates in Backtracking

When the input has duplicates (e.g., `[1, 2, 2]`), naive backtracking generates duplicate results. The fix:

1. **Sort the input** so duplicates are adjacent.
2. **Skip an element if it equals the previous element and we didn't use the previous element at this level.**

For subsets:

```go
sort.Ints(nums)
for i := start; i < len(nums); i++ {
    // skip duplicate at the same decision level
    if i > start && nums[i] == nums[i-1] {
        continue
    }
    current = append(current, nums[i])
    backtrack(i + 1)
    current = current[:len(current)-1]
}
```

**Why `i > start` and not `i > 0`:** The condition must only skip at the current branching level. When `i == start`, we're at the first choice of this level — we must consider it. When `i > start` and `nums[i] == nums[i-1]`, we've already explored a branch starting with this same value at this level — choosing it again would produce a duplicate subset.

**Example with `[1, 2, 2]`:**

Without skipping: subsets include both `[1, 2(a)]` and `[1, 2(b)]` — duplicates.

With skipping: at the level where `start = 1`, we try `2(a)` (index 1). When the loop advances to index 2, `nums[2] == nums[1]` and `2 > start=1`, so we skip `2(b)`. This eliminates the duplicate.

For permutations with duplicates, the skip condition is slightly different — use a `used` array and skip when `nums[i] == nums[i-1] && !used[i-1]`. This ensures only one ordering of identical elements is explored.

---

## 4. Implementation Checklist

### Function Signatures

```go
package dfs

// DFS returns the traversal order from start using recursive DFS.
func DFS(graph [][]int, start int, visited []bool) []int { ... }

// DFSIterative returns the traversal order from start using an explicit stack.
func DFSIterative(graph [][]int, start int) []int { ... }

// HasCycleDirected returns true if the directed graph contains a cycle.
// Uses three-color (white/gray/black) DFS.
func HasCycleDirected(graph [][]int, n int) bool { ... }

// Subsets returns all subsets of nums.
func Subsets(nums []int) [][]int { ... }

// SubsetsWithDuplicates returns all unique subsets of nums (may contain duplicates).
func SubsetsWithDuplicates(nums []int) [][]int { ... }

// Permutations returns all permutations of nums (all distinct).
func Permutations(nums []int) [][]int { ... }

// CombinationSum returns all unique combinations of candidates that sum to target.
// Each candidate may be reused unlimited times.
func CombinationSum(candidates []int, target int) [][]int { ... }

// WordSearch returns true if word exists in the board as a path of adjacent cells.
func WordSearch(board [][]byte, word string) bool { ... }
```

### Test Cases & Edge Cases

| Function | Must-Test Cases |
|----------|----------------|
| `DFS` | Linear chain 0→1→2→3 from node 0; complete graph (visits all); graph with cycle (doesn't infinite loop — visited check); disconnected graph (only reachable nodes); single node → `[0]`; compare order with iterative version. |
| `HasCycleDirected` | DAG (no cycle) → `false`; simple cycle A→B→C→A → `true`; self-loop A→A → `true`; diamond DAG A→B, A→C, B→D, C→D → `false`; disconnected graph with cycle in one component → `true`; single node no edges → `false`. |
| `Subsets` | Empty input → `[[]]`; `[1]` → `[[], [1]]` (2 subsets); `[1,2,3]` → 8 subsets; verify no duplicates; verify each subset is sorted (since input is processed in order). |
| `SubsetsWithDuplicates` | `[1,2,2]` → 6 unique subsets; `[0]` → `[[], [0]]`; `[1,1,1]` → `[[], [1], [1,1], [1,1,1]]` (4 subsets, not 8). |
| `Permutations` | `[1]` → `[[1]]`; `[1,2]` → `[[1,2], [2,1]]`; `[1,2,3]` → 6 permutations; verify each permutation has length n; verify no duplicates. |
| `CombinationSum` | `[2,3,6,7]` target 7 → `[[2,2,3], [7]]`; target 1 with `[2]` → `[]` (impossible); target 0 → `[[]]`; single candidate equals target → `[[candidate]]`. |
| `WordSearch` | Word exists as straight line → `true`; word requires turns → `true`; word not in board → `false`; word longer than total cells → `false`; single cell board, single char word → `true` or `false`; word path must not reuse a cell. |

---

## 5. Backtracking Template Library

### Template 1: Subsets (Include/Exclude)

```go
func Subsets(nums []int) [][]int {
    result := [][]int{}
    current := []int{}

    var backtrack func(start int)
    backtrack = func(start int) {
        // Every partial combination is a valid subset — record it
        tmp := make([]int, len(current))
        copy(tmp, current)
        result = append(result, tmp)

        for i := start; i < len(nums); i++ {
            current = append(current, nums[i])  // ← CHOOSE: include nums[i]
            backtrack(i + 1)                      // ← EXPLORE: recurse with next index
            current = current[:len(current)-1]    // ← UNCHOOSE: exclude nums[i]
        }
    }

    backtrack(0)
    return result
}
```

**How it works:** At each call, we record the current state as a valid subset. Then for each remaining element (from `start` to end), we include it and recurse. After returning, we remove it. The `start` parameter ensures we only move forward, avoiding duplicate subsets like `{1,2}` and `{2,1}`.

---

### Template 2a: Permutations (Used-Array Based)

```go
func Permutations(nums []int) [][]int {
    result := [][]int{}
    current := []int{}
    used := make([]bool, len(nums))

    var backtrack func()
    backtrack = func() {
        if len(current) == len(nums) {
            // Only complete permutations are valid — record at leaf
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }

        for i := 0; i < len(nums); i++ {
            if used[i] {
                continue // skip already-placed elements
            }

            used[i] = true                         // ← CHOOSE: mark as used
            current = append(current, nums[i])     // ← CHOOSE: place nums[i]

            backtrack()                             // ← EXPLORE

            current = current[:len(current)-1]     // ← UNCHOOSE: remove nums[i]
            used[i] = false                         // ← UNCHOOSE: mark as available
        }
    }

    backtrack()
    return result
}
```

**How it works:** At each level, we iterate all elements and pick any that isn't already used. The `used` array prevents reusing the same element. Both the append and the `used` flag must be undone during unchoose.

---

### Template 2b: Permutations (Swap-Based)

```go
func PermutationsSwap(nums []int) [][]int {
    result := [][]int{}

    var backtrack func(start int)
    backtrack = func(start int) {
        if start == len(nums) {
            tmp := make([]int, len(nums))
            copy(tmp, nums)
            result = append(result, tmp)
            return
        }

        for i := start; i < len(nums); i++ {
            nums[start], nums[i] = nums[i], nums[start]  // ← CHOOSE: swap into position
            backtrack(start + 1)                           // ← EXPLORE: fix this position, recurse
            nums[start], nums[i] = nums[i], nums[start]  // ← UNCHOOSE: swap back
        }
    }

    backtrack(0)
    return result
}
```

**How it works:** Instead of tracking `used`, we swap elements into the current position. Position `start` is being decided: try every element from `start` to end by swapping it into position `start`, then recurse for `start+1`. The swap-back restores the original order. More memory-efficient (no `used` array or `current` slice), but harder to extend for deduplication.

---

### Template 3: Combinations (k elements from n)

```go
func Combinations(n, k int) [][]int {
    result := [][]int{}
    current := []int{}

    var backtrack func(start int)
    backtrack = func(start int) {
        if len(current) == k {
            // Reached desired size — record
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }

        // Pruning: need (k - len(current)) more elements,
        // and only (n - i + 1) are available from i onward.
        // Stop when there aren't enough elements left.
        for i := start; i <= n-(k-len(current))+1; i++ {
            current = append(current, i)          // ← CHOOSE: include i
            backtrack(i + 1)                       // ← EXPLORE
            current = current[:len(current)-1]     // ← UNCHOOSE: exclude i
        }
    }

    backtrack(1) // 1-indexed: choosing from {1, 2, ..., n}
    return result
}
```

**How it works:** Same structure as subsets but with two differences: (1) we only record when we've picked exactly k elements, and (2) we prune the loop upper bound — if we need 3 more elements but only 2 remain, we stop early. The pruning bound `n - (k - len(current)) + 1` is the latest we can start and still fill k positions.

---

### Template 4: Constraint Satisfaction (N-Queens Style)

```go
func SolveNQueens(n int) [][]string {
    result := [][]string{}
    // Track which columns and diagonals are under attack
    cols := make([]bool, n)
    diag1 := make([]bool, 2*n-1) // row - col + (n-1) mapped to [0, 2n-2]
    diag2 := make([]bool, 2*n-1) // row + col mapped to [0, 2n-2]
    board := make([]int, n)       // board[row] = col where queen is placed

    var backtrack func(row int)
    backtrack = func(row int) {
        if row == n {
            // All queens placed — record the board
            solution := make([]string, n)
            for r := 0; r < n; r++ {
                rowBytes := make([]byte, n)
                for c := 0; c < n; c++ {
                    if board[r] == c {
                        rowBytes[c] = 'Q'
                    } else {
                        rowBytes[c] = '.'
                    }
                }
                solution[r] = string(rowBytes)
            }
            result = append(result, solution)
            return
        }

        for col := 0; col < n; col++ {
            d1 := row - col + n - 1
            d2 := row + col

            // Pruning: skip if column or diagonal is under attack
            if cols[col] || diag1[d1] || diag2[d2] {
                continue
            }

            // ← CHOOSE: place queen
            board[row] = col
            cols[col] = true
            diag1[d1] = true
            diag2[d2] = true

            backtrack(row + 1)  // ← EXPLORE: place next row's queen

            // ← UNCHOOSE: remove queen
            cols[col] = false
            diag1[d1] = false
            diag2[d2] = false
        }
    }

    backtrack(0)
    return result
}
```

**How it works:** We place one queen per row. For each row, we try every column. The constraint check (column and both diagonals) prunes invalid placements before we recurse. The `cols`, `diag1`, `diag2` arrays make constraint checking O(1) instead of scanning all previously placed queens. All three constraint arrays must be set in the choose step and cleared in the unchoose step.

---

## 6. Visual Diagrams

### 6.1 DFS Traversal: Stack, Discovery, and Finish Times

**Graph (adjacency list, directed):**

```
0: [1, 2]
1: [3]
2: [3]
3: [4]
4: []
```

```
Graph structure:

    0 ──→ 1
    |      |
    ↓      ↓
    2 ──→ 3
           |
           ↓
           4
```

**DFS from node 0 (recursive):**

```
Call Stack          Action                              Discovery  Finish
(grows right)                                           Time       Time
──────────────────────────────────────────────────────────────────────────
[0]                 Visit 0 (gray), explore neighbor 1   d[0]=1
[0, 1]              Visit 1 (gray), explore neighbor 3   d[1]=2
[0, 1, 3]           Visit 3 (gray), explore neighbor 4   d[3]=3
[0, 1, 3, 4]        Visit 4 (gray), no neighbors         d[4]=4
[0, 1, 3, 4]        4 done → black                                  f[4]=5
[0, 1, 3]           3 done → black                                  f[3]=6
[0, 1]              1 done → black                                  f[1]=7
[0]                 Back to 0, explore neighbor 2
[0, 2]              Visit 2 (gray), explore neighbor 3   d[2]=8
[0, 2]              3 is BLACK → skip (not a cycle)
[0, 2]              2 done → black                                  f[2]=9
[0]                 0 done → black                                  f[0]=10

Discovery order:  0, 1, 3, 4, 2
Finish order:     4, 3, 1, 2, 0  (reverse = topological order for DAGs)

Color timeline:
     t=1   t=2   t=3   t=4   t=5   t=6   t=7   t=8   t=9   t=10
  0: GRAY  gray  gray  gray  gray  gray  gray  gray  gray  BLACK
  1:       GRAY  gray  gray  gray  gray  BLACK
  2:                                      GRAY  BLACK
  3:             GRAY  gray  gray  BLACK
  4:                   GRAY  BLACK

Note: When DFS reaches node 3 from node 2 (at t=8), node 3 is BLACK.
      Black = already fully explored. This is a CROSS EDGE, not a cycle.
      If node 3 were GRAY, that would mean 3 is an ancestor of 2 in the
      current DFS path → CYCLE.
```

---

### 6.2 Decision Tree for Subsets of [1, 2, 3]

```
                              []
                 ┌────────────┴────────────┐
             include 1                  skip 1
                [1]                        []
          ┌──────┴──────┐           ┌──────┴──────┐
      include 2      skip 2    include 2      skip 2
        [1,2]         [1]        [2]            []
       ┌──┴──┐      ┌──┴──┐   ┌──┴──┐       ┌──┴──┐
      +3    skip   +3    skip +3    skip    +3    skip
    [1,2,3] [1,2] [1,3]  [1] [2,3]  [2]   [3]    []

 Leaves (all 2^3 = 8 subsets):
 [1,2,3]  [1,2]  [1,3]  [1]  [2,3]  [2]  [3]  []

 Note: In the actual implementation, we collect results at EVERY node,
 not just leaves. The tree traversal order (DFS, pre-order) gives:
 [], [1], [1,2], [1,2,3], [1,3], [2], [2,3], [3]
```

**How the code maps to this tree:**

```
backtrack(start=0):
  record []
  i=0: choose 1 → backtrack(start=1):
    record [1]
    i=1: choose 2 → backtrack(start=2):
      record [1,2]
      i=2: choose 3 → backtrack(start=3):
        record [1,2,3]
        (no more choices)
      unchoose 3
    unchoose 2
    i=2: choose 3 → backtrack(start=3):
      record [1,3]
    unchoose 3
  unchoose 1
  i=1: choose 2 → backtrack(start=2):
    record [2]
    i=2: choose 3 → backtrack(start=3):
      record [2,3]
    unchoose 3
  unchoose 2
  i=2: choose 3 → backtrack(start=3):
    record [3]
  unchoose 3
```

---

### 6.3 Decision Tree for Permutations of [1, 2, 3]

```
                                  []
                   ┌───────────────┼───────────────┐
               pick 1           pick 2           pick 3
                 [1]              [2]              [3]
            ┌────┴────┐     ┌────┴────┐     ┌────┴────┐
         pick 2    pick 3  pick 1   pick 3  pick 1   pick 2
         [1,2]     [1,3]   [2,1]    [2,3]   [3,1]    [3,2]
           |         |       |        |        |        |
        pick 3    pick 2  pick 3   pick 1   pick 2   pick 1
        [1,2,3]  [1,3,2] [2,1,3] [2,3,1]  [3,1,2]  [3,2,1]

  Leaves (all 3! = 6 permutations):
  [1,2,3]  [1,3,2]  [2,1,3]  [2,3,1]  [3,1,2]  [3,2,1]

  Note: Results are only collected at LEAVES (when len(current) == n).
  At each level, we choose from ALL elements, skipping those already used.

  Level 0: choose from {1, 2, 3}         → 3 branches
  Level 1: choose from remaining 2       → 2 branches each
  Level 2: choose from remaining 1       → 1 branch each
  Total leaves: 3 × 2 × 1 = 6 = 3!
```

---

### 6.4 Backtracking with Pruning: SubsetsWithDuplicates([1, 2, 2])

```
Input: [1, 2, 2] (sorted)

                                []
                ┌────────────────┼──────────────────┐
             i=0: pick 1       i=1: pick 2      i=2: pick 2
               [1]                [2]            ╳ SKIP!
          ┌─────┴─────┐      ┌────┴────┐        (i=2 > start=0
       i=1: +2    i=2: +2  i=2: +2  (done)      and nums[2]==nums[1])
        [1,2]      ╳       [2,2]
         |       SKIP!        |
      i=2: +2  (i=2>start=1  (done)
      [1,2,2]   and nums[2]
         |      ==nums[1])
       (done)

  Results collected: [], [1], [1,2], [1,2,2], [2], [2,2]
  Count: 6 unique subsets (not 8, because we pruned 2 duplicate branches)

  ╳ = pruned branch (skipped due to duplicate check)

  Without pruning we'd get duplicates:
  ┌──────────────────────────────────────────────────────────┐
  │  [1,2(a)] and [1,2(b)] would both appear                │
  │  [2(a)] and [2(b)] would both appear                     │
  │  These are eliminated by the skip condition:             │
  │  if i > start && nums[i] == nums[i-1] { continue }      │
  └──────────────────────────────────────────────────────────┘
```

---

## 7. Self-Assessment

Answer these without looking at your code or notes. If you struggle with any, revisit the relevant section.

### Question 1
**What happens if you forget the "unchoose" step in backtracking?**

<details>
<summary>Answer</summary>

The state accumulates across branches. Every subsequent branch sees modifications from all previous branches. For example, in subsets generation, if you append `1` and recurse, then don't remove `1` before appending `2`, you end up with `[1, 2]` instead of just `[2]` on what should be the "exclude 1, include 2" branch. The result is that only one path through the decision tree is explored correctly (the first one), and all subsequent paths are corrupted. You'll get wrong results — typically too few results, or results that contain elements they shouldn't. Forgetting to unchoose is the single most common backtracking bug.

</details>

### Question 2
**Why does three-color marking detect cycles but simple visited marking doesn't for directed graphs?**

<details>
<summary>Answer</summary>

In a directed graph, there are multiple types of edges in the DFS tree: tree edges, back edges, forward edges, and cross edges. A cycle exists if and only if there is a **back edge** — an edge from a node to one of its ancestors in the current DFS path.

A simple `visited` boolean can't distinguish between:
- A **back edge** to an ancestor (which is a cycle): the ancestor is on the current DFS stack.
- A **cross edge** to a node in a different, already-completed branch (which is NOT a cycle): the node was fully explored on a prior branch.

Three-color marking solves this. Gray nodes are ancestors on the current DFS path. Black nodes are fully explored. An edge to a gray node is a back edge (cycle). An edge to a black node is a cross/forward edge (not a cycle).

For **undirected** graphs, simple visited + parent check works because cross edges don't exist in undirected DFS — any edge to a visited node that isn't your parent is necessarily a back edge (cycle).

</details>

### Question 3
**Why do you need to copy the `current` slice before appending it to `result` in Go backtracking?**

<details>
<summary>Answer</summary>

Go slices are backed by arrays. When you `append(result, current)`, you're storing a **reference** to the same underlying array. As the backtracking continues and modifies `current` (appending and truncating), it may modify the underlying array that previous "saved" slices are pointing to. This corrupts already-saved results.

The fix is to allocate a new slice and copy:
```go
tmp := make([]int, len(current))
copy(tmp, current)
result = append(result, tmp)
```

This creates an independent copy that won't be affected by future modifications to `current`. This is a Go-specific gotcha — in languages like Python, `result.append(current[:])` or `result.append(list(current))` serves the same purpose.

</details>

### Question 4
**Given `candidates = [2, 3, 5]` and `target = 8`, how does sorting + breaking early prune the CombinationSum search tree?**

<details>
<summary>Answer</summary>

After sorting (already sorted: `[2, 3, 5]`), at each recursive level we can break out of the loop as soon as `candidates[i] > remaining target`:

- At the top level with remaining=8: try 2, 3, 5 (all ≤ 8).
- After choosing `2, 2, 2` with remaining=2: try 2 (works!), then try 3. Since 3 > 2, **break** — no need to try 5 either.
- After choosing `3, 3` with remaining=2: try 2? No, since we only look at candidates ≥ start. Try 3. Since 3 > 2, **break**.

Without sorting and breaking, we'd continue trying candidates that are too large to fit, exploring dead-end branches. With sorting, a single `break` eliminates all remaining candidates at once (they're all larger). This can reduce the effective branching factor significantly, especially when the target is small relative to the candidate values.

</details>

### Question 5
**In the iterative DFS implementation, you mark visited on dequeue (pop). In BFS, we said marking on enqueue is critical. Why is the DFS approach different?**

<details>
<summary>Answer</summary>

In BFS, marking on dequeue allows duplicate queue entries, blowing up the queue to O(E) size and potentially giving wrong distances. In BFS, we *can* mark on enqueue because each node should only be processed once at its correct (shortest) distance.

In iterative DFS, a node can appear on the stack multiple times (pushed by different neighbors). We use the `if visited[node] { continue }` check on pop instead. This is correct because:

1. DFS doesn't care about shortest distances, so processing a node via any path is fine.
2. The stack naturally handles this — later pushes sit on top and get popped first, but the `continue` check skips duplicate pops.

You *could* also mark on push for iterative DFS (and get a slightly different traversal order), but the mark-on-pop approach is simpler and more commonly used. The key difference from BFS is that correctness of DFS doesn't depend on *when* you first visit a node, whereas BFS correctness (shortest path) does.

</details>
