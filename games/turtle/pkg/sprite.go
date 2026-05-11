package pkg

import (
	"errors"
	"time"

	"github.com/zrcoder/yaq/pkg"
)

type Sprite struct {
	*Pen    `yaml:"pen"`
	*Block  `yaml:"block"`
	IsPen   bool   `yaml:"isPen"`
	Sprites string `yaml:"sprites"`
}

type Block struct {
	BgColor string `yaml:"bgColor"`
	filled  bool
}

func newBlock(color string) *Block {
	return &Block{BgColor: color}
}

type Pen struct {
	*Game
	*pkg.Position
	Color   string `yaml:"color"`
	Display string `yaml:"display"`
	IsUp    bool   `yaml:"isUp"`
}

func (p *Pen) setStateUp(up bool) {
	if !up {
		if p.currentLevel.Grid[p.Y][p.X] == nil {
			p.Send(errors.New("should paint on marked position"))
			return
		}
		p.paint(p.Position)
	}
	p.IsUp = up
}

func (p *Pen) up(steps int) {
	p.move(steps, pkg.Up)
}

func (p *Pen) down(steps int) {
	p.move(steps, pkg.Down)
}

func (p *Pen) left(steps int) {
	p.move(steps, pkg.Left)
}

func (p *Pen) right(steps int) {
	p.move(steps, pkg.Right)
}

func (p *Pen) upLeft(steps int) {
	p.move(steps, pkg.UpLeft)
}

func (p *Pen) upRight(steps int) {
	p.move(steps, pkg.UpRight)
}

func (p *Pen) downLeft(steps int) {
	p.move(steps, pkg.DownLeft)
}

func (p *Pen) downRight(steps int) {
	p.move(steps, pkg.DownRight)
}

func (p *Pen) move(steps int, dir pkg.Direction) {
	if p.IsUp {
		dst := p.Transform(pkg.Direction{Y: steps * dir.Y, X: steps * dir.X})
		if p.outRange(&dst) {
			p.Send(errors.New("can't move out of the world"))
			return
		}
		p.Position = &dst
		p.Send(pkg.MoveMsg{})
		return
	}
	for ; steps > 0; steps-- {
		if err := p.step(dir); err != nil {
			p.Send(err)
			return
		}
		p.Send(pkg.MoveMsg{})
		time.Sleep(300 * time.Millisecond)
	}
}

func (p *Pen) step(dir pkg.Direction) error {
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

func (p *Pen) paint(dst *pkg.Position) {
	y, x := dst.Y, dst.X
	if !p.currentLevel.Grid[y][x].filled {
		p.currentLevel.Grid[y][x].BgColor = p.Color
		p.currentLevel.Grid[y][x].filled = true
		p.totalPoses--
	}
	p.Position = dst
}
