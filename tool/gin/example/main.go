package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/errs"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx"
)

// 定义请求和响应结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// 模拟用户认证信息
type UserClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func main() {
	g := gin.Default()
	ginx.UseDefaultMiddleware(g)
	router := ginx.NewDefaultEngine(
		ginx.WithEnv(t_http.EnvDev),
		ginx.WithEngine(g),
	)

	// 1. 注册不需要请求体和用户认证的路由
	router.GET("/ping", ginx.Wrap(func(ctx *gin.Context) t_http.Response {
		return t_http.Response{
			HttpCode: http.StatusOK,
			CommonResp: t_http.CommonResp{
				Code:    0,
				Message: "pong",
				Data:    nil,
			},
		}
	}))

	// 2. 注册需要请求体但不需要用户认证的路由
	router.POST("/login", ginx.WrapReq(func(ctx *gin.Context, req LoginRequest) t_http.Response {
		// 模拟登录逻辑
		if req.Username == "admin" && req.Password == "123456" {
			return t_http.Response{
				HttpCode: http.StatusOK,
				CommonResp: t_http.CommonResp{
					Code:    0,
					Message: "登录成功",
					Data:    "token-string",
				},
			}
		}
		return t_http.Response{
			HttpCode: http.StatusUnauthorized,
			CommonResp: t_http.CommonResp{
				Code:    40100,
				Message: "用户名或密码错误",
				Data:    nil,
			},
		}
	}))

	// 3. 注册需要用户认证但不需要请求体的路由
	router.GET("/profile", ginx.WrapClaims(
		// 获取用户认证信息的函数
		func(ctx *gin.Context) (UserClaims, error) {
			// 这里应该从JWT或其他认证方式中获取用户信息
			// 模拟从Header中获取token并解析
			token := ctx.GetHeader("Authorization")
			if token == "" {
				return UserClaims{}, errs.NewErr("认证失败", "缺少Authorization头")
			}

			// 这里应该验证token的有效性
			// 模拟返回用户信息
			return UserClaims{
				UserID:   1,
				Username: "admin",
				Role:     "admin",
			}, nil
		},
		// 业务处理函数
		func(ctx *gin.Context, claims UserClaims) t_http.Response {
			// 模拟获取用户信息
			userInfo := UserInfo{
				ID:       claims.UserID,
				Username: claims.Username,
				Email:    claims.Username + "@example.com",
			}

			return t_http.Response{
				HttpCode: http.StatusOK,
				CommonResp: t_http.CommonResp{
					Code:    0,
					Message: "获取用户信息成功",
					Data:    userInfo,
				},
			}
		},
	))

	// 4. 注册既需要请求体又需要用户认证的路由
	router.POST("/update-profile", ginx.WrapClaimsAndReq(
		// 获取用户认证信息的函数
		func(ctx *gin.Context) (UserClaims, error) {
			// 同上，从Header中获取token并解析
			token := ctx.GetHeader("Authorization")
			if token == "" {
				return UserClaims{}, errs.NewErr("认证失败", "缺少Authorization头")
			}

			// 模拟返回用户信息
			return UserClaims{
				UserID:   1,
				Username: "admin",
				Role:     "admin",
			}, nil
		},
		// 业务处理函数
		func(ctx *gin.Context, req UserInfo, claims UserClaims) t_http.Response {
			// 模拟更新用户信息
			updatedInfo := UserInfo{
				ID:       claims.UserID,
				Username: claims.Username,
				Email:    req.Email, // 使用请求中的新邮箱
			}

			return t_http.Response{
				HttpCode: http.StatusOK,
				CommonResp: t_http.CommonResp{
					Code:    0,
					Message: "更新用户信息成功",
					Data:    updatedInfo,
				},
			}
		},
	))

	router.Run("0.0.0.0:8080")
}
