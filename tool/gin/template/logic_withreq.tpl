{{- define "logic_withreq" -}}
package {{.PackName}}

import (
	"net/http"

	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/handler"
)

func (s *{{.PackName}}Service) {{.Handler}}(ctx *gin.Context) {
	var req {{.Req}}
	err := handler.Bind(ctx, &req)
	if err != nil {
		handler.HandleResponse(ctx, t_http.Response{
			HttpCode: http.StatusBadRequest,
			Code:     handler.DefaultBindErrCode,
			Message:  "非法的参数: " + err.Error(),
			Data:     nil,
		})
		return
	}
	// TODO: 写你的业务逻辑
}

{{- end -}}
