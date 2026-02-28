# Day 8 — Sorting Algorithms: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Why It's Useful | Category | Time |
|---|----------|----------------|----------|------|
| 1 | [sorting.at — Sorting Algorithm Animations](https://sorting.at/) | Side-by-side visual comparison of sorting algorithms on the same input. Watch MergeSort's predictable splits vs QuickSort's pivot-dependent partitions in real time. Run all three (merge, quick, counting) simultaneously. | Interact | 10 min |
| 2 | [VisuAlgo — Sorting](https://visualgo.net/en/sorting) | Step-by-step animation with pseudocode highlighting. Especially valuable for understanding the merge step and Lomuto vs Hoare partition. You can input custom arrays. | Interact | 10 min |
| 3 | [The Comparison Sort Lower Bound (MIT OCW 6.006)](https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/resources/lecture-7-counting-sort-radix-sort-lower-bounds-for-sorting-and-searching/) | Erik Demaine's lecture on why comparison sorts can't beat O(n log n) — the decision tree argument explained clearly with diagrams. Also covers counting sort and radix sort. | Watch | 20 min |
| 4 | [QuickSort Average Case Analysis (Jeff Erickson)](https://jeffe.cs.illinois.edu/teaching/algorithms/book/Algorithms-JeffE.pdf) | Chapter 1 covers quicksort's expected O(n log n) analysis using linearity of expectation and the harmonic series. Rigorous but readable — the best written explanation of why the average case works out. | Read | 20 min |
| 5 | [Go sort package source — pdqsort](https://cs.opensource.google/go/go/+/refs/tags/go1.22.0:src/sort/zsortinterface.go) | The actual Go standard library sort implementation. Since Go 1.19, `sort.Slice` uses pdqsort (pattern-defeating quicksort). Read the code comments explaining the pivot strategy, fallback to heapsort, and insertion sort for small partitions. | Reference | 15 min |
| 6 | [Pattern-Defeating Quicksort (Orson Peters)](https://arxiv.org/abs/2106.05123) | The original pdqsort paper. Short and practical — explains how pdqsort detects adversarial patterns (already sorted, reverse sorted, organ-pipe) and adapts. This is what Go's sort actually does under the hood. | Read | 15 min |
| 7 | [Stability in Sorting — What It Is and Why It Matters (Programiz)](https://www.programiz.com/dsa/sorting-algorithm#stability) | Clear visual explanation of stability with concrete before/after examples showing how equal elements preserve their relative order. Short but fills the conceptual gap. | Read | 5 min |
| 8 | [Counting Sort — Back to Back SWE (YouTube)](https://www.youtube.com/watch?v=OKd534EWcdk) | Animated walkthrough of counting sort showing the counting array construction, prefix sum accumulation, and stable output placement. Demystifies why the backward pass preserves stability. | Watch | 10 min |

**Reading strategy:** Start with resources 1 and 2 for visual intuition — spend time on the merge step and partitioning animations. Watch resource 3 for the theoretical lower bound. Read 4 for the quicksort average-case proof. Skim 5 and 6 to understand what Go actually uses. Hit 7 and 8 before implementing counting sort.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 – 12:20 | Review Concepts (No Code)

| Time | Min | Activity |
|------|-----|----------|
| 12:00 | 8 | **Read the Day 8 section in OVERVIEW.md.** Study the complexity table. Internalize: MergeSort is always O(n log n) with O(n) extra space and stable. QuickSort is O(n log n) average, O(n^2) worst, in-place, not stable. CountingSort is O(n+k), not comparison-based. |
| 12:08 | 7 | **Open sorting.at or VisuAlgo.** Run MergeSort on [38, 27, 43, 3, 9, 82, 10]. Watch the recursive splitting and merging. Then run QuickSort on the same array — observe how pivot choice affects partition balance. Run QuickSort on an already-sorted array with first-element pivot to see the O(n^2) degradation. |
| 12:15 | 5 | **Study the ASCII diagrams in Section 6 below.** Trace the MergeSort recursion tree on paper. Then trace the QuickSort partition step. For CountingSort, follow the counting array → output construction. |

### 12:20 – 12:55 | Implement MergeSort (35 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 12:20 | 15 | **Implement MergeSort from scratch.** Create `sorting.go`. Write `MergeSort(arr []int) []int` — the recursive function that splits, recurses on each half, and calls `merge`. Write `merge(left, right []int) []int` — two-pointer merge of two sorted slices into a new result slice. Don't forget to flush the remaining tail of whichever slice isn't exhausted. |
| 12:35 | 5 | **Test MergeSort.** Write table-driven tests in `sorting_test.go`: empty, single element, already sorted, reverse sorted, all duplicates, mixed with negatives. Run `go test`. |
| 12:40 | 10 | **Trace through on paper.** Pick input `[5, 2, 8, 1, 9, 3]`. Draw the full recursion tree (splits and merges). At each merge step, write out both inputs and the merged output. Verify your code matches. |
| 12:50 | 5 | **Analyze complexity.** Convince yourself: depth is log n, each level does O(n) total merge work → O(n log n). Extra space: each merge allocates a new slice; total at any one time is O(n). MergeSort is stable because `merge` takes from `left` when values are equal. |

### 12:55 – 1:35 | Implement QuickSort (40 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 12:55 | 20 | **Implement QuickSort with Lomuto partition and randomized pivot.** Write `QuickSort(arr []int, lo, hi int)` — picks a random pivot in [lo, hi], swaps it to position hi, calls `partition`, recurses on both sides. Write `partition(arr []int, lo, hi int) int` — Lomuto scheme: pivot = arr[hi], maintain pointer `i` tracking the boundary of elements ≤ pivot, scan `j` from lo to hi-1, swap when arr[j] ≤ pivot, finally swap pivot into position i. Return i. |
| 13:15 | 5 | **Test QuickSort.** Same test cases as MergeSort. Add a test with 1000 random elements to stress-test the randomized pivot. |
| 13:20 | 10 | **Implement Hoare partition variant.** Write `partitionHoare(arr []int, lo, hi int) int` — two pointers converging from both ends. Pivot = arr[lo]. Left pointer advances while < pivot, right pointer retreats while > pivot, swap when both stop. Be careful: Hoare returns a split point where the pivot is NOT necessarily at the returned index. QuickSort recurses on [lo, p] and [p+1, hi]. |
| 13:30 | 5 | **Compare Lomuto vs Hoare.** Run both on the same input. Note: Lomuto does ~n/2 more swaps on average (every element ≤ pivot is swapped, even those already in place). Hoare converges from both ends so it does fewer swaps but has a subtler off-by-one to get right. |

### 1:35 – 1:50 | Implement CountingSort (15 minutes)

| Time | Min | Activity |
|------|-----|----------|
| 1:35 | 10 | **Implement CountingSort.** Write `CountingSort(arr []int, maxVal int) []int`. Build counting array of size maxVal+1, count occurrences, convert counts to prefix sums, build output array by scanning input right-to-left (for stability), decrement counts as you place elements. |
| 1:45 | 5 | **Test CountingSort.** Test with small range (0–9), all zeros, single element, large array with small range. Verify stability: create a struct-based test where equal keys have different payloads, confirm original order is preserved. |

### 1:50 – 2:00 | Recap

| Time | Min | Activity |
|------|-----|----------|
| 1:50 | 3 | Close all files. Write down from memory: MergeSort complexity (time, space, stable?). QuickSort (average, worst, space, stable?). CountingSort (time, space, constraint). |
| 1:53 | 3 | Write down: Why can't comparison sorts beat O(n log n)? What makes QuickSort O(n^2) in the worst case and how does randomization fix it? Why is MergeSort stable but QuickSort isn't? |
| 1:56 | 2 | Write down one gotcha you hit during each implementation. |
| 1:58 | 2 | Write down: What does Go's `sort.Slice` actually use? When would you use CountingSort instead of QuickSort? |

---

## 3. Core Concepts Deep Dive

### 3.1 Comparison Sort Lower Bound: O(n log n) and the Decision Tree Argument

Any comparison-based sorting algorithm can be modeled as a **decision tree** — a binary tree where:
- Each internal node is a comparison "is a[i] < a[j]?"
- Each leaf is a permutation of the input (one possible sorted output)
- The path from root to a leaf is the sequence of comparisons made for that input

**The key insight:** There are n! possible permutations of n elements. The decision tree must have at least n! leaves (one for each possible input ordering). A binary tree with n! leaves has height at least log₂(n!).

**Stirling's approximation:**
```
log₂(n!) = log₂(1) + log₂(2) + ... + log₂(n)
         ≥ (n/2) × log₂(n/2)
         = Θ(n log n)
```

More precisely: log₂(n!) ≈ n log₂(n) − n log₂(e) + O(log n) = Θ(n log n).

**Therefore:** Any comparison sort must make at least Ω(n log n) comparisons in the worst case. MergeSort and HeapSort achieve this bound, so they are **asymptotically optimal** among comparison sorts.

**Why counting sort escapes this bound:** It doesn't compare elements to each other. It uses element values as array indices — an operation that's not a comparison. The decision tree argument doesn't apply because counting sort isn't traversing a decision tree at all.

### 3.2 Stability: What It Means, When It Matters

A sorting algorithm is **stable** if elements with equal keys appear in the output in the same order they appeared in the input.

```
Input:  [(Alice, 85), (Bob, 90), (Carol, 85), (Dave, 90)]
Sort by score:

Stable:   [(Alice, 85), (Carol, 85), (Bob, 90), (Dave, 90)]
                ↑ Alice before Carol — same as input order

Unstable: [(Carol, 85), (Alice, 85), (Dave, 90), (Bob, 90)]
                ↑ Carol before Alice — original order NOT preserved
```

**When stability matters:**
1. **Multi-key sorting.** Sort by last name, then stable-sort by first name. The first sort's order is preserved within equal first names.
2. **Database ordering.** Users expect "sort by date, then by name" to be deterministic.
3. **Radix sort correctness.** Radix sort processes digits from least significant to most significant. Each digit-sort MUST be stable or earlier-digit ordering is destroyed.

**Which sorts are stable?**

| Algorithm | Stable? | Why / Why Not |
|-----------|---------|---------------|
| MergeSort | Yes | Merge step takes from left when equal |
| Insertion Sort | Yes | Only shifts elements that are strictly greater |
| Counting Sort | Yes | Right-to-left placement + prefix sums preserve order |
| QuickSort | No | Partition swaps can move equal elements past each other |
| HeapSort | No | Heap extraction doesn't preserve original order of equals |
| Selection Sort | No | Swaps can jump an element over equals |

**Go standard library:** `sort.Slice` is NOT stable. `sort.SliceStable` is stable (uses merge sort).

### 3.3 MergeSort: Divide-and-Conquer Analysis

**Algorithm:**
1. If array has 0 or 1 elements, it's already sorted — return.
2. Split array into two halves.
3. Recursively sort each half.
4. **Merge** the two sorted halves into a single sorted array.

**The merge step** is where all the real work happens:

```
Two sorted halves:    [1, 5, 8]  and  [2, 3, 9]

Use two pointers (i for left, j for right):
  Compare 1 vs 2 → take 1      result: [1]
  Compare 5 vs 2 → take 2      result: [1, 2]
  Compare 5 vs 3 → take 3      result: [1, 2, 3]
  Compare 5 vs 9 → take 5      result: [1, 2, 3, 5]
  Compare 8 vs 9 → take 8      result: [1, 2, 3, 5, 8]
  Right exhausted? No. Left exhausted? Yes.
  Flush right remainder → 9     result: [1, 2, 3, 5, 8, 9]
```

**Time analysis:**
- Recursion depth: log₂(n) (we halve each time)
- At each level, every element participates in exactly one merge
- Each merge is O(size of subarrays being merged)
- Total work per level = O(n)
- Total: O(n) × O(log n) = **O(n log n)** — every case (best, average, worst)

**Space analysis:**
- Each merge allocates a new slice of size left + right
- At any point in the recursion, the maximum total extra memory is O(n)
- Recursion stack depth: O(log n)
- Total: **O(n)** auxiliary space

**Why it's stable:** In the merge step, when `left[i] == right[j]`, we take from `left` first. Since `left` elements preceded `right` elements in the original array, equal elements maintain their relative order.

### 3.4 QuickSort: Lomuto vs Hoare Partitioning

**The Algorithm:**
1. Pick a **pivot** element.
2. **Partition** the array: elements < pivot go left, elements > pivot go right, pivot goes to its final position.
3. Recurse on the left and right partitions (excluding the pivot).

**Lomuto Partition (simpler, more swaps):**

```go
func partition(arr []int, lo, hi int) int {
    pivot := arr[hi]  // pivot is the last element
    i := lo           // i tracks the boundary: everything before i is ≤ pivot
    for j := lo; j < hi; j++ {
        if arr[j] <= pivot {
            arr[i], arr[j] = arr[j], arr[i]
            i++
        }
    }
    arr[i], arr[hi] = arr[hi], arr[i]  // place pivot at its final position
    return i
}
```

- Invariant: `arr[lo..i-1]` ≤ pivot, `arr[i..j-1]` > pivot, `arr[j..hi-1]` unprocessed.
- Pivot ends up at index `i` — its correct sorted position.
- Number of swaps: one for every element ≤ pivot, even those already in the correct position.

**Hoare Partition (fewer swaps, subtler):**

```go
func partitionHoare(arr []int, lo, hi int) int {
    pivot := arr[lo]  // pivot is the first element
    i, j := lo-1, hi+1
    for {
        for { i++; if arr[i] >= pivot { break } }
        for { j--; if arr[j] <= pivot { break } }
        if i >= j {
            return j
        }
        arr[i], arr[j] = arr[j], arr[i]
    }
}
```

- Two pointers converge from opposite ends.
- Swaps only happen when both pointers find misplaced elements — roughly half the swaps of Lomuto.
- **Critical difference:** Hoare returns a split point `j`. The pivot is NOT necessarily at `arr[j]`. You recurse on `[lo, j]` and `[j+1, hi]` — NOT `[lo, j-1]` and `[j+1, hi]`.

**Comparison:**

| Aspect | Lomuto | Hoare |
|--------|--------|-------|
| Pivot placement | Pivot is at its final sorted position | Pivot can be anywhere in [lo, j] |
| Recursion | [lo, p-1] and [p+1, hi] | [lo, p] and [p+1, hi] |
| Swaps (average) | ~n/2 | ~n/6 |
| All-duplicates input | O(n^2) — pivot always at one end | O(n log n) — pointers meet in middle |
| Code complexity | Simpler, easier to get right | Off-by-one errors are common |
| Interview recommendation | **Use Lomuto** — easier to code correctly | Know it exists, implement if asked |

**Randomized Pivot:**

```go
func QuickSort(arr []int, lo, hi int) {
    if lo >= hi { return }
    // Randomize to avoid O(n^2) on sorted/reverse-sorted input
    randIdx := lo + rand.Intn(hi-lo+1)
    arr[randIdx], arr[hi] = arr[hi], arr[randIdx]
    p := partition(arr, lo, hi)
    QuickSort(arr, lo, p-1)
    QuickSort(arr, p+1, hi)
}
```

**Average case O(n log n) — intuition:** A random pivot is "good" (splits the array between 25%-75%) about half the time. Even with 50% bad pivots, the expected recursion depth is O(log n), and each level does O(n) work. The formal proof uses linearity of expectation on the number of comparisons: E[comparisons] = 2n × H(n) ≈ 2n ln n = O(n log n), where H(n) is the nth harmonic number.

**Worst case O(n^2):** If the pivot is always the smallest or largest element (e.g., sorted input with first-element pivot), one partition has n-1 elements and the other has 0. Depth becomes n, total work is n + (n-1) + ... + 1 = O(n^2). Randomized pivot makes this astronomically unlikely — the probability of O(n^2) with random pivots is ~1/n!.

### 3.5 CountingSort: When and How

**Constraint:** Only works on non-negative integers (or values mappable to non-negative integers) within a known, small range [0, k].

**Algorithm:**
1. Create a counting array `count[0..k]` initialized to zero.
2. Count occurrences: for each element `x`, increment `count[x]`.
3. Convert to prefix sums: `count[i] += count[i-1]`. Now `count[x]` tells you how many elements are ≤ x.
4. Build output: scan input **right-to-left**, place each element at `output[count[x] - 1]`, decrement `count[x]`.

```
Input:  [4, 2, 2, 8, 3, 3, 1]     k = 8

Step 1 — Count:
  count: [0, 1, 2, 2, 1, 0, 0, 0, 1]
          0  1  2  3  4  5  6  7  8

Step 2 — Prefix sums:
  count: [0, 1, 3, 5, 6, 6, 6, 6, 7]

Step 3 — Place (right to left):
  arr[6] = 1 → output[count[1]-1] = output[0] = 1, count[1]-- → 0
  arr[5] = 3 → output[count[3]-1] = output[4] = 3, count[3]-- → 4
  arr[4] = 3 → output[count[3]-1] = output[3] = 3, count[3]-- → 3
  arr[3] = 8 → output[count[8]-1] = output[6] = 8, count[8]-- → 6
  arr[2] = 2 → output[count[2]-1] = output[2] = 2, count[2]-- → 2
  arr[1] = 2 → output[count[2]-1] = output[1] = 2, count[2]-- → 1
  arr[0] = 4 → output[count[4]-1] = output[5] = 4, count[4]-- → 5

Output: [1, 2, 2, 3, 3, 4, 8]
```

**Time:** O(n + k) — two passes over the input (count + place), one pass over the count array (prefix sums).

**Space:** O(n + k) — output array O(n), counting array O(k).

**Why right-to-left for stability:** Scanning right-to-left ensures that when two equal elements exist, the one appearing later in the input is placed at the higher index in the output. This preserves the original relative order of equal elements.

**When to use CountingSort:**
- Integer values in a small, known range (e.g., ages 0–150, ASCII characters 0–127, scores 0–100)
- When k = O(n), total time is O(n) — beating the comparison sort lower bound
- **Avoid when:** k >> n (the counting array wastes space), or values are floats/strings/complex objects

### 3.6 Go's sort.Slice Internals: pdqsort

Since Go 1.19, `sort.Slice` uses **pdqsort (pattern-defeating quicksort)** — a hybrid algorithm by Orson Peters.

**What pdqsort does:**

1. **Small subarrays (≤ 12 elements):** Uses **insertion sort**. The overhead of recursion and pivot selection isn't worth it for tiny arrays. Insertion sort's low constant factor and good cache behavior win here.

2. **Normal case:** Uses **quicksort with median-of-three pivot selection**. Pick the median of the first, middle, and last elements as the pivot. This avoids the worst case on already-sorted and reverse-sorted inputs.

3. **Bad pivot detection:** After each partition, check if the pivot was "bad" (one side has < 1/8 of the elements). If too many bad pivots accumulate, switch to **heapsort** — guaranteeing O(n log n) worst case.

4. **Pattern detection:** If the input is already sorted or nearly sorted, pdqsort detects this (partition produces an empty left side) and switches to **insertion sort** for that partition. This gives O(n) on already-sorted input.

5. **Equal elements:** If many equal elements cluster around the pivot, pdqsort uses a special "partition equal" step that groups equals together, avoiding unnecessary recursion over them.

**Why Go chose pdqsort:**
- O(n log n) guaranteed (heapsort fallback)
- O(n) on already-sorted input (pattern detection)
- Faster than pure quicksort on real-world data (adapts to patterns)
- In-place (O(log n) stack space)
- Significantly outperforms introsort (the previous common hybrid) in benchmarks

**Key takeaway:** Go's `sort.Slice` is NOT stable. Use `sort.SliceStable` when you need stability — it uses a bottom-up merge sort.

### 3.7 Practical Considerations: Sort vs. Heap vs. Bucket

Not every "find order" problem needs a full sort.

| Scenario | Best Approach | Why |
|----------|--------------|-----|
| Need fully sorted output | Sort — O(n log n) | That's literally what sort does |
| Need top-k elements | Heap — O(n log k) | Min-heap of size k; never sort the rest |
| Need the kth element only | Quickselect — O(n) avg | Partial quicksort; no need to sort either side fully |
| Need sorted output, small integer range | Counting sort — O(n+k) | Beats comparison sorts |
| Need sorted output, bounded-length keys | Radix sort — O(d(n+k)) | d passes of counting sort |
| Streaming data, need current top-k | Heap — O(log k) per element | Can't sort what you haven't seen yet |
| Data is nearly sorted (few inversions) | Insertion sort — O(n + inversions) | Adaptive; degrades gracefully |

**Rule of thumb:** If you only need partial ordering information (top-k, kth element, min/max), a heap or selection algorithm is better than sorting. If you need the complete ordering, sort.

---

## 4. Implementation Checklist

### Function Signatures

```go
package sorting

import "math/rand"

// MergeSort returns a new sorted slice. Does not modify the input.
// Time: O(n log n)  Space: O(n)  Stable: Yes
func MergeSort(arr []int) []int

// merge combines two sorted slices into one sorted slice.
func merge(left, right []int) []int

// QuickSort sorts arr[lo..hi] in-place using Lomuto partition
// with randomized pivot.
// Time: O(n log n) avg, O(n^2) worst   Space: O(log n)  Stable: No
func QuickSort(arr []int, lo, hi int)

// partition rearranges arr[lo..hi] around a pivot (arr[hi]).
// Returns the pivot's final index.
func partition(arr []int, lo, hi int) int

// CountingSort sorts non-negative integers with values in [0, maxVal].
// Returns a new sorted slice.
// Time: O(n + k)  Space: O(n + k)  Stable: Yes
func CountingSort(arr []int, maxVal int) []int
```

### Test Cases

```go
package sorting

import (
    "math/rand"
    "sort"
    "testing"
)

func TestMergeSort(t *testing.T) {
    tests := []struct {
        name  string
        input []int
        want  []int
    }{
        {"empty", []int{}, []int{}},
        {"single", []int{42}, []int{42}},
        {"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
        {"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
        {"all duplicates", []int{7, 7, 7, 7}, []int{7, 7, 7, 7}},
        {"mixed", []int{38, 27, 43, 3, 9, 82, 10}, []int{3, 9, 10, 27, 38, 43, 82}},
        {"negatives", []int{-3, 0, -1, 5, 2}, []int{-3, -1, 0, 2, 5}},
        {"two elements", []int{2, 1}, []int{1, 2}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MergeSort(tt.input)
            // compare got with tt.want
        })
    }
}

func TestQuickSort(t *testing.T) {
    // Same test cases as MergeSort, but QuickSort is in-place:
    // Copy input, call QuickSort(copy, 0, len(copy)-1), compare.
    // Additional: 1000 random elements — verify against sort.Ints.
}

func TestCountingSort(t *testing.T) {
    // Test with maxVal=9, input [4,2,2,8,3,3,1] → [1,2,2,3,3,4,8]
    // All zeros: [0,0,0] → [0,0,0]
    // Single element: [5] → [5]
    // Large array, small range: 10000 elements in [0, 10]
}

// TestQuickSortStress verifies correctness on random large inputs.
func TestQuickSortStress(t *testing.T) {
    for trial := 0; trial < 100; trial++ {
        n := rand.Intn(1000) + 1
        arr := make([]int, n)
        for i := range arr {
            arr[i] = rand.Intn(2000) - 1000
        }
        expected := make([]int, n)
        copy(expected, arr)
        sort.Ints(expected)

        QuickSort(arr, 0, len(arr)-1)
        // compare arr with expected
    }
}
```

**Edge cases to cover explicitly:**
- **Empty slice** — should return empty, not panic
- **Single element** — trivially sorted
- **Already sorted** — verifies QuickSort with randomized pivot doesn't degrade
- **Reverse sorted** — classic QuickSort worst case without randomization
- **All duplicates** — Lomuto partition degrades to O(n^2) here; this is a known weakness
- **Two elements** — boundary condition for partition logic

---

## 5. Sorting as a Preprocessing Step

Sorting is rarely the end goal. It's a **preprocessing step** that enables faster algorithms. Recognizing when sorting unlocks a simpler solution is a key pattern.

### 5.1 Two Pointers (Needs Sorted Input)

**Pattern:** Sort the array, then use converging pointers from both ends.

**Example — Two Sum (sorted variant):**
```go
sort.Ints(arr)
lo, hi := 0, len(arr)-1
for lo < hi {
    sum := arr[lo] + arr[hi]
    if sum == target { return true }
    if sum < target { lo++ } else { hi-- }
}
```

**Why sorting helps:** Without sorting, you need a hash map (O(n) space) or brute-force O(n^2). With sorting, two pointers gives O(n) time, O(1) extra space (if sorting in-place).

**Also applies to:** 3Sum, 4Sum, container with most water, trapping rain water.

### 5.2 Binary Search (Needs Sorted Input)

**Pattern:** Sort once O(n log n), then answer multiple queries in O(log n) each.

**Example — Count elements in a range [lo, hi]:**
```go
sort.Ints(arr)
left := sort.SearchInts(arr, lo)       // lower bound
right := sort.SearchInts(arr, hi+1)    // upper bound
count := right - left
```

**When this wins:** If you have q queries on the same array, sorting once + binary search per query is O(n log n + q log n), beating O(q × n) brute force.

### 5.3 Grouping Anagrams (Sort Each String as a Key)

**Pattern:** Sort the characters in each string to create a canonical key. Group by key.

```go
func GroupAnagrams(strs []string) [][]string {
    groups := make(map[string][]string)
    for _, s := range strs {
        key := sortString(s)  // "eat" → "aet", "tea" → "aet"
        groups[key] = append(groups[key], s)
    }
    // collect groups into result
}

func sortString(s string) string {
    b := []byte(s)
    sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
    return string(b)
}
```

**Why sorting:** Two strings are anagrams iff they have the same characters in the same quantities. Sorting both produces the same string → identical hash map key. Alternative: use a character-frequency array as the key (avoids sorting, but the sorted-string approach is simpler to code).

### 5.4 Interval Problems (Sort by Start or End Time)

**Pattern:** Sort intervals by start time (or end time), then scan linearly.

**Merge Intervals:**
```go
sort.Slice(intervals, func(i, j int) bool {
    return intervals[i][0] < intervals[j][0]
})
// Now overlapping intervals are adjacent — scan and merge
```

**Non-overlapping Intervals (max intervals to keep):**
```go
sort.Slice(intervals, func(i, j int) bool {
    return intervals[i][1] < intervals[j][1]  // sort by END time
})
// Greedy: pick earliest-ending interval, skip all overlapping ones
```

**Why sorting is essential:** Without sorting, you'd need O(n^2) pairwise comparison to find overlaps. After sorting by start time, overlapping intervals are guaranteed to be adjacent in the sorted order.

### 5.5 Meeting the Greedy Choice Property (Sort by Deadline, Profit, etc.)

**Pattern:** Many greedy algorithms require sorting to establish the order in which to make choices.

| Problem | Sort By | Greedy Choice |
|---------|---------|---------------|
| Activity selection | End time (ascending) | Pick earliest-ending compatible activity |
| Job scheduling with deadlines | Deadline (ascending) or profit (descending) | Schedule highest-profit jobs by deadline |
| Fractional knapsack | Value/weight ratio (descending) | Take items greedily by density |
| Minimum platforms | Start time + end time events | Sweep line: +1 at start, -1 at end |
| Assign cookies to children | Both arrays sorted ascending | Match smallest cookie to least-greedy child |

**The general principle:** Sorting reveals the optimal order for greedy choices. Without sorting, you can't guarantee the locally optimal choice is globally optimal.

---

## 6. Visual Diagrams

### 6.1 MergeSort Recursion Tree

Input: `[38, 27, 43, 3, 9, 82, 10]`

```
                        [38, 27, 43, 3, 9, 82, 10]
                       /                            \
              [38, 27, 43, 3]                   [9, 82, 10]
              /              \                  /           \
         [38, 27]         [43, 3]          [9, 82]         [10]
         /      \         /     \          /      \          |
       [38]    [27]    [43]    [3]       [9]    [82]       [10]
         \      /         \     /          \      /          |
         [27, 38]         [3, 43]          [9, 82]         [10]
              \              /                  \           /
          [3, 27, 38, 43]                    [9, 10, 82]
                       \                            /
                   [3, 9, 10, 27, 38, 43, 82]


Merge step detail (merging [3,27,38,43] and [9,10,82]):

  left:  [3, 27, 38, 43]     right: [9, 10, 82]
          ^                           ^
          i                           j

  Compare 3 < 9  → take 3    output: [3]
  Compare 27 > 9 → take 9    output: [3, 9]
  Compare 27 > 10→ take 10   output: [3, 9, 10]
  Compare 27 < 82→ take 27   output: [3, 9, 10, 27]
  Compare 38 < 82→ take 38   output: [3, 9, 10, 27, 38]
  Compare 43 < 82→ take 43   output: [3, 9, 10, 27, 38, 43]
  Left exhausted → flush 82  output: [3, 9, 10, 27, 38, 43, 82]

  Total comparisons at this merge: 6 (for 7 elements)
```

### 6.2 QuickSort Partition Step (Lomuto)

Input: `[3, 7, 8, 5, 2, 1, 9, 5, 4]`, pivot = last element = 4

```
  Initial state:  pivot = arr[8] = 4
  i = 0 (boundary of "≤ pivot" region)
  j scans left to right

  j=0: arr[0]=3 ≤ 4 → swap arr[0]↔arr[0], i++
       [3, 7, 8, 5, 2, 1, 9, 5, 4]
        ≤ |  ?  ?  ?  ?  ?  ?  ?  P
        i=1

  j=1: arr[1]=7 > 4 → skip
       [3, 7, 8, 5, 2, 1, 9, 5, 4]
        ≤ | >  ?  ?  ?  ?  ?  ?  P
        i=1

  j=2: arr[2]=8 > 4 → skip
  j=3: arr[3]=5 > 4 → skip
       [3, 7, 8, 5, 2, 1, 9, 5, 4]
        ≤ | >  >  >  ?  ?  ?  ?  P
        i=1

  j=4: arr[4]=2 ≤ 4 → swap arr[1]↔arr[4], i++
       [3, 2, 8, 5, 7, 1, 9, 5, 4]
        ≤  ≤ | >  >  >  ?  ?  ?  P
           i=2

  j=5: arr[5]=1 ≤ 4 → swap arr[2]↔arr[5], i++
       [3, 2, 1, 5, 7, 8, 9, 5, 4]
        ≤  ≤  ≤ | >  >  >  ?  ?  P
              i=3

  j=6: arr[6]=9 > 4 → skip
  j=7: arr[7]=5 > 4 → skip

  Final: swap arr[i]↔arr[hi] (pivot into position)
       [3, 2, 1, 4, 7, 8, 9, 5, 5]
        ≤  ≤  ≤  P| >  >  >  >  >
              pivot at index 3

  Result: pivot 4 is at its FINAL sorted position (index 3).
          Everything left of it is ≤ 4.
          Everything right of it is > 4.
          Recurse on [3,2,1] and [7,8,9,5,5].
```

### 6.3 CountingSort: Counting Array and Output Construction

Input: `[4, 2, 2, 8, 3, 3, 1]`, maxVal = 8

```
Step 1 — Count occurrences:

  Input:  4  2  2  8  3  3  1
          ↓  ↓  ↓  ↓  ↓  ↓  ↓

  Index:  0  1  2  3  4  5  6  7  8
  Count: [0, 1, 2, 2, 1, 0, 0, 0, 1]
              ↑  ↑  ↑  ↑           ↑
             one two two one      one
              1   2   3   4        8


Step 2 — Prefix sums (cumulative count):

  Index:  0  1  2  3  4  5  6  7  8
  Count: [0, 1, 3, 5, 6, 6, 6, 6, 7]
              ↑  ↑  ↑  ↑           ↑
         "1 element  "6 elements   "7 elements
          is ≤ 1"    are ≤ 4"      are ≤ 8"


Step 3 — Build output (scan input RIGHT to LEFT for stability):

  Output: [ _, _, _, _, _, _, _ ]
           0  1  2  3  4  5  6

  i=6: arr[6]=1, count[1]=1 → output[0]=1, count[1]=0
       [ 1, _, _, _, _, _, _ ]

  i=5: arr[5]=3, count[3]=5 → output[4]=3, count[3]=4
       [ 1, _, _, _, 3, _, _ ]

  i=4: arr[4]=3, count[3]=4 → output[3]=3, count[3]=3
       [ 1, _, _, 3, 3, _, _ ]

  i=3: arr[3]=8, count[8]=7 → output[6]=8, count[8]=6
       [ 1, _, _, 3, 3, _, 8 ]

  i=2: arr[2]=2, count[2]=3 → output[2]=2, count[2]=2
       [ 1, _, 2, 3, 3, _, 8 ]

  i=1: arr[1]=2, count[2]=2 → output[1]=2, count[2]=1
       [ 1, 2, 2, 3, 3, _, 8 ]

  i=0: arr[0]=4, count[4]=6 → output[5]=4, count[4]=5
       [ 1, 2, 2, 3, 3, 4, 8 ]  ✓ sorted!
```

---

## 7. Self-Assessment

Answer these from memory after your session. If you can't, that's tomorrow's priority.

### Q1: Why can't any comparison-based sort do better than O(n log n)?

<details>
<summary>Answer</summary>

Any comparison sort can be modeled as a binary decision tree where each internal node is a comparison and each leaf is a permutation. There are n! possible input permutations, so the tree needs at least n! leaves. A binary tree with n! leaves has height at least log₂(n!) = Θ(n log n). The height is the worst-case number of comparisons — so every comparison sort must make at least Ω(n log n) comparisons in the worst case.

Non-comparison sorts (counting sort, radix sort) bypass this bound because they use element values as indices rather than comparing elements to each other.
</details>

### Q2: When would you choose MergeSort over QuickSort?

<details>
<summary>Answer</summary>

Choose MergeSort when:
1. **Stability is required** — MergeSort preserves the relative order of equal elements; QuickSort does not.
2. **Worst-case guarantee matters** — MergeSort is always O(n log n); QuickSort is O(n^2) worst case (though randomized pivot makes this unlikely).
3. **Sorting linked lists** — MergeSort needs no random access (only sequential traversal for merge), making it ideal for linked lists. QuickSort's partition step requires random access.
4. **External sorting** — When data doesn't fit in memory, merge sort's sequential access pattern works well with disk I/O.

Choose QuickSort when:
1. **Memory is constrained** — QuickSort is in-place (O(log n) stack space); MergeSort needs O(n) extra space.
2. **Average-case performance matters most** — QuickSort has better cache locality and smaller constant factors, making it faster in practice despite the same asymptotic complexity.
</details>

### Q3: Explain the difference between Lomuto and Hoare partitioning. Which is better for all-duplicate arrays?

<details>
<summary>Answer</summary>

**Lomuto:** Uses one pointer scanning left-to-right plus a boundary pointer. Pivot goes to `arr[hi]`. After partition, pivot is at its final sorted position. Recurse on `[lo, p-1]` and `[p+1, hi]`. Simpler to implement.

**Hoare:** Uses two converging pointers from opposite ends. Pivot is `arr[lo]`. Pointers swap misplaced elements until they cross. After partition, pivot is NOT necessarily at the returned index. Recurse on `[lo, p]` and `[p+1, hi]`.

**All-duplicate arrays:** Hoare is much better. Lomuto's partition puts every element (all equal to pivot) in the left partition, producing an n-1/0 split — O(n^2). Hoare's converging pointers meet roughly in the middle, producing a balanced n/2/n/2 split — O(n log n).

This is a significant practical difference and one reason production implementations (including pdqsort) use Hoare-style partitioning or special-case equal elements.
</details>

### Q4: You have 10 million integers, each in the range [0, 1000]. What sorting algorithm do you use and why?

<details>
<summary>Answer</summary>

**Counting Sort.** The value range k = 1000 is tiny compared to n = 10,000,000.

Time: O(n + k) = O(10,000,000 + 1,000) ≈ O(n) — linear.
Space: O(n + k) — the output array dominates.

This is dramatically faster than any comparison sort, which would require O(n log n) ≈ 10,000,000 × 23 ≈ 230 million operations. Counting sort does roughly 20 million operations (two passes over n, one over k).

The small range [0, 1000] makes the counting array trivially sized (1001 integers = ~8 KB). This is the textbook case for counting sort.
</details>

### Q5: What is pdqsort and why does Go use it instead of plain quicksort or introsort?

<details>
<summary>Answer</summary>

**pdqsort (pattern-defeating quicksort)** is a hybrid sorting algorithm that combines:
- **Quicksort** with median-of-three pivot for the normal case
- **Insertion sort** for small subarrays (≤ 12 elements)
- **Heapsort** fallback when too many bad pivots are detected (guaranteeing O(n log n) worst case)
- **Pattern detection** that recognizes already-sorted runs and switches to insertion sort, achieving O(n) on sorted input

Go chose pdqsort over plain quicksort because plain quicksort has O(n^2) worst case. It chose pdqsort over introsort because pdqsort adapts better to real-world patterns — it's faster on nearly-sorted data, data with many duplicates, and other structured inputs that are common in practice. The heapsort fallback ensures the same O(n log n) worst-case guarantee as introsort.
</details>

---

## Complexity Reference (Quick Glance)

| Algorithm | Best | Average | Worst | Space | Stable | Notes |
|-----------|------|---------|-------|-------|--------|-------|
| MergeSort | O(n log n) | O(n log n) | O(n log n) | O(n) | Yes | Predictable, good for linked lists |
| QuickSort | O(n log n) | O(n log n) | O(n^2) | O(log n) | No | Fastest in practice; randomize pivot |
| CountingSort | O(n + k) | O(n + k) | O(n + k) | O(n + k) | Yes | k = value range; non-comparison |
| HeapSort | O(n log n) | O(n log n) | O(n log n) | O(1) | No | In-place, but poor cache locality |
| Insertion Sort | O(n) | O(n^2) | O(n^2) | O(1) | Yes | Best for small/nearly-sorted arrays |
| Go sort.Slice | O(n) | O(n log n) | O(n log n) | O(log n) | No | pdqsort: adaptive hybrid |
| Go sort.SliceStable | O(n log n) | O(n log n) | O(n log n) | O(n) | Yes | Bottom-up merge sort |
