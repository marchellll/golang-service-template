package common

type Config struct {
	ServiceName               string `validate:"required"`
	Host                      string
	Port                      string `validate:"required"`
	HealthcheckTimeoutSeconds int    `validate:"min=1"`
	DbConfig                  `validate:"required"`
	RedisConfig               `validate:"required"`
	TelemetryConfig           `validate:"required"`
	JWTConfig                 `validate:"required"`
	AllowedOrigins            string `validate:""` // Comma-separated list of allowed CORS origins
}

type TelemetryConfig struct {
	Enabled        bool   `validate:""`
	OtelEndpoint   string `validate:""`
	ServiceName    string `validate:""`
	ServiceVersion string `validate:""`
	Environment    string `validate:""`
	MetricsEnabled bool   `validate:""`
	TracingEnabled bool   `validate:""`
}

type RedisConfig struct {
	Address string `validate:"required"`
}

type DbConfig struct {
	Dialect string `validate:"required"`
	Host    string `validate:"required"`
	Port    string `validate:"required"`
	DBName  string `validate:"required"`

	Username string `validate:"required"`
	Password string `validate:"required"`
	SslMode  string `validate:""` // Optional, defaults to "require" if not set
}

type JWTConfig struct {
	Secret   string `validate:"required"`
	Issuer   string `validate:"required"`
	Audience string `validate:"required"`
}
