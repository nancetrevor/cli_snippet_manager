package models

import (
	"github.com/charmbracelet/bubbles/list"
)

func NewListModel(items []list.Item, delegate list.ItemDelegate) ListModel {
	l := list.New(items, delegate, 0, 0)
	l.Title = "My CLI Commands"
	return ListModel{List: l}
}
