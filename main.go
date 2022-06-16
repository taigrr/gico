package main

import (
	"fmt"
	"os"

	"github.com/taigrr/gitgraph/graph"
)

type DayCount [366]int

func main() {
	svg := graph.GetImage([]int{1, 2, 5, 6, 5, 4, 5, 8, 7, 43, 2, 3})
	f, err := os.Create("out.svg")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer f.Close()
	svg.WriteTo(f)
}
