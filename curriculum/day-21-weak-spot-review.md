# Day 21: Weak Spot Review & Next Steps

> **Time:** 2 hours | **Goal:** Identify gaps, drill weak spots, build a plan for continued growth | **Language:** Go

This is the final day. No new patterns. Today you consolidate everything from the past 20 days into durable skill. Be honest with yourself during the self-evaluation — the gaps you identify today are the ones that will cost you in a real interview.

---

## Self-Evaluation Checklist

For each topic day, answer three questions honestly. Mark **Yes** or **No**. Any "No" means that topic is a candidate for drilling today.

| Day | Topic | Identify pattern in < 2 min? | Code template from memory in < 5 min? | Handle edge cases without prompting? |
|-----|-------|:---:|:---:|:---:|
| 1 | Two Pointers | _ | _ | _ |
| 2 | Sliding Window | _ | _ | _ |
| 3 | Binary Search | _ | _ | _ |
| 4 | Stacks & Monotonic Stacks | _ | _ | _ |
| 5 | Queues, Deques & Monotonic Queues | _ | _ | _ |
| 6 | Linked Lists | _ | _ | _ |
| 7 | Hash Maps & Sets | _ | _ | _ |
| 8 | Recursion & Backtracking | _ | _ | _ |
| 9 | Binary Trees (DFS) | _ | _ | _ |
| 10 | Binary Trees (BFS) | _ | _ | _ |
| 11 | BSTs & Sorted Tree Problems | _ | _ | _ |
| 12 | Heaps & Priority Queues | _ | _ | _ |
| 13 | Tries | _ | _ | _ |
| 14 | Graphs (BFS/DFS) | _ | _ | _ |
| 15 | Graphs (Topological Sort, Union-Find) | _ | _ | _ |
| 16 | Dynamic Programming (1D) | _ | _ | _ |
| 17 | Dynamic Programming (2D & Knapsack) | _ | _ | _ |
| 18 | Greedy & Intervals | _ | _ | _ |

### Scoring

- **0-2 "No" answers total:** You're in strong shape. Use drilling time for speed practice.
- **3-6 "No" answers:** Targeted drilling today will close the gaps. Focus on the topics with the most "No" marks.
- **7+ "No" answers:** You need more than one day. Prioritize the 3 weakest topics today and schedule additional review sessions this week.

---

## Pattern Recognition Speed Test

Set a timer. For each problem description, identify the **pattern** and the **primary data structure** before revealing the answer. Target: 30 seconds per problem, 10 minutes total.

Read the problem, say your answer out loud, then check.

### Problem 1
> Remove duplicates from a sorted array in-place.

<details><summary>Answer</summary>

**Pattern:** Two pointers (read/write)
**Structure:** Array with slow/fast pointer

</details>

### Problem 2
> Find the longest substring with at most K distinct characters.

<details><summary>Answer</summary>

**Pattern:** Sliding window (variable width, find longest)
**Structure:** Hash map for character frequency + left/right pointers

</details>

### Problem 3
> Find the first position where you could insert a value into a sorted array.

<details><summary>Answer</summary>

**Pattern:** Binary search (left-bound / bisect-left)
**Structure:** Array with lo/hi pointers

</details>

### Problem 4
> Given daily temperatures, find how many days until a warmer day for each position.

<details><summary>Answer</summary>

**Pattern:** Monotonic stack (decreasing, pop on greater element)
**Structure:** Stack of indices

</details>

### Problem 5
> Find the maximum value in every contiguous window of size K.

<details><summary>Answer</summary>

**Pattern:** Monotonic deque (decreasing, sliding window max)
**Structure:** Deque of indices

</details>

### Problem 6
> Determine if a linked list has a cycle, and find where the cycle begins.

<details><summary>Answer</summary>

**Pattern:** Two pointers (slow/fast — Floyd's cycle detection)
**Structure:** Linked list

</details>

### Problem 7
> Given two arrays, find their intersection (each element appears as many times as it shows in both).

<details><summary>Answer</summary>

**Pattern:** Hash map frequency count
**Structure:** Hash map

</details>

### Problem 8
> Generate all valid combinations of N pairs of parentheses.

<details><summary>Answer</summary>

**Pattern:** Backtracking with pruning (open < n, close < open)
**Structure:** Recursion + string/slice builder

</details>

### Problem 9
> Find the diameter (longest path between any two nodes) of a binary tree.

<details><summary>Answer</summary>

**Pattern:** Binary tree DFS (post-order, track global max)
**Structure:** Binary tree, recursive depth function

</details>

### Problem 10
> Return the values of a binary tree in level order (grouped by level).

<details><summary>Answer</summary>

**Pattern:** BFS level-order traversal
**Structure:** Queue (slice in Go)

</details>

### Problem 11
> Given a BST and two node values, find their lowest common ancestor.

<details><summary>Answer</summary>

**Pattern:** BST property exploitation (both left -> go left, both right -> go right, split -> found)
**Structure:** BST, iterative or recursive

</details>

### Problem 12
> Find the K-th largest element in an unsorted array.

<details><summary>Answer</summary>

**Pattern:** Min-heap of size K (or quickselect)
**Structure:** `container/heap` with min-heap of size K

</details>

### Problem 13
> Design a system that adds words and searches for words where '.' matches any single character.

<details><summary>Answer</summary>

**Pattern:** Trie with DFS search (branch on wildcard)
**Structure:** Trie (TrieNode with children map/array)

</details>

### Problem 14
> Count the number of islands in a 2D grid of '1's and '0's.

<details><summary>Answer</summary>

**Pattern:** Graph DFS/BFS on a grid (flood fill, mark visited)
**Structure:** 2D grid, visited set or in-place marking

</details>

### Problem 15
> Given a list of course prerequisites, find a valid order to take all courses.

<details><summary>Answer</summary>

**Pattern:** Topological sort (Kahn's BFS or DFS)
**Structure:** Adjacency list + in-degree array + queue

</details>

### Problem 16
> Find the maximum sum of a non-empty contiguous subarray.

<details><summary>Answer</summary>

**Pattern:** 1D DP (Kadane's algorithm — dp[i] = max(nums[i], dp[i-1]+nums[i]))
**Structure:** Single variable (running max)

</details>

### Problem 17
> Given a set of coin denominations, find the fewest coins needed to make a target amount.

<details><summary>Answer</summary>

**Pattern:** DP (unbounded knapsack variant — 1D DP over amount)
**Structure:** 1D DP array, `dp[i] = min coins to make amount i`

</details>

### Problem 18
> Given a set of intervals, find the minimum number of conference rooms required.

<details><summary>Answer</summary>

**Pattern:** Greedy with sorting (sort by start, use min-heap for end times)
**Structure:** Sort + min-heap (or sweep line)

</details>

### Problem 19
> Merge K sorted linked lists into one sorted list.

<details><summary>Answer</summary>

**Pattern:** Heap (min-heap of K list heads, extract-min and push next)
**Structure:** `container/heap` with linked list nodes

</details>

### Problem 20
> Given a string and a dictionary of words, determine if the string can be segmented into dictionary words.

<details><summary>Answer</summary>

**Pattern:** 1D DP (dp[i] = can s[:i] be segmented) or BFS/backtracking with memoization
**Structure:** DP array + hash set for dictionary

</details>

### Score yourself

| Result | Meaning |
|--------|---------|
| 18-20 correct | Pattern recognition is solid. Focus on speed and edge cases. |
| 14-17 correct | Good foundation. Drill the patterns you missed. |
| 10-13 correct | Several gaps. Prioritize the missed patterns in your drilling session. |
| < 10 correct | Extend your review beyond today. Schedule 3-4 more sessions this week. |

---

## Weak Spot Drilling Protocol

Pick your 2-3 weakest topics from the self-evaluation. For each one, follow this exact protocol.

### Step 1: Re-read the Pattern Catalog (5 minutes)

Open the day's curriculum file and re-read only the pattern description, the template code, and the edge cases list. Do not re-read the full problem solutions. You are refreshing the mental model, not re-learning.

```
curriculum/day-XX-topic.md  ->  Pattern section + Template section + Edge cases
```

### Step 2: Code the Template from Memory (5 minutes)

Close the file. Open a blank Go file. Write the core template from memory.

```go
// Example: sliding window (variable width, find longest)
// Write this without looking. If you stall for more than 60 seconds
// on any part, that's the specific sub-skill to target.

func slidingWindowLongest(s string, k int) int {
    freq := make(map[byte]int)
    left, maxLen := 0, 0
    for right := 0; right < len(s); right++ {
        freq[s[right]]++
        for /* window invalid */ len(freq) > k {
            freq[s[left]]--
            if freq[s[left]] == 0 {
                delete(freq, s[left])
            }
            left++
        }
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}
```

Compare against the reference. Note every mistake. Those are your specific gaps.

### Step 3: Solve One Medium Problem (15 minutes, hard timer)

Pick one medium problem that uses this pattern. Set a 15-minute timer and solve it end-to-end: understand, plan, code, test.

Suggested problems by topic (LeetCode numbers):

| Topic | Problem |
|-------|---------|
| Two Pointers | 3Sum (#15) |
| Sliding Window | Minimum Window Substring (#76) |
| Binary Search | Search in Rotated Sorted Array (#33) |
| Stacks | Evaluate Reverse Polish Notation (#150) |
| Queues/Deques | Sliding Window Maximum (#239) |
| Linked Lists | Reorder List (#143) |
| Hash Maps | Group Anagrams (#49) |
| Backtracking | Combination Sum (#39) |
| Trees (DFS) | Path Sum III (#437) |
| Trees (BFS) | Binary Tree Zigzag Level Order (#103) |
| BSTs | Validate BST (#98) |
| Heaps | Top K Frequent Elements (#347) |
| Tries | Word Search II (#212) |
| Graphs (BFS/DFS) | Course Schedule (#207) |
| Graphs (Topo/UF) | Redundant Connection (#684) |
| DP (1D) | House Robber (#198) |
| DP (2D) | Unique Paths (#62) |
| Greedy/Intervals | Non-overlapping Intervals (#435) |

### Step 4: Evaluate and Move On

- **Solved in < 15 min with clean code:** This topic is recovered. Move to the next weak spot.
- **Solved in 15 min but messy:** Spend 5 more minutes cleaning up. Note the rough edges for future review.
- **Could not solve in 15 min:** This topic needs a full session. Schedule it for tomorrow. For now, read the solution, understand it, and move on to the next weak spot.

---

## Session Plan (2 Hours)

```
12:00 - 12:15  Self-Evaluation Checklist
               Fill in the table above. Be honest.
               Count your "No" answers. Identify your 2-3 weakest topics.

12:15 - 12:30  Pattern Recognition Speed Test
               Go through all 20 problems.
               Score yourself. Add any missed patterns to your weak list.

12:30 - 12:55  Drill Weak Spot #1 (25 min)
               5 min re-read + 5 min template + 15 min problem

12:55 - 13:20  Drill Weak Spot #2 (25 min)
               5 min re-read + 5 min template + 15 min problem

13:20 - 13:45  Drill Weak Spot #3 (25 min)
               5 min re-read + 5 min template + 15 min problem
               (If no third weak spot, use this time for a timed
               medium problem from any topic for speed practice.)

13:45 - 14:00  Build Your Continued Practice Plan
               Fill in the template below. This is your roadmap
               for the days and weeks after this curriculum.
```

---

## Continued Practice Plan Template

Copy this template and fill it in. This is your plan for after the 21 days.

### Daily Practice (10-15 minutes)

Use the warmup CLI tool every morning for spaced repetition:

```bash
warmup daily
```

This gives you 1-2 short problems from your SRS queue. The algorithm surfaces patterns you haven't practiced recently or ones you've struggled with. Do this before your morning coffee becomes a habit.

### Weekly Practice (3-4 hours total)

Follow the **5-problem weekly plan**:

| Day | Type | Purpose | Time |
|-----|------|---------|------|
| Monday | 1 Easy | Warm up, build confidence, reinforce basics | 15 min |
| Tuesday | 1 Medium (weak topic) | Drill a pattern from your weak list | 30 min |
| Wednesday | 1 Medium (random topic) | Maintain breadth across all patterns | 30 min |
| Thursday | 1 Hard | Stretch problem — combine multiple patterns | 45 min |
| Friday | 1 Review (from SRS) | Revisit a problem you solved before, re-solve from scratch | 20 min |

### Monthly Check-in

Once a month, redo the pattern recognition speed test from this page. Track your score over time. If a topic keeps appearing in your weak list, give it a dedicated 2-hour session.

### Pre-Interview Protocol

When you have an interview scheduled:

**One week before:**
- Do 2-3 full mock sessions using the Day 19/Day 20 format
- Time yourself strictly: 35 minutes per problem
- Practice explaining your thought process out loud
- Review your weak topics one more time

**Two days before:**
- Skim the pattern catalog for all 18 topics (not deep study, just recognition refresh)
- Do 2-3 easy problems to keep your hands warm
- Stop studying the night before — diminishing returns past this point

**The morning of:**
- One easy warmup problem (5-10 minutes)
- Review the Interview Day Cheat Sheet below
- That's it. Trust your preparation.

### Signs You're Ready

You are interview-ready when all five of these are true:

- [ ] You can identify the correct pattern within 2 minutes of reading a problem
- [ ] You can solve most medium problems in under 20 minutes
- [ ] You can explain your approach clearly before writing any code
- [ ] You can spot and handle edge cases without being prompted
- [ ] You can analyze time/space complexity and discuss tradeoffs on the fly

If any of these feel shaky, keep practicing. There is no shortcut.

---

## Interview Day Cheat Sheet

Print this or keep it on your phone. Review it 30 minutes before the interview.

### Before the Interview

- Review the pattern catalog (skim, not study)
- Sleep 7+ hours the night before
- Have water and a pen/paper nearby
- Test your mic, camera, and screen-sharing setup
- Have your IDE ready with a blank Go file

### During the Interview: The 35-Minute Framework

```
 0:00 -  3:00  CLARIFY
               - Restate the problem in your own words
               - Ask about input constraints (size, range, duplicates, negative numbers)
               - Ask about edge cases (empty input, single element, all same values)
               - Confirm the expected return type

 3:00 -  5:00  BRUTE FORCE
               - State the naive solution and its complexity
               - "The brute force would be O(n^2) using nested loops..."
               - This shows you understand the problem even if you optimize

 5:00 - 10:00  OPTIMIZE
               - Identify the pattern (say it out loud)
               - "This looks like a sliding window problem because..."
               - Walk through your approach on an example
               - State the time/space complexity of your approach
               - Get a verbal "go ahead" from the interviewer before coding

10:00 - 30:00  CODE
               - Write clean, readable code
               - Use descriptive variable names (left, right, not i, j)
               - Talk while you code: "Now I'm handling the case where..."
               - If you get stuck, say so: "I'm thinking about how to handle X"
               - Do not go silent for more than 15 seconds

30:00 - 35:00  TEST
               - Trace through your code with a small example
               - Test an edge case (empty, single element, max size)
               - Fix any bugs you find
               - State the final time and space complexity
```

### Communication Tips

**Do:**
- Think out loud at all times
- State your assumptions explicitly: "I'm assuming the input is non-empty"
- When choosing between approaches, explain the tradeoff: "I could use a hash map for O(n) time but O(n) space, or sort for O(n log n) time and O(1) space"
- If you realize a mistake, say: "Wait, I see a bug here — let me fix it"
- Ask the interviewer: "Does this approach make sense before I start coding?"

**Don't:**
- Code in silence
- Jump straight to code without explaining your plan
- Ignore edge cases until the interviewer asks
- Say "I've seen this problem before" (even if you have)
- Panic if you're stuck — take a breath, re-read the problem, think about what data structure would help

### Red Flags to Avoid

These are the things interviewers consistently flag as negative signals:

1. **Silent coding** — The interviewer can't evaluate your thinking if you don't share it
2. **No plan before coding** — Writing code without a stated approach suggests you're guessing
3. **Not testing** — Finishing code and saying "I think that's right" without tracing through an example
4. **Ignoring hints** — If the interviewer says "What about the case where X is empty?", they're telling you there's a bug
5. **Over-engineering** — Using a segment tree when a hash map will do. Start simple.
6. **Poor naming** — Single-letter variables everywhere make your code hard to discuss
7. **Not knowing complexity** — You should always be able to state time and space complexity

---

## Final Notes

You've spent 21 days building a systematic toolkit for coding interviews. The patterns don't change. The problems do, but they're all combinations of the same 18 building blocks you've practiced.

From here, the work is maintenance and confidence:
- Use spaced repetition to keep the patterns fresh
- Practice under timed conditions to build speed
- Do mock interviews to build comfort with the communication layer

The curriculum is done. The preparation is not — it's ongoing. But you now have the structure to make every hour of practice count.

Good luck.
