package graph

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taigrr/gico/ui/graph/help"
)

type Graph struct {
	Help help.Help
	// df   lipgloss.DoeFoot
}

func (g Graph) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, _ := g.Help.Update(msg)
		t, _ := h.(help.Help)
		g.Help = t
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return g, tea.Quit
		case "right", "l":
		case "up", "k":
		case "down", "j":
		case "left", "h":
		case "G":
		default:
			h, _ := g.Help.Update(msg)
			t, _ := h.(help.Help)
			g.Help = t
		}
	}
	return g, cmd
}

func (g Graph) Init() tea.Cmd {
	return nil
}

func (g Graph) View() string {
	return ""
}

func New() Graph {
	var g Graph
	g.Help = help.New()
	return g
}

//func (g Graph) UpdateDoeFoot(df lipgloss.DoeFoot) Graph {
//	g.df = df
//	g.Help = g.Help.UpdateDoeFoot(df)
//	return g
//}
