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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return s, tea.Quit
		case "ctrl+r":
			s.runCode()
		}
	}
	var cmd tea.Cmd
	_, cmd = s.Game.Update(msg)
	return s, cmd
}

func (s *star) runCode() {
	go func(code string) {
		_, err := ixgo.RunFile("main.gop", code, nil, 0)
		if err != nil {
			err = common.ParseBuildError(err, s.PreCode())
			s.Send(err)
		} else {
			s.MarkResult()
		}
	}(s.PreCode() + s.Editor.Value())
}
