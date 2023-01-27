package main

import (
	"fmt"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	gterm "github.com/taigrr/gitgraph/term"

	"github.com/taigrr/gico/types"
)

func main() {
	year := time.Now().Year()
	for i := year - 4; i <= year; i++ {
		GetYear(i)
	}
}

func GetYear(year int) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	data := types.NewDataSet()

	r, err := git.PlainOpenWithOptions(".", &(git.PlainOpenOptions{DetectDotGit: true}))
	if err != nil {
		fmt.Printf("gitfetch error: Could not find a git repository to open!\n")
		os.Exit(1)
	}
	ref, err := r.Head()
	// TODO handle this error
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})

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
	fmt.Print(gterm.GetYearUnicode(freq))
}
