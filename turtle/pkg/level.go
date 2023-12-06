package pkg

import (
	"errors"
	"strings"

	lp "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/yaq/common"
)

type Level struct {
	*Game
	preCode    string
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
			l.pen = sp.pen
			l.pen.Position = &common.Position{Y: y, X: x}
		} else {
			l.Grid[y][x] = newBlock(sp.Block.Color)
			l.totalPoses++
		}
	}
	for y, line := range lines {
		line = line[:min(l.Columns, len(line))]
		for x, key := range line {
			sp := l.Sprites[string(key)]
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
	l.pen.SetStateUp(l.pen.IsUp)
	return nil
}

func (l *Level) View() string {
	buf := strings.Builder{}
	sep := strings.Repeat("• ", l.Columns) + "•\n"
	buf.WriteString(sep)
	for y, row := range l.Grid {
		for x, sp := range row {
			bgColor := ""
			if sp != nil {
				bgColor = sp.Color
			}
			s := lp.NewStyle().Background(lp.Color(bgColor))
			display := "• "
			if l.pen.Y == y && l.pen.X == x {
				display = "•" + l.pen.Display
			}
			buf.WriteString(s.Render(display))
		}
		buf.WriteString("•\n")
		buf.WriteString(sep)
	}
	return buf.String()
}
