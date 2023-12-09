package common

import (
	"fmt"
	"strconv"
	"strings"
)

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

type (
	ErrMsg  = error
	MoveMsg struct{}
	ResMsg  struct{}
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

func ParseBuildError(err error, preCode string) error {
	msg := err.Error()
	i := strings.Index(msg, " ")
	if i == -1 {
		return err
	}
	pre, post := msg[:i], msg[i+1:]
	arr := strings.SplitN(pre, ":", 3)
	if len(arr) < 2 {
		return err
	}
	line, _ := strconv.Atoi(arr[1])
	line -= strings.Count(preCode, "\n")
	return fmt.Errorf("line %d: %s", line, post)
}
