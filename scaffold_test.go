package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testProblem() ScaffoldProblem {
	return ScaffoldProblem{
		Provider:    "leetcode",
		ProblemID:   "two-sum",
		Title:       "Two Sum",
		Difficulty:  "Easy",
		Description: "Given an array of integers nums and an integer target, return indices of the two numbers.",
		Example:     "Input: nums = [2,7,11,15], target = 9\nOutput: [0,1]",
		CodeSnippet: "func twoSum(nums []int, target int) []int {\n    \n}",
		TestInput:   "[2,7,11,15]\n9",
		Meta: &FuncMeta{
			Name:   "twoSum",
			Params: []ParamMeta{{Name: "nums", Type: "integer[]"}, {Name: "target", Type: "integer"}},
			Return: &ParamMeta{Type: "integer[]"},
		},
	}
}

func TestScaffold_EnsureScaffold_CreatesAllFiles(t *testing.T) {
	dir := t.TempDir()
	s, err := NewScaffold(dir, "go")
	if err != nil {
		t.Fatalf("NewScaffold: %v", err)
	}

	p := testProblem()
	path, err := s.EnsureScaffold(p)
	if err != nil {
		t.Fatalf("EnsureScaffold: %v", err)
	}

	// solution.go should exist and contain snippet + markers + func main.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("solution file should exist")
	}
	solContent, _ := os.ReadFile(path)
	sol := string(solContent)
	if !strings.Contains(sol, "func twoSum") {
		t.Error("solution file should contain the code snippet")
	}
	if !strings.Contains(sol, codeBeginMarker) {
		t.Error("solution file should contain code begin marker")
	}
	if !strings.Contains(sol, codeEndMarker) {
		t.Error("solution file should contain code end marker")
	}
	if !strings.Contains(sol, "package main") {
		t.Error("Go solution file should have package declaration")
	}
	if !strings.Contains(sol, "func main()") {
		t.Error("Go solution file should have func main()")
	}
	if !strings.Contains(sol, "func readLine") {
		t.Error("Go solution file should have readLine helper")
	}

	// question.md should exist with title and description.
	qPath := s.QuestionPath(p.Provider, p.ProblemID)
	if _, err := os.Stat(qPath); os.IsNotExist(err) {
		t.Fatal("question.md should exist")
	}
	qContent, _ := os.ReadFile(qPath)
	q := string(qContent)
	if !strings.Contains(q, "# Two Sum (Easy)") {
		t.Errorf("question.md should contain title header, got: %s", q)
	}
	if !strings.Contains(q, "leetcode.com/problems/two-sum") {
		t.Error("question.md should contain LeetCode link")
	}

	// testcases.txt should exist with input.
	tcPath := s.TestCasesPath(p.Provider, p.ProblemID)
	if _, err := os.Stat(tcPath); os.IsNotExist(err) {
		t.Fatal("testcases.txt should exist")
	}
	tcContent, _ := os.ReadFile(tcPath)
	tc := string(tcContent)
	if !strings.Contains(tc, "input:") {
		t.Error("testcases.txt should contain input: marker")
	}
	if !strings.Contains(tc, "[2,7,11,15]") {
		t.Error("testcases.txt should contain test input data")
	}
}

func TestScaffold_EnsureScaffold_PreservesExisting(t *testing.T) {
	dir := t.TempDir()
	s, err := NewScaffold(dir, "go")
	if err != nil {
		t.Fatalf("NewScaffold: %v", err)
	}

	p := testProblem()
	_, err = s.EnsureScaffold(p)
	if err != nil {
		t.Fatal(err)
	}

	solPath := s.SolutionPath(p.Provider, p.ProblemID)
	os.WriteFile(solPath, []byte("// user's custom solution\nfunc mySolution() {}"), 0o644)

	p.CodeSnippet = "func different() {}"
	_, err = s.EnsureScaffold(p)
	if err != nil {
		t.Fatal(err)
	}

	content, _ := os.ReadFile(solPath)
	if !strings.Contains(string(content), "mySolution") {
		t.Error("EnsureScaffold should not overwrite existing solution file")
	}
}

func TestScaffold_SolutionPath_Go(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewScaffold(dir, "go")
	expected := filepath.Join(dir, "leetcode", "two-sum", "solution.go")
	if got := s.SolutionPath("leetcode", "two-sum"); got != expected {
		t.Errorf("SolutionPath = %q, want %q", got, expected)
	}
}

func TestScaffold_SolutionPath_Python(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewScaffold(dir, "python")
	expected := filepath.Join(dir, "leetcode", "two-sum", "solution.py")
	if got := s.SolutionPath("leetcode", "two-sum"); got != expected {
		t.Errorf("SolutionPath = %q, want %q", got, expected)
	}
}

// --- buildSolutionFile tests ---

func TestBuildSolutionFile_Go_HasHarness(t *testing.T) {
	p := testProblem()
	content := buildSolutionFile(p, lookupLang("go"), "golang")

	checks := []struct {
		desc string
		want string
	}{
		{"package", "package main"},
		{"imports", "encoding/json"},
		{"begin marker", codeBeginMarker},
		{"end marker", codeEndMarker},
		{"snippet", "func twoSum"},
		{"func main", "func main()"},
		{"readLine helper", "func readLine"},
		{"stdin reader", "bufio.NewReader(os.Stdin)"},
		{"deserialize nums", "json.Unmarshal([]byte(readLine(stdin)), &nums)"},
		{"deserialize target", "json.Unmarshal([]byte(readLine(stdin)), &target)"},
		{"call twoSum", "twoSum(nums, target)"},
		{"serialize output", "json.Marshal(ans)"},
		{"output prefix", `output:`},
	}

	for _, c := range checks {
		if !strings.Contains(content, c.want) {
			t.Errorf("Go solution should contain %s (%q), got:\n%s", c.desc, c.want, content)
		}
	}

	// Verify ordering: begin marker -> snippet -> end marker -> func main
	beginIdx := strings.Index(content, codeBeginMarker)
	snippetIdx := strings.Index(content, "func twoSum")
	endIdx := strings.Index(content, codeEndMarker)
	mainIdx := strings.Index(content, "func main()")
	if !(beginIdx < snippetIdx && snippetIdx < endIdx && endIdx < mainIdx) {
		t.Error("code markers should surround snippet, func main should come after")
	}
}

func TestBuildSolutionFile_Go_VoidFunction(t *testing.T) {
	p := ScaffoldProblem{
		CodeSnippet: "func reverseString(s []byte) {\n    \n}",
		Meta: &FuncMeta{
			Name:   "reverseString",
			Params: []ParamMeta{{Name: "s", Type: "character[]"}},
			Return: nil, // void
		},
	}
	content := buildSolutionFile(p, lookupLang("go"), "golang")

	// Should call function without capturing return.
	if !strings.Contains(content, "reverseString(s)") {
		t.Error("should call void function without return capture")
	}
	// Should serialize the first param for output.
	if !strings.Contains(content, "json.Marshal(s)") {
		t.Error("void function should serialize first param as output")
	}
}

func TestBuildSolutionFile_Go_NoMeta(t *testing.T) {
	p := ScaffoldProblem{
		CodeSnippet: "func mystery() {\n}",
		Meta:        nil,
	}
	content := buildSolutionFile(p, lookupLang("go"), "golang")

	if !strings.Contains(content, "func main()") {
		t.Error("should still have func main() even without metadata")
	}
	if !strings.Contains(content, "TODO") {
		t.Error("should have TODO comment when harness can't be generated")
	}
}

func TestBuildSolutionFile_Go_Empty(t *testing.T) {
	p := ScaffoldProblem{}
	content := buildSolutionFile(p, lookupLang("go"), "golang")
	if !strings.Contains(content, "func main()") {
		t.Error("empty Go file should still have func main()")
	}
}

func TestBuildSolutionFile_Python_NoHarness(t *testing.T) {
	p := testProblem()
	content := buildSolutionFile(p, lookupLang("python"), "python3")

	if strings.Contains(content, "func main") {
		t.Error("Python file should not contain Go boilerplate")
	}
	if !strings.Contains(content, "# "+codeBeginMarker) {
		t.Error("Python markers should use # comment style")
	}
	if !strings.Contains(content, "func twoSum") {
		t.Error("should include the code snippet")
	}
}

func TestBuildSolutionFile_Cpp_HasIncludes(t *testing.T) {
	p := testProblem()
	content := buildSolutionFile(p, lookupLang("cpp"), "cpp")
	if !strings.Contains(content, "#include") {
		t.Error("C++ file should include common headers")
	}
	if !strings.Contains(content, codeBeginMarker) {
		t.Error("should include code markers")
	}
}

// --- buildGoMainFunc tests ---

func TestBuildGoMainFunc_NilMeta(t *testing.T) {
	content := buildGoMainFunc(nil)
	if !strings.Contains(content, "func main()") {
		t.Error("nil meta should produce fallback main")
	}
	if !strings.Contains(content, "TODO") {
		t.Error("nil meta should produce TODO comment")
	}
}

func TestBuildGoMainFunc_SystemDesign(t *testing.T) {
	meta := &FuncMeta{Name: "Constructor", SystemDesign: true}
	content := buildGoMainFunc(meta)
	if !strings.Contains(content, "TODO") {
		t.Error("system design should produce fallback main with TODO")
	}
}

// --- lcTypeToGo tests ---

func TestLcTypeToGo(t *testing.T) {
	tests := []struct {
		lcType string
		want   string
	}{
		{"integer", "int"},
		{"integer[]", "[]int"},
		{"integer[][]", "[][]int"},
		{"string", "string"},
		{"string[]", "[]string"},
		{"boolean", "bool"},
		{"double", "float64"},
		{"character[]", "[]byte"},
		{"TreeNode", "*TreeNode"},
		{"ListNode", "*ListNode"},
		{"unknown_type", "interface{}"},
	}
	for _, tt := range tests {
		t.Run(tt.lcType, func(t *testing.T) {
			got := lcTypeToGo(tt.lcType)
			if got != tt.want {
				t.Errorf("lcTypeToGo(%q) = %q, want %q", tt.lcType, got, tt.want)
			}
		})
	}
}

// --- buildQuestionFile tests ---

func TestBuildQuestionFile_LeetCode(t *testing.T) {
	p := testProblem()
	content := buildQuestionFile(p)
	if !strings.Contains(content, "# Two Sum (Easy)") {
		t.Error("should contain title with difficulty")
	}
	if !strings.Contains(content, "leetcode.com/problems/two-sum") {
		t.Error("should contain LeetCode link")
	}
}

func TestBuildQuestionFile_NonLeetCode(t *testing.T) {
	p := testProblem()
	p.Provider = "codewars"
	content := buildQuestionFile(p)
	if strings.Contains(content, "leetcode.com") {
		t.Error("non-leetcode provider should not have LeetCode link")
	}
}

// --- buildTestCasesFile tests ---

func TestBuildTestCasesFile_WithInput(t *testing.T) {
	content := buildTestCasesFile("[2,7,11,15]\n9")
	if !strings.Contains(content, "input:") || !strings.Contains(content, "output:") {
		t.Error("should contain input: and output: markers")
	}
	if !strings.Contains(content, "[2,7,11,15]") {
		t.Error("should contain the test input data")
	}
}

func TestBuildTestCasesFile_EmptyInput(t *testing.T) {
	content := buildTestCasesFile("")
	if !strings.Contains(content, "input:") || !strings.Contains(content, "output:") {
		t.Error("should contain markers even with empty input")
	}
}

// --- WARMUP_EDITOR / shellSplit / buildEditorCmd tests ---

func TestShellSplit(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"nvim file.go", []string{"nvim", "file.go"}},
		{"nvim -c 'split q.md'", []string{"nvim", "-c", "split q.md"}},
		{`nvim {{.CodeFile}} -c 'split {{.DescriptionFile}}' -c 'vsplit {{.TestCasesFile}}'`,
			[]string{"nvim", "{{.CodeFile}}", "-c", "split {{.DescriptionFile}}", "-c", "vsplit {{.TestCasesFile}}"}},
		{"code --new-window", []string{"code", "--new-window"}},
		{"", nil},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := shellSplit(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("shellSplit(%q) = %v (len %d), want %v (len %d)", tt.input, got, len(got), tt.want, len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("shellSplit(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestBuildEditorCmd_Templates(t *testing.T) {
	cmd := buildEditorCmd(
		"nvim {{.CodeFile}} -c 'split {{.DescriptionFile}}' -c 'vsplit {{.TestCasesFile}}'",
		"/tmp/solution.go",
		"/tmp/question.md",
		"/tmp/testcases.txt",
		"/tmp",
	)

	if cmd.Path == "" {
		t.Fatal("command path should not be empty")
	}

	args := cmd.Args // Args[0] is the command itself
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "/tmp/solution.go") {
		t.Error("should substitute CodeFile")
	}
	if !strings.Contains(joined, "split /tmp/question.md") {
		t.Error("should substitute DescriptionFile")
	}
	if !strings.Contains(joined, "vsplit /tmp/testcases.txt") {
		t.Error("should substitute TestCasesFile")
	}
}

func TestBuildEditorCmd_FolderTemplate(t *testing.T) {
	cmd := buildEditorCmd(
		"code {{.Folder}}",
		"/tmp/sol.go", "/tmp/q.md", "/tmp/tc.txt", "/tmp/mydir",
	)
	joined := strings.Join(cmd.Args, " ")
	if !strings.Contains(joined, "/tmp/mydir") {
		t.Error("should substitute Folder")
	}
}

// --- lookupLang tests ---

func TestLookupLang_Defaults(t *testing.T) {
	info := lookupLang("unknown-lang")
	if info.Extension != ".go" {
		t.Errorf("unknown lang should default to Go, got %q", info.Extension)
	}
}

func TestLookupLang_NormalizesInput(t *testing.T) {
	tests := []struct {
		input   string
		wantExt string
	}{
		{"go", ".go"}, {"golang", ".go"},
		{"python", ".py"}, {"python3", ".py"}, {"py", ".py"},
		{"cpp", ".cpp"}, {"c++", ".cpp"},
		{"java", ".java"},
		{"js", ".js"}, {"javascript", ".js"},
		{"ts", ".ts"}, {"typescript", ".ts"},
		{"rust", ".rs"}, {"c", ".c"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := lookupLang(tt.input).Extension; got != tt.wantExt {
				t.Errorf("lookupLang(%q).Extension = %q, want %q", tt.input, got, tt.wantExt)
			}
		})
	}
}

func TestLookupLang_HasCommentStyle(t *testing.T) {
	for name, info := range langTable {
		if info.CommentLine == "" {
			t.Errorf("langTable[%q] missing CommentLine", name)
		}
	}
}

func TestFindEditor(t *testing.T) {
	original := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", original)

	os.Setenv("EDITOR", "custom-editor")
	if got := findEditor(); got != "custom-editor" {
		t.Errorf("findEditor() = %q, want %q", got, "custom-editor")
	}

	os.Unsetenv("EDITOR")
	if got := findEditor(); got == "" {
		t.Error("findEditor() should return a fallback editor")
	}
}
