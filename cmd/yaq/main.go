package main

import (
	"log"

	"charm.land/huh/v2"
	"github.com/zrcoder/yaq"

	// import to register games
	_ "github.com/zrcoder/yaq/hanoi"
	_ "github.com/zrcoder/yaq/star"
	_ "github.com/zrcoder/yaq/turtle"
)

func main() {
	game := "hanoi"
	vimMode := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Yaq games").
				Options(
					huh.NewOption("Turtle Graphics", "turtle"),
					huh.NewOption("Code Star", "star"),
					huh.NewOption("Hanoi", "hanoi"),
				).
				Value(&game),
			huh.NewConfirm().Title("Vim mode").Value(&vimMode),
		),
	)
	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	yaq.Run(game, vimMode)
}
