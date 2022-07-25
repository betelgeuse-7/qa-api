package httphandlers

import (
	"log"
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
	userRepo             models.UserRepository
	jwtRepo              *jwtauth.TokenRepo
	logger               *logger.Logger
	domain, atCookieName string
	useHTTPS             bool
}

func (e *Engine) SetRESTRoutes(relationalDbConf *config.ConfigRelationalDB, jwtConf *config.ConfigJwt, useHTTPS bool) error {
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
	h := &Handler{userRepo: userRepo, jwtRepo: jwtRepo, logger: logger, domain: domain, atCookieName: "access-token", useHTTPS: useHTTPS}

	v1.POST("/login", h.Login)

	users := v1.Group("/users")
	users.POST("/", h.NewUser)

	secured := v1.Group("/secure")
	secured.Use(h.AuthTokenMiddleware)
	secured.GET("/x", func(ctx *gin.Context) {
		ctx.String(200, " HEEEEEEEEEEEEEEEEEEEEEEEEEY\nuser_id=%d\n", ctx.GetInt64(ContextUserIdKey))
	})
	return nil
}
