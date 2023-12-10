package pkg

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	poleWidth     = 1
	diskWidthUnit = 4

	poleCh   = "|"
	diskCh   = " "
	groundCh = "‾"
)

type Pile struct {
	*Game
	name    string
	disks   []*Disk
	overOne bool
}

func (p *Pile) empty() bool {
	return len(p.disks) == 0
}

func (p *Pile) push(d *Disk) {
	p.disks = append(p.disks, d)
}

func (p *Pile) pop() *Disk {
	n := len(p.disks)
	res := p.disks[n-1]
	p.disks = p.disks[:n-1]
	return res
}

func (p *Pile) top() *Disk {
	n := len(p.disks)
	return p.disks[n-1]
}

func (p *Pile) view() string {
	lines := make([]string, maxDisks+4)
	lines[0] = strings.Repeat(" ", maxDisks*diskWidthUnit)
	disks := p.disks
	writeDisk := func(i int) {
		lines[i] = disks[len(disks)-1].view
		disks = disks[:len(disks)-1]
	}
	if p.overOne {
		writeDisk(1)
	}
	for i := maxDisks; i > 0; i-- {
		j := maxDisks - i + 2
		if i == len(disks) {
			writeDisk(j)
		} else {
			lines[j] = poleCh
		}
	}
	lines[len(lines)-1] = p.name
	lines[len(lines)-2] = strings.Repeat(groundCh, maxDisks*diskWidthUnit)
	return lipgloss.NewStyle().Width(maxDisks * diskWidthUnit).Render(
		lipgloss.JoinVertical(lipgloss.Center, lines...),
	)
}
