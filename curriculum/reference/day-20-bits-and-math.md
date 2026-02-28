# Day 20 — Bit Manipulation & Math

---

## 1. Curated Learning Resources

| # | Resource | Category | Why It's Useful |
|---|----------|----------|-----------------|
| 1 | [Bit Twiddling Hacks (Stanford)](https://graphics.stanford.edu/~seander/bithacks.html) | Reference | The definitive reference for bit tricks. Bookmarkable. Covers everything from counting bits to computing abs without branching. Keep it open while solving problems. |
| 2 | [Go `math/bits` Package Docs](https://pkg.go.dev/math/bits) | Go stdlib | See what Go gives you for free: `OnesCount`, `Reverse`, `LeadingZeros`, `TrailingZeros`, `Len`. Implement them yourself today, then use this package in production. |
| 3 | [Bitwise Cmd — Interactive Binary Visualizer](https://bitwisecmd.com/) | Visualizer | Type an expression like `42 & 15` and see the binary representations and result side by side. Great for building intuition when a trick isn't clicking. |
| 4 | [Visualgo — Bitmask](https://visualgo.net/en/bitmask) | Visualizer | Animated visualization of bitmask operations and how they map to set operations. Helps connect abstract bit ops to concrete subset logic. |
| 5 | [XOR — The Magical Bitwise Operator (Florian Hartmann)](https://florian.github.io/xor-trick/) | Explanation | Clear walkthrough of XOR properties with proofs and examples. Covers the single-number trick, swap without temp, and linked-list XOR trick. Read before the session. |
| 6 | [Sieve of Eratosthenes Animation (Wikipedia)](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes#/media/File:Sieve_of_Eratosthenes_animation.gif) | Visualizer | Animated GIF of the sieve crossing out composites step by step. One look and the algorithm clicks. |
| 7 | [Two's Complement Explanation (Cornell CS 3110)](https://www.cs.cornell.edu/~tomf/notes/cps104/twoscomp.html) | Explanation | Concise, correct explanation of two's complement with examples. Essential for understanding why `n & (-n)` isolates the lowest set bit and why `^x` behaves differently on signed vs unsigned types. |
| 8 | [LeetCode Bit Manipulation Topic](https://leetcode.com/tag/bit-manipulation/) | Practice | Sorted by difficulty. After today's session, try: Single Number (136), Hamming Weight (191), Reverse Bits (190), Missing Number (268), Power of Two (231), Counting Bits (338). |

---

## 2. Detailed 2-Hour Session Plan

### Review Block (12:00 -- 12:20) — Internalize the Theory

| Time | Activity |
|------|----------|
| 12:00 -- 12:05 | Read the operator table in Section 3.1. For each operator, mentally compute `13 op 6` (binary: `1101` op `0110`). Write the results by hand on paper. Confirm: AND=`0100`(4), OR=`1111`(15), XOR=`1011`(11), NOT 13 on 4 bits=`0010`(2 as unsigned 4-bit). |
| 12:05 -- 12:12 | Study the **Bit Manipulation Trick Reference Card** (Section 5). For each trick, trace one example in binary on paper. Don't just read — write the bits. |
| 12:12 -- 12:17 | Read the **XOR properties** section (3.2). Convince yourself why XOR-ing all elements of `[4, 1, 2, 1, 2]` yields `4`. Write the accumulator at each step (Section 6 has the diagram). |
| 12:17 -- 12:20 | Read the **Kernighan's algorithm** explanation (3.3). Trace `n = 11` (`1011`) through the loop: `1011 & 1010 = 1010`, `1010 & 1001 = 1000`, `1000 & 0111 = 0000`. Three iterations, three set bits. |

### Implement Block (12:20 -- 1:20) — Code From Scratch

| Time | Problem | Key Insight | Target |
|------|---------|-------------|--------|
| 12:20 -- 12:30 | **SingleNumber** | XOR all elements. `a ^ a = 0` cancels pairs, leaving the unique one. | 5 min code, 5 min tests |
| 12:30 -- 12:40 | **HammingWeight** | Kernighan's: `n &= n - 1` clears the lowest set bit. Count iterations. | Also test with `uint32(0)` and `uint32(math.MaxUint32)` |
| 12:40 -- 12:52 | **ReverseBits** | Loop 32 times: extract LSB of input, shift result left, OR the bit in. | This one trips people up. Draw it on paper first. |
| 12:52 -- 1:02 | **MissingNumber** | XOR indices `0..n` with all array elements. Pairs cancel, missing number survives. | Alternative: sum formula `n*(n+1)/2 - sum(nums)`. Implement both. |
| 1:02 -- 1:08 | **IsPowerOfTwo** | `n > 0 && n & (n-1) == 0`. Powers of 2 have exactly one set bit. | Edge cases: 0, 1, negative numbers, `math.MinInt` |
| 1:08 -- 1:20 | **CountPrimes** (Sieve) | Sieve of Eratosthenes. Allocate `[]bool`, cross out multiples starting at `p*p`. | Watch out: count primes *less than* n, not *up to* n. |

### Solidify Block (1:20 -- 1:50) — Edge Cases, Variants, Math

| Time | Activity |
|------|----------|
| 1:20 -- 1:30 | **GCD** — Implement the Euclidean algorithm (both recursive and iterative). Test with `(0, 5)`, `(12, 8)`, `(17, 13)`, `(-12, 8)`. Then implement LCM using GCD. |
| 1:30 -- 1:38 | **Edge case sweep** — Go back through every function. Test: zero input, single-element arrays, `uint32` max, negative numbers for GCD. Fix any bugs. |
| 1:38 -- 1:45 | **Bitmask as set** exercise — Write a short function that uses a `uint32` to represent a set of letters. Implement `add(set, char)`, `contains(set, char)`, `remove(set, char)`. Use this to check if a string has all unique characters in O(1) space. |
| 1:45 -- 1:50 | **Variant**: try Single Number II — every element appears three times except one. XOR alone doesn't work. Think about counting bits modulo 3 in each position. (Stretch goal — don't worry if you don't finish.) |

### Recap Block (1:50 -- 2:00) — Write From Memory

| Time | Activity |
|------|----------|
| 1:50 -- 1:55 | Close all references. On paper or a blank file, write: (1) The complexity of Kernighan's, (2) Why `n & (n-1)` clears the lowest set bit, (3) The Sieve time complexity, (4) The GCD recurrence. |
| 1:55 -- 2:00 | Write one gotcha for each problem you implemented. If you got stuck anywhere, note *why* — that's tomorrow's warm-up. |

---

## 3. Core Concepts Deep Dive

### 3.1 Bitwise Operators in Go

Go has six bitwise operators. Five are binary (two operands), one is unary:

```
AND       a & b      1 only if both bits are 1
OR        a | b      1 if either bit is 1
XOR       a ^ b      1 if bits differ
AND NOT   a &^ b     clear bits in a that are set in b  (Go-specific operator)
LEFT SH   a << n     shift bits left by n (multiply by 2^n)
RIGHT SH  a >> n     shift bits right by n (divide by 2^n)
```

**Go's NOT is unusual.** Most languages use `~x` for bitwise NOT. Go uses `^x` (unary `^`). The same symbol `^` serves double duty:
- **Binary**: `a ^ b` = XOR
- **Unary**: `^x` = bitwise NOT (flips every bit)

The compiler tells them apart by context — `^x` has one operand, `a ^ b` has two.

```go
x := uint8(0b00001111)
fmt.Printf("%08b\n", ^x)    // 11110000  (NOT)
fmt.Printf("%08b\n", x^0b11001100) // 11000011  (XOR)
```

**AND NOT (`&^`)** is Go's unique operator. `a &^ b` clears bits: it's equivalent to `a & (^b)` in other languages. Useful for clearing specific bits without a separate NOT step.

```go
n := 0b1111
mask := 0b0110
fmt.Printf("%04b\n", n &^ mask) // 1001 — cleared bits 1 and 2
```

**Shift behavior on signed types:** `>>` on a signed integer performs *arithmetic* shift — it preserves the sign bit. On unsigned integers, `>>` performs *logical* shift — fills with zeros. For bit manipulation, prefer `uint32` or `uint64` to avoid surprises.

### 3.2 XOR Properties and Proofs

XOR has four properties that make it a workhorse for bit problems:

| Property | Statement | Why |
|----------|-----------|-----|
| Self-inverse | `a ^ a = 0` | Each bit: `0^0=0`, `1^1=0` |
| Identity | `a ^ 0 = a` | Each bit: `0^0=0`, `1^0=1` |
| Commutative | `a ^ b = b ^ a` | Bit-by-bit comparison is symmetric |
| Associative | `(a ^ b) ^ c = a ^ (b ^ c)` | Can reorder/regroup freely |

**Why XOR finds the unique element:**

Given `[a, b, a, c, b]`, compute `a ^ b ^ a ^ c ^ b`.

By commutativity and associativity, regroup: `(a ^ a) ^ (b ^ b) ^ c = 0 ^ 0 ^ c = c`.

Every paired element cancels to zero. The unique element survives. This works regardless of array order — commutativity and associativity let us rearrange freely.

**Key limitation:** This only works when every duplicate appears an *even* number of times and exactly one element appears an *odd* number of times. For three occurrences (Single Number II), you need to count bits modulo 3 instead.

### 3.3 Kernighan's Bit Counting: Why `n & (n-1)` Clears the Lowest Set Bit

Consider any integer `n` with its lowest set bit at position `k`:

```
n     = ...1???...1 0 0 0    (bit k is the rightmost 1)
                    ^--- k
```

Subtracting 1 flips bit `k` to 0 and all zeros below it to 1:

```
n     = ...1???...1 0 0 0
n - 1 = ...1???...0 1 1 1
```

AND-ing them: the bits above `k` are unchanged (`1 & 1 = 1`, `0 & 0 = 0`). Bit `k` becomes `1 & 0 = 0`. Bits below `k` become `0 & 1 = 0`.

```
n & (n-1) = ...1???...0 0 0 0    (lowest set bit cleared!)
```

**Kernighan's algorithm** exploits this: keep clearing the lowest set bit until `n` is zero. The number of iterations equals the number of set bits (population count).

```go
func HammingWeight(n uint32) int {
    count := 0
    for n != 0 {
        n &= n - 1  // clear lowest set bit
        count++
    }
    return count
}
```

Time: O(k) where k = number of set bits. Much better than checking all 32 bits when the number is sparse.

### 3.4 Bit Masks for Sets

An integer can represent a **set of elements** where bit `i` being set means "element `i` is in the set." For a universe of up to 32 elements, a `uint32` suffices. Up to 64, use `uint64`.

| Set Operation | Bit Expression | Example |
|---------------|---------------|---------|
| Empty set | `0` | `set := 0` |
| Singleton {i} | `1 << i` | `1 << 3` = `{3}` |
| Union A ∪ B | `A \| B` | `0b1010 \| 0b0110 = 0b1110` |
| Intersection A ∩ B | `A & B` | `0b1010 & 0b0110 = 0b0010` |
| Contains i? | `A & (1 << i) != 0` | |
| Add i | `A \| (1 << i)` | |
| Remove i | `A &^ (1 << i)` | |
| Toggle i | `A ^ (1 << i)` | |
| Complement | `^A` (within bit width) | |
| Set size | `bits.OnesCount(A)` | |
| Iterate subsets of S | `for sub := S; sub > 0; sub = (sub - 1) & S` | Classic bitmask DP loop |

**Bitmask DP** uses this to represent states. Example: Traveling Salesman on n cities. State = `dp[visited_mask][current_city]` where `visited_mask` is a bitmask of which cities have been visited. For n <= 20, this is feasible (2^20 ~ 1M states).

### 3.5 Two's Complement

Modern computers represent signed integers using **two's complement**:

- Positive numbers: standard binary.
- Negative numbers: flip all bits of the positive value and add 1.

For an 8-bit signed integer:
```
 5 in binary:  0000 0101
~5 (flip):     1111 1010
+1:            1111 1011  =  -5
```

So `-n = ^n + 1` (in Go: `^n + 1`, which is the same as `-n` for signed integers).

**Why `n & (-n)` isolates the lowest set bit:**

```
n     = ...1010 1000        (some number)
-n    = ...0101 1000        (two's complement: flip + 1)
n & -n = ...0000 1000       (only the lowest set bit survives)
```

The bits above the lowest set bit are flipped (complemented), so AND makes them 0. The lowest set bit itself is 1 in both `n` and `-n`. The bits below it are 0 in both.

This trick is used in **Fenwick trees** (Binary Indexed Trees) to navigate the tree structure, and for isolating the least significant bit in various algorithms.

### 3.6 GCD and the Euclidean Algorithm

The **greatest common divisor** of two integers `a` and `b` is the largest integer that divides both.

**Euclidean algorithm**: `gcd(a, b) = gcd(b, a % b)`, with base case `gcd(a, 0) = a`.

**Why it terminates:** Each step replaces the larger number with a remainder that is strictly smaller. The remainder is always non-negative and less than the divisor, so the sequence is strictly decreasing and must reach 0.

**Why it's correct:** If `d` divides both `a` and `b`, then `d` divides `a - q*b = a % b`. So the set of common divisors of `(a, b)` is the same as `(b, a % b)`. The GCD is preserved.

**Time complexity:** O(log(min(a, b))). At each step, the remainder is at most half the divisor (by a Fibonacci-based argument), so the number of steps is at most ~2*log_2(min(a, b)).

```go
// Recursive
func GCD(a, b int) int {
    if b == 0 {
        return a
    }
    return GCD(b, a%b)
}

// Iterative
func GCDIterative(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}
```

**Extended GCD** finds integers `x, y` such that `a*x + b*y = gcd(a, b)`. This is used for computing modular inverses: if `gcd(a, m) = 1`, then `a*x = 1 (mod m)` where `x` is from the extended GCD. Rarely needed in interviews but good to know exists.

**LCM via GCD:** `lcm(a, b) = |a * b| / gcd(a, b)`. Compute GCD first to avoid overflow.

### 3.7 Sieve of Eratosthenes

The sieve finds all primes up to `n` by iteratively marking composite numbers.

**Algorithm:**
1. Create a boolean array `isPrime[0..n]`, initialize all to `true`. Set `isPrime[0] = false`, `isPrime[1] = false`.
2. For each `p` from 2 to `sqrt(n)`:
   - If `isPrime[p]` is true (p is prime):
     - Mark all multiples of `p` starting from `p*p` as `false`.
3. All indices still marked `true` are prime.

**Why start at p^2?** All smaller multiples of `p` (i.e., `2p, 3p, ..., (p-1)p`) have already been crossed out by smaller primes. For example, `3 * 5 = 15` was already crossed out when processing prime 3 (as a multiple of 3). When we reach prime 5, we can start at `5 * 5 = 25`.

**Time complexity:** O(n log log n). The sum of `n/p` over all primes `p <= n` is `n * sum(1/p)` ~ `n * ln(ln(n))` by Mertens' theorem. This is nearly linear — the log log n factor grows extremely slowly.

**Space:** O(n) for the boolean array.

```go
func CountPrimes(n int) int {
    if n <= 2 {
        return 0
    }
    isPrime := make([]bool, n)
    for i := 2; i < n; i++ {
        isPrime[i] = true
    }
    for p := 2; p*p < n; p++ {
        if isPrime[p] {
            for multiple := p * p; multiple < n; multiple += p {
                isPrime[multiple] = false
            }
        }
    }
    count := 0
    for i := 2; i < n; i++ {
        if isPrime[i] {
            count++
        }
    }
    return count
}
```

---

## 4. Implementation Checklist

### Function Signatures

```go
package bits

// --- Bit Manipulation ---

// SingleNumber returns the element that appears exactly once.
// All other elements appear exactly twice.
// Approach: XOR all elements. Pairs cancel to 0.
// Time: O(n)  Space: O(1)
func SingleNumber(nums []int) int

// HammingWeight returns the number of 1-bits in n.
// Approach: Kernighan's — n &= n-1 clears lowest set bit.
// Time: O(k) where k = number of set bits  Space: O(1)
func HammingWeight(n uint32) int

// ReverseBits reverses the 32 bits of n.
// Approach: loop 32 times, extract LSB, shift result left, OR in.
// Time: O(32) = O(1)  Space: O(1)
func ReverseBits(n uint32) uint32

// MissingNumber finds the missing number in [0, n] given n numbers.
// Approach: XOR all indices 0..n with all array elements.
// Time: O(n)  Space: O(1)
func MissingNumber(nums []int) int

// IsPowerOfTwo returns true if n is a power of 2.
// Approach: n > 0 && n & (n-1) == 0.
// Time: O(1)  Space: O(1)
func IsPowerOfTwo(n int) bool

// --- Math ---

// CountPrimes returns the count of primes less than n.
// Approach: Sieve of Eratosthenes.
// Time: O(n log log n)  Space: O(n)
func CountPrimes(n int) int

// GCD returns the greatest common divisor of a and b.
// Approach: Euclidean algorithm.
// Time: O(log(min(a,b)))  Space: O(1) iterative
func GCD(a, b int) int
```

### Test Cases and Edge Cases

| Function | Must-Test Inputs | Why |
|----------|-----------------|-----|
| `SingleNumber` | `[2,2,1]` -> 1; `[4,1,2,1,2]` -> 4; `[1]` -> 1 | Basic, multiple pairs, single element |
| `HammingWeight` | `0` -> 0; `0b1011` -> 3; `0xFFFFFFFF` -> 32 | Zero, sparse, all bits set |
| `ReverseBits` | `0` -> 0; `0b1` -> `1 << 31`; `0b10100101000001111010011100` -> known value | Zero, single bit (must land at MSB), LeetCode example |
| `MissingNumber` | `[3,0,1]` -> 2; `[0]` -> 1; `[0,1]` -> 2; `[1]` -> 0 | Normal, single element, missing at end, missing 0 |
| `IsPowerOfTwo` | `0` -> false; `1` -> true; `16` -> true; `18` -> false; `-16` -> false | Zero trap, 2^0, power, non-power, negative |
| `CountPrimes` | `0` -> 0; `1` -> 0; `2` -> 0; `3` -> 1; `10` -> 4; `100` -> 25 | Boundaries, small n, known counts |
| `GCD` | `(12, 8)` -> 4; `(0, 5)` -> 5; `(5, 0)` -> 5; `(17, 13)` -> 1; `(-12, 8)` -> +/-4 | Normal, zero input, coprime, negative (handle with abs) |

### Sample Test Structure

```go
func TestSingleNumber(t *testing.T) {
    tests := []struct {
        nums []int
        want int
    }{
        {[]int{2, 2, 1}, 1},
        {[]int{4, 1, 2, 1, 2}, 4},
        {[]int{1}, 1},
        {[]int{-1, -1, -2}, -2},
    }
    for _, tt := range tests {
        got := SingleNumber(tt.nums)
        if got != tt.want {
            t.Errorf("SingleNumber(%v) = %d, want %d", tt.nums, got, tt.want)
        }
    }
}

func TestHammingWeight(t *testing.T) {
    tests := []struct {
        n    uint32
        want int
    }{
        {0, 0},
        {0b00001011, 3},
        {0xFFFFFFFF, 32},
        {1, 1},
        {1 << 31, 1},
    }
    for _, tt := range tests {
        got := HammingWeight(tt.n)
        if got != tt.want {
            t.Errorf("HammingWeight(%032b) = %d, want %d", tt.n, got, tt.want)
        }
    }
}

func TestReverseBits(t *testing.T) {
    tests := []struct {
        n    uint32
        want uint32
    }{
        {0, 0},
        {1, 1 << 31},
        {0xFFFFFFFF, 0xFFFFFFFF},
        {0b00000010100101000001111010011100, 0b00111001011110000010100101000000},
    }
    for _, tt := range tests {
        got := ReverseBits(tt.n)
        if got != tt.want {
            t.Errorf("ReverseBits(%032b) = %032b, want %032b", tt.n, got, tt.want)
        }
    }
}

func TestMissingNumber(t *testing.T) {
    tests := []struct {
        nums []int
        want int
    }{
        {[]int{3, 0, 1}, 2},
        {[]int{0}, 1},
        {[]int{0, 1}, 2},
        {[]int{1}, 0},
        {[]int{9, 6, 4, 2, 3, 5, 7, 0, 1}, 8},
    }
    for _, tt := range tests {
        got := MissingNumber(tt.nums)
        if got != tt.want {
            t.Errorf("MissingNumber(%v) = %d, want %d", tt.nums, got, tt.want)
        }
    }
}

func TestIsPowerOfTwo(t *testing.T) {
    tests := []struct {
        n    int
        want bool
    }{
        {0, false},
        {1, true},
        {2, true},
        {16, true},
        {18, false},
        {-16, false},
        {1 << 30, true},
    }
    for _, tt := range tests {
        got := IsPowerOfTwo(tt.n)
        if got != tt.want {
            t.Errorf("IsPowerOfTwo(%d) = %v, want %v", tt.n, got, tt.want)
        }
    }
}

func TestCountPrimes(t *testing.T) {
    tests := []struct {
        n    int
        want int
    }{
        {0, 0},
        {1, 0},
        {2, 0},
        {3, 1},
        {10, 4},
        {100, 25},
        {1000, 168},
    }
    for _, tt := range tests {
        got := CountPrimes(tt.n)
        if got != tt.want {
            t.Errorf("CountPrimes(%d) = %d, want %d", tt.n, got, tt.want)
        }
    }
}

func TestGCD(t *testing.T) {
    tests := []struct {
        a, b int
        want int
    }{
        {12, 8, 4},
        {0, 5, 5},
        {5, 0, 5},
        {17, 13, 1},
        {100, 75, 25},
    }
    for _, tt := range tests {
        got := GCD(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
        }
    }
}
```

---

## 5. Bit Manipulation Trick Reference Card

Keep this open while solving problems. Every trick uses O(1) time and space.

```
+---------------------------------------------------------------------+
|                  BIT MANIPULATION CHEAT SHEET (Go)                   |
+----------------------+----------------------+-----------------------+
| Operation            | Go Expression        | Example (n = 0b1010) |
+----------------------+----------------------+-----------------------+
| Check bit i is set   | n & (1 << i) != 0    | i=1: 1010 & 0010 = 2 |
| Set bit i            | n | (1 << i)         | i=0: 1010 | 0001=1011|
| Clear bit i          | n &^ (1 << i)        | i=1: 1010 &^0010=1000|
| Toggle bit i         | n ^ (1 << i)         | i=3: 1010 ^ 1000=0010|
| Clear lowest set bit | n & (n - 1)          | 1010 & 1001 = 1000   |
| Isolate lowest set   | n & (-n)             | 1010 & 0110 = 0010   |
|  bit                 |                      |                       |
| Is power of 2?       | n > 0 && n&(n-1)==0  | 1000 -> true          |
| All bits below i     | (1 << i) - 1         | i=3: 0111             |
| All bits i and above | ^((1 << i) - 1)      | i=2: ...11111100      |
| Count set bits       | Kernighan's loop     | See HammingWeight     |
| Swap a, b (no temp)  | a ^= b; b ^= a;     | Works but prefer Go's |
|                      | a ^= b               | a, b = b, a           |
| Bitwise NOT          | ^n                   | ^1010 = ...0101       |
| AND NOT (clear mask) | n &^ mask            | Go-specific operator  |
| Sign of int          | n >> 63 (for int64)  | 0 if pos, -1 if neg   |
| Absolute value       | (n ^ (n >> 63)) -    | Branchless, signed    |
|  (int64)             |   (n >> 63)          |                       |
| Check if even        | n & 1 == 0           | 1010 & 0001 = 0       |
| Multiply by 2^k      | n << k               | 3 << 2 = 12           |
| Divide by 2^k (uint) | n >> k               | 12 >> 2 = 3           |
+----------------------+----------------------+-----------------------+
| XOR PROPERTIES                                                      |
|   a ^ a = 0          a ^ 0 = a          Commutative & Associative   |
+---------------------------------------------------------------------+
| BITMASK SET OPERATIONS                                               |
|   Union:       A | B              Intersection:  A & B               |
|   Difference:  A &^ B             Symmetric Diff: A ^ B             |
|   Add elem i:  A | (1 << i)      Remove elem i: A &^ (1 << i)      |
|   Contains i:  A & (1 << i) != 0  Size:  bits.OnesCount(A)         |
|   Iterate subsets of S:  for sub := S; sub > 0; sub = (sub-1) & S   |
+---------------------------------------------------------------------+
| COMMON PATTERNS                                                      |
|   Find unique elem in pairs:       XOR all elements                  |
|   Missing number in [0..n]:        XOR indices with elements         |
|   Two unique elems in pairs:       XOR all -> diff bit -> split      |
|   Reverse bits:                    Loop 32x, extract LSB, shift      |
+---------------------------------------------------------------------+
```

---

## 6. Visual Diagrams

### 6.1 Binary Representation with Bit Positions

```
  Bit Position:   7    6    5    4    3    2    1    0
  Bit Value:     128   64   32   16   8    4    2    1
                +----+----+----+----+----+----+----+----+
  n = 42:       |  0 |  0 |  1 |  0 |  1 |  0 |  1 |  0 |  = 32 + 8 + 2 = 42
                +----+----+----+----+----+----+----+----+
                         ^         ^         ^
                         |         |         |
                      bit 5     bit 3     bit 1
                    (set)      (set)     (set)

  Operations on n = 42 (0b00101010):

  Check bit 3:   42 & (1<<3)  = 00101010 & 00001000 = 00001000 != 0  -> SET
  Check bit 2:   42 & (1<<2)  = 00101010 & 00000100 = 00000000 = 0   -> NOT SET
  Set bit 0:     42 | (1<<0)  = 00101010 | 00000001 = 00101011 = 43
  Clear bit 3:   42 &^(1<<3)  = 00101010 &^00001000 = 00100010 = 34
  Toggle bit 5:  42 ^ (1<<5)  = 00101010 ^ 00100000 = 00001010 = 10
```

### 6.2 Kernighan's Algorithm: Clearing Bits Step by Step

```
  Count set bits of n = 53  (binary: 00110101)

  Iteration 1:
    n       = 0 0 1 1 0 1 0 1     (53)
    n - 1   = 0 0 1 1 0 1 0 0     (52)
    n&(n-1) = 0 0 1 1 0 1 0 0     (52)  <- cleared bit 0
              count = 1                      ^

  Iteration 2:
    n       = 0 0 1 1 0 1 0 0     (52)
    n - 1   = 0 0 1 1 0 0 1 1     (51)
    n&(n-1) = 0 0 1 1 0 0 0 0     (48)  <- cleared bit 2
              count = 2              ^

  Iteration 3:
    n       = 0 0 1 1 0 0 0 0     (48)
    n - 1   = 0 0 1 0 1 1 1 1     (47)
    n&(n-1) = 0 0 1 0 0 0 0 0     (32)  <- cleared bit 4
              count = 3          ^

  Iteration 4:
    n       = 0 0 1 0 0 0 0 0     (32)
    n - 1   = 0 0 0 1 1 1 1 1     (31)
    n&(n-1) = 0 0 0 0 0 0 0 0     (0)   <- cleared bit 5
              count = 4      ^

  n == 0 -> STOP.  Answer: 4 set bits.

  Note: We only looped 4 times (not 8) -- once per set bit.
```

### 6.3 XOR Accumulator: Finding the Single Number

```
  Input: [4, 1, 2, 1, 2]
  Find the element that appears once.

  Step    Element    XOR Accumulator (binary)    Decimal
  ----    -------    ------------------------    -------
  init               0 0 0                       0
   1       4         1 0 0                       4        (0 ^ 4 = 4)
   2       1         1 0 1                       5        (4 ^ 1 = 5)
   3       2         1 1 1                       7        (5 ^ 2 = 7)
   4       1         1 1 0                       6        (7 ^ 1 = 6)  <- 1 cancels
   5       2         1 0 0                       4        (6 ^ 2 = 4)  <- 2 cancels

  Result: 4

  What happened:
    4 ^ 1 ^ 2 ^ 1 ^ 2
  = 4 ^ (1 ^ 1) ^ (2 ^ 2)     <- regroup by associativity/commutativity
  = 4 ^    0    ^    0
  = 4
```

### 6.4 Sieve of Eratosthenes: Crossing Out Grid

```
  Find all primes less than 30.

  Initial grid (all marked as potentially prime):

   2  3  4  5  6  7  8  9 10 11
  12 13 14 15 16 17 18 19 20 21
  22 23 24 25 26 27 28 29

  Step 1: p = 2 (prime). Cross out multiples of 2 starting at 2*2=4:
   2  3  x  5  x  7  x  x  x  11
   x  13  x  x  x  17  x  19  x  x
   x  23  x  x  x  x  x  29

  Step 2: p = 3 (prime). Cross out multiples of 3 starting at 3*3=9:
   2  3  x  5  x  7  x  x  x  11
   x  13  x  x  x  17  x  19  x  x
   x  23  x  x  x  x  x  29

  (9, 15, 21, 27 newly crossed; 6, 12, 18, 24 already crossed by 2)

  Step 3: p = 5 (prime). Cross out starting at 5*5=25:
   2  3  x  5  x  7  x  x  x  11
   x  13  x  x  x  17  x  19  x  x
   x  23  x  x  x  x  x  29

  (25 newly crossed; all other multiples of 5 already crossed)

  p = 6 -> 6*6 = 36 > 30. STOP.

  Remaining primes: 2, 3, 5, 7, 11, 13, 17, 19, 23, 29
  Count: 10 primes less than 30

  Full grid with final state:
  +----+----+----+----+----+----+----+----+----+----+
  |  2 |  3 |  x |  5 |  x |  7 |  x |  x |  x | 11 |
  |  P |  P |    |  P |    |  P |    |    |    |  P |
  +----+----+----+----+----+----+----+----+----+----+
  |  x | 13 |  x |  x |  x | 17 |  x | 19 |  x |  x |
  |    |  P |    |    |    |  P |    |  P |    |    |
  +----+----+----+----+----+----+----+----+----+----+
  |  x | 23 |  x |  x |  x |  x |  x | 29 |    |    |
  |    |  P |    |    |    |    |    |  P |    |    |
  +----+----+----+----+----+----+----+----+----+----+
    P = prime    x = composite (crossed out)
```

### 6.5 Reverse Bits: Step-by-Step

```
  Reverse bits of n = 0b1101 (13), shown for 8-bit simplicity:

  Input:   n = 0 0 0 0 1 1 0 1
  Output: result (building right to left)

  Step   Extract LSB    Shift result left   OR bit in     n >>= 1
  ----   -----------    ----------------    ----------    -------
   1     n&1 = 1        result=00000000     00000001      n=00000110
   2     n&1 = 0        result=00000010     00000010      n=00000011
   3     n&1 = 1        result=00000100     00000101      n=00000001
   4     n&1 = 1        result=00001010     00001011      n=00000000
   5     n&1 = 0        result=00010110     00010110      n=00000000
   6     n&1 = 0        result=00101100     00101100      n=00000000
   7     n&1 = 0        result=01011000     01011000      n=00000000
   8     n&1 = 0        result=10110000     10110000      n=00000000

  Result: 10110000 = 176

  Verify: input 00001101 reversed is 10110000

  Key: You MUST loop all 32 (or 8) times, even after n becomes 0,
       because trailing zeros in n become leading zeros in the result.
```

### 6.6 n & (n-1) vs n & (-n): Two Complementary Tricks

```
  n = 52 = 0b00110100

  n & (n-1): CLEARS the lowest set bit
  -----------------------------------------------
    n       = 0 0 1 1 0 1 0 0    (52)
    n - 1   = 0 0 1 1 0 0 1 1    (51)
    n&(n-1) = 0 0 1 1 0 0 0 0    (48)
                          ^
                          lowest set bit GONE

  n & (-n): ISOLATES the lowest set bit
  -----------------------------------------------
    n       = 0 0 1 1 0 1 0 0    (52)
    -n      = 1 1 0 0 1 1 0 0    (two's complement of 52)
    n&(-n)  = 0 0 0 0 0 1 0 0    (4)
                          ^
                          lowest set bit ONLY

  Use cases:
    n & (n-1) -> Kernighan's bit count, power-of-2 check
    n & (-n)  -> Fenwick tree navigation, isolating LSB
```

---

## 7. Self-Assessment

Answer these without looking at the material. If you can't, revisit the relevant section.

### Question 1: Kernighan's Core Insight
**Why does `n & (n-1)` clear the lowest set bit?** Trace it in binary for `n = 12` (`1100`). What is `n` after each iteration of Kernighan's algorithm?

<details>
<summary>Check your answer</summary>

When you subtract 1 from `n`, the lowest set bit flips to 0 and all lower bits flip to 1:
```
n     = 1100
n - 1 = 1011
n & (n-1) = 1000  <- bit 2 cleared
```
Next iteration:
```
n     = 1000
n - 1 = 0111
n & (n-1) = 0000  <- bit 3 cleared, done
```
Two iterations -> two set bits. Kernighan's runs in O(k) where k = number of set bits.
</details>

### Question 2: XOR Limitations
**Single Number uses XOR to find the unique element. What assumption must hold for this to work?** What happens if one element appears 3 times instead of 2? How would you adapt?

<details>
<summary>Check your answer</summary>

XOR cancellation requires each duplicate to appear an **even** number of times (specifically, the problem states exactly twice). `a ^ a = 0` only works for pairs.

If one element appears 3 times: `a ^ a ^ a = a` (it doesn't cancel). XOR alone fails.

Adaptation for "all appear 3 times except one": count the number of 1-bits in each position across all numbers. For each bit position, the count modulo 3 gives the bit of the unique number. This uses O(32) = O(1) space with a 32-element count array.
</details>

### Question 3: Bitmask DP Setup
**How would you use bitmasks to represent subsets in a DP problem?** Describe the state representation for the Traveling Salesman Problem on 5 cities. How many total states exist? What does `dp[0b10110][3]` represent?

<details>
<summary>Check your answer</summary>

State: `dp[mask][i]` = minimum cost to visit exactly the cities in `mask`, ending at city `i`.

For 5 cities, `mask` ranges from `0b00000` to `0b11111` (0 to 31), giving 32 possible subsets. Combined with 5 ending cities: 32 x 5 = 160 total states.

`dp[0b10110][3]` = minimum cost to visit cities {1, 2, 4} (bits 1, 2, 4 are set), ending at city 3. Wait -- city 3 corresponds to bit 3, which is NOT set in `0b10110`. This state is invalid: the ending city must be in the visited set. A valid example: `dp[0b10110][2]` = min cost visiting cities {1, 2, 4}, ending at city 2.
</details>

### Question 4: Sieve Optimization
**In the Sieve of Eratosthenes, why do we start crossing out at `p*p` instead of `2*p`?** What would happen if we started at `2*p`? Would the result be different?

<details>
<summary>Check your answer</summary>

Starting at `2*p` would produce the **same correct result** -- but would do redundant work. Every composite number `k*p` where `k < p` has a prime factor smaller than `p`, so it was already crossed out when we processed that smaller prime.

For example, when p = 7: `2x7=14` (crossed by 2), `3x7=21` (crossed by 3), `4x7=28` (crossed by 2), `5x7=35` (crossed by 5), `6x7=42` (crossed by 2). The first multiple of 7 not yet crossed is `7x7=49`.

Starting at `p*p` is a constant-factor optimization, not an asymptotic one. The time complexity remains O(n log log n) either way.
</details>

### Question 5: Go-Specific Bit Gotcha
**What is the difference between `^x` and `x ^ y` in Go? Why doesn't Go use `~` for bitwise NOT?** What happens if you try to compute `^uint32(0)` vs `^int32(0)`?

<details>
<summary>Check your answer</summary>

`^x` (unary) is bitwise NOT -- flips all bits. `x ^ y` (binary) is XOR. Go uses the same symbol for both; the compiler disambiguates by operand count.

Go doesn't use `~` because it follows a principle of minimal syntax -- `^` already handles both roles. (Plan 9 C, Go's ancestor, also used `^` for NOT.)

- `^uint32(0)` = `0xFFFFFFFF` (all 32 bits set, value 4294967295)
- `^int32(0)` = `-1` (all 32 bits set in two's complement, interpreted as signed = -1)

Both have the same bit pattern, but the type determines the interpretation. This matters for right shifts: `^int32(0) >> 1` = `-1` (arithmetic shift preserves sign), while `^uint32(0) >> 1` = `2147483647` (logical shift fills with zero).
</details>
