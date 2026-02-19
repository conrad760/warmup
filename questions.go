package main

// CuratedQuestion holds curated approach options and solution for a leetcode problem.
// Title, Description, Example, Difficulty are loaded from leetgo's database.
// Category overrides the database's topic tags when set.
type CuratedQuestion struct {
	Slug     string
	Category string // When set, overrides the category derived from leetgo's topic tags.
	Options  []Option
	Solution string
}

// curatedBank is the built-in set of approach options and solutions.
var curatedBank = []CuratedQuestion{
	// Arrays & Hashing
	{
		Slug:     "two-sum",
		Category: "Arrays & Hashing",
		Options: []Option{
			{Text: "Use a hash map to store each number's complement as you iterate — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Sort the array then use two pointers — O(n log n) time, O(1) space", Rating: Plausible},
			{Text: "Check every pair of elements with nested loops — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Binary search for each element's complement — O(n log n) time, O(1) space", Rating: Plausible},
			{Text: "Use a set to track seen values then check membership — this only finds whether a pair exists, not their indices — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Hash Map Complement
// Time: O(n) | Space: O(n)
func twoSum(nums []int, target int) []int {
    seen := make(map[int]int) // value -> index
    for i, num := range nums {
        complement := target - num
        if j, ok := seen[complement]; ok {
            return []int{j, i}
        }
        seen[num] = i
    }
    return nil
}`,
	},
	{
		Slug:     "group-anagrams",
		Category: "Arrays & Hashing",
		Options: []Option{
			{Text: "Sort each string's characters and use the sorted version as a hash map key — O(n * k log k) time, O(n * k) space", Rating: Optimal},
			{Text: "Use a character frequency count array as the hash key for each string — O(n * k) time, O(n * k) space", Rating: Optimal},
			{Text: "Compare every pair of strings by checking if one is an anagram of the other — O(n^2 * k) time, O(n * k) space", Rating: Suboptimal},
			{Text: "Sort the entire array of strings lexicographically and group adjacent anagrams — anagrams won't necessarily be adjacent after sorting — O(n * k log k) time, O(n * k) space", Rating: Wrong},
		},
		Solution: `// Pattern: Hash Map with Sorted Key
// Time: O(n * k log k) | Space: O(n * k)
// where n = number of strings, k = max string length
func groupAnagrams(strs []string) [][]string {
    groups := make(map[string][]string)
    for _, s := range strs {
        key := sortString(s)
        groups[key] = append(groups[key], s)
    }
    result := make([][]string, 0, len(groups))
    for _, group := range groups {
        result = append(result, group)
    }
    return result
}

func sortString(s string) string {
    runes := []rune(s)
    sort.Slice(runes, func(i, j int) bool {
        return runes[i] < runes[j]
    })
    return string(runes)
}`,
	},
	{
		Slug:     "top-k-frequent-elements",
		Category: "Arrays & Hashing",
		Options: []Option{
			{Text: "Count frequencies with a hash map, then use bucket sort where index = frequency — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Count frequencies then use a min-heap of size k — O(n log k) time, O(n) space", Rating: Plausible},
			{Text: "Count frequencies then sort by frequency descending — O(n log n) time, O(n) space", Rating: Suboptimal},
			{Text: "Use nested loops to find the most frequent element k times, removing each after finding — O(n * k) time, O(n) space", Rating: Suboptimal},
			{Text: "Sort the array then count consecutive runs to find top k — misses that sorted order != frequency order — O(n log n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bucket Sort by Frequency
// Time: O(n) | Space: O(n)
func topKFrequent(nums []int, k int) []int {
    freq := make(map[int]int)
    for _, n := range nums {
        freq[n]++
    }
    buckets := make([][]int, len(nums)+1)
    for num, count := range freq {
        buckets[count] = append(buckets[count], num)
    }
    result := make([]int, 0, k)
    for i := len(buckets) - 1; i >= 0 && len(result) < k; i-- {
        result = append(result, buckets[i]...)
    }
    return result[:k]
}`,
	},

	// Two Pointers
	{
		Slug:     "valid-palindrome",
		Category: "Two Pointers",
		Options: []Option{
			{Text: "Use two pointers from both ends, skip non-alphanumeric, compare lowercase — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Filter the string to keep only alphanumeric chars, lowercase it, then reverse and compare — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Convert to lowercase, use a stack to push first half then compare with second half — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Compare each character to its mirror position using index math — fails because non-alpha chars shift positions — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Two Pointers (inward)
// Time: O(n) | Space: O(1)
func isPalindrome(s string) bool {
    left, right := 0, len(s)-1
    for left < right {
        for left < right && !isAlphaNum(s[left]) {
            left++
        }
        for left < right && !isAlphaNum(s[right]) {
            right--
        }
        if toLower(s[left]) != toLower(s[right]) {
            return false
        }
        left++
        right--
    }
    return true
}`,
	},
	{
		Slug:     "3sum",
		Category: "Two Pointers",
		Options: []Option{
			{Text: "Sort the array, fix one element, use two pointers for the remaining pair, skip duplicates — O(n^2) time, O(1) space", Rating: Optimal},
			{Text: "Use three nested loops and a set to deduplicate results — O(n^3) time, O(n) space", Rating: Suboptimal},
			{Text: "Fix one element, use a hash set to find the complement pair — O(n^2) time, O(n) space", Rating: Plausible},
			{Text: "Sort the array and use binary search for the third element after fixing two — O(n^2 log n) time, O(1) space", Rating: Plausible},
			{Text: "Use divide and conquer to split the array and find triplets in halves — triplets can span halves making this incorrect — O(n^2) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Sort + Two Pointers
// Time: O(n^2) | Space: O(1) excluding output
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    var result [][]int
    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }
        lo, hi := i+1, len(nums)-1
        for lo < hi {
            sum := nums[i] + nums[lo] + nums[hi]
            if sum == 0 {
                result = append(result, []int{nums[i], nums[lo], nums[hi]})
                for lo < hi && nums[lo] == nums[lo+1] { lo++ }
                for lo < hi && nums[hi] == nums[hi-1] { hi-- }
                lo++; hi--
            } else if sum < 0 {
                lo++
            } else {
                hi--
            }
        }
    }
    return result
}`,
	},
	{
		Slug:     "container-with-most-water",
		Category: "Two Pointers",
		Options: []Option{
			{Text: "Two pointers starting at both ends, move the shorter side inward — O(n) time, O(1) space", Rating: Optimal},
			{Text: "For each line, binary search for the best partner by height among lines farther than some distance — O(n log n) time, O(1) space", Rating: Plausible},
			{Text: "Check every pair of lines and compute the area — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Use a stack to find next greater elements then compute areas — stack approach doesn't apply to this container geometry — O(n) time, O(n) space", Rating: Wrong},
			{Text: "Sort lines by height and greedily pick tallest pairs — sorting loses position information needed for width — O(n log n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Two Pointers (greedy inward)
// Time: O(n) | Space: O(1)
func maxArea(height []int) int {
    left, right := 0, len(height)-1
    maxWater := 0
    for left < right {
        h := min(height[left], height[right])
        area := h * (right - left)
        if area > maxWater { maxWater = area }
        if height[left] < height[right] {
            left++
        } else {
            right--
        }
    }
    return maxWater
}`,
	},

	// Sliding Window
	{
		Slug:     "best-time-to-buy-and-sell-stock",
		Category: "Sliding Window",
		Options: []Option{
			{Text: "Track the minimum price seen so far and compute max profit at each step — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Check every buy-sell pair — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Sort prices and take the difference between max and min — sorting loses the temporal buy-before-sell constraint — O(n log n) time, O(1) space", Rating: Wrong},
			{Text: "Use Kadane's algorithm on the daily price differences — O(n) time, O(1) space", Rating: Plausible},
		},
		Solution: `// Pattern: Sliding Window / Running Minimum
// Time: O(n) | Space: O(1)
func maxProfit(prices []int) int {
    minPrice := prices[0]
    maxProf := 0
    for _, price := range prices {
        if price < minPrice { minPrice = price }
        if price-minPrice > maxProf { maxProf = price - minPrice }
    }
    return maxProf
}`,
	},
	{
		Slug:     "longest-substring-without-repeating-characters",
		Category: "Sliding Window",
		Options: []Option{
			{Text: "Sliding window with a hash map storing each character's latest index — O(n) time, O(min(n,m)) space where m is charset size", Rating: Optimal},
			{Text: "Check every possible substring for uniqueness — O(n^3) time, O(min(n,m)) space", Rating: Suboptimal},
			{Text: "Sliding window with a hash set, shrink window from left when duplicate found — O(n) time, O(min(n,m)) space", Rating: Plausible},
			{Text: "Sort the string and find the longest run of unique characters — sorting destroys substring ordering — O(n log n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Sliding Window with Hash Map
// Time: O(n) | Space: O(min(n, m)) where m = charset size
func lengthOfLongestSubstring(s string) int {
    lastSeen := make(map[byte]int)
    maxLen, left := 0, 0
    for right := 0; right < len(s); right++ {
        if idx, ok := lastSeen[s[right]]; ok && idx >= left {
            left = idx + 1
        }
        lastSeen[s[right]] = right
        if right-left+1 > maxLen { maxLen = right - left + 1 }
    }
    return maxLen
}`,
	},
	{
		Slug:     "minimum-window-substring",
		Category: "Sliding Window",
		Options: []Option{
			{Text: "Sliding window with two frequency maps, expand right until valid then shrink left — O(n + m) time, O(m) space", Rating: Optimal},
			{Text: "Check every substring of s to see if it contains all of t — O(n^2 * m) time, O(m) space", Rating: Suboptimal},
			{Text: "For each character in t, find its position in s and combine intervals — doesn't handle duplicates in t correctly — O(n * m) time, O(m) space", Rating: Wrong},
			{Text: "Use two pointers with a counter tracking how many characters are still needed — O(n + m) time, O(m) space", Rating: Plausible},
		},
		Solution: `// Pattern: Sliding Window with Frequency Count
// Time: O(n + m) | Space: O(m)
func minWindow(s string, t string) string {
    if len(t) > len(s) { return "" }
    need := make(map[byte]int)
    for i := 0; i < len(t); i++ { need[t[i]]++ }
    window := make(map[byte]int)
    have, total := 0, len(need)
    bestLen, bestStart, left := len(s)+1, 0, 0
    for right := 0; right < len(s); right++ {
        c := s[right]
        window[c]++
        if need[c] > 0 && window[c] == need[c] { have++ }
        for have == total {
            if right-left+1 < bestLen {
                bestLen = right - left + 1
                bestStart = left
            }
            lc := s[left]
            window[lc]--
            if need[lc] > 0 && window[lc] < need[lc] { have-- }
            left++
        }
    }
    if bestLen > len(s) { return "" }
    return s[bestStart : bestStart+bestLen]
}`,
	},

	// Stack
	{
		Slug:     "valid-parentheses",
		Category: "Stack",
		Options: []Option{
			{Text: "Use a stack: push opening brackets, pop and match for closing brackets — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Use a stack but push the expected closing bracket for each opener, then check equality on pop — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Repeatedly remove adjacent matching pairs until string is empty or stuck — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Count opening and closing brackets of each type — counts match doesn't mean nesting is correct (e.g. \"][\") — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Stack
// Time: O(n) | Space: O(n)
func isValid(s string) bool {
    stack := []byte{}
    matching := map[byte]byte{')': '(', ']': '[', '}': '{'}
    for i := 0; i < len(s); i++ {
        c := s[i]
        if c == '(' || c == '[' || c == '{' {
            stack = append(stack, c)
        } else {
            if len(stack) == 0 || stack[len(stack)-1] != matching[c] {
                return false
            }
            stack = stack[:len(stack)-1]
        }
    }
    return len(stack) == 0
}`,
	},
	{
		Slug:     "min-stack",
		Category: "Stack",
		Options: []Option{
			{Text: "Use two stacks: one for values and one tracking the current minimum at each level — O(1) per operation, O(n) space", Rating: Optimal},
			{Text: "Store (value, currentMin) pairs in a single stack — O(1) per operation, O(n) space", Rating: Optimal},
			{Text: "Use a stack plus a min-heap to track the current minimum — O(log n) push/pop due to heap operations, O(n) space", Rating: Plausible},
			{Text: "Keep a single stack and scan for minimum on each getMin call — O(n) per getMin, O(n) space", Rating: Suboptimal},
			{Text: "Store only the minimum value in a variable — breaks when the minimum is popped and you need the next minimum — O(1) per operation, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Auxiliary Stack / Pair Stack
// Time: O(1) per operation | Space: O(n)
type MinStack struct {
    stack []entry
}
type entry struct { val, minVal int }

func (s *MinStack) Push(val int) {
    minVal := val
    if len(s.stack) > 0 && s.stack[len(s.stack)-1].minVal < val {
        minVal = s.stack[len(s.stack)-1].minVal
    }
    s.stack = append(s.stack, entry{val, minVal})
}
func (s *MinStack) Pop()        { s.stack = s.stack[:len(s.stack)-1] }
func (s *MinStack) Top() int    { return s.stack[len(s.stack)-1].val }
func (s *MinStack) GetMin() int { return s.stack[len(s.stack)-1].minVal }`,
	},
	{
		Slug:     "evaluate-reverse-polish-notation",
		Category: "Stack",
		Options: []Option{
			{Text: "Use a stack: push numbers, pop two operands when hitting an operator, push result — O(n) time, O(n) space", Rating: Optimal},
			{Text: "Convert RPN to infix expression then evaluate — unnecessarily complex, still O(n) time, O(n) space", Rating: Plausible},
			{Text: "Recursively parse from right to left, treating last token as root operator — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Scan left to right, find first operator adjacent to two numbers, evaluate in-place — O(n^2) time due to shifting, O(1) space", Rating: Suboptimal},
			{Text: "Evaluate left to right applying each operator to a running total — ignores operator precedence and operand ordering — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Stack-based Expression Evaluation
// Time: O(n) | Space: O(n)
func evalRPN(tokens []string) int {
    stack := []int{}
    for _, token := range tokens {
        switch token {
        case "+", "-", "*", "/":
            b := stack[len(stack)-1]
            a := stack[len(stack)-2]
            stack = stack[:len(stack)-2]
            switch token {
            case "+": stack = append(stack, a+b)
            case "-": stack = append(stack, a-b)
            case "*": stack = append(stack, a*b)
            case "/": stack = append(stack, a/b)
            }
        default:
            num, _ := strconv.Atoi(token)
            stack = append(stack, num)
        }
    }
    return stack[0]
}`,
	},

	// Binary Search
	{
		Slug:     "binary-search",
		Category: "Binary Search",
		Options: []Option{
			{Text: "Classic binary search: compare middle element, narrow to left or right half — O(log n) time, O(1) space", Rating: Optimal},
			{Text: "Linear scan through the array — O(n) time, O(1) space", Rating: Suboptimal},
			{Text: "Use a hash set for O(1) lookup — O(n) time to build, O(n) space — wasteful for a sorted array", Rating: Plausible},
			{Text: "Jump search: skip ahead by sqrt(n) then linear scan the block — O(sqrt(n)) time, O(1) space", Rating: Plausible},
			{Text: "Start from both ends and move inward — two pointer search doesn't exploit sorted property efficiently — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Binary Search
// Time: O(log n) | Space: O(1)
func search(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target { return mid }
        if nums[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
    }
    return -1
}`,
	},
	{
		Slug:     "search-in-rotated-sorted-array",
		Category: "Binary Search",
		Options: []Option{
			{Text: "Modified binary search: determine which half is sorted, then check if target lies in that half — O(log n) time, O(1) space", Rating: Optimal},
			{Text: "Find the pivot with binary search, then binary search the correct half — O(log n) time, O(1) space", Rating: Plausible},
			{Text: "Linear scan to find the target — O(n) time, O(1) space", Rating: Suboptimal},
			{Text: "Standard binary search without accounting for rotation — will miss elements in the rotated portion — O(log n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Modified Binary Search
// Time: O(log n) | Space: O(1)
func search(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target { return mid }
        if nums[lo] <= nums[mid] {
            if target >= nums[lo] && target < nums[mid] { hi = mid - 1 } else { lo = mid + 1 }
        } else {
            if target > nums[mid] && target <= nums[hi] { lo = mid + 1 } else { hi = mid - 1 }
        }
    }
    return -1
}`,
	},
	{
		Slug:     "find-minimum-in-rotated-sorted-array",
		Category: "Binary Search",
		Options: []Option{
			{Text: "Binary search comparing mid to right boundary to determine which half contains the minimum — O(log n) time, O(1) space", Rating: Optimal},
			{Text: "Find the pivot point (where nums[i] > nums[i+1]) with binary search, minimum is at pivot+1 — O(log n) time, O(1) space", Rating: Plausible},
			{Text: "Linear scan tracking the minimum — O(n) time, O(1) space", Rating: Suboptimal},
			{Text: "Check only the first and last elements — the minimum is always at one end — this is false for rotated arrays — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Binary Search on Rotated Array
// Time: O(log n) | Space: O(1)
func findMin(nums []int) int {
    lo, hi := 0, len(nums)-1
    for lo < hi {
        mid := lo + (hi-lo)/2
        if nums[mid] > nums[hi] { lo = mid + 1 } else { hi = mid }
    }
    return nums[lo]
}`,
	},

	// Linked List
	{
		Slug:     "reverse-linked-list",
		Category: "Linked List",
		Options: []Option{
			{Text: "Iteratively reverse pointers using prev/curr/next variables — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Recursively reverse: each call returns the new tail and rewires pointers — O(n) time, O(n) space (call stack)", Rating: Plausible},
			{Text: "Store all values in an array, then rebuild the list in reverse — O(n) time, O(n) space", Rating: Suboptimal},
			{Text: "Swap the first and last nodes, then move inward — singly linked list has no backward pointer so finding the last is O(n) each time — O(n^2) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Iterative Pointer Reversal
// Time: O(n) | Space: O(1)
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    curr := head
    for curr != nil {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    return prev
}`,
	},
	{
		Slug:     "merge-two-sorted-lists",
		Category: "Linked List",
		Options: []Option{
			{Text: "Use a dummy head, iterate both lists comparing values, append smaller node — O(n+m) time, O(1) space", Rating: Optimal},
			{Text: "Recursively merge: pick the smaller head, recurse on the rest — O(n+m) time, O(n+m) space (call stack)", Rating: Plausible},
			{Text: "Collect all values into an array, sort, rebuild the list — O((n+m) log(n+m)) time, O(n+m) space", Rating: Suboptimal},
			{Text: "Interleave nodes alternately from list1 and list2 — alternating doesn't maintain sorted order — O(n+m) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Iterative Merge with Dummy Head
// Time: O(n + m) | Space: O(1)
func mergeTwoLists(l1, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    tail := dummy
    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            tail.Next = l1; l1 = l1.Next
        } else {
            tail.Next = l2; l2 = l2.Next
        }
        tail = tail.Next
    }
    if l1 != nil { tail.Next = l1 } else { tail.Next = l2 }
    return dummy.Next
}`,
	},
	{
		Slug:     "linked-list-cycle",
		Category: "Linked List",
		Options: []Option{
			{Text: "Floyd's cycle detection: slow pointer moves 1 step, fast pointer moves 2 steps, they meet if cycle exists — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Use a hash set to store visited nodes, check for revisits — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Store each visited node address in a sorted list and binary search for duplicates — O(n log n) time, O(n) space", Rating: Suboptimal},
			{Text: "Traverse the list and count nodes; if count exceeds a threshold assume cycle — unreliable heuristic with no guarantee — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Floyd's Tortoise and Hare
// Time: O(n) | Space: O(1)
func hasCycle(head *ListNode) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast { return true }
    }
    return false
}`,
	},

	// Trees
	{
		Slug:     "invert-binary-tree",
		Category: "Trees",
		Options: []Option{
			{Text: "Recursive DFS: swap left and right children at each node, recurse on both — O(n) time, O(h) space", Rating: Optimal},
			{Text: "Iterative BFS with a queue: swap children level by level — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Collect all nodes in-order, rebuild a mirrored tree from the values — O(n) time, O(n) space", Rating: Suboptimal},
			{Text: "Only swap the values at left and right nodes without restructuring pointers — fails when subtrees have different shapes — O(n) time, O(h) space", Rating: Wrong},
		},
		Solution: `// Pattern: Recursive DFS
// Time: O(n) | Space: O(h) where h = tree height
func invertTree(root *TreeNode) *TreeNode {
    if root == nil { return nil }
    root.Left, root.Right = root.Right, root.Left
    invertTree(root.Left)
    invertTree(root.Right)
    return root
}`,
	},
	{
		Slug:     "maximum-depth-of-binary-tree",
		Category: "Trees",
		Options: []Option{
			{Text: "Recursive DFS: return 1 + max(depth(left), depth(right)) — O(n) time, O(h) space", Rating: Optimal},
			{Text: "Iterative BFS counting levels — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Iterative DFS with a stack tracking depth — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Count the total number of nodes and compute log2(n) — only works for complete binary trees — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Recursive DFS
// Time: O(n) | Space: O(h)
func maxDepth(root *TreeNode) int {
    if root == nil { return 0 }
    left := maxDepth(root.Left)
    right := maxDepth(root.Right)
    if left > right { return left + 1 }
    return right + 1
}`,
	},
	{
		Slug:     "binary-tree-level-order-traversal",
		Category: "Trees",
		Options: []Option{
			{Text: "BFS with a queue, processing one level at a time by tracking queue size — O(n) time, O(n) space", Rating: Optimal},
			{Text: "DFS with a depth parameter, appending to the correct level's slice — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Repeated traversals for each level, filtering by depth — O(n * h) time, O(n) space", Rating: Suboptimal},
			{Text: "In-order traversal and group consecutive nodes — in-order doesn't produce level-ordered groups — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: BFS Level-by-Level
// Time: O(n) | Space: O(n)
func levelOrder(root *TreeNode) [][]int {
    if root == nil { return nil }
    var result [][]int
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        levelSize := len(queue)
        level := make([]int, 0, levelSize)
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            level = append(level, node.Val)
            if node.Left != nil { queue = append(queue, node.Left) }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
    }
    return result
}`,
	},

	// Graphs
	{
		Slug:     "number-of-islands",
		Category: "Graphs",
		Options: []Option{
			{Text: "DFS/BFS from each unvisited land cell, mark connected land as visited — O(m*n) time, O(m*n) space", Rating: Optimal},
			{Text: "Union-Find: union adjacent land cells, count distinct components — O(m*n * alpha(m*n)) time, O(m*n) space", Rating: Plausible},
			{Text: "For each land cell, run a full BFS to find its island, use a separate visited grid instead of modifying the input — O(m*n) time, O(m*n) space — correct but uses unnecessary extra storage", Rating: Suboptimal},
			{Text: "Count all land cells and divide by average island size — no way to determine average island size without visiting all — O(m*n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: DFS Grid Traversal
// Time: O(m * n) | Space: O(m * n) for recursion stack
func numIslands(grid [][]byte) int {
    rows, cols := len(grid), len(grid[0])
    count := 0
    var dfs func(r, c int)
    dfs = func(r, c int) {
        if r < 0 || r >= rows || c < 0 || c >= cols || grid[r][c] != '1' { return }
        grid[r][c] = '0'
        dfs(r+1, c); dfs(r-1, c); dfs(r, c+1); dfs(r, c-1)
    }
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if grid[r][c] == '1' { count++; dfs(r, c) }
        }
    }
    return count
}`,
	},
	{
		Slug:     "clone-graph",
		Category: "Graphs",
		Options: []Option{
			{Text: "BFS/DFS with a hash map mapping original node to its clone — O(V+E) time, O(V) space", Rating: Optimal},
			{Text: "Serialize the graph to a string then deserialize into new nodes — O(V+E) time, O(V+E) space", Rating: Plausible},
			{Text: "Create all cloned nodes first, then iterate neighbors to wire edges — still needs a map, essentially the same as optimal — O(V+E) time, O(V) space", Rating: Plausible},
			{Text: "Copy nodes by traversing without tracking visited — will loop infinitely on cycles — O(?) time", Rating: Wrong},
		},
		Solution: `// Pattern: BFS/DFS with Hash Map Clone Tracking
// Time: O(V + E) | Space: O(V)
func cloneGraph(node *Node) *Node {
    if node == nil { return nil }
    cloned := map[*Node]*Node{}
    var dfs func(*Node) *Node
    dfs = func(n *Node) *Node {
        if c, ok := cloned[n]; ok { return c }
        cp := &Node{Val: n.Val}
        cloned[n] = cp
        for _, neighbor := range n.Neighbors {
            cp.Neighbors = append(cp.Neighbors, dfs(neighbor))
        }
        return cp
    }
    return dfs(node)
}`,
	},
	{
		Slug:     "course-schedule",
		Category: "Graphs",
		Options: []Option{
			{Text: "Topological sort using DFS with three states (unvisited, visiting, visited) to detect cycles — O(V+E) time, O(V+E) space", Rating: Optimal},
			{Text: "Kahn's algorithm: BFS with in-degree tracking, if all nodes processed then no cycle — O(V+E) time, O(V+E) space", Rating: Optimal},
			{Text: "Try every possible ordering of courses and check if prerequisites are satisfied — O(n!) time, O(n) space", Rating: Suboptimal},
			{Text: "Just check if any course is its own prerequisite — misses longer cycles like A->B->C->A — O(E) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Topological Sort (DFS cycle detection)
// Time: O(V + E) | Space: O(V + E)
func canFinish(numCourses int, prerequisites [][]int) bool {
    graph := make([][]int, numCourses)
    for _, p := range prerequisites { graph[p[1]] = append(graph[p[1]], p[0]) }
    state := make([]int, numCourses) // 0=unvisited, 1=visiting, 2=visited
    var hasCycle func(int) bool
    hasCycle = func(node int) bool {
        if state[node] == 1 { return true }
        if state[node] == 2 { return false }
        state[node] = 1
        for _, next := range graph[node] { if hasCycle(next) { return true } }
        state[node] = 2
        return false
    }
    for i := 0; i < numCourses; i++ { if hasCycle(i) { return false } }
    return true
}`,
	},

	// Dynamic Programming
	{
		Slug:     "climbing-stairs",
		Category: "Dynamic Programming",
		Options: []Option{
			{Text: "Bottom-up DP using two variables (Fibonacci-like): dp[i] = dp[i-1] + dp[i-2] — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Top-down recursion with memoization — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Pure recursion without memoization — O(2^n) time, O(n) space", Rating: Suboptimal},
			{Text: "Use matrix exponentiation on the Fibonacci recurrence — O(log n) time, O(1) space", Rating: Plausible},
			{Text: "The answer is always 2^(n-1) — this is only true if you could take any size step — O(1) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bottom-Up DP (Fibonacci)
// Time: O(n) | Space: O(1)
func climbStairs(n int) int {
    if n <= 2 { return n }
    prev, curr := 1, 2
    for i := 3; i <= n; i++ { prev, curr = curr, prev+curr }
    return curr
}`,
	},
	{
		Slug:     "house-robber",
		Category: "Dynamic Programming",
		Options: []Option{
			{Text: "Bottom-up DP: at each house choose max(rob current + dp[i-2], skip and take dp[i-1]), track with two variables — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Top-down recursion with memoization: for each house decide rob or skip — O(n) time, O(n) space", Rating: Plausible},
			{Text: "Try every subset of non-adjacent houses — O(2^n) time, O(n) space", Rating: Suboptimal},
			{Text: "Greedy: always rob the richest house, skip its neighbors, repeat — greedy doesn't yield optimal for all inputs — O(n log n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Bottom-Up DP (Space Optimized)
// Time: O(n) | Space: O(1)
func rob(nums []int) int {
    prev, curr := 0, 0
    for _, num := range nums {
        prev, curr = curr, max(curr, prev+num)
    }
    return curr
}`,
	},
	{
		Slug:     "longest-increasing-subsequence",
		Category: "Dynamic Programming",
		Options: []Option{
			{Text: "Patience sorting: maintain tails array, use binary search to place each element — O(n log n) time, O(n) space", Rating: Optimal},
			{Text: "DP where dp[i] = length of LIS ending at index i, check all previous elements — O(n^2) time, O(n) space", Rating: Plausible},
			{Text: "Generate all subsequences and check which are increasing — O(2^n) time, O(n) space", Rating: Suboptimal},
			{Text: "Sort the array and find the longest common subsequence with the original — O(n^2) time, O(n^2) space", Rating: Plausible},
			{Text: "Greedily pick the next smallest increasing element — greedy misses better subsequences starting later — O(n) time, O(n) space", Rating: Wrong},
		},
		Solution: `// Pattern: Patience Sorting with Binary Search
// Time: O(n log n) | Space: O(n)
func lengthOfLIS(nums []int) int {
    tails := []int{}
    for _, num := range nums {
        lo, hi := 0, len(tails)
        for lo < hi {
            mid := lo + (hi-lo)/2
            if tails[mid] < num { lo = mid + 1 } else { hi = mid }
        }
        if lo == len(tails) { tails = append(tails, num) } else { tails[lo] = num }
    }
    return len(tails)
}`,
	},

	// Backtracking
	{
		Slug:     "subsets",
		Category: "Backtracking",
		Options: []Option{
			{Text: "Backtracking: at each index decide to include or exclude the element, recurse forward — O(n * 2^n) time, O(n) space", Rating: Optimal},
			{Text: "Iterative: start with [[]], for each number add it to copies of all existing subsets — O(n * 2^n) time, O(n * 2^n) space", Rating: Plausible},
			{Text: "Bit manipulation: iterate from 0 to 2^n - 1, each bitmask represents a subset — O(n * 2^n) time, O(n * 2^n) space", Rating: Plausible},
			{Text: "Generate all permutations then remove duplicates — generates way more than needed — O(n!) time, O(n!) space", Rating: Suboptimal},
			{Text: "Sort and use sliding window of different sizes — sliding window picks contiguous elements only, misses non-contiguous subsets — O(n^2) time, O(n^2) space", Rating: Wrong},
		},
		Solution: `// Pattern: Backtracking (Include/Exclude)
// Time: O(n * 2^n) | Space: O(n) recursion depth
func subsets(nums []int) [][]int {
    var result [][]int
    var backtrack func(start int, current []int)
    backtrack = func(start int, current []int) {
        tmp := make([]int, len(current))
        copy(tmp, current)
        result = append(result, tmp)
        for i := start; i < len(nums); i++ {
            backtrack(i+1, append(current, nums[i]))
        }
    }
    backtrack(0, []int{})
    return result
}`,
	},
	{
		Slug:     "combination-sum",
		Category: "Backtracking",
		Options: []Option{
			{Text: "Backtracking with start index to avoid duplicates, reuse candidates by not advancing index — O(n^(T/M)) time, O(T/M) space where T=target, M=min candidate", Rating: Optimal},
			{Text: "DP building up all combinations for each sum from 1 to target — O(n * T^2) time, O(T^2) space", Rating: Plausible},
			{Text: "Generate all possible combinations up to target/min length, filter by sum — O(n^(T/M)) time, exponential space", Rating: Suboptimal},
			{Text: "Sort candidates and use two pointers to find pairs summing to target — only finds pairs, not arbitrary-length combinations — O(n log n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Backtracking with Reuse
// Time: O(n^(T/M)) | Space: O(T/M)
func combinationSum(candidates []int, target int) [][]int {
    var result [][]int
    var backtrack func(start, remaining int, combo []int)
    backtrack = func(start, remaining int, combo []int) {
        if remaining == 0 {
            tmp := make([]int, len(combo))
            copy(tmp, combo)
            result = append(result, tmp)
            return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] > remaining { continue }
            backtrack(i, remaining-candidates[i], append(combo, candidates[i]))
        }
    }
    backtrack(0, target, []int{})
    return result
}`,
	},
	{
		Slug:     "word-search",
		Category: "Backtracking",
		Options: []Option{
			{Text: "Backtracking DFS from each cell matching the first letter, mark visited cells and unmark on backtrack — O(m*n*4^L) time, O(L) space", Rating: Optimal},
			{Text: "BFS from each matching start cell with per-path visited state — correct but harder to implement than DFS and uses more memory for visited tracking — O(m*n*4^L) time, O(m*n*L) space", Rating: Plausible},
			{Text: "DFS from each cell without backtracking (permanently marking visited) — may block valid paths that reuse cells from a different starting direction — O(m*n) time, O(m*n) space", Rating: Suboptimal},
			{Text: "Check if the board contains all characters in the word with correct frequencies — character existence doesn't imply a valid connected path — O(m*n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Backtracking DFS on Grid
// Time: O(m * n * 4^L) | Space: O(L) where L = word length
func exist(board [][]byte, word string) bool {
    rows, cols := len(board), len(board[0])
    var dfs func(r, c, idx int) bool
    dfs = func(r, c, idx int) bool {
        if idx == len(word) { return true }
        if r < 0 || r >= rows || c < 0 || c >= cols || board[r][c] != word[idx] { return false }
        temp := board[r][c]
        board[r][c] = '#'
        found := dfs(r+1, c, idx+1) || dfs(r-1, c, idx+1) || dfs(r, c+1, idx+1) || dfs(r, c-1, idx+1)
        board[r][c] = temp
        return found
    }
    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ { if dfs(r, c, 0) { return true } }
    }
    return false
}`,
	},

	// Greedy
	{
		Slug:     "jump-game",
		Category: "Greedy",
		Options: []Option{
			{Text: "Greedy: track the farthest reachable index, iterate and update — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Greedy from the end: start at the last index, scan backward to see if any earlier index can reach it, shift target — O(n) time, O(1) space", Rating: Plausible},
			{Text: "DP: dp[i] = whether index i is reachable, check all previous indices — O(n^2) time, O(n) space", Rating: Plausible},
			{Text: "BFS treating each index as a node with edges to reachable indices — O(n^2) time, O(n) space", Rating: Suboptimal},
			{Text: "Check if the array contains any zeros — arrays can have zeros and still be solvable if you can jump over them — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Greedy (Farthest Reach)
// Time: O(n) | Space: O(1)
func canJump(nums []int) bool {
    farthest := 0
    for i := 0; i < len(nums); i++ {
        if i > farthest { return false }
        if i+nums[i] > farthest { farthest = i + nums[i] }
    }
    return true
}`,
	},
	{
		Slug:     "maximum-subarray",
		Category: "Greedy",
		Options: []Option{
			{Text: "Kadane's algorithm: track current sum, reset to current element when running sum is negative, track global max — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Divide and conquer: split array, find max in left, right, and crossing subarrays — O(n log n) time, O(log n) space", Rating: Plausible},
			{Text: "Check all possible subarrays — O(n^2) time, O(1) space", Rating: Suboptimal},
			{Text: "Sort the array and sum the largest elements — sorting destroys contiguity — O(n log n) time, O(1) space", Rating: Wrong},
			{Text: "Prefix sums: find max prefix sum minus min prefix sum seen before it — O(n) time, O(1) space", Rating: Plausible},
		},
		Solution: `// Pattern: Kadane's Algorithm
// Time: O(n) | Space: O(1)
func maxSubArray(nums []int) int {
    maxSum, curSum := nums[0], nums[0]
    for _, num := range nums[1:] {
        if curSum < 0 { curSum = num } else { curSum += num }
        if curSum > maxSum { maxSum = curSum }
    }
    return maxSum
}`,
	},
	{
		Slug:     "task-scheduler",
		Category: "Greedy",
		Options: []Option{
			{Text: "Greedy math: compute (maxFreq - 1) * (n + 1) + countOfMaxFreq, take max with total tasks — O(n) time, O(1) space", Rating: Optimal},
			{Text: "Use a max-heap to always schedule the most frequent available task, with a cooldown queue — O(n * m) time, O(m) space", Rating: Plausible},
			{Text: "Try all possible orderings and find the shortest valid one — O(m!) time where m = number of tasks", Rating: Suboptimal},
			{Text: "Just sum up all tasks — ignores idle time forced by the cooldown constraint — O(n) time, O(1) space", Rating: Wrong},
		},
		Solution: `// Pattern: Greedy Math Formula
// Time: O(n) | Space: O(1) — 26 letters at most
func leastInterval(tasks []byte, n int) int {
    freq := [26]int{}
    for _, t := range tasks { freq[t-'A']++ }
    maxFreq := 0
    for _, f := range freq { if f > maxFreq { maxFreq = f } }
    maxCount := 0
    for _, f := range freq { if f == maxFreq { maxCount++ } }
    result := (maxFreq-1)*(n+1) + maxCount
    if len(tasks) > result { return len(tasks) }
    return result
}`,
	},

	// Design
	{
		Slug:     "lru-cache",
		Category: "Design",
		Options: []Option{
			{Text: "Hash map + doubly linked list: map gives O(1) lookup, list maintains access order for O(1) eviction — O(1) per operation, O(capacity) space", Rating: Optimal},
			{Text: "Use an ordered map (e.g., LinkedHashMap) — if the language supports it — Go doesn't have one built in — O(1) per operation, O(capacity) space", Rating: Plausible},
			{Text: "Use a min-heap ordered by access time — O(log n) per operation instead of O(1) — O(capacity) space", Rating: Plausible},
			{Text: "Hash map with timestamps, scan for oldest on eviction — O(n) eviction, O(capacity) space", Rating: Suboptimal},
			{Text: "Use a hash map and evict a random key when full — doesn't track recency, so frequently used keys get evicted — O(1) per operation, O(capacity) space", Rating: Wrong},
		},
		Solution: `// Pattern: Hash Map + Doubly Linked List
// Time: O(1) per get/put | Space: O(capacity)
type LRUCache struct {
    capacity   int
    cache      map[int]*node
    head, tail *node
}
type node struct {
    key, val   int
    prev, next *node
}

func Constructor(capacity int) LRUCache {
    h, t := &node{}, &node{}
    h.next, t.prev = t, h
    return LRUCache{capacity, make(map[int]*node), h, t}
}

func (c *LRUCache) Get(key int) int {
    if n, ok := c.cache[key]; ok { c.remove(n); c.insertFront(n); return n.val }
    return -1
}
func (c *LRUCache) Put(key, value int) {
    if n, ok := c.cache[key]; ok { n.val = value; c.remove(n); c.insertFront(n); return }
    if len(c.cache) >= c.capacity { lru := c.tail.prev; c.remove(lru); delete(c.cache, lru.key) }
    n := &node{key: key, val: value}; c.cache[key] = n; c.insertFront(n)
}
func (c *LRUCache) remove(n *node)      { n.prev.next = n.next; n.next.prev = n.prev }
func (c *LRUCache) insertFront(n *node)  { n.next = c.head.next; n.prev = c.head; c.head.next.prev = n; c.head.next = n }`,
	},
	{
		Slug:     "implement-trie-prefix-tree",
		Category: "Design",
		Options: []Option{
			{Text: "Trie with nodes containing a children array and an end-of-word flag — O(L) per operation, O(N*L) space where N=words, L=avg length", Rating: Optimal},
			{Text: "Store words in a hash set, iterate all for prefix check — O(1) search, O(N*L) startsWith — O(N*L) space", Rating: Suboptimal},
			{Text: "Store words in a sorted array, use binary search for prefix — O(log N * L) per operation, O(N*L) space", Rating: Plausible},
			{Text: "Use a single string concatenating all words — search and prefix operations become O(total chars), impractical — O(N*L) space", Rating: Wrong},
		},
		Solution: `// Pattern: Trie (Prefix Tree)
// Time: O(L) per operation | Space: O(N * L)
type Trie struct {
    children [26]*Trie
    isEnd    bool
}

func (t *Trie) Insert(word string) {
    curr := t
    for _, ch := range word {
        idx := ch - 'a'
        if curr.children[idx] == nil { curr.children[idx] = &Trie{} }
        curr = curr.children[idx]
    }
    curr.isEnd = true
}
func (t *Trie) Search(word string) bool { n := t.find(word); return n != nil && n.isEnd }
func (t *Trie) StartsWith(prefix string) bool { return t.find(prefix) != nil }
func (t *Trie) find(s string) *Trie {
    curr := t
    for _, ch := range s {
        idx := ch - 'a'
        if curr.children[idx] == nil { return nil }
        curr = curr.children[idx]
    }
    return curr
}`,
	},
	{
		Slug:     "design-twitter",
		Category: "Design",
		Options: []Option{
			{Text: "Hash maps for followers and tweet lists, merge k sorted lists with a min-heap for feed — O(k log k) for feed where k = followees, O(N) space", Rating: Optimal},
			{Text: "Store all tweets in one list, filter by user's follow set on feed request — O(T) per feed where T = total tweets, O(T) space", Rating: Suboptimal},
			{Text: "Pre-compute feeds on every post (fan-out on write) — O(F) per post where F = followers, O(N*F) space", Rating: Plausible},
			{Text: "Store tweets in a BST sorted by time — BST doesn't help with filtering by followee set — O(T) per feed, O(T) space", Rating: Wrong},
		},
		Solution: `// Pattern: Hash Map + Merge K Sorted Lists
// Time: O(k log k) per getNewsFeed | Space: O(users + tweets)
type Twitter struct {
    time    int
    tweets  map[int][]tweet
    follows map[int]map[int]bool
}
type tweet struct { id, time int }

func (t *Twitter) PostTweet(userId, tweetId int) {
    t.time++
    t.tweets[userId] = append(t.tweets[userId], tweet{tweetId, t.time})
}
func (t *Twitter) GetNewsFeed(userId int) []int {
    var all []tweet
    users := []int{userId}
    for uid := range t.follows[userId] { users = append(users, uid) }
    for _, uid := range users {
        tw := t.tweets[uid]
        start := len(tw) - 10; if start < 0 { start = 0 }
        all = append(all, tw[start:]...)
    }
    sort.Slice(all, func(i, j int) bool { return all[i].time > all[j].time })
    result := make([]int, 0, 10)
    for i := 0; i < len(all) && i < 10; i++ { result = append(result, all[i].id) }
    return result
}
func (t *Twitter) Follow(er, ee int) {
    if t.follows[er] == nil { t.follows[er] = make(map[int]bool) }
    t.follows[er][ee] = true
}
func (t *Twitter) Unfollow(er, ee int) { delete(t.follows[er], ee) }`,
	},
}
