package main

import (
	"log/slog"

	"standard/user"
	"standard/web"
)

func main() {
	r := web.NewRouter()

	userStore := user.NewMemoryStore()
	userService := user.NewService(userStore)

	user.NewHandler(userService).RegisterRoutes(r)

	if err := r.ListenAndServe(":8080"); err != nil {
		slog.Error("Server Error", slog.Any("err", err))
	}
}
