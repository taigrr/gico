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
	mapTex        sync.RWMutex
	freqHashCache map[int]map[string]map[string]types.ExpFreq
	repoHashCache map[int]map[string]map[string]types.ExpRepo
	// the Repo Cache holds a list of all commits from HEAD back to parent
	// the key is the repo path
	// if the hash of the first commit / HEAD commit doesn't match the current HEAD,
	// then it can be discarded and reloaded
	repoCache map[string][]types.Commit
)

func init() {
	freqHashCache = make(map[int]map[string]map[string]types.ExpFreq)
	repoCache = make(map[string][]types.Commit)
}

func GetCachedRepo(path string, head string) ([]types.Commit, bool) {
	mapTex.RLock()
	defer mapTex.RUnlock()
	if commits, ok := repoCache[path]; !ok {
		return []types.Commit{}, false
	} else if len(commits) > 0 && commits[0].Hash == head {
		return commits, true
	}
	return []types.Commit{}, false
}

func IsRepoCached(path string, head string) bool {
	mapTex.RLock()
	defer mapTex.RUnlock()
	if commits, ok := repoCache[path]; !ok {
		return false
	} else {
		return len(commits) > 0 && commits[0].Hash == head
	}
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

func CacheRepo(year int, authors, repoPaths []string, commits []types.Commit) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.Lock()
	defer mapTex.Unlock()
	if _, ok := repoHashCache[year]; !ok {
		repoHashCache[year] = make(map[string]map[string]types.ExpRepo)
	}
	if _, ok := repoHashCache[year][a]; !ok {
		repoHashCache[year][a] = make(map[string]types.ExpRepo)
	}
	repoHashCache[year][a][r] = types.ExpRepo{Commits: commits, Created: time.Now()}
	go func() {
		time.Sleep(time.Hour * 1)
		mapTex.Lock()
		defer mapTex.Unlock()
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
