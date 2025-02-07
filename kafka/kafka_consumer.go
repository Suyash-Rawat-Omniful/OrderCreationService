package kafka

import (
	"context"
	"encoding/json"
	"service1/models"
	"time"

	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

// Implement message handler
type MessageHandler struct{}

// Process implements pubsub.IPubSubMessageHandler.
func (h *MessageHandler) Process(ctx context.Context, message *pubsub.Message) error {
	log.Printf("Received message: %s", string(message.Value))

	// Define a variable to hold the parsed data
	var order models.KafkaResponseOrderMessage
	err := json.Unmarshal(message.Value, &order)
	if err != nil {
		log.WithError(err).Error("Failed to parse Kafka message")
		return err
	}

	// Call WMS Inventory Checking logic from here

	return nil
}

func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
	// Process message
	return nil
}

// Initialize Kafka Consumer
func InitializeKafkaConsumer(ctx context.Context) {
	consumer := kafka.NewConsumer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithConsumerGroup("my-consumer-group"),
		kafka.WithClientID("my-consumer"),
		kafka.WithKafkaVersion("2.8.1"),
		kafka.WithRetryInterval(time.Second),
	)

	handler := &MessageHandler{}
	consumer.RegisterHandler("oms-service-topic2", handler)
	consumer.Subscribe(ctx)
}
