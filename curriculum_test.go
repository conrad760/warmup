package main

import (
	"testing"
)

func TestCurriculumDaysExist(t *testing.T) {
	// Verify all 18 topic days are present.
	for day := 1; day <= 18; day++ {
		if _, ok := curriculumDays[day]; !ok {
			t.Errorf("missing curriculum entry for day %d", day)
		}
	}
}

func TestCurriculumDaysHaveProblems(t *testing.T) {
	for day, cd := range curriculumDays {
		if len(cd.Problems) == 0 {
			t.Errorf("day %d (%s) has no problems", day, cd.Topic)
		}
		if cd.Day != day {
			t.Errorf("day %d has mismatched Day field: %d", day, cd.Day)
		}
		if cd.Topic == "" {
			t.Errorf("day %d has empty topic", day)
		}
	}
}

func TestCurriculumSlugsExistInCuratedBank(t *testing.T) {
	// Build a set of all curated problem IDs.
	allCurated := make(map[string]bool)
	for _, q := range curatedBank {
		allCurated[q.ProblemID] = true
	}
	for _, q := range curatedBankExtended {
		allCurated[q.ProblemID] = true
	}

	for day, cd := range curriculumDays {
		for _, slug := range cd.Problems {
			if !allCurated[slug] {
				t.Errorf("day %d (%s): slug %q not found in curated bank", day, cd.Topic, slug)
			}
		}
	}
}

func TestCurriculumNoDuplicateSlugsPerDay(t *testing.T) {
	for day, cd := range curriculumDays {
		seen := make(map[string]bool)
		for _, slug := range cd.Problems {
			if seen[slug] {
				t.Errorf("day %d (%s): duplicate slug %q", day, cd.Topic, slug)
			}
			seen[slug] = true
		}
	}
}

func TestCurriculumProblems(t *testing.T) {
	// Valid day.
	cd, err := CurriculumProblems(1)
	if err != nil {
		t.Fatalf("unexpected error for day 1: %v", err)
	}
	if cd.Day != 1 {
		t.Errorf("expected day 1, got %d", cd.Day)
	}
	if len(cd.Problems) == 0 {
		t.Error("expected non-empty problem list for day 1")
	}

	// Invalid days.
	for _, day := range []int{0, -1, 19, 22, 100} {
		_, err := CurriculumProblems(day)
		if err == nil {
			t.Errorf("expected error for day %d, got nil", day)
		}
	}
}

func TestListCurriculumDays(t *testing.T) {
	days := ListCurriculumDays()
	if len(days) != len(curriculumDays) {
		t.Errorf("expected %d days, got %d", len(curriculumDays), len(days))
	}
	// Verify sorted order.
	for i := 1; i < len(days); i++ {
		if days[i].Day <= days[i-1].Day {
			t.Errorf("days not sorted: day %d at index %d, day %d at index %d",
				days[i-1].Day, i-1, days[i].Day, i)
		}
	}
}
