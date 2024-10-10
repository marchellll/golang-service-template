package common

type RedisConfig struct {
	Address string `validate:"required"`
}

type Config struct {
	Host        string
	Port        string `validate:"required"`
	DbConfig    `validate:"required"`
	RedisConfig `validate:"required"`
}

type DbConfig struct {
	Dialect string `validate:"required"`
	Host    string `validate:"required"`
	Port    string `validate:"required"`
	DBName  string `validate:"required"`

	Username string `validate:"required"`
	Password string `validate:"required"`
}