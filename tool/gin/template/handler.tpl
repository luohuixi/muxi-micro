{{- define "handler" -}}
package {{.PackName}}

import (
	"github.com/gin-gonic/gin"
)

const (
	prefix = "{{.Prefix}}"
	group  = "{{.Group}}"
)

type {{.Name}}Handler interface {
    {{- range $service := .Service}}
	{{$service.Handler}}(gin.IRouter)
	{{- end -}}
	{{ "\n" }}
	RunGroup{{.Name}}(gin.IRouter, ...gin.HandlerFunc)
}

type {{.PackName}}Handler struct {
	s {{.Name}}Service
}

func New{{.Name}}Handler(s {{.Name}}Service) {{.Name}}Handler {
	return &{{.PackName}}Handler{
		s: s,
	}
}

func (h *{{.PackName}}Handler) RunGroup{{.Name}}(g gin.IRouter, middleware ...gin.HandlerFunc) {
	addr := prefix + "/" + group
	g.Group(addr, middleware...)
	{
	    {{- range $service := .Service}}
    	h.{{$service.Handler}}(g)
    	{{- end }}
	}
}

{{- $outer := . -}}
{{- range $service := .Service}}
{{ "\n" }}
func (h *{{$outer.PackName}}Handler) {{$service.Handler}}(g gin.IRouter) {
	g.{{$service.Method.Method}}("{{$service.Method.Route}}", h.s.{{$service.Handler}})
}
{{- end -}}

{{- end -}}
