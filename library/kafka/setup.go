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
	applicationConfig, err := configs.Get(constants.ApplicationConfig)
	if err != nil {
		log.Fatalf("Failed to get application config: %v", err)
	}

	kafkaHost := applicationConfig.GetString(constants.KafkaHostConfigKey)
	if kafkaHost == "" {
		fmt.Println("Kafka host is not set")
	}

	kafkaUsername := applicationConfig.GetString(constants.KafkaUsernameConfigKey)
	kafkaPassword := applicationConfig.GetString(constants.KafkaPasswordConfigKey)

	if kafkaHost == "" {
		kafkaHost = "kafka:9092"
	}

	if kafkaUsername == "" {
		kafkaUsername = ""
	}

	if kafkaPassword == "" {
		kafkaPassword = ""
	}

	fmt.Printf("kafkaHost: %v\n", kafkaHost)

	// kafka connection
	err = Connect([]string{kafkaHost}, kafkaUsername, kafkaPassword)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	log.Println("Kafka connection established")
}
