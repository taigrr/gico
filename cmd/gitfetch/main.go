package main

import (
	"fmt"
	"os"
	"time"

	"github.com/taigrr/gico/cmd/internal/cli"
	"github.com/taigrr/gico/commits"
	gterm "github.com/taigrr/gico/graph/term"
)

func main() {
	year := time.Now().Year()
	dir, err := os.Getwd()
	if err != nil {
		cli.Fatal("failed to get current directory", err)
	}
	r, err := commits.OpenRepo(dir)
	if err != nil {
		cli.Fatal("failed to open current directory as a git repo", err)
	}
	cs, err := r.GetCommitSet()
	if err != nil {
		cli.Fatal("failed to read commits", err)
	}
	cs = cs.FilterByYear(year)
	freq := cs.ToYearFreq()
	fmt.Print(gterm.GetYearUnicode(freq))
}
