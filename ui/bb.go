package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taigrr/gico/ui/graph"
)

func InteractiveGraph() {
	m := graph.New()
	df := lipgloss.NewDoeFoot()
	m = m.UpdateDoeFoot(df)

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
