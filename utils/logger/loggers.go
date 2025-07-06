package logger

import (
	"time"
	"users-service/constants"
	"users-service/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggerMiddlewareOptions is the set of configurable allowed for log
type LoggerMiddlewareOptions struct {
	NotLogQueryParams  bool
	NotLogHeaderParams bool
}

// Logger is the middleware to be used for logging the request and response information
// This should be the first middleware to be added, in case the recovery middleware is not being used.
// Otherwise, it should be the second one, just after the recovery middleware.
func Logger(options LoggerMiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetHeader(constants.RequestIDHeader)

		if id == "" {
			// get a unique id
			uid, err := uuid.NewUUID()
			if err == nil {
				id = uid.String()
			}
		}

		// apply the id in the context
		ctx.Set(constants.RequestId, id)

		// Start timer
		start := time.Now()
		log := logger.GetLogger(ctx)
		// log the initial request
		log.With(zap.Time(constants.StartTimeLogParam, start)).Info("startRequest")

		// Process request
		ctx.Next()

		log.With(
			zap.String(constants.RequestId, id),
			zap.Int(constants.StatusCode, ctx.Writer.Status()),
			zap.Int64("Latency", time.Since(start).Milliseconds()),
			zap.String(constants.ClientIPLogParam, ctx.ClientIP()),
			zap.String(constants.ErrorLogParam, ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
		).Info("summary")
	}
}
