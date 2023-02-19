package commits

import (
	"regexp"
	"sort"
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/taigrr/gico/types"
)

func (paths RepoSet) GetRepoAuthors() ([]string, error) {
	outChan := make(chan types.Commit, 10)
	var wg sync.WaitGroup
	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			repo, err := OpenRepo(path)
			if err != nil {
				return
			}
			cc, err := repo.GetCommitChan()
			if err != nil {
				return
			}
			for c := range cc {
				outChan <- c
			}
		}(p)
	}
	go func() {
		wg.Wait()
		close(outChan)
	}()
	authors := make(map[string]bool)
	for commit := range outChan {
		authors[commit.Author.Email] = true
		authors[commit.Author.Name] = true
	}
	a := []string{}
	for k := range authors {
		a = append(a, k)
	}
	sort.Strings(a)
	return a, nil
}

func (paths RepoSet) GetRepoCommits(year int, authors []string) ([][]types.Commit, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}

	commits := make([][]types.Commit, yearLength)
	for i := 0; i < yearLength; i++ {
		commits[i] = []types.Commit{}
	}
	cache, ok := GetCachedRepos(year, authors, paths)
	if ok {
		return cache, nil
	}
	outChan := make(chan types.Commit, 10)
	var wg sync.WaitGroup
	for _, p := range paths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			repo, err := OpenRepo(path)
			if err != nil {
				return
			}
			cc, err := repo.GetCommitChan()
			if err != nil {
				return
			}
			cc = FilterCChanByYear(cc, year)
			cc2, err := FilterCChanByAuthor(cc, authors)
			if err != nil {
				return
			}
			for c := range cc2 {
				outChan <- c
			}
		}(p)
	}
	go func() {
		wg.Wait()
		close(outChan)
	}()
	for commit := range outChan {
		d := commit.TimeStamp.YearDay() - 1
		commits[d] = append(commits[d], commit)
	}
	for i := 0; i < len(commits); i++ {
		sort.Slice(commits[i], func(w, j int) bool {
			return commits[i][w].TimeStamp.Before(commits[i][j].TimeStamp)
		})
	}
	CacheRepos(year, authors, paths, commits)
	return commits, nil
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
			defer wg.Done()
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
		freq[commit.TimeStamp.YearDay()-1]++
	}
	return freq
}

func (repo Repo) GetCommitChan() (chan types.Commit, error) {
	cc := make(chan types.Commit, 30)
	r := git.Repository(repo.Repo)
	ref, err := r.Head()
	if err != nil {
		close(cc)
		return cc, err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		close(cc)
		return cc, err
	}
	go func() {
		cIter.ForEach(func(c *object.Commit) error {
			ts := c.Author.When
			commit := types.Commit{
				Author: types.Author{
					Name:  c.Author.Name,
					Email: c.Author.Email,
				},
				Message: c.Message, TimeStamp: ts,
				Hash: c.Hash.String(), Repo: repo.Path,
				FilesChanged: 0, Added: 0, Deleted: 0,
			}
			// Too slow, commenting for now
			//			stats, err := c.Stats()
			//			if err != nil {
			//				for _, stat := range stats {
			//					commit.Added += stat.Addition
			//					commit.Deleted += stat.Deletion
			//					commit.FilesChanged++
			//				}
			//			}
			cc <- commit
			return nil
		})
		close(cc)
	}()
	return cc, nil
}

func FreqFromChan(cc chan types.Commit, year int) types.Freq {
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
		regset:
			for _, r := range regSet {
				if r.MatchString(commit.Author.Email) || r.MatchString(commit.Author.Name) {
					out <- commit
					break regset
				}
			}
		}
		close(out)
	}()
	return out, nil
}
