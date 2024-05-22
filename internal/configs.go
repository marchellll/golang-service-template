package internal

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Host string `mapstructure:"HOST"`
	Port string `mapstructure:"PORT"`
}

func NewConfig() Config {
	config := Config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}

	log.Printf("config: %+v", config)

	return config
}