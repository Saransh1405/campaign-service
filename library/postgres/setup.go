package postgres

import (
	"context"
	"users-service/constants"
	"users-service/logger"
	"users-service/models"
	"users-service/utils/configs"

	"go.uber.org/zap"
)

func InitPostgresDB(ctx context.Context) {

	var log logger.Logger

	PostgresConfig, err := configs.Get(constants.PostgresConfig)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	var postgresConfig models.PostgresConfig

	postgresConfig.Host = PostgresConfig.GetString(constants.PostgresHostKey)
	postgresConfig.Port = PostgresConfig.GetString(constants.PostgresPortKey)
	postgresConfig.User = PostgresConfig.GetString(constants.PostgresUserKey)
	postgresConfig.Password = PostgresConfig.GetString(constants.PostgresPasswordKey)
	postgresConfig.DBName = PostgresConfig.GetString(constants.PostgresDBNameKey)
	postgresConfig.SSLMode = PostgresConfig.GetString(constants.PostgresSSLModeKey)
	postgresConfig.TimeZone = PostgresConfig.GetString(constants.PostgresTimeZoneKey)

	// Connection for postgres
	ConnectDatabase(ctx, postgresConfig)

}
