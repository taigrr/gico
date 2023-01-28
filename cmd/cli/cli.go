package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/gitgraph/term"
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}

func main() {
	wfreq := []int{}
	yd := time.Now().YearDay()
	freq, _ := commits.GlobalFrequency(time.Now().Year(), []string{""})
	if yd < 7 {
		//	xfreq, _ := commits.GlobalFrequency(time.Now().Year()-1, []string{""})
	} else {
		// TODO fix bug for negative in first week of Jan
		for i := 0; i < 7; i++ {
			d := time.Now().YearDay() - 1 - 6 + i
			wfreq = append(wfreq, freq[d])
		}
	}
	fmt.Println("week:")
	fmt.Println(term.GetWeekUnicode(wfreq))
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("year:")
	fmt.Println(term.GetYearUnicode(freq))
}
