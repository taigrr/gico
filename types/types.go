package types

import (
	"time"
)

type (
	Month  string
	Commit struct {
		Deleted      int       `json:"deleted,omitempty"`
		Added        int       `json:"added,omitempty"`
		FilesChanged int       `json:"files_changed,omitempty"`
		Message      string    `json:"message,omitempty"`
		Hash         string    `json:"hash,omitempty"`
		TimeStamp    time.Time `json:"ts,omitempty"`
		Author       string    `json:"author,omitempty"`
		Repo         string    `json:"repo,omitempty"`
		Path         string    `json:"path,omitempty"`
	}
	DataSet map[time.Time]WorkDay
	Freq    []int
	ExpFreq struct {
		YearFreq Freq
		Created  time.Time
	}
	WorkDay struct {
		Day     time.Time `json:"day"`
		Count   int       `json:"count"`
		Commits []Commit  `json:"commits,omitempty"`
	}
)
