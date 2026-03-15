package main

import (
	"log"
	"net"
	"os"

	"github-harvester/internal/collector/adapter"
	handlers "github-harvester/internal/collector/handler"
	"github-harvester/internal/collector/service"
	"github-harvester/internal/pkg/pb"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("COLLECTOR_PORT")
	if port == "" {
		port = "8888"
	}

	githubAdapter := adapter.NewGitHubAdapter()
	collectorService := service.NewCollectorService(githubAdapter)
	grpcHandler := handlers.NewHandler(collectorService)

	server := grpc.NewServer()
	pb.RegisterCollectorServer(server, grpcHandler)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Collector gRPC server starting on port %s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
