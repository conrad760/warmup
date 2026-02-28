# DSA Interview Prep — 21-Day Curriculum

**Schedule:** Daily, 12:00 PM - 2:00 PM (2 hours)
**Language:** Go
**Focus:** Pattern recognition, decision frameworks, and interview problem-solving speed.

## How to Use This Curriculum

Each day's guide focuses on INTERVIEW READINESS:
- **Pattern Catalog** — every pattern for the topic with triggers and Go templates
- **Decision Framework** — how to pick the right approach in under 2 minutes
- **Common Interview Traps** — the mistakes interviewers expect
- **Thought Process Walkthrough** — full interview simulation on representative problems
- **Quick Drill** — micro-exercises to test fluency

**Rule:** If you can't identify the pattern and start coding within 5 minutes of reading a problem, that topic needs more work.

## Daily Session Structure

| Time | Block | What to Do |
|------|-------|------------|
| 12:00-12:20 | Review | Read the day's pattern catalog and decision framework |
| 12:20-1:20 | Practice | Work through the thought process walkthroughs, then the quick drill |
| 1:20-1:50 | Drill | Solve 1-2 additional problems from the topic using the patterns |
| 1:50-2:00 | Recap | From memory: list the patterns, their triggers, and one trap per pattern |

## Schedule (ordered by interview frequency)

### Week 1: High-Frequency Topics
| Day | Topic | Guide |
|-----|-------|-------|
| 1 | Arrays & Hashing Patterns | [day-01-arrays-and-hashing.md](day-01-arrays-and-hashing.md) |
| 2 | Two Pointers & Sliding Window | [day-02-two-pointers-and-sliding-window.md](day-02-two-pointers-and-sliding-window.md) |
| 3 | Stacks | [day-03-stacks.md](day-03-stacks.md) |
| 4 | Binary Search | [day-04-binary-search.md](day-04-binary-search.md) |
| 5 | Linked Lists | [day-05-linked-lists.md](day-05-linked-lists.md) |
| 6 | Binary Trees (DFS & BFS) | [day-06-binary-trees.md](day-06-binary-trees.md) |
| 7 | BST & Tree Construction | [day-07-bst-and-tree-construction.md](day-07-bst-and-tree-construction.md) |

### Week 2: Core Algorithm Patterns
| Day | Topic | Guide |
|-----|-------|-------|
| 8 | Graphs: BFS & DFS | [day-08-graphs-bfs-dfs.md](day-08-graphs-bfs-dfs.md) |
| 9 | Topological Sort | [day-09-topological-sort.md](day-09-topological-sort.md) |
| 10 | Heaps & Priority Queues | [day-10-heaps.md](day-10-heaps.md) |
| 11 | Dynamic Programming: 1D | [day-11-dp-1d.md](day-11-dp-1d.md) |
| 12 | Dynamic Programming: 2D | [day-12-dp-2d.md](day-12-dp-2d.md) |
| 13 | Backtracking | [day-13-backtracking.md](day-13-backtracking.md) |
| 14 | Greedy | [day-14-greedy.md](day-14-greedy.md) |

### Week 3: Remaining Patterns + Mock Practice
| Day | Topic | Guide |
|-----|-------|-------|
| 15 | Intervals & Sweep Line | [day-15-intervals.md](day-15-intervals.md) |
| 16 | Tries & Union-Find | [day-16-tries-and-union-find.md](day-16-tries-and-union-find.md) |
| 17 | Design Problems | [day-17-design.md](day-17-design.md) |
| 18 | Bit Manipulation & Math | [day-18-bits-and-math.md](day-18-bits-and-math.md) |
| 19 | Mock Practice 1 | [day-19-mock-practice-1.md](day-19-mock-practice-1.md) |
| 20 | Mock Practice 2 | [day-20-mock-practice-2.md](day-20-mock-practice-2.md) |
| 21 | Weak Spot Review | [day-21-weak-spot-review.md](day-21-weak-spot-review.md) |

## Deep-Dive Reference Material

The `reference/` subdirectory contains detailed guides on each data structure and algorithm — internal implementations, complexity proofs, Go standard library internals, and curated learning resources. Use these when you need to go deeper on a topic:

| Reference | Description |
|-----------|-------------|
| [day-01-hash-tables.md](reference/day-01-hash-tables.md) | Hash table internals, collision resolution, Go `map` implementation |
| [day-02-linked-lists.md](reference/day-02-linked-lists.md) | Singly/doubly linked list internals, sentinel nodes, Go idioms |
| [day-03-stacks-and-queues.md](reference/day-03-stacks-and-queues.md) | Stack/queue/deque internals, circular buffers, monotonic structures |
| [day-04-binary-trees.md](reference/day-04-binary-trees.md) | Binary tree representations, all four traversals (recursive + iterative) |
| [day-05-binary-search-trees.md](reference/day-05-binary-search-trees.md) | BST operations, AVL/Red-Black overview, validate/kth smallest |
| [day-06-heaps.md](reference/day-06-heaps.md) | Heap internals, sift-up/down, heapify proof, Go `container/heap` |
| [day-07-tries-and-union-find.md](reference/day-07-tries-and-union-find.md) | Trie structure, union-find with path compression and union by rank |
| [day-08-sorting.md](reference/day-08-sorting.md) | Merge sort, quick sort, counting sort, Go `sort` package internals |
| [day-09-binary-search.md](reference/day-09-binary-search.md) | Binary search variants, lower/upper bound, search on answer space |
| [day-10-two-pointers.md](reference/day-10-two-pointers.md) | Opposite ends, fast-slow, Floyd's cycle detection |
| [day-11-sliding-window.md](reference/day-11-sliding-window.md) | Fixed/variable window templates, frequency map techniques |
| [day-12-graphs-and-bfs.md](reference/day-12-graphs-and-bfs.md) | Graph representations, BFS, shortest path, connected components |
| [day-13-dfs-and-backtracking.md](reference/day-13-dfs-and-backtracking.md) | DFS, three-color marking, backtracking template, pruning |
| [day-14-topological-sort.md](reference/day-14-topological-sort.md) | Kahn's algorithm, DFS-based topo sort, cycle detection |
| [day-15-dp-foundations.md](reference/day-15-dp-foundations.md) | DP recipe, memoization vs tabulation, 1D DP patterns |
| [day-16-dp-sequences.md](reference/day-16-dp-sequences.md) | LCS, edit distance, palindromic subsequences, space optimization |
| [day-17-dp-2d-grid.md](reference/day-17-dp-2d-grid.md) | Grid DP, unique paths, min path sum, knapsack variants |
| [day-18-greedy.md](reference/day-18-greedy.md) | Greedy strategy, exchange arguments, interval scheduling |
| [day-19-intervals.md](reference/day-19-intervals.md) | Interval merge/insert, sweep line, meeting rooms |
| [day-20-bits-and-math.md](reference/day-20-bits-and-math.md) | Bit operations, XOR tricks, GCD, modular arithmetic |
| [day-21-design.md](reference/day-21-design.md) | LRU cache, min stack, design patterns for system-level problems |

## Quick Pattern Reference

| Signal | Pattern | Day |
|--------|---------|-----|
| "Find pairs/complement" | Hash map complement | 1 |
| "Frequency count / group by property" | Hash map frequency | 1 |
| "Check for duplicates" | Hash set | 1 |
| "Group anagrams / group by key" | Hash map with sorted/counted key | 1 |
| "Sorted array + find pair" | Two pointers (opposite ends) | 2 |
| "Triplets / k-sum" | Sort + fix one + two pointers | 2 |
| "Remove duplicates in-place" | Two pointers (fast-slow) | 2 |
| "Contiguous subarray max/min" | Sliding window (fixed) | 2 |
| "Longest substring with condition" | Sliding window (variable) | 2 |
| "Minimum window containing X" | Sliding window + frequency map | 2 |
| "Matching/nesting brackets" | Stack | 3 |
| "Next greater/smaller element" | Monotonic stack | 3 |
| "Evaluate expression / RPN" | Stack | 3 |
| "Decode string / nested structure" | Stack | 3 |
| "Min stack / special stack" | Auxiliary stack tracking invariant | 3 |
| "Sorted array + find target" | Binary search (standard) | 4 |
| "First/last occurrence" | Binary search (lower/upper bound) | 4 |
| "Minimum X such that condition" | Binary search on answer space | 4 |
| "Rotated sorted array" | Modified binary search | 4 |
| "Peak finding" | Binary search on condition | 4 |
| "Reverse a linked list" | Three-pointer iterative | 5 |
| "Merge sorted lists" | Dummy head + pointer walk | 5 |
| "Detect cycle" | Floyd's slow/fast pointers | 5 |
| "Find middle node" | Slow/fast pointers | 5 |
| "LRU cache" | Doubly linked list + hash map | 5 |
| "Tree depth/height" | DFS (recursive post-order) | 6 |
| "Level-by-level processing" | BFS with level-size loop | 6 |
| "Path sum / root-to-leaf" | DFS (pre-order with accumulator) | 6 |
| "Invert/mirror tree" | DFS or BFS (any order) | 6 |
| "Serialize/deserialize tree" | Pre-order DFS or level-order BFS | 6 |
| "Lowest common ancestor" | DFS with subtree search | 6 |
| "Validate BST" | In-order traversal or min/max range DFS | 7 |
| "Kth smallest in BST" | In-order traversal stopping at k | 7 |
| "Construct tree from traversals" | Recursive split using in-order index map | 7 |
| "BST iterator" | Controlled in-order with stack | 7 |
| "Shortest path (unweighted)" | BFS | 8 |
| "Connected components / islands" | DFS or BFS flood fill | 8 |
| "Minimum steps / transformations" | BFS (level-by-level) | 8 |
| "Multi-source shortest" | Multi-source BFS | 8 |
| "Cycle in directed graph" | DFS three-color marking | 8 |
| "Dependency ordering" | Topological sort (Kahn's BFS) | 9 |
| "Course schedule / build order" | Topological sort | 9 |
| "Alien dictionary" | Build graph from constraints + topo sort | 9 |
| "Top K elements" | Min-heap of size K | 10 |
| "Kth largest/smallest" | Heap (or quickselect) | 10 |
| "Merge K sorted lists" | Min-heap of list heads | 10 |
| "Running median" | Two heaps (max-heap + min-heap) | 10 |
| "Continuous min/max with removals" | Heap with lazy deletion | 10 |
| "Number of ways" | DP (count paths) | 11 |
| "Minimum cost / maximum value" | DP (optimize over choices) | 11 |
| "Fibonacci / climbing stairs" | DP: dp[i] = dp[i-1] + dp[i-2] | 11 |
| "Rob houses / take-or-skip" | DP: dp[i] = max(take + dp[i-2], dp[i-1]) | 11 |
| "Coin change / unbounded knapsack" | DP: dp[i] = min(dp[i-coin] + 1) | 11 |
| "Longest increasing subsequence" | DP or patience sorting (O(n log n)) | 11 |
| "Compare/align two strings" | 2D DP: dp[i][j] on prefixes | 12 |
| "Edit distance" | 2D DP: insert/delete/replace recurrence | 12 |
| "Longest common subsequence" | 2D DP: match or skip | 12 |
| "Grid paths / min cost path" | 2D DP: dp[r][c] from neighbors | 12 |
| "0/1 knapsack" | 2D DP: take-or-skip with weight constraint | 12 |
| "Generate all subsets" | Backtracking: include/exclude | 13 |
| "Generate all permutations" | Backtracking: swap or used-array | 13 |
| "Combination sum" | Backtracking with pruning | 13 |
| "Word search in grid" | Backtracking DFS on grid | 13 |
| "N-queens / Sudoku" | Backtracking with constraint checking | 13 |
| "Activity selection / max non-overlapping" | Greedy: sort by end time | 14 |
| "Jump game / reach the end" | Greedy: track farthest reachable | 14 |
| "Minimum platforms / resources" | Greedy with sorting | 14 |
| "Assign tasks / pair optimally" | Greedy: sort + pair strategy | 14 |
| "Merge overlapping intervals" | Sort by start + merge | 15 |
| "Insert interval" | Binary search or linear scan + merge | 15 |
| "Meeting rooms (can attend all?)" | Sort by start + check overlap | 15 |
| "Meeting rooms II (min rooms)" | Sweep line or min-heap | 15 |
| "Interval intersection" | Two pointers on sorted intervals | 15 |
| "Autocomplete / prefix matching" | Trie | 16 |
| "Word search II (multiple words)" | Trie + backtracking | 16 |
| "Connected components (dynamic)" | Union-Find | 16 |
| "Number of islands (union variant)" | Union-Find | 16 |
| "Accounts merge / equivalence" | Union-Find | 16 |
| "LRU/LFU cache" | Hash map + linked list / ordered structure | 17 |
| "Implement data structure with X" | Combine primitives (stack, heap, hash map) | 17 |
| "Design iterator / stream" | State machine with pointer | 17 |
| "Single number (find unique)" | XOR all elements | 18 |
| "Power of two / count bits" | Bit manipulation | 18 |
| "Missing number in range" | XOR or math (sum formula) | 18 |
| "GCD / prime factorization" | Euclidean algorithm | 18 |
| "Modular arithmetic" | Properties of mod under +, *, - | 18 |
