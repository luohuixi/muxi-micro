package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	DefaultOrigins      = []string{"*"}
	DefaultAllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	DefaultAllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	DefaultMaxAge       = 12 * time.Hour
)

type corsCfg struct {
	origins          []string
	allowMethods     []string
	allowHeaders     []string
	allowCredentials bool
	maxAge           time.Duration
}

type CorsOption func(cfg *corsCfg)

func WithCorsOrigins(origins []string) CorsOption {
	return func(cfg *corsCfg) {
		cfg.origins = origins
	}
}

func WithCorsAllowMethods(methods ...string) CorsOption {
	return func(cfg *corsCfg) {
		cfg.allowMethods = methods
	}
}

func WithCorsAllowHeaders(headers ...string) CorsOption {
	return func(cfg *corsCfg) {
		cfg.allowHeaders = headers
	}
}

func WithCorsAllowCredentials(allow bool) CorsOption {
	return func(cfg *corsCfg) {
		cfg.allowCredentials = allow
	}
}

func WithCorsMaxAge(d time.Duration) CorsOption {
	return func(cfg *corsCfg) {
		cfg.maxAge = d
	}
}

// 跨域中间件

func Cors(opts ...CorsOption) gin.HandlerFunc {
	cfg := &corsCfg{
		origins:          DefaultOrigins,
		allowMethods:     DefaultAllowMethods,
		allowHeaders:     DefaultAllowHeaders,
		allowCredentials: true,
		maxAge:           DefaultMaxAge,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cors.New(cors.Config{
		AllowOrigins:     cfg.origins,
		AllowMethods:     cfg.allowMethods,
		AllowHeaders:     cfg.allowHeaders,
		AllowCredentials: cfg.allowCredentials,
		MaxAge:           cfg.maxAge,
	})
}
