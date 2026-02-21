package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ReviewLog is the persistent store for all problem review history.
type ReviewLog struct {
	Reviews   map[string]*ProblemReview `json:"reviews"`
	StartedAt time.Time                 `json:"started_at"` // when first review was recorded
}

// ProblemReview tracks spaced repetition state for one problem.
type ProblemReview struct {
	ProblemID     string    `json:"slug"` // JSON key kept as "slug" for backward compat with existing reviews.json
	LastReviewAt  time.Time `json:"last_review_at"`
	NextReviewAt  time.Time `json:"next_review_at"`
	Interval      float64   `json:"interval"`      // days until next review
	EaseFactor    float64   `json:"ease_factor"`   // SM-2 ease (starts 2.5)
	Repetitions   int       `json:"repetitions"`   // consecutive good reps
	ApproachHist  []int     `json:"approach_hist"` // last 10 approach ratings (0=optimal..3=wrong)
	SubmitHist    []int     `json:"submit_hist"`   // last 10 submit results (0=accepted, 1=wrong, 2=error, 3=none)
	TotalReviews  int       `json:"total_reviews"`
	TotalOptimal  int       `json:"total_optimal"`  // times picked optimal approach
	TotalAccepted int       `json:"total_accepted"` // times submit was accepted
	CodingTimes   []int     `json:"coding_times"`   // last 10 coding durations in seconds (0 = didn't code)
}

const (
	SubmitAccepted = 0
	SubmitWrong    = 1
	SubmitError    = 2
	SubmitNone     = 3 // user didn't submit
)

// computeQuality derives a 0-5 SM-2 quality score from approach rating and submit result.
//
//	Optimal + Accepted = 5
//	Optimal + no submit = 4
//	Plausible + Accepted = 3
//	Plausible + no submit/wrong = 2
//	Suboptimal = 1
//	Wrong = 0
func computeQuality(approach Rating, submitResult int) int {
	switch approach {
	case Optimal:
		if submitResult == SubmitAccepted {
			return 5
		}
		return 4
	case Plausible:
		if submitResult == SubmitAccepted {
			return 3
		}
		return 2
	case Suboptimal:
		return 1
	default: // Wrong
		return 0
	}
}

func (pr *ProblemReview) updateSM2(quality int) {
	ef := pr.EaseFactor
	q := float64(quality)
	ef = ef + (0.1 - (5-q)*(0.08+(5-q)*0.02))
	if ef < 1.3 {
		ef = 1.3
	}
	pr.EaseFactor = ef

	if quality >= 3 {
		pr.Repetitions++
		switch pr.Repetitions {
		case 1:
			pr.Interval = 1
		case 2:
			pr.Interval = 6
		default:
			pr.Interval = pr.Interval * ef
		}
	} else {
		pr.Repetitions = 0
		pr.Interval = 1
	}

	now := time.Now()
	pr.LastReviewAt = now
	pr.NextReviewAt = now.Add(time.Duration(pr.Interval*24) * time.Hour)
}

// RecordReview records a single review event for a problem.
func (rl *ReviewLog) RecordReview(slug string, approach Rating, submitResult int, codingTime int) {
	pr, ok := rl.Reviews[slug]
	if !ok {
		pr = &ProblemReview{
			ProblemID:  slug,
			EaseFactor: 2.5,
			Interval:   0,
		}
		rl.Reviews[slug] = pr
	}

	pr.TotalReviews++
	if approach == Optimal {
		pr.TotalOptimal++
	}
	if submitResult == SubmitAccepted {
		pr.TotalAccepted++
	}

	pr.ApproachHist = append(pr.ApproachHist, int(approach))
	if len(pr.ApproachHist) > 10 {
		pr.ApproachHist = pr.ApproachHist[len(pr.ApproachHist)-10:]
	}

	pr.SubmitHist = append(pr.SubmitHist, submitResult)
	if len(pr.SubmitHist) > 10 {
		pr.SubmitHist = pr.SubmitHist[len(pr.SubmitHist)-10:]
	}

	pr.CodingTimes = append(pr.CodingTimes, codingTime)
	if len(pr.CodingTimes) > 10 {
		pr.CodingTimes = pr.CodingTimes[len(pr.CodingTimes)-10:]
	}

	quality := computeQuality(approach, submitResult)
	pr.updateSM2(quality)
}

// PickNextQuestion selects the next question using SRS priority:
// due (most overdue first), then new, then not-yet-due.
// It returns the index into questions, or -1 if none remain.
func (rl *ReviewLog) PickNextQuestion(questions []Question, sessionSeen map[int]bool) int {
	now := time.Now()

	type candidate struct {
		idx      int
		priority int     // 0=overdue, 1=new, 2=not-yet-due
		overdue  float64 // hours overdue (for sorting)
	}

	var candidates []candidate

	for i, q := range questions {
		if sessionSeen[i] {
			continue
		}

		pr, reviewed := rl.Reviews[q.ProblemID]
		if !reviewed {
			candidates = append(candidates, candidate{idx: i, priority: 1, overdue: 0})
		} else if now.After(pr.NextReviewAt) {
			overdue := now.Sub(pr.NextReviewAt).Hours()
			candidates = append(candidates, candidate{idx: i, priority: 0, overdue: overdue})
		} else {
			candidates = append(candidates, candidate{idx: i, priority: 2, overdue: -pr.NextReviewAt.Sub(now).Hours()})
		}
	}

	if len(candidates) == 0 {
		return -1
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].priority != candidates[j].priority {
			return candidates[i].priority < candidates[j].priority
		}
		return candidates[i].overdue > candidates[j].overdue
	})

	return candidates[0].idx
}

func reviewLogPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "warmup", "reviews.json")
}

// LoadReviewLog loads the review log from disk, or returns a fresh one.
func LoadReviewLog() *ReviewLog {
	rl := &ReviewLog{
		Reviews:   make(map[string]*ProblemReview),
		StartedAt: time.Now(),
	}

	data, err := os.ReadFile(reviewLogPath())
	if err != nil {
		return rl
	}

	if err := json.Unmarshal(data, rl); err != nil {
		return rl
	}

	if rl.Reviews == nil {
		rl.Reviews = make(map[string]*ProblemReview)
	}

	return rl
}

// Save writes the review log to disk.
func (rl *ReviewLog) Save() error {
	path := reviewLogPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(rl, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ParseSubmitResult extracts the submit verdict from provider output.
func ParseSubmitResult(output string) int {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "accepted") {
		return SubmitAccepted
	}
	if strings.Contains(lower, "wrong answer") {
		return SubmitWrong
	}
	if strings.Contains(lower, "compile error") || strings.Contains(lower, "runtime error") ||
		strings.Contains(lower, "time limit") || strings.Contains(lower, "memory limit") {
		return SubmitError
	}
	return SubmitNone
}

// SessionReport generates an end-of-session summary string.
func (rl *ReviewLog) SessionReport(sessionReviews []sessionEntry, totalQuestions int) string {
	if len(sessionReviews) == 0 {
		return "No problems reviewed this session."
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString("═══════════════════════════════════════════\n")
	b.WriteString("            SESSION REPORT\n")
	b.WriteString("═══════════════════════════════════════════\n\n")

	total := len(sessionReviews)
	optimal := 0
	plausible := 0
	submitted := 0
	accepted := 0
	coded := 0
	totalCodingSecs := 0
	for _, e := range sessionReviews {
		if e.approach == Optimal {
			optimal++
		}
		if e.approach == Plausible {
			plausible++
		}
		if e.submitResult != SubmitNone {
			submitted++
		}
		if e.submitResult == SubmitAccepted {
			accepted++
		}
		if e.codingTime > 0 {
			coded++
			totalCodingSecs += e.codingTime
		}
	}

	b.WriteString(fmt.Sprintf("  Problems reviewed:  %d\n", total))
	b.WriteString(fmt.Sprintf("  Optimal approach:   %d/%d (%.0f%%)\n", optimal, total, pct(optimal, total)))
	b.WriteString(fmt.Sprintf("  Plausible+:         %d/%d (%.0f%%)\n", optimal+plausible, total, pct(optimal+plausible, total)))
	if submitted > 0 {
		b.WriteString(fmt.Sprintf("  Submitted:          %d\n", submitted))
		b.WriteString(fmt.Sprintf("  Accepted:           %d/%d (%.0f%%)\n", accepted, submitted, pct(accepted, submitted)))
	}
	if coded > 0 {
		avgSecs := totalCodingSecs / coded
		b.WriteString(fmt.Sprintf("  Coded:              %d problems\n", coded))
		b.WriteString(fmt.Sprintf("  Avg coding time:    %d:%02d\n", avgSecs/60, avgSecs%60))
		b.WriteString(fmt.Sprintf("  Total coding time:  %d:%02d\n", totalCodingSecs/60, totalCodingSecs%60))
	}
	b.WriteString("\n")

	b.WriteString("  Problem                                    Approach   Submit     Time\n")
	b.WriteString("  ────────────────────────────────────────────────────────────────────\n")
	for _, e := range sessionReviews {
		name := e.title
		if len(name) > 40 {
			name = name[:37] + "..."
		}
		approachStr := e.approach.String()
		submitStr := "—"
		switch e.submitResult {
		case SubmitAccepted:
			submitStr = "Accepted"
		case SubmitWrong:
			submitStr = "Wrong"
		case SubmitError:
			submitStr = "Error"
		}
		timeStr := "—"
		if e.codingTime > 0 {
			timeStr = fmt.Sprintf("%d:%02d", e.codingTime/60, e.codingTime%60)
		}
		b.WriteString(fmt.Sprintf("  %-40s %-10s %-10s %s\n", name, approachStr, submitStr, timeStr))
	}

	b.WriteString("\n")
	now := time.Now()
	due := 0
	dueToday := 0
	mastered := 0
	for _, pr := range rl.Reviews {
		if now.After(pr.NextReviewAt) {
			due++
		}
		if pr.NextReviewAt.Before(now.Add(24 * time.Hour)) {
			dueToday++
		}
		if pr.Interval >= 21 && pr.Repetitions >= 3 {
			mastered++
		}
	}
	newCount := totalQuestions - len(rl.Reviews)
	if newCount < 0 {
		newCount = 0
	}
	b.WriteString(fmt.Sprintf("  Due now:     %d\n", due))
	b.WriteString(fmt.Sprintf("  Due today:   %d\n", dueToday))
	b.WriteString(fmt.Sprintf("  New unseen:  %d\n", newCount))
	b.WriteString(fmt.Sprintf("  Mastered:    %d\n", mastered))

	b.WriteString("\n═══════════════════════════════════════════\n")

	return b.String()
}

// LifetimeStats generates a comprehensive stats string for --stats mode.
func (rl *ReviewLog) LifetimeStats(questions []Question) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString("═══════════════════════════════════════════\n")
	b.WriteString("           DSA WARMUP STATS\n")
	b.WriteString("═══════════════════════════════════════════\n\n")

	totalQuestions := len(questions)
	totalReviewed := 0
	totalReviews := 0
	totalOptimal := 0
	totalAccepted := 0
	totalSubmitted := 0

	now := time.Now()
	dueNow := 0
	dueToday := 0
	mastered := 0

	// Build a set of problem IDs from the (potentially filtered) questions list
	// so all stats are scoped correctly when --category is active.
	questionIDs := make(map[string]bool, len(questions))
	for _, q := range questions {
		questionIDs[q.ProblemID] = true
	}

	type catStats struct {
		reviewed int
		total    int
		optimal  int
		reviews  int
		accepted int
	}
	cats := make(map[string]*catStats)

	for _, q := range questions {
		cat := q.Category
		if cats[cat] == nil {
			cats[cat] = &catStats{}
		}
		cats[cat].total++
	}

	type weakEntry struct {
		slug string
		ef   float64
	}
	var weakest []weakEntry

	for _, pr := range rl.Reviews {
		// Only count reviews for problems in the current question set.
		if !questionIDs[pr.ProblemID] {
			continue
		}

		totalReviewed++
		totalReviews += pr.TotalReviews
		totalOptimal += pr.TotalOptimal
		totalAccepted += pr.TotalAccepted
		for _, s := range pr.SubmitHist {
			if s != SubmitNone {
				totalSubmitted++
			}
		}

		if now.After(pr.NextReviewAt) {
			dueNow++
		}
		if pr.NextReviewAt.Before(now.Add(24 * time.Hour)) {
			dueToday++
		}
		if pr.Interval >= 21 && pr.Repetitions >= 3 {
			mastered++
		}

		for _, q := range questions {
			if q.ProblemID == pr.ProblemID {
				cat := q.Category
				if cats[cat] == nil {
					cats[cat] = &catStats{}
				}
				cats[cat].reviewed++
				cats[cat].optimal += pr.TotalOptimal
				cats[cat].reviews += pr.TotalReviews
				cats[cat].accepted += pr.TotalAccepted
				break
			}
		}

		if pr.TotalReviews >= 2 {
			weakest = append(weakest, weakEntry{pr.ProblemID, pr.EaseFactor})
		}
	}

	b.WriteString(fmt.Sprintf("  Total problems:   %d\n", totalQuestions))
	b.WriteString(fmt.Sprintf("  Reviewed:         %d/%d (%.0f%%)\n", totalReviewed, totalQuestions, pct(totalReviewed, totalQuestions)))
	b.WriteString(fmt.Sprintf("  Total reviews:    %d\n", totalReviews))
	b.WriteString(fmt.Sprintf("  Mastered:         %d (interval >= 21 days)\n", mastered))
	b.WriteString(fmt.Sprintf("  Due now:          %d\n", dueNow))
	b.WriteString(fmt.Sprintf("  Due today:        %d\n", dueToday))
	b.WriteString(fmt.Sprintf("  New (unseen):     %d\n", totalQuestions-totalReviewed))
	b.WriteString("\n")

	if totalReviews > 0 {
		b.WriteString(fmt.Sprintf("  Approach accuracy (optimal): %.0f%%\n", pct(totalOptimal, totalReviews)))
	}
	if totalSubmitted > 0 {
		b.WriteString(fmt.Sprintf("  Submit accuracy:             %.0f%%\n", pct(totalAccepted, totalSubmitted)))
	}
	b.WriteString("\n")

	b.WriteString("  Category                   Seen  Optimal%%  Reviews\n")
	b.WriteString("  ─────────────────────────────────────────────────\n")

	var catNames []string
	for name := range cats {
		catNames = append(catNames, name)
	}
	sort.Strings(catNames)

	for _, name := range catNames {
		cs := cats[name]
		optPct := 0.0
		if cs.reviews > 0 {
			optPct = pct(cs.optimal, cs.reviews)
		}
		displayName := name
		if len(displayName) > 25 {
			displayName = displayName[:22] + "..."
		}
		b.WriteString(fmt.Sprintf("  %-25s %2d/%-3d  %5.0f%%    %d\n",
			displayName, cs.reviewed, cs.total, optPct, cs.reviews))
	}

	if len(weakest) > 0 {
		sort.Slice(weakest, func(i, j int) bool {
			return weakest[i].ef < weakest[j].ef
		})
		b.WriteString("\n  Weakest problems (lowest ease factor):\n")
		shown := 10
		if len(weakest) < shown {
			shown = len(weakest)
		}
		for i := 0; i < shown; i++ {
			b.WriteString(fmt.Sprintf("    %.2f  %s\n", weakest[i].ef, weakest[i].slug))
		}
	}

	// Coding speed section — show problems with timing history
	type codingEntry struct {
		slug    string
		times   []int // non-zero coding times
		latest  int
		fastest int
		avg     int
	}
	var codingEntries []codingEntry
	totalCodingSessions := 0
	for _, pr := range rl.Reviews {
		if !questionIDs[pr.ProblemID] {
			continue
		}
		var nonZero []int
		for _, t := range pr.CodingTimes {
			if t > 0 {
				nonZero = append(nonZero, t)
			}
		}
		if len(nonZero) == 0 {
			continue
		}
		totalCodingSessions += len(nonZero)
		fastest := nonZero[0]
		sum := 0
		for _, t := range nonZero {
			sum += t
			if t < fastest {
				fastest = t
			}
		}
		codingEntries = append(codingEntries, codingEntry{
			slug:    pr.ProblemID,
			times:   nonZero,
			latest:  nonZero[len(nonZero)-1],
			fastest: fastest,
			avg:     sum / len(nonZero),
		})
	}
	if len(codingEntries) > 0 {
		sort.Slice(codingEntries, func(i, j int) bool {
			return codingEntries[i].slug < codingEntries[j].slug
		})
		b.WriteString("\n  Coding Speed:\n")
		b.WriteString(fmt.Sprintf("  Problems coded:     %d\n", len(codingEntries)))
		b.WriteString(fmt.Sprintf("  Total sessions:     %d\n\n", totalCodingSessions))
		b.WriteString("  Problem                          Latest   Best     Avg    Tries\n")
		b.WriteString("  ──────────────────────────────────────────────────────────────\n")
		for _, ce := range codingEntries {
			name := ce.slug
			if len(name) > 30 {
				name = name[:27] + "..."
			}
			b.WriteString(fmt.Sprintf("  %-30s  %3d:%02d   %3d:%02d   %3d:%02d    %d\n",
				name,
				ce.latest/60, ce.latest%60,
				ce.fastest/60, ce.fastest%60,
				ce.avg/60, ce.avg%60,
				len(ce.times)))
		}
	}

	b.WriteString("\n═══════════════════════════════════════════\n")
	return b.String()
}

type sessionEntry struct {
	title        string
	slug         string
	approach     Rating
	submitResult int
	codingTime   int // seconds spent coding (0 = didn't code)
}

func pct(num, denom int) float64 {
	if denom == 0 {
		return 0
	}
	return math.Round(float64(num) / float64(denom) * 100)
}
