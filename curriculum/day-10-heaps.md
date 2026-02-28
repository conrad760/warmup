# Day 10: Heaps & Priority Queues

> **Time budget:** 2 hours | **Prereqs:** Arrays, sorting intuition, basic graph traversal (Day 8)
> **Goal:** Recognize when a heap is the right tool, know the six core patterns, and write Go's `container/heap` boilerplate from memory — no hesitation, no fumbling with the interface.

---

## Go `container/heap` Boilerplate — Memorize This

Go doesn't have a built-in priority queue. You implement the `heap.Interface` (5 methods) on a slice type, then use `heap.Init`, `heap.Push`, and `heap.Pop` from the `container/heap` package. **Know this cold.** You'll write it in every heap problem.

### Min-Heap of ints

```go
import "container/heap"

type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool   { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{})  { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    val := old[n-1]
    *h = old[:n-1]
    return val
}
```

### Max-Heap of ints — flip `Less`

```go
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool   { return h[i] > h[j] } // only change
func (h MaxHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{})  { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    val := old[n-1]
    *h = old[:n-1]
    return val
}
```

### Heap of structs (common for weighted edges, freq pairs, etc.)

```go
type Item struct {
    val      int
    priority int
}

type ItemHeap []Item

func (h ItemHeap) Len() int            { return len(h) }
func (h ItemHeap) Less(i, j int) bool   { return h[i].priority < h[j].priority }
func (h ItemHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *ItemHeap) Push(x interface{})  { *h = append(*h, x.(Item)) }
func (h *ItemHeap) Pop() interface{} {
    old := *h
    n := len(old)
    val := old[n-1]
    *h = old[:n-1]
    return val
}
```

**Usage pattern (same every time):**
```go
h := &MinHeap{}     // or &MaxHeap{}, &ItemHeap{}
heap.Init(h)
heap.Push(h, 5)
heap.Push(h, 3)
top := heap.Pop(h).(int) // type assertion required — returns interface{}
peek := (*h)[0]           // peek without popping — index 0 is always the min/max
```

**Critical details:**
- `Push` and `Pop` on the type use pointer receivers (`*MinHeap`). `Len`, `Less`, `Swap` use value receivers.
- You call `heap.Push(h, val)` and `heap.Pop(h)` — the **package-level** functions, not the method directly. The package functions handle sift-up / sift-down for you.
- `Pop` returns `interface{}`. You must type-assert: `heap.Pop(h).(int)`.
- To peek at the top element without removing it: `(*h)[0]`.
- `heap.Init(h)` heapifies an existing slice in O(n). Use it when you pre-populate the slice.

---

## Pattern Catalog

### 1. Top-K Elements — Min-Heap of Size K

**Trigger:** "K largest elements," "K most frequent," "Kth largest element," any problem asking for the top/bottom K of something.

**Core idea:** Maintain a **min-heap of size K**. For every element, push it onto the heap. If the heap exceeds size K, pop the minimum. After processing all elements, the heap contains exactly the K largest. The heap root is the Kth largest.

**Why min-heap for K largest?** Counterintuitive but correct — the min-heap lets you efficiently evict the smallest of your K candidates. Whatever survives in the heap is guaranteed to be among the K largest.

**Go Template:**
```go
func topKLargest(nums []int, k int) []int {
    h := &MinHeap{}
    heap.Init(h)
    for _, num := range nums {
        heap.Push(h, num)
        if h.Len() > k {
            heap.Pop(h) // evict the smallest — doesn't belong in top K
        }
    }
    result := make([]int, h.Len())
    for i := len(result) - 1; i >= 0; i-- {
        result[i] = heap.Pop(h).(int)
    }
    return result
}
```

**Complexity:** O(n log k) time, O(k) space. Better than sorting (O(n log n)) when k << n.

**Watch out:**
- K largest → min-heap. K smallest → max-heap. Invert the intuition.
- If the problem asks for the Kth largest **single value**, just return `(*h)[0]` after processing — the heap root.
- If elements are structs (e.g., frequency pairs), define `Less` to compare the relevant field.

---

### 2. Merge K Sorted Lists/Arrays — Min-Heap of K Heads

**Trigger:** "Merge K sorted ___," "smallest range covering elements from K lists," any problem combining multiple sorted sequences.

**Core idea:** Put the head (first unconsumed element) of each list into a min-heap. Extract the minimum, then push the next element from that same list. Repeat until all lists are exhausted.

**Go Template (Merge K Sorted Lists — LC 23):**
```go
type ListNode struct {
    Val  int
    Next *ListNode
}

type NodeHeap []*ListNode

func (h NodeHeap) Len() int            { return len(h) }
func (h NodeHeap) Less(i, j int) bool   { return h[i].Val < h[j].Val }
func (h NodeHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *NodeHeap) Push(x interface{})  { *h = append(*h, x.(*ListNode)) }
func (h *NodeHeap) Pop() interface{} {
    old := *h
    n := len(old)
    val := old[n-1]
    *h = old[:n-1]
    return val
}

func mergeKLists(lists []*ListNode) *ListNode {
    h := &NodeHeap{}
    heap.Init(h)

    // Seed the heap with the head of each non-nil list
    for _, l := range lists {
        if l != nil {
            heap.Push(h, l)
        }
    }

    dummy := &ListNode{}
    cur := dummy
    for h.Len() > 0 {
        node := heap.Pop(h).(*ListNode)
        cur.Next = node
        cur = cur.Next
        if node.Next != nil {
            heap.Push(h, node.Next)
        }
    }
    return dummy.Next
}
```

**Complexity:** O(N log k) time where N = total nodes across all lists, O(k) space for the heap.

**Watch out:**
- Skip nil lists when seeding the heap. A nil push will panic.
- The heap always has at most K elements — one per list. This is what makes it efficient.
- For **arrays** instead of linked lists, track a `(value, listIndex, elementIndex)` tuple in the heap.

---

### 3. Running Median — Two Heaps (Max-Heap + Min-Heap)

**Trigger:** "Find median from data stream," "median of a sliding window," any problem needing the median after each insertion.

**Core idea:** Split the stream into two halves:
- `lo` — a **max-heap** holding the smaller half (root = largest of the small elements)
- `hi` — a **min-heap** holding the larger half (root = smallest of the large elements)

Keep them balanced: `lo.Len()` equals `hi.Len()` or `lo.Len()` equals `hi.Len() + 1`.

Median is either `lo`'s root (odd count) or the average of both roots (even count).

**Go Template (Find Median from Data Stream — LC 295):**
```go
type MedianFinder struct {
    lo *MaxHeap // smaller half — max-heap
    hi *MinHeap // larger half — min-heap
}

func Constructor() MedianFinder {
    lo := &MaxHeap{}
    hi := &MinHeap{}
    heap.Init(lo)
    heap.Init(hi)
    return MedianFinder{lo: lo, hi: hi}
}

func (mf *MedianFinder) AddNum(num int) {
    // Step 1: Always push to lo first
    heap.Push(mf.lo, num)

    // Step 2: Ensure lo's max <= hi's min (heap ordering property)
    if mf.hi.Len() > 0 && (*mf.lo)[0] > (*mf.hi)[0] {
        heap.Push(mf.hi, heap.Pop(mf.lo).(int))
    }

    // Step 3: Rebalance — lo can have at most 1 extra element
    if mf.lo.Len() > mf.hi.Len()+1 {
        heap.Push(mf.hi, heap.Pop(mf.lo).(int))
    } else if mf.hi.Len() > mf.lo.Len() {
        heap.Push(mf.lo, heap.Pop(mf.hi).(int))
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.lo.Len() > mf.hi.Len() {
        return float64((*mf.lo)[0])
    }
    return float64((*mf.lo)[0]+(*mf.hi)[0]) / 2.0
}
```

**Complexity:** O(log n) per `AddNum`, O(1) per `FindMedian`. Space O(n).

**Watch out:**
- The balancing logic has three parts: (1) push, (2) fix ordering violation, (3) fix size imbalance. Missing any one produces wrong results.
- Decide a convention and stick to it. Here: `lo` is allowed to have one extra. Some solutions allow `hi` to have one extra — either works, but be consistent.
- For the sliding-window median variant (LC 480), you also need lazy deletion — significantly harder. Know the two-heap version first.

---

### 4. K Closest / K Most Frequent — Heap as Selection

**Trigger:** "K closest points to origin," "K most frequent elements," problems where you compute a derived value (distance, frequency) and select the top K.

**Core idea:** Same as Pattern 1, but you first compute the relevant metric, then use a heap to select. Often involves a frequency map + heap, or a distance calculation + heap.

**Go Template (Top K Frequent Elements — LC 347):**
```go
type FreqItem struct {
    val  int
    freq int
}

type FreqHeap []FreqItem

func (h FreqHeap) Len() int            { return len(h) }
func (h FreqHeap) Less(i, j int) bool   { return h[i].freq < h[j].freq } // min-heap by freq
func (h FreqHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *FreqHeap) Push(x interface{})  { *h = append(*h, x.(FreqItem)) }
func (h *FreqHeap) Pop() interface{} {
    old := *h
    n := len(old)
    val := old[n-1]
    *h = old[:n-1]
    return val
}

func topKFrequent(nums []int, k int) []int {
    freq := map[int]int{}
    for _, n := range nums {
        freq[n]++
    }

    h := &FreqHeap{}
    heap.Init(h)
    for val, f := range freq {
        heap.Push(h, FreqItem{val: val, freq: f})
        if h.Len() > k {
            heap.Pop(h)
        }
    }

    result := make([]int, h.Len())
    for i := len(result) - 1; i >= 0; i-- {
        result[i] = heap.Pop(h).(FreqItem).val
    }
    return result
}
```

**Complexity:** O(n + u log k) where u = unique elements. Space O(u) for the freq map + O(k) for the heap.

**Watch out:**
- Build the frequency map **first**, then iterate the map (not the original array) to push into the heap. Otherwise you push duplicates.
- For K closest points, `Less` compares squared distances (avoid floating-point sqrt):
  ```go
  func (h PointHeap) Less(i, j int) bool {
      di := h[i].x*h[i].x + h[i].y*h[i].y
      dj := h[j].x*h[j].x + h[j].y*h[j].y
      return di > dj // max-heap — evict the farthest to keep K closest
  }
  ```
  K closest → max-heap (evict farthest). K farthest → min-heap (evict closest).

---

### 5. Task Scheduling with Cooldown — Max-Heap + Cooldown Queue

**Trigger:** "Task scheduler with cooldown interval," "reorganize string so no adjacent duplicates," problems where you must space out repeated elements.

**Core idea:** Greedily pick the most frequent remaining task (max-heap by frequency). After executing a task, it enters a cooldown queue with a timestamp for when it can re-enter the heap. Each time step, check if anything in the cooldown queue is ready.

**Go Template (Task Scheduler — LC 621):**
```go
type Pair struct {
    freq int
    ready int // earliest time this task can run again
}

func leastInterval(tasks []byte, n int) int {
    freq := map[byte]int{}
    for _, t := range tasks {
        freq[t]++
    }

    // Max-heap by frequency
    h := &MaxHeap{}
    heap.Init(h)
    for _, f := range freq {
        heap.Push(h, f)
    }

    time := 0
    cooldown := []Pair{} // queue of (freq, readyTime)

    for h.Len() > 0 || len(cooldown) > 0 {
        time++

        // Check if any task is ready to re-enter the heap
        if len(cooldown) > 0 && cooldown[0].ready <= time {
            heap.Push(h, cooldown[0].freq)
            cooldown = cooldown[1:]
        }

        if h.Len() > 0 {
            f := heap.Pop(h).(int)
            f--
            if f > 0 {
                cooldown = append(cooldown, Pair{freq: f, ready: time + n + 1})
            }
        }
        // else: idle cycle
    }
    return time
}
```

**Complexity:** O(N * log 26) ≈ O(N) where N = total tasks. At most 26 unique tasks.

**Watch out:**
- The cooldown queue is a plain FIFO queue, not a heap — tasks come off cooldown in the order they entered.
- `ready = time + n + 1`: if cooldown is `n`, and you execute at time `t`, the task is available at time `t + n + 1`.
- An idle cycle happens when the heap is empty but tasks are in cooldown. Time still advances.
- For "reorganize string" (LC 767), the pattern is similar but with cooldown = 1 (no adjacent duplicates).

---

### 6. Dijkstra's Shortest Path — Min-Heap of (distance, node)

**Trigger:** "Shortest path in a weighted graph," "minimum cost to reach ___," "network delay time," any weighted graph traversal seeking the minimum total weight.

**Core idea:** BFS doesn't work for weighted graphs. Instead, use a min-heap ordered by cumulative distance. Always process the node with the smallest known distance next. Once a node is popped, its shortest distance is finalized.

**Go Template (Network Delay Time — LC 743):**
```go
type Edge struct {
    to, weight int
}

type DistNode struct {
    dist, node int
}

type DistHeap []DistNode

func (h DistHeap) Len() int            { return len(h) }
func (h DistHeap) Less(i, j int) bool   { return h[i].dist < h[j].dist }
func (h DistHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *DistHeap) Push(x interface{})  { *h = append(*h, x.(DistNode)) }
func (h *DistHeap) Pop() interface{} {
    old := *h
    n := len(old)
    val := old[n-1]
    *h = old[:n-1]
    return val
}

func networkDelayTime(times [][]int, n int, k int) int {
    // Build adjacency list
    adj := make([][]Edge, n+1)
    for _, t := range times {
        adj[t[0]] = append(adj[t[0]], Edge{to: t[1], weight: t[2]})
    }

    dist := make([]int, n+1)
    for i := range dist {
        dist[i] = math.MaxInt64
    }
    dist[k] = 0

    h := &DistHeap{}
    heap.Init(h)
    heap.Push(h, DistNode{dist: 0, node: k})

    for h.Len() > 0 {
        cur := heap.Pop(h).(DistNode)
        if cur.dist > dist[cur.node] {
            continue // stale entry — already found a shorter path
        }
        for _, e := range adj[cur.node] {
            newDist := cur.dist + e.weight
            if newDist < dist[e.to] {
                dist[e.to] = newDist
                heap.Push(h, DistNode{dist: newDist, node: e.to})
            }
        }
    }

    maxDist := 0
    for i := 1; i <= n; i++ {
        if dist[i] == math.MaxInt64 {
            return -1 // unreachable node
        }
        if dist[i] > maxDist {
            maxDist = dist[i]
        }
    }
    return maxDist
}
```

**Complexity:** O((V + E) log V) time, O(V + E) space.

**Watch out:**
- **Stale entry check:** `if cur.dist > dist[cur.node] { continue }`. Without this, you reprocess nodes and potentially get TLE. This is the "lazy deletion" approach — cheaper than using a visited set and works identically.
- You can also use a `visited` bool array instead: skip if `visited[cur.node]` is true, set it true after popping.
- Dijkstra does NOT work with negative edge weights. If the interviewer introduces negative weights, pivot to Bellman-Ford.
- Don't confuse this with BFS. BFS uses a plain queue and works for unweighted graphs. Dijkstra uses a priority queue and works for non-negative weighted graphs.

---

## Decision Framework

Read the problem statement. Match the **first** rule that fits:

| Signal in problem | Heap type | Pattern |
|---|---|---|
| "K largest" / "Kth largest" | Min-heap of size K | Pattern 1 |
| "K smallest" / "Kth smallest" | Max-heap of size K | Pattern 1 (inverted) |
| "K most frequent" / "K closest" | Compute metric, then min-heap of size K | Pattern 4 |
| "Merge K sorted ___" | Min-heap of K elements | Pattern 2 |
| "Median from stream" / "running median" | Two heaps (max-heap lo, min-heap hi) | Pattern 3 |
| "Schedule tasks with cooldown" / "no adjacent duplicates" | Max-heap + cooldown queue | Pattern 5 |
| "Shortest path" / "minimum cost" (weighted graph) | Min-heap of (dist, node) | Pattern 6 |
| "Repeatedly need min or max" with dynamic insertions | Heap (type depends on need) | General |

**Heap vs. Sort:**
- If you only need the top K and n >> k: heap is O(n log k), sort is O(n log n). Heap wins.
- If you need the full sorted order: just sort. Don't complicate things with a heap.
- If you need repeated min/max extraction with interleaved insertions: heap. Sorting doesn't help because the data keeps changing.

**Heap vs. Quickselect:**
- Quickselect finds the Kth element in O(n) average but O(n²) worst case and is harder to code correctly.
- Heap is O(n log k) worst case and easier to code.
- In interviews, heap is almost always the safer choice unless the interviewer specifically asks about quickselect.

---

## Common Interview Traps

### 1. Top-K largest uses a MIN-heap
This trips up almost everyone the first time. You want the K largest, so you use a min-heap of size K. The min-heap evicts the smallest element, ensuring only the K largest survive. If you use a max-heap, you'd need to push everything and pop K times — that's O(n log n), no better than sorting.

### 2. Go's `container/heap` requires 5 methods
`Len`, `Less`, `Swap`, `Push`, `Pop`. You will be writing this boilerplate under time pressure. If you stumble on the interface, you burn 3-5 minutes. Practice until it's muscle memory:
- `Push` and `Pop` take pointer receivers and deal with `interface{}`
- `Len`, `Less`, `Swap` take value receivers
- `Pop` removes from the end of the slice (the heap package swaps the root to the end first, then calls your `Pop`)

### 3. Two-heap median: the balancing logic
Three things can go wrong: (1) You push to the wrong heap. (2) The heaps violate the ordering invariant (`lo`'s max > `hi`'s min). (3) The size balance is off. Always enforce all three in order: push → fix ordering → fix size.

### 4. Dijkstra: stale entries
When you update a node's distance and push a new entry, the old entry with the larger distance is still in the heap. When it gets popped later, you **must** skip it. Check: `if cur.dist > dist[cur.node] { continue }`. Missing this doesn't produce wrong answers (usually), but causes TLE on large inputs.

### 5. Forgetting type assertions on `heap.Pop`
`heap.Pop(h)` returns `interface{}`. Every. Single. Time. If you try to use the result directly without `.(int)` or `.(YourType)`, it won't compile. This is a pure time-waster in interviews.

### 6. Peeking without popping
`(*h)[0]` gives you the min (or max) without removing it. Don't call `heap.Pop` when you just want to look. This matters in the two-heap median pattern where you frequently compare the tops of both heaps.

### 7. Confusing `heap.Push`/`heap.Pop` with the method `Push`/`Pop`
Call the **package-level** functions: `heap.Push(h, val)` and `heap.Pop(h)`. These do the sift-up/sift-down. Your struct methods `Push` and `Pop` are helpers called by the package — you don't call them directly.

---

## Thought Process Walkthrough

### Problem 1: Top K Frequent Elements (LC 347)

**Problem:** Given an integer array `nums` and an integer `k`, return the `k` most frequent elements. You may return the answer in any order.

**Interview simulation — what to say out loud:**

**Step 1: Classify (15 seconds)**
> "K most frequent — this is the frequency + heap selection pattern. I'll build a frequency map, then use a min-heap of size K keyed by frequency."

**Step 2: Approach (30 seconds)**
> "First pass: count frequencies in a hash map. Second pass: iterate the map. For each (value, frequency) pair, push onto a min-heap ordered by frequency. If the heap exceeds size K, pop the minimum frequency element. After processing, the heap contains the K most frequent."

**Step 3: Edge cases (15 seconds)**
> "K equals the number of unique elements — return everything. All elements have the same frequency — any K elements work. Single element — return it."

**Step 4: Code**
```go
func topKFrequent(nums []int, k int) []int {
    freq := map[int]int{}
    for _, n := range nums {
        freq[n]++
    }

    h := &FreqHeap{}
    heap.Init(h)
    for val, f := range freq {
        heap.Push(h, FreqItem{val: val, freq: f})
        if h.Len() > k {
            heap.Pop(h)
        }
    }

    result := make([]int, h.Len())
    for i := len(result) - 1; i >= 0; i-- {
        result[i] = heap.Pop(h).(FreqItem).val
    }
    return result
}
```

**Step 5: Complexity analysis (say this)**
> "Time is O(n) for the frequency map, then O(u log k) for the heap operations where u is the number of unique elements. Total O(n + u log k). Space is O(u) for the map plus O(k) for the heap."

**Step 6: Test with example**
> Input: `nums = [1,1,1,2,2,3], k = 2`
>
> Frequency map: `{1:3, 2:2, 3:1}`
> Process 1 (freq 3): heap = `[(1,3)]`, size 1 ≤ 2
> Process 2 (freq 2): heap = `[(2,2), (1,3)]`, size 2 ≤ 2
> Process 3 (freq 1): heap = `[(3,1), (1,3), (2,2)]`, size 3 > 2 → pop min freq → evict (3,1)
> Heap = `[(2,2), (1,3)]`
> Result: `[1, 2]` ✓

**Step 7: Follow-up readiness**
- "What if K is very large?" → If k ≈ n, sorting is simpler and O(n log n) is fine.
- "Can you do better?" → Bucket sort by frequency: O(n) time. Create an array of size n+1, where index i holds all elements with frequency i. Walk backwards to collect the K most frequent. Mention this as the optimal approach.

---

### Problem 2: Find Median from Data Stream (LC 295)

**Problem:** Design a data structure that supports `addNum(int)` and `findMedian() float64` for a stream of integers.

**Interview simulation — what to say out loud:**

**Step 1: Classify (15 seconds)**
> "Running median — this is the two-heap pattern. Max-heap for the lower half, min-heap for the upper half."

**Step 2: Approach (45 seconds)**
> "I'll maintain two heaps: `lo` (max-heap, stores the smaller half) and `hi` (min-heap, stores the larger half). Invariant: `lo.Len()` equals `hi.Len()` or `lo.Len()` equals `hi.Len() + 1`. On each `addNum`:
> 1. Push to `lo`.
> 2. If `lo`'s max exceeds `hi`'s min, move `lo`'s root to `hi` (fix ordering).
> 3. If sizes are unbalanced, move from the larger heap to the smaller (fix size).
>
> For `findMedian`: if `lo` is bigger, return `lo`'s root. Otherwise average both roots."

**Step 3: Edge cases (15 seconds)**
> "First element — goes into `lo`, median is that element. Two elements — one in each heap, median is the average. All identical elements — works fine, both heaps hold the same value."

**Step 4: Code**
```go
type MedianFinder struct {
    lo *MaxHeap
    hi *MinHeap
}

func Constructor() MedianFinder {
    lo := &MaxHeap{}
    hi := &MinHeap{}
    heap.Init(lo)
    heap.Init(hi)
    return MedianFinder{lo: lo, hi: hi}
}

func (mf *MedianFinder) AddNum(num int) {
    heap.Push(mf.lo, num)

    // Fix ordering: lo's max must be <= hi's min
    if mf.hi.Len() > 0 && (*mf.lo)[0] > (*mf.hi)[0] {
        heap.Push(mf.hi, heap.Pop(mf.lo).(int))
    }

    // Fix size: lo can have at most 1 more than hi
    if mf.lo.Len() > mf.hi.Len()+1 {
        heap.Push(mf.hi, heap.Pop(mf.lo).(int))
    } else if mf.hi.Len() > mf.lo.Len() {
        heap.Push(mf.lo, heap.Pop(mf.hi).(int))
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.lo.Len() > mf.hi.Len() {
        return float64((*mf.lo)[0])
    }
    return float64((*mf.lo)[0]+(*mf.hi)[0]) / 2.0
}
```

**Step 5: Complexity analysis (say this)**
> "`AddNum` is O(log n) — at most 3 heap operations, each O(log n). `FindMedian` is O(1) — just peeking at roots. Space is O(n) total for both heaps."

**Step 6: Test with example**
> `addNum(1)`: lo = `[1]`, hi = `[]` → median = 1.0
> `addNum(2)`: push 2 to lo → lo = `[2,1]`, hi = `[]`
>   ordering: lo.max = 2, hi empty → skip
>   size: lo.Len() = 2 > hi.Len()+1 = 1 → move 2 to hi
>   lo = `[1]`, hi = `[2]` → median = (1+2)/2 = 1.5
> `addNum(3)`: push 3 to lo → lo = `[3,1]`, hi = `[2]`
>   ordering: lo.max = 3 > hi.min = 2 → move 3 to hi
>   lo = `[1]`, hi = `[2,3]`
>   size: hi.Len() = 2 > lo.Len() = 1 → move 2 to lo
>   lo = `[2,1]`, hi = `[3]` → median = 2.0 ✓

**Step 7: Follow-up readiness**
- "What if elements can be removed?" → Lazy deletion: mark removed elements, adjust counts, clean up when they surface at the root.
- "What about a sliding window median?" → Same two-heap structure plus lazy deletion with a hash map of pending removals (LC 480). Significantly harder.

---

## Time Targets

| Problem | Target | Notes |
|---|---|---|
| Kth Largest Element (LC 215) | 8 min | Min-heap of size K. Direct application. |
| Top K Frequent Elements (LC 347) | 10 min | Freq map + min-heap. Know bucket sort follow-up. |
| Merge K Sorted Lists (LC 23) | 12 min | Heap of list heads. Watch nil checks. |
| Find Median from Data Stream (LC 295) | 15 min | Two-heap balancing. Practice the addNum logic. |
| K Closest Points to Origin (LC 973) | 8 min | Max-heap of size K by squared distance. |
| Task Scheduler (LC 621) | 15 min | Max-heap + cooldown queue. Tricky simulation. |
| Network Delay Time (LC 743) | 12 min | Dijkstra. Stale entry check is key. |

If you're over these times, the bottleneck is almost always the `container/heap` boilerplate. Write it 5 times from memory until it takes under 60 seconds.

---

## Quick Drill

Five problems to solve **in order** during your practice session. Don't move on until the current one compiles and passes.

| # | Problem | Pattern | Target |
|---|---|---|---|
| 1 | [Kth Largest Element in an Array](https://leetcode.com/problems/kth-largest-element-in-an-array/) (LC 215) | Top-K (min-heap of size K) | 8 min |
| 2 | [Top K Frequent Elements](https://leetcode.com/problems/top-k-frequent-elements/) (LC 347) | Freq map + heap selection | 10 min |
| 3 | [Merge K Sorted Lists](https://leetcode.com/problems/merge-k-sorted-lists/) (LC 23) | Merge K sorted (min-heap of heads) | 12 min |
| 4 | [Find Median from Data Stream](https://leetcode.com/problems/find-median-from-data-stream/) (LC 295) | Two heaps (running median) | 15 min |
| 5 | [Network Delay Time](https://leetcode.com/problems/network-delay-time/) (LC 743) | Dijkstra (min-heap of dist, node) | 12 min |

**Total coding time:** ~57 min. Use remaining time for review and self-assessment.

---

## Self-Assessment

Answer these without looking at your notes. If you miss any, re-read the relevant section.

### 1. Why does "K largest" use a min-heap and not a max-heap?
**Expected answer:** A min-heap of size K keeps the K largest elements seen so far. When a new element arrives, if the heap is full, compare it to the root (the smallest of the K largest). If the new element is bigger, pop the root and push the new one. The min-heap efficiently evicts the smallest candidate. A max-heap would give you the overall maximum, not the K largest — you'd have to push everything and pop K times, which is O(n log n), no better than sorting.

### 2. Write the 5 methods of Go's `heap.Interface` for a min-heap of ints from memory.
**Expected answer:** `Len() int` returns `len(h)`. `Less(i, j int) bool` returns `h[i] < h[j]`. `Swap(i, j int)` swaps `h[i]` and `h[j]`. `Push(x interface{})` appends `x.(int)` to the slice. `Pop() interface{}` removes and returns the last element. Push and Pop use pointer receivers; Len, Less, Swap use value receivers.

### 3. In the two-heap median, what are the three steps of `AddNum` and why is each necessary?
**Expected answer:** (1) Push the new number to `lo` (max-heap). This is the default insertion point. (2) Fix the ordering invariant: if `lo`'s max > `hi`'s min, move `lo`'s root to `hi`. This ensures all elements in `lo` are ≤ all elements in `hi`. (3) Fix the size invariant: if either heap has more than 1 extra element, move its root to the other. This ensures the median is always accessible at the root(s).

### 4. In Dijkstra's algorithm, what happens if you skip the stale entry check?
**Expected answer:** The algorithm still produces correct shortest distances (because a stale entry has `cur.dist > dist[cur.node]`, so any relaxation it attempts will fail the `newDist < dist[e.to]` check). But it processes nodes multiple times unnecessarily, degrading performance from O((V+E) log V) toward O(V*E) in pathological cases. In interview problems with tight time limits, this causes TLE.

### 5. When should you use a heap vs. just sorting the array?
**Expected answer:** Use a heap when: (1) you only need the top K elements and K << n (heap is O(n log k) vs. sort's O(n log n)), (2) elements arrive in a stream and you can't wait for all of them, or (3) you need to interleave insertions and extractions (sort requires re-sorting after each insertion). Use sorting when you need the complete order or when K ≈ n.

### 6. What's the key difference between the `Push`/`Pop` methods you define and the `heap.Push`/`heap.Pop` package functions?
**Expected answer:** Your methods are low-level helpers: `Push` appends to the end of the slice, `Pop` removes the last element. They do NOT maintain heap order. The package functions `heap.Push` and `heap.Pop` call your methods AND perform sift-up / sift-down to maintain the heap invariant. Always call the package functions in your solution code; never call the methods directly.
