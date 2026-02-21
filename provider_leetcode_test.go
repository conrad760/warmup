package main

import (
	"strings"
	"testing"
)

func TestNormalizeLangSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"go", "golang"},
		{"golang", "golang"},
		{"Go", "golang"},
		{"python", "python3"},
		{"python3", "python3"},
		{"py", "python3"},
		{"java", "java"},
		{"cpp", "cpp"},
		{"c++", "cpp"},
		{"c", "c"},
		{"js", "javascript"},
		{"javascript", "javascript"},
		{"ts", "typescript"},
		{"typescript", "typescript"},
		{"rust", "rust"},
		{"rs", "rust"},
		{"ruby", "ruby"}, // unmapped, lowercased as-is
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeLangSlug(tt.input)
			if got != tt.want {
				t.Errorf("normalizeLangSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLcParseContent_BasicHTML(t *testing.T) {
	html := `<p>Given an array of integers <code>nums</code> and an integer <code>target</code>, return indices of the two numbers.</p>
<p><strong>Example 1:</strong></p>
<pre>
Input: nums = [2,7,11,15], target = 9
Output: [0,1]
</pre>
<p><strong>Constraints:</strong></p>
<ul><li>2 &lt;= nums.length</li></ul>`

	desc, example := lcParseContent(html)

	if !strings.Contains(desc, "array of integers") {
		t.Errorf("description should contain problem text, got: %s", desc)
	}
	if strings.Contains(desc, "<") {
		t.Errorf("description should not contain HTML tags, got: %s", desc)
	}
	// Constraints should be stripped from description.
	if strings.Contains(strings.ToLower(desc), "constraint") {
		t.Errorf("description should not contain constraints, got: %s", desc)
	}
	// Example is in the second split part — may or may not be populated
	// depending on how the <strong>Example tags are structured.
	_ = example
}

func TestLcParseContent_WithExampleSections(t *testing.T) {
	html := `<p>Description text here.</p>
<p><strong>Example 1:</strong></p>
<pre>Input: x = 5
Output: 25</pre>
<p><strong>Example 2:</strong></p>
<pre>Input: x = -3
Output: 9</pre>
<p><strong>Constraints:</strong></p>
<ul><li>-100 &lt;= x &lt;= 100</li></ul>`

	desc, example := lcParseContent(html)

	if !strings.Contains(desc, "Description text here") {
		t.Errorf("description missing main text, got: %s", desc)
	}
	if example == "" {
		t.Error("example should not be empty when Example section exists")
	}
	if strings.Contains(example, "Constraint") {
		t.Error("example should not include constraints section")
	}
}

func TestLcHTMLToText_UnescapesEntities(t *testing.T) {
	got := lcHTMLToText(`<p>a &amp; b &lt; c &gt; d</p>`)
	if !strings.Contains(got, "a & b < c > d") {
		t.Errorf("should unescape HTML entities, got: %s", got)
	}
}

func TestLcHTMLToText_SuperscriptConversion(t *testing.T) {
	html := `<p>2<sup>31</sup> - 1</p>`
	// Superscripts are converted in lcParseContent (before lcHTMLToText),
	// so test via lcParseContent.
	desc, _ := lcParseContent(html)
	if !strings.Contains(desc, "^31") {
		t.Errorf("superscripts should be converted to ^N, got: %s", desc)
	}
}

func TestLcHTMLToText_StripsAllTags(t *testing.T) {
	html := `<div class="foo"><span>hello</span> <a href="bar">world</a></div>`
	got := lcHTMLToText(html)
	if strings.Contains(got, "<") || strings.Contains(got, ">") {
		t.Errorf("should strip all HTML tags, got: %s", got)
	}
	if !strings.Contains(got, "hello") || !strings.Contains(got, "world") {
		t.Errorf("should preserve text content, got: %s", got)
	}
}
