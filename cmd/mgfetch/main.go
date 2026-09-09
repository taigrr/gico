package main

import (
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"

	"github.com/taigrr/gico/cmd/internal/cli"
	"github.com/taigrr/gico/commits"
)

type Repo git.Repository

func main() {
	year := time.Now().Year()
	aName, _ := commits.GetAuthorName()
	aEmail, _ := commits.GetAuthorEmail()
	authors := []string{aName, aEmail}
	mr, err := commits.GetRepos()
	if err != nil {
		cli.Fatal("failed to load repo list", err)
	}
	if len(mr) == 0 {
		cli.Fatal("found no repos", nil)
	}
	gfreq, err := mr.FrequencyChan(year, authors)
	if err != nil {
		cli.Fatal("failed to calculate commit frequency", err)
	}
	fmt.Print(gfreq.String())
}
