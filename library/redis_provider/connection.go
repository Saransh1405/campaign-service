package redis_provider

import (
	"context"
	"errors"
	"users-service/constants"
	"users-service/logger"
	"users-service/utils/configs"

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

	opt, _ := redis.ParseURL(redisUrl)

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
