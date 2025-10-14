package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/errs"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/log"
	"net/http"
)

const DefaultBindErrCode = 42201

var (
	ErrBindFailed = errs.NewErr("bind fail", "request bind failed")
)

// 解析参数通用函数
func Bind(ctx *gin.Context, req any) error {
	var err error
	// 根据请求方法选择合适的绑定方式
	if ctx.Request.Method == http.MethodGet {
		err = ctx.ShouldBindQuery(req) // 处理GET请求的查询参数
	} else {
		err = ctx.ShouldBind(req) // 处理POST、PUT请求的请求体数据
	}

	if err != nil {
		return ErrBindFailed.WithCause(err)
	}

	return nil
}

// HandleResponse 处理需要自定义业务码的请求
func HandleResponse(ctx *gin.Context, resp t_http.Response) {
	finalResp := t_http.FinalResp{
		Code:    resp.Code,
		Message: resp.Message,
		Data:    resp.Data,
		LogID:   log.GetLogID(ctx),
	}

	ctx.JSON(resp.HttpCode, finalResp)
}

// HandleSuccessResponseWithData 快速处理成功响应
func HandleSuccessResponseWithData(ctx *gin.Context, data any) {
	HandleResponse(ctx, t_http.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// HandleSuccessResponse 快速处理成功响应
func HandleSuccessResponse(ctx *gin.Context) {
	HandleResponse(ctx, t_http.Response{
		Code:    http.StatusOK,
		Message: "success",
	})
}
