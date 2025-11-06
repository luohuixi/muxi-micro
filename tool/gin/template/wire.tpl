{{- define "wire" -}}
//go:generate wire
//go:build wireinject

package main

import (
	"github.com/google/wire"
	{{- $outer := . -}}
    {{- range $name := .Name}}
    "{{$outer.Project}}/handler/{{$name.Min}}"
    {{- end}}
    "{{$outer.Project}}/router"
)

func InitApp() *App {
	wire.Build(
	    {{- range $name := .Name}}
    	{{$name.Min}}.New{{$name.Max}}Service,
    	{{$name.Min}}.New{{$name.Max}}Handler,
    	{{end}}
		router.RegisterGin,

		NewApp,
	)
	return &App{}
}
{{- end -}}
