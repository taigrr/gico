package commits

import (
	"errors"
	"os"

	git "github.com/go-git/go-git/v5"

	"github.com/taigrr/gico/types"
	"github.com/taigrr/mg/parse"
)

type (
	Repo      git.Repository
	CommitSet struct {
		Commits []types.Commit
		Year    int
	}
)

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

func GetMRRepos() (RepoSet, error) {
	mrconf, err := parse.LoadMRConfig()
	if err != nil {
		return RepoSet{}, err
	}
	paths := mrconf.GetRepoPaths()
	return RepoSet(paths), nil
}
