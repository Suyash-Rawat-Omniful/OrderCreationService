package main

import (
	"context"
	"fmt"
	"service1/database"
	"service1/router"

	"github.com/omniful/go_commons/http"
)

func main() {
	ctx := context.TODO()
	database.ConnectMongo(ctx)
	server := http.InitializeServer(":8080", 0, 0, 70)
	fmt.Print("starting server")
	err := router.Initialize(ctx, server)
	if err != nil {

	}
	err = server.StartServer("OMS")
	if err != nil {
		fmt.Println("Error in starting the server")
		return
	}
	fmt.Println("Server started")
}

// func runHttpServer(ctx context.Context) {
// 	server := http.InitializeServer(":8080", 10*time.Second, 10*time.Second, 70*time.Second)

// 	// Initialize middlewares and routes
// 	err := router.Initialize(ctx, server)
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		panic(err)
// 	}

// 	log.Infof("Starting server on port" + ":8080")

// 	err = server.StartServer("Tenant-service")
// 	if err != nil {
// 		log.Errorf(err.Error())
// 		panic(err)
// 	}

// 	<-shutdown.GetWaitChannel()
// }
