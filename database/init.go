package database

import (
	"context"
	"fmt"
	"os"
	"time"

	// redis "service1/Redis"

	"service1/kafka"

	"github.com/joho/godotenv"

	k "github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/sqs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
var Queue *sqs.Queue

func Init(c context.Context) {
	connectMongo(c)
	initializeSqs(c)
	initializeKafkaProducer()
	// go kafka.InitializeKafkaConsumer(c)
}

func getDatabaseUri() string {
	return "mongodb://127.0.0.1:27017/OMS"
}

func connectMongo(c context.Context) {
	fmt.Println("Connecting to mongo...")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(getDatabaseUri())
	var err error
	DB, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	err = DB.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Successfully connected to MongoDB!")
}
func initializeSqs(ctx context.Context) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	acc := os.Getenv("AWS_ACCOUNT")
	queue_name := os.Getenv("QUEUE_NAME")

	sqsConfig := sqs.GetSQSConfig(ctx, false, "ord", "eu-north-1", acc, "")
	url, err := sqs.GetUrl(ctx, sqsConfig, queue_name)
	fmt.Println(*url)
	if err != nil {
		fmt.Println(err)
	}
	queueInstance, err := sqs.NewStandardQueue(ctx, queue_name, sqsConfig)
	Queue = queueInstance
	if err != nil {
		fmt.Println(err)
	}

}
func initializeKafkaProducer() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	topic := os.Getenv("TOPIC")
	kafkaBrokers := make([]string, 1)
	kafkaBrokers[0] = "localhost:9092"
	kafkaClientID := topic
	kafkaVersion := "2.8.1"
	fmt.Print("kafka version is : ", kafkaVersion, "\n")

	producer := k.NewProducer(
		k.WithBrokers(kafkaBrokers),
		k.WithConsumerGroup("my-consumer-group"),
		k.WithClientID(kafkaClientID),
		k.WithKafkaVersion(kafkaVersion),
	)
	fmt.Println("Initialized Kafka Producer")
	kafka.SetProducer(producer)

}
