package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/taigrr/gico/commits"
	"github.com/taigrr/gico/graph/term"
)

func xmain() {
	n := time.Now()
	repoPaths, err := commits.GetMRRepos()
	if err != nil {
		panic(err)
	}
	freq, err := repoPaths.Frequency(n.Year(), []string{"Groot"})
	if err != nil {
		panic(err)
	}
	wfreq, err := repoPaths.GetWeekFreq([]string{"Groot"})
	if err != nil {
		panic(err)
	}
	fmt.Println("week:")
	fmt.Println(term.GetWeekUnicode(wfreq))
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("year:")
	fmt.Println(term.GetYearUnicode(freq))
}

type errMsg error

type model struct {
	SettingsModel any
	GraphModel    any
	CommitModel   any
	quitting      bool
	err           error
}

var quitKeys = key.NewBinding(
	key.WithKeys("q", "esc", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit

		}
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	if m.quitting {
		return "\n"
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
