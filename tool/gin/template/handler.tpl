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
	RegisterGroup{{.Name}}(gin.IRouter, ...gin.HandlerFunc)
}

type {{.PackName}}Handler struct {
	s {{.Name}}Service
}

func New{{.Name}}Handler(s {{.Name}}Service) {{.Name}}Handler {
	return &{{.PackName}}Handler{
		s: s,
	}
}

func (h *{{.PackName}}Handler) RegisterGroup{{.Name}}(g gin.IRouter, middleware ...gin.HandlerFunc) {
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
// {{$service.Handler}} {{$service.Doc.Summary}}
// @Summary {{$service.Doc.Summary}}
// @Description {{$service.Doc.Description}}
// @Tags {{$service.Doc.Tag}}
// @Accept {{$service.Doc.Accept}}
// @Produce {{$service.Doc.Produce}}
{{- range $param := $service.Doc.Param}}
// @Param {{$param}}
{{- end}}
// @Success {{$service.Doc.Success}}
{{- if $service.Doc.Failure}}
// @Failure {{$service.Doc.Failure}}
{{- end}}
// @Router {{$service.Doc.Router}}
func (h *{{$outer.PackName}}Handler) {{$service.Handler}}(g gin.IRouter) {
	g.{{$service.Method.Method}}("{{$service.Method.Route}}", h.s.{{$service.Handler}})
}
{{- end -}}

{{- end -}}
