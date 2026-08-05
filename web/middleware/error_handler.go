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
			c.Response.WriteError(err.Code, err.Message)
			return nil
		}

		slog.Error(
			"Internal Server Error",
			slog.Any("err", err),
		)

		c.Response.WriteError(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return nil
	}
}
