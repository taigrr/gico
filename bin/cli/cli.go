package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/taigrr/gitgraph/term"
)

func init() {
	rand.Seed(time.Now().UnixMilli())
}
func main() {
	freq := []int{}
	for i := 0; i < 7; i++ {
		freq = append(freq, rand.Int())
	}
	fmt.Println("week:")
	term.GetWeekUnicode(freq)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	freq = []int{}
	for i := 0; i < 365; i++ {
		freq = append(freq, rand.Int())
	}
	fmt.Println("year:")
	term.GetYearUnicode(freq)
}
