# Day 12 — Graph Representations & BFS: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Focus | Time |
|---|----------|-------|------|
| 1 | [Graph Algorithms for Technical Interviews — freeCodeCamp](https://www.youtube.com/watch?v=tWVWeAqZ0WU) | Full walkthrough of adjacency list construction, BFS, and DFS with step-by-step code. Starts from scratch with clear whiteboard visuals. Covers connected components and shortest path. | 20 min (watch first 20) |
| 2 | [BFS Shortest Path — William Fiset](https://www.youtube.com/watch?v=oDqjPvD54Ss) | Animated BFS showing the queue state at every step and how distances are computed layer-by-layer. The best visualization of why BFS guarantees shortest paths in unweighted graphs. | 12 min |
| 3 | [Rotten Oranges — NeetCode](https://www.youtube.com/watch?v=y704fEOx0s0) | Multi-source BFS explained visually on a grid. Shows how seeding the queue with multiple start nodes simulates simultaneous expansion. The canonical multi-source BFS problem. | 12 min |
| 4 | [Graph Representation — VisuAlgo](https://visualgo.net/en/graphds) | Interactive visualization of adjacency list, adjacency matrix, and edge list. Toggle between representations for the same graph. Manipulate edges and watch the structures update live. | 10 min |
| 5 | [BFS Interactive Visualization — VisuAlgo](https://visualgo.net/en/dfsbfs) | Step through BFS one node at a time. Watch the queue, visited set, and traversal order update in real time. Toggle between BFS and DFS on the same graph to compare behavior. | 10 min |
| 6 | [Number of Islands — NeetCode](https://www.youtube.com/watch?v=pV2kpPD66nE) | Grid BFS treating each cell as a graph node with 4-directional neighbors. Shows how to count connected components on a grid. Directly applies BFS and visited marking. | 10 min |
| 7 | [0-1 BFS — CP-Algorithms](https://cp-algorithms.com/graph/01_bfs.html) | Concise explanation of 0-1 BFS using a deque. Shows why pushing weight-0 edges to the front and weight-1 edges to the back gives correct shortest paths in O(V+E) without Dijkstra. | 10 min |
| 8 | [Word Ladder — NeetCode](https://www.youtube.com/watch?v=h9iTnkgv05E) | BFS on an implicit graph where nodes are words and edges connect words differing by one character. Excellent example of BFS applied to non-obvious graph structures. | 12 min |

---

## 2. Detailed 2-Hour Session Plan

### 12:00 — 12:20 | Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:06 | Read Section 3 below (graph terminology and representations). Study the three representations for the same graph. For each one, answer: "How do I check if edge (2,3) exists?" and "How do I iterate all neighbors of node 2?" |
| 12:06 - 12:12 | Study the BFS algorithm (Section 3.3). On paper, trace BFS starting from node 0 on the example graph in Section 6.1. Draw the queue state at each step. Verify your trace matches the diagram in Section 6.2. |
| 12:12 - 12:18 | Read multi-source BFS (Section 3.4) and 0-1 BFS (Section 3.5). Mentally trace the multi-source BFS diagram in Section 6.3 — two sources expanding simultaneously. Understand why seeding the queue with multiple nodes is equivalent to adding a virtual super-source connected to all real sources. |
| 12:18 - 12:20 | Review the complexity table. Say out loud: "BFS is O(V+E) time, O(V) space. Adjacency list is O(V+E) space. Adjacency matrix is O(V^2) space." Review the "when to use each representation" rules. |

### 12:20 — 1:20 | Implement (From Scratch)

| Time | Problem | Notes |
|------|---------|-------|
| 12:20 - 12:35 | `BuildGraph` | Build an adjacency list from an edge list. Handle both directed and undirected (add edge in both directions). Test: empty graph, single node, self-loop, disconnected nodes. Then build one version that takes `[][]int` edges and one that takes `[][2]int`. |
| 12:35 - 12:55 | `BFS` | Implement BFS returning traversal order. Use a `[]int` slice as queue. Mark visited on enqueue, not dequeue. Test: linear chain (0→1→2→3), complete graph, disconnected graph (should only visit reachable nodes), single node, graph with cycle. Write `ShortestPathUnweighted` next — BFS with a `dist []int` array initialized to -1. Return `dist[end]`. Test: path exists, no path exists (disconnected), start == end, graph with cycle. |
| 12:55 - 13:10 | `CountComponents` | Loop over all nodes. If not visited, run BFS and increment count. Test: fully connected graph → 1, all isolated nodes → n, two components, single node. |
| 13:10 - 13:20 | `OrangesRotting` (multi-source BFS) | Seed queue with all initially rotten oranges. BFS layer-by-layer. Count layers (minutes). At end, check if any fresh orange remains. Test: all rotten, all fresh (no rotten → -1), no oranges, single orange, fresh orange unreachable. |

### 1:20 — 1:50 | Solidify (Edge Cases & Variants)

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | Grid BFS: `ShortestPathGrid` — given a 2D grid of 0s and 1s, find shortest path from top-left to bottom-right through 0-cells. Convert grid position to BFS: neighbors are 4-directional `(r-1,c), (r+1,c), (r,c-1), (r,c+1)`. Use bounds checking. Test: 1x1 grid, no path, straight line, large open grid. |
| 1:30 - 1:40 | Go back to `BFS` and `ShortestPathUnweighted`. Add test for directed graph where path exists A→B but not B→A. Add test for graph with parallel edges. Verify `CountComponents` works on a directed graph (note: components in directed graphs are different — discuss weakly vs strongly connected). |
| 1:40 - 1:50 | Review all implementations. Verify every BFS marks visited on enqueue. Verify `BuildGraph` handles both directed and undirected. Extract common BFS skeleton into a comment and compare across your three BFS-based functions — they should share the same core loop. |

### 1:50 — 2:00 | Recap (From Memory)

Write down without looking:
1. The BFS algorithm in pseudocode (queue, visited, dequeue-then-enqueue-neighbors loop).
2. Why BFS finds shortest paths in unweighted graphs (it explores all nodes at distance d before distance d+1).
3. The three graph representations and when to use each.
4. How multi-source BFS works (seed queue with all sources at distance 0).
5. Why you must mark visited on enqueue, not dequeue.

---

## 3. Core Concepts Deep Dive

### 3.1 Graph Terminology

**Directed vs Undirected:** In a directed graph, edge (u, v) goes from u to v only. In an undirected graph, (u, v) means both u→v and v→u. When building an adjacency list for an undirected graph, add the edge in both directions.

**Weighted vs Unweighted:** Unweighted means every edge has implicit weight 1 (or equivalently, all edges are equal). BFS finds shortest paths in unweighted graphs. For weighted graphs, you need Dijkstra or Bellman-Ford.

**Connected vs Disconnected:** An undirected graph is connected if there is a path between every pair of vertices. A directed graph is **strongly connected** if there is a directed path between every pair. A single BFS from any node reaches all nodes in a connected undirected graph. In a disconnected graph, you need to run BFS from each unvisited node to find all components.

**Cyclic vs Acyclic:** A cycle is a path from a node back to itself. Trees are connected acyclic graphs. DAGs (directed acyclic graphs) have directed edges but no cycles — they admit topological orderings.

**Sparse vs Dense:** A graph with V vertices can have at most V*(V-1)/2 edges (undirected) or V*(V-1) edges (directed). Sparse graphs have E much less than V^2 (most real-world graphs). Dense graphs have E close to V^2 (rare). The distinction matters for representation choice: adjacency list is O(V+E) space, matrix is O(V^2).

**Degree:** In an undirected graph, the degree of a node is the number of edges incident to it. In a directed graph, **in-degree** is edges coming in, **out-degree** is edges going out. The sum of all degrees = 2 * |E| (each edge contributes to two nodes' degrees).

---

### 3.2 Three Representations Compared

#### Adjacency List

Each node stores a list of its neighbors. The most common representation for interview problems and sparse graphs.

```go
// Unweighted
graph := make([][]int, n)
for _, e := range edges {
    u, v := e[0], e[1]
    graph[u] = append(graph[u], v)
    graph[v] = append(graph[v], u) // omit for directed
}

// Weighted
type Edge struct {
    To, Weight int
}
graph := make([][]Edge, n)
for _, e := range edges {
    u, v, w := e[0], e[1], e[2]
    graph[u] = append(graph[u], Edge{v, w})
    graph[v] = append(graph[v], Edge{u, w}) // omit for directed
}
```

**Edge lookup:** O(degree) — scan the neighbor list. Not great for "does edge (u,v) exist?" queries, but you rarely need that in BFS/DFS problems.

**Iterate neighbors:** O(degree) — just range over `graph[u]`. This is the critical operation for BFS/DFS, and adjacency list makes it optimal.

**Space:** O(V + E). Each edge is stored once (directed) or twice (undirected).

**When to use:** Almost always. Sparse graphs, BFS, DFS, topological sort, any traversal-based algorithm.

#### Adjacency Matrix

A V×V boolean (or integer) matrix where `matrix[u][v]` indicates whether edge (u,v) exists (or its weight).

```go
matrix := make([][]bool, n)
for i := range matrix {
    matrix[i] = make([]bool, n)
}
for _, e := range edges {
    u, v := e[0], e[1]
    matrix[u][v] = true
    matrix[v][u] = true // omit for directed
}
```

**Edge lookup:** O(1) — just check `matrix[u][v]`.

**Iterate neighbors:** O(V) — must scan the entire row, even if the node has only 2 neighbors out of 10,000 nodes. This is the main drawback.

**Space:** O(V^2), regardless of the number of edges.

**When to use:** Dense graphs where E ≈ V^2 (rare in interviews). Also useful when you need O(1) edge existence checks (e.g., checking if two nodes are directly connected in a constraint). Floyd-Warshall all-pairs shortest path uses a matrix.

#### Edge List

A flat list of edges. Each edge is a pair (or triple with weight).

```go
type EdgeEntry struct {
    From, To, Weight int
}
edges := []EdgeEntry{
    {0, 1, 5},
    {0, 2, 3},
    {1, 2, 1},
}
```

**Edge lookup:** O(E) — linear scan.

**Iterate neighbors:** O(E) — scan all edges to find ones involving the node. Terrible for traversal.

**Space:** O(E). Very compact.

**When to use:** Kruskal's MST (sort edges by weight, process in order). Also useful as input format that you convert to adjacency list. You almost never traverse with an edge list directly.

---

### 3.3 BFS: The Queue-Based Algorithm

BFS explores nodes in order of their distance from the source. It processes all nodes at distance d before any at distance d+1. This **layer-by-layer** expansion is what makes it find shortest paths in unweighted graphs.

**Algorithm:**

```go
func BFS(graph [][]int, start int) []int {
    n := len(graph)
    visited := make([]bool, n)
    visited[start] = true
    queue := []int{start}
    order := []int{}

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)

        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true  // mark on ENQUEUE
                queue = append(queue, neighbor)
            }
        }
    }
    return order
}
```

**Why it finds shortest paths:** Consider the queue at any point during execution. It contains nodes from at most two consecutive layers — the current layer being processed and the next layer being enqueued. When we dequeue a node at distance d, its unvisited neighbors are at distance d+1. Because we process all distance-d nodes before any distance-(d+1) nodes, the first time we reach any node is via a shortest path.

**Critical detail — mark visited on ENQUEUE, not DEQUEUE:**

If you mark visited when dequeuing, the same node can be added to the queue multiple times by different neighbors before it gets dequeued. This wastes time and space (the queue can grow to O(E) instead of O(V)), and in the worst case can cause incorrect distance calculations. Marking on enqueue guarantees each node enters the queue exactly once.

```
WRONG:                              CORRECT:
queue.add(start)                    visited[start] = true
while queue not empty:              queue.add(start)
    node = queue.dequeue()          while queue not empty:
    if visited[node]: continue          node = queue.dequeue()
    visited[node] = true                for neighbor in adj[node]:
    for neighbor in adj[node]:              if not visited[neighbor]:
        queue.add(neighbor)  ←DUPS!             visited[neighbor] = true
                                                queue.add(neighbor)
```

**Level-by-level BFS:** To track the layer (distance), process all nodes at the current level before moving to the next:

```go
func BFSLevels(graph [][]int, start int) [][]int {
    n := len(graph)
    visited := make([]bool, n)
    visited[start] = true
    queue := []int{start}
    levels := [][]int{}

    for len(queue) > 0 {
        levelSize := len(queue)
        level := []int{}
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            level = append(level, node)
            for _, neighbor := range graph[node] {
                if !visited[neighbor] {
                    visited[neighbor] = true
                    queue = append(queue, neighbor)
                }
            }
        }
        levels = append(levels, level)
    }
    return levels
}
```

**Complexity:** O(V + E) time (each vertex dequeued once, each edge examined once). O(V) space (visited array + queue holding at most O(V) nodes).

---

### 3.4 Multi-Source BFS

Standard BFS starts from one source. Multi-source BFS starts from **multiple sources simultaneously**. Instead of seeding the queue with one node, seed it with all source nodes at distance 0.

**Why it works:** Conceptually, it is equivalent to adding a virtual "super-source" node connected to all real sources with weight-0 edges, then running standard BFS from the super-source. The real sources are all at distance 0, and BFS expands outward from all of them in parallel.

**Algorithm:**

```go
func MultiSourceBFS(graph [][]int, sources []int) []int {
    n := len(graph)
    dist := make([]int, n)
    for i := range dist {
        dist[i] = -1
    }

    queue := []int{}
    for _, s := range sources {
        dist[s] = 0
        queue = append(queue, s)
    }

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        for _, neighbor := range graph[node] {
            if dist[neighbor] == -1 {
                dist[neighbor] = dist[node] + 1
                queue = append(queue, neighbor)
            }
        }
    }
    return dist // dist[i] = shortest distance from any source to node i
}
```

**Classic application — Rotten Oranges:**
- Grid of oranges: 0 = empty, 1 = fresh, 2 = rotten.
- Each minute, rotten oranges infect adjacent fresh oranges (4-directional).
- Question: minimum minutes until all oranges are rotten (or -1 if impossible).
- Solution: seed queue with all initially rotten oranges. BFS layer-by-layer. Each layer = 1 minute. After BFS, check if any fresh oranges remain.

**Other applications:** Shortest distance from any police station to every cell on a map. Shortest distance from any gate in a maze. Any "minimum distance from the nearest X" problem.

---

### 3.5 0-1 BFS

Standard BFS works for unweighted graphs (all edges have weight 1). What if edges have weight 0 or 1? You could use Dijkstra, but that is O((V+E) log V). **0-1 BFS** solves this in O(V + E) using a **deque** instead of a regular queue.

**Key insight:** In BFS, the queue always has nodes in order of their distance (all distance-d nodes before distance-(d+1) nodes). With 0-weight edges, a neighbor might be at the *same* distance as the current node. To maintain sorted order in the deque:
- Weight-0 edge: push neighbor to the **front** of the deque (same distance as current).
- Weight-1 edge: push neighbor to the **back** of the deque (distance + 1).

```go
func ZeroOneBFS(graph [][]Edge, start int) []int {
    n := len(graph)
    dist := make([]int, n)
    for i := range dist {
        dist[i] = math.MaxInt
    }
    dist[start] = 0

    // Use a slice-based deque
    deque := []int{start}

    for len(deque) > 0 {
        node := deque[0]
        deque = deque[1:]

        for _, edge := range graph[node] {
            newDist := dist[node] + edge.Weight
            if newDist < dist[edge.To] {
                dist[edge.To] = newDist
                if edge.Weight == 0 {
                    deque = append([]int{edge.To}, deque...) // push front
                } else {
                    deque = append(deque, edge.To) // push back
                }
            }
        }
    }
    return dist
}
```

**Note on Go performance:** `append([]int{x}, deque...)` for push-front is O(n). For production code, use a proper doubly-linked list or ring buffer. For interview coding, this is acceptable and clear.

**When to use:** Any shortest-path problem where edges have only two possible weights (0 and 1). Example: minimum flips to connect two cells in a grid where some passages are open (cost 0) and others need a door opened (cost 1).

---

### 3.6 Bidirectional BFS

Instead of searching from just the start, search from **both the start and the end** simultaneously. Alternate expanding one level from each side. Stop when the two searches meet.

**Why it helps:** Standard BFS explores all nodes up to distance d, which is roughly O(b^d) nodes where b is the branching factor. Bidirectional BFS explores O(b^(d/2)) from each side, totaling O(2 * b^(d/2)) — exponentially less when b is large.

**Algorithm sketch:**

```go
func BidirectionalBFS(graph [][]int, start, end int) int {
    if start == end {
        return 0
    }
    visitedStart := map[int]int{start: 0} // node -> distance from start
    visitedEnd := map[int]int{end: 0}     // node -> distance from end
    queueStart := []int{start}
    queueEnd := []int{end}

    for len(queueStart) > 0 && len(queueEnd) > 0 {
        // Always expand the smaller frontier
        if len(queueStart) <= len(queueEnd) {
            if dist := expandLevel(graph, &queueStart, visitedStart, visitedEnd); dist >= 0 {
                return dist
            }
        } else {
            if dist := expandLevel(graph, &queueEnd, visitedEnd, visitedStart); dist >= 0 {
                return dist
            }
        }
    }
    return -1 // no path
}
```

**When to use:** Large search spaces with high branching factor where you know both start and end. Classic example: **Word Ladder** — find shortest transformation from "hit" to "cog" changing one letter at a time. The implicit graph has huge branching. Bidirectional BFS dramatically prunes the search.

**Limitation:** Requires knowing the end node. Does not apply to problems like "find shortest path to any node satisfying condition X" unless you can enumerate all targets.

---

## 4. Implementation Checklist

### Function Signatures

```go
package graphs

// BuildGraph constructs an adjacency list from a list of edges.
// If directed is false, each edge is added in both directions.
func BuildGraph(n int, edges [][]int, directed bool) [][]int { ... }

// BFS returns the traversal order starting from start.
func BFS(graph [][]int, start int) []int { ... }

// ShortestPathUnweighted returns the shortest distance from start to end
// in an unweighted graph. Returns -1 if no path exists.
func ShortestPathUnweighted(graph [][]int, start, end int) int { ... }

// CountComponents returns the number of connected components in an
// undirected graph.
func CountComponents(n int, graph [][]int) int { ... }

// OrangesRotting returns the minimum number of minutes until no fresh
// orange remains. Returns -1 if impossible. Grid values: 0=empty,
// 1=fresh, 2=rotten.
func OrangesRotting(grid [][]int) int { ... }
```

### Test Cases & Edge Cases

| Function | Must-Test Cases |
|----------|----------------|
| `BuildGraph` | Empty edge list → n isolated nodes; single edge directed vs undirected; self-loop `[1,1]`; duplicate edges; edges referencing node 0 and node n-1 (boundary). |
| `BFS` | Linear chain 0→1→2→3 from node 0 → `[0,1,2,3]`; complete graph; graph with cycle; single node; start node with no neighbors (isolated). |
| `ShortestPathUnweighted` | Path exists → correct distance; no path (disconnected) → `-1`; start == end → `0`; cycle in path (should still find shortest); directed graph where reverse path does not exist. |
| `CountComponents` | Fully connected → `1`; all isolated → `n`; two components of different sizes; single node → `1`; graph with self-loop (still 1 component if connected). |
| `OrangesRotting` | All rotten → `0`; all fresh with adjacent rotten → correct minutes; unreachable fresh → `-1`; no oranges at all → `0`; single fresh orange with no rotten → `-1`; single rotten orange → `0`; 1×1 grid. |

---

## 5. BFS Application Patterns

### Pattern 1: Shortest Path in Unweighted Graph

The foundational BFS application. Since BFS explores layer-by-layer, the first time you reach a node is via a shortest path (measured in number of edges).

**Template:** Run BFS from the source, maintaining a `dist[]` array. Return `dist[target]`.

**Examples:** Shortest path between two nodes in a social network, minimum moves for a knight on a chessboard, minimum operations to transform one state to another.

### Pattern 2: Level-Order Traversal of a Tree

A tree is just a connected acyclic graph. BFS on a tree is level-order traversal. No `visited` array needed if you track the parent (or just use visited — it is simpler).

**Relationship to Day 4:** `LevelOrder` from binary trees is BFS where each node has at most 2 children. Graph BFS generalizes this to arbitrary branching.

### Pattern 3: Multi-Source Shortest Distance

Seed the queue with all source nodes at distance 0. BFS outward. Each node's distance is the minimum distance to any source.

**Examples:** Rotten oranges, walls and gates, shortest distance from any building. Any time the problem says "minimum distance to the *nearest* X."

### Pattern 4: Word Ladder (BFS on Implicit Graph)

Not all BFS problems give you an explicit adjacency list. In Word Ladder, the graph is implicit:
- **Nodes:** dictionary words.
- **Edges:** two words are connected if they differ by exactly one letter.

You build neighbors on the fly during BFS rather than precomputing the full graph. For efficiency, either:
- Try all 26 letter substitutions at each position: O(L * 26) per word.
- Or precompute a pattern map: `"h*t" → [hit, hot, hat]` and connect words sharing a pattern.

**Key insight:** Whenever a problem says "find the minimum number of transformations/steps," think BFS on an implicit graph where states are nodes and valid transitions are edges.

### Pattern 5: Grid BFS

Treat a 2D grid as a graph. Each cell `(r, c)` is a node. Neighbors are the 4-directional (or 8-directional) adjacent cells.

**Neighbor generation:**

```go
dirs := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
for _, d := range dirs {
    nr, nc := r+d[0], c+d[1]
    if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
        // valid neighbor
    }
}
```

**Examples:** Shortest path in a maze, flood fill, number of islands (BFS from each unvisited land cell, count components).

**Grid BFS vs Graph BFS:** The only difference is how you represent nodes and generate neighbors. Instead of `graph[node]`, you compute `(nr, nc)` from directions. Instead of `visited[node]`, you use `visited[r][c]`. Everything else (queue, mark-on-enqueue, layer counting) is identical.

---

## 6. Visual Diagrams

### 6.1 A Small Graph in All Three Representations

**Graph:**

```
    0 --- 1
    |   / |
    |  /  |
    | /   |
    2 --- 3
          |
          4
```

Edges: (0,1), (0,2), (1,2), (1,3), (2,3), (3,4). Undirected. 5 nodes, 6 edges.

**Adjacency List:**

```
0: [1, 2]
1: [0, 2, 3]
2: [0, 1, 3]
3: [1, 2, 4]
4: [3]
```

Space: O(V + 2E) = O(5 + 12) = O(17) entries.
Check edge (2,3): scan list for node 2 → [0, 1, 3] → found at position 2. O(degree).
Iterate neighbors of node 1: [0, 2, 3]. O(degree).

**Adjacency Matrix:**

```
    0   1   2   3   4
0 [ .   1   1   .   . ]
1 [ 1   .   1   1   . ]
2 [ 1   1   .   1   . ]
3 [ .   1   1   .   1 ]
4 [ .   .   .   1   . ]
```

Space: O(V^2) = O(25) entries. Mostly empty for this sparse graph.
Check edge (2,3): matrix[2][3] = 1. O(1).
Iterate neighbors of node 1: scan row 1 → columns 0, 2, 3 are 1. O(V).

**Edge List:**

```
[(0,1), (0,2), (1,2), (1,3), (2,3), (3,4)]
```

Space: O(E) = O(6) entries. Most compact.
Check edge (2,3): linear scan all 6 edges. O(E).
Iterate neighbors of node 1: scan all 6 edges for any involving node 1. O(E).

---

### 6.2 BFS Traversal with Queue State

**Graph (adjacency list):**

```
0: [1, 2]
1: [0, 2, 3]
2: [0, 1, 3]
3: [1, 2, 4]
4: [3]
```

**BFS from node 0:**

```
Layer 0 (distance 0):
─────────────────────
  Dequeue: 0
  Neighbors of 0: [1, 2]
    1 not visited → mark visited, enqueue
    2 not visited → mark visited, enqueue

  Queue:  [1, 2]              Visited: {0, 1, 2}
  Order:  [0]

Layer 1 (distance 1):
─────────────────────
  Dequeue: 1
  Neighbors of 1: [0, 2, 3]
    0 visited → skip
    2 visited → skip
    3 not visited → mark visited, enqueue

  Queue:  [2, 3]              Visited: {0, 1, 2, 3}
  Order:  [0, 1]

  Dequeue: 2
  Neighbors of 2: [0, 1, 3]
    0 visited → skip
    1 visited → skip
    3 visited → skip

  Queue:  [3]                 Visited: {0, 1, 2, 3}
  Order:  [0, 1, 2]

Layer 2 (distance 2):
─────────────────────
  Dequeue: 3
  Neighbors of 3: [1, 2, 4]
    1 visited → skip
    2 visited → skip
    4 not visited → mark visited, enqueue

  Queue:  [4]                 Visited: {0, 1, 2, 3, 4}
  Order:  [0, 1, 2, 3]

Layer 3 (distance 3):
─────────────────────
  Dequeue: 4
  Neighbors of 4: [3]
    3 visited → skip

  Queue:  []                  Visited: {0, 1, 2, 3, 4}
  Order:  [0, 1, 2, 3, 4]

─── Done ───

Distances from node 0:
  0→0: 0   0→1: 1   0→2: 1   0→3: 2   0→4: 3

Traversal tree (shortest-path tree):
        0
       / \
      1   2
      |
      3
      |
      4
```

**Key observation:** Nodes 1 and 2 are both at distance 1 (same layer). Node 3 is at distance 2 even though it is a neighbor of both 1 and 2 — BFS reaches it first through 1 (but 2 also tries; it is already visited). Node 4 is at distance 3 (only reachable through 3).

---

### 6.3 Multi-Source BFS: Two Starting Nodes

**Grid (5×5):**

```
. . . . .       S = source (rotten orange)
. . . . .       . = fresh orange
S . . . .       X = empty cell
. . . . .
. . . . S
```

Sources: (2,0) and (4,4). BFS simultaneously from both.

```
Initial (t=0):         After t=1:             After t=2:
. . . . .              . . . . .              . . . . .
. . . . .              1 . . . .              1 2 . . .
0 . . . .              0 1 . . .              0 1 2 . .
. . . . .              1 . . . 1              1 2 . . 1
. . . . 0              . . . 1 0              . . . 1 0
                                               ↑       ↑
                                               from    from
                                               (2,0)   (4,4)

After t=3:             After t=4:
. . . . .              . . . . .
1 2 3 . .              1 2 3 4 .
0 1 2 3 .              0 1 2 3 4
1 2 3 . 1              1 2 3 2 1
. . . 1 0              . . 3 1 0
      ↑                      ↑
      fronts                 fronts
      approaching            met!

After t=5 (done):
5 4 3 4 5
1 2 3 4 3          dist[r][c] = shortest distance
0 1 2 3 2          from EITHER source
1 2 3 2 1
2 3 3 1 0

Key: Each cell shows its distance to the nearest source.
     Both sources expand simultaneously — the waves meet in the middle.
     This is identical to running BFS from a single virtual super-source
     connected to both (2,0) and (4,4) at distance 0.
```

---

## 7. Self-Assessment

Answer these without looking at your code or notes. If you struggle with any, revisit the relevant section.

### Question 1
**Why must you mark nodes visited when enqueuing, not when dequeuing?**

<details>
<summary>Answer</summary>

If you mark visited on dequeue, a node can be added to the queue multiple times before it gets processed. Consider nodes A and B that both have neighbor C. When BFS processes A, it enqueues C. When it processes B (before C is dequeued), it enqueues C again because C is not yet marked visited. This leads to duplicate work (O(E) queue size instead of O(V)), wasted processing, and can cause incorrect distance computations in multi-source or distance-tracking scenarios. Marking on enqueue ensures each node enters the queue exactly once, preserving O(V+E) time and O(V) space guarantees.

</details>

### Question 2
**When would you use an adjacency matrix over an adjacency list?**

<details>
<summary>Answer</summary>

Use an adjacency matrix when:
1. The graph is **dense** (E ≈ V^2), so the matrix space O(V^2) is comparable to the adjacency list space O(V+E).
2. You need **O(1) edge existence checks** — "is there an edge between u and v?" frequently during the algorithm (e.g., checking constraints).
3. The algorithm inherently operates on the matrix, like **Floyd-Warshall** all-pairs shortest path which iterates over all (i,j) pairs.
4. V is small (say V ≤ 1000) and the simpler indexing is worth the space overhead.

For most interview problems, the graph is sparse and you need to iterate neighbors efficiently, so adjacency list is the default choice.

</details>

### Question 3
**Explain multi-source BFS. How does seeding the queue with multiple nodes give correct shortest distances?**

<details>
<summary>Answer</summary>

Multi-source BFS initializes the queue with all source nodes at distance 0 instead of a single source. This is equivalent to adding a virtual super-source node connected to every real source with a weight-0 edge, then running standard BFS from the super-source. The first layer of the virtual BFS reaches all real sources at distance 0, and from there it expands normally. Because BFS explores by increasing distance, each non-source node is first reached from the nearest source, giving the correct shortest distance to any source. The key invariant — nodes at distance d are all processed before nodes at distance d+1 — still holds because all sources start at the same distance (0).

</details>

### Question 4
**You have a graph with edges of weight 0 or 1. Why can't you use standard BFS, and what do you use instead?**

<details>
<summary>Answer</summary>

Standard BFS assumes all edges have equal weight (1). With weight-0 edges, a neighbor across a 0-weight edge is at the **same** distance as the current node, not distance + 1. Standard BFS would place it at the back of the queue (as if it were distance + 1), processing it too late and potentially computing wrong distances.

**0-1 BFS** fixes this by using a **deque** instead of a queue. When traversing a weight-0 edge, push the neighbor to the **front** (same distance = process soon). When traversing a weight-1 edge, push to the **back** (distance + 1 = process later). This maintains the invariant that the deque is sorted by distance, giving correct shortest paths in O(V+E) time — the same as standard BFS, and faster than Dijkstra's O((V+E) log V).

</details>

### Question 5
**You're given a 2D grid and asked to find the shortest path from the top-left to the bottom-right. How do you model this as a graph BFS?**

<details>
<summary>Answer</summary>

Treat each cell (r, c) as a graph node. Edges connect each cell to its 4-directional neighbors (up, down, left, right) that are within bounds and passable (not walls). Run BFS from (0, 0). Use a `visited[rows][cols]` array. Generate neighbors using a direction array `{{-1,0},{1,0},{0,-1},{0,1}}`. The BFS distance to `(rows-1, cols-1)` is the shortest path length. This works because all edges have weight 1 (one step per move), making BFS optimal. If different cells have different traversal costs, you would need Dijkstra or 0-1 BFS instead.

</details>
