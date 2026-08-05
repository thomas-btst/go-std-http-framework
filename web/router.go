// Package web provides a simple HTTP router for handling REST API requests.
package web

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"
)

type Middleware func(HandlerFunc) HandlerFunc

type Router struct {
	*http.ServeMux
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

func (r *Router) AddRoute(method HTTPMethod, path string, handler HandlerFunc) {
	pattern := string(method) + " " + path
	lastHandler := handler
	for _, v := range slices.Backward(r.middlewares) {
		lastHandler = v(lastHandler)
	}
	r.HandleFunc(pattern, adapter(lastHandler))
}

func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) Group(prefix string, register func(subRouter *Router)) {
	subRouter := &Router{
		ServeMux:    http.NewServeMux(),
		middlewares: append([]Middleware(nil), r.middlewares...),
	}

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
