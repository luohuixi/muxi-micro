{{- define "router" -}}
package router

import (
    {{- $outer := . -}}
	{{- range $name := .Name}}
    "{{$outer.Project}}/handler/{{$name.Min}}"
    {{- end}}
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/engine"
)

func Run(addr string) error {
	g := engine.NewEngine()
	engine.UseDefaultMiddleware(g)
    {{range $name := .Name}}
	s{{$name.Min}} := {{$name.Min}}.New{{$name.Max}}Service()
	h{{$name.Min}} := {{$name.Min}}.New{{$name.Max}}Handler(s{{$name.Min}})
	h{{$name.Min}}.RunGroup{{$name.Max}}(g)
	{{end}}
	if err := g.Run(addr); err != nil {
		return err
	}

	return nil
}

{{- end -}}
