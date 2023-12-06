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
	Keys      KeyMap
	KeysHelp  help.Model
}

func (b *Base) Init(data []byte) {
	b.IndexData = data
	ta := textarea.New()
	ta.SetHeight(b.Rows)
	ta.SetWidth(b.Columns * 3)
	ta.Focus()
	b.Keys = getCommonKeys()
	b.KeysHelp = help.New()
	b.Editor = ta
}

func (b *Base) ErrorView(msg string) string {
	return dialog.Error(msg).Height(b.Rows).Width(b.Columns * 3).String()
}

func (b *Base) SucceedView(msg string) string {
	return dialog.Success(msg).Height(b.Rows).Width(b.Columns * 3).String()
}

func (b *Base) KeysView() string {
	return b.KeysHelp.View(b.Keys)
}
