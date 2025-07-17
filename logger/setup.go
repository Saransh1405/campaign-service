package logger

import (
	"campaign-service/constants"
	"campaign-service/utils/configs"
)

func InitLogger() {

	// Get a log config frokm yaml
	LoggerConfig, err := configs.Get(constants.LoggerConfig)
	if err != nil {
		panic(err)
	}

	// Setup logging from a config
	SetupLogging(LoggerConfig.GetString(constants.PathKey), LoggerConfig.GetString(constants.LogLevelKey))

}
