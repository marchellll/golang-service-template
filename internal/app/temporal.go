package app

import (
	"context"
	"golang-service-template/internal/common"
	"golang-service-template/internal/telemetry"

	"github.com/rs/zerolog"
	"github.com/samber/do"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
)

func NewTemporalClient(i *do.Injector) (client.Client, error) {
	logger := do.MustInvoke[zerolog.Logger](i)
	config := do.MustInvoke[common.Config](i)

	// If Temporal config is not provided, return nil
	if config.TemporalConfig.Address == "" {
		logger.Debug().Msg("Temporal address not configured, skipping Temporal client")
		return nil, nil
	}

	namespace := config.TemporalConfig.Namespace
	if namespace == "" {
		namespace = "default"
	}

	// Build client options
	clientOptions := client.Options{
		HostPort:  config.TemporalConfig.Address,
		Namespace: namespace,
	}

	// Add OpenTelemetry tracing interceptor if telemetry is enabled
	clientOptions.Interceptors = createTracingInterceptors(i, config, logger)

	c, err := client.Dial(clientOptions)

	if err != nil {
		logger.Error().Err(err).Msg("failed to create Temporal client")
		return nil, err
	}

	logger.Info().
		Str("address", config.TemporalConfig.Address).
		Str("namespace", namespace).
		Msg("Temporal client connected")

	return c, nil
}

// createTracingInterceptors creates OpenTelemetry tracing interceptors for Temporal
// if telemetry is enabled and configured. Returns an empty slice if telemetry is disabled or unavailable.
func createTracingInterceptors(i *do.Injector, config common.Config, logger zerolog.Logger) []interceptor.ClientInterceptor {
	telemetryInstance, err := do.Invoke[*telemetry.Telemetry](i)
	if err != nil || telemetryInstance == nil {
		return nil
	}

	if !config.TelemetryConfig.Enabled || !config.TelemetryConfig.TracingEnabled {
		return nil
	}

	tracerProvider := otel.GetTracerProvider()
	if tracerProvider == nil {
		return nil
	}

	propagator := otel.GetTextMapPropagator()
	if propagator == nil {
		propagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}

	tracerOptions := opentelemetry.TracerOptions{
		Tracer:               tracerProvider.Tracer("temporal"),
		TextMapPropagator:    propagator,
		DisableSignalTracing: false,
		DisableQueryTracing:  false,
	}

	tracingInterceptor, err := opentelemetry.NewTracingInterceptor(tracerOptions)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create Temporal tracing interceptor, continuing without tracing")
		return nil
	}

	logger.Info().Msg("Temporal OpenTelemetry tracing enabled")
	return []interceptor.ClientInterceptor{tracingInterceptor}
}

func ShutdownTemporalClient(ctx context.Context, c client.Client) error {
	if c == nil {
		return nil
	}
	c.Close()
	return nil
}
