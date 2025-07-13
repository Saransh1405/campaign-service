package redis_rate_limiter

import (
	"campaign-service/library/redis_provider"
	"campaign-service/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"go.uber.org/zap"
)

var RateLimiter *redis_rate.Limiter

func CreateInstance(log logger.Logger) error {
	RateLimiter = redis_rate.NewLimiter(redis_provider.RedisClient())
	log.Info("Created a new rate limiter instance.")

	return nil
}

// DbContext - Helper Functions
func GetInstance() *redis_rate.Limiter {
	return RateLimiter
}

func CheckLimiter(ctx *gin.Context, key string) bool {
	start := time.Now()

	// Get log with context
	log := logger.GetLogger(ctx)

	res, err := RateLimiter.Allow(ctx, key, redis_rate.PerMinute(redis_provider.RatePerMinute))
	if err != nil {
		log.With(zap.Error(err)).Error(err.Error())
		return false
	}

	// We may get allowed = 0 if it was not allowed
	log.Debug("redis : "+key, zap.Any("allowed", res.Allowed), zap.Any("remaining", res.Remaining))

	timeElapsed := time.Since(start)
	log.Debug(fmt.Sprintf("Redis rate limiter took %s", timeElapsed))

	return res.Allowed != 0
}
