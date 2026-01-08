{{- define "router" -}}
package router

import (
    {{- $outer := . -}}
	{{- range $name := .Name}}
    "YourPath/handler/{{$name.Min}}"
    {{- end}}
    "github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/engine"
)

func RegisterGin(
    {{- range $name := .Name}}
	h{{$name.Min}} {{$name.Min}}.{{$name.Max}}Handler,
	{{- end}}
) *gin.Engine {
	g := engine.NewEngine()
	engine.UseDefaultMiddleware(g)
    {{range $name := .Name}}
	h{{$name.Min}}.RegisterGroup{{$name.Max}}(g, /*middleware*/)
	{{- end}}

	return g
}

{{- end -}}
