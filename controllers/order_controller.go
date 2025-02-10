package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"os"
	"strconv"
	"time"

	"service1/database"
	oms_kafka "service1/kafka"
	"service1/models"
	orders "service1/service"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/csv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRequest struct {
	FilePath string `json:"file_path"`
}

func UploadOrders(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	if _, err := os.Stat(req.FilePath); os.IsNotExist(err) {
		fmt.Println("file not found : ", req.FilePath)
		c.JSON(400, gin.H{"error": "File not found"})
		return
	}

	orders.SetProducer(c, database.Queue, req.FilePath)

	c.JSON(200, gin.H{"message": "Orders uploaded successfully to the queue"})
}

func CallParseCSV(c *gin.Context) {
	ParseCSV("storage/orders.csv")
}

func ParseCSV(filePath string) {
	orders, err := performcsvopr(filePath)
	if err != nil {
		fmt.Println("\nfailed to parse csv with path : ", filePath)
		return
	}
	for _, order := range orders {
		if err := storeOrder(order); err != nil {
			fmt.Print("\nparseCSV : falied to save order")
			return
		}
	}

}
func performcsvopr(filePath string) ([]*models.Order, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	orderGroups := make(map[string]*models.Order)

	Csv, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),
		csv.WithSource(csv.Local),
		csv.WithLocalFileInfo(filePath),
		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}
	err = Csv.InitializeReader(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}
	for !Csv.IsEOF() {
		var records csv.Records
		records, err := Csv.ReadNextBatch()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Processing records:")
		fmt.Println(records)
		for _, record := range records {
			orderNo := record[0]
			customerName := record[1]
			skuIDStr := record[2]
			quantityStr := record[3]
			hubIDStr := record[4]

			skuID, err := strconv.Atoi(skuIDStr)
			if err != nil {
				fmt.Println("invalid sku_id", skuIDStr, ":", err)
				continue
			}
			hubID, err := strconv.Atoi(hubIDStr)
			if err != nil {
				fmt.Println("invalid hub_id", hubIDStr, ":", err)
				continue
			}
			if !isValidHub(uint(hubID)) {
				fmt.Println("invalid hub_id", hubIDStr, ":", err)
				continue
			}
			if !isValidSKU(uint(skuID)) {
				fmt.Println("sku -> ", skuIDStr, "doesn't exists")
				continue
			}

			quantity, err := strconv.Atoi(quantityStr)
			if err != nil {
				fmt.Print("invalid quantity ", quantityStr, ":", err)
				continue
			}

			orderKey := fmt.Sprintf("%s-%s", orderNo, customerName)
			order, exists := orderGroups[orderKey]
			if !exists {
				now := primitive.NewDateTimeFromTime(time.Now())
				order = &models.Order{
					ID:           primitive.NewObjectID(),
					CustomerName: customerName,
					OrderNo:      orderNo,
					OrderItems:   []models.OrderItem{},
					Status:       "on_hold",
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				orderGroups[orderKey] = order
			}

			orderItem := models.OrderItem{
				SKUID:    skuIDStr,
				Quantity: quantity,
				HubID:    hubIDStr,
			}
			order.OrderItems = append(order.OrderItems, orderItem)
		}
	}

	var orders []*models.Order
	for _, order := range orderGroups {
		orders = append(orders, order)
	}

	fmt.Println("Final orders:")
	for _, order := range orders {
		fmt.Printf("Order No: %s, Customer: %s, Total Items: %d\n", order.OrderNo, order.CustomerName, len(order.OrderItems))

		_, err := json.Marshal(order)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Print("calling store order\n")
		storeOrder(order)
	}

	return orders, nil
}

func storeOrder(order *models.Order) error {
	fmt.Print("stored order mei agaye\n")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DB := database.DB
	if DB == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}
	collection := DB.Database("OMS").Collection("orders")

	now := primitive.NewDateTimeFromTime(time.Now())
	order.ID = primitive.NewObjectID()
	order.CreatedAt = now
	order.UpdatedAt = now

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		log.Printf("Failed to insert order: %v", err)
		return err
	}

	insertedID := result.InsertedID.(primitive.ObjectID).Hex()
	fmt.Print("Order inserted with ID:", insertedID, "\n\n\n\n")

	x, _ := json.Marshal(insertedID)
	oms_kafka.PublishMessageToKafka(x, insertedID)

	return nil
}
