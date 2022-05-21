package cmd

import (
	"qa/conf"
	"qa/database/postgres"
	"qa/models"
	"qa/server"

	"github.com/go-chi/chi"
)

func RunUnsecure() error {
	config, err := conf.NewConfig(false)
	if err != nil {
		panic(err)
	}

	s := server.NewServer(config.Port, chi.NewRouter())
	s.SetupRoutes()

	db, err := postgres.NewPostgres(config).Connect()
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	models.RegisterPostgresDB(db)

	return s.StartUnsecure()
}
