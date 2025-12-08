package app

import (
	"strings"

	"github.com/samber/do"

	"golang-service-template/internal/common"
	"golang-service-template/internal/errz"
	"golang-service-template/internal/handler"
	"golang-service-template/internal/middleware"
	"golang-service-template/internal/telemetry"
	"net/http"

	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// splitAndTrim splits a string by delimiter and trims whitespace from each part
func splitAndTrim(s, delimiter string) []string {
	parts := strings.Split(s, delimiter)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// securityHeadersMiddleware adds security headers to responses
func securityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			c.Response().Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			// Strict-Transport-Security should only be set when using HTTPS
			// Uncomment if your service is behind HTTPS/TLS termination
			// c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			return next(c)
		}
	}
}

// routes and middlewares here
func addRoutes(
	e *echo.Echo,
	injector *do.Injector,
) {
	logger := do.MustInvoke[zerolog.Logger](injector)
	config := do.MustInvoke[common.Config](injector)

	// global middlewares
	e.Use(echo_middleware.Recover())                // First - handle panics
	
	// Configure CORS securely
	allowedOrigins := []string{"*"} // Default to allow all for development
	if config.AllowedOrigins != "" {
		// Parse comma-separated origins
		origins := []string{}
		for _, origin := range splitAndTrim(config.AllowedOrigins, ",") {
			if origin != "" {
				origins = append(origins, origin)
			}
		}
		if len(origins) > 0 {
			allowedOrigins = origins
		}
	}
	
	// CORS configuration: if using wildcard, don't allow credentials
	allowCredentials := len(allowedOrigins) == 1 && allowedOrigins[0] != "*"
	
	e.Use(echo_middleware.CORSWithConfig(echo_middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowCredentials: allowCredentials,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		MaxAge:           86400, // 24 hours
	}))
	
	// Body size limit to prevent DoS attacks
	e.Use(echo_middleware.BodyLimit("1M"))
	
	e.Use(middleware.RequestIDMiddleware())         // Third - generate request ID for tracing
	e.Use(middleware.TelemetryMiddleware(injector)) // Fourth - track all requests (including failed ones)
	e.Use(middleware.LoggerMiddleware(logger))      // Fifth - logger should capture request ID and telemetry context
	e.Use(errz.ErrorRendererMiddleware())           // Sixth - handle error rendering
	e.Use(middleware.ValidatorMiddleware())         // Seventh - set up request validation
	
	// Security headers
	e.Use(securityHeadersMiddleware())

	// routes
	addHealthzRoutes(injector, e)
	addTaskRoutes(injector, e)
	addMetricsRoutes(injector, e)

	// root route
	e.Any("/", echo.WrapHandler(http.NotFoundHandler()))
}

func addHealthzRoutes(injector *do.Injector, e *echo.Echo) {
	healthController := do.MustInvoke[handler.HealthzController](injector)
	e.GET("/healthz", healthController.GetHealthz())
	e.POST("/healthz", healthController.GetHealthz())
	e.GET("/readyz", healthController.GetReadyz())
	e.GET("/errorz", healthController.Errorz())
}

func addTaskRoutes(injector *do.Injector, e *echo.Echo) {
	taskGroup := e.Group("/tasks")

	// taskGroup.Use(echo_middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "username" && password == "password" {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	taskGroup.POST("", do.MustInvoke[handler.TaskController](injector).Create())
	taskGroup.GET("", do.MustInvoke[handler.TaskController](injector).Find())
	taskGroup.GET("/:id", do.MustInvoke[handler.TaskController](injector).GetById())
	taskGroup.PATCH("/:id", do.MustInvoke[handler.TaskController](injector).Update())
	taskGroup.DELETE("/:id", do.MustInvoke[handler.TaskController](injector).Delete())

	securedTaskGroup := e.Group("/secured/tasks")
	securedTaskGroup.Use(middleware.ValidateJWTMiddleware(
		do.MustInvoke[common.Config](injector),
		do.MustInvoke[zerolog.Logger](injector),
	))

	securedTaskGroup.GET("", do.MustInvoke[handler.TaskController](injector).FindByUserId())

}

func addMetricsRoutes(injector *do.Injector, e *echo.Echo) {
	// Get telemetry from dependency injection (optional, may not be available)
	if tel, err := do.Invoke[*telemetry.Telemetry](injector); err == nil && tel != nil {
		e.GET("/metrics", echo.WrapHandler(tel.GetMetricsHandler()))
	}
}
