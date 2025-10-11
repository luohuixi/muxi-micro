package main

import (
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/engine"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/handler"
	"github.com/muxi-Infra/muxi-micro/static"
	"net/http"

	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
)

// 定义请求和响应结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	g := gin.Default()
	engine.UseDefaultMiddleware(g)
	router := engine.NewEngine(
		engine.WithEnv(static.EnvDev),
	)

	// 1. 注册不需要请求体和用户认证的路由
	router.GET("/ping", handler.Wrap(func(ctx *gin.Context) t_http.Response {
		return t_http.Response{
			HttpCode: http.StatusOK,
			Code:     0,
			Message:  "pong",
			Data:     nil,
		}
	}))

	// 2. 注册需要请求体但不需要用户认证的路由
	router.POST("/login", handler.WrapReq(func(ctx *gin.Context, req LoginRequest) t_http.Response {
		// 模拟登录逻辑
		if req.Username == "admin" && req.Password == "123456" {
			return t_http.Response{
				HttpCode: http.StatusOK,
				Message:  "登录成功",
				Data:     "token-string",
			}
		}
		return t_http.Response{
			HttpCode: http.StatusUnauthorized,
			Code:     40100,
			Message:  "用户名或密码错误",
			Data:     nil,
		}
	}))

	router.Run("0.0.0.0:8080")
}
