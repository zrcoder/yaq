package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/zrcoder/yaq/pkg"

	// import to register games
	_ "github.com/zrcoder/yaq/games/hanoi"
	_ "github.com/zrcoder/yaq/games/star"
	_ "github.com/zrcoder/yaq/games/turtle"
)

func main() {
	game := "star"
	vimMode := false

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Yaq games").
				Options(
					huh.NewOption("Code Star", "star"),
					huh.NewOption("Hanoi", "hanoi"),
					huh.NewOption("Turtle Graphics", "turtle"),
				).
				Value(&game),
			huh.NewConfirm().Title("Vim mode").Value(&vimMode),
		),
	).WithViewHook(func(v tea.View) tea.View {
		v.AltScreen = true
		return v
	})

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	pkg.Run(game, vimMode)
}
