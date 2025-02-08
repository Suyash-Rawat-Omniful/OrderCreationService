package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"service1/database"
	"service1/models"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"github.com/omniful/go_commons/http"
	interservice_client "github.com/omniful/go_commons/interservice-client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Start() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	topic := os.Getenv("TOPIC")
	msgCnt := 0

	worker, err := ConnectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer started ")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCnt++
				var str string
				if err = json.Unmarshal(msg.Value, &str); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("recieved :", str, "\n\n")
				}
				fmt.Printf("Received order Count %d: | Topic(%s) | Message(%s) \n", msgCnt, string(msg.Topic), str)
				order := string(msg.Value)
				fmt.Printf("order: ", order, "\n\n\n")

				objectID, _ := primitive.ObjectIDFromHex(str)

				fmt.Print("objectID is ", objectID, "\n\n")

				collection := database.DB.Database("OMS").Collection("orders")
				var savedOrder models.Order
				err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&savedOrder)
				if err != nil {
					log.Printf("Order not found: %v", err)
					continue
				}
				var flag bool = true
				for ind := range savedOrder.OrderItems {
					SkuID := savedOrder.OrderItems[ind].SKUID
					Quantity := savedOrder.OrderItems[ind].Quantity
					skuIDStr := SkuID
					config := interservice_client.Config{
						ServiceName: "user-service",
						BaseURL:     "http://localhost:8081/api/V1/inventories/",
						Timeout:     5 * time.Second,
					}
					client, err := interservice_client.NewClientWithConfig(config)
					if err != nil {
						flag = false
						continue
					}
					url := config.BaseURL + "validate/" + skuIDStr
					body := map[string]interface{}{
						"sku_id":         SkuID,
						"given_quantity": Quantity,
					}
					bodyBytes, err := json.Marshal(body)
					if err != nil {
						flag = false
						continue
					}
					req := &http.Request{
						Url:     url,
						Body:    bytes.NewReader(bodyBytes),
						Timeout: 7 * time.Second,
						Headers: map[string][]string{
							"Content-Type": {"application/json"},
						},
					}
					fmt.Print("url is ", url, "\n")
					resp, _ := client.Get(req, body)
					if err != nil || resp == nil {
						fmt.Printf("Error making GET request to validate inventory: %v\n", err)
						flag = false
						continue
					} else {
						fmt.Print("response of the validate inventory is request is ", resp, "\n\n\n\n")
						if resp.IsSuccess() {

						}
						continue
					}
				}
				if flag {
					savedOrder.Status = "new_order"
				} else {
					savedOrder.Status = "failed inventory check"
				}
				first, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, savedOrder)
				if err != nil {
					fmt.Println("failed to update order status:", err)
				} else {
					fmt.Println("order status updated successfully:", first)
				}

			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	fmt.Println("Processed", msgCnt, "messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}

}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}
