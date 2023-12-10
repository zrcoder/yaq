package hanoi

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goplus/igop"

	"github.com/zrcoder/yaq"
	"github.com/zrcoder/yaq/common"
	_ "github.com/zrcoder/yaq/exported/github.com/zrcoder/yaq/hanoi/pkg"
	"github.com/zrcoder/yaq/hanoi/pkg"
)

func init() {
	h := &hanoi{Game: pkg.Instance}
	yaq.Register("hanoi", h)
}

type hanoi struct {
	*pkg.Game
}

func (h *hanoi) SetBase(base *yaq.Base) {
	h.Base = base
}

func (h *hanoi) Run() {
	h.SetSceneSize(20, 60)
	p := tea.NewProgram(h, tea.WithAltScreen())
	h.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (h *hanoi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return h, tea.Quit
		case tea.KeyCtrlR:
			h.runCode()
		}
	}
	var cmd tea.Cmd
	_, cmd = h.Game.Update(msg)
	return h, cmd
}

func (h *hanoi) runCode() {
	precode := `import . "github.com/zrcoder/yaq/hanoi/pkg"` + "\n"
	go func(code string) {
		_, err := igop.RunFile("main.gop", code, nil, 0)
		if err != nil {
			err = common.ParseBuildError(err, precode)
			h.Send(err)
		} else {
			h.MarkResult()
		}
	}(precode + h.Editor.Value())
}
