# Day 3: Stacks

> **Time budget:** 2 hours | **Goal:** Pattern recognition and speed, not fundamentals
> You already know push/pop/peek. Today is about seeing a problem and *instantly* knowing which stack pattern to reach for.

---

## Pattern Catalog

### Pattern 1: Matching / Balancing

**Trigger:** "When you see... nested pairs, matching brackets, valid parentheses, matching HTML tags, or any problem involving symmetric structure."

**Go Template:**
```go
func isValid(s string) bool {
    match := map[byte]byte{')': '(', ']': '[', '}': '{'}
    stack := []byte{}
    for i := 0; i < len(s); i++ {
        if s[i] == '(' || s[i] == '[' || s[i] == '{' {
            stack = append(stack, s[i])
        } else {
            if len(stack) == 0 || stack[len(stack)-1] != match[s[i]] {
                return false
            }
            stack = stack[:len(stack)-1]
        }
    }
    return len(stack) == 0 // <-- don't forget this
}
```

**Complexity:** O(n) time, O(n) space.

**Watch out:**
- Forgetting to check `len(stack) == 0` at the end. `"((("` has no mismatches but is still invalid.
- Off-by-one when dealing with multi-character tokens (e.g., HTML `<div>` vs `</div>`).
- Input with only closing brackets — the stack-empty check before pop prevents panic.

---

### Pattern 2: Monotonic Stack (Next Greater / Smaller Element)

**Trigger:** "When you see... 'next greater element,' 'next smaller element,' 'stock span,' 'days until warmer temperature,' 'largest rectangle in histogram,' or any problem that asks for the nearest element satisfying a comparison in one direction."

**Go Template (next greater element — decreasing stack of indices):**
```go
func nextGreaterElement(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result {
        result[i] = -1 // default: no greater element
    }
    stack := []int{} // stores indices, values are decreasing
    for i := 0; i < n; i++ {
        for len(stack) > 0 && nums[i] > nums[stack[len(stack)-1]] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[idx] = nums[i]
        }
        stack = append(stack, i)
    }
    return result
}
```

**Complexity:** O(n) time (each element pushed and popped at most once), O(n) space.

**Watch out:**
- **Increasing vs. decreasing:** A *decreasing* stack (top is smallest) finds **next greater**. An *increasing* stack (top is largest) finds **next smaller**. Mnemonic: the stack maintains the *opposite* order of what you're searching for.
- **Indices vs. values:** Almost always store indices. You can always look up the value, but you can't recover the index from the value.
- **Circular arrays:** Loop `2*n` times, use `i % n` as the index.
- **Largest rectangle in histogram:** Append a sentinel `0` height at the end so every bar gets popped and processed. Without it, you'll miss bars remaining in the stack.

---

### Pattern 3: Expression Evaluation

**Trigger:** "When you see... 'evaluate expression,' 'reverse Polish notation,' 'basic calculator,' 'infix to postfix,' or any arithmetic parsing problem."

**Go Template (RPN / postfix evaluation):**
```go
func evalRPN(tokens []string) int {
    stack := []int{}
    for _, t := range tokens {
        switch t {
        case "+", "-", "*", "/":
            b := stack[len(stack)-1]; stack = stack[:len(stack)-1]
            a := stack[len(stack)-1]; stack = stack[:len(stack)-1]
            switch t {
            case "+": stack = append(stack, a+b)
            case "-": stack = append(stack, a-b)
            case "*": stack = append(stack, a*b)
            case "/": stack = append(stack, a/b) // truncates toward zero in Go
            }
        default:
            num, _ := strconv.Atoi(t)
            stack = append(stack, num)
        }
    }
    return stack[0]
}
```

**Complexity:** O(n) time, O(n) space.

**Watch out:**
- **Operand order matters:** Pop `b` first, then `a`. The operation is `a op b`, not `b op a`. This is critical for subtraction and division.
- **Infix with parentheses (Basic Calculator):** Use two stacks (one for numbers, one for operators) OR recursively handle parentheses by treating `(` as a sub-problem boundary.
- **Negative numbers in input:** `"-3"` vs `"-"` — make sure your number-vs-operator check is robust.
- **Integer division truncation:** Go truncates toward zero, which matches LeetCode's spec. Python floors. Know your language.

---

### Pattern 4: Stack as Undo / History

**Trigger:** "When you see... 'minimum element at any time,' 'undo last operation,' 'browser back/forward,' 'snapshot state,' or any problem where you need to efficiently revert to a previous state."

**Go Template (Min Stack — store value/min pairs):**
```go
type MinStack struct {
    stack []struct{ val, min int }
}

func (s *MinStack) Push(val int) {
    curMin := val
    if len(s.stack) > 0 && s.stack[len(s.stack)-1].min < val {
        curMin = s.stack[len(s.stack)-1].min
    }
    s.stack = append(s.stack, struct{ val, min int }{val, curMin})
}

func (s *MinStack) Pop()    { s.stack = s.stack[:len(s.stack)-1] }
func (s *MinStack) Top() int { return s.stack[len(s.stack)-1].val }
func (s *MinStack) GetMin() int { return s.stack[len(s.stack)-1].min }
```

**Complexity:** O(1) for all operations, O(n) space.

**Watch out:**
- When you pop, the min changes. Each stack entry must carry *its own* min context. You cannot track min in a single variable.
- The "two-stack" variant (separate min stack) only pushes to the min stack when a new min is seen. It saves space but is trickier — the pair approach above is safer in interviews.
- **Browser history:** Two stacks (back and forward). Navigating to a new page clears the forward stack. Easy to forget that detail.

---

### Pattern 5: Stack-Based DFS

**Trigger:** "When you see... 'traverse iteratively,' 'avoid recursion,' 'iterative DFS,' 'flatten nested structure,' or when the interviewer says 'now do it without recursion.'"

**Go Template (iterative DFS on a binary tree — preorder):**
```go
func iterativeDFS(root *TreeNode) []int {
    if root == nil {
        return nil
    }
    result := []int{}
    stack := []*TreeNode{root}
    for len(stack) > 0 {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        result = append(result, node.Val)
        if node.Right != nil { // push right first so left is processed first
            stack = append(stack, node.Right)
        }
        if node.Left != nil {
            stack = append(stack, node.Left)
        }
    }
    return result
}
```

**Complexity:** O(n) time, O(h) space where h is tree height (O(n) worst case for skewed tree).

**Watch out:**
- **Push order:** Push right before left so that left is on top and processed first (LIFO).
- **Inorder iterative:** Requires a different approach — push all left children first, then process, then move right. Do not try to adapt the preorder template directly.
- **Graph DFS:** Add a `visited` map. Without it you'll loop forever on cycles.
- **Nested structures (flatten list):** The stack naturally handles arbitrary nesting depth — this is the whole point vs. recursion (explicit stack vs. call stack).

---

## Decision Framework

Read the problem statement. Walk through these questions in order:

```
START
 |
 v
Does the problem involve matching or validating nested pairs?
(brackets, tags, parentheses, symmetric structure)
 |
 YES --> Pattern 1: Matching / Balancing Stack
 |
 NO
 |
 v
Does it ask for the "next greater," "next smaller," "span,"
or "first element satisfying a comparison" in one direction?
 |
 YES --> Pattern 2: Monotonic Stack
 |
 NO
 |
 v
Does it involve evaluating, parsing, or converting
arithmetic expressions?
 |
 YES --> Pattern 3: Expression Evaluation Stack
 |
 NO
 |
 v
Does it need O(1) access to historical state, undo capability,
or tracking a running aggregate (min, max) that changes with pop?
 |
 YES --> Pattern 4: Stack as Undo / History
 |
 NO
 |
 v
Does it require DFS traversal, flattening nested structures,
or converting recursion to iteration?
 |
 YES --> Pattern 5: Stack-Based DFS
 |
 NO
 |
 v
Probably not a pure stack problem.
Consider: queue, deque, heap, or two-pointer approaches.
(But double-check: sometimes a stack hides inside a larger solution.)
```

**Quick disambiguation — Monotonic stack direction:**

```
Need NEXT GREATER element?  --> Maintain a DECREASING stack (top is smallest)
Need NEXT SMALLER element?  --> Maintain an INCREASING stack (top is largest)

Why? You pop when the current element "beats" the top.
- Decreasing stack: pop when current > top  --> you found the next greater for top.
- Increasing stack: pop when current < top  --> you found the next smaller for top.
```

---

## Common Interview Traps

### 1. Monotonic Stack: Increasing vs. Decreasing Confusion
The naming is counterintuitive. "Decreasing stack" means values decrease from bottom to top (the top is the smallest). When you see a value *larger* than the top, you pop — that's how you find the "next greater." **Tip:** Don't memorize the name. Instead, remember the invariant: "I pop when the current element violates the stack's order." Then think about what that pop *means*.

### 2. Valid Parentheses: Forgetting the Final Empty Check
Your loop handles mismatches, but `"((("` never triggers a mismatch. You *must* return `len(stack) == 0` at the end. This is the #1 bug interviewers see on this problem.

### 3. Largest Rectangle in Histogram: The Sentinel Trick
After processing all bars, there may still be bars in the stack that were never popped. The standard fix: append a bar of height `0` to the input. This forces everything to pop. Without it, you need a separate cleanup loop — easy to mess up under pressure.

```go
heights = append(heights, 0) // sentinel — process before main loop
```

### 4. Min Stack: Losing Min Context on Pop
If you track the global min in a single variable, popping the min element breaks everything — you don't know the *previous* min. Each stack frame must store its own min. The pair `{value, currentMin}` approach is the safest.

### 5. Stack Underflow: Always Check Before Pop
In Go, `stack[len(stack)-1]` panics on an empty slice. **Every** pop or peek must be guarded:
```go
if len(stack) == 0 {
    // handle: return false, skip, push directly, etc.
}
```
This is especially common in matching problems (input starts with a closing bracket) and monotonic stack edge cases (first element).

### 6. Operand Order in Expression Evaluation
When you pop two operands for an operator, the *first* pop is the *right* operand. `a - b` means pop `b`, then pop `a`, then compute `a - b`. Reversing this gives wrong answers for `-` and `/` but correct answers for `+` and `*`, making the bug intermittent and hard to catch.

---

## Thought Process Walkthrough

### Problem 1: Valid Parentheses (LeetCode 20)

> Given a string containing just `(`, `)`, `{`, `}`, `[`, `]`, determine if the input string is valid.

**Step 1 — Clarify (30 seconds)**

Say out loud to the interviewer:
- "So I need to check that every opening bracket has a corresponding closing bracket of the same type, and they're properly nested — not just balanced counts."
- "Can the string be empty?" (Yes — return true.)
- "Can it contain characters other than brackets?" (No, only brackets.)
- "What's the max length?" (Determines if O(n) space is acceptable — it always is here.)

**Step 2 — Brute Force (30 seconds)**

"The brute force would be to repeatedly remove adjacent matching pairs `()`, `[]`, `{}` until no more can be removed. If the string is empty, it's valid. That's O(n^2) in the worst case."

**Step 3 — Optimize (1 minute)**

"A stack makes this O(n). When I see an opener, I push it. When I see a closer, I check the top of the stack for a match. If it doesn't match, or the stack is empty, it's invalid. At the end, the stack must be empty."

**Step 4 — Code in Go (3-4 minutes)**

```go
func isValid(s string) bool {
    match := map[byte]byte{
        ')': '(',
        ']': '[',
        '}': '{',
    }
    stack := []byte{}

    for i := 0; i < len(s); i++ {
        ch := s[i]
        if ch == '(' || ch == '[' || ch == '{' {
            stack = append(stack, ch)
        } else {
            if len(stack) == 0 {
                return false
            }
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            if top != match[ch] {
                return false
            }
        }
    }
    return len(stack) == 0
}
```

**Step 5 — Test (1-2 minutes)**

Walk through examples verbally:

- `"()"` — push `(`, see `)`, pop `(`, match. Stack empty. Return true.
- `"([)]"` — push `(`, push `[`, see `)`, top is `[`, no match. Return false.
- `"((("` — push three times. Loop ends. Stack not empty. Return false.
- `""` — loop doesn't execute. Stack empty. Return true.
- `")"` — stack empty when we see `)`. Return false.

**Step 6 — Follow-ups**

- "What if the string includes non-bracket characters?" Skip them in the loop.
- "What if we need to find the longest valid substring?" That's a different problem (DP or stack-with-indices). Mention you know the approach.
- "Can you do it with O(1) space?" Only for single bracket types using a counter. For multiple types, a stack is necessary.

---

### Problem 2: Daily Temperatures (LeetCode 739)

> Given an array of daily temperatures, return an array where `answer[i]` is the number of days you have to wait after day `i` to get a warmer temperature. If no future day is warmer, put `0`.

**Step 1 — Clarify (30 seconds)**

- "So for each element, I need to find the distance to the next *strictly* greater element to the right."
- "What are the constraints?" (Up to 10^5 elements, temperatures 30-100.)
- "If there's no warmer day, I return 0 for that position."

**Step 2 — Brute Force (30 seconds)**

"For each day, scan forward until I find a warmer temperature. O(n^2) worst case — for example, a strictly decreasing array."

**Step 3 — Optimize (1-2 minutes)**

"This is a classic 'next greater element' problem. I recognize the monotonic stack pattern."

"I'll maintain a stack of *indices* whose temperatures are in decreasing order. When I encounter a temperature *warmer* than the stack's top, I pop — the current day is the answer for that popped day. The distance is `current_index - popped_index`."

"Each index is pushed once and popped at most once, so it's O(n) total."

**Step 4 — Code in Go (3-4 minutes)**

```go
func dailyTemperatures(temperatures []int) []int {
    n := len(temperatures)
    answer := make([]int, n) // default 0 — correct for "no warmer day"
    stack := []int{}         // indices; temperatures[stack] is decreasing

    for i := 0; i < n; i++ {
        for len(stack) > 0 && temperatures[i] > temperatures[stack[len(stack)-1]] {
            j := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            answer[j] = i - j
        }
        stack = append(stack, i)
    }
    // anything remaining in stack has answer 0 (already set)
    return answer
}
```

**Step 5 — Test (1-2 minutes)**

Input: `[73, 74, 75, 71, 69, 72, 76, 73]`

```
i=0: stack=[], push 0.          stack=[0]
i=1: 74>73, pop 0, ans[0]=1.    stack=[]
     push 1.                     stack=[1]
i=2: 75>74, pop 1, ans[1]=1.    stack=[]
     push 2.                     stack=[2]
i=3: 71<75, push 3.             stack=[2,3]
i=4: 69<71, push 4.             stack=[2,3,4]
i=5: 72>69, pop 4, ans[4]=1.    stack=[2,3]
     72>71, pop 3, ans[3]=2.    stack=[2]
     72<75, push 5.             stack=[2,5]
i=6: 76>72, pop 5, ans[5]=1.    stack=[2]
     76>75, pop 2, ans[2]=4.    stack=[]
     push 6.                     stack=[6]
i=7: 73<76, push 7.             stack=[6,7]

Remaining stack [6,7]: ans[6]=0, ans[7]=0.
```

Output: `[1,1,4,2,1,1,0,0]` — matches expected.

**Step 6 — Follow-ups**

- "Can you do it in O(1) extra space?" Yes — iterate from right to left, use the answer array itself to skip ahead. Mention the approach, code if asked.
- "What if it's circular (the array wraps around)?" Iterate `2*n` times, use `i % n`.
- "What if we need the *next smaller* element instead?" Change `>` to `<` and maintain an increasing stack.

---

## Time Targets

**For a 45-minute interview:**

```
0:00 - 0:03  Read problem, clarify inputs/outputs/edge cases    [3 min]
0:03 - 0:05  State brute force approach and its complexity       [2 min]
0:05 - 0:10  Identify the stack pattern, explain optimization    [5 min]
0:10 - 0:12  Outline the algorithm in plain English              [2 min]
0:12 - 0:22  Write clean code                                   [10 min]
0:22 - 0:28  Dry-run with a concrete example                    [6 min]
0:28 - 0:32  Analyze time/space complexity                       [4 min]
0:32 - 0:38  Handle edge cases, test with tricky input           [6 min]
0:38 - 0:45  Discuss follow-ups, trade-offs, alternative approaches [7 min]
```

**Key benchmarks for stack problems specifically:**
- Pattern recognition: **under 2 minutes.** If you don't recognize it's a monotonic stack problem by minute 5, you're in trouble.
- Coding the template: **under 10 minutes.** These templates are short. Practice until you can type them without thinking.
- Testing monotonic stack: **always trace through at least 4-5 elements** to show the push/pop behavior. Interviewers want to see you understand the flow.

---

## Quick Drill

Complete each in under 2 minutes. Write Go code on paper or a blank file — no IDE autocomplete.

### Drill 1: Check If Brackets Are Balanced (Matching Stack)
**Input:** `"{[()]}"`  **Output:** `true`
**Input:** `"{[(])}"`  **Output:** `false`

Write `func isBalanced(s string) bool` using a stack and a map of closing-to-opening brackets. Don't forget the final `len(stack) == 0` check.

### Drill 2: Next Greater Element (Monotonic Stack)
**Input:** `[2, 1, 2, 4, 3]`  **Output:** `[4, 2, 4, -1, -1]`

Write `func nextGreater(nums []int) []int` using a decreasing stack of indices. Default result value is `-1`.

### Drill 3: Evaluate RPN Expression (Stack Eval)
**Input:** `["2", "1", "+", "3", "*"]`  **Output:** `9`  (i.e., `(2+1)*3`)

Write `func evalRPN(tokens []string) int`. Remember: pop `b` first, then `a`, compute `a op b`.

### Drill 4: Min Stack — O(1) getMin (Pair Stack)
Design `MinStack` with `Push`, `Pop`, `Top`, `GetMin` — all O(1).

Store `{val, currentMin}` pairs. On push, `currentMin = min(val, previous top's min)`.

### Drill 5: Reverse a String Using a Stack (Basic Stack Usage)
**Input:** `"hello"`  **Output:** `"olleh"`

Write `func reverseString(s string) string`. Push all characters, then pop them into a `strings.Builder`. This is trivial — use it as a warm-up to confirm your stack mechanics are automatic.

---

## Self-Assessment

For each scenario below, name the pattern and the core stack invariant. Answers are at the bottom — don't peek.

**Q1:** "Problem says: 'For each building, find the first taller building to its right.'"

**Q2:** "Problem says: 'Given a string of mixed brackets, determine if it is properly nested.'"

**Q3:** "Problem says: 'Implement a data structure that returns the minimum element in O(1) and supports standard stack operations.'"

**Q4:** "Problem says: 'Evaluate the expression given in postfix notation.'"

**Q5:** "Problem says: 'Given a binary tree, return its preorder traversal without using recursion.'"

---

<details>
<summary><strong>Answers</strong></summary>

**A1:** Monotonic stack (Pattern 2). Maintain a decreasing stack of indices. Pop when the current building is taller than the top.

**A2:** Matching / Balancing stack (Pattern 1). Push openers, pop and match on closers. Verify stack is empty at end.

**A3:** Stack as undo/history (Pattern 4). Store `{value, currentMin}` pairs so each frame knows its own min.

**A4:** Expression evaluation (Pattern 3). Push operands, pop two and apply operator when you encounter an operator.

**A5:** Stack-based DFS (Pattern 5). Push root, pop and process, push right then left children. Explicit stack replaces the call stack.

</details>
