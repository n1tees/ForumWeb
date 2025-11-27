package server

import (
	"ForumWeb/internal/db"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func Run() error {

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	router := NewRouter(db.DB)
	router.RegisterRoutes(mux)

	s := &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	return s.start()
}

// gracefull shutdown
func (s *Server) start() error {

	errCh := make(chan error, 1)

	go func() {
		log.Println("HTTP server started on", s.httpServer.Addr)

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Println("shutdown signal:", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			return err
		}
		log.Println("server gracefully stopped")
		return nil

	case err := <-errCh:

		return err
	}
}
