package app

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConfig struct {
	Address string
}

func ConnectRedis(logger *zap.Logger, config Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Address, // TODO
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		logger.Fatal("failed to connect to redis", zap.Error(status.Err()))
	}

	return rdb
}