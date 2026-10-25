package challenges

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	UpDown key.Binding
	Select key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		UpDown: key.NewBinding(
			key.WithKeys("up", "down", "k", "j"),
			key.WithHelp("↑↓ j/k", "navigate"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}
