package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcAdapter "github.com/Aiya594/shipment-microservice/internal/adapters/grpc"
	"github.com/Aiya594/shipment-microservice/internal/adapters/postgres"
	"github.com/Aiya594/shipment-microservice/internal/app"
	cfgApp2 "github.com/Aiya594/shipment-microservice/internal/domain/configs/app"
	cfgDB "github.com/Aiya594/shipment-microservice/internal/domain/configs/db"
	"github.com/Aiya594/shipment-microservice/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize database connection
	cfg := cfgDB.LoadConfigsDB()
	ctx := context.Background()
	db, err := cfg.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := postgres.NewShipmentRepo(db)

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize service
	service := app.NewShipmentService(repo, logger)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	shipmentServer := grpcAdapter.NewShipmentServer(service, logger)
	proto.RegisterShipmentServiceServer(grpcServer, shipmentServer)

	// Start server
	cfgApp := cfgApp2.LoadConfigApp()

	lis, err := net.Listen("tcp", ":"+cfgApp.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logger.Info("Server succcessfully started", "port", cfgApp.Port)

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
