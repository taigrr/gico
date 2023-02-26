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
		AllAuthors       selectablelist
		SelectedAuthors  []string
		AllRepos         selectablelist
		SelectedRepos    []string
		cursor           SettingsCursor
		highlightedEntry int
		AuthorList       list.Model
		RepoList         list.Model
	}
)

type selectablelist []selectable

type selectable struct {
	text     string
	selected bool
}

func (i selectable) Title() string       { return i.text }
func (i selectable) FilterValue() string { return i.text }
func (i selectablelist) GetSelected() []string {
	selected := []string{}
	for _, v := range i {
		if v.selected {
			selected = append(selected, v.text)
		}
	}
	return selected
}

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

	m.AllRepos = []selectable{}
	for _, v := range allRepos {
		m.AllRepos = append(m.AllRepos, selectable{text: v, selected: true})
	}

	m.AllAuthors = []selectable{}
	for _, v := range allAuthors {
		m.AllAuthors = append(m.AllAuthors, selectable{text: v, selected: false})
	}
	m.SelectedRepos = m.AllRepos.GetSelected()
	email, _ := commits.GetAuthorEmail()
	if email != "" {
		m.SelectedAuthors = append(m.SelectedAuthors, email)
	}
	name, _ := commits.GetAuthorName()
	if name != "" {
		m.SelectedAuthors = append(m.SelectedAuthors, name)
	}
	for _, v := range m.SelectedAuthors {
	inner:
		for i, s := range m.AllAuthors {
			if s.text == v {
				m.AllAuthors[i].selected = true
				break inner
			}
		}
	}
	return m, nil
}
