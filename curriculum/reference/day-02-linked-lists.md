# Day 2 -- Linked Lists

**Date:** Week 1, Day 2
**Time Block:** 12:00 PM - 2:00 PM (2 hours)
**Goal:** Build a complete singly linked list and a doubly linked list with sentinel nodes from scratch in Go, internalize the core pointer manipulation patterns that recur in interview problems, and understand why linked lists exist alongside arrays at the hardware level.

---

## 1. Curated Learning Resources

### Interact

1. **"Linked List" -- VisuAlgo**
   https://visualgo.net/en/list
   Step through insertions, deletions, and reversals on singly and doubly linked lists. Watch how pointers rewire during each operation. Use this *before* coding to see the pointer mechanics animated -- it's far more effective than reading about them.

### Watch

2. **"Linked Lists in 4 Minutes" -- Michael Sambol**
   https://www.youtube.com/watch?v=F8AbOfQwl1c
   A concise visual recap of singly and doubly linked lists, covering insert, delete, and traversal. Good for a quick 4-minute refresher if the concepts feel stale.

### Read

3. **"Linked Lists" -- Programiz**
   https://www.programiz.com/dsa/linked-list
   Clean diagrams showing node structure, pointer wiring for insert/delete at head, tail, and middle. Covers singly, doubly, and circular variants. Useful as a 10-minute visual refresher before implementation.

4. **"Gallery of Processor Cache Effects" -- Igor Ostrovsky**
   https://igoro.com/archive/gallery-of-processor-cache-effects/
   Not linked-list specific, but essential context. Demonstrates how cache line behavior affects data structure performance. Examples 1 and 2 show why sequential array access is dramatically faster than pointer chasing -- the hardware reason linked lists lose to arrays on traversal.

5. **"Go Slices: Usage and Internals" -- Go Blog**
   https://go.dev/blog/slices-intro
   Understand how Go slices work under the hood (pointer, length, capacity, backing array). This is the contiguous-memory competitor to linked lists. Knowing both representations lets you make informed tradeoff decisions.

6. **"Memory Layout and Garbage Collection in Go"**
   https://go.dev/doc/gc-guide
   Go's garbage collector guide. Relevant because linked lists produce many small heap allocations that the GC must track. Understand the basics: Go uses a concurrent, tri-color mark-and-sweep collector. Linked list nodes are individually heap-allocated and can cause memory fragmentation.

### Reference

7. **Go Standard Library: `container/list`**
   https://pkg.go.dev/container/list
   Go's built-in doubly linked list implementation. Read the source (it's short -- ~200 lines) *after* you build your own. It uses a sentinel ring design: a single sentinel node where `head.Next` is the first element and `head.Prev` is the last. Compare this with the two-sentinel approach you'll implement.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 - 12:20 | Review & Internalize (20 min)

| Time | Activity |
|------|----------|
| 12:00 - 12:07 | Open VisuAlgo linked list visualization. Step through: (1) insert at head, (2) insert at tail, (3) delete from middle, (4) reverse. Do this for both singly and doubly linked lists. Watch where pointers break and reconnect. |
| 12:07 - 12:13 | Read the OVERVIEW.md Day 2 section. Study the complexity table. On paper, draw: a 3-node singly linked list, label head pointer, each node's `Val` and `Next` fields, and the final `nil`. Then draw a 3-node doubly linked list with sentinel head and tail nodes. |
| 12:13 - 12:18 | Read the "Core Concepts Deep Dive" section below (Section 3). Focus on: memory layout differences (array vs linked list), why sentinels eliminate edge cases, and the runner technique overview. |
| 12:18 - 12:20 | Quick mental rehearsal: write down (on paper, not code) the steps for reversing a singly linked list using prev/curr/next. If you can't do it without hesitation, you'll iron it out during implementation. |

**No code yet.** The goal is to enter implementation with a clear mental model of where pointers point at every step.

### 12:20 - 1:20 | Implement (60 min)

| Time | Activity |
|------|----------|
| 12:20 - 12:55 | **Build the singly linked list** (35 min). Start with `SNode` and `SList` structs, then implement in this order: `PushFront`, `PushBack`, `PopFront`, `Find`, `Delete`, `Reverse`, `Len`. Write a test for each method immediately after implementing it. See Section 4 for exact signatures and test cases. |
| 12:55 - 12:57 | **Break.** Stand up, stretch. You've been focused for 35 minutes. |
| 12:57 - 1:20 | **Build the doubly linked list with sentinels** (23 min). Start by initializing sentinel `head` and `tail` nodes in the constructor, wired to each other. Then implement: `PushFront`, `PushBack`, `PopFront`, `PopBack`, `Remove`, `MoveToFront`, `Len`. The sentinel pattern should make every method cleaner -- no nil checks. Write tests as you go. |

**Implementation order rationale:** The singly linked list forces you to confront all the edge cases (nil head, updating head on delete, finding the previous node). The doubly linked list with sentinels then shows you how those edge cases dissolve. You appreciate the sentinels more after struggling without them.

### 1:20 - 1:50 | Solidify (30 min)

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | **Edge case testing.** Add tests for: operations on empty list, single-element list, delete head, delete tail, reverse a 1-element list, reverse a 2-element list, PopFront/PopBack on single element. See the edge case checklist in Section 4. |
| 1:30 - 1:42 | **Pointer manipulation patterns.** Implement the five core patterns from Section 5: reverse (iterative), find middle (fast/slow), detect cycle (Floyd's), merge two sorted lists (dummy head), remove nth from end (two pointers with gap). These are standalone functions, not methods on your list struct. |
| 1:42 - 1:50 | **Pattern connection.** Glance at the LRU cache description in OVERVIEW.md Day 21. Notice it combines a hash map with a doubly linked list -- your `DList` with `Remove` and `MoveToFront` is exactly what an LRU cache needs. If time permits, sketch (don't fully implement) how you'd wire them together. |

### 1:50 - 2:00 | Recap (10 min)

| Time | Activity |
|------|----------|
| 1:50 - 1:55 | Close all references. Write from memory: the complexity of each singly linked list operation, the complexity of each doubly linked list operation, and *why* deleting a node in a singly linked list is O(n) even with a pointer to the node. |
| 1:55 - 2:00 | Write down one gotcha that tripped you up today. Examples: "I forgot to save `curr.Next` before overwriting it during reversal," or "sentinel nodes mean I never check `if head == nil`," or "Go's GC means I don't need to explicitly free removed nodes." |

---

## 3. Core Concepts Deep Dive

### Memory Layout: Array vs. Linked List

Understanding *why* arrays and linked lists have different performance characteristics requires looking at how they sit in memory.

**Array (Go slice):**
```
Memory (contiguous block):
┌─────┬─────┬─────┬─────┬─────┐
│  10 │  20 │  30 │  40 │  50 │
└─────┴─────┴─────┴─────┴─────┘
 0x00   0x08  0x10  0x18  0x20    ← addresses are sequential
```

All elements are packed next to each other. When the CPU reads element 0, the cache line (typically 64 bytes) pulls in elements 0 through ~7 as well. Iterating the array hits the cache on almost every access. This is called **spatial locality**.

**Linked list (heap allocated nodes):**
```
Memory (scattered across the heap):
  0x00: [Val:10 | Next:0x48] ──→  0x48: [Val:20 | Next:0x1A0] ──→  0x1A0: [Val:30 | Next:nil]
```

Each node is independently allocated by `new(SNode)` or `&SNode{}`. The allocator places them wherever there's free space -- they're scattered across the heap. Following `Next` pointers is **pointer chasing**: each dereference is likely a cache miss because the next node isn't in the same cache line.

**Performance implications:**

| Operation | Array | Linked List | Why |
|-----------|-------|-------------|-----|
| Sequential iteration | ~1 ns/element (cache hits) | ~5-20 ns/element (cache misses) | Cache line prefetching works for arrays, fails for pointer chasing |
| Random access by index | O(1) -- pointer arithmetic | O(n) -- must traverse | Arrays compute address: `base + index * size` |
| Insert at front | O(n) -- shift all elements | O(1) -- rewire one pointer | Arrays must move data; lists just update pointers |
| Insert at arbitrary position (given pointer) | O(n) -- shift elements after | O(1) -- rewire two pointers | Same reasoning as above |
| Memory overhead per element | 0 bytes (just the value) | 8-16 bytes (one or two pointers) | Each node carries pointer baggage |

**The practical takeaway:** Use a slice (array) by default in Go. It's faster for iteration, uses less memory, and is simpler. Reach for a linked list only when you need O(1) insert/remove at known positions (LRU cache, deques with node references) or when the data structure's semantics require it.

### Why Sentinel/Dummy Nodes Eliminate Edge Cases

A sentinel node is a dummy node that exists solely to simplify code. It holds no meaningful data. The real elements of the list sit between the sentinels.

**Without sentinels -- every operation must check for nil:**

```go
func (l *SList) PushFront(val int) {
    node := &SNode{Val: val, Next: l.head}
    l.head = node
    l.size++
}

func (l *SList) Delete(val int) bool {
    // Special case: delete the head
    if l.head != nil && l.head.Val == val {
        l.head = l.head.Next
        l.size--
        return true
    }
    // General case: find node before the target
    curr := l.head
    for curr != nil && curr.Next != nil {
        if curr.Next.Val == val {
            curr.Next = curr.Next.Next
            l.size--
            return true
        }
        curr = curr.Next
    }
    return false
}
```

The `Delete` method has **two separate code paths**: one for the head, one for everything else. This is because deleting a node requires access to the *previous* node (to rewire its `Next`), but the head has no previous node.

**With sentinels -- the head always has a predecessor:**

```go
type DList struct {
    head, tail *DNode // sentinels -- never hold real data
    size       int
}

func NewDList() *DList {
    h := &DNode{}
    t := &DNode{}
    h.Next = t
    t.Prev = h
    return &DList{head: h, tail: t}
}

func (l *DList) Remove(node *DNode) {
    // No nil checks. No special cases. Works for any node.
    node.Prev.Next = node.Next
    node.Next.Prev = node.Prev
    l.size--
}

func (l *DList) PushFront(val int) {
    node := &DNode{Val: val}
    // Insert between head sentinel and the first real node
    node.Prev = l.head
    node.Next = l.head.Next
    l.head.Next.Prev = node
    l.head.Next = node
    l.size++
}
```

`Remove` is two lines of pointer rewiring. No `if head == nil`, no `if node == head`, no `if node == tail`. The sentinels absorb all the boundary conditions. Every real node is guaranteed to have a non-nil `Prev` and `Next`.

**Cost:** Two extra node allocations per list. Negligible.

### The Runner (Two-Pointer) Technique

The "runner" technique uses two pointers traversing the list at different speeds or with a fixed gap. It appears in a surprising number of linked list problems.

**Three variants:**

1. **Fast/slow (different speeds):** Slow moves 1 step, fast moves 2 steps per iteration. When fast reaches the end, slow is at the middle. Used for: finding the middle node, detecting cycles (Floyd's), finding the start of a cycle.

2. **Fixed gap:** Advance the first pointer k steps, then advance both pointers at the same speed. When the first reaches the end, the second is k steps from the end. Used for: removing the nth node from the end.

3. **Two-list merge:** Two pointers each on a different list, advancing based on value comparison. Used for: merging sorted lists, finding the intersection point of two lists.

The key insight is that linked lists don't support random access, so any algorithm that needs to "look ahead" or "find relative position" must encode that information in the pointer movement pattern.

### Singly vs. Doubly Linked Lists in Practice

| Factor | Singly | Doubly |
|--------|--------|--------|
| Memory per node | 1 pointer (8 bytes) | 2 pointers (16 bytes) |
| Delete given node pointer | O(n) -- must find prev | O(1) -- prev is known |
| Traverse backward | Impossible | O(1) per step |
| Implementation complexity | Simpler | More pointers to maintain |
| Common uses | Stacks, simple chains, hash table chaining | LRU cache, deques, browser history, undo systems |

**Rule of thumb:** Use singly linked when you only traverse forward and only insert/delete at the head (stack behavior). Use doubly linked when you need O(1) removal of arbitrary nodes or bidirectional traversal.

### Go's Garbage Collector and Linked Lists

In C/C++, removing a node from a linked list requires explicitly freeing the memory (`free(node)` or `delete node`). Forget to free, and you leak memory. Free too early while something still points to it, and you get a use-after-free bug.

In Go, the garbage collector handles this:

```go
func (l *DList) Remove(node *DNode) {
    node.Prev.Next = node.Next
    node.Next.Prev = node.Prev
    l.size--
    // No free() needed. Once nothing references `node`, the GC reclaims it.
}
```

**However, be aware of:**

1. **GC pressure.** Each node is a separate heap allocation. A linked list with 1 million nodes means 1 million objects the GC must scan. Compared to a slice of 1 million integers (one allocation), linked lists create significantly more work for the garbage collector.

2. **Memory fragmentation.** Nodes allocated at different times end up scattered in memory. This hurts cache performance (as discussed above) and can lead to memory fragmentation -- free space exists but in small, non-contiguous chunks.

3. **Lingering references.** If removed nodes still have pointers to list elements (e.g., `node.Next` still points into the list), the GC can't collect the chain. Go's `container/list` explicitly nils out removed nodes' pointers to help the GC:
   ```go
   // From Go's container/list source:
   e.prev = nil
   e.next = nil
   e.list = nil
   ```
   This is defensive but good practice for long-lived lists with frequent removals.

---

## 4. Implementation Checklist

### A. Singly Linked List

```go
type SNode struct {
    Val  int
    Next *SNode
}

type SList struct {
    head *SNode
    size int
}

func NewSList() *SList
func (l *SList) PushFront(val int)          // Insert at head. O(1).
func (l *SList) PushBack(val int)           // Insert at tail. O(n) -- must traverse.
func (l *SList) PopFront() (int, bool)      // Remove and return head value. O(1).
func (l *SList) Find(val int) *SNode        // Return first node with val, or nil. O(n).
func (l *SList) Delete(val int) bool        // Remove first node with val. O(n).
func (l *SList) Reverse()                   // Reverse the list in place. O(n).
func (l *SList) Len() int                   // Return size. O(1).
func (l *SList) ToSlice() []int             // Collect values into a slice (for testing). O(n).
```

**Implementation notes:**
- `PushFront`: Create node, point its `Next` to current head, update head. Two operations.
- `PushBack`: Traverse to the last node (where `curr.Next == nil`), set `curr.Next` to the new node. Special case: if head is nil, the new node becomes head.
- `PopFront`: Save head's value, advance head to `head.Next`, return value. Check for empty list first.
- `Delete`: Two-path approach (head case + general case with prev pointer), or use the "pointer to pointer" trick to unify them.
- `Reverse`: Three-pointer approach: `prev`, `curr`, `next`. On each step: save `next = curr.Next`, set `curr.Next = prev`, advance `prev = curr`, advance `curr = next`. At the end, set `head = prev`.

**Tests to write:**

```go
func TestSList_PushFrontAndToSlice(t *testing.T)
// PushFront 1, 2, 3. ToSlice should return [3, 2, 1].

func TestSList_PushBack(t *testing.T)
// PushBack 1, 2, 3. ToSlice should return [1, 2, 3].

func TestSList_PopFront(t *testing.T)
// PushFront 1, 2, 3. PopFront should return 3. ToSlice should be [2, 1].

func TestSList_PopFrontEmpty(t *testing.T)
// PopFront on empty list returns (0, false).

func TestSList_Find(t *testing.T)
// PushBack 10, 20, 30. Find(20) returns node with Val 20. Find(99) returns nil.

func TestSList_Delete(t *testing.T)
// PushBack 1, 2, 3. Delete(2). ToSlice should be [1, 3]. Len should be 2.

func TestSList_DeleteHead(t *testing.T)
// PushBack 1, 2, 3. Delete(1). ToSlice should be [2, 3].

func TestSList_DeleteTail(t *testing.T)
// PushBack 1, 2, 3. Delete(3). ToSlice should be [1, 2].

func TestSList_DeleteOnlyElement(t *testing.T)
// PushFront(5). Delete(5). Len should be 0. ToSlice should be [].

func TestSList_DeleteNotFound(t *testing.T)
// PushBack 1, 2, 3. Delete(99) returns false. Len unchanged.

func TestSList_Reverse(t *testing.T)
// PushBack 1, 2, 3, 4. Reverse. ToSlice should be [4, 3, 2, 1].

func TestSList_ReverseEmpty(t *testing.T)
// Reverse on empty list. No panic. ToSlice returns [].

func TestSList_ReverseSingle(t *testing.T)
// PushFront(1). Reverse. ToSlice should be [1].

func TestSList_Len(t *testing.T)
// Verify Len after PushFront, PushBack, PopFront, Delete sequences.
```

**Edge cases to handle:**
- [ ] PushBack on empty list (head is nil)
- [ ] PopFront on empty list (return zero value + false)
- [ ] Delete the head node
- [ ] Delete the only node (list becomes empty)
- [ ] Delete a value not in the list
- [ ] Reverse an empty list (no-op, no panic)
- [ ] Reverse a single-element list (no-op)
- [ ] Reverse a two-element list (good sanity check)

---

### B. Doubly Linked List with Sentinel Nodes

```go
type DNode struct {
    Val        int
    Prev, Next *DNode
}

type DList struct {
    head, tail *DNode // sentinel nodes
    size       int
}

func NewDList() *DList                       // Init sentinels: head <-> tail.
func (l *DList) PushFront(val int) *DNode    // Insert after head sentinel. O(1). Return the node.
func (l *DList) PushBack(val int) *DNode     // Insert before tail sentinel. O(1). Return the node.
func (l *DList) PopFront() (int, bool)       // Remove first real node. O(1).
func (l *DList) PopBack() (int, bool)        // Remove last real node. O(1).
func (l *DList) Remove(node *DNode)          // Remove a specific node. O(1).
func (l *DList) MoveToFront(node *DNode)     // Remove node, re-insert at front. O(1).
func (l *DList) Len() int                    // Return size. O(1).
func (l *DList) ToSlice() []int              // Collect values front-to-back (for testing). O(n).
func (l *DList) ToSliceReverse() []int       // Collect values back-to-front (for testing). O(n).
```

**Implementation notes:**
- `NewDList`: Create head and tail sentinels. Wire `head.Next = tail` and `tail.Prev = head`. The list is "empty" but the sentinels always exist.
- `PushFront`: Create a new node. Insert it between `head` and `head.Next`. Four pointer assignments.
- `PushBack`: Create a new node. Insert it between `tail.Prev` and `tail`. Same four pointer assignments, different anchors.
- `Remove`: Set `node.Prev.Next = node.Next` and `node.Next.Prev = node.Prev`. That's it. No nil checks.
- `MoveToFront`: Call `Remove(node)`, then `PushFront` logic (or factor out an `insertAfter` helper).
- Return `*DNode` from `PushFront`/`PushBack` so callers can store the node reference for O(1) removal later (this is exactly what LRU cache needs).

**Helper method (recommended):**
```go
func (l *DList) insertAfter(after *DNode, node *DNode)
// Insert node immediately after the `after` node. Centralizes the 4-pointer wiring.
// PushFront calls insertAfter(l.head, node).
// PushBack calls insertAfter(l.tail.Prev, node).
```

**Tests to write:**

```go
func TestDList_PushFrontAndToSlice(t *testing.T)
// PushFront 1, 2, 3. ToSlice should return [3, 2, 1].

func TestDList_PushBack(t *testing.T)
// PushBack 1, 2, 3. ToSlice should return [1, 2, 3].

func TestDList_PopFront(t *testing.T)
// PushBack 1, 2, 3. PopFront returns 1. ToSlice is [2, 3].

func TestDList_PopBack(t *testing.T)
// PushBack 1, 2, 3. PopBack returns 3. ToSlice is [1, 2].

func TestDList_PopFrontEmpty(t *testing.T)
// PopFront on empty list returns (0, false).

func TestDList_PopBackEmpty(t *testing.T)
// PopBack on empty list returns (0, false).

func TestDList_Remove(t *testing.T)
// PushBack 1, 2, 3. Get the node for 2 (from PushBack return value).
// Remove that node. ToSlice is [1, 3]. Len is 2.

func TestDList_RemoveHead(t *testing.T)
// PushBack 1, 2, 3. Remove the first real node. ToSlice is [2, 3].

func TestDList_RemoveTail(t *testing.T)
// PushBack 1, 2, 3. Remove the last real node. ToSlice is [1, 2].

func TestDList_MoveToFront(t *testing.T)
// PushBack 1, 2, 3. MoveToFront the node with 3. ToSlice is [3, 1, 2].

func TestDList_MoveToFrontAlreadyFront(t *testing.T)
// PushBack 1, 2, 3. MoveToFront the node with 1. ToSlice is [1, 2, 3].

func TestDList_ToSliceReverse(t *testing.T)
// PushBack 1, 2, 3. ToSliceReverse should return [3, 2, 1].
// This verifies the Prev pointers are correct.

func TestDList_InterleavedOperations(t *testing.T)
// PushFront, PushBack, PopFront, PopBack, Remove in sequence.
// Verify ToSlice and ToSliceReverse match at each step.
```

**Edge cases to handle:**
- [ ] PopFront / PopBack on empty list (check if `head.Next == tail`)
- [ ] Remove followed by MoveToFront on a different node (no stale pointers)
- [ ] Single-element list: PopFront and PopBack should both work and leave empty list
- [ ] MoveToFront on a node that's already the front (should be a no-op -- verify with ToSlice)
- [ ] Verify backward traversal matches forward traversal (ToSliceReverse is the reverse of ToSlice)

---

## 5. Pointer Manipulation Patterns

These five patterns are the building blocks for linked list interview problems. Each is a standalone function operating on the `SNode` type defined above.

### Pattern 1: Reverse a Linked List (Iterative)

The most fundamental linked list operation. Uses three pointers: `prev`, `curr`, `next`.

```go
func ReverseList(head *SNode) *SNode {
    var prev *SNode
    curr := head
    for curr != nil {
        next := curr.Next  // 1. Save next before we overwrite it
        curr.Next = prev   // 2. Reverse the pointer
        prev = curr        // 3. Advance prev
        curr = next        // 4. Advance curr
    }
    return prev // prev is the new head
}
```

**Why it works:** At each step, we reverse one arrow. `prev` trails behind `curr` by one node. After processing all nodes, `prev` points to the last node we visited -- which is the new head.

**Common mistake:** Forgetting step 1. If you write `curr.Next = prev` before saving `curr.Next`, you lose the reference to the rest of the list.

### Pattern 2: Find the Middle Node (Fast/Slow Pointers)

```go
func FindMiddle(head *SNode) *SNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    return slow // middle node (upper middle for even-length lists)
}
```

**Why it works:** Fast moves at 2x speed. When fast reaches the end (has traversed N nodes), slow has traversed N/2 nodes -- the middle.

**For even-length lists:** `[1, 2, 3, 4]` -- when fast is at 4 (can't move 2 more), slow is at 3. This gives the upper middle. To get the lower middle (node 2), use `fast = head.Next` as the starting condition for fast -- but be consistent and know which version you need.

### Pattern 3: Detect and Find Cycle Start (Floyd's Algorithm)

**Phase 1: Detect the cycle.**
```go
func HasCycle(head *SNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            return true
        }
    }
    return false
}
```

**Phase 2: Find where the cycle begins.**
```go
func FindCycleStart(head *SNode) *SNode {
    slow, fast := head, head
    // Phase 1: detect meeting point
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            // Phase 2: find cycle start
            slow = head
            for slow != fast {
                slow = slow.Next
                fast = fast.Next // both move 1 step now
            }
            return slow // cycle start
        }
    }
    return nil // no cycle
}
```

**Why Phase 2 works:** Let the distance from head to cycle start be `a`, the distance from cycle start to meeting point be `b`, and the cycle length be `c`. At the meeting point, slow has traveled `a + b` steps and fast has traveled `a + b + k*c` steps (for some integer k). Since fast moves at 2x speed: `2(a + b) = a + b + k*c`, so `a + b = k*c`, meaning `a = k*c - b`. Starting one pointer at head and one at the meeting point, both moving 1 step: after `a` steps, the head pointer is at the cycle start, and the meeting pointer has traveled `a = k*c - b` additional steps from the meeting point (which is `b` into the cycle), putting it at `b + k*c - b = k*c` steps into the cycle -- back at the cycle start.

### Pattern 4: Merge Two Sorted Lists (Dummy Head Technique)

```go
func MergeSorted(l1, l2 *SNode) *SNode {
    dummy := &SNode{} // dummy head -- avoids special-casing the first node
    tail := dummy

    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            tail.Next = l1
            l1 = l1.Next
        } else {
            tail.Next = l2
            l2 = l2.Next
        }
        tail = tail.Next
    }

    // Attach whichever list has remaining nodes
    if l1 != nil {
        tail.Next = l1
    } else {
        tail.Next = l2
    }

    return dummy.Next // skip the dummy
}
```

**Why the dummy head:** Without it, you'd need a special case for the first node to initialize `head`. The dummy node acts as a sentinel for the result list -- you always append to `tail.Next`, and at the end, `dummy.Next` is the real head. This pattern appears in many list-building problems.

### Pattern 5: Remove Nth Node from End (Two Pointers with Gap)

```go
func RemoveNthFromEnd(head *SNode, n int) *SNode {
    dummy := &SNode{Next: head}
    first := dummy
    second := dummy

    // Advance first pointer n+1 steps ahead
    for i := 0; i <= n; i++ {
        first = first.Next
    }

    // Move both until first reaches the end
    for first != nil {
        first = first.Next
        second = second.Next
    }

    // second.Next is the node to remove
    second.Next = second.Next.Next
    return dummy.Next
}
```

**Why it works:** After advancing `first` by `n+1` steps, the gap between `first` and `second` is `n+1`. When `first` reaches nil (one past the last node), `second` is at the node *before* the target -- exactly where we need to be to rewire the pointer.

**Why the dummy:** If `n` equals the list length, the node to remove is the head. The dummy ensures `second` starts before the head, so `second.Next = second.Next.Next` works even when removing the first element.

---

## 6. Visual Diagrams

### Singly Linked List Structure

```
   head
    │
    ▼
 ┌──────────┐     ┌──────────┐     ┌──────────┐
 │ Val: 10  │     │ Val: 20  │     │ Val: 30  │
 │ Next: ───┼────►│ Next: ───┼────►│ Next: nil│
 └──────────┘     └──────────┘     └──────────┘
```

Each node holds a value and a pointer to the next node. The last node's `Next` is `nil`. The list struct holds a pointer to the head node.

### Doubly Linked List with Sentinel Nodes

```
 SENTINEL                                                    SENTINEL
  (head)                                                      (tail)
 ┌────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐     ┌────────┐
 │ Val: 0 │     │ Val: 10  │     │ Val: 20  │     │ Val: 30  │     │ Val: 0 │
 │ Prev:nil│◄───┤ Prev: ───│◄───┤ Prev: ───│◄───┤ Prev: ───│◄───┤Prev: ──│
 │ Next: ──┼───►│ Next: ───┼───►│ Next: ───┼───►│ Next: ───┼───►│Next:nil│
 └────────┘     └──────────┘     └──────────┘     └──────────┘     └────────┘

 Sentinels are always present, even when the list is empty:

 EMPTY LIST:
 ┌────────┐     ┌────────┐
 │ head   │◄───►│ tail   │
 │Next:───┼───► │Prev:───│
 └────────┘     └────────┘
```

The sentinels guarantee every real node has a non-nil `Prev` and `Next`. No edge cases for head/tail operations.

### Reversal: Step-by-Step Pointer Rewiring

Starting list: `1 -> 2 -> 3 -> nil`

```
Step 0 (initial):
  prev = nil    curr = 1    (next = undefined)

  nil    1 ──► 2 ──► 3 ──► nil
   ▲     ▲
  prev  curr

Step 1: next = curr.Next (save 2), curr.Next = prev (point 1 to nil), advance prev & curr
  nil ◄── 1    2 ──► 3 ──► nil
          ▲    ▲
         prev curr

Step 2: next = curr.Next (save 3), curr.Next = prev (point 2 to 1), advance prev & curr
  nil ◄── 1 ◄── 2    3 ──► nil
               ▲    ▲
              prev  curr

Step 3: next = curr.Next (save nil), curr.Next = prev (point 3 to 2), advance prev & curr
  nil ◄── 1 ◄── 2 ◄── 3    nil
                    ▲    ▲
                   prev  curr

curr == nil → loop ends. Return prev (node 3). New list: 3 -> 2 -> 1 -> nil.
```

### Floyd's Cycle Detection: Pointer Positions

```
List with cycle: 1 -> 2 -> 3 -> 4 -> 5 -> 3 (back to node 3)

                    ┌───────────────────┐
                    │                   │
                    ▼                   │
  1 ──► 2 ──► 3 ──► 4 ──► 5 ───────────┘

Phase 1: Detect meeting point
  Step 0: slow=1, fast=1
  Step 1: slow=2, fast=3
  Step 2: slow=3, fast=5
  Step 3: slow=4, fast=4   ← they meet! (both at node 4)

Phase 2: Find cycle start
  Reset slow to head:
  Step 0: slow=1, fast=4
  Step 1: slow=2, fast=5
  Step 2: slow=3, fast=3   ← they meet at node 3 = cycle start!

The "a = kc - b" math in action:
  a (head to cycle start) = 2 (nodes 1 → 2 → 3)
  b (cycle start to meeting point) = 1 (node 3 → 4)
  c (cycle length) = 3 (nodes 3 → 4 → 5 → 3)
  a = kc - b → 2 = 1*3 - 1 ✓
```

---

## 7. Self-Assessment

Answer these from memory at the end of the session. If you can't, that's your signal for what to revisit tomorrow.

### Question 1: Singly Linked List Deletion

**Why does deleting a node from a singly linked list require O(n) even if you have a pointer to the node?**

*Expected answer:* To remove a node, you need to set `previousNode.Next = node.Next`. In a singly linked list, each node only has a `Next` pointer -- there's no `Prev` pointer to find the predecessor. The only way to find the previous node is to traverse from the head until you find the node whose `Next` points to the target. This traversal is O(n). In a doubly linked list, `node.Prev` gives you the predecessor directly, making removal O(1).

(Note: There is a trick for singly linked lists -- copy the *next* node's value into the current node, then delete the next node. But this doesn't work if the node to delete is the tail, and it changes node identities, which can break external references.)

### Question 2: Sentinel Nodes

**What is the invariant that sentinel nodes maintain, and what happens to your code when you add them?**

*Expected answer:* The invariant is that every real node in the list is always bracketed by the sentinels: `head <-> [real nodes] <-> tail`. This means every real node has a non-nil `Prev` and non-nil `Next`. The practical consequence: insert and remove operations never need to check for nil or special-case the head/tail. A single code path handles all cases. The tradeoff is two extra node allocations per list and the discipline to never treat sentinels as data-carrying nodes (e.g., iteration should start at `head.Next` and stop at `tail`).

### Question 3: Cache Performance

**A linked list and an array both hold 1 million integers. Iterating through the array is ~5-10x faster in practice, even though both are O(n). Why?**

*Expected answer:* Array elements are stored in contiguous memory. When the CPU loads one element, the cache line (typically 64 bytes, holding ~8 ints) brings in the next several elements for free. Subsequent accesses are cache hits. Linked list nodes are scattered across the heap in non-contiguous locations. Each `Next` pointer dereference is likely a cache miss because the next node isn't in the same cache line. The CPU's hardware prefetcher can predict sequential array access patterns but cannot predict pointer-chasing patterns. The O(n) hides a constant factor that's ~5-10x larger for linked lists due to cache misses.

### Question 4: Floyd's Cycle Detection

**In Floyd's algorithm, after the slow and fast pointers meet inside a cycle, why does resetting one pointer to the head and advancing both at speed 1 lead them to meet at the cycle's start?**

*Expected answer:* Let `a` = distance from head to cycle start, `b` = distance from cycle start to meeting point, and `c` = cycle length. At the meeting point, slow traveled `a + b` and fast traveled `2(a + b)`. The difference `a + b` must be a multiple of `c` (fast lapped slow some number of times): `a + b = k*c`. Therefore `a = k*c - b`. Starting one pointer at head (needs `a` steps to reach cycle start) and one at the meeting point (which is `b` steps past cycle start): after `a` steps, the meeting-point pointer has traveled `a = k*c - b` more steps, putting it at position `b + (k*c - b) = k*c` from cycle start, which is exactly the cycle start (since `k*c` is a full number of laps). Both pointers meet at the cycle start.

### Question 5: Design Connection

**An LRU cache needs O(1) `Get(key)` and O(1) `Put(key, value)` with eviction of the least recently used entry when at capacity. Why does this require *both* a hash map and a doubly linked list? Why can't either alone suffice?**

*Expected answer:* A hash map alone gives O(1) lookup by key but has no concept of ordering -- you can't efficiently determine which entry was least recently used. A doubly linked list alone can maintain insertion/access order (most recent at front, least recent at back), and offers O(1) removal and insertion at known positions, but finding a node by key requires O(n) traversal. Combining them: the hash map maps each key to its corresponding DLL node (O(1) lookup), and the DLL maintains recency order (O(1) move-to-front on access, O(1) eviction from the back). `Get` looks up the node via the map, then calls `MoveToFront`. `Put` either updates and moves an existing node, or inserts a new node at the front and evicts from the back if over capacity.
