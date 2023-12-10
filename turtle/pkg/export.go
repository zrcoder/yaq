//go:generate igop export -outdir ../../exported .

package pkg

func PenUp() {
	Instance.pen.setStateUp(true)
}

func PenDown() {
	Instance.pen.setStateUp(false)
}

func Up(steps int) {
	Instance.pen.up(steps)
}

func Left(steps int) {
	Instance.pen.left(steps)
}

func Down(steps int) {
	Instance.pen.down(steps)
}

func Right(steps int) {
	Instance.pen.right(steps)
}

func UpLeft(steps int) {
	Instance.pen.upLeft(steps)
}

func UpRight(steps int) {
	Instance.pen.upRight(steps)
}

func DownLeft(steps int) {
	Instance.pen.downLeft(steps)
}

func DownRight(steps int) {
	Instance.pen.downRight(steps)
}
