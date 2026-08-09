package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"standard/web"
)

const defaultMaxAge = 86400

// CORSConfig defines the configuration for CORS middleware.
type CORSConfig struct {
	// Origins is a list of allowed origins. Use "*" to allow all origins.
	Origins []string
	// Methods is a list of allowed HTTP methods. Use "*" to allow all methods.
	Methods []web.HTTPMethod
	// Headers is a list of allowed HTTP headers. Use "*" to allow all headers.
	Headers []string
	// ExposeHeaders is a list of response headers accessible to the client.
	ExposeHeaders []string
	// AllowCredentials indicates whether the request can include user credentials
	// like cookies, HTTP authentication, or client-side SSL certificates.
	// When true, [CORSConfig.Origins], [CORSConfig.Methods], [CORSConfig.Headers] and [CORSConfig.ExposeHeaders] cannot contain the wildcard "*".
	AllowCredentials bool
	// MaxAge is the maximum time in seconds that the results of a preflight request can be cached.
	// Default is 86400 seconds (24 hours). A negative value means no "Access-Control-Max-Age" header will be sent.
	MaxAge int
}

var defaultMethods = []web.HTTPMethod{
	web.GET, web.POST, web.PUT, web.DELETE, web.PATCH, web.OPTIONS,
}

var defaultHeaders = []string{"Content-Type", "Authorization"}

// DefaultCORSConfig returns a permissive CORS configuration without credentials.
// It allows all origins, standard HTTP methods, and common JSON API headers.
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		Origins: []string{"*"},
		Methods: defaultMethods,
		Headers: defaultHeaders,
	}
}

// DefaultCORSConfigWithCredentials returns a CORS configuration with credentials enabled
// for the specified explicit origins.
//
// It allows standard HTTP methods, and common JSON API headers.
func DefaultCORSConfigWithCredentials(origins ...string) *CORSConfig {
	return &CORSConfig{
		Origins:          origins,
		AllowCredentials: true,
		Methods:          defaultMethods,
		Headers:          defaultHeaders,
	}
}

// CORS returns a middleware that handles Cross-Origin Resource Sharing (CORS).
//
// It intercepts preflight OPTIONS requests, sets the appropriate Access-Control-*
// response headers, and handles wildcard verification when credentials are enabled.
//
// If config is nil, an empty [CORSConfig] is used.
func CORS(config *CORSConfig) web.Middleware {
	if config == nil {
		config = &CORSConfig{}
	}

	methodsStr := make([]string, len(config.Methods))
	for i, method := range config.Methods {
		methodsStr[i] = strings.ToUpper(string(method))
	}

	verifyWildcard := func(paramName string, array []string) {
		for _, item := range array {
			if item == "*" {
				panic(fmt.Sprintf("CORS: wildcard '*' is forbidden in the '%s' configuration when AllowCredentials is true", paramName))
			}
		}
	}

	if config.AllowCredentials {
		verifyWildcard("origins", config.Origins)
		verifyWildcard("methods", methodsStr)
		verifyWildcard("headers", config.Headers)
		verifyWildcard("expose headers", config.ExposeHeaders)
	}

	origins := make(map[string]struct{})
	for _, origin := range config.Origins {
		origins[strings.ToLower(origin)] = struct{}{}
	}

	methods := strings.Join(methodsStr, ", ")
	headers := strings.Join(config.Headers, ", ")
	exposeHeaders := strings.Join(config.ExposeHeaders, ", ")

	maxAge := config.MaxAge
	switch {
	case maxAge == 0:
		maxAge = defaultMaxAge
	case maxAge < 0:
		maxAge = 0
	}

	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			incomingOrigin := c.Header.Get("Origin")
			if incomingOrigin == "" {
				return next(c)
			}

			var origin string
			if _, ok := origins[strings.ToLower(incomingOrigin)]; ok {
				origin = incomingOrigin
			} else if _, ok := origins["*"]; ok {
				origin = "*"
			}

			isPreflight := c.Method == http.MethodOptions && c.Header.Get("Access-Control-Request-Method") != ""

			if origin == "" {
				if isPreflight {
					return c.Response.NoContent()
				}

				return next(c)
			}

			header := c.Response.Header()
			header.Set("Access-Control-Allow-Origin", origin)
			if origin != "*" {
				header.Set("Vary", "Origin")
			}

			if config.AllowCredentials {
				header.Set("Access-Control-Allow-Credentials", "true")
			}

			if isPreflight {
				if methods != "" {
					header.Set("Access-Control-Allow-Methods", methods)
				}

				if headers != "" {
					header.Set("Access-Control-Allow-Headers", headers)
				}

				header.Set("Access-Control-Max-Age", strconv.Itoa(maxAge))

				return c.Response.NoContent()
			}

			if exposeHeaders != "" {
				header.Set("Access-Control-Expose-Headers", exposeHeaders)
			}

			return next(c)
		}
	}
}
