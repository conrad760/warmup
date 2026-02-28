# Day 10 — Two Pointers: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Focus | Time |
|---|----------|-------|------|
| 1 | [Two Pointers Technique — NeetCode](https://neetcode.io/courses/advanced-algorithms/3) | Video walkthrough of all three patterns with animated pointer movement. Good for visual learners. | 15 min |
| 2 | [Two Sum II — LeetCode Editorial](https://leetcode.com/problems/two-sum-ii-input-array-is-sorted/editorial/) | Step-by-step proof of correctness for opposite-ends on sorted arrays, with diagrams showing why you never skip valid pairs. | 10 min |
| 3 | [Container With Most Water — NeetCode](https://www.youtube.com/watch?v=UuiTKBwPgAo) | Visual explanation of why moving the shorter side is always safe. Includes the greedy proof. | 10 min |
| 4 | [Floyd's Cycle Detection — Brilliant.org](https://brilliant.org/wiki/floyds-cycle-detection-algorithm/) | Mathematical proof of why slow and fast pointers meet inside the cycle, and why resetting one to head finds the cycle start. | 15 min |
| 5 | [Dutch National Flag Problem — Wikipedia](https://en.wikipedia.org/wiki/Dutch_national_flag_problem) | Dijkstra's original 3-way partition algorithm. Clear loop invariant definition. | 10 min |
| 6 | [LeetCode Two Pointers Problem Set](https://leetcode.com/tag/two-pointers/) | Sorted by acceptance rate. Start with Easy tier: Remove Duplicates, Move Zeroes, Valid Palindrome. | Reference |
| 7 | [Go Slices: Usage and Internals — Go Blog](https://go.dev/blog/slices-intro) | Understanding slice headers, in-place modification, and why `nums[:slow]` works. Critical for write-pointer problems. | 10 min |
| 8 | [Trapping Rain Water — Back to Back SWE](https://www.youtube.com/watch?v=HmBbcDiJapY) | Visual walkthrough of the two-pointer approach to trapping rain water, comparing it to the prefix-max approach. | 12 min |

---

## 2. Detailed 2-Hour Session Plan

### 12:00 — 12:20 | Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:07 | Read the three patterns below (Section 3). For each one, mentally trace through a small example. No keyboard. |
| 12:07 - 12:12 | Study the proof sketch for two-sum sorted (Section 3.2). Convince yourself why no valid pair is skipped. Draw the pointer positions on paper. |
| 12:12 - 12:17 | Read the "write pointer" pattern and Dutch National Flag pattern. Trace remove-duplicates on `[1,1,2,2,3]` by hand. |
| 12:17 - 12:20 | Review the decision framework (Section 5). For each of these, name the pattern: "sorted array pair sum," "remove zeros in-place," "linked list cycle." |

### 12:20 — 1:20 | Implement (From Scratch)

| Time | Problem | Pattern | Notes |
|------|---------|---------|-------|
| 12:20 - 12:32 | `TwoSumSorted` | Opposite ends | Start here — it is the canonical example. Return 1-indexed. Write tests for: no solution, single pair, negative numbers. |
| 12:32 - 12:47 | `ThreeSum` | Opposite ends + iteration | Fix `nums[i]`, run two-sum on the remainder. The hard part is skipping duplicates at all three levels. |
| 12:47 - 13:00 | `ContainerWithMostWater` | Opposite ends | After implementing, write a comment explaining *why* you move the shorter side (the proof). |
| 13:00 - 13:08 | `RemoveDuplicates` | Same direction (write pointer) | `slow` is the write position, `fast` scans. Return new length. |
| 13:08 - 13:15 | `MoveZeroes` | Same direction (write pointer) | Write all non-zeros forward, then fill zeros. Or swap-based approach. |
| 13:15 - 13:20 | `LinkedListCycle` | Floyd's (fast-slow) | Check `fast != nil && fast.Next != nil` before advancing. |

### 1:20 — 1:50 | Solidify (Edge Cases & Variants)

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | `FindCycleStart` — After slow and fast meet, reset one to head, advance both by 1. They meet at cycle start. Add tests: no cycle, cycle at head, cycle in middle. |
| 1:30 - 1:42 | `TrappingRainWater` (two-pointer approach) — Track `leftMax` and `rightMax`. Process from whichever side has the smaller max. This is the hardest problem today. |
| 1:42 - 1:50 | Go back to `ThreeSum` and test: all zeros `[0,0,0,0]`, no valid triplets, all negatives, large duplicate runs. Verify your dedup logic handles every case. |

### 1:50 — 2:00 | Recap (From Memory)

Write down without looking:
1. The three two-pointer patterns and when to use each.
2. The complexity of every function you implemented today (all O(n) time, O(1) space — except ThreeSum which is O(n^2) time).
3. One gotcha per problem (e.g., "ThreeSum: skip duplicates *after* finding a match, not before").
4. The invariant that makes two-sum sorted correct in one sentence.

---

## 3. Core Concepts Deep Dive

### 3.1 The Three Patterns

#### Pattern A: Opposite Ends (Converging)

```
left = 0, right = len-1
while left < right:
    evaluate(left, right)
    if need_more:  left++
    else:          right--
```

**When to use:** Sorted input where you need pairs satisfying a condition (sum, difference, area). The sorted order provides a monotonic property that tells you which pointer to move.

**Key invariant:** At every step, the optimal answer either involves the current pair or lies strictly between the current pointers. You never need to revisit a pointer position you have passed.

#### Pattern B: Same Direction (Fast-Slow for Arrays)

```
slow = 0   // write position
for fast := 0; fast < len; fast++:
    if shouldKeep(arr[fast]):
        arr[slow] = arr[fast]
        slow++
return slow  // new length
```

**When to use:** In-place compaction, partitioning, or filtering. The slow pointer marks "where to write next," the fast pointer scans every element.

**Key invariant:** Everything before `slow` is the processed, valid output. Everything from `slow` to `fast` is "garbage" that has been consumed or will be overwritten.

#### Pattern C: Floyd's Cycle Detection (Fast-Slow for Linked Lists)

```
slow = head, fast = head
while fast != nil && fast.Next != nil:
    slow = slow.Next
    fast = fast.Next.Next
    if slow == fast:
        // cycle detected
```

**When to use:** Linked list cycle detection, finding the middle of a linked list, finding the start of a cycle, detecting duplicates in a constrained integer array (LeetCode 287).

**Key invariant:** If a cycle of length C exists, after entering the cycle, the distance between slow and fast decreases by 1 each step. They must meet within C steps of both being in the cycle.

---

### 3.2 WHY Two Pointers Gives O(n)

The brute force for pair problems is O(n^2): try all pairs. Two pointers achieves O(n) because **the sorted order creates an invariant that lets you safely eliminate entire groups of candidates with a single pointer move.**

Each step, at least one pointer moves inward. There are at most `n` inward moves total (left can move right at most n times, right can move left at most n times). So the total work is O(n).

The critical insight is not just efficiency — it is **correctness**. You must prove that moving a pointer never skips a valid answer.

#### Proof Sketch: Two-Sum on Sorted Array

**Given:** Sorted array `a[0] <= a[1] <= ... <= a[n-1]`, target `T`.
**Claim:** If `a[L] + a[R] == T` for some `L < R`, the two-pointer algorithm will find it.

**Proof by contradiction:**

Start with `left = 0`, `right = n-1`.

Consider the moment when `left = L` (we have advanced left to the correct position). We need to show that `right >= R` at this point — i.e., we have not moved right past R.

- `right` only decreases when `a[left] + a[right] > T`.
- For all `left' < L`: `a[left'] < a[L]`, so `a[left'] + a[R] < a[L] + a[R] = T`.
- When `left < L`, we have `a[left] + a[right] < T` for any `right <= R` (since `a[left] < a[L]`). So the algorithm moves `left++`, not `right--`.
- Therefore `right` cannot have decreased below `R` before `left` reaches `L`.

Symmetrically, if we reach `right = R` first, `left` cannot have passed `L`.

Thus the algorithm encounters the state `(L, R)` and returns the answer. QED.

**The intuition in one sentence:** Moving `left` right can only increase the sum; moving `right` left can only decrease it. If the sum is too small, the left element is too small for *any* remaining right element, so skip it.

---

### 3.3 The "Write Pointer" Pattern

This is a specific instance of same-direction two pointers used for **in-place array modification**.

```
slow (write position) ──►
                          [processed | unprocessed]
fast (read position)  ───────────────────────────►
```

**Canonical problems:**

| Problem | Keep Condition | Result |
|---------|---------------|--------|
| Remove Duplicates (sorted) | `nums[fast] != nums[slow-1]` | Unique elements, return new length |
| Move Zeroes | `nums[fast] != 0` | Non-zeros moved front, zeros at back |
| Remove Element | `nums[fast] != val` | All instances of `val` removed |

**Why this works:** We only need one pass. The fast pointer reads every element exactly once. The slow pointer writes at most once per element. Total work: O(n). Space: O(1) — we modify in-place.

**Go implementation pattern:**

```go
func RemoveDuplicates(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
    slow := 1 // first element is always kept
    for fast := 1; fast < len(nums); fast++ {
        if nums[fast] != nums[slow-1] {
            nums[slow] = nums[fast]
            slow++
        }
    }
    return slow
}
```

---

### 3.4 Dutch National Flag Problem (3-Way Partition)

Three pointers partition an array into three regions. Dijkstra designed this for quicksort with many duplicate keys.

**Problem:** Given an array with values 0, 1, 2, sort them in one pass.

**Invariant:**

```
[0..lo-1]  = all 0s     (red)
[lo..mid-1] = all 1s    (white)
[mid..hi]   = unknown
[hi+1..n-1] = all 2s    (blue)

 0  0  0  1  1  1  ?  ?  ?  2  2  2
         lo       mid       hi
```

**Algorithm:**

```go
func SortColors(nums []int) {
    lo, mid, hi := 0, 0, len(nums)-1
    for mid <= hi {
        switch nums[mid] {
        case 0:
            nums[lo], nums[mid] = nums[mid], nums[lo]
            lo++
            mid++
        case 1:
            mid++
        case 2:
            nums[mid], nums[hi] = nums[hi], nums[mid]
            hi--
            // do NOT advance mid — the swapped element is unknown
        }
    }
}
```

**Why `mid` does not advance on case 2:** We swapped an unknown value from position `hi` into position `mid`. We must examine it before moving on.

**Application:** Any 3-way partition. Used in quicksort when the pivot has many duplicates (partition into `< pivot`, `== pivot`, `> pivot`).

---

## 4. Implementation Checklist

### Function Signatures

```go
package twopointers

// -- Opposite Ends --

// TwoSumSorted returns 1-indexed positions of two numbers that add to target.
// Precondition: nums is sorted in non-decreasing order, exactly one solution exists.
func TwoSumSorted(nums []int, target int) (int, int) { ... }

// ThreeSum returns all unique triplets [a, b, c] where a + b + c = 0.
// Results must not contain duplicate triplets.
func ThreeSum(nums []int) [][]int { ... }

// ContainerWithMostWater returns the max area of water between two lines.
func ContainerWithMostWater(height []int) int { ... }

// TrappingRainWater returns total water trapped between bars.
func TrappingRainWater(height []int) int { ... }

// -- Same Direction (Write Pointer) --

// RemoveDuplicates removes duplicates in-place from sorted array.
// Returns new length. Elements beyond the new length are irrelevant.
func RemoveDuplicates(nums []int) int { ... }

// MoveZeroes moves all zeros to the end while maintaining relative order
// of non-zero elements. Modifies in-place.
func MoveZeroes(nums []int) { ... }

// -- Floyd's Cycle Detection --

type ListNode struct {
    Val  int
    Next *ListNode
}

// HasCycle returns true if the linked list contains a cycle.
func HasCycle(head *ListNode) bool { ... }

// FindCycleStart returns the node where the cycle begins, or nil if no cycle.
func FindCycleStart(head *ListNode) *ListNode { ... }
```

### Test Cases & Edge Cases

| Function | Must-Test Cases |
|----------|----------------|
| `TwoSumSorted` | Two-element array `[1,3], target=4`; negative numbers `[-3,-1,0,2,4], target=1`; target is double of one element `[1,2,3,4], target=6` |
| `ThreeSum` | All zeros `[0,0,0,0]` → `[[0,0,0]]`; no triplets `[1,2,3]` → `[]`; large duplicate runs `[-1,-1,-1,0,1,1,1]`; empty/short arrays |
| `ContainerWithMostWater` | Two elements; descending heights `[5,4,3,2,1]`; equal heights; single tall bar in middle |
| `TrappingRainWater` | Ascending only (0 water); descending only (0 water); V-shape `[3,0,3]` → 3; single/two elements → 0 |
| `RemoveDuplicates` | Empty array; all same `[1,1,1]` → 1; already unique `[1,2,3]` → 3; single element |
| `MoveZeroes` | No zeros; all zeros; zeros at front; zeros at back; alternating `[0,1,0,1]` |
| `HasCycle` | Nil head; single node no cycle; single node self-cycle; cycle in middle; tail points to head |
| `FindCycleStart` | No cycle → nil; cycle at head; cycle at middle node; long tail before cycle |

---

## 5. Pattern Decision Framework

```
                        ┌─────────────────────────┐
                        │ Can I use two pointers?  │
                        └────────────┬────────────┘
                                     │
                    ┌────────────────┼────────────────┐
                    ▼                ▼                 ▼
          ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
          │  Input is     │  │  Need in-place│  │ Linked list  │
          │  sorted (or   │  │  compaction / │  │ structure    │
          │  can sort it) │  │  partition?   │  │ problem?     │
          └──────┬───────┘  └──────┬───────┘  └──────┬───────┘
                 │                  │                  │
                 ▼                  ▼                  ▼
          Opposite Ends       Same Direction      Floyd's
          (Converging)        (Write Pointer)     (Fast-Slow)
                 │                  │                  │
                 ▼                  ▼                  ▼
       ┌─────────────────┐ ┌───────────────┐ ┌───────────────┐
       │ • Two Sum II    │ │ • Remove Dups │ │ • Has Cycle   │
       │ • Three Sum     │ │ • Move Zeroes │ │ • Cycle Start │
       │ • Container     │ │ • Remove Elem │ │ • Find Middle │
       │ • Trapping Rain │ │ • Dutch Flag  │ │ • Nth from End│
       │ • Pair with     │ │ • Partition   │ │ • Happy Number│
       │   given diff    │ │   Array       │ │ • Find Dup #  │
       └─────────────────┘ └───────────────┘ └───────────────┘
```

### Quick Decision Checklist

1. **Need pairs summing to target in sorted array?** → Opposite ends. O(n).
2. **Need pairs summing to target in unsorted array?** → Hash map is simpler and also O(n). Two pointers requires sorting first → O(n log n).
3. **Need triplets summing to target?** → Sort + fix one + opposite ends for the pair. O(n^2).
4. **Need to remove/filter elements in-place?** → Same direction write pointer. O(n), O(1) space.
5. **Need to partition into 2 groups in-place?** → Same direction (like Lomuto partition).
6. **Need to partition into 3 groups in-place?** → Dutch National Flag (3 pointers).
7. **Linked list: cycle, middle, or nth-from-end?** → Floyd's fast-slow.
8. **"Find the duplicate number" in `[1..n]` with O(1) space?** → Floyd's on implicit linked list (value as next-pointer).

---

## 6. Visual Diagrams

### 6.1 Opposite Ends: Two Sum Sorted

**Input:** `[1, 3, 5, 7, 11, 13]`, target = `12`

```
Step 1:  [1, 3, 5, 7, 11, 13]
          L                 R     sum = 1+13 = 14 > 12  → R--

Step 2:  [1, 3, 5, 7, 11, 13]
          L             R         sum = 1+11 = 12 = 12  → FOUND! (L=0, R=4)

Return (1, 5)  ← 1-indexed
```

**When sum is too small:**

```
Target = 16

Step 1:  [1, 3, 5, 7, 11, 13]
          L                 R     sum = 1+13 = 14 < 16  → L++
                                  (1 is too small for ANY right partner)

Step 2:  [1, 3, 5, 7, 11, 13]
             L              R     sum = 3+13 = 16 = 16  → FOUND!
```

### 6.2 Same Direction: Remove Duplicates from Sorted Array

**Input:** `[1, 1, 2, 2, 3]`

```
Initial: [1, 1, 2, 2, 3]
          S
          F

Step 1:  [1, 1, 2, 2, 3]     F=1: nums[1]=1 == nums[S-1]=1 → skip
          S  F

Step 2:  [1, 1, 2, 2, 3]     F=2: nums[2]=2 != nums[S-1]=1 → write!
          S     F

         [1, 2, 2, 2, 3]     nums[S] = nums[F], S++
             S  F

Step 3:  [1, 2, 2, 2, 3]     F=3: nums[3]=2 == nums[S-1]=2 → skip
             S     F

Step 4:  [1, 2, 2, 2, 3]     F=4: nums[4]=3 != nums[S-1]=2 → write!
             S        F

         [1, 2, 3, 2, 3]     nums[S] = nums[F], S++
                S     F

Done. Return S = 3. Array prefix: [1, 2, 3]
```

### 6.3 Container With Most Water: Why Move the Shorter Side

**Input:** `[3, 1, 6, 4, 2]`

```
         |
    |    |
    |    |  |
  | |    |  |  |
  | |    |  |  |
  3 1    6  4  2
  L               R

  Area = min(3,2) * 4 = 8

  Key insight: width will DECREASE by 1 no matter what.
  If we move L (the taller side, 3), we MIGHT get a shorter min-height.
  If we move R (the shorter side, 2), we MIGHT get a taller min-height.

  The shorter side is the bottleneck. Keeping it and shrinking width
  can ONLY make things worse or the same. Moving it away gives a
  CHANCE of improvement.

  So: always move the pointer pointing to the shorter line.

  Move R-- :
         |
    |    |
    |    |  |
  | |    |  |
  | |    |  |
  3 1    6  4
  L          R     Area = min(3,4) * 3 = 9   ← improved!

  Move L++ :
         |
    |    |
    |    |  |
    |    |  |
    |    |  |
    1    6  4
    L       R      Area = min(1,4) * 2 = 2

  Move L++ :
         |
         |  |
         |  |
         |  |
         |  |
         6  4
         L  R      Area = min(6,4) * 1 = 4

  Max area = 9
```

### 6.4 Floyd's Cycle Detection

**List:** `1 → 2 → 3 → 4 → 5 → 3` (cycle starts at node 3)

```
Phase 1: Detect the cycle

         1 → 2 → 3 → 4 → 5
                  ↑         │
                  └─────────┘

Step 0:  S=1, F=1
Step 1:  S=2, F=3            (slow +1, fast +2)
Step 2:  S=3, F=5            (slow +1, fast +2)
Step 3:  S=4, F=4            (slow +1, fast +2: 5→3→4)
                              fast went 5→3, then 3→4
Step 4:  S=5, F=3            (slow +1, fast +2: 4→5→3)
Step 5:  S=3, F=5
Step 6:  S=4, F=4            ← MEET at node 4

Phase 2: Find cycle start

Reset one pointer to head. Move both by 1.

         p1=1 (head), p2=4 (meeting point)
Step 1:  p1=2, p2=5
Step 2:  p1=3, p2=3          ← MEET at node 3 = cycle start!

WHY THIS WORKS:
  Let:  d = distance from head to cycle start
        k = distance from cycle start to meeting point
        C = cycle length

  At meeting: slow traveled d + k steps.
              fast traveled d + k + m*C steps (m full laps).
              fast = 2 * slow → d + k + m*C = 2(d + k)
              → m*C = d + k → d = m*C - k

  From meeting point, traveling d more steps:
              k + d = k + (m*C - k) = m*C → back at cycle start!
  From head, traveling d steps → at cycle start.
  So both arrive at cycle start simultaneously.
```

---

## 7. Self-Assessment

Answer these without looking at your code or notes. If you struggle with any, revisit the relevant section.

### Question 1
**Why does moving the shorter side in Container With Most Water never miss the optimal answer?**

<details>
<summary>Answer</summary>

The area is `min(height[L], height[R]) * (R - L)`. The width `(R - L)` decreases by 1 no matter which pointer you move. The only way to increase area is to increase the minimum height. If you keep the shorter side and move the taller side inward, the minimum height can only stay the same or decrease (it is still bounded by the shorter side). So keeping the shorter side can never produce a larger area than what you have now. You can safely discard it and move it inward. The optimal pair's taller side will never be eliminated because we only discard the shorter side.

</details>

### Question 2
**Can two pointers replace the hash map approach for unsorted two-sum?**

<details>
<summary>Answer</summary>

Not as a direct drop-in. Two pointers requires sorted input. If you sort first, the indices change, so you lose the original index information. For the problem "return the original indices of the two elements" (LeetCode 1), sorting destroys that information unless you carry the original indices along (sort an array of `(value, index)` pairs). The hash map approach is O(n) time and O(n) space on unsorted input and naturally preserves indices. Two pointers after sorting is O(n log n) time and O(1) extra space (or O(n) if you store index pairs). Choose based on whether you need original indices and whether the O(n) space matters.

</details>

### Question 3
**In the write-pointer pattern for Remove Duplicates, why does `slow` start at 1 instead of 0?**

<details>
<summary>Answer</summary>

The first element is always part of the result (there is no previous element it could be a duplicate of). So position 0 is already correctly placed. The `slow` pointer indicates "the next position to write to," which starts at index 1. The `fast` pointer also starts at 1 because we compare `nums[fast]` against `nums[slow-1]` (the last written element). Starting both at 0 would compare `nums[0]` against `nums[-1]`, which is invalid.

</details>

### Question 4
**In Floyd's algorithm, after detecting a cycle, why does resetting one pointer to the head and advancing both by 1 find the cycle start?**

<details>
<summary>Answer</summary>

Let `d` = distance from head to cycle start, `k` = distance from cycle start to meeting point, and `C` = cycle length. When they meet, slow has traveled `d + k` and fast has traveled `2(d + k)`. The difference `d + k` must be a multiple of `C` (fast did extra full laps): `d + k = mC` for some integer `m`. Rearranging: `d = mC - k`. Starting one pointer at head (needs `d` steps to reach cycle start) and one at the meeting point (needs `mC - k` steps to reach cycle start, since it is `k` past the start and must travel the remaining `mC - k` to complete `m` full laps back to the start): both travel the same distance `d` to reach the cycle start simultaneously.

</details>

### Question 5
**ThreeSum has O(n^2) time complexity. Why can't two pointers reduce it to O(n)?**

<details>
<summary>Answer</summary>

Two pointers reduces the inner pair search from O(n) to O(n) — which is the same! The two-pointer technique on the inner pair is already O(n), and you must do this for each of the O(n) fixed elements. The outer loop is unavoidable: you need to consider each element as a potential member of the triplet. There is no known way to solve 3Sum in better than O(n^2) time (and it is conjectured that O(n^2) is optimal). Two pointers eliminates the *constant factor* and *space* cost compared to using a hash set for the inner search, but does not change the asymptotic complexity.

</details>
