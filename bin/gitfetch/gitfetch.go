package main

import (
	"fmt"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/taigrr/gico"
	gterm "github.com/taigrr/gitgraph/term"
)

func main() {

	r, err := git.PlainOpenWithOptions(".", &(git.PlainOpenOptions{DetectDotGit: true}))
	if err != nil {
		fmt.Printf("gitfetch error: Could not find a git repository to open!\n")
		os.Exit(1)
	}
	ref, err := r.Head()
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})

	//	year := time.Now().Year()
	//	yearStart := time.Time{}
	//	yearStart.AddDate(year, 0, 0)
	data := gico.NewDataSet()
	err = cIter.ForEach(func(c *object.Commit) error {
		ts := c.Author.When
		commit := gico.Commit{Author: c.Author.Name, Message: c.Message, TimeStamp: ts}
		roundedTS := ts.Round(time.Hour * 24)
		wd, ok := data[roundedTS]
		if !ok {
			wd = gico.WorkDay{}
			wd.Commits = []gico.Commit{}
		}
		wd.Commits = append(wd.Commits, commit)
		wd.Count++
		wd.Day = roundedTS
		data[roundedTS] = wd
		return nil
	})
	freq := [366]int{}
	for k, v := range data {
		if k.Year() != time.Now().Year() {
			continue
		}
		// this is equivalent to adding len(commits) to the freq total, but
		// it's a stub for later when we do more here
		for range v.Commits {
			freq[k.YearDay()-1]++
		}
	}
	gterm.GetYearUnicode(freq[:])
}
