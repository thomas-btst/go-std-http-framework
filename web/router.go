// Package web provides a simple HTTP router for handling REST API requests.
package web

import (
	"net/http"
	"slices"
)

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
	var pattern string
	if method == "" {
		pattern = path
	} else {
		pattern = string(method) + " " + path
	}

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
