# Day 9 — Binary Search: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Why It's Useful | Category | Time |
|---|----------|----------------|----------|------|
| 1 | [Binary Search — Algorithm Visualizer](https://algorithm-visualizer.org/brute-force/binary-search) | Interactive step-by-step animation showing lo, mid, hi pointers narrowing on a sorted array. You can modify the array and target value and watch each iteration. Best for building the visual mental model before coding. | Interact | 10 min |
| 2 | [Extra, Extra — Read All About It: Nearly All Binary Searches and Mergesorts are Broken (Joshua Bloch, Google AI Blog)](https://research.google/blog/extra-extra-read-all-about-it-nearly-all-binary-searches-and-mergesorts-are-broken/) | The legendary 2006 post revealing that the `mid = (lo + hi) / 2` overflow bug existed in the JDK for 9 years. Demonstrates why `mid = lo + (hi - lo) / 2` matters. Written by the author of *Effective Java*. Bentley's *Programming Pearls* noted that 90% of professional programmers couldn't write a correct binary search — this post shows why. | Read | 10 min |
| 3 | [Go sort.Search documentation and source](https://pkg.go.dev/sort#Search) | Go's standard library binary search. `sort.Search(n, func(i int) bool) int` returns the smallest index i in [0, n) for which f(i) is true, assuming f partitions the range into false, false, ..., true, true. This is exactly the lower-bound template. Understanding this signature is key to using it idiomatically. | Reference | 10 min |
| 4 | [Binary Search on Answer — Codeforces EDU](https://codeforces.com/edu/course/2/lesson/6) | A structured course on the "binary search on answer" paradigm with 6 practice problems. Covers the conceptual shift from "searching an array" to "searching a range of possible answers" with a feasibility predicate. The best systematic treatment of this pattern. | Interact | 20 min |
| 5 | [Binary Search 101 — LeetCode Discuss (zhijun_liao)](https://leetcode.com/discuss/general-discussion/786126/Python-Powerful-Ultimate-Binary-Search-Template) | A popular post that distills binary search into one universal template (`lo < hi` with `hi = mid` or `lo = mid + 1`). Language-agnostic reasoning. Shows how lower bound, upper bound, rotated array, and search-on-answer all reduce to the same skeleton. | Read | 15 min |
| 6 | [Topcoder — Binary Search Tutorial](https://www.topcoder.com/thrive/articles/Binary%20Search) | Thorough walkthrough of the two templates (exact match vs boundary finding), with careful invariant analysis. Includes worked examples for "first true" and "last true" variants and explains how to avoid infinite loops through boundary update rules. | Read | 15 min |
| 7 | [Koko Eating Bananas — NeetCode (YouTube)](https://www.youtube.com/watch?v=U2SozAs9RzA) | Video walkthrough of the "binary search on answer" paradigm using LeetCode 875. Shows how to define the search space (min/max eating speed), write the feasibility function, and apply the lower-bound template. Concrete example before you implement it yourself. | Watch | 10 min |

**Reading strategy:** Start with resource 1 for visual intuition. Read resource 2 to understand why binary search is easy to get wrong — this will motivate careful attention to templates. Study resource 3 to understand Go's `sort.Search` semantics. Read resources 5 and 6 for template mastery. Skim resource 4 and watch resource 7 before implementing the search-on-answer problem.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 – 12:20 | Review Concepts (No Code)

| Time | Min | Activity |
|------|-----|----------|
| 12:00 | 8 | **Read the Day 9 section in OVERVIEW.md.** Study the complexity table. Internalize: standard binary search is O(log n) on sorted arrays. Lower/upper bound use `lo < hi` and converge to a boundary. Search on answer is O(log(range) * check). Know the two templates cold before touching code. |
| 12:08 | 7 | **Open the Algorithm Visualizer (resource 1).** Run a standard binary search on `[1, 3, 5, 7, 9, 11, 13, 15]` for target 7 — watch lo, mid, hi converge. Then search for target 6 (not present) — observe where lo and hi end up when the element is missing. Then mentally trace lower bound on `[1, 3, 5, 5, 5, 7, 9]` for target 5 — where should it return? |
| 12:15 | 5 | **Study the ASCII diagrams in Section 6 below.** Trace the standard binary search step by step. Then study the lower-bound vs upper-bound diagram on duplicates. Finally, study the rotated array diagram — identify which half is always sorted. |

### 12:20 – 12:45 | Implement Standard Binary Search + Lower/Upper Bound (25 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 12:20 | 8 | **Implement standard binary search.** Create `bsearch.go`. Write `BinarySearch(arr []int, target int) int` using the `lo <= hi` template. Return the index if found, -1 if not. Use `mid = lo + (hi-lo)/2`. Write tests in `bsearch_test.go`: found at start, middle, end; not found (too small, too large, between elements); empty array; single element. |
| 12:28 | 8 | **Implement LowerBound.** Write `LowerBound(arr []int, target int) int` — returns the smallest index where `arr[i] >= target`. Use the `lo < hi` template: if `arr[mid] >= target` then `hi = mid`, else `lo = mid + 1`. Returns `len(arr)` if all elements are smaller. Test: no duplicates, all duplicates of target, target not present (returns insertion point), empty array. |
| 12:36 | 5 | **Implement UpperBound.** Write `UpperBound(arr []int, target int) int` — returns the smallest index where `arr[i] > target`. Identical structure to LowerBound but change `>=` to `>`. Test: verify that `UpperBound - LowerBound` equals the count of target in the array. |
| 12:41 | 4 | **Verify against sort.Search.** Write a quick test that confirms `LowerBound(arr, target)` returns the same value as `sort.Search(len(arr), func(i int) bool { return arr[i] >= target })`. Understand why they're equivalent. |

### 12:45 – 1:15 | Rotated Array Variants (30 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 12:45 | 15 | **Implement SearchRotated.** Write `SearchRotated(arr []int, target int) int` — search for target in a rotated sorted array (no duplicates). The key insight: after computing mid, one half is always sorted. Check if target falls in the sorted half; if yes, search there; otherwise search the other half. Use `lo <= hi` since you return inside the loop on exact match. Test: target in left sorted portion, target in right sorted portion, target not present, array not rotated, single element, two elements. |
| 13:00 | 10 | **Implement FindMinRotated.** Write `FindMinRotated(arr []int) int` — find the minimum element in a rotated sorted array. Use `lo < hi` template. If `arr[mid] > arr[hi]`, the min is in the right half (`lo = mid + 1`); otherwise it's in the left half including mid (`hi = mid`). Test: rotated at various positions, not rotated (already sorted), single element, two elements. |
| 13:10 | 5 | **Trace through on paper.** Pick `[4, 5, 6, 7, 0, 1, 2]`, search for 0. Draw lo, mid, hi at each step. Then trace FindMinRotated on the same array. Verify your code matches. |

### 1:15 – 1:45 | Search on Answer + First Bad Version (30 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 1:15 | 5 | **Implement FirstBadVersion.** Write `FirstBadVersion(n int, isBad func(int) bool) int`. This is the purest form of the `lo < hi` template: if `isBad(mid)` then `hi = mid`, else `lo = mid + 1`. Return `lo`. Test with bad version at start, end, middle. |
| 1:20 | 20 | **Implement KokoEatingBananas (search on answer).** Write `MinEatingSpeed(piles []int, h int) int`. The answer space is [1, max(piles)]. For each candidate speed `mid`, write a feasibility function `canFinish(piles, mid, h) bool` that computes total hours needed at that speed. If feasible, `hi = mid`; else `lo = mid + 1`. Return `lo`. Test: single pile, all piles equal, h equals number of piles (must eat each in one hour), h is very large (speed 1 works). |
| 1:40 | 5 | **Reflect on the paradigm.** All five search-on-answer problems share the same shape: define a search range, write a feasibility predicate, apply lower-bound template. Write down this pattern in your own words. |

### 1:45 – 2:00 | Recap

| Time | Min | Activity |
|------|-----|----------|
| 1:45 | 4 | Close all files. Write down from memory: the two templates (`lo <= hi` for exact match, `lo < hi` for boundary finding). The overflow-safe mid calculation. When to use `hi = mid` vs `hi = mid - 1`. |
| 1:49 | 3 | Write down: How do you determine which half of a rotated array is sorted? What is the feasibility function pattern for search-on-answer? |
| 1:52 | 3 | Write down one gotcha you hit during each implementation. |
| 1:55 | 3 | Write down: What does `sort.Search` return when no element satisfies the predicate? How does `UpperBound(arr, t) - LowerBound(arr, t)` give you the count of `t`? |
| 1:58 | 2 | Review the bug taxonomy in Section 5 below. Did you hit any of these bugs during the session? |

---

## 3. Core Concepts Deep Dive

### 3.1 The Two Templates: `lo <= hi` vs `lo < hi`

There are two fundamental binary search templates, and using the wrong one is the #1 source of bugs.

**Template 1: `lo <= hi` — Exact Match (return inside the loop)**

```go
func BinarySearch(arr []int, target int) int {
    lo, hi := 0, len(arr)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if arr[mid] == target {
            return mid          // found — return immediately
        } else if arr[mid] < target {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return -1                   // not found
}
```

**When to use:** You want to find a specific value and return its index. The loop searches until it either finds the target or exhausts the space (`lo > hi`).

**Key properties:**
- `hi` starts at `len(arr) - 1` (valid index)
- Both updates exclude `mid`: `lo = mid + 1`, `hi = mid - 1`
- The loop terminates when `lo > hi` (search space is empty)
- If target appears multiple times, this returns *any* occurrence (not necessarily the first)

**Template 2: `lo < hi` — Boundary Finding (return after the loop)**

```go
func LowerBound(arr []int, target int) int {
    lo, hi := 0, len(arr)
    for lo < hi {
        mid := lo + (hi-lo)/2
        if arr[mid] >= target {
            hi = mid            // mid might be the answer
        } else {
            lo = mid + 1        // mid is definitely not the answer
        }
    }
    return lo                   // lo == hi — the boundary
}
```

**When to use:** You want to find a *boundary* — the first position where a condition becomes true. Lower bound, upper bound, first bad version, and search-on-answer all use this template.

**Key properties:**
- `hi` starts at `len(arr)` (one past the last valid index — the answer can be "past the end")
- Only `lo` excludes `mid` (`lo = mid + 1`); `hi` includes `mid` (`hi = mid`)
- The loop terminates when `lo == hi` — they converge to the boundary
- The invariant: the answer is always in `[lo, hi]`

**Why the asymmetry in updates matters:**
- When `arr[mid] >= target`: mid *could be* the answer (it satisfies the condition), so we keep it in the range: `hi = mid`
- When `arr[mid] < target`: mid is *definitely not* the answer, so we exclude it: `lo = mid + 1`
- This asymmetry is what makes the loop converge without infinite loops

**Decision table:**

| Situation | Template | `lo` init | `hi` init | Loop | Return |
|-----------|----------|-----------|-----------|------|--------|
| Find exact value | `lo <= hi` | 0 | n-1 | `lo <= hi` | inside loop or -1 |
| Find first `>= target` | `lo < hi` | 0 | n | `lo < hi` | `lo` after loop |
| Find first `> target` | `lo < hi` | 0 | n | `lo < hi` | `lo` after loop |
| Find first bad version | `lo < hi` | 1 | n | `lo < hi` | `lo` after loop |
| Search on answer | `lo < hi` | min_ans | max_ans | `lo < hi` | `lo` after loop |

### 3.2 Why `mid = lo + (hi - lo) / 2` Prevents Overflow

The naive formula `mid = (lo + hi) / 2` overflows when `lo + hi` exceeds the maximum integer value.

```
Example with 32-bit signed integers (max = 2,147,483,647):
  lo = 2,000,000,000
  hi = 2,100,000,000

  lo + hi = 4,100,000,000  ← exceeds 2^31 - 1, wraps to negative!
  (lo + hi) / 2 = negative number / 2 = negative  ← WRONG

  lo + (hi - lo) / 2
  = 2,000,000,000 + (100,000,000) / 2
  = 2,000,000,000 + 50,000,000
  = 2,050,000,000  ← CORRECT, no overflow
```

**In Go:** `int` is 64-bit on modern platforms, so overflow at 2^63 is unlikely for array indices. But the habit is important because:
1. It applies universally across languages (Java, C++ still use 32-bit ints for array indexing)
2. It applies to search-on-answer where `lo` and `hi` can be arbitrary large values
3. It's zero cognitive cost once you internalize it

**The Bloch post (2006):** This exact bug existed in Java's `Arrays.binarySearch` for 9 years (1997-2006) and in the C standard library. Jon Bentley's *Programming Pearls* binary search had the same bug. If experts get it wrong, use the safe form.

### 3.3 Lower Bound and Upper Bound: Precise Definitions

Given a sorted array `arr` and a target value:

**Lower Bound** (`bisect_left` in Python, `lower_bound` in C++):
- Returns the **smallest index `i`** such that `arr[i] >= target`
- Equivalently: the leftmost position where `target` could be inserted to keep the array sorted
- If `target` exists, returns the index of its **first occurrence**
- If `target` doesn't exist, returns the index where it *would* go

**Upper Bound** (`bisect_right` in Python, `upper_bound` in C++):
- Returns the **smallest index `i`** such that `arr[i] > target`
- Equivalently: the rightmost position where `target` could be inserted to keep the array sorted
- If `target` exists, returns **one past its last occurrence**
- If `target` doesn't exist, returns the same as lower bound

**The relationship:**

```
arr:    [1, 3, 5, 5, 5, 7, 9]
         0  1  2  3  4  5  6

target = 5:
  LowerBound → 2  (first 5 is at index 2)
  UpperBound → 5  (first element > 5 is at index 5)
  Count of 5 = UpperBound - LowerBound = 5 - 2 = 3  ✓

target = 4 (not present):
  LowerBound → 2  (first element >= 4 is 5 at index 2)
  UpperBound → 2  (first element > 4 is also 5 at index 2)
  Count of 4 = 2 - 2 = 0  ✓

target = 10 (larger than all):
  LowerBound → 7  (past the end — len(arr))
  UpperBound → 7
```

**How they differ by one comparison:**

```go
// LowerBound: first index where arr[i] >= target
if arr[mid] >= target { hi = mid } else { lo = mid + 1 }

// UpperBound: first index where arr[i] > target
if arr[mid] > target  { hi = mid } else { lo = mid + 1 }
//           ^^ only this changes: >= becomes >
```

The entire structural difference between lower bound and upper bound is `>=` vs `>` in the condition. Everything else is identical.

**Go's `sort.Search`** is lower bound generalized: it returns the smallest `i` in `[0, n)` for which `f(i)` is true, assuming `f` is monotone (false, false, ..., true, true). To get:
- Lower bound: `sort.Search(len(arr), func(i int) bool { return arr[i] >= target })`
- Upper bound: `sort.Search(len(arr), func(i int) bool { return arr[i] > target })`

### 3.4 Binary Search on Rotated Arrays

A rotated sorted array is a sorted array that has been rotated at some pivot:

```
Original:  [0, 1, 2, 4, 5, 6, 7]
Rotated:   [4, 5, 6, 7, 0, 1, 2]   (rotated at index 4)
```

**The key insight: when you compute `mid`, at least one half is always sorted.**

```
[4, 5, 6, 7, 0, 1, 2]
 lo          mid       hi
 
arr[lo]=4 <= arr[mid]=7  →  left half [4,5,6,7] is sorted
                             right half [0,1,2] contains the rotation point
```

**Search in rotated array (no duplicates):**

```go
func SearchRotated(arr []int, target int) int {
    lo, hi := 0, len(arr)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if arr[mid] == target {
            return mid
        }
        // Determine which half is sorted
        if arr[lo] <= arr[mid] {
            // Left half [lo..mid] is sorted
            if arr[lo] <= target && target < arr[mid] {
                hi = mid - 1    // target is in the sorted left half
            } else {
                lo = mid + 1    // target is in the other half
            }
        } else {
            // Right half [mid..hi] is sorted
            if arr[mid] < target && target <= arr[hi] {
                lo = mid + 1    // target is in the sorted right half
            } else {
                hi = mid - 1    // target is in the other half
            }
        }
    }
    return -1
}
```

**Why this works:** In a sorted half, you can check if the target falls within its range using a simple comparison. If it does, search there. If it doesn't, the target *must* be in the other half (if it exists at all).

**Find minimum in rotated array:**

```go
func FindMinRotated(arr []int) int {
    lo, hi := 0, len(arr)-1
    for lo < hi {
        mid := lo + (hi-lo)/2
        if arr[mid] > arr[hi] {
            lo = mid + 1    // min is in the right half (past mid)
        } else {
            hi = mid         // mid could be the min
        }
    }
    return arr[lo]
}
```

**Why compare with `arr[hi]` and not `arr[lo]`?** If the array is not rotated (`[1,2,3,4,5]`), comparing with `arr[lo]` would always see `arr[mid] >= arr[lo]` and keep pushing `lo` right, overshooting the minimum. Comparing with `arr[hi]` correctly identifies that when `arr[mid] <= arr[hi]`, the right portion is sorted and the min is at or before `mid`.

### 3.5 Binary Search on Answer Space

This is the paradigm shift: instead of searching through an array, you binary search over the range of possible *answers*.

**The pattern:**
1. Define the search space: what is the minimum possible answer? The maximum?
2. Write a **feasibility function**: given a candidate answer `mid`, can the problem be solved?
3. The feasibility function is monotonic: if `mid` works, then `mid + 1` also works (or vice versa)
4. Apply the lower-bound template to find the boundary where feasibility flips

**Template:**

```go
func searchOnAnswer(/* problem params */) int {
    lo, hi := minPossibleAnswer, maxPossibleAnswer
    for lo < hi {
        mid := lo + (hi-lo)/2
        if feasible(mid /*, problem params */) {
            hi = mid        // mid works — try smaller (if minimizing)
        } else {
            lo = mid + 1    // mid doesn't work — need bigger
        }
    }
    return lo
}
```

**Example — Koko Eating Bananas (LeetCode 875):**

Koko has `piles` of bananas and `h` hours. She picks an eating speed `k` (bananas/hour). Each hour she eats from one pile — if the pile has fewer than `k` bananas, she finishes it and waits. What's the minimum `k` to finish all piles in `h` hours?

```go
func MinEatingSpeed(piles []int, h int) int {
    lo, hi := 1, maxPile(piles)
    for lo < hi {
        mid := lo + (hi-lo)/2
        if canFinish(piles, mid, h) {
            hi = mid        // this speed works — try slower
        } else {
            lo = mid + 1    // too slow — need faster
        }
    }
    return lo
}

func canFinish(piles []int, speed, h int) bool {
    hours := 0
    for _, p := range piles {
        hours += (p + speed - 1) / speed  // ceil division
    }
    return hours <= h
}
```

**The feasibility function as a step function:**

```
speed:       1   2   3   4   5   6   7   8   ...
canFinish:   F   F   F   T   T   T   T   T   ...
                         ^
                    answer = 4 (first true)
```

**Other search-on-answer problems:**
- **Ship packages within D days** (LeetCode 1011): search on ship capacity
- **Split array largest sum** (LeetCode 410): search on the maximum subarray sum
- **Minimum days to make m bouquets** (LeetCode 1482): search on the number of days
- **Magnetic force between balls** (LeetCode 1552): search on the minimum distance

All follow the identical template. The only thing that changes is the feasibility function.

### 3.6 Go's `sort.Search`: Signature, Semantics, and Lower Bound Mapping

```go
func Search(n int, f func(int) bool) int
```

**Semantics:** `sort.Search` uses binary search to find the smallest index `i` in `[0, n)` for which `f(i)` is true, assuming that `f` is monotone: `f(i) == true` implies `f(i+1) == true`.

**If no index satisfies `f`, it returns `n`.**

**Implementation (simplified):**

```go
func Search(n int, f func(int) bool) int {
    lo, hi := 0, n
    for lo < hi {
        mid := lo + (hi-lo)/2
        if f(mid) {
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}
```

This is exactly the `lo < hi` boundary-finding template.

**Mapping to lower/upper bound:**

```go
// Lower bound: first index where arr[i] >= target
idx := sort.Search(len(arr), func(i int) bool {
    return arr[i] >= target
})

// Upper bound: first index where arr[i] > target
idx := sort.Search(len(arr), func(i int) bool {
    return arr[i] > target
})

// Exact match: lower bound, then check
idx := sort.Search(len(arr), func(i int) bool {
    return arr[i] >= target
})
if idx < len(arr) && arr[idx] == target {
    // found at idx
}

// Search on answer (e.g., minimum speed):
speed := sort.Search(maxSpeed, func(k int) bool {
    return canFinish(piles, k+1, h)  // +1 because Search is 0-indexed
})
// speed+1 is the answer
```

**Gotcha with `sort.Search` indexing:** The function parameter `f` receives indices `0` through `n-1`. If your answer space starts at 1 (like eating speed), you need to offset: either search `[0, max-1]` and add 1 to the result, or adjust inside the closure.

---

## 4. Implementation Checklist

### Function Signatures

```go
package bsearch

// BinarySearch returns the index of target in sorted arr, or -1 if not found.
// Uses lo <= hi template (exact match, return inside loop).
// Time: O(log n)  Space: O(1)
func BinarySearch(arr []int, target int) int

// LowerBound returns the smallest index i such that arr[i] >= target.
// Returns len(arr) if all elements are less than target.
// Uses lo < hi template (boundary finding, return after loop).
// Time: O(log n)  Space: O(1)
func LowerBound(arr []int, target int) int

// UpperBound returns the smallest index i such that arr[i] > target.
// Returns len(arr) if all elements are <= target.
// Time: O(log n)  Space: O(1)
func UpperBound(arr []int, target int) int

// SearchRotated returns the index of target in a rotated sorted array
// (no duplicates), or -1 if not found.
// Time: O(log n)  Space: O(1)
func SearchRotated(arr []int, target int) int

// FindMinRotated returns the minimum element in a rotated sorted array
// (no duplicates).
// Time: O(log n)  Space: O(1)
func FindMinRotated(arr []int) int

// FirstBadVersion returns the first bad version in [1, n].
// isBad(v) returns true iff version v is bad.
// All versions after the first bad one are also bad.
// Time: O(log n)  Space: O(1)
func FirstBadVersion(n int, isBad func(int) bool) int

// MinEatingSpeed returns the minimum eating speed k such that Koko can
// finish all piles within h hours. Each hour she eats min(k, pile) bananas
// from one pile.
// Time: O(n * log(max(piles)))  Space: O(1)
func MinEatingSpeed(piles []int, h int) int
```

### Test Cases

```go
package bsearch

import (
    "sort"
    "testing"
)

func TestBinarySearch(t *testing.T) {
    tests := []struct {
        name   string
        arr    []int
        target int
        want   int
    }{
        {"found at start", []int{1, 3, 5, 7, 9}, 1, 0},
        {"found at end", []int{1, 3, 5, 7, 9}, 9, 4},
        {"found in middle", []int{1, 3, 5, 7, 9}, 5, 2},
        {"not found - too small", []int{1, 3, 5, 7, 9}, 0, -1},
        {"not found - too large", []int{1, 3, 5, 7, 9}, 10, -1},
        {"not found - between elements", []int{1, 3, 5, 7, 9}, 4, -1},
        {"empty array", []int{}, 5, -1},
        {"single element - found", []int{42}, 42, 0},
        {"single element - not found", []int{42}, 7, -1},
        {"two elements - find first", []int{1, 2}, 1, 0},
        {"two elements - find second", []int{1, 2}, 2, 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := BinarySearch(tt.arr, tt.target)
            if got != tt.want {
                t.Errorf("BinarySearch(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
            }
        })
    }
}

func TestLowerBound(t *testing.T) {
    tests := []struct {
        name   string
        arr    []int
        target int
        want   int
    }{
        {"first occurrence", []int{1, 3, 5, 5, 5, 7, 9}, 5, 2},
        {"not present - insertion point", []int{1, 3, 5, 7, 9}, 4, 2},
        {"smaller than all", []int{1, 3, 5, 7, 9}, 0, 0},
        {"larger than all", []int{1, 3, 5, 7, 9}, 10, 5},
        {"empty array", []int{}, 5, 0},
        {"all same - target matches", []int{5, 5, 5, 5}, 5, 0},
        {"all same - target larger", []int{5, 5, 5, 5}, 6, 4},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := LowerBound(tt.arr, tt.target)
            if got != tt.want {
                t.Errorf("LowerBound(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
            }
        })
    }
}

func TestUpperBound(t *testing.T) {
    tests := []struct {
        name   string
        arr    []int
        target int
        want   int
    }{
        {"past last occurrence", []int{1, 3, 5, 5, 5, 7, 9}, 5, 5},
        {"not present", []int{1, 3, 5, 7, 9}, 4, 2},
        {"count via bounds", []int{1, 3, 5, 5, 5, 7, 9}, 5, 5},
        // UpperBound - LowerBound = 5 - 2 = 3 occurrences of 5
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := UpperBound(tt.arr, tt.target)
            if got != tt.want {
                t.Errorf("UpperBound(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
            }
        })
    }
}

func TestLowerBoundMatchesSortSearch(t *testing.T) {
    arr := []int{1, 2, 4, 4, 4, 6, 8, 10}
    for target := 0; target <= 12; target++ {
        got := LowerBound(arr, target)
        want := sort.Search(len(arr), func(i int) bool { return arr[i] >= target })
        if got != want {
            t.Errorf("target=%d: LowerBound=%d, sort.Search=%d", target, got, want)
        }
    }
}

func TestSearchRotated(t *testing.T) {
    tests := []struct {
        name   string
        arr    []int
        target int
        want   int
    }{
        {"found in left sorted", []int{4, 5, 6, 7, 0, 1, 2}, 5, 1},
        {"found in right sorted", []int{4, 5, 6, 7, 0, 1, 2}, 1, 5},
        {"found at pivot", []int{4, 5, 6, 7, 0, 1, 2}, 0, 4},
        {"not found", []int{4, 5, 6, 7, 0, 1, 2}, 3, -1},
        {"not rotated", []int{1, 2, 3, 4, 5}, 3, 2},
        {"single element found", []int{1}, 1, 0},
        {"single element not found", []int{1}, 0, -1},
        {"two elements", []int{2, 1}, 1, 1},
        {"empty", []int{}, 1, -1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := SearchRotated(tt.arr, tt.target)
            if got != tt.want {
                t.Errorf("SearchRotated(%v, %d) = %d, want %d", tt.arr, tt.target, got, tt.want)
            }
        })
    }
}

func TestFindMinRotated(t *testing.T) {
    tests := []struct {
        name string
        arr  []int
        want int
    }{
        {"rotated", []int{4, 5, 6, 7, 0, 1, 2}, 0},
        {"rotated at end", []int{2, 3, 4, 5, 1}, 1},
        {"not rotated", []int{1, 2, 3, 4, 5}, 1},
        {"single", []int{1}, 1},
        {"two elements rotated", []int{2, 1}, 1},
        {"two elements sorted", []int{1, 2}, 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := FindMinRotated(tt.arr)
            if got != tt.want {
                t.Errorf("FindMinRotated(%v) = %d, want %d", tt.arr, got, tt.want)
            }
        })
    }
}

func TestFirstBadVersion(t *testing.T) {
    tests := []struct {
        name    string
        n       int
        firstBad int
    }{
        {"bad at start", 5, 1},
        {"bad at end", 5, 5},
        {"bad in middle", 10, 4},
        {"single version bad", 1, 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            isBad := func(v int) bool { return v >= tt.firstBad }
            got := FirstBadVersion(tt.n, isBad)
            if got != tt.firstBad {
                t.Errorf("FirstBadVersion(%d) = %d, want %d", tt.n, got, tt.firstBad)
            }
        })
    }
}

func TestMinEatingSpeed(t *testing.T) {
    tests := []struct {
        name  string
        piles []int
        h     int
        want  int
    }{
        {"example 1", []int{3, 6, 7, 11}, 8, 4},
        {"example 2", []int{30, 11, 23, 4, 20}, 5, 30},
        {"example 3", []int{30, 11, 23, 4, 20}, 6, 23},
        {"single pile", []int{10}, 5, 2},
        {"h equals piles", []int{3, 6, 7, 11}, 4, 11},
        {"large h", []int{3, 6, 7, 11}, 100, 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MinEatingSpeed(tt.piles, tt.h)
            if got != tt.want {
                t.Errorf("MinEatingSpeed(%v, %d) = %d, want %d", tt.piles, tt.h, got, tt.want)
            }
        })
    }
}
```

**Edge cases to cover explicitly:**
- **Empty array** — should return -1 (search) or 0 (lower/upper bound), not panic
- **Single element** — boundary condition for all variants
- **Target at boundaries** — first element, last element
- **Target not present** — between elements, smaller than all, larger than all
- **All duplicates** — lower bound should return 0, upper bound should return n
- **Rotated array not actually rotated** — already sorted is a valid rotation (by 0)
- **Search on answer with trivial inputs** — speed 1 always works if h is large enough

---

## 5. The Binary Search Bug Taxonomy

### Bug 1: Integer Overflow in Mid Calculation

**Broken:**
```go
mid := (lo + hi) / 2
```

**What goes wrong:** When `lo + hi` exceeds `math.MaxInt64` (or `MaxInt32` in other languages), the addition wraps to a negative number. The resulting `mid` is negative, causing an out-of-bounds panic or an infinite loop.

**Fix:**
```go
mid := lo + (hi-lo)/2
```

**When it matters:** Primarily in search-on-answer where `lo` and `hi` can be billions. For array indices in Go (64-bit int), overflow is unlikely but not impossible if `hi` is set to `math.MaxInt64`.

---

### Bug 2: Wrong Loop Condition (`lo <= hi` vs `lo < hi`)

**Broken — using `lo <= hi` for boundary finding:**
```go
func LowerBound(arr []int, target int) int {
    lo, hi := 0, len(arr)
    for lo <= hi {              // BUG: should be lo < hi
        mid := lo + (hi-lo)/2
        if arr[mid] >= target { // PANIC: arr[len(arr)] is out of bounds
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}
```

**What goes wrong:** When `hi` is initialized to `len(arr)` and the loop runs while `lo <= hi`, mid can equal `len(arr)`, causing an index-out-of-range panic on `arr[mid]`. Even without the panic, the loop may not terminate correctly.

**Fix:** Use `lo < hi`. When `hi = len(arr)`, the loop condition `lo < hi` prevents accessing `arr[hi]` directly (mid will be at most `hi - 1` when `lo < hi`).

**Broken — using `lo < hi` for exact match:**
```go
func BinarySearch(arr []int, target int) int {
    lo, hi := 0, len(arr)-1
    for lo < hi {               // BUG: should be lo <= hi
        mid := lo + (hi-lo)/2
        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return -1                   // BUG: misses the case where lo == hi == target index
}
```

**What goes wrong:** When `lo == hi`, the loop exits without checking `arr[lo]`. If the target is at that position, you return -1 incorrectly.

**Fix:** Use `lo <= hi` so the single-element case is checked.

---

### Bug 3: Off-by-One in Boundary Update (`lo = mid` vs `lo = mid + 1`)

**Broken — using `lo = mid` with integer division rounding down:**
```go
func LowerBound(arr []int, target int) int {
    lo, hi := 0, len(arr)
    for lo < hi {
        mid := lo + (hi-lo)/2
        if arr[mid] >= target {
            hi = mid
        } else {
            lo = mid            // BUG: should be mid + 1
        }
    }
    return lo
}
```

**What goes wrong — infinite loop:**
```
arr = [1, 3], target = 3
Iteration 1: lo=0, hi=2, mid=1, arr[1]=3 >= 3 → hi=1
Iteration 2: lo=0, hi=1, mid=0, arr[0]=1 < 3  → lo=0  (mid is 0, lo stays 0)
Iteration 3: lo=0, hi=1, mid=0, arr[0]=1 < 3  → lo=0  ← INFINITE LOOP
```

When `hi = lo + 1`, `mid = lo + (hi-lo)/2 = lo + 0 = lo`. If `lo = mid`, lo never advances. The search space `[lo, hi)` never shrinks.

**Fix:** `lo = mid + 1`. Since we've determined that `mid` is *not* the answer (the condition was false), we exclude it.

**The rule:** With integer division rounding down (`mid` biased toward `lo`):
- `lo = mid` is dangerous — can cause infinite loops
- `lo = mid + 1` is safe — always shrinks the search space
- `hi = mid` is safe — always shrinks the search space (since `mid < hi` when `lo < hi`)
- `hi = mid - 1` is used only in the `lo <= hi` template where you've already checked `mid`

---

### Bug 4: Infinite Loop When `lo = mid` and `hi = lo + 1`

This is the same root cause as Bug 3, but it's common enough to call out separately. It appears in a different template variant:

**Broken — "find last" with wrong mid calculation:**
```go
// Trying to find the LAST index where arr[i] <= target
func LastLE(arr []int, target int) int {
    lo, hi := 0, len(arr)-1
    for lo < hi {
        mid := lo + (hi-lo)/2       // rounds DOWN toward lo
        if arr[mid] <= target {
            lo = mid                 // BUG: infinite loop when hi = lo + 1
        } else {
            hi = mid - 1
        }
    }
    return lo
}
```

**Fix — round UP when `lo = mid`:**
```go
func LastLE(arr []int, target int) int {
    lo, hi := 0, len(arr)-1
    for lo < hi {
        mid := lo + (hi-lo+1)/2     // rounds UP — crucial!
        if arr[mid] <= target {
            lo = mid                 // safe now: mid > lo when hi > lo
        } else {
            hi = mid - 1
        }
    }
    return lo
}
```

**The rule:** If your update sets `lo = mid` (without +1), use `mid = lo + (hi-lo+1)/2` (ceiling division) to ensure `mid > lo` and the space shrinks. If your update sets `hi = mid` (without -1), use the standard `mid = lo + (hi-lo)/2` (floor division).

---

### Bug 5: Not Handling Empty Arrays

**Broken:**
```go
func BinarySearch(arr []int, target int) int {
    lo, hi := 0, len(arr)-1         // hi = -1 for empty array
    for lo <= hi {                   // 0 <= -1 is false — loop doesn't execute ✓
        // ...
    }
    return -1                        // Returns -1 — correct for exact match
}

func FindMinRotated(arr []int) int {
    lo, hi := 0, len(arr)-1         // hi = -1 for empty array
    for lo < hi {
        mid := lo + (hi-lo)/2
        if arr[mid] > arr[hi] {     // PANIC: arr[-1] out of bounds
            lo = mid + 1
        } else {
            hi = mid
        }
    }
    return arr[lo]                   // PANIC: arr[0] out of bounds
}
```

**Fix:** Guard at the top of every function:
```go
func FindMinRotated(arr []int) int {
    if len(arr) == 0 {
        return -1  // or panic, depending on contract
    }
    // ... rest of implementation
}
```

**The rule:** Every binary search function should either:
1. Document that the input must be non-empty (and panic if violated), or
2. Handle `len(arr) == 0` as the first line

For `lo <= hi` with `hi = len(arr) - 1`, the empty case naturally returns -1 (loop doesn't execute). For `lo < hi` with `hi = len(arr)`, the empty case returns `lo = 0` which may or may not be correct for your use case.

---

### Summary Table

| Bug | Root Cause | Symptom | Prevention |
|-----|-----------|---------|------------|
| Overflow | `(lo + hi) / 2` | Negative mid, panic | Always `lo + (hi-lo)/2` |
| Wrong loop condition | `<=` vs `<` mismatch | Misses elements or panics | Match template to problem type |
| `lo = mid` without +1 | Search space doesn't shrink | Infinite loop | `lo = mid + 1` (floor div) or ceiling div |
| `lo = mid` with floor div | `mid == lo` when `hi = lo+1` | Infinite loop on 2 elements | Use ceiling: `lo + (hi-lo+1)/2` |
| Empty array | No length guard | Index out of range panic | Check `len(arr) == 0` first |

---

## 6. Visual Diagrams

### 6.1 Standard Binary Search: Narrowing on a Sorted Array

Search for target = 7 in `[1, 3, 5, 7, 9, 11, 13]`:

```
Iteration 1:
  [1, 3, 5, 7, 9, 11, 13]
   lo         mid          hi
   0          3            6
   arr[3] = 7 == target → FOUND at index 3 ✓

Now search for target = 9:

Iteration 1:
  [1, 3, 5, 7, 9, 11, 13]
   lo         mid          hi
   0          3            6
   arr[3] = 7 < 9 → lo = mid + 1 = 4

Iteration 2:
  [1, 3, 5, 7, 9, 11, 13]
                  lo  mid   hi
                  4   5     6
   arr[5] = 11 > 9 → hi = mid - 1 = 4

Iteration 3:
  [1, 3, 5, 7, 9, 11, 13]
                  lo
                  mid
                  hi
                  4
   arr[4] = 9 == target → FOUND at index 4 ✓

Now search for target = 6 (NOT PRESENT):

Iteration 1:
  [1, 3, 5, 7, 9, 11, 13]
   lo         mid          hi
   0          3            6
   arr[3] = 7 > 6 → hi = mid - 1 = 2

Iteration 2:
  [1, 3, 5, 7, 9, 11, 13]
   lo  mid  hi
   0   1    2
   arr[1] = 3 < 6 → lo = mid + 1 = 2

Iteration 3:
  [1, 3, 5, 7, 9, 11, 13]
           lo
           mid
           hi
           2
   arr[2] = 5 < 6 → lo = mid + 1 = 3

   Now lo (3) > hi (2) → loop exits → return -1 (not found)
```

### 6.2 Lower Bound vs Upper Bound on an Array with Duplicates

Array: `[1, 3, 5, 5, 5, 7, 9]`, target = 5

```
Lower Bound (first index where arr[i] >= 5):

  [1, 3, 5, 5, 5, 7, 9]
   lo         mid       hi        lo=0, hi=7 (past end)
   0          3         7
   arr[3]=5 >= 5 → hi = 3

  [1, 3, 5, 5, 5, 7, 9]
   lo  mid  hi
   0   1    3
   arr[1]=3 < 5 → lo = 2

  [1, 3, 5, 5, 5, 7, 9]
         lo hi
         mid
         2  3
   arr[2]=5 >= 5 → hi = 2

   lo == hi == 2 → return 2
                             ↓
  [1, 3, 5, 5, 5, 7, 9]
         ▲
   LowerBound = 2 (first 5)


Upper Bound (first index where arr[i] > 5):

  [1, 3, 5, 5, 5, 7, 9]
   lo         mid       hi        lo=0, hi=7
   0          3         7
   arr[3]=5, NOT > 5 → lo = 4

  [1, 3, 5, 5, 5, 7, 9]
                  lo mid hi
                  4  5   7
   arr[5]=7 > 5 → hi = 5

  [1, 3, 5, 5, 5, 7, 9]
                  lo hi
                  mid
                  4  5
   arr[4]=5, NOT > 5 → lo = 5

   lo == hi == 5 → return 5
                                   ↓
  [1, 3, 5, 5, 5, 7, 9]
                     ▲
   UpperBound = 5 (first element after all 5s)


  Count of 5 = UpperBound - LowerBound = 5 - 2 = 3 ✓

  Visual summary:
  Index:  0  1  2  3  4  5  6
  Value:  1  3  5  5  5  7  9
              ▲ LB         ▲ UB
              │             │
              └── 3 fives ──┘
```

### 6.3 Rotated Sorted Array: Two Sorted Halves

```
Original sorted array:
  [0, 1, 2, 4, 5, 6, 7]

Rotated at index 4 (rotation point = original index of minimum):
  [4, 5, 6, 7, 0, 1, 2]
   ├────────┤  ├──────┤
   left sorted  right sorted
   (>= pivot)   (< pivot)

The key insight: after computing mid, ONE half is always sorted.

Case 1: mid lands in left sorted half
  [4, 5, 6, 7, 0, 1, 2]
   lo      mid         hi
   arr[lo]=4 <= arr[mid]=7 → LEFT half is sorted
   ├─sorted─┤  ├─unsorted─┤

   If target is in [arr[lo], arr[mid]):
       search left  → hi = mid - 1
   Else:
       search right → lo = mid + 1


Case 2: mid lands in right sorted half
  [6, 7, 0, 1, 2, 4, 5]
   lo      mid         hi
   arr[lo]=6 > arr[mid]=1 → RIGHT half is sorted
   ├unsorted┤  ├─sorted──┤

   If target is in (arr[mid], arr[hi]]:
       search right → lo = mid + 1
   Else:
       search left  → hi = mid - 1


Finding the minimum — always compare mid with hi:
  [4, 5, 6, 7, 0, 1, 2]
   lo      mid         hi

   arr[mid]=7 > arr[hi]=2 → min is in right half → lo = mid + 1

  [4, 5, 6, 7, 0, 1, 2]
                  lo mid hi

   arr[mid]=1 <= arr[hi]=2 → min is at mid or left → hi = mid

  [4, 5, 6, 7, 0, 1, 2]
                  lo
                  hi
   lo == hi → return arr[lo] = 0 ✓
```

### 6.4 Search on Answer: The Feasibility Function as a Step Function

Koko Eating Bananas: piles = [3, 6, 7, 11], h = 8

```
For each candidate speed k, compute total hours needed:
  hours(k) = ceil(3/k) + ceil(6/k) + ceil(7/k) + ceil(11/k)

  k=1:  3 + 6 + 7 + 11 = 27 hours  (need <= 8)  → INFEASIBLE
  k=2:  2 + 3 + 4 + 6  = 15 hours                → INFEASIBLE
  k=3:  1 + 2 + 3 + 4  = 10 hours                → INFEASIBLE
  k=4:  1 + 2 + 2 + 3  = 8  hours                → FEASIBLE ✓
  k=5:  1 + 2 + 2 + 3  = 8  hours                → FEASIBLE ✓
  k=6:  1 + 1 + 2 + 2  = 6  hours                → FEASIBLE ✓
  ...
  k=11: 1 + 1 + 1 + 1  = 4  hours                → FEASIBLE ✓

The feasibility predicate forms a STEP FUNCTION:

  Speed (k):  1   2   3   4   5   6   7   8   9  10  11
  feasible:   F   F   F   T   T   T   T   T   T   T   T
                          ▲
                     first T = answer = 4

This is exactly a lower-bound search on the answer space [1, 11]:

  lo=1, hi=11
  ┌──────────────────────────────────────────────────┐
  │  F   F   F   T   T   T   T   T   T   T   T      │
  │  lo              mid                        hi   │
  │                  mid=6, feasible → hi=6          │
  │                                                  │
  │  F   F   F   T   T   T                          │
  │  lo      mid         hi                          │
  │          mid=3, infeasible → lo=4                │
  │                                                  │
  │              T   T   T                           │
  │              lo  mid hi                          │
  │              mid=5, feasible → hi=5              │
  │                                                  │
  │              T   T                               │
  │              lo  hi                              │
  │              mid=4, feasible → hi=4              │
  │                                                  │
  │              T                                   │
  │              lo=hi=4 → return 4 ✓                │
  └──────────────────────────────────────────────────┘
```

---

## 7. Self-Assessment

Answer these from memory after your session. If you can't, that's tomorrow's priority.

### Q1: Given a sorted array with duplicates, what's the difference between what lower bound and upper bound return?

<details>
<summary>Answer</summary>

**Lower bound** returns the index of the **first occurrence** of the target (or the insertion point if not present) — the smallest `i` where `arr[i] >= target`.

**Upper bound** returns **one past the last occurrence** of the target — the smallest `i` where `arr[i] > target`.

The count of the target in the array is `UpperBound - LowerBound`. If the target is not present, both return the same index and the count is 0.

The only code difference is the comparison: `>=` for lower bound, `>` for upper bound.
</details>

### Q2: How do you know which half of a rotated sorted array is sorted?

<details>
<summary>Answer</summary>

After computing `mid`, compare `arr[lo]` with `arr[mid]`:

- If `arr[lo] <= arr[mid]`: the **left half** `[lo..mid]` is sorted. The rotation point is in the right half.
- If `arr[lo] > arr[mid]`: the **right half** `[mid..hi]` is sorted. The rotation point is in the left half.

At least one half is always sorted because the array was originally fully sorted and then rotated at a single point. The rotation breaks exactly one contiguous sorted region into two.

Once you know which half is sorted, check if the target falls within that sorted range. If it does, search there. If not, search the other half.
</details>

### Q3: Why does `lo = mid` (without `+ 1`) cause an infinite loop, and when exactly does it happen?

<details>
<summary>Answer</summary>

It happens when `hi = lo + 1` (the search space has exactly two elements).

With floor division: `mid = lo + (hi - lo) / 2 = lo + 0 = lo`. Setting `lo = mid` means `lo` stays at `lo` — the search space `[lo, hi)` doesn't shrink.

The fix is one of:
1. Always use `lo = mid + 1` (exclude mid because the condition proved it's not the answer)
2. If you must use `lo = mid`, use ceiling division: `mid = lo + (hi - lo + 1) / 2` so that `mid > lo` when `hi > lo`

The standard lower-bound template avoids this entirely: `hi = mid` (safe with floor div) and `lo = mid + 1` (always shrinks).
</details>

### Q4: Explain the "binary search on answer" paradigm. What are the three things you need to define?

<details>
<summary>Answer</summary>

Instead of searching through a data structure, you binary search over the **range of possible answers** to find the optimal one.

You need three things:

1. **Search space bounds**: the minimum and maximum possible answers. (e.g., for Koko eating bananas: min speed = 1, max speed = max pile size)

2. **Feasibility function (predicate)**: given a candidate answer `mid`, a function that returns true/false for "can the problem be solved with this answer?" Must be monotonic — once true, stays true for all larger (or all smaller) values.

3. **Optimization direction**: are you minimizing (find first feasible) or maximizing (find last feasible)? This determines whether `hi = mid` (minimizing) or `lo = mid` (maximizing, with ceiling div).

The feasibility function is O(n) or O(n log n) typically. Total complexity is O(log(range) * cost_of_feasibility_check).
</details>

### Q5: What does Go's `sort.Search(n, f)` return when no index satisfies `f`? How do you use it for exact-match search?

<details>
<summary>Answer</summary>

When no index in `[0, n)` satisfies `f`, `sort.Search` returns `n` (one past the end of the range).

For exact-match search in a sorted slice:

```go
idx := sort.Search(len(arr), func(i int) bool {
    return arr[i] >= target
})
if idx < len(arr) && arr[idx] == target {
    // found at idx
} else {
    // not found
}
```

`sort.Search` finds the lower bound (first position where `arr[i] >= target`). You then check if that position actually contains the target. If `idx == len(arr)`, the target is larger than all elements. If `arr[idx] != target`, the target isn't in the array.
</details>

---

## Complexity Reference (Quick Glance)

| Variant | Time | Space | Template | Notes |
|---------|------|-------|----------|-------|
| Standard search | O(log n) | O(1) | `lo <= hi` | Return inside loop or -1 |
| Lower bound | O(log n) | O(1) | `lo < hi` | First `>= target`; `hi = len(arr)` |
| Upper bound | O(log n) | O(1) | `lo < hi` | First `> target`; `hi = len(arr)` |
| Search in rotated array | O(log n) | O(1) | `lo <= hi` | One half always sorted |
| Find min in rotated array | O(log n) | O(1) | `lo < hi` | Compare `arr[mid]` with `arr[hi]` |
| First bad version | O(log n) | O(1) | `lo < hi` | Pure predicate search |
| Search on answer | O(log R * C) | O(1) | `lo < hi` | R = range, C = feasibility check cost |
| `sort.Search` | O(log n) | O(1) | `lo < hi` | Returns `n` if no match |
