package ginx

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/zapx"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/middleware/cors"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/middleware/limiter"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/middleware/timeout"
	"github.com/muxi-Infra/muxi-micro/static"
)

type engineConfig struct {
	env  static.Env
	g    *gin.Engine
	l    logger.Logger
	name string
}

type EngineOption func(*engineConfig)

// 设置服务名称
func WithEnv(env static.Env) EngineOption {
	return func(cfg *engineConfig) {
		cfg.env = env
	}
}

// 手动控制gin的Engine
func WithEngine(g *gin.Engine) EngineOption {
	return func(cfg *engineConfig) {
		cfg.g = g
	}
}

func WithName(name string) EngineOption {
	return func(cfg *engineConfig) {
		cfg.name = name
	}
}
func WithLogger(l logger.Logger) EngineOption {
	return func(cfg *engineConfig) {
		cfg.l = l
	}
}

// 创建默认引擎，附带常用中间件和可选配置
func NewDefaultEngine(opts ...EngineOption) *gin.Engine {
	cfg := &engineConfig{
		env:  static.EnvProd,
		g:    gin.Default(),
		name: DefaultName,
		l:    zapx.NewZapLogger(zapx.WithDefaultZapCore()), //如果不配置logger的话就默认使用zap
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// 非生产环境注册 pprof
	if cfg.env != static.EnvProd {
		pprof.Register(cfg.g)
	}

	cfg.g.Use(
		GlobalLoggerMiddleware(cfg.l),
		GlobalNameMiddleware(cfg.name),
	)

	return cfg.g
}

func UseDefaultMiddleware(g *gin.Engine) {
	g.Use(
		cors.Cors(),
		limiter.Limiter(),
		timeout.Timeout(),
	)
}
