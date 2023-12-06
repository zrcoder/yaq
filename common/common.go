package common

const (
	IndexFile = "index.toml"
	TomlExt   = ".toml"
)

type State = string

var (
	Running State = ""
	Succeed State = "succeed"
	Failed  State = "failed"
)

type (
	Position  struct{ Y, X int }
	Direction = Position
)

var (
	Up        = Direction{Y: -1, X: 0}
	Left      = Direction{Y: 0, X: -1}
	Down      = Direction{Y: 1, X: 0}
	Right     = Direction{Y: 0, X: 1}
	UpLeft    = Direction{Y: -1, X: -1}
	UpRight   = Direction{Y: -1, X: 1}
	DownLeft  = Direction{Y: 1, X: -1}
	DownRight = Direction{Y: 1, X: 1}
)

func (p *Position) Transform(d Direction) Position {
	return Position{Y: p.Y + d.Y, X: p.X + d.X}
}
