package types

import (
	"fmt"
	"strings"
	"time"
)

type (
	Month  string
	Author struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	}
	Commit struct {
		Deleted      int       `json:"deleted,omitempty"`
		Added        int       `json:"added,omitempty"`
		FilesChanged int       `json:"files_changed,omitempty"`
		Message      string    `json:"message,omitempty"`
		Hash         string    `json:"hash,omitempty"`
		TimeStamp    time.Time `json:"ts,omitempty"`
		Author       Author    `json:"author,omitempty"`
		Repo         string    `json:"repo,omitempty"`
		Path         string    `json:"path,omitempty"`
	}
	DataSet map[time.Time]WorkDay
	Freq    []int
	ExpFreq struct {
		YearFreq Freq
		Created  time.Time
	}
	ExpRepos struct {
		Commits [][]Commit
		Created time.Time
	}
	ExpRepo struct {
		Commits []Commit
		Created time.Time
	}
	WorkDay struct {
		Day     time.Time `json:"day"`
		Count   int       `json:"count"`
		Commits []Commit  `json:"commits,omitempty"`
	}
)

func (c Commit) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\n",
		c.TimeStamp.Format("0"+time.Kitchen),
		c.Author, c.Repo,
		strings.TrimSpace(c.Message))
}
