package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	Port   string
	Router *chi.Mux
}

// return a new server instance
func NewServer(port int, router *chi.Mux) *Server {
	p := fmt.Sprintf(":%d", port)

	return &Server{
		Port:   p,
		Router: router,
	}
}

func (s *Server) StartUnsecure() error {
	log.Printf("[BOOT] SERVER STARTED LISTENING ON http://localhost%s\n", s.Port)

	return http.ListenAndServe(s.Port, s.Router)
}

func (s *Server) GET(pattern string, handlerFn http.HandlerFunc) {
	s.Router.Get(pattern, handlerFn)
}

func (s *Server) POST(pattern string, handlerFn http.HandlerFunc) {
	s.Router.Post(pattern, handlerFn)
}
