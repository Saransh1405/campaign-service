package kafka

import (
	"campaign-service/library/kafka/campaign"
	"log"

	"github.com/IBM/sarama"
)

func Connect(brokerList []string, KafkaUsername, KafkaPassword string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Only enable SASL if username and password are provided
	if KafkaUsername != "" && KafkaPassword != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = KafkaUsername
		config.Net.SASL.Password = KafkaPassword
	} else {
		// Disable SASL for local development
		config.Net.SASL.Enable = false
	}

	// Create campaign producer
	createCampaignProducer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create campaign producer: %v", err)
	}
	campaign.Producer = createCampaignProducer

	// Create update campaign producer
	updateCampaignProducer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create update campaign producer: %v", err)
	}
	campaign.Producer = updateCampaignProducer

	// Start consumers in background
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	go func() {
		// Start Campaign Kafka consumers
		if err := campaign.StartConsumers(brokerList, config); err != nil {
			log.Printf("Error starting campaign consumers: %v", err)
		}
	}()

	return nil
}
