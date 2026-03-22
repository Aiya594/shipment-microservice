package main

import (
	"context"
	"log"

	"github.com/Aiya594/shipment-microservice/internal/adapters/postgres"
	"github.com/Aiya594/shipment-microservice/internal/app"
	configs "github.com/Aiya594/shipment-microservice/internal/domain/configs/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize database connection
	cfg := configs.NewConfigsDB()
	ctx := context.Background()
	db, err := cfg.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := postgres.NewShipmentRepo(db)

	// Initialize service
	service := app.NewShipmentService(repo)

}
