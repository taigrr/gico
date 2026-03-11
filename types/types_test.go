package types

import (
	"strings"
	"testing"
	"time"
)

func TestNewDataSet(t *testing.T) {
	ds := NewDataSet()
	if ds == nil {
		t.Fatal("NewDataSet returned nil")
	}
	if len(ds) != 0 {
		t.Fatalf("expected empty DataSet, got %d entries", len(ds))
	}
}

func TestNewCommit(t *testing.T) {
	c := NewCommit("Alice", "alice@example.com", "initial commit", "myrepo", "/tmp/myrepo", 10, 2, 3)
	if c.Author.Name != "Alice" {
		t.Errorf("expected author name Alice, got %s", c.Author.Name)
	}
	if c.Author.Email != "alice@example.com" {
		t.Errorf("expected author email alice@example.com, got %s", c.Author.Email)
	}
	if c.Message != "initial commit" {
		t.Errorf("expected message 'initial commit', got %s", c.Message)
	}
	if c.Added != 10 {
		t.Errorf("expected Added=10, got %d", c.Added)
	}
	if c.Deleted != 2 {
		t.Errorf("expected Deleted=2, got %d", c.Deleted)
	}
	if c.FilesChanged != 3 {
		t.Errorf("expected FilesChanged=3, got %d", c.FilesChanged)
	}
	if c.Repo != "myrepo" {
		t.Errorf("expected Repo=myrepo, got %s", c.Repo)
	}
	if c.Path != "/tmp/myrepo" {
		t.Errorf("expected Path=/tmp/myrepo, got %s", c.Path)
	}
	if c.TimeStamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestCommitString(t *testing.T) {
	ts := time.Date(2025, 6, 15, 14, 30, 0, 0, time.UTC)
	c := Commit{
		TimeStamp: ts,
		Author:    Author{Name: "Bob", Email: "bob@example.com"},
		Repo:      "testrepo",
		Message:   "fix bug",
	}
	s := c.String()
	if !strings.Contains(s, "testrepo") {
		t.Errorf("expected string to contain repo name, got: %s", s)
	}
	if !strings.Contains(s, "fix bug") {
		t.Errorf("expected string to contain message, got: %s", s)
	}
}

func TestFreqMerge(t *testing.T) {
	tests := []struct {
		name string
		a, b Freq
		want Freq
	}{
		{
			name: "equal length",
			a:    Freq{1, 2, 3},
			b:    Freq{4, 5, 6},
			want: Freq{5, 7, 9},
		},
		{
			name: "a shorter than b",
			a:    Freq{1, 2},
			b:    Freq{3, 4, 5},
			want: Freq{4, 6, 5},
		},
		{
			name: "b shorter than a",
			a:    Freq{1, 2, 3},
			b:    Freq{4, 5},
			want: Freq{5, 7, 3},
		},
		{
			name: "empty a",
			a:    Freq{},
			b:    Freq{1, 2, 3},
			want: Freq{1, 2, 3},
		},
		{
			name: "empty b",
			a:    Freq{1, 2, 3},
			b:    Freq{},
			want: Freq{1, 2, 3},
		},
		{
			name: "both empty",
			a:    Freq{},
			b:    Freq{},
			want: Freq{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Merge(tt.b)
			if len(got) != len(tt.want) {
				t.Fatalf("expected length %d, got %d", len(tt.want), len(got))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("index %d: expected %d, got %d", i, tt.want[i], got[i])
				}
			}
		})
	}
}

func TestDataSetOperations(t *testing.T) {
	ds := NewDataSet()
	now := time.Now().Truncate(24 * time.Hour)
	ds[now] = WorkDay{
		Day:   now,
		Count: 5,
	}
	if wd, ok := ds[now]; !ok {
		t.Fatal("expected to find workday for today")
	} else if wd.Count != 5 {
		t.Errorf("expected count 5, got %d", wd.Count)
	}
}
