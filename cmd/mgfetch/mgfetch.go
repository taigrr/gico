package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gterm "github.com/taigrr/gico/gitgraph/term"

	"github.com/taigrr/gico/types"
	"github.com/taigrr/mg/parse"
)

func main() {
	year := time.Now().Year() - 1
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}

	gfreq := make(YearFreq, yearLength)

	mrconf, err := parse.LoadMRConfig()
	if err != nil {
		panic(err)
	}
	paths := mrconf.GetRepoPaths()
	for _, p := range paths {
		repo, err := OpenRepo(p)
		if err != nil {
			panic(err)
		}
		freq, err := repo.GetYear(year)
		gfreq = gfreq.Merge(freq)
		if err != nil {
			panic(err)
		}
	}
	fmt.Print(gfreq.String())
}

type Repo git.Repository

func (a YearFreq) Merge(b YearFreq) YearFreq {
	x := len(a)
	y := len(b)
	if x < y {
		x = y
	}
	c := make(YearFreq, x)
	copy(c, a)
	for i := 0; i < y; i++ {
		c[i] += b[i]
	}
	return c
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

type YearFreq []int

func (yf YearFreq) String() string {
	return gterm.GetYearUnicode(yf)
}

func (repo Repo) GetYear(year int) (YearFreq, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	freq := make([]int, yearLength)
	data := types.NewDataSet()
	r := git.Repository(repo)
	ref, err := r.Head()
	if err != nil {
		return freq, err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return freq, err
	}
	err = cIter.ForEach(func(c *object.Commit) error {
		ts := c.Author.When
		commit := types.Commit{Author: c.Author.Name, Message: c.Message, TimeStamp: ts}
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
		return nil
	})

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
	return freq, nil
}
