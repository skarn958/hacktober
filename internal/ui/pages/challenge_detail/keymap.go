package challenge_detail

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Spawn  key.Binding
	Top    key.Binding
	Bottom key.Binding
}

func defaultKeyMap() KeyMap {
	return KeyMap{
		Spawn: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "start challenge"),
		),
		Top: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "goto top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "goto bottom"),
		),
	}
}
