package mongoDb

import (
	"campaign-service/constants"
	"campaign-service/logger"
	"campaign-service/utils/configs"
	"errors"

	"go.uber.org/zap"
)

func InitMongoDB() {
	log := logger.GetLoggerWithoutContext()
	MongoConfig, err := configs.Get(constants.MongoConfig)

	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	mongoUrl := MongoConfig.GetString(constants.MongoUrlKey)
	mongoDatabase := MongoConfig.GetString(constants.MongoDatabaseKey)
	if mongoUrl == "" || mongoDatabase == "" {
		log.With(zap.Error(errors.New("configuration is missing for mongodb"))).Error("Configuration is missing for mongodb")
	}

	err = NewConnection(mongoDatabase, mongoUrl)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	//  	CreateIndex()

	//  	 // Start watching all collections
	//   if err := WatchCollections(context.Background()); err != nil {
	// 	log.With(zap.Error(err)).Error("Failed to initialize collection watchers")
	// 	}
	log.Info("MongoDB connection established")
}
