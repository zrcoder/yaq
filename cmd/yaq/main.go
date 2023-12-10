package main

import (
	"os"

	"github.com/zrcoder/yaq"
	// import to register games
	_ "github.com/zrcoder/yaq/hanoi"
	_ "github.com/zrcoder/yaq/star"
	_ "github.com/zrcoder/yaq/turtle"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	yaq.Run(path)
}
