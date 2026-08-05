package ui

import (
	"testing"
	"time"

	"charm.land/bubbles/v2/table"

	"github.com/taigrr/gico/types"
)

func TestCommitRowsForDay(t *testing.T) {
	commits := [][]types.Commit{
		nil,
		{
			{
				TimeStamp: time.Date(2026, time.August, 5, 9, 30, 0, 0, time.Local),
				Repo:      "/home/tai/code/foss/gico",
				Author:    types.Author{Name: "Tai"},
				Message:   "fix the graph",
			},
		},
	}

	got := commitRowsForDay(commits, 1)
	want := table.Row{"09:30AM", "gico", "Tai", "fix the graph"}
	if len(got) != 1 {
		t.Fatalf("commitRowsForDay returned %d rows, want 1", len(got))
	}
	for i := range want {
		if got[0][i] != want[i] {
			t.Fatalf("commitRowsForDay row[%d] = %q, want %q", i, got[0][i], want[i])
		}
	}
}

func TestCommitRowsForDayOutOfRange(t *testing.T) {
	if got := commitRowsForDay(nil, 0); got != nil {
		t.Fatalf("commitRowsForDay(nil, 0) = %v, want nil", got)
	}
	if got := commitRowsForDay(make([][]types.Commit, 1), -1); got != nil {
		t.Fatalf("commitRowsForDay(_, -1) = %v, want nil", got)
	}
	if got := commitRowsForDay(make([][]types.Commit, 1), 1); got != nil {
		t.Fatalf("commitRowsForDay(_, 1) = %v, want nil", got)
	}
}

func TestCommitLogViewOutOfRange(t *testing.T) {
	log := CommitLog{
		Commits: make([][]types.Commit, 1),
		YearDay: 1,
	}

	if got := log.View(); got != "No commits to display" {
		t.Fatalf("CommitLog.View() = %q, want no commits message", got)
	}
}
