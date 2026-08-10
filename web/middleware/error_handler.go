package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"standard/web"
)

func ErrorHandler(next web.HandlerFunc) web.HandlerFunc {
	return func(c *web.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		if err, ok := errors.AsType[*web.HTTPError](err); ok {
			return c.Response.HTTPError(err)
		}

		slog.Error(
			"Internal Server Error",
			slog.Any("err", err),
		)

		return c.Response.Error(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
	}
}
