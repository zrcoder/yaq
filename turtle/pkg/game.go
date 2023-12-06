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

var Instace = &Game{}

type Game struct {
	err          error
	currentLevel *Level
	*tea.Program
	*yaq.Base
	pen        *Pen
	Sprites    map[string]*Sprite `toml:"sprites"`
	state      common.State
	Levels     []string `toml:"levels"`
	totalPoses int
	Rows       int `toml:"rows"`
	Columns    int `toml:"columns"`
	levelIndex int
	loaded     bool
}

type errMsg error

func (g *Game) Init() tea.Cmd {
	err := g.load()
	if err != nil {
		g.err = err
		return nil
	}
	if len(g.Levels) == 0 {
		return func() tea.Msg { return errMsg(errors.New("no levels found")) }
	}
	g.levelIndex = 0
	g.loadCurrentLevel()
	return textarea.Blink
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if g.allFinished() {
		return g, nil
	}
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case errMsg:
		g.err = msg
		g.state = common.Failed
		g.Editor.Blur()
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
			if !g.Editor.Focused() {
				cmd = g.Editor.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}
	g.Editor, cmd = g.Editor.Update(msg)
	cmds = append(cmds, cmd)
	return g, tea.Batch(cmds...)
}

func (g *Game) View() string {
	if g.allFinished() {
		return dialog.Success("all challenges finished!").String()
	}

	title := fmt.Sprintf("%s > %s", g.Name, g.Levels[g.levelIndex])
	leftView := ""
	switch {
	case g.err != nil:
		leftView = g.ErrorView(g.err.Error())
	case !g.loaded:
		leftView = "loading"
	case g.state == common.Succeed:
		leftView = g.SucceedView("Well done")
	case g.state == common.Failed:
		leftView = g.ErrorView("failed")
	default:
		leftView = g.currentLevel.View()
	}
	leftView = lp.JoinVertical(lp.Left, title, leftView)
	rightView := lp.JoinVertical(lp.Left,
		style.Help.Render(g.currentLevel.Hint), "",
		g.Editor.View(), "",
		g.KeysView())
	return lp.JoinHorizontal(lp.Top, leftView, "   ", rightView)
}

func (g *Game) load() error {
	data, err := os.ReadFile(filepath.Join(g.CfgPath, common.IndexFile))
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, g)
}

func (g *Game) loadCurrentLevel() {
	lvl := &Level{}
	data, err := os.ReadFile(filepath.Join(g.CfgPath, g.Levels[g.levelIndex]+common.TomlExt))
	if err != nil {
		g.err = err
		return
	}
	err = toml.Unmarshal(data, lvl)
	if err != nil {
		g.err = err
		return
	}
	err = lvl.initialize()
	if err != nil {
		g.err = err
		return
	}
	g.currentLevel = lvl
	g.loaded = true
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
