package types

import (
	"time"

	gterm "github.com/taigrr/gico/graph/term"
)

func NewDataSet() DataSet {
	return make(DataSet)
}

func NewCommit(Author, Message, Repo, Path string, LOC int) Commit {
	return Commit{
		Message: Message,
		Author:  Author, LOC: LOC, TimeStamp: time.Now(),
		Repo: Repo, Path: Path,
	}
}

func (yf Freq) String() string {
	return gterm.GetYearUnicode(yf)
}

func (a Freq) Merge(b Freq) Freq {
	x := len(a)
	y := len(b)
	if x < y {
		x = y
	}
	c := make(Freq, x)
	copy(c, a)
	for i := 0; i < y; i++ {
		c[i] += b[i]
	}
	return c
}
