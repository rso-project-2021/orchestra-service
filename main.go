package main

import (
	"log"
	"orchestra-service/api"
	"orchestra-service/config"
	"orchestra-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// Load configuration settings.
	config, err := config.New(".")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Connect to gRPC server.
	conn, err := grpc.Dial(config.GrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
