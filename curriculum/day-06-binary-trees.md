# Day 6: Binary Trees (DFS & BFS)

> **Goal:** Given a binary tree problem, instantly recognize which traversal
> strategy and which information-flow direction (top-down vs bottom-up) to use,
> then write clean Go code under interview time pressure.
>
> **Time budget:** 2 hours | **Prereqs:** You know preorder/inorder/postorder/level-order mechanics

---

## Pattern Catalog

### Pattern 1: Top-Down DFS (Pass Info Down via Parameters)

**Trigger:** The problem says "check / accumulate something along the path from root toward leaves." The parent knows something the child needs (a running sum, a depth counter, a valid range).

**Go Template:**

```go
// Top-down: parent passes context to children.
// The "answer" is typically updated at leaf nodes or at every node.
func topDown(node *TreeNode, infoFromParent int) {
    if node == nil {
        return
    }
    // Use infoFromParent + node.Val to decide something
    // e.g., check a condition, update a global answer

    topDown(node.Left, updatedInfo)
    topDown(node.Right, updatedInfo)
}
```

**Classic problems:**
- **Max Depth** (pass depth down, track max at leaves)
- **Path Sum** (pass remaining target down, check == 0 at leaf)
- **Same Tree** (pass paired nodes down, compare at each step)
- **Validate BST** (pass valid range `[lo, hi]` down)

**Complexity:** O(n) time, O(h) space where h = tree height (recursion stack).

**Watch out:**
- You must handle the nil child before accessing `.Val`.
- "Leaf" means `node.Left == nil && node.Right == nil` -- don't confuse it with "nil node."
- If you need the answer returned up, top-down is the wrong pattern. Top-down pushes info *down* and usually writes to an external variable or collects at leaves.

---

### Pattern 2: Bottom-Up DFS (Build Answer from Children)

**Trigger:** The answer for a node depends on answers from its subtrees. You need to *combine* child results. Think: "If I knew the answer for left and right subtrees, can I compute it for this node?"

**Go Template:**

```go
// Bottom-up: children return info to the parent.
// The function's return value carries data UPWARD.
var globalAnswer int // often the real answer lives here, not in the return value

func bottomUp(node *TreeNode) int {
    if node == nil {
        return 0 // base case: what does an empty subtree contribute?
    }
    left := bottomUp(node.Left)
    right := bottomUp(node.Right)

    // Combine: update globalAnswer using left + right + node.Val
    // Return: what the PARENT needs (often different from globalAnswer)
    return someValueForParent
}
```

**Classic problems:**
- **Max Depth (alt)** -- `return 1 + max(left, right)`
- **Diameter** -- global tracks `left + right`; return `1 + max(left, right)` to parent
- **Balanced Check** -- return height, but return -1 as sentinel for "unbalanced"
- **Subtree Sum / Max Path Sum** -- return single-branch sum; global tracks split path

**Complexity:** O(n) time, O(h) space.

**Watch out:**
- **The return value is NOT the answer.** The return value is what the parent needs. The actual answer is often a global/closure variable updated at each node. This is the #1 source of bugs.
- Base case must be correct: for height, nil returns 0; for max-path-sum, nil returns 0 (not math.MinInt, since you want to *ignore* empty branches).

---

### Pattern 3: Level-Order BFS

**Trigger:** The problem mentions "level," "depth," "left-to-right," "right side," "zigzag," or "average at each depth." Anything that processes nodes *layer by layer*.

**Go Template:**

```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return nil
    }
    var result [][]int
    queue := []*TreeNode{root}

    for len(queue) > 0 {
        size := len(queue)          // SNAPSHOT the level size
        level := make([]int, 0, size)

        for i := 0; i < size; i++ { // process exactly this level
            node := queue[0]
            queue = queue[1:]
            level = append(level, node.Val)

            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        result = append(result, level)
    }
    return result
}
```

**Classic problems:**
- **Level Order Traversal** -- collect each level into a subslice
- **Right Side View** -- take the last element of each level (or BFS right-to-left, take first)
- **Zigzag** -- reverse odd-indexed levels
- **Average of Levels** -- sum each level, divide by size

**Complexity:** O(n) time, O(w) space where w = max width of tree (up to n/2 at the last level of a complete tree).

**Watch out:**
- **Snapshot `len(queue)` before the inner loop.** The queue grows inside the loop; if you use `len(queue)` as the loop bound directly, you mix levels.
- Go has no built-in queue. Using a slice with `queue[0]` / `queue[1:]` is fine for interviews but is O(n) total for dequeue. Mention this if asked about optimization.
- For "right side view," you do NOT need a separate data structure -- just check `i == size-1` inside the inner loop.

---

### Pattern 4: Path Tracking

**Trigger:** The problem asks for actual paths (not just existence), root-to-leaf enumerations, or "all paths that sum to X."

**Go Template:**

```go
func pathTrack(node *TreeNode, target int, path []int, result *[][]int) {
    if node == nil {
        return
    }
    path = append(path, node.Val)

    // Leaf check: is this a complete root-to-leaf path?
    if node.Left == nil && node.Right == nil {
        if target == node.Val { // or whatever condition
            // COPY the path -- critical!
            tmp := make([]int, len(path))
            copy(tmp, path)
            *result = append(*result, tmp)
        }
        return // backtracking happens via slice mechanics
    }

    pathTrack(node.Left, target-node.Val, path, result)
    pathTrack(node.Right, target-node.Val, path, result)
    // No explicit "remove last" needed IF you pass path by value-header (slice trick).
    // But if using a shared slice with manual backtrack: path = path[:len(path)-1]
}
```

**Classic problems:**
- **Root-to-Leaf Paths** (collect all)
- **Path Sum II** (collect paths that sum to target)
- **Binary Tree Paths** (return as strings)
- **Path Sum III** (prefix sum variant -- any start/end node, uses hashmap)

**Complexity:** O(n * h) time in worst case (copying paths), O(h) space for recursion.

**Watch out:**
- **You MUST copy the slice before appending to results.** Go slices share underlying arrays. If you do `*result = append(*result, path)` without copying, all entries will point to the same (mutated) backing array.
- For Path Sum III (paths starting anywhere), the pattern shifts to prefix-sum with a hashmap -- it's a different beast. Know that variant exists.
- Backtracking with slices in Go: `append(path, val)` may or may not create a new backing array. The safest interview approach: always copy when saving, or explicitly backtrack with `path = path[:len(path)-1]`.

---

### Pattern 5: Tree Construction

**Trigger:** "Build a tree from two traversals" or "serialize / deserialize a tree."

**Go Template (preorder + inorder):**

```go
func buildTree(preorder []int, inorder []int) *TreeNode {
    if len(preorder) == 0 {
        return nil
    }
    rootVal := preorder[0]
    root := &TreeNode{Val: rootVal}

    // Find root in inorder to split left/right
    mid := 0
    for i, v := range inorder {
        if v == rootVal {
            mid = i
            break
        }
    }
    // inorder[:mid]  = left subtree's inorder
    // inorder[mid+1:] = right subtree's inorder
    // preorder[1:1+mid] = left subtree's preorder (same count as inorder left)
    // preorder[1+mid:]  = right subtree's preorder

    root.Left = buildTree(preorder[1:1+mid], inorder[:mid])
    root.Right = buildTree(preorder[1+mid:], inorder[mid+1:])
    return root
}
```

**Serialize/Deserialize (preorder with nil markers):**

```go
// Serialize: preorder, write "null" for nil nodes.
// Deserialize: consume tokens one by one from a queue/index.
func serialize(root *TreeNode) string {
    if root == nil {
        return "null"
    }
    return fmt.Sprintf("%d,%s,%s", root.Val, serialize(root.Left), serialize(root.Right))
}

// Deserialize uses a pointer to track position in the token list.
func deserialize(data string) *TreeNode {
    tokens := strings.Split(data, ",")
    idx := 0
    var build func() *TreeNode
    build = func() *TreeNode {
        if tokens[idx] == "null" {
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

**Complexity:** O(n) with index map optimization, O(n^2) naive for construction. Serialize/deserialize is O(n).

**Watch out:**
- The key insight for preorder+inorder: `preorder[0]` is always the root. The position of that root in `inorder` tells you the *size* of the left subtree, which lets you split both arrays.
- For postorder+inorder: root is `postorder[last]`, then same split logic.
- You CANNOT reconstruct a unique tree from preorder+postorder alone (unless BST).
- Serialize/deserialize: the `idx` variable must be shared across recursive calls. In Go, use a closure or a pointer to an int.

---

### Pattern 6: LCA (Lowest Common Ancestor)

**Trigger:** "Find the lowest common ancestor of two nodes." Also useful as a subroutine in distance-between-nodes problems.

**Go Template:**

```go
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    // Base cases
    if root == nil {
        return nil
    }
    if root == p || root == q {
        return root // found one of the targets
    }

    left := lowestCommonAncestor(root.Left, p, q)
    right := lowestCommonAncestor(root.Right, p, q)

    // Split point: p is in one subtree, q is in the other
    if left != nil && right != nil {
        return root
    }
    // Both in same subtree
    if left != nil {
        return left
    }
    return right
}
```

**Why it works:** Post-order traversal. Each call returns:
- `nil` if neither p nor q is in this subtree
- `p` or `q` if one of them is found
- The LCA if both are found in different subtrees

The node where `left != nil && right != nil` is the split point -- that's the LCA.

**BST variant:** If the tree is a BST, compare values:

```go
func lcaBST(root, p, q *TreeNode) *TreeNode {
    for root != nil {
        if p.Val < root.Val && q.Val < root.Val {
            root = root.Left
        } else if p.Val > root.Val && q.Val > root.Val {
            root = root.Right
        } else {
            return root // split point
        }
    }
    return nil
}
```

**Complexity:** O(n) for general tree, O(h) for BST. Space: O(h) recursive, O(1) iterative BST.

**Watch out:**
- **When one node is the ancestor of the other**, the function still works: the ancestor node is found first and returned up. The other node is in its subtree, but the ancestor node is already the LCA. Trace through an example to convince yourself.
- This assumes both p and q exist in the tree. If that's not guaranteed, you need a two-pass approach (first verify existence).
- For BST-LCA, the iterative approach is preferred -- O(1) space, cleaner.

---

## Decision Framework

When you read a tree problem in an interview, run through this checklist:

```
1. Does the problem mention "level," "depth layer," "left to right"?
   --> YES: Level-order BFS (Pattern 3)

2. Does the answer for a node depend on info from its PARENT?
   (running sum, valid range, depth counter)
   --> YES: Top-down DFS (Pattern 1)

3. Does the answer for a node depend on answers from its CHILDREN?
   (height, diameter, is-balanced, subtree sum)
   --> YES: Bottom-up DFS (Pattern 2)

4. Does the problem ask for actual paths or all paths matching a condition?
   --> YES: Path tracking with backtracking (Pattern 4)

5. Does the problem ask to build/reconstruct a tree?
   --> YES: Tree construction / serialize-deserialize (Pattern 5)

6. Does the problem ask for an ancestor or a meeting point?
   --> YES: LCA pattern (Pattern 6)
```

**Quick rules of thumb:**
- "Return a property of the whole tree" (height, diameter, balanced) --> **bottom-up DFS**
- "Check each node against a condition passed from parent" (range check, path sum) --> **top-down DFS**
- "Level by level" or "left to right at each depth" --> **BFS**
- "Find an ancestor" --> **LCA pattern**
- "Construct a tree from traversal" --> **recursive subdivision using traversal properties**

**Hybrid problems exist.** Diameter uses bottom-up to compute heights but a global variable to track the max left+right sum. Max Path Sum is the same shape. Recognize the "return one thing to parent, track a different thing globally" pattern.

---

## Common Interview Traps

### Trap 1: Diameter -- The Path Goes THROUGH a Node

The diameter is the longest path between *any* two nodes. This path does NOT have to pass through the root. At each node, the "through this node" path length is `leftHeight + rightHeight`. But you return `1 + max(leftHeight, rightHeight)` to the parent (a single branch).

```
      1
     / \
    2   3       Diameter = 4 (path: 5->2->1->3->4, but NOT rooted at 1
   /     \        if the tree were deeper on one side)
  5       4
```

**The fix:** Global variable tracks the max across all nodes. The return value only carries single-branch height upward.

### Trap 2: Balanced Tree -- Recursive, Not One-Level Check

A balanced tree requires *every* subtree to have left/right heights differ by at most 1. Checking only the root's children is wrong.

```
        1
       / \
      2   3       Root: |2-1| = 1, looks balanced.
     / \    \     But node 2's subtree: |2-0| = 2, NOT balanced.
    4   5    6
   /
  7
```

**The fix:** Use bottom-up with a -1 sentinel. If any subtree returns -1, propagate -1 upward immediately.

```go
func checkBalanced(node *TreeNode) int {
    if node == nil {
        return 0
    }
    left := checkBalanced(node.Left)
    if left == -1 { return -1 }
    right := checkBalanced(node.Right)
    if right == -1 { return -1 }
    if abs(left - right) > 1 { return -1 }
    return 1 + max(left, right)
}
```

### Trap 3: LCA -- One Node Is the Ancestor of the Other

```
        3
       / \
      5   1
     / \
    6   2
```

LCA(5, 6) = 5. Node 5 IS the ancestor. The algorithm works because when we reach node 5, we return it immediately (it matches `root == p`). We never search deeper into its subtree for q. The left call at node 3 returns 5, the right call returns nil, so we return 5. Correct.

### Trap 4: BFS -- Forgetting to Snapshot Queue Length

```go
// WRONG: queue grows during iteration, mixing levels
for i := 0; i < len(queue); i++ { ... }

// RIGHT: snapshot the size before processing the level
size := len(queue)
for i := 0; i < size; i++ { ... }
```

### Trap 5: Confusing "Height" and "Depth"

| Term   | Direction  | Root value | Leaf value | Pattern    |
|--------|-----------|------------|------------|------------|
| Height | Bottom-up | max        | 0          | Bottom-up  |
| Depth  | Top-down  | 0          | max        | Top-down   |

"Max depth of a tree" and "height of a tree" are the same number, but they're computed in opposite directions. Interviewers sometimes test whether you know the distinction.

### Trap 6: Bottom-Up Return vs Global Answer

```go
// In diameter: the function returns HEIGHT (for the parent).
// The DIAMETER is tracked in a separate variable.

var diameter int

func height(node *TreeNode) int {
    // ...
    diameter = max(diameter, left + right) // global: the "through" path
    return 1 + max(left, right)            // return: single branch to parent
}
```

If you return `left + right` to the parent, you're saying "the parent can use both branches," which is wrong -- a path can't fork.

---

## Thought Process Walkthrough

### Problem 1: Binary Tree Maximum Depth (LC 104) -- The Warmup

**Interviewer says:** "Given the root of a binary tree, return its maximum depth. Maximum depth is the number of nodes along the longest path from root to a leaf."

**Step 1 -- Classify (10 seconds)**

"Maximum depth is a property of the whole tree that depends on subtree depths. That's bottom-up DFS."

(Could also do top-down by passing depth down and tracking max. Both work. Bottom-up is cleaner here.)

**Step 2 -- Recurrence (20 seconds)**

"If I know the max depth of left and right subtrees, the depth of this node is `1 + max(left, right)`. Base case: nil node has depth 0."

**Step 3 -- Code (2 minutes)**

```go
func maxDepth(root *TreeNode) int {
    if root == nil {
        return 0
    }
    left := maxDepth(root.Left)
    right := maxDepth(root.Right)
    return 1 + max(left, right)
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

**Step 4 -- Verify (1 minute)**

Trace through:
```
    3
   / \
  9  20
     / \
    15   7
```
- `maxDepth(9)` = 1 + max(0,0) = 1
- `maxDepth(15)` = 1, `maxDepth(7)` = 1
- `maxDepth(20)` = 1 + max(1,1) = 2
- `maxDepth(3)` = 1 + max(1,2) = 3. Correct.

**Step 5 -- Complexity**
- Time: O(n), visit every node once.
- Space: O(h), recursion stack. Worst case O(n) for skewed tree, O(log n) for balanced.

**Total time: ~4 minutes.** This is the warmup. It proves you can handle recursion cleanly.

---

### Problem 2: Diameter of Binary Tree (LC 543) -- The Real Test

**Interviewer says:** "Given the root of a binary tree, return the length of the diameter. The diameter is the longest path between any two nodes (measured in edges)."

**Step 1 -- Classify (15 seconds)**

"The diameter passes through some node. At that node, it's left height + right height. I need to check every node. This is bottom-up DFS with a global tracker."

Key insight to state aloud: "The longest path doesn't have to go through the root."

**Step 2 -- Identify the dual role (30 seconds)**

"My recursive function returns the HEIGHT of the subtree (what the parent needs). But at each node, I also compute `left + right` which is the diameter THROUGH this node. I track the max diameter globally."

Say to interviewer: "The function has two jobs: return height upward, and update diameter sideways. I'll use a closure variable for the diameter."

**Step 3 -- Code (3 minutes)**

```go
func diameterOfBinaryTree(root *TreeNode) int {
    diameter := 0

    var height func(node *TreeNode) int
    height = func(node *TreeNode) int {
        if node == nil {
            return 0
        }
        left := height(node.Left)
        right := height(node.Right)

        // The path through this node has length left + right (in edges)
        if left + right > diameter {
            diameter = left + right
        }

        // Return this node's height (in edges) to parent
        return 1 + max(left, right)
    }

    height(root)
    return diameter
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

**Step 4 -- Trace through a tricky case (2 minutes)**

```
        1
       / \
      2   3
     / \
    4   5
```

- `height(4)` = 1 + max(0,0) = 1. diameter = max(0, 0+0) = 0.
- `height(5)` = 1. diameter = 0.
- `height(2)` = 1 + max(1,1) = 2. diameter = max(0, 1+1) = 2.
- `height(3)` = 1. diameter = max(2, 0+0) = 2.
- `height(1)` = 1 + max(2,1) = 3. diameter = max(2, 2+1) = 3.

Answer: 3. Path is 4->2->1->3 (or 5->2->1->3). Correct.

**Step 5 -- Edge cases to mention**
- Single node: diameter = 0.
- Skewed tree (linked list): diameter = n-1.
- Diameter doesn't pass through root: already handled since we check every node.

**Step 6 -- Complexity**
- Time: O(n).
- Space: O(h).

**Total time: ~7 minutes.** This demonstrates you understand the dual-return pattern that separates strong candidates from average ones.

---

## Time Targets

| Problem | Target | Notes |
|---------|--------|-------|
| Max Depth (LC 104) | < 4 min | Pure bottom-up. Must be instant. |
| Same Tree (LC 100) | < 4 min | Top-down comparison. |
| Path Sum (LC 112) | < 5 min | Top-down with leaf check. |
| Level Order (LC 102) | < 6 min | BFS template. Know it cold. |
| Diameter (LC 543) | < 8 min | Bottom-up with global. Must explain dual role. |
| Balanced Tree (LC 110) | < 8 min | Bottom-up with -1 sentinel. |
| LCA (LC 236) | < 8 min | The split-point recursion. |
| Right Side View (LC 199) | < 8 min | BFS, take last of each level. |
| Build from Pre+In (LC 105) | < 12 min | Subdivision + index math. |
| Serialize/Deserialize (LC 297) | < 15 min | Preorder with nil markers. |

---

## Quick Drill (5 exercises)

Do these on paper or whiteboard first, then verify in code.

### Drill 1: Pattern Classification (2 min)

For each problem, write ONLY the pattern name (don't solve it):
1. "Return true if two trees are mirrors of each other" --> ?
2. "Return the sum of all left leaves" --> ?
3. "Return the values visible from the right side" --> ?
4. "Find the LCA of two nodes in a BST" --> ?
5. "Return all root-to-leaf paths as strings" --> ?

<details>
<summary>Answers</summary>

1. Top-down DFS (compare mirrored pairs: left.left with right.right)
2. Bottom-up DFS (or top-down passing "isLeft" flag -- either works)
3. Level-order BFS (right side view)
4. LCA pattern (BST variant -- iterative, compare values)
5. Path tracking (collect paths, format as strings)

</details>

### Drill 2: Spot the Bug (2 min)

```go
func isBalanced(root *TreeNode) bool {
    if root == nil {
        return true
    }
    left := height(root.Left)
    right := height(root.Right)
    return abs(left-right) <= 1
}
```

What's wrong?

<details>
<summary>Answer</summary>

It only checks balance at the root. A subtree deeper down could be unbalanced and this would miss it. Fix: either (a) recurse with `isBalanced(root.Left) && isBalanced(root.Right)` (O(n log n)), or (b) use the -1 sentinel bottom-up approach (O(n)).

</details>

### Drill 3: Write Diameter from Memory (5 min)

Close this guide. Write `diameterOfBinaryTree` from scratch. Check against the template above. Did you:
- Use a closure variable for the global answer?
- Return height (not diameter) to the parent?
- Handle nil correctly?

### Drill 4: BFS Variant (5 min)

Write "Zigzag Level Order Traversal" (LC 103). Hint: it's the standard BFS template with one twist -- reverse the level slice for odd-indexed levels.

<details>
<summary>Skeleton</summary>

```go
func zigzagLevelOrder(root *TreeNode) [][]int {
    if root == nil {
        return nil
    }
    var result [][]int
    queue := []*TreeNode{root}
    leftToRight := true

    for len(queue) > 0 {
        size := len(queue)
        level := make([]int, size)
        for i := 0; i < size; i++ {
            node := queue[0]
            queue = queue[1:]
            // Place value based on direction
            if leftToRight {
                level[i] = node.Val
            } else {
                level[size-1-i] = node.Val
            }
            if node.Left != nil { queue = append(queue, node.Left) }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
        leftToRight = !leftToRight
    }
    return result
}
```

</details>

### Drill 5: LCA Trace (3 min)

Given this tree, trace the LCA algorithm for nodes 6 and 4:
```
        3
       / \
      5   1
     / \ / \
    6  2 0  8
```

Write what each recursive call returns at each node.

<details>
<summary>Trace</summary>

- `LCA(6, ...)`: root==6==p, return 6
- `LCA(2, ...)`: neither 2==p nor 2==q; left=nil, right=nil -> return nil
- `LCA(5, ...)`: left=6 (non-nil), right=nil -> return 6 (LCA not found yet, propagate)
- `LCA(0, ...)`: return nil
- `LCA(8, ...)`: root==8? No. It's node 8, not 4. return nil. (Wait -- where is 4? It's not in this tree as drawn. Let's say we meant nodes 6 and 8.)
- `LCA(1, ...)`: left=nil, right=8 -> return 8
- `LCA(3, ...)`: left=6 (non-nil), right=8 (non-nil) -> return 3. **LCA(6,8) = 3.**

</details>

---

## Self-Assessment (5 questions)

Rate yourself honestly after today's session.

**1. Can you write the BFS level-order template from memory in under 2 minutes?**

If no: practice it 3 more times. It should be muscle memory. The snapshot-queue-length pattern appears in every BFS variant.

**2. Given a new tree problem, can you classify it into one of the 6 patterns within 30 seconds?**

If no: re-read the Decision Framework. Practice by scrolling through LeetCode's "Tree" tag and classifying 10 problems without solving them.

**3. Can you explain why diameter returns height to the parent but tracks diameter globally?**

If no: this is the core "dual role" insight of bottom-up DFS. Trace through the example again. The function serves two purposes -- and confusing them is the most common mistake.

**4. Do you handle Go slice aliasing correctly in path-tracking problems?**

If no: write a small test that demonstrates the bug. Create a path `[]int{1,2,3}`, append it to results, then modify the original. See what happens to results. This will burn the lesson in.

**5. Can you write LCA and explain why it handles the case where one node is an ancestor of the other?**

If no: trace through the case where p=5 and q=6 in the tree above. The key insight is that when we find p, we return it immediately without searching deeper -- we don't need to, because if q is below p, then p is already the LCA.

---

## Summary Cheat Sheet

```
PATTERN             | RETURN TO PARENT     | GLOBAL/ANSWER          | BASE (nil)
--------------------|----------------------|------------------------|----------
Top-down DFS        | (void -- no return)  | updated via params     | return
Bottom-up height    | 1 + max(L, R)        | n/a                    | 0
Bottom-up diameter  | 1 + max(L, R)        | max(L + R) across all  | 0
Bottom-up balanced  | height or -1         | final != -1            | 0
BFS level-order     | n/a                  | collected per level    | nil check
Path tracking       | (void -- backtrack)  | collected at leaves    | return
LCA                 | node or nil          | the split point        | nil
```

> **Key takeaway:** In tree interviews, the hard part is never the traversal -- it's
> recognizing which information flows UP vs DOWN, and keeping the return value
> separate from the answer. Master that distinction and every tree problem becomes
> a variation of the same 6 templates.
