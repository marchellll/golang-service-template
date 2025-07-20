package app

import (
	"golang-service-template/internal/common"

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

	_config := common.Config{
		Host:           getenv("HOST"),
		Port:           getenv("PORT"),

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
