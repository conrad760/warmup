# Day 4: Binary Search

> **Time budget:** 2 hours | **Goal:** Master every binary search variant interviewers throw at you, write bug-free code under pressure, and never get tripped up by off-by-one errors again.

Binary search is deceptively simple. Interviewers love it because a 15-line function has at least 4 places to introduce a subtle bug. This guide drills the variants until the templates are muscle memory.

---

## Pattern Catalog

### 1. Standard Search -- Exact Match in Sorted Array

**Trigger:** "Find target in a sorted array" or "return the index of X."

**Go Template:**

```go
func binarySearch(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target {
            return mid
        } else if nums[mid] < target {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return -1 // not found
}
```

**Complexity:** O(log n) time, O(1) space.

**Watch out:**
- Use `lo <= hi` (not `lo < hi`) -- you must check the final single-element window.
- `mid = lo + (hi-lo)/2` prevents integer overflow. Never write `(lo+hi)/2`.
- Both `lo` and `hi` move *past* mid (`mid+1`, `mid-1`). This is what guarantees termination.

---

### 2. Lower Bound (First Occurrence / bisect_left)

**Trigger:** "Find the first element >= target," "first occurrence of X," "leftmost position," "how many elements are less than X."

**Go Template:**

```go
// Returns the index of the first element >= target.
// If all elements < target, returns len(nums).
func lowerBound(nums []int, target int) int {
    lo, hi := 0, len(nums) // NOTE: hi = len(nums), not len(nums)-1
    for lo < hi {           // NOTE: strict < , not <=
        mid := lo + (hi-lo)/2
        if nums[mid] < target {
            lo = mid + 1 // mid is too small, exclude it
        } else {
            hi = mid // mid *might* be the answer, keep it
        }
    }
    return lo // lo == hi, both point at the insertion point
}
```

**Complexity:** O(log n) time, O(1) space.

**Watch out:**
- `hi` starts at `len(nums)` (one past the end) because the answer could be "after everything."
- Loop is `lo < hi`, NOT `lo <= hi`. With `lo <= hi` here you get an infinite loop when `lo == hi`.
- When `nums[mid] >= target`, set `hi = mid` (NOT `mid-1`). Mid itself might be the answer.
- To find "first occurrence of target": call `lowerBound`, then check `if lo < len(nums) && nums[lo] == target`.

---

### 3. Upper Bound (bisect_right)

**Trigger:** "Find the first element > target," "last occurrence of X," "how many elements are <= X."

**Go Template:**

```go
// Returns the index of the first element > target.
// All elements at indices [0, result) are <= target.
func upperBound(nums []int, target int) int {
    lo, hi := 0, len(nums)
    for lo < hi {
        mid := lo + (hi-lo)/2
        if nums[mid] <= target { // NOTE: <= not <
            lo = mid + 1
        } else {
            hi = mid
        }
    }
    return lo
}
```

**Complexity:** O(log n) time, O(1) space.

**Watch out:**
- The ONLY difference from lower bound is `<=` vs `<` in the comparison. Internalize this: `<` gives you bisect_left, `<=` gives you bisect_right.
- To find "last occurrence of target": `idx := upperBound(nums, target) - 1`, then check `if idx >= 0 && nums[idx] == target`.
- To count occurrences: `upperBound(target) - lowerBound(target)`.

---

### 4. Rotated Sorted Array -- Search & Find Minimum

**Trigger:** "Array was sorted then rotated," "find target in rotated array," "find the minimum element."

**Go Template (Search in Rotated Array):**

```go
func searchRotated(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target {
            return mid
        }
        // Determine which half is sorted
        if nums[lo] <= nums[mid] { // left half is sorted
            if nums[lo] <= target && target < nums[mid] {
                hi = mid - 1 // target is in the sorted left half
            } else {
                lo = mid + 1 // target is in the right half
            }
        } else { // right half is sorted
            if nums[mid] < target && target <= nums[hi] {
                lo = mid + 1 // target is in the sorted right half
            } else {
                hi = mid - 1 // target is in the left half
            }
        }
    }
    return -1
}
```

**Go Template (Find Minimum in Rotated Array):**

```go
func findMin(nums []int) int {
    lo, hi := 0, len(nums)-1
    for lo < hi {
        mid := lo + (hi-lo)/2
        if nums[mid] > nums[hi] {
            lo = mid + 1 // minimum is in the right half
        } else {
            hi = mid // mid itself could be the minimum
        }
    }
    return nums[lo]
}
```

**Complexity:** O(log n) time, O(1) space.

**Watch out:**
- In `searchRotated`: the condition `nums[lo] <= nums[mid]` uses `<=` (not `<`). When `lo == mid` (two-element window), the left "half" is one element and IS sorted.
- The "not rotated" edge case: if the array was rotated by its full length, it's back to sorted. The `nums[lo] <= nums[mid]` check handles this naturally -- but trace through it to convince yourself.
- For `findMin`: compare `nums[mid]` to `nums[hi]`, not `nums[lo]`. Comparing to `nums[lo]` is ambiguous when the array isn't rotated.
- With duplicates, worst case degrades to O(n). Handle `nums[mid] == nums[hi]` by doing `hi--`.

---

### 5. Search on Answer Space

**Trigger:** "Minimum X such that condition is met," "maximum speed/capacity," "can you do it in K days/trips?" -- any problem where you're optimizing a value subject to a constraint and the feasibility function is monotonic.

This is THE pattern most candidates miss. The key insight: you're not binary searching through an array. You're binary searching through the *space of possible answers*.

**Go Template:**

```go
// Generic "minimum value that satisfies condition" template.
func searchOnAnswer(lo, hi int, feasible func(int) bool) int {
    for lo < hi {
        mid := lo + (hi-lo)/2
        if feasible(mid) {
            hi = mid // mid works, but maybe something smaller works too
        } else {
            lo = mid + 1 // mid doesn't work, need bigger
        }
    }
    return lo
}
```

**Example -- Koko Eating Bananas (LC 875):**

```go
func minEatingSpeed(piles []int, h int) int {
    // lo = 1 (minimum possible speed: eat 1 banana/hour)
    // hi = max(piles) (maximum: eat the biggest pile in 1 hour)
    lo, hi := 1, 0
    for _, p := range piles {
        if p > hi {
            hi = p
        }
    }

    for lo < hi {
        mid := lo + (hi-lo)/2
        if canFinish(piles, mid, h) {
            hi = mid // this speed works, try slower
        } else {
            lo = mid + 1 // too slow, need faster
        }
    }
    return lo
}

func canFinish(piles []int, speed, h int) bool {
    hours := 0
    for _, p := range piles {
        hours += (p + speed - 1) / speed // ceiling division
    }
    return hours <= h
}
```

**Complexity:** O(n * log(max_answer)) time, O(1) space. The log factor is from the binary search over the answer space; n is from evaluating the feasibility function.

**Watch out:**
- Getting `lo` and `hi` of the search space wrong is the #1 mistake. Think hard about the absolute minimum and maximum possible answers. Off-by-one here means wrong answer.
- Ceiling division: `(a + b - 1) / b` is the Go idiom. Do not use `math.Ceil` with float conversion -- it introduces floating-point bugs.
- The feasibility function must be *monotonic*: if speed `k` works, then speed `k+1` also works. If this property doesn't hold, binary search on answer doesn't apply.
- For "maximize minimum" problems (e.g., split array largest sum minimized), flip the logic: `feasible` returns true if you CAN achieve the constraint, and you're looking for the minimum "largest sum" that's feasible.

---

### 6. Peak Finding

**Trigger:** "Find a peak element," "bitonic array," "mountain array."

**Go Template:**

```go
func findPeakElement(nums []int) int {
    lo, hi := 0, len(nums)-1
    for lo < hi {
        mid := lo + (hi-lo)/2
        if nums[mid] > nums[mid+1] {
            hi = mid // peak is at mid or to the left
        } else {
            lo = mid + 1 // peak is to the right
        }
    }
    return lo // lo == hi, that's the peak
}
```

**Complexity:** O(log n) time, O(1) space.

**Watch out:**
- `mid+1` is always valid because `lo < hi` guarantees `mid < hi`, so `mid+1 <= hi < len(nums)`.
- We're comparing *adjacent* elements, not comparing to a target. This is what makes it feel different from other binary search problems.
- This finds *a* peak, not necessarily *the* global maximum. That's usually what the problem asks for.
- For a "mountain array" (strictly increases then strictly decreases), same template finds the summit.

---

## Decision Framework

Use this flowchart when you see a new problem:

```
Is the search space explicitly sorted (array/matrix)?
|
+-- YES --> Do you need an exact match?
|           |
|           +-- YES --> Pattern 1: Standard Search
|           +-- NO  --> "First/last occurrence" or "insertion point"?
|                       |
|                       +-- YES --> Pattern 2 or 3: Lower/Upper Bound
|                       +-- NO  --> Is it rotated?
|                                   |
|                                   +-- YES --> Pattern 4: Rotated Array
|                                   +-- NO  --> Is it "find a peak"?
|                                               |
|                                               +-- YES --> Pattern 6: Peak Finding
|
+-- NO  --> Is there a monotonic feasibility function?
            (i.e., if answer X works, then X+1 also works)
            |
            +-- YES --> Pattern 5: Search on Answer Space
            +-- NO  --> Binary search probably doesn't apply.
```

**The meta-signal:** If the brute force is "try every value from 1 to MAX and pick the best valid one," and the check function flips from false to true (or true to false) at some threshold, that's binary search on answer. This converts O(MAX * n) brute force into O(log(MAX) * n).

---

## Common Interview Traps

### Trap 1: `lo <= hi` vs `lo < hi`

| Template | Loop condition | Why |
|---|---|---|
| Standard search (exact match) | `lo <= hi` | Must check the last single-element window; exits when `lo > hi` (not found) |
| Lower/upper bound | `lo < hi` | `lo` and `hi` converge to the answer; when `lo == hi`, that IS the answer |
| Search on answer | `lo < hi` | Same convergence logic as lower bound |
| Peak finding | `lo < hi` | Same convergence logic |

**Rule of thumb:** If you return `mid` inside the loop (exact match), use `lo <= hi`. If the answer is where `lo` and `hi` converge, use `lo < hi`.

### Trap 2: `lo = mid` Causing Infinite Loops

If you ever write `lo = mid` (without `+1`), check for infinite loops with a 2-element array. When `lo + 1 == hi`, `mid = lo`, so `lo = mid` doesn't change anything. Infinite loop.

Fix: if you need `lo = mid`, use `mid = lo + (hi-lo+1)/2` (round UP instead of down). But honestly, restructure your logic to avoid `lo = mid` if possible -- it's a trap.

### Trap 3: Integer Overflow in Mid Calculation

```go
// BAD -- overflows if lo + hi > math.MaxInt
mid := (lo + hi) / 2

// GOOD -- always safe
mid := lo + (hi-lo)/2
```

In Go with 64-bit ints, this rarely causes issues in practice, but interviewers WILL notice if you write the unsafe version. Always use the safe form. It's a free signal of experience.

### Trap 4: Rotated Array Edge Cases

The "not rotated" case (`[1, 2, 3, 4, 5]` rotated by 0) and the "fully rotated" case must both work. Trace your code through:
- `[1, 2, 3, 4, 5]` (not rotated)
- `[2, 1]` (two elements)
- `[1]` (single element)

### Trap 5: Search on Answer -- Bounds Wrong

For Koko eating bananas:
- `lo = 1` (NOT 0 -- division by zero)
- `hi = max(piles)` (NOT sum -- eating faster than the biggest pile wastes time but doesn't help)

For ship packages within D days:
- `lo = max(weights)` (NOT 0 -- ship must carry the heaviest single package)
- `hi = sum(weights)` (ship everything in one day)

Always ask: "What's the absolute minimum possible answer? What's the absolute maximum?"

### Trap 6: Off-by-One in "Last Valid" vs "First Invalid"

Lower bound finds "first position where condition is true." If you need "last position where condition is false," that's `lowerBound - 1`. Don't try to restructure the binary search to find it directly -- that's where bugs creep in. Find the boundary, then adjust.

---

## Thought Process Walkthrough

### Walkthrough 1: Search in Rotated Sorted Array (LC 33)

**Setting:** Interviewer says, "Given a rotated sorted array with distinct values, find the index of a target. O(log n)."

**Step 1 -- Clarify (30 seconds).**
- "Distinct values, correct? No duplicates?" (Yes -- duplicates make it harder, confirm the constraint.)
- "Sorted in ascending order before rotation?" (Yes.)
- "Return -1 if not found?" (Yes.)

**Step 2 -- State the approach (30 seconds).**
"I'll use binary search. The key insight is that in a rotated sorted array, at least one half around `mid` is always sorted. I'll identify which half is sorted, check if the target falls in that sorted range, and eliminate the other half."

**Step 3 -- Write the code. Talk through every decision.**

```go
func search(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    // Using lo <= hi because we return mid on exact match
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target {
            return mid
        }
        // Which half is sorted?
        if nums[lo] <= nums[mid] {
            // Left half [lo..mid] is sorted.
            // Is target in this sorted range?
            if nums[lo] <= target && target < nums[mid] {
                hi = mid - 1
            } else {
                lo = mid + 1
            }
        } else {
            // Right half [mid..hi] is sorted.
            if nums[mid] < target && target <= nums[hi] {
                lo = mid + 1
            } else {
                hi = mid - 1
            }
        }
    }
    return -1
}
```

**Step 4 -- Trace through examples (1-2 minutes).**

Example: `nums = [4, 5, 6, 7, 0, 1, 2]`, `target = 0`.

| Step | lo | hi | mid | nums[mid] | Sorted half | Action |
|------|----|----|-----|-----------|-------------|--------|
| 1 | 0 | 6 | 3 | 7 | left [4,5,6,7] | 0 not in [4,7), go right: lo=4 |
| 2 | 4 | 6 | 5 | 1 | right [1,2] | 0 not in (1,2], go left: hi=4 |
| 3 | 4 | 4 | 4 | 0 | found! | return 4 |

Edge case: `nums = [1, 3]`, `target = 3`.

| Step | lo | hi | mid | nums[mid] | Sorted half | Action |
|------|----|----|-----|-----------|-------------|--------|
| 1 | 0 | 1 | 0 | 1 | left [1] (lo<=mid) | 3 not in [1,1), go right: lo=1 |
| 2 | 1 | 1 | 1 | 3 | found! | return 1 |

**Step 5 -- State complexity.** "O(log n) time, O(1) space."

**Why interviewers like this problem:** The `nums[lo] <= nums[mid]` condition (with `<=`, not `<`) is the subtle part. Most bugs come from getting that wrong, or from getting the range checks wrong (`<=` vs `<` on each side). Tracing through a 2-element example catches both.

---

### Walkthrough 2: Koko Eating Bananas (LC 875)

**Setting:** "Koko has piles of bananas. She can eat at speed `k` bananas/hour, one pile at a time (if a pile has fewer than `k`, she finishes it and waits). She has `h` hours. Find the minimum integer `k` such that she can eat all bananas in `h` hours."

**Step 1 -- Clarify (30 seconds).**
- "h is always >= len(piles)?" (Yes, confirmed by constraints -- she has at least one hour per pile.)
- "piles[i] >= 1?" (Yes.)

**Step 2 -- Recognize the pattern (this is the critical moment).**
"The brute force is: try `k = 1, 2, 3, ...` up to `max(piles)`, and for each `k`, compute total hours. The first `k` that works is the answer. But the feasibility function is monotonic -- if speed `k` works, then `k+1` also works. So I can binary search on `k`."

Say this out loud. Interviewers want to hear you identify the monotonic property.

**Step 3 -- Define the search space.**
- "The minimum possible speed is 1 (eating 1 banana per hour)."
- "The maximum possible speed is `max(piles)` -- eating faster than the biggest pile doesn't help since she waits after finishing a pile."

**Step 4 -- Write the feasibility function first.**

```go
func hoursNeeded(piles []int, speed int) int {
    total := 0
    for _, p := range piles {
        total += (p + speed - 1) / speed // ceiling division
    }
    return total
}
```

Explain: "For each pile of size `p`, she needs `ceil(p / speed)` hours. I use the integer ceiling trick `(p + speed - 1) / speed` to avoid floating point."

**Step 5 -- Write the binary search.**

```go
func minEatingSpeed(piles []int, h int) int {
    lo, hi := 1, 0
    for _, p := range piles {
        if p > hi {
            hi = p
        }
    }

    // Find minimum speed where hoursNeeded <= h
    for lo < hi {
        mid := lo + (hi-lo)/2
        if hoursNeeded(piles, mid) <= h {
            hi = mid // this speed works, try slower
        } else {
            lo = mid + 1 // too slow
        }
    }
    return lo
}
```

**Step 6 -- Trace through an example.**

`piles = [3, 6, 7, 11]`, `h = 8`.

Search space: `lo = 1`, `hi = 11`.

| Step | lo | hi | mid | hours needed | Action |
|------|----|----|-----|-------------|--------|
| 1 | 1 | 11 | 6 | 1+1+2+2 = 6 <= 8 | works, hi=6 |
| 2 | 1 | 6 | 3 | 1+2+3+4 = 10 > 8 | too slow, lo=4 |
| 3 | 4 | 6 | 5 | 1+2+2+3 = 8 <= 8 | works, hi=5 |
| 4 | 4 | 5 | 4 | 1+2+2+3 = 8 <= 8 | works, hi=4 |
| 5 | lo==hi==4 | done | | | return 4 |

**Step 7 -- State complexity.**
"O(n * log(max_pile)) time where n = len(piles). O(1) space."

**Why interviewers love this:** It tests whether you can see that binary search applies when there's no array to search. The "search on answer space" insight separates candidates who've memorized patterns from those who understand the underlying principle.

---

## Time Targets

| Problem type | Target time | Notes |
|---|---|---|
| Standard binary search | 3-4 min | Should be automatic |
| Lower/upper bound | 4-5 min | Template should be memorized |
| Rotated sorted array | 8-10 min | Budget time for edge case tracing |
| Search on answer space | 10-12 min | Most time goes to defining the search space and feasibility function |
| Peak finding | 5-6 min | Quick once you see the pattern |

If you're exceeding these times, the template isn't in muscle memory yet. Drill more.

---

## Quick Drill (5 Exercises)

Complete these in order. Time yourself. Do not look at the templates above.

### Drill 1: First Bad Version (LC 278) -- 5 min
You have `n` versions `[1, 2, ..., n]` and a function `isBadVersion(version) bool`. Find the first bad version. Map this to a pattern, write the code, trace through `n=5, first_bad=4`.

<details>
<summary>Pattern & key decision</summary>

Lower bound. Search space is `[1, n]`. `lo < hi` loop. If `isBadVersion(mid)`, set `hi = mid`. Else `lo = mid + 1`. Return `lo`.

```go
func firstBadVersion(n int) int {
    lo, hi := 1, n
    for lo < hi {
        mid := lo + (hi-lo)/2
        if isBadVersion(mid) {
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}
```
</details>

### Drill 2: Find First and Last Position of Element in Sorted Array (LC 34) -- 8 min
Given a sorted array and target, return `[first, last]` index. If not found, return `[-1, -1]`. Use lower bound and upper bound together.

<details>
<summary>Pattern & key decision</summary>

Call lower bound to get `first`. Check if `nums[first] == target`. Call upper bound and subtract 1 to get `last`.

```go
func searchRange(nums []int, target int) []int {
    if len(nums) == 0 {
        return []int{-1, -1}
    }
    first := lowerBound(nums, target)
    if first == len(nums) || nums[first] != target {
        return []int{-1, -1}
    }
    last := upperBound(nums, target) - 1
    return []int{first, last}
}
```
</details>

### Drill 3: Find Minimum in Rotated Sorted Array (LC 153) -- 6 min
Distinct elements. Write the code, trace through `[3, 4, 5, 1, 2]` and `[1, 2, 3, 4, 5]` (not rotated).

<details>
<summary>Pattern & key decision</summary>

Compare `nums[mid]` to `nums[hi]`. If `nums[mid] > nums[hi]`, minimum is to the right (`lo = mid + 1`). Else `hi = mid`. Return `nums[lo]`.

```go
func findMin(nums []int) int {
    lo, hi := 0, len(nums)-1
    for lo < hi {
        mid := lo + (hi-lo)/2
        if nums[mid] > nums[hi] {
            lo = mid + 1
        } else {
            hi = mid
        }
    }
    return nums[lo]
}
```
</details>

### Drill 4: Capacity to Ship Packages Within D Days (LC 1011) -- 12 min
Weights array, D days. Find minimum ship capacity. This is search on answer. Define the search space bounds, write the feasibility function, write the binary search.

<details>
<summary>Pattern & key decision</summary>

`lo = max(weights)` (must carry the heaviest single package). `hi = sum(weights)` (ship everything in one day). Feasibility: greedily load packages; if current day exceeds capacity, start a new day.

```go
func shipWithinDays(weights []int, days int) int {
    lo, hi := 0, 0
    for _, w := range weights {
        if w > lo {
            lo = w
        }
        hi += w
    }

    for lo < hi {
        mid := lo + (hi-lo)/2
        if canShip(weights, mid, days) {
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}

func canShip(weights []int, capacity, days int) bool {
    daysNeeded, current := 1, 0
    for _, w := range weights {
        if current+w > capacity {
            daysNeeded++
            current = 0
        }
        current += w
    }
    return daysNeeded <= days
}
```
</details>

### Drill 5: Find Peak Element (LC 162) -- 5 min
Write from memory. Trace through `[1, 2, 1, 3, 5, 6, 4]`.

<details>
<summary>Pattern & key decision</summary>

Compare `nums[mid]` to `nums[mid+1]`. If descending, peak is at mid or left (`hi = mid`). If ascending, peak is to the right (`lo = mid + 1`).

```go
func findPeakElement(nums []int) int {
    lo, hi := 0, len(nums)-1
    for lo < hi {
        mid := lo + (hi-lo)/2
        if nums[mid] > nums[mid+1] {
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}
```

Trace: `lo=0, hi=6, mid=3, nums[3]=3 < nums[4]=5 -> lo=4`. `lo=4, hi=6, mid=5, nums[5]=6 > nums[6]=4 -> hi=5`. `lo=4, hi=5, mid=4, nums[4]=5 < nums[5]=6 -> lo=5`. `lo==hi==5`. Return 5. (nums[5]=6 is a peak.)
</details>

---

## Self-Assessment (5 Questions)

Answer these without looking at the guide. If you can't answer confidently, revisit that section.

**1.** You're writing a lower bound function and your loop is `for lo <= hi`. What goes wrong and why?

<details>
<summary>Answer</summary>

When `lo == hi`, you've found the answer, but the loop runs one more iteration. If `nums[mid] >= target`, you set `hi = mid`, which doesn't change anything (`hi` is already `mid` since `lo == hi == mid`). Infinite loop. The lower bound template requires `lo < hi` because the answer is where `lo` and `hi` converge.
</details>

**2.** In the rotated sorted array search, why is the check `nums[lo] <= nums[mid]` (with `<=`) instead of `nums[lo] < nums[mid]`?

<details>
<summary>Answer</summary>

When `lo == mid` (which happens in a 2-element window like `lo=3, hi=4, mid=3`), the left "half" is just one element and IS sorted. With strict `<`, you'd incorrectly conclude the left half is NOT sorted and fall into the wrong branch. The `=` case handles the single-element left half.
</details>

**3.** For the "capacity to ship packages within D days" problem, why is `lo = max(weights)` and not `lo = 1`?

<details>
<summary>Answer</summary>

The ship must be able to carry any single package. If `lo = 1` and the heaviest package weighs 100, any capacity below 100 is infeasible because you can't split a single package across days. Setting `lo = max(weights)` is the tightest correct lower bound, making the search faster and avoiding edge cases in the feasibility function.
</details>

**4.** You're solving a new problem and the brute force is: "try every integer value of X from 1 to 10^9 and return the smallest X where `check(X)` is true." How do you optimize this, and what property must `check` have?

<details>
<summary>Answer</summary>

Binary search on answer space, converting O(10^9 * cost_of_check) to O(log(10^9) * cost_of_check) = O(30 * cost_of_check). The `check` function must be monotonic: once it becomes true at some value X, it stays true for all values > X (or vice versa). Without monotonicity, binary search can't eliminate half the space at each step.
</details>

**5.** Write the ceiling division expression for `a / b` using only integer arithmetic in Go. Why not use `float64` conversion?

<details>
<summary>Answer</summary>

`(a + b - 1) / b`. Float conversion (`int(math.Ceil(float64(a) / float64(b)))`) is dangerous because `float64` has 53 bits of mantissa -- for large integers near `10^18`, you lose precision and get wrong answers. Integer arithmetic is exact and faster.
</details>

---

**End of Day 4.** If the self-assessment went well, you're interview-ready on binary search. If not, re-drill the templates until you can write them without thinking -- because in an interview, your brain needs to focus on the problem, not the template.
