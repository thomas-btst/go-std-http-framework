package main

import (
	"log/slog"

	"standard/user"
	"standard/web"
	"standard/web/middleware"
)

var corsConfig = middleware.DefaultCORSConfig()

func main() {
	r := web.NewRouter()

	r.UseGlobal(middleware.Logging)
	r.UseGlobal(middleware.CORS(corsConfig))
	r.Use(middleware.ErrorHandler)

	userStore := user.NewMemoryStore()
	userService := user.NewService(userStore)

	user.NewHandler(userService).RegisterRoutes(r)

	if err := r.ListenAndServe(":8080"); err != nil {
		slog.Error("Server Error", slog.Any("err", err))
	}
}
