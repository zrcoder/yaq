package pkg

import (
	"errors"
	"strings"

	lp "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/yaq/common"
)

type Level struct {
	*Game
	Layout     string `toml:"layout"`
	Code       string `toml:"code"`
	Hint       string `toml:"hint"`
	SuccessMsg string `toml:"successMsg"`
	Name       string
	Grid       [][]*Block
}

func (l *Level) initialize() error {
	l.Grid = make([][]*Block, l.Rows)
	for i := range l.Grid {
		l.Grid[i] = make([]*Block, l.Columns)
	}
	l.totalPoses = 0
	lines := strings.SplitN(l.Layout, "\n", l.Rows)
	setCommonSprite := func(sp *Sprite, y, x int) {
		if sp.IsPen {
			l.pen = sp.Pen
			l.pen.Position = &common.Position{Y: y, X: x}
		} else {
			l.Grid[y][x] = newBlock(sp.Block.BgColor)
			l.totalPoses++
		}
	}
	for y, line := range lines {
		line = line[:min(l.Columns, len(line))]
		for x, ch := range line {
			key := string(ch)
			sp := l.Sprites[key]
			if sp == nil {
				continue
			}
			if sp.Sprites == "" {
				setCommonSprite(sp, y, x)
				continue
			}
			for _, k := range sp.Sprites {
				s := l.Sprites[string(k)]
				if s != nil {
					setCommonSprite(s, y, x)
				}
			}
		}
	}
	if l.pen == nil {
		return errors.New("no pen found")
	}
	l.pen.Game = l.Game
	l.pen.setStateUp(l.pen.IsUp)
	l.Editor.SetValue(strings.TrimRight(l.Code, "\n"))
	return nil
}

func (l *Level) View() string {
	buf := strings.Builder{}
	for y, row := range l.Grid {
		for x, sp := range row {
			display := " • " // blank
			bgColor := ""
			if sp != nil { // block
				bgColor = sp.BgColor
				display = "   "
			}
			if l.pen.Y == y && l.pen.X == x { // pen
				if !l.pen.IsUp {
					bgColor = l.pen.Color
				}
				display = " " + l.pen.Display + " "
			}
			s := lp.NewStyle().Background(lp.Color(bgColor))
			buf.WriteString(s.Render(display))
		}
		buf.WriteString(" \n")
	}
	return buf.String()
}
