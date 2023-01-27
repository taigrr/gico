package types

import (
	"time"
)

func NewDataSet() DataSet {
	return make(DataSet)
}

func NewCommit(Author, Message, Repo, Path string, LOC int) Commit {
	ci := Commit{
		Message: Message,
		Author:  Author, LOC: LOC, TimeStamp: time.Now(),
		Repo: Repo, Path: Path,
	}
	return ci
}
