package commits

import (
	"testing"
	"time"

	"github.com/taigrr/gico/types"
)

func TestHashSliceDeterministic(t *testing.T) {
	a := hashSlice([]string{"foo", "bar", "baz"})
	b := hashSlice([]string{"baz", "foo", "bar"})
	if a != b {
		t.Error("hashSlice should be order-independent")
	}
}

func TestHashSliceDoesNotMutateInput(t *testing.T) {
	input := []string{"foo", "bar", "baz"}
	want := []string{"foo", "bar", "baz"}

	_ = hashSlice(input)

	for i := range want {
		if input[i] != want[i] {
			t.Fatalf("hashSlice mutated input at index %d: got %q, want %q", i, input[i], want[i])
		}
	}
}

func TestHashSliceDifferentInputs(t *testing.T) {
	a := hashSlice([]string{"foo", "bar"})
	b := hashSlice([]string{"foo", "baz"})
	if a == b {
		t.Error("different inputs should produce different hashes")
	}
}

func TestCacheRepoRoundTrip(t *testing.T) {
	path := "/test/cache-repo"
	head := "deadbeef123"
	commits := []types.Commit{
		{Hash: head, Author: types.Author{Name: "test"}, TimeStamp: time.Now()},
	}

	// Should miss before caching
	_, ok := GetCachedRepo(path, head)
	if ok {
		t.Error("expected cache miss before storing")
	}

	CacheRepo(path, commits)

	// Should hit after caching
	cached, ok := GetCachedRepo(path, head)
	if !ok {
		t.Error("expected cache hit after storing")
	}
	if len(cached) != 1 || cached[0].Hash != head {
		t.Error("cached data doesn't match stored data")
	}

	// Different head should miss
	_, ok = GetCachedRepo(path, "differenthead")
	if ok {
		t.Error("expected cache miss for different head")
	}
}

func TestCacheGraphRoundTrip(t *testing.T) {
	year := 2025
	authors := []string{"alice@example.com"}
	paths := []string{"/repo/one", "/repo/two"}
	freq := types.Freq{1, 2, 3, 0, 0}

	// Should miss before caching
	_, ok := GetCachedGraph(year, authors, paths)
	if ok {
		t.Error("expected cache miss before storing")
	}

	CacheGraph(year, authors, paths, freq)

	// Should hit
	cached, ok := GetCachedGraph(year, authors, paths)
	if !ok {
		t.Error("expected cache hit")
	}
	if len(cached) != len(freq) {
		t.Errorf("expected freq length %d, got %d", len(freq), len(cached))
	}
}

func TestCacheReposAuthorsRoundTrip(t *testing.T) {
	paths := []string{"/repo/author-test"}
	authors := []string{"alice", "bob"}

	_, ok := GetCachedReposAuthors(paths)
	if ok {
		t.Error("expected cache miss")
	}

	CacheReposAuthors(paths, authors)

	cached, ok := GetCachedReposAuthors(paths)
	if !ok {
		t.Error("expected cache hit")
	}
	if len(cached) != 2 {
		t.Errorf("expected 2 authors, got %d", len(cached))
	}
}
