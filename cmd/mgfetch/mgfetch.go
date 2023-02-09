package main

import (
	"fmt"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"

	"github.com/taigrr/gico/commits"
)

type Repo git.Repository

func main() {
	year := time.Now().Year()
	aName, _ := commits.GetAuthorName()
	aEmail, _ := commits.GetAuthorEmail()
	authors := []string{aName, aEmail}
	mr, err := commits.GetMRRepos()
	if err != nil {
		panic(err)
	}
	if len(mr) == 0 {
		fmt.Println("found no repos!")
		os.Exit(1)
	}
	gfreq, err := mr.FrequencyChan(year, authors)
	if err != nil {
		panic(err)
	}
	fmt.Print(gfreq.String())
}
