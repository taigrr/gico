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
	hashCache map[int]map[string]map[string]types.ExpFreq
)

func init() {
	hashCache = make(map[int]map[string]map[string]types.ExpFreq)
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
	if m1, ok := hashCache[year]; !ok {
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

func CacheGraph(year int, authors, repoPaths []string, freq types.Freq) {
	a := hashSlice(authors)
	r := hashSlice(repoPaths)
	mapTex.Lock()
	defer mapTex.Unlock()
	if _, ok := hashCache[year]; !ok {
		hashCache[year] = make(map[string]map[string]types.ExpFreq)
	}
	if _, ok := hashCache[year][a]; !ok {
		hashCache[year][a] = make(map[string]types.ExpFreq)
	}
	hashCache[year][a][r] = types.ExpFreq{YearFreq: freq, Created: time.Now()}
	go func() {
		time.Sleep(time.Hour * 1)
		mapTex.Lock()
		defer mapTex.Unlock()
		delete(hashCache[year][a], r)
	}()
}
