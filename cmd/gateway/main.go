package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github-harvester/docs"
	"github-harvester/internal/gateway/adapter"
	handlers "github-harvester/internal/gateway/handler"
	"github-harvester/internal/gateway/service"
	"github-harvester/internal/pkg/pb"

	httpSwagger "github.com/swaggo/http-swagger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title Github Harvester API
// @description     API Gateway for collecting GitHub stats.
func main() {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:8888"
	}

	if httpPort == "" || collectorAddr == "" {
		log.Fatal("HTTP_PORT and COLLECTOR_ADDR environment variables must be set")
	}

	conn, err := grpc.NewClient(collectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial collector at %s: %v", collectorAddr, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}()

	grpcClient := pb.NewCollectorClient(conn)
	adapter := adapter.NewCollectorAdapter(grpcClient)
	svc := service.NewHarvesterService(adapter)
	handler := handlers.NewRepoHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /repo/{owner}/{repo}", handler.GetRepo)
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)
	wrappedMux := LoggingMiddleware(mux)
	srv := &http.Server{
		Addr:         ":" + httpPort,
		Handler:      wrappedMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Gateway is running on port %s, connecting to collector at %s", httpPort, collectorAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen error: %s\n", err)
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HTTP] %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
