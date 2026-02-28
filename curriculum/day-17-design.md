# Day 17: Design Problems

> **Time:** 2 hours | **Level:** Refresher | **Language:** Go
>
> Design problems test whether you can combine data structures under pressure.
> There are only ~5 patterns that matter. Know them cold.

---

## Pattern Catalog

### Pattern 1: Hash Map + Doubly Linked List — LRU Cache

**Trigger:** "Design a cache with O(1) get and put, evict least recently used."

**Why two structures:** Hash map gives O(1) lookup. DLL gives O(1) removal and
insertion at ends. Neither alone provides both. The hash map stores
`key -> *Node`, and the node lives in the DLL. That pointer is the glue.

```go
type Node struct {
    key, val   int
    prev, next *Node
}

type LRUCache struct {
    cap        int
    cache      map[int]*Node
    head, tail *Node // sentinels
}

func Constructor(capacity int) LRUCache {
    head := &Node{}
    tail := &Node{}
    head.next = tail
    tail.prev = head
    return LRUCache{
        cap:   capacity,
        cache: make(map[int]*Node),
        head:  head,
        tail:  tail,
    }
}

func (l *LRUCache) remove(node *Node) {
    node.prev.next = node.next
    node.next.prev = node.prev
}

func (l *LRUCache) insertAfterHead(node *Node) {
    node.next = l.head.next
    node.prev = l.head
    l.head.next.prev = node
    l.head.next = node
}

func (l *LRUCache) Get(key int) int {
    node, ok := l.cache[key]
    if !ok {
        return -1
    }
    l.remove(node)
    l.insertAfterHead(node)
    return node.val
}

func (l *LRUCache) Put(key, value int) {
    if node, ok := l.cache[key]; ok {
        node.val = value
        l.remove(node)
        l.insertAfterHead(node)
        return
    }
    node := &Node{key: key, val: value}
    l.cache[key] = node
    l.insertAfterHead(node)
    if len(l.cache) > l.cap {
        victim := l.tail.prev
        l.remove(victim)
        delete(l.cache, victim.key) // MUST delete from map too
    }
}
```

**Complexity:** O(1) get, O(1) put.

**Watch out:**
- Put on an existing key must update the value AND move to front. Forgetting
  the move-to-front is the single most common bug.
- Eviction must `delete(l.cache, victim.key)`. Forgetting this leaks entries.
- The node stores `key` so that during eviction you can look up the map entry.
  Without the key on the node, you cannot delete from the map.
- Sentinel nodes (`head`, `tail`) eliminate every nil check. Always use them.

---

### Pattern 2: Pair Stack / Auxiliary Stack — Min Stack

**Trigger:** "Stack with O(1) push, pop, top, and getMin."

**Why it works:** Each stack frame records the min at that depth. When you pop,
the previous frame's min is automatically correct because it was computed
without knowledge of the popped element.

```go
type MinStack struct {
    stack []entry
}

type entry struct {
    val, min int
}

func Constructor() MinStack {
    return MinStack{}
}

func (s *MinStack) Push(val int) {
    curMin := val
    if len(s.stack) > 0 && s.stack[len(s.stack)-1].min < val {
        curMin = s.stack[len(s.stack)-1].min
    }
    s.stack = append(s.stack, entry{val, curMin})
}

func (s *MinStack) Pop() {
    s.stack = s.stack[:len(s.stack)-1]
}

func (s *MinStack) Top() int {
    return s.stack[len(s.stack)-1].val
}

func (s *MinStack) GetMin() int {
    return s.stack[len(s.stack)-1].min
}
```

**Complexity:** O(1) all operations. O(n) space.

**Watch out:**
- The pair approach (storing min alongside each value) is simpler and less
  error-prone than maintaining a separate auxiliary stack. Prefer it.
- Pop does not need to return anything or do any comparison — the previous
  entry already has the correct min.

---

### Pattern 3: Two Heaps — Find Median from Data Stream

**Trigger:** "Continuously add numbers, query median at any time."

**Why two heaps:** A max-heap holds the smaller half, a min-heap holds the
larger half. The median is at one or both tops. No single structure gives
O(log n) insert with O(1) median.

```go
import "container/heap"

// MaxHeap: stores the lower half. Top = largest of the small numbers.
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] } // max-heap
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
    old := *h
    val := old[len(old)-1]
    *h = old[:len(old)-1]
    return val
}

// MinHeap: stores the upper half. Top = smallest of the large numbers.
type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] } // min-heap
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
    old := *h
    val := old[len(old)-1]
    *h = old[:len(old)-1]
    return val
}

type MedianFinder struct {
    lo *MaxHeap // lower half
    hi *MinHeap // upper half
}

func Constructor() MedianFinder {
    lo := &MaxHeap{}
    hi := &MinHeap{}
    heap.Init(lo)
    heap.Init(hi)
    return MedianFinder{lo: lo, hi: hi}
}

func (mf *MedianFinder) AddNum(num int) {
    // Step 1: push to lo
    heap.Push(mf.lo, num)
    // Step 2: move lo's top to hi (ensures lo's max <= hi's min)
    heap.Push(mf.hi, heap.Pop(mf.lo))
    // Step 3: rebalance — lo should have >= hi's count
    if mf.lo.Len() < mf.hi.Len() {
        heap.Push(mf.lo, heap.Pop(mf.hi))
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.lo.Len() > mf.hi.Len() {
        return float64((*mf.lo)[0])
    }
    return float64((*mf.lo)[0]+(*mf.hi)[0]) / 2.0
}
```

**Complexity:** O(log n) AddNum, O(1) FindMedian.

**Watch out:**
- The 3-step add pattern (push lo, move top to hi, rebalance) is the cleanest.
  Memorize this exact sequence — it avoids every ordering edge case.
- `(*mf.lo)[0]` accesses the heap top without popping. This is the Go idiom.
- Go's `container/heap` distinction: you implement `Push`/`Pop` on the
  interface (append/truncate the slice), then call `heap.Push`/`heap.Pop`
  (the package functions) which call your methods plus sift. Never call
  your interface's Push/Pop directly — always go through `heap.Push`/`heap.Pop`.

---

### Pattern 4: Hash Map + Heap — LFU Cache (Concept Only)

**Trigger:** "Evict the least frequently used item; ties broken by LRU."

You will not be asked to implement this from scratch in 45 minutes. Know the
idea: maintain a hash map from key to node, and group nodes by frequency. Each
frequency bucket is itself an LRU list (a doubly linked list). Track the
current minimum frequency. On access, move the node from frequency `f` to
frequency `f+1`. If the old bucket is empty and was the min, increment min.

If an interviewer asks this, describe the architecture, then ask which part
they want you to code.

---

### Pattern 5: The General Approach — Combining Structures

Every design problem follows the same logic:

| Required operation | Structure that gives it |
|--------------------|------------------------|
| O(1) lookup by key | Hash map |
| O(1) ordered insert/remove | Doubly linked list |
| O(1) min/max | Heap, or pair-stack |
| O(log n) sorted insert | Balanced BST or heap |
| O(1) stack/queue ops | Slice or linked list |

When no single structure covers all required operations, you combine two.
The hash map almost always appears — it is the "index" into the other
structure.

---

## The Design Problem Methodology

Use this framework for ANY design problem. Say each step out loud in the
interview.

```
Step 1: LIST the required operations and their time complexity targets.
        "We need Get in O(1), Put in O(1), eviction of the LRU element."

Step 2: For EACH operation, what structure gives that complexity?
        "O(1) lookup → hash map. O(1) order tracking → linked list."

Step 3: If no single structure covers all ops, COMBINE two.
        "Hash map alone can't track order. DLL alone can't do O(1) lookup.
         Together they can."

Step 4: Identify the GLUE between structures.
        "Hash map values are pointers to DLL nodes. Each DLL node stores
         the key so we can delete from the map during eviction."

Step 5: DRAW the data flow before coding.
        "On Get: map lookup → find node → remove from DLL → reinsert at
         head → return value."
```

**DRAW BEFORE YOU CODE.** On the whiteboard or in comments. Tell the
interviewer: "Let me sketch the data flow for each operation before I
write code." This is how you demonstrate senior-level thinking. No
interviewer has ever penalized a candidate for planning too carefully.

---

## Common Interview Traps

### LRU Cache Traps

1. **Put on existing key must UPDATE and MOVE TO FRONT.**
   Most candidates remember to update the value. They forget to move the
   node. The test case that catches this:
   ```
   Put(1, 1), Put(2, 2), Get(1), Put(3, 3)
   // Evicts 2, not 1 — because Get(1) moved key 1 to front
   // If Put(1, newVal) doesn't move to front, key 1 gets wrongly evicted
   ```

2. **Eviction must delete from the hash map.**
   You remove the tail node from the DLL. You must also
   `delete(cache, victim.key)`. Without this, the map still holds a pointer
   to a dangling node, and `len(cache)` never decreases.

3. **Sentinel nodes eliminate nil checks.**
   Without sentinels, `remove` needs four nil checks (is node the head? the
   tail? the only node? etc.). With sentinels, `remove` is always the same
   two pointer swaps. There is no reason not to use them.

4. **The node must store the key.**
   During eviction you have the node but need to delete from the map by key.
   If the node doesn't store the key, you're stuck.

### Min Stack Traps

5. **Pop removes the current min context automatically.**
   No comparison needed on pop. The previous entry's `min` field was computed
   when that entry was pushed, and it is correct for all elements below it.

### Median Finder Traps

6. **Rebalance after every insert.**
   The invariant is `lo.Len() == hi.Len()` or `lo.Len() == hi.Len() + 1`.
   The 3-step pattern (push lo, move top to hi, rebalance) enforces this
   without any conditional branches on the value.

7. **Go's `container/heap`: interface methods vs package functions.**
   ```go
   // WRONG — bypasses sift-up/sift-down:
   mf.lo.Push(5)

   // CORRECT — maintains heap invariant:
   heap.Push(mf.lo, 5)
   ```
   Your type's `Push` just appends to the slice. `heap.Push` calls your
   `Push` then sifts. If you call your `Push` directly, the heap property
   breaks silently. This is Go's most confusing standard library API.

---

## Thought Process Walkthrough

### Walkthrough 1: LRU Cache (LC 146)

**Interviewer says:** "Design a data structure that follows the constraints
of a Least Recently Used cache."

**Your response, out loud:**

> "The operations are Get(key) and Put(key, value), both in O(1). On
> capacity overflow, evict the least recently used key."
>
> "O(1) lookup by key means a hash map. But I also need to track access
> order and evict from one end in O(1) — that's a doubly linked list."
>
> "The map stores key → pointer to DLL node. The DLL is ordered by recency:
> most recent at head, least recent at tail."
>
> "On Get: look up the node in the map, remove it from its current DLL
> position, reinsert at head, return value."
>
> "On Put with existing key: same as Get but also update the value."
>
> "On Put with new key: create node, insert at head, add to map. If over
> capacity, remove tail.prev (the node before the tail sentinel), delete
> its key from the map."
>
> "I'll use sentinel head and tail nodes so remove is always two pointer
> swaps with no nil checks."

Then write the code from the Pattern 1 section above. Practice until you
can write it in under 12 minutes.

**Key decision points the interviewer is evaluating:**
- Did you identify both structures and WHY you need each?
- Did you explain the glue (map values are node pointers)?
- Did you handle Put-on-existing-key correctly?
- Did you delete from the map on eviction?
- Did you store the key in the node?

---

### Walkthrough 2: Find Median from Data Stream (LC 295)

**Interviewer says:** "Design a data structure that supports adding integers
and finding the median at any point."

**Your response, out loud:**

> "AddNum needs to be efficient — O(log n). FindMedian should be O(1)."
>
> "If I keep all numbers sorted, the median is the middle element. But
> inserting into a sorted array is O(n). A heap gives O(log n) insert,
> but only exposes the min or max, not the middle."
>
> "Insight: if I split numbers into a lower half and upper half, the
> median is at the boundary. A max-heap for the lower half gives me the
> largest small number. A min-heap for the upper half gives me the
> smallest large number. The median is one or the average of both tops."
>
> "Invariant: lo has the same size as hi, or exactly one more. So median
> is either lo's top (odd count) or the average of both tops (even count)."
>
> "On AddNum I'll use a 3-step process: push to lo, move lo's top to hi,
> then if hi is larger, move hi's top back to lo. This maintains both
> the ordering and the size invariant without branching on the value."

Then write the code from the Pattern 3 section above.

**Key decision points the interviewer is evaluating:**
- Did you explain why a single heap isn't enough?
- Did you state the size invariant clearly?
- Can you implement Go's heap interface correctly (Less direction for
  max-heap vs min-heap)?
- Did you use `heap.Push`/`heap.Pop` (not the interface methods directly)?

---

## Time Targets

| Problem | Target | What it proves |
|---------|--------|----------------|
| LRU Cache | 15 min | You know the #1 design pattern cold |
| Min Stack | 8 min | You understand the pair-stack trick |
| Find Median from Data Stream | 18 min | You can implement Go heaps under pressure |

If LRU takes you more than 15 minutes, it is not ready for an interview.
This is the most commonly asked design problem. It must be automatic.

---

## Practice Plan (2 hours)

| Block | Minutes | Activity |
|-------|---------|----------|
| 1 | 5 | Read the methodology section. Internalize the 5 steps. |
| 2 | 25 | LRU Cache from memory. No references. If stuck, re-read Pattern 1, then start over. Repeat until you hit 15 min. |
| 3 | 10 | Min Stack from memory. Should be fast. |
| 4 | 25 | Median Finder from memory. Focus on the heap interface boilerplate — you need to write it without thinking. |
| 5 | 10 | Review the traps section. For each trap, can you explain WHY it's a bug? |
| 6 | 25 | LRU Cache again, cold. Time yourself. This is your benchmark. |
| 7 | 10 | Write out the methodology steps from memory. Practice saying them out loud as if explaining to an interviewer. |
| 8 | 10 | Self-assessment below. |

---

## Quick Drill: Flash Recall

Answer without looking up. Then verify.

1. In LRU Cache, what does the hash map store as its value?
2. Why does the DLL node need to store the key?
3. In the 3-step median insert, what is the order of operations?
4. What is the difference between `mf.lo.Push(x)` and `heap.Push(mf.lo, x)`?
5. In Min Stack, what happens to getMin after a pop?
6. What are sentinel nodes and why do they simplify LRU?
7. For LRU Put with an existing key, what two things must you do?

<details>
<summary>Answers</summary>

1. A pointer to the DLL node (`*Node`).
2. During eviction you have the node but need to `delete(cache, key)` from
   the map. Without the key on the node, you cannot perform the deletion.
3. Push to lo, move lo's top to hi, if hi is larger move hi's top to lo.
4. `mf.lo.Push(x)` just appends to the slice — the heap property is broken.
   `heap.Push(mf.lo, x)` appends then sifts up to restore the invariant.
5. The previous top-of-stack entry already has the correct min for all
   elements below it. No recalculation needed.
6. Dummy head and tail nodes that are never removed. They guarantee that
   `node.prev` and `node.next` are never nil during remove/insert, so you
   never need nil checks or special cases.
7. Update the value AND move the node to the front of the DLL. Forgetting
   the move-to-front is the most common bug.

</details>

---

## Self-Assessment

After completing the practice plan, score yourself honestly:

| Skill | Yes | No |
|-------|-----|----|
| I can write LRU Cache from memory in under 15 min | | |
| I can explain WHY the hash map and DLL are both needed | | |
| I can write Go's heap interface (5 methods) without looking it up | | |
| I can explain the 3-step median insert and why it works | | |
| I can list 3 common LRU Cache bugs from memory | | |
| I use the 5-step methodology naturally when approaching a new design problem | | |

**If any "No":** That specific item is your priority for tomorrow's warm-up.
Design problems are pass/fail — you either know the pattern or you don't.
There is no partial credit for a broken LRU Cache.
