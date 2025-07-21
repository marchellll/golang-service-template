# 🔍 Telemetry & Observability
| **Grafana** | http://localhost:3001 | Dashboards & visualization | admin/admin123 |This servi- Go to http://localhost:3001 (admin/admin123)e includes a comprehensive observability stack using **OpenTelemetry** for monitoring metrics - **URL**: http://localhost:3001nd distributed tracing. The system integrates **Jaeger**, **Prometheus**, and **Grafana   - **Dashboards**: http://localhost:3001 (Grafana)* using Docker Compose profiles for a complete observability solution.

## 🚀 Quick Start

### Option 1: Basic Development (No Observability)
```bash
docker-compose up
```

### Option 2: With Full Observability Stack
```bash
docker-compose --profile observability up
```

### Option 3: Everything (Including Future Services)
```bash
docker-compose --profile full up
```

## 📊 Access Points

When running with observability profile:

| Service | URL | Purpose | Credentials |
|---------|-----|---------|-------------|
| **Your API** | http://localhost:8080 | Main application | - |
| **Metrics Endpoint** | http://localhost:8080/metrics | Raw Prometheus metrics | - |
| **Prometheus** | http://localhost:9090 | Metrics collection & queries | - |
| **Jaeger** | http://localhost:16686 | Distributed tracing | - |
| **Grafana** | http://localhost:33000 | Dashboards & visualization | admin/admin123 |

## Configuration

Telemetry is configured through environment variables:

```bash
# Enable/disable telemetry
TELEMETRY_ENABLED=true

# Enable/disable metrics collection (Prometheus)
TELEMETRY_METRICS_ENABLED=true

# Enable/disable distributed tracing (Jaeger)
TELEMETRY_TRACING_ENABLED=true

# OpenTelemetry endpoint for Jaeger
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318/v1/traces

# Service information
SERVICE_NAME=golang-service-template
SERVICE_VERSION=1.0.0
ENVIRONMENT=development
```

## Architecture

### OpenTelemetry Integration
- **Metrics**: Prometheus-compatible metrics via OpenTelemetry SDK
- **Tracing**: Jaeger-compatible traces via OTLP HTTP exporter
- **Propagation**: Full distributed tracing context propagation
- **Auto-instrumentation**: HTTP requests automatically traced and measured

### Docker Compose Profiles
The observability stack can be started optionally:

```bash
# Development only
docker-compose up

# With observability stack
docker-compose --profile observability up
# or
docker-compose --profile full up
```

## 📈 What You Get

### **Metrics (Prometheus)**
- **HTTP Request Metrics**: Count, duration, status codes by method, path, and status
- **Business Metrics**: Task operations, success/error rates
- **Custom Metrics**: Dynamic creation of any counter or histogram metrics
- **System Metrics**: Go runtime metrics via OpenTelemetry

### **Tracing (Jaeger)**
- **Distributed Traces**: See request flow across services
- **Span Details**: HTTP requests, database operations, business logic
- **Error Tracking**: Failed operations with full context
- **Context Propagation**: Automatic trace correlation across service boundaries

### **Dashboards (Grafana)**
- **Pre-configured Datasources**: Prometheus + Jaeger
- **Custom Dashboards**: Ready for your metrics
- **Alerting**: Set up alerts on your metrics
- **Correlation**: Link metrics and traces together

## 🔧 Development Workflow

### 1. **Start with Observability**
```bash
docker-compose --profile observability up -d
```

### 2. **Generate Some Traffic**
```bash
# Make some requests
curl http://localhost:8080/healthz
curl http://localhost:8080/tasks
curl -X POST http://localhost:8080/tasks -d '{"description":"test"}'
```

### 3. **Explore Your Data**

**Metrics in Prometheus:**
- Go to http://localhost:9090
- Try queries like:
  - `http_requests_total`
  - `rate(http_requests_total[5m])`
  - `task_create_total`

**Traces in Jaeger:**
- Go to http://localhost:16686
- Select "golang-service-template" service
- Click "Find Traces"

**Dashboards in Grafana:**
- Go to http://localhost:33000 (admin/admin123)
- Create dashboards using Prometheus data
- Add Jaeger traces for correlation

## Features

### 1. HTTP Request Tracking
The telemetry middleware automatically tracks:
- Request count by method, path, and status (`http_requests_total`)
- Request duration by method and path (`http_request_duration_seconds`)
- Detailed trace spans with attributes (user agent, remote address, etc.)

### 2. Generic Service Metrics
Services can record any metrics using flexible methods:
- **Counters**: `Increment(ctx, metricName, attrs...)`
- **Durations**: `RecordDuration(ctx, metricName, startTime, attrs...)`
- **Traces**: `CreateSpan(ctx, spanName, attrs...)`
- **Errors**: `RecordError(ctx, err)`

### 3. Automatic Metric Creation
- Metrics are created dynamically when first used
- No need to pre-define metrics
- Consistent naming patterns encouraged

## Usage Examples

### Basic Service Implementation
```go
func (s *service) CreateUser(ctx context.Context, user User) (*User, error) {
    start := time.Now()

    // Create span with descriptive name
    ctx, span := s.telemetry.CreateSpan(ctx, "user_create",
        attribute.String("operation", "create"),
        attribute.String("user.type", "admin"))
    defer span.End()

    // Business logic...
    result, err := s.doBusinessLogic(ctx, user)

    if err != nil {
        // Record error metrics
        s.telemetry.Increment(ctx, "user_create_total",
            attribute.String("status", "error"))
        s.telemetry.RecordDuration(ctx, "user_create_duration_seconds",
            start,
            attribute.String("status", "error"))
        s.telemetry.RecordError(ctx, err)
        return nil, err
    }

    // Record success metrics
    s.telemetry.Increment(ctx, "user_create_total",
        attribute.String("status", "success"))
    s.telemetry.RecordDuration(ctx, "user_create_duration_seconds",
        start,
        attribute.String("status", "success"))

    return result, nil
}
```

### Automatic HTTP Tracking
The telemetry middleware is automatically applied to all HTTP routes and tracks:
- Request counts and response times
- Distributed trace context propagation
- Error rates and status codes

## Metric Naming Conventions

Follow these patterns for consistent metric names:

### Counters (use with `Increment`)
- `{service}_{operation}_total` - Count of operations
- `{service}_{operation}_errors_total` - Count of errors

Examples:
```go
s.telemetry.Increment(ctx, "task_create_total",
    attribute.String("status", "success"))
s.telemetry.Increment(ctx, "user_login_total",
    attribute.String("method", "oauth"))
```

### Histograms (use with `RecordDuration`)
- `{service}_{operation}_duration_seconds` - Operation duration

Examples:
```go
s.telemetry.RecordDuration(ctx, "task_create_duration_seconds",
    start, attribute.String("status", "success"))
s.telemetry.RecordDuration(ctx, "database_query_duration_seconds",
    start, attribute.String("table", "users"))
```

### Span Names (use with `CreateSpan`)
- `{service}_{operation}` - e.g., `task_create`, `user_get`, `order_process`

Examples:
```go
ctx, span := s.telemetry.CreateSpan(ctx, "task_get",
    attribute.String("task.id", id))
ctx, span := s.telemetry.CreateSpan(ctx, "payment_process",
    attribute.String("amount", "100.00"))
```

## 🎯 Example Queries

### Prometheus Queries
```promql
# Request rate per second
rate(http_requests_total[5m])

# Request duration 95th percentile
histogram_quantile(0.95, http_request_duration_seconds_bucket)

# Error rate
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# Task operation success rate
rate(task_create_total{status="success"}[5m]) / rate(task_create_total[5m])

# Business metrics examples
sum(rate(task_create_total[5m])) by (status)
histogram_quantile(0.99, task_create_duration_seconds_bucket)
```

## Observability StackWhen using the observability profile, the following services are available:

### Jaeger (Distributed Tracing)
- **URL**: http://localhost:16686
- **Purpose**: View distributed traces across services
- **Features**: Trace search, service map, performance analysis

### Prometheus (Metrics Collection)
- **URL**: http://localhost:9090
- **Purpose**: Metrics storage and querying
- **Endpoint**: `GET /metrics` on your service

### Grafana (Visualization)
- **URL**: http://localhost:33000
- **Credentials**: admin/admin123
- **Purpose**: Dashboards and alerting
- **Pre-configured**: Prometheus datasource included

## 📁 File Structure
```
observability/
├── prometheus.yml              # Prometheus configuration
└── grafana/
    ├── provisioning/
    │   ├── datasources/
    │   │   └── datasources.yml # Auto-configure Prometheus & Jaeger
    │   └── dashboards/
    │       └── dashboards.yml  # Dashboard provider config
    └── dashboards/             # Put your .json dashboards here
```

## Generic Telemetry Methods

The telemetry package provides flexible, reusable methods:

### 1. `CreateSpan(ctx, name, attrs...)`
Creates a trace span with any name you want.

### 2. `Increment(ctx, metricName, attrs...)`
Increments a counter metric (automatically creates if doesn't exist).

### 3. `RecordDuration(ctx, metricName, startTime, attrs...)`
Records a duration in a histogram metric. Duration is calculated automatically from start time.

### 4. `RecordError(ctx, err)`
Records an error in the current span.

### 5. `GetMetricsHandler()`
Returns the Prometheus metrics HTTP handler for `/metrics` endpoint.

## Migration from Old Methods

If you're upgrading from business-specific methods:

| Old Method | New Method |
|------------|------------|
| `CreateTaskSpan(ctx, "create", id)` | `CreateSpan(ctx, "task_create", attribute.String("task.id", id))` |
| `RecordTaskOperation(ctx, "create", "success", dur)` | `Increment(ctx, "task_create_total", attribute.String("status", "success"))` + `RecordDuration(ctx, "task_create_duration_seconds", start, ...)` |

## Extension Points

The current implementation provides a production-ready foundation that can be extended with:

1. **Custom Exporters**: Add additional OTEL exporters for other systems
2. **Custom Metrics**: Application-specific business metrics
3. **Alerting**: Connect Grafana to notification systems
4. **Service Mesh Integration**: Istio/Envoy integration for network-level observability
5. **Log Correlation**: Connect structured logs with trace IDs

## 🔄 Production Setup

For production, consider updating:

1. **Security**: Enable HTTPS and authentication
2. **Storage**: Configure persistent volumes for data retention
3. **Retention**: Set appropriate data retention policies
4. **Scaling**: Use external Prometheus/Jaeger instances
5. **Alerting**: Configure Alertmanager for notifications
6. **High Availability**: Multi-instance deployments

## Benefits

- **Production-Ready**: Full OpenTelemetry implementation with real exporters
- **Zero-Impact When Disabled**: No performance overhead when telemetry is turned off
- **Automatic Instrumentation**: HTTP requests tracked automatically
- **Flexible**: Generic methods work for any business domain
- **Nil-Safe**: All methods handle disabled/missing telemetry gracefully
- **Standards-Compliant**: Uses OpenTelemetry standards for interoperability
- **Docker-Integrated**: Complete observability stack via Docker Compose
- **Development-Friendly**: Easy to enable/disable during development
- **Real Observability**: Industry standard tools (Jaeger, Prometheus, Grafana)
- **Easy Development**: One command to get full observability stack
- **Extensible**: Easy to add more services and metrics
- **Correlation**: Link metrics and traces together for complete visibility

## 📚 Learn More

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/)
- [Prometheus Query Language](https://prometheus.io/docs/prometheus/latest/querying/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [Grafana Dashboards](https://grafana.com/docs/grafana/latest/dashboards/)

## Quick Start

1. **Enable telemetry** in your environment:
   ```bash
   export TELEMETRY_ENABLED=true
   export TELEMETRY_METRICS_ENABLED=true
   export TELEMETRY_TRACING_ENABLED=true
   ```

2. **Start with observability stack**:
   ```bash
   docker-compose --profile full up
   ```

3. **Add telemetry to your service methods**:
   ```go
   start := time.Now()
   ctx, span := s.telemetry.CreateSpan(ctx, "my_operation")
   defer span.End()

   // Your business logic here...

   s.telemetry.Increment(ctx, "my_operation_total",
       attribute.String("status", "success"))
   s.telemetry.RecordDuration(ctx, "my_operation_duration_seconds", start)
   ```

4. **View results**:
   - **Traces**: http://localhost:16686 (Jaeger)
   - **Metrics**: http://localhost:9090 (Prometheus)
   - **Dashboards**: http://localhost:33000 (Grafana)

The telemetry system gives you complete observability into your service with minimal code changes!
