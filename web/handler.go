package web

import (
	"log/slog"
	"net/http"
)

type HandlerFunc func(*Context) error

type Middleware func(HandlerFunc) HandlerFunc

func (r *Router) adapter(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		context := newContext(req, w, r.validator)

		err := handler(context)
		if err == nil {
			return
		}

		slog.Error(
			"CRITICAL: Unhandled error reached adapter",
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.Any("err", err),
		)

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
