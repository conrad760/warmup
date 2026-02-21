package main

import "fmt"

// MockProvider is a provider with hardcoded problems for testing without network access.
// Useful for development, CI, and demonstrating the provider interface.
type MockProvider struct{}

func init() {
	RegisterProvider("mock", func() Provider { return &MockProvider{} })
}

func (m *MockProvider) Name() string { return "mock" }

func (m *MockProvider) FetchProblem(id string, lang string) (*ProblemData, error) {
	p, ok := mockProblems[id]
	if !ok {
		return nil, fmt.Errorf("mock provider: problem %q not found", id)
	}

	// Return a copy so callers can't mutate the source.
	result := p
	// Deep copy Meta to prevent shared state.
	if p.Meta != nil {
		metaCopy := *p.Meta
		paramsCopy := make([]ParamMeta, len(p.Meta.Params))
		copy(paramsCopy, p.Meta.Params)
		metaCopy.Params = paramsCopy
		if p.Meta.Return != nil {
			retCopy := *p.Meta.Return
			metaCopy.Return = &retCopy
		}
		result.Meta = &metaCopy
	}
	return &result, nil
}

// mockProblems holds a small set of hardcoded problems for testing.
var mockProblems = map[string]ProblemData{
	"fizz-buzz": {
		ID:    "fizz-buzz",
		Title: "Fizz Buzz",
		Description: `Given an integer n, return a string array answer (1-indexed) where:
answer[i] == "FizzBuzz" if i is divisible by 3 and 5.
answer[i] == "Fizz" if i is divisible by 3.
answer[i] == "Buzz" if i is divisible by 5.
answer[i] == i (as a string) if none of the above conditions are true.`,
		Examples: `Input: n = 3
Output: ["1","2","Fizz"]`,
		Difficulty:  "Easy",
		Tags:        []string{"Math", "String", "Simulation"},
		CodeSnippet: "func fizzBuzz(n int) []string {\n    \n}",
		TestInput:   "3",
		Meta: &FuncMeta{
			Name:   "fizzBuzz",
			Params: []ParamMeta{{Name: "n", Type: "integer"}},
			Return: &ParamMeta{Type: "string[]"},
		},
	},
	"reverse-string": {
		ID:    "reverse-string",
		Title: "Reverse String",
		Description: `Write a function that reverses a string.
The input string is given as an array of characters s.
You must do this by modifying the input array in-place with O(1) extra memory.`,
		Examples: `Input: s = ["h","e","l","l","o"]
Output: ["o","l","l","e","h"]`,
		Difficulty:  "Easy",
		Tags:        []string{"Two Pointers", "String"},
		CodeSnippet: "func reverseString(s []byte) {\n    \n}",
		TestInput:   `["h","e","l","l","o"]`,
		Meta: &FuncMeta{
			Name:   "reverseString",
			Params: []ParamMeta{{Name: "s", Type: "character[]"}},
			Return: nil, // void — in-place modification
		},
	},
	"valid-palindrome": {
		ID:    "valid-palindrome",
		Title: "Valid Palindrome",
		Description: `A phrase is a palindrome if, after converting all uppercase letters
into lowercase letters and removing all non-alphanumeric characters, it reads
the same forward and backward. Alphanumeric characters include letters and numbers.
Given a string s, return true if it is a palindrome, or false otherwise.`,
		Examples: `Input: s = "A man, a plan, a canal: Panama"
Output: true
Explanation: "amanaplanacanalpanama" is a palindrome.`,
		Difficulty:  "Easy",
		Tags:        []string{"Two Pointers", "String"},
		CodeSnippet: "func isPalindrome(s string) bool {\n    \n}",
		TestInput:   `"A man, a plan, a canal: Panama"`,
		Meta: &FuncMeta{
			Name:   "isPalindrome",
			Params: []ParamMeta{{Name: "s", Type: "string"}},
			Return: &ParamMeta{Type: "boolean"},
		},
	},
	"binary-search": {
		ID:    "binary-search",
		Title: "Binary Search",
		Description: `Given an array of integers nums which is sorted in ascending order,
and an integer target, write a function to search target in nums. If target
exists, then return its index. Otherwise, return -1.
You must write an algorithm with O(log n) runtime complexity.`,
		Examples: `Input: nums = [-1,0,3,5,9,12], target = 9
Output: 4
Explanation: 9 exists in nums and its index is 4.`,
		Difficulty:  "Easy",
		Tags:        []string{"Array", "Binary Search"},
		CodeSnippet: "func search(nums []int, target int) int {\n    \n}",
		TestInput:   "[-1,0,3,5,9,12]\n9",
		Meta: &FuncMeta{
			Name:   "search",
			Params: []ParamMeta{{Name: "nums", Type: "integer[]"}, {Name: "target", Type: "integer"}},
			Return: &ParamMeta{Type: "integer"},
		},
	},
	"merge-sort": {
		ID:    "merge-sort",
		Title: "Merge Sort Implementation",
		Description: `Implement the merge sort algorithm. Given an array of integers,
sort them in ascending order using the divide-and-conquer merge sort approach.`,
		Examples: `Input: nums = [5,2,3,1]
Output: [1,2,3,5]`,
		Difficulty:  "Medium",
		Tags:        []string{"Array", "Sorting", "Divide and Conquer"},
		CodeSnippet: "func sortArray(nums []int) []int {\n    \n}",
		TestInput:   "[5,2,3,1]",
		Meta: &FuncMeta{
			Name:   "sortArray",
			Params: []ParamMeta{{Name: "nums", Type: "integer[]"}},
			Return: &ParamMeta{Type: "integer[]"},
		},
	},
}
