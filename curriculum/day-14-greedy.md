# Day 14: Greedy Algorithms

> **Time:** 2 hours | **Level:** Refresher | **Language:** Go

The hard part of greedy isn't the code — it's **recognizing that greedy works**. Every
greedy problem has a simple implementation once you see the insight. The real interview
skill is distinguishing greedy from DP, and being able to articulate *why* the greedy
choice is safe.

---

## Pattern Catalog

---

### Pattern 1: Running State / Kadane's

**Trigger:** "Maximum subarray," "best profit from one buy/sell," "maximum sum ending at each position."

**Why greedy works:** A negative running sum can never help a future subarray — resetting to zero (or to the current element) is always at least as good as carrying the deficit forward.

**Template (Kadane's — Maximum Subarray, LC 53):**

```go
func maxSubArray(nums []int) int {
    maxSum := nums[0]
    curSum := nums[0]

    for i := 1; i < len(nums); i++ {
        // KEY DECISION: extend the current subarray, or start fresh?
        // If curSum is negative, it can only hurt — start fresh.
        if curSum < 0 {
            curSum = nums[i]
        } else {
            curSum += nums[i]
        }

        if curSum > maxSum {
            maxSum = curSum
        }
    }

    return maxSum
}
```

**Template (Best Time to Buy and Sell Stock, LC 121):**

```go
func maxProfit(prices []int) int {
    minPrice := prices[0]
    maxProfit := 0

    for i := 1; i < len(prices); i++ {
        profit := prices[i] - minPrice
        if profit > maxProfit {
            maxProfit = profit
        }
        if prices[i] < minPrice {
            minPrice = prices[i]
        }
    }

    return maxProfit
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- All-negative array: Kadane's must return the max single element, NOT zero. Initialize `maxSum = nums[0]`, not `maxSum = 0`.
- Buy/sell stock: you can't sell before buying. Track the minimum price *so far*, not the global minimum.

---

### Pattern 2: Farthest Reach

**Trigger:** "Can you reach the end?" "Minimum jumps to reach the end." Any problem where you advance through positions and choose how far to go.

**Why greedy works:** If you can reach position `i`, and from `i` you can reach `i + nums[i]`, then you can reach everything in between. Tracking just the farthest reachable position captures all the information you need — you never need to reconsider an earlier, shorter reach.

**Template (Jump Game, LC 55):**

```go
func canJump(nums []int) bool {
    farthest := 0

    for i := 0; i < len(nums); i++ {
        if i > farthest {
            return false          // stuck — can't reach position i
        }
        if i+nums[i] > farthest {
            farthest = i + nums[i]
        }
    }

    return true
}
```

**Template (Jump Game II — minimum jumps, LC 45):**

```go
func jump(nums []int) int {
    jumps := 0
    curEnd := 0               // farthest we can go with current number of jumps
    farthest := 0             // farthest we can go with one more jump

    // Note: iterate to len-2, not len-1 (we don't jump FROM the last index)
    for i := 0; i < len(nums)-1; i++ {
        if i+nums[i] > farthest {
            farthest = i + nums[i]
        }

        if i == curEnd {      // exhausted current jump's range — must jump again
            jumps++
            curEnd = farthest
        }
    }

    return jumps
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- Jump Game II is NOT simple greedy — it's a BFS-style level processing. `curEnd` is the boundary of the current "level." When you hit it, you increment jumps and extend the boundary to `farthest`. Think of it as BFS where each level = one jump.
- Off-by-one: iterate to `len(nums)-2` in Jump Game II. If you include the last index, you might count an extra jump.

---

### Pattern 3: Interval Scheduling

**Trigger:** "Maximum non-overlapping intervals," "minimum removals to eliminate overlaps," "meeting rooms," "merge intervals" (though merge is more sweep than greedy).

**Why greedy works:** Picking the interval that ends earliest leaves the most room for future intervals. An interval that ends later can only block more — never fewer — future options.

**Template (Non-Overlapping Intervals / Activity Selection, LC 435):**

```go
func eraseOverlapIntervals(intervals [][]int) int {
    // CRITICAL: sort by END time, not start time
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][1] < intervals[j][1]
    })

    count := 0               // number of intervals to remove
    prevEnd := intervals[0][1]

    for i := 1; i < len(intervals); i++ {
        if intervals[i][0] < prevEnd {
            // Overlap — remove this interval (it ends later or equal)
            count++
        } else {
            // No overlap — keep this interval, update prevEnd
            prevEnd = intervals[i][1]
        }
    }

    return count
}
```

**Complexity:** O(n log n) time (sorting dominates), O(1) extra space.

**Watch out:**
- **Sorting by START time does NOT work** for activity selection. Counterexample: `[1,100], [2,3], [4,5]` — sorting by start picks `[1,100]` first, which blocks everything. Sorting by end picks `[2,3]` then `[4,5]` — two intervals instead of one.
- For "minimum meeting rooms" (LC 253), greedy doesn't directly apply — use a min-heap or sweep line.
- When the problem says "minimum removals," it's equivalent to "maximum non-overlapping" — `removals = total - maxNonOverlapping`.

---

### Pattern 4: Task Scheduling

**Trigger:** "Minimum time to execute tasks with cooldown," "task scheduler," "rearrange string k distance apart."

**Why greedy works:** The most frequent task dictates the minimum length. You must space out its occurrences by the cooldown period, and everything else fills the gaps. No rearrangement can beat this bound.

**Template (Task Scheduler, LC 621):**

```go
func leastInterval(tasks []byte, n int) int {
    freq := [26]int{}
    for _, t := range tasks {
        freq[t-'A']++
    }

    // Find the max frequency and how many tasks share it
    maxFreq := 0
    maxCount := 0
    for _, f := range freq {
        if f > maxFreq {
            maxFreq = f
            maxCount = 1
        } else if f == maxFreq {
            maxCount++
        }
    }

    // Formula: (maxFreq - 1) chunks of size (n + 1), plus the maxCount tasks in last chunk
    //
    // Example: tasks=[A,A,A,B,B,B], n=2
    // A _ _ | A _ _ | A  →  (3-1) chunks of size 3, plus last chunk
    // A B _ | A B _ | A B →  fill B into gaps
    // Result: max(formula, len(tasks))
    //
    //   (maxFreq-1) * (n+1) + maxCount
    //       2       *   3    +    2     = 8
    //   len(tasks) = 6
    //   answer = max(8, 6) = 8

    result := (maxFreq-1)*(n+1) + maxCount
    if len(tasks) > result {
        return len(tasks)     // enough variety to fill all gaps — no idle time
    }
    return result
}
```

**Complexity:** O(n) time (where n = number of tasks), O(1) space (26-letter alphabet).

**Watch out:**
- The `max(formula, len(tasks))` is essential. When there are many distinct tasks, the idle time disappears and the answer is simply the total number of tasks.
- Don't confuse this with a simulation — the formula approach is O(1) per task, no heap needed.

---

### Pattern 5: Gas Station / Circular Tour

**Trigger:** "Circular route," "gas stations," "starting point for a round trip."

**Why greedy works:** If the total gas >= total cost, a solution must exist. And if you fail to reach station `j` starting from station `i`, then starting from any station between `i` and `j` also fails (they'd have even less surplus). So you can skip ahead to `j+1`.

**Template (Gas Station, LC 134):**

```go
func canCompleteCircuit(gas []int, cost []int) int {
    totalSurplus := 0
    currentSurplus := 0
    start := 0

    for i := 0; i < len(gas); i++ {
        totalSurplus += gas[i] - cost[i]
        currentSurplus += gas[i] - cost[i]

        if currentSurplus < 0 {
            // Can't reach i+1 from start — reset
            start = i + 1
            currentSurplus = 0
        }
    }

    if totalSurplus < 0 {
        return -1                 // impossible regardless of start
    }
    return start
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- The `totalSurplus < 0` check is necessary. Without it, you might return a start index even when no valid tour exists.
- Don't try to simulate the circular route literally — the one-pass approach handles it.
- The greedy reset (`start = i + 1`) skips ALL stations from old start through `i`. This is the key insight interviewers probe.

---

### Pattern 6: Assign / Distribute Greedily

**Trigger:** "Assign cookies to children," "boats to save people," "pair elements to minimize/maximize." Problems where you match items from two groups.

**Why greedy works:** After sorting, the optimal pairing is always between adjacent or extreme elements. You can prove by exchange argument: swapping any greedy assignment for a non-greedy one never improves the result.

**Template (Assign Cookies, LC 455):**

```go
func findContentChildren(g []int, s []int) int {
    sort.Ints(g)              // children's greed factors
    sort.Ints(s)              // cookie sizes

    child, cookie := 0, 0
    for child < len(g) && cookie < len(s) {
        if s[cookie] >= g[child] {
            child++               // this child is satisfied — move to next child
        }
        cookie++                  // either used this cookie or it's too small — move on
    }

    return child
}
```

**Template (Boats to Save People, LC 881):**

```go
func numRescueBoats(people []int, limit int) int {
    sort.Ints(people)
    lo, hi := 0, len(people)-1
    boats := 0

    for lo <= hi {
        if people[lo]+people[hi] <= limit {
            lo++                  // lightest person pairs with heaviest
        }
        hi--                      // heaviest always takes a boat
        boats++
    }

    return boats
}
```

**Complexity:** O(n log n) time (sorting), O(1) extra space.

**Watch out:**
- The two-pointer approach for boats requires sorting first.
- In assign cookies, always advance the cookie pointer — a too-small cookie is useless and must be skipped.
- These problems look like they might need DP or backtracking. They don't. The sort + greedy combination is provably optimal via exchange argument.

---

## Greedy vs DP Decision Framework

**This is the single most valuable skill for greedy problems in interviews.** The code is always simple. Knowing *when* to write it is everything.

### The Litmus Test

> "Can I construct a counterexample where making the locally optimal choice leads to a globally suboptimal result?"

- **If yes → you need DP** (or some other approach).
- **If no → greedy is likely correct.** Try to state why in one sentence (the "greedy choice property").

### Three Examples Where Greedy Works (and Why)

**1. Jump Game (LC 55)**
Greedy choice: always track the farthest reachable position.
Why it works: reaching farther is never worse than reaching shorter. If position `j > i` is reachable, everything between `i` and `j` is also reachable.
Counterexample attempt: none — extending reach has no cost.

**2. Activity Selection / Non-Overlapping Intervals (LC 435)**
Greedy choice: pick the interval that ends earliest.
Why it works: ending earlier leaves the most room for future intervals. An interval ending at time 5 never blocks more than one ending at time 3.
Counterexample attempt: none — earlier end is always >= as good.

**3. Huffman Coding**
Greedy choice: always merge the two lowest-frequency nodes.
Why it works: low-frequency symbols should have long codes (deep in tree). Merging two large frequencies first would waste short codes on rare symbols.
Counterexample attempt: none — provable via exchange argument.

### Three Examples Where Greedy FAILS (and DP is Needed)

**1. Coin Change (LC 322) — arbitrary denominations**
Greedy choice attempt: always pick the largest coin.
Counterexample: coins = `[1, 3, 4]`, amount = `6`.
Greedy: `4 + 1 + 1` = 3 coins. Optimal: `3 + 3` = 2 coins.
Why greedy fails: a large coin now can prevent a better combination later. Subproblems overlap.

**2. Longest Increasing Subsequence (LC 300)**
Greedy choice attempt: always extend with the smallest possible next element.
Counterexample: `[3, 1, 8, 2, 5]`.
Greedy starting from `1`: picks `1, 2, 5` (length 3). But needs to consider `1, 8` branch too.
Why greedy fails: the "best" next element depends on what comes after — requires looking ahead.
(Note: the patience sorting / binary search approach is not pure greedy — it maintains state across multiple subsequences.)

**3. 0/1 Knapsack**
Greedy choice attempt: pick the item with the best value/weight ratio.
Counterexample: capacity = `5`, items = `[(value=6, weight=3), (value=5, weight=2), (value=5, weight=2)]`.
Greedy picks `(6,3)` then one `(5,2)` → total = 11, weight = 5.
Optimal: both `(5,2)` items → total = 10? No — greedy actually wins here.
Better counterexample: capacity = `4`, items = `[(value=5, weight=3), (value=4, weight=2), (value=4, weight=2)]`.
Greedy: `(5,3)` → value 5, weight 3, can't fit another. Optimal: `(4,2) + (4,2)` = value 8.
Why greedy fails: taking one item affects which combinations of items you can take. The fractional knapsack is greedy; the 0/1 version is not.

### Patterns That Are ALWAYS Greedy

Memorize these — don't waste time considering DP for them:

| Pattern                        | Why It's Always Greedy                                    |
|--------------------------------|-----------------------------------------------------------|
| Interval scheduling by end time | Earlier ending is provably optimal (exchange argument)    |
| Kadane's max subarray          | Negative prefix always hurts — reset is safe              |
| Gas station circular tour      | If you fail at `j` from `i`, all starts `i..j` also fail |
| Farthest reach / jump game     | Farther is always >= as good as shorter                   |
| Buy/sell stock (one transaction)| Tracking running min is sufficient                       |

### Patterns That Are NEVER Greedy

Don't be tempted:

| Pattern                          | Why Greedy Fails                                  | Use Instead     |
|----------------------------------|---------------------------------------------------|-----------------|
| Coin change (arbitrary coins)    | Large coin now blocks better combo                | DP              |
| 0/1 Knapsack                    | Taking one item constrains future choices          | DP              |
| Longest increasing subsequence   | Best next element depends on future               | DP + bin search |
| Edit distance                   | Each operation affects all subsequent alignment    | 2D DP           |
| Partition into equal subsets     | One assignment affects feasibility of remainder    | DP / backtrack  |

### Quick Decision Flowchart

```
Problem asks for optimal (max/min) over a sequence?
│
├─ Can I reduce the problem to a single pass with a simple rule?
│   ├─ YES: Does the greedy choice provably never hurt?
│   │   ├─ YES → GREEDY (state the greedy choice property)
│   │   └─ NO / UNSURE → find a counterexample
│   │       ├─ Counterexample found → DP
│   │       └─ No counterexample → likely greedy, proceed with caution
│   └─ NO: Are there overlapping subproblems?
│       ├─ YES → DP
│       └─ NO → possibly divide and conquer or other
│
└─ Problem asks for all solutions / count of solutions?
    └─ Almost never greedy → backtracking or DP
```

---

## Common Interview Traps

### 1. Kadane's — All-Negative Array

```go
// WRONG — returns 0 for [-3, -2, -1]
maxSum := 0
curSum := 0
for _, n := range nums {
    curSum = max(0, curSum+n)     // resets to 0, but answer should be -1
    maxSum = max(maxSum, curSum)
}

// RIGHT — initialize from first element
maxSum := nums[0]
curSum := nums[0]
for i := 1; i < len(nums); i++ {
    if curSum < 0 {
        curSum = nums[i]
    } else {
        curSum += nums[i]
    }
    maxSum = max(maxSum, curSum)
}
```

If the problem guarantees at least one positive number, resetting to 0 works. Otherwise, it doesn't. Ask the interviewer.

### 2. Jump Game II — It's BFS, Not Simple Greedy

The "minimum jumps" problem looks like you should always jump as far as possible. That's wrong.

```
[2, 3, 1, 1, 4]
     ^--- jumping to index 1 (value 3) is better than jumping to index 2 (value 1)
          even though index 2 is farther
```

The correct approach tracks a "level boundary" (`curEnd`), like BFS levels. Within each level, you find the farthest reachable position, then advance the boundary. This is greedy in the sense that you pick the best next level, but the structure is BFS.

### 3. Interval Scheduling — Sort by END, Not START

This is the most common greedy mistake. Prove it to yourself:

```
Intervals: [1,100], [2,3], [4,5]

Sort by start: [1,100], [2,3], [4,5]
  Pick [1,100] → blocks everything else → 1 interval

Sort by end:   [2,3], [4,5], [1,100]
  Pick [2,3] → Pick [4,5] → [1,100] overlaps → 2 intervals ✓
```

### 4. Gas Station — The totalSurplus Check

```go
// WRONG — might return a start index even when no solution exists
for i := 0; i < len(gas); i++ {
    currentSurplus += gas[i] - cost[i]
    if currentSurplus < 0 {
        start = i + 1
        currentSurplus = 0
    }
}
return start  // BUG: what if total gas < total cost?

// RIGHT — add the global feasibility check
if totalSurplus < 0 {
    return -1
}
return start
```

### 5. "It Feels Greedy" Isn't a Proof

Interviewers will ask: *"Why does greedy work here?"* You need a one-sentence answer. Practice these:

| Problem                      | One-Sentence Justification                                           |
|------------------------------|----------------------------------------------------------------------|
| Max Subarray (Kadane's)      | A negative running sum can only decrease future sums — discard it.   |
| Jump Game                    | Reaching farther has no cost, so it's never suboptimal.              |
| Non-Overlapping Intervals    | Ending earlier leaves maximal room for future intervals.             |
| Gas Station                  | Failing at `j` from `i` means every start in `[i,j]` also fails.    |
| Task Scheduler               | The most frequent task forces minimum idle time; others fill gaps.   |
| Boats to Save People         | Pairing heaviest with lightest maximizes the use of each boat.       |

If you can't state the justification in one sentence, reconsider whether greedy actually works.

---

## Thought Process Walkthrough

### Walkthrough 1: Jump Game (LC 55)

> Given `nums = [2, 3, 1, 1, 4]`, can you reach the last index?

**Step 1 — Recognize the pattern.**
"Can you reach the end by jumping forward" → farthest reach pattern (Pattern 2).

**Step 2 — State the greedy choice.**
"At each position, update the farthest index I can reach. If I ever find myself at a position beyond farthest, I'm stuck."

**Step 3 — Justify to the interviewer.**
"Reaching farther is never worse than reaching shorter — if I can reach index 5, I can also reach indices 0-4. So tracking a single `farthest` variable captures all reachable positions."

**Step 4 — Trace through the example.**

```
nums = [2, 3, 1, 1, 4]
        ^
i=0: farthest = max(0, 0+2) = 2     can reach indices 0-2
i=1: 1 <= 2 ✓  farthest = max(2, 1+3) = 4   can reach indices 0-4
i=2: 2 <= 4 ✓  farthest = max(4, 2+1) = 4   unchanged
i=3: 3 <= 4 ✓  farthest = max(4, 3+1) = 4   unchanged
i=4: 4 <= 4 ✓  reached the end → return true
```

**Step 5 — Trace a failing case.**

```
nums = [3, 2, 1, 0, 4]
i=0: farthest = 3
i=1: farthest = max(3, 3) = 3
i=2: farthest = max(3, 3) = 3
i=3: farthest = max(3, 3) = 3       stuck at 3, need to reach 4
i=4: 4 > 3 → return false
```

**Step 6 — Code it.** (See Pattern 2 template above.)

**Step 7 — Complexity.** O(n) time, O(1) space. Single pass.

---

### Walkthrough 2: Non-Overlapping Intervals (LC 435)

> Given `intervals = [[1,2], [2,3], [3,4], [1,3]]`, return the minimum number to remove so the rest don't overlap.

**Step 1 — Recognize the pattern.**
"Remove minimum intervals to eliminate overlaps" → interval scheduling (Pattern 3). Equivalent to: keep maximum non-overlapping intervals, answer = total − kept.

**Step 2 — State the greedy choice.**
"Sort by end time. Greedily keep each interval that doesn't overlap with the previously kept one. If it overlaps, remove it (increment count)."

**Step 3 — Justify to the interviewer.**
"Among all intervals I could keep, the one ending earliest leaves the most room for future intervals. Picking anything else could only block more — never fewer — future options."

**Step 4 — Trace through the example.**

```
Input:   [[1,2], [2,3], [3,4], [1,3]]
Sorted by end: [[1,2], [2,3], [1,3], [3,4]]

Start: prevEnd = 2 (from [1,2]), removals = 0

i=1: [2,3] → start 2 >= prevEnd 2 → no overlap → keep, prevEnd = 3
i=2: [1,3] → start 1 < prevEnd 3  → overlap → remove, removals = 1
i=3: [3,4] → start 3 >= prevEnd 3 → no overlap → keep, prevEnd = 4

Answer: 1 removal
Kept: [1,2], [2,3], [3,4] — 3 non-overlapping intervals ✓
```

**Step 5 — Why sort by end, not start?**

If we sorted by start: `[[1,2], [1,3], [2,3], [3,4]]`
Picking `[1,2]` first → `[1,3]` overlaps (remove) → `[2,3]` ok → `[3,4]` ok → 1 removal. Happens to work here, but:

Counterexample: `[[1,10], [2,3], [4,5], [6,7]]`
Sort by start: pick `[1,10]` → blocks everything → keep 1.
Sort by end: pick `[2,3]`, then `[4,5]`, then `[6,7]` → keep 3.

**Step 6 — Code it.** (See Pattern 3 template above.)

**Step 7 — Complexity.** O(n log n) time for sorting, O(1) extra space.

---

## Time Targets

| Problem                      | LC #  | Target | Notes                                       |
|------------------------------|-------|--------|---------------------------------------------|
| Maximum Subarray             | 53    | 5 min  | Kadane's — must be instant                  |
| Best Time to Buy/Sell Stock  | 121   | 5 min  | Running min — trivial once you see it       |
| Jump Game                    | 55    | 5 min  | Farthest reach one-liner loop               |
| Jump Game II                 | 45    | 8 min  | BFS-style level processing — trickier       |
| Non-Overlapping Intervals    | 435   | 8 min  | Sort by end + single pass                   |
| Task Scheduler               | 621   | 10 min | Formula derivation takes a moment           |
| Gas Station                  | 134   | 8 min  | Two-variable tracking + feasibility check   |
| Assign Cookies               | 455   | 5 min  | Sort + two pointers                         |
| Boats to Save People         | 881   | 5 min  | Sort + two pointers from both ends          |

---

## Quick Drill (30 minutes)

Do these without looking at templates. Time yourself.

1. **Maximum Subarray** (LC 53) — Write Kadane's. Test with `[-2, 1, -3, 4, -1, 2, 1, -5, 4]` AND `[-3, -2, -1]`. Target: 5 minutes.

2. **Jump Game** (LC 55) — Write farthest reach. Test with a failing case. Target: 5 minutes.

3. **Non-Overlapping Intervals** (LC 435) — Sort by END time. Say out loud why start-time sorting fails. Target: 8 minutes.

4. **Gas Station** (LC 134) — Write the one-pass solution. Don't forget the `totalSurplus` check. Target: 8 minutes.

After each one, check:
- Did I handle edge cases (all negative, impossible tour, no intervals)?
- Can I state in one sentence WHY greedy works for this problem?
- Could I produce a counterexample if someone suggested DP?

---

## Self-Assessment

### Can I explain these from memory?

| Question                                                            | Confident? |
|---------------------------------------------------------------------|------------|
| What is the greedy choice property? Can I state it in one sentence? |            |
| Why does Kadane's reset when `curSum < 0`?                          |            |
| Why sort by end time for interval scheduling, not start time?       |            |
| Why is Jump Game II BFS-like rather than simple greedy?             |            |
| What is the task scheduler formula and why does it work?            |            |
| How does the gas station reset argument work?                       |            |
| Coin change with arbitrary coins — why does greedy fail?            |            |
| 0/1 knapsack — why does greedy fail but fractional knapsack works?  |            |
| Can I name 3 problems that are ALWAYS greedy?                       |            |
| Can I name 3 problems where greedy FAILS and DP is needed?         |            |

### Red flags that you need more practice:
- You reach for DP on a problem that's pure greedy (Jump Game, Kadane's).
- You assume greedy works without being able to articulate why.
- You sort intervals by start time instead of end time.
- Kadane's returns 0 on an all-negative array.
- You can't produce a counterexample for coin change greedy.

### Green lights — you're ready:
- For every greedy problem, you can state the greedy choice property in one sentence.
- You can identify at least 3 "greedy imposters" (problems that feel greedy but need DP) and produce counterexamples.
- You sort intervals by end time without thinking about it.
- You can explain the Jump Game II BFS structure clearly.
- You hit the time targets above consistently.
