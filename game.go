package yaq

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"

	"github.com/zrcoder/yaq/common"
)

type Game interface {
	SetBase(*Base)
	Run()
}

var supportedGames = map[string]Game{}

func Register(mode string, g Game) {
	supportedGames[mode] = g
}

func Get(mode string) (Game, error) {
	g, ok := supportedGames[mode]
	if ok {
		return g, nil
	}
	return nil, fmt.Errorf("not supported game mode: %s", mode)
}

func Run(path string) {
	base := &Base{CfgPath: path}
	data, err := os.ReadFile(filepath.Join(path, common.IndexFile))
	fatalIfError(err)
	err = toml.Unmarshal(data, base)
	fatalIfError(err)
	game, err := Get(base.Mode)
	fatalIfError(err)
	base.Init(data)
	game.SetBase(base)
	game.Run()
}

func fatalIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
