# Go Std HTTP Framework (`go-std-http-framework`)

A lightweight, zero-router-dependency HTTP framework and REST API template built directly on top of Go's standard library (`net/http` with Go 1.22+ dynamic routing).

Designed to bridge the gap between pure `net/http` and full-featured web frameworks (like Gin or Echo), **`go-std-http-framework`** provides an ergonomic developer experience (`Context`, request binding, validation, middleware chaining, CORS, JSON responses) while remaining fast, clean, and minimal.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation & Setup](#installation--setup)
  - [Run the Application](#run-the-application)
- [Core Concepts & Code Examples](#core-concepts--code-examples)
  - [Server Initialization (`main.go`)](#server-initialization-maingo)
  - [Routing & HTTP Method Helpers](#routing--http-method-helpers)
  - [Request Binding & Automatic Validation](#request-binding--automatic-validation)
  - [Response Helpers](#response-helpers)
  - [Middleware Pipeline & Custom Middlewares](#middleware-pipeline--custom-middlewares)
  - [Error Handling](#error-handling)
- [Clean Architecture Domain Example](#clean-architecture-domain-example)
  - [REST API Reference](#rest-api-reference)
- [Testing](#testing)
- [Documentation](#documentation)
- [Code Quality & Linting](#code-quality--linting)
- [Author](#author)
- [License](#license)

## Features

- **Native Routing Engine**: Built directly on Go 1.22+ `http.ServeMux`.
- **HTTP Method Helpers**: Dedicated router methods (`GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, `OPTIONS`, `Any`, and `AddRoute`).
- **Route Grouping**: Prefix-based sub-routers (`/users`, `/api/v1`) with inherited route middleware and validator instances.
- **Two-Tier Middleware Pipeline**: Separate support for **Global** middleware (wrapping the entire HTTP server) and **Route-specific** middleware (executed per endpoint).
- **Built-in Middleware Suite**:
  - `Logging`: Structured request logger with execution time, status code, and path using `log/slog`.
  - `CORS`: Fully configurable CORS middleware supporting preflight OPTIONS requests, allowed origins, headers, methods, credentials, and cache max-age.
  - `ErrorHandler`: Centralized error interceptor converting `HTTPError` and unhandled errors into uniform JSON responses.
- **Custom Middleware Support**: Straightforward middleware interface (`func(web.HandlerFunc) web.HandlerFunc`) for custom authentication, rate limiting, and request transformation.
- **Ergonomic `web.Context`**: Wraps `*http.Request` and `*web.Response`, providing helpers for typed path parameter extraction (`PathInt`, `PathString`), JSON decoding, and validation.
- **Robust Struct Validation**: Integrated validation via `go-playground/validator/v10` with translated English error messages, struct JSON tag mapping, and pluggable `Validator` interface.
- **Fluent Response Writer (`web.Response`)**: Native helpers for `JSON()`, `NoContent()`, `Status()`, `HTTPError()`, and `Error()`.
- **Graceful Shutdown Server**: `web.NewServer()` encapsulates `http.Server` and intercepts SIGINT/SIGTERM signals to perform clean shutdown with a 10-second timeout.
- **Production-Ready Clean Architecture**: Example `user` domain showcasing complete separation of concerns (Store with `sync.RWMutex`, Service with `bcrypt` password hashing, Handler, DTOs, and mapping).

## Technologies

- **Go Standard Library**:
  - `net/http`: Core HTTP server, `http.ServeMux` dynamic pattern routing, request and response management.
  - `log/slog`: High-performance structured logging.
  - `context`: Request lifecycle management, timeout handling, and graceful shutdown signal propagation.
  - `encoding/json`: JSON serialization and deserialization for API communication.
  - `sync`: Thread-safe concurrency primitives (`sync.RWMutex`) for in-memory data storage.
- **Third-Party Libraries**:
  - [`go-playground/validator/v10`](https://github.com/go-playground/validator): Struct and field validation engine for incoming request payloads.
  - [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt): Secure password hashing.
- **Tooling & Documentation**:
  - `pkgsite`: Standard Go package documentation server.

## Project Structure

```text
.
├── main.go                     # Application entrypoint & dependency injection
├── go.mod                      # Module definition and dependencies
├── go.sum                      # Checksums for dependencies
├── web/                        # Core HTTP Framework package
│   ├── context.go              # Context struct, Bind(), Body(), PathInt(), PathString()
│   ├── errors.go               # HTTPError custom error type with status and details
│   ├── handler.go              # HandlerFunc and Middleware definitions, adapter
│   ├── methods.go              # HTTPMethod constants and helper methods (GET, POST, etc.)
│   ├── response.go             # Response wrapper with JSON(), NoContent(), Status()
│   ├── router.go               # Router, Route Grouping & middleware execution
│   ├── server.go               # Server wrapper with built-in Graceful Shutdown
│   ├── validator.go            # Validator interface and DefaultValidator (v10)
│   └── middleware/             # Built-in middlewares
│       ├── cors.go             # CORS middleware and configuration presets
│       ├── error_handler.go    # Centralized HTTP error interceptor
│       └── logging.go          # Structured request logger with slog
└── user/                       # Example domain (Clean Architecture)
    ├── handler.go              # HTTP Handlers and route registration
    ├── mapping.go              # DTO to domain model mapping helpers
    ├── requests.go             # Request payload DTOs with validation tags
    ├── responses.go            # Response DTOs (sanitized, no sensitive data)
    ├── service.go              # Business logic with bcrypt password hashing
    ├── store.go                # Store interface & thread-safe in-memory implementation
    └── user.go                 # Domain entity model
```

## Getting Started

### Prerequisites

- Go 1.22 or higher installed on your system.
- Git for cloning the repository.

### Installation & Setup

Clone the repository and install dependencies:

```bash
git clone https://github.com/thomas-btst/go-std-http-framework.git
cd go-std-http-framework
go mod download
```

### Run the Application

```bash
go run main.go
```

The server will start listening on `http://localhost:8080`.

## Core Concepts & Code Examples

### Server Initialization (`main.go`)

```go
package main

import (
	"log/slog"

	"github.com/thomas-btst/go-std-http-framework/user"
	"github.com/thomas-btst/go-std-http-framework/web"
	"github.com/thomas-btst/go-std-http-framework/web/middleware"
)

var corsConfig = middleware.DefaultCORSConfig()

func main() {
	r := web.NewRouter()

	// Global middlewares (executed on every incoming request)
	r.UseGlobal(middleware.Logging)
	r.UseGlobal(middleware.CORS(corsConfig))

	// Route middleware (executed for registered routes)
	r.Use(middleware.ErrorHandler)

	// Dependency Injection & Domain Route Registration
	userStore := user.NewMemoryStore()
	userService := user.NewService(userStore)
	user.NewHandler(userService).RegisterRoutes(r)

	// Start server with Graceful Shutdown
	server := web.NewServer(":8080", r)
	if err := server.Start(); err != nil {
		slog.Error("Server Error", slog.Any("err", err))
	}
}
```

### Routing & HTTP Method Helpers

Register routes using dedicated HTTP method helpers or grouped prefixes:

```go
r := web.NewRouter()

// Direct HTTP method helpers
r.GET("/health", func(c *web.Context) error {
	return c.Response.JSON(http.StatusOK, map[string]string{"status": "healthy"})
})

// Route Grouping
r.Group("/api/v1", func(sub *web.Router) {
	sub.GET("/items", listItemsHandler)
	sub.POST("/items", createItemHandler)
	sub.GET("/items/{id}", getItemHandler)
	sub.PUT("/items/{id}", updateItemHandler)
	sub.DELETE("/items/{id}", deleteItemHandler)
})
```

### Request Binding & Automatic Validation

`c.Bind(&req)` decodes the JSON request body and runs struct validation automatically:

```go
type CreateItemRequest struct {
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}

func CreateItemHandler(c *web.Context) error {
	var req CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return err // Automatically converted to 422 Unprocessable Entity with details
	}

	id, err := c.PathInt("id")
	if err != nil {
		return err // Automatically converted to 400 Bad Request
	}

	return c.Response.JSON(http.StatusCreated, map[string]any{
		"id":    id,
		"name":  req.Name,
		"email": req.Email,
	})
}
```

### Response Helpers

`web.Response` provides clean utilities for sending API responses:

```go
// Send JSON with status code
c.Response.JSON(http.StatusOK, user)

// Send 204 No Content
c.Response.NoContent()

// Send custom HTTP error
c.Response.Error(http.StatusForbidden, "Access forbidden")
```

### Middleware Pipeline & Custom Middlewares

The framework supports both global and route-level middlewares, built-in CORS configurations, and custom middleware handlers.

#### Built-in CORS Configuration

```go
// Permissive default CORS (allows all origins)
r.UseGlobal(middleware.CORS(middleware.DefaultCORSConfig()))

// CORS with credentials and specific origins
corsConfig := middleware.DefaultCORSConfigWithCredentials("https://example.com")
r.UseGlobal(middleware.CORS(corsConfig))
```

#### Writing Custom Middlewares

Middlewares implement the `web.Middleware` type signature (`func(web.HandlerFunc) web.HandlerFunc`):

```go
// Example: Header authentication middleware
func RequireAuth(next web.HandlerFunc) web.HandlerFunc {
	return func(c *web.Context) error {
		authHeader := c.Header.Get("Authorization")
		if authHeader == "" {
			return web.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		// Proceed to next handler
		return next(c)
	}
}

// Attach to specific route groups or globally
r.Use(RequireAuth)
```

### Error Handling

Return `web.HTTPError` directly from any handler function. The `ErrorHandler` middleware formats it into a uniform JSON response:

```go
func GetUser(c *web.Context) error {
	id, err := c.PathInt("id")
	if err != nil {
		return err
	}

	user, err := userService.Get(c.Context(), id)
	if err != nil {
		return web.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.Response.JSON(http.StatusOK, user)
}
```

Standard error response:

```json
{
  "status": 404,
  "message": "User not found"
}
```

Validation error response (422 Unprocessable Entity):

```json
{
  "status": 422,
  "message": "Request body validation failed",
  "details": {
    "name": "name must be at least 2 characters in length",
    "password": "password is a required field"
  }
}
```

## Clean Architecture Domain Example

The included `user/` package demonstrates Clean Architecture best practices:

- **Domain Entity (`user/user.go`)**: Core business model.
- **DTOs & Mapping (`user/requests.go`, `user/responses.go`, `user/mapping.go`)**: Clean separation between network payloads and internal domain objects.
- **Store Layer (`user/store.go`)**: Interface-based repository with a thread-safe in-memory implementation (`sync.RWMutex`).
- **Service Layer (`user/service.go`)**: Business logic orchestrator with bcrypt password hashing.
- **Handler Layer (`user/handler.go`)**: HTTP transport mapping endpoints to service calls.

### REST API Reference

The `user` domain exposes the following CRUD endpoints:

| Method | Endpoint | Description | Request Body | Success Status |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/users/` | List all users | None | `200 OK` |
| `GET` | `/users/{id}` | Get user by ID | None | `200 OK` |
| `POST` | `/users/` | Create a new user | `{"name": "...", "password": "..."}` | `201 Created` |
| `PUT` | `/users/{id}` | Update an existing user | `{"name": "...", "password": "..."}` | `204 No Content` |
| `DELETE` | `/users/{id}` | Delete a user by ID | None | `204 No Content` |

## Testing

Run unit and package tests using Go's built-in testing tool:

```bash
# Run all tests
go test -v ./...

# Run tests with race detection enabled
go test -v -race ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

Generate and view package documentation locally using `pkgsite`:

```bash
go install golang.org/x/pkgsite/cmd/pkgsite@latest
pkgsite -http=:6060
```

Visit `http://localhost:6060/github.com/thomas-btst/go-std-http-framework` to browse the API documentation.

## Code Quality & Linting

Maintain code health, consistency, and static analysis with the following tools:

### Formatting & Vetting

Ensure standard Go formatting and check for suspicious constructs:

```bash
# Format all files
go fmt ./...

# Static analysis with Go vet
go vet ./...
```

### Linter (`golangci-lint`)

Run `golangci-lint` to catch potential bugs, style issues, and enforce best practices:

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## Author

Name : Thomas BATISTA  
Website : [thomas-batista.fr](https://thomas-batista.fr)

## License

© Thomas BATISTA. All rights reserved.
