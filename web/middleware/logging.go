// Package middleware provides logging middleware for HTTP requests.
package middleware

import (
	"log/slog"
	"time"

	"standard/web"
)

func Logging(next web.HandlerFunc) web.HandlerFunc {
	return func(c *web.Context) {
		start := time.Now()

		next(c)

		slog.Info("HTTP Request",
			slog.String("method", c.Method),
			slog.String("path", c.URL.Path),
			slog.Int("status", c.Response.StatusCode),
			slog.Duration("duration", time.Since(start)),
		)
	}
}
