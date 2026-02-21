package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// Code markers used to delimit user solution code from boilerplate.
// These match leetgo's convention so tooling/habits transfer.
const (
	codeBeginMarker = "@lc code=begin"
	codeEndMarker   = "@lc code=end"
)

// langInfo holds language-specific metadata for scaffolding solution files.
type langInfo struct {
	Extension   string // file extension including dot (e.g. ".go")
	Header      string // boilerplate prepended before the code markers
	EmptyFile   string // content when no code snippet is available
	CommentLine string // single-line comment prefix (e.g. "//", "#")
}

var langTable = map[string]langInfo{
	"go": {
		Extension:   ".go",
		Header:      "",
		EmptyFile:   "package main\n\nfunc main() {}\n\n// TODO: implement your solution here\n",
		CommentLine: "//",
	},
	"python3": {
		Extension:   ".py",
		Header:      "",
		EmptyFile:   "# TODO: implement your solution here\n",
		CommentLine: "#",
	},
	"java": {
		Extension:   ".java",
		Header:      "",
		EmptyFile:   "// TODO: implement your solution here\n",
		CommentLine: "//",
	},
	"cpp": {
		Extension:   ".cpp",
		Header:      "#include <vector>\n#include <string>\nusing namespace std;\n\n",
		EmptyFile:   "#include <vector>\n#include <string>\nusing namespace std;\n\n// TODO: implement your solution here\n",
		CommentLine: "//",
	},
	"javascript": {
		Extension:   ".js",
		Header:      "",
		EmptyFile:   "// TODO: implement your solution here\n",
		CommentLine: "//",
	},
	"typescript": {
		Extension:   ".ts",
		Header:      "",
		EmptyFile:   "// TODO: implement your solution here\n",
		CommentLine: "//",
	},
	"rust": {
		Extension:   ".rs",
		Header:      "",
		EmptyFile:   "// TODO: implement your solution here\n",
		CommentLine: "//",
	},
	"c": {
		Extension:   ".c",
		Header:      "#include <stdlib.h>\n\n",
		EmptyFile:   "#include <stdlib.h>\n\n// TODO: implement your solution here\n",
		CommentLine: "//",
	},
}

// lookupLang returns the langInfo for the given language slug, defaulting to Go.
func lookupLang(lang string) langInfo {
	slug := normalizeLangSlug(lang)
	if info, ok := langTable[slug]; ok {
		return info
	}
	return langTable["go"]
}

// ScaffoldProblem contains the data needed to scaffold all files for a problem.
type ScaffoldProblem struct {
	Provider    string
	ProblemID   string
	Title       string
	Difficulty  string
	Description string
	Example     string
	CodeSnippet string
	TestInput   string
	Meta        *FuncMeta // function metadata for test harness generation
}

// Scaffold manages the workspace directory where solution files are stored.
//
// Layout per problem (mirrors leetgo):
//
//	<baseDir>/<provider>/<problem-id>/
//	  solution.<ext>    — code markers + snippet + test harness
//	  question.md       — problem description
//	  testcases.txt     — input:/output: test case pairs
type Scaffold struct {
	baseDir  string
	lang     langInfo
	langSlug string // normalized language slug (e.g. "golang", "python3")
}

// NewScaffold creates a scaffold manager rooted at the given directory.
// If baseDir is empty, defaults to ~/.config/warmup/workspace.
func NewScaffold(baseDir string, lang string) (*Scaffold, error) {
	if baseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home dir: %w", err)
		}
		baseDir = filepath.Join(home, ".config", "warmup", "workspace")
	}
	slug := normalizeLangSlug(lang)
	return &Scaffold{baseDir: baseDir, lang: lookupLang(lang), langSlug: slug}, nil
}

// SolutionDir returns the directory for a specific problem.
func (s *Scaffold) SolutionDir(provider, problemID string) string {
	return filepath.Join(s.baseDir, provider, problemID)
}

// SolutionPath returns the full path to the solution file for a problem.
func (s *Scaffold) SolutionPath(provider, problemID string) string {
	return filepath.Join(s.SolutionDir(provider, problemID), "solution"+s.lang.Extension)
}

// QuestionPath returns the full path to the question description file.
func (s *Scaffold) QuestionPath(provider, problemID string) string {
	return filepath.Join(s.SolutionDir(provider, problemID), "question.md")
}

// TestCasesPath returns the full path to the test cases file.
func (s *Scaffold) TestCasesPath(provider, problemID string) string {
	return filepath.Join(s.SolutionDir(provider, problemID), "testcases.txt")
}

// Exists returns whether a solution file already exists for this problem.
func (s *Scaffold) Exists(provider, problemID string) bool {
	_, err := os.Stat(s.SolutionPath(provider, problemID))
	return err == nil
}

// EnsureScaffold creates all scaffold files for a problem if the solution file doesn't exist.
// If the solution file already exists, all files are left untouched (preserves user's code).
// Returns the path to the solution file.
func (s *Scaffold) EnsureScaffold(p ScaffoldProblem) (string, error) {
	solPath := s.SolutionPath(p.Provider, p.ProblemID)

	if s.Exists(p.Provider, p.ProblemID) {
		return solPath, nil
	}

	dir := s.SolutionDir(p.Provider, p.ProblemID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create solution dir: %w", err)
	}

	// Write solution file.
	solContent := buildSolutionFile(p, s.lang, s.langSlug)
	if err := os.WriteFile(solPath, []byte(solContent), 0o644); err != nil {
		return "", fmt.Errorf("write solution file: %w", err)
	}

	// Write question.md (best-effort — don't fail scaffold if this fails).
	qContent := buildQuestionFile(p)
	qPath := s.QuestionPath(p.Provider, p.ProblemID)
	if err := os.WriteFile(qPath, []byte(qContent), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write %s: %v\n", qPath, err)
	}

	// Write testcases.txt (best-effort).
	tcContent := buildTestCasesFile(p.TestInput)
	tcPath := s.TestCasesPath(p.Provider, p.ProblemID)
	if err := os.WriteFile(tcPath, []byte(tcContent), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write %s: %v\n", tcPath, err)
	}

	return solPath, nil
}

// OpenInEditor opens files in the user's preferred editor.
//
// Checks WARMUP_EDITOR first, which supports leetgo-style templates:
//
//	{{.CodeFile}}        — path to solution file
//	{{.DescriptionFile}} — path to question.md
//	{{.TestCasesFile}}   — path to testcases.txt
//	{{.Folder}}          — path to the problem directory
//
// Example: WARMUP_EDITOR="nvim {{.CodeFile}} -c 'split {{.DescriptionFile}}' -c 'vsplit {{.TestCasesFile}}'"
//
// Falls back to $EDITOR with just the solution file, then to common editors.
func (s *Scaffold) OpenInEditor(provider, problemID string) *exec.Cmd {
	codePath := s.SolutionPath(provider, problemID)
	descPath := s.QuestionPath(provider, problemID)
	tcPath := s.TestCasesPath(provider, problemID)
	folder := s.SolutionDir(provider, problemID)

	if warmupEditor := os.Getenv("WARMUP_EDITOR"); warmupEditor != "" {
		return buildEditorCmd(warmupEditor, codePath, descPath, tcPath, folder)
	}

	editor := findEditor()
	return exec.Command(editor, codePath)
}

// buildEditorCmd parses a WARMUP_EDITOR template string and builds the exec.Cmd.
func buildEditorCmd(editorTemplate, codePath, descPath, tcPath, folder string) *exec.Cmd {
	data := struct {
		CodeFile        string
		DescriptionFile string
		TestCasesFile   string
		Folder          string
	}{codePath, descPath, tcPath, folder}

	// Split the template into tokens using shell-like splitting.
	tokens := shellSplit(editorTemplate)
	if len(tokens) == 0 {
		return exec.Command("vi", codePath)
	}

	// Apply template substitution to each token.
	args := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if strings.Contains(tok, "{{") {
			tmpl, err := template.New("").Parse(tok)
			if err == nil {
				var buf strings.Builder
				if err := tmpl.Execute(&buf, data); err == nil {
					tok = buf.String()
				}
			}
		}
		args = append(args, tok)
	}

	return exec.Command(args[0], args[1:]...)
}

// shellSplit performs basic POSIX-style shell tokenization.
// Handles single quotes, double quotes, and unquoted tokens.
func shellSplit(s string) []string {
	var tokens []string
	var current strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"' && !inSingle:
			inDouble = !inDouble
		case c == ' ' && !inSingle && !inDouble:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

// findEditor returns the editor command to use.
func findEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	for _, e := range []string{"nvim", "vim", "vi", "nano"} {
		if _, err := exec.LookPath(e); err == nil {
			return e
		}
	}
	return "vi"
}

// --- Code extraction ---

// ExtractCode reads the solution file and returns the code between @lc markers.
// If no markers are found, returns the entire file content (for non-scaffolded files).
func (s *Scaffold) ExtractCode(provider, problemID string) (string, error) {
	solPath := s.SolutionPath(provider, problemID)
	data, err := os.ReadFile(solPath)
	if err != nil {
		return "", fmt.Errorf("read solution file: %w", err)
	}

	content := string(data)
	return extractCodeBetweenMarkers(content), nil
}

// extractCodeBetweenMarkers extracts code between @lc code=begin and @lc code=end markers.
// Returns the full content if markers aren't found.
func extractCodeBetweenMarkers(content string) string {
	lines := strings.Split(content, "\n")

	beginIdx := -1
	endIdx := -1
	for i, line := range lines {
		if strings.Contains(line, codeBeginMarker) {
			beginIdx = i
		}
		if strings.Contains(line, codeEndMarker) {
			endIdx = i
			break
		}
	}

	if beginIdx == -1 || endIdx == -1 || beginIdx >= endIdx {
		return content
	}

	// Extract lines between markers (exclusive of marker lines).
	extracted := lines[beginIdx+1 : endIdx]
	return strings.TrimSpace(strings.Join(extracted, "\n"))
}

// ReadTestInput reads the test input from testcases.txt.
// Returns the content between "input:" and "output:" markers.
func (s *Scaffold) ReadTestInput(provider, problemID string) (string, error) {
	tcPath := s.TestCasesPath(provider, problemID)
	data, err := os.ReadFile(tcPath)
	if err != nil {
		return "", fmt.Errorf("read test cases file: %w", err)
	}

	content := string(data)
	return parseTestCasesInput(content), nil
}

// parseTestCasesInput extracts the input section from testcases.txt content.
// Format: "input:\n<data>\noutput:\n<data>"
func parseTestCasesInput(content string) string {
	lines := strings.Split(content, "\n")

	inInput := false
	var inputLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "input:" {
			inInput = true
			continue
		}
		if trimmed == "output:" {
			break
		}
		if inInput {
			inputLines = append(inputLines, line)
		}
	}

	return strings.TrimSpace(strings.Join(inputLines, "\n"))
}

// --- File builders ---

// buildSolutionFile generates the complete solution file with code markers and test harness.
func buildSolutionFile(p ScaffoldProblem, lang langInfo, langSlug string) string {
	if p.CodeSnippet == "" {
		return lang.EmptyFile
	}

	// For Go, generate a full file with func main() test harness.
	if langSlug == "golang" {
		return buildGoSolutionFile(p)
	}

	// For other languages, generate code markers without a test harness.
	var b strings.Builder
	b.WriteString(lang.Header)
	b.WriteString(lang.CommentLine + " " + codeBeginMarker + "\n\n")
	b.WriteString(p.CodeSnippet)
	b.WriteString("\n\n")
	b.WriteString(lang.CommentLine + " " + codeEndMarker + "\n")
	return b.String()
}

// buildGoSolutionFile generates a Go solution file with:
//   - package main + imports
//   - @lc code=begin / @lc code=end markers around the user's solution
//   - func main() test harness that reads from stdin, calls the solution, prints output
//   - inline readLine helper
func buildGoSolutionFile(p ScaffoldProblem) string {
	var b strings.Builder

	// Determine what imports we need.
	imports := []string{`"bufio"`, `"encoding/json"`, `"fmt"`, `"os"`, `"strings"`}

	b.WriteString("package main\n\nimport (\n")
	for _, imp := range imports {
		b.WriteString("\t" + imp + "\n")
	}
	b.WriteString(")\n\n")

	// Code markers + user solution.
	b.WriteString("// " + codeBeginMarker + "\n\n")
	b.WriteString(p.CodeSnippet)
	b.WriteString("\n\n// " + codeEndMarker + "\n\n")

	// Generate func main() test harness.
	b.WriteString(buildGoMainFunc(p.Meta))

	// Inline helper: readLine.
	b.WriteString(goReadLineHelper)

	return b.String()
}

// buildGoMainFunc generates the func main() test harness from FuncMeta.
func buildGoMainFunc(meta *FuncMeta) string {
	if meta == nil || meta.SystemDesign {
		return goFallbackMain
	}

	var b strings.Builder
	b.WriteString("func main() {\n")
	b.WriteString("\tstdin := bufio.NewReader(os.Stdin)\n\n")

	// Deserialize each parameter from stdin.
	for _, param := range meta.Params {
		goType := lcTypeToGo(param.Type)
		b.WriteString(fmt.Sprintf("\tvar %s %s\n", param.Name, goType))
		b.WriteString(fmt.Sprintf("\tjson.Unmarshal([]byte(readLine(stdin)), &%s)\n\n", param.Name))
	}

	// Build the function call.
	paramNames := make([]string, len(meta.Params))
	for i, param := range meta.Params {
		paramNames[i] = param.Name
	}
	call := fmt.Sprintf("%s(%s)", meta.Name, strings.Join(paramNames, ", "))

	if meta.Return != nil && meta.Return.Type != "" {
		// Has return value: capture and serialize.
		b.WriteString(fmt.Sprintf("\tans := %s\n", call))
		b.WriteString("\tout, _ := json.Marshal(ans)\n")
		b.WriteString("\tfmt.Println(\"\\noutput:\", string(out))\n")
	} else {
		// Void function (in-place modification): call and serialize the first param.
		b.WriteString(fmt.Sprintf("\t%s\n", call))
		if len(meta.Params) > 0 {
			b.WriteString(fmt.Sprintf("\tout, _ := json.Marshal(%s)\n", meta.Params[0].Name))
			b.WriteString("\tfmt.Println(\"\\noutput:\", string(out))\n")
		}
	}

	b.WriteString("}\n\n")
	return b.String()
}

// lcTypeToGo maps LeetCode type strings to Go types.
func lcTypeToGo(lcType string) string {
	switch strings.ToLower(lcType) {
	case "integer", "int":
		return "int"
	case "integer[]", "int[]":
		return "[]int"
	case "integer[][]", "int[][]":
		return "[][]int"
	case "string":
		return "string"
	case "string[]":
		return "[]string"
	case "string[][]":
		return "[][]string"
	case "boolean", "bool":
		return "bool"
	case "boolean[]":
		return "[]bool"
	case "double", "float", "long":
		return "float64"
	case "double[]", "float[]":
		return "[]float64"
	case "character[]", "char[]":
		return "[]byte"
	case "character[][]", "char[][]":
		return "[][]byte"
	case "treenode":
		return "*TreeNode"
	case "listnode":
		return "*ListNode"
	default:
		return "interface{}"
	}
}

const goReadLineHelper = `func readLine(r *bufio.Reader) string {
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}
`

const goFallbackMain = `func main() {
	// TODO: test harness not auto-generated for this problem type.
	// Read input from stdin, call your solution, and print the result.
	// See testcases.txt for the expected input/output format.
}

`

// buildQuestionFile generates question.md content matching leetgo's format.
func buildQuestionFile(p ScaffoldProblem) string {
	var b strings.Builder

	// Header: # Title (Difficulty)
	b.WriteString(fmt.Sprintf("# %s (%s)\n\n", p.Title, p.Difficulty))

	// LeetCode link (if provider is leetcode).
	if p.Provider == "leetcode" || p.Provider == "" {
		b.WriteString(fmt.Sprintf("[LeetCode](https://leetcode.com/problems/%s/)\n\n", p.ProblemID))
	}

	// Description.
	if p.Description != "" {
		b.WriteString(p.Description)
		b.WriteString("\n")
	}

	// Example.
	if p.Example != "" {
		b.WriteString("\n## Example\n\n")
		b.WriteString("```\n")
		b.WriteString(p.Example)
		b.WriteString("\n```\n")
	}

	return b.String()
}

// buildTestCasesFile generates testcases.txt in leetgo's input:/output: format.
func buildTestCasesFile(testInput string) string {
	if testInput == "" {
		return "input:\n\noutput:\n"
	}

	var b strings.Builder
	b.WriteString("input:\n")
	b.WriteString(strings.TrimSpace(testInput))
	b.WriteString("\noutput:\n")

	return b.String()
}
