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
	v1 := r.Group("api/v1")

	v1.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello client!"})
	})

	users := v1.Group("/users")
	users.POST("/", NewUser)
}
