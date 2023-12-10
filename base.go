package yaq

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/zrcoder/rdor/pkg/dialog"
)

type Base struct {
	Editor    textarea.Model
	CfgPath   string
	Name      string `toml:"name"`
	Mode      string `toml:"mode"`
	IndexData []byte
	Rows      int `toml:"rows"`
	Columns   int `toml:"columns"`
	height    int
	width     int
	Keys      KeyMap
	KeysHelp  help.Model
}

func (b *Base) Init(data []byte) {
	b.IndexData = data
	b.Keys = getCommonKeys()
	b.KeysHelp = help.New()
	ta := textarea.New()
	ta.Focus()
	b.Editor = ta
}

func (b *Base) SetSceneSize(height, width int) {
	b.height = height
	b.width = width
	b.Editor.SetHeight(height)
	b.Editor.SetWidth(width)
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
