package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItem struct {
	SKUID    string `bson:"sku_id,omitempty" json:"sku_id"`
	Quantity int    `bson:"quantity" json:"quantity"`
	HubID    string `bson:"hub_id,omitempty" json:"hub_id"`
}

type Order struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	CustomerName string              `bson:"customer_name,omitempty" json:"customer_name"`
	OrderNo      string              `bson:"order_no" json:"order_no"`
	OrderItems   []OrderItem         `bson:"order_items" json:"order_items"`
	Status       string              `bson:"status" json:"status"` // Order status (on_hold initially)
	CreatedAt    primitive.DateTime  `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt    primitive.DateTime  `bson:"updated_at,omitempty" json:"updated_at"`
	DeletedAt    *primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at"`
}
