package kafka

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/pubsub"
)

type KafkaProducer struct {
	Producer *kafka.ProducerClient
}

var producerInstance *KafkaProducer

func SetProducer(producer *kafka.ProducerClient) {
	if producerInstance == nil {
		producerInstance = &KafkaProducer{}
	}
	producerInstance.Producer = producer
}

func getProducer() *kafka.ProducerClient {
	if producerInstance != nil {
		return producerInstance.Producer
	}
	return nil
}

func PublishMessageToKafka(bytesOrderItem []byte, orderID string) {
	ctx := context.Background()
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	topic := os.Getenv("TOPIC")
	msg := &pubsub.Message{
		Topic: topic,
		Key:   orderID,
		Value: bytesOrderItem,
		Headers: map[string]string{
			"custom-header": "value",
		},
	}

	producer := getProducer()
	err := producer.Publish(ctx, msg)
	if err != nil {
		log.Println("Error publishing message to kafka")
	} else {
		log.Println("Message published to kafka")
	}
}
