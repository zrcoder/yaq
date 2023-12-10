package pkg

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lp "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/yaq/common"
)

type Level struct {
	err      error
	overDisk *Disk
	buf      *strings.Builder
	*Game
	Code        string `tomal:"code"`
	name        string
	Hint        string `tomal:"hint"`
	successInfo string
	state       common.State
	piles       []*Pile
	steps       int
	earnedStars int
	Disks       int `tomal:"disks"`
}

func (l *Level) initialize() error {
	l.steps = 0
	l.overDisk = nil
	l.buf = &strings.Builder{}
	l.piles = make([]*Pile, 3)
	pileNames := []string{"A", "B", "C"}
	for i := range l.piles {
		l.piles[i] = &Pile{Game: l.Game, name: pileNames[i]}
	}
	l.shuffleDiskStyles()
	disks := make([]*Disk, l.Disks)
	for i := 1; i <= l.Disks; i++ {
		disks[l.Disks-i] = newDisk(i, l.diskStyles[i-1])
	}
	l.piles[0].disks = disks
	l.Editor.SetValue(strings.TrimRight(l.Code, "\n"))
	return nil
}

func (l *Level) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case common.ErrMsg:
		l.err = msg
		l.state = common.Failed
		return nil
	case tea.KeyMsg:
		l.err = nil
		switch l.state {
		case common.Succeed:
			l.state = common.Running
			l.gotoNextLevel()
		case common.Failed:
			l.state = common.Running
			l.resetLevel()
		default:
		}
	}
	var cmd tea.Cmd
	l.Editor, cmd = l.Editor.Update(msg)
	return cmd
}

func (l *Level) view() string {
	title := fmt.Sprintf("%s > %s", l.Name, l.name)
	title += style.Help.Render(fmt.Sprintf("\tsteps: %d", l.steps))
	l.buf.Reset()
	l.writePoles()

	mainView := ""
	switch l.state {
	case common.Running:
		views := make([]string, len(l.piles))
		for i, p := range l.piles {
			views[i] = p.view()
		}
		mainView = lp.JoinVertical(lp.Center,
			views[2],
			lp.JoinHorizontal(lp.Top, views[0], views[1]),
		)
	case common.Succeed:
		mainView = l.SucceedViewWithStars(l.successInfo, totalStars, l.earnedStars)
	case common.Failed:
		mainView = l.ErrorView("failed")
	}

	leftView := lp.JoinVertical(lp.Left,
		title,
		mainView,
	)
	rightView := lp.JoinVertical(lp.Left,
		style.Help.Render(l.Hint), "",
		l.Editor.View(), "",
		l.KeysView(),
	)
	return lp.JoinHorizontal(lp.Top, leftView, "     ", rightView)
}

func (l *Level) writePoles() {
	views := make([]string, len(l.piles))
	for i, p := range l.piles {
		views[i] = p.view()
	}
	poles := lp.JoinHorizontal(
		lp.Top,
		views...,
	)
	l.buf.WriteString(poles)
	l.writeBlankLine()
}

func (l *Level) writeLine(s string) {
	l.buf.WriteString(s)
	l.writeBlankLine()
}

func (l *Level) writeBlankLine() {
	l.buf.WriteByte('\n')
}

func (l *Level) markResult() {
	if l.success() {
		l.state = common.Succeed
		l.setSuccessView()
	} else {
		l.state = common.Failed
	}
	l.Send(common.ResMsg{})
}

func (l *Level) setSuccessView() {
	minSteps := 1<<l.Disks - 1
	totalStars := 5
	if l.steps == minSteps {
		l.successInfo = "Fantastic! you earned all the stars!"
		l.earnedStars = totalStars
		return
	}
	l.successInfo = fmt.Sprintf("Done! Taken %d steps, can you complete it in %d step(s)? ", l.steps, minSteps)
	l.earnedStars = 3
	if l.steps-minSteps > minSteps/2 {
		l.earnedStars = 1
	}
}

func (l *Level) pick(i int) {
	curPile := l.piles[i]
	if l.overDisk == nil && curPile.empty() {
		return
	}
	if l.overDisk == nil {
		curPile.overOne = true
		l.overDisk = curPile.top()
		return
	}
	if !curPile.empty() && l.overDisk.id > curPile.top().id {
		l.Send(errCantMove)
		return
	}
	if !curPile.empty() && l.overDisk == curPile.top() {
		curPile.overOne = false
		l.overDisk = nil
		return
	}
	for _, p := range l.piles {
		if p.overOne {
			l.steps++
			curPile.push(p.pop())
			p.overOne = false
			l.overDisk = nil
		}
	}
	time.Sleep(500 * time.Millisecond)
}

func (l *Level) success() bool {
	last := l.piles[len(l.piles)-1]
	return len(last.disks) == l.Disks
}
