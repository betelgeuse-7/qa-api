package server

import (
	"qa/handlers"
	//mw "qa/server/middlewares"

	"github.com/go-chi/chi/middleware"
)

func (s *Server) SetupRoutes() {
	s.Router.Use(
		middleware.Logger,
		//mw.JWTAuthorization,
	)

	s.POST("/api/register", handlers.NewUser)
}
