// export by github.com/goplus/ixgo/cmd/qexp

package pkg

import (
	q "github.com/zrcoder/yaq/hanoi/pkg"

	"reflect"

	"github.com/goplus/ixgo"
)

func init() {
	ixgo.RegisterPackage(&ixgo.Package{
		Name: "pkg",
		Path: "github.com/zrcoder/yaq/hanoi/pkg",
		Deps: map[string]string{
			"charm.land/bubbles/v2/textarea":          "textarea",
			"charm.land/bubbletea/v2":                 "tea",
			"charm.land/lipgloss/v2":                  "lipgloss",
			"errors":                                  "errors",
			"fmt":                                     "fmt",
			"github.com/zrcoder/rdor/pkg/style":       "style",
			"github.com/zrcoder/rdor/pkg/style/color": "color",
			"github.com/zrcoder/yaq":                  "yaq",
			"github.com/zrcoder/yaq/common":           "common",
			"gopkg.in/yaml.v3":                        "yaml",
			"math/rand":                               "rand",
			"os":                                      "os",
			"path/filepath":                           "filepath",
			"strings":                                 "strings",
			"time":                                    "time",
		},
		Interfaces: map[string]reflect.Type{},
		NamedTypes: map[string]reflect.Type{
			"Disk":  reflect.TypeOf((*q.Disk)(nil)).Elem(),
			"Game":  reflect.TypeOf((*q.Game)(nil)).Elem(),
			"Level": reflect.TypeOf((*q.Level)(nil)).Elem(),
			"Pile":  reflect.TypeOf((*q.Pile)(nil)).Elem(),
		},
		AliasTypes: map[string]reflect.Type{},
		Vars: map[string]reflect.Value{
			"Instance": reflect.ValueOf(&q.Instance),
		},
		Funcs: map[string]reflect.Value{
			"A": reflect.ValueOf(q.A),
			"B": reflect.ValueOf(q.B),
			"C": reflect.ValueOf(q.C),
		},
		TypedConsts:   map[string]ixgo.TypedConst{},
		UntypedConsts: map[string]ixgo.UntypedConst{},
	})
}
