package campaign

import (
	"campaign-service/models"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

func SendDataToKafka(payloadData *models.CampaignEvent, topic string) error {
	jsonString, err := json.Marshal(payloadData)
	if err != nil {
		return fmt.Errorf("failed to marshal payload data: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonString),
	}

	_, _, err = Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return nil
}
