package turtle

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goplus/igop"

	"github.com/zrcoder/yaq"
	"github.com/zrcoder/yaq/common"
	_ "github.com/zrcoder/yaq/exported/github.com/zrcoder/yaq/turtle/pkg"
	"github.com/zrcoder/yaq/turtle/pkg"
)

func init() {
	yaq.Register("turtle", &turtle{Game: pkg.Instance})
}

type turtle struct{ *pkg.Game }

func (t *turtle) SetBase(base *yaq.Base) {
	t.Base = base
}

func (t *turtle) Run() {
	t.Base.SetSceneSize(t.Rows, t.Columns*3)
	p := tea.NewProgram(t, tea.WithAltScreen())
	t.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (t *turtle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s *turtle) runCode() {
	preCode := `import . "github.com/zrcoder/yaq/turtle/pkg"` + "\n"
	go func(code string) {
		_, err := igop.RunFile("main.gop", code, nil, 0)
		if err != nil {
			err = common.ParseBuildError(err, preCode)
			s.Send(err)
		} else {
			s.MarkResult()
		}
	}(preCode + s.Editor.Value())
}
