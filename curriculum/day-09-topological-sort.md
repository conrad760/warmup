# Day 9: Topological Sort

> **Time budget:** 2 hours | **Prereqs:** Directed graphs, adjacency lists, BFS/DFS (Day 8)
> **Goal:** See any dependency-ordering problem, pick Kahn's or DFS topo sort in 30 seconds, and code it clean — including the hard variant (Alien Dictionary).

---

## Pattern Catalog

### 1. Kahn's Algorithm (BFS In-Degree) — Course Schedule, Task Ordering

**Trigger:** "Can you finish all courses?", "find a valid task ordering given prerequisites," any problem where items have dependencies and you need a linear ordering.

**Go Template:**
```go
func topoSortKahn(numNodes int, edges [][]int) ([]int, bool) {
    adj := make([][]int, numNodes)
    inDeg := make([]int, numNodes)
    for _, e := range edges {
        // e[0] depends on e[1]: edge from e[1] -> e[0]
        adj[e[1]] = append(adj[e[1]], e[0])
        inDeg[e[0]]++
    }

    // Seed queue with all zero in-degree nodes
    queue := []int{}
    for i := 0; i < numNodes; i++ {
        if inDeg[i] == 0 {
            queue = append(queue, i)
        }
    }

    order := []int{}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)
        for _, nb := range adj[node] {
            inDeg[nb]--
            if inDeg[nb] == 0 {
                queue = append(queue, nb)
            }
        }
    }

    if len(order) != numNodes {
        return nil, false // cycle detected — not all nodes processed
    }
    return order, true
}
```

**Complexity:** O(V + E) time, O(V + E) space.

**Watch out:**
- Edge direction matters. "Course B requires course A" means edge A → B (A must come first). Read the problem carefully to get this right.
- Cycle detection is built in: if `len(order) < numNodes`, there's a cycle. No extra work needed.
- Disconnected nodes with no edges have in-degree 0 and get seeded into the queue automatically. They appear in the output. Don't forget them.

---

### 2. DFS-Based Topo Sort (Post-Order Reversal) — Alternative Approach

**Trigger:** Same as Pattern 1. Use this when you're more comfortable with DFS or when the problem naturally lends itself to recursive exploration (e.g., you need to process dependencies before the dependent node).

**Go Template:**
```go
func topoSortDFS(numNodes int, edges [][]int) ([]int, bool) {
    adj := make([][]int, numNodes)
    for _, e := range edges {
        adj[e[1]] = append(adj[e[1]], e[0])
    }

    // 0 = unvisited, 1 = in current path (gray), 2 = done (black)
    color := make([]int, numNodes)
    order := []int{}
    hasCycle := false

    var dfs func(node int)
    dfs = func(node int) {
        if hasCycle {
            return
        }
        color[node] = 1 // mark gray — currently on recursion stack
        for _, nb := range adj[node] {
            if color[nb] == 1 {
                hasCycle = true // back edge → cycle
                return
            }
            if color[nb] == 0 {
                dfs(nb)
            }
        }
        color[node] = 2 // mark black — done
        order = append(order, node) // post-order
    }

    for i := 0; i < numNodes; i++ {
        if color[i] == 0 {
            dfs(i)
        }
    }

    if hasCycle {
        return nil, false
    }

    // Reverse post-order = topological order
    for l, r := 0, len(order)-1; l < r; l, r = l+1, r-1 {
        order[l], order[r] = order[r], order[l]
    }
    return order, true
}
```

**Complexity:** O(V + E) time, O(V + E) space (adjacency list + recursion stack).

**Watch out:**
- You **must** reverse the post-order. A common bug is to forget the reversal and return nodes in the wrong order.
- Three-color marking is essential for cycle detection. Two colors (visited/not visited) can't distinguish a back edge (cycle) from a cross edge (already-finished node).
- Don't forget to loop over all nodes (`for i := 0; i < numNodes`). Disconnected components won't be visited otherwise.

---

### 3. Cycle Detection in Directed Graph — Three-Color DFS or Kahn's Leftover Check

**Trigger:** "Is it possible to finish all courses?", "detect a circular dependency," any problem that asks whether a valid ordering exists.

**Go Template (Kahn's approach — simplest):**
```go
func canFinish(numCourses int, prerequisites [][]int) bool {
    adj := make([][]int, numCourses)
    inDeg := make([]int, numCourses)
    for _, p := range prerequisites {
        adj[p[1]] = append(adj[p[1]], p[0])
        inDeg[p[0]]++
    }

    queue := []int{}
    for i := 0; i < numCourses; i++ {
        if inDeg[i] == 0 {
            queue = append(queue, i)
        }
    }

    processed := 0
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        processed++
        for _, nb := range adj[node] {
            inDeg[nb]--
            if inDeg[nb] == 0 {
                queue = append(queue, nb)
            }
        }
    }

    return processed == numCourses
}
```

**Complexity:** O(V + E) time, O(V + E) space.

**Watch out:**
- This is identical to Kahn's algorithm but you only need the boolean — you don't need the actual order. Still, return the order if the problem asks for it (Course Schedule II).
- The DFS three-color approach (Pattern 2) also detects cycles. Pick whichever you're faster at coding.
- Self-loops (`edge [a, a]`) are immediate cycles. Kahn's handles this: the node will have in-degree ≥ 1 from itself and will never enter the queue.

---

### 4. Build Order from Constraints — Alien Dictionary (Build Graph from Pairwise Comparisons)

**Trigger:** "Given a sorted list, determine the underlying ordering rules," "derive an alphabet from sorted words," any problem where the graph isn't given directly and you must construct it from sequential comparisons.

**Go Template (Alien Dictionary):**
```go
func alienOrder(words []string) string {
    // Step 1: Initialize adjacency list and in-degree for ALL characters
    adj := map[byte][]byte{}
    inDeg := map[byte]int{}
    for _, w := range words {
        for i := 0; i < len(w); i++ {
            if _, ok := inDeg[w[i]]; !ok {
                inDeg[w[i]] = 0
                adj[w[i]] = []byte{}
            }
        }
    }

    // Step 2: Build graph by comparing adjacent word pairs
    for i := 0; i < len(words)-1; i++ {
        w1, w2 := words[i], words[i+1]
        // TRAP: if w1 is longer and w2 is a prefix of w1, input is invalid
        if len(w1) > len(w2) && strings.HasPrefix(w1, w2) {
            return ""
        }
        for j := 0; j < min(len(w1), len(w2)); j++ {
            if w1[j] != w2[j] {
                adj[w1[j]] = append(adj[w1[j]], w2[j])
                inDeg[w2[j]]++
                break // only the FIRST difference matters
            }
        }
    }

    // Step 3: Kahn's algorithm on the constructed graph
    queue := []byte{}
    for ch, deg := range inDeg {
        if deg == 0 {
            queue = append(queue, ch)
        }
    }

    result := []byte{}
    for len(queue) > 0 {
        ch := queue[0]
        queue = queue[1:]
        result = append(result, ch)
        for _, nb := range adj[ch] {
            inDeg[nb]--
            if inDeg[nb] == 0 {
                queue = append(queue, nb)
            }
        }
    }

    if len(result) != len(inDeg) {
        return "" // cycle
    }
    return string(result)
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**Complexity:** O(C) where C = total number of characters across all words (for graph construction + topo sort).

**Watch out:**
- Only compare **adjacent** word pairs. Comparing all pairs produces wrong edges.
- Only the **first** differing character in each pair gives you an edge. After the first mismatch, remaining characters tell you nothing.
- The prefix trap: `["abc", "ab"]` is invalid — a longer word cannot come before its own prefix. You must detect this and return `""`.
- Every character that appears in any word must be in the graph, even if it participates in no edges. Initialize in-degree for all characters before building edges.
- Duplicate edges: if the same edge `(a → b)` appears from multiple word pairs, in-degree gets incremented multiple times. Either deduplicate edges with a set or accept that Kahn's still works correctly (it does — the extra in-degree will be decremented the right number of times). But deduplication is cleaner.

---

## Decision Framework

Read the problem statement. Match the **first** rule that fits:

| Signal in problem | Technique |
|---|---|
| "Prerequisites" or "dependencies" between tasks | Kahn's topo sort (Pattern 1) |
| "Is it possible to finish all ___?" | Cycle detection via topo sort (Pattern 3) |
| "Find a valid ordering" given constraints | Kahn's or DFS topo sort (Pattern 1 or 2) |
| "Determine order from sorted list / pairwise comparisons" | Build graph + topo sort (Pattern 4) |
| Need **lexicographically smallest** valid ordering | Kahn's with min-heap instead of queue |
| "Parallel scheduling" / "minimum number of semesters" | Kahn's (BFS level count = critical path length) |

**Kahn's vs DFS — when to pick which:**
- Default to **Kahn's**. It's iterative (no stack overflow risk), cycle detection is built in (check count), and it's easier to modify (swap queue for min-heap for lex order).
- Use **DFS** when you're more fluent with it or when the problem is naturally recursive (e.g., evaluating expressions with dependencies).

---

## Common Interview Traps

### 1. Not handling disconnected nodes
Nodes with no incoming or outgoing edges still must appear in the output. Kahn's handles this automatically (they start with in-degree 0). DFS handles it if you loop over **all** nodes.

### 2. Alien Dictionary: comparing non-adjacent words
Only compare `words[i]` with `words[i+1]`. Comparing `words[0]` with `words[2]` can produce incorrect edges that contradict the actual ordering.

### 3. Alien Dictionary: the prefix trap
`["apple", "app"]` means `"apple"` comes before `"app"` in the alien order. But in **any** lexicographic order, a prefix must come before a longer word with the same prefix. So this input is invalid — return `""`. Many candidates miss this.

### 4. Self-loops = immediate cycle
If a constraint says `a` depends on `a`, that's a cycle. Kahn's handles it (in-degree never reaches 0). DFS handles it (gray node visited again). But make sure your graph construction doesn't silently create self-loops from bad input parsing.

### 5. Multiple valid orderings
Most problems accept any valid topological order. But if the problem asks for **lexicographically smallest**, replace the queue with a min-heap:
```go
// Replace queue with min-heap for lexicographic order
h := &IntHeap{}
heap.Init(h)
for i := 0; i < n; i++ {
    if inDeg[i] == 0 {
        heap.Push(h, i)
    }
}
for h.Len() > 0 {
    node := heap.Pop(h).(int)
    // ... same as before, but use heap.Push instead of append
}
```

### 6. Getting edge direction backwards
"Course 1 requires course 0" means edge `0 → 1` (0 must come first), NOT `1 → 0`. This is the most silent, frustrating bug. Always write out one example edge on your whiteboard before coding the graph construction.

---

## Thought Process Walkthrough

### Problem 1: Course Schedule II (LC 210) — Standard Topo Sort

**Problem:** There are `numCourses` courses labeled `0` to `numCourses - 1`. You are given `prerequisites` where `prerequisites[i] = [ai, bi]` means you must take course `bi` before course `ai`. Return a valid ordering to finish all courses. If impossible, return an empty array.

**Interview simulation — what to say out loud:**

**Step 1: Classify (15 seconds)**
> "Prerequisites, valid ordering — this is textbook topological sort. I'll use Kahn's algorithm."

**Step 2: Approach (30 seconds)**
> "Build an adjacency list and in-degree array. For each prerequisite `[a, b]`, add edge `b → a` and increment `inDeg[a]`. Seed a queue with all zero in-degree nodes. Process BFS-style: dequeue a node, add it to the result, decrement in-degree of its neighbors, enqueue any neighbor that hits zero. If the result has fewer than `numCourses` nodes, there's a cycle — return empty."

**Step 3: Edge cases (15 seconds)**
> "No prerequisites — every node has in-degree 0, all get queued, any permutation works. Single course — return `[0]`. Cycle — return `[]`."

**Step 4: Code**
```go
func findOrder(numCourses int, prerequisites [][]int) []int {
    adj := make([][]int, numCourses)
    inDeg := make([]int, numCourses)
    for _, p := range prerequisites {
        adj[p[1]] = append(adj[p[1]], p[0])
        inDeg[p[0]]++
    }

    queue := []int{}
    for i := 0; i < numCourses; i++ {
        if inDeg[i] == 0 {
            queue = append(queue, i)
        }
    }

    order := []int{}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)
        for _, nb := range adj[node] {
            inDeg[nb]--
            if inDeg[nb] == 0 {
                queue = append(queue, nb)
            }
        }
    }

    if len(order) != numCourses {
        return []int{}
    }
    return order
}
```

**Step 5: Complexity analysis (say this)**
> "Time is O(V + E) — we visit every node once and process every edge once. Space is O(V + E) for the adjacency list. V is numCourses, E is len(prerequisites)."

**Step 6: Test with example**
> Input: `numCourses = 4, prerequisites = [[1,0],[2,0],[3,1],[3,2]]`
>
> Adjacency: `0→[1,2], 1→[3], 2→[3]`
> In-degree: `[0, 1, 1, 2]`
> Queue starts: `[0]`
> Process 0 → order: `[0]`, decrement 1→0, 2→0 → queue: `[1, 2]`
> Process 1 → order: `[0, 1]`, decrement 3→1 → queue: `[2]`
> Process 2 → order: `[0, 1, 2]`, decrement 3→0 → queue: `[3]`
> Process 3 → order: `[0, 1, 2, 3]`
> len(order) == 4 == numCourses ✓

---

### Problem 2: Alien Dictionary (LC 269) — Graph Construction + Topo Sort

**Problem:** Given a list of words sorted in alien dictionary order, derive the order of characters in the alien alphabet. Return the characters in order. If no valid order exists, return `""`.

**Interview simulation — what to say out loud:**

**Step 1: Classify (20 seconds)**
> "I'm given a sorted list and need to figure out the underlying ordering. This is the 'build order from constraints' pattern. I need to construct a graph from the sorted words, then run topological sort."

**Step 2: Approach (60 seconds)**
> "Two phases. Phase 1: Build the graph. Compare each adjacent pair of words. Find the first position where they differ — that gives me an edge: `word1[j] → word2[j]` (word1's char comes before word2's char). I stop after the first difference because subsequent positions don't give valid information."
>
> "I also need to handle the prefix edge case: if word1 is longer than word2 and word2 is a prefix of word1, the input is invalid."
>
> "Phase 2: Run Kahn's topo sort on the graph. If the result doesn't include all characters, there's a cycle."
>
> "Important: I need to track **all** unique characters, not just the ones that appear in edges. Characters with no ordering constraints still belong in the output."

**Step 3: Edge cases (20 seconds)**
> "Single word — return all its unique characters in any order. Two identical words — no edges, return characters in any order. The prefix trap: `['abc', 'ab']` is invalid, return `''`. Cycle in constraints — return `''`."

**Step 4: Code**
```go
func alienOrder(words []string) string {
    adj := map[byte]map[byte]bool{}
    inDeg := map[byte]int{}

    // Initialize all characters
    for _, w := range words {
        for i := 0; i < len(w); i++ {
            if adj[w[i]] == nil {
                adj[w[i]] = map[byte]bool{}
            }
            if _, ok := inDeg[w[i]]; !ok {
                inDeg[w[i]] = 0
            }
        }
    }

    // Build edges from adjacent word pairs
    for i := 0; i < len(words)-1; i++ {
        w1, w2 := words[i], words[i+1]
        minLen := len(w1)
        if len(w2) < minLen {
            minLen = len(w2)
        }
        // Prefix trap check
        if len(w1) > len(w2) {
            isPrefix := true
            for j := 0; j < len(w2); j++ {
                if w1[j] != w2[j] {
                    isPrefix = false
                    break
                }
            }
            if isPrefix {
                return ""
            }
        }
        for j := 0; j < minLen; j++ {
            if w1[j] != w2[j] {
                if !adj[w1[j]][w2[j]] { // avoid duplicate edges
                    adj[w1[j]][w2[j]] = true
                    inDeg[w2[j]]++
                }
                break
            }
        }
    }

    // Kahn's
    queue := []byte{}
    for ch, deg := range inDeg {
        if deg == 0 {
            queue = append(queue, ch)
        }
    }

    result := []byte{}
    for len(queue) > 0 {
        ch := queue[0]
        queue = queue[1:]
        result = append(result, ch)
        for nb := range adj[ch] {
            inDeg[nb]--
            if inDeg[nb] == 0 {
                queue = append(queue, nb)
            }
        }
    }

    if len(result) != len(inDeg) {
        return ""
    }
    return string(result)
}
```

**Step 5: Complexity analysis (say this)**
> "Let C = total characters across all words, U = unique characters. Graph construction is O(C) — we scan through each word pair comparing characters. Kahn's is O(U + E) where E is the number of edges. Total: O(C). Space is O(U + E) for the graph."

**Step 6: Test with example**
> Input: `["wrt", "wrf", "er", "ett", "rftt"]`
>
> Compare "wrt" vs "wrf": first diff at index 2 → `t → f`
> Compare "wrf" vs "er": first diff at index 0 → `w → e`
> Compare "er" vs "ett": first diff at index 1 → `r → t`
> Compare "ett" vs "rftt": first diff at index 0 → `e → r`
>
> Edges: `t→f, w→e, r→t, e→r`
> All chars: `{w, r, t, f, e}`
> In-degree: `w:0, r:1, t:1, f:1, e:1`
> Queue: `[w]`
> Process w → e's inDeg→0 → queue: `[e]`, result: `[w]`
> Process e → r's inDeg→0 → queue: `[r]`, result: `[w,e]`
> Process r → t's inDeg→0 → queue: `[t]`, result: `[w,e,r]`
> Process t → f's inDeg→0 → queue: `[f]`, result: `[w,e,r,t]`
> Process f → result: `[w,e,r,t,f]`
> len(result) == 5 == len(inDeg) ✓ → return `"wertf"`

**Step 7: Follow-up readiness**
- "If asked for lexicographic smallest: use a min-heap instead of a queue in Kahn's."
- "If asked to validate whether a given dictionary is consistent: build the graph and check for cycles."

---

## Time Targets

| Problem | Target | Notes |
|---|---|---|
| Course Schedule (LC 207) | 8 min | Pure cycle detection — Kahn's or DFS. |
| Course Schedule II (LC 210) | 10 min | Same as above but return the order. |
| Alien Dictionary (LC 269) | 18 min | Graph construction is the hard part, not the topo sort. |
| Parallel Courses (LC 1136) | 12 min | Kahn's + count BFS levels for critical path. |
| Minimum Height Trees (LC 310) | 15 min | Leaf-stripping variant of topo sort. |

If you're over these times, isolate the bottleneck: is it graph construction or the topo sort itself? Drill that part separately.

---

## Quick Drill

Five problems to solve **in order** during your practice session. Don't move on until the current one compiles and passes.

| # | Problem | Pattern | Target |
|---|---|---|---|
| 1 | [Course Schedule](https://leetcode.com/problems/course-schedule/) (LC 207) | Cycle detection (Kahn's) | 8 min |
| 2 | [Course Schedule II](https://leetcode.com/problems/course-schedule-ii/) (LC 210) | Kahn's topo sort | 10 min |
| 3 | [Alien Dictionary](https://leetcode.com/problems/alien-dictionary/) (LC 269) | Build graph + topo sort | 18 min |
| 4 | [Minimum Height Trees](https://leetcode.com/problems/minimum-height-trees/) (LC 310) | Leaf-stripping (topo sort variant) | 15 min |
| 5 | [Parallel Courses](https://leetcode.com/problems/parallel-courses/) (LC 1136) | Kahn's + level count | 12 min |

**Total coding time:** ~63 min. Use remaining time for review and self-assessment.

---

## Self-Assessment

Answer these without looking at your notes. If you miss any, re-read the relevant section.

### 1. What are the two standard approaches for topological sort, and how does each detect cycles?
**Expected answer:** (1) Kahn's algorithm (BFS + in-degree): if the final order contains fewer than V nodes, there's a cycle — some nodes never reach in-degree 0. (2) DFS post-order reversal: use three-color marking (white/gray/black). If you visit a gray node during DFS, you've found a back edge, which means a cycle.

### 2. In Kahn's algorithm, what goes into the initial queue and why?
**Expected answer:** All nodes with in-degree 0. These are nodes with no prerequisites — they can be processed first. If the graph has no such nodes and there are nodes to process, the entire graph is one or more cycles.

### 3. In the Alien Dictionary problem, why do you only compare adjacent words?
**Expected answer:** The list is sorted. Adjacent words give you the tightest constraints. Comparing non-adjacent words (e.g., word 0 vs word 2) can produce edges that are redundant at best or incorrect at worst, because the ordering between non-adjacent words is a transitive consequence of adjacent comparisons, not a direct one.

### 4. What is the "prefix trap" in Alien Dictionary and how do you detect it?
**Expected answer:** If `words[i]` is longer than `words[i+1]` and `words[i+1]` is a prefix of `words[i]`, the input is invalid. In any lexicographic ordering, a shorter prefix must come before the longer word. This situation means the given sorted order is self-contradictory. Detect it by checking `len(w1) > len(w2)` and verifying all characters of `w2` match the start of `w1`.

### 5. When would you use a min-heap instead of a queue in Kahn's algorithm?
**Expected answer:** When the problem asks for the **lexicographically smallest** valid topological ordering. A regular queue gives any valid order. A min-heap always picks the smallest available node next, guaranteeing lexicographic minimality. The time complexity changes from O(V + E) to O(V log V + E).

### 6. Why do you reverse the order in DFS-based topo sort?
**Expected answer:** DFS appends a node to the result in **post-order** — after all its descendants are processed. This means dependencies appear **after** the nodes that depend on them. Reversing the post-order gives the correct topological order where dependencies come first.
