package pkg

import (
	"embed"
	"errors"
	"fmt"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	lp "charm.land/lipgloss/v2"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/yaq/config/turtle"
	"gopkg.in/yaml.v3"

	"github.com/zrcoder/yaq"
	"github.com/zrcoder/yaq/common"
)

var Instance = &Game{}

type Game struct {
	err          error
	currentLevel *Level
	*tea.Program
	*yaq.Base
	pen        *Pen
	Name       string             `yaml:"name"`
	Sprites    map[string]*Sprite `yaml:"sprites"`
	Levels     []string           `yaml:"levels"`
	state      common.State
	totalPoses int
	levelIndex int
	loaded     bool
}

func (g *Game) Init() tea.Cmd {
	err := g.load()
	if err != nil {
		return func() tea.Msg { return err }
	}
	g.levelIndex = 0
	err = g.loadCurrentLevel()
	if err != nil {
		return func() tea.Msg { return err }
	}
	return textarea.Blink
}

func (g *Game) Update(msg tea.Msg) tea.Cmd {
	if g.allFinished() {
		return nil
	}

	switch msg := msg.(type) {
	case common.ErrMsg:
		g.err = msg
		g.state = common.Failed
		return nil
	case tea.KeyMsg:
		g.err = nil
		switch g.state {
		case common.Succeed:
			g.state = common.Running
			g.gotoNextLevel()
		case common.Failed:
			g.state = common.Running
			g.resetLevel()
		default:
		}
	}

	return g.EditorUpdate(msg)
}

func (g *Game) View() string {
	view := ""
	switch {
	case g.allFinished():
		view = dialog.Success("all challenges finished!").String()
	case g.err != nil:
		view = g.ErrorView(g.err.Error())
	case !g.loaded:
		view = g.LoadingView()
	case g.state == common.Succeed:
		view = g.SucceedView("Well done")
	case g.state == common.Failed:
		view = g.ErrorView("failed")
	default:
		view = g.currentLevel.View()
	}
	title := ""
	if !g.allFinished() {
		title := fmt.Sprintf("%s > %s", g.Name, g.Levels[g.levelIndex])
		title += style.Help.Render(fmt.Sprintf("\tLeft: %d\n", g.totalPoses))
	}
	view = lp.JoinVertical(lp.Left,
		title,
		view,
	)
	return view
}

func (g *Game) Hint() string {
	if !g.loaded {
		return ""
	}
	return g.currentLevel.Hint
}

func (g *Game) MarkResult() {
	if g.succeed() {
		g.state = common.Succeed
	} else {
		g.state = common.Failed
	}
	g.Send(common.ResMsg{})
}

func (g *Game) FS() embed.FS {
	return turtle.FS
}

func (g *Game) load() error {
	return yaml.Unmarshal(g.IndexData, g)
}

func (g *Game) loadCurrentLevel() error {
	g.loaded = false
	if len(g.Levels) == 0 {
		return errors.New("no levels found")
	}
	if g.allFinished() {
		return nil
	}

	g.currentLevel = &Level{Game: g}
	data, err := g.FS().ReadFile(g.Levels[g.levelIndex] + common.YamlExt)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, g.currentLevel)
	if err != nil {
		return err
	}
	err = g.currentLevel.initialize()
	if err != nil {
		return err
	}
	g.loaded = true
	return nil
}

func (g *Game) allFinished() bool {
	return g.levelIndex == len(g.Levels)
}

func (g *Game) gotoNextLevel() {
	if g.allFinished() {
		return
	}
	g.levelIndex++
	g.resetLevel()
}

func (g *Game) resetLevel() {
	g.loadCurrentLevel()
}

func (g *Game) outRange(p *common.Position) bool {
	y, x := p.Y, p.X
	return y < 0 || y >= g.Rows || x < 0 || x >= g.Columns
}

func (g *Game) succeed() bool {
	return g.totalPoses == 0
}
