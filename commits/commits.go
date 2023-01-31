package commits

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/taigrr/gico/types"
	"github.com/taigrr/mg/parse"
)

type Repo git.Repository

type RepoSet []string

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
		time.Sleep(time.Minute * 15)
		mapTex.Lock()
		defer mapTex.Unlock()
		delete(hashCache[year][a], r)
	}()
}

func GetMRRepos() (RepoSet, error) {
	mrconf, err := parse.LoadMRConfig()
	if err != nil {
		return RepoSet{}, err
	}
	paths := mrconf.GetRepoPaths()
	return RepoSet(paths), nil
}

func (paths RepoSet) GetWeekFreq(authors []string) (types.Freq, error) {
	now := time.Now()
	year := now.Year()
	freq, err := paths.FrequencyChan(year, authors)
	if err != nil {
		return types.Freq{}, err
	}
	today := now.YearDay() - 1
	if today < 6 {
		curYear := year - 1
		curFreq, err := paths.FrequencyChan(curYear, authors)
		if err != nil {
			return types.Freq{}, err
		}
		freq = append(curFreq, freq...)
		today += 365
		if curYear%4 == 0 {
			today++
		}
	}

	week := freq[today-6 : today+1]
	return week, nil
}

func (paths RepoSet) FrequencyChan(year int, authors []string) (types.Freq, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	cache, ok := GetCachedGraph(year, authors, paths)
	if ok {
		return cache, nil
	}
	outChan := make(chan types.Commit, 10)
	var wg sync.WaitGroup
	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			repo, err := OpenRepo(path)
			if err != nil {
				return
			}
			cc, err := repo.GetCommitChan()
			if err != nil {
				return
			}
			cc = FilterCChanByYear(cc, year)
			cc, err = FilterCChanByAuthor(cc, authors)
			if err != nil {
				return
			}
			for c := range cc {
				outChan <- c
			}
			wg.Done()
		}(p)
	}
	go func() {
		wg.Wait()
		close(outChan)
	}()
	freq := YearFreqFromChan(outChan, year)
	CacheGraph(year, authors, paths, freq)
	return freq, nil
}

func YearFreqFromChan(cc chan types.Commit, year int) types.Freq {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	freq := make([]int, yearLength)
	for commit := range cc {
		if commit.TimeStamp.Year() != year {
			continue
		}
		freq[commit.TimeStamp.YearDay()-1]++
	}
	return freq
}

func (paths RepoSet) GlobalFrequency(year int, authors []string) (types.Freq, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	gfreq := make(types.Freq, yearLength)
	for _, p := range paths {
		repo, err := OpenRepo(p)
		if err != nil {
			return types.Freq{}, err
		}
		commits, err := repo.GetCommitSet()
		if err != nil {
			return types.Freq{}, err
		}
		commits = commits.FilterByYear(year)
		commits, err = commits.FilterByAuthorRegex(authors)
		if err != nil {
			return types.Freq{}, err
		}
		freq := commits.ToYearFreq()
		gfreq = gfreq.Merge(freq)
	}
	return gfreq, nil
}

func (repo Repo) GetCommitSet() (CommitSet, error) {
	cs := CommitSet{}
	commits := []types.Commit{}
	r := git.Repository(repo)
	ref, err := r.Head()
	if err != nil {
		return cs, err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return cs, err
	}
	cIter.ForEach(func(c *object.Commit) error {
		ts := c.Author.When
		commit := types.Commit{Author: c.Author.Name, Message: c.Message, TimeStamp: ts}
		commits = append(commits, commit)
		return nil
	})
	cs.Commits = commits
	return cs, nil
}

func (cs CommitSet) ToYearFreq() types.Freq {
	year := cs.Year
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	freq := make([]int, yearLength)
	data := types.NewDataSet()
	for _, commit := range cs.Commits {
		ts := commit.TimeStamp
		roundedTS := ts.Round(time.Hour * 24)
		wd, ok := data[roundedTS]
		if !ok {
			wd = types.WorkDay{}
			wd.Commits = []types.Commit{}
		}
		wd.Commits = append(wd.Commits, commit)
		wd.Count++
		wd.Day = roundedTS
		data[roundedTS] = wd
	}
	for k, v := range data {
		if k.Year() != year {
			continue
		}
		// this is equivalent to adding len(commits) to the freq total, but
		// it's a stub for later when we do more here
		for range v.Commits {
			freq[k.YearDay()]++
		}
	}
	return freq
}

func (cs CommitSet) FilterByAuthorRegex(authors []string) (CommitSet, error) {
	regSet := [](*regexp.Regexp){}
	for _, a := range authors {
		r, err := regexp.Compile(a)
		if err != nil {
			return CommitSet{}, err
		}
		regSet = append(regSet, r)
	}
	newCS := CommitSet{Year: cs.Year}
	for _, commit := range cs.Commits {
		for _, r := range regSet {
			if r.MatchString(commit.Author) {
				newCS.Commits = append(newCS.Commits, commit)
				break
			}
		}
	}
	return newCS, nil
}

func (cs CommitSet) FilterByYear(year int) CommitSet {
	newCS := CommitSet{Year: year}
	for _, commit := range cs.Commits {
		if commit.TimeStamp.Year() == year {
			newCS.Commits = append(newCS.Commits, commit)
		}
	}
	return newCS
}
