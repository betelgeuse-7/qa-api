package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type AppConfig struct {
	RelationalDB ConfigRelationalDB
	Auth         ConfigAuth
	HttpServer   ConfigHttpServer
}

func (a *AppConfig) Parse(file string) error {
	if ext := filepath.Ext(file); ext != ".json" {
		return fmt.Errorf("expected a json file, got %s", ext)
	}
	bx, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bx, a)
}

type ConfigRelationalDB struct {
	Name, Host, User, Ssl, DbName string
	Port                          uint
}

type ConfigAuth struct {
	Jwt ConfigJwt
}

type ConfigJwt struct {
	SecretKeyFile, Expiry string
}

type ConfigHttpServer struct {
	HttpVersion string
	Port        string
	UseTLS      bool
	DevMode     bool
}
