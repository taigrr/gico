package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/taigrr/gico/types"
	"github.com/taigrr/mg/parse"
)

type Repo git.Repository

func main() {
	year := time.Now().Year() - 1
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	gfreq := make(types.YearFreq, yearLength)
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
		commits, err := repo.GetCommitSetForYear(year)
		if err != nil {
			panic(err)
		}
		commits, err = commits.FilterByAuthorRegex([]string{"Groot"})
		if err != nil {
			panic(err)
		}
		freq := commits.ToYearFreq()
		gfreq = gfreq.Merge(freq)
	}
	fmt.Print(gfreq.String())
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

func (repo Repo) GetCommitSetForYear(year int) (CommitSet, error) {
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
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
	err = cIter.ForEach(func(c *object.Commit) error {
		ts := c.Author.When
		commit := types.Commit{Author: c.Author.Name, Message: c.Message, TimeStamp: ts}
		if commit.TimeStamp.Year() == year {
			commits = append(commits, commit)
		}
		return nil
	})
	cs.Commits = commits
	cs.Year = year
	return cs, nil
}
