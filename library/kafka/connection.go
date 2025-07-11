package kafka

import (
	"github.com/IBM/sarama"
)

func Connect(brokerList []string, KafkaUsername, KafkaPassword string) error {
	// var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.SASL.Enable = true
	config.Net.SASL.User = KafkaUsername
	config.Net.SASL.Password = KafkaPassword

	// // Create cart producer
	// cartProducer, err := sarama.NewSyncProducer(brokerList, config)
	// if err != nil {
	// 	log.Fatalf("Failed to create cart producer: %v", err)
	// }
	// cart.Producer = cartProducer
	// // Create order producer
	// orderProducer, err := sarama.NewSyncProducer(brokerList, config)
	// if err != nil {
	// 	log.Fatalf("Failed to create order producer: %v", err)
	// }
	// orders.Producer = orderProducer
	// // Create streams producer

	// streamsProducer, err := sarama.NewSyncProducer(brokerList, config)
	// if err != nil {
	// 	log.Fatalf("Failed to create streams producer: %v", err)
	// }

	// kafkaStreams.Producer = streamsProducer
	// Start consumers in background
	// Add Consumer Group configs
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// go func() {
	// 	// Start Cart Kafka consumers
	// 	if err := cart.StartConsumers(config, brokerList); err != nil {
	// 		log.Printf("Error starting cart consumers: %v", err)
	// 	}

	// 	// Start Kafka streams consumers
	// 	if err := StartStreamsConsumers(config, brokerList); err != nil {
	// 		log.Printf("Error starting cart consumers: %v", err)
	// 	}
	// }()

	// go func() {
	// 	// Start Order Kafka consumers
	// 	if err := orders.StartConsumers(config, brokerList); err != nil {
	// 		log.Printf("Error starting order consumers: %v", err)
	// 	}
	// }()

	return nil
}
