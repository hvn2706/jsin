package server

import (
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"jsin/config"
	"jsin/logger"
	"jsin/server/middleware"
	"net/http"
)

type Server struct {
	Router *chi.Mux
}

type Option func(s *Server)

func NewServer(cfg config.Config, opts ...Option) *Server {
	s := &Server{
		Router: chi.NewRouter(),
	}
	for _, opt := range opts {
		opt(s)
	}
	s.initRoutes(cfg)
	return s
}

func (s *Server) initRoutes(cfg config.Config) {
	// 1. Health API
	healthRouter := chi.NewRouter()
	s.Router.Mount("/health", healthRouter)
	healthRouter.Get("/ready", ready)
	healthRouter.Get("/liveness", liveness)

	// 2. Public API
	apiV1Router := chi.NewRouter()
	s.Router.Mount("/api/v1", apiV1Router)

	// 3. Public API with API key
	apikeyRouter := chi.NewRouter()
	apikeyRouter.Use(
		middleware.ApiKeyValidateMiddleware(cfg),
	)
	s.Router.Mount("/api/v2", apikeyRouter)
}

func (s *Server) Serve(cfg config.ServerListen) error {
	logger.Infof("Listening on port %v", cfg.Port)
	address := fmt.Sprintf(":%v", cfg.Port)
	return http.ListenAndServe(address, s.Router)
}
