package engine

import (
	"context"
	"log/slog"

	"github.com/AliAlbhrani/StudentsArchive/env"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

var _ = func() bool {
	RDB = redis.NewClient(&redis.Options{
		Addr:     env.REDIS_HOST,
		Password: env.REDIS_PASSWORD,
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		return false
	}
	slog.Info("connected to redis")
	return true
}()
