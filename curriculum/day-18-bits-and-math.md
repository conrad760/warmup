# Day 18: Bit Manipulation & Math

> **Time:** 2 hours | **Level:** Refresher | **Language:** Go
>
> These topics are lower frequency but fast to review. When they appear,
> you either know the trick or you don't — there's no way to derive XOR
> properties under pressure. Drill the patterns, then move on.

---

## Pattern Catalog

### Pattern 1: XOR to Find Unique — Single Number

**Trigger:** "Every element appears twice except one. Find the unique element."
Any problem where pairing/cancellation is involved.

```go
func singleNumber(nums []int) int {
    result := 0
    for _, n := range nums {
        result ^= n
    }
    return result
}
```

**Complexity:** O(n) time, O(1) space.

**Why it works:** XOR is self-inverse: `a ^ a = 0` and `a ^ 0 = a`. All
paired elements cancel, leaving only the unique one.

**Watch out:**
- Only works when exactly one element is unique and all others appear
  exactly twice. If elements appear three times, you need a different
  approach (bit counting mod 3).
- Order doesn't matter — XOR is commutative and associative.

---

### Pattern 2: Kernighan's Bit Count — Number of 1 Bits

**Trigger:** "Count the number of set bits." "Hamming weight." Any problem
that needs popcount.

```go
func hammingWeight(n uint32) int {
    count := 0
    for n != 0 {
        n &= n - 1 // clears the lowest set bit
        count++
    }
    return count
}
```

**Complexity:** O(k) where k = number of set bits. Worst case O(32).

**Why it works:** `n - 1` flips the lowest set bit and all bits below it.
ANDing with `n` clears exactly that lowest set bit. Each iteration
removes one set bit, so the loop runs exactly k times.

**Watch out:**
- Go also has `math/bits.OnesCount(uint(n))` — mention it to show you
  know the stdlib, but implement Kernighan's to show you understand it.
- The parameter type matters: `uint32` vs `int`. The problem usually
  specifies 32-bit unsigned.

---

### Pattern 3: Bit Masking for Subsets

**Trigger:** "Generate all subsets" as an alternative to backtracking.
Works when n is small (n <= 20).

```go
func subsets(nums []int) [][]int {
    n := len(nums)
    total := 1 << n // 2^n
    result := make([][]int, 0, total)

    for mask := 0; mask < total; mask++ {
        subset := []int{}
        for i := 0; i < n; i++ {
            if mask&(1<<i) != 0 {
                subset = append(subset, nums[i])
            }
        }
        result = append(result, subset)
    }
    return result
}
```

**Complexity:** O(n * 2^n) time and space.

**Why it works:** Each integer from 0 to 2^n-1 is a binary mask where
bit i indicates whether `nums[i]` is included. This bijection between
integers and subsets means iterating integers gives every subset exactly
once.

**Watch out:**
- `1 << n` must use an integer type wide enough. For n=31 with `int32`,
  you overflow. Go's `int` is 64-bit on modern platforms, so usually fine.
- Bitmask iteration gives subsets in a different order than backtracking.
  If the problem needs lexicographic order, sort or use backtracking.

---

### Pattern 4: Power of 2 Check

**Trigger:** "Is n a power of 2?"

```go
func isPowerOfTwo(n int) bool {
    return n > 0 && n&(n-1) == 0
}
```

**Complexity:** O(1).

**Why it works:** A power of 2 has exactly one set bit. `n & (n-1)` clears
the lowest set bit. If the result is 0, there was only one set bit.

**Watch out:**
- `n > 0` is essential. `n = 0` would pass `n & (n-1) == 0` but 0 is not
  a power of 2.
- In Go, operator precedence: `&` binds tighter than `==`, so
  `n&(n-1) == 0` parses correctly, but adding parens is clearer.

---

### Pattern 5: Missing Number — XOR or Sum

**Trigger:** "Array contains n distinct numbers from 0 to n. Find the
missing one."

```go
// XOR approach
func missingNumber(nums []int) int {
    xor := len(nums)
    for i, v := range nums {
        xor ^= i ^ v
    }
    return xor
}

// Sum approach (simpler to explain)
func missingNumberSum(nums []int) int {
    n := len(nums)
    expected := n * (n + 1) / 2
    actual := 0
    for _, v := range nums {
        actual += v
    }
    return expected - actual
}
```

**Complexity:** O(n) time, O(1) space for both.

**Why XOR works:** XOR 0..n with all array elements. Every number present
in both cancels. The missing number has no pair, so it remains.

**Why sum works:** Gauss's formula gives the expected sum. Subtract the
actual sum. The difference is the missing number.

**Watch out:**
- Sum approach can overflow for very large n with 32-bit integers. Go's
  `int` is 64-bit, so this is rarely a problem, but mention it.
- XOR approach avoids overflow entirely — a good reason to prefer it.

---

### Pattern 6: GCD / Euclidean Algorithm

**Trigger:** "Simplify a fraction." "GCD of two numbers." Any problem
involving greatest common divisor or least common multiple.

```go
func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

func lcm(a, b int) int {
    return a / gcd(a, b) * b // divide first to avoid overflow
}
```

**Complexity:** O(log(min(a, b))).

**Watch out:**
- For LCM, compute `a / gcd(a, b) * b` not `a * b / gcd(a, b)`. The
  multiplication `a * b` can overflow even on 64-bit.
- GCD with negative numbers: take absolute values first, or the result
  may be negative depending on inputs.
- Go does not have a built-in GCD (unlike Python). You must write it.

---

### Pattern 7: Sieve of Eratosthenes — Count Primes

**Trigger:** "Count the number of primes less than n." Any bulk prime
generation problem.

```go
func countPrimes(n int) int {
    if n < 2 {
        return 0
    }
    isComposite := make([]bool, n)
    count := 0

    for i := 2; i < n; i++ {
        if isComposite[i] {
            continue
        }
        count++
        // Start crossing out at i*i, not 2*i
        for j := i * i; j < n; j += i {
            isComposite[j] = true
        }
    }
    return count
}
```

**Complexity:** O(n log log n) time, O(n) space.

**Watch out:**
- Start the inner loop at `i * i`, not `2 * i`. All smaller multiples
  have already been marked by smaller primes. Starting at `2 * i` is
  correct but wastes time — interviewers expect the `i * i` optimization.
- `i * i` can overflow for large `i` on 32-bit. Check bounds or use
  `j := i * i` only when `i * i < n`.
- The problem usually says "less than n", not "less than or equal to n".
  Read carefully.

---

### Pattern 8: Modular Arithmetic

**Trigger:** "Return the answer modulo 10^9 + 7." Common in counting and
DP problems.

```go
const MOD = 1_000_000_007

// Apply mod after every addition and multiplication
func addMod(a, b int) int {
    return (a + b) % MOD
}

func mulMod(a, b int) int {
    return (a % MOD) * (b % MOD) % MOD
}

// Modular exponentiation: base^exp % mod
func powMod(base, exp, mod int) int {
    result := 1
    base %= mod
    for exp > 0 {
        if exp%2 == 1 {
            result = result * base % mod
        }
        exp /= 2
        base = base * base % mod
    }
    return result
}
```

**Key rules:**
- `(a + b) % m = ((a % m) + (b % m)) % m`
- `(a * b) % m = ((a % m) * (b % m)) % m`
- Subtraction: `((a - b) % m + m) % m` to avoid negative results.
- Division: multiply by the modular inverse. `a / b % m = a * powMod(b, m-2, m) % m`
  (only when m is prime, via Fermat's little theorem).

**Watch out:**
- Apply mod at every step, not just at the end. Intermediate products
  overflow even 64-bit integers.
- Go's `%` can return negative values for negative operands. Always add
  `+ MOD` before the final `% MOD` when subtraction is involved.

---

## Decision Framework

| Signal in the problem | Pattern to apply |
|----|-----|
| "Find the one unique element" / elements cancel in pairs | XOR (Pattern 1) |
| "Count set bits" / "hamming distance" / "hamming weight" | Kernighan's (Pattern 2) |
| "Generate all subsets" (n <= 20, alternative to backtracking) | Bitmask iteration (Pattern 3) |
| "Is it a power of 2?" | `n & (n-1)` (Pattern 4) |
| "Missing number in 0..n" | XOR or sum (Pattern 5) |
| "GCD" / "simplify fraction" / "LCM" | Euclidean (Pattern 6) |
| "Count primes up to n" / "is n prime" (bulk) | Sieve (Pattern 7) |
| "Answer modulo 10^9+7" | Modular arithmetic (Pattern 8) |

If none of these triggers match, you probably don't need bit manipulation
or math tricks — look at other patterns first.

---

## Quick Reference Card: Bit Operations in Go

All operations assume `n` is an unsigned or non-negative integer and bit
positions are 0-indexed from the right.

```
Operation            Go expression             Example (n = 0b1010)
─────────────────────────────────────────────────────────────────────
Check bit i          n & (1 << i) != 0          bit 1: 1010 & 0010 → != 0 ✓
Set bit i            n | (1 << i)               set 2: 1010 | 0100 → 1110
Clear bit i          n &^ (1 << i)              clr 3: 1010 &^ 1000 → 0010
Toggle bit i         n ^ (1 << i)               tog 1: 1010 ^ 0010 → 1000
Lowest set bit       n & (-n)                   1010 & 0110 → 0010
Clear lowest set     n & (n - 1)                1010 & 1001 → 1000
Is power of 2        n > 0 && n&(n-1) == 0      1000 → true
All 1s mask (k bits) (1 << k) - 1               k=4: 1111
Count set bits       Kernighan's loop            (see Pattern 2)
NOT (bitwise)        ^n                          ^1010 = ...0101
XOR                  a ^ b                       1010 ^ 1100 = 0110
```

**Go-specific syntax:**
- `&^` is "AND NOT" (bit clear). Go has this as a dedicated operator;
  most languages use `& ~`.
- `^` is both unary NOT and binary XOR. Unary: `^n`. Binary: `a ^ b`.
- `>>` is arithmetic right shift for signed, logical for unsigned.
  Use `uint` types to guarantee logical shift.

---

## Common Interview Traps

### Bit Manipulation Traps

1. **Go uses `^` for both NOT (unary) and XOR (binary).**
   `^n` is bitwise NOT. `a ^ b` is XOR. Easy to confuse when reading
   code quickly. If you need to flip all bits, write `^n`. If you need
   to XOR two values, write `a ^ b`. Add a comment when the intent
   could be ambiguous.

2. **XOR trick only works with exact pairing.**
   If the problem says "one element appears once, rest appear three
   times," XOR alone won't work. You need bit counting per position
   (sum each bit position mod 3).

3. **Arithmetic vs logical right shift.**
   `>>` on a signed `int` in Go is arithmetic (preserves sign bit).
   `>> `on `uint` is logical (fills with zeros). If you're shifting
   a negative signed integer, the sign bit propagates. Use `uint`
   casts when you need logical shift behavior.

4. **Reverse bits: process ALL 32 bits.**
   A common mistake is looping only until the highest set bit. You
   must reverse all 32 positions. The bit at position 0 must end up
   at position 31, even if the upper bits are zero.

   ```go
   func reverseBits(n uint32) uint32 {
       var result uint32
       for i := 0; i < 32; i++ {
           result = (result << 1) | (n & 1)
           n >>= 1
       }
       return result
   }
   ```

5. **Operator precedence with `&`, `|`, `^`.**
   In Go, bitwise operators have higher precedence than `==` and `!=`,
   so `n & (n-1) == 0` works. But `n & 1 == 0` also works as expected
   because `&` binds before `==`. Still, use parentheses for clarity:
   `(n & 1) == 0`.

### Math Traps

6. **Sieve: start crossing out at p*p, not 2*p.**
   All multiples below p*p have already been marked by smaller primes.
   Starting at 2*p is correct but slower, and interviewers may ask why
   you didn't optimize.

7. **Integer overflow in sum formula.**
   `n * (n + 1) / 2` can overflow if n is large and you're using 32-bit
   integers. Go's `int` is 64-bit, but mention the concern in an
   interview to show awareness.

8. **Modular subtraction can go negative.**
   `(a - b) % MOD` can be negative in Go. Always use `((a - b) % MOD + MOD) % MOD`.

---

## Thought Process Walkthrough

### Walkthrough 1: Single Number (LC 136)

**Interviewer says:** "Given a non-empty array of integers where every
element appears twice except for one, find that single one. Do it in
O(n) time and O(1) space."

**Your response, out loud:**

> "The O(1) space constraint rules out a hash set. The O(n) time rules
> out sorting."
>
> "XOR has the property that a ^ a = 0 and a ^ 0 = a. If I XOR all
> elements together, every pair cancels to 0, and I'm left with the
> unique element."
>
> "I'll initialize result to 0, iterate through the array, XOR each
> element into result, and return it."

```go
func singleNumber(nums []int) int {
    result := 0
    for _, n := range nums {
        result ^= n
    }
    return result
}
```

> "Time O(n), space O(1). This works because XOR is commutative and
> associative, so the order doesn't matter."

**What the interviewer is checking:**
- Did you immediately recognize the XOR trick?
- Can you explain WHY it works (self-inverse property)?
- Did you state the constraint that all other elements appear exactly
  twice?

This should take under 3 minutes including explanation.

---

### Walkthrough 2: Number of 1 Bits (LC 191)

**Interviewer says:** "Write a function that takes an unsigned integer
and returns the number of '1' bits (Hamming weight)."

**Your response, out loud:**

> "I could check each of the 32 bits with a mask, but there's a faster
> approach: Kernighan's trick."
>
> "`n & (n - 1)` clears the lowest set bit. I can count how many times
> I do this before n becomes 0. That count equals the number of set
> bits."
>
> "For example, n = 0b1100. n-1 = 0b1011. n & (n-1) = 0b1000. One set
> bit cleared, count = 1. Repeat: 0b1000 & 0b0111 = 0. Count = 2. Done."

```go
func hammingWeight(n uint32) int {
    count := 0
    for n != 0 {
        n &= n - 1
        count++
    }
    return count
}
```

> "Time O(k) where k is the number of set bits. Space O(1)."

**What the interviewer is checking:**
- Do you know Kernighan's trick, or do you use the naive 32-iteration
  approach?
- Can you trace through an example to explain why `n & (n-1)` works?

This should take under 3 minutes.

---

## Time Targets

| Problem | Target | What it proves |
|---------|--------|----------------|
| Single Number | 3 min | XOR is automatic |
| Number of 1 Bits | 3 min | You know Kernighan's |
| Reverse Bits | 5 min | You can work with bit positions |
| Power of Two | 2 min | One-liner recognition |
| Missing Number | 3 min | XOR or sum — pick and explain |
| Count Primes | 8 min | You know the sieve |
| Subsets (bitmask) | 7 min | You can use bitmask iteration |

These are fast problems. If any takes more than double the target, drill
that specific pattern.

---

## Practice Plan (2 hours)

| Block | Minutes | Activity |
|-------|---------|----------|
| 1 | 10 | Read the pattern catalog. For each pattern, trace through one example by hand. |
| 2 | 5 | Memorize the quick reference card. Cover it, then write out all operations from memory. |
| 3 | 15 | Single Number + Number of 1 Bits + Power of Two + Missing Number. All from memory. These are one-pattern problems — should be fast. |
| 4 | 15 | Reverse Bits from memory. This is the trickiest bit problem. Make sure you process all 32 bits. |
| 5 | 10 | Subsets using bitmask iteration. Compare mentally with your backtracking version from Day 13. |
| 6 | 15 | Count Primes (sieve). Write it, then explain the `i*i` optimization out loud. |
| 7 | 10 | GCD, LCM, modular exponentiation. Write all three from memory. |
| 8 | 10 | Review the traps section. For each trap, write a one-line code example that demonstrates the bug. |
| 9 | 15 | Pick any two problems from the drill below and solve them timed. |
| 10 | 15 | Self-assessment. Review anything you got wrong. |

---

## Quick Drill: Flash Recall

Answer without looking up. Then verify.

1. What is `a ^ a`? What is `a ^ 0`?
2. What does `n & (n - 1)` do?
3. How do you check if bit 5 is set in Go?
4. What is Go's bit-clear operator and how does it differ from `& ~`?
5. Why start the sieve inner loop at `i * i` instead of `2 * i`?
6. What is `gcd(48, 18)` computed step by step?
7. Why is `((a - b) % MOD + MOD) % MOD` necessary instead of `(a - b) % MOD`?
8. When does the XOR-find-unique trick NOT work?
9. How many iterations does Kernighan's take for `n = 0b10000000`?
10. What is the difference between `^n` and `a ^ b` in Go?

<details>
<summary>Answers</summary>

1. `a ^ a = 0`. `a ^ 0 = a`. XOR is self-inverse and 0 is the identity.
2. Clears the lowest set bit of n.
3. `n & (1 << 5) != 0`.
4. `&^` (AND NOT). It's equivalent to `& ^` or `& ~` in other languages,
   but Go provides it as a single operator. `a &^ b` clears all bits in
   `a` that are set in `b`.
5. All multiples of `i` below `i * i` have already been marked as
   composite by smaller primes (e.g., `3 * 4` was already marked when
   processing prime 2).
6. `gcd(48, 18)` → `gcd(18, 48%18)` → `gcd(18, 12)` → `gcd(12, 6)` →
   `gcd(6, 0)` → 6.
7. Go's `%` operator can return negative values when the dividend is
   negative. Adding `MOD` before the final `%` ensures the result is
   non-negative.
8. When elements appear more than twice, or when there are multiple
   unique elements. XOR cancellation requires exact pairs.
9. One iteration. There is exactly one set bit, so `n & (n-1)` clears
   it immediately.
10. `^n` is unary bitwise NOT (flips all bits). `a ^ b` is binary XOR
    (flips bits where a and b differ). Same operator symbol, different
    arity.

</details>

---

## Self-Assessment

After completing the practice plan, score yourself honestly:

| Skill | Yes | No |
|-------|-----|----|
| I can explain XOR self-inverse and apply it to single-number problems instantly | | |
| I can write Kernighan's bit count and trace through an example | | |
| I can write bitmask subset iteration without reference | | |
| I know the Go-specific bit operators (`&^`, unary `^`, shift behavior) | | |
| I can write the Sieve of Eratosthenes with the i*i optimization | | |
| I can write GCD, LCM, and modular exponentiation from memory | | |
| I can list 3 bit manipulation traps from memory | | |
| Every problem in the time targets table is at or below target | | |

**If any "No":** Drill that specific pattern once more. These are
pattern-recognition problems — once you've seen the trick, it should be
instant recall. If it's not instant, you haven't drilled it enough.
