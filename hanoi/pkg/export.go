//go:generate go tool qexp  -outdir ../../exported .

package pkg

func A() {
	Instance.currentLevel.pick(0)
}

func B() {
	Instance.currentLevel.pick(1)
}

func C() {
	Instance.currentLevel.pick(2)
}
