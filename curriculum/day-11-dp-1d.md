# Day 11: Dynamic Programming — 1D

> **Time:** 2 hours | **Level:** Refresher | **Language:** Go
>
> DP is the #1 interview killer. The difference between candidates who solve DP
> and those who don't is almost never intelligence — it's having a **repeatable
> recipe** and the discipline to follow it out loud. Today is about drilling
> that recipe until it's automatic.

---

## The DP Recipe (Interview Version)

This is the single most important section. Memorize these five steps. Every
1D DP problem you will ever see in an interview follows this skeleton.

### The 5 Steps

```
Step 1: DEFINE STATE    → "dp[i] represents ___"
Step 2: RECURRENCE      → "dp[i] = ___ because ___"
Step 3: BASE CASES      → "dp[0] = ___ because ___"
Step 4: ITERATION ORDER → left-to-right? right-to-left?
Step 5: SPACE OPTIMIZE  → can I drop the array for O(1)?
```

### How to Verbalize This to an Interviewer

Say these words out loud before you write a single line of code:

> "Let me define my state first. I'll let dp[i] represent the maximum amount
> I can rob from the first i houses. Then the recurrence is: for each house i,
> I either rob it and add dp[i-2], or skip it and take dp[i-1]. The base cases
> are dp[0] = nums[0] and dp[1] = max(nums[0], nums[1]). I'll iterate left to
> right since each state depends on smaller indices. And since I only look back
> two steps, I can optimize to O(1) space."

This takes 20 seconds and tells the interviewer:
- You have a framework (you won't flail)
- You understand the problem before coding
- You can reason about correctness
- You're thinking about optimization

**Practice saying this out loud for every problem today.** Silent DP practice
is half-value practice.

---

## Decision Framework: Is This a DP Problem?

Before applying the recipe, you need to recognize DP. Look for these signals:

| Signal | Example |
|---|---|
| "Number of ways to ___" | Climbing stairs, decode ways |
| "Minimum/maximum cost to ___" | Coin change, min cost climbing stairs |
| "Can you ___?" (feasibility) | Word break, can you reach end |
| Choices at each step + overlapping subproblems | House robber, paint houses |
| Greedy seems right but a counterexample exists | Coin change with coins [1,3,4] amount=6 |

**Anti-signals** (probably NOT 1D DP):
- Input is a tree or graph → likely DFS/BFS (or tree DP, not 1D)
- You need all valid combinations listed → backtracking
- Greedy provably works (e.g., activity selection with unit weights)

---

## Pattern Catalog

### Pattern 1: Fibonacci-Style

**Trigger:** dp[i] depends only on dp[i-1] and dp[i-2]. The problem feels
like "how many ways to reach step i?"

**The Recipe Applied:**
```
State:      dp[i] = number of ways to reach step i (or best value at step i)
Recurrence: dp[i] = dp[i-1] + dp[i-2]        (ways — add)
            dp[i] = min(dp[i-1], dp[i-2]) + cost[i]  (optimization — min/max)
Base cases: dp[0] = 1 (or 0), dp[1] = 1 (or derived)
Order:      left to right
Space:      O(1) — only need two previous values
```

**Problems:** Climbing Stairs, Min Cost Climbing Stairs, House Robber (variant),
Decode Ways, Tribonacci.

**Go Template:**
```go
func climbStairs(n int) int {
    if n <= 2 {
        return n
    }
    prev2, prev1 := 1, 2
    for i := 3; i <= n; i++ {
        curr := prev1 + prev2
        prev2 = prev1
        prev1 = curr
    }
    return prev1
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- Decode Ways looks Fibonacci but has conditional transitions (zero digits).
  The recurrence has `if` guards, not a clean sum.
- Off-by-one: does dp[0] mean "0 stairs" or "1st stair"? Define it clearly
  before coding.

---

### Pattern 2: Best Ending Here

**Trigger:** You need the best (max/min/longest) subarray or subsequence
ending at each index. The key insight: the answer at position i *extends*
or *restarts* from position i.

**The Recipe Applied:**
```
State:      dp[i] = best value of subarray/subsequence ending at index i
Recurrence: dp[i] = max(dp[i-1] + nums[i], nums[i])       (Kadane's)
            dp[i] = max(dp[j] + 1) for all j < i where nums[j] < nums[i]  (LIS)
Base cases: dp[0] = nums[0] (Kadane's) or dp[i] = 1 for all i (LIS)
Order:      left to right
Space:      O(1) for Kadane's, O(n) for LIS
```

**Problems:** Maximum Subarray (Kadane's), Longest Increasing Subsequence,
Maximum Product Subarray, Word Break.

**Go Template (Kadane's):**
```go
func maxSubArray(nums []int) int {
    // Step 1: dp[i] = max sum of subarray ending at i
    // Step 2: dp[i] = max(dp[i-1] + nums[i], nums[i])
    // Step 3: dp[0] = nums[0]
    // Step 4: left to right
    // Step 5: only need previous value → O(1)
    
    best := nums[0]
    curr := nums[0]
    for i := 1; i < len(nums); i++ {
        curr = max(curr+nums[i], nums[i]) // extend or restart
        best = max(best, curr)
    }
    return best
}
```

**Go Template (LIS — O(n^2)):**
```go
func lengthOfLIS(nums []int) int {
    n := len(nums)
    dp := make([]int, n) // dp[i] = length of LIS ending at i
    for i := range dp {
        dp[i] = 1 // base: every element is an LIS of length 1
    }
    result := 1
    for i := 1; i < n; i++ {
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] {
                dp[i] = max(dp[i], dp[j]+1)
            }
        }
        result = max(result, dp[i])
    }
    return result
}
```

**Complexity:** Kadane's O(n)/O(1). LIS O(n^2)/O(n).

**Watch out:**
- Maximum Product Subarray: you need to track both `maxEndingHere` and
  `minEndingHere` because a negative times a negative flips sign.
- LIS: the O(n^2) solution is perfectly fine for interviews. Don't jump to
  the O(n log n) patience sorting approach unless explicitly asked.
- Word Break: "ending here" means the substring `s[0..i]` is breakable.
  The inner loop checks all split points.

---

### Pattern 3: Unbounded Choices

**Trigger:** You have a set of choices (coins, square numbers, etc.) and can
use each choice unlimited times. You want min/max/count to reach a target.

**The Recipe Applied:**
```
State:      dp[i] = min coins (or ways, or count) to make amount i
Recurrence: dp[i] = min(dp[i-coin] + 1) for each coin   (minimize)
            dp[i] = sum(dp[i-coin]) for each coin        (count ways)
Base case:  dp[0] = 0 (need 0 coins/squares to make 0)
Order:      left to right (small amounts before large)
Space:      O(target) — can't reduce further
```

**Problems:** Coin Change, Coin Change II, Perfect Squares.

**Go Template (Coin Change):**
```go
func coinChange(coins []int, amount int) int {
    // Step 1: dp[i] = min coins to make amount i
    // Step 2: dp[i] = min(dp[i-c] + 1) for each coin c where i-c >= 0
    // Step 3: dp[0] = 0 (zero coins for zero amount)
    // Step 4: left to right
    // Step 5: already 1D, can't reduce
    
    dp := make([]int, amount+1)
    for i := 1; i <= amount; i++ {
        dp[i] = math.MaxInt32 // NOT 0 — this is a min problem
    }
    
    for i := 1; i <= amount; i++ {
        for _, c := range coins {
            if c <= i && dp[i-c] != math.MaxInt32 {
                dp[i] = min(dp[i], dp[i-c]+1)
            }
        }
    }
    
    if dp[amount] == math.MaxInt32 {
        return -1
    }
    return dp[amount]
}
```

**Complexity:** O(amount * len(coins)) time, O(amount) space.

**Watch out:**
- **The #1 coin change bug:** initializing dp with 0 instead of infinity.
  For minimization, unreachable states must be infinity, not zero.
- dp array size is `amount+1`, not `amount`. dp[0] is your base case.
- Check `dp[i-c] != MaxInt32` before adding 1, or you'll overflow.
- Coin Change II (count ways) uses a different loop order to avoid counting
  permutations as distinct combinations.

---

### Pattern 4: Decision at Each Step (Take or Skip)

**Trigger:** At each element, you make a binary decision: take it (with some
constraint) or skip it. Classic "include/exclude" pattern.

**The Recipe Applied:**
```
State:      dp[i] = best result considering elements 0..i
Recurrence: dp[i] = max(dp[i-2] + val[i],   ← take i (skip i-1)
                        dp[i-1])              ← skip i
Base cases: dp[0] = val[0], dp[1] = max(val[0], val[1])
Order:      left to right
Space:      O(1) — only need two previous values
```

**Problems:** House Robber, House Robber II (circular), Delete and Earn,
Paint Houses.

**Go Template (House Robber):**
```go
func rob(nums []int) int {
    // Step 1: dp[i] = max money robbing from houses 0..i
    // Step 2: dp[i] = max(dp[i-2] + nums[i], dp[i-1])
    // Step 3: dp[0] = nums[0], dp[1] = max(nums[0], nums[1])
    // Step 4: left to right
    // Step 5: O(1) — two variables
    
    n := len(nums)
    if n == 1 {
        return nums[0]
    }
    
    prev2 := nums[0]
    prev1 := max(nums[0], nums[1])
    for i := 2; i < n; i++ {
        curr := max(prev2+nums[i], prev1)
        prev2 = prev1
        prev1 = curr
    }
    return prev1
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- **Common mistake:** thinking "skip" means dp[i-2]. No — skip means you
  inherit dp[i-1] (the best up to the previous house). Take means dp[i-2]
  + current value (because you can't take adjacent).
- House Robber II (circular): run House Robber twice — once on nums[0..n-2],
  once on nums[1..n-1], return the max.
- Delete and Earn reduces to House Robber after frequency counting.

---

### Pattern 5: State Machine DP

**Trigger:** You have multiple distinct states and transitions between them.
Often involves "cooldown", "transaction limits", or "holding vs not holding."

**The Recipe Applied:**
```
State:      Multiple dp values per index:
            hold[i]  = best profit on day i while holding stock
            free[i]  = best profit on day i while not holding stock
            cool[i]  = best profit on day i in cooldown
Recurrence: hold[i] = max(hold[i-1], free[i-1] - price[i])  (keep holding or buy)
            free[i] = max(free[i-1], cool[i-1])              (stay free or exit cooldown)
            cool[i] = hold[i-1] + price[i]                   (sell)
Base cases: hold[0] = -price[0], free[0] = 0, cool[0] = 0
Order:      left to right
Space:      O(1) — track current values of each state
```

**Problems:** Best Time to Buy/Sell Stock with Cooldown, with Transaction Fee,
with K Transactions.

**Go Template (Buy/Sell with Cooldown):**
```go
func maxProfit(prices []int) int {
    if len(prices) <= 1 {
        return 0
    }
    
    hold := -prices[0]  // bought on day 0
    free := 0           // never bought
    cool := 0           // impossible on day 0, but 0 is safe
    
    for i := 1; i < len(prices); i++ {
        newHold := max(hold, free-prices[i])
        newFree := max(free, cool)
        newCool := hold + prices[i]
        hold, free, cool = newHold, newFree, newCool
    }
    
    return max(free, cool)
}
```

**Complexity:** O(n) time, O(1) space.

**Watch out:**
- You must compute all new states from old states before overwriting. Use
  temp variables (newHold, newFree, newCool) or you corrupt the computation.
- Draw the state machine on paper first. Label edges with actions (buy, sell,
  hold, cooldown). The code writes itself from the diagram.
- With K transactions, add a transaction dimension: dp[i][k][0/1]. This
  becomes 2D DP but the state machine logic is the same.

---

## Common Interview Traps (Summary)

| Trap | What Goes Wrong | Fix |
|---|---|---|
| Coin change init | dp filled with 0, min() returns 0 for everything | Init dp[1..] = math.MaxInt32 |
| dp array size | dp has length n, but you need dp[0] as base case | Allocate n+1 |
| Word break inner loop | Substring slicing off by one, wrong loop direction | dp[i] = can we form s[0:i]; inner j checks s[j:i] |
| House robber "skip" | Writing dp[i] = dp[i-2] for skip instead of dp[i-1] | Skip means "best up to previous" = dp[i-1] |
| LIS over-optimization | Spending 20 min on O(n log n) when O(n^2) passes | Ask interviewer — O(n^2) is usually fine |
| State machine overwrites | Updating hold before computing free, using stale values | Compute all new states into temps, then assign |
| Not checking impossible | dp[i-c]+1 when dp[i-c] == MaxInt32 causes overflow | Guard with `dp[i-c] != math.MaxInt32` |

---

## Thought Process Walkthrough 1: House Robber

> **Interviewer:** "Given an array of non-negative integers representing money
> in each house, find the maximum you can rob without robbing adjacent houses."

### How to talk through this (say out loud):

**1. Clarify:** "So I can't rob two houses next to each other. I want to
maximize total money. No wrapping — it's a line, not a circle."

**2. Recognize the pattern:** "At each house I have a decision: rob it or
skip it. This is a take-or-skip DP."

**3. Apply the recipe out loud:**

> "Let me define dp[i] as the maximum money I can rob from houses 0 through i.
>
> For the recurrence: if I rob house i, I get nums[i] plus the best I could do
> through house i-2 (can't rob i-1). If I skip house i, I keep dp[i-1].
> So dp[i] = max(dp[i-2] + nums[i], dp[i-1]).
>
> Base cases: dp[0] = nums[0] (only one house, rob it). dp[1] = max(nums[0],
> nums[1]) (two houses, pick the richer one).
>
> I iterate left to right. And since I only look back two positions, I can
> use O(1) space with two variables."

**4. Code it:**

```go
func rob(nums []int) int {
    n := len(nums)
    if n == 1 {
        return nums[0]
    }
    prev2 := nums[0]
    prev1 := max(nums[0], nums[1])
    for i := 2; i < n; i++ {
        curr := max(prev2+nums[i], prev1)
        prev2 = prev1
        prev1 = curr
    }
    return prev1
}
```

**5. Verify:** "For [2,7,9,3,1]: prev2=2, prev1=7. i=2: max(2+9,7)=11.
i=3: max(7+3,11)=11. i=4: max(11+1,11)=12. Answer is 12. That's robbing
houses 0,2,4 → 2+9+1=12. Correct."

**6. Complexity:** "O(n) time, O(1) space."

Total speaking time: ~2 minutes. This is exactly the right pace.

---

## Thought Process Walkthrough 2: Coin Change

> **Interviewer:** "Given coin denominations and a target amount, find the
> minimum number of coins to make that amount. Return -1 if impossible."

### How to talk through this:

**1. Clarify:** "Unlimited coins of each denomination. I want minimum count,
not the actual coins. Coins are positive integers."

**2. Recognize the pattern:** "Minimum cost to reach a target with unlimited
choices from a set — this is unbounded choices DP."

**3. Why not greedy?** (Say this — it shows depth.) "Greedy would pick the
largest coin first, but that fails. With coins [1,3,4] and amount=6, greedy
gives 4+1+1=3 coins, but 3+3=2 coins is better."

**4. Apply the recipe:**

> "dp[i] = minimum number of coins to make amount i.
>
> Recurrence: for each coin c, if I use coin c, then dp[i] = dp[i-c] + 1.
> I take the minimum over all coins. dp[i] = min(dp[i-c] + 1) for all c
> where c <= i and dp[i-c] is reachable.
>
> Base case: dp[0] = 0. Zero coins to make zero.
>
> I need to initialize everything else to infinity since this is a minimum
> problem.
>
> I iterate from 1 to amount, left to right."

**5. Code it:**

```go
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := 1; i <= amount; i++ {
        dp[i] = math.MaxInt32
    }
    for i := 1; i <= amount; i++ {
        for _, c := range coins {
            if c <= i && dp[i-c] != math.MaxInt32 {
                dp[i] = min(dp[i], dp[i-c]+1)
            }
        }
    }
    if dp[amount] == math.MaxInt32 {
        return -1
    }
    return dp[amount]
}
```

**6. Verify:** "Coins [1,3,4], amount=6. dp[0]=0. dp[1]=1. dp[2]=2.
dp[3]=min(dp[2]+1, dp[0]+1)=1. dp[4]=min(dp[3]+1, dp[1]+1, dp[0]+1)=1.
dp[5]=min(dp[4]+1, dp[2]+1, dp[1]+1)=2. dp[6]=min(dp[5]+1, dp[3]+1,
dp[2]+1)=2. Answer: 2 (two coins of 3). Correct."

**7. Complexity:** "O(amount * len(coins)) time, O(amount) space."

---

## Practice Problems with Time Targets

### Warm-up (10 min each)
| # | Problem | Pattern | Target |
|---|---|---|---|
| 1 | [Climbing Stairs](https://leetcode.com/problems/climbing-stairs/) | Fibonacci | 5 min |
| 2 | [Min Cost Climbing Stairs](https://leetcode.com/problems/min-cost-climbing-stairs/) | Fibonacci | 8 min |

### Core (15 min each)
| # | Problem | Pattern | Target |
|---|---|---|---|
| 3 | [House Robber](https://leetcode.com/problems/house-robber/) | Take/Skip | 10 min |
| 4 | [Coin Change](https://leetcode.com/problems/coin-change/) | Unbounded | 12 min |
| 5 | [Maximum Subarray](https://leetcode.com/problems/maximum-subarray/) | Best Ending Here | 8 min |
| 6 | [Word Break](https://leetcode.com/problems/word-break/) | Best Ending Here | 15 min |
| 7 | [Decode Ways](https://leetcode.com/problems/decode-ways/) | Fibonacci (conditional) | 15 min |

### Stretch (20 min each)
| # | Problem | Pattern | Target |
|---|---|---|---|
| 8 | [House Robber II](https://leetcode.com/problems/house-robber-ii/) | Take/Skip + trick | 15 min |
| 9 | [Longest Increasing Subsequence](https://leetcode.com/problems/longest-increasing-subsequence/) | Best Ending Here | 15 min |
| 10 | [Best Time to Buy and Sell Stock with Cooldown](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-cooldown/) | State Machine | 20 min |

---

## Quick Drill: Recipe Speed Round

For each problem below, don't code — just write the 5-step recipe as fast
as you can. Target: 60 seconds each.

1. **Perfect Squares** — dp[i] = ? recurrence = ? base = ?
2. **Paint Houses (3 colors)** — dp[i] = ? recurrence = ? base = ?
3. **Maximum Product Subarray** — dp[i] = ? (hint: two arrays) recurrence = ? base = ?
4. **Delete and Earn** — reduce to which problem? dp[i] = ? recurrence = ?

<details>
<summary>Answers</summary>

1. dp[i] = min perfect squares summing to i. dp[i] = min(dp[i-j*j]+1) for
   all j where j*j <= i. dp[0] = 0.

2. dp[i][c] = min cost painting houses 0..i where house i is color c.
   dp[i][c] = cost[i][c] + min(dp[i-1][other colors]). dp[0][c] = cost[0][c].

3. dpMax[i] = max product ending at i, dpMin[i] = min product ending at i.
   dpMax[i] = max(nums[i], dpMax[i-1]*nums[i], dpMin[i-1]*nums[i]).
   dpMin[i] = min(nums[i], dpMax[i-1]*nums[i], dpMin[i-1]*nums[i]).
   Base: dpMax[0] = dpMin[0] = nums[0].

4. Reduce to House Robber. Bucket values by number, sum each bucket.
   dp[i] = max(dp[i-2] + earn[i], dp[i-1]) where earn[i] = i * count[i].

</details>

---

## Self-Assessment Checklist

After your practice session, score yourself honestly:

- [ ] **Recipe fluency:** Can I state the 5 steps for any new DP problem
      in under 90 seconds without looking at notes?
- [ ] **Pattern recognition:** Given a problem statement, can I identify
      which of the 5 patterns applies within 30 seconds?
- [ ] **Verbalization:** Did I practice saying the recipe OUT LOUD, not
      just thinking it? (This is the #1 differentiator in interviews.)
- [ ] **Trap awareness:** Can I list the 3 most common DP bugs from memory?
- [ ] **Coin change from scratch:** Can I write correct coin change in
      under 5 minutes without any reference?
- [ ] **House robber from scratch:** Can I write correct house robber in
      under 4 minutes without any reference?

### If you're struggling:
- Go back to the recipe. Every DP failure is a recipe failure.
- Solve House Robber 5 times from scratch until it's muscle memory.
- Then solve Coin Change 5 times. These two problems cover the two most
  important patterns (take/skip and unbounded choices).
- Once those are automatic, the rest are variations.

### Key takeaway:
**DP interviews are not about being clever. They're about having a systematic
process and communicating it clearly.** The recipe IS the process. Practice
it until you can recite it in your sleep.
