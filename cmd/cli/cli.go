package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/types"
)

const (
	settings Cursor = iota
	graph
	commitLog

	authors SettingsCursor = iota
	repos
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
		Year    int
		YearDay int
		Commits [][]types.Commit
		Authors []string
		Repos   []string
		Table   table.Model
	}
	SettingsCursor int
	Settings       struct {
		AllAuthors      []string
		SelectedAuthors []string
		AllRepos        []string
		SelectedRepos   []string
		cursor          SettingsCursor
		AuthorList      list.Model
		RepoList        list.Model
	}
	Graph struct {
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
		case "j", "k", "b", "n", "pgdown", "pgup":
		default:
			mr := commits.RepoSet(m.Repos)
			cis, _ := mr.GetRepoCommits(m.Year, m.Authors)
			m.Commits = cis
			commits := m.Commits[m.YearDay]
			rows := []table.Row{}
			for _, c := range commits {
				repo := filepath.Base(c.Repo)
				r := table.Row{
					c.TimeStamp.Format("0" + time.Kitchen),
					repo,
					c.Author.Name,
					c.Message,
				}
				rows = append(rows, r)
			}
			m.Table.SetRows(rows)
			m.Table.SetCursor(0)
		}
	}
	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func newTable() table.Model {
	t := table.New()
	t.SetColumns([]table.Column{
		{Title: "Time", Width: 8},
		{Title: "Repository", Width: 20},
		{Title: "Author", Width: 15},
		{Title: "Message", Width: 40},
	})
	t.KeyMap.LineUp = key.NewBinding(key.WithKeys("k"),
		key.WithHelp("k", "move up one commit"))
	t.KeyMap.LineDown = key.NewBinding(key.WithKeys("j"),
		key.WithHelp("j", "move down one commit"))
	t.Focus()
	return t
}

func (m CommitLog) Init() tea.Cmd {
	return nil
}

func (m CommitLog) View() string {
	if len(m.Commits) == 0 {
		return "No commits to display"
	}

	if len(m.Commits[m.YearDay]) == 0 {
		return "No commits to display"
	}
	var b strings.Builder
	b.WriteString("\nCommit Log\n\n")
	b.WriteString(m.Table.View())
	return b.String()
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
	aName, err := commits.GetAuthorName()
	if err != nil {
		return m, err
	}
	aEmail, err := commits.GetAuthorEmail()
	if err != nil {
		return m, err
	}
	mr, err := commits.GetMRRepos()
	if err != nil {
		return m, err
	}
	m.Authors = []string{aName, aEmail}
	m.Repos = mr
	m.Year = year
	m.YearDay = today
	m.Table = newTable()
	{
		cis, _ := mr.GetRepoCommits(m.Year, m.Authors)
		m.Commits = cis
		commits := m.Commits[m.YearDay]
		rows := []table.Row{}
		for _, c := range commits {
			repo := filepath.Base(c.Repo)
			r := table.Row{
				c.TimeStamp.Format("0" + time.Kitchen),
				repo,
				c.Author.Name,
				c.Message,
			}
			rows = append(rows, r)
		}
		m.Table.SetRows(rows)
	}
	m.Commits, err = mr.GetRepoCommits(m.Year, m.Authors)
	if err != nil {
		return m, err
	}
	m.Table.SetCursor(0)
	return m, err
}

func NewSettings() (Settings, error) {
	var m Settings
	var err error
	m.cursor = authors
	m.AllRepos, err = commits.GetMRRepos()
	if err != nil {
		return m, err
	}
	m.SelectedRepos = m.AllRepos
	m.AllAuthors, err = commits.RepoSet(m.AllRepos).GetRepoAuthors()
	if err != nil {
		return m, err
	}
	email, _ := commits.GetAuthorEmail()
	if email != "" {
		m.SelectedAuthors = append(m.SelectedAuthors, email)
	}
	name, _ := commits.GetAuthorName()
	if name != "" {
		m.SelectedAuthors = append(m.SelectedAuthors, name)
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
