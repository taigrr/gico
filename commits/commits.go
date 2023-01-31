package commits

import (
	"regexp"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/taigrr/gico/types"
)

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

func (paths RepoSet) Frequency(year int, authors []string) (types.Freq, error) {
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
