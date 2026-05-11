package hanoi

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/goplus/ixgo"

	_ "github.com/zrcoder/yaq/exported/github.com/zrcoder/yaq/games/hanoi/pkg"
	"github.com/zrcoder/yaq/games/hanoi/pkg"
	yaq "github.com/zrcoder/yaq/pkg"
)

func init() {
	yaq.Register("hanoi", &hanoi{Game: pkg.Instance})
}

type hanoi struct {
	*pkg.Game
}

func (h *hanoi) SetBase(base *yaq.Base) {
	h.Base = base
	h.Base.RunCodeAction = h.runCode
}

func (h *hanoi) Run() {
	h.SetSceneSize(20, 60)
	p := tea.NewProgram(h)
	h.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (h *hanoi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmdBase := h.Base.Update(msg)
	cmd := h.Game.Update(msg)
	return h, tea.Batch(cmdBase, cmd)
}
func (g *hanoi) View() tea.View {
	return g.Base.View(g.Game.View(), g.Hint())
}

func (h *hanoi) runCode() {
	precode := `import . "github.com/zrcoder/yaq/games/hanoi/pkg"` + "\n"
	go func(code string) {
		_, err := ixgo.RunFile("main.xgo", code, nil, 0)
		if err != nil {
			err = yaq.ParseBuildError(err, precode)
			h.Send(err)
		} else {
			h.MarkResult()
		}
	}(precode + h.EditorValue())
}
