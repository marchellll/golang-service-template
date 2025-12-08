package app

import (
	"golang-service-template/internal/common"
	"golang-service-template/internal/telemetry"

	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func NewTelemetry(i *do.Injector) (*telemetry.Telemetry, error) {
	config := do.MustInvoke[common.Config](i)
	logger := do.MustInvoke[zerolog.Logger](i)

	return telemetry.NewTelemetry(config.TelemetryConfig, logger)
}
