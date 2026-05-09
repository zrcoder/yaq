package yaq

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	lp "charm.land/lipgloss/v2"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/vtea"
)

type Base struct {
	Editor    *vtea.Model
	CfgPath   string
	Name      string `yaml:"name"`
	Mode      string `yaml:"mode"`
	IndexData []byte
	Rows      int `yaml:"rows"`
	Columns   int `yaml:"columns"`
	height    int
	width     int
	Keys      KeyMap
	KeysHelp  help.Model
}

func (b *Base) Init(data []byte) {
	b.IndexData = data
	b.Keys = getCommonKeys()
	b.KeysHelp = help.New()
	ta := vtea.New(vtea.WithFileName("x.sh"))
	b.Editor = ta
}

func (b *Base) SetSceneSize(height, width int) {
	b.height = height
	b.width = width
	b.Editor.SetSize(width, height)
}

func (b *Base) ErrorView(msg string) string {
	return dialog.Error(msg).Height(b.height).Width(b.width).String()
}

func (b *Base) SucceedView(msg string) string {
	return dialog.Success(msg).Height(b.height).Width(b.width).String()
}

func (b *Base) SucceedViewWithStars(msg string, total, stars int) string {
	return dialog.Success(msg).Height(b.height).Width(b.width).Stars(total, stars).String()
}

func (b *Base) LoadingView() string {
	return dialog.Success("loading...").Height(b.height).Width(b.width).String()
}

func (b *Base) KeysView() string {
	return b.KeysHelp.View(b.Keys)
}

func (b *Base) View(leftView, rightView string) tea.View {
	view := tea.NewView(lp.JoinHorizontal(lp.Top, leftView, "  ", rightView))
	view.AltScreen = true
	return view
}
