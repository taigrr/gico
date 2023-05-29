package main

import (
	"fmt"
	"os"
	"time"

	"github.com/taigrr/gico/commits"
	gterm "github.com/taigrr/gico/graph/term"
)

func main() {
	year := time.Now().Year()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	r, err := commits.OpenRepo(dir)
	if err != nil {
		fmt.Println("Error opening current directory as a git repo")
		os.Exit(1)
	}
	cs, err := r.GetCommitSet()
	if err != nil {
		panic(err)
	}
	cs = cs.FilterByYear(year)
	freq := cs.ToYearFreq()
	fmt.Print(gterm.GetYearUnicode(freq))
}
