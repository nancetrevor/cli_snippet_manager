package models

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type customKeyMap struct {
	Add    key.Binding
	Delete key.Binding
}

var keys = customKeyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add command"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete command"),
	),
}

func NewListModel(items []list.Item, delegate list.ItemDelegate) ListModel {
	l := list.New(items, delegate, 0, 0)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Add,
			keys.Delete,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.Add,
			keys.Delete,
		}
	}
	l.Title = "Saved Snips"
	return ListModel{List: l}
}
