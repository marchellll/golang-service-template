package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

// RequestIDMiddleware adds a unique request ID to each request
// This should be used early in the middleware chain for proper request tracing
func RequestIDMiddleware() echo.MiddlewareFunc {
	return echo_middleware.RequestIDWithConfig(echo_middleware.RequestIDConfig{
		Generator: func() string {
			// Generate a UUID for better uniqueness than the default
			return uuid.New().String()
		},
	})
}
