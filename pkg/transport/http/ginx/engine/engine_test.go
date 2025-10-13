package engine

import (
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	"github.com/muxi-Infra/muxi-micro/static"
	"testing"
)

func TestDefaultEngine(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		g := NewEngine(
			WithGinEngine(gin.New()),
			WithEnv(static.EnvDev),
			WithLogger(logx.NewStdLogger()),
			WithName("test"),
		)
		UseDefaultMiddleware(g)
	})

}
