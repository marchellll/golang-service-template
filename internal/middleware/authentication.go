package middleware

import (
	"context"
	"golang-service-template/internal/common"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/validator"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// this middleware validates JWT tokens
// no user validation is done here
func ValidateJWTMiddleware(config common.Config, logger zerolog.Logger) echo.MiddlewareFunc {
	if config.Secret == "" {
		logger.Fatal().Msg("JWT_SECRET is required but not configured")
	}

	keyFunc := func(ctx context.Context) (any, error) {
		return []byte(config.Secret), nil
	}

	jwtValidator, err := validator.New(
		validator.WithKeyFunc(keyFunc),
		validator.WithAlgorithm(validator.HS256),
		validator.WithIssuer(config.Issuer),
		validator.WithAudiences([]string{config.Audience}),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create jwt validator")
	}

	jwtMW, err := jwtmiddleware.New(
		jwtmiddleware.WithValidator(jwtValidator),
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create jwt middleware")
	}

	middleware := echo.WrapMiddleware(jwtMW.CheckJWT)

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
		claims, err := jwtmiddleware.GetClaims[*validator.ValidatedClaims](c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "failed to get validated claims")
		}

		c.Set(ContextKeyUserId, claims.RegisteredClaims.Subject)

		return next(c)
	}
}
