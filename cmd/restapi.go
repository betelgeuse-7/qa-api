package cmd

import (
	"log"

	"github.com/betelgeuse-7/qa/config"
	"github.com/betelgeuse-7/qa/httphandlers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func RunQARestAPI(conf *config.AppConfig) {
	if !(conf.HttpServer.DevMode) {
		// in release/prod mode
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	if !(conf.HttpServer.UseTLS) {
		// if useTLS is false, use H2C (HTTP/2 without TLS)
		// but browsers don't support H2C, and, even if we use
		// H2C, we still communicate with web browsers via HTTP/1.1
		//
		// So, this doesn't actually change anything.
		r.UseH2C = true
	}
	e := httphandlers.NewEngine(r)
	if err := e.SetRESTRoutes(&conf.RelationalDB, &conf.Auth.Jwt, conf.HttpServer.UseTLS); err != nil {
		log.Printf("[ERROR] cmd/restapi.go: couldn't set REST routes: %s\n", err.Error())
		return
	}
	log.Fatalln(r.Run(conf.HttpServer.Port))
}
