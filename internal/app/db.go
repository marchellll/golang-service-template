package app

import (
	"fmt"
	"golang-service-template/internal/common"

	mysql_drv "github.com/go-sql-driver/mysql"
	"github.com/samber/do"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rs/zerolog"
	gormzerolog "github.com/vitaliy-art/gorm-zerolog"
)

func ConnectDB(i *do.Injector) (*gorm.DB, error) {

	logger := do.MustInvoke[zerolog.Logger](i)
	config := do.MustInvoke[common.Config](i)

	var dialector gorm.Dialector

	switch config.DbConfig.Dialect {
	case "mysql":
		dbConfig := mysql_drv.NewConfig()
		dbConfig.Addr = config.DbConfig.Host + ":" + config.DbConfig.Port
		dbConfig.DBName = config.DbConfig.DBName
		dbConfig.User = config.DbConfig.Username
		dbConfig.Passwd = config.DbConfig.Password
		dbConfig.Net = "tcp"
		// https://stackoverflow.com/questions/29341590/how-to-parse-time-from-database/29343013#29343013
		dbConfig.ParseTime = true

		dialector = mysql.Open(dbConfig.FormatDSN())
	case "postgres":
		sslmode := config.DbConfig.SslMode
		if sslmode == "" {
			sslmode = "require" // Default to secure
		}
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.DbConfig.Host, config.DbConfig.Port, config.DbConfig.Username, config.DbConfig.Password, config.DbConfig.DBName, sslmode,
		)

		dialector = postgres.Open(dsn)
	default:
		logger.Fatal().Str("dialect", config.DbConfig.Dialect).Msg("unsupported database dialect")
		return nil, fmt.Errorf("unsupported database dialect: %s", config.DbConfig.Dialect)
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormzerolog.NewGormLogger(),
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to DB")
	}

	// Add OpenTelemetry instrumentation for database tracing
	if err := gormDB.Use(otelgorm.NewPlugin()); err != nil {
		logger.Fatal().Err(err).Msg("failed to add GORM OpenTelemetry plugin")
	}

	logger.Info().Msg("Database OpenTelemetry instrumentation enabled")

	return gormDB, nil
}
