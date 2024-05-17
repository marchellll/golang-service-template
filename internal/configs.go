package internal

type Config struct {
	Host string
	Port string
}

func NewConfig() Config {
	// TODO: use dotenv

	return Config{
		Host: "localhost",
		Port: "8080",
	}
}