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
// {{$service.Handler}} {{index $service.Doc "summary" 0}}
{{- range $summary := index $service.Doc "summary"}}
// @Summary {{$summary}}
{{- end}}
{{- range $description := index $service.Doc "description"}}
// @Description {{$description}}
{{- end}}
{{- range $tag := index $service.Doc "tag"}}
// @Tags {{$tag}}
{{- end}}
{{- range $accept := index $service.Doc "accept"}}
// @Accept {{$accept}}
{{- end}}
{{- range $produce := index $service.Doc "produce"}}
// @Produce {{$produce}}
{{- end}}
{{- range $param := index $service.Doc "param"}}
// @Param {{$param}}
{{- end}}
{{- range $success := index $service.Doc "success"}}
// @Success {{$success}}
{{- end}}
{{- range $failure := index $service.Doc "failure"}}
// @Failure {{$failure}}
{{- end}}
{{- range $router := index $service.Doc "router"}}
// @Router {{$router}}
{{- end}}
func (h *{{$outer.PackName}}Handler) {{$service.Handler}}(g gin.IRouter) {
	g.{{$service.Method.Method}}("{{$service.Method.Route}}", h.s.{{$service.Handler}})
}
{{- end -}}

{{- end -}}
