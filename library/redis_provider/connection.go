package redis_provider

import (
	"campaign-service/constants"
	"campaign-service/logger"
	"campaign-service/utils/configs"
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var Client *redis.Client
var RatePerMinute int = 2

func NewConnection(ctx context.Context, log logger.Logger) error {
	redisConfig, err := configs.Get(constants.RedisConfig)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	redisUrl := redisConfig.GetString(constants.RedisUrlKey)
	if redisUrl == "" {
		redisUrl = "redis:6379"
	}
	RatePerMinute = redisConfig.GetInt(constants.RedisRatePerMinute)
	if RatePerMinute == 0 {
		RatePerMinute = 2
	}

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Error("Failed to parse Redis URL", zap.Error(err), zap.String("redisUrl", redisUrl))
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	if opt == nil {
		log.Error("Redis options are nil after parsing URL", zap.String("redisUrl", redisUrl))
		return errors.New("Redis options are nil after parsing URL")
	}

	Client = redis.NewClient(opt)

	_, err2 := Client.Ping(ctx).Result()
	if err2 != nil {
		return err2
	}

	log.Info("Connected to Redis.")

	return nil
}

func RedisClient() *redis.Client {
	return Client
}
