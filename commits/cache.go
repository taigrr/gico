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
	mapTex    sync.RWMutex
	hashCache map[int]map[string]map[string]types.ExpYearFreq
)

func init() {
	hashCache = make(map[int]map[string]map[string]types.ExpYearFreq)
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

func GetCachedGraph(year int, authors []string, repoPaths []string) (types.YearFreq, bool) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.RLock()
	defer mapTex.RUnlock()
	if m1, ok := hashCache[year]; !ok {
		return types.YearFreq{}, false
	} else {
		if m2, ok := m1[a]; !ok {
			return types.YearFreq{}, false
		} else {
			if freq, ok := m2[r]; !ok {
				return types.YearFreq{}, false
			} else {
				if freq.Created.Before(time.Now().Add(-15 * time.Minute)) {
					return types.YearFreq{}, false
				} else {
					return freq.YearFreq, true
				}
			}
		}
	}
}

func CacheGraph(year int, authors, repoPaths []string, freq types.YearFreq) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.Lock()
	defer mapTex.Unlock()
	if _, ok := hashCache[year]; !ok {
		hashCache[year] = make(map[string]map[string]types.ExpYearFreq)
	}
	if _, ok := hashCache[year][a]; !ok {
		hashCache[year][a] = make(map[string]types.ExpYearFreq)
	}
	hashCache[year][a][r] = types.ExpYearFreq{YearFreq: freq, Created: time.Now()}
	go func() {
		time.Sleep(time.Minute * 15)
		mapTex.Lock()
		defer mapTex.Unlock()
		delete(hashCache[year][a], r)
	}()
}
