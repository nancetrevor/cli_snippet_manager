package models

import (
	"github.com/lfizzikz/snip/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	Title string
	List  list.Model
}

func (m ListModel) Init() tea.Cmd { return nil }

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := ui.DocStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return m.List.View()
}
