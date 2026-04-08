package commits

import (
	"testing"
	"time"

	"github.com/taigrr/gico/types"
)

func makeCommit(name, email string, ts time.Time) types.Commit {
	return types.Commit{
		Author:    types.Author{Name: name, Email: email},
		TimeStamp: ts,
		Hash:      "abc123",
		Repo:      "/test/repo",
	}
}

func TestCommitSetFilterByYear(t *testing.T) {
	cs := CommitSet{
		Commits: []types.Commit{
			makeCommit("a", "a@x.com", time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)),
			makeCommit("b", "b@x.com", time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)),
			makeCommit("c", "c@x.com", time.Date(2024, 12, 31, 23, 59, 0, 0, time.UTC)),
			makeCommit("d", "d@x.com", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		},
	}

	filtered := cs.FilterByYear(2024)
	if filtered.Year != 2024 {
		t.Errorf("expected Year=2024, got %d", filtered.Year)
	}
	if len(filtered.Commits) != 2 {
		t.Errorf("expected 2 commits for 2024, got %d", len(filtered.Commits))
	}

	filtered = cs.FilterByYear(2025)
	if len(filtered.Commits) != 1 {
		t.Errorf("expected 1 commit for 2025, got %d", len(filtered.Commits))
	}

	filtered = cs.FilterByYear(2000)
	if len(filtered.Commits) != 0 {
		t.Errorf("expected 0 commits for 2000, got %d", len(filtered.Commits))
	}
}

func TestCommitSetFilterByAuthorRegex(t *testing.T) {
	cs := CommitSet{
		Commits: []types.Commit{
			makeCommit("Alice", "alice@example.com", time.Now()),
			makeCommit("Bob", "bob@example.com", time.Now()),
			makeCommit("Charlie", "charlie@example.com", time.Now()),
		},
	}

	// Filter by name regex
	filtered, err := cs.FilterByAuthorRegex([]string{"^Ali"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered.Commits) != 1 {
		t.Errorf("expected 1 commit matching ^Ali, got %d", len(filtered.Commits))
	}

	// Filter by email regex
	filtered, err = cs.FilterByAuthorRegex([]string{"bob@"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered.Commits) != 1 {
		t.Errorf("expected 1 commit matching bob@, got %d", len(filtered.Commits))
	}

	// Multiple patterns
	filtered, err = cs.FilterByAuthorRegex([]string{"Alice", "Charlie"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filtered.Commits) != 2 {
		t.Errorf("expected 2 commits matching Alice|Charlie, got %d", len(filtered.Commits))
	}

	// Invalid regex
	_, err = cs.FilterByAuthorRegex([]string{"[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestCommitSetToYearFreq(t *testing.T) {
	// Regular year (2025)
	cs := CommitSet{
		Year: 2025,
		Commits: []types.Commit{
			makeCommit("a", "a@x.com", time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)),
			makeCommit("a", "a@x.com", time.Date(2025, 1, 1, 14, 0, 0, 0, time.UTC)),
			makeCommit("b", "b@x.com", time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC)),
		},
	}

	freq := cs.ToYearFreq()
	if len(freq) != 365 {
		t.Errorf("expected 365 days for 2025, got %d", len(freq))
	}
	if freq[0] != 2 {
		t.Errorf("expected 2 commits on Jan 1, got %d", freq[0])
	}
	// March 15 = day 74 (31+28+15=74, index 73)
	if freq[73] != 1 {
		t.Errorf("expected 1 commit on Mar 15, got %d", freq[73])
	}

	// Leap year (2024)
	csLeap := CommitSet{
		Year: 2024,
		Commits: []types.Commit{
			makeCommit("a", "a@x.com", time.Date(2024, 12, 31, 10, 0, 0, 0, time.UTC)),
		},
	}
	freqLeap := csLeap.ToYearFreq()
	if len(freqLeap) != 366 {
		t.Errorf("expected 366 days for 2024, got %d", len(freqLeap))
	}
}

func TestYearFreqFromChan(t *testing.T) {
	cc := make(chan types.Commit, 5)
	cc <- makeCommit("a", "a@x.com", time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC))
	cc <- makeCommit("a", "a@x.com", time.Date(2025, 1, 1, 14, 0, 0, 0, time.UTC))
	cc <- makeCommit("b", "b@x.com", time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC))
	close(cc)

	freq := YearFreqFromChan(cc, 2025)
	if len(freq) != 365 {
		t.Errorf("expected 365, got %d", len(freq))
	}
	if freq[0] != 2 {
		t.Errorf("expected 2 on Jan 1, got %d", freq[0])
	}
}

func TestFreqFromChan(t *testing.T) {
	cc := make(chan types.Commit, 5)
	cc <- makeCommit("a", "a@x.com", time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC))
	cc <- makeCommit("b", "b@x.com", time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)) // wrong year
	cc <- makeCommit("c", "c@x.com", time.Date(2025, 12, 31, 23, 0, 0, 0, time.UTC))
	close(cc)

	freq := FreqFromChan(cc, 2025)
	if len(freq) != 365 {
		t.Errorf("expected 365, got %d", len(freq))
	}
	if freq[0] != 1 {
		t.Errorf("expected 1 on Jan 1 (skip wrong year), got %d", freq[0])
	}
	if freq[364] != 1 {
		t.Errorf("expected 1 on Dec 31, got %d", freq[364])
	}
}

func TestFilterCChanByYear(t *testing.T) {
	in := make(chan types.Commit, 5)
	in <- makeCommit("a", "a@x.com", time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC))
	in <- makeCommit("b", "b@x.com", time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC))
	in <- makeCommit("c", "c@x.com", time.Date(2025, 7, 4, 10, 0, 0, 0, time.UTC))
	close(in)

	out := FilterCChanByYear(in, 2025)
	count := 0
	for range out {
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 commits for 2025, got %d", count)
	}
}

func TestFilterCChanByAuthor(t *testing.T) {
	in := make(chan types.Commit, 5)
	in <- makeCommit("Alice", "alice@x.com", time.Now())
	in <- makeCommit("Bob", "bob@x.com", time.Now())
	in <- makeCommit("Charlie", "charlie@x.com", time.Now())
	close(in)

	out, err := FilterCChanByAuthor(in, []string{"Alice", "Charlie"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count := 0
	for range out {
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 commits, got %d", count)
	}
}

func TestOpenRepoNonExistent(t *testing.T) {
	_, err := OpenRepo("/nonexistent/path")
	if err == nil {
		t.Error("expected error for nonexistent path")
	}
}

func TestOpenRepoNonDirectory(t *testing.T) {
	// Use a known file (not a directory)
	_, err := OpenRepo("/dev/null")
	if err == nil {
		t.Error("expected error for non-directory path")
	}
}
