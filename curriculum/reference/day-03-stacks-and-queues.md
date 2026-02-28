# Day 3 — Stacks, Queues & Deques: Deep Dive

**Date slot:** Week 1, Day 3  
**Time:** 12:00 PM – 2:00 PM (2 hours)  
**Prerequisites:** Day 1 (Hash Tables), Day 2 (Linked Lists)

---

## 1. Curated Learning Resources

| # | Title | URL | Category | Why |
|---|-------|-----|----------|-----|
| 1 | **Stacks & Queues — UC San Diego (Coursera excerpt)** | https://www.coursera.org/lecture/data-structures/stacks-UdKzQ | Watch | Clear 10-minute explanation of LIFO/FIFO semantics with call stack motivation |
| 2 | **Circular Buffer Visualization — BetterExplained** | https://betterexplained.com/articles/circular-buffers/ | Read/Interact | Best visual walkthrough of front/back pointer wrapping with modular arithmetic |
| 3 | **VisuAlgo — Stack, Queue, Deque** | https://visualgo.net/en/list | Interact | Step-through animations for push/pop/enqueue/dequeue on array and linked-list implementations |
| 4 | **Monotonic Stack — LeetCode Explore** | https://leetcode.com/explore/learn/card/queue-stack/ | Interact | Problem sets with built-in hints; covers monotonic stack pattern progressively |
| 5 | **Contiguous Stacks: How Go Manages Goroutine Stacks** | https://blog.cloudflare.com/how-stacks-are-handled-in-go/ | Read | Go-specific: explains segmented stacks (pre-1.4) vs contiguous stack copying (1.4+), dynamic growth/shrinkage |
| 6 | **Monotonic Stack Explained with Animations** | https://itnext.io/monotonic-stack-identify-pattern-3cb048f16835 | Read | Pattern-focused: next greater element, stock span, largest rectangle — with step-by-step diagrams |
| 7 | **Go Slices: Usage and Internals** | https://go.dev/blog/slices-intro | Reference | Official Go blog on slice header, capacity, append semantics — critical for array-backed stack understanding |
| 8 | **Implement Queue using Stacks — NeetCode** | https://www.youtube.com/watch?v=3Et9MrMc02A | Watch | Visual walkthrough of the two-stack queue with amortized analysis |

---

## 2. Detailed 2-Hour Session Plan

### 12:00 – 12:20 | Review (20 min)

| Time | Activity |
|------|----------|
| 12:00 – 12:05 | Read the complexity table from the overview. Internalize: all core stack/queue ops are O(1). Write it from memory. |
| 12:05 – 12:10 | Read Section 3 (Core Concepts) below. Focus on *why* naive array queue is O(n) dequeue and how the circular buffer fixes it. |
| 12:10 – 12:15 | Study the circular buffer ASCII diagram (Section 6). Trace an enqueue/dequeue sequence on paper. |
| 12:15 – 12:20 | Read the Go-specific notes: how `append` works for stacks, how goroutine stacks grow. No code yet. |

### 12:20 – 1:20 | Implement (60 min)

| Time | Activity | Notes |
|------|----------|-------|
| 12:20 – 12:35 | **Array-backed Stack** | Implement `Push`, `Pop`, `Peek`, `IsEmpty`, `Len`. Use a Go slice. Write tests for empty pop (panic or error), push/pop sequences, peek. |
| 12:35 – 12:55 | **Circular Buffer Queue** | Implement `Enqueue`, `Dequeue`, `Peek`, `IsEmpty`, `Len`, `resize`. Start with capacity 4. Track `front`, `back`, `size`, `cap`. Write tests for wrap-around, resize trigger, dequeue-all-then-enqueue. |
| 12:55 – 1:10 | **Deque (Circular Buffer)** | Implement `PushFront`, `PushBack`, `PopFront`, `PopBack`, `PeekFront`, `PeekBack`. Reuse circular buffer logic. Test interleaved front/back operations. |
| 1:10 – 1:20 | **Run all tests, fix bugs** | Aim for 100% pass. Common bugs: off-by-one in modular arithmetic, forgetting to unwrap on resize. |

**Break: 1:20 – 1:25 (5 min)** — Stand up, rest your eyes.

### 1:25 – 1:50 | Solidify (25 min)

| Time | Activity |
|------|----------|
| 1:25 – 1:35 | **Min Stack** — Implement `Push`, `Pop`, `Top`, `GetMin` all O(1). Use `(val, min)` pairs. Test with sequences that change the min multiple times. |
| 1:35 – 1:45 | **Two-Stack Queue** — Implement a queue using two stacks. Write tests proving FIFO order. Trace the amortized O(1) argument on paper. |
| 1:45 – 1:50 | **Monotonic Stack sketch** — Read the monotonic stack template in Section 5. Trace the "next greater element" example by hand on `[2, 1, 2, 4, 3]`. |

### 1:50 – 2:00 | Recap (10 min)

| Time | Activity |
|------|----------|
| 1:50 – 1:55 | Close all references. From memory, write the complexity of every operation for Stack, Queue (circular), and Deque. |
| 1:55 – 2:00 | Write down: (1) one gotcha about circular buffers, (2) when you would reach for a monotonic stack, (3) the amortized argument for dynamic array push. |

---

## 3. Core Concepts Deep Dive

### Stack: The Call Stack Analogy

A stack is LIFO — the last item pushed is the first popped. The most concrete example is the **function call stack**: when `main()` calls `foo()` calls `bar()`, each call pushes a frame. When `bar()` returns, its frame is popped, and execution resumes in `foo()`.

**How Go's goroutine stacks work:**

Go's approach to goroutine stacks has evolved:

1. **Segmented stacks (Go < 1.4):** Each goroutine started with a small stack (~4 KB). When more space was needed, a new segment was allocated and linked. Problem: "hot split" — if a function call repeatedly crossed the segment boundary, it would allocate/free segments in a tight loop, killing performance.

2. **Contiguous stacks (Go >= 1.4, current):** Each goroutine starts with a small stack (currently 2–8 KB). When the stack needs to grow, Go allocates a new, **larger** contiguous stack (typically 2x), **copies** the entire old stack to the new one, and updates all pointers. This eliminates hot-split issues at the cost of an occasional copy.

3. **Dynamic growth/shrinkage:** Stacks grow on demand (checked at function preambles via a "stack check" inserted by the compiler). They can also shrink during garbage collection if usage has dropped. This is why Go can cheaply spawn millions of goroutines — each starts tiny.

This is fundamentally different from C/pthreads where each thread gets a fixed ~1-8 MB stack allocated upfront.

### Queue: The Circular Buffer Trick

**The problem with a naive array queue:**

```
Enqueue(1), Enqueue(2), Enqueue(3):  [1, 2, 3]
Dequeue() -> 1:                       [_, 2, 3]  -- now what?
```

Option A: Shift everything left — `[2, 3, _]`. This is O(n) per dequeue.  
Option B: Just move the front pointer. But then you waste space at the front that never gets reclaimed.

**The fix — circular buffer with modular arithmetic:**

Use a fixed-size array with two indices: `front` (where the next dequeue reads) and `back` (where the next enqueue writes). When either index reaches the end of the array, it wraps to 0:

```
next_index = (current_index + 1) % capacity
```

Dequeue advances `front`; enqueue advances `back`. Both are O(1). The "wasted" space at the front gets reused when `back` wraps around.

**The full/empty ambiguity:** When `front == back`, is the buffer full or empty? Both states look the same. Solutions:
- **Track `size` separately** (cleanest — we use this approach).
- Waste one slot: buffer is full when `(back + 1) % cap == front`. Max usable capacity = cap - 1.
- Use a boolean flag.

**Resize:** When the buffer is full and you need to enqueue, allocate a new array (2x capacity), copy elements **in logical order** from `front` to `back` (wrapping around), reset `front = 0` and `back = size`. You cannot `copy()` the raw array because elements may be non-contiguous in memory.

### Deque: Double-Ended Operations

A deque (double-ended queue) supports push and pop at **both** ends in O(1). The circular buffer implementation is identical to the queue, plus:

- `PushFront`: decrement `front` (with wrap: `front = (front - 1 + cap) % cap`), write to `data[front]`.
- `PopBack`: decrement `back` (with wrap), read from `data[back]`.

**Go's slice as a poor man's deque (and why it's suboptimal):**

```go
// PushBack: O(1) amortized — fine
s = append(s, val)

// PopBack: O(1) — fine
s = s[:len(s)-1]

// PushFront: O(n) — must shift everything right
s = append([]int{val}, s...)

// PopFront: O(n) — must shift everything left
s = s[1:]  // technically O(1) but leaks memory; copy is O(n)
```

`PushFront` and `PopFront` are O(n) because all elements must be shifted. A proper circular buffer deque makes all four operations O(1) amortized.

### Amortized Analysis of Dynamic Array Resizing

When a dynamic array (Go slice) is full and you `append`, Go allocates a new backing array (roughly 2x the old capacity) and copies all elements. That single append costs O(n). But it happens rarely.

**The accounting argument:**

- Charge each `Push` operation $2 instead of $1.
- $1 pays for the actual write.
- $1 is "saved" to pay for the future copy.
- By the time the array is full (n elements), you have $n saved — exactly enough to pay for copying n elements to the new array.

**Result:** n pushes cost O(n) total work, so each push is **amortized O(1)**.

This is why Go's `append` is efficient despite occasional resizing. The Go runtime uses a growth factor of ~2x for small slices and ~1.25x for large slices (>256 elements as of Go 1.18+).

---

## 4. Implementation Checklist

### 4a. Array-Backed Stack

```go
type Stack struct {
    data []int
}

func NewStack() *Stack
func (s *Stack) Push(val int)       // append to data
func (s *Stack) Pop() (int, bool)   // remove & return top; false if empty
func (s *Stack) Peek() (int, bool)  // return top without removing; false if empty
func (s *Stack) IsEmpty() bool
func (s *Stack) Len() int
```

**Tests to write:**
- Push 5 elements, pop all — verify LIFO order
- Pop from empty stack returns `(0, false)`
- Peek doesn't change the stack
- Interleaved push/pop: Push(1), Push(2), Pop(), Push(3), Pop(), Pop() -> 2, 3, 1
- Len tracks correctly through push and pop

**Edge cases:**
- Pop/Peek on empty stack (return zero value + false, or panic — pick one convention)
- Large number of pushes to trigger multiple slice growths

### 4b. Circular Buffer Queue (with Resize)

```go
type CircularQueue struct {
    data       []int
    front      int   // index of the front element
    back       int   // index of the next write position
    size       int   // number of elements currently in the queue
    cap        int   // capacity of the backing array
}

func NewCircularQueue(initialCap int) *CircularQueue
func (q *CircularQueue) Enqueue(val int)        // add to back; resize if full
func (q *CircularQueue) Dequeue() (int, bool)   // remove from front; false if empty
func (q *CircularQueue) Peek() (int, bool)      // look at front; false if empty
func (q *CircularQueue) IsEmpty() bool
func (q *CircularQueue) Len() int
func (q *CircularQueue) resize()                // double capacity, unwrap elements
```

**Tests to write:**
- Enqueue/dequeue maintains FIFO order
- Wrap-around: start with cap 4, enqueue 3, dequeue 2, enqueue 3 more — forces `back` to wrap
- Resize trigger: enqueue beyond capacity, verify all elements still dequeue in order
- Dequeue all elements, then enqueue again — verify front/back reset correctly
- Interleaved enqueue/dequeue stress test (100+ operations)

**Edge cases:**
- Dequeue/Peek on empty queue
- Single-element queue: enqueue, peek, dequeue
- Resize when front > back (elements wrap around the array boundary)

### 4c. Deque (Circular Buffer, Double-Ended)

```go
type Deque struct {
    data       []int
    front      int
    back       int
    size       int
    cap        int
}

func NewDeque(initialCap int) *Deque
func (d *Deque) PushFront(val int)              // decrement front (wrapping), write
func (d *Deque) PushBack(val int)               // write at back, increment (wrapping)
func (d *Deque) PopFront() (int, bool)          // read at front, increment (wrapping)
func (d *Deque) PopBack() (int, bool)           // decrement back (wrapping), read
func (d *Deque) PeekFront() (int, bool)
func (d *Deque) PeekBack() (int, bool)
func (d *Deque) IsEmpty() bool
func (d *Deque) Len() int
func (d *Deque) resize()
```

**Tests to write:**
- PushFront only, then PopBack all — verify order
- PushBack only, then PopFront all — same as queue
- Interleaved: PushFront(1), PushBack(2), PushFront(3), PopFront() -> 3, PopBack() -> 2, PopFront() -> 1
- Resize from front-heavy operations (front wraps behind index 0)
- Use as a stack (PushBack/PopBack only)
- Use as a queue (PushBack/PopFront only)

**Edge cases:**
- Pop from either end on empty deque
- Single element — pop from either end should work
- Alternating front/back pushes to exercise both wrap directions

---

## 5. Advanced Patterns

### 5a. Monotonic Stack

A **monotonic stack** maintains elements in strictly increasing (or decreasing) order from bottom to top. Before pushing a new element, you pop all elements that violate the monotonic property.

**When to use:** Whenever you need to find the **next greater element**, **next smaller element**, **previous greater/smaller element**, or need to compute spans.

**The pattern (next greater element):**

For each element, we want to know the first element to its right that is larger.

1. Iterate left to right.
2. Maintain a stack of **indices** (not values) in decreasing order of their values.
3. For each new element, while the stack is non-empty and the current element is greater than the element at the stack's top index, pop the top — the current element is its "next greater."
4. Push the current index.

**Go template — Next Greater Element:**

```go
// NextGreaterElement returns an array where result[i] is the next element
// greater than nums[i] to its right, or -1 if none exists.
func NextGreaterElement(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result {
        result[i] = -1
    }

    stack := []int{} // stack of indices; values at these indices are decreasing

    for i := 0; i < n; i++ {
        // Pop all indices whose values are less than the current element.
        // The current element is the "next greater" for each of them.
        for len(stack) > 0 && nums[stack[len(stack)-1]] < nums[i] {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[top] = nums[i]
        }
        stack = append(stack, i)
    }

    return result
}
```

**Applications:**
- **Stock Span:** For each day, how many consecutive days before it had a price <= today? Use a monotonic *decreasing* stack. When you pop, the span is the distance between the current index and the new stack top.
- **Largest Rectangle in Histogram:** For each bar, find how far left and right it can extend (i.e., next smaller element on both sides). Area = height * width. Use a monotonic *increasing* stack.

```go
// LargestRectangleInHistogram finds the area of the largest rectangle
// that can be formed within the histogram bars.
func LargestRectangleInHistogram(heights []int) int {
    n := len(heights)
    stack := []int{} // indices of bars in increasing height order
    maxArea := 0

    for i := 0; i <= n; i++ {
        // Use 0 as a sentinel height for index n to flush the stack
        h := 0
        if i < n {
            h = heights[i]
        }

        for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
            height := heights[stack[len(stack)-1]]
            stack = stack[:len(stack)-1]

            width := i
            if len(stack) > 0 {
                width = i - stack[len(stack)-1] - 1
            }

            area := height * width
            if area > maxArea {
                maxArea = area
            }
        }
        stack = append(stack, i)
    }

    return maxArea
}
```

### 5b. Monotonic Deque — Sliding Window Maximum

For a sliding window of size k, maintain a deque of **indices** where the corresponding values are in decreasing order. The front of the deque is always the index of the maximum in the current window.

**The pattern:**
1. For each new element at index `i`:
   - Remove indices from the **back** of the deque while their values are <= `nums[i]` (they'll never be the max while `nums[i]` is in the window).
   - Push `i` to the back.
   - Remove the **front** if it's outside the window (`front index <= i - k`).
   - If `i >= k-1`, the front of the deque is the max for this window position.

```go
func MaxSlidingWindow(nums []int, k int) []int {
    result := []int{}
    deque := []int{} // indices, values are decreasing

    for i := 0; i < len(nums); i++ {
        // Remove from back: elements smaller than current can never be max
        for len(deque) > 0 && nums[deque[len(deque)-1]] <= nums[i] {
            deque = deque[:len(deque)-1]
        }
        deque = append(deque, i)

        // Remove from front: element is outside the window
        if deque[0] <= i-k {
            deque = deque[1:]
        }

        // Window is fully formed
        if i >= k-1 {
            result = append(result, nums[deque[0]])
        }
    }

    return result
}
```

### 5c. Stack-Based Expression Evaluation

**Infix to Postfix (Shunting Yard):**
- Numbers go directly to output.
- Operators: pop from operator stack to output while the top has >= precedence, then push current operator.
- `(` pushes to stack; `)` pops to output until `(` is found.

**RPN (Postfix) Evaluation:**
- Numbers push to operand stack.
- Operators pop two operands, compute, push result.

```go
func EvalRPN(tokens []string) int {
    stack := []int{}
    for _, tok := range tokens {
        switch tok {
        case "+", "-", "*", "/":
            b := stack[len(stack)-1]
            a := stack[len(stack)-2]
            stack = stack[:len(stack)-2]
            switch tok {
            case "+":
                stack = append(stack, a+b)
            case "-":
                stack = append(stack, a-b)
            case "*":
                stack = append(stack, a*b)
            case "/":
                stack = append(stack, a/b) // truncates toward zero in Go
            }
        default:
            num, _ := strconv.Atoi(tok)
            stack = append(stack, num)
        }
    }
    return stack[0]
}
```

### 5d. Two Stacks as a Queue

Maintain two stacks: `inStack` (for enqueue) and `outStack` (for dequeue). When `outStack` is empty and a dequeue is requested, pour all elements from `inStack` into `outStack` — this reverses the order, giving FIFO.

**Amortized O(1):** Each element is moved from `inStack` to `outStack` exactly once. So across n operations, total work is O(n), giving amortized O(1) per operation.

```go
type TwoStackQueue struct {
    inStack  []int
    outStack []int
}

func NewTwoStackQueue() *TwoStackQueue {
    return &TwoStackQueue{}
}

func (q *TwoStackQueue) Enqueue(val int) {
    q.inStack = append(q.inStack, val)
}

func (q *TwoStackQueue) Dequeue() (int, bool) {
    if len(q.outStack) == 0 {
        if len(q.inStack) == 0 {
            return 0, false
        }
        // Pour inStack into outStack (reverses order)
        for len(q.inStack) > 0 {
            top := q.inStack[len(q.inStack)-1]
            q.inStack = q.inStack[:len(q.inStack)-1]
            q.outStack = append(q.outStack, top)
        }
    }
    val := q.outStack[len(q.outStack)-1]
    q.outStack = q.outStack[:len(q.outStack)-1]
    return val, true
}

func (q *TwoStackQueue) Peek() (int, bool) {
    if len(q.outStack) == 0 {
        if len(q.inStack) == 0 {
            return 0, false
        }
        for len(q.inStack) > 0 {
            top := q.inStack[len(q.inStack)-1]
            q.inStack = q.inStack[:len(q.inStack)-1]
            q.outStack = append(q.outStack, top)
        }
    }
    return q.outStack[len(q.outStack)-1], true
}

func (q *TwoStackQueue) IsEmpty() bool {
    return len(q.inStack) == 0 && len(q.outStack) == 0
}
```

**The reverse (two queues as a stack):** Enqueue to `q1`. To pop, dequeue all but the last element from `q1` into `q2`, then dequeue the last element (that's the "top"). Swap `q1` and `q2`. This makes pop O(n) — less practical but a good exercise.

---

## 6. Visual Diagrams

### 6a. Circular Buffer — Front/Back Pointers Wrapping

Initial state: capacity = 6, empty

```
Index:    0     1     2     3     4     5
        +-----+-----+-----+-----+-----+-----+
Data:   |     |     |     |     |     |     |
        +-----+-----+-----+-----+-----+-----+
         ^
       front
       back
       size = 0
```

After Enqueue(A), Enqueue(B), Enqueue(C), Enqueue(D):

```
Index:    0     1     2     3     4     5
        +-----+-----+-----+-----+-----+-----+
Data:   |  A  |  B  |  C  |  D  |     |     |
        +-----+-----+-----+-----+-----+-----+
         ^                       ^
       front                   back
       size = 4
```

After Dequeue() -> A, Dequeue() -> B:

```
Index:    0     1     2     3     4     5
        +-----+-----+-----+-----+-----+-----+
Data:   |     |     |  C  |  D  |     |     |
        +-----+-----+-----+-----+-----+-----+
                      ^           ^
                    front       back
                    size = 2
```

After Enqueue(E), Enqueue(F), Enqueue(G):

```
Index:    0     1     2     3     4     5
        +-----+-----+-----+-----+-----+-----+
Data:   |     |     |  C  |  D  |  E  |  F  |
        +-----+-----+-----+-----+-----+-----+
                      ^                       ^
                    front                   back (wraps!)
                    size = 4

  ...then Enqueue(G):

Index:    0     1     2     3     4     5
        +-----+-----+-----+-----+-----+-----+
Data:   |  G  |     |  C  |  D  |  E  |  F  |
        +-----+-----+-----+-----+-----+-----+
               ^     ^
             back  front
             size = 5
```

**Wrap-around!** `back` went past index 5 and wrapped to 0 via `(5 + 1) % 6 = 0`, then to 1 after writing G.

After one more Enqueue(H) — buffer is full, **resize triggered**:

```
Before resize (full, size == cap == 6):
Index:    0     1     2     3     4     5
        +-----+-----+-----+-----+-----+-----+
Data:   |  G  |  H  |  C  |  D  |  E  |  F  |
        +-----+-----+-----+-----+-----+-----+
                      ^
                front AND back would collide -> size == cap -> resize!

After resize (new cap = 12, elements unwrapped in logical order):
Index:    0     1     2     3     4     5     6  ...  11
        +-----+-----+-----+-----+-----+-----+-----+-----+
Data:   |  C  |  D  |  E  |  F  |  G  |  H  |     | ... |
        +-----+-----+-----+-----+-----+-----+-----+-----+
         ^                                     ^
       front=0                               back=6
       size = 6, cap = 12
```

### 6b. Monotonic Stack — Next Greater Element

Input: `[2, 1, 2, 4, 3]`

Trace (stack holds indices; values shown for clarity):

```
i=0, val=2
  Stack: []              -> nothing to pop
  Push 0                 -> Stack: [0(2)]
  result: [-1, -1, -1, -1, -1]

i=1, val=1
  Stack: [0(2)]          -> 1 < 2, don't pop
  Push 1                 -> Stack: [0(2), 1(1)]
  result: [-1, -1, -1, -1, -1]

i=2, val=2
  Stack: [0(2), 1(1)]    -> 2 > 1, pop index 1 => result[1] = 2
  Stack: [0(2)]          -> 2 >= 2? Not strictly greater, don't pop
  Push 2                 -> Stack: [0(2), 2(2)]
  result: [-1,  2, -1, -1, -1]

i=3, val=4
  Stack: [0(2), 2(2)]    -> 4 > 2, pop index 2 => result[2] = 4
  Stack: [0(2)]          -> 4 > 2, pop index 0 => result[0] = 4
  Stack: []              -> empty, stop
  Push 3                 -> Stack: [3(4)]
  result: [ 4,  2,  4, -1, -1]

i=4, val=3
  Stack: [3(4)]          -> 3 < 4, don't pop
  Push 4                 -> Stack: [3(4), 4(3)]
  result: [ 4,  2,  4, -1, -1]

Done. Indices 3 and 4 remain on the stack -> no next greater element.

Final result: [4, 2, 4, -1, -1]
```

**Key insight:** Each element is pushed once and popped at most once. Total work across all iterations is O(n), not O(n^2).

### 6c. Two-Stack Queue — Enqueue and Dequeue

```
Operation: Enqueue(1), Enqueue(2), Enqueue(3)

  inStack:  [1, 2, 3]  (3 is on top)
  outStack: []

  Logical queue front -> 1    back -> 3

----------------------------------------------

Operation: Dequeue()

  outStack is empty -> pour inStack into outStack:

  inStack:  []
  outStack: [3, 2, 1]  (1 is on top — FIFO order!)

  Pop from outStack -> returns 1

  inStack:  []
  outStack: [3, 2]

----------------------------------------------

Operation: Enqueue(4), Enqueue(5)

  inStack:  [4, 5]     (new elements go here)
  outStack: [3, 2]     (existing elements stay)

  Logical queue: front -> 2, 3, 4, 5 -> back

----------------------------------------------

Operation: Dequeue()

  outStack is NOT empty -> pop directly -> returns 2

  inStack:  [4, 5]
  outStack: [3]

----------------------------------------------

Operation: Dequeue(), Dequeue()

  Pop outStack -> 3
  outStack now empty
  Pop outStack -> empty, pour inStack: outStack becomes [5, 4]
  Pop outStack -> 4

  inStack:  []
  outStack: [5]
```

**The invariant:** Elements in `outStack` are always older than elements in `inStack`. A pour only happens when `outStack` is empty, preserving FIFO order.

---

## 7. Self-Assessment

Answer these without looking at your notes. If you can't answer confidently, revisit that section tomorrow.

### Q1: Circular Buffer Full vs. Empty

> In a circular buffer queue that uses `front` and `back` indices (without a separate `size` field), both `front == back` when the buffer is full AND when it's empty. Name **two** strategies to distinguish full from empty, and explain the tradeoff of each.

<details>
<summary>Answer</summary>

**Strategy 1: Waste one slot.** The buffer is full when `(back + 1) % cap == front`. You can never use all `cap` slots — usable capacity is `cap - 1`. Tradeoff: wastes one slot of memory, but no extra bookkeeping variable.

**Strategy 2: Track `size` separately.** Maintain an integer `size` that increments on enqueue and decrements on dequeue. Full when `size == cap`, empty when `size == 0`. Tradeoff: uses an extra integer, but all `cap` slots are usable and the logic is simpler.

(A third option is a boolean `isFull` flag, but this is essentially a degenerate form of tracking size.)
</details>

### Q2: Amortized O(1) for Dynamic Array Push

> Explain why appending to a dynamic array (Go slice) is amortized O(1), even though individual appends can cost O(n) when a resize occurs. Use the "banker's method" or "aggregate method" to justify your answer.

<details>
<summary>Answer</summary>

**Aggregate method:** After n pushes starting from an empty array with doubling strategy, the total copy cost across all resizes is 1 + 2 + 4 + 8 + ... + n/2 + n = 2n - 1 = O(n). Each push also costs O(1) for the write itself. So total cost for n pushes is O(n) + O(n) = O(n), giving O(n)/n = **O(1) amortized per push**.

**Banker's method:** Charge $3 per push. $1 pays for the write. $2 is deposited as credit on the new element. When a resize happens (doubling from capacity k to 2k), the k/2 elements that were inserted since the last resize each have $2 of credit, providing $k total — exactly enough to pay for copying k elements. Since no resize ever goes into "debt," the $3 charge per operation is sufficient, proving amortized O(1).
</details>

### Q3: When to Use a Monotonic Stack

> You're given an array of daily temperatures and need to find, for each day, how many days you have to wait until a warmer temperature (or 0 if none). Which pattern solves this in O(n) time, and why does a hash map or sorting approach fail here?

<details>
<summary>Answer</summary>

**Monotonic stack** (decreasing from bottom to top). For each new temperature, pop all stack entries with a lower or equal temperature — the current day is the "next warmer day" for each of them, and the wait time is `current_index - popped_index`.

**Why hash map fails:** There's no key to look up — you're not searching for a specific value, you're searching for the *next* value satisfying a comparison condition at a *later position*. A hash map gives O(1) exact-match lookup, not ordered/positional queries.

**Why sorting fails:** Sorting destroys positional information. You need to know "the next greater element *to the right*," which depends on original indices. Sorting would require reconstructing positions afterward, and doesn't naturally give you the "next" relationship.

The monotonic stack works because it maintains a set of "unresolved" elements in order. Each element enters and leaves the stack exactly once, so total time is O(n).
</details>

### Q4: Circular Buffer Resize

> You have a circular buffer with capacity 8 containing elements at these positions:

```
Index: 0  1  2  3  4  5  6  7
Data:  F  G  _  _  _  C  D  E
       ^back          ^front
       size=5
```

> Draw the state after a resize to capacity 16. What would go wrong if you used a naive `copy(newData, oldData)` instead of unwrapping?

<details>
<summary>Answer</summary>

After proper unwrap + resize:
```
Index: 0  1  2  3  4  5  6  7  8  9  ...  15
Data:  C  D  E  F  G  _  _  _  _  _  ...  _
       ^front         ^back
       front=0, back=5, size=5, cap=16
```

If you used naive `copy(newData, oldData)`:
```
Index: 0  1  2  3  4  5  6  7  8  ...  15
Data:  F  G  _  _  _  C  D  E  _  ...  _
```

The problem: `front` is still 5 and `back` is still 2. The logical order is C, D, E, F, G, but the raw copy preserves the physical layout. Subsequent operations would still work *if* you keep the old `front`/`back` values, but you've wasted the opportunity to defragment. Worse, if resize code assumes `front=0` after resize (which is the convention), dequeue would return `newData[0]` = F instead of C, **breaking FIFO order**.
</details>

### Q5: Design Tradeoff

> You need a data structure that supports `Push`, `Pop`, and `GetMin` all in O(1) time. Why can't you just use a single stack with a `min` variable? What breaks, and what's the fix?

<details>
<summary>Answer</summary>

A single `min` variable fails on `Pop`: if you pop the current minimum, you don't know what the **previous** minimum was. You'd have to scan the entire stack to find the new min, which is O(n).

Example: Push(3), Push(1), Push(2). `min = 1`. Now Pop() removes 2 — min is still 1, fine. Pop() removes 1 — what's the new min? You'd need to scan and find 3, which is O(n).

**The fix:** Store `(value, currentMin)` pairs on the stack. When pushing value `v`, the pair is `(v, min(v, previousTop.min))`. When popping, the min is automatically restored from the pair below. Each pair's `min` field is the minimum of all elements at or below that position — no scanning needed.

Alternative: Use an auxiliary "min stack" that only pushes when a new minimum is encountered (and pops when the main stack pops a value equal to the min stack's top).
</details>
