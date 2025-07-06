package mongoDb

import (
	"context"
	"errors"
	"fmt"
	"time"
	"users-service/constants"
	"users-service/logger"
	"users-service/utils/configs"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Connection - Connection structure
type Connection struct {
	Conn     *mongo.Client
	ConnDB   *mongo.Database
	Database string
}

// Client - MongoDB Connection
var Client *Connection

// NewConnection - new connection of amqp
func NewConnection(log logger.Logger) error {

	MongoConfig, err := configs.Get(constants.MongoConfig)
	if err != nil {
		log.With(zap.Error(err)).Error(constants.BindingFailedErrr)
	}

	mongoUrl := MongoConfig.GetString(constants.MongoUrlKey)
	mongoDatabase := MongoConfig.GetString(constants.MongoDatabaseKey)

	if mongoUrl == "" || mongoDatabase == "" {
		return errors.New("COnfiguration is missing for mongodb")
	}

	mongoClient := &Connection{
		Conn:     nil,
		ConnDB:   nil,
		Database: mongoDatabase,
	}
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(mongoUrl)

	mongoClient.Conn, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	err = mongoClient.Conn.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	// @TODO Why are we printing mongodb url in logs?
	log.Info(fmt.Sprintf("Connected to MONGODB_URL : %s.", mongoUrl))

	mongoClient.ConnDB = mongoClient.Conn.Database(mongoClient.Database)
	Client = mongoClient

	return nil
}

// SetupTimeSeriesCollection setup mongo TimeSeries collections
func SetupTimeSeriesCollection(ctx context.Context, log logger.Logger) error {
	// Check if collection exists
	collections, err := Client.ConnDB.ListCollectionNames(ctx, bson.D{{}})
	if err != nil {
		return err
	}

	// Create the collection if it does not exist
	exists := false
	for _, name := range collections {
		if name == constants.MongoAPIUsagesCollection {
			exists = true
			break
		}
	}

	// Create the collection only if it doesn't exist
	if !exists {
		// Setup api usages analytics mongodb collection
		tso := options.TimeSeries().SetTimeField("timestamp").SetMetaField("apiKeyId")
		opts := options.CreateCollection().SetTimeSeriesOptions(tso)

		err2 := Client.ConnDB.CreateCollection(ctx, constants.MongoAPIUsagesCollection, opts)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

// GetCollection - Helper Functions
func GetCollection(collectionName string) *mongo.Collection {
	return Client.ConnDB.Collection(collectionName)
}

// DbContext - Helper Functions
func DbContext(i time.Duration) (context.Context, context.CancelFunc) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), i*time.Second)
	return ctx, ctxCancel
}
