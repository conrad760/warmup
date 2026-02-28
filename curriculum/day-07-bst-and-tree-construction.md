# Day 7: BST & Tree Construction

> **Time budget:** 2 hours
> **Prereqs:** Binary tree traversals, recursion, pointer manipulation in Go
> **Goal:** Internalize the BST-specific patterns interviewers love and build trees from scratch under pressure.

---

## Pattern Catalog

### 1. BST Validation (Min/Max Range Propagation)

**Trigger:** "Is this a valid BST?" or any problem where you need to confirm BST invariants hold.

**Core idea:** Every node has an allowed range `(min, max)`. The root starts with `(-inf, +inf)`. Going left tightens the max; going right tightens the min.

**Go Template:**

```go
func isValidBST(root *TreeNode) bool {
    return validate(root, math.MinInt64, math.MaxInt64)
}

func validate(node *TreeNode, min, max int) bool {
    if node == nil {
        return true
    }
    if node.Val <= min || node.Val >= max {
        return false
    }
    return validate(node.Left, min, node.Val) &&
        validate(node.Right, node.Val, max)
}
```

**Complexity:** O(n) time, O(h) space (call stack, where h = tree height).

**Watch out:**
- Using `int` min/max works for LeetCode constraints, but in an interview clarify the value range. If values can be `math.MinInt64`, pass `*int` pointers and use `nil` for unbounded.
- `<=` and `>=`, not `<` and `>`. BSTs typically require strict inequality (no duplicates), but ask your interviewer.

---

### 2. BST Search & Insert

**Trigger:** "Find a value" or "insert into a BST" -- any problem that leverages the sorted property to halve the search space.

**Core idea:** At each node, compare the target to `node.Val` and go left or right. No need to explore both subtrees.

**Go Template:**

```go
// Search
func searchBST(root *TreeNode, target int) *TreeNode {
    cur := root
    for cur != nil {
        if target == cur.Val {
            return cur
        } else if target < cur.Val {
            cur = cur.Left
        } else {
            cur = cur.Right
        }
    }
    return nil // not found
}

// Insert (returns new root)
func insertIntoBST(root *TreeNode, val int) *TreeNode {
    if root == nil {
        return &TreeNode{Val: val}
    }
    if val < root.Val {
        root.Left = insertIntoBST(root.Left, val)
    } else {
        root.Right = insertIntoBST(root.Right, val)
    }
    return root
}
```

**Complexity:** O(h) time and space for recursive insert, O(h) time and O(1) space for iterative search.

**Watch out:**
- Iterative insert is trickier because you need to track the parent. Recursive is cleaner in interviews.
- Always return the root from insert so the caller can update its pointer.

---

### 3. BST Delete (Three Cases)

**Trigger:** "Delete a node from a BST" or any mutation that must preserve BST ordering.

**Core idea:** Find the node, then handle three cases:
1. **Leaf:** Just remove it (return nil to parent).
2. **One child:** Bypass the node (return the non-nil child).
3. **Two children:** Find the in-order successor (smallest in right subtree), copy its value, then delete the successor from the right subtree.

**Go Template:**

```go
func deleteNode(root *TreeNode, key int) *TreeNode {
    if root == nil {
        return nil
    }
    if key < root.Val {
        root.Left = deleteNode(root.Left, key)
    } else if key > root.Val {
        root.Right = deleteNode(root.Right, key)
    } else {
        // Found the node to delete
        // Case 1 & 2: leaf or one child
        if root.Left == nil {
            return root.Right
        }
        if root.Right == nil {
            return root.Left
        }
        // Case 3: two children
        // Find in-order successor (leftmost in right subtree)
        successor := root.Right
        for successor.Left != nil {
            successor = successor.Left
        }
        root.Val = successor.Val
        // Delete the successor from the right subtree
        root.Right = deleteNode(root.Right, successor.Val)
    }
    return root
}
```

**Complexity:** O(h) time, O(h) space.

**Watch out:**
- The most common bug: copying the successor's value but forgetting to actually delete the successor node from the right subtree. That leaves a duplicate.
- You can use in-order predecessor (rightmost in left subtree) instead -- mention this to your interviewer to show depth.
- Cases 1 and 2 collapse elegantly: "if left is nil return right" handles both leaf (returns nil) and right-child-only cases.

---

### 4. In-Order Traversal = Sorted Output

**Trigger:** "Kth smallest," "kth largest," "BST iterator," "closest value," or anything that needs elements in sorted order from a BST.

**Core idea:** In-order traversal of a BST visits nodes in ascending order. You don't need to collect all values -- just count or stop early.

**Go Template (Kth Smallest):**

```go
func kthSmallest(root *TreeNode, k int) int {
    count := 0
    result := 0

    var inorder func(node *TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil || count >= k {
            return
        }
        inorder(node.Left)
        count++
        if count == k {
            result = node.Val
            return
        }
        inorder(node.Right)
    }

    inorder(root)
    return result
}
```

**Go Template (BST Iterator -- iterative with explicit stack):**

```go
type BSTIterator struct {
    stack []*TreeNode
}

func Constructor(root *TreeNode) BSTIterator {
    it := BSTIterator{}
    it.pushLeft(root)
    return it
}

func (it *BSTIterator) pushLeft(node *TreeNode) {
    for node != nil {
        it.stack = append(it.stack, node)
        node = node.Left
    }
}

func (it *BSTIterator) Next() int {
    top := it.stack[len(it.stack)-1]
    it.stack = it.stack[:len(it.stack)-1]
    it.pushLeft(top.Right)
    return top.Val
}

func (it *BSTIterator) HasNext() bool {
    return len(it.stack) > 0
}
```

**Complexity:** Kth smallest: O(h + k) time, O(h) space. Iterator: O(1) amortized per `Next()`, O(h) space.

**Watch out:**
- Kth smallest: is k 1-indexed or 0-indexed? LeetCode uses 1-indexed. Always clarify.
- Kth largest = (n - k + 1)th smallest, or do reverse in-order (right, root, left).
- The early termination `count >= k` check is critical for performance -- without it you traverse the entire tree.

---

### 5. BST from Sorted Input (Sorted Array to Balanced BST)

**Trigger:** "Convert sorted array to BST," "build a balanced BST," or any construction from sorted data.

**Core idea:** The middle element becomes the root (this guarantees balance). Recurse on the left half for the left subtree and the right half for the right subtree.

**Go Template:**

```go
func sortedArrayToBST(nums []int) *TreeNode {
    return build(nums, 0, len(nums)-1)
}

func build(nums []int, lo, hi int) *TreeNode {
    if lo > hi {
        return nil
    }
    mid := lo + (hi-lo)/2
    node := &TreeNode{Val: nums[mid]}
    node.Left = build(nums, lo, mid-1)
    node.Right = build(nums, mid+1, hi)
    return node
}
```

**Complexity:** O(n) time, O(log n) space (balanced recursion).

**Watch out:**
- `mid = lo + (hi - lo) / 2` not `(lo + hi) / 2` to avoid integer overflow (habit from binary search, interviewers notice).
- Inclusive bounds: `lo` to `hi` inclusive. Base case is `lo > hi`, NOT `lo >= hi` (that would skip single-element subtrees).
- If the problem says "sorted linked list to BST," you can still use this pattern but with a slow/fast pointer to find the middle, or use an in-order simulation trick for O(n) time.

---

### 6. Serialize / Deserialize a Binary Tree

**Trigger:** "Serialize and deserialize a binary tree," "encode/decode tree to string."

**Core idea:** Two main approaches:
- **Preorder with null markers:** Visit root, then left, then right. Write "null" for nil children. Deserialization reads tokens left-to-right and recurses.
- **BFS (level-order):** Queue-based. Good for visualizing but trickier to implement cleanly.

Preorder is simpler to code under pressure.

**Go Template (Preorder):**

```go
type Codec struct{}

func (c *Codec) serialize(root *TreeNode) string {
    var tokens []string
    var preorder func(node *TreeNode)
    preorder = func(node *TreeNode) {
        if node == nil {
            tokens = append(tokens, "N")
            return
        }
        tokens = append(tokens, strconv.Itoa(node.Val))
        preorder(node.Left)
        preorder(node.Right)
    }
    preorder(root)
    return strings.Join(tokens, ",")
}

func (c *Codec) deserialize(data string) *TreeNode {
    tokens := strings.Split(data, ",")
    idx := 0

    var build func() *TreeNode
    build = func() *TreeNode {
        if idx >= len(tokens) || tokens[idx] == "N" {
            idx++
            return nil
        }
        val, _ := strconv.Atoi(tokens[idx])
        idx++
        node := &TreeNode{Val: val}
        node.Left = build()
        node.Right = build()
        return node
    }

    return build()
}
```

**Complexity:** O(n) time and space for both serialize and deserialize.

**Watch out:**
- The `idx` variable must be shared across recursive calls. In Go, use a closure (as above) or pass `*int`. Do NOT pass `int` by value -- it won't update.
- Choose a delimiter that can't appear in values. Comma works for integers.
- Choose a null marker that can't be a valid value. "N" or "null" works.
- Edge case: empty tree. `serialize(nil)` should return `"N"`, and `deserialize("N")` should return `nil`.

---

### 7. LCA in a BST

**Trigger:** "Lowest common ancestor in a BST" -- simpler than the general binary tree version.

**Core idea:** Exploit the BST ordering. If both target values are less than current node, LCA is in the left subtree. If both are greater, LCA is in the right subtree. Otherwise, current node is the LCA (the split point).

**Go Template:**

```go
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    cur := root
    for cur != nil {
        if p.Val < cur.Val && q.Val < cur.Val {
            cur = cur.Left
        } else if p.Val > cur.Val && q.Val > cur.Val {
            cur = cur.Right
        } else {
            return cur // split point = LCA
        }
    }
    return nil
}
```

**Complexity:** O(h) time, O(1) space (iterative).

**Watch out:**
- This only works for BSTs. The general binary tree LCA requires checking both subtrees recursively -- O(n).
- The problem guarantees both p and q exist in the tree. If that's not guaranteed, you'd need to verify.
- The "split point" includes the case where one node IS the ancestor of the other (e.g., p.Val == cur.Val). The condition falls into the `else` branch correctly.

---

## Decision Framework

Use this to quickly identify which pattern to apply:

```
Problem mentions "BST"?
  |
  +--> Think: sorted property, in-order = sorted output
  |
  What is the operation?
  |
  +-- "Is it valid?" -----------> Pattern 1: min/max range validation
  |
  +-- "Find / Insert" ----------> Pattern 2: BST search & insert (O(h))
  |
  +-- "Delete a node" ----------> Pattern 3: three-case delete
  |
  +-- "Kth smallest/largest" ---> Pattern 4: in-order + early termination
  |   "Iterator / next element"
  |   "Closest value"
  |
  +-- "Build BST from ___" -----> Pattern 5: recursive midpoint subdivision
  |   "Sorted array to BST"
  |
  +-- "Serialize / Deserialize"-> Pattern 6: preorder + null markers
  |
  +-- "Lowest common ancestor"-> Pattern 7: compare both values, go left/right/done
```

**Key mental shortcuts:**
- "Validate BST" -> min/max range, NOT just checking immediate children
- "Kth smallest/largest" -> in-order traversal with a counter, stop at k
- "Build a tree from ___" -> find root, split remaining input, recurse
- "LCA in a BST" -> compare both values to current node and follow the split

---

## Common Interview Traps

### Trap 1: Validate BST -- Local vs. Global Constraint

**Wrong approach:** Check `node.Left.Val < node.Val && node.Right.Val > node.Val` at each node.

**Why it fails:**
```
      5
     / \
    1    6
        / \
       3   7     <-- 3 is less than 5, violates BST!
```
Node 6 correctly has 3 < 6, but 3 is in the right subtree of 5 and must be > 5. Local checks miss this.

**Fix:** Propagate `(min, max)` range from ancestors. Node 3 gets range `(5, 6)` and fails since 3 < 5.

---

### Trap 2: Delete -- Orphaned Successor

After finding a node with two children, you find the in-order successor, copy its value... and then forget to delete the successor from its original position. Now the tree has a duplicate.

**Fix:** After `root.Val = successor.Val`, always do `root.Right = deleteNode(root.Right, successor.Val)`.

---

### Trap 3: Kth Smallest -- Off-by-One

Some problems use 1-indexed k (LeetCode 230), others might use 0-indexed. If you increment your counter after visiting a node, check whether you compare `count == k` or `count == k-1`.

**Fix:** Always clarify indexing with the interviewer. Then trace through with k=1 on a single-node tree to sanity-check your code.

---

### Trap 4: Sorted Array to BST -- Bounds Error

Using `mid = (lo + hi) / 2` with `hi` as exclusive upper bound but treating it as inclusive (or vice versa). This creates unbalanced trees or misses elements.

**Fix:** Pick one convention and stick with it. Inclusive `[lo, hi]` is cleanest for this problem: base case `lo > hi`, mid picks `lo + (hi-lo)/2`, recurse on `[lo, mid-1]` and `[mid+1, hi]`.

---

### Trap 5: Serialize -- Shared State in Deserialization

In Go, passing the token index as a plain `int` to recursive calls means each call gets its own copy. The index never advances globally.

**Fix:** Use a closure that captures the index variable (as shown in Pattern 6), or pass `*int`.

---

## Thought Process Walkthrough

### Simulation 1: Validate BST (LeetCode 98)

> **Interviewer:** "Given the root of a binary tree, determine if it is a valid BST."

**Step 1 -- Clarify (30 seconds)**

"To confirm: a valid BST means every node in the left subtree is strictly less than the node, and every node in the right subtree is strictly greater. No duplicates allowed. The tree can be empty, and an empty tree is a valid BST. Correct?"

**Step 2 -- Identify the pattern (15 seconds)**

"This is BST validation. The key insight is that checking only immediate children is insufficient -- I need to propagate constraints from ancestors. I'll use min/max range propagation."

**Step 3 -- Describe the approach (1 minute)**

"I'll write a recursive helper that takes a node and the valid range `(min, max)` for that node. The root starts with `(-infinity, +infinity)`. When I go left, the max becomes the parent's value. When I go right, the min becomes the parent's value. If a node's value falls outside its range, return false."

**Step 4 -- Code (3-4 minutes)**

```go
func isValidBST(root *TreeNode) bool {
    return validate(root, math.MinInt64, math.MaxInt64)
}

func validate(node *TreeNode, min, max int) bool {
    if node == nil {
        return true
    }
    if node.Val <= min || node.Val >= max {
        return false
    }
    return validate(node.Left, min, node.Val) &&
        validate(node.Right, node.Val, max)
}
```

**Step 5 -- Trace through an example (1 minute)**

```
Tree:    5
        / \
       1    6
            / \
           3   7

validate(5, -inf, +inf)  -> 5 in range, continue
  validate(1, -inf, 5)   -> 1 in range, continue
    validate(nil, ...)    -> true
    validate(nil, ...)    -> true
  validate(6, 5, +inf)   -> 6 in range, continue
    validate(3, 5, 6)    -> 3 <= 5? YES -> return false!
```

"Correctly catches that 3 violates the global constraint."

**Step 6 -- Complexity**

"Time: O(n), we visit every node once. Space: O(h) for the recursion stack, where h is the height. Worst case O(n) for a skewed tree, O(log n) for a balanced tree."

**Step 7 -- Edge cases**

"Empty tree returns true. Single node returns true. All left or all right (skewed tree) still works because the range tightens on each call. If node values can be `math.MinInt64` or `math.MaxInt64`, I'd switch to using `*int` pointers with nil representing unbounded, but for typical interview constraints `int64` sentinels are fine."

---

### Simulation 2: Kth Smallest Element in a BST (LeetCode 230)

> **Interviewer:** "Given the root of a BST and an integer k, return the kth smallest value (1-indexed) in the tree."

**Step 1 -- Clarify (30 seconds)**

"So k=1 means the smallest element, k=2 is the second smallest, and so on. And k is guaranteed to be valid (1 <= k <= number of nodes). Correct?"

**Step 2 -- Identify the pattern (15 seconds)**

"In-order traversal of a BST produces values in ascending order. I just need to do an in-order traversal and stop at the kth element."

**Step 3 -- Describe the approach (45 seconds)**

"I'll do a recursive in-order traversal with a counter. Each time I 'visit' a node (between left and right recursive calls), I increment the counter. When the counter hits k, I record the value and short-circuit the rest of the traversal. The closure captures both the counter and result."

**Step 4 -- Code (3 minutes)**

```go
func kthSmallest(root *TreeNode, k int) int {
    count := 0
    result := 0

    var inorder func(node *TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil || count >= k {
            return
        }
        inorder(node.Left)
        count++
        if count == k {
            result = node.Val
            return
        }
        inorder(node.Right)
    }

    inorder(root)
    return result
}
```

**Step 5 -- Trace through an example (1 minute)**

```
Tree:    3
        / \
       1    4
        \
         2

k = 2 (want the 2nd smallest, which is 2)

inorder(3)
  inorder(1)
    inorder(nil) -> return
    count=1, k=2, not equal
    inorder(2)
      inorder(nil) -> return
      count=2, k=2, MATCH -> result=2, return
    return
  count >= k, skip the rest
return result = 2  (correct)
```

**Step 6 -- Complexity**

"Time: O(h + k). We descend h levels to reach the leftmost node, then visit k nodes. Space: O(h) for the recursion stack."

**Step 7 -- Follow-up discussion**

"If asked 'what if the BST is modified frequently and we need kth smallest often?' I'd augment each node with a count of nodes in its left subtree. Then finding kth smallest becomes O(h) without any traversal -- compare k to the left subtree count and go left or right accordingly."

---

## Time Targets

| Problem | Target | Notes |
|---|---|---|
| Validate BST | 8 min | Should be automatic by now |
| BST Search/Insert | 5 min | Simple O(h) traversal |
| BST Delete | 12 min | Three cases; take time to get successor right |
| Kth Smallest in BST | 8 min | In-order + counter |
| Sorted Array to BST | 8 min | Recursive midpoint split |
| Serialize/Deserialize | 15 min | More code; preorder approach is fastest to write |
| LCA in BST | 5 min | Three-way comparison, iterative |

---

## Quick Drill

Complete these five exercises in order. Time yourself.

### Drill 1: Validate BST (8 min)
Write `isValidBST` from scratch. Test against this tree (should return false):
```
      10
     /  \
    5    15
        /  \
       6   20
```
Node 6 is in the right subtree of 10 but is less than 10.

### Drill 2: Kth Smallest (8 min)
Write `kthSmallest` using in-order traversal with early termination. Test with k=3 on:
```
      5
     / \
    3    6
   / \
  2    4
 /
1
```
Expected output: 3.

### Drill 3: Sorted Array to Balanced BST (8 min)
Write `sortedArrayToBST` for input `[1, 2, 3, 4, 5, 6, 7]`. The root should be 4. Verify the tree height is 3 (balanced).

### Drill 4: Delete Node in BST (12 min)
Write `deleteNode`. Test by deleting 3 from:
```
      5
     / \
    3    6
   / \    \
  2    4    7
```
Node 3 has two children. The in-order successor is 4. After deletion, 4 should replace 3.

### Drill 5: LCA in BST (5 min)
Write iterative `lowestCommonAncestor`. Test with p=2, q=4 on the tree above (after drill 4's deletion). LCA should be 4 (since 2 is now left child of 4).

---

## Self-Assessment

Answer these without looking at the guide. If you can't answer confidently, revisit that pattern.

1. **Why is checking only `left.Val < root.Val < right.Val` insufficient for BST validation?** What specific type of invalid tree does it miss?

2. **In BST delete with two children, why do we pick the in-order successor specifically?** Could we pick a different node? What property must that node have?

3. **What is the time complexity of kth smallest in a BST?** Why is it O(h + k) and not just O(k)?

4. **When building a BST from a sorted array, what happens if you always pick the first element as the root instead of the middle?** What does the resulting tree look like?

5. **BST LCA vs. General Binary Tree LCA: what is the complexity difference and why?** What BST property makes BST LCA faster?

### Answer Key

1. Local-only checks miss violations where a node satisfies its parent's constraint but violates an ancestor's. Example: a node in the right subtree of the root that is less than the root. The min/max range must be propagated from all ancestors.

2. The in-order successor is the smallest value greater than the deleted node. It can safely replace the deleted node without violating BST ordering with respect to either subtree. You could alternatively use the in-order predecessor (largest value smaller than the node). The chosen node must be adjacent in the sorted order.

3. O(h + k): you must first descend h levels to reach the leftmost node (the smallest), then visit k nodes in order. You can't jump directly to the kth node without traversing down first.

4. You get a completely right-skewed tree (a linked list). Every subsequent element is larger, so each one becomes the right child. Height = n, all operations become O(n). Picking the middle guarantees O(log n) height.

5. BST LCA is O(h) time, O(1) space (iterative). General tree LCA is O(n) time, O(h) space (must search both subtrees). The BST's ordering property lets you decide at each node whether to go left, right, or stop -- you never need to explore both subtrees.

---

*End of Day 7. If BST delete and serialize/deserialize feel shaky, spend an extra 15 minutes drilling just those two -- they have the most moving parts and the most interview bugs.*
