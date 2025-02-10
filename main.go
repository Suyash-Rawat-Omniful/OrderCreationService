package main

import (
	"context"
	"fmt"
	"service1/consumer"
	"service1/database"
	"service1/listeners"
	"service1/router"

	"github.com/omniful/go_commons/http"
)

func main() {
	ctx := context.TODO()
	database.Init(ctx)
	go consumer.Start()
	go listeners.StartConsumer()

	//to setup the server (go_commons se hai ye)
	server := http.InitializeServer(":8080", 0, 0, 70)
	fmt.Print("starting server\n")

	//setup the routes(khud banaya hai)
	err := router.Initialize(ctx, server)
	if err != nil {

	}

	//server start karne ke liye(go_commons ka part hai)
	err = server.StartServer("OMS")
	if err != nil {
		fmt.Println("Error in starting the server")
		return
	}
	fmt.Println("Server started")
}
