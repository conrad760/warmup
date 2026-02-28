# Day 21 — Design Problems: LRU Cache, Min Stack, Median Finder

---

## 1. Curated Learning Resources

| # | Resource | Category | Why It's Useful |
|---|----------|----------|-----------------|
| 1 | [LRU Cache — NeetCode Video](https://www.youtube.com/watch?v=7ABFKPK2hD4) | Video | The clearest walkthrough of the hash map + doubly linked list combination. Draws the sentinel node approach and traces Get/Put step by step. Watch before coding. |
| 2 | [VisuAlgo — Linked List](https://visualgo.net/en/list) | Visualizer | Interactive DLL animations. Use this to convince yourself that Remove + InsertAfterHead are both O(1) when you have a direct pointer to the node. Essential mental model for LRU. |
| 3 | [Go `container/heap` Package Docs](https://pkg.go.dev/container/heap) | Go stdlib | You'll implement `heap.Interface` for the Median Finder's two heaps. Read the `IntHeap` example at the bottom — it's the canonical pattern you'll extend with a max-heap wrapper. |
| 4 | [LRU Cache Visualization — University of San Francisco](https://www.cs.usfca.edu/~galles/visualization/LRU.html) | Visualizer | Step through Put and Get operations and watch the DLL reorder. Shows eviction of the tail node when capacity is exceeded. |
| 5 | [Design Patterns in Go — Struct Embedding & Interfaces](https://go.dev/doc/effective_go#embedding) | Go patterns | Effective Go's section on embedding and interfaces. Relevant for understanding how to compose structs cleanly (e.g., MaxHeap wrapping MinHeap, LRUCache embedding a DList). |
| 6 | [Two Heaps Pattern — educative.io](https://www.educative.io/courses/grokking-the-coding-interview/3Y9jm7XRrXO) | Pattern guide | Explains the two-heap pattern for running median and its variants (sliding window median, IPO problem). Shows why a single heap fails and how balancing works. |
| 7 | [Min Stack — LeetCode Editorial](https://leetcode.com/problems/min-stack/editorial/) | Explanation | Walks through the auxiliary-stack approach and the single-stack-of-pairs approach. Good diagrams showing the min tracking through push/pop sequences. |
| 8 | [System Design — Cache Eviction Policies](https://en.wikipedia.org/wiki/Cache_replacement_policies) | Context | LRU in the real world: CPU caches, CDN caches, database buffer pools, Redis `maxmemory-policy`. Understanding the production context helps the design intuition stick. |

**Reading strategy:** Resources 1 and 4 first for LRU visual intuition (15 min). Resource 7 for Min Stack (5 min). Resource 6 for the two-heap pattern (5 min). Resource 3 when you implement the Median Finder. Resources 5 and 8 for enrichment after the session.

---

## 2. Detailed 2-Hour Session Plan

### Review Block (12:00 -- 12:20) — Internalize the Design Mindset

| Time | Activity |
|------|----------|
| 12:00 -- 12:05 | Read Section 5 (The Design Problem Methodology) below. Internalize the 5-step framework. Don't code yet — just think about the approach. |
| 12:05 -- 12:12 | **LRU Cache on paper.** List the required operations: `Get(key) -> O(1)`, `Put(key, value) -> O(1)`, evict LRU on capacity overflow. Ask: what gives O(1) lookup? (hash map.) What gives O(1) ordered insertion/removal? (DLL.) How to connect them? (map values are node pointers.) Draw the sentinel head/tail with 3 nodes in between. Trace a Get (move node to front) and a Put that triggers eviction (remove tail.Prev, delete from map). |
| 12:12 -- 12:16 | **Min Stack on paper.** List operations: `Push O(1)`, `Pop O(1)`, `Top O(1)`, `GetMin O(1)`. Why does a single `min` variable break? (Pop the current min — what's the new min?) Solution: each stack frame stores `(value, currentMin)`. Draw a stack of 3 pairs after pushing `5, 2, 7`. Verify GetMin returns 2, then pop `7` — GetMin still returns 2. Pop `2` — GetMin returns 5. |
| 12:16 -- 12:20 | **Median Finder on paper.** List operations: `AddNum O(log n)`, `FindMedian O(1)`. Why one heap fails: a min-heap gives you the smallest, not the middle. Two heaps: max-heap (lower half) and min-heap (upper half). The median is at one or both roots. Draw the state after adding `[5, 2, 8, 1]`. Verify the balancing invariant. |

### Implement Block (12:20 -- 1:20) — LRU Cache First (The Flagship)

| Time | Problem | Key Insight | Target |
|------|---------|-------------|--------|
| 12:20 -- 12:25 | **DNode + sentinel setup** | Define `DNode` with `key, val, prev, next`. Create `NewLRUCache` with sentinel `head` and `tail` linked to each other. The sentinel pattern eliminates all nil checks. | Struct definitions + constructor |
| 12:25 -- 12:32 | **DLL helpers**: `addAfterHead(node)`, `removeNode(node)`, `moveToFront(node)` | `addAfterHead` inserts between head sentinel and head.Next. `removeNode` patches prev/next pointers. `moveToFront` = remove + addAfterHead. | 3 helper methods, test each in isolation |
| 12:32 -- 12:42 | **LRUCache.Get(key)** | Look up in map. If miss, return -1. If hit, moveToFront, return val. | Test: put a key, get it, get a nonexistent key |
| 12:42 -- 12:58 | **LRUCache.Put(key, value)** | If key exists: update value, moveToFront. If key doesn't exist: create node, addAfterHead, add to map. If over capacity: remove `tail.Prev` (the LRU node), delete its key from the map. **Critical**: the DNode must store the key so you can delete it from the map during eviction. | Test: put 3 items in cap-2 cache, verify first is evicted. Put existing key — verify value updates and moves to front. |
| 12:58 -- 1:08 | **MinStack** | Use `[]struct{val, min int}`. On Push, new min = `min(val, current top's min)`. Pop just shrinks the slice. Top and GetMin read from the top element's fields. | Full implementation + tests |
| 1:08 -- 1:20 | **MedianFinder** | Max-heap for lower half (negate values in a min-heap), min-heap for upper half. `AddNum`: push to lo, rebalance so `lo.Peek() <= hi.Peek()`, then balance sizes. `FindMedian`: if odd total, return lo's root. If even, average both roots. | Use your Day 6 MinHeap or `container/heap`. Full implementation + tests. |

### Solidify Block (1:20 -- 1:50) — Edge Cases and Variants

| Time | Activity |
|------|----------|
| 1:20 -- 1:30 | **LRU Cache edge case sweep.** Test: (1) Get before any Put. (2) Put same key twice — second should update, not add duplicate. (3) Capacity 1 — every put evicts the previous. (4) Put→Get→Put→Get cycle that exercises both eviction and move-to-front. (5) Verify map size never exceeds capacity. |
| 1:30 -- 1:38 | **MinStack edge case sweep.** Test: (1) Push decreasing sequence `5, 3, 1` — GetMin at each step. (2) Push then pop back to empty — verify IsEmpty or panic behavior. (3) Push duplicates of the current min: `2, 2, 2` — pop all, min should update correctly. (4) All same values. |
| 1:38 -- 1:46 | **MedianFinder edge case sweep.** Test: (1) Single element — median is that element. (2) Two elements — average. (3) Add values in sorted order `1, 2, 3, 4, 5` — verify median at each step. (4) Add values in reverse order. (5) All duplicate values. (6) Negative numbers and zero. |
| 1:46 -- 1:50 | **Variant thinking** (no coding required). How would you modify LRU Cache to be an LFU Cache? (Add a frequency counter + map from frequency to DLL.) How would you make MedianFinder support `RemoveNum`? (Lazy deletion with a hash map of pending removals.) |

### Recap Block (1:50 -- 2:00) — Write From Memory

| Time | Activity |
|------|----------|
| 1:50 -- 1:54 | Close all references. From memory, write: (1) The data structures used in LRU Cache and why each is necessary. (2) Why a single min variable fails for MinStack. (3) The balancing invariant for MedianFinder's two heaps. |
| 1:54 -- 1:57 | Write the complexity of every operation for all three designs. |
| 1:57 -- 2:00 | Write one gotcha for each design problem — the bug you actually hit or would most likely hit in an interview. |

---

## 3. Core Concepts Deep Dive

### 3.1 The Design Problem Approach

Design problems differ from standard algorithm problems. Instead of finding a clever algorithm, you're **engineering a data structure** that meets specific complexity guarantees across multiple operations simultaneously.

The fundamental question is always: **"What combination of structures gives me O(target) for every required operation?"**

No single textbook data structure gives O(1) for everything. Hash maps have no order. Linked lists have no random access. Heaps have no O(1) arbitrary lookup. The art is in combining two structures so each covers the other's weakness.

### 3.2 LRU Cache: Hash Map + Doubly Linked List

**Why a hash map alone fails:** A hash map gives O(1) lookup by key, but it has no concept of ordering. You can't tell which key was least recently used. There's no way to find and evict the oldest entry without scanning all entries.

**Why a DLL alone fails:** A doubly linked list maintains insertion/access order perfectly, and removing a node given a pointer is O(1). But finding a node by key requires O(n) traversal — you'd have to walk the entire list.

**Why combining them works:** The hash map gives O(1) key → node pointer lookup. The DLL gives O(1) reordering (move a node to the front when accessed) and O(1) eviction (remove the tail node). Each structure compensates for the other's weakness.

**The "glue" is the node pointer.** The hash map doesn't store values directly — it stores pointers to DLL nodes. This cross-reference is what makes the combination work:

```
map[key] ──────────► DNode{key, val, prev, next}
                          ↑                ↓
                     DNode{...}  ←───  DNode{...}
```

**Sentinel nodes for clean code.** Without sentinels, every insert and remove operation needs `if head == nil` and `if tail == nil` checks. With dummy `head` and `tail` nodes that are always present:
- The list is never truly empty (the sentinels are always there).
- `addAfterHead` always inserts between `head` and `head.Next` — no nil check needed.
- `removeLRU` always removes `tail.Prev` — no nil check needed.
- The actual data lives between the sentinels: `head <-> [node1] <-> [node2] <-> ... <-> tail`.

**Why the DNode must store the key:** When you evict the LRU node (the one before `tail`), you need to delete its entry from the hash map. To do that, you need the key. If the node only stored the value, you'd have no way to find the corresponding map entry without a reverse lookup (which would be O(n) or require another hash map).

### 3.3 Min Stack: Why a Single Variable Breaks

Consider a stack with a single `min` variable:

```
Push 5: stack=[5], min=5          ✓
Push 2: stack=[5,2], min=2        ✓
Push 7: stack=[5,2,7], min=2      ✓
Pop 7:  stack=[5,2], min=2        ✓  (7 wasn't the min)
Pop 2:  stack=[5], min=???        ✗  (2 was the min — what's the new min?)
```

When you pop the current minimum, the previous minimum is lost. You'd need to scan the entire remaining stack to find it — O(n).

**The pair-stack solution:** Store `(value, currentMin)` at each level. The `currentMin` is the minimum of all elements at or below this position:

```
Push 5: stack=[(5, 5)]                   GetMin=5
Push 2: stack=[(5, 5), (2, 2)]           GetMin=2
Push 7: stack=[(5, 5), (2, 2), (7, 2)]   GetMin=2
Pop:    stack=[(5, 5), (2, 2)]           GetMin=2
Pop:    stack=[(5, 5)]                   GetMin=5  ← automatically restored!
```

Each Pop restores the previous min without any computation. The cost is ~2x memory (storing min alongside each value), but all operations remain O(1).

**Alternative — auxiliary min stack:** Instead of storing pairs, maintain a second stack that only tracks minimums. Push to the aux stack only when the new value is <= the current minimum. Pop from the aux stack only when the popped value equals the current minimum. This saves space when minimums change infrequently, but the pair approach is simpler and more interview-friendly.

### 3.4 Median Finder: The Two-Heap Approach

**Why a single heap fails:** A min-heap gives you the smallest element, not the median. A sorted array gives you the median in O(1) but insertion is O(n). We need O(log n) insert and O(1) median.

**The insight:** If you split all numbers into a lower half and an upper half at the median, then:
- The **maximum** of the lower half is either the median (odd count) or one of two values you average (even count).
- The **minimum** of the upper half is the other candidate.

A max-heap on the lower half gives you its maximum in O(1). A min-heap on the upper half gives you its minimum in O(1). Both heaps support O(log n) insertion.

**Invariants you must maintain:**

1. **Ordering**: Every element in `maxLo` ≤ every element in `minHi`. In practice, just ensure `maxLo.Peek() <= minHi.Peek()`.
2. **Size balance**: `maxLo.Len() == minHi.Len()` (even total) or `maxLo.Len() == minHi.Len() + 1` (odd total, median is `maxLo.Peek()`).

**AddNum algorithm:**

```
1. Push num to maxLo (the lower-half max-heap).
2. Push maxLo.Pop() to minHi (move the largest of the lower half to the upper half).
   This guarantees the ordering invariant: everything in maxLo ≤ everything in minHi.
3. If minHi.Len() > maxLo.Len(), push minHi.Pop() back to maxLo.
   This guarantees the size invariant: maxLo is always >= minHi in size.
```

This three-step process is cleaner than if-else branching because it always maintains both invariants regardless of input.

**FindMedian:**

```
If maxLo.Len() > minHi.Len():
    return maxLo.Peek()           // odd total — median is the middle element
Else:
    return (maxLo.Peek() + minHi.Peek()) / 2.0   // even total — average of two middle
```

### 3.5 Recurring Design Patterns

**Pattern 1: Hash map + ordered structure.** When you need both O(1) random access by key AND some ordering guarantee. Examples:
- LRU Cache: hash map + DLL (access order)
- LFU Cache: hash map + frequency map + DLL per frequency
- O(1) insert/delete/getRandom: hash map + array (swap-to-end for O(1) delete)

**Pattern 2: Two heaps for running statistics.** When you need the median (or any order statistic) of a dynamically changing set. The two heaps partition the data at the statistic you care about. Examples:
- Median Finder: max-heap lower + min-heap upper
- Sliding window median: same structure with lazy deletion
- IPO problem: max-heap for profits, min-heap for capital thresholds

**Pattern 3: Hash map + heap for priority with lookup.** When you need to both prioritize elements AND look them up by key. Examples:
- Task scheduler with cancellation
- LFU cache (hash map for O(1) key lookup, min-heap or DLL for eviction order)
- Dijkstra's algorithm (priority queue + visited set)

### 3.6 Go-Specific Design Considerations

**Struct methods in Go:** Go doesn't have classes, but struct methods with pointer receivers serve the same purpose. All design problem methods should use pointer receivers (`func (c *LRUCache) Get(...)`) since they mutate state.

**`container/heap` interface for Median Finder:** Go's `container/heap` requires you to implement 5 methods (Len, Less, Swap, Push, Pop). For a max-heap, just flip `Less`:

```go
type MaxIntHeap []int
func (h MaxIntHeap) Less(i, j int) bool { return h[i] > h[j] } // reversed!
```

Or use the negation trick: wrap a min-heap and negate values on Push/Pop/Peek.

**Pointers vs. values for DLL nodes:** Always use pointers (`*DNode`) in the hash map and DLL. Storing DNode by value would create copies — updating one wouldn't affect the other, breaking the cross-reference that makes the whole structure work.

**Struct embedding for composition:** Go's struct embedding can simplify wrapper types:

```go
type MaxHeap struct {
    MinHeap  // embed — inherits Len, siftUp, siftDown
}
// Override just Push/Pop/Peek to negate values
```

But for clarity in interviews, explicit fields are usually better than embedding.

---

## 4. Implementation Checklist

### LRU Cache — Struct Definitions and Method Signatures

```go
package design

// DNode is a doubly linked list node that stores both key and value.
// The key is needed so we can delete the map entry during eviction.
type DNode struct {
    key, val   int
    prev, next *DNode
}

// LRUCache combines a hash map with a doubly linked list.
// The list is ordered by access recency: most recent is right after head,
// least recent is right before tail. head and tail are sentinel nodes.
type LRUCache struct {
    cap        int
    cache      map[int]*DNode
    head, tail *DNode // sentinel nodes — never removed
}

// NewLRUCache creates an LRU cache with the given capacity.
// Initializes sentinel head and tail, linked to each other.
// Time: O(1)  Space: O(capacity)
func NewLRUCache(capacity int) *LRUCache {
    head := &DNode{}
    tail := &DNode{}
    head.next = tail
    tail.prev = head
    return &LRUCache{
        cap:   capacity,
        cache: make(map[int]*DNode),
        head:  head,
        tail:  tail,
    }
}

// Get retrieves the value for the key, or -1 if not found.
// On hit, moves the node to the front (most recently used).
// Time: O(1)
func (c *LRUCache) Get(key int) int

// Put inserts or updates the key-value pair.
// If the key exists, updates the value and moves to front.
// If the key is new and cache is at capacity, evicts the LRU entry (tail.prev).
// Time: O(1)
func (c *LRUCache) Put(key, value int)

// --- internal DLL helpers ---

// addAfterHead inserts node right after the head sentinel.
func (c *LRUCache) addAfterHead(node *DNode)

// removeNode detaches a node from the list (does not delete from map).
func (c *LRUCache) removeNode(node *DNode)

// moveToFront detaches a node and re-inserts it after head.
func (c *LRUCache) moveToFront(node *DNode)
```

### Min Stack — Struct Definition and Method Signatures

```go
// MinStack supports push, pop, top, and getMin all in O(1).
// Each stack frame stores the value and the running minimum
// at that stack depth.
type MinStack struct {
    stack []struct{ val, min int }
}

// NewMinStack creates an empty MinStack.
func NewMinStack() *MinStack

// Push adds val to the stack and updates the running min.
// Time: O(1)
func (s *MinStack) Push(val int)

// Pop removes the top element.
// Time: O(1)
func (s *MinStack) Pop()

// Top returns the top element without removing it.
// Time: O(1)
func (s *MinStack) Top() int

// GetMin returns the minimum element in the stack.
// Time: O(1)
func (s *MinStack) GetMin() int
```

### Median Finder — Struct Definition and Method Signatures

```go
import "container/heap"

// MaxIntHeap implements heap.Interface as a max-heap.
type MaxIntHeap []int

func (h MaxIntHeap) Len() int           { return len(h) }
func (h MaxIntHeap) Less(i, j int) bool { return h[i] > h[j] }  // max-heap
func (h MaxIntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxIntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MaxIntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// MinIntHeap implements heap.Interface as a min-heap.
type MinIntHeap []int

func (h MinIntHeap) Len() int           { return len(h) }
func (h MinIntHeap) Less(i, j int) bool { return h[i] < h[j] }  // min-heap
func (h MinIntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinIntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MinIntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// MedianFinder maintains two heaps to find the running median.
//   lo: max-heap holding the smaller half of the numbers
//   hi: min-heap holding the larger half of the numbers
// Invariant: lo.Len() == hi.Len() or lo.Len() == hi.Len() + 1
type MedianFinder struct {
    lo *MaxIntHeap // max-heap: lower half
    hi *MinIntHeap // min-heap: upper half
}

// NewMedianFinder creates an empty MedianFinder.
func NewMedianFinder() *MedianFinder

// AddNum adds a number to the data structure.
// Time: O(log n)
func (mf *MedianFinder) AddNum(num int)

// FindMedian returns the median of all added numbers.
// Time: O(1)
func (mf *MedianFinder) FindMedian() float64
```

### Test Cases

```go
// ─── LRU Cache Tests ───

func TestLRUCache_BasicGetPut(t *testing.T) {
    // LeetCode 146 example
    cache := NewLRUCache(2)
    cache.Put(1, 1)
    cache.Put(2, 2)
    assert(t, cache.Get(1) == 1, "Get(1) should return 1")
    cache.Put(3, 3)                              // evicts key 2
    assert(t, cache.Get(2) == -1, "Get(2) should return -1 after eviction")
    cache.Put(4, 4)                              // evicts key 1
    assert(t, cache.Get(1) == -1, "Get(1) should return -1 after eviction")
    assert(t, cache.Get(3) == 3, "Get(3) should return 3")
    assert(t, cache.Get(4) == 4, "Get(4) should return 4")
}

func TestLRUCache_UpdateExistingKey(t *testing.T) {
    cache := NewLRUCache(2)
    cache.Put(1, 10)
    cache.Put(1, 20)                             // update, not new entry
    assert(t, cache.Get(1) == 20, "Get(1) should return updated value 20")
    assert(t, len(cache.cache) == 1, "map should have exactly 1 entry")
}

func TestLRUCache_GetRefreshesRecency(t *testing.T) {
    cache := NewLRUCache(2)
    cache.Put(1, 1)
    cache.Put(2, 2)
    cache.Get(1)                                 // makes key 1 most recent
    cache.Put(3, 3)                              // should evict key 2, not key 1
    assert(t, cache.Get(2) == -1, "Get(2) should return -1 (evicted)")
    assert(t, cache.Get(1) == 1, "Get(1) should return 1 (refreshed)")
}

func TestLRUCache_CapacityOne(t *testing.T) {
    cache := NewLRUCache(1)
    cache.Put(1, 1)
    cache.Put(2, 2)                              // evicts key 1
    assert(t, cache.Get(1) == -1, "Get(1) should return -1")
    assert(t, cache.Get(2) == 2, "Get(2) should return 2")
}

func TestLRUCache_PutExistingKeyMovesToFront(t *testing.T) {
    cache := NewLRUCache(2)
    cache.Put(1, 1)
    cache.Put(2, 2)
    cache.Put(1, 10)                             // update key 1 — moves to front
    cache.Put(3, 3)                              // should evict key 2 (LRU), not key 1
    assert(t, cache.Get(2) == -1, "Get(2) evicted")
    assert(t, cache.Get(1) == 10, "Get(1) should return 10")
    assert(t, cache.Get(3) == 3, "Get(3) should return 3")
}

// ─── Min Stack Tests ───

func TestMinStack_BasicOperations(t *testing.T) {
    s := NewMinStack()
    s.Push(-2)
    s.Push(0)
    s.Push(-3)
    assert(t, s.GetMin() == -3, "GetMin should be -3")
    s.Pop()
    assert(t, s.Top() == 0, "Top should be 0")
    assert(t, s.GetMin() == -2, "GetMin should be -2 after popping -3")
}

func TestMinStack_DecreasingSequence(t *testing.T) {
    s := NewMinStack()
    s.Push(3)
    assert(t, s.GetMin() == 3, "min=3")
    s.Push(2)
    assert(t, s.GetMin() == 2, "min=2")
    s.Push(1)
    assert(t, s.GetMin() == 1, "min=1")
    s.Pop()
    assert(t, s.GetMin() == 2, "min restored to 2")
    s.Pop()
    assert(t, s.GetMin() == 3, "min restored to 3")
}

func TestMinStack_DuplicateMins(t *testing.T) {
    s := NewMinStack()
    s.Push(2)
    s.Push(2)
    s.Push(2)
    s.Pop()
    assert(t, s.GetMin() == 2, "min still 2 after popping one copy")
    s.Pop()
    assert(t, s.GetMin() == 2, "min still 2 after popping second copy")
}

func TestMinStack_IncreasingThenDecrease(t *testing.T) {
    s := NewMinStack()
    s.Push(5)
    s.Push(7)
    s.Push(3)
    s.Push(8)
    assert(t, s.GetMin() == 3, "min=3")
    s.Pop() // remove 8
    assert(t, s.GetMin() == 3, "min still 3")
    s.Pop() // remove 3
    assert(t, s.GetMin() == 5, "min restored to 5")
}

// ─── Median Finder Tests ───

func TestMedianFinder_OddCount(t *testing.T) {
    mf := NewMedianFinder()
    mf.AddNum(1)
    assertFloat(t, mf.FindMedian(), 1.0)
    mf.AddNum(2)
    assertFloat(t, mf.FindMedian(), 1.5)
    mf.AddNum(3)
    assertFloat(t, mf.FindMedian(), 2.0)
}

func TestMedianFinder_SortedInput(t *testing.T) {
    mf := NewMedianFinder()
    for _, v := range []int{1, 2, 3, 4, 5} {
        mf.AddNum(v)
    }
    assertFloat(t, mf.FindMedian(), 3.0)
}

func TestMedianFinder_ReverseSortedInput(t *testing.T) {
    mf := NewMedianFinder()
    for _, v := range []int{5, 4, 3, 2, 1} {
        mf.AddNum(v)
    }
    assertFloat(t, mf.FindMedian(), 3.0)
}

func TestMedianFinder_Duplicates(t *testing.T) {
    mf := NewMedianFinder()
    mf.AddNum(5)
    mf.AddNum(5)
    mf.AddNum(5)
    assertFloat(t, mf.FindMedian(), 5.0)
}

func TestMedianFinder_NegativeNumbers(t *testing.T) {
    mf := NewMedianFinder()
    mf.AddNum(-1)
    mf.AddNum(-2)
    assertFloat(t, mf.FindMedian(), -1.5)
    mf.AddNum(-3)
    assertFloat(t, mf.FindMedian(), -2.0)
}

func TestMedianFinder_StepByStep(t *testing.T) {
    // Trace from Section 6.3 visual diagram
    mf := NewMedianFinder()
    mf.AddNum(5)
    assertFloat(t, mf.FindMedian(), 5.0)
    mf.AddNum(2)
    assertFloat(t, mf.FindMedian(), 3.5)
    mf.AddNum(8)
    assertFloat(t, mf.FindMedian(), 5.0)
    mf.AddNum(1)
    assertFloat(t, mf.FindMedian(), 3.5)
    mf.AddNum(4)
    assertFloat(t, mf.FindMedian(), 4.0)
}
```

---

## 5. The Design Problem Methodology

A step-by-step framework for approaching any design problem in an interview or practice session.

### Step 1: List All Required Operations and Their Complexity Targets

Before thinking about implementation, write down every operation the problem demands and the time complexity it must achieve. This is your specification.

**LRU Cache worked example:**

```
Required operations:
  Get(key)        -> return value or -1       Target: O(1)
  Put(key, value) -> insert or update         Target: O(1)
                     evict LRU if at capacity Target: O(1)
```

### Step 2: For Each Operation, Ask "What Data Structure Gives This Complexity?"

Go through each operation and think about which fundamental data structure provides it:

```
O(1) lookup by key                → hash map
O(1) insert at a known position   → linked list (with pointer to position)
O(1) remove a known node          → doubly linked list (need prev pointer)
O(1) find the LRU element         → tail of an ordered list (DLL tail)
O(1) update recency on access     → move node to head of DLL
```

### Step 3: If No Single Structure Covers All Operations, Combine Two

No single structure from Step 2 covers everything:

```
Hash map alone:  ✓ O(1) lookup    ✗ No ordering (can't find LRU)
DLL alone:       ✓ O(1) reorder   ✗ O(n) lookup (must traverse)
```

**Combine them.** Hash map handles lookup. DLL handles ordering. Together they cover all operations.

### Step 4: Identify the "Glue" Between Structures

The two structures must reference each other to cooperate. Ask: **"What does the hash map store as its value?"**

```
cache: map[key] → *DNode    (pointer into the DLL)
DNode: {key, val, prev, next}
```

The map value is a **pointer to a DLL node**. This is the glue. When `Get(key)` is called:
1. Hash map gives you the node pointer in O(1).
2. You follow the pointer to the DLL node.
3. You call `moveToFront(node)` which is O(1) because you have direct access to prev/next.

The DNode stores the **key** — this is the reverse glue. When evicting:
1. DLL gives you the LRU node (`tail.prev`) in O(1).
2. You read `node.key` to know which map entry to delete.
3. You call `delete(cache, node.key)` in O(1).

### Step 5: Handle Edge Cases in the Interaction Between Structures

The combination introduces coupling. Changes to one structure must be reflected in the other:

```
Edge case                           What must happen
───────────                         ──────────────────
Put(key, val) on existing key       Update DNode val AND moveToFront
                                    Do NOT create a new node or map entry

Eviction during Put                 Remove DNode from DLL AND delete from map
                                    (if you forget the map delete, you leak entries)

Get(key) on hit                     Move node to front in DLL
                                    (if you forget this, eviction order breaks)

Put(key, val) at capacity           Remove tail.prev from DLL
                                    Delete tail.prev.key from map
                                    THEN insert the new node
                                    (order matters — remove before insert
                                     if capacity is 1)
```

### Applying This Framework to Other Design Problems

**Min Stack:**
1. Operations: Push O(1), Pop O(1), Top O(1), GetMin O(1).
2. Push/Pop/Top → stack. GetMin O(1) → need to track the running minimum.
3. Stack alone doesn't give O(1) GetMin (would need to scan). Combine: stack + min tracking.
4. Glue: each stack frame pairs `(value, currentMin)`. The min is derived from the previous frame's min.
5. Edge case: popping the current min automatically restores the previous min from the frame below.

**Median Finder:**
1. Operations: AddNum O(log n), FindMedian O(1).
2. FindMedian O(1) → need the middle element(s) accessible instantly. AddNum O(log n) → heap insertion speed.
3. One heap gives one extreme (min or max), not the middle. Two heaps splitting at the median give both middle candidates.
4. Glue: balancing invariant. After every AddNum, rebalance so `|lo.Len() - hi.Len()| <= 1` and `lo.Peek() <= hi.Peek()`.
5. Edge cases: first AddNum (only one heap has elements), even vs. odd total count (average vs. single root).

---

## 6. Visual Diagrams

### 6.1 LRU Cache: Hash Map + Doubly Linked List

```
Capacity = 3.  State after Put(1,A), Put(2,B), Put(3,C):

    Hash Map                    Doubly Linked List (MRU → LRU)
  ┌─────────────┐
  │ key │ value  │         ┌──────┐    ┌───────┐    ┌───────┐    ┌───────┐    ┌──────┐
  ├─────┼────────┤         │ HEAD │◄──►│ 3 : C │◄──►│ 2 : B │◄──►│ 1 : A │◄──►│ TAIL │
  │  1  │   ─────┼────┐    │sntnel│    │       │    │       │    │       │    │sntnel│
  │  2  │   ─────┼──┐ │    └──────┘    └───────┘    └───────┘    └───────┘    └──────┘
  │  3  │   ─────┼┐ │ │                   ▲              ▲             ▲
  └─────┴────────┘│ │ │                   │              │             │
                  │ │ └───────────────────┼──────────────┼─────────────┘
                  │ └─────────────────────┼──────────────┘
                  └───────────────────────┘

  Most recently used ← near HEAD       Least recently used → near TAIL


─── Get(2): hit ── Move node(2,B) to front ──

    Before: HEAD ↔ (3,C) ↔ (2,B) ↔ (1,A) ↔ TAIL

    Step 1: Remove (2,B) from current position
            HEAD ↔ (3,C) ↔ (1,A) ↔ TAIL       (2,B) detached

    Step 2: Insert (2,B) right after HEAD
            HEAD ↔ (2,B) ↔ (3,C) ↔ (1,A) ↔ TAIL

    After:  HEAD ↔ (2,B) ↔ (3,C) ↔ (1,A) ↔ TAIL
            map unchanged (still points to same node, just moved in list)


─── Put(4,D): cache full (cap=3) ── Evict LRU, insert new ──

    Before: HEAD ↔ (2,B) ↔ (3,C) ↔ (1,A) ↔ TAIL       map has keys {1,2,3}

    Step 1: Identify LRU = TAIL.prev = node(1,A)

    Step 2: Remove node(1,A) from DLL
            HEAD ↔ (2,B) ↔ (3,C) ↔ TAIL

    Step 3: Delete key 1 from map (node stores key=1 for this purpose)
            map now has keys {2,3}

    Step 4: Create new node(4,D), add to map
            map now has keys {2,3,4}

    Step 5: Insert node(4,D) after HEAD
            HEAD ↔ (4,D) ↔ (2,B) ↔ (3,C) ↔ TAIL

    After:  HEAD ↔ (4,D) ↔ (2,B) ↔ (3,C) ↔ TAIL       map has keys {2,3,4}
```

### 6.2 Min Stack: State After Push/Pop Operations

```
  Operations: Push(5), Push(3), Push(7), Push(1), Pop(), Pop()

  ┌──────────────────────────────────────────────────────────────────┐
  │  After each operation, showing (val, min) pairs:                 │
  │                                                                  │
  │  Push(5):      Push(3):      Push(7):      Push(1):             │
  │                                              ┌────────┐         │
  │                               ┌────────┐     │ (1, 1) │ ← top  │
  │                ┌────────┐     │ (7, 3) │     ├────────┤         │
  │  ┌────────┐    │ (3, 3) │     ├────────┤     │ (7, 3) │         │
  │  │ (5, 5) │    ├────────┤     │ (3, 3) │     ├────────┤         │
  │  │        │    │ (5, 5) │     ├────────┤     │ (3, 3) │         │
  │  └────────┘    └────────┘     │ (5, 5) │     ├────────┤         │
  │  GetMin=5      GetMin=3       └────────┘     │ (5, 5) │         │
  │                               GetMin=3       └────────┘         │
  │                                              GetMin=1           │
  │                                                                  │
  │  Pop(): remove (1,1)          Pop(): remove (7,3)               │
  │  ┌────────┐                   ┌────────┐                        │
  │  │ (7, 3) │ ← top             │ (3, 3) │ ← top                 │
  │  ├────────┤                   ├────────┤                        │
  │  │ (3, 3) │                   │ (5, 5) │                        │
  │  ├────────┤                   └────────┘                        │
  │  │ (5, 5) │                   GetMin=3                          │
  │  └────────┘                                                     │
  │  GetMin=3                     ← min automatically restored!     │
  └──────────────────────────────────────────────────────────────────┘

  Key insight: each frame's "min" field records the minimum of
  ALL elements at or below that frame. Pop never needs to recompute.
```

### 6.3 Median Finder: Two Heaps After Several AddNum Operations

```
  Operations: AddNum(5), AddNum(2), AddNum(8), AddNum(1), AddNum(4)

  ─── AddNum(5) ───
  maxLo (lower half):  [5]          minHi (upper half):  []
  Median = 5.0                      (odd count: maxLo root)

  ─── AddNum(2) ───
  Step: push 2 to lo → lo=[5,2], pop lo max(5) to hi → lo=[2], hi=[5]
  maxLo: [2]                         minHi: [5]
  Median = (2 + 5) / 2 = 3.5        (even count: average of roots)

  ─── AddNum(8) ───
  Step: push 8 to lo → lo=[8,2], pop lo max(8) to hi → lo=[2], hi=[5,8]
        hi larger → pop hi min(5) to lo → lo=[5,2], hi=[8]
  maxLo: [5, 2]                      minHi: [8]

         maxLo (max-heap)               minHi (min-heap)
            5                               8
           /
          2
  Median = 5.0                       (odd count: maxLo root)

  ─── AddNum(1) ───
  Step: push 1 to lo → lo=[5,2,1], pop lo max(5) to hi → lo=[2,1], hi=[5,8]
  maxLo: [2, 1]                      minHi: [5, 8]

         maxLo (max-heap)               minHi (min-heap)
            2                               5
           /                               /
          1                               8
  Median = (2 + 5) / 2 = 3.5         (even count: average of roots)

  ─── AddNum(4) ───
  Step: push 4 to lo → lo=[4,2,1], pop lo max(4) to hi → lo=[2,1], hi=[4,5,8]
        hi larger → pop hi min(4) to lo → lo=[4,2,1], hi=[5,8]
  maxLo: [4, 2, 1]                   minHi: [5, 8]

         maxLo (max-heap)               minHi (min-heap)
            4                               5
           / \                             /
          2   1                           8
  Median = 4.0                        (odd count: maxLo root)

  ─── Verification ───
  Sorted: [1, 2, 4, 5, 8]
  Median of 5 elements = element at index 2 = 4  ✓

  Invariants check:
    maxLo.Peek() = 4 ≤ minHi.Peek() = 5    ✓ (ordering)
    maxLo.Len() = 3, minHi.Len() = 2       ✓ (size: lo = hi + 1)
```

---

## 7. Self-Assessment

Answer these without looking at the material. If you can't answer confidently, revisit the relevant section.

### Question 1: What breaks if LRU Cache Put doesn't move existing keys to the front?

<details>
<summary>Check your answer</summary>

If you update an existing key's value without calling `moveToFront`, the node stays at its old position in the DLL. This means a recently-accessed key can be near the tail — the "least recently used" position. A subsequent Put that triggers eviction will remove `tail.prev`, which might be the key you just updated. The cache would evict recently-used entries while keeping stale ones.

Example:
```
cache = LRU(2)
Put(1, A)        → HEAD ↔ (1,A) ↔ TAIL
Put(2, B)        → HEAD ↔ (2,B) ↔ (1,A) ↔ TAIL
Put(1, A')       → If we DON'T move to front: HEAD ↔ (2,B) ↔ (1,A') ↔ TAIL
Put(3, C)        → Evicts tail.prev = (1,A') — WRONG! Key 1 was just accessed!
```

The fix: `Put` on an existing key must both update the value AND call `moveToFront`.
</details>

### Question 2: Why can't you use a single heap for the Median Finder?

<details>
<summary>Check your answer</summary>

A min-heap gives O(1) access to the minimum, and a max-heap gives O(1) access to the maximum. Neither gives O(1) access to the median (the middle element).

To find the median in a single heap, you'd need to pop n/2 elements, read the top, then push them all back — that's O(n log n) per query.

The two-heap approach works because it splits the data at the median. The max-heap's root is the largest element in the lower half, and the min-heap's root is the smallest in the upper half. These two roots are exactly the candidates for the median, both accessible in O(1).
</details>

### Question 3: In the MinStack, what happens if you push the same value as the current minimum?

<details>
<summary>Check your answer</summary>

You must still record the minimum in the new frame. If you push value `2` when the current min is already `2`, the new frame is `(2, 2)`. This seems redundant, but it's essential.

If you skip recording the min for duplicate values, then popping this frame would incorrectly change the min to whatever the frame below says — which might not be `2`. The pair approach handles this naturally: `min(2, currentMin=2) = 2`, so the new frame is `(2, 2)` and everything works.

With the auxiliary-stack variant, you must push to the aux stack when `val <= auxTop` (note: less than OR EQUAL, not strictly less than). Forgetting the equality causes a bug when duplicate minimums are popped.
</details>

### Question 4: Why does the DNode in LRU Cache need to store the key?

<details>
<summary>Check your answer</summary>

During eviction, you identify the LRU node via the DLL: it's `tail.prev`. You then need to delete this entry from the hash map. But `delete(map, key)` requires the key. If the node only stored the value, you'd have no way to determine which map key corresponds to this node without an expensive reverse lookup.

The node stores `(key, val)` specifically so that eviction can:
1. Get the LRU node: `victim = tail.prev`
2. Delete from map: `delete(cache, victim.key)`
3. Remove from DLL: `removeNode(victim)`

All in O(1). Without the key in the node, step 2 would require either O(n) map scan or a second reverse map, both of which defeat the purpose.
</details>

### Question 5: What happens to the Median Finder if you don't rebalance the heaps after each AddNum?

<details>
<summary>Check your answer</summary>

Without rebalancing, one heap could grow much larger than the other. For example, if you always push to `maxLo` first, it would accumulate all elements while `minHi` stays empty.

The consequence for `FindMedian`: the formula `(lo.Peek() + hi.Peek()) / 2` requires both heaps to be non-empty for even counts. Even for odd counts, `lo.Peek()` would be the maximum of ALL elements, not the median.

More subtly, even if both heaps have elements, unbalanced sizes mean the "split point" isn't at the median. If `lo` has 5 elements and `hi` has 1, `lo.Peek()` is the max of 5 elements — the 5th out of 6, not the median.

The balancing invariant (`|lo.Len() - hi.Len()| <= 1`) ensures the split is always at the median position.
</details>

---

## 8. Week 3 Retrospective

This is the final day of the 3-week curriculum. Here's how to assess your progress and continue growing.

### How to Continue After the 3 Weeks

**Spaced repetition with the warmup tool.** The problems you built over 21 days are now your rotation set. Revisit them on a schedule:
- **Day 1-7 structures:** Re-implement one per week from scratch. If you can write an LRU Cache or a Trie in under 15 minutes without references, it's internalized.
- **Day 8-14 patterns:** Pick one pattern per week and solve 2-3 new LeetCode problems using it. The warmup tool can queue these.
- **Day 15-21 advanced topics:** These need the most repetition. Revisit DP problems every 3-4 days until the recurrences are automatic.

**Weekly reviews.** Every Sunday, spend 30 minutes:
1. Pick one topic at random. Explain it aloud as if teaching someone.
2. Write the key complexity table from memory.
3. Implement the core operation (e.g., sift-down, DLL remove, partition) without references.

### Signs You've Internalized a Topic vs. Need More Practice

**Internalized:**
- You can implement it from scratch in under 15 minutes without hesitation.
- You can explain *why* each design choice was made, not just *what* the code does.
- When you see a new problem, you immediately recognize which pattern applies.
- You can describe the complexity of every operation without looking it up.
- You can explain what breaks if you skip a step (e.g., "if you don't move to front on Get, eviction order breaks").

**Needs more practice:**
- You remember the general idea but get stuck on details (e.g., sentinel node setup, sift-down comparison logic).
- You can solve the problem you practiced but struggle with variants.
- You need to look up the complexity table.
- You can't explain *why* — only *what*.
- You make the same off-by-one or edge case mistakes repeatedly.

### Suggested Next Steps

**Mock interviews (start within 1 week).** The gap between solving problems alone and solving them under pressure with someone watching is enormous. Practice with:
- A friend or study partner (trade interviewer/interviewee roles).
- Platforms like Pramp, interviewing.io, or Exercism.
- Record yourself solving a problem and watch it back — you'll notice verbal habits and hesitations.

**Contest problems (optional, for sharpening speed).** LeetCode weekly contests or Codeforces Div 2. The time pressure forces you to recognize patterns quickly. Don't worry about ranking — focus on solving A and B cleanly.

**System design (if targeting senior+ roles).** The data structure knowledge from these 3 weeks directly feeds system design:
- LRU Cache → cache layer design, Redis eviction policies
- Heaps → priority-based scheduling, rate limiters
- Tries → autocomplete systems, IP routing tables
- Graphs → social networks, dependency resolution, service mesh routing
- Hash maps → distributed hash tables, consistent hashing

**Build something.** The deepest way to internalize data structures is to use them in a real project. Ideas:
- An in-memory key-value store with TTL and LRU eviction (uses hash maps, DLLs, heaps for TTL).
- A simple search engine with a trie-based autocomplete.
- A task scheduler with priority queues and dependency resolution (heaps + topological sort).

### The One Rule Going Forward

**If you can't implement it without looking at references, you don't know it yet.** That's not failure — that's signal. Spend more time on it. The goal isn't to have seen every problem; it's to have a toolkit of patterns you can deploy from memory under pressure.

---

## Complexity Reference (Quick Glance)

| Design | Operation | Time | Space |
|--------|-----------|------|-------|
| LRU Cache | Get | O(1) | O(capacity) |
| LRU Cache | Put | O(1) | O(capacity) |
| Min Stack | Push | O(1) | O(n) |
| Min Stack | Pop | O(1) | O(n) |
| Min Stack | Top | O(1) | O(n) |
| Min Stack | GetMin | O(1) | O(n) |
| Median Finder | AddNum | O(log n) | O(n) |
| Median Finder | FindMedian | O(1) | O(n) |
