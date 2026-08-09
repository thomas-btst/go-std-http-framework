package web

import "net/http"

type HTTPMethod string

const (
	GET     HTTPMethod = http.MethodGet
	POST    HTTPMethod = http.MethodPost
	PUT     HTTPMethod = http.MethodPut
	DELETE  HTTPMethod = http.MethodDelete
	HEAD    HTTPMethod = http.MethodHead
	OPTIONS HTTPMethod = http.MethodOptions
	PATCH   HTTPMethod = http.MethodPatch
	CONNECT HTTPMethod = http.MethodConnect
	TRACE   HTTPMethod = http.MethodTrace
)

func (r *Router) GET(path string, h HandlerFunc) {
	r.AddRoute(GET, path, h)
}

func (r *Router) POST(path string, h HandlerFunc) {
	r.AddRoute(POST, path, h)
}

func (r *Router) PUT(path string, h HandlerFunc) {
	r.AddRoute(PUT, path, h)
}

func (r *Router) DELETE(path string, h HandlerFunc) {
	r.AddRoute(DELETE, path, h)
}

func (r *Router) PATCH(path string, h HandlerFunc) {
	r.AddRoute(PATCH, path, h)
}

func (r *Router) HEAD(path string, h HandlerFunc) {
	r.AddRoute(HEAD, path, h)
}

func (r *Router) OPTIONS(path string, h HandlerFunc) {
	r.AddRoute(OPTIONS, path, h)
}

func (r *Router) Any(path string, h HandlerFunc) {
	r.AddRoute("", path, h)
}
