package kafka

import (
	"log"
	"users-service/constants"
	"users-service/utils/configs"
)

// kafka command
// kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic create_campaign

func NewConnection() {
	// get application configs
	applicationConfig, err := configs.Get(constants.ApplicationConfig)
	if err != nil {
		log.Fatalf("Failed to get application config: %v", err)
	}

	kafkaHost := applicationConfig.GetString(constants.KafkaHostConfigKey)
	if kafkaHost == "" {
		log.Fatalf("Kafka host is not set")
	}
	kafkaUsername := applicationConfig.GetString(constants.KafkaUsernameConfigKey)
	kafkaPassword := applicationConfig.GetString(constants.KafkaPasswordConfigKey)

	// kafka connection
	err = Connect([]string{kafkaHost}, kafkaUsername, kafkaPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	log.Println("Kafka connection established")
}
