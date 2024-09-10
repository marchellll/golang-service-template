package app

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host string
	Port string `validate:"required"`
	DbConfig `validate:"required"`
	RedisConfig `validate:"required"`
}

func NewConfig(getenv func(string) string) Config {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)


	_config := Config{
		Host: getenv("HOST"),
		Port: getenv("PORT"),
		DbConfig: DbConfig{
			Dialect: getenv("DB_DIALECT"),
			Host: getenv("DB_HOST"),
			Port: getenv("DB_PORT"),
			DBName: getenv("DB_DBNAME"),
			Username: getenv("DB_USERNAME"),
			Password: getenv("DB_PASSWORD"),
		},
		RedisConfig: RedisConfig{
			Address: getenv("REDIS_ADDRESS"),
		},
	}

	err := validate.Struct(_config)

	if err == nil {
		return _config
	}

	errs := err.(validator.ValidationErrors)
	log.Panic().Err(err).Any("missing config", errs.Translate(trans)).Msg("failed to validate config: ")

	// log.Printf("config: %+v", config)

	return _config
}