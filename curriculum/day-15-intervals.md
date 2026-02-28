# Day 15: Intervals & Sweep Line

> **Time:** 2 hours | **Level:** Refresher | **Language:** Go

Interval problems are among the most frequently asked in interviews, and nearly all of
them reduce to one of six patterns. The code is never long — sorting + a single pass
covers most cases. The real skill is choosing the right sort order, recognizing when a
heap or sweep line is needed, and handling the edge cases around touching boundaries.

This day bridges directly from Day 14's interval scheduling greedy (sort by end time) and
Day 10's heap boilerplate (min-heap of end times for meeting rooms).

---

## Pattern Catalog

---

### Pattern 1: Merge Intervals

**Trigger:** "Merge all overlapping intervals," "return non-overlapping intervals that cover all input intervals."

**Core idea:** Sort by start time. Walk through the sorted list. If the current interval overlaps with the last merged interval (current start <= last merged end), extend the last merged interval's end. Otherwise, start a new merged interval.

**Template (Merge Intervals, LC 56):**

```go
func merge(intervals [][]int) [][]int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    merged := [][]int{intervals[0]}

    for i := 1; i < len(intervals); i++ {
        last := merged[len(merged)-1]

        if intervals[i][0] <= last[1] {
            // Overlap — extend the end of the last merged interval
            if intervals[i][1] > last[1] {
                last[1] = intervals[i][1]
            }
        } else {
            // No overlap — start a new interval
            merged = append(merged, intervals[i])
        }
    }

    return merged
}
```

**Complexity:** O(n log n) time (sorting), O(n) space for the result.

**Watch out:**
- Compare current interval's **start** to last merged interval's **end** — not start-to-start. This is the #1 bug.
- When extending, take `max(last[1], intervals[i][1])`. The current interval might be entirely contained within the last merged interval (e.g., `[1,10]` and `[2,5]`).
- Touching intervals `[1,2]` and `[2,3]`: these are typically considered overlapping (merged to `[1,3]`). The `<=` in the condition handles this. Read the problem statement — some problems use strict `<`.

---

### Pattern 2: Insert Interval

**Trigger:** "Insert a new interval into a sorted list of non-overlapping intervals and merge if necessary."

**Core idea:** Three phases: (1) add all intervals that end before the new interval starts, (2) merge all intervals that overlap with the new interval, (3) add all intervals that start after the new interval ends.

**Template (Insert Interval, LC 57):**

```go
func insert(intervals [][]int, newInterval []int) [][]int {
    result := [][]int{}
    i := 0
    n := len(intervals)

    // Phase 1: all intervals that end before newInterval starts
    for i < n && intervals[i][1] < newInterval[0] {
        result = append(result, intervals[i])
        i++
    }

    // Phase 2: merge overlapping intervals into newInterval
    for i < n && intervals[i][0] <= newInterval[1] {
        if intervals[i][0] < newInterval[0] {
            newInterval[0] = intervals[i][0]
        }
        if intervals[i][1] > newInterval[1] {
            newInterval[1] = intervals[i][1]
        }
        i++
    }
    result = append(result, newInterval)

    // Phase 3: all intervals that start after newInterval ends
    for i < n {
        result = append(result, intervals[i])
        i++
    }

    return result
}
```

**Complexity:** O(n) time (input is already sorted), O(n) space for the result.

**Watch out:**
- Phase 1 condition: `intervals[i][1] < newInterval[0]` — strictly less than. An interval ending at 5 overlaps with a new interval starting at 5.
- Phase 2 condition: `intervals[i][0] <= newInterval[1]` — an interval starting at 7 overlaps with a new interval ending at 7.
- Don't re-sort the input. The problem guarantees intervals are already sorted and non-overlapping. Sorting would waste O(n log n) and signals to the interviewer that you missed the constraint.

---

### Pattern 3: Meeting Rooms / Maximum Overlap (Heap Approach)

**Trigger:** "Minimum number of meeting rooms," "minimum platforms," "maximum number of overlapping intervals at any point."

**Core idea:** Sort by start time. Use a min-heap tracking end times of active meetings. For each new meeting, if its start >= the earliest end in the heap, pop that ended meeting (the room is freed). Push the current meeting's end time. The max heap size at any point is the answer.

**Template (Meeting Rooms II, LC 253):**

```go
func minMeetingRooms(intervals [][]int) int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    h := &MinHeap{}     // min-heap of end times (see Day 10 boilerplate)
    heap.Init(h)
    maxRooms := 0

    for _, iv := range intervals {
        // Free up rooms whose meetings have ended
        for h.Len() > 0 && (*h)[0] <= iv[0] {
            heap.Pop(h)
        }

        // This meeting needs a room
        heap.Push(h, iv[1])

        if h.Len() > maxRooms {
            maxRooms = h.Len()
        }
    }

    return maxRooms
}
```

**Complexity:** O(n log n) time (sort + heap operations), O(n) space for the heap.

**Watch out:**
- The heap condition is `(*h)[0] <= iv[0]`, meaning a meeting ending at time 5 frees the room for a meeting starting at time 5. If the problem says "a room is available only strictly after the previous meeting ends," use `<` instead. **Read the problem carefully.**
- Only pop meetings with end time <= current start. Don't pop ALL meetings — some might still be running.
- The answer is `maxRooms`, not `h.Len()` at the end. The heap size decreases as meetings end, but you need the peak.

---

### Pattern 4: Non-Overlapping Interval Selection (Greedy)

**Trigger:** "Maximum non-overlapping intervals," "minimum removals to eliminate overlaps."

This was covered in Day 14 (Pattern 3: Interval Scheduling). Recap for completeness:

**Core idea:** Sort by **end** time. Greedily keep each interval that doesn't overlap with the last kept interval.

**Template (Non-Overlapping Intervals, LC 435):**

```go
func eraseOverlapIntervals(intervals [][]int) int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][1] < intervals[j][1]
    })

    removals := 0
    prevEnd := intervals[0][1]

    for i := 1; i < len(intervals); i++ {
        if intervals[i][0] < prevEnd {
            removals++                  // overlap — remove this interval
        } else {
            prevEnd = intervals[i][1]   // no overlap — keep it
        }
    }

    return removals
}
```

**Complexity:** O(n log n) time, O(1) extra space.

**Watch out:**
- **Sort by END time, not start time.** This is the most common interval greedy mistake. See Day 14 for the full counterexample.
- `minimum removals = total - maximum non-overlapping`. Either formulation works.

---

### Pattern 5: Sweep Line (Event-Based)

**Trigger:** "Maximum number of overlapping intervals at any point," "minimum platforms needed," "sky line," "rectangle area union," any problem where you need to track how many intervals are active at each point in time.

**Core idea:** Decompose each interval `[start, end]` into two events: `+1` at `start` (something begins) and `-1` at `end` (something ends). Sort all events by time. Sweep through, maintaining a running count. The maximum count during the sweep is the answer.

This is an alternative to the heap approach in Pattern 3. Both solve the same problems; sweep line is often simpler to code and easier to reason about.

**Template (Meeting Rooms II via Sweep Line):**

```go
func minMeetingRooms(intervals [][]int) int {
    events := make([][2]int, 0, len(intervals)*2)

    for _, iv := range intervals {
        events = append(events, [2]int{iv[0], 1})   // meeting starts: +1
        events = append(events, [2]int{iv[1], -1})  // meeting ends:   -1
    }

    // Sort by time; if times are equal, process ends (-1) before starts (+1)
    sort.Slice(events, func(i, j int) bool {
        if events[i][0] != events[j][0] {
            return events[i][0] < events[j][0]
        }
        return events[i][1] < events[j][1] // -1 before +1
    })

    maxRooms := 0
    active := 0

    for _, e := range events {
        active += e[1]
        if active > maxRooms {
            maxRooms = active
        }
    }

    return maxRooms
}
```

**Complexity:** O(n log n) time (sorting events), O(n) space for events.

**Watch out:**
- **Tie-breaking is critical.** When an end event and a start event have the same time, which goes first?
  - If a room freed at time 5 can be reused by a meeting starting at time 5 → process ends before starts (`-1` before `+1`). This is the common case for meeting rooms.
  - If starting at time 5 means you need the room before the other meeting leaves at 5 → process starts before ends. This is less common but shows up in some problems.
  - **Always ask yourself: does an interval ending at time T conflict with one starting at time T?** Then encode the answer in the tie-breaker.
- Sweep line generalizes beyond 1D intervals. For rectangle union area, you sweep a vertical line across x-coordinates and track active y-ranges. The principle is the same.

**Template (Sweep Line for counting overlaps at each point — generic):**

```go
func maxOverlap(intervals [][]int) int {
    diff := map[int]int{}

    for _, iv := range intervals {
        diff[iv[0]]++
        diff[iv[1]]--
    }

    // Collect and sort the keys
    times := make([]int, 0, len(diff))
    for t := range diff {
        times = append(times, t)
    }
    sort.Ints(times)

    maxCount := 0
    count := 0
    for _, t := range times {
        count += diff[t]
        if count > maxCount {
            maxCount = count
        }
    }

    return maxCount
}
```

This "difference map" variant is handy when you want ends and starts at the same time to be processed together (aggregated per timestamp).

---

### Pattern 6: Interval Intersection (Two Pointers)

**Trigger:** "Intersection of two sorted interval lists," "common free time," "find overlapping portions between two schedules."

**Core idea:** Use two pointers, one for each sorted interval list. At each step, compute the overlap of the current pair (if any). Advance the pointer whose interval ends first — the other interval might still overlap with the next interval in the first list.

**Template (Interval List Intersections, LC 986):**

```go
func intervalIntersection(firstList [][]int, secondList [][]int) [][]int {
    result := [][]int{}
    i, j := 0, 0

    for i < len(firstList) && j < len(secondList) {
        // Compute the overlap
        lo := max(firstList[i][0], secondList[j][0])
        hi := min(firstList[i][1], secondList[j][1])

        if lo <= hi {
            result = append(result, []int{lo, hi})
        }

        // Advance the pointer with the earlier end
        if firstList[i][1] < secondList[j][1] {
            i++
        } else {
            j++
        }
    }

    return result
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**Complexity:** O(m + n) time, O(1) extra space (excluding output).

**Watch out:**
- The overlap is `[max(a.start, b.start), min(a.end, b.end)]`. It exists only if `lo <= hi`.
- Always advance the pointer with the **earlier end time**. The interval that ends later might still intersect with the next interval from the other list.
- Both lists must be sorted and non-overlapping (within each list). If they're not, merge them first (Pattern 1).

---

## Decision Framework

Read the problem statement. Match the **first** rule that fits:

| Signal in problem | Pattern | Sort order |
|---|---|---|
| "Merge overlapping intervals" | Pattern 1: Merge | Sort by start |
| "Insert into non-overlapping intervals" | Pattern 2: Insert | Already sorted (don't re-sort) |
| "How many rooms / resources needed" | Pattern 3: Meeting rooms (heap) or Pattern 5: Sweep line | Sort by start (heap) or by event time (sweep) |
| "Maximum non-overlapping" / "minimum removals" | Pattern 4: Greedy | Sort by **end** |
| "Max overlap at any point" / "most concurrent" | Pattern 5: Sweep line | Sort by event time |
| "Intersection of two interval lists" | Pattern 6: Two pointers | Both lists pre-sorted |

**Heap vs. Sweep Line for "max overlap" problems:**

| Criterion | Heap (Pattern 3) | Sweep Line (Pattern 5) |
|---|---|---|
| Conceptual | "Track active meetings" | "Count starts and ends" |
| Code length | Longer (need heap boilerplate) | Shorter (just sort events) |
| Handles follow-ups | Better for "which meetings are in each room" | Better for "count at each timestamp" |
| When to choose | If you need to track *which* intervals are active | If you only need a count |

In a pure "how many rooms" interview question, sweep line is usually faster to code. But know both — interviewers sometimes ask for the heap version as a follow-up or vice versa.

---

## Common Interview Traps

### 1. Touching Intervals: Overlapping or Not?

```
[1,2] and [2,3] — do these overlap?
```

It depends on the problem:
- **Merge Intervals (LC 56):** Yes, merge them → `[1,3]`. Condition: `intervals[i][0] <= last[1]`.
- **Non-Overlapping Intervals (LC 435):** Typically no — `[1,2]` and `[2,3]` can coexist. Condition: `intervals[i][0] < prevEnd`.
- **Meeting Rooms:** Depends on problem wording. "A meeting ending at 2 and one starting at 2 don't need separate rooms" → `<=` in the pop condition.

**Rule:** Don't assume. Check the problem's definition of overlap, and encode it in your comparison operator (`<` vs `<=`).

### 2. Merge Intervals: Wrong Comparison

```go
// WRONG — comparing starts
if intervals[i][0] <= last[0] { ... }

// RIGHT — comparing current start to last merged END
if intervals[i][0] <= last[1] { ... }
```

After sorting by start, all remaining intervals have start >= the current one. The question is whether the current interval's start falls within the last merged interval's range, which is determined by the last merged interval's **end**.

### 3. Meeting Rooms: Forgetting to Track Maximum

```go
// WRONG — returns heap size at the end (which may be smaller than the peak)
return h.Len()

// RIGHT — track the maximum heap size seen during the sweep
if h.Len() > maxRooms {
    maxRooms = h.Len()
}
return maxRooms
```

### 4. Sweep Line: Wrong Tie-Breaking

```
Meeting A: [1, 5]
Meeting B: [5, 10]

Events: (1, +1), (5, -1), (5, +1), (10, -1)
```

If you process the `+1` at time 5 before the `-1`, you'll count 2 rooms at time 5. If you process `-1` first, you'll count 1 room. The correct answer depends on whether room reuse at the exact boundary is allowed.

**Default assumption for meeting rooms:** room is freed at end time, so `-1` before `+1` → 1 room suffices.

### 5. Insert Interval: Sorting Pre-Sorted Input

```go
// WRONG — wastes O(n log n) and signals you missed the constraint
sort.Slice(intervals, ...)

// RIGHT — the input is already sorted, just scan through
for i < n && intervals[i][1] < newInterval[0] { ... }
```

### 6. Interval Intersection: Wrong Pointer Advance

```go
// WRONG — advancing the pointer with the earlier start
if firstList[i][0] < secondList[j][0] {
    i++
}

// RIGHT — advancing the pointer with the earlier END
if firstList[i][1] < secondList[j][1] {
    i++
}
```

The interval that ends first can't overlap with anything else in the other list. The interval that ends later might still overlap with the next interval from the first list.

---

## Thought Process Walkthrough

### Walkthrough 1: Merge Intervals (LC 56)

> Given `intervals = [[1,3],[2,6],[8,10],[15,18]]`, merge all overlapping intervals.

**Step 1 — Recognize the pattern.**
"Merge overlapping intervals" → Pattern 1. Sort by start, merge adjacent.

**Step 2 — State the approach.**
"Sort by start time. Initialize merged list with the first interval. For each subsequent interval, compare its start to the last merged interval's end. If start <= end, extend. Otherwise, new interval."

**Step 3 — Trace through the example.**

```
Sorted: [[1,3],[2,6],[8,10],[15,18]]   (already sorted)

merged = [[1,3]]

i=1: [2,6] → start 2 <= last end 3 → overlap → extend end to max(3,6) = 6
     merged = [[1,6]]

i=2: [8,10] → start 8 > last end 6 → no overlap → new interval
     merged = [[1,6],[8,10]]

i=3: [15,18] → start 15 > last end 10 → no overlap → new interval
     merged = [[1,6],[8,10],[15,18]]

Result: [[1,6],[8,10],[15,18]] ✓
```

**Step 4 — Test an edge case: fully contained interval.**

```
Input: [[1,10],[2,5],[3,7]]
Sorted: [[1,10],[2,5],[3,7]]

merged = [[1,10]]

i=1: [2,5] → start 2 <= end 10 → extend end to max(10,5) = 10 → unchanged
i=2: [3,7] → start 3 <= end 10 → extend end to max(10,7) = 10 → unchanged

Result: [[1,10]] ✓  (both smaller intervals were absorbed)
```

**Step 5 — Code it.** (See Pattern 1 template above.)

**Step 6 — Complexity.** O(n log n) time for sorting, O(n) space for the result.

---

### Walkthrough 2: Meeting Rooms II (LC 253) — Heap Approach

> Given `intervals = [[0,30],[5,10],[15,20]]`, find the minimum number of meeting rooms required.

**Step 1 — Recognize the pattern.**
"Minimum meeting rooms" → Pattern 3 (heap) or Pattern 5 (sweep line). Let's use the heap approach — it's the most commonly asked version.

**Step 2 — State the approach.**
"Sort by start time. Use a min-heap of end times representing active meetings. For each new meeting, if its start >= the earliest end in the heap, pop that room (meeting ended, room freed). Push the current meeting's end time. Track the maximum heap size."

**Step 3 — Justify the approach to the interviewer.**
"The min-heap always holds the end times of currently active meetings. The root is the earliest ending meeting. If the new meeting starts after (or exactly when) that meeting ends, we reuse the room. Otherwise, we need a new room. The maximum simultaneous heap size is the answer."

**Step 4 — Trace through the example.**

```
Sorted by start: [[0,30],[5,10],[15,20]]

i=0: [0,30]
     heap is empty, push end=30
     heap = [30], maxRooms = 1

i=1: [5,10]
     peek: 30. start 5 < 30 → can't reuse → push end=10
     heap = [10, 30], maxRooms = 2

i=2: [15,20]
     peek: 10. start 15 >= 10 → reuse! pop 10, push end=20
     heap = [20, 30], maxRooms = 2

Answer: 2 ✓
```

**Step 5 — Now trace using sweep line (for comparison).**

```
Events: (0,+1), (30,-1), (5,+1), (10,-1), (15,+1), (20,-1)

Sorted (ends before starts on tie):
  (0,+1), (5,+1), (10,-1), (15,+1), (20,-1), (30,-1)

Sweep:
  t=0:  active = 0+1 = 1, max = 1
  t=5:  active = 1+1 = 2, max = 2
  t=10: active = 2-1 = 1
  t=15: active = 1+1 = 2, max = 2
  t=20: active = 2-1 = 1
  t=30: active = 1-1 = 0

Answer: 2 ✓ (same answer, simpler code)
```

**Step 6 — Code the heap version.**

```go
func minMeetingRooms(intervals [][]int) int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })

    h := &MinHeap{}
    heap.Init(h)
    maxRooms := 0

    for _, iv := range intervals {
        for h.Len() > 0 && (*h)[0] <= iv[0] {
            heap.Pop(h)
        }
        heap.Push(h, iv[1])
        if h.Len() > maxRooms {
            maxRooms = h.Len()
        }
    }

    return maxRooms
}
```

**Step 7 — Complexity.** O(n log n) time. Sorting is O(n log n). Each interval is pushed and popped from the heap at most once: O(n log n) total heap operations. Space O(n) for the heap.

**Step 8 — Follow-up readiness.**
- "What if I want to know which meetings go in which room?" → Assign room IDs: when you pop an end time, record which room was freed. Use a struct `{endTime, roomID}` in the heap.
- "Can you do it in O(n) if intervals have bounded timestamps?" → Yes, use a difference array instead of a map.

---

## Time Targets

| Problem | LC # | Target | Notes |
|---|---|---|---|
| Merge Intervals | 56 | 8 min | The foundation — must be fast |
| Insert Interval | 57 | 8 min | Three-phase scan, no sort needed |
| Meeting Rooms | 252 | 5 min | Just sort + check adjacent |
| Meeting Rooms II | 253 | 10 min | Heap or sweep line — know both |
| Non-Overlapping Intervals | 435 | 8 min | Sort by end, greedy (Day 14 recap) |
| Interval List Intersections | 986 | 8 min | Two pointers, advance by end |
| Minimum Interval to Include Query | 1851 | 15 min | Sweep line + heap — harder |

---

## Quick Drill (30 minutes)

Do these without looking at templates. Time yourself.

1. **Merge Intervals** (LC 56) — Sort by start, merge in one pass. Test with `[[1,4],[0,4]]` (unsorted input) and `[[1,10],[2,3],[4,5]]` (fully contained). Target: 8 minutes.

2. **Insert Interval** (LC 57) — Three-phase scan. Test with inserting `[5,7]` into `[[1,3],[6,9]]`. Target: 8 minutes.

3. **Meeting Rooms II** (LC 253) — Write BOTH the heap version and the sweep line version. Compare the two. Target: 12 minutes total (6 each).

4. **Interval List Intersections** (LC 986) — Two pointers. Test with `[[0,2],[5,10]]` and `[[1,5],[8,12]]`. Target: 8 minutes.

After each one, check:
- Did I use the correct sort order (start vs. end)?
- Did I handle touching intervals correctly for *this specific problem*?
- Could I articulate the approach in one sentence before coding?

---

## Self-Assessment

### Can I explain these from memory?

| Question | Confident? |
|---|---|
| Merge intervals: why sort by start and compare current start to last END? | |
| Insert interval: why is it O(n) and not O(n log n)? | |
| Meeting rooms heap: why a min-heap of end times and not start times? | |
| Meeting rooms: what does the heap size represent at any given moment? | |
| Sweep line: how do I decide whether to process ends or starts first on a tie? | |
| Interval intersection: why advance the pointer with the earlier END, not START? | |
| Non-overlapping intervals: why sort by END time? (counterexample?) | |
| When to use heap vs. sweep line for "max overlap" problems? | |

### Red flags that you need more practice:
- You sort by start time for the greedy "maximum non-overlapping" problem.
- You compare start-to-start instead of start-to-end in merge intervals.
- You can't write the sweep line tie-breaking rule without looking it up.
- Meeting Rooms II takes you more than 10 minutes (including heap boilerplate).
- You re-sort already-sorted input in Insert Interval.

### Green lights — you're ready:
- You can state the sort order and key comparison for each pattern without hesitation.
- You know when touching intervals overlap vs. don't, and you encode it correctly in `<` vs `<=`.
- You can solve Meeting Rooms II with both heap and sweep line, and explain the trade-offs.
- You hit the time targets above consistently.
- Given a new interval problem, you can classify it into one of the six patterns within 30 seconds.
