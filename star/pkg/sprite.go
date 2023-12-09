package pkg

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/zrcoder/yaq/common"
)

type sprite struct {
	*scene
	*common.Position
	key       string
	Name      string `toml:"name"`
	Group     string `toml:"group"`
	Display   string `toml:"display"`
	Color     string `toml:"color"`
	BgColor   string `toml:"bgColor"`
	Foods     string `toml:"foods"`
	Forbbiden string `toml:"forbbiden"`
	Freinds   string `toml:"freinds"`
	Sprites   string `toml:"sprites"`
	count     int
	IsPlayer  bool `toml:"isPlayer"`
	CanMove   bool `toml:"canMove"`
}

func (s *sprite) Up(steps int) {
	s.move(common.Up, steps)
}

func (s *sprite) Left(steps int) {
	s.move(common.Left, steps)
}

func (s *sprite) Down(steps int) {
	s.move(common.Down, steps)
}

func (s *sprite) Right(steps int) {
	s.move(common.Right, steps)
}

func (s *sprite) UpLeft(steps int) {
	s.move(common.UpLeft, steps)
}

func (s *sprite) UpRight(steps int) {
	s.move(common.UpRight, steps)
}

func (s *sprite) DownLeft(steps int) {
	s.move(common.DownLeft, steps)
}

func (s *sprite) DownRight(steps int) {
	s.move(common.DownRight, steps)
}

func (s *sprite) move(dir common.Direction, steps int) {
	if !s.CanMove {
		name := s.Name
		if s.count > 0 {
			name = s.Group
		}
		err := fmt.Errorf("%s can't move", name)
		s.Send(err)
		return
	}

	for ; steps > 0; steps-- {
		err := s.step(dir)
		if err != nil {
			s.Send(err)
			return
		}
		s.Send(common.MoveMsg{})
		time.Sleep(300 * time.Millisecond)
	}
}

func (s *sprite) step(dir common.Direction) error {
	dstPos := s.Transform(dir)
	if s.outRange(dstPos) {
		return errors.New("can't move out of the world")
	}
	grid := s.currentLevel().grid
	y, x := s.Y, s.X
	srcSps := grid[y][x]
	idx := slices.Index(srcSps, s)
	toMove := grid[y][x][idx:]
	dstSps := grid[dstPos.Y][dstPos.X]
	canCross, name := s.crossCheck(dstSps)
	if !canCross {
		return fmt.Errorf("can't cross %s", name)
	}
	playerMoving := false
	for _, sp := range srcSps {
		if sp == s.player {
			playerMoving = true
			break
		}
	}
	if playerMoving && len(dstSps) > 0 {
		n := len(dstSps)
		top := dstSps[n-1]
		if strings.Contains(s.player.Foods, top.key) {
			s.totalStars--
			grid[dstPos.Y][dstPos.X] = dstSps[:n-1]
		}
	}
	grid[y][x] = grid[y][x][:idx]
	grid[dstPos.Y][dstPos.X] = append(grid[dstPos.Y][dstPos.X], toMove...)
	for _, sp := range toMove {
		sp.Y = dstPos.Y
		sp.X = dstPos.X
	}
	return nil
}

func (s *sprite) crossCheck(sps []*sprite) (bool, string) {
	if len(sps) == 0 {
		return !strings.Contains(s.Forbbiden, " "), "blank"
	}

	for i := len(sps) - 1; i >= 0; i-- {
		d := sps[i]
		if strings.Contains(s.Freinds, d.key) {
			return true, ""
		}
		if strings.Contains(s.Forbbiden, d.key) {
			return false, d.Name
		}
	}

	return true, ""
}

func (s *sprite) copy() *sprite {
	dst := *s
	return &dst
}
