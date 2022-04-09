package server

import (
	"qa/handlers"
	mw "qa/server/middlewares"

	"github.com/go-chi/chi/middleware"
)

func (s *Server) SetupRoutes() {
	s.Router.Use(
		middleware.Logger,
		mw.ContentTypeJSON,
	)

	s.GET("/api/user/{handle}", handlers.GetUser)
	s.POST("/api/register", handlers.NewUser)
	s.DELETE("/api/user/{handle}", handlers.DeleteUser)
	s.PUT("/api/user/{handle}", handlers.UpdateUser)

	s.POST("/api/question", handlers.NewQuestion)
}
