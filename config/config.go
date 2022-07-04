package config

import (
	"encoding/json"
	"errors"
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
	err = json.Unmarshal(bx, a)
	if err != nil {
		return err
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) == 0 {
		return errors.New("env variable 'JWT_SECRET' is not set")
	}
	a.Auth.Jwt.SecretKey = []byte(jwtSecret)
	return nil
}

type ConfigRelationalDB struct {
	Name, Host, User, Ssl, DbName string
	Port                          uint
}

type ConfigAuth struct {
	Jwt ConfigJwt
}

type ConfigJwt struct {
	SecretKey []byte
}

type ConfigHttpServer struct {
	HttpVersion string
	Port        string
	UseTLS      bool
	DevMode     bool
}
