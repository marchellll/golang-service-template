package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// LoggerMiddleware is a middleware that logs the request and response.
func LoggerMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			startTime := time.Now()

			logger.Debug().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Str("host", req.Host).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", req.UserAgent()).
				Str("request_id", res.Header().Get(echo.HeaderXRequestID)).
				Str("correlation_id", res.Header().Get(echo.HeaderXCorrelationID)).
				Msg("request received")


			err := next(c)

			duration := time.Since(startTime)

			logger.Debug().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Dur("latency", duration).
				Str("latency_human", duration.String()).
				Str("host", req.Host).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", req.UserAgent()).
				Str("request_id", res.Header().Get(echo.HeaderXRequestID)).
				Str("correlation_id", res.Header().Get(echo.HeaderXCorrelationID)).
				Stack().
				Err(err).
				Msg("request processed")

			return err
		}
	}
}
