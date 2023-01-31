package commits

import (
	"regexp"
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/taigrr/gico/types"
	"github.com/taigrr/mg/parse"
)

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

func YearFreqFromChan(cc chan types.Commit, year int) types.YearFreq {
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
