package pkg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	lp "github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"

	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/yaq"
	"github.com/zrcoder/yaq/common"
)

var Instance = &Game{}

type Game struct {
	err error
	*tea.Program
	*yaq.Base
	player      *Sprite
	successInfo string
	state       common.State
	scenes      []*Scene `yaml:"_"`
	SceneNames  []string `yaml:"scenes"`
	sceneIndex  int
	totalStars  int
	loaded      bool
}

func (g *Game) Init() tea.Cmd {
	if err := g.load(); err != nil {
		return func() tea.Msg { return err }
	}
	g.Editor.SetHeight(g.Rows)
	g.Editor.SetWidth(g.Columns * 3)
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
	g.successInfo = "Well done!"
	if !g.allFinished() {
		g.successInfo = g.currentLevel().SuccessMsg
	}
	var cmd tea.Cmd
	g.Editor, cmd = g.Editor.Update(msg)
	return g, cmd
}

func (g *Game) View() string {
	if g.allFinished() {
		return dialog.Success("all challenges finished!").String()
	}

	title := fmt.Sprintf("%s > %s > %s", g.Name, g.currentScene().name, g.currentLevel().name)
	title += style.Help.Render(fmt.Sprintf("\tLeft: %d\n", g.totalStars))
	leftView := ""
	switch {
	case g.err != nil:
		leftView = g.ErrorView(g.err.Error())
	case !g.loaded:
		leftView = g.LoadingView()
	case g.state == common.Succeed:
		leftView = g.SucceedView(g.successInfo)
	case g.state == common.Failed:
		leftView = g.ErrorView("failed")
	default:
		leftView = g.currentLevel().view()
	}
	leftView = lp.JoinVertical(lp.Left, title, leftView)
	rightView := lp.JoinVertical(lp.Left,
		style.Help.Render(g.currentLevel().Hint), "",
		g.Editor.View(), "",
		g.KeysView())
	return lp.JoinHorizontal(lp.Top, leftView, "  ", rightView)
}

func (g *Game) PreCode() string {
	return g.currentLevel().preCode
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
	if err := yaml.Unmarshal(g.IndexData, g); err != nil {
		return err
	}
	return g.loadScenes()
}

func (g *Game) loadScenes() error {
	if len(g.SceneNames) == 0 {
		return errors.New("no scenes in config")
	}

	g.scenes = make([]*Scene, len(g.SceneNames))
	for i, name := range g.SceneNames {
		s := &Scene{}
		path := filepath.Join(g.CfgPath, name, common.IndexFile)
		if data, err := os.ReadFile(path); err != nil {
			return err
		} else if err = yaml.Unmarshal(data, s); err != nil {
			return err
		}
		s.Game = g
		s.name = name
		s.bgColors[0] = s.BgColor1
		s.bgColors[1] = s.BgColor2
		for k, sp := range s.Sprites {
			sp.key = k
		}
		g.scenes[i] = s
	}
	g.sceneIndex = 0
	return g.currentScene().loadLevels()
}

func (g *Game) currentScene() *Scene { return g.scenes[g.sceneIndex] }
func (g *Game) currentLevel() *Level { return g.currentScene().currentLevel() }

func (g *Game) succeed() bool {
	return g.err == nil && g.totalStars == 0
}

func (g *Game) gotoNextLevel() {
	if g.allFinished() {
		return
	}
	g.state = common.Running
	if len(g.currentScene().LevelNames) == g.currentScene().levelIndex+1 {
		g.sceneIndex++
		if g.sceneIndex == len(g.scenes) {
			return
		}
		g.currentScene().levelIndex = 0
		if err := g.currentScene().loadLevels(); err != nil {
			g.err = err
		}
		return
	}
	g.currentScene().levelIndex++
	g.err = g.currentScene().loadCurrentLevel()
}

func (g *Game) resetLevel() {
	g.state = common.Running
	g.err = g.currentLevel().initialize()
}

func (g *Game) allFinished() bool {
	return len(g.scenes) == 0 ||
		g.sceneIndex == len(g.scenes) ||
		len(g.currentScene().levels) == 0 ||
		g.sceneIndex+1 == len(g.scenes) && g.currentScene().levelIndex == len(g.currentScene().levels)
}
