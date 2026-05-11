package turtle

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/goplus/ixgo"

	_ "github.com/zrcoder/yaq/exported/github.com/zrcoder/yaq/games/turtle/pkg"
	"github.com/zrcoder/yaq/games/turtle/pkg"
	yaq "github.com/zrcoder/yaq/pkg"
)

func init() {
	yaq.Register("turtle", &turtle{Game: pkg.Instance})
}

type turtle struct{ *pkg.Game }

func (t *turtle) SetBase(base *yaq.Base) {
	t.Base = base
	t.Base.RunCodeAction = t.runCode
}

func (t *turtle) Run() {
	t.Base.SetSceneSize(t.Rows, t.Columns*3)
	p := tea.NewProgram(t)
	t.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (t *turtle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmdBase := t.Base.Update(msg)
	cmd := t.Game.Update(msg)
	return t, tea.Batch(cmdBase, cmd)
}

func (t *turtle) View() tea.View {
	return t.Base.View(t.Game.View(), t.Hint())
}

func (t *turtle) runCode() {
	preCode := `import . "github.com/zrcoder/yaq/games/turtle/pkg"` + "\n"
	go func(code string) {
		_, err := ixgo.RunFile("main.xgo", code, nil, 0)
		if err != nil {
			err = yaq.ParseBuildError(err, preCode)
			t.Send(err)
		} else {
			t.MarkResult()
		}
	}(preCode + t.EditorValue())
}
