package cmd

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/betelgeuse-7/qa/config"
	"github.com/betelgeuse-7/qa/httphandlers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	_ROOT_PRIVILEGED_PORTS_END = uint(1025)
	_PORTS_END                 = uint(65535)
	_USAGE                     = "\t*---------------------------------------------------------------------*\n" +
		"\t\t-tls, --tls\t\t\t\tUse TLS\n" +
		"\t\t-port, --port, -port=, --port=\t\tPort number\n" +
		"\t\t-cert, --cert, -cert=, --cert=\t\tSSL certificate path\n" +
		"\t\t-key, --key, -key=, --key=\t\tSSL key path\n" +
		"\t\t-help, --help\t\t\t\tGet help\n" +
		"\t*---------------------------------------------------------------------*\n"
)

func RunQARestAPI(conf *config.AppConfig) {
	flags, err := __start()
	if err != nil {
		log.Printf("[ERROR] cmd/restapi.go: couldn't set REST routes: %s", err.Error())
		return
	}
	if flags.help {
		fmt.Print(_USAGE)
		return
	}
	conf.HttpServer.UseTLS = flags.useTLS
	conf.HttpServer.Port = flags.port
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
	if conf.HttpServer.UseTLS {
		log.Fatalln(r.RunTLS(conf.HttpServer.Port, flags.sslCert, flags.sslKey))
	}
	log.Fatalln(r.Run(conf.HttpServer.Port))
}

type flags struct {
	useTLS, help          bool
	port, sslCert, sslKey string
}

func __start() (flags, error) {
	var res flags
	help := flag.Bool("help", false, "help")
	useTLS := flag.Bool("tls", false, "use tls")
	portUint := flag.Uint("port", 8000, "port to listen on")
	certLocation := flag.String("cert", "", "ssl certificate location")
	keyLocation := flag.String("key", "", "ssl key location")
	flag.Parse()
	if !(*portUint >= _ROOT_PRIVILEGED_PORTS_END && *portUint <= _PORTS_END) {
		return res, fmt.Errorf("invalid port (out of range [1025,65535]): '%d'", *portUint)
	}
	port := fmt.Sprintf(":%d", *portUint)
	res.useTLS = *useTLS
	res.port = port
	res.sslCert = *certLocation
	res.sslKey = *keyLocation
	res.help = *help
	if res.useTLS {
		var errStr strings.Builder
		if len(res.sslCert) == 0 {
			errStr.WriteString("\n\t- Missing SSL certificate path")
		}
		if len(res.sslKey) == 0 {
			errStr.WriteString("\n\t- Missing SSL key path\n")
		}
		if len(errStr.String()) != 0 {
			return res, fmt.Errorf("%s", errStr.String())
		}
	}
	return res, nil
}
