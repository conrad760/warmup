package main

import (
	"fmt"
	"sort"
)

// CurriculumDay describes one day's problem set from the 21-day curriculum.
type CurriculumDay struct {
	Day      int
	Topic    string
	Problems []string // LeetCode problem slugs
}

// curriculumDays maps day number → problem set.  Only topic days (1-18) are
// included; days 19-21 are mock practice and review with no fixed problem set.
// Each slug list contains every LeetCode problem explicitly referenced in the
// corresponding curriculum/day-NN-*.md file that also exists in the curated
// question bank.
var curriculumDays = map[int]CurriculumDay{
	1: {
		Day:   1,
		Topic: "Arrays & Hashing",
		Problems: []string{
			"two-sum",
			"group-anagrams",
			"valid-anagram",
			"top-k-frequent-elements",
			"contains-duplicate",
		},
	},
	2: {
		Day:   2,
		Topic: "Two Pointers & Sliding Window",
		Problems: []string{
			"3sum",
			"minimum-window-substring",
			"longest-substring-without-repeating-characters",
			"container-with-most-water",
			"valid-palindrome",
			"two-sum-ii-input-array-is-sorted",
		},
	},
	3: {
		Day:   3,
		Topic: "Stacks",
		Problems: []string{
			"valid-parentheses",
			"daily-temperatures",
			"evaluate-reverse-polish-notation",
			"min-stack",
			"largest-rectangle-in-histogram",
		},
	},
	4: {
		Day:   4,
		Topic: "Binary Search",
		Problems: []string{
			"binary-search",
			"search-in-rotated-sorted-array",
			"find-minimum-in-rotated-sorted-array",
			"koko-eating-bananas",
		},
	},
	5: {
		Day:   5,
		Topic: "Linked Lists",
		Problems: []string{
			"reverse-linked-list",
			"linked-list-cycle",
			"merge-two-sorted-lists",
			"remove-nth-node-from-end-of-list",
			"reorder-list",
		},
	},
	6: {
		Day:   6,
		Topic: "Binary Trees",
		Problems: []string{
			"maximum-depth-of-binary-tree",
			"same-tree",
			"binary-tree-level-order-traversal",
			"construct-binary-tree-from-preorder-and-inorder-traversal",
			"serialize-and-deserialize-binary-tree",
		},
	},
	7: {
		Day:   7,
		Topic: "BST & Tree Construction",
		Problems: []string{
			"validate-binary-search-tree",
			"kth-smallest-element-in-a-bst",
			"lowest-common-ancestor-of-a-binary-search-tree",
			"binary-tree-maximum-path-sum",
			"construct-binary-tree-from-preorder-and-inorder-traversal",
		},
	},
	8: {
		Day:   8,
		Topic: "Graphs: BFS & DFS",
		Problems: []string{
			"number-of-islands",
			"clone-graph",
			"rotting-oranges",
		},
	},
	9: {
		Day:   9,
		Topic: "Topological Sort",
		Problems: []string{
			"course-schedule",
			"course-schedule-ii",
			"alien-dictionary",
		},
	},
	10: {
		Day:   10,
		Topic: "Heaps",
		Problems: []string{
			"merge-k-sorted-lists",
			"find-median-from-data-stream",
			"top-k-frequent-elements",
			"task-scheduler",
		},
	},
	11: {
		Day:   11,
		Topic: "DP: 1D",
		Problems: []string{
			"climbing-stairs",
			"coin-change",
			"decode-ways",
			"house-robber",
			"house-robber-ii",
			"longest-increasing-subsequence",
			"maximum-subarray",
			"word-break",
		},
	},
	12: {
		Day:   12,
		Topic: "DP: 2D",
		Problems: []string{
			"unique-paths",
			"longest-common-subsequence",
		},
	},
	13: {
		Day:   13,
		Topic: "Backtracking",
		Problems: []string{
			"subsets",
			"combination-sum",
			"word-search",
		},
	},
	14: {
		Day:   14,
		Topic: "Greedy",
		Problems: []string{
			"maximum-subarray",
			"best-time-to-buy-and-sell-stock",
			"jump-game",
			"jump-game-ii",
			"non-overlapping-intervals",
			"task-scheduler",
		},
	},
	15: {
		Day:   15,
		Topic: "Intervals",
		Problems: []string{
			"merge-intervals",
			"insert-interval",
			"meeting-rooms",
			"meeting-rooms-ii",
			"non-overlapping-intervals",
		},
	},
	16: {
		Day:   16,
		Topic: "Tries & Union Find",
		Problems: []string{
			"implement-trie-prefix-tree",
			"number-of-islands",
			"graph-valid-tree",
			"number-of-connected-components-in-an-undirected-graph",
		},
	},
	17: {
		Day:   17,
		Topic: "Design",
		Problems: []string{
			"lru-cache",
			"min-stack",
			"find-median-from-data-stream",
			"design-twitter",
		},
	},
	18: {
		Day:   18,
		Topic: "Bits & Math",
		Problems: []string{
			"number-of-1-bits",
			"reverse-bits",
			"missing-number",
			"counting-bits",
			"sum-of-two-integers",
		},
	},
}

// CurriculumProblems returns the problem slugs for the given day.
func CurriculumProblems(day int) (CurriculumDay, error) {
	cd, ok := curriculumDays[day]
	if !ok {
		return CurriculumDay{}, fmt.Errorf("no curriculum entry for day %d (valid: 1-18)", day)
	}
	return cd, nil
}

// ListCurriculumDays returns all curriculum days sorted by day number.
func ListCurriculumDays() []CurriculumDay {
	days := make([]CurriculumDay, 0, len(curriculumDays))
	for _, cd := range curriculumDays {
		days = append(days, cd)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Day < days[j].Day
	})
	return days
}
