package main

import (
	"fmt"
	"os"

	"github.com/lfizzikz/snip/internal/handlers"
	"github.com/lfizzikz/snip/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

const (
	screenList screen = iota
	screenAdd
)

type rootModel struct {
	active        screen
	list          models.ListModel
	add           models.AddModel
	confirmDelete bool
	deleteIndex   int
	deleteName    string

	width  int
	height int
}

func (m rootModel) Init() tea.Cmd {
	m.active = screenList
	return m.list.Init()
}
func (m *rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case models.SubmitCommandMsg:
		m.list.AddCommand(msg.Cmd)
		if err := handlers.AppendCommands([]models.Command{msg.Cmd}); err != nil {
			return m, nil
		}
		m.active = screenList
		return m, nil
	}

	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.height = ws.Height
	}
	if m.confirmDelete {
		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "y", "enter":
				m.list.RemoveAt(m.deleteIndex)
				_ = handlers.DeleteCommandByName(m.deleteName)

				m.confirmDelete = false
				return m, nil
			case "n", "esc":
				m.confirmDelete = false
				return m, nil
			}
		}
		return m, nil
	}

	switch m.active {

	case screenList:
		updated, cmd := m.list.Update(msg)
		m.list = *updated.(*models.ListModel)

		if key, ok := msg.(tea.KeyMsg); ok {
			switch key.String() {
			case "a":
				m.add = models.InitialModel()
				m.active = screenAdd
				return m, m.add.Init()
			case "d":
				m.confirmDelete = true
				m.deleteIndex = m.list.Index()

				if it, ok := m.list.List.SelectedItem().(interface{ Title() string }); ok {
					m.deleteName = it.Title()
				} else {
					m.deleteName = ""
				}
				return m, nil
			}
		}
		return m, cmd

	case screenAdd:
		updated, cmd := m.add.Update(msg)
		m.add = updated.(models.AddModel)

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "esc" {
			m.active = screenList
			return m, nil
		}
		return m, cmd
	}
	return m, nil
}
func (m rootModel) View() string {
	var base string
	switch m.active {
	case screenList:
		base = m.list.View()
	case screenAdd:
		base = m.add.View()
	default:
		return ""
	}

	if !m.confirmDelete {
		return base
	}

	dialog := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		Render(
			"Delete \"" + m.deleteName + "\"?\n\n" +
				"          y/n          ",
		)
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

}

func main() {
	commands, err := handlers.GetSavedCommands()
	if err != nil {
		fmt.Printf("error loading commands %v", err)
		os.Exit(1)
	}
	items := []list.Item{}
	for _, cmd := range commands {
		items = append(items, cmd)
	}

	delegate := models.NewItemDelegate()
	m := rootModel{
		active: screenList,
		list:   models.NewListModel(items, delegate),
		add:    models.InitialModel(),
	}

	p := tea.NewProgram(&m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
