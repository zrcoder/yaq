package pkg

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	lp "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/yaq/common"
)

type level struct {
	*scene
	preCode    string
	Layout     string `toml:"layout"`
	Code       string `toml:"code"`
	Hint       string `toml:"hint"`
	SuccessMsg string `toml:"successMsg"`
	name       string
	grid       [][][]*sprite
	helpItems  []*sprite
}

func (l *level) initialize() error {
	l.grid = make([][][]*sprite, l.Rows)
	for i := range l.grid {
		l.grid[i] = make([][]*sprite, l.Columns)
	}
	lines := strings.Split(l.Layout, "\n")
	l.clearSpritesCount()
	for y := range l.grid {
		if y >= len(lines) {
			break
		}
		line := lines[y]
		line = line[:min(len(line), l.Columns)]
		for x, ch := range line {
			sp := l.Sprites[string(ch)]
			if sp == nil {
				continue
			}
			if len(sp.Sprites) == 0 {
				l.grid[y][x] = []*sprite{l.genSprite(sp, y, x)}
				continue
			}
			sps := make([]*sprite, len(sp.Sprites))
			for i, ch := range sp.Sprites {
				s := l.Sprites[string(ch)]
				if s == nil {
					continue
				}
				sps[i] = l.genSprite(s, y, x)
			}
			l.grid[y][x] = sps
		}
	}
	if l.player == nil {
		return errors.New("yaq is not found")
	}
	return l.calculate()
}

func (l *level) genSprite(sp *sprite, y, x int) *sprite {
	sp.count++
	sp = sp.copy()
	if sp.IsPlayer {
		l.player = sp
	}
	sp.scene = l.scene
	sp.Position = &common.Position{Y: y, X: x}
	return sp
}

func (l *level) calculate() error {
	l.totalStars = 0
	buf := strings.Builder{}
	buf.WriteString(`import . "github.com/zrcoder/yaq/star/pkg"`)
	buf.WriteString("\n\n")
	l.helpItems = nil
	helpMap := map[string]*sprite{}
	for _, sp := range l.Sprites {
		if sp.count > 1 {
			if sp.Group == "" {
				return fmt.Errorf("no group name for sprite: %s", sp.key)
			}
			helpMap[sp.key] = sp
			buf.WriteString(fmt.Sprintf("var %s = make([]*Sprite, 0, %d)\n", sp.Group, sp.count))
		}
	}
	appendBuf := strings.Builder{}
	for y, row := range l.grid {
		for x, sps := range row {
			for i, sp := range sps {
				if l.Sprites[sp.key].count == 1 {
					if sp.Name == "" {
						return fmt.Errorf("no name for sprite: %s", sp.key)
					}
					helpMap[sp.key] = sp
					buf.WriteString(fmt.Sprintf("var %s = GetSprite(%d, %d, %d)\n", sp.Name, y, x, i))
				} else if l.Sprites[sp.key].count > 1 {
					appendBuf.WriteString(fmt.Sprintf("%s = append(%s, GetSprite(%d, %d, %d))\n", sp.Group, sp.Group, y, x, i))
				}
				if strings.Contains(l.player.Foods, sp.key) {
					l.totalStars++
				}
			}
		}
	}
	l.helpItems = make([]*sprite, 0, len(helpMap))
	for _, sp := range helpMap {
		l.helpItems = append(l.helpItems, sp)
	}
	l.preCode = buf.String() + "\n" + appendBuf.String()
	l.Editor.SetValue(strings.TrimRight(l.Code, "\n"))
	sort.Slice(l.helpItems, func(i, j int) bool { return l.helpItems[i].Name < l.helpItems[j].Name })
	return nil
}

func (l *level) view() string {
	buf := strings.Builder{}
	for y, row := range l.grid {
		for x, sps := range row {
			bgColor := l.bgColors[(y+x)%2]
			if spsbg := spsBgColor(sps); spsbg != "" {
				bgColor = spsbg
			}
			s := lp.NewStyle().Background(lp.Color(bgColor))
			if len(sps) > 0 {
				n := len(sps)
				sp := sps[n-1]
				s = s.Foreground(lp.Color(sp.Color))
				buf.WriteString(s.Render(" " + sp.Display + " "))
			} else {
				buf.WriteString(s.Render("   "))
			}
		}
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	buf.WriteString(l.helpInfo())
	return buf.String()
}

func spsBgColor(sps []*sprite) string {
	res := ""
	for i := len(sps) - 1; i >= 0; i-- {
		if sps[i].BgColor != "" {
			res = sps[i].BgColor
			break
		}
	}
	return res
}

func (l *level) helpInfo() string {
	const n = 4
	buf := strings.Builder{}
	for i, sp := range l.helpItems {
		display := lp.NewStyle().
			Background(lp.Color(sp.BgColor)).
			Foreground(lp.Color(sp.Color)).
			Render(" " + sp.Display + " ")
		buf.WriteString(display)
		buf.WriteString(": ")
		if sp.count == 1 {
			buf.WriteString(sp.Name)
		} else {
			buf.WriteString(sp.Group)
		}
		if (i+1)%n == 0 {
			buf.WriteString("\n")
		} else {
			buf.WriteString("\t")
		}
	}
	return strings.TrimSuffix(buf.String(), "\n")
}
