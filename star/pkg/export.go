//go:generate igop export -outdir ../../exported .

package pkg

func Up(steps int) {
	Instance.player.Up(steps)
}

func Left(steps int) {
	Instance.player.Left(steps)
}

func Down(steps int) {
	Instance.player.Down(steps)
}

func Right(steps int) {
	Instance.player.Right(steps)
}

func UpLeft(steps int) {
	Instance.player.UpLeft(steps)
}

func UpRight(steps int) {
	Instance.player.UpRight(steps)
}

func DownLeft(steps int) {
	Instance.player.DownLeft(steps)
}

func DownRight(steps int) {
	Instance.player.DownRight(steps)
}

func GetSprite(y, x, i int) *Sprite {
	return Instance.currentLevel().grid[y][x][i]
}
