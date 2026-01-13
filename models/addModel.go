package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddModel struct {
	name        string
	description string
	usage       string
	inputs      []textinput.Model
	focused     int
	err         error
}

type Command struct {
	OName  string `json:"name"`
	ODesc  string `json:"desc"`
	OUsage string `json:"usage"`
}

func (i Command) Title() string       { return i.OName }
func (i Command) Description() string { return i.ODesc }
func (i Command) Usage() string       { return i.OUsage }
func (i Command) FilterValue() string { return i.OName }

type SubmitCommandMsg struct {
	Cmd Command
}

type (
	errMsg error
)

const (
	name = iota
	desc
	usage
)
const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

func InitialModel() AddModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[name] = textinput.New()
	inputs[name].Placeholder = "cd"
	inputs[name].Focus()
	inputs[name].CharLimit = 85
	inputs[name].Width = 95
	inputs[name].Prompt = ""

	inputs[desc] = textinput.New()
	inputs[desc].Placeholder = "for changing directories"
	inputs[desc].CharLimit = 85
	inputs[desc].Width = 95
	inputs[desc].Prompt = ""

	inputs[usage] = textinput.New()
	inputs[usage].Placeholder = "cd {folder name}"
	inputs[usage].CharLimit = 85
	inputs[usage].Width = 95
	inputs[usage].Prompt = ""

	return AddModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}
func (m AddModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				c := Command{
					OName:  m.inputs[name].Value(),
					ODesc:  m.inputs[desc].Value(),
					OUsage: m.inputs[usage].Value(),
				}
				return m, func() tea.Msg { return SubmitCommandMsg{Cmd: c} }
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}
func (m AddModel) View() string {
	return fmt.Sprintf(
		` Input new command

 %s
 %s

 %s  
 %s  

 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Command name:"),
		m.inputs[name].View(),
		inputStyle.Width(30).Render("Command description:"),
		m.inputs[desc].View(),
		inputStyle.Width(30).Render("Command usage:"),
		m.inputs[usage].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *AddModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *AddModel) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
