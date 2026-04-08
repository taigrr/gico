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
	ExpAuthors struct {
		Authors []string
		Created time.Time
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

// IsLeapYear returns true if year is a leap year per the Gregorian calendar.
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// YearLength returns 366 for leap years and 365 otherwise.
func YearLength(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}

func (c Commit) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s",
		c.TimeStamp.Format("0"+time.Kitchen),
		c.Author, c.Repo,
		strings.TrimSpace(c.Message))
}
