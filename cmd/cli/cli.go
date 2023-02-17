package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/types"
)

const (
	settings Cursor = iota
	graph
	commitLog
)

type (
	Cursor int
	errMsg error
	model  struct {
		SettingsModel  Settings
		GraphModel     Graph
		CommitLogModel CommitLog
		HelpModel      help.Model
		Bindings       []key.Binding
		quitting       bool
		cursor         Cursor
		err            error
	}
	CommitLog struct {
		Year     int
		YearDay  int
		Commits  [][]types.Commit
		Selected int
		Authors  []string
		Repos    []string
	}
	Settings struct{}
	Graph    struct {
		Selected int
		Year     int
		Repos    []string
		Authors  []string
	}
)

var (
	quitKeys = key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("", "press q to quit"),
	)
	settingsKey = key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("", "press ctrl+g to open settings"),
	)
)

func initialModel() (model, error) {
	var m model
	var err error
	m.GraphModel, err = NewGraph()
	if err != nil {
		return m, err
	}
	m.CommitLogModel, err = NewCommitLog()
	if err != nil {
		return m, err
	}
	m.cursor = graph
	m.HelpModel = help.New()
	m.Bindings = []key.Binding{quitKeys, settingsKey}
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
	return fmt.Sprintf("This is the settings view %s", "fmt")
}

func (m CommitLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			if m.Selected < len(m.Commits)-1 {
				m.Selected++
			}
		case "k":
			if m.Selected > 0 {
				m.Selected--
			}
		default:
			mr := commits.RepoSet(m.Repos)
			cis, err := mr.GetRepoCommits(m.Year, m.Authors)
			if err != nil {
				m.Commits = cis
			}
		}
	}
	return m, nil
}

func (m CommitLog) Init() tea.Cmd {
	return nil
}

func (m CommitLog) View() string {
	if len(m.Commits) == 0 {
		return "No commits to display"
	}
	return fmt.Sprintf("%v", m.Commits[m.YearDay])
	// return fmt.Sprintf("This is the Commit Log, selected %v", "sd")
}

func YearLen(year int) int {
	yearLen := 365
	if year%4 == 0 {
		yearLen++
	}
	return yearLen
}

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

func NewCommitLog() (CommitLog, error) {
	var m CommitLog
	now := time.Now()
	today := now.YearDay() - 1
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
	m.Commits, err = mr.GetRepoCommits(m.Year, m.Authors)
	return m, err
}

func NewGraph() (Graph, error) {
	var m Graph
	now := time.Now()
	today := now.YearDay() - 1
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch cmd := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(cmd, settingsKey) {
			switch m.cursor {
			case settings:
				m.cursor = graph
			default:
				m.cursor = settings
			}
		}
		if key.Matches(cmd, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		}
	case errMsg:
		m.err = cmd
		return m, nil

	default:
	}
	switch m.cursor {
	// multiple cursors defined for extensibility, but only graph is used
	case graph, commitLog:
		tmp, _ := m.GraphModel.Update(msg)
		m.GraphModel, _ = tmp.(Graph)
		tmpC, cmd := m.CommitLogModel.Update(msg)
		m.CommitLogModel, _ = tmpC.(CommitLog)
		m.CommitLogModel.YearDay = m.GraphModel.Selected
		return m, cmd
	case settings:
		tmp, cmd := m.SettingsModel.Update(msg)
		m.SettingsModel, _ = tmp.(Settings)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.quitting {
		return ""
	}
	return lipgloss.JoinVertical(lipgloss.Top, m.GraphModel.View(), m.CommitLogModel.View(), m.HelpModel.ShortHelpView(m.Bindings))
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
