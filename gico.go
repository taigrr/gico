package gico

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/taigrr/gico/ui"
)

type Month string

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
		ui.InteractiveGraph()
	case "loadRepo":
		loadRepo()
	default:
		printHelp()
	}
}

type Commit struct {
	LOC       int       `json:"loc,omitempty"`
	Message   string    `json:"message,omitempty"`
	TimeStamp time.Time `json:"ts,omitempty"`
	Author    string    `json:"author,omitempty"`
	Repo      string    `json:"repo,omitempty"`
	Path      string    `json:"path,omitempty"`
}

type DataSet map[time.Time]WorkDay

func NewDataSet() DataSet {
	return make(DataSet)
}

type WorkDay struct {
	Day     time.Time `json:"day"`
	Count   int       `json:"count"`
	Commits []Commit  `json:"commits,omitempty"`
}

func NewCommit(Author, Message, Repo, Path string, LOC int) Commit {
	ci := Commit{Message: Message,
		Author: Author, LOC: LOC, TimeStamp: time.Now(),
		Repo: Repo, Path: Path}
	return ci
}

func loadRepo() {

}

func readCommitDB() DataSet {
	ds := DataSet{}
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
