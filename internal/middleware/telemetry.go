package middleware

import (
	"strconv"
	"time"

	"golang-service-template/internal/telemetry"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"go.opentelemetry.io/otel/attribute"
)

// TelemetryMiddleware adds telemetry tracking to HTTP requests
func TelemetryMiddleware(injector *do.Injector) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get telemetry from dependency injection
			tel, err := do.Invoke[*telemetry.Telemetry](injector)
			if err != nil {
				// If telemetry is not available, just continue without it
				return next(c)
			}

			start := time.Now()

			// Create span for the HTTP request (automatically handles enabled/disabled)
			ctx, span := tel.CreateSpan(c.Request().Context(), "http_request",
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.path", c.Path()),
				attribute.String("http.user_agent", c.Request().UserAgent()),
				attribute.String("http.remote_addr", c.RealIP()),
			)
			defer span.End()

			// Update the request context
			c.SetRequest(c.Request().WithContext(ctx))

			// Process request
			err = next(c)

			// Collect request information
			method := c.Request().Method
			path := c.Path()
			status := strconv.Itoa(c.Response().Status)

			// Record metrics (automatically handles enabled/disabled)
			tel.RecordHTTPRequest(ctx, method, path, status, start)

			// Record error if any
			if err != nil {
				tel.RecordError(ctx, err)
				span.SetAttributes(attribute.String("http.status_code", "500"))
			} else {
				span.SetAttributes(attribute.String("http.status_code", status))
			}

			return err
		}
	}
}
