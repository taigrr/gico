package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/graph/term"
)

func xmain() {
	n := time.Now()
	repoPaths, err := commits.GetMRRepos()
	if err != nil {
		panic(err)
	}
	freq, err := repoPaths.Frequency(n.Year(), []string{"Groot"})
	if err != nil {
		panic(err)
	}
	wfreq, err := repoPaths.GetWeekFreq([]string{"Groot"})
	if err != nil {
		panic(err)
	}
	fmt.Println("week:")
	fmt.Println(term.GetWeekUnicode(wfreq))
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("year:")
	fmt.Println(term.GetYearUnicode(freq))
}

type errMsg error

type model struct {
	SettingsModel  Settings
	GraphModel     Graph
	CommitLogModel CommitLog
	HelpModel      Help
	quitting       bool
	err            error
}

type Help struct{}

type CommitLog struct{}

type Settings struct{}

type Graph struct {
	Selected int
	Year     int
	Repos    []string
	Authors  []string
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func initialModel() (model, error) {
	var m model
	var err error
	m.GraphModel, err = NewGraph()
	if err != nil {
		return m, err
	}
	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Settings) Init() tea.Cmd {
	return nil
}

func (m Settings) View() string {
	return ""
}

func (m CommitLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m CommitLog) Init() tea.Cmd {
	return nil
}

func (m CommitLog) View() string {
	return ""
}

func (m Graph) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.Key:
		switch msg.String() {
		case "up":
		case "left":
			if m.Selected > 6 {
				m.Selected -= 7
			} else {
				// TODO calculate the square for this day last year
			}
		}
	}
	return m, nil
}

func NewGraph() (Graph, error) {
	var m Graph
	now := time.Now()
	today := now.YearDay()
	year := now.Year()
	aName, _ := commits.GetAuthorName()
	aEmail, _ := commits.GetAuthorEmail()
	authors := []string{aName, aEmail}
	mr, err := commits.GetMRRepos()
	if err != nil {
		return m, err
	}
	m.Repos = mr
	m.Authors = authors
	m.Year = year
	m.Selected = today
	return m, nil
}

func (m Graph) Init() tea.Cmd {
	return nil
}

func (m Graph) View() string {
	mr := commits.RepoSet(m.Repos)
	gfreq, _ := mr.FrequencyChan(m.Year, m.Authors)
	return gfreq.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit

		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.quitting {
		return "\n"
	}
	return m.GraphModel.View()
}

func main() {
	m, err := initialModel()
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
