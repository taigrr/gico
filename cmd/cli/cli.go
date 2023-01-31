package main

import (
	"fmt"
	"time"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/graph/term"
)

func main() {
	n := time.Now()
	repoPaths, err := commits.GetMRRepos()
	if err != nil {
		panic(err)
	}
	freq, err := repoPaths.Frequency(n.Year(), []string{"Groot"})
	if err != nil {
		panic(err)
	}
	wfreq, err := repoPaths.GetWeekFreq([]string{"Groot"})
	if err != nil {
		panic(err)
	}
	fmt.Println("week:")
	fmt.Println(term.GetWeekUnicode(wfreq))
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("year:")
	fmt.Println(term.GetYearUnicode(freq))
}
