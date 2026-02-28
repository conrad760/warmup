# Day 19 — Intervals & Sweep Line: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Why It's Useful | Time |
|---|----------|----------------|------|
| 1 | [NeetCode — Merge Intervals (video)](https://www.youtube.com/watch?v=44H3cEC2fFM) | Visual walkthrough of the sort-then-scan pattern with clear timeline diagrams. Establishes the mental model for all interval problems. | 12 min |
| 2 | [NeetCode — Meeting Rooms II (video)](https://www.youtube.com/watch?v=FdzJmTCVyJU) | Shows both the min-heap approach and the sweep line alternative side by side. Great for seeing why the heap tracks "active meetings." | 15 min |
| 3 | [NeetCode — Non-overlapping Intervals (video)](https://www.youtube.com/watch?v=nONCGxWoUfM) | Explains why sorting by END time is essential for the greedy selection, with the exchange argument intuitively presented. | 10 min |
| 4 | [CP-Algorithms: Sweep Line](https://cp-algorithms.com/geometry/sweep_line.html) | Formal treatment of the sweep line technique. Goes beyond intervals into computational geometry, but the core "events on a timeline" section is directly applicable. | 20 min |
| 5 | [LeetCode Editorial — 56. Merge Intervals](https://leetcode.com/problems/merge-intervals/editorial/) | The official editorial includes a connected-components approach (O(n^2)) vs. the sorting approach (O(n log n)), helping you see why sorting is the right move. | 10 min |
| 6 | [Back to Back SWE — Interval Scheduling Maximization](https://www.youtube.com/watch?v=hVhOeaONg1Y) | Walks through the formal exchange argument proof for why earliest-finish-time greedy is optimal. If you want rigor, this is it. | 15 min |
| 7 | [Tushar Roy — Insert Interval](https://www.youtube.com/watch?v=rlKp2ZK0gbo) | Step-by-step visual of the three-phase approach (before, overlap, after) for inserting into a sorted interval list without re-sorting. | 10 min |
| 8 | [Visualgo — Sorting Algorithms](https://visualgo.net/en/sorting) | Not interval-specific, but useful for reinforcing why O(n log n) sort is the prerequisite step. Animate the sort that kicks off every interval solution. | 5 min |

**Reading order:** Start with resource #1 (merge intervals — the foundation), then #7 (insert), #2 (meeting rooms heap), #3 (non-overlapping greedy + exchange argument), #4 (sweep line formalism). Skim #5 and #6 if time allows.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 - 12:20 — Review & Internalize (20 min)

| Time | Activity |
|------|----------|
| 12:00 - 12:05 | Read the core concepts section below. No code. Draw the overlap condition `a < d AND c < b` on paper. Convince yourself it handles all cases. |
| 12:05 - 12:10 | Study the interval problem taxonomy (Section 5). For each category, write down the sorting key (start vs. end) and the core technique (scan, heap, greedy, sweep). |
| 12:10 - 12:15 | Trace through the ASCII diagrams (Section 6) by hand. For the merge example, walk through the sorted intervals and simulate the merge logic. |
| 12:15 - 12:20 | Review the complexity table from the OVERVIEW. Commit to memory: all interval problems are O(n log n) due to the sort. The scan/heap step is O(n). |

### 12:20 - 1:20 — Implement from Scratch (60 min)

| Time | Problem | Key Technique | Notes |
|------|---------|---------------|-------|
| 12:20 - 12:35 | **Merge Intervals** | Sort by start, extend-or-append | This is the foundation. Get it clean: sort, iterate, compare `current.start <= last.end`, extend with `max(last.end, current.end)`. If you nail this, everything else is a variation. |
| 12:35 - 12:50 | **Insert Interval** | Three-phase scan | Don't re-sort. Walk the list: (1) add all intervals ending before newInterval starts, (2) merge all overlapping intervals with newInterval, (3) add the rest. This tests your overlap condition under pressure. |
| 12:50 - 1:05 | **Meeting Rooms II** | Min-heap of end times | Sort by start. For each meeting: if heap's min end <= current start, pop (room freed). Push current end. Answer = max heap size. This is where you connect "max simultaneous" to "min resources." |
| 1:05 - 1:20 | **Non-overlapping Intervals** (Erase Overlap) | Sort by END, greedy pick | Sort by end time. Track `lastEnd`. For each interval: if `start >= lastEnd`, keep it, update `lastEnd`. Else skip it (one erasure). Count erasures. Think about WHY end-time sorting works (exchange argument). |

### 1:20 - 1:50 — Solidify & Extend (30 min)

| Time | Activity |
|------|----------|
| 1:20 - 1:35 | **Sweep Line implementation** — Implement `SweepLineMaxOverlap`. Convert each interval to two events: `(start, +1)` and `(end, -1)`. Sort events (break ties: -1 before +1 if endpoints are exclusive, +1 before -1 if inclusive — know which your problem needs). Sweep and track running count. Verify it gives the same answer as Meeting Rooms II. |
| 1:35 - 1:45 | **Edge case gauntlet** — Run each implementation against: (1) touching intervals `[1,2],[2,3]`, (2) nested intervals `[1,10],[2,5],[3,4]`, (3) single interval, (4) already sorted input, (5) all overlapping `[1,5],[2,6],[3,7]`, (6) no overlaps `[1,2],[3,4],[5,6]`. |
| 1:45 - 1:50 | **Compare approaches** — For Meeting Rooms II, write down the heap approach vs. sweep line. When would you choose one over the other? (Heap: when you need to track which meetings are in each room. Sweep: when you only need the count.) |

### 1:50 - 2:00 — Recap (10 min)

Write down from memory:
1. The overlap condition for two intervals.
2. The complexity of each problem (all O(n log n)).
3. When to sort by start vs. end.
4. One gotcha (e.g., touching intervals, or why sort-by-start fails for non-overlapping selection).

---

## 3. Core Concepts Deep Dive

### The Interval Toolkit: Sort + Scan

Nearly every interval problem starts the same way:
1. **Sort** the intervals (by start or by end — knowing which is critical).
2. **Scan** left to right, maintaining some state (merged list, heap, counter, last-picked end).

The sort costs O(n log n). The scan is O(n). Total: **O(n log n)** dominated by the sort. This is true for merge, insert, meeting rooms, non-overlapping intervals, and sweep line.

### Overlap Detection: The One Condition

Two intervals `[a, b]` and `[c, d]` overlap when:

```
a < d  AND  c < b
```

Equivalently, they do NOT overlap when `a >= d OR c >= b` (one ends before the other starts).

For problems where touching counts as overlapping (e.g., merge intervals), use `<=`:
```
a <= d  AND  c <= b
```

This single condition is the atom from which every interval algorithm is built. Memorize it.

**Visual proof:**

```
Case 1: Overlap            Case 2: No overlap
a---------b                a----b
    c---------d                      c----d
a < d? YES                 a < d? YES
c < b? YES                 c < b? NO   --> no overlap
--> overlap
```

### Merge Intervals

**Algorithm:**
1. Sort intervals by start time.
2. Initialize result with the first interval.
3. For each subsequent interval:
   - If `current.start <= last.end` → overlap → extend: `last.end = max(last.end, current.end)`
   - Else → no overlap → append current to result.

**Why sorting by start is sufficient:** After sorting by start, we only need to check if the current interval's start falls within (or touches) the last merged interval. If it does, the intervals overlap and we extend. If it doesn't, since all future intervals have even larger starts, none of them can overlap with the last merged interval either. The scan is correct and complete.

**Why sorting by end is wrong for merge:** If you sort by end, an interval with a small start but large end could appear late in the sorted order, missing earlier overlaps.

### Meeting Rooms: Max Simultaneous = Min Rooms

The key insight: **the maximum number of meetings happening at any point in time equals the minimum number of rooms needed.** This is a direct application of the pigeonhole principle — if 3 meetings overlap at time t, you need at least 3 rooms at time t.

**Heap approach (sort by start):**
- Sort meetings by start time.
- Maintain a min-heap of end times (each entry = one active meeting's end time).
- For each meeting:
  - While the heap's minimum end time <= current meeting's start: pop (that room is free).
  - Push the current meeting's end time.
- The maximum heap size at any point during the scan is the answer.

Why this works: the heap represents "currently active meetings." When we process a new meeting, we first free up any rooms whose meetings have ended. Then we add the new meeting. The heap size is the number of rooms in use right now.

### Non-overlapping Interval Selection: Sort by END

**Problem:** Given intervals, remove the minimum number of intervals so that no two overlap. Equivalently: find the **maximum set of non-overlapping intervals**.

**Algorithm:**
1. Sort by **end time** (ascending).
2. Track `lastEnd = -infinity`.
3. For each interval: if `start >= lastEnd`, keep it, update `lastEnd = end`. Else, skip it.
4. Answer = total intervals - kept intervals.

**Why sort by end time (the exchange argument):**

Suppose the greedy picks interval G (earliest end) and an optimal solution picks interval O instead (O ends later). We can **swap O for G** in the optimal solution:
- G ends no later than O, so G leaves at least as much room for future intervals.
- Everything the optimal solution picked after O still fits after G.
- Therefore the swap doesn't reduce the number of picked intervals.

This exchange argument proves the greedy choice (earliest end) is always safe. **Sorting by start time does NOT work** because an interval with an early start but very late end blocks many future intervals:

```
Sort by start (WRONG):       Sort by end (CORRECT):
[1-----------10]  picked     [1--3]  picked
  [2--4] skipped             [2--4]  skipped (overlaps)
     [5--7] skipped             [5--7]  picked
                             [1-----------10]  skipped
Greedy by start: 1 interval  Greedy by end: 2 intervals
```

### Sweep Line: Events on a Timeline

The sweep line technique converts intervals into discrete events:

1. For each interval `[s, e]`, create two events: `(s, +1)` and `(e, -1)`.
2. Sort all events by coordinate. **Tie-breaking matters:**
   - If `[1,3]` and `[3,5]` should NOT count as overlapping (meeting rooms): process `-1` before `+1` at the same coordinate (one ends before the other starts).
   - If they SHOULD count as overlapping (merge intervals): process `+1` before `-1`.
3. Sweep left to right, maintaining a running count. The maximum count is the maximum number of simultaneous intervals.

Sweep line is more general than the heap approach — it extends naturally to weighted intervals, 2D sweep, and computational geometry problems.

### When to Sort by Start vs. End

| Problem | Sort By | Why |
|---------|---------|-----|
| Merge Intervals | **Start** | You need to process intervals left-to-right and extend the current merge window. Start-sorted ensures you see overlaps in order. |
| Insert Interval | **Start** (assumed pre-sorted) | Same merge logic — walk left to right. |
| Meeting Rooms I & II | **Start** | You process meetings in chronological order, checking what's still active. |
| Non-overlapping Selection | **End** | Greedy needs the earliest-finishing interval to maximize remaining room. |
| Sweep Line | **Event coordinate** | Neither start nor end of the interval — you sort the events themselves. |

**Rule of thumb:** If you're **merging or counting active intervals**, sort by start. If you're **selecting a maximum non-overlapping set**, sort by end.

---

## 4. Implementation Checklist

### Function Signatures

```go
package intervals

import (
    "container/heap"
    "sort"
)

// MergeIntervals merges all overlapping intervals.
// Input:  [[1,3],[2,6],[8,10],[15,18]]
// Output: [[1,6],[8,10],[15,18]]
func MergeIntervals(intervals [][]int) [][]int {
    // 1. Sort by start
    // 2. Scan: extend or append
}

// InsertInterval inserts a new interval into a sorted, non-overlapping list
// and merges if necessary.
// Input:  intervals=[[1,3],[6,9]], newInterval=[2,5]
// Output: [[1,5],[6,9]]
func InsertInterval(intervals [][]int, newInterval []int) [][]int {
    // Three phases: before, overlap, after
}

// MinMeetingRooms returns the minimum number of conference rooms needed.
// Uses a min-heap of end times.
// Input:  [[0,30],[5,10],[15,20]]
// Output: 2
func MinMeetingRooms(intervals [][]int) int {
    // 1. Sort by start
    // 2. Min-heap of end times
    // 3. For each meeting: pop expired, push current end
    // 4. Track max heap size
}

// EraseOverlapIntervals returns the minimum number of intervals to remove
// so that no two intervals overlap.
// Input:  [[1,2],[2,3],[3,4],[1,3]]
// Output: 1
func EraseOverlapIntervals(intervals [][]int) int {
    // 1. Sort by END time
    // 2. Greedy: keep if start >= lastEnd
    // 3. Return total - kept
}

// SweepLineMaxOverlap returns the maximum number of intervals overlapping
// at any point in time.
// Input:  [[0,30],[5,10],[15,20]]
// Output: 2
func SweepLineMaxOverlap(intervals [][]int) int {
    // 1. Create events: (coord, +1/-1)
    // 2. Sort by coord (break ties: -1 before +1)
    // 3. Sweep, track running count, record max
}

// ---------- Min-heap for Meeting Rooms II ----------

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}
```

### Test Cases

Write tests covering these scenarios for every function:

| Scenario | Example Input | Why It Matters |
|----------|---------------|----------------|
| **Touching intervals** | `[[1,2],[2,3]]` | Tests boundary: overlap or not? Merge: yes. Meeting rooms: depends on problem definition. |
| **Nested intervals** | `[[1,10],[2,5],[3,4]]` | The inner intervals are fully contained. Merge should produce `[[1,10]]`. Meeting rooms: 3 rooms at peak. |
| **Single interval** | `[[5,8]]` | Base case. Should return as-is for merge, 1 for meeting rooms, 0 erasures. |
| **Already sorted** | `[[1,3],[4,6],[7,9]]` | Sort is a no-op. Verifies logic doesn't depend on sort shuffling things. |
| **All overlapping** | `[[1,5],[2,6],[3,7]]` | Merge: `[[1,7]]`. Meeting rooms: 3. Erasures: 2. |
| **No overlaps** | `[[1,2],[3,4],[5,6]]` | Merge: unchanged. Meeting rooms: 1. Erasures: 0. |
| **Reverse sorted** | `[[5,8],[3,6],[1,4]]` | Tests that sort works correctly before the scan. |
| **Empty input** | `[]` | Guard clause. Return empty / 0. |

```go
package intervals

import (
    "reflect"
    "testing"
)

func TestMergeIntervals(t *testing.T) {
    tests := []struct {
        name     string
        input    [][]int
        expected [][]int
    }{
        {"touching", [][]int{{1, 2}, {2, 3}}, [][]int{{1, 3}}},
        {"nested", [][]int{{1, 10}, {2, 5}, {3, 4}}, [][]int{{1, 10}}},
        {"single", [][]int{{5, 8}}, [][]int{{5, 8}}},
        {"no overlap", [][]int{{1, 2}, {3, 4}, {5, 6}}, [][]int{{1, 2}, {3, 4}, {5, 6}}},
        {"all overlap", [][]int{{1, 5}, {2, 6}, {3, 7}}, [][]int{{1, 7}}},
        {"standard", [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][]int{{1, 6}, {8, 10}, {15, 18}}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MergeIntervals(tt.input)
            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}

func TestMinMeetingRooms(t *testing.T) {
    tests := []struct {
        name     string
        input    [][]int
        expected int
    }{
        {"nested", [][]int{{1, 10}, {2, 5}, {3, 4}}, 3},
        {"single", [][]int{{5, 8}}, 1},
        {"no overlap", [][]int{{1, 2}, {3, 4}, {5, 6}}, 1},
        {"all overlap", [][]int{{1, 5}, {2, 6}, {3, 7}}, 3},
        {"standard", [][]int{{0, 30}, {5, 10}, {15, 20}}, 2},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MinMeetingRooms(tt.input)
            if got != tt.expected {
                t.Errorf("got %d, want %d", got, tt.expected)
            }
        })
    }
}

func TestEraseOverlapIntervals(t *testing.T) {
    tests := []struct {
        name     string
        input    [][]int
        expected int
    }{
        {"standard", [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}, 1},
        {"all overlap", [][]int{{1, 5}, {2, 6}, {3, 7}}, 2},
        {"no overlap", [][]int{{1, 2}, {3, 4}, {5, 6}}, 0},
        {"single", [][]int{{5, 8}}, 0},
        {"nested", [][]int{{1, 10}, {2, 3}, {4, 5}}, 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := EraseOverlapIntervals(tt.input)
            if got != tt.expected {
                t.Errorf("got %d, want %d", got, tt.expected)
            }
        })
    }
}

func TestSweepLineMaxOverlap(t *testing.T) {
    tests := []struct {
        name     string
        input    [][]int
        expected int
    }{
        {"standard", [][]int{{0, 30}, {5, 10}, {15, 20}}, 2},
        {"all overlap", [][]int{{1, 5}, {2, 6}, {3, 7}}, 3},
        {"no overlap", [][]int{{1, 2}, {3, 4}, {5, 6}}, 1},
        {"nested", [][]int{{1, 10}, {2, 5}, {3, 4}}, 3},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := SweepLineMaxOverlap(tt.input)
            if got != tt.expected {
                t.Errorf("got %d, want %d", got, tt.expected)
            }
        })
    }
}
```

---

## 5. The Interval Problem Taxonomy

### Category 1: Merge / Consolidate

**Question:** Combine overlapping intervals into non-overlapping ones.

| Aspect | Detail |
|--------|--------|
| **Sort by** | Start time |
| **Approach** | Sort, scan left to right, extend or create new merged interval |
| **Key operation** | `max(last.end, current.end)` to extend |
| **Problems** | Merge Intervals (LC 56), Insert Interval (LC 57) |
| **Complexity** | O(n log n) time, O(n) space |

### Category 2: Scheduling / Resource Allocation

**Question:** Can all tasks be done? How many resources are needed?

| Aspect | Detail |
|--------|--------|
| **Sort by** | Start time |
| **Approach** | Min-heap of end times (or sweep line for counting only) |
| **Key insight** | Max simultaneous overlap = min resources needed |
| **Problems** | Meeting Rooms I (LC 252), Meeting Rooms II (LC 253) |
| **Complexity** | O(n log n) time, O(n) space |

### Category 3: Selection / Maximum Non-overlapping Set

**Question:** What's the maximum number of non-overlapping intervals you can keep?

| Aspect | Detail |
|--------|--------|
| **Sort by** | **End time** (critical difference) |
| **Approach** | Greedy: pick earliest-ending interval that doesn't conflict |
| **Key insight** | Earliest end leaves maximum room for future picks (exchange argument) |
| **Problems** | Non-overlapping Intervals (LC 435), Activity Selection |
| **Complexity** | O(n log n) time, O(1) space (beyond sort) |

### Category 4: Coverage

**Question:** Does a set of intervals fully cover a given range?

| Aspect | Detail |
|--------|--------|
| **Sort by** | Start time |
| **Approach** | Sort, scan, check that each interval starts at or before the current coverage boundary, extend the boundary to `max(boundary, end)` |
| **Key insight** | Gaps in coverage appear when `current.start > boundary` |
| **Problems** | Video Stitching (LC 1024), Minimum Number of Arrows to Burst Balloons (LC 452) |
| **Complexity** | O(n log n) time, O(1) space |

### Category 5: Intersection

**Question:** Find the common overlap among intervals.

| Aspect | Detail |
|--------|--------|
| **Sort by** | Start time (both lists) |
| **Approach** | Two pointers, one per list. The intersection of `[a,b]` and `[c,d]` is `[max(a,c), min(b,d)]` if it's non-empty. Advance the pointer with the smaller end. |
| **Key insight** | The overlap condition `max(a,c) <= min(b,d)` determines if an intersection exists |
| **Problems** | Interval List Intersections (LC 986) |
| **Complexity** | O(m + n) time (already sorted), O(1) space |

### Quick Decision Chart

```
Got intervals?
    |
    +--> Need to combine overlapping? --> MERGE (sort by start)
    |
    +--> How many resources/rooms? --> SCHEDULING (sort by start, heap)
    |
    +--> Max non-overlapping set? --> SELECTION (sort by END, greedy)
    |
    +--> Does it cover a range? --> COVERAGE (sort by start, track boundary)
    |
    +--> Common overlap of two lists? --> INTERSECTION (two pointers)
    |
    +--> Just count max overlap? --> SWEEP LINE (events, sort by coord)
```

---

## 6. Visual Diagrams

### Diagram 1: Merge Intervals

```
Input intervals (unsorted):
  [8,10]        ████
  [2,6]    ██████
  [15,18]              ████
  [1,3]   ███

Timeline:  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18

Step 1 — Sort by start:
  [1,3]   ███
  [2,6]    ██████
  [8,10]              ████
  [15,18]                                ████

Step 2 — Scan and merge:
  Process [1,3]:  result = [[1,3]]
  Process [2,6]:  2 <= 3 (overlap!) → extend to [1,6]
                  result = [[1,6]]
  Process [8,10]: 8 > 6 (no overlap) → append
                  result = [[1,6],[8,10]]
  Process [15,18]: 15 > 10 (no overlap) → append
                  result = [[1,6],[8,10],[15,18]]

Output:
  [1,6]   ████████
  [8,10]              ████
  [15,18]                                ████

Timeline:  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15 16 17 18
```

### Diagram 2: Meeting Rooms (Heap Approach)

```
Meetings sorted by start:
  A: [0 ================================ 30]
  B:      [5 ==== 10]
  C:                    [15 ==== 20]

Timeline:  0     5     10    15    20    25    30

Processing with min-heap (tracking end times):

  t=0:  Process A[0,30]   heap: [30]           rooms=1
  t=5:  Process B[5,10]   heap min=30 > 5      rooms=2
        (room not free)   heap: [10, 30]
  t=15: Process C[15,20]  heap min=10 <= 15    rooms=2
        (pop 10, room freed) heap: [30]
        push 20           heap: [20, 30]

  Max heap size = 2  →  Answer: 2 rooms

Room assignment:
  Room 1: |====== A [0,30] ==============================|
  Room 2:      |= B [5,10] =|    |= C [15,20] =|

Timeline:  0     5     10    15    20    25    30
```

### Diagram 3: Sweep Line Events

```
Intervals: [0,30], [5,10], [15,20]

Events:
  (0, +1)  (5, +1)  (10, -1)  (15, +1)  (20, -1)  (30, -1)

Sweep:
              +1        +1        -1          +1        -1         -1
  count:  0 ──→ 1 ──→ 2 ──→ 1 ──→ 2 ──→ 1 ──→ 0
              ↑                         ↑
              |                         |
  Event:     0    5       10      15       20        30

  Number   2 |         ██
  line     1 |   ██████    ███████    ████████
  (count)  0 |███                               ██████████
             +---+----+----+----+----+----+----+-->
              0   5   10   15   20   25   30

  Max count = 2  →  Answer: 2 simultaneous intervals
```

### Diagram 4: Non-overlapping Selection (Sort by End, Greedy)

```
Input intervals:
  [1,10]  ████████████████████
  [2,3]     ██
  [4,6]        ██████
  [5,7]          ██████
  [8,9]                  ██

Step 1 — Sort by END time:
  [2,3]     ██                        end=3
  [4,6]        ██████                 end=6
  [5,7]          ██████               end=7
  [8,9]                  ██           end=9
  [1,10] ████████████████████         end=10

Step 2 — Greedy scan (pick if start >= lastEnd):
  lastEnd = -∞

  [2,3]:   2 >= -∞  → PICK ✓    lastEnd = 3
  [4,6]:   4 >= 3   → PICK ✓    lastEnd = 6
  [5,7]:   5 < 6    → SKIP ✗
  [8,9]:   8 >= 6   → PICK ✓    lastEnd = 9
  [1,10]:  1 < 9    → SKIP ✗

Result: 3 kept, 2 removed

     1  2  3  4  5  6  7  8  9  10
     .  ██ .  ██████ .  .  ██ .  .    ← kept (3 intervals)
     .  .  .  .  █████. .  .  .  .    ← skipped [5,7]
     ████████████████████████████     ← skipped [1,10]

Note: Sorting by START would have picked [1,10] first,
      blocking everything. Only 1 interval kept instead of 3.
```

---

## 7. Self-Assessment

Answer these after your session. If you can't answer confidently, revisit the relevant section.

### Question 1: Why Sort by End for Selection?

> Why does sorting by end time work for the non-overlapping interval selection problem, but sorting by start time doesn't?

**What a good answer includes:**
- Earliest end leaves maximum room for future intervals.
- The exchange argument: replacing any chosen interval with an earlier-ending one never hurts.
- A concrete counterexample: `[1,10], [2,3], [4,5]`. Sort by start picks `[1,10]` first (1 interval). Sort by end picks `[2,3]` then `[4,5]` (2 intervals).

### Question 2: Sweep Line ↔ Meeting Rooms

> What's the relationship between the sweep line maximum count and the answer to Meeting Rooms II? Can you always use one in place of the other?

**What a good answer includes:**
- They compute the same value: the maximum number of simultaneously active intervals.
- Sweep line is O(n log n) and simpler when you only need the count.
- The heap approach is also O(n log n) but additionally lets you track *which* meetings are in each room (the heap entries can carry meeting identity).
- Both are interchangeable for the count. The heap is more powerful when you need to assign resources to specific intervals.

### Question 3: The Overlap Condition

> Two intervals `[a,b]` and `[c,d]` overlap when `a < d AND c < b`. Derive this condition from first principles. What changes if touching endpoints count as overlapping?

**What a good answer includes:**
- Two intervals do NOT overlap when one ends before the other starts: `b <= c OR d <= a`.
- Negating (De Morgan's): they DO overlap when `b > c AND d > a`, which rearranges to `a < d AND c < b`.
- For touching-counts-as-overlapping: change `<` to `<=` in the overlap check (or equivalently, change `<=` to `<` in the non-overlap check).

### Question 4: Insert Without Re-sorting

> Insert Interval (LC 57) runs in O(n) even though the input is already sorted. Why don't we need to re-sort after inserting and merging?

**What a good answer includes:**
- The input list is already sorted and non-overlapping.
- The three-phase approach processes intervals in order: (1) all before the new interval, (2) all overlapping with the new interval (merged together), (3) all after.
- Since we process in sorted order and merge is associative, the output is guaranteed to be sorted and non-overlapping. No re-sort needed.

### Question 5: Heap Size vs. Running Maximum

> In Meeting Rooms II, why is the answer the *maximum* heap size during the scan, not the *final* heap size?

**What a good answer includes:**
- The heap can shrink during the scan as meetings end. The final heap size only reflects meetings still active at the very end.
- The peak simultaneous count may occur in the middle of the timeline.
- Example: `[1,5],[2,3],[4,6]`. At time 2: heap = `{3,5}` (size 2). At time 4: we pop 3, push 6, heap = `{5,6}` (size 2). At time 5: still 2. But if we had `[1,3],[2,4],[5,6]`: peak is 2, final heap has only `[6]` (size 1).

---

## Appendix: Complexity Summary

| Problem | Time | Space | Sort Key |
|---------|------|-------|----------|
| Merge Intervals | O(n log n) | O(n) | Start |
| Insert Interval | O(n) | O(n) | Pre-sorted by start |
| Meeting Rooms I | O(n log n) | O(1) | Start |
| Meeting Rooms II (heap) | O(n log n) | O(n) | Start |
| Meeting Rooms II (sweep) | O(n log n) | O(n) | Event coordinate |
| Non-overlapping Intervals | O(n log n) | O(1) | **End** |
| Interval List Intersection | O(m + n) | O(1) | Pre-sorted by start |

**The universal truth:** Every interval problem is O(n log n) because of the sort. The scan/heap/greedy step after sorting is O(n). If the input is pre-sorted, you drop to O(n).
