package pkg

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lp "charm.land/lipgloss/v2"
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
	l.SetEditorValue(strings.TrimRight(l.Code, "\n"))
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
	return l.EditorUpdate(msg)
}

func (l *Level) view() string {
	l.buf.Reset()
	l.writePoles()

	switch l.state {
	case common.Running:
		views := make([]string, len(l.piles))
		for i, p := range l.piles {
			views[i] = p.view()
		}
		return lp.JoinVertical(lp.Center,
			views[2],
			lp.JoinHorizontal(lp.Top, views[0], views[1]),
		)
	case common.Succeed:
		return l.SucceedViewWithStars(l.successInfo, totalStars, l.earnedStars)
	case common.Failed:
		return l.ErrorView("failed")
	}
	return ""
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

	sleepDuration := time.Duration(0)
	defer func() {
		if sleepDuration > 0 {
			l.Send(common.MoveMsg{})
			time.Sleep(sleepDuration)
		}
	}()

	if l.overDisk == nil {
		curPile.overOne = true
		l.overDisk = curPile.top()
		sleepDuration = 200 * time.Millisecond
		return
	}
	if !curPile.empty() && l.overDisk.id > curPile.top().id {
		l.Send(errCantMove)
		return
	}
	if !curPile.empty() && l.overDisk == curPile.top() {
		curPile.overOne = false
		l.overDisk = nil
		sleepDuration = 200 * time.Millisecond
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
	sleepDuration = 500 * time.Millisecond
}

func (l *Level) success() bool {
	last := l.piles[len(l.piles)-1]
	return len(last.disks) == l.Disks
}
