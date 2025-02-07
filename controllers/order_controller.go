package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"service1/database"
	oms_kafka "service1/kafka"
	"service1/models"
	orders "service1/service"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/pubsub"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderRequest represents the API request payload
type OrderRequest struct {
	FilePath string `json:"file_path"`
}

// UploadOrders handles order creation from a CSV file
func UploadOrders(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate file existence
	if _, err := os.Stat(req.FilePath); os.IsNotExist(err) {
		fmt.Println("file not found : ", req.FilePath)
		c.JSON(400, gin.H{"error": "File not found"})
		return
	}

	orders.SetProducer(c, database.Queue, req.FilePath)

	// Parse CSV file and create orders
	// orders, err := performcsvopr(req.FilePath)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	// // Store each order in MongoDB (optional, depending on the flow)
	// for _, order := range orders {
	// 	if err := storeOrder(order); err != nil {
	// 		c.JSON(500, gin.H{"error": "Failed to save order"})
	// 		return
	// 	}
	// }

	c.JSON(200, gin.H{"message": "Orders uploaded successfully to the queue"})
}

func CallParseCSV(c *gin.Context) {
	ParseCSV("storage/orders.csv")
}

func ParseCSV(filePath string) {
	orders, err := performcsvopr(filePath)
	if err != nil {
		// c.JSON(500, gin.H{"error": err.Error()})
		fmt.Println("\nfailed to parse csv with path : ", filePath)
		return
	}

	// Store each order in MongoDB (optional, depending on the flow)
	for _, order := range orders {
		if err := storeOrder(order); err != nil {
			fmt.Print("\nparseCSV : falied to save order")
			return
		}
	}

}

// performcsvopr reads the CSV file and creates orders from the CSV records
func performcsvopr(filePath string) ([]*models.Order, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Map to group items by order_no and customer_name
	orderGroups := make(map[string]*models.Order)

	// Initialize the CSV reader (based on your previous implementation)
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

	// Process the records and group them by order_no and customer_name
	for !Csv.IsEOF() {
		var records csv.Records
		records, err := Csv.ReadNextBatch()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Processing records:")
		fmt.Println(records)
		for _, record := range records {
			orderNo := record[0]      // order_no
			customerName := record[1] // customer_name
			skuIDStr := record[2]     // sku_id
			quantityStr := record[3]  // quantity
			hubIDStr := record[4]

			skuID, err := strconv.Atoi(skuIDStr)
			if err != nil {
				fmt.Println("invalid sku_id", skuIDStr, ":", err)
				continue
			}
			if !isValidSKU(uint(skuID)) {
				// Skip this record if the SKU ID is not valid
				//print this sku was not valid
				fmt.Println("sku -> ", skuIDStr, "doesn't exists")
				continue
			}

			hubID, err := strconv.Atoi(hubIDStr)
			if err != nil {
				fmt.Println("invalid hub_id", hubIDStr, ":", err)
				continue
			}
			if !isValidHub(uint(hubID)) {
				// Skip this record if the SKU ID is not valid
				//print this sku was not valid
				fmt.Println("hub -> ", skuIDStr, "doesn't exists")
				continue
			}

			// Convert quantity to integer
			quantity, err := strconv.Atoi(quantityStr)
			if err != nil {
				fmt.Print("invalid quantity ", quantityStr, ":", err)
				continue
			}

			// Check if the order group for this order_no and customer_name already exists
			orderKey := fmt.Sprintf("%s-%s", orderNo, customerName)
			order, exists := orderGroups[orderKey]
			if !exists {
				// If order doesn't exist, create a new order
				now := primitive.NewDateTimeFromTime(time.Now())
				order = &models.Order{
					ID: primitive.NewObjectID(),
					// SellerID:     sellerID,
					// HubID:        hubID,
					CustomerName: customerName,
					OrderNo:      orderNo,
					OrderItems:   []models.OrderItem{}, // Start with an empty slice of items
					Status:       "on_hold",
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				// Add the new order to the map
				orderGroups[orderKey] = order
			}

			// Create a new OrderItem and append it to the order's OrderItems
			orderItem := models.OrderItem{
				SKUID:    skuIDStr,
				Quantity: quantity,
			}
			order.OrderItems = append(order.OrderItems, orderItem)
		}
	}

	// Convert the map of orders into a slice
	var orders []*models.Order
	for _, order := range orderGroups {
		orders = append(orders, order)
	}

	fmt.Println("Final orders:")
	for _, order := range orders {
		fmt.Printf("Order No: %s, Customer: %s, Total Items: %d\n", order.OrderNo, order.CustomerName, len(order.OrderItems))
		var msg pubsub.Message
		var event string = "publishing order of " + order.CustomerName + " and order number is " + order.OrderNo + "(has " + strconv.Itoa(len(order.OrderItems)) + ")"
		headers := map[string]string{
			"event": event,
		}
		msg.Topic = "Bulk_orders32"
		msg.Value = []byte(fmt.Sprintf("Order No: %s, Customer: %s, Total Items: %d\n", order.OrderNo, order.CustomerName, len(order.OrderItems)))
		msg.Key = order.OrderNo
		msg.Headers = headers
		// ctx := context.WithValue(context.Background(), "request_id", "req-123")
		// var stri string = (fmt.Sprintf("Order No: %s, Customer: %s, Total Items: %d\n", order.OrderNo, order.CustomerName, len(order.OrderItems)))
		value, err := json.Marshal(order)
		if err != nil {
			fmt.Println(err)
			continue
		}
		oms_kafka.PublishMessageToKafka(value, order.OrderNo)
		// if err != nil {
		// 	fmt.Print("message not sent to topic\n")
		// } else {
		// 	fmt.Print("message sent to topic\n")
		// }
	}

	return orders, nil
}

// storeOrder inserts a single order into MongoDB
func storeOrder(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Database("OMS").Collection("orders")

	_, err := collection.InsertOne(ctx, order)
	return err
}

func isValidSKU(skuID uint) bool {
	// Construct the URL for the GET request
	url := fmt.Sprintf("http://localhost:8081/api/V1/skus/validate/%d", skuID)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request to validate SKU: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return false
	}

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		// SKU is valid, return true
		fmt.Printf("SKU %d is valid\n", skuID)
		return true
	} else {
		// SKU is not valid, return false
		fmt.Printf("Invalid SKU ID: %d\nResponse: %s\n", skuID, body)
		return false
	}
}

func isValidHub(hubID uint) bool {
	// Construct the URL for the GET request
	url := fmt.Sprintf("http://localhost:8081/api/V1/hubs/validate/%d", hubID)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making GET request to validate Hub: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return false
	}

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		// SKU is valid, return true
		fmt.Printf("Hub %d is valid\n", hubID)
		return true
	} else {
		// SKU is not valid, return false
		fmt.Printf("Invalid hub ID: %d\nResponse: %s\n", hubID, body)
		return false
	}
}
