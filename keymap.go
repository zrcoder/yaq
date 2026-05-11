package yaq

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Quit         key.Binding
	Run          key.Binding
	SwitchEditor key.Binding
}

func getCommonKeys() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl", "c"),
			key.WithHelp("ctrl+c", "quit")),
		Run: key.NewBinding(key.WithKeys("ctrl", "r"),
			key.WithHelp("ctrl+r", "run codes")),
		SwitchEditor: key.NewBinding(key.WithKeys("ctrl", "e"),
			key.WithHelp("ctrl+e", "switch editor")),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.SwitchEditor, k.Run, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
