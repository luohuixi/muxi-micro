package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muxi-Infra/muxi-micro/pkg/transport/http/ginx/engine"
	"github.com/muxi-Infra/muxi-micro/static"
)

func main() {
	g := gin.Default()
	engine.UseDefaultMiddleware(g)
	router := engine.NewEngine(
		engine.WithEnv(static.EnvDev),
	)

	router.Run("0.0.0.0:8080")
}
