package star

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/goplus/igop"
	_ "github.com/goplus/igop/gopbuild"

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
	p := tea.NewProgram(s, tea.WithAltScreen())
	s.Program = p
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *star) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyCtrlR:
			s.runCode()
		}
	}
	var cmd tea.Cmd
	_, cmd = s.Game.Update(msg)
	return s, cmd
}

func (s *star) runCode() {
	go func(code string) {
		_, err := igop.RunFile("main.gop", code, nil, 0)
		if err != nil {
			err = common.ParseBuildError(err, s.PreCode())
			s.Send(err)
		} else {
			s.MarkResult()
		}
	}(s.PreCode() + s.Editor.Value())
}
