package types

import (
	"time"
)

type Month string

type Commit struct {
	LOC       int       `json:"loc,omitempty"`
	Message   string    `json:"message,omitempty"`
	TimeStamp time.Time `json:"ts,omitempty"`
	Author    string    `json:"author,omitempty"`
	Repo      string    `json:"repo,omitempty"`
	Path      string    `json:"path,omitempty"`
}

type DataSet map[time.Time]WorkDay

type WorkDay struct {
	Day     time.Time `json:"day"`
	Count   int       `json:"count"`
	Commits []Commit  `json:"commits,omitempty"`
}
