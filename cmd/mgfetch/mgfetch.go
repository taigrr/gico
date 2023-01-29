package main

import (
	"fmt"
	"time"

	git "github.com/go-git/go-git/v5"

	"github.com/taigrr/gico/commits"
)

type Repo git.Repository

func main() {
	year := time.Now().Year() - 1
	authors := []string{"Groot"}
	gfreq, err := commits.GlobalFrequencyChan(year, authors)
	if err != nil {
		panic(err)
	}
	fmt.Print(gfreq.String())
}
