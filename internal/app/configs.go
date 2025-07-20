package app

import (
	"golang-service-template/internal/common"
	"strconv"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func NewConfig(getenv func(string) string) common.Config {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en") // en should be get from header
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	// Parse HealthcheckTimeoutSeconds with default value
	healthcheckTimeout := 55 // default to 55 seconds
	if timeoutStr := getenv("HEALTHCHECK_TIMEOUT_SECONDS"); timeoutStr != "" {
		if parsed, err := strconv.Atoi(timeoutStr); err == nil {
			healthcheckTimeout = parsed
		} else {
			log.Warn().Str("value", timeoutStr).Msg("invalid HEALTHCHECK_TIMEOUT_SECONDS, using default 55")
		}
	}

	_config := common.Config{
		ServiceName:               getenv("SERVICE_NAME"),
		Host:                      getenv("HOST"),
		Port:                      getenv("PORT"),
		HealthcheckTimeoutSeconds: healthcheckTimeout,

		DbConfig: common.DbConfig{
			Dialect:  getenv("DB_DIALECT"),
			Host:     getenv("DB_HOST"),
			Port:     getenv("DB_PORT"),
			DBName:   getenv("DB_DBNAME"),
			Username: getenv("DB_USERNAME"),
			Password: getenv("DB_PASSWORD"),
		},
		RedisConfig: common.RedisConfig{
			Address: getenv("REDIS_ADDRESS"),
		},
	}

	err := validate.Struct(_config)

	if err == nil {
		log.Trace().Any("config", _config.DbConfig).Msg("config validated")
		return _config
	}

	errs := err.(validator.ValidationErrors)
	log.Panic().Err(err).Any("missing config", errs.Translate(trans)).Msg("failed to validate config")

	// This line is never reached due to log.Panic() above, but Go requires a return statement
	// for the function signature. In practice, the program will terminate at log.Panic().
	return common.Config{}
}
