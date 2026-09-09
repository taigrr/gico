package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/taigrr/gico/cmd/internal/cli"
	"github.com/taigrr/gico/ui"
)

func main() {
	m, err := ui.InitialModel()
	if err != nil {
		cli.Fatal("failed to initialize gico", err)
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		cli.Fatal("failed to run gico", err)
	}
}
