package main

import (
	"campaign-service/api"
	"campaign-service/constants"
	"campaign-service/library/kafka"
	"campaign-service/library/postgres"
	"campaign-service/library/redis_provider"
	"campaign-service/logger"
	"campaign-service/utils"
	"context"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"campaign-service/utils/configs"
	"campaign-service/utils/flags"
	"campaign-service/utils/httpclient"
	"campaign-service/utils/localization"
	loggerMiddleware "campaign-service/utils/logger"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	godotenv.Load()

	initConfigs()

	initCustomValidations()

	logger.InitLogger()

	initHTTPClient()

	postgres.InitPostgresDB(ctx)

	err := redis_provider.NewConnection(ctx, logger.GetLoggerWithoutContext())
	if err != nil {
		logger.GetLoggerWithoutContext().With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
	}

	// mongoDb.InitMongoDB()

	kafka.NewConnection()

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

	applicationConfig, err := configs.Get(constants.ApplicationConfig)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

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

	logMiddleware := loggerMiddleware.LoggerMiddlewareOptions{
		NotLogQueryParams:  true,
		NotLogHeaderParams: true,
	}

	logger := logger.GetLoggerWithoutContext()

	applicationConfig, err1 := configs.Get(constants.ApplicationConfig)
	if err1 != nil {
		logger.With(zap.Error(err1)).Error(constants.BindingFailedErrr)
	}

	router := api.GetRouter(localization.LoadBundle(""), loggerMiddleware.Logger(logMiddleware), applicationConfig)

	serverPort := applicationConfig.GetInt(constants.ServerPort)
	logger.Info(fmt.Sprintf("Running Server on port : %v", serverPort))

	err := router.Run(fmt.Sprintf(":%d", serverPort))
	if err != nil {
		logger.With(zap.Error(err)).Error(constants.ExternalServiceFailureError)
	}
}
