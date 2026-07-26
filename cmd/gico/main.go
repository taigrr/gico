package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/taigrr/gico/ui"
)

func main() {
	m, err := ui.InitialModel()
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
