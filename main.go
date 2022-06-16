package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/taigrr/gitgraph/graph"
)

type DayCount [366]int

func main() {
	freq := []int{}
	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < 366; i++ {
		freq = append(freq, rand.Int())
	}
	svg := graph.GetYearImage(freq)
	f, err := os.Create("out.svg")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer f.Close()
	svg.WriteTo(f)
}
