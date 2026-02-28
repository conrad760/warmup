# Day 5: Linked Lists

> **Goal:** Make pointer manipulation automatic under interview pressure.
> **Time budget:** 2 hours | **Language:** Go
> **Prerequisite:** You know what a linked list is. Today is about drilling the patterns until your hands type them without thinking.

```go
// Standard definition — burn this into memory
type ListNode struct {
    Val  int
    Next *ListNode
}
```

---

## Pattern Catalog

### 1. Dummy Head (Sentinel Node)

**Trigger:** Any problem where the head of the result list is unknown ahead of time, or where you might need to insert/delete at the head.

**Why it matters:** Eliminates every `if head == nil` special case. Your code becomes one uniform loop.

**Go Template:**
```go
func buildList(/* params */) *ListNode {
    dummy := &ListNode{}
    tail := dummy

    for /* condition */ {
        tail.Next = &ListNode{Val: /* value */}
        tail = tail.Next
    }

    return dummy.Next // real head
}
```

**Complexity:** O(n) time, O(1) extra space (the dummy node is constant).

**Watch out:**
- Always return `dummy.Next`, never `dummy`.
- `tail` must advance — forgetting `tail = tail.Next` creates an infinite structure of length 1.
- The dummy node's `Val` is garbage; never read it.

---

### 2. Reverse a Linked List (Three-Pointer Dance)

**Trigger:** Problem says "reverse," or you need the second half of a list reversed for reorder/palindrome checks.

**Go Template:**
```go
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode // starts as nil — becomes the new tail
    curr := head

    for curr != nil {
        next := curr.Next // 1. save next BEFORE breaking the link
        curr.Next = prev  // 2. reverse the pointer
        prev = curr       // 3. advance prev
        curr = next       // 4. advance curr
    }

    return prev // new head
}
```

**Drill the sequence aloud:** "Save next. Reverse pointer. Advance prev. Advance curr."

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- **The #1 interview bug:** forgetting to save `next` before overwriting `curr.Next`. Once you write `curr.Next = prev`, the original next is gone forever.
- `prev` starts as `nil` — this is correct because the old head's `Next` should become `nil` (it's the new tail).
- After the loop, `curr` is `nil` and `prev` points to the new head.

---

### 3. Fast-Slow Pointers

**Trigger:** "Find the middle," "detect a cycle," "find where the cycle starts," or any problem where you need to relate positions without knowing the length.

**Go Template — Find Middle:**
```go
func findMiddle(head *ListNode) *ListNode {
    slow, fast := head, head

    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }

    return slow // middle node (second middle if even length)
}
```

For splitting a list in half (needed for reorder/palindrome), you often want the node *before* the middle so you can cut the link:

```go
func findMiddlePrev(head *ListNode) *ListNode {
    slow, fast := head, head.Next // note: fast starts one ahead

    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }

    return slow // last node of first half
}
```

**Go Template — Detect Cycle:**
```go
func hasCycle(head *ListNode) bool {
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

**Go Template — Find Cycle Start (Floyd's Phase 2):**
```go
func detectCycle(head *ListNode) *ListNode {
    slow, fast := head, head

    // Phase 1: detect meeting point
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            // Phase 2: find entrance
            slow = head
            for slow != fast {
                slow = slow.Next
                fast = fast.Next // both move at speed 1 now
            }
            return slow // cycle start
        }
    }

    return nil // no cycle
}
```

**Complexity:** O(n) time, O(1) space for all variants.

**Watch out:**
- The termination condition is `fast != nil && fast.Next != nil` — both checks are needed. If you only check `fast.Next != nil`, you nil-dereference when `fast` itself is `nil` (even-length list).
- For "find middle," `slow` lands on the second middle node for even-length lists. Some problems want the first middle — adjust by starting `fast = head.Next`.
- Phase 2 of cycle detection: reset `slow` to `head`, move both pointers one step at a time. This is a mathematical proof, not intuition — just memorize it.

---

### 4. Merge Two Sorted Lists

**Trigger:** "Merge two sorted linked lists" or as a subroutine in merge sort on lists.

**Go Template:**
```go
func mergeTwoLists(l1, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
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

    return dummy.Next
}
```

**Complexity:** O(n + m) time, O(1) space (reusing existing nodes).

**Watch out:**
- Use `<=` not `<` for stability (preserves relative order of equal elements).
- Don't forget to attach the remaining tail — one list will run out first.
- This is dummy head pattern + comparison. Two patterns composed.

---

### 5. Remove Nth Node From End

**Trigger:** "Remove the nth node from the end of the list."

**Go Template:**
```go
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    dummy := &ListNode{Next: head}
    fast, slow := dummy, dummy

    // Advance fast by n+1 steps to create the gap
    for i := 0; i <= n; i++ {
        fast = fast.Next
    }

    // Move both until fast hits nil
    for fast != nil {
        slow = slow.Next
        fast = fast.Next
    }

    // slow is now the node BEFORE the one to remove
    slow.Next = slow.Next.Next

    return dummy.Next
}
```

**Complexity:** O(n) time, O(1) space, single pass.

**Watch out:**
- The dummy node is critical here. Without it, removing the head node (when n equals list length) requires a special case.
- The gap is `n+1`, not `n`. You need `slow` to land on the node *before* the target so you can unlink it.
- If the problem guarantees `n` is valid, you don't need bounds checking. In an interview, ask: "Can I assume n is always valid?"

---

### 6. Reorder / Rearrange (Split + Reverse + Interleave)

**Trigger:** "Reorder list" (L0 -> Ln -> L1 -> Ln-1 -> ...), "check if palindrome," or any problem that needs both ends of the list simultaneously.

**Go Template — Reorder List:**
```go
func reorderList(head *ListNode) {
    if head == nil || head.Next == nil {
        return
    }

    // Step 1: Find the middle (end of first half)
    slow, fast := head, head.Next
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }

    // Step 2: Split and reverse second half
    second := slow.Next
    slow.Next = nil // cut the list
    second = reverseList(second)

    // Step 3: Interleave (merge alternating)
    first := head
    for second != nil {
        tmp1 := first.Next
        tmp2 := second.Next
        first.Next = second
        second.Next = tmp1
        first = tmp1
        second = tmp2
    }
}
```

**Go Template — Palindrome Check:**
```go
func isPalindrome(head *ListNode) bool {
    if head == nil || head.Next == nil {
        return true
    }

    // Step 1: Find middle
    slow, fast := head, head.Next
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }

    // Step 2: Reverse second half
    second := reverseList(slow.Next)

    // Step 3: Compare
    first := head
    for second != nil {
        if first.Val != second.Val {
            return false
        }
        first = first.Next
        second = second.Next
    }

    return true
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- This pattern composes three sub-patterns: fast-slow (find middle), reversal, and merge/interleave. Practice each in isolation first.
- The split step (`slow.Next = nil`) is easy to forget. Without it, you have a cycle or a corrupted list.
- In the interleave step, you need to save *both* next pointers before rewiring. Two `tmp` variables, not one.
- For palindrome: if the interviewer cares about restoring the original list, reverse the second half again after comparison.

---

## Decision Framework

When you see the problem, pattern-match instantly:

| Signal in the problem | Pattern to reach for |
|---|---|
| "Reverse" or "reverse a portion" | Three-pointer reversal |
| "Middle" or "cycle" or "cycle start" | Fast-slow pointers |
| "Merge two sorted lists" | Dummy head + comparison merge |
| "Nth from end" or "remove from end" | Two pointers with n-gap |
| "Rearrange," "reorder," or "palindrome" | Split at middle + reverse second half + merge/interleave |
| Head of result unknown / insert-at-head edge case | Dummy head (sentinel) |

**Compound patterns:** Many medium/hard problems combine 2-3 of these. Reorder List = fast-slow + reverse + interleave. Recognize the sub-patterns and compose them.

**When you're stuck:** Ask yourself: "Do I need to know the length? Can fast-slow help? Would a dummy head remove my edge case? Do I need to process from both ends (split + reverse)?"

---

## Common Interview Traps

### 1. Forgetting to save `next` before overwriting `curr.Next`
```go
// BUG: curr.Next is gone after this line
curr.Next = prev
curr = curr.Next // this is now prev, not the original next!

// FIX: always save first
next := curr.Next
curr.Next = prev
curr = next
```

### 2. Wrong fast-slow termination
```go
// BUG: nil dereference when fast is nil (even-length list)
for fast.Next != nil {

// FIX: check fast first
for fast != nil && fast.Next != nil {
```
Go short-circuits `&&`, so if `fast` is nil, `fast.Next` is never evaluated.

### 3. Losing the head reference
```go
// BUG: you moved head and can't return the start of the list
for head != nil {
    head = head.Next
}
return head // this is nil

// FIX: use a separate traversal pointer
curr := head
for curr != nil {
    curr = curr.Next
}
return head // still points to the start
```

### 4. Not handling single-node or empty-list edge cases
Always add this guard at the top of your function if the problem allows empty/single-node inputs:
```go
if head == nil || head.Next == nil {
    return head
}
```
State this explicitly to the interviewer: "I'll handle the empty and single-node cases first."

### 5. Cycle detection — forgetting Phase 2
Detecting a cycle is not enough if the problem asks "where does the cycle start." You need Floyd's two-phase algorithm. Phase 1 finds the meeting point; Phase 2 resets one pointer to head and walks both at speed 1.

### 6. Off-by-one in "remove nth from end"
The dummy node fixes this. Without it, removing the first node is a special case:
```go
// Without dummy — needs special case
if n == length {
    return head.Next
}

// With dummy — uniform logic, no special case
dummy := &ListNode{Next: head}
```

---

## Thought Process Walkthrough

### Walkthrough 1: Reverse Linked List

**This must be automatic. Target: under 3 minutes including explanation.**

> **Interviewer:** "Reverse a singly linked list."

**Step 1 — Clarify (15 seconds):**
"Just to confirm — singly linked list, reverse it in place, return the new head? Can the list be empty or have one node? Great, I'll handle those."

**Step 2 — State the approach (20 seconds):**
"I'll use three pointers: prev, curr, and next. I walk through the list, reversing each pointer. prev starts at nil because the old head becomes the new tail pointing to nil."

**Step 3 — Code (90 seconds):**
```go
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    curr := head

    for curr != nil {
        next := curr.Next // save
        curr.Next = prev  // reverse
        prev = curr       // advance prev
        curr = next       // advance curr
    }

    return prev
}
```

**Step 4 — Trace through an example (30 seconds):**
"For list 1 -> 2 -> 3:
- Iteration 1: save next=2, 1->nil, prev=1, curr=2
- Iteration 2: save next=3, 2->1, prev=2, curr=3
- Iteration 3: save next=nil, 3->2, prev=3, curr=nil
- Return prev which is 3. List is now 3 -> 2 -> 1 -> nil."

**Step 5 — Complexity (10 seconds):**
"O(n) time, O(1) space. Single pass, no extra allocation."

**Step 6 — Edge cases (10 seconds):**
"Empty list: curr is nil immediately, return prev which is nil. Correct. Single node: one iteration, returns that node. Correct."

**Total: ~2.5 minutes.**

---

### Walkthrough 2: Reorder List

> **Interviewer:** "Given a linked list L0 -> L1 -> ... -> Ln, reorder it to L0 -> Ln -> L1 -> Ln-1 -> L2 -> Ln-2 -> ..."

**Step 1 — Clarify (20 seconds):**
"So I interleave from the front and back. Can the list be empty or have one node? I'll handle those as no-ops. I modify in place, not returning a new list?"

**Step 2 — Identify the pattern (30 seconds):**
"I can't efficiently access the end of a singly linked list. But if I reverse the second half, then the 'end' becomes the start of a second list, and I just interleave two lists.

Three steps:
1. Find the middle using fast-slow pointers.
2. Reverse the second half.
3. Merge/interleave the two halves."

**Step 3 — Code (3-4 minutes):**
```go
func reorderList(head *ListNode) {
    if head == nil || head.Next == nil {
        return
    }

    // Step 1: Find middle
    slow, fast := head, head.Next
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    // slow is the last node of the first half

    // Step 2: Reverse second half
    var prev *ListNode
    curr := slow.Next
    slow.Next = nil // cut the list in two
    for curr != nil {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    // prev is the head of the reversed second half

    // Step 3: Interleave
    first, second := head, prev
    for second != nil {
        tmp1 := first.Next
        tmp2 := second.Next
        first.Next = second
        second.Next = tmp1
        first = tmp1
        second = tmp2
    }
}
```

**Step 4 — Trace through an example (45 seconds):**
"For 1 -> 2 -> 3 -> 4 -> 5:
- Find middle: slow stops at 3. Cut: first half is 1->2->3, second is 4->5.
- Reverse second half: 5->4->nil.
- Interleave:
  - first=1, second=5: 1->5->2, first=2, second=4
  - first=2, second=4: 2->4->3, first=3, second=nil
  - second is nil, stop.
- Result: 1->5->2->4->3. Correct."

**Step 5 — Complexity:**
"O(n) time total — each step is O(n). O(1) space — all pointer manipulation."

**Step 6 — Edge cases:**
"Empty: returns immediately. Single node: returns. Two nodes 1->2: middle is 1, second half is 2, reversed is 2, interleave gives 1->2. Correct, no change needed."

**Total: ~6 minutes. This is a solid medium-problem pace.**

---

## Time Targets

| Task | Target Time |
|---|---|
| Reverse a linked list (cold, from memory) | < 3 min |
| Detect cycle + find cycle start | < 5 min |
| Merge two sorted lists | < 5 min |
| Remove nth from end | < 5 min |
| Reorder list (full: split + reverse + merge) | < 8 min |
| Palindrome linked list | < 8 min |

**Today's 2-hour breakdown:**

| Block | Duration | Activity |
|---|---|---|
| 1 | 15 min | Read this guide. Internalize the patterns. |
| 2 | 15 min | Code reverse list 3 times from scratch without looking. Time yourself. |
| 3 | 20 min | Code find-middle, detect-cycle, find-cycle-start from scratch. |
| 4 | 20 min | Code merge-two-sorted-lists and remove-nth-from-end. |
| 5 | 30 min | Code reorder-list and palindrome-check (these compose the sub-patterns). |
| 6 | 10 min | Quick drill (below) — do all 5 without looking at templates. |
| 7 | 10 min | Self-assessment (below). Be honest about gaps. |

---

## Quick Drill

Do these without looking at the templates above. Write each on paper or in an empty file. If you get stuck, stop, identify which sub-pattern you blanked on, re-read that one section, and try again.

**Exercise 1:** Write `reverseList` from scratch. You should be able to do this in under 2 minutes now. If you can't, repeat until you can.

**Exercise 2:** Write a function that returns `true` if a linked list has a cycle. Use fast-slow pointers.

**Exercise 3:** Write `mergeTwoLists` for two sorted linked lists. Use a dummy head.

**Exercise 4:** Write `removeNthFromEnd` in a single pass. Use the gap technique with a dummy node.

**Exercise 5:** Write `isPalindrome` for a linked list. This requires: find middle, reverse second half, compare. Do it in O(1) space.

**Grading yourself:**
- Compiled and correct on first try, no peeking: you own this pattern.
- Small bug (off-by-one, forgot to save next): review that specific trap, then redo.
- Had to look at the template: that pattern needs more reps. Schedule 10 minutes tomorrow to drill it again.

---

## Self-Assessment

Answer these without looking at the guide. If you can't answer confidently, that's your study target for tomorrow.

**Q1:** In the three-pointer reversal, what are the initial values of `prev`, `curr`, and `next`? Why does `prev` start as `nil`?

> `prev = nil` (old head becomes new tail, pointing to nil), `curr = head`, `next` is assigned inside the loop before each reversal step.

**Q2:** What is the termination condition for fast-slow pointer traversal, and why do you need *both* checks?

> `fast != nil && fast.Next != nil`. You need both because: if the list has even length, `fast` becomes `nil` (check 1 catches it). If odd length, `fast.Next` becomes `nil` (check 2 catches it). Missing either causes a nil dereference.

**Q3:** In "remove nth from end," why do you advance fast by `n+1` steps instead of `n`?

> Because you need `slow` to land on the node *before* the target. The gap of `n+1` means when `fast` hits `nil`, `slow` is one position before the node to delete, so you can do `slow.Next = slow.Next.Next`.

**Q4:** In reorder list, what happens if you forget `slow.Next = nil` after finding the middle?

> The first half still has a pointer into the second half. After you reverse the second half, you have a corrupted structure — the first half's tail still points into the (now reversed) second half, creating wrong connections or a cycle.

**Q5:** Why does Floyd's cycle detection Phase 2 work — what do you reset and how do you move?

> Reset `slow` to `head`. Move both `slow` and `fast` one step at a time. They meet at the cycle entrance. The math: if the distance from head to cycle start is `a`, and the meeting point is `b` steps into the cycle, then `a = c - b` where `c` is the cycle length. So both pointers travel the same distance `a` to reach the cycle start.

---

**End of Day 5.** If reversal isn't automatic yet, do 5 more reps tomorrow before moving on. Everything else in linked lists builds on top of it.
