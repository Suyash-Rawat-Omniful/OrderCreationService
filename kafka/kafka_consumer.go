package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"service1/models"
	"time"

	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/pubsub"
)

type MessageHandler struct{}

func (h *MessageHandler) Process(ctx context.Context, message *pubsub.Message) error {
	log.Printf("Received message: %s", string(message.Value))
	var order models.Order
	err := json.Unmarshal(message.Value, &order)
	if err != nil {
		log.WithError(err).Error("Failed to parse Kafka message")
		return err
	}

	return nil
}

func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
	// Process message
	return nil
}

func InitializeKafkaConsumer(ctx context.Context) {
	fmt.Print("starting consumer\n")
	consumer := kafka.NewConsumer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithConsumerGroup("my-consumer-group"),
		kafka.WithClientID("oms-service-topic2"),
		kafka.WithKafkaVersion("2.8.1"),
		kafka.WithRetryInterval(time.Second),
	)

	handler := &MessageHandler{}
	err := consumer.RegisterHandler("oms-service-topic2", handler)
	log.Error(err)
	fmt.Print("startedconsumer\n")
	consumer.Subscribe(ctx)
	fmt.Print("staadfadfa\n")
}
