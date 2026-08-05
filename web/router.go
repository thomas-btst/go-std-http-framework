// Package web provides a simple HTTP router for handling REST API requests.
package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Router struct {
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

func (r *Router) AddRoute(method HTTPMethod, path string, handler HandlerFunc) {
	pattern := string(method) + " " + path
	r.HandleFunc(pattern, adapter(handler))
}

func (r *Router) Group(prefix string, register func(subRouter *Router)) {
	subRouter := &Router{ServeMux: http.NewServeMux()}

	register(subRouter)

	r.Handle(prefix+"/", http.StripPrefix(prefix, subRouter))
}

func (r *Router) ListenAndServe(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("Starting server", slog.String("addr", addr))
		serverErrors <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-quit:
		slog.Info("Shutdown signal received, shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
