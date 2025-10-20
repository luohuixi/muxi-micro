{{- define "service" -}}
package {{.PackName}}

import (
	"github.com/gin-gonic/gin"
)

type {{.Name}}Service interface {
    {{- range $service := .Service}}
	{{$service.Handler}}(ctx *gin.Context)
	{{- end}}
}

type {{.PackName}}Service struct {
	// TODO: 添加依赖注入逻辑
}

func New{{.Name}}Service() {{.Name}}Service {
	return &{{.PackName}}Service{}
}

{{- end -}}
