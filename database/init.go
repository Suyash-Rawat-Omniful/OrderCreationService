package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func getDatabaseUri() string {
	return "mongodb://127.0.0.1:27017/OMS"
}

func ConnectMongo(c context.Context) {
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
