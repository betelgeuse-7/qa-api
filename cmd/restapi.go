package cmd

import (
	"log"

	"github.com/betelgeuse-7/qa/config"
	"github.com/betelgeuse-7/qa/httphandlers"
	"github.com/betelgeuse-7/qa/service/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func RunQARestAPI(httpServerConf config.ConfigHttpServer, relationalDbConf *config.ConfigRelationalDB) {
	if !(httpServerConf.DevMode) {
		// in release/prod mode
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	if !(httpServerConf.UseTLS) {
		// if useTLS is false, use H2C (HTTP/2 without TLS)
		// but browsers don't support H2C, and, even if we use
		// H2C, we still communicate with web browsers via HTTP/1.1
		//
		// So, this doesn't actually change anything.
		r.UseH2C = true
	}
	e := httphandlers.NewEngine(r)
	if err := e.SetRESTRoutes(relationalDbConf); err != nil {
		logger.Error("cmd/restapi.go: couldn't set REST routes: " + err.Error() + "\n")
		return
	}
	log.Fatalln(r.Run(httpServerConf.Port))
}
