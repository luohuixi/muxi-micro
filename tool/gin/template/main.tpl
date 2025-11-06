{{- define "main" -}}
package main

import (
	"github.com/gin-gonic/gin"
)

type App struct {
	g *gin.Engine
}

func NewApp(g *gin.Engine) *App {
	return &App{g: g}
}

func (a *App) Run() {
	a.g.Run("0.0.0.0:8080")
}

func main() {
	//app := InitApp()
	//app.Run()
}

{{- end -}}