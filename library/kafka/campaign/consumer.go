package campaign

import (
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils/campaigns"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const CreateCampaignTopic = "create_campaign"
const UpdateCampaignTopic = "update_campaign"

func StartConsumers(brokers []string, config *sarama.Config) error {
	go ConsumerCreateCampaign(brokers, config)
	go ConsumerUpdateCampaign(brokers, config)
	return nil
}

func ConsumerCreateCampaign(brokers []string, config *sarama.Config) error {
	log := logger.GetLoggerWithoutContext()

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Error("failed to create consumer", zap.Error(err))
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(CreateCampaignTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Error("failed to create partition consumer", zap.Error(err))
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	const maxRetries = 3
	var event models.CampaignEvent

	for {
		select {
		case err := <-partitionConsumer.Errors():
			log.Error("error from partition consumer", zap.Error(err))
		case msg := <-partitionConsumer.Messages():
			var lastError error
			success := false

			for attempt := 0; attempt < maxRetries; attempt++ {
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				campaignData, err := json.Marshal(event.Campaign)
				if err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				var campaign models.Campaign
				if err := json.Unmarshal(campaignData, &campaign); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				if err := campaigns.InsertIntoPostgres(campaign); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				success = true
				break
			}

			if !success {
				log.Error("failed to process campaign event", zap.Error(lastError))
				continue
			}

			log.Info("campaign event processed successfully", zap.Any("event", event))
		}
	}
}

func ConsumerUpdateCampaign(brokers []string, config *sarama.Config) error {
	log := logger.GetLoggerWithoutContext()

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Error("failed to create consumer", zap.Error(err))
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(UpdateCampaignTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Error("failed to create partition consumer", zap.Error(err))
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	const maxRetries = 3
	var event models.CampaignEvent

	for {
		select {
		case err := <-partitionConsumer.Errors():
			log.Error("error from partition consumer", zap.Error(err))
		case msg := <-partitionConsumer.Messages():
			var lastError error
			success := false

			for attempt := 0; attempt < maxRetries; attempt++ {
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				if err := campaigns.UpdateCampaignInPostgres(event.Campaign, event.UpdateFields); err != nil {
					lastError = err
					time.Sleep(time.Duration(attempt+1) * time.Second)
					continue
				}

				success = true
				break
			}

			if !success {
				log.Error("failed to process campaign event", zap.Error(lastError))
				continue
			}

			log.Info("campaign event processed successfully", zap.Any("event", event))
		}
	}
}
