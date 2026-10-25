package keymap

import "charm.land/bubbles/v2/key"

type GlobalKeyMap struct {
	Quit          key.Binding
	Help          key.Binding
	Back          key.Binding
	StopContainer key.Binding
	EnterFlag     key.Binding
}

func DefaultGlobalKeyMap() GlobalKeyMap {
	return GlobalKeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "ctrl+d"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		StopContainer: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "stop container"),
		),
		EnterFlag: key.NewBinding(
			key.WithKeys("ctrl+f"),
			key.WithHelp("ctrl+f", "enter flag"),
		),
	}
}
