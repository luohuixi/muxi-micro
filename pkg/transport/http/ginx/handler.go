package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/errs"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
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
		err = ctx.ShouldBind(req) // 处理POST、PUT等请求的请求体数据
	}

	if err != nil {
		return ErrBindFailed.WithCause(err)
	}

	return nil
}

// WrapReq 。用于处理有请求体的请求
// ctx表示上下文,req表示请求结构体
func WrapReq[Req any](fn func(*gin.Context, Req) t_http.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//检查前置中间件是否存在错误,如果存在应当直接返回
		if len(ctx.Errors) > 0 {
			return
		}
		//解析参数
		var req Req
		err := Bind(ctx, &req)
		if err != nil {
			HandleResponse(ctx, t_http.Response{
				HttpCode: http.StatusBadRequest,
				Code:     DefaultBindErrCode,
				Message:  "非法的参数: " + err.Error(),
				Data:     nil,
			})
			return
		}
		// 调用业务逻辑函数
		resp := fn(ctx, req)
		HandleResponse(ctx, resp)
		return
	}
}

// Wrap 。用于处理没有请求体的请求
// ctx表示上下文
func Wrap(fn func(*gin.Context)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//检查前置中间件是否存在错误,如果存在应当直接返回
		if len(ctx.Errors) > 0 {
			return
		}
		fn(ctx)
		return
	}
}

// 处理需要自定义业务码的请求
func HandleResponse(ctx *gin.Context, resp t_http.Response) {
	finalResp := t_http.FinalResp{
		Code:    resp.Code,
		Message: resp.Message,
		Data:    resp.Data,
		LogID:   GetLogID(ctx),
	}

	ctx.JSON(resp.HttpCode, finalResp)
}
