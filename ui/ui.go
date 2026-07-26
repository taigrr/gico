package ui

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

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
)

var quitKeys = key.NewBinding(
	key.WithKeys("q", "ctrl+c"),
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

// YearLen returns the number of days in a year.
// Deprecated: Use types.YearLength instead.
func YearLen(year int) int {
	return types.YearLength(year)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch cmd := msg.(type) {

	case tea.KeyPressMsg:
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
	var b tea.BatchMsg
	switch m.cursor {
	// multiple cursors defined for extensibility, but only graph is used
	case graph, commitLog:
		m.GraphModel.Authors = m.SettingsModel.SelectedAuthors
		m.GraphModel.Repos = m.SettingsModel.SelectedRepos

		m.CommitLogModel.Authors = m.SettingsModel.SelectedAuthors
		m.CommitLogModel.Repos = m.SettingsModel.SelectedRepos
		var c tea.Cmd
		m.GraphModel, c = m.GraphModel.Update(msg)
		b = append(b, c)

		m.CommitLogModel.Year = m.GraphModel.Year
		if m.CommitLogModel.YearDay != m.GraphModel.Selected {
			m.CommitLogModel.YearDay = m.GraphModel.Selected
			m.CommitLogModel.Table.SetCursor(0)
		}
		var cmd tea.Cmd
		m.CommitLogModel, cmd = m.CommitLogModel.Update(msg)
		b = append(b, cmd)
		// Forward non-key messages (e.g. window resize) to the settings model
		// so its lists stay sized, but do not let key presses reach the
		// settings toggle logic while the graph is shown.
		if _, isKey := msg.(tea.KeyPressMsg); !isKey {
			var scmd tea.Cmd
			m.SettingsModel, scmd = m.SettingsModel.Update(msg)
			b = append(b, scmd)
		}
		return m, tea.Batch(b...)
	case settings:
		var scmd tea.Cmd
		m.SettingsModel, scmd = m.SettingsModel.Update(msg)
		b = append(b, scmd)
		return m, tea.Batch(b...)

	}
	return m, nil
}

func (m model) View() tea.View {
	if m.err != nil {
		return tea.NewView(m.err.Error())
	}
	if m.quitting {
		return tea.NewView("")
	}
	if m.cursor == settings {
		return tea.NewView(m.SettingsModel.View())
	}
	return tea.NewView(lipgloss.JoinVertical(lipgloss.Top,
		m.GraphModel.View(),
		m.CommitLogModel.View(),
		m.HelpModel.ShortHelpView(m.Bindings),
	))
}
