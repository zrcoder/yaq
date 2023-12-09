// export by github.com/goplus/igop/cmd/qexp

package pkg

import (
	q "github.com/zrcoder/yaq/turtle/pkg"

	"reflect"

	"github.com/goplus/igop"
)

func init() {
	igop.RegisterPackage(&igop.Package{
		Name: "pkg",
		Path: "github.com/zrcoder/yaq/turtle/pkg",
		Deps: map[string]string{
			"errors": "errors",
			"fmt":    "fmt",
			"github.com/charmbracelet/bubbles/textarea": "textarea",
			"github.com/charmbracelet/bubbletea":        "tea",
			"github.com/charmbracelet/lipgloss":         "lipgloss",
			"github.com/pelletier/go-toml/v2":           "toml",
			"github.com/zrcoder/rdor/pkg/dialog":        "dialog",
			"github.com/zrcoder/rdor/pkg/style":         "style",
			"github.com/zrcoder/yaq":                    "yaq",
			"github.com/zrcoder/yaq/common":             "common",
			"os":                                        "os",
			"path/filepath":                             "filepath",
			"strings":                                   "strings",
			"time":                                      "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Block":  reflect.TypeOf((*q.Block)(nil)).Elem(),
			"Game":   reflect.TypeOf((*q.Game)(nil)).Elem(),
			"Pen":    reflect.TypeOf((*q.Pen)(nil)).Elem(),
			"Sprite": reflect.TypeOf((*q.Sprite)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"Instance": reflect.ValueOf(&q.Instance),
		},
		Funcs: map[string]reflect.Value{
			"Down":      reflect.ValueOf(q.Down),
			"DownLeft":  reflect.ValueOf(q.DownLeft),
			"DownRight": reflect.ValueOf(q.DownRight),
			"Left":      reflect.ValueOf(q.Left),
			"PenDown":   reflect.ValueOf(q.PenDown),
			"PenUp":     reflect.ValueOf(q.PenUp),
			"Right":     reflect.ValueOf(q.Right),
			"Up":        reflect.ValueOf(q.Up),
			"UpLeft":    reflect.ValueOf(q.UpLeft),
			"UpRight":   reflect.ValueOf(q.UpRight),
		},
		TypedConsts:   map[string]igop.TypedConst{},
		UntypedConsts: map[string]igop.UntypedConst{},
	})
}
