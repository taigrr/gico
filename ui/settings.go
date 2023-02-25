package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/taigrr/gico/commits"
)

const (
	authors SettingsCursor = iota
	repos
)

type (
	SettingsCursor int
	Settings       struct {
		AllAuthors       map[string]bool
		SelectedAuthors  []string
		AllRepos         map[string]bool
		SelectedRepos    []string
		cursor           SettingsCursor
		highlightedEntry int
		AuthorList       list.Model
		RepoList         list.Model
	}
)

var settingsKey = key.NewBinding(
	key.WithKeys("ctrl+g"),
	key.WithHelp("", "press ctrl+g to open settings"),
)

func (m Settings) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.cursor {
	case authors:
		var cmd tea.Cmd
		m.AuthorList, cmd = m.AuthorList.Update(msg)
		return m, cmd
	case repos:
		var cmd tea.Cmd
		m.RepoList, cmd = m.RepoList.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Settings) Init() tea.Cmd {
	return nil
}

func (m Settings) View() string {
	return fmt.Sprintf("This is the settings view %s", "fmt")
}

func NewSettings() (Settings, error) {
	var m Settings
	var err error
	m.cursor = authors
	allRepos, err := commits.GetMRRepos()
	if err != nil {
		return m, err
	}
	allAuthors, err := commits.RepoSet(allRepos).GetRepoAuthors()
	if err != nil {
		return m, err
	}

	m.AllRepos = make(map[string]bool)
	for _, v := range allRepos {
		m.AllRepos[v] = true
	}

	m.AllAuthors = make(map[string]bool)
	for _, v := range allAuthors {
		m.AllAuthors[v] = false
	}
	m.SelectedRepos = allRepos
	email, _ := commits.GetAuthorEmail()
	if email != "" {
		m.SelectedAuthors = append(m.SelectedAuthors, email)
	}
	name, _ := commits.GetAuthorName()
	if name != "" {
		m.SelectedAuthors = append(m.SelectedAuthors, name)
	}
	for _, v := range m.SelectedRepos {
		m.AllAuthors[v] = true
	}
	return m, nil
}
