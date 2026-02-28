# Day 1: Arrays & Hashing Patterns

> **Time block:** 12:00 - 2:00 PM (2 hours)
> **Goal:** Sharpen pattern recognition and solve array/hashing problems on autopilot.
> **Assumption:** You already know how slices, maps, and hash tables work in Go. This is about *seeing the pattern in the first 60 seconds* of reading a problem.

---

## Pattern Catalog

### 1. Frequency Counting

**Trigger:** "When you see *count occurrences*, *most frequent*, *top-K*, *find duplicates*, or *majority element*..."

**Template:**
```go
func frequencyCount(nums []int) map[int]int {
    freq := make(map[int]int)
    for _, n := range nums {
        freq[n]++
    }
    // Now query freq for whatever the problem asks:
    // - freq[x] > 1          → duplicate exists
    // - find max(freq[x])    → most frequent element
    // - sort by freq, take K → top-K frequent
    return freq
}
```

**Complexity:** O(n) time, O(n) space.

**Watch out:** Top-K does NOT require full sorting. Use a bucket sort trick: create buckets where index = frequency, then walk buckets from high to low. This gives O(n) instead of O(n log n).

```go
func topKFrequent(nums []int, k int) []int {
    freq := make(map[int]int)
    for _, n := range nums {
        freq[n]++
    }
    // Bucket sort: index = frequency, value = list of nums with that freq
    buckets := make([][]int, len(nums)+1)
    for num, count := range freq {
        buckets[count] = append(buckets[count], num)
    }
    var result []int
    for i := len(buckets) - 1; i >= 0 && len(result) < k; i-- {
        result = append(result, buckets[i]...)
    }
    return result[:k]
}
```

---

### 2. Complement / Two-Sum Pattern

**Trigger:** "When you see *find a pair that sums to*, *two numbers that add to target*, or *complement*..."

**Template:**
```go
func twoSum(nums []int, target int) (int, int) {
    seen := make(map[int]int) // value → index
    for i, n := range nums {
        complement := target - n
        if j, ok := seen[complement]; ok {
            return j, i
        }
        seen[n] = i
    }
    return -1, -1 // no pair found
}
```

**Complexity:** O(n) time, O(n) space.

**Watch out:** Store the value *after* checking for the complement, not before. If you store first, you might match an element with itself (e.g., target=6, nums=[3,3] — you'd incorrectly return index 0 twice if you store before checking).

Also: the problem might ask for **indices** vs. **values** vs. **boolean**. Read carefully. Returning indices requires storing index in the map. Returning values doesn't.

---

### 3. Grouping by Computed Key

**Trigger:** "When you see *group*, *categorize*, *partition*, or *anagram*..."

**Template:**
```go
func groupByKey(strs []string) map[string][]string {
    groups := make(map[string][]string)
    for _, s := range strs {
        key := computeKey(s) // The art is choosing this function
        groups[key] = append(groups[key], s)
    }
    // Convert map values to result
    result := make([][]string, 0, len(groups))
    for _, group := range groups {
        result = append(result, group)
    }
    return groups
}
```

**Key function choices:**
- **Anagrams:** Sort the characters → `"eat"` → `"aet"`. Or use a character frequency array as key: `[26]int` → convert to string.
- **Same digits:** Sort digits.
- **Same remainder:** `n % k` as the key.

**Complexity:** O(n * k) time where k = cost of computing the key (e.g., k = word length for sorting chars), O(n) space.

**Watch out:** In Go, you cannot use `[26]int` directly as a map key type... but actually you CAN. Arrays (not slices) are comparable in Go. This is a useful trick:

```go
func groupAnagrams(strs []string) [][]string {
    groups := make(map[[26]byte][]string)
    for _, s := range strs {
        var key [26]byte
        for _, c := range s {
            key[c-'a']++
        }
        groups[key] = append(groups[key], s)
    }
    result := make([][]string, 0, len(groups))
    for _, v := range groups {
        result = append(result, v)
    }
    return result
}
```

This avoids the O(k log k) sort cost per word, giving O(n * k) total where k = average word length.

---

### 4. Prefix Sum

**Trigger:** "When you see *subarray sum equals K*, *number of subarrays with sum*, *range sum query*, or *contiguous subarray*..."

**Template:**
```go
func subarraySumEqualsK(nums []int, k int) int {
    prefixCount := map[int]int{0: 1} // CRITICAL: empty prefix
    sum := 0
    count := 0
    for _, n := range nums {
        sum += n
        // If (sum - k) was a previous prefix sum, then the subarray
        // between that point and here sums to k
        if c, ok := prefixCount[sum-k]; ok {
            count += c
        }
        prefixCount[sum]++
    }
    return count
}
```

**Complexity:** O(n) time, O(n) space.

**Watch out:** The `{0: 1}` initialization is the #1 mistake. Without it, you miss subarrays starting at index 0. The empty prefix (sum = 0 before any element) is a valid prefix sum. Forgetting this is an instant bug.

**Mental model:** The prefix sum at index i is `sum(nums[0..i])`. The sum of subarray `nums[j+1..i]` is `prefix[i] - prefix[j]`. If that equals k, then `prefix[j] = prefix[i] - k`. So you look up `sum - k` in your map of previously seen prefix sums.

---

### 5. Index as Hash Key (In-Place Marking)

**Trigger:** "When you see *find missing number*, *find duplicate in range [1, n]*, *O(1) space constraint*, or *values are bounded by array length*..."

**Template:**
```go
func findMissing(nums []int) int {
    // Values are in range [1, n] where n = len(nums)
    // Use the array itself as a hash table: mark index (val-1) as visited
    for _, n := range nums {
        idx := abs(n) - 1
        if idx < len(nums) {
            nums[idx] = -abs(nums[idx]) // mark as seen by negating
        }
    }
    for i, n := range nums {
        if n > 0 {
            return i + 1 // this index was never marked → missing number
        }
    }
    return len(nums) + 1
}

func abs(x int) int {
    if x < 0 { return -x }
    return x
}
```

**Complexity:** O(n) time, **O(1) space** — that's the whole point.

**Watch out:** This mutates the input array. Always ask the interviewer: "Is it okay to modify the input?" If not, you need a different approach. Also, this only works when values are bounded (typically [1, n] or [0, n-1]).

---

### 6. Sliding Window + Hash Map (Bridge to Day 2)

**Trigger:** "When you see *longest substring without repeating*, *minimum window containing*, *at most K distinct*..."

**Template (preview):**
```go
func lengthOfLongestSubstring(s string) int {
    lastSeen := make(map[byte]int) // char → most recent index
    maxLen := 0
    left := 0
    for right := 0; right < len(s); right++ {
        if idx, ok := lastSeen[s[right]]; ok && idx >= left {
            left = idx + 1 // shrink window past the duplicate
        }
        lastSeen[s[right]] = right
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}
```

**Complexity:** O(n) time, O(min(n, charset)) space.

**Watch out:** The `idx >= left` check is critical. Without it, you shrink the window to a stale position from before your current window. This is covered in depth on Day 2.

---

## Decision Framework

Read the problem statement. Ask these questions in order:

```
START
  │
  ├─ Does it ask about PAIRS / TWO NUMBERS summing to a target?
  │    YES → Complement / Two-Sum Pattern (#2)
  │           Use map[value]index, look up (target - current)
  │
  ├─ Does it ask about FREQUENCIES, COUNTS, DUPLICATES, or TOP-K?
  │    YES → Frequency Map (#1)
  │           Build map[value]count, then query it
  │           If Top-K → consider bucket sort over heap
  │
  ├─ Does it ask to GROUP or CATEGORIZE elements by some property?
  │    YES → Group by Computed Key (#3)
  │           Choose a canonical key function, build map[key][]value
  │
  ├─ Does it ask about SUBARRAY SUM = K or RANGE SUMS?
  │    YES → Prefix Sum + Hash Map (#4)
  │           Running sum, store in map, look up (sum - k)
  │           REMEMBER the {0: 1} base case
  │
  ├─ Does it ask about MISSING/DUPLICATE in range [1,n] with O(1) space?
  │    YES → Index as Hash Key (#5)
  │           Use sign-flipping or cyclic sort
  │
  ├─ Does it ask about SUBSTRINGS with a character constraint?
  │    YES → Sliding Window + Hash Map (#6, see Day 2)
  │
  └─ None of the above?
       → Re-read the problem. Can you REFORMULATE it as one of the above?
         "Find if two elements differ by K" → complement pattern with target=K
         "Longest subarray with equal 0s and 1s" → prefix sum (treat 0 as -1)
```

---

## Common Interview Traps

### Trap 1: Self-Matching in Two-Sum
The problem says "find two *different* elements." You store `nums[0]` in the map, then immediately match it when processing `nums[0]` again. **Fix:** Always check complement *before* inserting the current element.

### Trap 2: Forgetting the Empty Prefix in Prefix Sum
You initialize `prefixCount := map[int]int{}` instead of `map[int]int{0: 1}`. Every subarray starting at index 0 is now invisible. **Fix:** Drill this initialization until it's muscle memory.

### Trap 3: Assuming Sorted Input
You see "find pairs that sum to K" and jump to two-pointer. But the input is unsorted, and sorting would lose the original indices (which the problem asks you to return). **Fix:** Default to hash map for unsorted input. Only use two-pointer if input is sorted or you don't need indices.

### Trap 4: "Return Indices" vs. "Return Values" vs. "Return Count"
Three different problems with three different implementations. Returning indices means your map stores `value → index`. Returning values means your map might just be a set. Returning count means you might store `value → count`. **Fix:** Re-read the return type before writing any code.

### Trap 5: Handling Negative Numbers in Prefix Sum
You think "sliding window can solve subarray sum = K" — but it can't when negatives are allowed, because shrinking the window doesn't guarantee the sum decreases. **Fix:** If the array can have negatives, prefix sum is the right pattern, not sliding window.

### Trap 6: Integer Overflow with Large Prefix Sums
In Go, `int` is 64-bit on most platforms, so this is less of an issue than in Java/C++. But if you're summing a million numbers each up to 10^9, you hit ~10^15. **Fix:** Know that Go's `int` is 64-bit (safe up to ~9.2 * 10^18). Mention this awareness to your interviewer.

### Trap 7: Map Iteration Order in Go
You group elements into a map and then iterate to build the result. Go map iteration order is randomized. If the problem requires a specific output order (e.g., "return groups in order of first appearance"), you need a separate slice to track insertion order. **Fix:** Maintain a `keys` slice alongside the map when order matters.

---

## Thought Process Walkthrough

### Problem A: Two Sum

> Given an array of integers `nums` and an integer `target`, return the indices of the two numbers that add up to `target`. Each input has exactly one solution, and you may not use the same element twice.

#### Step 1 — Clarify (ask the interviewer)
- "Can the array contain negative numbers?" (Yes → rules out certain shortcuts.)
- "Can there be duplicate values?" (Yes → my map needs to handle overwrites carefully.)
- "Is there always exactly one solution?" (Yes → I don't need to handle zero or multiple results.)
- "Should I return 0-indexed or 1-indexed?" (Clarify, don't assume.)

#### Step 2 — Brute Force
Check every pair:
```go
for i := 0; i < len(nums); i++ {
    for j := i + 1; j < len(nums); j++ {
        if nums[i]+nums[j] == target {
            return []int{i, j}
        }
    }
}
```
**Complexity:** O(n^2) time, O(1) space.
Say: *"This works but we can do better. For each element, I'm searching for its complement linearly. I can trade space for time by storing seen elements in a hash map for O(1) lookups."*

#### Step 3 — Optimize
Insight: For each element `nums[i]`, I need `target - nums[i]` to exist somewhere earlier in the array. A hash map lets me check that in O(1).

Walk through an example: `nums = [2, 7, 11, 15], target = 9`
- i=0: need 7, map is empty, store {2: 0}
- i=1: need 2, map has 2 at index 0 → return [0, 1]

#### Step 4 — Code
```go
func twoSum(nums []int, target int) []int {
    seen := make(map[int]int) // value → index
    for i, n := range nums {
        complement := target - n
        if j, ok := seen[complement]; ok {
            return []int{j, i}
        }
        seen[n] = i
    }
    return nil // problem guarantees a solution, so this won't be reached
}
```

#### Step 5 — Test
Walk through these cases out loud:
- **Happy path:** `[2, 7, 11, 15], target=9` → `[0, 1]`
- **Duplicates:** `[3, 3], target=6` → `[0, 1]` (complement check before insert handles this)
- **Negatives:** `[-1, 4, 5, 2], target=1` → `[-1 + 2 = 1]` → `[0, 3]`
- **Single pair at end:** `[1, 2, 3, 4], target=7` → `[2, 3]`

#### Step 6 — Follow-ups the Interviewer Might Ask
- *"What if there are multiple pairs?"* → Return all of them; change `return` to `append`.
- *"What if the array is sorted?"* → Two-pointer approach, O(1) space.
- *"What about three-sum?"* → Sort + for each element, run two-pointer on the rest. O(n^2).
- *"What if we need the count of pairs instead of indices?"* → Use `map[int]int` as a frequency map.

---

### Problem B: Group Anagrams

> Given an array of strings `strs`, group the anagrams together. Return the answer in any order.

#### Step 1 — Clarify
- "Are all strings lowercase English letters?" (Yes → 26-letter alphabet, simplifies key.)
- "Can strings be empty?" (Yes → `""` is its own anagram group.)
- "Does output order matter?" (No → simplifies things.)

#### Step 2 — Brute Force
For each string, compare it with every other string to check if they're anagrams (sort both, compare). Group matches together.
**Complexity:** O(n^2 * k log k) where k = max string length. Terrible.

Say: *"Two strings are anagrams if they have the same character frequencies. I can use a canonical form — either sorted characters or a frequency count — as a map key to group them in one pass."*

#### Step 3 — Optimize
**Option A:** Sort each string's characters → use sorted string as key. Cost per string: O(k log k).
**Option B:** Count character frequencies → use `[26]byte` array as key. Cost per string: O(k). Better.

Walk through: `["eat", "tea", "tan", "ate", "nat", "bat"]`
- "eat" → key: {a:1, e:1, t:1} → group 1
- "tea" → key: {a:1, e:1, t:1} → group 1
- "tan" → key: {a:1, n:1, t:1} → group 2
- "ate" → key: {a:1, e:1, t:1} → group 1
- "nat" → key: {a:1, n:1, t:1} → group 2
- "bat" → key: {a:1, b:1, t:1} → group 3

#### Step 4 — Code
```go
func groupAnagrams(strs []string) [][]string {
    groups := make(map[[26]byte][]string)

    for _, s := range strs {
        var key [26]byte
        for i := 0; i < len(s); i++ {
            key[s[i]-'a']++
        }
        groups[key] = append(groups[key], s)
    }

    result := make([][]string, 0, len(groups))
    for _, group := range groups {
        result = append(result, group)
    }
    return result
}
```

#### Step 5 — Test
- **Standard case:** `["eat","tea","tan","ate","nat","bat"]` → `[["eat","tea","ate"],["tan","nat"],["bat"]]`
- **Empty strings:** `[""]` → `[[""]]`
- **Single char:** `["a"]` → `[["a"]]`
- **No anagrams:** `["abc","def"]` → `[["abc"],["def"]]`
- **All same:** `["ab","ba","ab"]` → `[["ab","ba","ab"]]`

#### Step 6 — Follow-ups
- *"What if strings contain unicode?"* → Can't use `[26]byte`. Use `map[rune]int` and serialize to a string key, or sort the runes.
- *"What if the list is huge and doesn't fit in memory?"* → External sort by key, then group consecutive.
- *"Can you do this in a streaming fashion?"* → Yes, emit groups at the end; or use a trie-based approach for incremental grouping.

---

## Time Targets

For a standard 45-minute coding interview:

| Phase | Time | What You're Doing |
|---|---|---|
| **Clarify & examples** | 0:00 - 3:00 | Ask 2-3 clarifying questions. Write 1-2 small examples. |
| **Brute force** | 3:00 - 5:00 | State it verbally. Give the complexity. Don't code it. |
| **Optimal approach** | 5:00 - 10:00 | Explain the insight. Walk through your example with the optimal approach. Get interviewer buy-in before coding. |
| **Code** | 10:00 - 28:00 | Write clean code. Narrate as you go. Use descriptive variable names. |
| **Test & debug** | 28:00 - 35:00 | Trace through with a small example. Check edge cases. Fix bugs. |
| **Complexity analysis** | 35:00 - 37:00 | State time and space. Justify each. |
| **Follow-ups** | 37:00 - 45:00 | Interviewer asks variations. Discuss tradeoffs. |

**The #1 time sink:** Starting to code before you have a clear plan. Spending 5 extra minutes on approach saves 10 minutes of debugging.

---

## Quick Drill

Do each of these in under 2 minutes. If you can't, you need more reps.

### Drill 1: Contains Duplicate
> "Given an array, find if any value appears twice."

**Pattern:** Frequency counting (or just a set).
**One-liner logic:**
```go
func containsDuplicate(nums []int) bool {
    seen := make(map[int]bool)
    for _, n := range nums {
        if seen[n] { return true }
        seen[n] = true
    }
    return false
}
```

### Drill 2: Group Anagrams Key
> "Group strings that are anagrams — what's the key function?"

**Key:** Character frequency as `[26]byte`. Arrays are comparable in Go and can be map keys.
```go
var key [26]byte
for i := 0; i < len(s); i++ { key[s[i]-'a']++ }
// Use key as map key
```

### Drill 3: Two Sum in 5 Lines
> "Find two numbers that add to target — hash map approach."

```go
seen := make(map[int]int)
for i, n := range nums {
    if j, ok := seen[target-n]; ok { return []int{j, i} }
    seen[n] = i
}
return nil
```

### Drill 4: Subarray Sum Equals K
> "Count subarrays with sum exactly K — prefix sum approach."

```go
prefix := map[int]int{0: 1} // don't forget this
sum, count := 0, 0
for _, n := range nums {
    sum += n
    count += prefix[sum-k] // Go returns 0 for missing keys
    prefix[sum]++
}
return count
```

### Drill 5: First Non-Repeating Character
> "Find the first non-repeating character in a string."

**Pattern:** Frequency map, then second pass.
```go
func firstUniq(s string) byte {
    freq := [26]int{}
    for i := 0; i < len(s); i++ { freq[s[i]-'a']++ }
    for i := 0; i < len(s); i++ {
        if freq[s[i]-'a'] == 1 { return s[i] }
    }
    return ' '
}
```

---

## Self-Assessment

Answer these *without looking at the patterns above*. If you can't answer confidently in 15 seconds, review that pattern.

**1.** "A problem says 'find all pairs that sum to K in an unsorted array.' What's your first instinct and why?"

> **Expected:** Hash map complement pattern. For each element x, check if (K - x) exists in the map. O(n) time vs. O(n^2) brute force. NOT two-pointer — that requires sorted input and loses index info.

**2.** "A problem asks for the longest subarray with sum exactly K, and the array may contain negative numbers. Why is prefix sum better than sliding window here?"

> **Expected:** Sliding window assumes that adding elements increases the sum and removing elements decreases it. With negatives, that invariant breaks — expanding the window might decrease the sum. Prefix sum doesn't rely on monotonicity: you just look up (currentSum - K) in the map.

**3.** "You're asked to find the missing number in an array of n integers from range [1, n+1]. You're told O(1) space. What are two different approaches?"

> **Expected:** (a) Math: expected sum `n*(n+1)/2` minus actual sum. (b) XOR: XOR all values with all indices — duplicates cancel out, leaving the missing number. (c) Index marking (sign flip) if the range allows it. Know at least two.

**4.** "You need to group strings by anagram. You choose to sort each string as the key. Your interviewer says 'can you do better?' What do you say?"

> **Expected:** Replace O(k log k) sort-based key with O(k) frequency-count key using `[26]byte`. In Go, fixed-size arrays are comparable and usable as map keys directly, so no serialization needed.

**5.** "A problem asks 'find the most frequent element.' You've built the frequency map. What's the fastest way to find the max — and what's the time complexity?"

> **Expected:** Single pass over the frequency map tracking the max count and corresponding element. O(n) total. Don't sort the map entries (that's O(n log n)). If the problem asks for top-K, use bucket sort (O(n)) not a heap (O(n log k)), unless K is very small.

---

## Study Session Schedule (2 Hours)

| Time | Activity |
|---|---|
| 12:00 - 12:10 | Read through the Pattern Catalog above. Flag any pattern that feels rusty. |
| 12:10 - 12:15 | Run through all 5 Quick Drills. Time yourself. |
| 12:15 - 12:40 | **Solve: Two Sum** (LeetCode 1). Write it fresh. Target: 8 minutes. Then do **Valid Anagram** (LeetCode 242) in 5 minutes. |
| 12:40 - 13:10 | **Solve: Group Anagrams** (LeetCode 49). Full interview simulation — clarify, brute force, optimize, code, test. Target: 20 minutes. |
| 13:10 - 13:40 | **Solve: Subarray Sum Equals K** (LeetCode 560). This one trips people up. Practice the prefix sum pattern. Target: 20 minutes. |
| 13:40 - 13:50 | **Solve: Top K Frequent Elements** (LeetCode 347). Practice the bucket sort optimization. Target: 10 minutes. |
| 13:50 - 13:55 | Run through the Self-Assessment questions. Be honest about gaps. |
| 13:55 - 14:00 | Review any pattern you struggled with. Write the template from memory once. |

---

*Tomorrow (Day 2): Two Pointers & Sliding Window — we'll build on the hash map foundation from today and add the window mechanics.*
