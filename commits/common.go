package commits

import (
	"errors"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"

	"github.com/taigrr/mg/parse"

	"github.com/taigrr/gico/types"
)

type (
	Repo struct {
		Repo git.Repository
		Path string
	}
	CommitSet struct {
		Commits []types.Commit
		Year    int
	}
	RepoSet []string
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
	return Repo{Repo: *r, Path: directory}, err
}

func GetRepos() (RepoSet, error) {
	mgconf, err := parse.LoadMGConfig()
	if err != nil {
		mrconf, err := parse.LoadMRConfig()
		if err != nil {
			return RepoSet{}, err
		}
		paths := mrconf.GetRepoPaths()
		return RepoSet(paths), nil
	}
	paths := mgconf.GetRepoPaths()
	return RepoSet(paths), nil
}

func GetAuthorName() (string, error) {
	conf, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		return "", err
	}
	return conf.User.Name, nil
}

func GetAuthorEmail() (string, error) {
	conf, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		return "", err
	}
	return conf.User.Email, nil
}
