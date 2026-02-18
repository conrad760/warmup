package main

// curatedBankExtended contains the remaining Blind 75 problems.
var curatedBankExtended = []CuratedQuestion{
	// Arrays & Hashing
	{
		Slug: "contains-duplicate",
		Options: []Option{
			{Text: "Use a hash set — add each element and check if it already exists — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Sort the array then check adjacent elements for duplicates — O(n log n) time, O(1) space", Rating: Plausible},
			{Text: "Compare every pair of elements with nested loops — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Check if the length equals the number of unique elements using a bit vector of fixed size — fails when values exceed bit vector range — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Hash Set
// Time: O(n) | Space: O(n)
func containsDuplicate(nums []int) bool {
    seen := make(map[int]bool)
    for _, n := range nums {
        if seen[n] {
            return true
        }
        seen[n] = true
    }
    return false
}`,
	},
	{
		Slug: "valid-anagram",
		Options: []Option{
			{Text: "Count character frequencies with a fixed-size array, compare counts — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Sort both strings and compare — O(n log n) time, O(n) space", Rating: Plausible},
			{Text: "Use a hash map to count frequencies of the first string, decrement for the second — O(n) time, O(1) space (26 letters)", Rating: Optimal},
			{Text: "For each character in s, search for and remove it from t — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Compare the sums of ASCII values — different strings can have the same sum, e.g. 'ac' vs 'bb' — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Frequency Count Array
// Time: O(n) | Space: O(1) — fixed 26-letter alphabet
func isAnagram(s string, t string) bool {
    if len(s) != len(t) {
        return false
    }
    var count [26]int
    for i := 0; i < len(s); i++ {
        count[s[i]-'a']++
        count[t[i]-'a']--
    }
    for _, c := range count {
        if c != 0 {
            return false
        }
    }
    return true
}`,
	},
	{
		Slug: "product-of-array-except-self",
		Options: []Option{
			{Text: "Two-pass approach: build prefix products left-to-right, then suffix products right-to-left into the result — O(n) time, O(1) extra space", Rating: Optimal},
			{Text: "Use two separate prefix and suffix product arrays, multiply corresponding entries — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Compute total product, then divide by each element — fails when zeros are present — O(n) time, O(1) space", Rating: Wrong},
			{Text: "For each element, multiply all other elements in a nested loop — O(n^2) time, O(1) space", Rating: Suboptimal},
		},
		Solution: `// Pattern: Prefix/Suffix Products
// Time: O(n) | Space: O(1) extra (output array not counted)
func productExceptSelf(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    result[0] = 1
    for i := 1; i < n; i++ {
        result[i] = result[i-1] * nums[i-1]
    }
    suffix := 1
    for i := n - 2; i >= 0; i-- {
        suffix *= nums[i+1]
        result[i] *= suffix
    }
    return result
}`,
	},
	{
		Slug: "encode-and-decode-strings",
		Options: []Option{
			{Text: "Length-prefix encoding: store each string as its length + delimiter + content — O(n) time for encode/decode, O(n) space", Rating: Optimal},
			{Text: "Join strings with a non-ASCII delimiter and split on decode — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Use escape characters for a chosen delimiter — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Concatenate all strings directly — impossible to determine word boundaries on decode — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Length-Prefix Encoding
// Time: O(n) | Space: O(n) where n = total characters across all strings
func encode(strs []string) string {
    var buf []byte
    for _, s := range strs {
        buf = append(buf, []byte(fmt.Sprintf("%d#", len(s)))...)
        buf = append(buf, s...)
    }
    return string(buf)
}

func decode(s string) []string {
    var result []string
    i := 0
    for i < len(s) {
        j := i
        for s[j] != '#' {
            j++
        }
        length, _ := strconv.Atoi(s[i:j])
        result = append(result, s[j+1:j+1+length])
        i = j + 1 + length
    }
    return result
}`,
	},
	{
		Slug: "longest-consecutive-sequence",
		Options: []Option{
			{Text: "Use a hash set, for each number that is a sequence start (n-1 not in set), count the streak — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Sort the array then find the longest run of consecutive elements — O(n log n) time, O(1) space", Rating: Plausible},
			{Text: "Union-Find: union consecutive numbers and track component sizes — O(n * alpha(n)) time, O(n) space", Rating: Plausible},
			{Text: "For each number, search outward for n+1, n+2, ... in the array — without a set this is O(n) per search making it O(n^2) overall — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Find min and max, then check each integer in that range — fails when the range is huge but array is sparse — O(max-min) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Hash Set with Sequence Start Detection
// Time: O(n) | Space: O(n)
func longestConsecutive(nums []int) int {
    set := make(map[int]bool, len(nums))
    for _, n := range nums {
        set[n] = true
    }
    best := 0
    for n := range set {
        if !set[n-1] { // start of a sequence
            length := 1
            for set[n+length] {
                length++
            }
            if length > best {
                best = length
            }
        }
    }
    return best
}`,
	},

	// Two Pointers
	{
		Slug: "two-sum-ii-input-array-is-sorted",
		Options: []Option{
			{Text: "Two pointers from both ends: move left pointer right if sum too small, right pointer left if too large — O(n) time, O(1) space", Rating: Optimal},
			{Text: "For each element, binary search for its complement — O(n log n) time, O(1) space", Rating: Plausible},
			{Text: "Use a hash map like regular two-sum — O(n) time, O(n) space — wasteful since array is sorted", Rating: Plausible},
			{Text: "Check every pair with nested loops — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Use binary search on the entire array for the target — the target is a sum of two elements, not a single element to find — O(log n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Two Pointers (inward)
// Time: O(n) | Space: O(1)
func twoSum(numbers []int, target int) []int {
    left, right := 0, len(numbers)-1
    for left < right {
        sum := numbers[left] + numbers[right]
        if sum == target {
            return []int{left + 1, right + 1}
        } else if sum < target {
            left++
        } else {
            right--
        }
    }
    return nil
}`,
	},

	// Sliding Window
	{
		Slug: "longest-repeating-character-replacement",
		Options: []Option{
			{Text: "Sliding window tracking character frequencies — window is valid when windowSize - maxFreq <= k — O(n) time, O(1) space", Rating: Optimal},
			{Text: "For each character, use a sliding window to find the longest substring needing at most k replacements — O(26 * n) time, O(1) space", Rating: Plausible},
			{Text: "Check every substring and count the most frequent character — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Greedily replace the least frequent characters globally — ignores that replacements must be contiguous in a substring — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Sliding Window with Frequency Count
// Time: O(n) | Space: O(1) — 26 letters
func characterReplacement(s string, k int) int {
    var count [26]int
    maxFreq, maxLen, left := 0, 0, 0
    for right := 0; right < len(s); right++ {
        count[s[right]-'A']++
        if count[s[right]-'A'] > maxFreq {
            maxFreq = count[s[right]-'A']
        }
        for (right - left + 1) - maxFreq > k {
            count[s[left]-'A']--
            left++
        }
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}`,
	},

	// Binary Search
	{
		Slug: "koko-eating-bananas",
		Options: []Option{
			{Text: "Binary search on eating speed k from 1 to max(piles) — check if Koko can finish within h hours at speed mid — O(n log m) time, O(1) space", Rating: Optimal},
			{Text: "Binary search but with a tighter upper bound of ceil(sum(piles)/h) — O(n log(sum/h)) time, O(1) space — correct but doesn't improve worst case meaningfully", Rating: Plausible},
			{Text: "Linear search from speed 1 upward until Koko can finish in time — O(n * m) time, O(1) space", Rating: Suboptimal},
			{Text: "Compute total bananas / h as the speed — ignores that each pile is eaten independently with ceiling division — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Binary Search on Answer
// Time: O(n log m) | Space: O(1) where m = max pile size
func minEatingSpeed(piles []int, h int) int {
    lo, hi := 1, 0
    for _, p := range piles {
        if p > hi {
            hi = p
        }
    }
    for lo < hi {
        mid := lo + (hi-lo)/2
        hours := 0
        for _, p := range piles {
            hours += (p + mid - 1) / mid
        }
        if hours <= h {
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}`,
	},
	{
		Slug: "time-based-key-value-store",
		Options: []Option{
			{Text: "Hash map of key to sorted list of (timestamp, value) pairs, binary search for get — O(1) set, O(log n) get, O(n) space", Rating: Optimal},
			{Text: "Hash map of key to another map of timestamp to value, iterate timestamps for get — O(1) set, O(n) get, O(n) space", Rating: Suboptimal},
			{Text: "Store all entries in a single sorted list and linear scan — O(n) per get, O(n) space", Rating: Suboptimal},
			{Text: "Use a hash map with only the latest value per key — loses historical timestamps, get with earlier timestamp fails — O(1) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Hash Map + Binary Search
// Time: O(1) set, O(log n) get | Space: O(n)
type TimeMap struct {
    store map[string][]entry
}
type entry struct {
    timestamp int
    value     string
}

func Constructor() TimeMap {
    return TimeMap{store: make(map[string][]entry)}
}

func (t *TimeMap) Set(key string, value string, timestamp int) {
    t.store[key] = append(t.store[key], entry{timestamp, value})
}

func (t *TimeMap) Get(key string, timestamp int) string {
    entries := t.store[key]
    lo, hi := 0, len(entries)-1
    result := ""
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if entries[mid].timestamp <= timestamp {
            result = entries[mid].value
            lo = mid + 1
        } else {
            hi = mid - 1
        }
    }
    return result
}`,
	},

	// Linked List
	{
		Slug: "reorder-list",
		Options: []Option{
			{Text: "Find middle, reverse second half, merge two halves alternately — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Store all nodes in an array, use two pointers to reorder — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Use a deque to pop from front and back alternately — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Recursively swap the last node to be next after the first — finding the last node each time is O(n) per step — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Reverse the entire list then interleave with original — reversing destroys the original list, can't interleave — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Find Middle + Reverse + Merge
// Time: O(n) | Space: O(1)
func reorderList(head *ListNode) {
    if head == nil || head.Next == nil {
        return
    }
    // Find middle
    slow, fast := head, head
    for fast.Next != nil && fast.Next.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    // Reverse second half
    var prev *ListNode
    curr := slow.Next
    slow.Next = nil
    for curr != nil {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    // Merge two halves
    first, second := head, prev
    for second != nil {
        tmp1, tmp2 := first.Next, second.Next
        first.Next = second
        second.Next = tmp1
        first = tmp1
        second = tmp2
    }
}`,
	},
	{
		Slug: "remove-nth-node-from-end-of-list",
		Options: []Option{
			{Text: "Two pointers: advance fast pointer n steps ahead, then move both until fast reaches end — O(n) time, O(1) space", Rating: Optimal},
			{Text: "First pass to count length, second pass to remove at position length - n — O(n) time, O(1) space", Rating: Plausible},
			{Text: "Store all nodes in an array, remove the target by index — O(n) time, O(n) space — wastes memory copying all nodes", Rating: Suboptimal},
			{Text: "Use recursion and count positions on the way back up the call stack — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Remove the nth node from the start instead of the end — off by one in direction, removes wrong node — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Two Pointers with Gap
// Time: O(n) | Space: O(1)
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    dummy := &ListNode{Next: head}
    fast, slow := dummy, dummy
    for i := 0; i <= n; i++ {
        fast = fast.Next
    }
    for fast != nil {
        fast = fast.Next
        slow = slow.Next
    }
    slow.Next = slow.Next.Next
    return dummy.Next
}`,
	},
	{
		Slug: "merge-k-sorted-lists",
		Options: []Option{
			{Text: "Use a min-heap to always extract the smallest node across all lists — O(N log k) time, O(k) space", Rating: Optimal},
			{Text: "Divide and conquer: repeatedly merge pairs of lists — O(N log k) time, O(1) space", Rating: Optimal},
			{Text: "Merge lists one by one sequentially — O(N * k) time, O(1) space", Rating: Suboptimal},
			{Text: "Collect all values, sort, rebuild the list — O(N log N) time, O(N) space", Rating: Plausible},
			{Text: "Round-robin: take one node from each list in turns — doesn't maintain sorted order — O(N) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Divide and Conquer Merge
// Time: O(N log k) | Space: O(1) where N = total nodes, k = number of lists
func mergeKLists(lists []*ListNode) *ListNode {
    if len(lists) == 0 {
        return nil
    }
    for len(lists) > 1 {
        var merged []*ListNode
        for i := 0; i < len(lists); i += 2 {
            if i+1 < len(lists) {
                merged = append(merged, mergeTwoLists(lists[i], lists[i+1]))
            } else {
                merged = append(merged, lists[i])
            }
        }
        lists = merged
    }
    return lists[0]
}

func mergeTwoLists(l1, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    tail := dummy
    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            tail.Next = l1
            l1 = l1.Next
        } else {
            tail.Next = l2
            l2 = l2.Next
        }
        tail = tail.Next
    }
    if l1 != nil {
        tail.Next = l1
    } else {
        tail.Next = l2
    }
    return dummy.Next
}`,
	},

	{
		Slug: "delete-the-middle-node-of-a-linked-list",
		Options: []Option{
			{Text: "Fast and slow pointers: advance fast two steps and slow one step, use a prev pointer to delete the middle node when fast reaches the end — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Two-pass: first pass counts the length, second pass deletes the node at index n/2 — O(n) time, O(1) space", Rating: Plausible},
			{Text: "Store all nodes in an array, remove the middle by index, rebuild links — O(n) time, O(n) space — wastes memory copying all nodes", Rating: Suboptimal},
			{Text: "Use recursion counting positions on the way back to find and skip the middle node — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Advance a single pointer by half the first node's value to find the middle — the node values have no relation to list length or position — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Fast and Slow Pointers
// Time: O(n) | Space: O(1)
func deleteMiddle(head *ListNode) *ListNode {
    if head.Next == nil {
        return nil
    }
    slow, fast := head, head.Next.Next
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    slow.Next = slow.Next.Next
    return head
}`,
	},

	// Trees
	{
		Slug: "same-tree",
		Options: []Option{
			{Text: "Recursive DFS: compare values and recurse on both children — O(n) time, O(h) space", Rating: Optimal},
			{Text: "Iterative BFS with two queues, compare level by level — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Serialize both trees to strings and compare — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Compare in-order traversals of both trees — different trees can have the same in-order traversal — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Recursive DFS
// Time: O(n) | Space: O(h) where h = tree height
func isSameTree(p *TreeNode, q *TreeNode) bool {
    if p == nil && q == nil {
        return true
    }
    if p == nil || q == nil || p.Val != q.Val {
        return false
    }
    return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
}`,
	},
	{
		Slug: "subtree-of-another-tree",
		Options: []Option{
			{Text: "For each node in the main tree, check if the subtree rooted there equals subRoot using recursive comparison — O(m * n) time, O(h) space", Rating: Optimal},
			{Text: "Serialize both trees and check if one string contains the other — O(m + n) time, O(m + n) space", Rating: Plausible},
			{Text: "Hash each subtree using Merkle-style hashing, compare hash of subRoot to all subtree hashes — O(m + n) time, O(m) space", Rating: Plausible},
			{Text: "Compare only the root values and immediate children — misses deeper structural differences — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Recursive DFS with Tree Matching
// Time: O(m * n) | Space: O(h) where m, n = sizes of root and subRoot
func isSubtree(root *TreeNode, subRoot *TreeNode) bool {
    if root == nil {
        return false
    }
    if isSame(root, subRoot) {
        return true
    }
    return isSubtree(root.Left, subRoot) || isSubtree(root.Right, subRoot)
}

func isSame(a, b *TreeNode) bool {
    if a == nil && b == nil {
        return true
    }
    if a == nil || b == nil || a.Val != b.Val {
        return false
    }
    return isSame(a.Left, b.Left) && isSame(a.Right, b.Right)
}`,
	},
	{
		Slug: "lowest-common-ancestor-of-a-binary-search-tree",
		Options: []Option{
			{Text: "Exploit BST property: if both values are smaller go left, both larger go right, otherwise current node is LCA — O(h) time, O(1) space", Rating: Optimal},
			{Text: "Recursive version of the same BST property traversal — O(h) time, O(h) space", Rating: Plausible},
			{Text: "Find paths from root to both nodes, compare paths to find divergence point — O(h) time, O(h) space", Rating: Plausible},
			{Text: "Use the generic binary tree LCA algorithm (check left and right subtrees) — correct but ignores BST property, always visits O(n) nodes — O(n) time, O(h) space", Rating: Suboptimal},
			{Text: "Return the node with the smaller value — the LCA is not necessarily the smaller of the two nodes — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: BST Property Traversal
// Time: O(h) | Space: O(1)
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    curr := root
    for curr != nil {
        if p.Val < curr.Val && q.Val < curr.Val {
            curr = curr.Left
        } else if p.Val > curr.Val && q.Val > curr.Val {
            curr = curr.Right
        } else {
            return curr
        }
    }
    return nil
}`,
	},
	{
		Slug: "validate-binary-search-tree",
		Options: []Option{
			{Text: "Recursive DFS passing min/max bounds — each node must be within (min, max) — O(n) time, O(h) space", Rating: Optimal},
			{Text: "In-order traversal and check that the result is strictly increasing — O(n) time, O(n) space", Rating: Plausible},
			{Text: "In-order traversal with a prev pointer, check each node > prev — O(n) time, O(h) space", Rating: Optimal},
			{Text: "Check only that each node's value is greater than its left child and less than its right child — misses violations deeper in the tree — O(n) time, O(h) space", Rating: Wrong},
		},
		Solution: `// Pattern: Recursive DFS with Min/Max Bounds
// Time: O(n) | Space: O(h)
func isValidBST(root *TreeNode) bool {
    return validate(root, math.MinInt64, math.MaxInt64)
}

func validate(node *TreeNode, minVal, maxVal int) bool {
    if node == nil {
        return true
    }
    if node.Val <= minVal || node.Val >= maxVal {
        return false
    }
    return validate(node.Left, minVal, node.Val) && validate(node.Right, node.Val, maxVal)
}`,
	},
	{
		Slug: "kth-smallest-element-in-a-bst",
		Options: []Option{
			{Text: "In-order traversal (iterative with stack), return the kth element visited — O(H + k) time, O(H) space", Rating: Optimal},
			{Text: "Recursive in-order traversal collecting all elements, return the kth — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Augment BST nodes with subtree sizes for O(H) lookup — O(H) time, O(n) space for augmentation", Rating: Plausible},
			{Text: "Convert BST to sorted array via in-order, return arr[k-1] — O(n) time, O(n) space", Rating: Suboptimal},
			{Text: "Do a level-order traversal and pick the kth element — level order doesn't produce sorted order — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Iterative In-Order Traversal
// Time: O(H + k) | Space: O(H)
func kthSmallest(root *TreeNode, k int) int {
    stack := []*TreeNode{}
    curr := root
    for curr != nil || len(stack) > 0 {
        for curr != nil {
            stack = append(stack, curr)
            curr = curr.Left
        }
        curr = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        k--
        if k == 0 {
            return curr.Val
        }
        curr = curr.Right
    }
    return -1
}`,
	},
	{
		Slug: "construct-binary-tree-from-preorder-and-inorder-traversal",
		Options: []Option{
			{Text: "Recursion: first element of preorder is root, find it in inorder to split left/right subtrees, use a hash map for O(1) lookup — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Same recursive approach but linear search in inorder each time — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Iterative approach using a stack to build the tree from preorder and inorder — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Build a BST from preorder alone — this only works for BSTs, not general binary trees — O(n log n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Recursive Build with Index Map
// Time: O(n) | Space: O(n)
func buildTree(preorder []int, inorder []int) *TreeNode {
    inMap := make(map[int]int)
    for i, v := range inorder {
        inMap[v] = i
    }
    preIdx := 0
    var build func(lo, hi int) *TreeNode
    build = func(lo, hi int) *TreeNode {
        if lo > hi {
            return nil
        }
        rootVal := preorder[preIdx]
        preIdx++
        node := &TreeNode{Val: rootVal}
        mid := inMap[rootVal]
        node.Left = build(lo, mid-1)
        node.Right = build(mid+1, hi)
        return node
    }
    return build(0, len(inorder)-1)
}`,
	},
	{
		Slug: "binary-tree-maximum-path-sum",
		Options: []Option{
			{Text: "DFS returning max single-path gain from each node, update global max with left + node + right at each step — O(n) time, O(h) space", Rating: Optimal},
			{Text: "For each node, compute max downward path from it via separate DFS calls, then combine left + node + right — O(n^2) time from redundant traversals, O(h) space", Rating: Plausible},
			{Text: "Enumerate all paths between every pair of nodes — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Only consider root-to-leaf paths — misses paths that don't go through the root or don't end at leaves — O(n) time, O(h) space", Rating: Wrong},
		},
		Solution: `// Pattern: DFS with Global Max Update
// Time: O(n) | Space: O(h)
func maxPathSum(root *TreeNode) int {
    maxSum := math.MinInt64
    var dfs func(*TreeNode) int
    dfs = func(node *TreeNode) int {
        if node == nil {
            return 0
        }
        left := max(dfs(node.Left), 0)
        right := max(dfs(node.Right), 0)
        currentPath := node.Val + left + right
        if currentPath > maxSum {
            maxSum = currentPath
        }
        return node.Val + max(left, right)
    }
    dfs(root)
    return maxSum
}`,
	},
	{
		Slug: "serialize-and-deserialize-binary-tree",
		Options: []Option{
			{Text: "Preorder DFS with null markers: serialize to comma-separated values, deserialize recursively consuming tokens — O(n) time, O(n) space", Rating: Optimal},
			{Text: "BFS level-order with null markers — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Use both preorder and inorder arrays to serialize/deserialize — requires two traversals and no duplicate values — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Only serialize non-null values — without null markers, tree structure is ambiguous and cannot be reconstructed — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Preorder DFS with Null Markers
// Time: O(n) | Space: O(n)
type Codec struct{}

func (c *Codec) serialize(root *TreeNode) string {
    var sb []string
    var dfs func(*TreeNode)
    dfs = func(node *TreeNode) {
        if node == nil {
            sb = append(sb, "N")
            return
        }
        sb = append(sb, strconv.Itoa(node.Val))
        dfs(node.Left)
        dfs(node.Right)
    }
    dfs(root)
    return strings.Join(sb, ",")
}

func (c *Codec) deserialize(data string) *TreeNode {
    tokens := strings.Split(data, ",")
    idx := 0
    var dfs func() *TreeNode
    dfs = func() *TreeNode {
        if tokens[idx] == "N" {
            idx++
            return nil
        }
        val, _ := strconv.Atoi(tokens[idx])
        idx++
        node := &TreeNode{Val: val}
        node.Left = dfs()
        node.Right = dfs()
        return node
    }
    return dfs()
}`,
	},

	// Heap / Priority Queue
	{
		Slug: "find-median-from-data-stream",
		Options: []Option{
			{Text: "Two heaps: max-heap for lower half, min-heap for upper half, balance sizes — O(log n) addNum, O(1) findMedian, O(n) space", Rating: Optimal},
			{Text: "Maintain a sorted array with binary search insertion — O(n) addNum (shifting), O(1) findMedian, O(n) space", Rating: Plausible},
			{Text: "Store all numbers in an unsorted array, sort on each findMedian call — O(1) addNum, O(n log n) findMedian, O(n) space", Rating: Suboptimal},
			{Text: "Keep a running average — the average is not the median — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Two Heaps (Max-Heap + Min-Heap)
// Time: O(log n) addNum, O(1) findMedian | Space: O(n)
type MedianFinder struct {
    lo MaxHeap // max-heap for lower half
    hi MinHeap // min-heap for upper half
}

func Constructor() MedianFinder {
    return MedianFinder{}
}

func (mf *MedianFinder) AddNum(num int) {
    heap.Push(&mf.lo, num)
    heap.Push(&mf.hi, heap.Pop(&mf.lo).(int))
    if mf.lo.Len() < mf.hi.Len() {
        heap.Push(&mf.lo, heap.Pop(&mf.hi).(int))
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.lo.Len() > mf.hi.Len() {
        return float64(mf.lo[0])
    }
    return float64(mf.lo[0]+mf.hi[0]) / 2.0
}

type MaxHeap []int
func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool   { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{})  { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{}    { old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x }

type MinHeap []int
func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool   { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)        { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{})  { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{}    { old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x }`,
	},

	// Graphs
	{
		Slug: "pacific-atlantic-water-flow",
		Options: []Option{
			{Text: "BFS/DFS from ocean borders inward: find cells reachable from Pacific and Atlantic separately, return the intersection — O(m*n) time, O(m*n) space", Rating: Optimal},
			{Text: "DFS from each cell with memoization caching which oceans each cell can reach — O(m*n) time, O(m*n) space — correct but more complex than border-inward approach", Rating: Plausible},
			{Text: "DFS from each cell checking if it can reach both oceans without memoization — O((m*n)^2) time, O(m*n) space", Rating: Suboptimal},
			{Text: "Only check cells on the border — misses interior cells that can flow to both oceans — O(m+n) time, O(m+n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Multi-Source BFS/DFS from Borders
// Time: O(m * n) | Space: O(m * n)
func pacificAtlantic(heights [][]int) [][]int {
    if len(heights) == 0 {
        return nil
    }
    rows, cols := len(heights), len(heights[0])
    pacific := make([][]bool, rows)
    atlantic := make([][]bool, rows)
    for i := range pacific {
        pacific[i] = make([]bool, cols)
        atlantic[i] = make([]bool, cols)
    }
    dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
    var dfs func(r, c int, reachable [][]bool)
    dfs = func(r, c int, reachable [][]bool) {
        reachable[r][c] = true
        for _, d := range dirs {
            nr, nc := r+d[0], c+d[1]
            if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !reachable[nr][nc] && heights[nr][nc] >= heights[r][c] {
                dfs(nr, nc, reachable)
            }
        }
    }
    for r := 0; r < rows; r++ {
        dfs(r, 0, pacific)
        dfs(r, cols-1, atlantic)
    }
    for c := 0; c < cols; c++ {
        dfs(0, c, pacific)
        dfs(rows-1, c, atlantic)
    }
    var result [][]int
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if pacific[r][c] && atlantic[r][c] {
                result = append(result, []int{r, c})
            }
        }
    }
    return result
}`,
	},
	{
		Slug: "course-schedule-ii",
		Options: []Option{
			{Text: "Topological sort using Kahn's algorithm (BFS with in-degree tracking) — O(V+E) time, O(V+E) space", Rating: Optimal},
			{Text: "DFS-based topological sort with three-state cycle detection, append to result in post-order — O(V+E) time, O(V+E) space", Rating: Optimal},
			{Text: "Repeatedly scan for a node with in-degree 0, remove it and update neighbors, repeat — O(V*(V+E)) time, O(V+E) space — correct but re-scanning each time is wasteful vs a queue", Rating: Plausible},
			{Text: "Try all permutations and check if each ordering is valid — O(V!) time, O(V) space", Rating: Suboptimal},
			{Text: "Sort courses by number of prerequisites — courses with equal prerequisites can still have ordering constraints between them — O(V log V) time, O(V) space", Rating: Wrong},
		},
		Solution: `// Pattern: Topological Sort (Kahn's Algorithm / BFS)
// Time: O(V + E) | Space: O(V + E)
func findOrder(numCourses int, prerequisites [][]int) []int {
    graph := make([][]int, numCourses)
    inDegree := make([]int, numCourses)
    for _, p := range prerequisites {
        graph[p[1]] = append(graph[p[1]], p[0])
        inDegree[p[0]]++
    }
    queue := []int{}
    for i := 0; i < numCourses; i++ {
        if inDegree[i] == 0 {
            queue = append(queue, i)
        }
    }
    var order []int
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)
        for _, next := range graph[node] {
            inDegree[next]--
            if inDegree[next] == 0 {
                queue = append(queue, next)
            }
        }
    }
    if len(order) != numCourses {
        return []int{}
    }
    return order
}`,
	},
	{
		Slug: "graph-valid-tree",
		Options: []Option{
			{Text: "Check edges == n-1 and all nodes are connected via BFS/DFS — O(V+E) time, O(V+E) space", Rating: Optimal},
			{Text: "Union-Find: union each edge, if a union finds both nodes in the same set there's a cycle — O(V+E * alpha(V)) time, O(V) space", Rating: Optimal},
			{Text: "DFS checking for back edges (cycle detection) and that the graph is connected — O(V+E) time, O(V+E) space", Rating: Plausible},
			{Text: "Just check that the number of edges equals n-1 — this is necessary but not sufficient; the graph could be disconnected — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Union-Find (Graph = n-1 edges + connected)
// Time: O(V + E * alpha(V)) | Space: O(V)
func validTree(n int, edges [][]int) bool {
    if len(edges) != n-1 {
        return false
    }
    parent := make([]int, n)
    for i := range parent {
        parent[i] = i
    }
    var find func(int) int
    find = func(x int) int {
        if parent[x] != x {
            parent[x] = find(parent[x])
        }
        return parent[x]
    }
    for _, e := range edges {
        px, py := find(e[0]), find(e[1])
        if px == py {
            return false
        }
        parent[px] = py
    }
    return true
}`,
	},
	{
		Slug: "number-of-connected-components-in-an-undirected-graph",
		Options: []Option{
			{Text: "Union-Find: union each edge, count distinct roots — O(V + E * alpha(V)) time, O(V) space", Rating: Optimal},
			{Text: "BFS/DFS from each unvisited node, count the number of traversals — O(V+E) time, O(V+E) space", Rating: Optimal},
			{Text: "Floyd-Warshall to compute all-pairs reachability, then count groups of mutually reachable nodes — O(V^3) time, O(V^2) space — correct but massively overkill", Rating: Plausible},
			{Text: "For each pair of nodes, run DFS to check connectivity, group into components — O(V^2 * (V+E)) time, O(V+E) space", Rating: Suboptimal},
			{Text: "Count nodes with degree 0 and assume the rest form one component — multiple components can all have edges — O(V+E) time, O(V) space", Rating: Wrong},
		},
		Solution: `// Pattern: Union-Find
// Time: O(V + E * alpha(V)) | Space: O(V)
func countComponents(n int, edges [][]int) int {
    parent := make([]int, n)
    for i := range parent {
        parent[i] = i
    }
    var find func(int) int
    find = func(x int) int {
        if parent[x] != x {
            parent[x] = find(parent[x])
        }
        return parent[x]
    }
    components := n
    for _, e := range edges {
        px, py := find(e[0]), find(e[1])
        if px != py {
            parent[px] = py
            components--
        }
    }
    return components
}`,
	},

	// Advanced Graphs
	{
		Slug: "alien-dictionary",
		Options: []Option{
			{Text: "Build a directed graph from adjacent word comparisons, then topological sort (BFS or DFS) — O(C) time where C = total chars, O(1) space (26 letters)", Rating: Optimal},
			{Text: "Compare all pairs of words to find orderings, then topological sort — O(N^2 * L) time, O(1) space", Rating: Suboptimal},
			{Text: "Sort the characters by their first occurrence in the word list — first occurrence doesn't determine order — O(C) time, O(1) space", Rating: Wrong},
			{Text: "Use DFS-based topological sort with cycle detection for invalid orderings — O(C) time, O(1) space", Rating: Optimal},
		},
		Solution: `// Pattern: Topological Sort from Adjacent Word Comparisons
// Time: O(C) where C = total characters | Space: O(1) — at most 26 letters
func alienOrder(words []string) string {
    graph := make(map[byte]map[byte]bool)
    inDegree := make(map[byte]int)
    for _, w := range words {
        for i := 0; i < len(w); i++ {
            if graph[w[i]] == nil {
                graph[w[i]] = make(map[byte]bool)
            }
            inDegree[w[i]] += 0 // ensure key exists
        }
    }
    for i := 0; i < len(words)-1; i++ {
        w1, w2 := words[i], words[i+1]
        minLen := len(w1)
        if len(w2) < minLen {
            minLen = len(w2)
        }
        if len(w1) > len(w2) && w1[:minLen] == w2[:minLen] {
            return "" // invalid: prefix comes after longer word
        }
        for j := 0; j < minLen; j++ {
            if w1[j] != w2[j] {
                if !graph[w1[j]][w2[j]] {
                    graph[w1[j]][w2[j]] = true
                    inDegree[w2[j]]++
                }
                break
            }
        }
    }
    queue := []byte{}
    for ch := range inDegree {
        if inDegree[ch] == 0 {
            queue = append(queue, ch)
        }
    }
    var result []byte
    for len(queue) > 0 {
        ch := queue[0]
        queue = queue[1:]
        result = append(result, ch)
        for next := range graph[ch] {
            inDegree[next]--
            if inDegree[next] == 0 {
                queue = append(queue, next)
            }
        }
    }
    if len(result) != len(inDegree) {
        return "" // cycle detected
    }
    return string(result)
}`,
	},

	// 1-D Dynamic Programming
	{
		Slug: "coin-change",
		Options: []Option{
			{Text: "Bottom-up DP: dp[i] = min coins to make amount i, try each coin — O(amount * n) time, O(amount) space", Rating: Optimal},
			{Text: "Top-down recursion with memoization — O(amount * n) time, O(amount) space", Rating: Plausible},
			{Text: "BFS treating each amount as a node, edges are coin denominations — O(amount * n) time, O(amount) space", Rating: Plausible},
			{Text: "Greedy: always use the largest coin possible — fails for cases like coins=[1,3,4] amount=6 (greedy gives 4+1+1=3 coins but 3+3=2) — O(amount) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bottom-Up DP
// Time: O(amount * n) | Space: O(amount)
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := range dp {
        dp[i] = amount + 1
    }
    dp[0] = 0
    for i := 1; i <= amount; i++ {
        for _, coin := range coins {
            if coin <= i && dp[i-coin]+1 < dp[i] {
                dp[i] = dp[i-coin] + 1
            }
        }
    }
    if dp[amount] > amount {
        return -1
    }
    return dp[amount]
}`,
	},
	{
		Slug: "house-robber-ii",
		Options: []Option{
			{Text: "Run house-robber on nums[0:n-1] and nums[1:n], take the max — O(n) time, O(1) space", Rating: Optimal},
			{Text: "DP with a flag tracking whether the first house was robbed — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Try every valid subset of non-adjacent, non-circular houses — O(2^n) time, O(n) space", Rating: Suboptimal},
			{Text: "Just run standard house-robber ignoring the circular constraint — may rob both first and last house which are adjacent in a circle — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Two-Pass House Robber (Circular)
// Time: O(n) | Space: O(1)
func rob(nums []int) int {
    if len(nums) == 1 {
        return nums[0]
    }
    return max(robRange(nums, 0, len(nums)-2), robRange(nums, 1, len(nums)-1))
}

func robRange(nums []int, lo, hi int) int {
    prev, curr := 0, 0
    for i := lo; i <= hi; i++ {
        prev, curr = curr, max(curr, prev+nums[i])
    }
    return curr
}`,
	},
	{
		Slug: "decode-ways",
		Options: []Option{
			{Text: "Bottom-up DP: dp[i] = ways to decode s[0:i], check single digit and two-digit validity — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Top-down recursion with memoization — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Generate all possible decodings recursively without memoization — O(2^n) time, O(n) space", Rating: Suboptimal},
			{Text: "Simply count the number of valid single and double digit substrings — doesn't account for how choices at each position affect subsequent positions — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bottom-Up DP (Fibonacci-like)
// Time: O(n) | Space: O(1)
func numDecodings(s string) int {
    if s[0] == '0' {
        return 0
    }
    prev2, prev1 := 1, 1
    for i := 1; i < len(s); i++ {
        curr := 0
        if s[i] != '0' {
            curr = prev1
        }
        twoDigit := (s[i-1]-'0')*10 + (s[i] - '0')
        if twoDigit >= 10 && twoDigit <= 26 {
            curr += prev2
        }
        prev2, prev1 = prev1, curr
    }
    return prev1
}`,
	},
	{
		Slug: "palindromic-substrings",
		Options: []Option{
			{Text: "Expand around each center (both odd and even length) and count palindromes — O(n^2) time, O(1) space", Rating: Optimal},
			{Text: "DP table where dp[i][j] = whether s[i:j+1] is a palindrome — O(n^2) time, O(n^2) space", Rating: Plausible},
			{Text: "Manacher's algorithm adapted to count palindromes — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Check every substring by reversing and comparing — O(n^3) time, O(n) space", Rating: Suboptimal},
			{Text: "Count characters with even frequency — frequency doesn't determine palindromic substrings — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Expand Around Center
// Time: O(n^2) | Space: O(1)
func countSubstrings(s string) int {
    count := 0
    for i := 0; i < len(s); i++ {
        count += expandCount(s, i, i)   // odd length
        count += expandCount(s, i, i+1) // even length
    }
    return count
}

func expandCount(s string, left, right int) int {
    count := 0
    for left >= 0 && right < len(s) && s[left] == s[right] {
        count++
        left--
        right++
    }
    return count
}`,
	},
	{
		Slug: "word-break",
		Options: []Option{
			{Text: "Bottom-up DP: dp[i] = true if s[0:i] can be segmented, check all word endings — O(n^2 * m) time, O(n) space", Rating: Optimal},
			{Text: "BFS treating each valid prefix endpoint as a node — O(n^2 * m) time, O(n) space", Rating: Plausible},
			{Text: "Top-down recursion with memoization — O(n^2 * m) time, O(n) space", Rating: Plausible},
			{Text: "Greedy: always match the longest dictionary word from the current position — fails when a shorter match enables the rest to be segmented — O(n * m) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bottom-Up DP
// Time: O(n^2 * m) | Space: O(n) where m = max word length
func wordBreak(s string, wordDict []string) bool {
    wordSet := make(map[string]bool)
    for _, w := range wordDict {
        wordSet[w] = true
    }
    dp := make([]bool, len(s)+1)
    dp[0] = true
    for i := 1; i <= len(s); i++ {
        for j := 0; j < i; j++ {
            if dp[j] && wordSet[s[j:i]] {
                dp[i] = true
                break
            }
        }
    }
    return dp[len(s)]
}`,
	},

	// 2-D Dynamic Programming
	{
		Slug: "unique-paths",
		Options: []Option{
			{Text: "DP with a 1D array: dp[j] = number of paths to column j in current row — O(m*n) time, O(n) space", Rating: Optimal},
			{Text: "Math: compute C(m+n-2, m-1) using combinatorics — O(m+n) time, O(1) space", Rating: Optimal},
			{Text: "2D DP table where dp[i][j] = paths to cell (i,j) — O(m*n) time, O(m*n) space", Rating: Plausible},
			{Text: "Recursive exploration of all paths without memoization — O(2^(m+n)) time, O(m+n) space", Rating: Suboptimal},
			{Text: "BFS counting all paths — BFS visits each cell once and doesn't count multiple paths — O(m*n) time, O(m*n) space", Rating: Wrong},
		},
		Solution: `// Pattern: 1-D DP (Space Optimized)
// Time: O(m * n) | Space: O(n)
func uniquePaths(m int, n int) int {
    dp := make([]int, n)
    for j := range dp {
        dp[j] = 1
    }
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            dp[j] += dp[j-1]
        }
    }
    return dp[n-1]
}`,
	},
	{
		Slug: "longest-common-subsequence",
		Options: []Option{
			{Text: "2D DP: dp[i][j] = LCS of text1[0:i] and text2[0:j] — O(m*n) time, O(m*n) space", Rating: Optimal},
			{Text: "Space-optimized DP using two rows — O(m*n) time, O(min(m,n)) space", Rating: Optimal},
			{Text: "Top-down recursion with memoization — O(m*n) time, O(m*n) space", Rating: Plausible},
			{Text: "Generate all subsequences of both strings and find the longest common one — O(2^m * 2^n) time", Rating: Suboptimal},
			{Text: "Find the longest common substring instead — substrings must be contiguous, subsequences don't — O(m*n) time, O(m*n) space", Rating: Wrong},
		},
		Solution: `// Pattern: 2-D DP Table
// Time: O(m * n) | Space: O(m * n)
func longestCommonSubsequence(text1 string, text2 string) int {
    m, n := len(text1), len(text2)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if text1[i-1] == text2[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }
    return dp[m][n]
}`,
	},

	// Greedy
	{
		Slug: "jump-game-ii",
		Options: []Option{
			{Text: "Greedy BFS: track the farthest reachable position and count jumps at each level boundary — O(n) time, O(1) space", Rating: Optimal},
			{Text: "DP where dp[i] = minimum jumps to reach index i — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "BFS treating each index as a node with edges to reachable indices — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Always jump to the index with the maximum value — max value doesn't mean farthest reach from that position — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Greedy BFS (Level-by-Level)
// Time: O(n) | Space: O(1)
func jump(nums []int) int {
    jumps, curEnd, farthest := 0, 0, 0
    for i := 0; i < len(nums)-1; i++ {
        if i+nums[i] > farthest {
            farthest = i + nums[i]
        }
        if i == curEnd {
            jumps++
            curEnd = farthest
        }
    }
    return jumps
}`,
	},

	// Intervals
	{
		Slug: "insert-interval",
		Options: []Option{
			{Text: "Linear scan: add all intervals before the overlap, merge overlapping intervals, add all after — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Binary search to find insertion point, then merge overlapping neighbors — O(n) time (merging), O(n) space", Rating: Plausible},
			{Text: "Add the new interval to the list, sort, then merge all overlapping intervals — O(n log n) time, O(n) space", Rating: Suboptimal},
			{Text: "Insert at the position where start fits and ignore overlaps — fails to merge overlapping intervals — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Linear Scan and Merge
// Time: O(n) | Space: O(n)
func insert(intervals [][]int, newInterval []int) [][]int {
    var result [][]int
    i := 0
    // Add all intervals before the new interval
    for i < len(intervals) && intervals[i][1] < newInterval[0] {
        result = append(result, intervals[i])
        i++
    }
    // Merge overlapping intervals
    for i < len(intervals) && intervals[i][0] <= newInterval[1] {
        if intervals[i][0] < newInterval[0] {
            newInterval[0] = intervals[i][0]
        }
        if intervals[i][1] > newInterval[1] {
            newInterval[1] = intervals[i][1]
        }
        i++
    }
    result = append(result, newInterval)
    // Add remaining intervals
    for i < len(intervals) {
        result = append(result, intervals[i])
        i++
    }
    return result
}`,
	},
	{
		Slug: "merge-intervals",
		Options: []Option{
			{Text: "Sort by start time, iterate and merge overlapping intervals — O(n log n) time, O(n) space", Rating: Optimal},
			{Text: "Use a timeline/sweep line approach marking starts and ends — O(n log n) time, O(n) space", Rating: Plausible},
			{Text: "Compare every pair of intervals and merge if overlapping — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Merge intervals in the order they appear without sorting — non-adjacent intervals may overlap while adjacent ones don't, so merging sequentially misses overlaps — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Sort + Linear Merge
// Time: O(n log n) | Space: O(n)
func merge(intervals [][]int) [][]int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })
    var result [][]int
    for _, interval := range intervals {
        if len(result) > 0 && result[len(result)-1][1] >= interval[0] {
            if interval[1] > result[len(result)-1][1] {
                result[len(result)-1][1] = interval[1]
            }
        } else {
            result = append(result, interval)
        }
    }
    return result
}`,
	},
	{
		Slug: "non-overlapping-intervals",
		Options: []Option{
			{Text: "Greedy: sort by end time, always keep the interval that ends earliest — count removals — O(n log n) time, O(1) space", Rating: Optimal},
			{Text: "Sort by start time and greedily remove the interval with the later end when overlap occurs — O(n log n) time, O(1) space", Rating: Optimal},
			{Text: "DP similar to longest increasing subsequence on intervals — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Remove the shortest intervals first — short intervals aren't necessarily the ones causing overlaps — O(n log n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Greedy (Sort by End Time)
// Time: O(n log n) | Space: O(1)
func eraseOverlapIntervals(intervals [][]int) int {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][1] < intervals[j][1]
    })
    count, prevEnd := 0, intervals[0][1]
    for i := 1; i < len(intervals); i++ {
        if intervals[i][0] < prevEnd {
            count++
        } else {
            prevEnd = intervals[i][1]
        }
    }
    return count
}`,
	},
	{
		Slug: "meeting-rooms",
		Options: []Option{
			{Text: "Sort by start time, check if any meeting starts before the previous one ends — O(n log n) time, O(1) space", Rating: Optimal},
			{Text: "Compare every pair of meetings for overlap — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Check if total meeting time exceeds available time — overlapping meetings don't necessarily exceed total time — O(n) time, O(1) space", Rating: Wrong},
			{Text: "Sort by end time and check consecutive overlaps — sorting by end time also works for pairwise overlap detection — O(n log n) time, O(1) space", Rating: Plausible},
		},
		Solution: `// Pattern: Sort + Linear Scan
// Time: O(n log n) | Space: O(1)
func canAttendMeetings(intervals [][]int) bool {
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })
    for i := 1; i < len(intervals); i++ {
        if intervals[i][0] < intervals[i-1][1] {
            return false
        }
    }
    return true
}`,
	},
	{
		Slug: "meeting-rooms-ii",
		Options: []Option{
			{Text: "Sort start and end times separately, use two pointers to count overlapping meetings — O(n log n) time, O(n) space", Rating: Optimal},
			{Text: "Use a min-heap tracking meeting end times, pop if earliest end <= current start — O(n log n) time, O(n) space", Rating: Optimal},
			{Text: "Sort by start time and count active meetings with a sweep — O(n log n) time, O(n) space", Rating: Plausible},
			{Text: "Count the maximum number of meetings that start at the same time — ignores meetings that overlap without starting simultaneously — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Sweep Line (Sort Start/End Separately)
// Time: O(n log n) | Space: O(n)
func minMeetingRooms(intervals [][]int) int {
    n := len(intervals)
    starts := make([]int, n)
    ends := make([]int, n)
    for i, iv := range intervals {
        starts[i] = iv[0]
        ends[i] = iv[1]
    }
    sort.Ints(starts)
    sort.Ints(ends)
    rooms, endPtr := 0, 0
    for i := 0; i < n; i++ {
        if starts[i] < ends[endPtr] {
            rooms++
        } else {
            endPtr++
        }
    }
    return rooms
}`,
	},

	// Math & Geometry
	{
		Slug: "rotate-image",
		Options: []Option{
			{Text: "Transpose the matrix then reverse each row — O(n^2) time, O(1) space", Rating: Optimal},
			{Text: "Rotate four cells at a time layer by layer from outside in — O(n^2) time, O(1) space", Rating: Optimal},
			{Text: "Create a new matrix and copy rotated positions — O(n^2) time, O(n^2) space", Rating: Plausible},
			{Text: "Reverse each row then transpose — this gives a counter-clockwise rotation, not clockwise — O(n^2) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Transpose + Reverse Rows
// Time: O(n^2) | Space: O(1)
func rotate(matrix [][]int) {
    n := len(matrix)
    // Transpose
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
        }
    }
    // Reverse each row
    for i := 0; i < n; i++ {
        for lo, hi := 0, n-1; lo < hi; lo, hi = lo+1, hi-1 {
            matrix[i][lo], matrix[i][hi] = matrix[i][hi], matrix[i][lo]
        }
    }
}`,
	},
	{
		Slug: "spiral-matrix",
		Options: []Option{
			{Text: "Layer-by-layer traversal: shrink boundaries (top, bottom, left, right) after each direction — O(m*n) time, O(1) extra space", Rating: Optimal},
			{Text: "Simulate with direction vectors and a visited matrix — O(m*n) time, O(m*n) space — wastes memory when boundary tracking uses O(1)", Rating: Suboptimal},
			{Text: "Recursively peel off the outer ring and process the inner matrix — O(m*n) time, O(min(m,n)) space", Rating: Plausible},
			{Text: "Traverse row by row alternating direction — this gives a zigzag, not a spiral — O(m*n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Boundary Shrinking
// Time: O(m * n) | Space: O(1) extra
func spiralOrder(matrix [][]int) []int {
    var result []int
    top, bottom := 0, len(matrix)-1
    left, right := 0, len(matrix[0])-1
    for top <= bottom && left <= right {
        for c := left; c <= right; c++ {
            result = append(result, matrix[top][c])
        }
        top++
        for r := top; r <= bottom; r++ {
            result = append(result, matrix[r][right])
        }
        right--
        if top <= bottom {
            for c := right; c >= left; c-- {
                result = append(result, matrix[bottom][c])
            }
            bottom--
        }
        if left <= right {
            for r := bottom; r >= top; r-- {
                result = append(result, matrix[r][left])
            }
            left++
        }
    }
    return result
}`,
	},
	{
		Slug: "set-matrix-zeroes",
		Options: []Option{
			{Text: "Use first row and first column as markers, with two flags for their own zero status — O(m*n) time, O(1) space", Rating: Optimal},
			{Text: "Record zero positions in two sets (rows and columns), then zero out — O(m*n) time, O(m+n) space", Rating: Plausible},
			{Text: "Create a copy of the matrix, scan original, modify copy — O(m*n) time, O(m*n) space", Rating: Suboptimal},
			{Text: "Zero out rows and columns immediately when a zero is found — newly set zeros trigger cascading zeroing of unrelated rows/columns — O(m*n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: In-Place Markers (First Row/Column)
// Time: O(m * n) | Space: O(1)
func setZeroes(matrix [][]int) {
    rows, cols := len(matrix), len(matrix[0])
    firstRowZero := false
    firstColZero := false
    for c := 0; c < cols; c++ {
        if matrix[0][c] == 0 {
            firstRowZero = true
        }
    }
    for r := 0; r < rows; r++ {
        if matrix[r][0] == 0 {
            firstColZero = true
        }
    }
    for r := 1; r < rows; r++ {
        for c := 1; c < cols; c++ {
            if matrix[r][c] == 0 {
                matrix[r][0] = 0
                matrix[0][c] = 0
            }
        }
    }
    for r := 1; r < rows; r++ {
        for c := 1; c < cols; c++ {
            if matrix[r][0] == 0 || matrix[0][c] == 0 {
                matrix[r][c] = 0
            }
        }
    }
    if firstRowZero {
        for c := 0; c < cols; c++ {
            matrix[0][c] = 0
        }
    }
    if firstColZero {
        for r := 0; r < rows; r++ {
            matrix[r][0] = 0
        }
    }
}`,
	},

	// Bit Manipulation
	{
		Slug: "number-of-1-bits",
		Options: []Option{
			{Text: "Brian Kernighan's trick: n &= n-1 clears the lowest set bit, count iterations — O(k) time where k = number of set bits, O(1) space", Rating: Optimal},
			{Text: "Check each of the 32 bits using right shift and bitwise AND — O(32) time, O(1) space", Rating: Plausible},
			{Text: "Convert to binary string and count '1' characters — O(32) time, O(32) space — unnecessary string allocation for a bit problem", Rating: Suboptimal},
			{Text: "Use modulo 2 and division by 2 in a loop — same as bit checking but less idiomatic — O(32) time, O(1) space", Rating: Plausible},
			{Text: "Return the integer value itself — the value of a number is not its popcount — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Brian Kernighan's Bit Trick
// Time: O(k) where k = number of set bits | Space: O(1)
func hammingWeight(n uint32) int {
    count := 0
    for n != 0 {
        n &= n - 1
        count++
    }
    return count
}`,
	},
	{
		Slug: "counting-bits",
		Options: []Option{
			{Text: "DP using the relation dp[i] = dp[i >> 1] + (i & 1) — O(n) time, O(n) space", Rating: Optimal},
			{Text: "DP using dp[i] = dp[i & (i-1)] + 1 (Brian Kernighan relation) — O(n) time, O(n) space", Rating: Optimal},
			{Text: "For each number 0 to n, count bits individually — O(n log n) time, O(n) space", Rating: Suboptimal},
			{Text: "Use the pattern that the answer repeats with an offset every power of 2 — correct but harder to implement — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Return i % 2 for each number — only counts the last bit, not all set bits — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: DP with Bit Shift
// Time: O(n) | Space: O(n)
func countBits(n int) []int {
    dp := make([]int, n+1)
    for i := 1; i <= n; i++ {
        dp[i] = dp[i>>1] + (i & 1)
    }
    return dp
}`,
	},
	{
		Slug: "reverse-bits",
		Options: []Option{
			{Text: "Iterate 32 times: extract the last bit of n, shift result left and OR the bit — O(1) time, O(1) space", Rating: Optimal},
			{Text: "Divide and conquer: swap halves, quarters, etc. using bitmasks — O(1) time, O(1) space", Rating: Plausible},
			{Text: "Convert to binary string, reverse the string, convert back — O(32) time, O(32) space — unnecessary string allocation for a bit problem", Rating: Suboptimal},
			{Text: "XOR the number with 0xFFFFFFFF — this flips bits (NOT), not reverses them — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bit-by-Bit Reversal
// Time: O(1) — always 32 iterations | Space: O(1)
func reverseBits(num uint32) uint32 {
    var result uint32
    for i := 0; i < 32; i++ {
        result = (result << 1) | (num & 1)
        num >>= 1
    }
    return result
}`,
	},
	{
		Slug: "missing-number",
		Options: []Option{
			{Text: "XOR all indices 0..n with all array elements — duplicates cancel, leaving the missing number — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Gauss formula: compute n*(n+1)/2 minus the array sum — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Sort the array and find the first index where nums[i] != i — O(n log n) time, O(1) space — slower than needed and modifies input", Rating: Suboptimal},
			{Text: "Use a hash set of all numbers, check which 0..n is missing — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Return n if the last element isn't n, else return 0 — the missing number could be any value in [0,n] — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: XOR Cancellation
// Time: O(n) | Space: O(1)
func missingNumber(nums []int) int {
    xor := len(nums)
    for i, n := range nums {
        xor ^= i ^ n
    }
    return xor
}`,
	},
	{
		Slug: "sum-of-two-integers",
		Options: []Option{
			{Text: "Use bitwise operations: XOR for sum without carry, AND + left shift for carry, repeat until carry is 0 — O(32) time, O(1) space", Rating: Optimal},
			{Text: "Use repeated increment/decrement by 1 — O(|b|) time, O(1) space", Rating: Suboptimal},
			{Text: "Convert to binary strings and add digit by digit — O(32) time, O(32) space", Rating: Plausible},
			{Text: "Use log and exponent: log(e^a * e^b) = a + b — floating point precision makes this unreliable for integers — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bit Manipulation (XOR + AND for Carry)
// Time: O(1) — at most 32 iterations | Space: O(1)
func getSum(a int, b int) int {
    for b != 0 {
        carry := a & b
        a = a ^ b
        b = carry << 1
    }
    return a
}`,
	},
}
