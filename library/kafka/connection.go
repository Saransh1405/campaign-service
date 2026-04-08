package kafka

import (
	"campaign-service/library/kafka/activity"
	"campaign-service/library/kafka/campaign"
	"log"

	"github.com/IBM/sarama"
)

func Connect(brokerList []string, KafkaUsername, KafkaPassword string) error {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	if KafkaUsername != "" && KafkaPassword != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = KafkaUsername
		config.Net.SASL.Password = KafkaPassword
	} else {
		config.Net.SASL.Enable = false
	}

	createCampaignProducer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create campaign producer: %v", err)
	}
	campaign.Producer = createCampaignProducer

	updateCampaignProducer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create update campaign producer: %v", err)
	}
	campaign.Producer = updateCampaignProducer

	campaignActivityProducer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Failed to create campaign activity producer: %v", err)
	}
	activity.ActivityProducer = campaignActivityProducer

	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	go func() {
		if err := campaign.StartConsumers(brokerList, config); err != nil {
			log.Printf("Error starting campaign consumers: %v", err)
		}

		if err := activity.StartActivityConsumer(brokerList, config); err != nil {
			log.Printf("Error starting campaign activity consumers: %v", err)
		}
	}()

	return nil
}
