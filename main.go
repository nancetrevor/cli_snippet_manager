package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lfizzikz/snip/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenList screen = iota
	screenAdd
)

type rootModel struct {
	active screen
	list   models.ListModel
	add    models.AddModel
}

func (m rootModel) Init() tea.Cmd {
	m.active = screenList
	return m.list.Init()
}
func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case models.SubmitCommandMsg:
		if err := saveCommands(msg.Cmd); err != nil {
			return m, nil
		}
		m.active = screenList
		return m, nil
	}
	switch m.active {
	case screenList:
		updated, cmd := m.list.Update(msg)
		m.list = updated.(models.ListModel)

		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "a" {
			m.active = screenAdd
			return m, m.add.Init()
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
	switch m.active {
	case screenList:
		return m.list.View()
	case screenAdd:
		return m.add.View()
	default:
		return ""
	}
}

func main() {
	commands, err := getSavedCommands()
	if err != nil {
		fmt.Println("error loading commands")
		os.Exit(1)
	}
	space := GetLongestDescription(commands)
	items := []list.Item{}
	for _, cmd := range commands {
		cmd.ODesc = cmd.ODesc + strings.Repeat(" ", space-len(cmd.OUsage)) + cmd.OUsage
		items = append(items, cmd)
	}

	delegate := list.NewDefaultDelegate()
	m := rootModel{
		active: screenList,
		list:   models.NewListModel(items, delegate),
		add:    models.InitialModel(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func GetLongestDescription(list []models.Command) int {
	longest := 0
	for _, c := range list {
		if len(c.ODesc) > longest {
			longest = len(c.ODesc)
		}
	}
	longest += 2
	return longest
}

func saveCommands(f models.Command) error {
	items, err := getSavedCommands()
	if err != nil {
		if os.IsNotExist(err) {
			items = []models.Command{}
		} else {
			return err
		}
	}

	commandFile, err := CommandFileSavePath()
	if err != nil {
		return err
	}

	items = append(items, f)

	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(commandFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(commandFile, out, 0600)
}

func getSavedCommands() ([]models.Command, error) {
	c, err := CommandFileSavePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(c)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Command{}, nil
		}
		return nil, err
	}

	var items []models.Command
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func CommandFileSavePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "commands.json")
	return p, nil
}
