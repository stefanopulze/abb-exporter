package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	http   *http.Server
}

func NewServer() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) Start() {
	s.http = &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}
	go func() {
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("cannot start http server: %v", err))
		}
	}()
	slog.Info("Server http started :8080")
}

func (s *Server) Stop() error {
	if s.http != nil {
		slog.Info("shutting down http server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		return s.http.Shutdown(shutdownCtx)
	}
	return nil
}
