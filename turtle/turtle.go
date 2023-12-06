package turtle

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zrcoder/yaq"
	"github.com/zrcoder/yaq/turtle/pkg"
)

func init() {
	yaq.Register("turtle", &Turtle{Game: pkg.Instace})
}

type Turtle struct{ *pkg.Game }

func (t *Turtle) SetBase(base *yaq.Base) {
	t.Base = base
}

func (t *Turtle) Run() {
	fmt.Println("TODO: turtle graphics")
	p := tea.NewProgram(t, tea.WithAltScreen())
	t.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (t *Turtle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return t, tea.Quit
		case tea.KeyCtrlR:
			go t.runCode()
		}
	}
	var cmd tea.Cmd
	_, cmd = t.Game.Update(msg)
	return t, cmd
}

func (t *Turtle) runCode() {
}
