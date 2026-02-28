# Day 20: Mixed Pattern Practice — Session 2

> **Time:** 2 hours | **Format:** Timed mock interview | **Language:** Go
>
> This is your second mock session. The problems are drawn from different
> topics than Day 19 and are slightly harder. The goal: prove that your
> pattern recognition transfers to unfamiliar problem statements under
> time pressure. Your Day 20 score minus your Day 19 score is your
> **rate of improvement.**

---

## Session Format

**4 problems. 25 minutes each. 100 minutes total. 20 minutes for review and reflection.**

For each problem, follow this exact sequence:

```
0:00 - 2:00   Read the problem. Clarify edge cases out loud.
2:00 - 5:00   Classify the pattern. Decide on your approach.
5:00 - 8:00   Explain the approach out loud (as if to an interviewer).
8:00 - 20:00  Code the solution in Go.
20:00 - 25:00 Test with examples. Fix bugs. State time/space complexity.
```

**Rules — simulate real interview conditions:**

1. **No notes, no references.** Close every browser tab. Close this guide
   after reading each problem statement. Reopen only to check hints or
   the solution after your attempt.
2. **Time yourself.** Use a phone timer. When 25 minutes is up, stop. A
   partial solution is a valid data point — it tells you what to drill.
3. **Talk out loud.** Every thought goes through your mouth. If you're
   silent for more than 30 seconds, you're not practicing the right skill.
   In a real interview, silence is a red flag.
4. **Write on a single file or blank editor.** No autocomplete, no LSP,
   no running the code until you've finished writing. You won't have
   these in a whiteboard interview.

**Score each problem immediately after attempting it** using the rubric at
the bottom. Be brutally honest — inflated scores help no one.

---

## Interview Time Management Tips

This is your second mock. You now have data from Day 19 about how you
allocate time. Here are the three most important time management skills
for real interviews:

### 1. When to Abandon an Approach (the 10-Minute Rule)

If you've been coding an approach for more than 10 minutes without
meaningful progress — meaning you can't see a clear path to a working
solution — stop and pivot. Here's how:

- **Say it out loud:** "I've been going down this path for a while and
  I'm not making progress. Let me step back and reconsider."
- **Don't scrap everything silently.** The interviewer wants to see your
  reasoning. Explain *why* the current approach is failing: "This brute
  force is O(n³), which won't pass the constraint. Let me think about
  what structure I'm not exploiting."
- **Return to pattern recognition.** Re-read the problem constraints.
  The constraint itself often encodes the expected complexity: n ≤ 10⁵
  means O(n log n) or better; n ≤ 20 means exponential is fine.

### 2. How to Communicate Uncertainty

In a real interview, you will get stuck. The difference between a hire
and a no-hire is often how you handle being stuck:

- **Never go silent.** Say "I'm thinking about..." and narrate your
  thought process. Even wrong ideas demonstrate problem-solving ability.
- **Ask directed questions.** "I'm considering whether a monotonic stack
  would work here — does the ordering property matter?" is better than
  "I don't know what to do."
- **Propose alternatives explicitly.** "I see two approaches: a heap-based
  solution that's O(n log k), or a quickselect that's O(n) average. I'm
  going to go with the heap because it's easier to get right under time
  pressure. Would you like me to discuss the tradeoffs first?"

### 3. State Brute Force First — Even When You Know the Optimal

Always start by stating the brute force solution, even if you immediately
plan to optimize. Here's why:

- It proves you understand the problem. Many candidates jump to an
  "optimal" solution that actually solves a different problem.
- It establishes a baseline complexity that you're improving on. "The
  naive approach checks all pairs in O(n²). I can do better with a
  hash map because..."
- If you run out of time on the optimal solution, you can fall back to
  a working brute force. A correct O(n²) solution beats an incomplete
  O(n) solution every time.
- It gives the interviewer a chance to nudge you: "That's right. Can you
  think of a way to avoid the inner loop?" A good interviewer *wants*
  you to succeed.

---

## Problem 1: Minimum Window Containing All Characters

**Difficulty:** Medium
**Time limit:** 25 minutes
**Topic area:** Sliding Window (Day 2)

### Problem Statement

Given two strings `s` and `t`, return the **minimum-length substring** of
`s` that contains every character of `t` (including duplicates). If no
such substring exists, return the empty string `""`.

If there are multiple substrings of the same minimum length, return the
one that appears first (leftmost).

**Examples:**

```
Input:  s = "ADOBECODEBANC", t = "ABC"
Output: "BANC"
Explanation: The substring "BANC" (index 9-12) contains 'A', 'B', and 'C'.
             "ADOBEC" (index 0-5) also works but is longer.

Input:  s = "a", t = "a"
Output: "a"

Input:  s = "a", t = "aa"
Output: ""
Explanation: t requires two 'a's but s only has one.

Input:  s = "cabwefgewcwaefgcf", t = "cae"
Output: "cwae"
Explanation: Substring at index 8-11.

Input:  s = "bba", t = "ab"
Output: "ba"
```

**Constraints:**
- `1 <= len(s), len(t) <= 10^5`
- `s` and `t` consist of uppercase and lowercase English letters.

Your algorithm must run in **O(n)** time where n = len(s).

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This is a **variable-length sliding window** problem. The trigger phrases:
"minimum-length substring" + "contains all characters." You need a window
that expands to include all required characters, then contracts to find the
minimum.

Maps to: **Day 2 — Two Pointers & Sliding Window, Pattern 3** (variable
window / shrinkable window).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

1. Build a frequency map `need` from `t` — how many of each character are
   required.
2. Track how many characters are **fully satisfied** with a counter `formed`.
   The target is `len(need)` (number of unique characters in t).
3. Use two pointers `left` and `right`:
   - Expand `right`: add `s[right]` to a `window` frequency map. If the
     count of `s[right]` in `window` equals the count in `need`, increment
     `formed`.
   - When `formed == len(need)`, the window is valid. Try to **shrink from
     the left**: record the window if it's the smallest so far, then remove
     `s[left]` from the window map and advance `left`. If removing it drops
     below the required count, decrement `formed`.
4. Return the smallest window recorded.

The key insight: the `formed` counter lets you check window validity in O(1)
instead of comparing entire frequency maps.

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

```go
func minWindow(s string, t string) string {
    if len(s) < len(t) {
        return ""
    }

    // Frequency map of characters needed from t
    need := make(map[byte]int)
    for i := 0; i < len(t); i++ {
        need[t[i]]++
    }

    window := make(map[byte]int)    // current window character counts
    formed := 0                      // # of unique chars in t fully satisfied
    required := len(need)            // # of unique chars in t

    // Result tracking
    bestLen := len(s) + 1
    bestStart := 0

    left := 0
    for right := 0; right < len(s); right++ {
        // Expand: add s[right] to window
        ch := s[right]
        window[ch]++

        // Check if this character is now fully satisfied
        if cnt, ok := need[ch]; ok && window[ch] == cnt {
            formed++
        }

        // Contract: shrink from the left while window is valid
        for formed == required {
            // Update best if this window is smaller
            windowLen := right - left + 1
            if windowLen < bestLen {
                bestLen = windowLen
                bestStart = left
            }

            // Remove s[left] from window
            leftCh := s[left]
            window[leftCh]--
            if cnt, ok := need[leftCh]; ok && window[leftCh] < cnt {
                formed--
            }
            left++
        }
    }

    if bestLen > len(s) {
        return ""
    }
    return s[bestStart : bestStart+bestLen]
}
```

**Complexity:**
- Time: O(n) where n = len(s). Each character is visited at most twice —
  once by `right`, once by `left`. Map operations are O(1) for fixed
  alphabet sizes.
- Space: O(k) where k = number of unique characters in s and t. With
  English letters, k ≤ 52.

**Trace through example 1:** `s = "ADOBECODEBANC", t = "ABC"`
```
need = {A:1, B:1, C:1}, required = 3

right=0 (A): window={A:1}, formed=1
right=1 (D): window={A:1,D:1}, formed=1
right=2 (O): formed=1
right=3 (B): window={...,B:1}, formed=2
right=4 (E): formed=2
right=5 (C): window={...,C:1}, formed=3 → VALID
  shrink: window=ADOBEC (len=6), record [0,5]
  remove A: formed=2, left=1 → stop shrinking

right=9 (A):  formed=3 → VALID
right=10 (N): still valid
  shrink: window=ODEBANC → still valid → ODEBAN? no, check...
  Eventually: window=BANC (len=4), record [9,12]
  remove B: formed=2, left=10 → stop

bestLen=4, bestStart=9 → "BANC"
```

**Edge cases handled:**
- `len(s) < len(t)` → impossible, return "".
- `t` has duplicate characters → `need` tracks counts, not just presence.
- No valid window exists → `bestLen` stays at `len(s)+1`, return "".

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Pattern** | Variable-length sliding window (shrinkable) |
| **Day** | Day 2 — Two Pointers & Sliding Window, Pattern 3 |
| **Trigger** | "Minimum-length substring" + "contains all" = shrinkable window. The optimization target (minimum) means you expand to satisfy the constraint, then shrink to optimize. |
| **Key insight** | The `formed` counter tracks how many unique characters are fully satisfied, enabling O(1) validity checks instead of comparing full frequency maps on every step. |
| **Common mistakes** | (1) Using `window[ch] >= need[ch]` to increment `formed` — this increments `formed` multiple times for the same character. Must use `==` (exact match). (2) Forgetting that `t` can have duplicate characters, so presence isn't enough — you need counts. (3) Off-by-one in the shrink loop — shrinking should happen while valid, not after. |

---

## Problem 2: Find the Smallest Divisor Given a Threshold

**Difficulty:** Medium
**Time limit:** 25 minutes
**Topic area:** Binary Search on Answer (Day 4)

### Problem Statement

Given an array of positive integers `nums` and a positive integer
`threshold`, find the **smallest positive integer divisor** such that the
sum of all division results is less than or equal to `threshold`.

Each division result is **rounded up** to the nearest integer. Formally,
for divisor `d`, the result for element `nums[i]` is `ceil(nums[i] / d)`.

**Examples:**

```
Input:  nums = [1, 2, 5, 9], threshold = 6
Output: 5
Explanation:
  d=5: ceil(1/5) + ceil(2/5) + ceil(5/5) + ceil(9/5) = 1+1+1+2 = 5 ≤ 6 ✓
  d=4: ceil(1/4) + ceil(2/4) + ceil(5/4) + ceil(9/4) = 1+1+2+3 = 7 > 6 ✗
  The smallest valid divisor is 5.

Input:  nums = [44, 22, 33, 11, 1], threshold = 5
Output: 44
Explanation: Each element divided by 44 and ceiled → 1+1+1+1+1 = 5 ≤ 5 ✓

Input:  nums = [21212, 10101, 12121], threshold = 1000000
Output: 1
Explanation: d=1 means each element is itself. Sum = 43434 ≤ 1000000.

Input:  nums = [2, 3, 5, 7, 11], threshold = 11
Output: 3
Explanation:
  d=3: ceil(2/3) + ceil(3/3) + ceil(5/3) + ceil(7/3) + ceil(11/3)
     = 1 + 1 + 2 + 3 + 4 = 11 ≤ 11 ✓
  d=2: 1+2+3+4+6 = 16 > 11 ✗
```

**Constraints:**
- `1 <= len(nums) <= 5 × 10^4`
- `1 <= nums[i] <= 10^6`
- `len(nums) <= threshold <= 10^6`
- The answer is guaranteed to exist (threshold ≥ len(nums) ensures d=max(nums)
  always works since every ceil becomes 1).

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This is **binary search on the answer** (also called "binary search on
the solution space"). The trigger: you're searching for the **smallest
value** that satisfies a condition, and the condition is **monotonic** —
if divisor `d` works, then `d+1` also works (larger divisors produce
smaller sums).

Maps to: **Day 4 — Binary Search, Pattern 3** (search on answer /
minimize a value subject to a feasibility check).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

1. **Search space:** The divisor ranges from 1 to max(nums).
   - d=1 gives the largest possible sum (sum of all elements).
   - d=max(nums) gives the smallest possible sum (all 1s, which is len(nums)).
2. **Feasibility function:** Given a divisor `d`, compute the sum of
   `ceil(nums[i] / d)` for all i. Check if sum ≤ threshold.
3. **Binary search:** Find the smallest `d` where feasible(d) is true.
   - If feasible(mid), search left: `hi = mid`.
   - If not feasible(mid), search right: `lo = mid + 1`.
4. **Ceiling division without floating point:** `ceil(a / b) = (a + b - 1) / b`
   using integer arithmetic.

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

```go
func smallestDivisor(nums []int, threshold int) int {
    // Find max element (upper bound of search space)
    maxVal := 0
    for _, n := range nums {
        if n > maxVal {
            maxVal = n
        }
    }

    // Binary search on the divisor
    lo, hi := 1, maxVal

    for lo < hi {
        mid := lo + (hi-lo)/2

        if feasible(nums, mid, threshold) {
            hi = mid // mid works — try smaller
        } else {
            lo = mid + 1 // mid too small — need bigger divisor
        }
    }

    return lo
}

func feasible(nums []int, divisor, threshold int) bool {
    sum := 0
    for _, n := range nums {
        sum += (n + divisor - 1) / divisor // ceil division
        if sum > threshold {
            return false // early exit optimization
        }
    }
    return true
}
```

**Complexity:**
- Time: O(n × log(max(nums))). Binary search does O(log(max(nums)))
  iterations. Each iteration computes the sum in O(n). With max up to 10⁶,
  that's about 20 × 50,000 = 1,000,000 operations.
- Space: O(1). No extra data structures.

**Trace through example 1:** `nums = [1,2,5,9], threshold = 6`
```
maxVal = 9. Search space: [1, 9].

lo=1, hi=9:
  mid=5: sum = ceil(1/5)+ceil(2/5)+ceil(5/5)+ceil(9/5) = 1+1+1+2 = 5 ≤ 6 → hi=5
lo=1, hi=5:
  mid=3: sum = 1+1+2+3 = 7 > 6 → lo=4
lo=4, hi=5:
  mid=4: sum = 1+1+2+3 = 7 > 6 → lo=5
lo=5, hi=5: done → return 5
```

**Edge cases handled:**
- All elements equal 1 → any divisor works, returns 1.
- Single element → search space is [1, nums[0]].
- Threshold equals len(nums) → divisor must be max(nums).
- Early exit in feasibility check avoids integer overflow on large inputs.

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Pattern** | Binary search on the answer (minimize a value subject to a monotonic feasibility check) |
| **Day** | Day 4 — Binary Search, Pattern 3 |
| **Trigger** | "Find the smallest X such that [condition]" + the condition is monotonic (larger divisor → smaller sum). Whenever you're optimizing a value and can check feasibility, binary search on the answer applies. |
| **Key insight** | Ceiling division without floats: `ceil(a/b) = (a + b - 1) / b`. This avoids floating-point precision issues entirely. Also: early exit in the feasibility check when the sum already exceeds the threshold. |
| **Common mistakes** | (1) Using floating-point division for ceiling — introduces precision errors on large inputs. (2) Searching [1, sum(nums)] instead of [1, max(nums)] — correct but wastefully large search space. (3) Forgetting the `lo < hi` termination — using `lo <= hi` with `hi = mid` causes infinite loops. (4) Not recognizing the monotonic property that makes binary search valid. |

---

## Problem 3: Course Schedule — Ordering

**Difficulty:** Medium
**Time limit:** 25 minutes
**Topic area:** Graph — Topological Sort (Day 9)

### Problem Statement

There are a total of `numCourses` courses you must take, labeled from `0`
to `numCourses - 1`. Some courses have prerequisites: if `prerequisites[i]
= [a, b]`, you must take course `b` before course `a`.

Return **a valid ordering** of courses such that all prerequisites are
satisfied. If there are multiple valid orderings, return any one. If it is
**impossible** to finish all courses (there is a cycle), return an empty
array.

**Examples:**

```
Input:  numCourses = 4, prerequisites = [[1,0],[2,0],[3,1],[3,2]]
Output: [0, 1, 2, 3]  (or [0, 2, 1, 3] — both valid)
Explanation:
  0 has no prerequisites.
  1 requires 0. 2 requires 0.
  3 requires both 1 and 2.
  So 0 must come first, 3 must come last.

Input:  numCourses = 2, prerequisites = [[1,0]]
Output: [0, 1]

Input:  numCourses = 2, prerequisites = [[1,0],[0,1]]
Output: []
Explanation: Cycle: 0→1→0. Impossible.

Input:  numCourses = 1, prerequisites = []
Output: [0]

Input:  numCourses = 3, prerequisites = []
Output: [0, 1, 2]
Explanation: No dependencies — any ordering works.

Input:  numCourses = 5, prerequisites = [[1,0],[2,0],[3,1],[4,2],[4,3]]
Output: [0, 1, 2, 3, 4]  (or [0, 2, 1, 3, 4], etc.)
Explanation: Multiple valid topological orderings exist.
```

**Constraints:**
- `1 <= numCourses <= 2000`
- `0 <= len(prerequisites) <= 5000`
- `prerequisites[i].length == 2`
- `0 <= a, b < numCourses`
- `a != b` (no self-loops, but cycles of length ≥ 2 are possible)
- All pairs `[a, b]` are unique.

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This is a **topological sort** problem. The trigger: "ordering" +
"prerequisites" (dependencies) + "detect if impossible" (cycle detection).
Prerequisites define directed edges, and a valid ordering is a topological
ordering of the resulting DAG.

Maps to: **Day 9 — Graphs, Pattern 3** (Topological Sort / Kahn's BFS).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

**Kahn's Algorithm (BFS-based topological sort):**

1. Build an adjacency list and an **in-degree** array.
   - For each `[a, b]`: add edge `b → a` to the adjacency list. Increment
     `inDegree[a]`.
2. Initialize a queue with all nodes that have in-degree 0 (no prerequisites).
3. BFS loop:
   - Dequeue a node, append it to the result.
   - For each neighbor: decrement its in-degree. If in-degree becomes 0,
     enqueue it.
4. After the loop, if `len(result) == numCourses`, the ordering is valid.
   Otherwise, a cycle exists — return empty.

**Why this works:** Nodes with in-degree 0 have no unprocessed dependencies.
Processing them "removes" their outgoing edges, potentially freeing other
nodes. If a cycle exists, no node in the cycle ever reaches in-degree 0,
so they're never processed.

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

```go
func findOrder(numCourses int, prerequisites [][]int) []int {
    // Build adjacency list and in-degree array
    adj := make([][]int, numCourses)
    inDegree := make([]int, numCourses)

    for _, prereq := range prerequisites {
        course, pre := prereq[0], prereq[1]
        adj[pre] = append(adj[pre], course)
        inDegree[course]++
    }

    // Initialize queue with all zero in-degree nodes
    queue := make([]int, 0)
    for i := 0; i < numCourses; i++ {
        if inDegree[i] == 0 {
            queue = append(queue, i)
        }
    }

    // BFS: process nodes in topological order
    result := make([]int, 0, numCourses)

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        result = append(result, node)

        for _, neighbor := range adj[node] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    // Cycle detection: if not all nodes processed, a cycle exists
    if len(result) != numCourses {
        return []int{}
    }

    return result
}
```

**Alternative DFS-based topological sort:**

```go
func findOrder(numCourses int, prerequisites [][]int) []int {
    adj := make([][]int, numCourses)
    for _, prereq := range prerequisites {
        course, pre := prereq[0], prereq[1]
        adj[pre] = append(adj[pre], course)
    }

    // 0 = unvisited, 1 = in current path, 2 = fully processed
    state := make([]int, numCourses)
    result := make([]int, 0, numCourses)
    hasCycle := false

    var dfs func(node int)
    dfs = func(node int) {
        if hasCycle {
            return
        }
        state[node] = 1 // mark as in-progress

        for _, neighbor := range adj[node] {
            if state[neighbor] == 1 {
                hasCycle = true
                return
            }
            if state[neighbor] == 0 {
                dfs(neighbor)
            }
        }

        state[node] = 2 // mark as done
        result = append(result, node) // post-order
    }

    for i := 0; i < numCourses; i++ {
        if state[i] == 0 {
            dfs(i)
        }
    }

    if hasCycle {
        return []int{}
    }

    // DFS post-order gives reverse topological order — reverse it
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    return result
}
```

**Complexity (both approaches):**
- Time: O(V + E) where V = numCourses, E = len(prerequisites).
  Every node and edge is processed exactly once.
- Space: O(V + E) for the adjacency list. O(V) for in-degree array / state
  array. O(V) for the queue/stack.

**Trace through example 1:** `numCourses=4, prereqs=[[1,0],[2,0],[3,1],[3,2]]`
```
adj: 0→[1,2], 1→[3], 2→[3]
inDegree: [0, 1, 1, 2]

Queue starts with: [0]  (only node with inDegree 0)

Process 0: result=[0]. Neighbors 1,2:
  inDegree[1]=0 → enqueue. inDegree[2]=0 → enqueue.
  Queue: [1, 2]

Process 1: result=[0,1]. Neighbor 3:
  inDegree[3]=1. Queue: [2]

Process 2: result=[0,1,2]. Neighbor 3:
  inDegree[3]=0 → enqueue. Queue: [3]

Process 3: result=[0,1,2,3]. No neighbors.

len(result)=4 == numCourses → valid. Return [0,1,2,3].
```

**Edge cases handled:**
- No prerequisites → all in-degrees are 0, all nodes enqueue immediately.
  Result is simply [0, 1, ..., n-1].
- Cycle → nodes in the cycle never reach in-degree 0. `len(result) <
  numCourses` triggers the empty return.
- Single course → in-degree 0, immediately returned as [0].
- Disconnected components → each component's zero-in-degree nodes are
  enqueued independently.

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Pattern** | Topological sort (Kahn's BFS) |
| **Day** | Day 9 — Graphs, Pattern 3 |
| **Trigger** | "Ordering" + "prerequisites" (dependency relationships) + "is it possible" (cycle detection). Whenever tasks have dependency constraints and you need a valid execution order, think topological sort. |
| **Key insight** | In-degree tracking is the engine: nodes with in-degree 0 are "ready." Processing a node removes its outgoing edges. Cycle detection falls out naturally — if the result is shorter than expected, a cycle exists. |
| **Common mistakes** | (1) Edge direction confusion: `[a, b]` means b→a (b must come before a), not a→b. Getting this backwards produces a valid topological sort of the *wrong* graph. (2) DFS approach: forgetting to reverse the post-order result. (3) DFS cycle detection: using a boolean `visited` instead of a three-state system (unvisited / in-progress / done) — a boolean can't distinguish back edges from cross edges. |

---

## Problem 4: Meeting Rooms — Minimum Rooms Required with Allocation

**Difficulty:** Hard
**Time limit:** 25 minutes
**Topic area:** Intervals + Heap (Day 15 + Day 5)

### Problem Statement

You are given an array of meeting time intervals `intervals` where
`intervals[i] = [start_i, end_i]` represents a meeting from time
`start_i` to `end_i` (exclusive — a meeting ending at time 10 does not
conflict with one starting at time 10).

Return the **minimum number of conference rooms** required to hold all
meetings.

Additionally, return an **allocation array** where `allocation[i]` is the
room number (0-indexed) assigned to the i-th meeting (in the original input
order). Rooms should be **reused greedily** — when a room becomes free,
it should be the one assigned to the next meeting that needs it (specifically,
if multiple rooms are free, assign the one with the **lowest room number**).

**Examples:**

```
Input:  intervals = [[0,30],[5,10],[15,20]]
Output: rooms = 2, allocation = [0, 1, 1]
Explanation:
  Meeting 0 [0,30]:  Room 0 (only room, it's free).
  Meeting 1 [5,10]:  Room 0 is busy until 30. Allocate Room 1.
  Meeting 2 [15,20]: Room 0 busy until 30. Room 1 free (ended at 10). Reuse Room 1.

Input:  intervals = [[7,10],[2,4]]
Output: rooms = 1, allocation = [0, 0]
Explanation: Sort by start time → [2,4] then [7,10]. Both fit in Room 0.
  But allocation is in ORIGINAL order: meeting 0 is [7,10] → Room 0,
  meeting 1 is [2,4] → Room 0.

Input:  intervals = [[0,5],[5,10],[10,15]]
Output: rooms = 1, allocation = [0, 0, 0]
Explanation: No overlaps (end time = start time is not a conflict).

Input:  intervals = [[1,5],[2,6],[3,7],[4,8]]
Output: rooms = 4, allocation = [0, 1, 2, 3]
Explanation: All meetings overlap with each other. Need 4 rooms.

Input:  intervals = [[1,3],[2,4],[5,7],[6,8]]
Output: rooms = 2, allocation = [0, 1, 0, 1]
Explanation:
  Sorted: [1,3],[2,4],[5,7],[6,8].
  [1,3] → Room 0.
  [2,4] → Room 0 busy. Room 1.
  [5,7] → Room 0 free (ended at 3). Room 0.
  [6,8] → Room 0 busy. Room 1 free (ended at 4). Room 1.
  Original order indices: [1,3] was index 0, [2,4] was index 1,
  [5,7] was index 2, [6,8] was index 3. → [0, 1, 0, 1].

Input:  intervals = []
Output: rooms = 0, allocation = []
```

**Constraints:**
- `0 <= len(intervals) <= 10^4`
- `0 <= start_i < end_i <= 10^6`

---

<details>
<summary><strong>Category Hint</strong> (only look if stuck after 5 minutes)</summary>

This is an **intervals + min-heap** problem. The base version (just counting
rooms) is the classic "Meeting Rooms II" problem using a min-heap of end
times. The allocation part adds complexity: you need a second heap to track
which room numbers are available, always picking the lowest.

Maps to: **Day 15 — Intervals, Pattern 3** (merge/overlap counting) combined
with **Day 5 — Heap** (priority queue for greedy allocation).

</details>

---

<details>
<summary><strong>Approach Hint</strong> (only look if stuck after 10 minutes)</summary>

1. Sort meetings by start time (keep track of original indices).
2. Maintain two heaps:
   - **Busy heap (min-heap by end time):** stores `(endTime, roomNumber)`.
     The meeting ending soonest is at the top.
   - **Free heap (min-heap by room number):** stores available room numbers.
     The lowest-numbered free room is at the top.
3. For each meeting (in sorted order):
   - Pop all entries from the busy heap whose `endTime ≤ currentStart`.
     Push their room numbers onto the free heap.
   - If the free heap is non-empty, pop the lowest room number and assign it.
   - Otherwise, allocate a new room (room number = total rooms so far).
   - Push `(endTime, assignedRoom)` onto the busy heap.
   - Record `allocation[originalIndex] = assignedRoom`.
4. The total number of rooms allocated is the answer.

Why two heaps: the busy heap tells you *when* rooms become free. The free
heap tells you *which* free room to pick (lowest number, for greedy reuse).

</details>

---

<details>
<summary><strong>Full Solution in Go</strong> (only look after attempting)</summary>

```go
import (
    "container/heap"
    "sort"
)

// BusyRoom represents a room that's currently occupied
type BusyRoom struct {
    endTime  int
    roomNum  int
}

// BusyHeap is a min-heap ordered by end time
type BusyHeap []BusyRoom
func (h BusyHeap) Len() int            { return len(h) }
func (h BusyHeap) Less(i, j int) bool  { return h[i].endTime < h[j].endTime }
func (h BusyHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *BusyHeap) Push(x interface{}) { *h = append(*h, x.(BusyRoom)) }
func (h *BusyHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// FreeHeap is a min-heap of available room numbers
type FreeHeap []int
func (h FreeHeap) Len() int            { return len(h) }
func (h FreeHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h FreeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *FreeHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *FreeHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func minMeetingRooms(intervals [][]int) (int, []int) {
    n := len(intervals)
    if n == 0 {
        return 0, []int{}
    }

    // Create index-tracking wrapper and sort by start time
    type meeting struct {
        start, end, origIdx int
    }
    meetings := make([]meeting, n)
    for i, iv := range intervals {
        meetings[i] = meeting{iv[0], iv[1], i}
    }
    sort.Slice(meetings, func(i, j int) bool {
        if meetings[i].start == meetings[j].start {
            return meetings[i].end < meetings[j].end
        }
        return meetings[i].start < meetings[j].start
    })

    allocation := make([]int, n)
    busy := &BusyHeap{}
    free := &FreeHeap{}
    heap.Init(busy)
    heap.Init(free)
    totalRooms := 0

    for _, m := range meetings {
        // Free up all rooms whose meetings have ended
        for busy.Len() > 0 && (*busy)[0].endTime <= m.start {
            freed := heap.Pop(busy).(BusyRoom)
            heap.Push(free, freed.roomNum)
        }

        var room int
        if free.Len() > 0 {
            // Reuse the lowest-numbered free room
            room = heap.Pop(free).(int)
        } else {
            // Allocate a new room
            room = totalRooms
            totalRooms++
        }

        allocation[m.origIdx] = room
        heap.Push(busy, BusyRoom{m.end, room})
    }

    return totalRooms, allocation
}
```

**Complexity:**
- Time: O(n log n). Sorting is O(n log n). Each meeting involves at most
  one push/pop on each heap, each O(log n). Total: O(n log n).
- Space: O(n) for the heaps and the allocation array.

**Trace through example 5:** `intervals = [[1,3],[2,4],[5,7],[6,8]]`
```
Sorted meetings (already sorted by start):
  {1,3,orig=0}, {2,4,orig=1}, {5,7,orig=2}, {6,8,orig=3}

Meeting {1,3}:
  busy empty, free empty → new room 0. totalRooms=1.
  busy: [{end:3, room:0}]

Meeting {2,4}:
  busy top endTime=3 > start=2 → no freeing.
  free empty → new room 1. totalRooms=2.
  busy: [{3,0}, {4,1}]

Meeting {5,7}:
  busy top endTime=3 ≤ start=5 → free room 0. Pop.
  busy top endTime=4 ≤ start=5 → free room 1. Pop.
  free: [0, 1].
  Pop free → room 0. busy: [{7,0}]

Meeting {6,8}:
  busy top endTime=7 > start=6 → no freeing.
  free: [1]. Pop → room 1. busy: [{7,0}, {8,1}]

totalRooms=2. allocation = [0, 1, 0, 1].
```

**Edge cases handled:**
- Empty input → immediate return.
- No overlaps → one room reused for everything.
- All overlapping → each meeting gets its own room.
- End time equals start time → not a conflict (end is exclusive).
- Allocation preserves original input order via `origIdx`.

**Simplification (if the interviewer only asks for the room count):**

If you don't need the allocation array, you can simplify significantly
by using a single min-heap of end times:

```go
func minMeetingRooms(intervals [][]int) int {
    if len(intervals) == 0 {
        return 0
    }
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    // Min-heap of end times
    h := &IntHeap{intervals[0][1]}
    heap.Init(h)

    for i := 1; i < len(intervals); i++ {
        if (*h)[0] <= intervals[i][0] {
            heap.Pop(h) // reuse this room
        }
        heap.Push(h, intervals[i][1])
    }

    return h.Len()
}
```

Start with this simpler version, then extend if the interviewer asks
for room assignments.

</details>

---

### Pattern Debrief

| Aspect | Detail |
|--------|--------|
| **Patterns** | Interval scheduling + Min-heap for greedy allocation (two patterns combined) |
| **Days** | Day 15 — Intervals, Pattern 3 + Day 5 — Heaps |
| **Trigger** | "Minimum number of [resources] for overlapping [intervals]" = interval sweep + heap. The allocation requirement adds a second heap for room ID management. |
| **Key insight** | Two heaps serve different purposes: the busy-heap answers "which rooms are free NOW?" (sorted by end time), while the free-heap answers "which free room should I pick?" (sorted by room number). Separating these concerns makes the code clean and correct. |
| **Common mistakes** | (1) Forgetting to sort by start time — the sweep line approach requires processing events in chronological order. (2) Using a single heap and losing room identity — you can count rooms with one heap, but you can't track allocations. (3) Go's `container/heap` boilerplate — in an interview, mention that you know the interface (`Len`, `Less`, `Swap`, `Push`, `Pop`) and write it quickly. Practice the boilerplate until it's muscle memory. (4) Off-by-one on overlap detection: `end <= start` means no conflict (exclusive end), not `end < start`. |

---

## Scoring Rubric

**Score each problem immediately after your 25 minutes are up.** Be honest.

| Criterion | Points | How to Score |
|-----------|--------|-------------|
| **Pattern identified within 5 min** | +2 | Did you name the correct pattern/data structure before minute 5? |
| **Correct approach within 10 min** | +2 | Could you explain the full algorithm (recurrence, traversal type, etc.) before minute 10? |
| **Working solution within 25 min** | +3 | Does your code produce correct output for all given examples? Partial credit: +1 if the logic is right but has 1-2 bugs. |
| **Clean code (no bugs on first test)** | +1 | When you traced through your first example, did it work without edits? |
| **Optimal time complexity** | +1 | Is your solution the intended complexity? (O(n) for P1, O(n log M) for P2, O(V+E) for P3, O(n log n) for P4) |
| **Edge cases handled** | +1 | Did you explicitly handle or mention: empty input, single element, impossible cases, duplicates? |

**Maximum: 10 points per problem. 40 points total.**

### Score Interpretation

| Score | Assessment | Action |
|-------|-----------|--------|
| 36-40 | Interview ready | You're sharp. Day 21 is polish — focus on communication and speed. |
| 28-35 | Strong foundation, gaps in execution | Review the problems you dropped points on. Pattern recognition is there; drill implementation. |
| 20-27 | Pattern recognition needs work | Go back to the specific days for topics you missed. Redo those days' drills. |
| Below 20 | Core gaps remain | Focus Day 21 on re-drilling the weakest 2-3 topics rather than another mixed session. |

---

## Scorecard

Copy this and fill it in:

```
Problem 1 — Minimum Window Containing All Characters
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:    ___ / 2
  Working solution:      ___ / 3
  Clean first test:      ___ / 1
  Optimal complexity:    ___ / 1
  Edge cases:            ___ / 1
  Subtotal:              ___ / 10
  Time used:             ___ min

Problem 2 — Smallest Divisor Given a Threshold
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:    ___ / 2
  Working solution:      ___ / 3
  Clean first test:      ___ / 1
  Optimal complexity:    ___ / 1
  Edge cases:            ___ / 1
  Subtotal:              ___ / 10
  Time used:             ___ min

Problem 3 — Course Schedule Ordering
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:    ___ / 2
  Working solution:      ___ / 3
  Clean first test:      ___ / 1
  Optimal complexity:    ___ / 1
  Edge cases:            ___ / 1
  Subtotal:              ___ / 10
  Time used:             ___ min

Problem 4 — Meeting Rooms with Allocation
  Pattern in 5 min:      ___ / 2
  Approach in 10 min:    ___ / 2
  Working solution:      ___ / 3
  Clean first test:      ___ / 1
  Optimal complexity:    ___ / 1
  Edge cases:            ___ / 1
  Subtotal:              ___ / 10
  Time used:             ___ min

SESSION 2 TOTAL:         ___ / 40
```

---

## Reflection Template

**Complete this after the session. Spend the full 20 minutes on review and
reflection. Writing forces clarity — don't skip this.**

### 1. Which problem took the longest? Why?

```
Problem #___. Time used: ___ min.
Root cause (circle one):
  [ ] Didn't recognize the pattern
  [ ] Recognized pattern but couldn't translate to code
  [ ] Had the code but got stuck on bugs/edge cases
  [ ] Ran out of time during implementation
  [ ] Go-specific issue (heap boilerplate, syntax, etc.)

Specific detail:


```

### 2. Which pattern did you recognize fastest?

```
Problem #___. Recognized in ~___ seconds.
What was the trigger word/phrase in the problem statement?


```

### 3. Where did you get stuck — recognition, implementation, or edge cases?

```
Breakdown across all 4 problems:
  Pattern recognition:  ___ problems with no hesitation
  Implementation:       ___ problems where I knew what to do but struggled to code it
  Edge cases:           ___ problems where my first code missed a case

Weakest area:


```

### 4. Compare to Day 19

```
Day 19 total score:     ___ / 40
Day 20 total score:     ___ / 40
Improvement:            ___

Day 19 time breakdown (avg min per problem): ___
Day 20 time breakdown (avg min per problem): ___

Day 19 weakest problem: #___ (pattern: _____________)
Day 20 weakest problem: #___ (pattern: _____________)

Did the weak area change? If the same patterns are weak across both
sessions, that's a clear signal for Day 21 focus.


```

### 5. Which patterns are now automatic vs still require thinking?

```
AUTOMATIC (recognized in <30 sec, coded without hesitation):
  -
  -
  -

FAMILIAR (recognized in 1-2 min, some implementation friction):
  -
  -
  -

STILL SHAKY (took >3 min to recognize, or struggled to code):
  -
  -
  -
```

### 6. Top 3 weak areas for Day 21 review

```
These are the three things I most need to drill before the final session:

1. ________________________________ (from Day ___)
   Why: ________________________________

2. ________________________________ (from Day ___)
   Why: ________________________________

3. ________________________________ (from Day ___)
   Why: ________________________________

Plan for Day 21: Focus review time on these three areas. For each, redo
the core drill from that day until it's automatic.
```

### 7. If this were a real interview, would I have passed?

```
Realistically (be honest):
  Problem 1: [ ] Passed  [ ] Struggled  [ ] Failed
  Problem 2: [ ] Passed  [ ] Struggled  [ ] Failed
  Problem 3: [ ] Passed  [ ] Struggled  [ ] Failed
  Problem 4: [ ] Passed  [ ] Struggled  [ ] Failed

A typical bar: solve 2 mediums cleanly in 45 min, or 1 medium + meaningful
progress on a hard. Did you meet that bar?

Compared to Day 19, am I closer to the bar? Farther?


```

---

## Session Schedule (2 Hours)

| Time | Activity |
|------|----------|
| 0:00 - 0:05 | Read the session format, rules, and time management tips. Set up timer and blank editor. |
| 0:05 - 0:30 | **Problem 1** — Minimum Window Containing All Characters (25 min) |
| 0:30 - 0:32 | Score Problem 1. Briefly review the solution if needed. Reset editor. |
| 0:32 - 0:57 | **Problem 2** — Smallest Divisor Given a Threshold (25 min) |
| 0:57 - 0:59 | Score Problem 2. Brief review. Reset. |
| 0:59 - 1:24 | **Problem 3** — Course Schedule Ordering (25 min) |
| 1:24 - 1:26 | Score Problem 3. Brief review. Reset. |
| 1:26 - 1:51 | **Problem 4** — Meeting Rooms with Allocation (25 min) |
| 1:51 - 1:53 | Score Problem 4. Brief review. |
| 1:53 - 2:00 | **Reflection.** Fill out the entire reflection template including Day 19 comparison. Identify your Day 21 focus areas. |

---

*Tomorrow (Day 21): Final Review & Weak-Spot Drill. Use today's reflection
to build a targeted plan. Your top 3 weak areas become your Day 21
curriculum. The goal: walk into any interview knowing you've closed every
gap you identified.*
