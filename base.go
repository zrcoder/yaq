package yaq

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	lp "charm.land/lipgloss/v2"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/vtea"
)

type Base struct {
	editor        textarea.Model
	vimEditor     *vtea.Model
	vimMode       bool
	CfgPath       string
	Name          string `yaml:"name"`
	Mode          string `yaml:"mode"`
	IndexData     []byte
	Rows          int `yaml:"rows"`
	Columns       int `yaml:"columns"`
	height        int
	width         int
	Keys          KeyMap
	KeysHelp      help.Model
	RunCodeAction func()
}

func (b *Base) Init(data []byte) {
	b.IndexData = data
	b.Keys = getCommonKeys()
	b.KeysHelp = help.New()
	b.editor = textarea.New()
	b.editor.Focus()
	b.vimEditor = vtea.New(vtea.WithFileName("x.sh"))
}

func (b *Base) SetSceneSize(height, width int) {
	b.height = height
	b.width = width
	b.editor.SetWidth(width)
	b.editor.SetHeight(height)
	b.vimEditor.SetSize(width, height)
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

func (b *Base) EditorValue() string {
	if b.vimMode {
		return b.vimEditor.Value()
	}
	return b.editor.Value()
}

func (b *Base) EditorView() string {
	if b.vimMode {
		return b.vimEditor.View()
	}
	return b.editor.View()
}

func (b *Base) SetEditorValue(s string) {
	if b.vimMode {
		b.vimEditor.SetValue(s)
	} else {
		b.editor.SetValue(s)
	}
}

func (b *Base) SetEditorSize(width, height int) {
	if b.vimMode {
		b.vimEditor.SetSize(width, height)
	} else {
		b.editor.SetWidth(width)
		b.editor.SetHeight(height)
	}
}

func (b *Base) EditorUpdate(msg tea.Msg) (cmd tea.Cmd) {
	if b.vimMode {
		b.vimEditor, cmd = b.vimEditor.Update(msg)
		return cmd
	}
	b.editor, cmd = b.editor.Update(msg)
	return
}

func (b *Base) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit
		case "ctrl+r":
			if b.RunCodeAction != nil {
				b.RunCodeAction()
			}
		case "ctrl+e":
			b.switchEditor()
		}
	}
	return nil
}

func (b *Base) switchEditor() {
	if b.vimMode {
		b.vimMode = false
		b.editor.SetValue(b.vimEditor.Value())
		b.editor.Focus()
	} else {
		b.vimMode = true
		b.vimEditor.SetValue(b.editor.Value())
	}
}
