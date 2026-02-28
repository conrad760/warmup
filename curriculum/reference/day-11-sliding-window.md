# Day 11 — Sliding Window: Deep Dive

---

## 1. Curated Learning Resources

| # | Resource | Focus | Time |
|---|----------|-------|------|
| 1 | [Sliding Window — NeetCode](https://neetcode.io/courses/advanced-algorithms/4) | Video walkthrough of fixed and variable window patterns with animated pointer movement. Covers the template and when to expand vs contract. | 15 min |
| 2 | [Minimum Window Substring — NeetCode](https://www.youtube.com/watch?v=jSto0O4AJbM) | Visual explanation of the have/need counter technique. Step-by-step trace showing when `have == need` triggers contraction. The definitive video for this problem. | 15 min |
| 3 | [Sliding Window Technique — Algorithms Made Easy](https://www.youtube.com/watch?v=MK-NZ4hN7rs) | Clear visualization of how both pointers move and why total work is O(n). Good for building intuition before coding. | 10 min |
| 4 | [Longest Substring Without Repeating Characters — LeetCode Editorial](https://leetcode.com/problems/longest-substring-without-repeating-characters/editorial/) | All three approaches (brute force, sliding window with set, sliding window with map) side by side. Shows the optimization from O(n^2) to O(n). | 10 min |
| 5 | [LeetCode Sliding Window Problem Set](https://leetcode.com/tag/sliding-window/) | Sorted by acceptance rate. Start with Easy: Max Avg Subarray. Then Medium: Longest Substring No Repeat, Permutation in String. Then Hard: Min Window Substring. | Reference |
| 6 | [Go Strings, Bytes, Runes — Go Blog](https://go.dev/blog/strings) | How Go represents strings as byte slices, why `s[i]` gives a byte not a rune, and how `range` iterates runes. Critical for window problems on strings. | 10 min |
| 7 | [Sliding Window for Beginners — LeetCode Discuss](https://leetcode.com/discuss/study-guide/657507/Sliding-Window-for-Beginners-Problems-or-Template-or-Sample-Solutions) | Community guide with a general template and 15+ categorized problems. Good for extra practice after the session. | 10 min |
| 8 | [Frequency Counting Patterns in Go — Effective Go Maps](https://go.dev/doc/effective_go#maps) | How Go maps handle missing keys (zero value), which makes `map[byte]int` natural for frequency counting. No need to check existence before incrementing. | 5 min |

---

## 2. Detailed 2-Hour Session Plan

### 12:00 — 12:20 | Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:06 | Read Section 3 below. Understand the difference between fixed and variable windows. Mentally trace the fixed window on `[2,1,5,1,3,2], k=3`: slide through and track the sum at each position. |
| 12:06 - 12:12 | Study the variable window template (Section 3.2). Trace "longest substring without repeating characters" on `"abcabcbb"` by hand: expand right, when you hit a repeat, contract left until the window is valid again. Write down the window boundaries at each step on paper. |
| 12:12 - 12:18 | Read the have/need pattern (Section 3.4). Trace minimum window substring on `s="ADOBECODEBANC", t="ABC"` for just the first 8 characters. Track the `need` map, `window` map, `have` counter, and `need` counter. This is the hardest concept today — invest time here. |
| 12:18 - 12:20 | Review the variants taxonomy (Section 5). For each variant, say out loud: "Fixed size? Variable find longest? Variable find shortest? Count-based?" and name one example problem. |

### 12:20 — 1:20 | Implement (From Scratch)

| Time | Problem | Variant | Notes |
|------|---------|---------|-------|
| 12:20 - 12:32 | `MaxSumSubarrayK` | Fixed window | Warmup. Initialize sum of first k elements. Slide: add `arr[right]`, subtract `arr[right-k]`. Track max. Write tests for: k=1, k=len, all negative, all same values. |
| 12:32 - 12:50 | `LongestSubstringNoRepeat` | Variable (find longest) | Use `map[byte]int` storing last-seen index. When you see a repeat, jump `left` to `max(left, lastSeen[c]+1)`. Update max length at each step. Test: `""`, `"a"`, `"abcabc"`, `"bbbbb"`, `"pwwkew"`. |
| 12:50 - 13:10 | `MinWindowSubstring` | Variable (find shortest) | This is the hard one. Build `need` map from `t`. Maintain `window` map and `have`/`need` counters. Expand right always. When `have == need`, contract left to shrink. Track best start/length. Spend extra time here — this is the most important problem today. |
| 13:10 - 13:20 | `LongestRepeatingReplacement` | Variable (find longest, non-shrinkable) | Track frequency of each char in window. `maxFreq` = most frequent char in window. Window is valid when `windowLen - maxFreq <= k`. Key insight: `maxFreq` never needs to decrease (the window only grows or slides). |

### 1:20 — 1:50 | Solidify (Edge Cases & Variants)

| Time | Activity |
|------|----------|
| 1:20 - 1:32 | `PermutationInString` — check if any permutation of `s1` exists as a substring in `s2`. Fixed window of size `len(s1)`. Use two frequency arrays of size 26 and compare. Or use a sliding window with a `matches` counter tracking how many of the 26 chars have equal frequency. Test: `s1` longer than `s2`, single char, all same chars. |
| 1:32 - 1:42 | Go back to `MinWindowSubstring` and test edge cases: `t` longer than `s`, `t` has repeated chars (`s="aa", t="aa"`), `s == t`, `t` is a single char, no valid window exists. Make sure your have/need logic handles duplicate chars in `t` correctly. |
| 1:42 - 1:50 | Refactor: extract the variable window template into a comment block. Compare your three variable-window implementations. They should all follow the same expand-right/contract-left structure with different state tracking and contraction conditions. |

### 1:50 — 2:00 | Recap (From Memory)

Write down without looking:
1. The variable window template (expand right, contract left, update answer).
2. Why sliding window is O(n) even with the inner while loop (amortized: left pointer moves at most n times total).
3. The have/need technique for minimum window substring in one paragraph.
4. The difference between "find longest" (expand as far as possible, answer is max window) and "find shortest" (expand until valid, then shrink to find minimum).
5. One gotcha per problem implemented today.

---

## 3. Core Concepts Deep Dive

### 3.1 Fixed vs Variable Window: When to Use Each

**Fixed-size window:** The problem specifies the window size k. You compute an initial result for the first k elements, then slide by adding one element on the right and removing one on the left. Every window you examine has exactly k elements.

**When to use:** The problem says "subarray of size k" or "every k consecutive elements." Examples: max sum of k elements, average of each contiguous subarray of size k, max of each window of size k.

**Variable-size window:** The problem does NOT fix the window size. Instead, you have a validity condition. You expand the window until some condition is met (or broken), then contract to restore it (or optimize it). The window size changes throughout.

**When to use:** The problem says "longest/shortest substring/subarray satisfying X." Examples: longest substring without repeating characters (expand until a repeat, contract past it), minimum window containing all target characters (expand until all chars present, contract to minimize).

**Decision rule:** If k is given, it is fixed. If you are optimizing the window size, it is variable.

---

### 3.2 The Variable Window Template

```
left := 0
for right := 0; right < n; right++ {
    // 1. EXPAND: add arr[right] (or s[right]) into window state

    // 2. CONTRACT: while the window is invalid (or, for "find shortest",
    //             while the window is valid and we want to shrink it)
    for <window_invalid_condition> {
        // remove arr[left] from window state
        left++
    }

    // 3. UPDATE: record the answer from the current valid window
    //    (for "find longest": answer = max(answer, right-left+1))
    //    (for "find shortest": update inside the contraction loop)
}
```

**Three steps, every time:** expand, contract, update. The specifics of each step change per problem, but the skeleton is always the same.

**Why total work is O(n) even with the inner while loop:**

This is the most common question about sliding window. It looks like O(n^2) — a for loop with a while loop inside. But it is O(n) amortized.

The key insight: **the left pointer only moves to the right.** It never moves backward. Across the entire execution of the outer for loop (all n iterations), the left pointer moves at most n times total. Each element is added to the window exactly once (when `right` reaches it) and removed from the window at most once (when `left` passes it). So the total number of operations across both pointers is at most 2n.

**Amortized argument in detail:**

```
Total operations = (number of times right advances) + (number of times left advances)
                 = n + (at most n)
                 = O(n)
```

The inner while loop does NOT reset for each iteration of the outer loop. It picks up where it left off because `left` never goes backward. This is the same amortized argument as for two pointers.

---

### 3.3 Window State Management

Different problems require different data to track inside the window. Here are the common patterns:

| State to Track | Data Structure | Example Problem |
|---------------|----------------|-----------------|
| Running sum | Single `int` variable | Max sum of subarray of size k |
| Character frequency | `map[byte]int` or `[26]int` | Permutation in string, min window substring |
| Last-seen index | `map[byte]int` | Longest substring without repeating characters |
| Count of distinct chars | `int` counter + frequency map | Longest substring with at most k distinct characters |
| Max frequency char | `int` tracking max freq | Longest repeating character replacement |
| "Satisfied" char count | `have` / `need` int counters | Minimum window substring |

**Go-specific considerations:**

- `map[byte]int` — for ASCII problems, `s[i]` is a `byte`. Incrementing a missing key works because Go maps return zero value: `m[key]++` works even if `key` is not in the map.
- `[26]int` — for lowercase-only problems, use an array indexed by `s[i] - 'a'`. Faster than a map and avoids allocation.
- `[128]int` — for general ASCII. Covers all printable characters.
- When decrementing: after `windowMap[c]--`, check if the count hit 0. If so, optionally `delete(windowMap, c)` to keep the map clean (useful when checking `len(windowMap)` for distinct character count).

---

### 3.4 The "Have/Need" Pattern

This is the key insight for **Minimum Window Substring** and similar "find a window containing all required elements" problems.

**Setup:**

```go
need := map[byte]int{}   // required character frequencies (from target string t)
for i := 0; i < len(t); i++ {
    need[t[i]]++
}

window := map[byte]int{} // current window's character frequencies
have := 0                // how many characters are FULLY satisfied
needCount := len(need)   // total distinct characters we need to satisfy
```

**What "fully satisfied" means:** Character `c` is fully satisfied when `window[c] >= need[c]`. The counter `have` tracks how many distinct characters have reached their required count.

**Critical rule for incrementing `have`:**

```go
window[s[right]]++
if need[s[right]] > 0 && window[s[right]] == need[s[right]] {
    have++
}
```

You only increment `have` at the **exact moment** `window[c]` reaches `need[c]` — not when it exceeds it. If `window[c]` goes from 2 to 3 but `need[c]` is 2, that is not a new satisfaction event.

**Critical rule for decrementing `have`:**

```go
if need[s[left]] > 0 && window[s[left]] == need[s[left]] {
    have--
}
window[s[left]]--
left++
```

You decrement `have` **before** decrementing the window count, at the exact moment when removing `s[left]` would drop `window[c]` below `need[c]`.

**When to contract:**

```go
for have == needCount {
    // current window is valid — record it if it's the smallest so far
    if right-left+1 < bestLen {
        bestLen = right - left + 1
        bestStart = left
    }
    // try to shrink from the left
    // (decrement have if needed, decrement window, advance left)
}
```

**Why this works:** We expand until we have a valid window (all chars satisfied). Then we greedily shrink from the left to find the smallest valid window anchored at the current `right`. We keep shrinking until the window becomes invalid, then resume expanding.

**Full picture of the state machine:**

```
Expand right → window gains a character → maybe have increases
                                           ↓
                                    have == needCount?
                                    YES → record answer, contract left
                                          window loses a character
                                          maybe have decreases
                                          loop: still have == needCount?
                                    NO  → continue expanding
```

---

### 3.5 Shrinkable vs Non-Shrinkable Windows

Some variable window problems need the window to **both grow and shrink** as it slides. Others work with a window that **only grows or slides (never shrinks).**

**Shrinkable (grow and shrink):**
- **Minimum window substring:** You must shrink to find the smallest valid window. The answer requires exploring smaller windows.
- **Longest substring without repeating characters:** When a repeat is found, you must shrink past the previous occurrence.

Template: the inner `for` loop contracts `left` **as long as** the condition holds (or until valid).

**Non-shrinkable (only grow or slide):**
- **Longest repeating character replacement:** The window never truly shrinks. When the window becomes invalid (`windowLen - maxFreq > k`), you slide it by moving `left` forward by 1 — but you do NOT decrease `maxFreq`. The window size either stays the same or increases.

Why this works for longest-repeating-replacement: The answer is determined by the maximum valid window size ever seen. Once we have found a window of size W, there is no point in checking windows smaller than W. So when the window becomes invalid, we just slide it (maintaining size) rather than shrinking it. `maxFreq` is a historical maximum — it may overcount for the current window, but that only means we might fail to grow the window, never that we record a wrong answer.

**How to tell which you need:**

| Goal | Window Behavior | Contract Condition |
|------|-----------------|--------------------|
| Find **shortest** valid | Must shrink aggressively | `for` loop: shrink while valid |
| Find **longest** valid | May use non-shrinkable trick | `if` statement: slide when invalid |
| Count valid subarrays | Must shrink to count all valid windows | `for` loop: shrink while valid, count at each step |

---

## 4. Implementation Checklist

### Function Signatures

```go
package slidingwindow

// MaxSumSubarrayK returns the maximum sum of any contiguous subarray of size k.
// Precondition: 1 <= k <= len(arr), len(arr) >= 1.
func MaxSumSubarrayK(arr []int, k int) int { ... }

// LongestSubstringNoRepeat returns the length of the longest substring
// without repeating characters.
func LongestSubstringNoRepeat(s string) int { ... }

// MinWindowSubstring returns the smallest substring of s that contains
// all characters of t (including duplicates). Returns "" if no such window.
func MinWindowSubstring(s, t string) string { ... }

// LongestRepeatingReplacement returns the length of the longest substring
// containing the same letter after replacing at most k characters.
func LongestRepeatingReplacement(s string, k int) int { ... }

// PermutationInString returns true if any permutation of s1 is a substring of s2.
func PermutationInString(s1, s2 string) bool { ... }
```

### Test Cases & Edge Cases

| Function | Must-Test Cases |
|----------|----------------|
| `MaxSumSubarrayK` | `k == 1` (max element); `k == len(arr)` (entire array); all negative values `[-3,-2,-1], k=2` → `-3`; all same values; single element array |
| `LongestSubstringNoRepeat` | Empty string → `0`; single char `"a"` → `1`; all same `"aaaa"` → `1`; all unique `"abcd"` → `4`; repeat at end `"abcb"` → `3`; `"pwwkew"` → `3` |
| `MinWindowSubstring` | `t` longer than `s` → `""`; no valid window → `""`; `s == t` → `s`; `t` has repeated chars `s="aa", t="aa"` → `"aa"`; `t` is single char; multiple valid windows (return any smallest); `s="ADOBECODEBANC", t="ABC"` → `"BANC"` |
| `LongestRepeatingReplacement` | `k == 0` (no replacements); `k >= len(s)` (entire string); single char string; all same chars; `"AABABBA", k=1` → `4` |
| `PermutationInString` | `s1` longer than `s2` → `false`; exact match `s1 == s2`; permutation at start, middle, end of `s2`; no permutation exists; `s1` single char; `s1` and `s2` same length |

---

## 5. Sliding Window Variants Taxonomy

### Variant 1: Fixed-Size Window

**Pattern:** The window size k is given. Slide a window of exactly k elements across the array.

**Template:**
```go
// Compute initial window
windowState := computeInitial(arr[0:k])
answer := evaluate(windowState)

for right := k; right < len(arr); right++ {
    // Add the new element entering on the right
    add(arr[right], &windowState)
    // Remove the element leaving on the left
    remove(arr[right-k], &windowState)
    // Update the answer
    answer = best(answer, evaluate(windowState))
}
```

**Canonical example: Max Sum of Subarray of Size K**

```
Input:  [2, 1, 5, 1, 3, 2], k = 3
         ─────
           ─────
             ─────
               ─────
Sums:    8, 7, 9, 6
Answer:  9
```

**Other problems:** Average of subarrays of size k, maximum of each window of size k (use a deque), number of distinct elements in each window.

---

### Variant 2: Variable-Size — Find Longest

**Pattern:** Find the longest subarray/substring satisfying a condition. Expand right as far as possible. Contract left only when the window becomes invalid.

**Template:**
```go
left := 0
answer := 0

for right := 0; right < n; right++ {
    // Expand: add arr[right] to window state

    for windowIsInvalid() {
        // Contract: remove arr[left], left++
        left++
    }

    // The window [left..right] is valid — update answer
    answer = max(answer, right-left+1)
}
```

**Canonical example: Longest Substring Without Repeating Characters**

```
Input: "abcabcbb"

Window: [a]bcabcbb        len=1  max=1
        [ab]cabcbb        len=2  max=2
        [abc]abcbb        len=3  max=3
        a[bca]bcbb        'a' repeated → contract past first 'a', len=3  max=3
        ab[cab]cbb        'b' repeated → contract past first 'b', len=3  max=3
        abc[abc]bb        'c' repeated → contract past first 'c', len=3  max=3
        abca[bcb]b        'b' repeated → contract past first 'b', len=2  max=3
        abcab[cb]b        expand, len=2  max=3
        abcabc[b]b        'b' repeated → contract, len=1  max=3

Answer: 3
```

**Other problems:** Longest substring with at most k distinct characters, longest repeating character replacement, max consecutive ones III.

---

### Variant 3: Variable-Size — Find Shortest

**Pattern:** Find the shortest subarray/substring satisfying a condition. Expand right until valid, then contract left aggressively to minimize.

**Template:**
```go
left := 0
answer := math.MaxInt

for right := 0; right < n; right++ {
    // Expand: add arr[right] to window state

    for windowIsValid() {
        // Current window satisfies the condition — record it
        answer = min(answer, right-left+1)
        // Contract: remove arr[left], left++
        left++
    }
}

if answer == math.MaxInt {
    return -1 // or "" — no valid window found
}
```

**Canonical example: Minimum Window Substring**

```
Input: s = "ADOBECODEBANC", t = "ABC"

Expand until window contains A, B, C:
  [ADOBEC]ODEBANC    have=3, len=6 → record, contract
  A[DOBEC]ODEBANC    lost 'A', have=2 → stop contracting, resume expanding
  A[DOBECODEBA]NC    have=3 again, len=10 → record (not better)
  ...keep contracting...
  ADOBECODE[BANC]    have=3, len=4 → record (best!)

Answer: "BANC"
```

**Other problems:** Minimum size subarray sum (shortest subarray with sum >= target), minimum window containing all elements of a set.

---

### Variant 4: Count-Based — Count Subarrays

**Pattern:** Count the number of subarrays satisfying a condition. Often uses the "at most K" trick: `countExactlyK = countAtMost(K) - countAtMost(K-1)`.

**Template:**
```go
func countAtMost(arr []int, k int) int {
    left := 0
    count := 0

    for right := 0; right < len(arr); right++ {
        // Expand: add arr[right] to window state

        for windowExceedsK() {
            // Contract: remove arr[left], left++
            left++
        }

        // Every subarray ending at right and starting at any index in [left..right]
        // is valid. There are (right - left + 1) such subarrays.
        count += right - left + 1
    }
    return count
}
```

**Canonical example: Subarrays with at most K distinct characters**

```
Input: [1, 2, 1, 2, 3], k = 2

At each right position, count subarrays ending here with at most 2 distinct values:
  right=0: [1]                              → 1 subarray
  right=1: [2], [1,2]                       → 2 subarrays
  right=2: [1], [2,1], [1,2,1]              → 3 subarrays
  right=3: [2], [1,2], [2,1,2], [1,2,1,2]  → 4 subarrays
  right=4: [3], [2,3]                       → 2 subarrays (contracted past [1])

Total at-most-2: 12
Total at-most-1: 5
Exactly-2: 12 - 5 = 7
```

**Other problems:** Number of subarrays with exactly k distinct integers, count of subarrays with sum equal to target (less common — usually prefix sums is better here).

---

## 6. Visual Diagrams

### 6.1 Fixed Window Sliding Across an Array

**Input:** `[2, 1, 5, 1, 3, 2]`, `k = 3`

```
Step 1:  [2, 1, 5, 1, 3, 2]
          =========
          sum = 2+1+5 = 8        max = 8

Step 2:  [2, 1, 5, 1, 3, 2]
             =========
          -2      +1
          sum = 8-2+1 = 7        max = 8

Step 3:  [2, 1, 5, 1, 3, 2]
                =========
             -1      +3
          sum = 7-1+3 = 9        max = 9  ← new max

Step 4:  [2, 1, 5, 1, 3, 2]
                   =========
                -5      +2
          sum = 9-5+2 = 6        max = 9

Answer: 9

The window always has exactly k=3 elements.
Each step: subtract the element that leaves, add the element that enters.
Total operations: n - k + 1 slides = O(n).
```

### 6.2 Variable Window: Longest Substring Without Repeating Characters

**Input:** `"abcbda"`

```
Step 1:  a b c b d a
         L
         R
         window: {a:0}              len=1  max=1

Step 2:  a b c b d a
         L R
         window: {a:0, b:1}        len=2  max=2

Step 3:  a b c b d a
         L   R
         window: {a:0, b:1, c:2}   len=3  max=3

Step 4:  a b c b d a
         L     R
         'b' already in window at index 1
         → move L to max(L, 1+1) = 2
             L R
         window: {c:2, b:3}        len=2  max=3

Step 5:  a b c b d a
             L   R
         window: {c:2, b:3, d:4}   len=3  max=3

Step 6:  a b c b d a
             L     R
         'a' not in window (last seen at 0, but 0 < L=2)
         window: {c:2, b:3, d:4, a:5}  len=4  max=4  ← new max!

Answer: 4 ("cbda")

Key: left pointer NEVER moves backward.
     Total: right moves n times, left moves at most n times → O(n).
```

### 6.3 Have/Need State: Minimum Window Substring Step by Step

**Input:** `s = "ADBANC"`, `t = "ABC"`

```
Setup:  need = {A:1, B:1, C:1}    needCount = 3

─── Expand right ──────────────────────────────────────────────────

R=0, char='A':
   window = {A:1}
   A: window[A]==need[A] → have++     have=1, need=3
   have != needCount → keep expanding

   s:  [A] D  B  A  N  C
        L
        R

R=1, char='D':
   window = {A:1, D:1}
   D not in need → have unchanged    have=1, need=3

   s:  [A  D] B  A  N  C
        L  R

R=2, char='B':
   window = {A:1, D:1, B:1}
   B: window[B]==need[B] → have++     have=2, need=3

   s:  [A  D  B] A  N  C
        L     R

R=3, char='A':
   window = {A:2, D:1, B:1}
   A: window[A]=2 != need[A]=1 → no change   have=2, need=3

   s:  [A  D  B  A] N  C
        L        R

R=4, char='N':
   window = {A:2, D:1, B:1, N:1}
   N not in need → no change          have=2, need=3

   s:  [A  D  B  A  N] C
        L           R

R=5, char='C':
   window = {A:2, D:1, B:1, N:1, C:1}
   C: window[C]==need[C] → have++     have=3, need=3

   s:  [A  D  B  A  N  C]
        L              R

   ★ have == needCount! Window "ADBANC" is valid, len=6.
   Record best = (start=0, len=6)

─── Contract left ─────────────────────────────────────────────────

   Now shrink from left:

   Remove s[0]='A': window[A]=2, need[A]=1
     window[A]==need[A]? 2==1? No → have stays 3
     window[A]-- → window = {A:1, D:1, B:1, N:1, C:1}
     left=1

   s:   A [D  B  A  N  C]
           L           R

   ★ have still == 3! Window "DBANC" len=5 < 6.
   Record best = (start=1, len=5)

   Remove s[1]='D': D not in need → have stays 3
     window = {A:1, B:1, N:1, C:1}
     left=2

   s:   A  D [B  A  N  C]
              L        R

   ★ have still == 3! Window "BANC" len=4 < 5.
   Record best = (start=2, len=4)

   Remove s[2]='B': window[B]=1, need[B]=1
     window[B]==need[B]? 1==1? YES → have-- → have=2
     window[B]-- → window = {A:1, N:1, C:1}
     left=3

   s:   A  D  B [A  N  C]
                  L     R

   have=2 != 3 → stop contracting. No more right to expand.

─── Done ──────────────────────────────────────────────────────────

Answer: "BANC" (start=2, len=4)

Key observations:
  • have increments ONLY when window[c] reaches need[c] (exact match)
  • have decrements ONLY when window[c] is about to drop below need[c]
  • Extra copies (A appears twice) do NOT affect have
  • Contraction is aggressive: keep shrinking while window is valid
```

---

## 7. Self-Assessment

Answer these without looking at your code or notes. If you struggle with any, revisit the relevant section.

### Question 1
**Why is the sliding window O(n) even though there's a while loop inside the for loop?**

<details>
<summary>Answer</summary>

The left pointer only moves to the right — it never resets or moves backward. Across the entire execution of the outer for loop, the left pointer advances at most n times total. Each element is added to the window exactly once (when right reaches it) and removed at most once (when left passes it). The total work is at most 2n operations, which is O(n). The inner while loop's iterations are "paid for" by the left pointer advancement, not multiplied by the outer loop.

</details>

### Question 2
**What is the difference between finding the longest vs shortest valid window?**

<details>
<summary>Answer</summary>

For the **longest** valid window: you expand right freely and only contract left when the window becomes *invalid*. You update the answer *after* contracting (the window is valid and as large as possible). The contraction condition is `for windowIsInvalid()`.

For the **shortest** valid window: you expand right until the window becomes *valid*, then contract left *while the window is still valid*, recording the answer at each valid position. The contraction condition is `for windowIsValid()`. You update the answer *inside* the contraction loop.

The key distinction: longest contracts to restore validity, shortest contracts to exploit validity.

</details>

### Question 3
**In the have/need pattern, why do you increment `have` only when `window[c] == need[c]` and not when `window[c] >= need[c]`?**

<details>
<summary>Answer</summary>

`have` counts the number of *distinct characters* that are fully satisfied. Each character should contribute at most 1 to `have`. If you incremented `have` on every `window[c] >= need[c]`, then adding a third 'A' when only 1 is needed would increment `have` again, overcounting. The `==` check ensures `have` increments exactly once per character — at the precise moment the count reaches the requirement. Similarly, `have` decrements only when the count is about to drop *below* the requirement, using the check `window[c] == need[c]` before decrementing the window count.

</details>

### Question 4
**In Longest Repeating Character Replacement, why does `maxFreq` not need to decrease when the window shrinks?**

<details>
<summary>Answer</summary>

The answer is the maximum valid window size. A window is valid when `windowSize - maxFreq <= k`. Once we have found a valid window of size W, we only care about finding windows *larger* than W. If the current window is invalid and we slide it (move left by 1), the window size stays the same or grows — it never shrinks below the best we have seen. Since we only need a *larger* window to improve the answer, and a larger window needs an equal or greater `maxFreq` to be valid, the historical maximum of `maxFreq` is sufficient. An overestimated `maxFreq` might make us think an invalid window could be valid, but that just means the window stays the same size (no improvement). It never leads to recording an incorrect answer because we only update the answer when the window actually is valid.

</details>

### Question 5
**You're given a problem: "Find the number of subarrays with exactly K distinct integers." How do you apply sliding window?**

<details>
<summary>Answer</summary>

Direct sliding window does not easily count "exactly K" because both expanding and contracting can change whether the count equals K. The trick is: `exactlyK(K) = atMostK(K) - atMostK(K-1)`. The `atMostK` function uses a standard variable-size sliding window that contracts when the distinct count exceeds K. At each position of `right`, the number of valid subarrays ending at `right` is `right - left + 1` (every subarray starting from any index in `[left..right]` has at most K distinct values). Sum these counts to get `atMostK`. Subtract to get the exact count.

</details>
