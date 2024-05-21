package internal

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

func NewConfig() Config {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %+v", config)

	return config
}