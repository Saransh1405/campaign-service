package activity

import (
	"campaign-service/logger"
	"campaign-service/models"
	"campaign-service/utils/activity"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

const ActivityTopic = "campaign_activity"

func StartActivityConsumer(brokers []string, config *sarama.Config) error {
	go ConsumerActivity(brokers, config)
	return nil
}

func ConsumerActivity(brokers []string, config *sarama.Config) error {
	log := logger.GetLoggerWithoutContext()

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Error("failed to create consumer", zap.Error(err))
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(ActivityTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Error("failed to create partition consumer", zap.Error(err))
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	const maxRetries = 3
	var event models.ActivityEvent

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

				switch event.Action {
				case "join":
					activity.ConsumeJoinActivity(event.Participant)
				case "leave":
					activity.ConsumeLeaveActivity(event.Participant)
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
