package telemetry

import (
	"context"
	"fmt"
	"golang-service-template/internal/common"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	config common.TelemetryConfig
	logger zerolog.Logger

	// OpenTelemetry components
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	tracer         trace.Tracer
	meter          metric.Meter

	// Prometheus components
	promRegistry *prometheus.Registry
	promHandler  http.Handler

	// Metrics - using generic approach
	counters   map[string]metric.Int64Counter
	histograms map[string]metric.Float64Histogram

	// Shutdown function
	shutdownFuncs []func(context.Context) error
}

// NewTelemetry initializes telemetry with OpenTelemetry and Prometheus
func NewTelemetry(config common.TelemetryConfig, logger zerolog.Logger) (*Telemetry, error) {
	logger.Info().
		Bool("enabled", config.Enabled).
		Str("endpoint", config.OtelEndpoint).
		Bool("metrics", config.MetricsEnabled).
		Bool("tracing", config.TracingEnabled).
		Msg("Initializing OpenTelemetry")

	t := &Telemetry{
		config:     config,
		logger:     logger,
		counters:   make(map[string]metric.Int64Counter),
		histograms: make(map[string]metric.Float64Histogram),
	}

	if !config.Enabled {
		logger.Info().Msg("Telemetry is disabled")
		return t, nil
	}

	// Create resource
	res, err := t.createResource()
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup tracing
	if config.TracingEnabled {
		if err := t.setupTracing(res); err != nil {
			return nil, fmt.Errorf("failed to setup tracing: %w", err)
		}
	}

	// Setup metrics
	if config.MetricsEnabled {
		if err := t.setupMetrics(res); err != nil {
			return nil, fmt.Errorf("failed to setup metrics: %w", err)
		}
	}

	// Initialize metrics
	if err := t.initMetrics(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Set up propagators
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info().Msg("OpenTelemetry initialized successfully")
	return t, nil
}

func (t *Telemetry) createResource() (*resource.Resource, error) {
	return resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(t.config.ServiceName),
			semconv.ServiceVersionKey.String(t.config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(t.config.Environment),
		),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithProcess(),
	)
}

func (t *Telemetry) setupTracing(res *resource.Resource) error {
	var exporter sdktrace.SpanExporter
	var err error

	if t.config.OtelEndpoint != "" {
		// Use OTLP exporter for Jaeger
		exporter, err = otlptracehttp.New(context.Background(),
			otlptracehttp.WithEndpoint(t.config.OtelEndpoint),
			otlptracehttp.WithInsecure(), // Use HTTPS in production
		)
		if err != nil {
			return err
		}
		t.logger.Info().Str("endpoint", t.config.OtelEndpoint).Msg("Using OTLP trace exporter")
	} else {
		// Use console exporter for development
		t.logger.Info().Msg("Using console trace exporter (no endpoint configured)")
		// In development, we'll just log traces
	}

	if exporter != nil {
		t.tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
		)

		otel.SetTracerProvider(t.tracerProvider)
		t.shutdownFuncs = append(t.shutdownFuncs, t.tracerProvider.Shutdown)
	}

	t.tracer = otel.Tracer("golang-service-template",
		trace.WithInstrumentationVersion("1.0.0"),
		trace.WithInstrumentationAttributes(
			attribute.String("service.name", t.config.ServiceName),
		),
	)

	return nil
}

func (t *Telemetry) setupMetrics(res *resource.Resource) error {
	// Create Prometheus registry
	t.promRegistry = prometheus.NewRegistry()

	// Create Prometheus exporter
	promExporter, err := otelprometheus.New(
		otelprometheus.WithRegisterer(t.promRegistry),
		otelprometheus.WithoutUnits(),
	)
	if err != nil {
		return err
	}

	// Create meter provider
	t.meterProvider = sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(promExporter),
	)

	otel.SetMeterProvider(t.meterProvider)
	t.shutdownFuncs = append(t.shutdownFuncs, t.meterProvider.Shutdown)

	t.meter = otel.Meter("golang-service-template",
		metric.WithInstrumentationVersion("1.0.0"),
		metric.WithInstrumentationAttributes(
			attribute.String("service.name", t.config.ServiceName),
		),
	)

	// Create Prometheus handler
	t.promHandler = promhttp.HandlerFor(t.promRegistry, promhttp.HandlerOpts{})

	t.logger.Info().Msg("Prometheus metrics exporter configured")
	return nil
}

func (t *Telemetry) initMetrics() error {
	if !t.config.MetricsEnabled || t.meter == nil {
		return nil
	}

	// Pre-initialize common HTTP metrics only
	commonCounters := []string{
		"http_requests_total",
	}

	commonHistograms := []string{
		"http_request_duration_seconds",
	}

	// Initialize common counters
	for _, name := range commonCounters {
		if _, err := t.getOrCreateCounter(name); err != nil {
			return err
		}
	}

	// Initialize common histograms
	for _, name := range commonHistograms {
		if _, err := t.getOrCreateHistogram(name); err != nil {
			return err
		}
	}

	return nil
}

// getOrCreateCounter gets or creates a counter metric
func (t *Telemetry) getOrCreateCounter(name string) (metric.Int64Counter, error) {
	if counter, exists := t.counters[name]; exists {
		return counter, nil
	}

	counter, err := t.meter.Int64Counter(name, metric.WithDescription("Counter metric: "+name))
	if err != nil {
		return nil, err
	}

	t.counters[name] = counter
	return counter, nil
}

// getOrCreateHistogram gets or creates a histogram metric
func (t *Telemetry) getOrCreateHistogram(name string) (metric.Float64Histogram, error) {
	if histogram, exists := t.histograms[name]; exists {
		return histogram, nil
	}

	histogram, err := t.meter.Float64Histogram(name, metric.WithDescription("Histogram metric: "+name))
	if err != nil {
		return nil, err
	}

	t.histograms[name] = histogram
	return histogram, nil
}

// GetMetricsHandler returns the Prometheus metrics handler
func (t *Telemetry) GetMetricsHandler() http.Handler {
	if !t.config.MetricsEnabled || t.promHandler == nil {
		return http.NotFoundHandler()
	}
	return t.promHandler
}

// Increment increments a counter metric
func (t *Telemetry) Increment(ctx context.Context, metricName string, attrs ...attribute.KeyValue) {
	if !t.config.Enabled || !t.config.MetricsEnabled {
		return
	}

	counter, err := t.getOrCreateCounter(metricName)
	if err != nil {
		t.logger.Error().Err(err).Str("metric", metricName).Msg("Failed to get counter")
		return
	}

	counter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordDuration records a duration in a histogram metric
func (t *Telemetry) RecordDuration(ctx context.Context, metricName string, startTime time.Time, attrs ...attribute.KeyValue) {
	if !t.config.Enabled || !t.config.MetricsEnabled {
		return
	}

	duration := time.Since(startTime).Seconds()
	histogram, err := t.getOrCreateHistogram(metricName)
	if err != nil {
		t.logger.Error().Err(err).Str("metric", metricName).Msg("Failed to get histogram")
		return
	}

	histogram.Record(ctx, duration, metric.WithAttributes(attrs...))
}

// RecordHTTPRequest records HTTP request metrics using generic methods
func (t *Telemetry) RecordHTTPRequest(ctx context.Context, method, path, status string, startTime time.Time) {
	if !t.config.Enabled || !t.config.MetricsEnabled {
		return
	}

	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.String("status", status),
	}

	// Use generic methods
	t.Increment(ctx, "http_requests_total", attrs...)
	t.RecordDuration(ctx, "http_request_duration_seconds", startTime,
		attribute.String("method", method),
		attribute.String("path", path),
	)
}

// CreateSpan creates a new trace span with automatic checks
func (t *Telemetry) CreateSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	if !t.config.Enabled || !t.config.TracingEnabled || t.tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}
	return t.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// RecordError records an error in the current span
func (t *Telemetry) RecordError(ctx context.Context, err error) {
	if t == nil || !t.config.Enabled || !t.config.TracingEnabled || err == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// Shutdown gracefully shuts down telemetry
func (t *Telemetry) Shutdown(ctx context.Context) error {
	t.logger.Info().Msg("Shutting down telemetry")

	for _, shutdown := range t.shutdownFuncs {
		if err := shutdown(ctx); err != nil {
			t.logger.Error().Err(err).Msg("Failed to shutdown telemetry component")
		}
	}

	return nil
}
