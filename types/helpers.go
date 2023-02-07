package types

import (
	"time"

	gterm "github.com/taigrr/gico/graph/term"
)

func NewDataSet() DataSet {
	return make(DataSet)
}

func NewCommit(Author, Message, Repo, Path string, Added, Deleted, FilesChanged int) Commit {
	ci := Commit{
		Message: Message, Added: Added, Deleted: Deleted,
		Author: Author, FilesChanged: FilesChanged, TimeStamp: time.Now(),
		Repo: Repo, Path: Path,
	}
	return ci
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
