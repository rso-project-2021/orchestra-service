package main

import (
	"context"
	"log"
	"orchestra-service/api"
	"orchestra-service/config"
	"orchestra-service/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// Load configuration settings.
	config, err := config.New(".")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Set gRPC server timeout period.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	// Connect to gRPC server.
	conn, err := grpc.DialContext(ctx, config.GrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial gRPC server: ", err)
	}
	defer conn.Close()

	client := proto.NewLocationServiceClient(conn)

	// Create a server and setup routes.
	server, err := api.NewServer(config, client)
	if err != nil {
		log.Fatal("Failed to create a server: ", err)
	}

	// Start a server.
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("Failed to start a server: ", err)
	}

}
