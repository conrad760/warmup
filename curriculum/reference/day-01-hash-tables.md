# Day 1 -- Hash Tables

**Date:** Week 1, Day 1
**Time Block:** 12:00 PM - 2:00 PM (2 hours)
**Goal:** Build two hash map implementations from scratch in Go (chaining and open addressing), understand the internals of Go's built-in `map`, and internalize the interview patterns where hash tables are the key insight.

---

## 1. Curated Learning Resources

### Read

1. **"Hash Tables" -- VisuAlgo**
   https://visualgo.net/en/hashtable
   **Category:** Interact
   Interactive step-by-step visualization of chaining and open addressing (linear probing, quadratic probing, double hashing). Step through insertions and deletions to build intuition for how probe chains form and break. Use this *before* coding to ground the concept visually.

2. **"How the Go runtime implements maps efficiently (without generics)" -- Dave Cheney**
   https://dave.cheney.net/2018/05/29/how-the-go-runtime-implements-maps-efficiently-without-generics
   **Category:** Read
   Explains Go's actual map implementation: bucket arrays, overflow buckets, the tophash optimization, and the incremental evacuation strategy during growth. This is the single best resource for understanding what `map[string]int` does under the hood.

3. **Go Source: `runtime/map.go`**
   https://github.com/golang/go/blob/master/src/runtime/map.go
   **Category:** Reference
   The real implementation. You don't need to read all ~1600 lines, but skim the `hmap` struct definition, `makemap`, `mapaccess1`, `mapassign`, and `growWork`. Reading production source code shows you how the constants (bucket size = 8, load factor = 6.5) were chosen.

4. **"Hash Table" -- Programiz**
   https://www.programiz.com/dsa/hash-table
   **Category:** Read
   Clean, well-illustrated walkthrough of hash tables with chaining and open addressing. Good as a 10-minute refresher if the concepts feel rusty. Covers collision resolution strategies with diagrams.

5. **"Hashing" -- MIT OpenCourseWare 6.006 Lecture Notes**
   https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/resources/mit6_006s20_lec4/
   **Category:** Read
   Formal treatment of universal hashing, the birthday paradox applied to collisions, and the amortized analysis of table doubling. This is where the math behind load factor thresholds lives. Skim the proofs, absorb the intuitions.

6. **"Implementing a Hash Map in Go from Scratch"**
   https://medium.com/@mkeeler_34486/implementing-a-hash-map-in-go-from-scratch-b3e59d0e0e0e
   **Category:** Read
   A practical step-by-step implementation in Go with chaining. Useful as a reference *after* you've attempted your own implementation -- compare your design decisions.

7. **FNV Hash -- Go Standard Library**
   https://pkg.go.dev/hash/fnv
   **Category:** Reference
   Documentation for `hash/fnv`, the hash function you'll use in your implementation. FNV-1a is fast, simple, and has good distribution for string keys. Know how to use `fnv.New32a()` with `io.WriteString`.

---

## 2. Detailed 2-Hour Session Plan

### 12:00 - 12:20 | Review & Internalize (20 min)

| Time | Activity |
|------|----------|
| 12:00 - 12:08 | Open VisuAlgo hash table visualization. Step through 5-6 insertions with chaining, then switch to linear probing. Watch what happens when you delete in open addressing without tombstones (probe chains break). |
| 12:08 - 12:15 | Read the OVERVIEW.md Day 1 section. Re-read the complexity table. Draw the bucket structure on paper (or whiteboard): an array of slots, each slot holding a chain. |
| 12:15 - 12:20 | Skim the Dave Cheney article on Go's map internals. Focus on: the `hmap` struct, bucket size of 8, tophash array, and the two growth strategies (same-size for too many overflows vs. double-size for high load). |

**No code yet.** The goal is to enter the implementation phase with a clear mental model.

### 12:20 - 1:20 | Implement (60 min)

| Time | Activity |
|------|----------|
| 12:20 - 12:55 | **Build the chaining hash map** (35 min). Start with the struct and constructor, then `Put`, `Get`, `Delete`, `Len`, and `resize`. Write tests alongside each method. See Section 4 for exact signatures and test cases. |
| 12:55 - 13:20 | **Build the open addressing hash map** (25 min). Reuse the same API. Focus on linear probing, tombstone handling for Delete, and probe termination. This should go faster since the API is identical -- the interesting part is the probing logic. |

**Break point:** If you finish the chaining map by 12:55 and feel solid, take a 2-minute stretch before starting open addressing.

### 1:20 - 1:50 | Solidify (30 min)

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | **Stress-test both implementations.** Insert 10,000 random keys, verify all are retrievable. Delete half, verify the other half still works. Compare your map's behavior against Go's built-in `map` as an oracle. |
| 1:30 - 1:40 | **Edge cases.** Add tests for: empty map get/delete, inserting duplicate keys (update in place), resize triggers, getting after deleting in open addressing (tombstone chain). |
| 1:40 - 1:50 | **Pattern practice.** Implement Two Sum using your chaining hash map (not Go's built-in map). Then implement a frequency counter. This connects the data structure to how it's used in interviews. |

### 1:50 - 2:00 | Recap (10 min)

| Time | Activity |
|------|----------|
| 1:50 - 1:55 | Close all references. Write from memory: the complexity of Put, Get, Delete (average and worst), the load factor threshold, and why resize must re-hash all keys. |
| 1:55 - 2:00 | Write down one thing that surprised you or that you'd forget by tomorrow. Examples: "tombstones are needed because zeroing a slot breaks the probe chain" or "Go's map uses a tophash byte array for fast comparison before checking full keys." |

---

## 3. Core Concepts Deep Dive

### How Go's Built-In `map` Works Internally

Go's `map` is a hash table with chaining, but the implementation is more nuanced than textbook chaining:

**Structure:**
```
hmap {
    count      int           // number of live entries (len(m))
    B          uint8         // log2 of number of buckets (so 2^B buckets)
    hash0      uint32        // hash seed (randomized per map instance)
    buckets    unsafe.Pointer // array of 2^B bmap structs
    oldbuckets unsafe.Pointer // previous bucket array during growth
    ...
}

bmap (bucket) {
    tophash [8]uint8         // top 8 bits of hash for each entry
    keys    [8]keyType       // packed keys
    values  [8]valueType     // packed values
    overflow *bmap           // pointer to next overflow bucket
}
```

**Key design decisions:**

1. **Bucket size = 8.** Each bucket holds 8 key-value pairs before overflowing. This is a cache-line optimization: checking 8 tophash bytes (8 bytes total) fits in a single cache line and is faster than pointer-chasing through a linked list.

2. **Tophash optimization.** Before comparing full keys, Go checks the top 8 bits of the hash against the `tophash` array. This is a fast filter: if the tophash doesn't match, the full key comparison (which may involve string comparison) is skipped. Most bucket entries are rejected by tophash alone.

3. **Memory layout.** Keys are stored contiguously, then values contiguously (not interleaved as `key1, val1, key2, val2`). This avoids padding waste when key and value have different alignment requirements.

4. **Two growth strategies:**
   - **Double-size growth:** When load factor exceeds 6.5 (count / 2^B > 6.5), allocate 2^(B+1) buckets and incrementally evacuate entries.
   - **Same-size growth:** When there are too many overflow buckets (indicating poor distribution or many deletions leaving sparse buckets), allocate a new array of the same size and re-compact entries.

5. **Incremental evacuation.** During growth, old buckets aren't evacuated all at once (that would cause a latency spike). Instead, each `mapaccess` or `mapassign` call evacuates 1-2 old buckets. Both `oldbuckets` and `buckets` are checked during this transition period.

6. **Hash seed randomization.** Each map gets a random `hash0` seed at creation. This prevents hash-flooding attacks where an adversary crafts keys that all collide, degrading the map to O(n) per operation.

### The Math Behind Load Factor 0.75

The load factor alpha = n / m (entries / buckets) controls the tradeoff between space and time:

**For chaining:**
- The expected length of each chain is alpha.
- Expected cost of a successful search: 1 + alpha/2 (scan half the chain on average).
- Expected cost of an unsuccessful search: 1 + alpha (scan the whole chain).
- At alpha = 0.75, expected unsuccessful search cost is 1.75 comparisons -- fast enough.
- At alpha = 1.0, you average 2 comparisons per lookup. At alpha = 2.0, it's 3. The degradation is linear but graceful.

**For open addressing (linear probing):**
- Expected probes for a successful search: ~(1/2)(1 + 1/(1 - alpha)).
- Expected probes for an unsuccessful search: ~(1/2)(1 + 1/(1 - alpha)^2).
- At alpha = 0.75: successful ~2.5 probes, unsuccessful ~8.5 probes. Getting expensive.
- At alpha = 0.5: successful ~1.5 probes, unsuccessful ~2.5 probes. Much better.
- Open addressing degrades faster than chaining as alpha rises because **clustering** worsens -- runs of occupied slots grow and merge.

**Why 0.75?** It's a pragmatic choice, not a mathematical constant. Java's HashMap, Python's dict (which uses open addressing but resizes at 2/3), and most textbooks converge on the 0.65-0.80 range. The reasoning: below 0.5, you waste too much memory; above 0.8, collisions spike. 0.75 sits in the sweet spot where average lookup is ~1.5 comparisons for chaining and memory utilization is reasonable.

**Go's choice of 6.5:** This seems high, but remember Go packs 8 entries per bucket. The effective per-slot load is 6.5/8 = 0.8125 -- right in the standard range.

### Chaining vs. Open Addressing

| Factor | Chaining | Open Addressing |
|--------|----------|-----------------|
| **Collision handling** | Append to list/slice at bucket | Probe to next slot (linear, quadratic, double hash) |
| **Cache performance** | Poor -- pointer chasing through heap-allocated nodes | Good -- sequential memory access during probing |
| **Load factor tolerance** | Degrades gracefully past 1.0 | Performance cliff as alpha approaches 1.0 |
| **Deletion** | Simple -- remove node from chain | Complex -- requires tombstones to preserve probe chains |
| **Memory overhead** | Pointers per node (or slice headers per bucket) | Wastes empty slots; needs tombstone flags |
| **Implementation simplicity** | Simpler | More subtle edge cases |
| **When to choose** | When deletions are frequent, when load factor may spike, or when simplicity matters | When cache performance matters, when keys are small, when the dataset size is predictable |

**In practice:** Most general-purpose hash maps (Go, Java, Python < 3.6) use chaining or a chaining variant. Open addressing shines in specialized contexts: CPU caches, embedded systems, or when you control the workload (e.g., compiler symbol tables with known-size alphabets).

### Hash Function Properties

A good hash function for a hash table must have:

1. **Deterministic:** The same key always produces the same hash. This seems obvious, but floating-point keys can violate this if not handled carefully (NaN != NaN).

2. **Uniform distribution:** Hash values should spread evenly across the output range. Poor distribution means some buckets get overloaded while others sit empty. Test by histogramming bucket occupancy.

3. **Avalanche effect:** A small change in input (flipping one bit of the key) should change roughly half the bits of the output. Without this, similar keys (e.g., "user1", "user2", "user3") cluster into adjacent buckets.

4. **Fast to compute:** The hash function runs on every Put, Get, and Delete. It must be O(key length) and have low constant factors. Cryptographic hashes (SHA-256) have great distribution but are too slow for hash tables. FNV-1a and xxHash are designed for speed.

**FNV-1a (the one you'll use):**
```
hash = offset_basis (2166136261 for 32-bit)
for each byte in key:
    hash = hash XOR byte
    hash = hash * FNV_prime (16777619 for 32-bit)
return hash
```

It's simple, fast, and has good distribution for typical string keys. The XOR-then-multiply order (FNV-1a) has slightly better avalanche properties than multiply-then-XOR (FNV-1).

---

## 4. Implementation Checklist

### A. Chaining Hash Map

```go
type Entry struct {
    Key   string
    Value int
}

type ChainingMap struct {
    buckets [][]Entry
    size    int
    cap     int
}

func NewChainingMap(initialCap int) *ChainingMap
func (m *ChainingMap) Put(key string, value int)
func (m *ChainingMap) Get(key string) (int, bool)
func (m *ChainingMap) Delete(key string) bool
func (m *ChainingMap) Len() int
func (m *ChainingMap) LoadFactor() float64
func (m *ChainingMap) resize()
func (m *ChainingMap) hash(key string) int
```

**Implementation notes:**
- `hash(key)` should use `hash/fnv` (FNV-1a 32-bit), then `% m.cap` to get the bucket index.
- `Put`: hash the key, scan the bucket for an existing entry with the same key (update in place), otherwise append. After insertion, check load factor and resize if > 0.75.
- `Get`: hash the key, scan the bucket linearly, return `(value, true)` or `(0, false)`.
- `Delete`: hash the key, scan the bucket, remove the entry. Order within a bucket doesn't matter, so swap with the last element and shrink (avoids shifting).
- `resize`: allocate a new bucket array of 2x capacity. Re-hash every entry from every bucket into the new array. Update `m.cap`.

**Tests to write:**

```go
func TestChainingMap_PutAndGet(t *testing.T)
// Insert 5 key-value pairs. Get each one. Verify correct values.

func TestChainingMap_UpdateExistingKey(t *testing.T)
// Put("a", 1), then Put("a", 2). Get("a") should return 2. Len() should be 1.

func TestChainingMap_GetMissing(t *testing.T)
// Get on a key that was never inserted. Should return (0, false).

func TestChainingMap_Delete(t *testing.T)
// Insert, delete, verify Get returns false. Verify Len decrements.

func TestChainingMap_DeleteMissing(t *testing.T)
// Delete a key that doesn't exist. Should return false. Len unchanged.

func TestChainingMap_Resize(t *testing.T)
// Start with capacity 4. Insert enough keys to trigger resize.
// Verify all keys are still retrievable after resize.

func TestChainingMap_ManyKeys(t *testing.T)
// Insert 1000 keys ("key0" through "key999"). Verify all retrievable.
// Delete 500 of them. Verify remaining 500 still work.

func TestChainingMap_EmptyMapOperations(t *testing.T)
// Get, Delete, Len on a freshly created map.
```

**Edge cases to handle:**
- [ ] Initial capacity of 0 or negative (default to a reasonable minimum like 8)
- [ ] Updating an existing key does not change `size`
- [ ] Delete returns false for non-existent key
- [ ] Resize re-hashes all entries (old index != new index because `hash % cap` changes)
- [ ] Works correctly after multiple resize cycles
- [ ] Empty string as key (valid -- it should hash normally)

---

### B. Open Addressing Hash Map (Linear Probing)

```go
type OAEntry struct {
    Key     string
    Value   int
    occupied bool
    deleted  bool
}

type OpenAddressingMap struct {
    table []OAEntry
    size  int
    cap   int
}

func NewOpenAddressingMap(initialCap int) *OpenAddressingMap
func (m *OpenAddressingMap) Put(key string, value int)
func (m *OpenAddressingMap) Get(key string) (int, bool)
func (m *OpenAddressingMap) Delete(key string) bool
func (m *OpenAddressingMap) Len() int
func (m *OpenAddressingMap) resize()
func (m *OpenAddressingMap) hash(key string) int
```

**Implementation notes:**
- Each slot has three states: **empty** (`!occupied && !deleted`), **occupied** (`occupied && !deleted`), **tombstone** (`!occupied && deleted`, or use `deleted` flag on an occupied entry -- pick a convention and be consistent).
- `Put`: hash the key, probe forward. If you find the key, update in place. If you find an empty slot, insert. If you find a tombstone, you *can* reuse it but must continue probing to check for the key further along the chain (otherwise you'd create duplicates).
- `Get`: hash the key, probe forward. Skip tombstones (they're not empty, so don't stop). Stop at an empty slot or when you've checked all slots.
- `Delete`: find the key, mark it as a tombstone. Do **not** shift subsequent entries (that's robin hood hashing, a different technique). Do **not** just clear the slot (breaks probe chains).
- `resize`: allocate a new table of 2x capacity. Re-insert all occupied (non-tombstone) entries. Tombstones are discarded -- this is how the table self-cleans.

**Tests to write:**

```go
func TestOAMap_PutAndGet(t *testing.T)
// Same as chaining -- basic insert and retrieve.

func TestOAMap_LinearProbing(t *testing.T)
// Force a collision: use a small capacity (4), insert keys that hash
// to the same index. Verify all are retrievable.

func TestOAMap_DeleteWithTombstone(t *testing.T)
// Insert A, B, C where B and C collide with A.
// Delete B. Verify C is still findable (probe chain crosses tombstone).

func TestOAMap_UpdateExistingKey(t *testing.T)
// Put same key twice, verify value is updated and Len is unchanged.

func TestOAMap_InsertIntoTombstone(t *testing.T)
// Insert A, delete A (creates tombstone), insert B that hashes to same slot.
// B should occupy the tombstone slot. Verify A is gone and B is present.

func TestOAMap_Resize(t *testing.T)
// Fill past 75% load factor. Verify resize triggers and all keys survive.
// Verify tombstones are cleaned up (no tombstones in new table).

func TestOAMap_FullProbeLoop(t *testing.T)
// Ensure Get on a missing key terminates even in a nearly-full table.
```

**Edge cases to handle:**
- [ ] Three distinct slot states: empty, occupied, tombstone
- [ ] Get must skip tombstones, not stop at them
- [ ] Put must check for existing key *past* tombstones before inserting into a tombstone slot
- [ ] Resize discards tombstones (only re-insert occupied entries)
- [ ] Probe loop termination: stop when you hit an empty slot or have scanned all slots
- [ ] Load factor threshold should be lower than chaining (0.5-0.7 is safer for open addressing)
- [ ] Capacity must always be >= 1 to avoid division by zero in modular hashing

---

## 5. Common Interview Patterns Using Hash Tables

### Pattern 1: Complement / Pair Finding

**The idea:** You need to find two elements that combine to a target value. Instead of the O(n^2) nested loop, store each element in a hash map and check if its complement exists.

**Template:**
```go
seen := map[int]int{} // value -> index
for i, num := range nums {
    complement := target - num
    if j, ok := seen[complement]; ok {
        return []int{j, i}
    }
    seen[num] = i
}
```

**Where it shows up:** Two Sum, pair sums, finding if any two numbers differ by K.

### Pattern 2: Frequency Counting

**The idea:** Count occurrences of each element, then use the counts to answer the question.

**Template:**
```go
freq := map[string]int{}
for _, item := range items {
    freq[item]++
}
// Now iterate freq to find most/least common, check counts, etc.
```

**Where it shows up:** Top K frequent elements, valid anagram, first unique character, majority element, ransom note.

### Pattern 3: Grouping by Computed Key

**The idea:** Compute a canonical key for each element. Elements with the same key belong to the same group.

**Template:**
```go
groups := map[string][]string{}
for _, word := range words {
    key := computeKey(word) // e.g., sort the letters
    groups[key] = append(groups[key], word)
}
```

**Where it shows up:** Group anagrams (key = sorted letters), group shifted strings (key = difference pattern), group by frequency distribution.

### Pattern 4: Deduplication / Existence Tracking

**The idea:** Use a hash set (map with empty struct values) to track what you've already seen.

**Template:**
```go
seen := map[int]struct{}{}
for _, num := range nums {
    if _, exists := seen[num]; exists {
        // duplicate found
    }
    seen[num] = struct{}{}
}
```

**Where it shows up:** Contains duplicate, longest consecutive sequence, happy number (cycle detection), intersection of two arrays.

### Pattern 5: Two-Pass vs. One-Pass with Hash Map

**The idea:** Some problems can be solved in two passes (build the map, then query it) or one pass (query and build simultaneously). The one-pass approach is often trickier but more elegant.

**Two-pass example (Two Sum):**
```go
// Pass 1: build map
index := map[int]int{}
for i, num := range nums {
    index[num] = i
}
// Pass 2: query
for i, num := range nums {
    if j, ok := index[target - num]; ok && j != i {
        return []int{i, j}
    }
}
```

**One-pass example (Two Sum):**
```go
// Build and query simultaneously
seen := map[int]int{}
for i, num := range nums {
    if j, ok := seen[target - num]; ok {
        return []int{j, i}
    }
    seen[num] = i
}
```

**Tradeoff:** One-pass uses the same time and space complexity but avoids iterating twice. It also naturally handles the "don't use the same element twice" constraint since you only see elements that came *before* the current one.

**Where it shows up:** Two Sum (one-pass vs. two-pass), subarray sum equals K (prefix sum + map), continuous subarray sum.

### Pattern 6: Prefix Sum + Hash Map

**The idea:** Store prefix sums in a hash map to find subarrays with a target sum in O(n). If `prefix[j] - prefix[i] == target`, then the subarray `[i+1..j]` sums to target.

**Template:**
```go
prefixCount := map[int]int{0: 1} // base case: empty prefix
sum := 0
count := 0
for _, num := range nums {
    sum += num
    if c, ok := prefixCount[sum - target]; ok {
        count += c
    }
    prefixCount[sum]++
}
```

**Where it shows up:** Subarray sum equals K, contiguous array (equal 0s and 1s), longest subarray with sum K.

---

## 6. Self-Assessment

Answer these from memory at the end of the session. If you can't, that's your signal for what to revisit tomorrow.

### Question 1: Resize Mechanics

**What happens during a hash table resize, and why can't you just copy the backing array to a larger one?**

*Expected answer:* When the load factor exceeds the threshold, you allocate a new backing array (typically 2x the size). You cannot simply copy the array because each entry's bucket index is computed as `hash(key) % capacity`. When capacity changes, `hash(key) % new_capacity` produces a different index than `hash(key) % old_capacity` for most entries. You must re-hash every key and re-insert it into the new array.

### Question 2: Tombstones in Open Addressing

**Why do you need tombstones when deleting from an open addressing hash table? What goes wrong if you just clear the slot?**

*Expected answer:* In open addressing, entries that collided during insertion were placed in subsequent slots via probing. When you search for a key, you follow the same probe sequence and stop at the first empty slot. If you clear a deleted entry's slot (making it empty), you break the probe chain -- any entry that was inserted *after* the deleted one (and probed past it) becomes unreachable. A tombstone marks the slot as "deleted but keep probing," preserving the chain's integrity. Tombstones are cleaned up during resize.

### Question 3: Load Factor Tradeoff

**What happens if you set the load factor threshold too low (e.g., 0.25)? Too high (e.g., 0.95)?**

*Expected answer:* Too low (0.25): the table resizes very aggressively, wasting memory. Most buckets are empty. Amortized insertion cost increases because resizing (which is O(n)) happens more frequently relative to the number of insertions. Too high (0.95): collisions become very frequent. In chaining, chains grow long (average chain length ~0.95). In open addressing, clustering gets severe -- the expected number of probes for an unsuccessful search at 0.95 load is ~200 for linear probing. Performance degrades toward O(n) per operation.

### Question 4: Chaining vs. Open Addressing

**Give one scenario where chaining is clearly better than open addressing, and one where open addressing wins.**

*Expected answer:* Chaining is better when deletions are frequent -- chaining handles deletion naturally (just remove the node), while open addressing requires tombstones that accumulate and degrade performance until the next resize. Open addressing is better when the key-value pairs are small and the table fits in cache -- probing through contiguous memory is faster than chasing pointers to heap-allocated nodes. CPU cache hit rates dominate in practice.

### Question 5: Hash Map Patterns

**You're given an unsorted array of integers and need to find the length of the longest consecutive sequence (e.g., [100, 4, 200, 1, 3, 2] -> 4, because [1, 2, 3, 4]). How does a hash set give you O(n) time?**

*Expected answer:* Insert all elements into a hash set. For each element, check if `element - 1` is in the set. If it is, skip this element (it's not the start of a sequence). If it's not, this element is the start of a consecutive sequence -- count upward (`element + 1`, `element + 2`, ...) while the next value exists in the set. Each element is visited at most twice (once in the outer loop, once as part of a sequence count), so total work is O(n). The hash set gives O(1) existence checks, making this possible.
