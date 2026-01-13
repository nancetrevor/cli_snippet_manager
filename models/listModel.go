package models

import (
	"fmt"
	"io"

	"github.com/lfizzikz/snip/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListModel struct {
	Title string
	List  list.Model
}

func (m ListModel) Init() tea.Cmd { return nil }

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *ListModel) AddCommand(c Command) {
	m.List.InsertItem(len(m.List.Items()), c)
}

func (m *ListModel) Index() int              { return m.List.Index() }
func (m *ListModel) SelectedItem() list.Item { return m.List.SelectedItem() }
func (m *ListModel) RemoveAt(i int) {
	m.List.RemoveItem(i)
}

type itemUsage interface {
	Usage() string
}

type itemDisplay interface {
	Title() string
	Description() string
}

type itemDelegate struct {
	titleStyle lipgloss.Style
	descStyle  lipgloss.Style
	usageStyle lipgloss.Style
}

func NewItemDelegate() itemDelegate {
	return itemDelegate{
		titleStyle: lipgloss.NewStyle().Bold(true),
		descStyle:  lipgloss.NewStyle(),
		usageStyle: lipgloss.NewStyle(),
	}
}

func (d itemDelegate) Height() int                               { return 2 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	di, ok := listItem.(itemDisplay)
	title := ""
	desc := ""
	if ok {
		title = d.titleStyle.Render(di.Title())
		desc = d.descStyle.Render(di.Description())
	}

	iu, ok := listItem.(itemUsage)
	usage := ""
	if ok {
		usage = d.usageStyle.Render(iu.Usage())
	}

	titleText := di.Title()
	descText := di.Description()

	isSelected := index == m.Index()

	selectedRow := lipgloss.NewStyle().
		Background(lipgloss.Color("#36385A")).
		Foreground(lipgloss.Color("#E6E6F0"))

	gap := 8
	gapStr := lipgloss.NewStyle().Width(gap).Render("")
	if isSelected {
		gapStr = selectedRow.Width(gap).Height(2).Render("")

	}
	leftW := max(lipgloss.Width(titleText), lipgloss.Width(descText))
	available := m.Width() - leftW - gap
	rightW := max(available, 0)

	titleStyle := d.titleStyle
	descStyle := d.descStyle
	if isSelected {
		titleStyle = d.titleStyle.Background(lipgloss.Color("#36385A")).Foreground(lipgloss.Color("#E6E6F0"))
		descStyle = d.descStyle.Background(lipgloss.Color("#36385A")).Foreground(lipgloss.Color("#E6E6F0"))
	}

	title = titleStyle.Width(leftW).Render(titleText)
	desc = descStyle.Width(leftW).Render(descText)
	left := lipgloss.JoinVertical(lipgloss.Left, title, desc)

	right := lipgloss.NewStyle().
		Width(rightW).
		Height(2).
		Align(lipgloss.Right).
		Render(usage)
	if isSelected {
		right = selectedRow.
			Width(rightW).
			Height(2).
			Align(lipgloss.Right).
			Render(usage)
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		left,
		gapStr,
		right,
	)

	if index == m.Index() {
		row = selectedRow.
			Width(m.Width() + gap).
			Height(2).
			Render(row)
	}

	fmt.Fprint(w, row)
}
