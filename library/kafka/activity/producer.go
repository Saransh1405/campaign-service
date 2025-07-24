package activity

import (
	"campaign-service/models"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

var ActivityProducer sarama.SyncProducer

func SendActivityDataToKafka(payloadData *models.ActivityEvent, topic string) error {
	jsonString, err := json.Marshal(payloadData)
	if err != nil {
		return fmt.Errorf("failed to marshal payload data: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonString),
	}

	_, _, err = ActivityProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return nil
}
