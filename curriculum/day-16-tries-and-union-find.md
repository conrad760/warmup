# Day 16: Tries & Union-Find

> **Format:** 2-hour refresher | **Language:** Go | **Focus:** Interview patterns and decision-making

---

## Pattern Catalog

### Pattern 1: Prefix Search / Autocomplete (Trie)

**Trigger:** "Implement a data structure that supports insert, search, and prefix lookup." Any problem mentioning prefix matching, autocomplete, or dictionary lookups.

**Go Template:**

```go
type TrieNode struct {
    children [26]*TrieNode
    isEnd    bool
}

type Trie struct {
    root *TrieNode
}

func NewTrie() *Trie {
    return &Trie{root: &TrieNode{}}
}

func (t *Trie) Insert(word string) {
    cur := t.root
    for _, ch := range word {
        idx := ch - 'a'
        if cur.children[idx] == nil {
            cur.children[idx] = &TrieNode{}
        }
        cur = cur.children[idx]
    }
    cur.isEnd = true
}

func (t *Trie) Search(word string) bool {
    node := t.find(word)
    return node != nil && node.isEnd
}

func (t *Trie) StartsWith(prefix string) bool {
    return t.find(prefix) != nil
}

// find returns the node at the end of the path, or nil.
func (t *Trie) find(s string) *TrieNode {
    cur := t.root
    for _, ch := range s {
        idx := ch - 'a'
        if cur.children[idx] == nil {
            return nil
        }
        cur = cur.children[idx]
    }
    return cur
}
```

**Complexity:** O(L) insert, search, startsWith where L = length of the word/prefix. Space O(N*L) worst case for N words.

**Watch out:**
- `Search` checks `isEnd`. `StartsWith` does NOT. Mixing these up is the #1 trie bug in interviews.
- If the character set isn't lowercase a-z, use a `map[rune]*TrieNode` instead of `[26]*TrieNode`.

---

### Pattern 2: Word Search II (Trie + Backtracking on Grid)

**Trigger:** "Given a board of characters and a list of words, find all words that can be formed by adjacent cells." Any problem combining a dictionary with grid traversal.

**Go Template:**

```go
func findWords(board [][]byte, words []string) []string {
    trie := NewTrie()
    for _, w := range words {
        trie.Insert(w)
    }

    rows, cols := len(board), len(board[0])
    var result []string
    dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

    var backtrack func(r, c int, node *TrieNode, path []byte)
    backtrack = func(r, c int, node *TrieNode, path []byte) {
        ch := board[r][c]
        idx := ch - 'a'
        next := node.children[idx]
        if next == nil {
            return
        }

        path = append(path, ch)

        if next.isEnd {
            result = append(result, string(path))
            next.isEnd = false // de-duplicate: mark found word
        }

        board[r][c] = '#' // mark visited
        for _, d := range dirs {
            nr, nc := r+d[0], c+d[1]
            if nr >= 0 && nr < rows && nc >= 0 && nc < cols && board[nr][nc] != '#' {
                backtrack(nr, nc, next, path)
            }
        }
        board[r][c] = ch // restore

        // Optimization: prune dead-end trie branches
        if isLeafNode(next) {
            node.children[idx] = nil
        }
    }

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            backtrack(r, c, trie.root, nil)
        }
    }
    return result
}

func isLeafNode(node *TrieNode) bool {
    for _, child := range node.children {
        if child != nil {
            return false
        }
    }
    return true
}
```

**Complexity:** O(M * N * 4^L) worst case where M*N is the grid size and L is max word length. The trie prunes most branches in practice.

**Watch out:**
- Prune trie nodes after finding words. Interviewers specifically look for this optimization.
- Set `isEnd = false` after collecting a word to avoid duplicates in the result.
- Restore the board cell after backtracking (the `board[r][c] = ch` line).

---

### Pattern 3: Longest Common Prefix via Trie

**Trigger:** "Find the longest common prefix among a set of strings." Also useful for problems asking about shared prefixes across a dictionary.

**Go Template:**

```go
func longestCommonPrefix(strs []string) string {
    if len(strs) == 0 {
        return ""
    }
    trie := NewTrie()
    for _, s := range strs {
        if s == "" {
            return "" // empty string means LCP is ""
        }
        trie.Insert(s)
    }

    // Walk down the trie until a node has != 1 child or is an end-of-word.
    var prefix []byte
    cur := trie.root
    for {
        count, nextIdx := 0, -1
        for i, child := range cur.children {
            if child != nil {
                count++
                nextIdx = i
            }
        }
        // Stop if branch (>1 child), dead end (0 children), or a word ends here.
        if count != 1 || cur.isEnd {
            break
        }
        prefix = append(prefix, byte(nextIdx)+'a')
        cur = cur.children[nextIdx]
    }
    return string(prefix)
}
```

**Complexity:** O(S) where S = sum of all character lengths. Space O(S).

**Watch out:**
- Must stop at `isEnd` nodes. If one word is `"app"` and another is `"apple"`, the LCP is `"app"` -- you stop because `"app"` ends there even though the branch continues.
- An empty string in the input means the LCP is `""`.

---

### Pattern 4: Connected Components (Union-Find)

**Trigger:** "How many connected components?" "Number of islands." "Friend circles / provinces." Any problem where you merge groups and count distinct groups.

**Go Template:**

```go
type UnionFind struct {
    parent []int
    rank   []int
    count  int // number of distinct components
}

func NewUnionFind(n int) *UnionFind {
    uf := &UnionFind{
        parent: make([]int, n),
        rank:   make([]int, n),
        count:  n,
    }
    for i := range uf.parent {
        uf.parent[i] = i
    }
    return uf
}

func (uf *UnionFind) Find(x int) int {
    if uf.parent[x] != x {
        uf.parent[x] = uf.Find(uf.parent[x]) // path compression
    }
    return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) bool {
    rx, ry := uf.Find(x), uf.Find(y)
    if rx == ry {
        return false // already in the same set
    }
    // union by rank
    if uf.rank[rx] < uf.rank[ry] {
        rx, ry = ry, rx
    }
    uf.parent[ry] = rx
    if uf.rank[rx] == uf.rank[ry] {
        uf.rank[rx]++
    }
    uf.count--
    return true
}
```

**Number of Islands using Union-Find:**

```go
func numIslands(grid [][]byte) int {
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
            // Union with right and down neighbors only (avoids double work)
            if r+1 < rows && grid[r+1][c] == '1' {
                uf.Union(r*cols+c, (r+1)*cols+c)
            }
            if c+1 < cols && grid[r][c+1] == '1' {
                uf.Union(r*cols+c, r*cols+c+1)
            }
        }
    }
    return uf.count - waterCount
}
```

**Complexity:** O(M * N * alpha(M*N)) which is effectively O(M*N). Space O(M*N).

**Watch out:**
- Subtract `waterCount` from `uf.count` at the end. Every cell starts as its own component, including water cells.
- Only union with right and down neighbors to avoid redundant union calls.

---

### Pattern 5: Cycle Detection in Undirected Graph

**Trigger:** "Does adding this edge create a cycle?" "Find the redundant connection." Any undirected graph where you process edges one at a time and need to detect when a cycle forms.

**Go Template:**

```go
// Returns the edge that, when added, creates a cycle.
// Edges are 1-indexed: [[1,2],[1,3],[2,3]]
func findRedundantConnection(edges [][]int) []int {
    n := len(edges)
    uf := NewUnionFind(n + 1) // 1-indexed nodes

    for _, e := range edges {
        if !uf.Union(e[0], e[1]) {
            return e // find(u) == find(v) already, so this edge creates a cycle
        }
    }
    return nil // no redundant edge
}
```

**Complexity:** O(N * alpha(N)) ~ O(N). Space O(N).

**Watch out:**
- The cycle edge is the first edge where `Union` returns `false` (both endpoints already share a root).
- Node indexing: many problems use 1-indexed nodes. Size your parent array accordingly (`n+1`).

---

### Pattern 6: Dynamic Connectivity / Equivalence Grouping

**Trigger:** "Merge accounts that share an email." "Group items by some equivalence relation." Problems where relationships arrive incrementally and you need to query or merge groups.

**Go Template (Accounts Merge):**

```go
func accountsMerge(accounts [][]string) [][]string {
    // Map each email to a unique id
    emailToID := map[string]int{}
    emailToName := map[string]string{}
    id := 0

    for _, acc := range accounts {
        name := acc[0]
        for _, email := range acc[1:] {
            if _, ok := emailToID[email]; !ok {
                emailToID[email] = id
                id++
            }
            emailToName[email] = name
        }
    }

    uf := NewUnionFind(id)

    // Union all emails within the same account
    for _, acc := range accounts {
        firstID := emailToID[acc[1]]
        for _, email := range acc[2:] {
            uf.Union(firstID, emailToID[email])
        }
    }

    // Group emails by root
    groups := map[int][]string{}
    for email, eid := range emailToID {
        root := uf.Find(eid)
        groups[root] = append(groups[root], email)
    }

    // Build result
    var result [][]string
    for _, emails := range groups {
        sort.Strings(emails)
        name := emailToName[emails[0]]
        result = append(result, append([]string{name}, emails...))
    }
    return result
}
```

**Complexity:** O(N * K * alpha(NK)) for N accounts with K emails each, plus O(NK log(NK)) for sorting. Space O(NK).

**Watch out:**
- Map emails to integer IDs first, then use standard Union-Find on integers. Don't try to union strings directly.
- Sort emails within each group before returning (problem requirement that's easy to forget).

---

## Decision Framework

```
What does the problem ask?
|
|-- "Prefix matching", "autocomplete", "dictionary lookup"
|     --> Trie (Pattern 1)
|
|-- "Find all words from a list in a 2D grid"
|     --> Trie + Backtracking (Pattern 2)
|     (Not brute-force DFS per word -- that's O(words * grid * 4^L))
|
|-- "Longest common prefix"
|     --> Trie (Pattern 3) or horizontal scan (simpler for this specific case)
|
|-- "How many connected components?" / "Are these connected?"
|     --> Union-Find (Pattern 4) if edges arrive incrementally
|     --> BFS/DFS if you already have the full graph
|
|-- "Does this edge create a cycle?" / "Find redundant edge"
|     --> Union-Find cycle detection (Pattern 5)
|
|-- "Group items by equivalence" / "Merge accounts"
|     --> Union-Find (Pattern 6)
```

**When to pick Union-Find over DFS/BFS:**
- Edges arrive one at a time (online/incremental).
- You need to repeatedly query "are X and Y connected?" as unions happen.
- You need to count components as merges happen.

**When DFS/BFS is simpler:**
- You have the full graph upfront and just need one traversal.
- You need shortest paths (BFS) -- Union-Find can't do this.

---

## Common Interview Traps

### Trie Traps

| Trap | What goes wrong | Fix |
|------|----------------|-----|
| `Search` vs `StartsWith` confusion | `Search("app")` returns true for `"apple"` | `Search` must check `node.isEnd == true` |
| Word Search II without trie pruning | TLE on large inputs | Remove trie leaf nodes after finding words |
| Word Search II duplicates | Same word found multiple times | Set `isEnd = false` after collecting the word |
| Hardcoded `[26]` array | Crashes on uppercase or digits | Use `map[rune]*TrieNode` if charset varies |

### Union-Find Traps

| Trap | What goes wrong | Fix |
|------|----------------|-----|
| No path compression | O(n) per Find instead of O(alpha(n)) | Add `uf.parent[x] = uf.Find(uf.parent[x])` |
| No union by rank | Tree degenerates to linked list | Always attach shorter tree under taller root |
| Decrement count when roots match | Component count goes negative | Only decrement in `Union` when `rx != ry` |
| Off-by-one on node indexing | Wrong unions, panic | Check if nodes are 0-indexed or 1-indexed |
| Number of Islands water cells | Count is too high | Track and subtract water cells from total |

---

## Thought Process Walkthrough

### Walkthrough 1: Implement Trie (LeetCode 208)

> "Implement a trie with insert, search, and startsWith methods."

**Step 1 -- Clarify.**
- Character set? Lowercase a-z. So `[26]*TrieNode` is fine.
- What should `search("app")` return if only `"apple"` was inserted? `false`.
- What should `startsWith("app")` return? `true`.

**Step 2 -- Design the node.**
```
TrieNode:
    children [26]*TrieNode
    isEnd    bool
```
No need to store the character itself -- the index in the parent's `children` array encodes it.

**Step 3 -- Insert.**
Walk character by character. Create nodes as needed. Mark the last node `isEnd = true`.

**Step 4 -- Search.**
Walk character by character. If any child is nil, return `false`. At the end, return `node.isEnd`.

**Step 5 -- StartsWith.**
Same as Search but return `true` as long as we reach the end of the prefix (don't check `isEnd`).

**Step 6 -- Test mentally.**
- Insert `"apple"`, search `"apple"` -> true. Search `"app"` -> false (isEnd is false at 'p'). StartsWith `"app"` -> true.
- Insert `"app"`, search `"app"` -> true now.

**Step 7 -- Code it.** (See Pattern 1 template above.)

**Complexity check:** All operations O(L). Space O(total characters inserted). State this out loud.

---

### Walkthrough 2: Number of Islands via Union-Find (LeetCode 200)

> "Given a 2D grid of '1's (land) and '0's (water), count the number of islands."

**Step 1 -- Recognize the pattern.**
Connected components problem. DFS works, but the interviewer asks for a Union-Find approach.

**Step 2 -- Map 2D to 1D.**
Cell `(r, c)` maps to index `r * cols + c`. This gives each cell a unique integer ID for Union-Find.

**Step 3 -- Initialize.**
Create a Union-Find of size `rows * cols`. Every cell starts as its own component. Track how many cells are water.

**Step 4 -- Process the grid.**
For each land cell `(r, c)`:
- If the cell below `(r+1, c)` is also land, union them.
- If the cell to the right `(r, c+1)` is also land, union them.
- Only check right and down to avoid double-processing.

**Step 5 -- Compute answer.**
`uf.count - waterCount`. Each water cell is its own component in the Union-Find, but we don't want to count those.

**Step 6 -- Trace through a small example.**
```
Grid:       1 1 0
            0 1 0
            0 0 1

Initial count: 9 (all cells are separate components)
Water cells: 5

Process (0,0): union with (0,1) -> count=8, union with (1,0) fails (water)
Process (0,1): union with (0,2) fails (water), union with (1,1) -> count=7
Process (1,1): no land neighbors right or down
Process (2,2): no neighbors to check

Answer: 7 - 5 = 2 islands. Correct.
```

**Step 7 -- Code it.** (See Pattern 4 template above.)

**Complexity check:** O(M*N * alpha(M*N)) ~ O(M*N). Space O(M*N). State this.

---

## Time Targets

| Problem | Target | Red Flag |
|---------|--------|----------|
| Implement Trie (LC 208) | 10 min | > 15 min |
| Word Search II (LC 212) | 25 min | > 35 min |
| Number of Islands UF (LC 200) | 12 min | > 18 min |
| Redundant Connection (LC 684) | 10 min | > 15 min |
| Accounts Merge (LC 721) | 20 min | > 30 min |
| Longest Common Prefix (LC 14) | 8 min | > 12 min |

---

## Quick Drill (30 min)

Do these in order. Write code from memory, not copy-paste.

1. **(5 min) Union-Find from scratch.** Write `NewUnionFind`, `Find` with path compression, `Union` with rank. This must be muscle memory.

2. **(5 min) Trie from scratch.** Write `Insert`, `Search`, `StartsWith`. Again, muscle memory.

3. **(8 min) Redundant Connection (LC 684).** Process edges, return the first one where both endpoints already share a root.

4. **(12 min) Implement Trie, then use it for prefix count.** Extend the trie so each node stores how many words pass through it. Given a prefix, return the count. (Common follow-up question.)

**Prefix-count extension hint:**

```go
type TrieNode struct {
    children  [26]*TrieNode
    isEnd     bool
    prefixCnt int // increment on every insert that passes through
}

func (t *Trie) Insert(word string) {
    cur := t.root
    for _, ch := range word {
        idx := ch - 'a'
        if cur.children[idx] == nil {
            cur.children[idx] = &TrieNode{}
        }
        cur = cur.children[idx]
        cur.prefixCnt++
    }
    cur.isEnd = true
}

func (t *Trie) CountPrefix(prefix string) int {
    cur := t.root
    for _, ch := range prefix {
        idx := ch - 'a'
        if cur.children[idx] == nil {
            return 0
        }
        cur = cur.children[idx]
    }
    return cur.prefixCnt
}
```

---

## Self-Assessment

### After the drill, check honestly:

**Trie:**
- [ ] Can I write Insert/Search/StartsWith in under 5 minutes without reference?
- [ ] Do I correctly distinguish Search (needs `isEnd`) from StartsWith (doesn't)?
- [ ] Can I explain why Word Search II uses a trie instead of searching each word independently?
- [ ] Do I know the trie pruning optimization for Word Search II and why it matters?

**Union-Find:**
- [ ] Can I write Find (with path compression) and Union (with rank) in under 3 minutes?
- [ ] Do I remember to only decrement `count` when the roots actually differ?
- [ ] Can I explain when Union-Find beats DFS/BFS and when it doesn't?
- [ ] Can I map a 2D grid problem into Union-Find without hesitation?

**Decision-making:**
- [ ] Given a new problem, can I identify whether it's a Trie problem, Union-Find problem, or neither within 60 seconds?
- [ ] Can I articulate the time complexity of Union-Find operations including why alpha(n) is effectively constant?

### If you failed any checkbox:
- Re-read the corresponding pattern section.
- Re-do the drill item without looking at the template.
- The goal is recognition speed. These two data structures are templates you plug in, not things you derive from scratch each time.
