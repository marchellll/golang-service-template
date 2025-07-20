package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	maxBodyLogSize = 10 * 1024 // 10KB limit for body logging
)

// isBodyLoggingEnabled checks if body logging is enabled via environment variable
// Set ENABLE_BODY_LOGGING=true to enable body logging (default: false)
func isBodyLoggingEnabled() bool {
	enabled := os.Getenv("ENABLE_BODY_LOGGING")
	if enabled == "" {
		return false // Default to disabled for security/performance
	}

	result, err := strconv.ParseBool(enabled)
	if err != nil {
		return false // Default to disabled if invalid value
	}

	return result
}

// bodyDumpResponseWriter wraps http.ResponseWriter to capture response body
type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// shouldLogBody determines if we should log the body based on content type
func shouldLogBody(contentType string) bool {
	// Only log text-based content types
	textTypes := []string{
		"application/json",
		"application/xml",
		"text/",
		"application/x-www-form-urlencoded",
	}

	contentType = strings.ToLower(contentType)
	for _, textType := range textTypes {
		if strings.Contains(contentType, textType) {
			return true
		}
	}
	return false
}

// truncateBody truncates body if it's too large and masks sensitive data
func truncateBody(body string, contentType string) string {
	// Don't log if content type suggests binary data
	if !shouldLogBody(contentType) {
		return "[binary content]"
	}

	if len(body) > maxBodyLogSize {
		return body[:maxBodyLogSize] + "... [truncated]"
	}

	// Mask sensitive fields in JSON
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		// Simple masking for common sensitive fields
		sensitiveFields := []string{"password", "token", "secret", "key", "credential"}
		for _, field := range sensitiveFields {
			if strings.Contains(strings.ToLower(body), `"`+field+`"`) {
				// This is a simple approach - in production you might want more sophisticated masking
				body = strings.ReplaceAll(body, `"`+field+`":"`, `"`+field+`":"[MASKED]","masked_original":"`)
			}
		}
	}

	return body
}

// LoggerMiddleware is a middleware that logs the request and response.
func LoggerMiddleware(logger zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			startTime := time.Now()

			// Read and log request body (only if enabled)
			var requestBody string
			bodyLoggingEnabled := isBodyLoggingEnabled()

			if bodyLoggingEnabled && req.Body != nil {
				bodyBytes, err := io.ReadAll(req.Body)
				if err == nil {
					contentType := req.Header.Get("Content-Type")
					requestBody = truncateBody(string(bodyBytes), contentType)
					// Restore the body for downstream handlers
					req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}

			logEvent := logger.Debug().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("host", req.Host).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", req.UserAgent()).
				Str("request_id", res.Header().Get(echo.HeaderXRequestID)).
				Str("correlation_id", res.Header().Get(echo.HeaderXCorrelationID))

			// Only add request body if logging is enabled
			if bodyLoggingEnabled {
				logEvent = logEvent.Str("request_body", requestBody)
			}

			logEvent.Msg("request received")

			// Capture response body (only if enabled)
			var resBody *bytes.Buffer
			if bodyLoggingEnabled {
				resBody = new(bytes.Buffer)
				mw := io.MultiWriter(c.Response().Writer, resBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
				c.Response().Writer = writer
			}

			err := next(c)

			duration := time.Since(startTime)

			// Process response body for logging (only if enabled)
			var responseBody string
			if bodyLoggingEnabled && resBody != nil {
				responseContentType := res.Header().Get("Content-Type")
				responseBody = truncateBody(resBody.String(), responseContentType)
			}

			// Determine log level based on HTTP status code and error presence
			lvl := zerolog.InfoLevel
			status := res.Status

			// Set error level for 4XX and 5XX status codes, or if there's an error
			if status >= 400 || err != nil {
				lvl = zerolog.ErrorLevel
			} else if status >= 300 {
				lvl = zerolog.WarnLevel
			}

			logEvent = logger.WithLevel(lvl).
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
				Err(err)

			// Only add body fields if logging is enabled
			if bodyLoggingEnabled {
				logEvent = logEvent.
					Str("request_body", requestBody).
					Str("response_body", responseBody)
			}

			logEvent.Msg("request processed")

			return err
		}
	}
}
