# Day 8: Graphs — BFS & DFS

> **Time budget:** 2 hours | **Prereqs:** BFS/DFS basics, queue/stack mechanics, adjacency lists
> **Goal:** See a problem, classify it in 30 seconds, pick the right traversal, and code it without hesitation.

---

## Pattern Catalog

### 1. Grid DFS — Islands, Flood Fill, Connected Regions

**Trigger:** "Count the number of islands," "fill a region," "find connected components in a 2D grid."

**Go Template:**
```go
func dfs(grid [][]byte, r, c int) {
    if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) {
        return
    }
    if grid[r][c] != '1' {
        return
    }
    grid[r][c] = '0' // mark visited by mutating grid
    dfs(grid, r+1, c)
    dfs(grid, r-1, c)
    dfs(grid, r, c+1)
    dfs(grid, r, c-1)
}

// In main loop: for each cell == '1', call dfs(), increment count.
```

**Complexity:** O(R * C) time, O(R * C) stack space worst case.

**Watch out:**
- Mutating the input grid avoids a separate `visited` set, but mention this trade-off to your interviewer. Offer to restore the grid afterward if asked.
- Stack overflow on very large grids (e.g., 1000x1000). Mention you'd switch to iterative DFS or BFS if that's a concern.

---

### 2. Grid BFS — Shortest Path in Grid, Minimum Steps

**Trigger:** "Shortest path in a grid," "minimum number of moves," any grid problem asking for a distance.

**Go Template:**
```go
type point struct{ r, c int }

func bfs(grid [][]int, sr, sc int) int {
    rows, cols := len(grid), len(grid[0])
    dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
    queue := []point{{sr, sc}}
    grid[sr][sc] = -1 // mark visited ON ENQUEUE
    steps := 0

    for len(queue) > 0 {
        size := len(queue) // level-by-level
        for i := 0; i < size; i++ {
            p := queue[0]
            queue = queue[1:]
            if isTarget(p) {
                return steps
            }
            for _, d := range dirs {
                nr, nc := p.r+d[0], p.c+d[1]
                if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] != -1 {
                    grid[nr][nc] = -1 // mark visited HERE, not on dequeue
                    queue = append(queue, point{nr, nc})
                }
            }
        }
        steps++
    }
    return -1
}
```

**Complexity:** O(R * C) time, O(R * C) space for the queue.

**Watch out:**
- **Mark visited on ENQUEUE, not dequeue.** This is the #1 BFS bug. If you mark on dequeue, the same cell gets added to the queue multiple times and your runtime blows up.
- Use level-by-level processing (`size := len(queue)`) to track distance. Without it you have no way to count steps.

---

### 3. Graph DFS — Clone Graph, Connected Components, Path Existence

**Trigger:** "Clone/deep copy a graph," "are nodes A and B connected?", "find all connected components," "traverse all reachable nodes."

**Go Template (Clone Graph):**
```go
// Node is: type Node struct { Val int; Neighbors []*Node }
func cloneGraph(node *Node) *Node {
    if node == nil {
        return nil
    }
    cloned := make(map[*Node]*Node) // old -> new
    var dfs func(n *Node) *Node
    dfs = func(n *Node) *Node {
        if c, ok := cloned[n]; ok {
            return c // already cloned — stops infinite loops
        }
        copy := &Node{Val: n.Val}
        cloned[n] = copy // register BEFORE recursing
        for _, nb := range n.Neighbors {
            copy.Neighbors = append(copy.Neighbors, dfs(nb))
        }
        return copy
    }
    return dfs(node)
}
```

**Complexity:** O(V + E) time, O(V) space for the visited/cloned map.

**Watch out:**
- You **must** insert into the cloned map before recursing into neighbors. Otherwise cycles cause infinite recursion.
- For connected components, remember to loop over **all** nodes — not just start from node 0. Disconnected components are a classic miss.

---

### 4. Graph BFS — Shortest Path Unweighted, Word Ladder

**Trigger:** "Shortest path between two nodes," "fewest edges," any unweighted graph distance question.

**Go Template:**
```go
func shortestPath(adj map[int][]int, src, dst int) int {
    visited := map[int]bool{src: true}
    queue := []int{src}
    steps := 0

    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]
            if node == dst {
                return steps
            }
            for _, nb := range adj[node] {
                if !visited[nb] {
                    visited[nb] = true // mark on enqueue
                    queue = append(queue, nb)
                }
            }
        }
        steps++
    }
    return -1
}
```

**Complexity:** O(V + E) time, O(V) space.

**Watch out:**
- BFS guarantees shortest path **only for unweighted graphs.** If there are weights, you need Dijkstra. State this explicitly in interviews.
- For undirected graphs, make sure you add edges in **both** directions when building the adjacency list.

---

### 5. Multi-Source BFS — All Sources Simultaneously

**Trigger:** "Rotten oranges spread to adjacent fresh ones," "walls and gates — fill distance from nearest gate," any problem with **multiple starting points** spreading outward.

**Go Template (Rotten Oranges):**
```go
func orangesRotting(grid [][]int) int {
    rows, cols := len(grid), len(grid[0])
    dirs := [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
    type point struct{ r, c int }
    var queue []point
    fresh := 0

    // Seed queue with ALL sources
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == 2 {
                queue = append(queue, point{r, c})
            } else if grid[r][c] == 1 {
                fresh++
            }
        }
    }

    minutes := 0
    for len(queue) > 0 && fresh > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            p := queue[0]
            queue = queue[1:]
            for _, d := range dirs {
                nr, nc := p.r+d[0], p.c+d[1]
                if nr >= 0 && nr < rows && nc >= 0 && nc < cols && grid[nr][nc] == 1 {
                    grid[nr][nc] = 2 // mark on enqueue
                    fresh--
                    queue = append(queue, point{nr, nc})
                }
            }
        }
        minutes++
    }
    if fresh > 0 {
        return -1
    }
    return minutes
}
```

**Complexity:** O(R * C) time and space.

**Watch out:**
- The key insight: enqueue **all** sources before processing anything. This gives you the "spreading from all points simultaneously" behavior for free.
- Don't forget the edge case: if `fresh == 0` at the start, answer is `0` (nothing to rot). The template handles this because the loop body never executes.

---

### 6. Implicit Graph BFS — Generate Neighbors On-the-Fly

**Trigger:** "Word ladder — transform one word to another," "open the lock — reach target combination," any problem where the **graph isn't given** and you generate neighbors from a set of rules.

**Go Template (Word Ladder):**
```go
func ladderLength(beginWord string, endWord string, wordList []string) int {
    wordSet := make(map[string]bool, len(wordList))
    for _, w := range wordList {
        wordSet[w] = true
    }
    if !wordSet[endWord] {
        return 0
    }

    visited := map[string]bool{beginWord: true}
    queue := []string{beginWord}
    steps := 1

    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            word := queue[0]
            queue = queue[1:]
            if word == endWord {
                return steps
            }
            // Generate all neighbors: change each char to a-z
            buf := []byte(word)
            for j := 0; j < len(buf); j++ {
                orig := buf[j]
                for ch := byte('a'); ch <= byte('z'); ch++ {
                    buf[j] = ch
                    next := string(buf)
                    if wordSet[next] && !visited[next] {
                        visited[next] = true
                        queue = append(queue, next)
                    }
                }
                buf[j] = orig
            }
        }
        steps++
    }
    return 0
}
```

**Complexity:** O(M^2 * N) where M = word length, N = word list size. Each word generates 26*M neighbors, and creating each string costs O(M).

**Watch out:**
- Generating neighbors is the hard part. For word ladder it's 26 * word_length neighbors per node. Know this cost.
- Delete from `wordSet` instead of using a separate `visited` set — same effect, saves memory. (Both approaches are fine in interviews; pick whichever you explain more clearly.)
- **Bidirectional BFS** is the advanced optimization. Mention it if you finish early but don't code it unless asked.

---

## Decision Framework

Read the problem statement. Match the **first** rule that fits:

| Signal in problem | Technique |
|---|---|
| Grid + "count regions/islands/components" | DFS flood fill (Pattern 1) |
| Grid + "shortest path / minimum steps" | BFS on grid (Pattern 2) |
| Grid + "spread from multiple sources" | Multi-source BFS (Pattern 5) |
| "Clone" or "deep copy" a graph | DFS with old-to-new hash map (Pattern 3) |
| "Shortest path" in unweighted graph | BFS (Pattern 4) |
| "Minimum transformations / steps" + no explicit graph | BFS on implicit graph (Pattern 6) |
| "Does path exist between A and B" | DFS or BFS — either works (Pattern 3) |
| Multiple starting points / simultaneous spread | Multi-source BFS (Pattern 5) |

**The one rule you must never break:**
Shortest path + unweighted = BFS. Always. No exceptions. If you reach for DFS here, stop yourself.

---

## Common Interview Traps

### 1. Marking visited on DEQUEUE instead of ENQUEUE
```
BAD:  node := queue[0]; queue = queue[1:]; visited[node] = true
GOOD: visited[nb] = true; queue = append(queue, nb)
```
If you mark on dequeue, multiple copies of the same node pile up in the queue. Time complexity goes from O(V+E) to O(V^2) or worse. This is the single most common BFS bug.

### 2. Grid: forgetting to check bounds before accessing `grid[r][c]`
Always check `r >= 0 && r < rows && c >= 0 && c < cols` **before** reading the cell. Out-of-bounds panics are instant interview failures in Go.

### 3. Graph: not handling disconnected components
```go
// WRONG: only processes component containing node 0
dfs(0)

// RIGHT: process every node
for i := 0; i < n; i++ {
    if !visited[i] {
        dfs(i)
        count++
    }
}
```

### 4. Clone graph: infinite loop from cycles
You must check if a node is already in your cloned map **before** recursing. Insert the clone into the map **before** processing neighbors. This breaks the cycle.

### 5. BFS: not tracking distance
If you just pop and push without level-by-level processing, you have no way to know how many steps you've taken. Use the `size := len(queue)` pattern to process one level at a time.

### 6. Directed vs undirected edges
When building an adjacency list for an **undirected** graph:
```go
adj[u] = append(adj[u], v)
adj[v] = append(adj[v], u) // don't forget this
```
Forgetting the reverse edge is silent — your code runs but gives wrong answers.

---

## Thought Process Walkthrough

### Problem 1: Number of Islands (Grid DFS)

**Problem:** Given a 2D grid of `'1'`s (land) and `'0'`s (water), count the number of islands. An island is surrounded by water and is formed by connecting adjacent lands horizontally or vertically.

**Interview simulation — what to say out loud:**

**Step 1: Classify (15 seconds)**
> "This is a grid problem asking me to count connected regions. That's DFS flood fill — Pattern 1 in my head."

**Step 2: Approach (30 seconds)**
> "I'll scan every cell. When I find a '1', I increment my island count and DFS to sink the entire island by marking cells as '0'. This way I never count the same island twice. I don't need a separate visited set because I'm modifying the grid in place."

**Step 3: Edge cases (20 seconds)**
> "Empty grid — return 0. Single cell grid. All water. All land (one giant island). I'll handle the empty grid check up front."

**Step 4: Code**
```go
func numIslands(grid [][]byte) int {
    if len(grid) == 0 {
        return 0
    }
    rows, cols := len(grid), len(grid[0])
    count := 0

    var dfs func(r, c int)
    dfs = func(r, c int) {
        if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] != '1' {
            return
        }
        grid[r][c] = '0'
        dfs(r+1, c)
        dfs(r-1, c)
        dfs(r, c+1)
        dfs(r, c-1)
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == '1' {
                count++
                dfs(r, c)
            }
        }
    }
    return count
}
```

**Step 5: Complexity analysis (say this)**
> "Time is O(R * C) — each cell is visited at most once. Space is O(R * C) in the worst case for the recursion stack, if the entire grid is land and we spiral through it. If the interviewer is concerned about stack depth, I can convert to iterative DFS with an explicit stack or use BFS."

**Step 6: Test with example**
> Walk through a small 3x3 grid, tracing which cells get sunk.

---

### Problem 2: Word Ladder (Implicit Graph BFS)

**Problem:** Given `beginWord`, `endWord`, and a dictionary `wordList`, find the length of the shortest transformation sequence from `beginWord` to `endWord`, such that only one letter can be changed at a time and each transformed word must exist in the word list.

**Interview simulation — what to say out loud:**

**Step 1: Classify (15 seconds)**
> "Shortest transformation sequence — that's 'minimum steps.' There's no explicit graph, so I'm building one implicitly. This is BFS on an implicit graph — Pattern 6."

**Step 2: Approach (60 seconds)**
> "Each word is a node. Two words are connected if they differ by exactly one character. I need the shortest path from beginWord to endWord."
>
> "For neighbor generation, I'll try changing each position to every letter a–z. That's 26 * word_length candidates per node. I check each candidate against a hash set of the word list."
>
> "I'll use standard BFS with level-by-level processing to track the number of steps."

**Step 3: Edge cases (20 seconds)**
> "endWord not in wordList — return 0 immediately. beginWord equals endWord — the problem guarantees they differ, but I'd handle it. No valid path exists — return 0 after BFS exhausts."

**Step 4: Code**
```go
func ladderLength(beginWord string, endWord string, wordList []string) int {
    wordSet := make(map[string]bool, len(wordList))
    for _, w := range wordList {
        wordSet[w] = true
    }
    if !wordSet[endWord] {
        return 0
    }

    visited := map[string]bool{beginWord: true}
    queue := []string{beginWord}
    steps := 1

    for len(queue) > 0 {
        size := len(queue)
        for i := 0; i < size; i++ {
            word := queue[0]
            queue = queue[1:]
            if word == endWord {
                return steps
            }
            buf := []byte(word)
            for j := 0; j < len(buf); j++ {
                orig := buf[j]
                for ch := byte('a'); ch <= byte('z'); ch++ {
                    if ch == orig {
                        continue
                    }
                    buf[j] = ch
                    next := string(buf)
                    if wordSet[next] && !visited[next] {
                        visited[next] = true
                        queue = append(queue, next)
                    }
                }
                buf[j] = orig
            }
        }
        steps++
    }
    return 0
}
```

**Step 5: Complexity analysis (say this)**
> "Let M = word length, N = number of words in the list. For each word we dequeue, we generate 26 * M neighbors, and each string comparison/creation is O(M). So total work is O(N * M * 26 * M) = O(N * M^2). Space is O(N * M) for the visited set and queue."

**Step 6: Follow-up readiness**
- "If the interviewer asks to optimize: bidirectional BFS — start BFS from both ends and meet in the middle. Reduces the branching factor."
- "If asked for all shortest paths (Word Ladder II): BFS to find distances, then DFS/backtrack to reconstruct all paths."

---

## Time Targets

| Problem | Target | Notes |
|---|---|---|
| Number of Islands | 8 min | Pure template. Should be automatic. |
| Flood Fill | 6 min | Same as islands — even simpler. |
| Clone Graph | 10 min | DFS + hash map. Watch for cycle handling. |
| Rotting Oranges | 10 min | Multi-source BFS. Count fresh, seed all rotten. |
| Word Ladder | 15 min | Implicit graph. Neighbor generation is the core. |

If you're over these times, drill the specific pattern template until it's muscle memory.

---

## Quick Drill

Five problems to solve **in order** during your practice session. Don't move on until the current one compiles and passes.

| # | Problem | Pattern | Target |
|---|---|---|---|
| 1 | [Number of Islands](https://leetcode.com/problems/number-of-islands/) (LC 200) | Grid DFS | 8 min |
| 2 | [Flood Fill](https://leetcode.com/problems/flood-fill/) (LC 733) | Grid DFS | 6 min |
| 3 | [Rotting Oranges](https://leetcode.com/problems/rotting-oranges/) (LC 994) | Multi-source BFS | 10 min |
| 4 | [Clone Graph](https://leetcode.com/problems/clone-graph/) (LC 133) | Graph DFS + hash map | 10 min |
| 5 | [Word Ladder](https://leetcode.com/problems/word-ladder/) (LC 127) | Implicit graph BFS | 15 min |

**Total coding time:** ~49 min. Use remaining time for review and self-assessment.

---

## Self-Assessment

Answer these without looking at your notes. If you miss any, re-read the relevant section.

### 1. When do you use BFS vs DFS on a grid?
**Expected answer:** BFS when you need shortest path / minimum steps. DFS when you need to explore/count/mark connected regions. BFS guarantees shortest path on unweighted graphs; DFS does not.

### 2. Where exactly do you mark a node as visited in BFS, and why?
**Expected answer:** On **enqueue**, not on dequeue. If you mark on dequeue, the same node can be added to the queue multiple times by different neighbors, wasting time and potentially causing incorrect results.

### 3. How does multi-source BFS differ from regular BFS?
**Expected answer:** You seed the queue with **all** source nodes before starting the loop. Each level of BFS then expands from all sources simultaneously. This gives you the minimum distance from the nearest source to every reachable cell.

### 4. In clone graph, what causes an infinite loop and how do you prevent it?
**Expected answer:** Cycles in the graph cause infinite recursion if you don't track already-cloned nodes. Prevent it by inserting the new node into a `cloned` map **before** recursing into its neighbors. On re-entry, check the map and return the existing clone.

### 5. What makes a graph "implicit" and how does that change your BFS code?
**Expected answer:** An implicit graph is one where nodes and edges aren't given directly — you generate neighbors by applying transformation rules (e.g., change one letter, turn one dial). The BFS structure is identical, but instead of iterating over an adjacency list, you have a neighbor-generation function. The main cost is generating and validating those neighbors.
