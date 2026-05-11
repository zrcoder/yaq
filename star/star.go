package star

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/goplus/ixgo"
	_ "github.com/goplus/ixgo/xgobuild"

	"github.com/zrcoder/yaq"
	"github.com/zrcoder/yaq/common"
	_ "github.com/zrcoder/yaq/exported/github.com/zrcoder/yaq/star/pkg"
	"github.com/zrcoder/yaq/star/pkg"
)

func init() {
	s := &star{Game: pkg.Instance}
	// star mode is the default game of yaq
	yaq.Register("", s)
	yaq.Register("star", s)
}

type star struct {
	*pkg.Game
}

func (s *star) SetBase(base *yaq.Base) {
	s.Base = base
	s.Base.RunCodeAction = s.runCode
}

func (s *star) Run() {
	s.SetSceneSize(s.Rows, s.Columns*3)
	p := tea.NewProgram(s)
	s.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *star) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmdBase := s.Base.Update(msg)
	_, cmd := s.Game.Update(msg)
	return s, tea.Batch(cmdBase, cmd)
}

func (s *star) runCode() {
	go func(code string) {
		_, err := ixgo.RunFile("main.xgo", code, nil, 0)
		if err != nil {
			err = common.ParseBuildError(err, s.PreCode())
			s.Send(err)
		} else {
			s.MarkResult()
		}
	}(s.PreCode() + s.EditorValue())
}
