package yaq

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding
	Run  key.Binding
}

func getCommonKeys() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl", "c"),
			key.WithHelp("ctrl+c", "quit")),
		Run: key.NewBinding(key.WithKeys("ctrl", "r"),
			key.WithHelp("ctrl+r", "run codes")),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Run, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
