package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/taigrr/gico/ui"
)

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
	// TODO use cobra-cli instead of switch case
	case "inc", "increment", "add":
		increment()
	case "graph":
		printGraph()
	case "interactive":
		ui.InteractiveGraph()
	case "loadRepo":
		loadRepo()
	default:
		printHelp()
	}
}

func NewDataSet() types.DataSet {
	return make(types.DataSet)
}

func NewCommit(Author, Message, Repo, Path string, LOC int) types.Commit {
	ci := types.Commit{
		Message: Message,
		Author:  Author, LOC: LOC, TimeStamp: time.Now(),
		Repo: Repo, Path: Path,
	}
	return ci
}

func loadRepo() {
}

func readCommitDB() types.DataSet {
	ds := types.DataSet{}
	return ds
}

func printHelp() {
	fmt.Println("help:")
}

func increment() {
	commits := readCommitDB()
	// crea
	fmt.Printf("%v\n", commits)
}

func printGraph() {
	fmt.Println("printGraph")
}
