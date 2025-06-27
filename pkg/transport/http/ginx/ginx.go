package ginx

import (
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/errs"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/middleware/cors"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/middleware/limiter"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/middleware/timeout"

	"net/http"
)

var shouldBindErr = errs.NewErr("bind fail", "ctx shouldBind failed")

var (
	defaultBindErrCode      = 42201
	defaultGetClaimsErrCode = 40101
)

type engineConfig struct {
	env         t_http.Env
	enablePprof bool
	g           *gin.Engine
}

type EngineOption func(*engineConfig)

// 设置运行环境
func WithEnv(env t_http.Env) EngineOption {
	return func(cfg *engineConfig) {
		cfg.env = env
		// dev 和 test 环境默认开启 pprof
		if env == t_http.EnvDev || env == t_http.EnvTest {
			cfg.enablePprof = true
		}
	}
}

// 手动控制是否开启 pprof
func WithPprof(enable bool) EngineOption {
	return func(cfg *engineConfig) {
		cfg.enablePprof = enable
	}
}

// 手动控制gin的Engine
func WithEngine(g *gin.Engine) EngineOption {
	return func(cfg *engineConfig) {
		cfg.g = g
	}
}

// 创建默认引擎，附带常用中间件和可选配置
func NewDefaultEngine(opts ...EngineOption) *gin.Engine {
	cfg := &engineConfig{
		env:         t_http.EnvProd,
		enablePprof: false,
		g:           gin.Default(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// 非生产环境注册 pprof
	if cfg.enablePprof {
		pprof.Register(cfg.g)
	}

	return cfg.g
}

func UseDefaultMiddleware(g *gin.Engine) {
	g.Use(
		cors.Cors(),
		limiter.Limiter(),
		timeout.Timeout(),
	)
}

func SetBindErrCode(errCode int) {
	defaultBindErrCode = errCode
}

func SetGetClaimsErrCode(errCode int) {
	defaultGetClaimsErrCode = errCode
}

// 解析参数通用函数
func bind(ctx *gin.Context, req any) error {
	var err error
	// 根据请求方法选择合适的绑定方式
	if ctx.Request.Method == http.MethodGet {
		err = ctx.ShouldBindQuery(req) // 处理GET请求的查询参数
	} else {
		err = ctx.ShouldBind(req) // 处理POST、PUT等请求的请求体数据
	}

	if err != nil {
		return shouldBindErr.WithCause(err)
	}

	return nil
}

func WrapClaimsAndReq[Req any, UserClaims any](getClaims func(ctx context.Context) (UserClaims, error), fn func(*gin.Context, Req, UserClaims) t_http.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 检查前置中间件是否存在错误,如果存在应当直接返回
		if len(ctx.Errors) > 0 {
			return
		}

		//解析请求
		var req Req
		err := bind(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, t_http.CommonResp{
				Code:    defaultBindErrCode,
				Message: "非法的参数: " + err.Error(),
				Data:    nil,
			})
			return
		}

		//获取uc参数
		uc, err := getClaims(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, t_http.CommonResp{
				Code:    defaultGetClaimsErrCode,
				Message: "登陆状态异常:" + err.Error(),
				Data:    nil,
			})
			return
		}

		//执行函数
		res := fn(ctx, req, uc)
		ctx.JSON(res.HttpCode, res.CommonResp)
		return
	}
}

// WrapReq 。用于处理有请求体的请求
// ctx表示上下文,req表示请求结构体,Resp表示响应结构体(这里全部填web.Response)
func WrapReq[Req any](fn func(*gin.Context, Req) t_http.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//检查前置中间件是否存在错误,如果存在应当直接返回
		if len(ctx.Errors) > 0 {
			return
		}

		//解析参数
		var req Req
		err := bind(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, t_http.CommonResp{
				Code:    defaultBindErrCode,
				Message: "非法的参数: " + err.Error(),
				Data:    nil,
			})
			return
		}

		// 调用业务逻辑函数
		res := fn(ctx, req)
		ctx.JSON(res.HttpCode, res.CommonResp)
		return
	}
}

// Wrap 。用于处理没有请求体的请求
// ctx表示上下文,Resp表示响应结构体(这里全部填web.Response)
func Wrap(fn func(*gin.Context) t_http.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//检查前置中间件是否存在错误,如果存在应当直接返回
		if len(ctx.Errors) > 0 {
			return
		}

		res := fn(ctx)
		ctx.JSON(res.HttpCode, res.CommonResp)
		return
	}
}

// WrapClaims 用于处理有用户验证但是没有请求体的请求
// ctx表示上下文,Resp表示响应结构体(这里全部填web.Response),UserClaims表示用户信息
func WrapClaims[UserClaims any](getClaims func(ctx context.Context) (UserClaims, error), fn func(*gin.Context, UserClaims) t_http.Response) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//检查前置中间件是否存在错误,如果存在应当直接返回
		if len(ctx.Errors) > 0 {
			return
		}

		//获取uc参数
		uc, err := getClaims(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, t_http.CommonResp{
				Code:    defaultGetClaimsErrCode,
				Message: "登陆状态异常:" + err.Error(),
				Data:    nil,
			})
			return
		}

		//执行函数
		res := fn(ctx, uc)
		ctx.JSON(res.HttpCode, res.CommonResp)
		return
	}
}
