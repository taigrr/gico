package commits

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/taigrr/gico/types"
)

var (
	mapTex          sync.RWMutex
	freqHashCache   map[int]map[string]map[string]types.ExpFreq
	repoHashCache   map[int]map[string]map[string]types.ExpRepos
	authorHashCache map[string]types.ExpAuthors
	repoCache       = make(map[string]types.ExpRepo)
	// the Repo Cache holds a list of all commits from HEAD back to parent
	// the key is the repo path
	// if the hash of the first commit / HEAD commit doesn't match the current HEAD,
	// then it can be discarded and reloaded
)

func init() {
	freqHashCache = make(map[int]map[string]map[string]types.ExpFreq)
	repoHashCache = make(map[int]map[string]map[string]types.ExpRepos)
	repoCache = make(map[string]types.ExpRepo)
	authorHashCache = make(map[string]types.ExpAuthors)
}

func hashSlice(in []string) string {
	sort.Strings(in)
	sb := strings.Builder{}
	for _, s := range in {
		sb.WriteString(s)
	}
	h := md5.New()
	h.Write([]byte(sb.String()))
	b := h.Sum(nil)
	return fmt.Sprintf("%x\n", b)
}

func GetCachedGraph(year int, authors []string, repoPaths []string) (types.Freq, bool) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.RLock()
	defer mapTex.RUnlock()
	if m1, ok := freqHashCache[year]; !ok {
		return types.Freq{}, false
	} else {
		if m2, ok := m1[a]; !ok {
			return types.Freq{}, false
		} else {
			if freq, ok := m2[r]; !ok {
				return types.Freq{}, false
			} else {
				if freq.Created.Before(time.Now().Add(-15 * time.Minute)) {
					return types.Freq{}, false
				} else {
					return freq.YearFreq, true
				}
			}
		}
	}
}

func GetCachedRepo(path string, head string) ([]types.Commit, bool) {
	mapTex.RLock()
	defer mapTex.RUnlock()
	if commits, ok := repoCache[path]; !ok {
		return []types.Commit{}, false
	} else if len(commits.Commits) > 0 && commits.Commits[0].Hash == head {
		return commits.Commits, true
	}
	return []types.Commit{}, false
}

func CacheRepo(path string, commits []types.Commit) {
	mapTex.Lock()
	defer mapTex.Unlock()
	repoCache[path] = types.ExpRepo{Commits: commits, Created: time.Now()}
	go func() {
		time.Sleep(time.Hour * 1)
		mapTex.Lock()
		defer mapTex.Unlock()
		delete(repoCache, path)
	}()
}

func CacheReposAuthors(paths []string, authors []string) {
	r := hashSlice(paths)
	mapTex.Lock()
	defer mapTex.Unlock()
	authorHashCache[r] = types.ExpAuthors{Authors: authors, Created: time.Now()}
	go func() {
		time.Sleep(time.Hour * 1)
		mapTex.Lock()
		defer mapTex.Unlock()
		delete(authorHashCache, r)
	}()
}

func GetCachedReposAuthors(paths []string) ([]string, bool) {
	r := hashSlice(paths)
	mapTex.RLock()
	defer mapTex.RUnlock()
	if m1, ok := authorHashCache[r]; !ok {
		return []string{}, false
	} else if m1.Created.Before(time.Now().Add(time.Minute * -15)) {
		return []string{}, false
	} else {
		return m1.Authors, true
	}
}

func GetCachedRepos(year int, authors, repoPaths []string) ([][]types.Commit, bool) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.RLock()
	defer mapTex.RUnlock()
	if m1, ok := repoHashCache[year]; !ok {
		return [][]types.Commit{{}}, false
	} else {
		if m2, ok := m1[a]; !ok {
			return [][]types.Commit{{}}, false
		} else {
			if commits, ok := m2[r]; !ok {
				return [][]types.Commit{{}}, false
			} else {
				if commits.Created.Before(time.Now().Add(-15 * time.Minute)) {
					return [][]types.Commit{{}}, false
				} else {
					return commits.Commits, true
				}
			}
		}
	}
}

func CacheRepos(year int, authors, repoPaths []string, commits [][]types.Commit) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.Lock()
	defer mapTex.Unlock()
	if _, ok := repoHashCache[year]; !ok {
		repoHashCache[year] = make(map[string]map[string]types.ExpRepos)
	}
	if _, ok := repoHashCache[year][a]; !ok {
		repoHashCache[year][a] = make(map[string]types.ExpRepos)
	}
	repoHashCache[year][a][r] = types.ExpRepos{Commits: commits, Created: time.Now()}
	go func() {
		time.Sleep(time.Hour * 1)
		mapTex.Lock()
		defer mapTex.Unlock()
		// optimization, check if the creation time has changed since the last usage
		delete(repoHashCache[year][a], r)
	}()
}

func CacheGraph(year int, authors, repoPaths []string, freq types.Freq) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.Lock()
	defer mapTex.Unlock()
	if _, ok := freqHashCache[year]; !ok {
		freqHashCache[year] = make(map[string]map[string]types.ExpFreq)
	}
	if _, ok := freqHashCache[year][a]; !ok {
		freqHashCache[year][a] = make(map[string]types.ExpFreq)
	}
	freqHashCache[year][a][r] = types.ExpFreq{YearFreq: freq, Created: time.Now()}
	go func() {
		time.Sleep(time.Hour * 1)
		mapTex.Lock()
		defer mapTex.Unlock()
		delete(freqHashCache[year][a], r)
	}()
}
