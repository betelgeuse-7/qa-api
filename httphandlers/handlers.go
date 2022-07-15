package httphandlers

import (
	"log"
	"net/http"
	"os"

	"github.com/betelgeuse-7/qa/config"
	"github.com/betelgeuse-7/qa/service/jwtauth"
	"github.com/betelgeuse-7/qa/service/logger"
	"github.com/betelgeuse-7/qa/service/sqlbuild"
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
	userRepo                           models.UserRepository
	jwtRepo                            *jwtauth.TokenRepo
	logger                             *logger.Logger
	domain, rtCookieName, atCookieName string
}

func (e *Engine) SetRESTRoutes(relationalDbConf *config.ConfigRelationalDB, jwtConf *config.ConfigJwt) error {
	r := e.ginEngine
	v1 := r.Group("api/v1")
	pg, err := postgres.New(relationalDbConf)
	if err != nil {
		return err
	}
	err = pg.Connect()
	if err != nil {
		return err
	}
	sqlbuilder := sqlbuild.New()
	userRepo := models.NewUserRepo(pg.Db, sqlbuilder)
	jwtRepo := jwtauth.NewTokenRepo(jwtConf)
	logger := logger.NewLogger(log.Default())
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "127.0.0.1"
		log.Println("[INFO] Server domain is not set. Set to '127.0.0.1' by default")
	}
	h := &Handler{userRepo: userRepo, jwtRepo: jwtRepo, logger: logger, domain: domain, atCookieName: "access-token", rtCookieName: "refresh-token"}
	v1.Use(h.AuthTokenMiddleware)
	v1.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello client!"})
	})
	users := v1.Group("/users")
	users.POST("/", h.NewUser)
	return nil
}
