package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Month string

var months = []Month{"Jan", "Feb", "Mar",
	"Apr", "May", "Jun",
	"Jul", "Aug", "Sep",
	"Oct", "Nov", "Dec"}

var days [366]int

func init() {
	// parse configs
	// choose action from CLI
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		printGraph()
		os.Exit(0)
	}
	switch args[0] {
	case "inc", "increment", "add":
		increment()
	case "graph":
		printGraph()
	case "interactive":
		interactiveGraph()
	case "loadRepo":
		loadRepo()
	default:
		printHelp()
	}
}

type Commit struct {
	LOC       int       `json:"loc"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"ts"`
	Author    string    `json:"author"`
}

type DataSet map[time.Time]WorkDay

type WorkDay struct {
	Day     time.Time `json:"day"`
	Count   int       `json:"count"`
	Message string    `json:"message,omitempty"`
	Commits []Commit  `json:"commits,omitempty"`
}

func loadRepo() {

}

func readCommitDB() DataSet {
	ds := DataSet{}
	ds = append(ds,
	return ds
}
func printHelp() {
	fmt.Println("help:")
}

func increment() {
	commits := readCommitDB()
	
	fmt.Printf("%v\n", commits)
}

func printGraph() {
	fmt.Println("printGraph")
}
