package yaq

import (
	"embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/zrcoder/yaq/common"
)

func Run(gameMode string, vimMode bool) {
	game, err := get(gameMode)
	fatalIfError(err)

	data, err := game.FS().ReadFile(common.IndexFile)
	fatalIfError(err)

	base := &Base{}
	err = yaml.Unmarshal(data, base)
	fatalIfError(err)

	base.vimMode = vimMode
	base.Init(data)
	game.SetBase(base)
	game.Run()
}

type Game interface {
	SetBase(*Base)
	Run()
	FS() embed.FS
}

var supportedGames = map[string]Game{}

func Register(mode string, g Game) {
	supportedGames[mode] = g
}

func get(mode string) (Game, error) {
	g, ok := supportedGames[mode]
	if ok {
		return g, nil
	}
	return nil, fmt.Errorf("not supported game mode: %s", mode)
}

func fatalIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
