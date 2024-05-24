package app

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Host string
	Port string
	MySQLConfig
	RedisConfig
}

func NewConfig() Config {
	config := Config{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		MySQLConfig: MySQLConfig{
			Host: os.Getenv("MYSQL_HOST"),
			Port: os.Getenv("MYSQL_PORT"),
			DBName: os.Getenv("MYSQL_DBNAME"),
			Username: os.Getenv("MYSQL_USERNAME"),
			Password: os.Getenv("MYSQL_PASSWORD"),
		},
		RedisConfig: RedisConfig{
			Address: os.Getenv("REDIS_ADDRESS"),
		},
	}

	// log.Printf("config: %+v", config)

	return config
}