package main

import (
	"github.com/betelgeuse-7/qa/cmd"
	"github.com/betelgeuse-7/qa/config"
)

func main() {
	conf := &config.AppConfig{}
	if err := conf.Parse("./config/conf.json"); err != nil {
		panic("Parsing conf.json: " + err.Error() + "\n")
	}
	cmd.RunQARestAPI(conf)
}
