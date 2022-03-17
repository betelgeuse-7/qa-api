package main

import (
	"qa/conf"
	"qa/database"
	"qa/models"
	"qa/server"

	"github.com/go-chi/chi"
)

func main() {
	config, err := conf.NewConfig(false)
	if err != nil {
		panic(err)
	}

	s := server.NewServer(config.Port, chi.NewRouter())
	s.SetupRoutes()

	db, err := database.NewDatabase(database.POSTGRES, config).Connect()
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	models.RegisterPostgresDB(db)

	if err := s.StartUnsecure(); err != nil {
		panic(err)
	}
}
