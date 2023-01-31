package main

import (
	"fmt"
	"time"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/graph/term"
)

func main() {
	wfreq := []int{}
	yd := time.Now().YearDay()
	repoPaths, err := commits.GetMRRepos()
	if err != nil {
		panic(err)
	}
	freq, err := repoPaths.GlobalFrequency(time.Now().Year(), []string{""})
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
