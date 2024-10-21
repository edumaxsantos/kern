package main

import "github.com/charmbracelet/bubbles/key"

type MyKeyMap struct {
	InputView key.Binding
}

func (k MyKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.InputView}
}

func (k MyKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.InputView},
	}
}

var keys = MyKeyMap{
	InputView: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "Open pin definition"),
	),
}
