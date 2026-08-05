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
	mux               *http.ServeMux
	routeMiddlewares  []Middleware
	globalMiddlewares []Middleware
}

func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) Use(middleware ...Middleware) {
	r.routeMiddlewares = append(r.routeMiddlewares, middleware...)
}

func (r *Router) UseGlobal(middleware ...Middleware) {
	r.globalMiddlewares = append(r.globalMiddlewares, middleware...)
}

func (r *Router) AddRoute(method HTTPMethod, path string, handler HandlerFunc) {
	pattern := string(method) + " " + path
	lastHandler := handler
	for _, m := range slices.Backward(r.routeMiddlewares) {
		lastHandler = m(lastHandler)
	}
	r.mux.HandleFunc(pattern, adapter(lastHandler))
}

func (r *Router) Group(prefix string, register func(subRouter *Router)) {
	subRouter := &Router{
		mux:              http.NewServeMux(),
		routeMiddlewares: append([]Middleware(nil), r.routeMiddlewares...),
	}

	register(subRouter)

	r.mux.Handle(prefix+"/", http.StripPrefix(prefix, subRouter))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	lastHandler := func(c *Context) error {
		r.mux.ServeHTTP(c.Response, c.Request)
		return nil
	}

	for _, m := range slices.Backward(r.globalMiddlewares) {
		lastHandler = m(lastHandler)
	}

	adapter(lastHandler)(w, req)
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
