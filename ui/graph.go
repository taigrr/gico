package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/taigrr/gico/commits"
)

func shiftSelectionByWeeks(year, selected, weeks int) (int, int) {
	targetDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, selected+(weeks*7))
	return targetDate.Year(), targetDate.YearDay() - 1
}

type (
	Graph struct {
		Selected int
		Year     int
		Repos    []string
		Authors  []string
	}
)

func (m Graph) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentYear := m.Year
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
			m.Year, m.Selected = shiftSelectionByWeeks(m.Year, m.Selected, -1)
		case "right":
			m.Year, m.Selected = shiftSelectionByWeeks(m.Year, m.Selected, 1)
		}
	}
	if m.Year != currentYear {
		go func(year int, repos, authors []string) {
			mr := commits.RepoSet(repos)
			mr.FrequencyChan(year-1, authors)
			mr.FrequencyChan(year+1, authors)
		}(m.Year, m.Repos, m.Authors)
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
