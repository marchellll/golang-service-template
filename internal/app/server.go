package app

import (
	"context"
	"golang-service-template/internal/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func RunNewServer(
	injector *do.Injector,
) func(ctx context.Context) error {
	logger := do.MustInvoke[zerolog.Logger](injector)
	config := do.MustInvoke[common.Config](injector)

	e := echo.New()

	addRoutes(e, injector)
	addSocketIoRoutes(e, injector)

	logger.Info().Str("port", config.Port).Msg("starting server")

	go func() {
		if err := e.Start(":" + config.Port); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	return func(ctx context.Context) error {
		return e.Shutdown(ctx)
	}
}
