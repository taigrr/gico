package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/taigrr/gico/commits"
)

type (
	Graph struct {
		Selected int
		Year     int
		Repos    []string
		Authors  []string
	}
)

func (m Graph) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	yearLen := YearLen(m.Year)
	prevYearLen := YearLen(m.Year - 1)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			if m.Selected%7 != 6 {
				m.Selected++
			}
		case "up":
			if m.Selected%7 != 0 {
				m.Selected--
			}
		case "left":
			if m.Selected > 6 {
				m.Selected -= 7
			} else {
				// TODO calculate the square for this day last year
				m.Selected -= 7
				m.Selected += prevYearLen
				m.Year--
				go func() {
					mr := commits.RepoSet(m.Repos)
					mr.FrequencyChan(m.Year-1, m.Authors)
				}()
			}
		case "right":
			if m.Selected < yearLen-7 {
				m.Selected += 7
			} else {
				m.Selected += 7
				m.Selected -= yearLen
				m.Year++
				go func() {
					mr := commits.RepoSet(m.Repos)
					mr.FrequencyChan(m.Year+1, m.Authors)
				}()
			}
		}
	}
	return m, nil
}

func NewGraph(authors, repos []string) (Graph, error) {
	var m Graph
	now := time.Now()
	today := now.YearDay() - 1
	year := now.Year()
	m.Repos = repos
	m.Authors = authors
	m.Year = year
	m.Selected = today
	return m, nil
}

func (m Graph) Init() tea.Cmd {
	go func() {
		mr := commits.RepoSet(m.Repos)
		mr.FrequencyChan(m.Year-1, m.Authors)
		mr.FrequencyChan(m.Year+1, m.Authors)
	}()
	return nil
}

func (m Graph) View() string {
	mr := commits.RepoSet(m.Repos)
	gfreq, _ := mr.FrequencyChan(m.Year, m.Authors)
	return gfreq.StringSelected(m.Selected)
}
