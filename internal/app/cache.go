package app

import (
	"context"
	"golang-service-template/internal/common"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/samber/do"
)

func ConnectRedis(i *do.Injector) (*redis.Client, error) {
	logger := do.MustInvoke[zerolog.Logger](i)
	config := do.MustInvoke[common.Config](i)

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Address, // TODO
		Password: "",             // no password set
		DB:       0,              // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		logger.Fatal().Err(status.Err()).Msg("failed to connect to redis")
	}

	return rdb, nil
}
