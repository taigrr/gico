package ui

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/types"
)

type (
	CommitLog struct {
		Year    int
		YearDay int
		Commits [][]types.Commit
		Authors []string
		Repos   []string
		Table   table.Model
	}
)

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
