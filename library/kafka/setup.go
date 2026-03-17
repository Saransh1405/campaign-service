package kafka

import (
	"campaign-service/constants"
	"campaign-service/utils/configs"
	"fmt"
	"log"
)

// kafka command
// kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic campaign_activity

func NewConnection() {
	// get application configs
	kafkaConfig, err := configs.Get(constants.KafkaConfig)
	if err != nil {
		log.Fatalf("Failed to get application config: %v", err)
	}

	kafkaHost := kafkaConfig.GetString(constants.KafkaHostKey)
	if kafkaHost == "" {
		fmt.Println("Kafka host is not set")
	}

	kafkaUsername := kafkaConfig.GetString(constants.KafkaUsernameKey)
	kafkaPassword := kafkaConfig.GetString(constants.KafkaPasswordKey)

	if kafkaHost == "" {
		kafkaHost = "kafka:9092"
	}

	if kafkaUsername == "" {
		kafkaUsername = ""
	}

	if kafkaPassword == "" {
		kafkaPassword = ""
	}

	// kafka connection
	err = Connect([]string{kafkaHost}, kafkaUsername, kafkaPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	log.Println("Kafka connection established")
}
