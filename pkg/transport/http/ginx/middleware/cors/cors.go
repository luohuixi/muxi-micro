package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
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

type Option func(cfg *corsCfg)

func WithCorsOrigins(origins []string) Option {
	return func(cfg *corsCfg) {
		cfg.origins = origins
	}
}

func WithCorsAllowMethods(methods ...string) Option {
	return func(cfg *corsCfg) {
		cfg.allowMethods = methods
	}
}

func WithCorsAllowHeaders(headers ...string) Option {
	return func(cfg *corsCfg) {
		cfg.allowHeaders = headers
	}
}

func WithCorsMaxAge(d time.Duration) Option {
	return func(cfg *corsCfg) {
		cfg.maxAge = d
	}
}

// 跨域中间件
func WithCorsAllowCredentials(allow bool) Option {
	return func(cfg *corsCfg) {
		cfg.allowCredentials = allow
	}
}

func Cors(opts ...Option) gin.HandlerFunc {
	cfg := &corsCfg{
		origins:          DefaultOrigins,
		allowMethods:     DefaultAllowMethods,
		allowHeaders:     DefaultAllowHeaders,
		allowCredentials: false,
		maxAge:           DefaultMaxAge,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// 兼容校验
	if cfg.allowCredentials {
		for _, o := range cfg.origins {
			if o == "*" {
				log.Println("[cors] ⚠️ 警告: AllowCredentials=true 时不能使用 AllowOrigins=[\"*\"].")
			}
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     cfg.origins,
		AllowMethods:     cfg.allowMethods,
		AllowHeaders:     cfg.allowHeaders,
		AllowCredentials: cfg.allowCredentials,
		MaxAge:           cfg.maxAge,
	})
}
