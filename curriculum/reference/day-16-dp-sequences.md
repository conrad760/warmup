# Day 16 — Dynamic Programming: Sequences

---

## 1. Curated Learning Resources

| # | Resource | Format | Why This One |
|---|----------|--------|-------------|
| 1 | [Back To Back SWE — Edit Distance (Dynamic Programming)](https://www.youtube.com/watch?v=MiqoA-yF-0M) | Video (18 min) | The single best visual walkthrough of edit distance. Draws the full DP table, fills every cell, and explains how each of the three operations maps to a specific neighbor cell. Pause at 8:00 and try to fill the next row yourself. |
| 2 | [NeetCode — Longest Common Subsequence](https://www.youtube.com/watch?v=Ua0GhsJSlWM) | Video (18 min) | Builds the LCS recurrence from scratch with a clear 2D grid visualization. Shows why a match extends the diagonal and a mismatch takes the max of two neighbors. Watch this first — LCS is the foundation for everything else today. |
| 3 | [VisuAlgo — Edit Distance Visualization](https://visualgo.net/en/editdistance) | Interactive | Step through the edit distance DP table cell by cell. Toggle between operations (insert, delete, replace) and see how the backtracking path recovers the actual edit script. Use during the review block and again while implementing. |
| 4 | [University of San Francisco — LCS Visualization](https://www.cs.usfca.edu/~galles/visualization/DPLCS.html) | Interactive | Enter two strings and watch the LCS table fill in real time with arrows showing dependencies. Then click "find LCS" to see the backtracking path highlighted. Essential for building intuition before you code recovery. |
| 5 | [Abdul Bari — Edit Distance (Dynamic Programming)](https://www.youtube.com/watch?v=We3YDTzNXEk) | Video (16 min) | Whiteboard-style, slower pace. Excellent at explaining the base cases (transforming to/from empty strings) and why each operation adds exactly 1 to the cost. Good if the Back To Back SWE video goes too fast. |
| 6 | [Tushar Roy — Longest Palindromic Subsequence](https://www.youtube.com/watch?v=_nCsPn7_OgI) | Video (12 min) | Walks through both approaches: the interval DP formulation `dp[i][j]` on substrings AND the reduction to LCS of string and its reverse. Seeing both approaches side by side cements the connection. |
| 7 | [MIT 6.006 — DP III: Parenthesization, Edit Distance, Knapsack](https://www.youtube.com/watch?v=ocZMDMZwhCY) | Video (52 min) | Erik Demaine covers edit distance as a suffix DP problem. The key insight: he frames every two-string DP as "what happens to the first character of each suffix?" This mental model generalizes to any sequence DP. Watch the edit distance section (starts ~12:00). |
| 8 | [Aditya Verma — Space Optimization in LCS](https://www.youtube.com/watch?v=hR3s9M-2hQo) | Video (15 min) | Dedicated walkthrough of the two-row trick for LCS. Explains why each row only depends on the previous row, how to swap rows, and the subtle diagonal issue when you try the same trick on edit distance. Watch after you implement the full table version. |

**Suggested order:** Resource 2 (NeetCode LCS) first to set the foundation, then 1 (Back To Back SWE edit distance) during the review block. Keep 4 (USF visualization) open while implementing LCS. Resources 6 and 7 during the palindrome section. Resource 8 after your full-table versions work.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 - 12:20 — Review (No Code)

| Time | Activity |
|------|----------|
| 12:00 - 12:06 | Read Section 3 (Core Concepts) end to end. Focus on the two-sequence DP pattern: the state is `dp[i][j]` involving prefixes of both strings. Internalize *why* a match extends the diagonal (both prefixes shrink by one character). |
| 12:06 - 12:12 | **Hand-trace the LCS table for "ABCBDAB" and "BDCAB" on paper.** Use the table in Section 5 as a reference, but try to fill each cell yourself before looking. Focus on the first 3 rows to get the pattern, then speed through the rest. Draw arrows showing which cell each value came from (diagonal for match, left or up for mismatch). |
| 12:12 - 12:17 | **Hand-trace the Edit Distance table for "kitten" and "sitting" on paper.** Fill at least the first 3 rows yourself. For each cell, write which operation was used (M/R/I/D). Pay attention to the base cases: `dp[i][0] = i` (delete everything), `dp[0][j] = j` (insert everything). |
| 12:17 - 12:20 | **Read the "Connection to real tools" paragraph** in Section 3. Understand: `diff` uses LCS, spell check uses edit distance, bioinformatics uses both. These are not toy problems — they power tools you use daily. |

### 12:20 - 1:20 — Implement

| Time | Activity |
|------|----------|
| 12:20 - 12:35 | **LCS — full 2D table.** Write `LCS(s1, s2 string) int`. Allocate the `(m+1) x (n+1)` table, fill it row by row. Test with: `("ABCBDAB", "BDCAB")` -> 4, `("", "abc")` -> 0, `("abc", "abc")` -> 3, `("abc", "def")` -> 0. Verify against your hand-traced table. |
| 12:35 - 12:47 | **Edit Distance — full 2D table.** Write `EditDistance(s1, s2 string) int`. The recurrence is nearly identical to LCS but with three operations on mismatch. Test with: `("kitten", "sitting")` -> 3, `("", "abc")` -> 3, `("abc", "abc")` -> 0, `("intention", "execution")` -> 5. |
| 12:47 - 12:57 | **Longest Palindromic Subsequence.** Write `LongestPalindromicSubseq(s string) int`. Use the LCS reduction: LPS(s) = LCS(s, reverse(s)). Test with: `("bbbab")` -> 4, `("cbbd")` -> 2, `("a")` -> 1, `("abcba")` -> 5 (already a palindrome). |
| 12:57 - 1:10 | **LCS Space Optimized.** Write `LCSOptimized(s1, s2 string) int`. Use only two rows (`prev` and `curr`), swapping after each outer iteration. Ensure you iterate over the shorter string in the inner loop for O(min(m,n)) space. Verify results match the full-table version on all test cases. |
| 1:10 - 1:20 | **Recover the actual LCS (Bonus).** Write `RecoverLCS(s1, s2 string) string`. Build the full 2D table, then backtrack from `dp[m][n]`: if characters match, prepend and go diagonal; else go in the direction of the larger neighbor. Test: `("ABCBDAB", "BDCAB")` -> `"BCAB"` (one valid LCS). |

### 1:20 - 1:50 — Solidify

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | **Edge cases for all problems.** Both strings empty. One string empty. Identical strings. Completely different strings (no common characters). Single-character strings. Unicode? (For today, stick with ASCII — Go's `s[i]` works fine.) |
| 1:30 - 1:40 | **Variant: try the interval DP formulation for palindromic subsequence.** Write it using `dp[i][j]` = LPS of `s[i..j]`. Fill by increasing substring length. Compare: is this easier or harder to think about than the LCS reduction? (Both are valid; knowing both deepens understanding.) |
| 1:40 - 1:50 | **Think about recovery for edit distance.** You don't need to implement it, but trace through the "kitten" -> "sitting" table and write down the sequence of operations. At each cell, you know whether you matched, replaced, inserted, or deleted. This is how `diff` tools produce their output. |

### 1:50 - 2:00 — Recap (From Memory)

Write down without looking:

1. The LCS recurrence (match case and mismatch case).
2. The edit distance recurrence (match case and the three operations).
3. Why space optimization works for LCS (what does each row depend on?).
4. Why recovering the LCS requires the full 2D table.
5. The relationship between LPS and LCS.
6. One gotcha you hit during implementation.

---

## 3. Core Concepts Deep Dive

### The Two-Sequence DP Pattern

Whenever you're comparing, aligning, or transforming two sequences, the state is almost always:

> `dp[i][j]` = answer for the prefix `s1[0..i-1]` and the prefix `s2[0..j-1]`

The dimensions of the table are `(len(s1)+1) x (len(s2)+1)`. The `+1` is for the empty prefix — the base case row and column.

At each cell `(i, j)`, you look at the characters `s1[i-1]` and `s2[j-1]` (the last characters of the current prefixes) and ask: "Do these characters interact? If so, how?" The answer determines the recurrence.

This pattern covers:
- **LCS**: do the characters match?
- **Edit distance**: are the characters the same, or do we need an operation?
- **Sequence alignment**: match, mismatch penalty, gap penalty
- **Shortest common supersequence**: combine LCS with reconstruction

The key mental model: you are processing both strings simultaneously, one character at a time, from left to right. At each step, you decide what to do with the current pair of characters.

---

### Longest Common Subsequence (LCS)

**State:** `dp[i][j]` = length of the LCS of `s1[0..i-1]` and `s2[0..j-1]`.

**Recurrence:**
```
if s1[i-1] == s2[j-1]:
    dp[i][j] = dp[i-1][j-1] + 1      // match: extend diagonal
else:
    dp[i][j] = max(dp[i-1][j],        // skip character from s1
                   dp[i][j-1])         // skip character from s2
```

**Base case:** `dp[0][j] = 0` and `dp[i][0] = 0` — the LCS of any string with the empty string is 0.

**Why the match case extends the diagonal:** If the last characters of both prefixes are the same, they *must* be part of any optimal LCS for these prefixes. Including them in the LCS means both prefixes shrink by one character (moving diagonally), and we add 1 to the length. Proof: if you could do better without using this match, you could also do at least as well with it — contradiction.

**Why the mismatch case takes the max of two neighbors:** If the characters don't match, at least one of them is NOT part of the LCS for these prefixes. So we consider two options:
- Drop the last character of `s1` (look at `dp[i-1][j]`).
- Drop the last character of `s2` (look at `dp[i][j-1]`).

We take the better of the two. Note: we DON'T consider `dp[i-1][j-1]` in the mismatch case because it's always `<= max(dp[i-1][j], dp[i][j-1])` — it can never be strictly better (adding a character to a prefix can only maintain or increase the LCS).

---

### Edit Distance (Levenshtein Distance)

**State:** `dp[i][j]` = minimum number of operations to convert `s1[0..i-1]` into `s2[0..j-1]`.

**Three operations and how each maps to a cell in the table:**

```
if s1[i-1] == s2[j-1]:
    dp[i][j] = dp[i-1][j-1]            // match: no operation needed
else:
    dp[i][j] = 1 + min(
        dp[i-1][j-1],                   // REPLACE s1[i-1] with s2[j-1]
        dp[i-1][j],                     // DELETE  s1[i-1]
        dp[i][j-1]                      // INSERT  s2[j-1] after position i
    )
```

**Base cases:**
- `dp[i][0] = i` — converting a string of length `i` to empty requires `i` deletions.
- `dp[0][j] = j` — converting empty to a string of length `j` requires `j` insertions.

**The beautiful symmetry:** Each of the three operations corresponds to moving to a specific neighboring cell:
- **Replace** = diagonal `(i-1, j-1)`: both prefixes shrink by 1 (we handled both characters).
- **Delete** = up `(i-1, j)`: s1 shrinks by 1 (we removed a character), s2 stays (still need to produce it).
- **Insert** = left `(i, j-1)`: s2 shrinks by 1 (we produced a character by inserting), s1 stays (its characters are still there).

This directional mapping is what makes the 2D table so powerful — the geometry of the table IS the algorithm.

---

### Recovering the Solution: Backtracking Through the DP Table

Computing the *length* of the LCS (or the edit distance) is step one. Recovering the *actual subsequence* (or the edit operations) is a separate skill — and it's the one that trips people up in interviews.

**How to recover the LCS:**

Start at `dp[m][n]` and walk backward:
1. If `s1[i-1] == s2[j-1]`: this character is in the LCS. Record it and move to `dp[i-1][j-1]`.
2. If `dp[i-1][j] >= dp[i][j-1]`: the LCS came from dropping `s1[i-1]`. Move to `dp[i-1][j]`.
3. Else: the LCS came from dropping `s2[j-1]`. Move to `dp[i][j-1]`.

Prepend (or collect and reverse) the recorded characters to get the actual LCS.

**Critical insight:** Recovery requires the FULL 2D table. The space-optimized version (two rows) doesn't preserve enough history to trace back. This is a common interview follow-up: "now find the actual LCS, not just the length."

**How to recover edit operations:**

Same idea — start at `dp[m][n]` and follow the path of optimal decisions back to `dp[0][0]`. At each step, you know which operation was used based on which neighbor was the source:
- Diagonal with match → no-op (characters are already equal).
- Diagonal with `+1` → replace.
- Up with `+1` → delete from s1.
- Left with `+1` → insert into s1.

---

### Space Optimization: The Two-Row Trick

**Why it works for LCS:** Look at the recurrence — `dp[i][j]` depends only on three cells: `dp[i-1][j-1]`, `dp[i-1][j]`, and `dp[i][j-1]`. All of these are either in the current row or the previous row. So we only need two rows.

```go
prev := make([]int, n+1)
curr := make([]int, n+1)
for i := 1; i <= m; i++ {
    for j := 1; j <= n; j++ {
        if s1[i-1] == s2[j-1] {
            curr[j] = prev[j-1] + 1
        } else {
            curr[j] = max(prev[j], curr[j-1])
        }
    }
    prev, curr = curr, prev  // swap (reuse old prev as new curr)
    // clear curr for next iteration? No — we'll overwrite every cell.
}
```

Space: O(min(m, n)) if you iterate over the shorter string in the inner loop.

**The diagonal issue in edit distance:** Edit distance also needs `dp[i-1][j-1]` — the diagonal. When you're filling `curr[j]` left to right, you've already overwritten `curr[j-1]` (which was `prev[j-1]` from the perspective of the previous iteration). You need to save the diagonal value before overwriting:

```go
for i := 1; i <= m; i++ {
    prev_diag := prev[0]    // save dp[i-1][0] before we start
    curr[0] = i             // base case for edit distance
    for j := 1; j <= n; j++ {
        temp := curr[j]     // this will be next iteration's diagonal
        if s1[i-1] == s2[j-1] {
            curr[j] = prev_diag
        } else {
            curr[j] = 1 + min(prev_diag, prev[j], curr[j-1])
        }
        prev_diag = temp
    }
    // ... swap prev and curr
}
```

Or more commonly, keep `prev` and `curr` as separate rows and handle the diagonal explicitly.

---

### Longest Palindromic Subsequence (LPS)

**Approach 1: Reduce to LCS.** The longest palindromic subsequence of `s` is exactly the LCS of `s` and `reverse(s)`. This works because a palindrome reads the same forwards and backwards — so any common subsequence between `s` and its reverse is a palindrome, and the longest such is the LPS.

```go
func LongestPalindromicSubseq(s string) int {
    rev := reverse(s)
    return LCS(s, rev)
}
```

Time: O(n^2). Space: O(n^2), or O(n) with the two-row trick.

**Approach 2: Interval DP.** Define `dp[i][j]` = length of the LPS of `s[i..j]`.

```
if s[i] == s[j]:
    dp[i][j] = dp[i+1][j-1] + 2    // both endpoints are in the palindrome
else:
    dp[i][j] = max(dp[i+1][j],      // drop left endpoint
                   dp[i][j-1])       // drop right endpoint
```

Base cases: `dp[i][i] = 1` (single character is a palindrome of length 1). `dp[i][i-1] = 0` (empty substring, needed when `j = i+1` and characters don't match).

**Traversal order:** Fill by increasing substring length (`len` from 2 to `n`), setting `i = 0..n-len` and `j = i+len-1`. This ensures `dp[i+1][j-1]`, `dp[i+1][j]`, and `dp[i][j-1]` are already computed.

Both approaches give the same answer. The LCS reduction is easier to code (you already have LCS). The interval DP approach builds different intuition and is more generalizable to other interval DP problems.

---

### Connection to Real Tools

These are not academic exercises. Sequence DP powers tools you use every day:

- **`diff` and `git diff`:** The core algorithm is LCS. `diff` finds the longest common subsequence of lines between two files. Lines NOT in the LCS are the additions and deletions you see in the diff output. The `patience diff` algorithm in git is a refinement that uses LCS on unique matching lines first.

- **Spell check and autocorrect:** Edit distance. When you type "teh", your spell checker computes the edit distance to every word in the dictionary (using optimized methods like BK-trees) and suggests words within distance 1-2. `"the"` is distance 1 (transpose, or delete + insert).

- **DNA sequence alignment (bioinformatics):** A generalized edit distance with different costs for substitutions, insertions, and deletions (called the Needleman-Wunsch algorithm for global alignment and Smith-Waterman for local alignment). The BLAST tool used in genomics is an optimized version of sequence DP.

- **Merge conflict resolution:** The 3-way merge used in git is built on top of LCS/diff of three versions of a file.

---

## 4. Implementation Checklist

### Function Signatures

```go
package dp

// --- Longest Common Subsequence ---

// LCS returns the length of the longest common subsequence of s1 and s2
// using a full (m+1) x (n+1) DP table.
// Time: O(m*n), Space: O(m*n)
func LCS(s1, s2 string) int

// LCSOptimized returns the length of the longest common subsequence
// using only two rows for O(min(m,n)) space.
// Time: O(m*n), Space: O(min(m,n))
func LCSOptimized(s1, s2 string) int

// RecoverLCS returns an actual longest common subsequence string
// (not just its length) by backtracking through the full DP table.
// Time: O(m*n), Space: O(m*n)
func RecoverLCS(s1, s2 string) string

// --- Edit Distance ---

// EditDistance returns the minimum number of operations (insert, delete,
// replace) to convert s1 into s2.
// Time: O(m*n), Space: O(m*n)
func EditDistance(s1, s2 string) int

// --- Longest Palindromic Subsequence ---

// LongestPalindromicSubseq returns the length of the longest palindromic
// subsequence of s, using the LCS-of-string-and-reverse reduction.
// Time: O(n^2), Space: O(n^2)
func LongestPalindromicSubseq(s string) int
```

### Test Cases & Edge Cases

**LCS:**

| Input | Expected | Notes |
|-------|----------|-------|
| `("ABCBDAB", "BDCAB")` | `4` | LCS is "BCAB" or "BDAB" |
| `("", "abc")` | `0` | Empty string |
| `("abc", "")` | `0` | Empty string |
| `("", "")` | `0` | Both empty |
| `("abc", "abc")` | `3` | Identical strings |
| `("abc", "def")` | `0` | No common characters |
| `("AGGTAB", "GXTXAYB")` | `4` | LCS is "GTAB" |
| `("a", "a")` | `1` | Single matching character |
| `("a", "b")` | `0` | Single non-matching character |

**LCSOptimized:** Same test cases as LCS — verify identical results.

**RecoverLCS:**

| Input | Expected (one valid LCS) | Notes |
|-------|--------------------------|-------|
| `("ABCBDAB", "BDCAB")` | `"BCAB"` or `"BDAB"` | Multiple valid LCS exist |
| `("", "abc")` | `""` | Empty string |
| `("abc", "abc")` | `"abc"` | Identical |
| `("abc", "def")` | `""` | No common chars |

**EditDistance:**

| Input | Expected | Notes |
|-------|----------|-------|
| `("kitten", "sitting")` | `3` | k->s (replace), e->i (replace), ""->g (insert) |
| `("", "abc")` | `3` | 3 insertions |
| `("abc", "")` | `3` | 3 deletions |
| `("", "")` | `0` | Both empty |
| `("abc", "abc")` | `0` | Identical |
| `("intention", "execution")` | `5` | Classic example |
| `("a", "b")` | `1` | Single replace |
| `("ab", "a")` | `1` | Single delete |
| `("a", "ab")` | `1` | Single insert |
| `("horse", "ros")` | `3` | LeetCode #72 example |

**LongestPalindromicSubseq:**

| Input | Expected | Notes |
|-------|----------|-------|
| `"bbbab"` | `4` | "bbbb" |
| `"cbbd"` | `2` | "bb" |
| `"a"` | `1` | Single character |
| `"abcba"` | `5` | Already a palindrome |
| `"abcd"` | `1` | No palindrome longer than 1 |
| `"aab"` | `2` | "aa" |
| `""` | `0` | Empty string |
| `"character"` | `5` | "carac" |

---

## 5. The 2D DP Table Walkthrough

### LCS of "ABCBDAB" and "BDCAB"

The full `(8 x 6)` table. `dp[i][j]` = LCS length of `s1[0..i-1]` and `s2[0..j-1]`.

Arrows show which cell each value came from: `↖` = diagonal (match), `←` = left, `↑` = up.

```
           ""    B     D     C     A     B
       ┌─────┬─────┬─────┬─────┬─────┬─────┐
  ""   │  0  │  0  │  0  │  0  │  0  │  0  │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  A    │  0  │  0← │  0← │  0← │  1↖ │  1← │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  B    │  0  │  1↖ │  1← │  1← │  1↑ │  2↖ │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  C    │  0  │  1↑ │  1↑ │  2↖ │  2← │  2↑ │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  B    │  0  │  1↖ │  1↑ │  2↑ │  2↑ │  3↖ │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  D    │  0  │  1↑ │  2↖ │  2↑ │  2↑ │  3↑ │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  A    │  0  │  1↑ │  2↑ │  2↑ │  3↖ │  3↑ │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  B    │  0  │  1↖ │  2↑ │  2↑ │  3↑ │  4↖ │
       └─────┴─────┴─────┴─────┴─────┴─────┘

Answer: dp[7][5] = 4
```

**Cell-by-cell computation (selected cells to build intuition):**

```
dp[1][1]: s1[0]='A', s2[0]='B' → mismatch → max(dp[0][1], dp[1][0]) = max(0,0) = 0
dp[1][4]: s1[0]='A', s2[3]='A' → match!   → dp[0][3] + 1 = 0 + 1 = 1
dp[2][1]: s1[1]='B', s2[0]='B' → match!   → dp[1][0] + 1 = 0 + 1 = 1
dp[2][5]: s1[1]='B', s2[4]='B' → match!   → dp[1][4] + 1 = 1 + 1 = 2
dp[4][5]: s1[3]='B', s2[4]='B' → match!   → dp[3][4] + 1 = 2 + 1 = 3
dp[7][5]: s1[6]='B', s2[4]='B' → match!   → dp[6][4] + 1 = 3 + 1 = 4
```

**Backtracking to recover the LCS (see Section 6 for the visual path):**

```
Start at dp[7][5] = 4
  s1[6]='B' == s2[4]='B' → match! Record 'B', move to dp[6][4] = 3
  s1[5]='A' == s2[3]='A' → match! Record 'A', move to dp[5][3] = 2
  s1[4]='D' != s2[2]='C' → dp[4][3]=2 >= dp[5][2]=2, move up to dp[4][3] = 2
  s1[3]='B' != s2[2]='C' → dp[3][3]=2 > dp[4][2]=1, move up to dp[3][3] = 2
  s1[2]='C' == s2[2]='C' → match! Record 'C', move to dp[2][2] = 1
  s1[1]='B' != s2[1]='D' → dp[1][2]=0 < dp[2][1]=1, move left... wait.
  Actually: dp[1][2]=0, dp[2][1]=1 → dp[2][1] is larger, move up to dp[1][1] = 0
  s1[0]='A' != s2[0]='B' → dp[0][1]=0 = dp[1][0]=0, move up to dp[0][1] = 0
  i=0, done.

  Collected (reverse order): B, A, C, B → reversed: "BCAB" ✓
```

---

### Edit Distance of "kitten" and "sitting"

The full `(7 x 8)` table. `dp[i][j]` = min operations to convert `s1[0..i-1]` to `s2[0..j-1]`.

Each cell shows the value and the operation that produced it:
- `M` = Match (diagonal, cost 0)
- `R` = Replace (diagonal, cost 1)
- `D` = Delete from s1 (up, cost 1)
- `I` = Insert into s1 (left, cost 1)

```
              ""     s      i      t      t      i      n      g
          ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┐
    ""    │  0   │  1 I │  2 I │  3 I │  4 I │  5 I │  6 I │  7 I │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    k     │  1 D │  1 R │  2 R │  3 R │  4 R │  5 R │  6 R │  7 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    i     │  2 D │  2 R │  1 M │  2 I │  3 I │  4 M │  5 I │  6 I │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    t     │  3 D │  3 R │  2 D │  1 M │  2 M │  3 I │  4 R │  5 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    t     │  4 D │  4 R │  3 D │  2 D │  1 M │  2 I │  3 R │  4 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    e     │  5 D │  5 R │  4 D │  3 D │  2 D │  2 R │  3 R │  4 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    n     │  6 D │  6 R │  5 D │  4 D │  3 D │  3 D │  2 M │  3 I │
          └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┘

Answer: dp[6][7] = 3
```

**Cell-by-cell computation (key cells):**

```
Base cases:
  dp[0][j] = j  (insert j characters to build s2 prefix from empty)
  dp[i][0] = i  (delete i characters from s1 prefix)

dp[1][1]: 'k' vs 's' → mismatch → 1 + min(dp[0][0], dp[0][1], dp[1][0])
          = 1 + min(0, 1, 1) = 1 (Replace k→s)

dp[2][2]: 'i' vs 'i' → match! → dp[1][1] = 1 (Match, no cost)

dp[3][3]: 't' vs 't' → match! → dp[2][2] = 1 (Match, no cost)

dp[4][4]: 't' vs 't' → match! → dp[3][3] = 1 (Match, no cost)

dp[5][5]: 'e' vs 'i' → mismatch → 1 + min(dp[4][4], dp[4][5], dp[5][4])
          = 1 + min(1, 2, 2) = 2 (Replace e→i)

dp[6][6]: 'n' vs 'n' → match! → dp[5][5] = 2 (Match, no cost)

dp[6][7]: 'n' vs 'g' → mismatch → 1 + min(dp[5][6], dp[5][7], dp[6][6])
          = 1 + min(3, 4, 2) = 3 (Insert g)
```

**The actual edit sequence (backtracking from dp[6][7]):**

```
dp[6][7]=3: Insert 'g'     → move to dp[6][6]=2
dp[6][6]=2: Match 'n'='n'  → move to dp[5][5]=2
dp[5][5]=2: Replace 'e'→'i'→ move to dp[4][4]=1
dp[4][4]=1: Match 't'='t'  → move to dp[3][3]=1
dp[3][3]=1: Match 't'='t'  → move to dp[2][2]=1
dp[2][2]=1: Match 'i'='i'  → move to dp[1][1]=1
dp[1][1]=1: Replace 'k'→'s'→ move to dp[0][0]=0

Operations (in order):
  1. Replace 'k' with 's'  → "sitten"
  2. Replace 'e' with 'i'  → "sittin"
  3. Insert 'g' at end     → "sitting"

Total: 3 operations ✓
```

---

## 6. Visual Diagrams

### LCS Table with Backtracking Path

The path from `dp[7][5]` back to `dp[0][0]` to recover the LCS. Cells on the path are marked with `*`. Diagonal moves on matches (marked `↖`) contribute characters to the LCS.

```
           ""    B     D     C     A     B
       ┌─────┬─────┬─────┬─────┬─────┬─────┐
  ""   │ *0  │  0  │  0  │  0  │  0  │  0  │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  A    │  0  │ *0  │  0  │  0  │  1  │  1  │
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  B    │  0  │  1  │ *1  │  1  │  1  │  2  │  ← 'B' (↖ match at dp[2][1])
       ├─────┼─────┼─────┼─────┼─────┼─────┤         was from dp[1][1]→dp[2][1]?
  C    │  0  │  1  │  1  │ *2  │  2  │  2  │  ← 'C' (↖ match)
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  B    │  0  │  1  │  1  │ *2  │  2  │  3  │  ↑ (no match, move up)
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  D    │  0  │  1  │  2  │ *2  │  2  │  3  │  ↑ (no match, move up)
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  A    │  0  │  1  │  2  │  2  │ *3  │  3  │  ← 'A' (↖ match)
       ├─────┼─────┼─────┼─────┼─────┼─────┤
  B    │  0  │  1  │  2  │  2  │  3  │ *4  │  ← 'B' (↖ match)
       └─────┴─────┴─────┴─────┴─────┴─────┘

  Path: (7,5)↖ → (6,4)↖ → (5,3)↑ → (4,3)↑ → (3,3)↖ → (2,2)↑ → (1,1)↑ → (0,1)
                                                                     ↑ could also go
                                                                       to (0,0) or (1,0)
  Characters collected at ↖ moves: B, A, C, B
  Reversed: B C A B → "BCAB"

  Note: There are other valid LCS paths (e.g., "BDAB"). The specific
  path depends on tie-breaking when dp[i-1][j] == dp[i][j-1].
```

---

### Edit Distance Table with Operations

Every cell annotated with the operation used: M=Match, R=Replace, D=Delete, I=Insert.

```
              ""     s      i      t      t      i      n      g
          ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┐
    ""    │  0   │  1 I │  2 I │  3 I │  4 I │  5 I │  6 I │  7 I │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    k     │  1 D │  1 R │  2 R │  3 R │  4 R │  5 R │  6 R │  7 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    i     │  2 D │  2 R │  1 M │  2 I │  3 I │  4 M │  5 I │  6 I │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    t     │  3 D │  3 R │  2 D │  1 M │  2 M │  3 I │  4 R │  5 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    t     │  4 D │  4 R │  3 D │  2 D │  1 M │  2 I │  3 R │  4 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    e     │  5 D │  5 R │  4 D │  3 D │  2 D │  2 R │  3 R │  4 R │
          ├──────┼──────┼──────┼──────┼──────┼──────┼──────┼──────┤
    n     │  6 D │  6 R │  5 D │  4 D │  3 D │  3 D │  2 M │  3 I │
          └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┘

  Backtracking path (optimal alignment):

    (6,7) ← I   Insert 'g'
      ↑
    (6,6) ← M   Match 'n'='n'
      ↖
    (5,5) ← R   Replace 'e'→'i'
      ↖
    (4,4) ← M   Match 't'='t'
      ↖
    (3,3) ← M   Match 't'='t'
      ↖
    (2,2) ← M   Match 'i'='i'
      ↖
    (1,1) ← R   Replace 'k'→'s'
      ↖
    (0,0)  done

  Alignment:
    s1:  k  i  t  t  e  n  _
    s2:  s  i  t  t  i  n  g
    op:  R  M  M  M  R  M  I

  3 operations total (2 replacements + 1 insertion)
```

---

### LCS Dependency Diagram

This shows the three cells that `dp[i][j]` depends on, and which cell is used in which case.

```
  When s1[i-1] == s2[j-1] (match):

              j-1      j
           ┌───────┬───────┐
    i-1    │dp[i-1]│       │
           │ [j-1] │       │      dp[i][j] = dp[i-1][j-1] + 1
           │   ↖   │       │      (diagonal: both prefixes shrink by 1)
           ├───────┼───────┤
    i      │       │dp[i]  │
           │       │ [j]   │
           └───────┴───────┘

  When s1[i-1] != s2[j-1] (mismatch):

              j-1      j
           ┌───────┬───────┐
    i-1    │  (x)  │dp[i-1]│
           │       │ [j]   │      dp[i][j] = max(dp[i-1][j],    ← "skip s1[i-1]"
           │       │   ↑   │                      dp[i][j-1])   ← "skip s2[j-1]"
           ├───────┼───────┤
    i      │dp[i]  │dp[i]  │
           │ [j-1] │ [j]   │      (x) dp[i-1][j-1] is NOT considered because
           │   ←   │       │      it's always ≤ max of the other two
           └───────┴───────┘


  Edit Distance — all three neighbors matter on mismatch:

              j-1      j
           ┌───────┬───────┐
    i-1    │dp[i-1]│dp[i-1]│
           │ [j-1] │ [j]   │      dp[i][j] = 1 + min(dp[i-1][j-1],  ← Replace
           │   ↖R  │   ↑D  │                        dp[i-1][j],      ← Delete
           ├───────┼───────┤                         dp[i][j-1])     ← Insert
    i      │dp[i]  │dp[i]  │
           │ [j-1] │ [j]   │
           │   ←I  │       │
           └───────┴───────┘
```

---

## 7. Self-Assessment

Answer these without looking at the notes. If you can't, revisit the relevant section.

### Question 1
**Why does recovering the LCS require the full 2D table (not the space-optimized version)?**

<details>
<summary>Answer</summary>

The space-optimized version only keeps two rows (current and previous). When you finish filling the table, you have the final answer `dp[m][n]` but you've discarded all rows except the last two. Recovery requires backtracking from `dp[m][n]` all the way to `dp[0][0]`, visiting cells across ALL rows. Without the full table, you can't trace backward — you don't know which cell each value came from.

The same issue applies to edit distance recovery: you need the full table to reconstruct the sequence of operations.

If you need both space optimization AND recovery, you'd need the Hirschberg algorithm (divide-and-conquer on the table), which achieves O(n) space and O(mn) time but is significantly more complex.
</details>

### Question 2
**What's the relationship between LCS length and edit distance? Can you derive one from the other?**

<details>
<summary>Answer</summary>

For two strings of length `m` and `n` with LCS length `L`:
- The minimum number of deletions + insertions (no replacements allowed) to convert one string to the other is: `(m - L) + (n - L) = m + n - 2L`.

This is because:
- You must delete the `m - L` characters from `s1` that aren't in the LCS.
- You must insert the `n - L` characters from `s2` that aren't in the LCS.

With replacements allowed (standard Levenshtein distance), the relationship is more complex. The edit distance is always `≤ m + n - 2L` because replacements can save operations (one replacement instead of one deletion + one insertion). But the edit distance is always `≥ max(m, n) - L` because you need at least that many operations.

In the special case where `m == n`, the edit distance (with replacements) equals `m - L` if no insertions/deletions are cheaper than replacements. But this doesn't hold in general.
</details>

### Question 3
**In the LCS recurrence, when the characters DON'T match, why don't we also consider `dp[i-1][j-1]`? (i.e., why is it `max(dp[i-1][j], dp[i][j-1])` and not `max(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])`?)**

<details>
<summary>Answer</summary>

Because `dp[i-1][j-1] <= dp[i-1][j]` and `dp[i-1][j-1] <= dp[i][j-1]`. Always.

Think about it: `dp[i-1][j]` considers the prefixes `s1[0..i-2]` and `s2[0..j-1]`. `dp[i-1][j-1]` considers `s1[0..i-2]` and `s2[0..j-2]`. The second string prefix is shorter in `dp[i-1][j-1]`, so its LCS can only be equal or shorter. Adding more characters to a prefix can never decrease the LCS.

So including `dp[i-1][j-1]` in the max wouldn't change the result — it's dominated by both of the other two values. Leaving it out is a simplification, not an approximation.

This is a subtle but important point that shows you truly understand the recurrence, not just memorize it.
</details>

### Question 4
**Why does the space-optimized edit distance need special handling for the diagonal, while LCS does not?**

<details>
<summary>Answer</summary>

Both LCS and edit distance use `dp[i-1][j-1]` (the diagonal). The difference is in when the current row overwrites information you still need.

In LCS with two rows, when you compute `curr[j]`, the diagonal `dp[i-1][j-1]` is `prev[j-1]`. This value is in the *previous* row and hasn't been touched yet — `prev` is read-only during the inner loop. So no special handling is needed.

But if you try to do edit distance with a *single row* (overwriting in place instead of using two rows), then by the time you compute `dp[j]`, the value `dp[j-1]` has already been overwritten with the current row's value. The old `dp[j-1]` (which is the diagonal `dp[i-1][j-1]`) is gone. You need to save it in a temporary variable before overwriting.

With two explicit rows (`prev` and `curr`), this issue doesn't arise for either problem — but many optimized implementations use a single array, and that's where the diagonal becomes tricky.
</details>

### Question 5
**You're given a string and asked for the minimum number of characters to insert to make it a palindrome. How does this relate to LPS?**

<details>
<summary>Answer</summary>

The minimum insertions to make a string of length `n` into a palindrome is:

`n - LPS(s)`

where `LPS(s)` is the length of the longest palindromic subsequence.

Why? The LPS is already a palindrome — it's the longest part of the string that "works." The remaining `n - LPS(s)` characters are the ones that aren't part of any palindromic subsequence, and each one needs exactly one insertion (its mirror character) to become part of a palindrome.

Example: `s = "abcd"`, LPS = 1 (any single character). Minimum insertions = 4 - 1 = 3. Result: `"abcdcba"` or `"dcbabcd"`, etc.

Example: `s = "aebcbda"`, LPS = 5 (`"abcba"`). Minimum insertions = 7 - 5 = 2.

This is a classic reduction that shows up in interviews — recognizing the connection to LPS is the entire solve.
</details>
