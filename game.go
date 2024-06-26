package yaq

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/zrcoder/yaq/common"
)

func Run(path string) {
	base := &Base{CfgPath: path}
	data, err := os.ReadFile(filepath.Join(path, common.IndexFile))
	fatalIfError(err)
	err = yaml.Unmarshal(data, base)
	fatalIfError(err)
	game, err := get(base.Mode)
	fatalIfError(err)
	base.Init(data)
	game.SetBase(base)
	game.Run()
}

type Game interface {
	SetBase(*Base)
	Run()
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
