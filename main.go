package main

import (
	"context"
	"fmt"
	"users-service/api"
	"users-service/constants"
	"users-service/library/kafka"
	"users-service/library/postgres"
	"users-service/library/redis_provider"
	"users-service/logger"
	"users-service/utils"

	"users-service/utils/configs"
	"users-service/utils/localization"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"time"
	"users-service/utils/flags"
	"users-service/utils/httpclient"
	loggerMiddleware "users-service/utils/logger"

	// grpcLib "users-service/util/grpc"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Load configuration
	godotenv.Load()
	initConfigs()

	// Setup custom binding validations
	initCustomValidations()

	// setup a logger
	logger.InitLogger()

	// setup http client
	initHTTPClient()

	// Connect a postgres
	postgres.InitPostgresDB(ctx)

	// Connect a kafka
	kafka.NewConnection()

	// Connect a redis
	err := redis_provider.NewConnection(ctx, logger.GetLoggerWithoutContext())
	if err != nil {
		logger.GetLoggerWithoutContext().With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
	}

	// Start router and Use middleware
	startRouter(ctx)
}

func initConfigs() {
	// init configs
	configs.Init(flags.BaseConfigPath())
}

func initCustomValidations() {
	// init custom validations
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(constants.EnumKey, utils.ValidateEnum)
	}
}

func initHTTPClient() {

	log := logger.GetLoggerWithoutContext()
	// get application configs
	applicationConfig, err := configs.Get(constants.ApplicationConfig)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	// init http client
	err = httpclient.Init(httpclient.Config{
		ConnectTimeout: time.Millisecond *
			applicationConfig.GetDuration(constants.HTTPConnectTimeoutInMillisKey),
		KeepAliveDuration: time.Millisecond *
			applicationConfig.GetDuration(constants.HTTPKeepAliveDurationInMillisKey),
		MaxIdleConnections: applicationConfig.GetInt(constants.HTTPMaxIdleConnectionsKey),
		IdleConnectionTimeout: time.Millisecond *
			applicationConfig.GetDuration(constants.HTTPIdleConnectionTimeoutInMillisKey),
		TLSHandshakeTimeout: time.Millisecond *
			applicationConfig.GetDuration(constants.HTTPTlsHandshakeTimeoutInMillisKey),
		ExpectContinueTimeout: time.Millisecond *
			applicationConfig.GetDuration(constants.HTTPExpectContinueTimeoutInMillisKey),
		Timeout: time.Millisecond *
			applicationConfig.GetDuration(constants.HTTPTimeoutInMillisKey),
	})
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}
}

func startRouter(ctx context.Context) {

	// set a middleware options for a logger
	logMiddleware := loggerMiddleware.LoggerMiddlewareOptions{
		NotLogQueryParams:  true,
		NotLogHeaderParams: true,
	}

	// get logger
	logger := logger.GetLoggerWithoutContext()

	// get application config
	applicationConfig, err1 := configs.Get(constants.ApplicationConfig)
	if err1 != nil {
		logger.With(zap.Error(err1)).Error(constants.BindingFailedErrr)
	}

	// get router
	router := api.GetRouter(localization.LoadBundle(""), loggerMiddleware.Logger(logMiddleware), applicationConfig)

	// now start router
	serverPort := applicationConfig.GetInt(constants.ServerPort)
	logger.Info(fmt.Sprintf("Running Server on port : %v", serverPort))

	err := router.Run(fmt.Sprintf(":%d", serverPort))
	if err != nil {
		logger.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
	}
}
