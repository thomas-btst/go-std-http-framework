package web

import (
	"log/slog"
	"net/http"
)

type HandlerFunc func(*Context) error

func adapter(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := newContext(r, w)

		err := handler(context)
		if err == nil {
			return
		}

		slog.Error(
			"CRITICAL: Unhandled error reached adapter",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Any("err", err),
		)

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)
	}
}
