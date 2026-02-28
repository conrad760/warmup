# Day 4 — Binary Trees & Traversals: Deep Dive

**Date:** Day 4 of 21-day DSA refresher
**Time:** 12:00 PM - 2:00 PM (2 hours)
**Language:** Go
**Prerequisite:** You built stacks and queues yesterday — you'll use both today.

---

## 1. Curated Learning Resources

| # | Title | URL | Category | Why It's Here |
|---|-------|-----|----------|---------------|
| 1 | **Binary Tree Visualizer** | https://visualgo.net/en/bst | Interactive | Step through traversals visually. Toggle between pre/in/post/level-order and watch the animation. |
| 2 | **Tree Traversals (Inorder, Preorder, Postorder) — Abdul Bari** | https://www.youtube.com/watch?v=BHB0B1jFKQc | Video | Clear 15-minute walkthrough of all three DFS traversals with call-stack tracing. |
| 3 | **Iterative Inorder Traversal — NeetCode** | https://www.youtube.com/watch?v=g_S5WuasWUE | Video | Explains the "push all lefts" pattern with stack visualization. |
| 4 | **Morris Traversal Explained** | https://www.educative.io/answers/what-is-morris-traversal | Article | Threaded binary tree concept with diagrams. The clearest explanation of how temporary links work. |
| 5 | **Iterative Postorder Traversal — Two Approaches** | https://leetcode.com/problems/binary-tree-postorder-traversal/editorial/ | Article | LeetCode editorial comparing two-stack and one-stack (last-visited) approaches side by side. |
| 6 | **Binary Tree Bootcamp (Back to Back SWE)** | https://www.youtube.com/watch?v=BHB0B1jFKQc | Video | Covers tree terminology, recursive thinking, and complexity analysis. |
| 7 | **Go Data Structures: Binary Tree** | https://golangbyexample.com/binary-tree-in-go/ | Go-specific | Idiomatic Go implementation patterns for tree nodes, including nil handling. |
| 8 | **LeetCode Tree Tag — Easy** | https://leetcode.com/tag/tree/ | Practice | Filter by Easy for today's utility functions (MaxDepth, Invert, Symmetric). |

---

## 2. Detailed 2-Hour Session Plan

### Review Block (12:00 - 12:20)

| Time | Activity |
|------|----------|
| 12:00 - 12:07 | Read through the terminology section below. Quiz yourself: what's the difference between *complete* and *full*? Between *depth* and *height*? Draw a tree that is full but not complete. |
| 12:07 - 12:14 | Study the traversal order table. On paper, draw a 7-node tree and write the visit order for all 4 traversals. Don't code yet — just trace. |
| 12:14 - 12:20 | Read the iterative in-order and post-order concept sections below. Understand the *why* before touching the keyboard. |

### Implement Block (12:20 - 1:20)

| Time | Activity | Notes |
|------|----------|-------|
| 12:20 - 12:25 | Define `TreeNode` struct and a helper to build a sample tree | Manual construction: `root := &TreeNode{Val: 1, Left: ...}` |
| 12:25 - 12:35 | **Recursive traversals**: InOrder, PreOrder, PostOrder | These should be trivial. ~3 min each. Use a helper with `*[]int` result. |
| 12:35 - 12:40 | **Level-order traversal** (BFS) | Use a `[]*TreeNode` slice as queue. Return `[][]int`. |
| 12:40 - 12:50 | **Iterative in-order** | The "push all lefts" pattern. This is the key one. |
| 12:50 - 1:05 | **Iterative post-order** | Implement both: (1) two-stack trick, (2) one-stack with `lastVisited`. Compare. |
| 1:05 - 1:15 | **Iterative pre-order** | Easiest iterative — push right then left. Do this one last as a cooldown. |
| 1:15 - 1:20 | Run all traversals on the same tree, verify outputs match recursive versions | Print results side by side. |

### Solidify Block (1:20 - 1:50)

| Time | Activity |
|------|----------|
| 1:20 - 1:30 | Implement `MaxDepth(root)` — recursive one-liner, then iterative BFS version counting levels. |
| 1:30 - 1:40 | Implement `InvertTree(root)` — swap left and right recursively. |
| 1:40 - 1:50 | Implement `IsSymmetric(root)` — write a helper `isMirror(left, right)`. Test with symmetric and asymmetric trees. Edge cases: nil root, single node, skewed tree. |

### Recap Block (1:50 - 2:00)

| Time | Activity |
|------|----------|
| 1:50 - 1:55 | From memory, write the complexity of each traversal (time and space). Write the iterative in-order algorithm in pseudocode without looking. |
| 1:55 - 2:00 | Write down one gotcha you hit today and one thing that clicked. |

---

## 3. Core Concepts Deep Dive

### Tree Terminology Refresher

```
             1              <- depth 0
           /   \
          2     3           <- depth 1
         / \     \
        4   5     6         <- depth 2
       /
      7                     <- depth 3
```

| Term | Definition | This Tree |
|------|-----------|-----------|
| **Depth** of a node | Number of edges from root to that node | Node 7 has depth 3 |
| **Height** of a tree | Maximum depth of any node (= depth of deepest leaf) | Height = 3 |
| **Height** of a node | Number of edges on longest path from that node to a leaf | Node 2 has height 2 |
| **Complete** | Every level full except possibly the last, filled left-to-right | This tree is NOT complete (level 2 missing node after 5) |
| **Full** | Every node has 0 or 2 children | NOT full (node 3 has 1 child, node 4 has 1 child) |
| **Perfect** | Full + all leaves at the same depth | NOT perfect |
| **Balanced** | Left/right subtree heights differ by at most 1 (at every node) | NOT balanced (left subtree height 3, right subtree height 1) |

Key relationship: A **perfect** tree is always **complete** and **full**. A **complete** tree is not necessarily full. A **full** tree is not necessarily complete.

Node count facts:
- Perfect binary tree of height h: `2^(h+1) - 1` nodes
- Complete binary tree: between `2^h` and `2^(h+1) - 1` nodes at height h

---

### Why Recursive Traversals Map to the Call Stack

Every recursive call pushes a frame onto the call stack. When the function returns, the frame pops. This is exactly a stack — and it's why iterative traversals use an explicit stack to mimic recursion.

```go
func inorder(node *TreeNode, result *[]int) {
    if node == nil { return }
    inorder(node.Left, result)   // push left frame
    *result = append(*result, node.Val) // visit after left returns
    inorder(node.Right, result)  // push right frame
}
```

The call stack for the sample tree during in-order:
```
Call inorder(1)
  Call inorder(2)
    Call inorder(4)
      Call inorder(7)
        Call inorder(nil) → return
        Visit 7
        Call inorder(nil) → return
      Visit 4
      Call inorder(nil) → return
    Visit 2
    Call inorder(5)
      Call inorder(nil) → return
      Visit 5
      Call inorder(nil) → return
  Visit 1
  Call inorder(3)
    Call inorder(nil) → return
    Visit 3
    Call inorder(6)
      Call inorder(nil) → return
      Visit 6
      Call inorder(nil) → return
```

The recursive structure naturally handles "go deep first, then come back." The explicit stack in iterative versions does the same thing — you just manage the push/pop yourself.

---

### Iterative In-Order: The "Push All Lefts" Pattern

**The algorithm:**
```
1. Start at root.
2. Push current node and move to left child. Repeat until nil.
3. Pop from stack, visit the node.
4. Move to the right child. Go to step 2.
5. Stop when stack is empty AND current is nil.
```

**Why it works:**

In-order means: left subtree → node → right subtree.

The "push all lefts" phase walks to the leftmost node while saving every ancestor on the stack. When you hit nil, the top of the stack is the leftmost unvisited node — exactly what in-order says to visit first.

After visiting a node, you move to its right child. If the right child exists, the same "push all lefts" phase processes the right subtree's leftmost path. If right is nil, the next pop gives the parent, which is the next in-order node.

**The invariant**: At any point, the stack holds the chain of unvisited ancestors whose left subtrees have been fully processed. The next node to visit is always the stack top.

Think of it this way: the stack stores the "return addresses" — the nodes you still owe a visit to after you finish exploring their left subtrees.

```go
func InorderIterative(root *TreeNode) []int {
    var result []int
    var stack []*TreeNode
    curr := root

    for curr != nil || len(stack) > 0 {
        // Phase 1: Push all lefts
        for curr != nil {
            stack = append(stack, curr)
            curr = curr.Left
        }
        // Phase 2: Pop and visit
        curr = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        result = append(result, curr.Val)
        // Phase 3: Move to right subtree
        curr = curr.Right
    }
    return result
}
```

---

### Iterative Post-Order: Two Approaches Compared

Post-order (left → right → root) is the hardest to do iteratively because you visit a node *after* both its children. When you pop a node from the stack, you don't know if you've already processed its right child or not.

#### Approach 1: Two-Stack Trick

**Idea:** Pre-order is root → left → right. If you change it to root → right → left (push left before right), and then *reverse* the result, you get left → right → root — which is post-order.

```go
func PostorderTwoStack(root *TreeNode) []int {
    if root == nil { return nil }
    var result []int
    stack := []*TreeNode{root}

    for len(stack) > 0 {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        result = append(result, node.Val)  // "visit" in reverse post-order
        if node.Left != nil {              // push left FIRST
            stack = append(stack, node.Left)
        }
        if node.Right != nil {             // so right is processed FIRST
            stack = append(stack, node.Right)
        }
    }
    // Reverse result
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    return result
}
```

**Pros:** Simple, easy to remember.
**Cons:** Uses O(n) extra space for the result reversal (or a second stack). Not "true" post-order processing — you're visiting in a different order and reversing.

#### Approach 2: One Stack with `lastVisited`

**Idea:** Track the last visited node. When you peek at the stack top:
- If it has a right child that you haven't visited yet → move to the right child.
- Otherwise → pop and visit it.

```go
func PostorderLastVisited(root *TreeNode) []int {
    var result []int
    var stack []*TreeNode
    var lastVisited *TreeNode
    curr := root

    for curr != nil || len(stack) > 0 {
        // Push all lefts
        for curr != nil {
            stack = append(stack, curr)
            curr = curr.Left
        }
        // Peek
        top := stack[len(stack)-1]
        // If right child exists and hasn't been visited
        if top.Right != nil && top.Right != lastVisited {
            curr = top.Right
        } else {
            // Visit this node
            result = append(result, top.Val)
            lastVisited = top
            stack = stack[:len(stack)-1]
        }
    }
    return result
}
```

**Pros:** True post-order processing (visits in the correct order as you go). Useful when you need to *process* nodes in post-order (e.g., evaluating expression trees, computing subtree sizes) rather than just collecting values.
**Cons:** Harder to remember. The `lastVisited` check is subtle.

#### When to Use Which

| Approach | Use When |
|----------|----------|
| Two-stack/reverse | You just need the post-order sequence (most interview problems) |
| `lastVisited` | You need to process nodes in actual post-order (compute subtree aggregates, free memory, etc.) |

---

### Morris Traversal: O(1) Space In-Order

Standard iterative traversal uses O(h) stack space. Morris traversal achieves O(1) space by temporarily modifying the tree.

**Core idea:** Use the nil right pointers of in-order predecessor nodes as "threads" back to the current node — creating temporary links that let you return to a node after processing its left subtree, without a stack.

**Algorithm:**
```
curr = root
while curr != nil:
    if curr.Left == nil:
        visit(curr)
        curr = curr.Right        // follow the right pointer (might be a thread)
    else:
        // Find in-order predecessor (rightmost node in left subtree)
        pred = curr.Left
        while pred.Right != nil AND pred.Right != curr:
            pred = pred.Right

        if pred.Right == nil:
            // First visit: create thread
            pred.Right = curr     // temporary link back to curr
            curr = curr.Left      // go left
        else:
            // Second visit: remove thread, visit curr
            pred.Right = nil      // restore tree
            visit(curr)
            curr = curr.Right
```

**Why it works:** The predecessor's right pointer is normally nil (it's a leaf or has no right child in the left subtree). Morris "borrows" this nil pointer to store a return link. When you follow the thread back, you know you've already processed the left subtree, so you visit the node and move right.

**Complexity:** O(n) time (each edge is traversed at most 3 times), O(1) space.

**Trade-off:** Temporarily mutates the tree. Not safe for concurrent access. The tree is restored to its original form by the end.

Implementation is optional for Day 4, but understanding the concept demonstrates mastery of tree structure. It occasionally appears in interviews as a follow-up: "Can you do in-order traversal in O(1) space?"

---

### Relationship to Graph Traversals

| Tree Traversal | Graph Equivalent | Data Structure | Key Property |
|---------------|------------------|----------------|--------------|
| Pre-order | DFS (pre-visit) | Stack | Visit node before children — useful for top-down problems |
| In-order | (BST-specific) | Stack | Produces sorted output in BSTs |
| Post-order | DFS (post-visit) | Stack | Visit node after children — useful for bottom-up problems |
| Level-order | BFS | Queue | Visit nodes by distance from root — shortest path |

DFS on a tree is just pre/in/post-order. There's no "visited" set needed because trees are acyclic. DFS on a general graph needs visited tracking to avoid infinite loops.

BFS on a tree is level-order traversal. Same queue-based approach as graph BFS. The "level-by-level" structure comes from processing all nodes at distance d before any at distance d+1.

---

## 4. Implementation Checklist

### TreeNode Definition

```go
type TreeNode struct {
    Val         int
    Left, Right *TreeNode
}
```

### Recursive Traversals

```go
func InorderRecursive(root *TreeNode) []int
func PreorderRecursive(root *TreeNode) []int
func PostorderRecursive(root *TreeNode) []int
func LevelOrder(root *TreeNode) [][]int
```

### Iterative Traversals

```go
func InorderIterative(root *TreeNode) []int    // push-all-lefts pattern
func PreorderIterative(root *TreeNode) []int   // stack: push right, then left
func PostorderIterative(root *TreeNode) []int  // two-stack OR lastVisited
```

### Utility Functions

```go
func MaxDepth(root *TreeNode) int       // max(left depth, right depth) + 1
func InvertTree(root *TreeNode) *TreeNode  // swap left and right recursively
func IsSymmetric(root *TreeNode) bool   // helper: isMirror(left, right *TreeNode) bool
```

### Test Cases & Edge Cases

| Test Case | Expected Behavior |
|-----------|------------------|
| `nil` root | Return empty result / 0 / true / nil |
| Single node | Return `[val]` / depth 0 / unchanged / true |
| Left-skewed tree (linked list) | Traversals degenerate to linear scan. Stack grows to O(n). |
| Right-skewed tree | Same as above. |
| Perfect tree (3 levels) | All traversals produce distinct orderings. Verify each. |
| Tree with duplicate values | Traversals still work — values repeat in output. |

### Suggested Test Tree

```go
//         1
//        / \
//       2   3
//      / \   \
//     4   5   6

root := &TreeNode{Val: 1,
    Left: &TreeNode{Val: 2,
        Left:  &TreeNode{Val: 4},
        Right: &TreeNode{Val: 5},
    },
    Right: &TreeNode{Val: 3,
        Right: &TreeNode{Val: 6},
    },
}
```

Expected outputs:
- In-order: `[4, 2, 5, 1, 3, 6]`
- Pre-order: `[1, 2, 4, 5, 3, 6]`
- Post-order: `[4, 5, 2, 6, 3, 1]`
- Level-order: `[[1], [2, 3], [4, 5, 6]]`

---

## 5. Traversal Pattern Recognition

### In-Order (Left → Root → Right)

**When it's useful:**
- **BST sorted output** — In-order on a BST visits nodes in ascending order. This is the fundamental BST property.
- **Kth smallest element** — In-order traversal, stop at the kth visit.
- **Validate BST** — In-order should produce strictly increasing sequence.
- **Convert BST to sorted doubly linked list** — In-order traversal while linking nodes.

**Signal in problem:** "sorted order", "kth smallest", "validate BST", "BST to list".

### Pre-Order (Root → Left → Right)

**When it's useful:**
- **Serialization/deserialization** — Pre-order + null markers uniquely defines a tree. Deserialize by reading root, then recursively building left and right.
- **Copying/cloning a tree** — Visit root first, create copy, then recurse on children.
- **Building expression trees** — Prefix notation maps directly to pre-order.
- **Constructing tree from preorder + inorder** — Pre-order's first element is always the root.

**Signal in problem:** "serialize", "clone/copy tree", "construct tree from traversal", "flatten to linked list".

### Post-Order (Left → Right → Root)

**When it's useful:**
- **Deletion / freeing nodes** — Delete children before parent (bottom-up).
- **Computing subtree aggregates** — Height, size, sum, diameter. You need children's answers before computing the parent's.
- **Evaluating expression trees** — Evaluate operands before applying operator.
- **Lowest Common Ancestor** — Recurse into subtrees first, then decide at the current node.

**Signal in problem:** "height", "diameter", "subtree sum", "evaluate expression", "delete tree", "LCA".

### Level-Order (BFS)

**When it's useful:**
- **Level-by-level processing** — Average of each level, largest value in each row.
- **Zigzag traversal** — BFS with alternating left-right and right-left per level.
- **Right side view** — Last node of each BFS level.
- **Minimum depth** — BFS finds the first leaf (shortest path from root).
- **Connect next pointers** — Process all nodes at the same level to wire them together.

**Signal in problem:** "level", "row", "depth-wise", "minimum depth", "right side view", "zigzag".

### Quick Decision Table

| I need to... | Use |
|-------------|-----|
| Get sorted values from BST | In-order |
| Serialize a tree | Pre-order |
| Compute heights/sizes bottom-up | Post-order |
| Find minimum depth or process by level | Level-order |
| Find kth smallest in BST | In-order (stop at k) |
| Evaluate an expression tree | Post-order |
| Clone a tree | Pre-order |
| Get right-side view | Level-order |

---

## 6. Visual Diagrams

### Sample Tree

```
              1
            /   \
           2     3
          / \     \
         4   5     6
```

### Traversal Orders on This Tree

```
IN-ORDER (left, root, right):

    4 → 2 → 5 → 1 → 3 → 6

    Visit leftmost first, work back up and right.

PRE-ORDER (root, left, right):

    1 → 2 → 4 → 5 → 3 → 6

    Visit root first, then dive left, then right.

POST-ORDER (left, right, root):

    4 → 5 → 2 → 6 → 3 → 1

    Visit leaves first, root last.

LEVEL-ORDER (BFS):

    Level 0: [1]
    Level 1: [2, 3]
    Level 2: [4, 5, 6]

    Left to right at each depth.
```

### Iterative In-Order: Stack State Step-by-Step

Tree:
```
        1
       / \
      2   3
     / \
    4   5
```

```
Step | Action            | Stack (bottom→top) | curr  | Output
-----+-------------------+--------------------+-------+--------
  1  | Push 1, go left   | [1]                | 2     |
  2  | Push 2, go left   | [1, 2]             | 4     |
  3  | Push 4, go left   | [1, 2, 4]          | nil   |
  4  | Pop 4, visit      | [1, 2]             | nil*  | [4]
  5  | curr=4.Right=nil  | [1, 2]             | nil   | [4]
  6  | Pop 2, visit      | [1]                | nil*  | [4, 2]
  7  | curr=2.Right=5    | [1]                | 5     | [4, 2]
  8  | Push 5, go left   | [1, 5]             | nil   | [4, 2]
  9  | Pop 5, visit      | [1]                | nil*  | [4, 2, 5]
 10  | curr=5.Right=nil  | [1]                | nil   | [4, 2, 5]
 11  | Pop 1, visit      | []                 | nil*  | [4, 2, 5, 1]
 12  | curr=1.Right=3    | []                 | 3     | [4, 2, 5, 1]
 13  | Push 3, go left   | [3]                | nil   | [4, 2, 5, 1]
 14  | Pop 3, visit      | []                 | nil*  | [4, 2, 5, 1, 3]
 15  | curr=3.Right=nil  | []                 | nil   | [4, 2, 5, 1, 3]
 16  | Stack empty, done  |                    |       | [4, 2, 5, 1, 3]
```

*After popping, `curr` is set to the popped node temporarily, then reassigned to `node.Right`.

**Key observation:** The stack never holds more than `h` nodes (tree height). For this balanced tree of height 2, the max stack size is 3 (step 3). For a skewed tree of n nodes, the stack grows to n.

### Iterative Post-Order (lastVisited): Stack State

```
Step | Action              | Stack      | lastVisited | Output
-----+---------------------+------------+-------------+--------
  1  | Push 1→2→4          | [1,2,4]    | nil         |
  2  | Peek 4, no right    | [1,2,4]    | nil         |
  3  | Pop 4, visit        | [1,2]      | 4           | [4]
  4  | Peek 2, right=5≠4   | [1,2]      | 4           | [4]
  5  | Push 5 (curr=5)     | [1,2,5]    | 4           | [4]
  6  | Peek 5, no right    | [1,2,5]    | 4           |
  7  | Pop 5, visit        | [1,2]      | 5           | [4,5]
  8  | Peek 2, right=5==5  | [1,2]      | 5           |
  9  | Pop 2, visit        | [1]        | 2           | [4,5,2]
 10  | Peek 1, right=3≠2   | [1]        | 2           |
 11  | Push 3 (curr=3)     | [1,3]      | 2           |
 12  | Peek 3, no right*   | [1,3]      | 2           |
 13  | Pop 3, visit        | [1]        | 3           | [4,5,2,3]
 14  | Peek 1, right=3==3  | [1]        | 3           |
 15  | Pop 1, visit        | []         | 1           | [4,5,2,3,1]
```

*Node 3 has no left child, so we don't push anything extra. In the full tree (with node 6), step 12 would push to 3's right subtree first.

---

## 7. Self-Assessment

Answer these without looking at the material above. If you can't, that's your signal for tomorrow.

### Question 1
**Why is iterative post-order harder than iterative pre-order?**

<details>
<summary>Answer</summary>

Pre-order visits the root *before* its children — when you pop a node from the stack, you visit it immediately and push its children. You never need to return to it.

Post-order visits the root *after* both children. When you pop a node, you don't know if its right subtree has been processed yet. You might need to push it back. This ambiguity requires either a `lastVisited` tracker or the two-stack reversal trick.

In short: pre-order's "visit on pop" is simple. Post-order's "visit on pop only if right subtree is done" requires extra bookkeeping.
</details>

### Question 2
**What's the space complexity difference between recursive and iterative traversals? Is iterative always better?**

<details>
<summary>Answer</summary>

Both are O(h) where h is the tree height — recursive uses the call stack, iterative uses an explicit stack. The space usage is the same.

Iterative is NOT always better. For balanced trees, both use O(log n). For skewed trees, both use O(n). Iterative avoids potential stack overflow for very deep trees (Go's goroutine stack starts small but grows, so this is less of an issue in Go than in languages with fixed stack sizes). The only way to beat O(h) is Morris traversal at O(1) space.

Level-order (BFS) is different: it uses O(w) space where w is the max width of the tree, which can be O(n) for a complete tree's bottom level.
</details>

### Question 3
**You're processing a tree where each node represents a file or directory, and you need to compute the total size of each directory (including all subdirectories). Which traversal order do you use and why?**

<details>
<summary>Answer</summary>

Post-order. You need to know the sizes of all children (subdirectories and files) before you can compute the parent directory's total size. Post-order processes children before the parent, so by the time you visit a directory node, all its descendants have already been computed.

This is the classic "bottom-up aggregation" pattern that post-order is designed for.
</details>

### Question 4
**In the iterative in-order traversal, what does the stack represent at any given moment? Why do we "push all lefts"?**

<details>
<summary>Answer</summary>

The stack holds the chain of ancestor nodes whose left subtrees are currently being explored (or have been explored) but the nodes themselves haven't been visited yet.

We "push all lefts" because in-order says: visit the leftmost node first. Pushing all left children walks us to the leftmost node while bookmarking every ancestor along the way. When we pop, we visit that node, then move to its right subtree. If the right subtree is nil, the next pop gives the parent — which is the next node in in-order sequence because its left subtree (which we just finished) has been fully processed.
</details>

### Question 5
**Given a binary tree with n nodes, what is the maximum and minimum possible height? How does height affect traversal space complexity?**

<details>
<summary>Answer</summary>

- **Minimum height:** `floor(log2(n))` — achieved by a complete/perfect binary tree. A perfectly balanced tree packs the most nodes into the fewest levels.
- **Maximum height:** `n - 1` — achieved by a completely skewed tree (every node has only one child, essentially a linked list).

Traversal space complexity is O(h):
- Balanced tree: O(log n) space — very efficient.
- Skewed tree: O(n) space — the stack holds every node.

This is why balanced BSTs (AVL, Red-Black) matter: they guarantee O(log n) height, which keeps all operations efficient.
</details>
