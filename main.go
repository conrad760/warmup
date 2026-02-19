package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Rating classifies how good an approach is for a given problem.
type Rating int

const (
	Optimal Rating = iota
	Plausible
	Suboptimal
	Wrong
)

func (r Rating) String() string {
	switch r {
	case Optimal:
		return "OPTIMAL"
	case Plausible:
		return "PLAUSIBLE"
	case Suboptimal:
		return "SUBOPTIMAL"
	case Wrong:
		return "WRONG"
	}
	return ""
}

func (r Rating) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Rating) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	case "OPTIMAL":
		*r = Optimal
	case "PLAUSIBLE":
		*r = Plausible
	case "SUBOPTIMAL":
		*r = Suboptimal
	case "WRONG":
		*r = Wrong
	default:
		return fmt.Errorf("invalid rating: %s", s)
	}
	return nil
}

type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
)

func (d Difficulty) String() string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	}
	return ""
}

func (d Difficulty) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Difficulty) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "easy":
		*d = Easy
	case "medium":
		*d = Medium
	case "hard":
		*d = Hard
	default:
		return fmt.Errorf("invalid difficulty: %s", s)
	}
	return nil
}

type Option struct {
	Text   string
	Rating Rating
}

type Question struct {
	Title        string
	Difficulty   Difficulty
	Category     string
	Description  string
	Example      string
	Options      []Option
	Solution     string
	LeetcodeSlug string
}

type Stats struct {
	Total      int
	Optimal    int
	Plausible  int
	Suboptimal int
	Wrong      int
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type leetgoResultMsg struct {
	output string
	err    error
}

type model struct {
	questions      []Question
	unseen         []int
	currentIdx     int
	shuffledOpts   []Option
	cursor         int
	selected       int // -1 = not yet answered
	revealed       bool
	showSolution   bool
	timer          int
	timerRunning   bool
	leetgoWorkdir  string    // leetgo workspace directory for "Try It"
	triedIt        bool      // whether user has opened leetgo for current question
	leetgoRunning  bool      // background test/submit in progress
	leetgoOutput   string    // captured stdout+stderr from last test/submit
	leetgoSpinner  int       // spinner frame index
	leetgoAction   string    // "test" or "submit" — for status label
	leetgoStarted  time.Time // wall-clock time when leetgo coding started
	leetgoElapsed  int       // frozen elapsed seconds (set when timer stops)
	leetgoTimerOn  bool      // whether the leetgo coding timer is running
	lastMaxScroll  int       // max scroll offset from last render
	timerExpired   bool
	stats          Stats
	width          int
	height         int
	pulseOn        bool // for timer pulse effect
	scrollOffset   int  // viewport scroll offset
	reviewLog      *ReviewLog
	sessionEntries []sessionEntry
	sessionSeen    map[int]bool
	submitResult   int    // current question's submit result (SubmitNone initially)
	sessionReport  string // filled on quit, rendered after program exits
	leetgoMessage  string // transient message shown when user tries unavailable action
	categoryFilter string // active category filter (empty = all categories)
}

const defaultTimer = 300

// leetgoCodingSecs returns the total elapsed seconds on the leetgo coding timer.
func (m model) leetgoCodingSecs() int {
	if !m.leetgoTimerOn {
		return m.leetgoElapsed
	}
	return m.leetgoElapsed + int(time.Since(m.leetgoStarted).Seconds())
}

func (m *model) pickQuestion() {
	idx := -1
	if m.reviewLog != nil {
		idx = m.reviewLog.PickNextQuestion(m.questions, m.sessionSeen)
	}
	if idx == -1 {
		if len(m.unseen) == 0 {
			m.unseen = make([]int, len(m.questions))
			for i := range m.questions {
				m.unseen[i] = i
			}
		}
		ri := rand.Intn(len(m.unseen))
		idx = m.unseen[ri]
		m.unseen = append(m.unseen[:ri], m.unseen[ri+1:]...)
	}
	m.currentIdx = idx
	m.sessionSeen[idx] = true
	q := m.questions[m.currentIdx]
	m.shuffledOpts = make([]Option, len(q.Options))
	copy(m.shuffledOpts, q.Options)
	rand.Shuffle(len(m.shuffledOpts), func(i, j int) {
		m.shuffledOpts[i], m.shuffledOpts[j] = m.shuffledOpts[j], m.shuffledOpts[i]
	})
	m.cursor = 0
	m.selected = -1
	m.revealed = false
	m.showSolution = false
	m.timer = defaultTimer
	m.timerExpired = false
	m.scrollOffset = 0
	m.submitResult = SubmitNone
}

func (m *model) saveCurrentReview() {
	if m.reviewLog == nil || m.selected == -1 {
		return
	}
	q := m.questions[m.currentIdx]
	approach := m.shuffledOpts[m.selected].Rating

	codingTime := 0
	if m.triedIt {
		codingTime = m.leetgoCodingSecs()
	}
	m.reviewLog.RecordReview(q.LeetcodeSlug, approach, m.submitResult, codingTime)

	m.sessionEntries = append(m.sessionEntries, sessionEntry{
		title:        q.Title,
		slug:         q.LeetcodeSlug,
		approach:     approach,
		submitResult: m.submitResult,
		codingTime:   codingTime,
	})

	_ = m.reviewLog.Save()
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.lastMaxScroll = viewMaxScroll
		if m.leetgoRunning {
			m.leetgoSpinner++
		}
		if m.timerRunning && m.timer > 0 {
			m.timer--
			if m.timer <= 30 {
				m.pulseOn = !m.pulseOn
			}
			if m.timer == 0 {
				m.timerRunning = false
				m.timerExpired = true
			}
		}
		return m, tickCmd()

	case leetgoResultMsg:
		m.leetgoRunning = false
		m.leetgoOutput = stripANSI(msg.output)
		if msg.err != nil && m.leetgoOutput == "" {
			m.leetgoOutput = fmt.Sprintf("Error: %v", msg.err)
		}
		if m.leetgoAction == "submit" {
			m.submitResult = ParseSubmitResult(m.leetgoOutput)
			if m.submitResult == SubmitAccepted && m.leetgoTimerOn {
				m.leetgoElapsed = m.leetgoCodingSecs()
				m.leetgoTimerOn = false
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.revealed {
				m.saveCurrentReview()
			}
			if m.reviewLog != nil {
				m.sessionReport = m.reviewLog.SessionReport(m.sessionEntries, len(m.questions))
			}
			return m, tea.Quit

		case "up", "k":
			if m.selected == -1 {
				if m.cursor > 0 {
					m.cursor--
				} else if m.scrollOffset > 0 {
					m.scrollOffset--
				}
			} else {
				if m.scrollOffset > 0 {
					m.scrollOffset--
				}
			}
			return m, nil

		case "down", "j":
			if m.selected == -1 {
				if m.cursor < len(m.shuffledOpts)-1 {
					m.cursor++
				} else if m.scrollOffset < m.lastMaxScroll {
					m.scrollOffset++
				}
			} else {
				if m.scrollOffset < m.lastMaxScroll {
					m.scrollOffset++
				}
			}
			return m, nil

		case "enter", " ":
			if m.selected == -1 && len(m.shuffledOpts) > 0 {
				m.selected = m.cursor
				m.revealed = true
				m.timerRunning = false
				rating := m.shuffledOpts[m.selected].Rating
				m.stats.Total++
				switch rating {
				case Optimal:
					m.stats.Optimal++
				case Plausible:
					m.stats.Plausible++
				case Suboptimal:
					m.stats.Suboptimal++
				case Wrong:
					m.stats.Wrong++
				}
				m.scrollOffset = 0
			}
			return m, nil

		case "s":
			if m.revealed {
				m.showSolution = !m.showSolution
				if !m.showSolution {
					m.scrollOffset = 0
				}
			}
			return m, nil

		case "t":
			if m.canTryIt() && !m.leetgoRunning {
				if !m.triedIt {
					m.leetgoElapsed = 0
					m.leetgoStarted = time.Now()
					m.leetgoTimerOn = true
				}
				m.triedIt = true
				m.leetgoMessage = ""
				return m, m.tryItCmd()
			}
			if m.wantsTryIt() && m.leetgoWorkdir == "" {
				m.leetgoMessage = "leetgo workspace not configured. To enable:\n" +
					"  1. mkdir ~/leetcode && cd ~/leetcode && leetgo init\n" +
					"  2. Restart with: dsa-warmup -workspace ~/leetcode"
			}
			return m, nil

		case "T":
			if m.triedIt && !m.leetgoRunning {
				m.leetgoRunning = true
				m.leetgoAction = "test"
				m.leetgoOutput = ""
				m.leetgoSpinner = 0
				m.leetgoMessage = ""
				return m, m.leetgoRunCmd("test")
			}
			if m.wantsTryIt() && m.leetgoWorkdir == "" {
				m.leetgoMessage = "leetgo workspace not configured. To enable:\n" +
					"  1. mkdir ~/leetcode && cd ~/leetcode && leetgo init\n" +
					"  2. Restart with: dsa-warmup -workspace ~/leetcode"
			}
			return m, nil

		case "S":
			if m.triedIt && !m.leetgoRunning {
				m.leetgoRunning = true
				m.leetgoAction = "submit"
				m.leetgoOutput = ""
				m.leetgoSpinner = 0
				m.leetgoMessage = ""
				return m, m.leetgoRunCmd("submit")
			}
			if m.wantsTryIt() && m.leetgoWorkdir == "" {
				m.leetgoMessage = "leetgo workspace not configured. To enable:\n" +
					"  1. mkdir ~/leetcode && cd ~/leetcode && leetgo init\n" +
					"  2. Restart with: dsa-warmup -workspace ~/leetcode"
			}
			return m, nil

		case "n":
			if m.revealed && !m.leetgoRunning {
				m.saveCurrentReview()
				m.pickQuestion()
				m.timerRunning = true
				m.triedIt = false
				m.leetgoElapsed = 0
				m.leetgoTimerOn = false
				m.leetgoOutput = ""
				m.leetgoAction = ""
				m.leetgoMessage = ""
			}
			return m, nil

		case "p":
			if !m.timerExpired && m.selected == -1 {
				m.timerRunning = !m.timerRunning
			}
			return m, nil

		case "r":
			if m.selected == -1 {
				m.timer = defaultTimer
				m.timerExpired = false
				m.timerRunning = true
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) canTryIt() bool {
	if !m.revealed || m.selected == -1 || m.leetgoWorkdir == "" {
		return false
	}
	rating := m.shuffledOpts[m.selected].Rating
	if rating != Optimal && rating != Plausible {
		return false
	}
	return m.questions[m.currentIdx].LeetcodeSlug != ""
}

// wantsTryIt reports whether the user earned Try It but may lack a workspace.
func (m model) wantsTryIt() bool {
	if !m.revealed || m.selected == -1 {
		return false
	}
	rating := m.shuffledOpts[m.selected].Rating
	if rating != Optimal && rating != Plausible {
		return false
	}
	return m.questions[m.currentIdx].LeetcodeSlug != ""
}

func (m model) problemExists(slug string) bool {
	if m.leetgoWorkdir == "" {
		return false
	}
	pattern := m.leetgoWorkdir + "/go/*." + slug + "/solution.go"
	matches, err := filepath.Glob(pattern)
	return err == nil && len(matches) > 0
}

// tryItCmd opens the problem in leetgo, using "edit" if already scaffolded
// or "pick" to scaffold for the first time.
func (m model) tryItCmd() tea.Cmd {
	slug := m.questions[m.currentIdx].LeetcodeSlug
	var cmd *exec.Cmd
	if m.problemExists(slug) {
		cmd = exec.Command("leetgo", "edit", slug, "-l", "go")
	} else {
		cmd = exec.Command("leetgo", "pick", slug, "-l", "go")
	}
	if m.leetgoWorkdir != "" {
		cmd.Dir = m.leetgoWorkdir
	}
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return nil
	})
}

// leetgoRunCmd runs a leetgo subcommand in the background and returns a leetgoResultMsg.
func (m model) leetgoRunCmd(subcmd string) tea.Cmd {
	slug := m.questions[m.currentIdx].LeetcodeSlug
	workdir := m.leetgoWorkdir
	return func() tea.Msg {
		cmd := exec.Command("leetgo", subcmd, slug, "-l", "go")
		if workdir != "" {
			cmd.Dir = workdir
		}
		out, err := cmd.CombinedOutput()
		return leetgoResultMsg{output: string(out), err: err}
	}
}

var (
	colorWhite  = lipgloss.Color("#e4e4e7")
	colorDim    = lipgloss.Color("#71717a")
	colorBlue   = lipgloss.Color("#6366f1")
	colorPurple = lipgloss.Color("#a855f7")
	colorGreen  = lipgloss.Color("#22c55e")
	colorYellow = lipgloss.Color("#eab308")
	colorOrange = lipgloss.Color("#f97316")
	colorRed    = lipgloss.Color("#ef4444")

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBlue).
			PaddingLeft(1).
			PaddingRight(1)

	statsStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	timerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite)

	timerPulseStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRed)

	timerExpiredStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorRed)

	colorSoftWhite = lipgloss.Color("#d4d4d8")
	colorMidGrey   = lipgloss.Color("#a1a1aa")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorMidGrey)

	descStyle = lipgloss.NewStyle().
			Foreground(colorSoftWhite)

	exampleStyle = lipgloss.NewStyle().
			Foreground(colorSoftWhite).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDim)

	optionStyle = lipgloss.NewStyle().
			Foreground(colorWhite)

	optionSelectedStyle = lipgloss.NewStyle().
				Foreground(colorBlue).
				Bold(true)

	bannerOptimalStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorGreen)

	bannerPlausibleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorYellow)

	bannerSuboptimalStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorOrange)

	bannerWrongStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorRed)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	badgeEasy = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(colorGreen).
			Bold(true).
			Padding(0, 1)

	badgeMedium = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(colorYellow).
			Bold(true).
			Padding(0, 1)

	badgeHard = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(colorRed).
			Bold(true).
			Padding(0, 1)

	badgeCategory = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(colorPurple).
			Padding(0, 1)

	ratingOptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGreen)

	ratingPlausibleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorYellow)

	ratingSuboptStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorOrange)

	ratingWrongStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorRed)

	progressFull  = lipgloss.NewStyle().Foreground(colorBlue)
	progressEmpty = lipgloss.NewStyle().Foreground(colorDim)
)

// viewMaxScroll is set by View() so Update() can clamp scroll offsets.
// This is needed because View() has a value receiver and can't mutate the model.
var viewMaxScroll int

func diffBadge(d Difficulty) string {
	switch d {
	case Easy:
		return badgeEasy.Render(" Easy ")
	case Medium:
		return badgeMedium.Render(" Medium ")
	case Hard:
		return badgeHard.Render(" Hard ")
	}
	return ""
}

func ratingLabel(r Rating) string {
	switch r {
	case Optimal:
		return ratingOptStyle.Render("[OPTIMAL]")
	case Plausible:
		return ratingPlausibleStyle.Render("[PLAUSIBLE]")
	case Suboptimal:
		return ratingSuboptStyle.Render("[SUBOPTIMAL]")
	case Wrong:
		return ratingWrongStyle.Render("[WRONG]")
	}
	return ""
}

func (m model) View() string {
	maxW := 84
	if m.width < maxW {
		maxW = m.width
	}
	textW := maxW - 6
	if textW < 40 {
		textW = 40
	}
	optW := textW - 10

	var b strings.Builder
	q := m.questions[m.currentIdx]

	headerText := "DSA Warmup"
	if m.categoryFilter != "" {
		headerText = "DSA Warmup — " + m.categoryFilter
	}
	header := headerStyle.Render(headerText)
	statStr := fmt.Sprintf("  #%d  ", m.stats.Total)
	statStr += lipgloss.NewStyle().
		Foreground(colorGreen).
		Render(fmt.Sprintf("O:%d", m.stats.Optimal))
	statStr += " "
	statStr += lipgloss.NewStyle().
		Foreground(colorYellow).
		Render(fmt.Sprintf("P:%d", m.stats.Plausible))
	statStr += " "
	statStr += lipgloss.NewStyle().
		Foreground(colorOrange).
		Render(fmt.Sprintf("S:%d", m.stats.Suboptimal))
	statStr += " "
	statStr += lipgloss.NewStyle().Foreground(colorRed).Render(fmt.Sprintf("W:%d", m.stats.Wrong))
	b.WriteString(header + statsStyle.Render(statStr))
	b.WriteString("\n")

	mins := m.timer / 60
	secs := m.timer % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)

	barWidth := 20
	filled := 0
	if defaultTimer > 0 {
		filled = (m.timer * barWidth) / defaultTimer
	}
	if filled < 0 {
		filled = 0
	}
	bar := progressFull.Render(strings.Repeat("█", filled)) +
		progressEmpty.Render(strings.Repeat("░", barWidth-filled))

	var timerLine string
	if m.timerExpired {
		timerLine = timerExpiredStyle.Render("TIME'S UP") + " " + bar
	} else if m.timer <= 30 && m.pulseOn {
		timerLine = timerPulseStyle.Render(timeStr) + " " + bar
	} else {
		timerLine = timerStyle.Render(timeStr) + " " + bar
	}

	pauseIndicator := ""
	if !m.timerRunning && !m.timerExpired && m.selected == -1 {
		pauseIndicator = lipgloss.NewStyle().Foreground(colorYellow).Render(" PAUSED")
	}
	b.WriteString(timerLine + pauseIndicator)
	b.WriteString("\n")

	var badges string
	if m.reviewLog != nil {
		if pr, ok := m.reviewLog.Reviews[q.LeetcodeSlug]; ok {
			days := pr.Interval
			reps := pr.Repetitions
			if days >= 21 && reps >= 3 {
				badges = lipgloss.NewStyle().Foreground(colorGreen).Render("★ mastered")
			} else if reps > 0 {
				badges = lipgloss.NewStyle().
					Foreground(colorYellow).
					Render(fmt.Sprintf("↻ %dd", int(days)))
			} else {
				badges = lipgloss.NewStyle().Foreground(colorOrange).Render("↻ learning")
			}
		} else {
			badges = lipgloss.NewStyle().Foreground(colorDim).Render("new")
		}
	}
	if badges != "" {
		b.WriteString(badges)
		b.WriteString("\n")
	}
	b.WriteString(titleStyle.Render(q.Title))
	b.WriteString("\n")

	b.WriteString(descStyle.Render(wrapText(q.Description, textW)))
	b.WriteString("\n")

	if q.Example != "" {
		exLabel := lipgloss.NewStyle().Bold(true).Foreground(colorDim).Render("Example")
		wrapped := wrapLines(q.Example, textW-6)
		exBody := formatExample(wrapped, textW-6)
		b.WriteString(exLabel + "\n")
		b.WriteString(exampleStyle.Render(exBody))
		b.WriteString("\n")
	}

	b.WriteString(
		lipgloss.NewStyle().Foreground(colorWhite).Bold(true).Render("Choose an approach:"),
	)
	b.WriteString("\n")

	letters := []string{"A", "B", "C", "D", "E", "F"}
	for i, opt := range m.shuffledOpts {
		letter := letters[i]
		if m.selected == -1 {
			wrapped := wrapText(opt.Text, optW)
			wrappedLines := strings.Split(wrapped, "\n")
			for li, wl := range wrappedLines {
				if li == 0 {
					if i == m.cursor {
						prefix := lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("> ")
						b.WriteString(
							prefix + optionSelectedStyle.Render(fmt.Sprintf("%s) %s", letter, wl)),
						)
					} else {
						b.WriteString("  " + optionStyle.Render(fmt.Sprintf("%s) %s", letter, wl)))
					}
				} else {
					if i == m.cursor {
						b.WriteString("  " + optionSelectedStyle.Render("   "+wl))
					} else {
						b.WriteString("  " + optionStyle.Render("   "+wl))
					}
				}
				b.WriteString("\n")
			}
		} else {
			marker := "  "
			if i == m.selected {
				marker = lipgloss.NewStyle().Foreground(colorBlue).Bold(true).Render("> ")
			}
			label := ratingLabel(opt.Rating)
			var lineStyle lipgloss.Style
			switch opt.Rating {
			case Optimal:
				lineStyle = lipgloss.NewStyle().Foreground(colorGreen)
			case Plausible:
				lineStyle = lipgloss.NewStyle().Foreground(colorYellow)
			case Suboptimal:
				lineStyle = lipgloss.NewStyle().Foreground(colorOrange)
			case Wrong:
				lineStyle = lipgloss.NewStyle().Foreground(colorRed)
			}
			wrapped := wrapText(opt.Text, optW)
			wrappedLines := strings.Split(wrapped, "\n")
			for li, wl := range wrappedLines {
				if li == 0 {
					b.WriteString(marker + lineStyle.Render(fmt.Sprintf("%s) %s", letter, wl)))
					b.WriteString(" " + label)
				} else {
					b.WriteString("  " + lineStyle.Render("   "+wl))
				}
				b.WriteString("\n")
			}
		}
	}

	if m.revealed {
		b.WriteString(diffBadge(q.Difficulty) + " " + badgeCategory.Render(q.Category))
		b.WriteString("\n")
		rating := m.shuffledOpts[m.selected].Rating
		switch rating {
		case Optimal:
			b.WriteString(bannerOptimalStyle.Render(">> Nailed it!"))
		case Plausible:
			b.WriteString(bannerPlausibleStyle.Render(">> Good instinct -- see the optimal below"))
		case Suboptimal:
			b.WriteString(bannerSuboptimalStyle.Render(">> Works but too slow for interviews"))
		case Wrong:
			b.WriteString(bannerWrongStyle.Render(">> That approach won't work here"))
		}
		b.WriteString("\n")
	}

	if m.showSolution {
		codeW := maxW - 6
		wrappedCode := wrapCodeLines(q.Solution, codeW)
		highlighted := highlightGo(wrappedCode)
		solStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPurple).
			Padding(0, 1)
		b.WriteString(solStyle.Render(highlighted))
		b.WriteString("\n")
	}

	if m.canTryIt() && !m.triedIt {
		tryItStyle := lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		b.WriteString(tryItStyle.Render("Press [t] to open in leetgo"))
		b.WriteString("\n")
	} else if m.wantsTryIt() && m.leetgoWorkdir == "" && !m.triedIt {
		noWorkStyle := lipgloss.NewStyle().Foreground(colorRed)
		b.WriteString(noWorkStyle.Render("[t] Try It — leetgo workspace not configured"))
		b.WriteString("\n")
	}

	if m.triedIt {
		lgElapsed := m.leetgoCodingSecs()
		lgMins := lgElapsed / 60
		lgSecs := lgElapsed % 60
		lgTimeStr := fmt.Sprintf("%02d:%02d", lgMins, lgSecs)
		var lgTimerLine string
		if !m.leetgoTimerOn && lgElapsed > 0 {
			// Timer stopped (accepted) — show final time in green
			lgTimerLine = lipgloss.NewStyle().
				Foreground(colorGreen).
				Bold(true).
				Render("Coding: " + lgTimeStr)
		} else {
			lgTimerLine = lipgloss.NewStyle().
				Foreground(colorBlue).
				Bold(true).
				Render("Coding: " + lgTimeStr)
		}
		b.WriteString(lgTimerLine)
		b.WriteString("\n")
	}

	if m.leetgoMessage != "" {
		msgStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorRed).
			Foreground(colorRed).
			Padding(0, 1).
			MaxWidth(maxW - 2)
		b.WriteString(msgStyle.Render(m.leetgoMessage))
		b.WriteString("\n")
	}

	if m.leetgoRunning {
		frame := spinnerFrames[m.leetgoSpinner%len(spinnerFrames)]
		slug := q.LeetcodeSlug
		action := "Testing"
		if m.leetgoAction == "submit" {
			action = "Submitting"
		}
		spinStyle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true)
		b.WriteString(spinStyle.Render(fmt.Sprintf("%s %s %s...", frame, action, slug)))
		b.WriteString("\n")
	} else if m.leetgoOutput != "" {
		outStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGreen).
			Padding(0, 1).
			MaxWidth(maxW - 2)
		label := "Test"
		if m.leetgoAction == "submit" {
			label = "Submit"
		}
		labelStyle := lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
		b.WriteString(labelStyle.Render(fmt.Sprintf("── %s Result ──", label)))
		b.WriteString("\n")
		wrapped := wrapLines(strings.TrimSpace(m.leetgoOutput), maxW-6)
		b.WriteString(outStyle.Render(wrapped))
		b.WriteString("\n")
	}

	var helpParts []string
	if m.selected == -1 {
		helpParts = append(
			helpParts,
			"[Enter] Select",
			"[p] Pause",
			"[r] Reset",
			"[j/k] Scroll",
			"[q] Quit",
		)
	} else {
		noWorkRedStyle := lipgloss.NewStyle().Foreground(colorRed)
		if m.canTryIt() && !m.triedIt {
			helpParts = append(helpParts, "[t] Try It")
		} else if m.wantsTryIt() && m.leetgoWorkdir == "" {
			helpParts = append(
				helpParts,
				noWorkRedStyle.Render("[t] Try It"),
				noWorkRedStyle.Render("[T] Test"),
				noWorkRedStyle.Render("[S] Submit"),
			)
		}
		if m.triedIt {
			if m.leetgoRunning {
				dimStyle := lipgloss.NewStyle().Foreground(colorDim)
				helpParts = append(
					helpParts,
					dimStyle.Render("[T] Test"),
					dimStyle.Render("[S] Submit"),
				)
			} else {
				helpParts = append(helpParts, "[t] Edit", "[T] Test", "[S] Submit")
			}
		}
		helpParts = append(helpParts, "[n] Next")
		if m.showSolution {
			helpParts = append(helpParts, "[s] Hide Solution")
		} else {
			helpParts = append(helpParts, "[s] Show Solution")
		}
		helpParts = append(helpParts, "[j/k] Scroll", "[q] Quit")
	}
	b.WriteString(helpStyle.Render(strings.Join(helpParts, "  ")))
	b.WriteString("\n")

	raw := b.String()
	lines := strings.Split(raw, "\n")
	totalLines := len(lines)
	viewHeight := m.height - 2
	if viewHeight < 10 {
		viewHeight = 10
	}

	if totalLines <= viewHeight {
		viewMaxScroll = 0
		padded := make([]string, len(lines))
		for i, l := range lines {
			padded[i] = "  " + l
		}
		return strings.Join(padded, "\n")
	}

	maxOffset := totalLines - viewHeight
	viewMaxScroll = maxOffset
	offset := m.scrollOffset
	if offset > maxOffset {
		offset = maxOffset
	}
	if offset < 0 {
		offset = 0
	}

	visible := lines[offset:]
	if len(visible) > viewHeight {
		visible = visible[:viewHeight]
	}

	pct := 0
	if maxOffset > 0 {
		pct = (offset * 100) / maxOffset
	}
	indicatorStyle := lipgloss.NewStyle().Foreground(colorYellow).Bold(true)
	if pct < 100 {
		visible[len(visible)-1] = indicatorStyle.Render(
			fmt.Sprintf("  ▼ %d%% | j/k to scroll ▼", pct),
		)
	} else {
		visible[len(visible)-1] = indicatorStyle.Render("  ✓ End | k to scroll up")
	}

	padded := make([]string, len(visible))
	for i, l := range visible {
		padded[i] = "  " + l
	}
	return strings.Join(padded, "\n")
}

func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		maxWidth = 80
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	var lines []string
	line := words[0]
	for _, w := range words[1:] {
		if len(line)+1+len(w) > maxWidth {
			lines = append(lines, line)
			line = w
		} else {
			line += " " + w
		}
	}
	lines = append(lines, line)
	return strings.Join(lines, "\n")
}

// wrapCodeLines hard-wraps long code lines, indenting continuations.
func wrapCodeLines(code string, maxWidth int) string {
	if maxWidth <= 0 {
		maxWidth = 76
	}
	lines := strings.Split(code, "\n")
	var result []string
	for _, line := range lines {
		if len(line) <= maxWidth {
			result = append(result, line)
			continue
		}
		for len(line) > maxWidth {
			breakAt := maxWidth
			for i := maxWidth - 1; i > maxWidth/2; i-- {
				if line[i] == ' ' || line[i] == ',' || line[i] == '{' || line[i] == '(' {
					breakAt = i + 1
					break
				}
			}
			result = append(result, line[:breakAt])
			line = "    " + strings.TrimLeft(line[breakAt:], " ")
		}
		if line != "" {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// wrapLines wraps each line independently, preserving original line breaks.
func wrapLines(text string, maxWidth int) string {
	if maxWidth <= 0 {
		maxWidth = 80
	}
	inputLines := strings.Split(text, "\n")
	var result []string
	for _, line := range inputLines {
		if len(line) <= maxWidth {
			result = append(result, line)
			continue
		}
		words := strings.Fields(line)
		if len(words) == 0 {
			result = append(result, line)
			continue
		}
		cur := words[0]
		for _, w := range words[1:] {
			if len(cur)+1+len(w) > maxWidth {
				result = append(result, cur)
				cur = "  " + w // indent continuation lines
			} else {
				cur += " " + w
			}
		}
		result = append(result, cur)
	}
	return strings.Join(result, "\n")
}

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRe.ReplaceAllString(s, "")
}

// formatExample styles example text to look like LeetCode's format.
// It highlights Input:, Output:, and Explanation: labels.
func formatExample(text string, maxWidth int) string {
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(colorWhite)
	valStyle := lipgloss.NewStyle().Foreground(colorSoftWhite)

	lines := strings.Split(text, "\n")
	var out []string
	for _, line := range lines {
		styled := false
		for _, label := range []string{"Input:", "Output:", "Explanation:"} {
			idx := strings.Index(line, label)
			if idx >= 0 {
				prefix := line[:idx]
				rest := line[idx+len(label):]
				out = append(out, prefix+labelStyle.Render(label)+valStyle.Render(rest))
				styled = true
				break
			}
		}
		if !styled {
			out = append(out, valStyle.Render(line))
		}
	}
	return strings.Join(out, "\n")
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func highlightGo(code string) string {
	keywords := []string{
		"func", "return", "if", "else", "for", "range", "var", "type",
		"struct", "interface", "package", "import", "const", "switch",
		"case", "default", "break", "continue", "defer", "go", "map",
		"chan", "select", "nil", "true", "false", "make", "append",
		"len", "cap", "new", "delete", "copy",
	}
	types := []string{
		"int", "string", "bool", "byte", "float64", "float32",
		"int64", "int32", "uint", "error", "rune",
	}

	kwStyle := lipgloss.NewStyle().Foreground(colorBlue).Bold(true)
	typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#06b6d4"))
	commentStyle := lipgloss.NewStyle().Foreground(colorDim)
	stringStyle := lipgloss.NewStyle().Foreground(colorGreen)
	numStyle := lipgloss.NewStyle().Foreground(colorOrange)
	funcNameStyle := lipgloss.NewStyle().Foreground(colorPurple)

	codeLines := strings.Split(code, "\n")
	var result []string

	for _, line := range codeLines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "//") {
			result = append(result, commentStyle.Render(line))
			continue
		}

		commentIdx := strings.Index(line, "//")
		codePart := line
		commentPart := ""
		if commentIdx >= 0 {
			inStr := false
			for ci := 0; ci < commentIdx; ci++ {
				if line[ci] == '"' {
					inStr = !inStr
				}
			}
			if !inStr {
				codePart = line[:commentIdx]
				commentPart = line[commentIdx:]
			}
		}

		var highlighted strings.Builder
		i := 0
		runes := []byte(codePart)
		for i < len(runes) {
			if runes[i] == '"' {
				j := i + 1
				for j < len(runes) && runes[j] != '"' {
					if runes[j] == '\\' {
						j++
					}
					j++
				}
				if j < len(runes) {
					j++
				}
				highlighted.WriteString(stringStyle.Render(string(runes[i:j])))
				i = j
				continue
			}
			if runes[i] == '\'' {
				j := i + 1
				for j < len(runes) && runes[j] != '\'' {
					if runes[j] == '\\' {
						j++
					}
					j++
				}
				if j < len(runes) {
					j++
				}
				highlighted.WriteString(stringStyle.Render(string(runes[i:j])))
				i = j
				continue
			}
			if runes[i] >= '0' && runes[i] <= '9' && (i == 0 || !isAlphaHL(runes[i-1])) {
				j := i
				for j < len(runes) && ((runes[j] >= '0' && runes[j] <= '9') || runes[j] == '.') {
					j++
				}
				if j < len(runes) && isAlphaHL(runes[j]) {
					highlighted.WriteByte(runes[i])
					i++
				} else {
					highlighted.WriteString(numStyle.Render(string(runes[i:j])))
					i = j
				}
				continue
			}
			if isAlphaHL(runes[i]) {
				j := i
				for j < len(runes) && (isAlphaHL(runes[j]) || (runes[j] >= '0' && runes[j] <= '9')) {
					j++
				}
				word := string(runes[i:j])
				if isKW(word, keywords) {
					if word == "func" && j < len(runes) && runes[j] == ' ' {
						highlighted.WriteString(kwStyle.Render(word))
						highlighted.WriteByte(' ')
						j++
						k := j
						for k < len(runes) && (isAlphaHL(runes[k]) || (runes[k] >= '0' && runes[k] <= '9')) {
							k++
						}
						if k > j {
							highlighted.WriteString(funcNameStyle.Render(string(runes[j:k])))
							j = k
						}
					} else {
						highlighted.WriteString(kwStyle.Render(word))
					}
				} else if isKW(word, types) {
					highlighted.WriteString(typeStyle.Render(word))
				} else {
					highlighted.WriteString(word)
				}
				i = j
				continue
			}
			highlighted.WriteByte(runes[i])
			i++
		}

		if commentPart != "" {
			highlighted.WriteString(commentStyle.Render(commentPart))
		}
		result = append(result, highlighted.String())
	}
	return strings.Join(result, "\n")
}

func isAlphaHL(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isKW(word string, list []string) bool {
	for _, kw := range list {
		if word == kw {
			return true
		}
	}
	return false
}

func main() {
	questionsFile := flag.String(
		"questions",
		"",
		"Path to JSON file with additional curated questions (slug + approaches)",
	)
	workspace := flag.String(
		"workspace",
		"",
		"Path to leetgo workspace directory (auto-detected if not set)",
	)
	showStats := flag.Bool("stats", false, "Show lifetime stats and exit")
	categoryFilter := flag.String("category", "", "Focus on a specific category (case-insensitive partial match)")
	listCategories := flag.Bool("categories", false, "List available categories and exit")
	flag.Parse()

	if _, err := exec.LookPath("leetgo"); err != nil {
		fmt.Fprintln(os.Stderr, "leetgo is not installed.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(
			os.Stderr,
			"This tool uses leetgo to fetch LeetCode problems and run test/submit.",
		)
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Setup:")
		fmt.Fprintln(os.Stderr, "  1. Install leetgo:")
		fmt.Fprintln(os.Stderr, "       go install github.com/j178/leetgo@latest")
		fmt.Fprintln(os.Stderr, "  2. Create a workspace:")
		fmt.Fprintln(os.Stderr, "       mkdir ~/leetcode && cd ~/leetcode && leetgo init")
		fmt.Fprintln(os.Stderr, "  3. Download the question database:")
		fmt.Fprintln(os.Stderr, "       leetgo cache update")
		fmt.Fprintln(os.Stderr, "  4. Run this tool:")
		fmt.Fprintln(os.Stderr, "       dsa-warmup -workspace ~/leetcode")
		os.Exit(1)
	}

	dbPath, err := findLeetgoDatabase()
	if err != nil {
		workdir := *workspace
		if workdir == "" {
			workdir = findLeetgoWorkspace()
		}
		if workdir != "" {
			fmt.Println("Question database not found. Running leetgo cache update...")
			cacheCmd := exec.Command("leetgo", "cache", "update")
			cacheCmd.Dir = workdir
			cacheCmd.Stdout = os.Stdout
			cacheCmd.Stderr = os.Stderr
			if err := cacheCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "\nleetgo cache update failed: %v\n", err)
				os.Exit(1)
			}
			dbPath, err = findLeetgoDatabase()
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Question database not found.")
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(
				os.Stderr,
				"The leetgo question cache is needed to load problem descriptions.",
			)
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "Setup:")
			fmt.Fprintln(os.Stderr, "  1. Create a leetgo workspace (if you haven't already):")
			fmt.Fprintln(os.Stderr, "       mkdir ~/leetcode && cd ~/leetcode && leetgo init")
			fmt.Fprintln(os.Stderr, "  2. Download the question database:")
			fmt.Fprintln(os.Stderr, "       cd ~/leetcode && leetgo cache update")
			fmt.Fprintln(os.Stderr, "  3. Run this tool:")
			fmt.Fprintln(os.Stderr, "       dsa-warmup -workspace ~/leetcode")
			os.Exit(1)
		}
	}

	workdir := *workspace
	if workdir == "" {
		workdir = findLeetgoWorkspace()
	}
	if workdir == "" {
		fmt.Fprintln(os.Stderr, "Warning: no leetgo workspace found.")
		fmt.Fprintln(os.Stderr, "  [t] Try It, [T] Test, and [S] Submit will be disabled.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "  To enable, either:")
		fmt.Fprintln(os.Stderr, "    - Run from a directory containing leetgo.yaml")
		fmt.Fprintln(os.Stderr, "    - Pass -workspace <path> to your leetgo workspace")
		fmt.Fprintln(
			os.Stderr,
			"    - Create one: mkdir ~/leetcode && cd ~/leetcode && leetgo init",
		)
		fmt.Fprintln(os.Stderr, "")
	}

	if workdir != "" {
		fmt.Print("Updating leetgo cache...")
		cacheCmd := exec.Command("leetgo", "cache", "update")
		cacheCmd.Dir = workdir
		if out, err := cacheCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, " failed: %s\n", strings.TrimSpace(string(out)))
		} else {
			fmt.Println(" done")
		}
	}

	allCurated := append(curatedBank, curatedBankExtended...)
	questions, err := loadQuestions(dbPath, allCurated)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading questions: %v\n", err)
		os.Exit(1)
	}

	if *questionsFile != "" {
		extra, err := loadQuestionsFromJSON(dbPath, *questionsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to load %s: %v\n", *questionsFile, err)
		} else {
			questions = append(questions, extra...)
			fmt.Printf("Loaded %d additional questions from %s\n", len(extra), *questionsFile)
		}
	}

	// Collect unique categories from loaded questions.
	categorySet := make(map[string]int)
	for _, q := range questions {
		categorySet[q.Category]++
	}
	var categoryNames []string
	for name := range categorySet {
		categoryNames = append(categoryNames, name)
	}
	sort.Strings(categoryNames)

	if *listCategories {
		fmt.Println("Available categories:")
		for _, name := range categoryNames {
			fmt.Printf("  %-30s (%d problems)\n", name, categorySet[name])
		}
		return
	}

	// Filter by category if requested.
	activeCategoryFilter := ""
	if *categoryFilter != "" {
		needle := strings.ToLower(*categoryFilter)
		var matched []string
		for _, name := range categoryNames {
			if strings.ToLower(name) == needle {
				matched = []string{name}
				break
			}
			if strings.Contains(strings.ToLower(name), needle) {
				matched = append(matched, name)
			}
		}
		if len(matched) == 0 {
			fmt.Fprintf(os.Stderr, "No category matching %q. Available categories:\n", *categoryFilter)
			for _, name := range categoryNames {
				fmt.Fprintf(os.Stderr, "  %-30s (%d problems)\n", name, categorySet[name])
			}
			os.Exit(1)
		}
		if len(matched) > 1 {
			fmt.Fprintf(os.Stderr, "Ambiguous category %q matches multiple:\n", *categoryFilter)
			for _, name := range matched {
				fmt.Fprintf(os.Stderr, "  %-30s (%d problems)\n", name, categorySet[name])
			}
			os.Exit(1)
		}
		activeCategoryFilter = matched[0]
		var filtered []Question
		for _, q := range questions {
			if q.Category == activeCategoryFilter {
				filtered = append(filtered, q)
			}
		}
		questions = filtered
		fmt.Printf("Focused on %q: %d problems\n", activeCategoryFilter, len(questions))
	}

	reviewLog := LoadReviewLog()

	if *showStats {
		fmt.Print(reviewLog.LifetimeStats(questions))
		return
	}

	now := time.Now()
	dueCount := 0
	reviewedCount := 0
	slugSet := make(map[string]bool, len(questions))
	for _, q := range questions {
		slugSet[q.LeetcodeSlug] = true
	}
	for slug, pr := range reviewLog.Reviews {
		if !slugSet[slug] {
			continue
		}
		reviewedCount++
		if now.After(pr.NextReviewAt) {
			dueCount++
		}
	}
	newCount := len(questions) - reviewedCount

	fmt.Printf("Loaded %d questions from leetgo database\n", len(questions))
	fmt.Printf("Due: %d | New: %d | Reviewed: %d\n", dueCount, newCount, reviewedCount)
	if workdir != "" {
		fmt.Printf("leetgo workspace: %s\n", workdir)
	}
	time.Sleep(1 * time.Second)

	m := model{
		questions:      questions,
		selected:       -1,
		timer:          defaultTimer,
		width:          80,
		height:         40,
		leetgoWorkdir:  workdir,
		reviewLog:      reviewLog,
		sessionSeen:    make(map[int]bool),
		submitResult:   SubmitNone,
		categoryFilter: activeCategoryFilter,
	}
	m.unseen = make([]int, len(m.questions))
	for i := range m.questions {
		m.unseen[i] = i
	}
	m.pickQuestion()
	m.timerRunning = true

	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if fm, ok := finalModel.(model); ok && fm.sessionReport != "" {
		fmt.Print(fm.sessionReport)
	}
}
