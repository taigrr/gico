package help

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Help struct {
	df lipgloss.DoeFoot
}

func (h Help) Update(m tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}
func (h Help) Init() tea.Cmd {
	return nil
}

func (h Help) View() string {
	return ""
}

func New() Help {
	return Help{}
}

func (h Help) UpdateDoeFoot(df lipgloss.DoeFoot) Help {
	h.df = df
	return h
}
