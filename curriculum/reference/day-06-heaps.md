# Day 6 — Heaps & Priority Queues: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Why It's Useful | Time |
|---|----------|----------------|------|
| 1 | [VisuAlgo — Binary Heap](https://visualgo.net/en/heap) | Interactive insert/extract-min animations. Step through sift-up and sift-down one swap at a time. | 10 min |
| 2 | [Building a Heap — O(n) Proof (YouTube, Back To Back SWE)](https://www.youtube.com/watch?v=MiyLo8adrWw) | Clear walkthrough of why heapify is O(n) with the summation argument. | 12 min |
| 3 | [Go `container/heap` package docs](https://pkg.go.dev/container/heap) | The official interface spec. Read the `IntHeap` example at the bottom — it's the canonical pattern. | 10 min |
| 4 | [Heap Visualization — USFCA](https://www.cs.usfca.edu/~galles/visualization/Heap.html) | Another visual tool. Good for watching heapify on an unsorted array (bottom-up build). | 5 min |
| 5 | [Priority Queue Applications (CP-algorithms)](https://cp-algorithms.com/data_structures/priority_queue.html) | Concise reference on when to use heaps: Dijkstra, merge-K, top-K, scheduling. | 10 min |
| 6 | [CLRS — Chapter 6: Heapsort](https://mitpress.mit.edu/9780262046305/introduction-to-algorithms/) | The textbook treatment of heap operations and the O(n) build-heap proof with the formal summation. If you have access, sections 6.1–6.3 are all you need. | 20 min |
| 7 | [D-ary Heaps — Wikipedia](https://en.wikipedia.org/wiki/D-ary_heap) | Short, clear explanation of the generalization from binary to d-ary heaps and the tradeoffs. | 5 min |

**Reading strategy:** Hit resources 1 and 4 first for visual intuition, then 2 or 6 for the O(n) proof, then 3 when you're ready to see how Go wraps it all.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 – 12:20 | Review (No Code)

| Min | Activity |
|-----|----------|
| 12:00 | Read the Day 6 section in OVERVIEW.md. Internalize the complexity table. |
| 12:05 | Open VisuAlgo or USFCA heap visualizer. Insert 10 random values one by one. Watch sift-up. Then extract-min repeatedly. Watch sift-down. |
| 12:10 | Study the array representation diagram (Section 6 below). Convince yourself of the index math: parent = `(i-1)/2`, left = `2*i+1`, right = `2*i+2`. |
| 12:15 | Read through the O(n) heapify proof (Section 3). Trace it on paper with a 7-element array. |

### 12:20 – 1:20 | Implement (From Scratch, No References)

| Min | Activity |
|-----|----------|
| 12:20 | Create `heap.go`. Define the `MinHeap` struct with a `data []int` field. |
| 12:25 | Implement `Len()`, `Peek()`, and helper index functions (`parent`, `leftChild`, `rightChild`). |
| 12:30 | Implement `siftUp(i int)` — swap with parent while smaller. |
| 12:37 | Implement `Push(val int)` — append to data, then siftUp on the last index. |
| 12:42 | Implement `siftDown(i int)` — compare with both children, swap with the smaller child. |
| 12:50 | Implement `Pop() int` — swap root with last element, shrink slice, siftDown(0). |
| 12:58 | Implement `Heapify(data []int) *MinHeap` — iterate from last non-leaf down to 0, calling siftDown. |
| 13:08 | Write a `main()` or test that pushes 10 values, peeks, pops all, and verifies sorted output. |
| 13:15 | Fix any bugs. Run through one more test by hand: Push 5 values, Pop 2, Push 3 more, Pop all. |

### 1:20 – 1:50 | Solidify (Edge Cases, Variants, Priority Queue)

| Min | Activity |
|-----|----------|
| 1:20 | Write tests for edge cases: empty heap Pop (should panic or return error), single element, all duplicate values, already-sorted input, reverse-sorted input. |
| 1:28 | Implement a `PriorityQueue` wrapper using your MinHeap — items have a `value` and `priority`. Lower priority = dequeued first. |
| 1:35 | Implement the top-K pattern: given a stream of ints, maintain a min-heap of size K and return the K largest seen so far. |
| 1:42 | (Stretch) Sketch the two-heap median pattern: max-heap for lower half, min-heap for upper half. Implement `AddNum` and `FindMedian`. |

### 1:50 – 2:00 | Recap (From Memory)

| Min | Activity |
|-----|----------|
| 1:50 | Close all files. Write down the complexity of Push, Pop, Peek, and Heapify. |
| 1:53 | Write down why Heapify is O(n) in one sentence. |
| 1:55 | Write down the index formulas for parent, left child, right child. |
| 1:57 | Write down one gotcha you hit during implementation. |

---

## 3. Core Concepts Deep Dive

### 3.1 Array Representation of a Complete Binary Tree

A heap is a **complete binary tree** — every level is full except possibly the last, which is filled left-to-right. This property means there are no gaps, so we can store it in a flat array with no pointers:

```
For node at index i (0-indexed):
  Parent:      (i - 1) / 2
  Left child:  2*i + 1
  Right child: 2*i + 2
```

**Why it works:** A complete binary tree has a unique level-order traversal. The array IS that traversal. Because there are no holes (completeness), the index math perfectly maps parent-child relationships without wasted space. Every element in `data[0..n-1]` is a real node.

**Why not use pointers?** Cache locality. A flat array is stored contiguously in memory, so traversing parent→child is just an index computation, not a pointer chase. This makes heaps extremely fast in practice.

### 3.2 Sift-Up vs. Sift-Down

| | Sift-Up | Sift-Down |
|---|---------|-----------|
| **Direction** | Leaf → Root | Root → Leaf |
| **Used in** | `Push` (insert) | `Pop` (extract), `Heapify` |
| **Mechanism** | Compare with parent; swap if smaller (min-heap) | Compare with both children; swap with the smaller child |
| **Terminates when** | Node >= parent, or node is at root | Node <= both children, or node is a leaf |
| **Worst case** | O(log n) — travel from leaf to root | O(log n) — travel from root to leaf |

**Why sift-down for heapify and not sift-up?** Starting from the bottom and sifting down gives O(n). Starting from the top and sifting up each element gives O(n log n). The asymmetry is explained next.

### 3.3 Heapify is O(n), NOT O(n log n) — The Proof

Building a heap by calling sift-down on every non-leaf node from bottom to top is O(n). This surprises people because each sift-down is O(log n), so the naive bound is O(n log n). The key insight: **most nodes are near the bottom and sift a short distance.**

In a complete binary tree with n nodes:
- ~n/2 nodes are leaves (sift distance 0 — we skip them)
- ~n/4 nodes are at height 1 (sift distance at most 1)
- ~n/8 nodes are at height 2 (sift distance at most 2)
- ...
- 1 node is at height log n (the root, sift distance at most log n)

The total work is:

```
        ⌊log n⌋
  T(n) =   Σ    (nodes at height h) × h
          h=0

        ⌊log n⌋
       =   Σ    ⌊n / 2^(h+1)⌋ × h
          h=0

              ⌊log n⌋
       ≤ n ×    Σ     h / 2^(h+1)
                h=0

              ∞
       ≤ n ×  Σ   h / 2^(h+1)
              h=0

              ∞
       = n/2 × Σ   h × x^h      where x = 1/2
               h=0

The series  Σ h·x^h = x / (1-x)²
evaluated at x = 1/2: (1/2) / (1/2)² = 2

So: T(n) ≤ n/2 × 2 = n

Therefore: T(n) = O(n)
```

**Intuition:** The sum converges because the number of nodes grows exponentially (2× per level) while the sift distance only grows linearly (+1 per level). The exponential growth dominates — the vast majority of nodes barely move.

If you built the heap by repeatedly calling sift-**up** (inserting one element at a time), you'd get O(n log n) because the ~n/2 leaf-level nodes each sift up O(log n) levels. Bottom-up heapify avoids this by sifting the heavy (high) nodes down a short distance.

### 3.4 Min-Heap vs. Max-Heap

The **only** difference is the comparison direction:

```
Min-heap: parent <= children      (root is minimum)
Max-heap: parent >= children      (root is maximum)
```

In code, you flip one comparison operator. Everything else — the index math, sift-up, sift-down, heapify — is structurally identical.

**Practical Go trick:** Implement a min-heap. To get a max-heap, negate values on insert and negate again on extract. Or use a comparator function:

```go
type Heap struct {
    data []int
    less func(a, b int) bool // min-heap: a < b, max-heap: a > b
}
```

### 3.5 Go's `container/heap` Interface

Go's standard library provides `container/heap`, which works through an interface:

```go
type Interface interface {
    sort.Interface          // Len() int, Less(i, j int) bool, Swap(i, j int)
    Push(x any)             // add x to the end
    Pop() any               // remove and return the last element
}
```

**How it works:**
- You define a type (usually a slice wrapper) and implement all 5 methods.
- `heap.Push(h, val)` calls your `Push` to append, then sifts up internally.
- `heap.Pop(h)` sifts down internally, then calls your `Pop` to remove the last element.
- `heap.Init(h)` calls heapify (bottom-up sift-down) on your data.

**Why `any` (formerly `interface{}`):** Go didn't have generics until 1.18. The `container/heap` package predates generics, so it uses the empty interface for element types. You type-assert on extraction: `val := heap.Pop(h).(int)`.

**Canonical example:**

```go
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }  // min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
    *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}
```

**Why implement your own first:** The `container/heap` package hides the sift-up and sift-down logic. You need to understand those mechanics before relying on the abstraction. Build it from scratch on Day 6; use `container/heap` in later days for convenience.

### 3.6 D-ary Heaps

A **d-ary heap** generalizes the binary heap so each node has `d` children instead of 2.

**Index math for d-ary heap (0-indexed):**
```
Parent of i:     (i - 1) / d
k-th child of i: d*i + k + 1    (for k = 0, 1, ..., d-1)
```

**Tradeoffs:**

| Property | Binary (d=2) | D-ary (d>2) |
|----------|-------------|-------------|
| Tree height | log₂ n | log_d n (shorter) |
| Sift-up cost | O(log₂ n) | O(log_d n) — **faster** |
| Sift-down cost | O(log₂ n) | O(d · log_d n) — **slower** (must compare d children) |
| Best for | General purpose | Decrease-key heavy workloads (e.g., Dijkstra) |

**When d > 2 helps:** Algorithms like Dijkstra call decrease-key (sift-up) far more often than extract-min (sift-down). A 4-ary heap reduces the tree height, making sift-up faster, at the cost of slightly more comparisons in sift-down. In practice, d=4 is often optimal due to cache line sizes.

**When d=2 is fine:** For simple priority queues where push and pop are equally frequent, binary heaps are simpler and fast enough.

---

## 4. Implementation Checklist

### Core API

```go
package heap

// MinHeap is a min-heap backed by a slice.
type MinHeap struct {
    data []int
}

// NewMinHeap creates an empty min-heap.
func NewMinHeap() *MinHeap

// Push adds a value to the heap.              O(log n)
func (h *MinHeap) Push(val int)

// Pop removes and returns the minimum value.  O(log n)
// Panics if the heap is empty.
func (h *MinHeap) Pop() int

// Peek returns the minimum value without removing it.  O(1)
// Panics if the heap is empty.
func (h *MinHeap) Peek() int

// Len returns the number of elements in the heap.      O(1)
func (h *MinHeap) Len() int

// --- internal ---

// siftUp restores the heap property by moving data[i] toward the root.
func (h *MinHeap) siftUp(i int)

// siftDown restores the heap property by moving data[i] toward the leaves.
func (h *MinHeap) siftDown(i int)

// Heapify builds a min-heap in-place from an unsorted slice.  O(n)
func Heapify(data []int) *MinHeap
```

### Priority Queue Wrapper

```go
type Item struct {
    Value    string
    Priority int // lower = higher priority (dequeued first)
}

type PriorityQueue struct {
    heap []Item
}

func (pq *PriorityQueue) Push(item Item)
func (pq *PriorityQueue) Pop() Item
func (pq *PriorityQueue) Peek() Item
func (pq *PriorityQueue) Len() int
```

### Test Plan

```go
func TestMinHeap(t *testing.T) {
    // Basic operations
    // - Push several values, Peek should always return the minimum
    // - Pop should return values in ascending order

    // Edge cases
    // - Empty heap: Len() == 0
    // - Single element: Push one, Peek, Pop, Len() == 0
    // - All same values: Push 5,5,5,5 — Pop should return 5 four times
    // - Already sorted input: Push 1,2,3,4,5 — verify heap property
    // - Reverse sorted input: Push 5,4,3,2,1 — Peek should be 1
    // - Interleaved push/pop: Push 3, Push 1, Pop (1), Push 2, Pop (2), Pop (3)

    // Heapify
    // - Heapify on an unsorted slice, then Pop all — should come out sorted
    // - Heapify on empty slice — Len() == 0
    // - Heapify on single element slice
    // - Heapify on already-a-heap input

    // Stress (optional)
    // - Push 10,000 random values, Pop all, verify sorted order
}
```

### Implementation Tips

1. **siftUp:** Start at index `i`. While `i > 0` and `data[i] < data[parent(i)]`, swap and move `i` to `parent(i)`.

2. **siftDown:** Start at index `i`. Loop: find the smallest among `data[i]`, `data[left(i)]`, `data[right(i)]`. If `i` is already smallest, stop. Otherwise swap with the smallest child and continue from that child's index. Guard against children being out of bounds.

3. **Pop:** Swap `data[0]` with `data[n-1]`, truncate slice (`data = data[:n-1]`), then `siftDown(0)`.

4. **Heapify:** Set `h.data = data`. Loop `i` from `len(data)/2 - 1` down to `0`, calling `siftDown(i)`.

---

## 5. Heap Application Patterns

### 5.1 Top-K Elements (Min-Heap of Size K)

**Problem:** Given n elements, find the K largest.

**Approach:** Maintain a min-heap of size K. For each element: if the heap has fewer than K items, push it. Otherwise, if the element is larger than the heap's minimum, pop the min and push the new element. At the end, the heap contains the K largest elements.

**Why min-heap for largest-K?** The min-heap's root is the *smallest* of the K candidates — the gatekeeper. Any new element must beat the gatekeeper to enter. This is O(n log K), which beats sorting O(n log n) when K << n.

```go
// TopK returns the k largest elements from nums (unordered).
func TopK(nums []int, k int) []int {
    h := NewMinHeap()
    for _, num := range nums {
        if h.Len() < k {
            h.Push(num)
        } else if num > h.Peek() {
            h.Pop()
            h.Push(num)
        }
    }
    result := make([]int, h.Len())
    for i := range result {
        result[i] = h.Pop()
    }
    return result
}
```

**Complexity:** O(n log K) time, O(K) space.

### 5.2 Merge K Sorted Lists

**Problem:** Given K sorted lists, merge them into one sorted list.

**Approach:** Push the head of each list into a min-heap (keyed by value). Pop the min, add it to the result, and push the next element from that same list. Repeat until the heap is empty.

**Complexity:** O(N log K) where N = total elements across all lists.

### 5.3 Running Median (Two-Heap Pattern)

**Problem:** Given a stream of numbers, support `AddNum(num)` and `FindMedian()`.

**Approach:** Maintain two heaps:
- `maxLo` — a **max-heap** holding the lower half of the numbers
- `minHi` — a **min-heap** holding the upper half

**Invariants:**
1. Every element in `maxLo` ≤ every element in `minHi`
2. Sizes differ by at most 1: `|maxLo.Len() - minHi.Len()| <= 1`

```go
type MedianFinder struct {
    lo *MaxHeap // max-heap: lower half (stores negated values in a MinHeap)
    hi *MinHeap // min-heap: upper half
}

func NewMedianFinder() *MedianFinder {
    return &MedianFinder{lo: NewMaxHeap(), hi: NewMinHeap()}
}

func (mf *MedianFinder) AddNum(num int) {
    // Always push to lo first, then rebalance
    mf.lo.Push(num)

    // Ensure lo's max <= hi's min
    if mf.hi.Len() > 0 && mf.lo.Peek() > mf.hi.Peek() {
        mf.hi.Push(mf.lo.Pop())
    } else {
        // Rebalance: lo can have at most 1 more than hi
        // Move lo's max to hi if lo is too big
    }

    // Balance sizes: lo.Len() == hi.Len() or lo.Len() == hi.Len() + 1
    if mf.lo.Len() > mf.hi.Len()+1 {
        mf.hi.Push(mf.lo.Pop())
    } else if mf.hi.Len() > mf.lo.Len() {
        mf.lo.Push(mf.hi.Pop())
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.lo.Len() > mf.hi.Len() {
        return float64(mf.lo.Peek())
    }
    return float64(mf.lo.Peek()+mf.hi.Peek()) / 2.0
}
```

**Complexity:** O(log n) per `AddNum`, O(1) per `FindMedian`.

**Max-heap trick in Go (negate values):**

```go
// MaxHeap wraps MinHeap by negating values.
type MaxHeap struct {
    h *MinHeap
}

func NewMaxHeap() *MaxHeap         { return &MaxHeap{h: NewMinHeap()} }
func (m *MaxHeap) Push(val int)    { m.h.Push(-val) }
func (m *MaxHeap) Pop() int        { return -m.h.Pop() }
func (m *MaxHeap) Peek() int       { return -m.h.Peek() }
func (m *MaxHeap) Len() int        { return m.h.Len() }
```

### 5.4 Dijkstra's Shortest Path

**Problem:** Find shortest paths from a source node in a weighted graph (non-negative weights).

**Approach:** Use a min-heap of `(distance, node)` pairs. Start with `(0, source)`. Pop the minimum-distance node, relax its neighbors — if a shorter path is found, push the new `(distance, neighbor)` pair.

**Complexity:** O((V + E) log V) with a binary heap.

### 5.5 Kth Largest Element in a Stream

**Problem:** Design a class that, given k and an initial array, can return the kth largest element after each new insertion.

**Approach:** Maintain a min-heap of size K. The root is always the Kth largest. When a new element arrives, if it's larger than the root, pop and push. The new root is the answer.

This is the same as Top-K but applied continuously to a stream.

---

## 6. Visual Diagrams

### 6.1 Min-Heap: Tree and Array Representation

```
          Tree View                        Array View

            1                     index: [ 0 ][ 1 ][ 2 ][ 3 ][ 4 ][ 5 ][ 6 ]
          /   \                   value: [ 1 ][ 3 ][ 2 ][ 7 ][ 6 ][ 5 ][ 4 ]
        3       2
       / \     / \                Index math:
      7   6   5   4              parent(i) = (i-1)/2    parent(5) = 2  ✓ data[2]=2 ≤ data[5]=5
                                 left(i)   = 2i+1       left(1)   = 3  ✓ data[1]=3 ≤ data[3]=7
                                 right(i)  = 2i+2       right(1)  = 4  ✓ data[1]=3 ≤ data[4]=6
```

### 6.2 Sift-Up During Insert (Push 0 into the heap above)

```
Step 0: Append 0 at index 7

          1                       [ 1 ][ 3 ][ 2 ][ 7 ][ 6 ][ 5 ][ 4 ][ 0 ]
        /   \                                                            ^
      3       2                                                         i=7
     / \     / \
    7   6   5   4
   /
  0    ← new element at index 7, parent = (7-1)/2 = 3, data[3]=7

Step 1: 0 < 7 → swap(7, 3)

          1                       [ 1 ][ 3 ][ 2 ][ 0 ][ 6 ][ 5 ][ 4 ][ 7 ]
        /   \                                       ^
      3       2                                    i=3
     / \     / \
    0   6   5   4                parent = (3-1)/2 = 1, data[1]=3
   /
  7

Step 2: 0 < 3 → swap(3, 1)

          1                       [ 1 ][ 0 ][ 2 ][ 3 ][ 6 ][ 5 ][ 4 ][ 7 ]
        /   \                             ^
      0       2                          i=1
     / \     / \
    3   6   5   4                parent = (1-1)/2 = 0, data[0]=1
   /
  7

Step 3: 0 < 1 → swap(1, 0)

          0                       [ 0 ][ 1 ][ 2 ][ 3 ][ 6 ][ 5 ][ 4 ][ 7 ]
        /   \                       ^
      1       2                    i=0    ← reached root, DONE
     / \     / \
    3   6   5   4
   /
  7
```

### 6.3 Sift-Down During Extract-Min (Pop from the heap above)

```
Step 0: Save root (0). Swap root with last element (7). Remove last.

          7                       [ 7 ][ 1 ][ 2 ][ 3 ][ 6 ][ 5 ][ 4 ]
        /   \                       ^
      1       2                    i=0
     / \     / \
    3   6   5   4                children: left=1 (val 1), right=2 (val 2)
                                 smallest child = 1 at index 1

Step 1: 7 > 1 → swap(0, 1)

          1                       [ 1 ][ 7 ][ 2 ][ 3 ][ 6 ][ 5 ][ 4 ]
        /   \                             ^
      7       2                          i=1
     / \     / \
    3   6   5   4                children: left=3 (val 3), right=4 (val 6)
                                 smallest child = 3 at index 3

Step 2: 7 > 3 → swap(1, 3)

          1                       [ 1 ][ 3 ][ 2 ][ 7 ][ 6 ][ 5 ][ 4 ]
        /   \                                       ^
      3       2                                    i=3
     / \     / \
    7   6   5   4                left child = 2*3+1 = 7, out of bounds → leaf, DONE

Final heap is valid. Returned value: 0
```

### 6.4 Heapify Process (Bottom-Up on Unsorted Array)

```
Input array: [4, 10, 3, 5, 1]

Tree form:         4
                 /   \
               10     3
              /  \
             5    1

Last non-leaf index = n/2 - 1 = 5/2 - 1 = 1

─── i=1: siftDown(1) ── node 10, children: 5, 1
    1 < 10 → swap(1, 4)
                   4                   [ 4 ][ 1 ][ 3 ][ 5 ][10 ]
                 /   \
                1     3
              /  \
             5   10

─── i=0: siftDown(0) ── node 4, children: 1, 3
    1 < 4 → swap(0, 1)
                   1                   [ 1 ][ 4 ][ 3 ][ 5 ][10 ]
                 /   \
                4     3
              /  \
             5   10

    continue sifting 4 at i=1, children: 5, 10
    4 < 5 and 4 < 10 → STOP

Final heap:        1                   [ 1 ][ 4 ][ 3 ][ 5 ][10 ]
                 /   \
                4     3                ✓ Valid min-heap
              /  \
             5   10

Total swaps: 2 (for 5 elements — well under 5 log 5 ≈ 11.6)
```

---

## 7. Self-Assessment

Answer these from memory after your session. If you can't, that's tomorrow's priority.

### Q1: Why does heapify start from the last non-leaf node and work upward, rather than starting from the root?

<details>
<summary>Answer</summary>

Heapify uses sift-down. Sift-down assumes the subtrees below are already valid heaps. By starting at the bottom (last non-leaf = index `n/2 - 1`) and working up, each node's children have already been heapified by the time we process it. Starting from the root would violate this precondition — the children wouldn't be heaps yet, so sift-down wouldn't restore the property correctly.

Leaves are already trivially valid heaps (no children), so we skip them.
</details>

### Q2: Why is a heap better than a sorted array for implementing a priority queue?

<details>
<summary>Answer</summary>

A sorted array gives O(1) peek and O(1) pop (from the end), but **O(n) insert** — you have to find the insertion point (O(log n) via binary search) and then shift elements to make room (O(n)).

A heap gives O(1) peek, O(log n) pop, and **O(log n) insert**. The insert advantage is decisive: priority queues need frequent inserts, and O(log n) vs. O(n) is the difference between usable and unusable at scale.
</details>

### Q3: You're implementing sift-down for a min-heap and your node has two children. You compare the node with only the left child and swap if needed. What bug does this introduce?

<details>
<summary>Answer</summary>

You might swap with the left child even though the right child is smaller. After the swap, the right child is smaller than the new parent, violating the heap property. **You must compare both children and swap with the smaller one.** This ensures the new parent is smaller than both children.

Example: parent=5, left=3, right=2. Swapping with left gives parent=3, left=5, right=2. But 2 < 3, so the heap property is broken at the root.
</details>

### Q4: In the O(n) heapify proof, why can't you achieve O(n) by calling sift-up on each element from index 0 to n-1 instead?

<details>
<summary>Answer</summary>

Sift-up moves a node *toward the root*, and its cost is proportional to the node's *depth*. The ~n/2 nodes at the bottom level have depth O(log n), so sifting them all up costs ~(n/2) × log n = **O(n log n)**.

Sift-down moves a node *toward the leaves*, and its cost is proportional to the node's *height*. The ~n/2 nodes at the bottom have height 0 (they're leaves — skipped entirely), while only 1 node (the root) sifts down O(log n). The sum converges to O(n) because the few nodes that travel far are vastly outnumbered by the many nodes that barely move.
</details>

### Q5: You need a data structure that supports: insert, delete-max, AND delete-min, all in O(log n). Can a single heap do this? What would you use instead?

<details>
<summary>Answer</summary>

A single min-heap gives O(log n) insert and O(log n) delete-min, but **delete-max is O(n)** — the max could be any leaf, and finding it requires scanning all ~n/2 leaves.

A single max-heap has the symmetric problem. For both operations in O(log n), you'd need either:
- A **min-max heap** (alternating min/max levels)
- Two heaps (min and max) cross-referenced with a hash map for lazy deletion
- A **balanced BST** (e.g., Go's `btree` or a red-black tree), which gives O(log n) for min, max, insert, and delete

This is a good reminder that heaps are specialized: they're optimal for one-sided priority access, not both.
</details>

---

## Complexity Reference (Quick Glance)

| Operation | Time | Notes |
|-----------|------|-------|
| Push | O(log n) | Append + sift-up |
| Pop | O(log n) | Swap root/last + sift-down |
| Peek | O(1) | Return `data[0]` |
| Heapify | **O(n)** | Bottom-up sift-down |
| Search | O(n) | Heaps are NOT search structures |
| Space | O(n) | Flat array, no pointers |
