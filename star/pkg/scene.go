package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/zrcoder/yaq/common"
)

const tomlExt = ".toml"

type scene struct {
	*Game
	Sprites    map[string]*sprite `toml:"sprites"`
	name       string
	bgColors   [2]string
	BgColor1   string   `toml:"bgColor1"`
	BgColor2   string   `toml:"bgColor2"`
	LevelNames []string `toml:"levels"`
	levels     []*level `toml:"_"`
	levelIndex int
}

func (s *scene) loadLevels() error {
	if len(s.LevelNames) == 0 {
		return fmt.Errorf("no levels in scene %s", s.name)
	}

	s.levels = make([]*level, len(s.LevelNames))
	for i, name := range s.LevelNames {
		l := &level{}
		if data, err := os.ReadFile(filepath.Join(s.CfgPath, s.name, name+tomlExt)); err != nil {
			return err
		} else if err := toml.Unmarshal(data, l); err != nil {
			return err
		}
		l.name = name
		l.scene = s
		s.levels[i] = l
	}
	s.levelIndex = 0
	return s.loadCurrentLevel()
}

func (s *scene) loadCurrentLevel() error {
	if len(s.levels) == 0 {
		return fmt.Errorf("no levels found for scend %s", s.name)
	}
	if data, err := os.ReadFile(filepath.Join(s.CfgPath, s.name, s.levels[s.levelIndex].name+tomlExt)); err != nil {
		return err
	} else if err := toml.Unmarshal(data, s.currentLevel()); err != nil {
		return err
	}
	return s.currentLevel().initialize()
}

func (s *scene) currentLevel() *level {
	return s.levels[s.levelIndex]
}

func (s *scene) outRange(pos common.Position) bool {
	return pos.Y < 0 || pos.Y >= s.Rows || pos.X < 0 || pos.X >= s.Columns
}

func (s *scene) clearSpritesCount() {
	for _, sp := range s.Sprites {
		sp.count = 0
	}
}
