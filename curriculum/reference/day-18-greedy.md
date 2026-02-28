# Day 18 — Greedy Algorithms

## 1. Curated Learning Resources

| # | Resource | Why It's Useful |
|---|----------|-----------------|
| 1 | [Greedy Algorithms — Abdul Bari (YouTube)](https://www.youtube.com/watch?v=ARvQcqJ_-NY) | Clear, slow-paced walkthrough of greedy choice property with activity selection as the running example. Good warm-up before the session. |
| 2 | [Exchange Argument Proof Technique — Jeff Erickson's Algorithms (Chapter 4)](https://jeffe.cs.illinois.edu/teaching/algorithms/book/04-greedy.pdf) | The gold-standard reference for how to formally prove a greedy algorithm is correct. Read sections 4.1-4.3. |
| 3 | [Kadane's Algorithm Visualized — NeetCode](https://www.youtube.com/watch?v=5WZl3MMT0Eg) | Step-by-step visual of the running sum and where it resets. Makes the "why" behind the reset intuitive. |
| 4 | [Jump Game I & II — NeetCode](https://www.youtube.com/watch?v=Yan0cv2cLy8) | Shows the "farthest reach" approach and the BFS-level interpretation of Jump Game II side by side. |
| 5 | [When to Use Greedy vs DP — Algorithms with Attitude](https://www.youtube.com/watch?v=lFYPPMxmKWE) | Directly addresses the greedy-vs-DP decision, with examples where greedy intuition fails. |
| 6 | [Task Scheduler (LeetCode 621) Explanation — Larry's LeetCode](https://www.youtube.com/watch?v=ySTQCRya6B0) | Breaks down the formula derivation and the slot-filling grid visualization. |
| 7 | [Gas Station Problem Walkthrough — Back to Back SWE](https://www.youtube.com/watch?v=lJwbPZGo05A) | Focuses on the circular-route invariant and why one pass suffices. |
| 8 | [cp-algorithms: Greedy Algorithms](https://cp-algorithms.com/others/greedy.html) | Concise reference covering exchange arguments and common greedy patterns. Good for quick review after the session. |

---

## 2. Detailed 2-Hour Session Plan

### 12:00–12:20 — Review & Conceptualize (20 min)

| Time | Activity |
|------|----------|
| 12:00–12:07 | Read through [Section 3: Core Concepts](#3-core-concepts-deep-dive) below. Focus on the greedy choice property, optimal substructure, and the exchange argument proof. Don't code yet — just internalize when and why greedy works. |
| 12:07–12:12 | Study the [Greedy vs DP Decision Framework](#5-greedy-vs-dp-decision-framework). For each example, ask yourself: "Could I construct a counterexample where the greedy choice fails?" |
| 12:12–12:17 | Review the [ASCII diagrams](#6-visual-diagrams) for Kadane's and Jump Game. Trace through the examples by hand on paper. |
| 12:17–12:20 | Read all five function signatures in the [Implementation Checklist](#4-implementation-checklist). Mentally sketch the approach for each before writing any code. |

### 12:20–1:20 — Implement (60 min)

| Time | Problem | Key Idea | Target |
|------|---------|----------|--------|
| 12:20–12:32 | **MaxSubArray (Kadane's)** | Track `currentSum`. If `currentSum < 0`, reset to `nums[i]`. Otherwise `currentSum += nums[i]`. Update `maxSum` at each step. | Get working with all-negative edge case. This is the cleanest greedy — one variable, one decision. |
| 12:32–12:44 | **CanJump (Jump Game)** | Scan left to right, maintain `farthest = max(farthest, i + nums[i])`. If `i > farthest`, return false. | Builds the "farthest reach" invariant. Simple boolean outcome. |
| 12:44–12:58 | **MinJumps (Jump Game II)** | BFS interpretation: maintain `curEnd` (end of current level) and `farthest` (end of next level). When `i` reaches `curEnd`, increment jumps and set `curEnd = farthest`. | Hardest of the five. Take time to understand why this is BFS on a virtual graph. |
| 12:58–1:10 | **LeastInterval (Task Scheduler)** | Count frequencies. The most frequent task forces `(maxFreq - 1)` gaps of size `n + 1`. Fill idle slots with remaining tasks. Answer is `max(len(tasks), (maxFreq-1)*(n+1) + countOfMax)`. | The formula is the greedy insight. Make sure you understand WHY the most frequent task is the bottleneck. |
| 1:10–1:20 | **CanCompleteCircuit (Gas Station)** | If `totalGas >= totalCost`, a solution exists. Track `tank` (running surplus). Whenever `tank < 0`, reset start to `i+1` and reset `tank`. | The reset logic is the greedy choice: any station that causes a deficit can't be the start. |

### 1:20–1:50 — Solidify (30 min)

| Time | Activity |
|------|----------|
| 1:20–1:30 | **Edge cases**: Run all-negative arrays through Kadane's. Test `[0,0,0]` in CanJump. Test single-element arrays. Test a task list with one unique task (e.g., `['A','A','A']` with `n=2`). Test a gas station where `totalGas < totalCost`. |
| 1:30–1:40 | **Variants**: (1) Modify Kadane's to return the actual subarray indices, not just the sum. (2) For Jump Game II, track the actual jump path. (3) For Gas Station, verify by simulating the full circuit from your returned start index. |
| 1:40–1:50 | **Greedy proofs**: For each of the five problems, write one sentence explaining why the greedy choice is safe. If you can't, that's the signal to dig deeper into the exchange argument. |

### 1:50–2:00 — Recap (10 min)

| Time | Activity |
|------|----------|
| 1:50–1:55 | From memory, write down the complexity of each problem (all are O(n) time, O(1) space). Write the one-line greedy insight for each. |
| 1:55–2:00 | Answer at least 2 of the [Self-Assessment](#7-self-assessment) questions without looking at notes. If you can't, flag the topic for tomorrow. |

---

## 3. Core Concepts Deep Dive

### What Makes a Problem Greedy

A problem is solvable by greedy when it has **both**:

1. **Greedy choice property**: There exists a locally optimal choice that is part of some globally optimal solution. You don't need to solve subproblems first to make the choice — you can commit now.

2. **Optimal substructure**: After making the greedy choice, the remaining problem is a smaller instance of the same type, and its optimal solution combines with the greedy choice to form the global optimum.

DP also requires optimal substructure, but DP explores all choices (building a table). Greedy commits to one choice and never revisits.

### The Exchange Argument

The exchange argument is the standard technique for proving greedy correctness. The structure:

1. **Assume** an optimal solution `OPT` that differs from the greedy solution `G`.
2. **Find** the first point where they differ — `OPT` made choice X, `G` made choice Y.
3. **Show** that swapping X for Y in `OPT` produces a solution that is at least as good. (This is the key step.)
4. **Conclude** by induction that `G` is optimal (you can transform any `OPT` into `G` without losing quality).

**Worked example: Activity Selection (Interval Scheduling)**

> **Problem**: Given activities with start/end times, select the maximum number of non-overlapping activities.
>
> **Greedy strategy**: Always pick the activity that finishes earliest.
>
> **Proof by exchange argument**:
>
> Let `OPT = {a₁, a₂, ..., aₖ}` be an optimal solution sorted by finish time, and let `G = {g₁, g₂, ..., gₘ}` be the greedy solution sorted by finish time.
>
> We want to show `m ≥ k` (greedy picks at least as many activities as any optimal solution).
>
> **Step 1**: Consider the first activity. Greedy picks `g₁`, which finishes earliest among all activities. If `OPT` also picks `g₁`, great — they agree. If `OPT` picks `a₁ ≠ g₁`, then `finish(g₁) ≤ finish(a₁)` (because greedy picks the earliest finisher).
>
> **Step 2 (Exchange)**: Replace `a₁` with `g₁` in `OPT`. Since `g₁` finishes no later than `a₁`, it can't create any new conflicts with `a₂, a₃, ...`. So the modified `OPT` is still valid and has the same size.
>
> **Step 3 (Induction)**: Now `OPT` and `G` agree on the first choice. Repeat the argument for the remaining activities. At each step, we can swap without losing any activities.
>
> **Conclusion**: `|G| ≥ |OPT|`. Since `OPT` is optimal, `|G| = |OPT|`. Greedy is optimal. ∎

### Greedy vs DP: When Greedy Fails

The fundamental question: **"Does the locally optimal choice ever lead to a globally suboptimal solution?"**

**Greedy works — Coin Change with US denominations `[25, 10, 5, 1]`**:
- For amount 30: Greedy picks 25, then 5. Total = 2 coins. This is optimal.
- Why? Each larger denomination is a multiple of smaller ones (or close enough). The greedy choice of "use the biggest coin that fits" never overshoots.

**Greedy fails — Coin Change with arbitrary denominations `[1, 3, 4]`**:
- For amount 6: Greedy picks 4, then 1, then 1. Total = 3 coins.
- Optimal: 3 + 3 = 2 coins.
- Why? 4 is not a "dominant" denomination — using it blocks a better combination. The locally optimal choice (pick 4) leads to a globally suboptimal result. You need DP to explore both possibilities.

### Kadane's Algorithm

**The insight**: At each index `i`, decide whether to extend the current subarray or start fresh.

```
currentSum = max(nums[i], currentSum + nums[i])
```

**Why resetting when `currentSum < 0` is correct**: If the sum of all elements up to index `i` is negative, then *any* subarray starting at or before index `i` and extending past it would be improved by dropping the prefix up to `i`. A negative prefix can never contribute to a maximum sum — it can only drag it down. So we reset and start a new subarray at `i`.

**All-negative arrays**: Kadane's handles this correctly because at each step, `currentSum = max(nums[i], currentSum + nums[i])`. When all values are negative, `nums[i] > currentSum + nums[i]` always (adding a negative to a negative makes it worse), so `currentSum` keeps resetting to each individual element. `maxSum` tracks the largest (least negative) of these.

### Jump Game: The "Farthest Reach" Invariant

**Invariant**: As you scan from left to right, maintain `farthest` — the maximum index you can reach from any index you've visited so far.

```
farthest = max(farthest, i + nums[i])   for each i where i <= farthest
```

If at any point `i > farthest`, you're stuck. If `farthest >= n-1`, you can reach the end.

**Why greedy works**: You never need to consider which specific path gets you to the end — only whether the end is reachable. The farthest reach is monotonically non-decreasing, and it captures all reachable indices in one pass.

### Jump Game II: The BFS Interpretation

Think of the array as a graph. From index `i`, you can jump to any index in `[i+1, i+nums[i]]`. Finding the minimum jumps is finding the shortest path — which is BFS.

**Each "level" of BFS** is the range of indices reachable with one more jump:
- Level 0: `{0}` (start)
- Level 1: all indices reachable from index 0 → `{1, ..., nums[0]}`
- Level 2: all indices reachable from level 1 that aren't in level 0 or 1
- ...

The greedy implementation tracks `curEnd` (right boundary of current level) and `farthest` (right boundary of next level). When `i` passes `curEnd`, you've moved to the next level (one more jump).

### Task Scheduler: The Bottleneck Formula

The most frequent task creates the bottleneck because it requires the most cooldown gaps.

**Derivation**: Let `maxFreq` be the highest frequency among all tasks, and `countOfMax` be how many tasks share that frequency.

Imagine laying out the most frequent task with `n` idle slots between each occurrence:

```
A _ _ A _ _ A
```

This creates `(maxFreq - 1)` gaps, each of width `n`. Total slots for the frame: `(maxFreq - 1) * (n + 1)`. Then add `countOfMax` for the final group (the last occurrence of each max-frequency task).

```
result = (maxFreq - 1) * (n + 1) + countOfMax
```

But if there are enough diverse tasks to fill all idle slots, no idle time is needed, and the answer is just `len(tasks)`.

**Final answer**: `max(len(tasks), (maxFreq - 1) * (n + 1) + countOfMax)`

### Gas Station: The Running Surplus

**Two key observations**:
1. If `sum(gas) >= sum(cost)`, a valid starting station exists (and is unique).
2. If starting from station `s` causes the tank to go negative at station `k`, then no station in `[s, k]` can be the valid start. (Because your tank was non-negative at `s` — starting from any later station gives you strictly less fuel at station `k`.)

So: scan once, track `tank`. Whenever `tank < 0`, reset `start = i + 1` and `tank = 0`. The last reset gives the answer.

---

## 4. Implementation Checklist

### Function Signatures

```go
package greedy

// MaxSubArray returns the largest sum of any contiguous subarray.
// Kadane's algorithm: O(n) time, O(1) space.
func MaxSubArray(nums []int) int { ... }

// CanJump returns true if you can reach the last index.
// Farthest-reach greedy: O(n) time, O(1) space.
func CanJump(nums []int) bool { ... }

// MinJumps returns the minimum number of jumps to reach the last index.
// BFS-style greedy: O(n) time, O(1) space.
// Assumes the last index is always reachable.
func MinJumps(nums []int) int { ... }

// LeastInterval returns the minimum number of intervals (including idle)
// needed to schedule all tasks with cooldown n between same tasks.
// O(n) time, O(1) space (26-letter alphabet).
func LeastInterval(tasks []byte, n int) int { ... }

// CanCompleteCircuit returns the starting gas station index for a
// circular route, or -1 if impossible. O(n) time, O(1) space.
func CanCompleteCircuit(gas, cost []int) int { ... }
```

### Tests & Edge Cases

```go
package greedy

import "testing"

func TestMaxSubArray(t *testing.T) {
    tests := []struct {
        name string
        nums []int
        want int
    }{
        {"mixed", []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
        {"all negative", []int{-3, -2, -5, -1, -4}, -1},
        {"single element positive", []int{5}, 5},
        {"single element negative", []int{-3}, -3},
        {"all positive", []int{1, 2, 3, 4}, 10},
        {"reset in middle", []int{2, -8, 3, -1, 4}, 6},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := MaxSubArray(tt.nums); got != tt.want {
                t.Errorf("MaxSubArray(%v) = %d, want %d", tt.nums, got, tt.want)
            }
        })
    }
}

func TestCanJump(t *testing.T) {
    tests := []struct {
        name string
        nums []int
        want bool
    }{
        {"reachable", []int{2, 3, 1, 1, 4}, true},
        {"stuck at zero", []int{3, 2, 1, 0, 4}, false},
        {"single element", []int{0}, true},
        {"all zeros except first", []int{5, 0, 0, 0, 0, 0}, true},
        {"all zeros", []int{0, 0, 0}, false},
        {"large first jump", []int{10, 0, 0, 0, 0}, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := CanJump(tt.nums); got != tt.want {
                t.Errorf("CanJump(%v) = %v, want %v", tt.nums, got, tt.want)
            }
        })
    }
}

func TestMinJumps(t *testing.T) {
    tests := []struct {
        name string
        nums []int
        want int
    }{
        {"two jumps", []int{2, 3, 1, 1, 4}, 2},
        {"one jump", []int{2, 1}, 1},
        {"already there", []int{0}, 0},
        {"three jumps", []int{1, 1, 1, 1}, 3},
        {"optimal choice", []int{2, 3, 0, 1, 4}, 2},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := MinJumps(tt.nums); got != tt.want {
                t.Errorf("MinJumps(%v) = %d, want %d", tt.nums, got, tt.want)
            }
        })
    }
}

func TestLeastInterval(t *testing.T) {
    tests := []struct {
        name  string
        tasks []byte
        n     int
        want  int
    }{
        {"standard", []byte{'A', 'A', 'A', 'B', 'B', 'B'}, 2, 8},
        {"no cooldown", []byte{'A', 'A', 'A', 'B', 'B', 'B'}, 0, 6},
        {"single task type", []byte{'A', 'A', 'A'}, 2, 7},
        {"enough variety", []byte{'A', 'A', 'A', 'B', 'B', 'B', 'C', 'C', 'D'}, 2, 9},
        {"one task", []byte{'A'}, 5, 1},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := LeastInterval(tt.tasks, tt.n); got != tt.want {
                t.Errorf("LeastInterval(%v, %d) = %d, want %d", tt.tasks, tt.n, got, tt.want)
            }
        })
    }
}

func TestCanCompleteCircuit(t *testing.T) {
    tests := []struct {
        name string
        gas  []int
        cost []int
        want int
    }{
        {"standard", []int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}, 3},
        {"impossible", []int{2, 3, 4}, []int{3, 4, 3}, -1},
        {"single station", []int{5}, []int{4}, 0},
        {"start at zero", []int{3, 1, 1}, []int{1, 2, 2}, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := CanCompleteCircuit(tt.gas, tt.cost); got != tt.want {
                t.Errorf("CanCompleteCircuit(%v, %v) = %d, want %d",
                    tt.gas, tt.cost, got, tt.want)
            }
        })
    }
}
```

---

## 5. Greedy vs DP Decision Framework

### The Decision Flowchart

```
1. Can you construct a counterexample where the greedy choice fails?
   YES → Use DP.
   NO  → Greedy might work. Continue checking.

2. Does sorting + a single pass solve it?
   YES → Likely greedy. (Interval scheduling, activity selection.)

3. Does the problem have "maximum non-overlapping" or "minimum number of" structure?
   YES → Likely greedy. (Interval partitioning, jump game.)

4. Does the problem have "number of ways" or "all possible" structure?
   YES → Likely DP. (Climbing stairs, subset sum, word break.)

5. Is there a clear exchange argument — can you prove swapping any non-greedy
   choice for the greedy choice never hurts?
   YES → Greedy is correct.
   NO / UNSURE → Fall back to DP.
```

### Examples Where Greedy Intuition Is Wrong

**1. Coin Change with arbitrary denominations**
- Coins: `[1, 3, 4]`, Amount: `6`
- Greedy (largest first): 4 + 1 + 1 = 3 coins
- Optimal (DP): 3 + 3 = 2 coins
- **Why greedy fails**: Picking the largest coin doesn't account for how remaining denominations combine. No exchange argument works — swapping 4 for 3 does improve the solution.

**2. Longest Increasing Subsequence (LIS)**
- Array: `[3, 1, 8, 2, 5]`
- Greedy (always extend with next larger): picks `3, 8` → length 2
- Optimal (DP): `1, 2, 5` → length 3
- **Why greedy fails**: Committing to a larger value early closes off longer subsequences. The choice at each step depends on the entire future, not just the next element.

**3. 0/1 Knapsack**
- Items: `[(weight=10, value=60), (weight=20, value=100), (weight=30, value=120)]`, Capacity: `50`
- Greedy (best value/weight ratio first): item 1 (ratio 6) + item 2 (ratio 5) = value 160, weight 30
- Optimal (DP): item 2 + item 3 = value 220, weight 50
- **Why greedy fails**: The fractional knapsack is greedy, but 0/1 knapsack is not. You can't take partial items, so the ratio heuristic doesn't account for how items combine to fill capacity.

**4. Word Break**
- String: `"catsanddog"`, Dict: `["cats", "cat", "sand", "and", "dog"]`
- Greedy (longest match first): `"cats" + "anddog"` — stuck, `"anddog"` not in dict
- Optimal (DP): `"cat" + "sand" + "dog"` — valid
- **Why greedy fails**: Taking the longest prefix match can consume characters needed for later valid splits. Multiple decompositions must be explored.

### Quick Reference Table

| Signal | Greedy | DP |
|--------|--------|----|
| "Minimum/maximum with one-pass" | ✓ | |
| "Number of ways" | | ✓ |
| "All possible results" | | ✓ |
| Sort + scan suffices | ✓ | |
| Current choice depends on future | | ✓ |
| Clear exchange argument exists | ✓ | |
| Overlapping subproblems, many choices | | ✓ |

---

## 6. Visual Diagrams

### Kadane's Algorithm — Running Sum with Reset

```
Array:     [-2,  1, -3,  4, -1,  2,  1, -5,  4]

Step-by-step:
Index:       0    1    2    3    4    5    6    7    8
Value:      -2    1   -3    4   -1    2    1   -5    4
curSum:     -2    1   -2    4    3    5    6    1    5
maxSum:     -2    1    1    4    4    5    6    6    6
             ↑         ↑
           reset     reset would make it -5,
           (start    but max(nums[2], curSum+nums[2])
            fresh)   = max(-3, -2) = -2, so we keep going

Detail at index 0:  curSum = max(-2, 0 + -2) = -2    maxSum = -2
Detail at index 1:  curSum = max( 1, -2 +  1) =  1    maxSum =  1  ← reset! (1 > -1)
Detail at index 3:  curSum = max( 4, -2 +  4) =  4    maxSum =  4  ← reset! (4 > 2)
                     ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
                     The subarray [4, -1, 2, 1] gives the max sum of 6.
                     It starts fresh at index 3 because the prefix sum was negative.
```

### All-Negative Example

```
Array:     [-3, -2, -5, -1, -4]

Index:       0    1    2    3    4
Value:      -3   -2   -5   -1   -4
curSum:     -3   -2   -5   -1   -4
maxSum:     -3   -2   -2   -1   -1
             ↑    ↑    ↑    ↑    ↑
           Each step resets because extending makes it worse.
           maxSum correctly tracks the least-negative value: -1
```

### Jump Game — Farthest Reach Expanding

```
Array:      [2,  3,  1,  1,  4]
Index:       0   1   2   3   4

Scan left to right, track farthest reachable index:

  i=0: farthest = max(0, 0+2) = 2       can reach up to index 2
       [===========]
        0   1   2   3   4

  i=1: farthest = max(2, 1+3) = 4       can reach up to index 4 ✓ DONE
       [=====================]
        0   1   2   3   4

  Result: true (farthest >= 4)


Failing example: [3,  2,  1,  0,  4]

  i=0: farthest = max(0, 0+3) = 3
       [================]
        0   1   2   3   4

  i=1: farthest = max(3, 1+2) = 3       no improvement
  i=2: farthest = max(3, 2+1) = 3       no improvement
  i=3: farthest = max(3, 3+0) = 3       stuck! farthest=3, need 4
       [================]
        0   1   2   3   X

  Result: false (farthest < 4)
```

### Jump Game II — BFS Levels on the Array

```
Array:      [2,  3,  1,  1,  4]
Index:       0   1   2   3   4

Level 0 (0 jumps):  indices {0}
  From 0: can reach 1..2
  curEnd=0, farthest=2

  |  0  |  1     2  |  3     4  |
  |L0   |           |           |
  |start|           |           |

Level 1 (1 jump):   indices {1, 2}
  From 1: can reach 2..4
  From 2: can reach 3..3
  curEnd=2, farthest=4

  |  0  |  1     2  |  3     4  |
  |L0   |  L1   L1  |           |
  |     |  →→→→→→→→→→→→→→→→→→→ |  (farthest = 4, covers end!)

Level 2 (2 jumps):  indices {3, 4} — includes target!

  |  0  |  1     2  |  3     4  |
  |L0   |  L1   L1  |  L2   L2  |

  Answer: 2 jumps

Code trace:
  i=0: farthest = max(0, 0+2) = 2
       i == curEnd(0) → jumps=1, curEnd=2
  i=1: farthest = max(2, 1+3) = 4
  i=2: i == curEnd(2) → jumps=2, curEnd=4
       curEnd(4) >= n-1(4) → return 2
```

### Task Scheduler — Slot-Filling Grid

```
Tasks: [A, A, A, B, B, B],  n = 2

Frequencies:  A:3  B:3    maxFreq=3, countOfMax=2

Step 1: Lay out the most frequent task with n-wide gaps:

  | A | _ | _ | A | _ | _ | A |
    ^       ^   ^       ^   ^
    slot 1      slot 2      (no gap after last)

  Frame size = (maxFreq - 1) * (n + 1) + countOfMax
             = (3 - 1)      * (2 + 1)  + 2
             = 6 + 2 = 8

Step 2: Fill idle slots with remaining tasks:

  | A | B | _ | A | B | _ | A | B |
                ^           ^
              idle         idle... wait, B also has maxFreq!

  Actually with countOfMax = 2:
  | A | B | idle | A | B | idle | A | B |
    1   2    3     4   5    6     7   8

  Total intervals = 8


Example with enough variety to fill all idle slots:
Tasks: [A, A, A, B, B, B, C, C, D],  n = 2

  | A | B | C | A | B | C | A | B | D |
    1   2   3   4   5   6   7   8   9

  Formula gives: (3-1)*(2+1) + 2 = 8
  But len(tasks) = 9 > 8
  Answer: max(9, 8) = 9  (no idle slots needed!)
```

---

## 7. Self-Assessment

Answer these without looking at your notes. If you can't, revisit the corresponding section.

**1. Why does Kadane's algorithm handle all-negative arrays correctly?**

> *Expected*: Because at each step, `currentSum = max(nums[i], currentSum + nums[i])`. When all values are negative, `nums[i]` is always greater than `currentSum + nums[i]` (adding a negative to a negative makes it worse), so `currentSum` resets to each element individually. `maxSum` tracks the largest (least negative) single element, which is the correct answer.

**2. Can you explain why greedy works for Jump Game but not for Coin Change with arbitrary denominations?**

> *Expected*: Jump Game has the greedy choice property — extending the farthest reach is always safe because reaching further can never hurt (it only opens up more options). There's no way that "not jumping as far as possible" gives you access to indices that jumping far doesn't. For Coin Change, picking the largest denomination that fits can block better combinations (e.g., `[1,3,4]` for amount 6: greedy picks 4+1+1=3 coins, but 3+3=2 coins is better). The local choice constrains the remaining problem in a way that may be globally suboptimal.

**3. In Jump Game II, what does each "BFS level" represent, and why does this give the minimum number of jumps?**

> *Expected*: Each BFS level represents the set of indices reachable with exactly `k` jumps. Level 0 is `{0}`, level 1 is all indices reachable from level 0, and so on. BFS explores all positions at distance `k` before any at distance `k+1`, so the first time we reach the last index, we've found the minimum number of jumps. The greedy implementation tracks this with `curEnd` (right boundary of current level) and `farthest` (right boundary of next level).

**4. Why does the Task Scheduler formula use `max(len(tasks), formula)`? When does each case dominate?**

> *Expected*: The formula `(maxFreq - 1) * (n + 1) + countOfMax` computes the minimum time assuming there ARE idle slots (the most frequent task forces gaps). But if there are enough distinct tasks to fill all gaps, no idle time exists and the answer is simply `len(tasks)`. The formula dominates when task variety is low (lots of idle time). `len(tasks)` dominates when variety is high (tasks fill all gaps).

**5. For Gas Station, why is it safe to reset the starting station to `i+1` when the running tank goes negative at station `i`?**

> *Expected*: If you started at station `s` and your tank went negative at station `i`, it means the segment `[s, i]` has a net fuel deficit. Starting from any station `j` in `[s+1, i]` is even worse, because you'd have the same deficit from `[j, i]` minus the (non-negative) surplus you had accumulated from `[s, j-1]`. So no station in `[s, i]` can be the answer — skip them all and try `i+1`.
