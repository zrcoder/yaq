package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	lp "github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/style"
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
	Name       string             `toml:"name"`
	Sprites    map[string]*Sprite `toml:"sprites"`
	Levels     []string           `toml:"levels"`
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

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if g.allFinished() {
		return g, nil
	}

	switch msg := msg.(type) {
	case common.ErrMsg:
		g.err = msg
		g.state = common.Failed
		return g, nil
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

	var cmd tea.Cmd
	g.Editor, cmd = g.Editor.Update(msg)
	return g, cmd
}

func (g *Game) View() string {
	if g.allFinished() {
		return dialog.Success("all challenges finished!").String()
	}

	title := fmt.Sprintf("%s > %s", g.Name, g.Levels[g.levelIndex])
	title += style.Help.Render(fmt.Sprintf("\tLeft: %d\n", g.totalPoses))
	leftView := ""
	switch {
	case g.err != nil:
		leftView = g.ErrorView(g.err.Error())
	case !g.loaded:
		leftView = g.LoadingView()
	case g.state == common.Succeed:
		leftView = g.SucceedView("Well done")
	case g.state == common.Failed:
		leftView = g.ErrorView("failed")
	default:
		leftView = g.currentLevel.View()
	}
	leftView = lp.JoinVertical(lp.Left,
		title,
		leftView,
	)
	hintView := ""
	if g.loaded {
		hintView = style.Help.Render(g.currentLevel.Hint)
	}
	rightView := lp.JoinVertical(lp.Left,
		hintView, "",
		g.Editor.View(), "",
		g.KeysView())
	return lp.JoinHorizontal(lp.Top, leftView, "   ", rightView)
}

func (g *Game) MarkResult() {
	if g.succeed() {
		g.state = common.Succeed
	} else {
		g.state = common.Failed
	}
	g.Send(common.ResMsg{})
}

func (g *Game) load() error {
	return toml.Unmarshal(g.IndexData, g)
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
	data, err := os.ReadFile(filepath.Join(g.CfgPath, g.Levels[g.levelIndex]+common.TomlExt))
	if err != nil {
		return err
	}
	err = toml.Unmarshal(data, g.currentLevel)
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
