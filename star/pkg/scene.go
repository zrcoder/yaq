package pkg

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zrcoder/yaq/common"
	"gopkg.in/yaml.v3"
)

type Scene struct {
	*Game
	Sprites    map[string]*Sprite `yaml:"sprites"`
	name       string
	bgColors   [2]string
	BgColor1   string   `yaml:"bgColor1"`
	BgColor2   string   `yaml:"bgColor2"`
	LevelNames []string `yaml:"levels"`
	levels     []*Level `yaml:"_"`
	levelIndex int
}

func (s *Scene) loadLevels() error {
	if len(s.LevelNames) == 0 {
		return fmt.Errorf("no levels in scene %s", s.name)
	}

	s.levels = make([]*Level, len(s.LevelNames))
	for i, name := range s.LevelNames {
		l := &Level{}
		if data, err := os.ReadFile(filepath.Join(s.CfgPath, s.name, name+common.YamlExt)); err != nil {
			return err
		} else if err := yaml.Unmarshal(data, l); err != nil {
			return err
		}
		l.name = name
		l.Scene = s
		s.levels[i] = l
	}
	s.levelIndex = 0
	return s.loadCurrentLevel()
}

func (s *Scene) loadCurrentLevel() error {
	if len(s.levels) == 0 {
		return fmt.Errorf("no levels found for scend %s", s.name)
	}
	if data, err := os.ReadFile(filepath.Join(s.CfgPath, s.name, s.levels[s.levelIndex].name+common.YamlExt)); err != nil {
		return err
	} else if err := yaml.Unmarshal(data, s.currentLevel()); err != nil {
		return err
	}
	return s.currentLevel().initialize()
}

func (s *Scene) currentLevel() *Level {
	return s.levels[s.levelIndex]
}

func (s *Scene) outRange(pos common.Position) bool {
	return pos.Y < 0 || pos.Y >= s.Rows || pos.X < 0 || pos.X >= s.Columns
}

func (s *Scene) clearSpritesCount() {
	for _, sp := range s.Sprites {
		sp.count = 0
	}
}
