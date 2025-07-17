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

// Client - Redis Connection
var Client *redis.Client
var RatePerMinute int = 2 // Rate per minute Default to 2

func NewConnection(ctx context.Context, log logger.Logger) error {
	// get application configs
	applicationConfig, err := configs.Get(constants.ApplicationConfig)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	// RedisConfig := applicationConfig.Get(constants.RedisConfig)

	redisUrl := applicationConfig.GetString(constants.RedisUrlKey)
	RatePerMinute = applicationConfig.GetInt(constants.RedisRatePerMinute)
	if redisUrl == "" {
		return errors.New("configuration is missing for redis")
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

	// Ping to check if redis connection is working
	_, err2 := Client.Ping(ctx).Result()
	if err2 != nil {
		return err2
	}

	log.Info("Connected to Redis.")

	return nil
}

// RedisClient - Helper Functions
func RedisClient() *redis.Client {
	return Client
}
