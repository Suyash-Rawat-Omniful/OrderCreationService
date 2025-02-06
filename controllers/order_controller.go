package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"service1/database"
	"service1/models"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/csv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderRequest represents the API request payload
type OrderRequest struct {
	FilePath     string             `json:"file_path"`
	SellerID     primitive.ObjectID `json:"seller_id"`
	HubID        primitive.ObjectID `json:"hub_id"`
	CustomerName string             `json:"customer_name"`
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
		c.JSON(400, gin.H{"error": "File not found"})
		return
	}

	// Parse CSV file and create orders
	orders, err := performcsvopr(req.FilePath, req.SellerID, req.HubID, req.CustomerName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Store each order in MongoDB (optional, depending on the flow)
	for _, order := range orders {
		if err := storeOrder(order); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save order"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Orders uploaded successfully", "total_orders": len(orders)})
}

// performcsvopr reads the CSV file and creates orders from the CSV records
func performcsvopr(filePath string, sellerID, hubID primitive.ObjectID, customerName string) ([]*models.Order, error) {
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

		for _, record := range records {
			orderNo := record[0]      // order_no
			customerName := record[1] // customer_name
			skuID := record[2]        // sku_id
			quantityStr := record[3]  // quantity

			// Convert quantity to integer
			quantity, err := strconv.Atoi(quantityStr)
			if err != nil {
				return nil, fmt.Errorf("invalid quantity %s: %v", quantityStr, err)
			}

			// Check if the order group for this order_no and customer_name already exists
			orderKey := fmt.Sprintf("%s-%s", orderNo, customerName)
			order, exists := orderGroups[orderKey]
			if !exists {
				// If order doesn't exist, create a new order
				now := primitive.NewDateTimeFromTime(time.Now())
				order = &models.Order{
					ID:           primitive.NewObjectID(),
					SellerID:     sellerID,
					HubID:        hubID,
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
				SKUID:    skuID,
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
