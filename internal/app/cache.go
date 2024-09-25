package app

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/samber/do"
)

type RedisConfig struct {
	Address string `validate:"required"`
}

func ConnectRedis(i *do.Injector) (*redis.Client, error) {
	logger := do.MustInvoke[zerolog.Logger](i)
	config := do.MustInvoke[Config](i)

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Address, // TODO
		Password: "",                         // no password set
		DB:       0,                          // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		logger.Fatal().Err(status.Err()).Msg("failed to connect to redis")
	}

	return rdb, nil
}
