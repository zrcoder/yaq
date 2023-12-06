package pkg

import (
	"errors"

	"github.com/zrcoder/yaq/common"
)

type Pen struct {
	*Game
	*common.Position
	Color   string `toml:"color"`
	Display string `toml:"display"`
	IsUp    bool   `toml:"isUp"`
}

func (p *Pen) SetStateUp(up bool) {
	if !up {
		if p.currentLevel.Grid[p.Y][p.X] == nil {
			p.Send(errors.New("should paint on marked position"))
			return
		}
		p.paint(p.Position)
	}
	p.IsUp = up
}

func (p *Pen) Up(steps int) {
	p.move(steps, common.Up)
}

func (p *Pen) Down(steps int) {
	p.move(steps, common.Down)
}

func (p *Pen) Left(steps int) {
	p.move(steps, common.Left)
}

func (p *Pen) Right(steps int) {
	p.move(steps, common.Right)
}

func (p *Pen) UpLeft(steps int) {
	p.move(steps, common.UpLeft)
}

func (p *Pen) UpRight(steps int) {
	p.move(steps, common.UpRight)
}

func (p *Pen) DownLeft(steps int) {
	p.move(steps, common.DownLeft)
}

func (p *Pen) DownRight(steps int) {
	p.move(steps, common.DownRight)
}

func (p *Pen) move(steps int, dir common.Direction) {
	if p.IsUp {
		dst := p.Transform(common.Direction{Y: steps * dir.Y, X: steps * dir.X})
		if p.outRange(&dst) {
			p.Send(errMsg(errors.New("can't move out of the world")))
			return
		}
		p.Position = &dst
		return
	}
	for ; steps > 0; steps-- {
		if err := p.step(dir); err != nil {
			p.Send(errMsg(err))
			return
		}
	}
}

func (p *Pen) step(dir common.Direction) error {
	dst := p.Transform(dir)
	if p.outRange(&dst) {
		return errors.New("can't draw out of the world")
	}
	y, x := dst.Y, dst.X
	if p.currentLevel.Grid[y][x] == nil {
		return errors.New("should paint on marked position")
	}
	p.paint(&dst)
	return nil
}

func (p *Pen) paint(dst *common.Position) {
	y, x := dst.Y, dst.X
	p.currentLevel.Grid[y][x].Color = p.Color
	p.totalPoses--
	p.Position = dst
}
