{{- define "logic_noreq" -}}
package {{.PackName}}

import (
	"github.com/gin-gonic/gin"
)

func (s *{{.PackName}}Service) {{.Handler}}(ctx *gin.Context) {
	// TODO: 写你的业务逻辑
}

{{- end -}}