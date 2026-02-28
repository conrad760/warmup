# Day 2: Two Pointers & Sliding Window

> **Time block:** 2 hours
> **Goal:** Lock in pattern recognition so you identify the right technique within 30 seconds of reading a problem.
> **Assumption:** You already know how two pointers and sliding window work. This is about *speed, traps, and interview fluency*.

---

## Pattern Catalog

### 1. Opposite Ends (Converging Pointers)

**Trigger:** "When you see *sorted array*, *pair with target sum*, *container/area*, *palindrome check*, or *two elements satisfying a condition*..."

**Template:**
```go
func twoPointer(nums []int, target int) (int, int) {
    left, right := 0, len(nums)-1
    for left < right {
        sum := nums[left] + nums[right]
        if sum == target {
            return left, right
        } else if sum < target {
            left++
        } else {
            right--
        }
    }
    return -1, -1
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- Input MUST be sorted (or the problem structure must give you a reason to move one pointer — like Container With Most Water, where you move the shorter side).
- If the problem asks for indices in the original unsorted array, you can't sort without losing them. Use a hash map instead.
- Use `left < right` (strict), not `left <= right`, unless you want to process a single middle element.

---

### 2. Same Direction (Fast-Slow / Write Pointer)

**Trigger:** "When you see *remove duplicates in-place*, *move zeros*, *partition array*, *remove element*, or *O(1) extra space compaction*..."

**Template:**
```go
func removeDuplicates(nums []int) int {
    if len(nums) == 0 { return 0 }
    w := 1 // write pointer — next position to write
    for r := 1; r < len(nums); r++ {
        if nums[r] != nums[w-1] { // condition: different from last written
            nums[w] = nums[r]
            w++
        }
    }
    return w // new length
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- The write pointer `w` always points to the *next available slot*. The valid portion of the array is `nums[0:w]`.
- For "remove duplicates allowing at most K," compare `nums[r]` with `nums[w-K]`.
- For "move zeros," swap `nums[w]` and `nums[r]` instead of overwriting — this preserves non-zero order AND puts zeros at the end in one pass.

```go
func moveZeroes(nums []int) {
    w := 0
    for r := 0; r < len(nums); r++ {
        if nums[r] != 0 {
            nums[w], nums[r] = nums[r], nums[w]
            w++
        }
    }
}
```

---

### 3. Floyd's Cycle Detection (Tortoise and Hare)

**Trigger:** "When you see *linked list cycle*, *find the start of a cycle*, *find duplicate in read-only array [1,n]*, or *O(1) space cycle detection*..."

**Template:**
```go
func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast { return true }
    }
    return false
}
```

**Finding the cycle start (phase 2):**
```go
func detectCycleStart(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            // Phase 2: move one pointer to head, advance both by 1
            slow = head
            for slow != fast {
                slow = slow.Next
                fast = fast.Next
            }
            return slow
        }
    }
    return nil
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- The `fast != nil && fast.Next != nil` check order matters — check `fast` before `fast.Next` to avoid nil dereference.
- For "Find the Duplicate Number" (LeetCode 287): treat the array as a linked list where `nums[i]` is the "next" pointer. Index 0 is the head. Values in range [1,n] means index 0 is never pointed to, so it's guaranteed to be outside the cycle — the cycle start is the duplicate value.
- The math proof: if the distance from head to cycle start is `a`, and the distance from cycle start to meeting point is `b`, then `a = c - b` where `c` is the cycle length. That's why phase 2 works.

---

### 4. Three-Pointer / Fix One + Two Pointers

**Trigger:** "When you see *3Sum*, *3Sum closest*, *triplets with sum*, or *three elements satisfying a condition on a sorted array*..."

**Template:**
```go
func threeSum(nums []int, target int) [][]int {
    sort.Ints(nums)
    var result [][]int
    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue } // skip dup at level 1
        lo, hi := i+1, len(nums)-1
        for lo < hi {
            sum := nums[i] + nums[lo] + nums[hi]
            if sum == target {
                result = append(result, []int{nums[i], nums[lo], nums[hi]})
                lo++; hi--
                for lo < hi && nums[lo] == nums[lo-1] { lo++ }  // skip dup level 2
                for lo < hi && nums[hi] == nums[hi+1] { hi-- }  // skip dup level 3
            } else if sum < target {
                lo++
            } else {
                hi--
            }
        }
    }
    return result
}
```

**Complexity:** O(n^2) time (sort is O(n log n), dominated by the nested loops), O(1) extra space (ignoring output).

**Watch out:**
- **You must skip duplicates at ALL THREE levels:** (1) the outer `i` loop, (2) incrementing `lo` after a match, (3) decrementing `hi` after a match. Missing any one produces duplicate triplets.
- The outer skip is `i > 0 && nums[i] == nums[i-1]` — NOT `nums[i] == nums[i+1]`. The latter skips valid first elements.
- For 3Sum Closest, you don't need duplicate skipping (you're tracking a single best), but you track `closestSum` and update when `abs(sum - target) < abs(closestSum - target)`.

---

### 5. Fixed-Size Sliding Window

**Trigger:** "When you see *subarray of size K*, *maximum sum of K consecutive elements*, *average of each window of size K*..."

**Template:**
```go
func maxSumSizeK(nums []int, k int) int {
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += nums[i]
    }
    maxSum := windowSum
    for r := k; r < len(nums); r++ {
        windowSum += nums[r] - nums[r-k] // slide: add right, remove left
        if windowSum > maxSum {
            maxSum = windowSum
        }
    }
    return maxSum
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- Window size is `right - left + 1`. When the window is `[r-k+1, r]`, the element leaving is `nums[r-k]`.
- Edge case: `k > len(nums)`. Guard against this.
- This only works when the window size is FIXED. If the problem says "at most K" or "minimum subarray," that's a variable window.

---

### 6. Variable Window — Find Longest

**Trigger:** "When you see *longest substring without repeating*, *longest with at most K distinct characters*, or any *maximize* a contiguous subarray/substring subject to a constraint..."

**Template:**
```go
func longestWithAtMostKDistinct(s string, k int) int {
    freq := make(map[byte]int)
    maxLen := 0
    left := 0
    for right := 0; right < len(s); right++ {
        freq[s[right]]++
        for len(freq) > k { // WHILE not "if" — shrink until valid
            freq[s[left]]--
            if freq[s[left]] == 0 { delete(freq, s[left]) }
            left++
        }
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}
```

**Complexity:** O(n) time (each element is added and removed at most once), O(K) space.

**Watch out:**
- Use `for` (while-loop) to shrink, NEVER `if`. The window may need to shrink multiple positions.
- Update the answer AFTER shrinking — you want `maxLen` of a valid window.
- For "longest substring without repeating characters," you can use a `map[byte]int` storing the last seen index and jump `left` directly instead of incrementing one by one. That's an optimization, not a requirement.
- `delete(freq, s[left])` when the count hits 0 — otherwise `len(freq)` stays inflated.

---

### 7. Variable Window — Find Shortest

**Trigger:** "When you see *minimum window*, *smallest subarray with sum >= S*, or any *minimize* a contiguous subarray/substring subject to a constraint..."

**Template:**
```go
func minSubarrayLen(target int, nums []int) int {
    minLen := math.MaxInt32
    sum := 0
    left := 0
    for right := 0; right < len(nums); right++ {
        sum += nums[right]
        for sum >= target { // WHILE valid, try to shrink
            if right-left+1 < minLen {
                minLen = right - left + 1
            }
            sum -= nums[left]
            left++
        }
    }
    if minLen == math.MaxInt32 { return 0 }
    return minLen
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- The loop direction is REVERSED from "find longest":
  - **Longest:** expand always, shrink when INVALID, answer = max of valid windows.
  - **Shortest:** expand always, shrink while VALID (record answer each time you're valid), answer = min.
- Again: `for` not `if` for the shrink loop.
- Return 0 (or -1) if no valid window was ever found. Don't forget this edge case.

---

### 8. Window + Frequency Map

**Trigger:** "When you see *permutation in string*, *find all anagrams*, *minimum window substring*, or any problem matching character frequencies across two strings..."

**Template:**
```go
func findAnagrams(s string, p string) []int {
    if len(p) > len(s) { return nil }
    need := [26]int{}
    window := [26]int{}
    for i := 0; i < len(p); i++ {
        need[p[i]-'a']++
    }
    var result []int
    for r := 0; r < len(s); r++ {
        window[s[r]-'a']++
        if r >= len(p) { // slide: remove leftmost element of previous window
            window[s[r-len(p)]-'a']--
        }
        if window == need { // arrays are comparable in Go
            result = append(result, r-len(p)+1)
        }
    }
    return result
}
```

**Complexity:** O(n * 26) → O(n) time, O(1) space (fixed 26-size arrays).

**Watch out:**
- In Go, `[26]int` arrays are comparable with `==`. This is a huge advantage — no need to manually compare character by character.
- For the more general "Minimum Window Substring" variant, you can't use a fixed window — it's a variable window with a frequency map and a `have`/`need` counter (see walkthrough below).
- Don't confuse this with "fixed-size window." Find Anagrams IS fixed-size (window = len(p)), but Minimum Window Substring is variable-size.

---

## Decision Framework

Read the problem statement. Ask these questions in order:

```
START
  │
  ├─ Is the input SORTED (or should you sort it)?
  │    │
  │    ├─ Finding a PAIR? → Opposite Ends (#1)
  │    │     Two pointers from both ends, O(n)
  │    │
  │    └─ Finding a TRIPLET? → Fix one + Two Pointers (#4)
  │          Sort, outer loop fixes one, inner two-pointer, O(n²)
  │
  ├─ Is it an IN-PLACE array compaction or partition?
  │    YES → Write Pointer (#2)
  │          Slow=write position, fast=reader, O(n) time O(1) space
  │
  ├─ Does it mention CYCLE in linked list or duplicate in [1,n]?
  │    YES → Floyd's Tortoise and Hare (#3)
  │          Two phases: detect collision, then find cycle start
  │
  ├─ Is it a CONTIGUOUS SUBARRAY / SUBSTRING problem?
  │    │
  │    ├─ Fixed window size given (size K)?
  │    │    → Fixed Window (#5)
  │    │
  │    ├─ Asks for LONGEST / MAXIMUM satisfying a constraint?
  │    │    → Variable Window, expand always, shrink when INVALID (#6)
  │    │      answer = max(window size when valid)
  │    │
  │    ├─ Asks for SHORTEST / MINIMUM satisfying a constraint?
  │    │    → Variable Window, expand always, shrink while VALID (#7)
  │    │      answer = min(window size when valid)
  │    │
  │    └─ Character matching between two strings?
  │         → Window + Frequency Map (#8)
  │           Track counts in window vs. target, use have/need counter
  │
  └─ None of the above?
       → Re-read. Can you reformulate?
         "Is subsequence" → same-direction two pointers
         "Trapping rain water" → two pointers from ends (or stack)
         "Dutch national flag" → three-way partition (variant of write pointer)
```

---

## Common Interview Traps

### Trap 1: 3Sum — Forgetting to Skip Duplicates at ALL THREE Levels
You skip duplicates on the outer loop but forget the inner `lo++`/`hi--` skips after finding a match. Result: duplicate triplets in your output, which most judges auto-fail.

**Fix:** After appending a triplet, advance `lo` and `hi`, THEN skip while equal to previous. Three skip points total.

### Trap 2: Sliding Window — Using `if` Instead of `for` for Contraction
You write `if len(freq) > k { left++ }` instead of `for len(freq) > k { left++ }`. The window might need to shrink by MORE than one position to become valid again.

**Fix:** ALWAYS use `for` (while-loop) for the contraction step. Never `if`. Make this a rule with zero exceptions.

### Trap 3: Container With Most Water — Moving the Wrong Pointer
You move the taller side instead of the shorter one. The greedy insight is: moving the shorter side *might* find a taller line and increase area; moving the taller side *can only* decrease or maintain the height constraint.

**Fix:** Always `if height[left] < height[right] { left++ } else { right-- }`. Move the SHORTER one.

### Trap 4: Minimum Window Substring — Incrementing `have` at the Wrong Time
You increment `have` every time you add a character that appears in `t`, even if you already have enough of it. This inflates `have` beyond `need` and breaks the contraction logic.

**Fix:** Only increment `have` when `windowFreq[c] == needFreq[c]` — the EXACT moment you've satisfied that character's requirement, not before and not after.

### Trap 5: Off-by-One on Window Size
You compute window size as `right - left` instead of `right - left + 1`. A window from index 3 to index 5 has 3 elements, not 2.

**Fix:** Window size = `right - left + 1`. Tattoo this on your brain.

### Trap 6: Two Pointers on Unsorted Input
You see "find a pair summing to K" and reach for two pointers. But the input is unsorted and the problem asks for original indices. Sorting destroys index information.

**Fix:** Two pointers requires sorted input (or a structure like Container With Most Water where the array represents heights at positions). If you need original indices, use a hash map.

### Trap 7: Variable Window — Updating Answer at the Wrong Time
- **Find longest:** You update `maxLen` INSIDE the contraction loop (while the window is invalid). You should update AFTER the loop, when the window is valid.
- **Find shortest:** You update `minLen` OUTSIDE the contraction loop (when the window might be oversized). You should update INSIDE the loop, while the window is still valid.

**Fix:**
- Longest: expand → shrink until valid → record answer (outside shrink loop).
- Shortest: expand → while valid: record answer, then shrink (inside shrink loop).

---

## Thought Process Walkthrough

### Problem 1: 3Sum (Two Pointers)

> Given an integer array `nums`, return all triplets `[nums[i], nums[j], nums[k]]` such that `i != j`, `i != k`, `j != k`, and `nums[i] + nums[j] + nums[k] == 0`. The solution set must not contain duplicate triplets.

#### Step 1 — Clarify (30 seconds)

Say to the interviewer:
- "Can the array contain duplicates?" (Yes — that's what makes this problem tricky.)
- "Should I return values or indices?" (Values — so sorting is fine.)
- "Can I sort the input?" (Yes — since we return values, original order doesn't matter.)
- "Can the result include the same element twice?" (No — `i != j != k`.)

#### Step 2 — Brute Force (1 minute)

"The brute force is three nested loops checking all triplets. That's O(n³) time. I also need a set to deduplicate, which is messy."

Don't code this. Just state it and move on.

#### Step 3 — Optimize (3 minutes)

"I can reduce this to O(n²). The insight is: if I sort the array, then for each fixed element `nums[i]`, the problem reduces to Two Sum on the remaining subarray `nums[i+1:]` with target `-nums[i]`. Since the subarray is sorted, I use opposite-end two pointers instead of a hash map — this gives O(1) space and naturally handles the sorted order for deduplication."

"For deduplication: I skip `nums[i]` if it equals `nums[i-1]`. After finding a valid triplet, I skip forward past duplicate `lo` values and backward past duplicate `hi` values."

Walk through an example: `nums = [-1, 0, 1, 2, -1, -4]`

After sorting: `[-4, -1, -1, 0, 1, 2]`

- `i=0, nums[i]=-4`: target=4. `lo=1, hi=5`. Max sum = `-1+2=1 < 4`. No triplet.
- `i=1, nums[i]=-1`: target=1. `lo=2, hi=5`.
  - `-1+2=1` → match! Triplet: `[-1,-1,2]`. `lo→3, hi→4`.
  - `0+1=1` → match! Triplet: `[-1,0,1]`. `lo→4, hi→3`. Done.
- `i=2, nums[i]=-1`: same as `nums[1]`, skip.
- `i=3, nums[i]=0`: target=0. `lo=4, hi=5`. `1+2=3 > 0`, hi--. `lo >= hi`, done.

Result: `[[-1,-1,2], [-1,0,1]]`. Correct.

#### Step 4 — Code (8-10 minutes)

```go
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    var result [][]int

    for i := 0; i < len(nums)-2; i++ {
        // Early termination: if smallest value > 0, no triplet can sum to 0
        if nums[i] > 0 {
            break
        }
        // Skip duplicate values for i
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }

        lo, hi := i+1, len(nums)-1
        target := -nums[i]

        for lo < hi {
            sum := nums[lo] + nums[hi]
            if sum == target {
                result = append(result, []int{nums[i], nums[lo], nums[hi]})
                lo++
                hi--
                // Skip duplicates for lo
                for lo < hi && nums[lo] == nums[lo-1] {
                    lo++
                }
                // Skip duplicates for hi
                for lo < hi && nums[hi] == nums[hi+1] {
                    hi--
                }
            } else if sum < target {
                lo++
            } else {
                hi--
            }
        }
    }

    return result
}
```

#### Step 5 — Test (3 minutes)

Trace through these cases out loud:

- **Standard:** `[-1,0,1,2,-1,-4]` → `[[-1,-1,2],[-1,0,1]]` ✓
- **All zeros:** `[0,0,0]` → `[[0,0,0]]` — works because duplicate skip is `i > 0`, so first zero is processed.
- **No solution:** `[1,2,3]` → `[]` — `nums[0]=1 > 0`, early break.
- **All same:** `[-1,-1,-1,-1]` → target=1, max inner sum = `-1+(-1)=-2 < 1`. No match. `[]`.
- **Two elements only:** `[0,0]` → outer loop condition `i < len(nums)-2` → `i < 0`. No iteration. `[]`.

#### Step 6 — Complexity

- **Time:** O(n²). Sorting is O(n log n). Outer loop is O(n). Inner two-pointer is O(n). Total: O(n log n + n²) = O(n²).
- **Space:** O(1) extra (ignoring the output array and sort stack). If the sort is in-place (Go's sort.Ints is), then O(log n) stack space.

#### Step 7 — Follow-ups

- *"What about 4Sum?"* → Same idea: fix two elements, two-pointer on the rest. O(n³). Generalizes to kSum with O(n^(k-1)).
- *"What if you can't sort?"* → Use hash sets for the inner loop. Harder to deduplicate. Sorting is almost always preferred when you return values.
- *"3Sum closest?"* → Same structure, but track `closestSum` where `abs(sum - target)` is minimized. No need for deduplication.

---

### Problem 2: Minimum Window Substring (Sliding Window — Hardest Common Variant)

> Given two strings `s` and `t`, return the minimum window substring of `s` such that every character in `t` (including duplicates) is included in the window. If no such window exists, return `""`.

#### Step 1 — Clarify (30 seconds)

- "Can `t` have duplicate characters?" (Yes — if `t = "AA"`, the window must contain at least two A's.)
- "Are characters only lowercase letters?" (No — the problem says uppercase and lowercase.)
- "If multiple windows have the same minimum length, return any one?" (Return the first one found.)
- "Can `t` be longer than `s`?" (Yes — return `""` immediately.)

#### Step 2 — Brute Force (1 minute)

"I could check every substring of `s` and verify if it contains all characters of `t`. That's O(n² * m) where n = len(s) and m = len(t). Way too slow."

#### Step 3 — Optimize (4 minutes)

"I'll use a variable-size sliding window. The key insight: I expand the right pointer until the window is valid (contains all chars of `t`), then I shrink from the left to find the minimum valid window."

"State I need to track:"
- `needFreq`: frequency map of characters in `t`. Built once.
- `windowFreq`: frequency map of characters in the current window. Updated as I expand/shrink.
- `have`: number of characters where `windowFreq[c] >= needFreq[c]`. This is the count of **satisfied** unique characters.
- `need`: number of unique characters in `t` (= `len(needFreq)`).
- When `have == need`, the window is valid.

"Critical detail: I only increment `have` when `windowFreq[c]` goes from below `needFreq[c]` to exactly `needFreq[c]`. Not when it goes above."

Walk through: `s = "ADOBECODEBANC", t = "ABC"`

`needFreq = {A:1, B:1, C:1}`, `need = 3`

- Expand right until valid: `ADOBEC` (indices 0-5). `have = 3`. Record window "ADOBEC" (len 6).
- Shrink left: remove A → `have` drops to 2. Window invalid. Stop shrinking.
- Continue expanding: `ADOBECODE` ... `ADOBECODEBA` (have=3 again at index 10).
- Shrink: `DOBECODEBA` → still valid (D isn't needed). `OBECODEBA` → still valid. `BECODEBA` → still valid. `ECODEBA` → still valid. `CODEBA` → still valid (len 6, same as best). `ODEBA` → missing C, have drops. Stop.
- Continue expanding: `ODEBANC` → have=3. Shrink: `DEBANC` → valid. `EBANC` → valid. `BANC` → valid (len 4, new best!). `ANC` → missing B. Stop.

Result: `"BANC"`. Correct.

#### Step 4 — Code (10-12 minutes)

```go
func minWindow(s string, t string) string {
    if len(t) > len(s) {
        return ""
    }

    needFreq := make(map[byte]int)
    for i := 0; i < len(t); i++ {
        needFreq[t[i]]++
    }

    windowFreq := make(map[byte]int)
    have, need := 0, len(needFreq)
    bestLeft, bestLen := 0, math.MaxInt32

    left := 0
    for right := 0; right < len(s); right++ {
        // Expand: add s[right] to window
        c := s[right]
        windowFreq[c]++
        // Only increment "have" at the exact moment we satisfy this char
        if needFreq[c] > 0 && windowFreq[c] == needFreq[c] {
            have++
        }

        // Shrink: while window is valid, try to minimize
        for have == need {
            // Update best
            windowSize := right - left + 1
            if windowSize < bestLen {
                bestLen = windowSize
                bestLeft = left
            }
            // Remove s[left] from window
            d := s[left]
            windowFreq[d]--
            if needFreq[d] > 0 && windowFreq[d] < needFreq[d] {
                have--
            }
            left++
        }
    }

    if bestLen == math.MaxInt32 {
        return ""
    }
    return s[bestLeft : bestLeft+bestLen]
}
```

#### Step 5 — Test (3 minutes)

- **Standard:** `s="ADOBECODEBANC", t="ABC"` → `"BANC"` ✓
- **Exact match:** `s="ABC", t="ABC"` → `"ABC"`. Window becomes valid at right=2, shrinks to len 3, then breaks.
- **No valid window:** `s="ABC", t="D"` → `""`. `have` never reaches `need`.
- **Duplicate chars in t:** `s="AAAB", t="AA"` → `"AA"`. needFreq={A:2}. Window must have 2 A's. First valid at right=1 → "AA" (len 2). ✓
- **Single char:** `s="A", t="A"` → `"A"`. ✓
- **t longer than s:** `s="A", t="AB"` → `""`. Early return. ✓

#### Step 6 — Complexity

- **Time:** O(n + m). Each character in `s` is visited at most twice (once by `right`, once by `left`). Building `needFreq` is O(m).
- **Space:** O(m + |Σ|) where |Σ| is the alphabet size. In practice O(1) if the alphabet is fixed (e.g., 128 ASCII chars).

#### Step 7 — Follow-ups

- *"What if you need all minimum windows, not just the first?"* → Track all windows with length == bestLen instead of just the first.
- *"What if `t` can contain characters not in `s`?"* → It still works. `have` will never reach `need`. Return `""`.
- *"Can you do this with a `[128]int` array instead of a map?"* → Yes, and it's faster. Replace `map[byte]int` with `[128]int`. No map overhead. In interviews, mention this as an optimization.
- *"What if the problem asks for the minimum window containing at least K distinct characters?"* → Different problem. Use variable window (find shortest) with `len(freq) >= K` as the valid condition.

---

## Time Targets

For a standard 45-minute coding interview:

| Phase | Time | What You're Doing |
|---|---|---|
| **Clarify & examples** | 0:00 - 3:00 | Ask 2-3 questions. Write a small example. Confirm input/output format. |
| **Brute force** | 3:00 - 5:00 | State it verbally with complexity. "We can do better." |
| **Optimal approach** | 5:00 - 10:00 | Explain the pattern. Walk through your example step by step. Get interviewer's nod. |
| **Code** | 10:00 - 28:00 | Clean code. Narrate every decision. Name variables well. |
| **Test & debug** | 28:00 - 35:00 | Trace through a small example. Then edge cases. Fix any bugs. |
| **Complexity analysis** | 35:00 - 37:00 | Time and space. Justify the O(n) "each element visited at most twice" argument for sliding window. |
| **Follow-ups** | 37:00 - 45:00 | Variations, tradeoffs, what if constraints change? |

**The #1 time sink for two-pointer problems:** Duplicate handling in 3Sum. Practice until the three skip locations are muscle memory.

**The #1 time sink for sliding window:** Getting the `have`/`need` counter logic wrong in Minimum Window Substring. Practice until the "increment only when `windowFreq[c] == needFreq[c]`" rule is automatic.

---

## Quick Drill

Do each in under 2 minutes. If you can't, you need more reps on that pattern.

### Drill 1: Two Sum on Sorted Array — Opposite Ends in 5 Lines

> "Given a sorted array and a target, return the two indices (1-indexed) that sum to target."

```go
func twoSum(nums []int, target int) []int {
    l, r := 0, len(nums)-1
    for l < r {
        sum := nums[l] + nums[r]
        if sum == target { return []int{l + 1, r + 1} }
        if sum < target { l++ } else { r-- }
    }
    return nil
}
```

### Drill 2: Remove Duplicates from Sorted Array — Write Pointer

> "Remove duplicates in-place, return the new length."

```go
func removeDuplicates(nums []int) int {
    if len(nums) == 0 { return 0 }
    w := 1
    for r := 1; r < len(nums); r++ {
        if nums[r] != nums[w-1] {
            nums[w] = nums[r]
            w++
        }
    }
    return w
}
```

### Drill 3: Max Sum Subarray of Size K — Fixed Window

> "Find the maximum sum of any contiguous subarray of size k."

```go
func maxSum(nums []int, k int) int {
    sum := 0
    for i := 0; i < k; i++ { sum += nums[i] }
    best := sum
    for r := k; r < len(nums); r++ {
        sum += nums[r] - nums[r-k]
        if sum > best { best = sum }
    }
    return best
}
```

### Drill 4: Is Palindrome — Opposite Ends

> "Check if a string is a palindrome (ignore non-alphanumeric, case-insensitive)."

```go
func isPalindrome(s string) bool {
    l, r := 0, len(s)-1
    for l < r {
        for l < r && !isAlnum(s[l]) { l++ }
        for l < r && !isAlnum(s[r]) { r-- }
        if toLower(s[l]) != toLower(s[r]) { return false }
        l++; r--
    }
    return true
}
```

### Drill 5: Longest Substring With All Unique Chars — Variable Window

> "Find the length of the longest substring without repeating characters."

```go
func lengthOfLongestSubstring(s string) int {
    last := make(map[byte]int)
    maxLen, left := 0, 0
    for r := 0; r < len(s); r++ {
        if idx, ok := last[s[r]]; ok && idx >= left {
            left = idx + 1
        }
        last[s[r]] = r
        if r-left+1 > maxLen { maxLen = r - left + 1 }
    }
    return maxLen
}
```

---

## Self-Assessment

Answer each in under 15 seconds without looking at the guide. If you hesitate, revisit that pattern.

**1.** "Problem says: 'find the smallest window containing all characters of T in S.' What pattern? What state do you track?"

> **Expected:** Variable window — find shortest (#7) combined with frequency map (#8). Track `needFreq` (built from T), `windowFreq` (current window counts), `have` (number of fully satisfied characters), and `need` (number of unique chars in T). Expand right until `have == need`, then shrink left while `have == need`, recording the minimum window at each valid position.

**2.** "The input is a sorted array. The problem asks 'find two numbers that add up to target and return their values.' Two pointers or hash map? Why?"

> **Expected:** Two pointers. The input is sorted (perfect for converging pointers) and we need values not indices (sorting won't lose information). O(n) time, O(1) space — strictly better than the hash map's O(n) space.

**3.** "You're given an array and told to move all zeros to the end while maintaining relative order of non-zero elements. What pattern? What are the pointers doing?"

> **Expected:** Write pointer (same-direction, pattern #2). `w` is the write pointer (next position for a non-zero), `r` is the reader scanning left to right. When `nums[r] != 0`, swap `nums[w]` and `nums[r]`, increment `w`. This naturally pushes zeros to the end.

**4.** "A problem says 'find the longest substring with at most 2 distinct characters.' You're expanding the window and the number of distinct characters just hit 3. What do you do?"

> **Expected:** Enter the shrink loop: `for len(freq) > 2`, decrement `freq[s[left]]`, delete from map if count hits 0, increment `left`. Keep shrinking until `len(freq) <= 2`. Then (outside the loop) update `maxLen = max(maxLen, right-left+1)`. The answer is recorded when the window is valid.

**5.** "You're implementing 3Sum and your solution returns duplicate triplets. Where are the three places you need to add duplicate-skipping logic?"

> **Expected:** (1) Outer loop: `if i > 0 && nums[i] == nums[i-1] { continue }`. (2) After finding a match, advance `lo` past duplicates: `for lo < hi && nums[lo] == nums[lo-1] { lo++ }`. (3) After finding a match, retreat `hi` past duplicates: `for lo < hi && nums[hi] == nums[hi+1] { hi-- }`. Missing ANY one of these three produces duplicates.

---

## Study Session Schedule (2 Hours)

| Time | Activity |
|---|---|
| 0:00 - 0:10 | Read the Pattern Catalog. Flag any pattern that isn't instant recall. |
| 0:10 - 0:15 | Run all 5 Quick Drills. Time each one. Target: under 2 min each. |
| 0:15 - 0:35 | **Solve: 3Sum** (LeetCode 15). Full interview simulation. Target: 20 minutes. Focus on getting the three duplicate skips right on the first try. |
| 0:35 - 1:05 | **Solve: Minimum Window Substring** (LeetCode 76). Full interview sim. Target: 25 minutes. This is the boss fight. |
| 1:05 - 1:20 | **Solve: Longest Substring Without Repeating Characters** (LeetCode 3). Target: 10 minutes. If > 10 min, your variable window mechanics need work. |
| 1:20 - 1:35 | **Solve: Container With Most Water** (LeetCode 11). Target: 10 minutes. Practice explaining the greedy argument for why you move the shorter pointer. |
| 1:35 - 1:45 | **Solve: Find the Duplicate Number** (LeetCode 287). Target: 10 minutes. Apply Floyd's cycle detection to an array. |
| 1:45 - 1:55 | Run through the Self-Assessment questions. Be honest about gaps. |
| 1:55 - 2:00 | Write the template for your weakest pattern from memory. Once. No peeking. |

---

*Tomorrow (Day 3): We'll build on pointer mechanics and move into Stack & Queue patterns — monotonic stacks, bracket matching, and sliding window maximums.*
