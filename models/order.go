package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// OrderItem represents an individual item in an order
type OrderItem struct {
	SKUID    string `bson:"sku_id,omitempty" json:"sku_id"` // SKU identifier
	Quantity int    `bson:"quantity" json:"quantity"`       // Quantity of item ordered
}

// Order represents the main order entity
type Order struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"` // Order ID
	// SellerID     primitive.ObjectID  `bson:"seller_id,omitempty" json:"seller_id"`         // Seller reference
	// HubID        primitive.ObjectID  `bson:"hub_id,omitempty" json:"hub_id"`               // Hub reference
	CustomerName string              `bson:"customer_name,omitempty" json:"customer_name"` // Customer name
	OrderNo      string              `bson:"order_no" json:"order_no"`                     // Unique order number
	OrderItems   []OrderItem         `bson:"order_items" json:"order_items"`               // Array of order items
	Status       string              `bson:"status" json:"status"`                         // Order status (on_hold initially)
	CreatedAt    primitive.DateTime  `bson:"created_at,omitempty" json:"created_at"`       // Timestamp of order creation
	UpdatedAt    primitive.DateTime  `bson:"updated_at,omitempty" json:"updated_at"`       // Timestamp of last update
	DeletedAt    *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at"`       // Timestamp of soft deletion (optional)
}
