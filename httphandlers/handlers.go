package httphandlers

import (
	"net/http"

	"github.com/betelgeuse-7/qa/config"
	"github.com/betelgeuse-7/qa/storage/models"
	"github.com/betelgeuse-7/qa/storage/postgres"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	ginEngine *gin.Engine
}

func NewEngine(engine *gin.Engine) *Engine {
	return &Engine{ginEngine: engine}
}

type Handler struct {
	userRepo models.UserRepository
}

func (e *Engine) SetRESTRoutes(relationalDbConf *config.ConfigRelationalDB) error {
	r := e.ginEngine
	v1 := r.Group("api/v1")
	pg, err := postgres.New(relationalDbConf)
	if err != nil {
		return err
	}
	userRepo := models.NewUserRepo(pg.Db)
	h := &Handler{userRepo: userRepo}
	v1.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello client!"})
	})

	users := v1.Group("/users")
	users.POST("/", h.NewUser)
	return nil
}
