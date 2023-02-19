package commits

import (
	"log"
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
			log.Printf("skipping repo %s\n", repo.Path)
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

func (repo Repo) GetHead() (string, error) {
	r := git.Repository(repo.Repo)
	ref, err := r.Head()
	if err != nil {
		return "", err
	}
	return ref.String(), nil
}

func (repo Repo) GetCommitSet() (CommitSet, error) {
	cs := CommitSet{}
	commits := []types.Commit{}
	r := git.Repository(repo.Repo)
	ref, err := r.Head()
	if err != nil {
		return cs, err
	}
	if cachedRepo, ok := GetCachedRepo(repo.Path, ref.String()); ok {
		return CommitSet{Commits: cachedRepo}, nil
	}
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return cs, err
	}
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
		// this is too slow for now, so skipping
		//		stats, err := c.Stats()
		//		if err != nil {
		//			for _, stat := range stats {
		//				commit.Added += stat.Addition
		//				commit.Deleted += stat.Deletion
		//				commit.FilesChanged++
		//			}
		//		}
		commits = append(commits, commit)
		return nil
	})
	cs.Commits = commits
	// CacheRepo(repo.Path, cs.Commits)
	return cs, nil
}

func (cs CommitSet) ToYearFreq() types.Freq {
	year := cs.Year
	yearLength := 365
	if year%4 == 0 {
		yearLength++
	}
	freq := make([]int, yearLength)
	for _, v := range cs.Commits {
		freq[v.TimeStamp.YearDay()-1]++
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
	regset:
		for _, r := range regSet {
			if r.MatchString(commit.Author.Name) || r.MatchString(commit.Author.Email) {
				newCS.Commits = append(newCS.Commits, commit)
				break regset
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
