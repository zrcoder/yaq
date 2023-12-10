package pkg

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Disk struct {
	view  string
	id    int
	width int
}

func newDisk(id int, sty lipgloss.Style) *Disk {
	view := sty.Render(strings.Repeat(diskCh, id*diskWidthUnit))
	width, _ := lipgloss.Size(view)
	return &Disk{
		id:    id,
		view:  view,
		width: width,
	}
}
