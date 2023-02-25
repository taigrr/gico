package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
)

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func InitialModel() (model, error) {
	var m model
	var err error
	m.SettingsModel, err = NewSettings()
	if err != nil {
		return m, err
	}
	m.GraphModel, err = NewGraph(m.SettingsModel.SelectedAuthors, m.SettingsModel.SelectedRepos)
	if err != nil {
		return m, err
	}
	m.CommitLogModel, err = NewCommitLog()
	if err != nil {
		return m, err
	}
	m.cursor = graph
	m.HelpModel = help.New()
	m.Bindings = []key.Binding{
		quitKeys,
		settingsKey,
		m.CommitLogModel.Table.KeyMap.LineDown,
		m.CommitLogModel.Table.KeyMap.LineUp,
		m.CommitLogModel.Table.KeyMap.PageUp,
		m.CommitLogModel.Table.KeyMap.PageDown,
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func YearLen(year int) int {
	yearLen := 365
	if year%4 == 0 {
		yearLen++
	}
	return yearLen
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
		m.GraphModel.Authors = m.SettingsModel.SelectedAuthors
		m.GraphModel.Repos = m.SettingsModel.SelectedRepos

		m.CommitLogModel.Authors = m.SettingsModel.SelectedAuthors
		m.CommitLogModel.Repos = m.SettingsModel.SelectedRepos

		tmp, _ := m.GraphModel.Update(msg)
		m.GraphModel, _ = tmp.(Graph)

		m.CommitLogModel.Year = m.GraphModel.Year
		if m.CommitLogModel.YearDay != m.GraphModel.Selected {
			m.CommitLogModel.YearDay = m.GraphModel.Selected
			m.CommitLogModel.Table.SetCursor(0)
		}
		tmpC, cmd := m.CommitLogModel.Update(msg)
		m.CommitLogModel, _ = tmpC.(CommitLog)
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
	mainView := lipgloss.JoinVertical(lipgloss.Top, m.GraphModel.View(), m.CommitLogModel.View(), m.HelpModel.ShortHelpView(m.Bindings))
	if m.cursor == settings {
		return lipgloss.JoinHorizontal(lipgloss.Left, m.SettingsModel.View(), mainView)
	}
	return mainView
}

func main() {
	m, err := InitialModel()
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
