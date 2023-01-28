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
)

func main() {
	year := time.Now().Year()
	repo, err := OpenRepo(".")
	if err != nil {
		panic(err)
	}
	str, err := repo.GetYear(year)
	if err != nil {
		panic(err)
	}
	fmt.Print(str)
}

type Repo git.Repository

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

func (repo Repo) GetYear(year int) (string, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	data := types.NewDataSet()
	r := git.Repository(repo)
	ref, err := r.Head()
	if err != nil {
		return "", err
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return "", err
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

	freq := make([]int, yearLength)
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
	return gterm.GetYearUnicode(freq), nil
}
