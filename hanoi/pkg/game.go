package pkg

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pelletier/go-toml/v2"
	"github.com/zrcoder/rdor/pkg/style/color"
	"github.com/zrcoder/yaq/common"

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
	Name         string   `toml:"name"`
	Levels       []string `toml:"levels"`
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

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.ErrMsg:
		g.err = msg
		return g, nil
	}
	if g.currentLevel == nil {
		return g, nil
	}
	return g, g.currentLevel.update(msg)
}

func (g *Game) View() string {
	if g.err != nil {
		return g.ErrorView(g.err.Error())
	}
	if g.currentLevel == nil {
		return g.LoadingView()
	}
	return g.currentLevel.view()
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
	err := toml.Unmarshal(g.IndexData, g)
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
	data, err := os.ReadFile(filepath.Join(g.CfgPath, g.Levels[g.levelIndex]+common.TomlExt))
	if err != nil {
		return err
	}
	g.currentLevel = &Level{}
	err = toml.Unmarshal(data, g.currentLevel)
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

func (g *Game) helpInfo() string {
	return "Our goal is to move all disks from pile `1` to pile `3`."
}

func (g *Game) shuffleDiskStyles() {
	rand.Shuffle(len(g.diskStyles), func(i, j int) {
		g.diskStyles[i], g.diskStyles[j] = g.diskStyles[j], g.diskStyles[i]
	})
}
