# Day 14 — Topological Sort & Cycle Detection

---

## 1. Curated Learning Resources

| # | Resource | Format | Why This One |
|---|----------|--------|-------------|
| 1 | [Topological Sort — WilliamFiset](https://www.youtube.com/watch?v=eL-KzMXSXXI) | Video (18 min) | Animated walkthrough of both Kahn's and DFS-based approaches with clear visuals on a sample DAG. |
| 2 | [Kahn's Algorithm Step-by-Step — Abdul Bari](https://www.youtube.com/watch?v=tggiFvaxjrY) | Video (12 min) | Hand-traces in-degree tracking and queue state at each step. Best single resource for internalizing Kahn's. |
| 3 | [DFS-Based Topological Sort — Back To Back SWE](https://www.youtube.com/watch?v=AfSk24UTFS8) | Video (15 min) | Focuses on why reversed post-order works and how three-color marking detects cycles. |
| 4 | [Topological Sort — CP-Algorithms](https://cp-algorithms.com/graph/topological-sort.html) | Article | Concise reference with pseudocode for both algorithms, proofs of correctness, and implementation notes. |
| 5 | [Course Schedule I — NeetCode](https://www.youtube.com/watch?v=EgI5nU9etnU) | Video (10 min) | Directly maps cycle detection to the LeetCode problem. Shows DFS approach with clean code. |
| 6 | [Course Schedule II — NeetCode](https://www.youtube.com/watch?v=Akt3glAwyfY) | Video (10 min) | Extends Course Schedule I to produce the actual ordering. Shows how to adapt the solution. |
| 7 | [Topological Sort in Real Systems — Marc Brooker](https://brooker.co.za/blog/2012/11/07/topological-sort.html) | Blog | Short post on how topological sort drives real build systems and dependency resolution. Good for motivation. |
| 8 | [Alien Dictionary — LeetCode Discussion](https://leetcode.com/problems/alien-dictionary/editorial/) | Editorial | Stretch problem that combines graph construction from constraints with topological sort. Read after implementing the core problems. |

**Suggested order:** Resources 1-3 during the review block, 4 as a reference to keep open, 5-6 during implementation, 7-8 as follow-up.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 - 12:20 — Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:05 | Read through the core concepts section below. Focus on *why* topological sort only works on DAGs. |
| 12:05 - 12:12 | **Trace Kahn's algorithm by hand.** Use the DAG in Section 6. On paper: write the in-degree array, initialize the queue, and step through. At each step, write the queue state and the output so far. |
| 12:12 - 12:18 | **Trace DFS-based topo sort by hand.** Same DAG. Label each node with discovery/finish times. Write the post-order, then reverse it. Verify it matches a valid topological order. |
| 12:18 - 12:20 | Compare the two traces. Note: both produce valid orderings but they may differ. Convince yourself both are correct by checking every edge u->v has u before v. |

### 12:20 - 1:20 — Implement

| Time | Activity |
|------|----------|
| 12:20 - 12:30 | **Build adjacency list + in-degree array from edge list.** Write the `buildGraph` helper. This is shared infrastructure for everything that follows. |
| 12:30 - 12:48 | **Implement `KahnsTopoSort`.** Build in-degree array, seed queue with zero-in-degree nodes, BFS loop decrementing in-degrees and enqueuing newly-zero nodes. Return the order and whether a cycle was detected (len(result) < n). |
| 12:48 - 12:52 | **Test Kahn's.** Run against the example DAG, a graph with a cycle, a single node, and disconnected nodes. |
| 12:52 - 1:08 | **Implement `DFSTopoSort`.** Three-color marking (0=white, 1=gray, 2=black). For each unvisited node, run DFS. On entering a node, mark gray. On finishing, mark black and append to stack. If you visit a gray node, there's a cycle. Reverse the stack at the end. |
| 1:08 - 1:12 | **Test DFS-based.** Same test cases as Kahn's. Verify both produce valid (though possibly different) orderings. |
| 1:12 - 1:20 | **Implement `CanFinish` (Course Schedule I).** Thin wrapper around either algorithm — return true if no cycle is detected. Then implement `FindOrder` (Course Schedule II) — return the topological order, or empty slice if a cycle exists. |

### 1:20 - 1:50 — Solidify

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | **Edge cases.** Test: self-loops, graph where every node is independent (no edges), graph that's a single long chain, graph with multiple connected components. |
| 1:30 - 1:40 | **Lexicographically smallest ordering.** Modify Kahn's to use a min-heap instead of a plain queue. Implement and test. |
| 1:40 - 1:50 | **Compare approaches.** Write down: when would you choose Kahn's vs DFS-based? (Kahn's: more intuitive, natural cycle detection, easy to modify for lex order. DFS: more compact, natural for recursive thinkers, easily extends to longest-path-in-DAG.) |

### 1:50 - 2:00 — Recap (From Memory)

Write down without looking:

1. Kahn's algorithm in 4 steps.
2. DFS-based topo sort in 3 steps.
3. How each detects cycles.
4. Time and space complexity of both: O(V + E) time, O(V + E) space.
5. One gotcha you hit during implementation.

---

## 3. Core Concepts Deep Dive

### What is a DAG?

A **Directed Acyclic Graph** is a directed graph with no cycles. "Acyclic" means there is no path from any node back to itself following edge directions. Topological sort is defined only for DAGs because:

- A topological order requires that for every edge u -> v, u appears before v.
- If there's a cycle u -> ... -> u, then u must appear before itself, which is a contradiction.
- Therefore: **a topological ordering exists if and only if the graph is a DAG.**

A DAG with V vertices always has at least one node with in-degree 0 (if every node had in-degree >= 1, you could follow incoming edges backward indefinitely, which in a finite graph means a cycle).

### Kahn's Algorithm (BFS-Based)

**Intuition:** Repeatedly peel off nodes that have no remaining prerequisites.

**How it works:**
1. Compute the in-degree of every node.
2. Put all nodes with in-degree 0 into a queue — these have no dependencies.
3. Dequeue a node u, add it to the result. For each outgoing edge u -> v, decrement v's in-degree. If v's in-degree hits 0, enqueue it.
4. Repeat until the queue is empty.

**Why it's correct:** A node with in-degree 0 has no unsatisfied dependencies. Removing it and its edges can't violate the ordering of any remaining node, because no remaining node depended on a node that hasn't been processed yet. By induction, every node is placed into the result after all its predecessors.

**Cycle detection:** If the result contains fewer than V nodes when the queue empties, some nodes still have non-zero in-degree. Those nodes form (or depend on) a cycle — they can never reach in-degree 0 because they're waiting on each other.

### DFS-Based Topological Sort

**Intuition:** A node should appear in the ordering only after all nodes it can reach have been placed.

**How it works:**
1. Run DFS on every unvisited node.
2. Use three colors: white (unvisited), gray (in current DFS path), black (finished).
3. When a node finishes (all descendants fully explored), append it to a list.
4. Reverse the list at the end. This reversed post-order is a valid topological sort.

**Why reversed post-order works:** If there's an edge u -> v, then in a DFS:
- If v is visited before u's DFS call finishes, v will finish before u (it's deeper in the recursion). So v gets a smaller post-order number. After reversing, u comes before v.
- If v is already black when u visits it, v already finished and has a smaller post-order number. After reversing, u comes before v.
- If v is gray when u visits it, that's a back edge — meaning a cycle, which we reject.

In all valid cases, u ends up before v after reversal.

**Cycle detection:** If you encounter a gray node during DFS, you've found a back edge — a path from the gray node to itself through the current DFS path. This means a cycle exists.

### Kahn's vs DFS-Based: Comparison

| Aspect | Kahn's (BFS) | DFS-Based |
|--------|-------------|-----------|
| **Approach** | Peel off zero-in-degree nodes | Reversed post-order |
| **Data structure** | Queue + in-degree array | Recursion stack + color array |
| **Cycle detection** | Natural: result size < V | Requires three-color marking |
| **Intuition** | "Remove nodes with no prerequisites" | "Finish dependencies first" |
| **Lex smallest order** | Easy: replace queue with min-heap | Harder to adapt |
| **Code length** | Slightly longer (in-degree setup) | More compact |
| **Iterative?** | Naturally iterative | Naturally recursive |

**When to choose which:**
- **Kahn's** when you need cycle detection as a byproduct, lex-smallest ordering, or an iterative solution.
- **DFS** when you're already doing DFS for other reasons (e.g., longest path in DAG) or prefer recursive code.

### Multiple Valid Orderings

Most DAGs have more than one valid topological order. For example, if nodes A and B both have in-degree 0, either can come first. The algorithms produce one valid order depending on implementation details (queue ordering, DFS start order).

**Lexicographically smallest order:** Replace the queue in Kahn's with a min-heap (priority queue). This always processes the smallest-labeled zero-in-degree node first.

### Longest Path in a DAG

Process nodes in topological order. For each node u, relax all outgoing edges:
```
dist[v] = max(dist[v], dist[u] + weight(u, v))
```

This works because topological order guarantees u is processed before v. Unlike general graphs, longest path in a DAG is solvable in O(V + E) — no NP-hardness here.

**Application:** Critical path analysis in project management. The longest path through the dependency graph determines the minimum project duration.

---

## 4. Implementation Checklist

### Function Signatures

```go
package topo

// KahnsTopoSort returns a topological ordering using Kahn's algorithm.
// Returns the order and true if a cycle is detected.
func KahnsTopoSort(n int, edges [][]int) ([]int, bool) { ... }

// DFSTopoSort returns a topological ordering using DFS post-order reversal.
// Returns the order and true if a cycle is detected.
func DFSTopoSort(n int, edges [][]int) ([]int, bool) { ... }

// CanFinish returns true if all courses can be completed (no cycles).
// prereqs[i] = [course, prerequisite] means prerequisite -> course.
func CanFinish(numCourses int, prereqs [][]int) bool { ... }

// FindOrder returns a valid course order, or an empty slice if impossible.
func FindOrder(numCourses int, prereqs [][]int) []int { ... }

// AlienDictionary returns the character ordering given a sorted list of words
// in an alien language, or "" if the ordering is invalid.
func AlienDictionary(words []string) string { ... }
```

### Helper Functions

```go
// buildGraph constructs an adjacency list from an edge list.
// edges[i] = [from, to] means from -> to.
func buildGraph(n int, edges [][]int) [][]int { ... }

// buildGraphAndInDegree constructs adjacency list and in-degree array.
func buildGraphAndInDegree(n int, edges [][]int) ([][]int, []int) { ... }
```

### Test Cases

| Test Case | Input Description | Expected Behavior |
|-----------|-------------------|-------------------|
| Simple DAG | 4 nodes, edges: 0->1, 0->2, 1->3, 2->3 | Valid topo order (e.g., [0,1,2,3] or [0,2,1,3]) |
| Graph with cycle | 3 nodes, edges: 0->1, 1->2, 2->0 | Cycle detected, no valid order |
| Disconnected graph | 4 nodes, edges: 0->1 (nodes 2,3 isolated) | All 4 nodes in result, isolated nodes can appear anywhere |
| Single node | 1 node, no edges | Order: [0] |
| Self-loop | 2 nodes, edges: 0->0 | Cycle detected |
| Linear chain | 4 nodes, edges: 0->1, 1->2, 2->3 | Exactly one valid order: [0,1,2,3] |
| Diamond shape | 4 nodes, edges: 0->1, 0->2, 1->3, 2->3 | Multiple valid orders |
| Empty graph | 0 nodes, no edges | Empty result |
| All independent | 3 nodes, no edges | Any permutation of [0,1,2] |

### Validation Helper

```go
// isValidTopoOrder checks if order is a valid topological sort of the graph.
func isValidTopoOrder(n int, edges [][]int, order []int) bool {
    if len(order) != n {
        return false
    }
    pos := make(map[int]int)
    for i, v := range order {
        pos[v] = i
    }
    for _, e := range edges {
        if pos[e[0]] >= pos[e[1]] {
            return false // u must come before v
        }
    }
    return true
}
```

---

## 5. Real-World Applications

### Build Systems (make, Bazel, Gradle)

Source files have dependencies: `main.go` imports `utils.go` which imports `types.go`. The build system constructs a dependency graph and topologically sorts it to determine compilation order. If `types.go` hasn't been compiled, `utils.go` can't be compiled. A cycle in the dependency graph means the build is impossible (circular imports).

Bazel and similar systems also use topological order to maximize parallelism: all nodes at the same "depth" (with all dependencies satisfied) can be compiled concurrently.

### Package Managers (npm, go mod, pip)

When you run `go mod download` or `npm install`, the package manager builds a dependency graph of all transitive dependencies and resolves them in topological order. If package A depends on B and C, but B also depends on C, the manager installs C first, then B, then A. Cycle detection catches circular dependency bugs.

### Spreadsheet Cell Evaluation

In a spreadsheet, cell A1 might contain `=B1+C1`, and B1 might contain `=D1*2`. The spreadsheet engine builds a DAG of cell references and evaluates cells in topological order, ensuring every cell is computed after all cells it references. If A1 references B1 and B1 references A1, that's a circular reference error — a cycle in the DAG.

### Course Prerequisites (University Scheduling)

This is literally the Course Schedule problem. Courses are nodes, prerequisites are directed edges. A topological sort gives a valid semester-by-semester plan. If the prerequisite graph has a cycle (Course A requires B, B requires C, C requires A), no valid schedule exists — the registrar has a bug.

### Data Pipeline Orchestration (Apache Airflow)

Airflow represents workflows as DAGs (it's in the name). Each task node (extract data, transform, load, run tests) has dependencies on other tasks. The scheduler topologically sorts the DAG to determine execution order. Tasks with no unsatisfied dependencies are dispatched to workers in parallel — the same principle as Kahn's algorithm, where all zero-in-degree nodes can run concurrently.

---

## 6. Visual Diagrams

### A DAG with Prerequisites

```
    Course prerequisites:
    0: Intro to CS
    1: Data Structures (requires 0)
    2: Algorithms (requires 1)
    3: Databases (requires 0)
    4: Web Dev (requires 1, 3)
    5: Capstone (requires 2, 4)

         0
        / \
       v   v
       1   3
      / \ /
     v   v
     2   4
      \ /
       v
       5

    Edges: 0->1, 0->3, 1->2, 1->4, 3->4, 2->5, 4->5

    Valid topological orders:
      [0, 1, 3, 2, 4, 5]
      [0, 3, 1, 2, 4, 5]
      [0, 1, 2, 3, 4, 5]
      [0, 3, 1, 4, 2, 5]
      ... and more
```

### Kahn's Algorithm Step-by-Step

```
    Graph: 0->1, 0->3, 1->2, 1->4, 3->4, 2->5, 4->5

    Initial in-degrees:
    Node:      0  1  2  3  4  5
    In-degree: 0  1  1  1  2  2

    Step 1: Queue = [0]         Output = []
            Dequeue 0, add to output.
            Decrement neighbors: 1 (1->0), 3 (1->0)
            Enqueue 1, 3 (in-degree became 0)

    Node:      0  1  2  3  4  5
    In-degree: -  0  1  0  2  2

    Step 2: Queue = [1, 3]      Output = [0]
            Dequeue 1, add to output.
            Decrement neighbors: 2 (1->0), 4 (2->1)
            Enqueue 2 (in-degree became 0)

    Node:      0  1  2  3  4  5
    In-degree: -  -  0  0  1  2

    Step 3: Queue = [3, 2]      Output = [0, 1]
            Dequeue 3, add to output.
            Decrement neighbors: 4 (1->0)
            Enqueue 4 (in-degree became 0)

    Node:      0  1  2  3  4  5
    In-degree: -  -  0  -  0  2

    Step 4: Queue = [2, 4]      Output = [0, 1, 3]
            Dequeue 2, add to output.
            Decrement neighbors: 5 (2->1)
            Nothing enqueued.

    Node:      0  1  2  3  4  5
    In-degree: -  -  -  -  0  1

    Step 5: Queue = [4]         Output = [0, 1, 3, 2]
            Dequeue 4, add to output.
            Decrement neighbors: 5 (1->0)
            Enqueue 5 (in-degree became 0)

    Node:      0  1  2  3  4  5
    In-degree: -  -  -  -  -  0

    Step 6: Queue = [5]         Output = [0, 1, 3, 2, 4]
            Dequeue 5, add to output.
            No neighbors.

    Queue = []                  Output = [0, 1, 3, 2, 4, 5]

    len(output) == 6 == n  =>  No cycle. Valid topological order.
```

### DFS-Based Topological Sort

```
    Same graph: 0->1, 0->3, 1->2, 1->4, 3->4, 2->5, 4->5

    Colors: W = white (unvisited), G = gray (in path), B = black (done)

    Start DFS from node 0:

    Call stack         Colors                   Post-order stack
    ─────────────      ──────────────────────   ─────────────────
    dfs(0)             0:G 1:W 2:W 3:W 4:W 5:W  []
      dfs(1)           0:G 1:G 2:W 3:W 4:W 5:W  []
        dfs(2)         0:G 1:G 2:G 3:W 4:W 5:W  []
          dfs(5)       0:G 1:G 2:G 3:W 4:W 5:G  []
          return 5     0:G 1:G 2:G 3:W 4:W 5:B  [5]
        return 2       0:G 1:G 2:B 3:W 4:W 5:B  [5, 2]
        dfs(4)         0:G 1:G 2:B 3:W 4:G 5:B  [5, 2]
          visit 5: already B, skip
        return 4       0:G 1:G 2:B 3:W 4:B 5:B  [5, 2, 4]
      return 1         0:G 1:B 2:B 3:W 4:B 5:B  [5, 2, 4, 1]
      dfs(3)           0:G 1:B 2:B 3:G 4:B 5:B  [5, 2, 4, 1]
        visit 4: already B, skip
      return 3         0:G 1:B 2:B 3:B 4:B 5:B  [5, 2, 4, 1, 3]
    return 0           0:B 1:B 2:B 3:B 4:B 5:B  [5, 2, 4, 1, 3, 0]

    Reversed post-order:  [0, 3, 1, 4, 2, 5]

    Verify every edge u->v has u before v:
      0->1: pos 0 < pos 2  ✓
      0->3: pos 0 < pos 1  ✓
      1->2: pos 2 < pos 4  ✓
      1->4: pos 2 < pos 3  ✓
      3->4: pos 1 < pos 3  ✓
      2->5: pos 4 < pos 5  ✓
    Valid topological order.
```

### Cycle Detection with Kahn's Algorithm

```
    Graph WITH a cycle: 0->1, 1->2, 2->3, 3->1

         0 ──> 1 ──> 2
               ^     |
               |     v
               └──── 3

    Initial in-degrees:
    Node:      0  1  2  3
    In-degree: 0  2  1  1
                  ^
                  └── in-degree 2 because edges from 0 AND 3

    Step 1: Queue = [0]         Output = []
            Dequeue 0, add to output.
            Decrement neighbor 1: in-degree 2 -> 1
            Node 1 still has in-degree 1, NOT enqueued.

    Node:      0  1  2  3
    In-degree: -  1  1  1

    Queue is now EMPTY.          Output = [0]

    len(output) = 1  !=  n = 4

    CYCLE DETECTED.

    Nodes 1, 2, 3 all still have in-degree > 0.
    They form the cycle 1 -> 2 -> 3 -> 1 and can never
    reach in-degree 0 because they depend on each other.
```

---

## 7. Self-Assessment

Answer these without looking at the material above. If you struggle, that's your signal for tomorrow's focus.

**Q1: Can a graph have more than one valid topological ordering? When?**

<details>
<summary>Answer</summary>

Yes. A DAG has multiple valid topological orderings whenever there are two or more nodes that could come next at some point — i.e., when the graph's ordering constraints are not fully determined. For example, if nodes A and B both have in-degree 0 and no edge between them, either can come first. A DAG has exactly one valid topological order only when it has a unique Hamiltonian path (every consecutive pair in the ordering is connected by an edge).

</details>

**Q2: How would you modify Kahn's algorithm to produce the lexicographically smallest topological order?**

<details>
<summary>Answer</summary>

Replace the plain FIFO queue with a min-heap (priority queue). When multiple nodes have in-degree 0 simultaneously, the min-heap always dequeues the smallest-labeled one first. Everything else about the algorithm stays the same.

```go
// Use container/heap or your own min-heap from Day 6
// instead of a plain []int queue
```

</details>

**Q3: In the DFS-based topological sort, why do we need three colors instead of just a visited boolean? What goes wrong with only two states?**

<details>
<summary>Answer</summary>

With only `visited` / `not visited`, you can't distinguish between a node that's currently on the DFS stack (gray) and a node that's fully processed (black). If you visit an already-visited node, you don't know whether it's an ancestor in the current path (a back edge = cycle) or a node fully processed from a different branch (a cross/forward edge = no cycle). Three colors let you detect back edges specifically: visiting a gray node means a cycle; visiting a black node is safe.

</details>

**Q4: You have a DAG and want to find the longest path from a source node to any other node. Describe the approach and its time complexity.**

<details>
<summary>Answer</summary>

Topologically sort the DAG first. Then process nodes in topological order, relaxing outgoing edges:

```
Initialize dist[source] = 0, dist[all others] = -infinity
For each node u in topological order:
    For each edge u -> v with weight w:
        dist[v] = max(dist[v], dist[u] + w)
```

Time complexity: O(V + E) — same as the topological sort plus one pass over all edges. This is efficient because topological order guarantees all paths to u are finalized before processing u's outgoing edges.

</details>

**Q5: A colleague argues that you can detect cycles in an undirected graph using Kahn's algorithm. Are they right? What would you need to change?**

<details>
<summary>Answer</summary>

Not directly. Kahn's algorithm is designed for directed graphs where in-degree has a clear meaning. In an undirected graph, every edge contributes to the degree of both endpoints, and a tree (acyclic) with internal nodes would have degree >= 2, making the in-degree check meaningless in the directed sense.

For undirected cycle detection, you'd use either:
- **DFS** with parent tracking: if you visit a neighbor that's already visited and it's not the parent of the current node, there's a cycle.
- **Union-Find**: for each edge (u, v), if `Find(u) == Find(v)`, there's a cycle.

That said, you *can* adapt a degree-peeling approach for undirected graphs (iteratively remove degree-1 nodes — if all nodes are removed, no cycle; if some remain, they form a cycle), but this isn't standard Kahn's.

</details>
