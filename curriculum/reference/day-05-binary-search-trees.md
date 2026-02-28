# Day 5 — Binary Search Trees: Deep Dive

## 1. Curated Learning Resources

| # | Resource | Type | Why It's Useful |
|---|----------|------|-----------------|
| 1 | [VisuAlgo — BST Visualization](https://visualgo.net/en/bst) | Interactive | Step-by-step animations of insert, delete (all 3 cases), and search. Toggle between operations and watch pointers move. Best single resource for building intuition. |
| 2 | [Binary Search Tree — Go Implementation (yourbasic.org)](https://yourbasic.org/algorithms/binary-search-tree/) | Article + Go code | Clean, idiomatic Go BST implementation covering insert, search, and traversal. Good reference for Go-specific patterns (returning modified subtree roots). |
| 3 | [LeetCode 98: Validate BST — Editorial](https://leetcode.com/problems/validate-binary-search-tree/editorial/) | Problem + explanation | Walks through the common mistake (checking only immediate children) and the correct min/max range approach. |
| 4 | [Back to Back SWE — Delete a Node From a BST](https://www.youtube.com/watch?v=gcULXE7ViZw) | Video (15 min) | Focused deep-dive into the 3 deletion cases with whiteboard diagrams. Clear treatment of the in-order successor replacement. |
| 5 | [Self-Balancing BSTs — AVL and Red-Black Trees Overview (Abdul Bari)](https://www.youtube.com/watch?v=vRwi_UcZGjU) | Video (20 min) | Conceptual overview of why self-balancing matters, what rotations do, and how AVL differs from Red-Black. Not an implementation guide — exactly the right depth for interviews. |
| 6 | [Go Data Structures: Binary Search Tree (flaviocopes.com)](https://flaviocopes.com/golang-data-structure-binary-search-tree/) | Tutorial + Go code | Full Go BST with Insert, Search, Min, Max, and traversals. Good for comparing your implementation approach. |
| 7 | [LeetCode 230: Kth Smallest Element in a BST](https://leetcode.com/problems/kth-smallest-element-in-a-bst/) | Problem | Practice the in-order traversal = sorted order property. Try both recursive (counter) and iterative (stack) solutions. |
| 8 | [Red-Black Trees in 5 Minutes (Spanning Tree)](https://www.youtube.com/watch?v=qvZGUFHWChY) | Video (5 min) | If you want a rapid-fire conceptual overview of Red-Black trees — the 5 rules, why they guarantee O(log n), where they're used in practice (Java TreeMap, C++ std::map). |

---

## 2. Detailed 2-Hour Session Plan

### 12:00–12:20 — Review & Internalize (20 min)

| Time | Activity |
|------|----------|
| 12:00–12:08 | Read the BST invariant section in OVERVIEW.md. Draw a 7-node BST on paper (e.g., insert 8, 3, 10, 1, 6, 14, 4). Verify the invariant holds at every node. Write out the in-order traversal — confirm it's sorted. |
| 12:08–12:14 | Study the complexity table. For each operation (search, insert, delete, min, max), trace the path through your drawn tree and count the comparisons. Understand why it's O(h), not O(log n) — the distinction matters. |
| 12:14–12:20 | Read through the three deletion cases on paper. For your drawn tree, walk through deleting: a leaf (4), a node with one child (14), and a node with two children (3). Identify the in-order successor in each two-children case. |

### 12:20–1:20 — Implement from Scratch (60 min)

| Time | Activity |
|------|----------|
| 12:20–12:28 | Define `TreeNode` struct and `BST` struct with `Root`. Implement `Insert` — use the "return modified subtree root" pattern: `func insertNode(root *TreeNode, val int) *TreeNode`. Test: insert 8, 3, 10, 1, 6, 14, 4 and verify with in-order traversal. |
| 12:28–12:34 | Implement `Search(val int) bool`. Both recursive and iterative versions. Test: search for existing values (6) and missing values (7). |
| 12:34–12:38 | Implement `Min()` and `Max()`. Follow left/right pointers to the bottom. Handle empty tree edge case. |
| 12:38–12:58 | **Implement `Delete` — the core challenge.** Take it case by case. Start with the helper `func deleteNode(root *TreeNode, val int) *TreeNode`. Case 1 (leaf): return nil. Case 2 (one child): return the child. Case 3 (two children): find in-order successor, copy its value, recursively delete the successor. Test each case individually. This will take the most time — that's expected. |
| 12:58–1:08 | Implement `IsValidBST()`. Use the min/max range approach: `func isValid(node *TreeNode, min, max int) bool`. Start with `math.MinInt64` and `math.MaxInt64`. Test with a valid BST and a manually constructed invalid one (swap two nodes). |
| 1:08–1:20 | Implement `KthSmallest(k int) int`. Use iterative in-order traversal with an explicit stack, decrementing k. When k reaches 0, return the current value. Test with k=1 (should be min), k=n (should be max), k out of bounds. |

### 1:20–1:50 — Solidify (30 min)

| Time | Activity |
|------|----------|
| 1:20–1:30 | **Edge cases and stress testing.** Test delete on: root node with two children, single-node tree, deleting a value that doesn't exist. Test IsValidBST on: empty tree (valid), single node (valid), tree where left child is valid but deeper left subtree violates invariant. |
| 1:30–1:38 | Implement `InorderSuccessor(root *TreeNode, target int) *TreeNode`. Two cases: if the node has a right subtree, it's the min of that subtree. Otherwise, trace back from the root tracking the last left turn. |
| 1:38–1:46 | Implement `SortedArrayToBST(nums []int) *TreeNode`. Pick the middle element as root, recursively build left and right subtrees from the two halves. This produces a height-balanced BST. |
| 1:46–1:50 | Implement `LowestCommonAncestor(root *TreeNode, p, q int) *TreeNode`. Exploit the BST property: if both values are less than root, go left. If both are greater, go right. Otherwise, root is the LCA. |

### 1:50–2:00 — Recap (10 min)

| Time | Activity |
|------|----------|
| 1:50–1:55 | From memory, write the complexity of: Search, Insert, Delete, Min, Max, IsValidBST, KthSmallest. Note both average (balanced) and worst (skewed). |
| 1:55–2:00 | Write down one gotcha for each: Delete (must actually remove successor from its original position), IsValidBST (checking only immediate children is wrong), KthSmallest (1-indexed vs 0-indexed). |

---

## 3. Core Concepts Deep Dive

### The BST Invariant and Why It Enables O(log n)

The BST invariant: for every node `N`, all values in `N.Left` are strictly less than `N.Val`, and all values in `N.Right` are strictly greater than `N.Val`.

This invariant turns the tree into a binary decision structure — at each node you can eliminate an entire subtree from consideration. If the tree is balanced (each level roughly halves the remaining nodes), you get O(log n) operations. This is the same principle as binary search on a sorted array, but embedded in a linked structure that supports efficient insertion and deletion.

The critical word is **balanced**. The invariant alone gives O(h) operations, where h is the tree height. Only when h = O(log n) do you get O(log n) operations.

### Delete Operation: All Three Cases

Deletion is the most complex BST operation because removing a node must preserve the invariant.

**Case 1 — Leaf node (no children):**
Simply remove the node. Return nil to the parent.

```
Delete 4:
        8                    8
       / \                  / \
      3   10       →       3   10
     / \    \             / \    \
    1   6   14           1   6   14
       /
      4
```

**Case 2 — One child:**
Replace the node with its single child. The child's subtree already satisfies the invariant relative to the deleted node's parent.

```
Delete 10:
        8                    8
       / \                  / \
      3   10       →       3   14
     / \    \             / \
    1   6   14           1   6
```

**Case 3 — Two children:**
This is the tricky one. You can't just remove the node — both subtrees need a parent.

**Strategy:** Replace the node's value with its **in-order successor** (the smallest value in the right subtree), then delete the successor from the right subtree. The successor has at most one child (a right child), so deleting it is Case 1 or Case 2.

```
Delete 3 (two children):
  1. Find in-order successor of 3 → it's 4 (min of right subtree rooted at 6)
  2. Copy 4 into the node where 3 was
  3. Delete 4 from the right subtree (leaf — Case 1)

        8                    8
       / \                  / \
      3   10       →       4   10
     / \    \             / \    \
    1   6   14           1   6   14
       /
      4
```

**In-order successor vs. predecessor — why either works:**
- **Successor**: smallest value in the right subtree. Always >= all left subtree values, and <= all right subtree values (after removing it from the right subtree). Invariant preserved.
- **Predecessor**: largest value in the left subtree. Same logic, mirrored.

Both produce a valid BST. Convention varies — most implementations use the successor. Using the predecessor is equally correct. In a real self-balancing tree, alternating between the two can help maintain balance.

### BST Validation: Why Checking Only Immediate Children Fails

A common mistake:

```go
// WRONG: only checks immediate children
func isValidBST(node *TreeNode) bool {
    if node == nil { return true }
    if node.Left != nil && node.Left.Val >= node.Val { return false }
    if node.Right != nil && node.Right.Val <= node.Val { return false }
    return isValidBST(node.Left) && isValidBST(node.Right)
}
```

This passes the following invalid tree:

```
        5
       / \
      1   7
         / \
        3   8     ← 3 is less than 5, but it's in the RIGHT subtree of 5!
```

Each node looks valid relative to its immediate parent, but 3 violates the invariant with respect to the root (5).

**The correct approach — propagate min/max ranges:**

```go
func isValid(node *TreeNode, min, max int) bool {
    if node == nil { return true }
    if node.Val <= min || node.Val >= max { return false }
    return isValid(node.Left, min, node.Val) &&
           isValid(node.Right, node.Val, max)
}

// Initial call:
isValid(root, math.MinInt64, math.MaxInt64)
```

Every node must fall within a valid range, and the range narrows as you descend:

```
        5 (range: -∞, +∞)
       / \
      1   7
(range: -∞, 5)  (range: 5, +∞)
         / \
        3   8
   (5, 7) ← 3 is NOT in (5, 7), so INVALID!
```

**Alternative approach:** In-order traversal should produce strictly increasing values. Track the previous value and verify each new value is greater.

### Balanced vs. Unbalanced: What Sorted Insertions Do

If you insert elements in sorted order (1, 2, 3, 4, 5, 6, 7), the BST degenerates into a linked list:

```
Sorted insertions: 1, 2, 3, 4, 5

1
 \
  2
   \
    3
     \
      4
       \
        5

Height = 4 (n - 1)
Search for 5: 4 comparisons = O(n)
```

The same values inserted in a balanced order (4, 2, 6, 1, 3, 5, 7) produce:

```
Balanced insertions: 4, 2, 6, 1, 3, 5, 7

        4
       / \
      2   6
     / \ / \
    1  3 5  7

Height = 2 (log₂ 7 ≈ 2.8)
Search for any value: ≤ 2 comparisons = O(log n)
```

This is why self-balancing trees exist — they prevent this degeneration regardless of insertion order.

### Self-Balancing Trees: AVL and Red-Black (Conceptual)

You won't implement these in interviews, but you should know what they are and why they matter.

**AVL Trees (Adelson-Velsky and Landis, 1962)**
- **Invariant:** For every node, the heights of the left and right subtrees differ by at most 1.
- **Mechanism:** After every insert/delete, check the balance factor (left height - right height). If it's -2 or +2, perform **rotations** (single or double) to restore balance.
- **Four rotation cases:** Left-Left, Right-Right, Left-Right, Right-Left.
- **Guarantees:** Strictly balanced. Height is always ≤ 1.44 * log₂(n+2). Search is fast.
- **Trade-off:** More rotations on insert/delete compared to Red-Black. Better for read-heavy workloads.

**Red-Black Trees**
- **Invariant:** Each node is colored red or black. Five rules:
  1. Every node is red or black.
  2. Root is black.
  3. Every nil leaf is black.
  4. Red nodes have only black children (no two reds in a row).
  5. Every path from root to nil leaf has the same number of black nodes.
- **Guarantees:** Height ≤ 2 * log₂(n+1). Less strict than AVL, so slightly taller trees but fewer rotations on mutations.
- **Trade-off:** Better for write-heavy workloads. Simpler to implement (relatively).

**Where they're used:**
- Java's `TreeMap` and `TreeSet` — Red-Black tree.
- C++'s `std::map` and `std::set` — typically Red-Black tree.
- Go's built-in `map` — **not** a BST. It's a hash map. Go's standard library has no built-in ordered map (though `btree` packages exist in the ecosystem).
- Linux kernel's scheduling — Red-Black tree (CFS scheduler).
- Databases (B-trees, a generalization) — disk-friendly balanced trees.

### In-Order Successor and Predecessor

**In-order successor** of a node N: the node with the smallest value greater than N.Val.

Two cases:
1. **N has a right subtree:** Successor is the leftmost node in the right subtree (the minimum of the right subtree).
2. **N has no right subtree:** Successor is the nearest ancestor for which N is in the left subtree. Walk up from N; the first time you go "up and to the left," that ancestor is the successor.

```
        20
       /  \
     10    30
    /  \
   5   15
      /  \
    12    18

Successor of 15: → has right subtree → min of right subtree = 18
Successor of 18: → no right subtree → walk up: 18→15 (right child, keep going),
                   15→10 (right child, keep going), 10→20 (LEFT child, stop) → 20
Successor of 5:  → no right subtree → walk up: 5→10 (LEFT child, stop) → 10
```

**In-order predecessor** of N: the node with the largest value less than N.Val. Mirror logic — max of left subtree, or nearest ancestor where N is in the right subtree.

**Why they matter:** Successor/predecessor are used in BST deletion (Case 3), range queries, and iterator implementations (moving to the next/previous element in sorted order).

---

## 4. Implementation Checklist

### Struct Definitions

```go
type TreeNode struct {
    Val         int
    Left, Right *TreeNode
}

type BST struct {
    Root *TreeNode
}
```

### Function Signatures

```go
// Core operations
func (b *BST) Insert(val int)
func (b *BST) Search(val int) bool
func (b *BST) Delete(val int)
func (b *BST) Min() (int, bool)          // value, exists
func (b *BST) Max() (int, bool)          // value, exists

// Validation & queries
func (b *BST) IsValidBST() bool
func (b *BST) KthSmallest(k int) (int, bool)  // value, found
func InorderSuccessor(root *TreeNode, target int) *TreeNode

// Helpers (internal, lowercase)
func insertNode(root *TreeNode, val int) *TreeNode
func deleteNode(root *TreeNode, val int) *TreeNode
func findMin(root *TreeNode) *TreeNode
func isValid(node *TreeNode, min, max int) bool
```

### Test Cases and Edge Cases

**Insert:**
- Insert into empty tree → becomes root
- Insert smaller value → goes left
- Insert larger value → goes right
- Insert duplicate → define behavior (reject, or go left/right — document your choice)
- Insert sorted sequence → verify it works (even though it creates a skewed tree)

**Search:**
- Search in empty tree → false
- Search for root value → true
- Search for leaf value → true
- Search for value not in tree → false
- Search after delete → false

**Delete:**
- Delete leaf node → parent's child becomes nil
- Delete node with one child (left only) → replaced by left child
- Delete node with one child (right only) → replaced by right child
- Delete node with two children → replaced by in-order successor, successor removed
- Delete root node (all three cases)
- Delete from single-node tree → tree becomes empty
- Delete value not in tree → tree unchanged
- Delete all nodes one by one → tree is empty, all operations still work

**Min / Max:**
- Empty tree → return false for exists
- Single node → min == max == root.Val
- Skewed left tree → min is deepest, max is root
- After deleting the current min/max → new min/max is correct

**IsValidBST:**
- Empty tree → true
- Single node → true
- Valid balanced BST → true
- Invalid: left child > parent → false
- Invalid: value in left subtree > ancestor (the "deep violation" case) → false
- Tree with `math.MinInt64` or `math.MaxInt64` values → handle boundary

**KthSmallest:**
- k = 1 → minimum value
- k = n (total nodes) → maximum value
- k > n → return false for found
- k <= 0 → return false for found
- Single-node tree, k = 1 → root value

**InorderSuccessor:**
- Node with right subtree → leftmost of right subtree
- Node without right subtree → nearest ancestor via left turn
- Successor of max node → nil (no successor)
- Single-node tree → nil

---

## 5. BST-Specific Patterns

### In-Order Traversal = Sorted Output

This is the fundamental BST property and the basis for nearly every BST problem.

```go
func inorder(node *TreeNode, result *[]int) {
    if node == nil { return }
    inorder(node.Left, result)
    *result = append(*result, node.Val)
    inorder(node.Right, result)
}
```

The output is guaranteed to be in ascending order for a valid BST. This means:
- **Kth smallest** = kth element in the in-order traversal.
- **Validate BST** = check that in-order traversal is strictly increasing.
- **Convert BST to sorted array** = in-order traversal.
- **Find closest value** = in-order traversal with early termination.

### LCA in a BST (Simpler Than General Tree LCA)

In a general binary tree, LCA requires checking both subtrees and returning the node where p and q split. In a BST, the ordering property gives you a direct path:

```go
func lowestCommonAncestor(root *TreeNode, p, q int) *TreeNode {
    node := root
    for node != nil {
        if p < node.Val && q < node.Val {
            node = node.Left           // both in left subtree
        } else if p > node.Val && q > node.Val {
            node = node.Right          // both in right subtree
        } else {
            return node                // split point — this is the LCA
        }
    }
    return nil
}
```

Time: O(h). No need to visit both subtrees or use recursion with merge logic.

### Range Queries and Pruning

The BST invariant lets you skip entire subtrees when searching within a range [lo, hi]:

```go
func rangeQuery(node *TreeNode, lo, hi int, result *[]int) {
    if node == nil { return }
    if node.Val > lo {
        rangeQuery(node.Left, lo, hi, result)   // only go left if there might be values >= lo
    }
    if node.Val >= lo && node.Val <= hi {
        *result = append(*result, node.Val)
    }
    if node.Val < hi {
        rangeQuery(node.Right, lo, hi, result)  // only go right if there might be values <= hi
    }
}
```

This prunes branches that can't contain values in the range. In a balanced BST, range queries retrieving k results run in O(log n + k) time.

### Converting Sorted Array to Balanced BST

Given a sorted array, build a height-balanced BST by always choosing the middle element as the root:

```go
func sortedArrayToBST(nums []int) *TreeNode {
    if len(nums) == 0 { return nil }
    mid := len(nums) / 2
    return &TreeNode{
        Val:   nums[mid],
        Left:  sortedArrayToBST(nums[:mid]),
        Right: sortedArrayToBST(nums[mid+1:]),
    }
}
```

This is the reverse of "in-order traversal = sorted array." It guarantees a balanced tree (height = floor(log₂ n)).

---

## 6. Visual Diagrams

### BST Deletion — All Three Cases

Starting tree:

```
           15
          /  \
        10    20
       /  \   / \
      5   12 18  25
         /     \
        11      19
```

**Case 1 — Delete leaf node (11):**

```
           15                          15
          /  \                        /  \
        10    20          →         10    20
       /  \   / \                  /  \   / \
      5   12 18  25               5   12 18  25
         /     \                         \
        11      19                        19

  Node 11 has no children → simply remove it.
```

**Case 2 — Delete node with one child (18):**

```
           15                          15
          /  \                        /  \
        10    20          →         10    20
       /  \   / \                  /  \   / \
      5   12 18  25               5   12 19  25
               \
                19

  Node 18 has one child (19) → replace 18 with 19.
```

**Case 3 — Delete node with two children (10):**

```
           15                          15
          /  \                        /  \
        10    20          →         11    20
       /  \   / \                  /  \   / \
      5   12 18  25               5   12 18  25
         /     \                         \
        11      19                        19

  Node 10 has two children.
  In-order successor of 10 = 11 (leftmost in right subtree of 10).
  Copy 11 into 10's position, delete original 11 (Case 1).
```

### Skewed vs. Balanced BST

Same seven values: 1, 2, 3, 4, 5, 6, 7

**Sorted insertions (1, 2, 3, 4, 5, 6, 7) → Skewed (linked list):**

```
  1
   \
    2
     \
      3
       \
        4
         \
          5
           \
            6
             \
              7

  Height: 6 (n - 1)
  Search for 7: 6 comparisons
  All operations: O(n)
```

**Balanced insertions (4, 2, 6, 1, 3, 5, 7) → Balanced:**

```
           4
          / \
        2     6
       / \   / \
      1   3 5   7

  Height: 2 (floor(log₂ 7))
  Search for any value: ≤ 2 comparisons
  All operations: O(log n)
```

**Using `sortedArrayToBST([1,2,3,4,5,6,7])` → Balanced:**

```
  Pick middle (4) as root
  Left half [1,2,3] → pick 2, with 1 left, 3 right
  Right half [5,6,7] → pick 6, with 5 left, 7 right

           4
          / \
        2     6
       / \   / \
      1   3 5   7

  Same balanced structure.
```

### Validate BST — Range Propagation

```
           8   ← range: (-∞, +∞)   ✓ (8 is in range)
          / \
         3   10  ← left range: (-∞, 8), right range: (8, +∞)
        / \    \
       1   6   14  ← ranges narrow further
      (-∞,3) (3,8) (10,+∞)

  Each node must fall within its inherited range.
  Range tightens at every level:
    - Going LEFT: upper bound becomes parent's value
    - Going RIGHT: lower bound becomes parent's value

  Invalid example — 9 placed in left subtree of 8:

           8   ← range: (-∞, +∞)
          / \
         3   10
        / \
       1   9  ← range should be (3, 8) but 9 > 8 → INVALID!

  Checking only immediate children:
    9 > 3 ✓ (valid as right child of 3)
  But checking range:
    9 NOT IN (3, 8) ✗ → correctly caught!
```

---

## 7. Self-Assessment

Answer these without looking at your code or notes. If you can't answer confidently, revisit that section.

### Question 1
**After deleting a node with two children using the in-order successor, is the resulting tree still a valid BST? Why?**

<details>
<summary>Answer</summary>

Yes. The in-order successor is the smallest value greater than the deleted node. By definition, it is: (a) greater than everything in the deleted node's left subtree, and (b) less than or equal to everything remaining in the deleted node's right subtree (since it was the minimum there). After copying the successor's value into the deleted node's position and removing the successor from the right subtree (which is a Case 1 or Case 2 delete), the invariant is preserved at every node. The same argument holds symmetrically for the in-order predecessor.
</details>

### Question 2
**You insert the values [1, 2, 3, 4, 5] into an empty BST in that order. What does the tree look like? What is the time complexity of searching for 5? How would you fix this?**

<details>
<summary>Answer</summary>

The tree is a right-skewed chain (linked list): 1→2→3→4→5. Searching for 5 requires 5 comparisons — O(n). To fix this, either: (a) insert in a balanced order (e.g., 3, 1, 4, 2, 5), (b) build from a sorted array using the middle-element technique (`sortedArrayToBST`), or (c) use a self-balancing BST (AVL or Red-Black) that automatically rebalances after each insertion.
</details>

### Question 3
**Your `IsValidBST` function checks `node.Left.Val < node.Val` and `node.Right.Val > node.Val` at each node and recursively validates both subtrees. A colleague claims this is sufficient. Construct a counterexample.**

<details>
<summary>Answer</summary>

```
      5
     / \
    1   7
       / \
      3   8
```

At each node, the immediate child check passes: 1 < 5, 7 > 5, 3 < 7, 8 > 7. But the tree is invalid because 3 is in the right subtree of 5, yet 3 < 5. The correct approach propagates a (min, max) range down the tree: when going right from 5, the range becomes (5, +∞). Node 7 is in (5, +∞) ✓. Going left from 7, the range becomes (5, 7). Node 3 is NOT in (5, 7) ✗.
</details>

### Question 4
**In your BST Delete implementation for the two-children case, you find the in-order successor and copy its value. What happens if you forget to then delete the successor from the right subtree?**

<details>
<summary>Answer</summary>

The successor's value now appears twice in the tree — once in its original position and once in the position of the deleted node. This violates the BST invariant (assuming no duplicates are allowed). Even if duplicates are allowed, the in-order traversal would contain a repeated value that wasn't originally present, corrupting the tree's logical content. Always complete the second step: after copying the successor's value, recursively delete the successor from the right subtree.
</details>

### Question 5
**Go's built-in `map` is a hash map, not a BST. Java's `TreeMap` is a Red-Black tree. When would you choose a tree-based map over a hash map?**

<details>
<summary>Answer</summary>

Choose a tree-based map (like `TreeMap`) when you need **ordered operations**: iterating keys in sorted order, finding the nearest key to a given value (floor/ceiling), range queries (all keys between A and B), or finding the min/max key. A hash map gives O(1) average for get/put but provides no ordering. A balanced BST gives O(log n) for get/put but also O(log n) for ordered queries that a hash map can't do at all (or would require O(n log n) by sorting the keys each time). If you only need get/put/delete, use a hash map. If you need any kind of ordered access, use a tree.
</details>
