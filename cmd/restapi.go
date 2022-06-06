package cmd

import (
	"log"

	"github.com/betelgeuse-7/qa/config"
	"github.com/betelgeuse-7/qa/httphandlers"
	"github.com/gin-gonic/gin"
)

func RunQARestAPI(httpServerConf config.ConfigHttpServer) {
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
	e.SetRESTRoutes()
	log.Fatalln(r.Run(httpServerConf.Port))
}
