package ui

import (
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectionMsg struct {
	IsSelected bool
}

type selectableDelegate struct {
	IsActiveList bool
}

func (s selectableDelegate) Height() int { return 1 }

func (s selectableDelegate) Spacing() int {
	return 1
}

func (s selectableDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if msg, ok := msg.(SelectionMsg); ok {
		s.IsActiveList = msg.IsSelected
		m.SetDelegate(s)
	}
	return nil
}

func (s selectableDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	x, ok := item.(selectable)
	if !ok {
		return
	}
	str := ""
	if x.selected {
		str += " [X] "
	} else {
		str += " [ ] "
	}
	str += x.text
	if s.IsActiveList && m.Index() == index {
		sty := list.NewDefaultItemStyles()
		str = sty.SelectedTitle.Render(str)
	}
	w.Write([]byte(str))
}

type delegateKeyMap struct {
	toggle key.Binding
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "choose"),
		),
	}
}
