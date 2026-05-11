package pkg

import (
	"embed"
	"errors"
	"fmt"
	"math/rand"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	lp "charm.land/lipgloss/v2"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"
	"github.com/zrcoder/yaq/common"
	"github.com/zrcoder/yaq/config/hanoi"
	"gopkg.in/yaml.v3"

	"github.com/zrcoder/yaq"
)

const (
	maxDisks   = 7
	totalStars = 5
)

var errCantMove = errors.New("can not move the disk above a smaller one")

var Instance = &Game{}

type Game struct {
	*yaq.Base
	*tea.Program
	currentLevel *Level
	err          error
	Name         string   `yaml:"name"`
	Levels       []string `yaml:"levels"`
	diskStyles   []lipgloss.Style
	pileWidth    int
	levelIndex   int
	loaded       bool
}

func (g *Game) Init() tea.Cmd {
	g.pileWidth = diskWidthUnit*(maxDisks) + poleWidth
	g.diskStyles = []lipgloss.Style{
		lipgloss.NewStyle().Background(color.Red),
		lipgloss.NewStyle().Background(color.Orange),
		lipgloss.NewStyle().Background(color.Yellow),
		lipgloss.NewStyle().Background(color.Green),
		lipgloss.NewStyle().Background(color.Blue),
		lipgloss.NewStyle().Background(color.Indigo),
		lipgloss.NewStyle().Background(color.Violet),
	}
	err := g.load()
	if err != nil {
		return func() tea.Msg { return err }
	}
	return textarea.Blink
}

func (g *Game) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		g.err = nil
	case common.ErrMsg:
		g.err = msg
		return nil
	}
	if g.currentLevel == nil {
		return nil
	}
	return g.currentLevel.update(msg)
}

func (g *Game) View() string {
	if g.currentLevel == nil {
		return g.LoadingView()
	}

	view := ""
	switch {
	case g.err != nil:
		view = g.ErrorView(g.err.Error())
	default:
		view = g.currentLevel.view()
	}
	title := fmt.Sprintf("%s > %s", g.Name, g.currentLevel.name)
	title += style.Help.Render(fmt.Sprintf("\tsteps: %d\n", g.currentLevel.steps))
	view = lp.JoinVertical(lp.Left, title, view)
	return view
}

func (g *Game) FS() embed.FS {
	return hanoi.FS
}

func (g *Game) Hint() string {
	if !g.loaded {
		return ""
	}
	return g.currentLevel.Hint
}

func (g *Game) MarkResult() {
	g.currentLevel.markResult()
}

func (g *Game) resetLevel() {
	g.currentLevel.initialize()
}

func (g *Game) gotoNextLevel() {
	g.levelIndex++
	g.loadCurrentLevel()
}

func (g *Game) load() error {
	err := yaml.Unmarshal(g.IndexData, g)
	if err != nil {
		return err
	}
	g.levelIndex = 0
	return g.loadCurrentLevel()
}

func (g *Game) loadCurrentLevel() error {
	if len(g.Levels) == 0 {
		return errors.New("no levels found")
	}
	if g.allLevelsFinished() {
		return errors.New("all levels finised")
	}
	data, err := g.FS().ReadFile(g.Levels[g.levelIndex] + common.YamlExt)
	if err != nil {
		return err
	}
	g.currentLevel = &Level{}
	err = yaml.Unmarshal(data, g.currentLevel)
	if err != nil {
		return err
	}
	g.currentLevel.Game = g
	g.currentLevel.name = g.Levels[g.levelIndex]
	return g.currentLevel.initialize()
}

func (g *Game) allLevelsFinished() bool {
	return g.levelIndex == len(g.Levels)
}

func (g *Game) shuffleDiskStyles() {
	rand.Shuffle(len(g.diskStyles), func(i, j int) {
		g.diskStyles[i], g.diskStyles[j] = g.diskStyles[j], g.diskStyles[i]
	})
}
