package middleware

import (
	"context"
	"golang-service-template/internal/common"
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
)

// this middleware validates JWT tokens
// no user validation is done here
func ValidateJWTMiddleware(config common.Config, logger zerolog.Logger) echo.MiddlewareFunc {
	keyFunc := func(ctx context.Context) (interface{}, error) {
		// Our token must be signed using this data.
		return []byte("secret"), nil // TODO: replace with actual secret from config
	}

	// Set up the validator.
	jwtValidator, err := validator.New(
		keyFunc,
		validator.HS256,
		"http://example.com/", // TODO: replace with actual issuer URL from config
		[]string{"audience"},  // TODO: replace with actual audience
	)

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create jwt validator")
	}

	// new middleware that checks the JWT token
	middlewareStruct := jwtmiddleware.New(jwtValidator.ValidateToken)

	middleware := echo.WrapMiddleware(middlewareStruct.CheckJWT)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return middleware(SetUserIdMiddleware(next))(c)
		}
	}
}

const ContextKeyUserId = "context_key_user_id"

// this middleware sets the user ID in the context
// it must be used after ValidateJWTMiddleware
func SetUserIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the user ID from the JWT token.
		claims, ok := c.Request().Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "failed to get validated claims")
		}

		// Add the user ID to the echo context.
		c.Set(ContextKeyUserId, claims.RegisteredClaims.Subject)

		return next(c)
	}
}
