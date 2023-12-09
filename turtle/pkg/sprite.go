package pkg

import (
	"errors"
	"time"

	"github.com/zrcoder/yaq/common"
)

type Sprite struct {
	*Pen    `toml:"pen"`
	*Block  `toml:"block"`
	IsPen   bool   `toml:"isPen"`
	Sprites string `toml:"sprites"`
}

type Block struct {
	BgColor string `toml:"bgColor"`
	filled  bool
}

func newBlock(color string) *Block {
	return &Block{BgColor: color}
}

type Pen struct {
	*Game
	*common.Position
	Color   string `toml:"color"`
	Display string `toml:"display"`
	IsUp    bool   `toml:"isUp"`
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
	p.move(steps, common.Up)
}

func (p *Pen) down(steps int) {
	p.move(steps, common.Down)
}

func (p *Pen) left(steps int) {
	p.move(steps, common.Left)
}

func (p *Pen) right(steps int) {
	p.move(steps, common.Right)
}

func (p *Pen) upLeft(steps int) {
	p.move(steps, common.UpLeft)
}

func (p *Pen) upRight(steps int) {
	p.move(steps, common.UpRight)
}

func (p *Pen) downLeft(steps int) {
	p.move(steps, common.DownLeft)
}

func (p *Pen) downRight(steps int) {
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
		time.Sleep(300 * time.Millisecond)
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
	if !p.currentLevel.Grid[y][x].filled {
		p.currentLevel.Grid[y][x].BgColor = p.Color
		p.currentLevel.Grid[y][x].filled = true
		p.totalPoses--
	}
	p.Position = dst
}
