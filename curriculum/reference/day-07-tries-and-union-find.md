# Day 7 — Tries & Union-Find: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Why It's Useful | Time |
|---|----------|----------------|------|
| 1 | [VisuAlgo — Trie](https://visualgo.net/en/trie) | Interactive insert/search/delete animations. Step through character-by-character node creation and see the shared-prefix structure emerge. | 10 min |
| 2 | [VisuAlgo — Union-Find](https://visualgo.net/en/ufds) | Animated path compression and union by rank. Watch the forest flatten in real time as Find operations compress paths. | 10 min |
| 3 | [Union-Find in 5 Minutes (YouTube, WilliamFiset)](https://www.youtube.com/watch?v=ibjEGG7ylHk) | Clear, fast walkthrough of path compression + union by rank with visual tree diagrams. Good for locking in the mental model. | 6 min |
| 4 | [The Inverse Ackermann Function (Jeff Erickson's Algorithms)](https://jeffe.cs.illinois.edu/teaching/algorithms/book/Algorithms-JeffE.pdf) | Chapter on Union-Find with a rigorous but readable explanation of why amortized cost is O(α(n)) and why α(n) ≤ 4 for any practical input. Pages on the Ackermann function itself build real intuition. | 20 min |
| 5 | [Trie Data Structure — USFCA Visualization](https://www.cs.usfca.edu/~galles/visualization/Trie.html) | Another interactive trie visualizer. Good for watching how delete works (pruning empty branches). | 5 min |
| 6 | [Tries in Go — Junmin Lee (Medium)](https://medium.com/the-andela-way/building-a-trie-in-go-5fb9f1a0cd9c) | Go-specific implementation walkthrough with idiomatic patterns (array children, rune-based variant). | 10 min |
| 7 | [CP-Algorithms — Disjoint Set Union](https://cp-algorithms.com/data_structures/disjoint_set_union.html) | Comprehensive reference covering both optimizations, the amortized analysis sketch, and practical applications (Kruskal's, offline RMQ). Includes C++ but the logic translates directly to Go. | 15 min |
| 8 | [CLRS — Chapter 21: Data Structures for Disjoint Sets](https://mitpress.mit.edu/9780262046305/introduction-to-algorithms/) | The textbook treatment of union by rank + path compression with the formal amortized proof. Sections 21.1–21.3 are all you need. | 20 min |

**Reading strategy:** Start with resources 1 and 2 for visual intuition on both structures. Then 3 for a quick Union-Find refresher video. Read 7 for the detailed Union-Find reference. Hit 6 for Go trie patterns before you start coding. Save 4 and 8 for the inverse Ackermann deep dive if time permits.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 – 12:55 | Trie (55 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 12:00 | 10 | **Review concepts (no code).** Read the Day 7 Trie section in OVERVIEW.md. Open VisuAlgo Trie. Insert the words "cat", "car", "card", "care", "caret". Watch how shared prefixes form a single path. Then search for "ca" (prefix exists, not a word) and "cart" (prefix exists partially, not stored). |
| 12:10 | 5 | Study the ASCII diagram in Section 3.1 below. Trace through the trie structure on paper. Convince yourself that StartsWith("car") is true but Search("car") depends on `isEnd`. |
| 12:15 | 20 | **Implement from scratch.** Create `trie.go`. Define `TrieNode` with `children [26]*TrieNode` and `isEnd bool`. Define `Trie` with a `root *TrieNode`. Implement `Insert(word)`, `Search(word)`, `StartsWith(prefix)` — all three share the same traversal logic. |
| 12:35 | 10 | Implement `Delete(word)` — the tricky one. Use a recursive approach: walk to the end, unmark `isEnd`, then on the way back up, remove nodes that have no children and aren't end-of-word for something else. |
| 12:45 | 5 | Implement `CountWordsWithPrefix(prefix)` — traverse to the prefix node, then DFS to count all `isEnd` nodes in the subtree. |
| 12:50 | 5 | **Test.** Write tests: insert "apple", "app", "ape", "bat". Search for each. Search for "ap" (false). StartsWith "ap" (true). Delete "app", verify "apple" still exists. CountWordsWithPrefix("ap") should return 2 after deleting "app". |

### 12:55 – 1:50 | Union-Find (55 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 12:55 | 10 | **Review concepts (no code).** Read the Day 7 Union-Find section in OVERVIEW.md. Open VisuAlgo Union-Find. Create 8 elements. Union(0,1), Union(2,3), Union(0,2). Watch the forest structure. Then call Find(3) and watch path compression flatten the tree. |
| 13:05 | 5 | Study the ASCII diagrams in Section 4.3 and 4.4 below. Trace path compression on paper: follow the path from a deep node to root, then redraw with all nodes pointing directly to root. |
| 13:10 | 15 | **Implement from scratch.** Create `unionfind.go`. Define `UnionFind` with `parent []int`, `rank []int`, `count int`. Implement `NewUnionFind(n)` (each element is its own parent, all ranks 0, count = n). Implement `Find(x)` with path compression. Implement `Union(x, y)` with union by rank. Implement `Connected(x, y)` and `Count()`. |
| 13:25 | 10 | **Test.** 5 elements. Union(0,1), Union(2,3). Verify Connected(0,1) true, Connected(0,2) false, Count() == 3. Union(1,3). Verify Connected(0,3) true, Count() == 2. Union all into one set, Count() == 1. |
| 13:35 | 10 | Implement **Number of Islands** using Union-Find (alternative to DFS). Grid of '1' and '0'. Initialize UF for all land cells. For each land cell, union with its right and bottom neighbors if they're also land. The answer is the final count of sets (minus water cells). |
| 13:45 | 5 | Implement **Redundant Connection** — given edges of a graph that would be a tree plus one extra edge, find the extra edge. Process edges in order; the first edge where both endpoints are already connected (via Find) is the answer. |

### 1:50 – 2:00 | Recap (10 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 1:50 | 3 | Close all files. Write down from memory: Trie Insert/Search/StartsWith complexity. Union-Find Find/Union complexity with both optimizations. |
| 1:53 | 3 | Write down: Why is trie Delete trickier than Insert? What is path compression in one sentence? What does α(n) mean and why don't we care about its exact value? |
| 1:56 | 2 | Write down one gotcha you hit during each implementation. |
| 1:58 | 2 | Write down: When would you choose Union-Find over BFS/DFS? When would you choose a trie over a hash map? |

---

## 3. Core Concepts Deep Dive — Trie

### 3.1 Why Tries Beat Hash Maps for Prefix Operations

A hash map gives you O(1) lookup for exact keys, but it has **no concept of prefix relationships**. To find all words starting with "car", you'd have to iterate every key in the map — O(n) where n is the number of stored words.

A trie stores strings character-by-character in a tree where shared prefixes share nodes. To find all words starting with "car", you walk 3 nodes down (c → a → r) and then explore the subtree — O(L + k) where L is the prefix length and k is the number of matching results.

| Operation | Hash Map | Trie |
|-----------|----------|------|
| Exact lookup | O(L) hash + O(1) amortized | O(L) traversal |
| Prefix search (all matches) | O(n × L) — scan all keys | O(L + k) — walk prefix, then DFS |
| Autocomplete (top-k by prefix) | O(n × L) | O(L + subtree size) |
| Longest common prefix | O(n × L) | O(L) — walk until branch |
| Lexicographic ordering | O(n log n) sort | O(total chars) — DFS gives sorted order for free |

**Rule of thumb:** If the problem involves prefixes, use a trie. If it involves exact key lookup only, a hash map is simpler and faster.

### 3.2 Space Analysis: When Tries Are Efficient vs. Wasteful

Each trie node using `[26]*TrieNode` allocates space for 26 pointers regardless of how many children actually exist. On a 64-bit system, that's 26 × 8 = **208 bytes per node**.

**When tries are space-efficient:**
- Many words share long common prefixes (e.g., dictionary words, URLs, file paths)
- The alphabet is small (lowercase English = 26)
- The stored strings are long relative to the alphabet size

**When tries are wasteful:**
- Few words, short strings, or minimal prefix overlap — each character gets its own node chain with 25 empty pointers per node
- Large alphabet (Unicode) — `[65536]*TrieNode` per node is absurd
- Storing a small number of long, dissimilar strings (e.g., random UUIDs) — hash map wins easily

**Concrete comparison:** Storing 100,000 English dictionary words:
- Hash map: ~100K entries × (key bytes + value + overhead) ≈ several MB
- Trie (array children): ~500K nodes × 208 bytes ≈ ~100 MB (most pointers are nil)
- Trie (map children): ~500K nodes × ~50 bytes ≈ ~25 MB (only allocated children)

For competitive programming and interviews, the `[26]*TrieNode` approach is standard and acceptable. In production systems, you'd use a more compact representation.

### 3.3 Children Representation: array[26] vs map[rune]

```go
// Option A: Fixed-size array (lowercase English only)
type TrieNode struct {
    children [26]*TrieNode
    isEnd    bool
}
// Access: node.children[ch - 'a']

// Option B: Hash map (any character set)
type TrieNode struct {
    children map[rune]*TrieNode
    isEnd    bool
}
// Access: node.children[ch]
```

| Aspect | `[26]*TrieNode` | `map[rune]*TrieNode` |
|--------|-----------------|---------------------|
| **Lookup speed** | O(1) — direct index | O(1) amortized — hash lookup |
| **Constant factor** | Very fast (single array offset) | Slower (hash, bucket lookup, possible collision) |
| **Memory per node** | 208 bytes (26 × 8) fixed | ~8 bytes per actual child + map overhead |
| **Sparse nodes** | Wastes memory (25 nil pointers if 1 child) | Efficient — only stores existing children |
| **Dense nodes** | Efficient — no per-child overhead | Map overhead exceeds array for >~10 children |
| **Character set** | Only a–z (or any fixed small alphabet) | Any rune — Unicode, mixed case, digits, etc. |
| **Cache locality** | Good — contiguous array | Poor — map buckets scattered in memory |
| **Code simplicity** | Simpler, no nil-map init needed | Must initialize map: `make(map[rune]*TrieNode)` |

**Recommendation for interviews:** Use `[26]*TrieNode` unless the problem explicitly involves characters outside lowercase English. It's faster, simpler, and what interviewers expect.

**Recommendation for production:** Use `map[rune]*TrieNode` or a compressed trie (radix tree) for memory efficiency with variable character sets.

### 3.4 Trie Variants (Conceptual Overview)

**Compressed Trie (Radix Tree / Patricia Tree)**

Instead of one node per character, a compressed trie collapses chains of single-child nodes into a single node storing a substring.

```
Standard Trie:              Compressed Trie:
      (root)                     (root)
       |                          |
       c                        "ca"
       |                        /    \
       a                     "r"     "t"
      / \                    / \
     r   t                "d"  "e"
    / \                         |
   d   e                      "t"
       |
       t

Words: car, card, care, caret, cat
```

**Benefit:** Dramatically fewer nodes when strings share long prefixes but diverge at specific points. Used in IP routing tables, HTTP routers (e.g., `httprouter` in Go), and file system paths.

**Ternary Search Trie (TST)**

Each node stores a single character and has three children: less-than, equal-to, greater-than. It combines the prefix-sharing of a trie with the space efficiency of a BST.

```
TST for "cat", "car", "cup":

         c
        /|\
       . = .
         a
        /|\
       . = .
         t
        /|\
       r = .
```

**Benefit:** Better space usage than a standard trie when the alphabet is large but the key set is sparse. Lookup is O(L × log σ) where σ is the alphabet size, vs O(L) for a standard trie — the log factor comes from the BST at each level.

**When to know these:** Interviews rarely ask you to implement compressed tries or TSTs, but they may come up in system design discussions (e.g., "how would you build an autocomplete system at scale?"). Know they exist and their tradeoffs.

### 3.5 Wildcard Search: Handling '.' Matching Any Character

Some problems (e.g., LeetCode "Design Add and Search Words Data Structure") require searching a trie where '.' matches any single character.

**Approach:** When you hit a '.', you can't follow a single child — you must explore **all non-nil children** via DFS.

```go
func (t *Trie) SearchWithWildcard(word string) bool {
    return t.dfs(t.root, word, 0)
}

func (t *Trie) dfs(node *TrieNode, word string, index int) bool {
    if node == nil {
        return false
    }
    if index == len(word) {
        return node.isEnd
    }

    ch := word[index]
    if ch == '.' {
        // Try every non-nil child
        for _, child := range node.children {
            if t.dfs(child, word, index+1) {
                return true
            }
        }
        return false
    }

    // Normal character — follow the specific child
    return t.dfs(node.children[ch-'a'], word, index+1)
}
```

**Complexity:** Worst case O(26^L) if every character is '.', but in practice it's much better because most branches are nil. Average case is closer to O(L × branching factor).

**Key insight:** Wildcard search turns a trie lookup from a linear walk into a DFS. This is the same transformation that makes regex matching on tries powerful — at each wildcard, you branch to all possibilities and prune dead ends.

---

## 4. Core Concepts Deep Dive — Union-Find

### 4.1 The Naive Approach and Why It's Slow

The simplest Union-Find uses a `parent[]` array where each element points to its parent, and the root points to itself.

**Naive Find:** Follow parent pointers until you reach the root.
**Naive Union:** Make one root the child of the other.

```
Union(0,1), Union(1,2), Union(2,3), Union(3,4):

    0           Naive Find(4): 4 → 3 → 2 → 1 → 0
    |           That's O(n) per Find!
    1
    |
    2
    |
    3
    |
    4
```

Without any optimization, the tree can degenerate into a linked list. Find becomes O(n), and a sequence of n union-find operations becomes O(n^2). This is no better than scanning an array.

### 4.2 Path Compression: Flattening the Tree During Find

Path compression modifies Find so that every node visited on the way to the root is **redirected to point directly to the root**.

```go
func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x]) // path compression
    }
    return uf.parent[x]
}
```

**That one line is the entire optimization.** After `Find(4)`, the tree transforms:

```
Before Find(4):        After Find(4):

    0                       0
    |                     / | \ \
    1                    1  2  3  4
    |
    2
    |
    3
    |
    4
```

Every node on the path now points directly to the root. Future Find operations on any of these nodes are O(1).

**Iterative version (avoids recursion stack for very deep trees):**

```go
func (uf *UnionFind) Find(x int) int {
    root := x
    for root != uf.parent[root] {
        root = uf.parent[root]
    }
    // Compress: point everything on the path to root
    for x != root {
        next := uf.parent[x]
        uf.parent[x] = root
        x = next
    }
    return root
}
```

### 4.3 Union by Rank vs. Union by Size

Both heuristics ensure we always attach the **shorter tree under the taller tree**, keeping the overall height low.

**Union by Rank:**
- `rank[x]` is an upper bound on the height of x's subtree.
- When unioning two roots, the lower-rank root becomes a child of the higher-rank root.
- If ranks are equal, pick either as the new root and **increment its rank by 1**.

```go
func (uf *UnionFind) Union(x, y int) {
    rootX, rootY := uf.Find(x), uf.Find(y)
    if rootX == rootY {
        return // already in the same set
    }
    // Attach smaller-rank tree under larger-rank tree
    if uf.rank[rootX] < uf.rank[rootY] {
        uf.parent[rootX] = rootY
    } else if uf.rank[rootX] > uf.rank[rootY] {
        uf.parent[rootY] = rootX
    } else {
        uf.parent[rootY] = rootX
        uf.rank[rootX]++
    }
    uf.count--
}
```

**Union by Size:**
- `size[x]` is the number of elements in x's subtree.
- Attach the smaller set under the larger set.
- Update the size of the new root.

```go
// Alternative: union by size
if uf.size[rootX] < uf.size[rootY] {
    uf.parent[rootX] = rootY
    uf.size[rootY] += uf.size[rootX]
} else {
    uf.parent[rootY] = rootX
    uf.size[rootX] += uf.size[rootY]
}
```

**Comparison:**

| Aspect | Union by Rank | Union by Size |
|--------|--------------|---------------|
| Storage | `rank[]` (upper bound on height) | `size[]` (exact count) |
| Guarantees | Tree height ≤ log n | Tree height ≤ log n |
| With path compression | Both give O(α(n)) amortized | Both give O(α(n)) amortized |
| Practical difference | Slightly simpler (only increments by 1 in the equal case) | `size[]` is sometimes more useful (e.g., "how many elements in this component?") |

**Recommendation:** Use rank for pure connectivity problems. Use size when you need to know component sizes.

### 4.4 The Inverse Ackermann Function α(n) — Why It's "Effectively O(1)"

With both path compression and union by rank, the amortized cost per operation is O(α(n)), where α is the **inverse Ackermann function**.

**What is the Ackermann function?** It's a function A(m, n) that grows incomprehensibly fast:

```
A(0, n) = n + 1
A(1, n) = n + 2
A(2, n) = 2n + 3
A(3, n) = 2^(n+3) - 3
A(4, n) = 2^2^2^...^2 (a tower of 2s, n+3 high) - 3
A(5, n) = ... (iterated towers, beyond human comprehension)
```

Some concrete values:
- A(4, 0) = 13
- A(4, 1) = 65533
- A(4, 2) = 2^65536 - 3 (a number with ~19,729 digits)

**The inverse:** α(n) = the smallest m such that A(m, 1) ≥ n.

- α(n) = 0 for n ≤ 2
- α(n) = 1 for n = 3
- α(n) = 2 for n ≤ 7
- α(n) = 3 for n ≤ 61
- α(n) = 4 for n ≤ 2^2^2^2^16 (a number so large it dwarfs the number of atoms in the observable universe, which is ~10^80)

**For all practical purposes, α(n) ≤ 4.** The number of atoms in the universe is roughly 10^80. You'd need n > 2^65536 to get α(n) = 5. This will never happen with real data.

**So why not just say O(1)?** Because O(α(n)) is technically not O(1) — it's a function that does grow, just unimaginably slowly. Computer scientists are precise about this distinction. But from an engineering perspective, **treat it as O(1)**.

### 4.5 Why Union-Find Beats BFS/DFS for Dynamic Connectivity

**The problem:** Given a graph where edges are added one at a time, after each addition, answer "are nodes u and v connected?"

| Approach | Add Edge | Query "Connected?" | Total for m edges + q queries |
|----------|----------|-------------------|-------------------------------|
| BFS/DFS from scratch | O(1) — just store the edge | O(V + E) — run BFS/DFS | O(q × (V + E)) |
| Rebuild components | O(V + E) — rebuild | O(1) — lookup component ID | O(m × (V + E)) |
| Union-Find | O(α(n)) | O(α(n)) | O((m + q) × α(n)) ≈ O(m + q) |

Union-Find is purpose-built for this: **incremental edge additions with connectivity queries**. Each Union and Find is nearly O(1), so processing m edges and q queries is nearly linear.

**When BFS/DFS is still better:**
- You need shortest paths (BFS gives you this; Union-Find doesn't)
- You need to enumerate all nodes in a component (Union-Find only answers "same set?")
- The graph has edge deletions (Union-Find doesn't support efficient un-union — use link-cut trees or rebuild)

---

## 5. Implementation Checklist

### Trie — Full API

```go
package trie

// TrieNode represents a single node in the trie.
type TrieNode struct {
    children [26]*TrieNode
    isEnd    bool
}

// Trie is a prefix tree for lowercase English words.
type Trie struct {
    root *TrieNode
}

// NewTrie creates an empty trie.
func NewTrie() *Trie

// Insert adds a word to the trie.                               O(L)
func (t *Trie) Insert(word string)

// Search returns true if the exact word exists in the trie.     O(L)
func (t *Trie) Search(word string) bool

// StartsWith returns true if any word in the trie has           O(L)
// the given prefix.
func (t *Trie) StartsWith(prefix string) bool

// Delete removes a word from the trie, pruning nodes            O(L)
// that are no longer needed by other words.
// Returns true if the word was found and deleted.
func (t *Trie) Delete(word string) bool

// CountWordsWithPrefix returns the number of complete words     O(L + subtree)
// in the trie that start with the given prefix.
func (t *Trie) CountWordsWithPrefix(prefix string) int
```

### Trie — Implementation Tips

1. **Insert:** Walk character by character. If `node.children[ch-'a']` is nil, create a new node. At the end, set `isEnd = true`.

2. **Search vs. StartsWith:** Identical traversal. The only difference is the return: Search checks `isEnd` on the final node; StartsWith just checks that the traversal didn't hit nil.

3. **Delete (the tricky one):** Use recursion. Walk to the end of the word. If the final node has children, just unmark `isEnd` (other words pass through it). If it has no children, delete it and propagate upward — each ancestor is deleted if it has no other children and `isEnd` is false.

```go
// Delete helper: returns true if the current node should be removed.
func (t *Trie) deleteHelper(node *TrieNode, word string, depth int) bool {
    if node == nil {
        return false // word not in trie
    }
    if depth == len(word) {
        if !node.isEnd {
            return false // word not in trie
        }
        node.isEnd = false
        // Remove this node if it has no children
        return !hasChildren(node)
    }

    idx := word[depth] - 'a'
    shouldDelete := t.deleteHelper(node.children[idx], word, depth+1)
    if shouldDelete {
        node.children[idx] = nil
        // Remove this node too if it has no other children and isn't end-of-word
        return !node.isEnd && !hasChildren(node)
    }
    return false
}

func hasChildren(node *TrieNode) bool {
    for _, child := range node.children {
        if child != nil {
            return true
        }
    }
    return false
}
```

4. **CountWordsWithPrefix:** Traverse to the prefix's last node, then DFS the subtree counting `isEnd` nodes.

```go
func (t *Trie) CountWordsWithPrefix(prefix string) int {
    node := t.root
    for i := 0; i < len(prefix); i++ {
        idx := prefix[i] - 'a'
        if node.children[idx] == nil {
            return 0
        }
        node = node.children[idx]
    }
    return countWords(node)
}

func countWords(node *TrieNode) int {
    if node == nil {
        return 0
    }
    count := 0
    if node.isEnd {
        count = 1
    }
    for _, child := range node.children {
        count += countWords(child)
    }
    return count
}
```

### Union-Find — Full API

```go
package unionfind

// UnionFind tracks disjoint sets with path compression and union by rank.
type UnionFind struct {
    parent []int
    rank   []int
    count  int // number of disjoint sets
}

// NewUnionFind creates a Union-Find structure for n elements (0 to n-1).
// Each element starts as its own set.
func NewUnionFind(n int) *UnionFind

// Find returns the root representative of x's set.             O(α(n))
// Applies path compression.
func (uf *UnionFind) Find(x int) int

// Union merges the sets containing x and y.                    O(α(n))
// Uses union by rank. Returns false if x and y
// were already in the same set.
func (uf *UnionFind) Union(x, y int) bool

// Connected returns true if x and y are in the same set.       O(α(n))
func (uf *UnionFind) Connected(x, y int) bool

// Count returns the number of disjoint sets.                   O(1)
func (uf *UnionFind) Count() int
```

### Union-Find — Implementation Tips

1. **NewUnionFind:** Initialize `parent[i] = i` for all i. All ranks to 0. Count = n.

```go
func NewUnionFind(n int) *UnionFind {
    parent := make([]int, n)
    rank := make([]int, n)
    for i := range parent {
        parent[i] = i
    }
    return &UnionFind{parent: parent, rank: rank, count: n}
}
```

2. **Find with path compression** — the one-liner recursive version:

```go
func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x])
    }
    return uf.parent[x]
}
```

3. **Union with rank** — always check roots first:

```go
func (uf *UnionFind) Union(x, y int) bool {
    rootX, rootY := uf.Find(x), uf.Find(y)
    if rootX == rootY {
        return false // already in the same set
    }
    if uf.rank[rootX] < uf.rank[rootY] {
        uf.parent[rootX] = rootY
    } else if uf.rank[rootX] > uf.rank[rootY] {
        uf.parent[rootY] = rootX
    } else {
        uf.parent[rootY] = rootX
        uf.rank[rootX]++
    }
    uf.count--
    return true
}
```

4. **Connected and Count** are trivial:

```go
func (uf *UnionFind) Connected(x, y int) bool {
    return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFind) Count() int {
    return uf.count
}
```

### Test Plan

```go
func TestTrie(t *testing.T) {
    // Basic insert and search
    // - Insert "apple", Search("apple") → true, Search("app") → false
    // - Insert "app", Search("app") → true

    // StartsWith
    // - StartsWith("app") → true, StartsWith("apl") → false
    // - StartsWith("apple") → true (exact word is a prefix of itself)
    // - StartsWith("applx") → false

    // Delete
    // - Delete "app", Search("app") → false, Search("apple") → true
    //   (shared prefix nodes must survive)
    // - Delete "apple", Search("apple") → false, StartsWith("app") → false
    //   (orphaned nodes pruned)
    // - Delete a word that doesn't exist → returns false, trie unchanged

    // CountWordsWithPrefix
    // - Insert "car", "card", "care", "caret", "cat"
    // - CountWordsWithPrefix("car") → 3 (car, card, care, caret... wait:
    //   car, card, care — caret starts with "care" which starts with "car")
    //   Actually: car, card, care, caret → 4
    // - CountWordsWithPrefix("care") → 2 (care, caret)
    // - CountWordsWithPrefix("cat") → 1
    // - CountWordsWithPrefix("z") → 0

    // Edge cases
    // - Empty string: Insert(""), Search("") → depends on design decision
    // - Single character: Insert("a"), Search("a") → true, Search("b") → false
    // - All same characters: Insert("aaa"), Insert("aa"), Insert("a")
}

func TestUnionFind(t *testing.T) {
    // Basic union and find
    // - 5 elements: Union(0,1), Find(0) == Find(1), Count() == 4
    // - Union(2,3), Connected(0,2) → false, Count() == 3
    // - Union(1,3), Connected(0,3) → true, Count() == 2

    // Idempotent union
    // - Union(0,1) again → returns false, Count unchanged

    // All elements in one set
    // - Union everything, Count() == 1

    // Path compression verification
    // - Create a long chain: Union(0,1), Union(1,2), Union(2,3), Union(3,4)
    // - Find(4): parent[4] should now point directly to root

    // Edge cases
    // - Single element: NewUnionFind(1), Find(0) == 0, Count() == 1
    // - Self-union: Union(0,0) → returns false, Count unchanged
    // - Large n: NewUnionFind(100000), union random pairs, verify consistency
}
```

---

## 6. Application Patterns

### 6.1 Trie Patterns

#### Autocomplete System

Given a list of words and their frequencies, support `input(char)` that returns the top 3 words with the current prefix, ordered by frequency.

**Approach:** Build a trie. At each node, maintain a sorted list of the top 3 (word, frequency) pairs in that node's subtree. On insert, propagate frequency info upward. On query, traverse to the prefix node and return its top-3 list.

```go
type AutocompleteNode struct {
    children [26]*AutocompleteNode
    top3     []WordFreq // maintained during insert
    isEnd    bool
    freq     int
}

type WordFreq struct {
    Word string
    Freq int
}
```

#### Word Search II (Trie + Backtracking)

**Problem:** Given an m×n board of characters and a list of words, find all words that can be formed by adjacent cells (no cell reused per word).

**Approach:** Build a trie from the word list. DFS/backtrack on the board, but instead of checking each word separately (slow), walk the trie simultaneously. At each board cell, follow the corresponding trie edge. If you reach an `isEnd` node, you found a word. Prune branches where the trie has no matching child.

```go
func FindWords(board [][]byte, words []string) []string {
    trie := NewTrie()
    for _, w := range words {
        trie.Insert(w)
    }

    var result []string
    rows, cols := len(board), len(board[0])

    var dfs func(r, c int, node *TrieNode, path []byte)
    dfs = func(r, c int, node *TrieNode, path []byte) {
        if r < 0 || r >= rows || c < 0 || c >= cols {
            return
        }
        ch := board[r][c]
        if ch == '#' || node.children[ch-'a'] == nil {
            return // visited or no trie edge
        }

        node = node.children[ch-'a']
        path = append(path, ch)

        if node.isEnd {
            result = append(result, string(path))
            node.isEnd = false // avoid duplicate results
        }

        board[r][c] = '#' // mark visited
        dfs(r-1, c, node, path)
        dfs(r+1, c, node, path)
        dfs(r, c-1, node, path)
        dfs(r, c+1, node, path)
        board[r][c] = ch // restore
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            dfs(r, c, trie.root, nil)
        }
    }
    return result
}
```

**Optimization:** After finding a word, prune the trie node if it has no children. This prevents re-exploring dead branches for the rest of the board scan.

#### Longest Common Prefix

**Approach:** Insert all strings into a trie. Starting from the root, walk down as long as each node has exactly one child and `isEnd` is false. The path traversed is the longest common prefix.

```go
func LongestCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }
    trie := NewTrie()
    for _, s := range strs {
        trie.Insert(s)
    }

    var prefix []byte
    node := trie.root
    for {
        // Count non-nil children
        childCount := 0
        nextIdx := -1
        for i, child := range node.children {
            if child != nil {
                childCount++
                nextIdx = i
            }
        }
        if childCount != 1 || node.isEnd {
            break // branch point or end of a word
        }
        prefix = append(prefix, byte(nextIdx)+'a')
        node = node.children[nextIdx]
    }
    return string(prefix)
}
```

### 6.2 Union-Find Patterns

#### Number of Islands (Union-Find Alternative to DFS)

```go
func NumIslands(grid [][]byte) int {
    if len(grid) == 0 {
        return 0
    }
    rows, cols := len(grid), len(grid[0])
    uf := NewUnionFind(rows * cols)

    waterCount := 0
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == '0' {
                waterCount++
                continue
            }
            // Union with right neighbor
            if c+1 < cols && grid[r][c+1] == '1' {
                uf.Union(r*cols+c, r*cols+c+1)
            }
            // Union with bottom neighbor
            if r+1 < rows && grid[r+1][c] == '1' {
                uf.Union(r*cols+c, (r+1)*cols+c)
            }
        }
    }
    return uf.Count() - waterCount
}
```

**Key trick:** Initialize UF for all cells, then subtract water cells from the count. Alternatively, only initialize UF for land cells using a mapping.

#### Accounts Merge

**Problem:** Given a list of accounts where each account is [name, email1, email2, ...], merge accounts that share any email.

**Approach:** Each email maps to a Union-Find index. For each account, union all emails in that account together. Then group emails by their root and reconstruct the merged accounts.

```go
func AccountsMerge(accounts [][]string) [][]string {
    emailToID := make(map[string]int)
    emailToName := make(map[string]string)
    id := 0

    // Assign an ID to each unique email
    for _, account := range accounts {
        name := account[0]
        for _, email := range account[1:] {
            if _, exists := emailToID[email]; !exists {
                emailToID[email] = id
                id++
            }
            emailToName[email] = name
        }
    }

    uf := NewUnionFind(id)

    // Union all emails within the same account
    for _, account := range accounts {
        firstID := emailToID[account[1]]
        for _, email := range account[2:] {
            uf.Union(firstID, emailToID[email])
        }
    }

    // Group emails by root
    rootToEmails := make(map[int][]string)
    for email, emailID := range emailToID {
        root := uf.Find(emailID)
        rootToEmails[root] = append(rootToEmails[root], email)
    }

    // Build result
    var result [][]string
    for _, emails := range rootToEmails {
        sort.Strings(emails)
        name := emailToName[emails[0]]
        result = append(result, append([]string{name}, emails...))
    }
    return result
}
```

#### Redundant Connection (Cycle Detection)

**Problem:** Given a tree with one extra edge (n nodes, n edges), find the edge that, if removed, makes the graph a valid tree.

```go
func FindRedundantConnection(edges [][]int) []int {
    n := len(edges)
    uf := NewUnionFind(n + 1) // nodes are 1-indexed

    for _, edge := range edges {
        if !uf.Union(edge[0], edge[1]) {
            return edge // already connected — this edge creates a cycle
        }
    }
    return nil // shouldn't reach here if input is valid
}
```

**Why this works:** A tree with n nodes has exactly n-1 edges. The input has n edges. Process them in order; the first edge connecting two already-connected nodes is the redundant one.

#### Kruskal's MST

**Problem:** Find the minimum spanning tree of a weighted graph.

```go
func KruskalMST(n int, edges [][]int) (int, [][]int) {
    // edges[i] = [u, v, weight]
    sort.Slice(edges, func(i, j int) bool {
        return edges[i][2] < edges[j][2]
    })

    uf := NewUnionFind(n)
    totalWeight := 0
    var mstEdges [][]int

    for _, edge := range edges {
        u, v, w := edge[0], edge[1], edge[2]
        if uf.Union(u, v) {
            totalWeight += w
            mstEdges = append(mstEdges, edge)
            if len(mstEdges) == n-1 {
                break // MST complete
            }
        }
    }
    return totalWeight, mstEdges
}
```

**Why Union-Find is perfect here:** Kruskal's processes edges in weight order, adding each edge if it doesn't create a cycle. "Does this edge create a cycle?" is exactly "are these two nodes already connected?" — the core Union-Find query.

---

## 7. Visual Diagrams

### 7.1 Trie Storing Words with Shared Prefixes

Words: **"cat", "car", "card", "care", "caret", "do", "dog"**

```
                         (root)
                        /      \
                       c        d
                       |        |
                       a        o ·
                      / \       |
                     t·  r·     g ·
                        / \
                       d·  e·
                           |
                           t·

Legend:  · = isEnd (a complete word ends here)

Paths:
  root → c → a → t·             = "cat"
  root → c → a → r·             = "car"
  root → c → a → r → d·         = "card"
  root → c → a → r → e·         = "care"
  root → c → a → r → e → t·     = "caret"
  root → d → o·                  = "do"
  root → d → o → g·             = "dog"

Shared prefixes:
  "ca"   — shared by cat, car, card, care, caret (5 words, 2 nodes)
  "car"  — shared by car, card, care, caret (4 words, 3 nodes)
  "care" — shared by care, caret (2 words, 4 nodes)
  "do"   — shared by do, dog (2 words, 2 nodes)
```

### 7.2 Union-Find Forest: Before and After Path Compression

**Setup:** 8 elements (0-7). Operations performed:
Union(0,1), Union(1,2), Union(3,4), Union(4,5), Union(2,5), Union(6,7)

```
After unions (before any extra Find calls):

    parent: [0, 0, 0, 3, 3, 0, 6, 6]
    rank:   [2, 0, 0, 1, 0, 0, 1, 0]

    Set 1 (root=0):             Set 2 (root=6):

          0                          6
        / | \                        |
       1  2  5                       7
             |
             3
             |
             4

    Count = 2


Now call Find(4):

    Path: 4 → 3 → 5 → 0    (following parent pointers to root)

    Path compression redirects every node on the path to point to 0:

    BEFORE Find(4):              AFTER Find(4):

          0                             0
        / | \                      / / | \ \
       1  2  5                    1 2  3  4  5
             |
             3                      6
             |                      |
             4                      7

    parent: [0, 0, 0, 0, 0, 0, 6, 6]
                         ^  ^
                      changed!

    Future Find(3) and Find(4) are now O(1) — direct parent is root.
```

### 7.3 Union by Rank: Step-by-Step Example

```
Start: 6 elements, each in its own set.

    parent: [0, 1, 2, 3, 4, 5]
    rank:   [0, 0, 0, 0, 0, 0]

    0   1   2   3   4   5          count = 6


─── Union(0, 1): rank[0] == rank[1] == 0 → pick 0 as root, rank[0]++

    parent: [0, 0, 2, 3, 4, 5]
    rank:   [1, 0, 0, 0, 0, 0]

      0    2   3   4   5           count = 5
      |
      1


─── Union(2, 3): rank[2] == rank[3] == 0 → pick 2 as root, rank[2]++

    parent: [0, 0, 2, 2, 4, 5]
    rank:   [1, 0, 1, 0, 0, 0]

      0    2    4   5              count = 4
      |    |
      1    3


─── Union(4, 5): rank[4] == rank[5] == 0 → pick 4 as root, rank[4]++

    parent: [0, 0, 2, 2, 4, 4]
    rank:   [1, 0, 1, 0, 1, 0]

      0    2    4                  count = 3
      |    |    |
      1    3    5


─── Union(0, 2): rank[0] == rank[2] == 1 → pick 0 as root, rank[0]++

    parent: [0, 0, 0, 2, 4, 4]
    rank:   [2, 0, 1, 0, 1, 0]

        0       4                  count = 2
       / \      |
      1   2     5
          |
          3


─── Union(0, 4): rank[0] = 2 > rank[4] = 1 → attach 4 under 0

    parent: [0, 0, 0, 2, 0, 4]
    rank:   [2, 0, 1, 0, 1, 0]

           0                       count = 1
         / | \
        1  2   4
           |   |
           3   5

    Height = 2, with 6 elements.
    Without union by rank, height could be 5 (a chain).
    Union by rank guarantees height ≤ log₂(n) = log₂(6) ≈ 2.58 ✓
```

---

## 8. Self-Assessment

Answer these from memory after your session. If you can't, that's tomorrow's priority.

### Q1: What's the time complexity of deleting a word from a trie, and why is it trickier than insert?

<details>
<summary>Answer</summary>

**Time complexity:** O(L) where L is the word length — same as insert.

**Why it's trickier:** Insert only creates nodes; it never needs to look backward. Delete must decide for each node on the path whether to remove it. A node should only be removed if:
1. It's not the end of another word (`isEnd == false`), AND
2. It has no children (no other word passes through it).

This requires bottom-up decision-making — you walk down to the end of the word, then work back up, removing nodes that are no longer needed. A recursive approach handles this naturally: the return value signals "I should be deleted" up the call stack. If you just unmark `isEnd` without pruning, you waste memory. If you prune too aggressively, you delete nodes that other words depend on.

Example: Trie contains "apple" and "app". Deleting "apple" must remove nodes for 'l' and 'e', but leave the 'a', 'p', 'p' nodes intact (they belong to "app"). Deleting "app" only unmarks `isEnd` on the second 'p' — no nodes are removed because "apple" still needs them.
</details>

### Q2: Why is path compression alone not enough — why do you also need union by rank?

<details>
<summary>Answer</summary>

Path compression only kicks in during **Find** operations. If you perform many Union operations before any Find, the trees can still become deep.

Consider Union(0,1), Union(1,2), Union(2,3), ..., Union(n-2, n-1) without any Find calls in between, and without union by rank. Each Union attaches the new element under the previous root. The result is a chain of depth n-1. The first Find operation after this will cost O(n).

With path compression alone, the amortized cost per operation is O(log n) — better than O(n) but not nearly-constant.

With **both** optimizations:
- Union by rank keeps trees shallow (height ≤ log n) regardless of Find frequency
- Path compression flattens whatever depth remains during Find calls
- Together, they achieve O(α(n)) amortized per operation — effectively O(1)

The two optimizations are complementary: rank prevents the tree from getting deep in the first place; compression flattens any remaining depth retroactively.
</details>

### Q3: You insert "app", "apple", and "application" into a trie. How many nodes exist (excluding root)? How many have isEnd = true?

<details>
<summary>Answer</summary>

Trace the paths:
```
root → a → p → p(·) → l → e(·) → — (this is "apple" done)
                            ↓ (no — "application" continues from 'l')

Actually, let's be precise:
  "app"         → a, p, p·
  "apple"       → a, p, p, l, e·
  "application" → a, p, p, l, i, c, a, t, i, o, n·
```

Shared prefix: "app" is shared by all three (3 nodes: a, p, p). Then "appl" is shared by "apple" and "application" (1 more node: l). Then they diverge:
- "apple" adds: e (1 node)
- "application" adds: i, c, a, t, i, o, n (7 nodes)

**Total nodes (excluding root):** 3 + 1 + 1 + 7 = **12 nodes**

**Nodes with isEnd = true:** 3 — the second 'p' (for "app"), 'e' (for "apple"), and 'n' (for "application").
</details>

### Q4: In a Union-Find with n = 10^18 elements (hypothetically), what is α(n)? Would using Union-Find still be practical?

<details>
<summary>Answer</summary>

α(10^18) = **4**.

Since α(n) = 4 for any n up to approximately 2^(2^(2^(2^16))), and 10^18 is astronomically smaller than that threshold, α(n) = 4.

Would it be practical? The Union-Find **data structure** would be fine — each operation is essentially O(1). The bottleneck would be **memory**: storing the `parent[]` and `rank[]` arrays for 10^18 elements would require ~16 × 10^18 bytes = 16 exabytes. That's the real constraint, not the algorithm's time complexity.

This is a good reminder: α(n) is always ≤ 4 for real-world inputs. The inverse Ackermann function is a theoretical curiosity, not a practical concern.
</details>

### Q5: You're building an autocomplete system. A user types a prefix and you need to return the top 5 matching words by frequency. What data structure would you use and why?

<details>
<summary>Answer</summary>

**A trie with frequency-augmented nodes.**

Store at each node a sorted list (or min-heap of size 5) of the highest-frequency words in its subtree. When the user types a prefix, traverse to the prefix node in O(L) time and return the pre-computed top-5 list in O(1).

**Why not a hash map?** To find all words with a given prefix, you'd scan every key — O(n × L). With a trie, you walk directly to the prefix in O(L).

**Why not a trie + DFS at query time?** DFS through the entire subtree at query time could be O(total characters in all matching words) — too slow if thousands of words share the prefix. Pre-computing top-k at each node trades insert-time work for O(1) query-time retrieval.

**Tradeoff:** Maintaining the top-5 at each node adds O(L × 5) work per insert (propagate up the path). This is acceptable because inserts are infrequent (dictionary updates) while queries are very frequent (every keystroke).

Alternative for production: a **ternary search trie** or **compressed trie** for memory efficiency, possibly combined with a separate frequency index.
</details>

---

## Complexity Reference (Quick Glance)

| Structure | Operation | Time | Notes |
|-----------|-----------|------|-------|
| Trie | Insert | O(L) | L = word length |
| Trie | Search | O(L) | Returns false if prefix exists but isn't a complete word |
| Trie | StartsWith | O(L) | Returns true if the path exists (regardless of isEnd) |
| Trie | Delete | O(L) | Must prune empty branches bottom-up |
| Trie | CountWordsWithPrefix | O(L + S) | S = subtree size |
| Trie | Space | O(N × L) | N = number of words; sharing prefixes reduces this in practice |
| Union-Find | Find | O(α(n)) ≈ O(1) | With path compression |
| Union-Find | Union | O(α(n)) ≈ O(1) | With union by rank/size |
| Union-Find | Connected | O(α(n)) ≈ O(1) | Two Find calls |
| Union-Find | Count | O(1) | Maintained during Union |
| Union-Find | Space | O(n) | parent[] + rank[] arrays |
