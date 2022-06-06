package httphandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Engine struct {
	ginEngine *gin.Engine
}

func NewEngine(engine *gin.Engine) *Engine {
	return &Engine{ginEngine: engine}
}

func (e *Engine) SetRESTRoutes() {
	r := e.ginEngine

	r.GET("/api/v1", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello client!",
		})
	})
}
