package commits

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
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

func GlobalFrequencyChan(year int, authors []string) (types.YearFreq, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	mrconf, err := parse.LoadMRConfig()
	if err != nil {
		return types.YearFreq{}, err
	}
	paths := mrconf.GetRepoPaths()
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

func YearFreqFromChan(cc chan types.Commit, year int) types.YearFreq {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	freq := make([]int, yearLength)
	data := types.NewDataSet()
	for commit := range cc {
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
			freq[k.YearDay()-1]++
		}
	}
	return freq
}

func GlobalFrequency(year int, authors []string) (types.YearFreq, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	gfreq := make(types.YearFreq, yearLength)
	mrconf, err := parse.LoadMRConfig()
	if err != nil {
		return types.YearFreq{}, err
	}
	paths := mrconf.GetRepoPaths()
	for _, p := range paths {
		repo, err := OpenRepo(p)
		if err != nil {
			return types.YearFreq{}, err
		}
		commits, err := repo.GetCommitSet()
		if err != nil {
			return types.YearFreq{}, err
		}
		commits = commits.FilterByYear(year)
		commits, err = commits.FilterByAuthorRegex(authors)
		if err != nil {
			return types.YearFreq{}, err
		}
		freq := commits.ToYearFreq()
		gfreq = gfreq.Merge(freq)
	}
	return gfreq, nil
}

func OpenRepo(directory string) (Repo, error) {
	if s, err := os.Stat(directory); err != nil {
		return Repo{}, err
	} else {
		if !s.IsDir() {
			return Repo{}, errors.New("received path to non-directory for git repo")
		}
	}
	r, err := git.PlainOpenWithOptions(directory, &(git.PlainOpenOptions{DetectDotGit: true}))
	return Repo(*r), err
}

type CommitSet struct {
	Commits []types.Commit
	Year    int
}

func (cs CommitSet) ToYearFreq() types.YearFreq {
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
			freq[k.YearDay()-1]++
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

func FilterCChanByYear(in chan types.Commit, year int) chan types.Commit {
	out := make(chan types.Commit, 30)
	go func() {
		for commit := range in {
			if commit.TimeStamp.Year() == year {
				out <- commit
			}
		}
		close(out)
	}()
	return out
}

func FilterCChanByAuthor(in chan types.Commit, authors []string) (chan types.Commit, error) {
	out := make(chan types.Commit, 30)
	regSet := [](*regexp.Regexp){}
	for _, a := range authors {
		r, err := regexp.Compile(a)
		if err != nil {
			close(out)
			return out, err
		}
		regSet = append(regSet, r)
	}
	go func() {
		for commit := range in {
			for _, r := range regSet {
				if r.MatchString(commit.Author) {
					out <- commit
					break
				}
			}
		}
		close(out)
	}()
	return out, nil
}

func (repo Repo) GetCommitChan() (chan types.Commit, error) {
	cc := make(chan types.Commit, 30)
	r := git.Repository(repo)
	ref, err := r.Head()
	if err != nil {
		return cc, err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return cc, err
	}
	go func() {
		cIter.ForEach(func(c *object.Commit) error {
			ts := c.Author.When
			commit := types.Commit{Author: c.Author.Name, Message: c.Message, TimeStamp: ts}
			cc <- commit
			return nil
		})
		close(cc)
	}()
	return cc, nil
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
