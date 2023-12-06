package pkg

type Block struct {
	Color string
}

type Sprite struct {
	*Block
	*Pen
	Sprites string `toml:"sprites"`
	IsPen   bool   `toml:"isPen"`
}

func newBlock(color string) *Block {
	return &Block{Color: color}
}
