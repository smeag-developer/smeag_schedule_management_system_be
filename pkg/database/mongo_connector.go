package database

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	models "nxt_match_event_manager_api/internal/models/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func MongoInitConnector(dbConfig *models.DBconfig) (*mongo.Client, error) {

	client, err := ConnectWithRetry(dbConfig.Uri, 5, 2*time.Second)

	if err != nil {
		log.Fatalf("Failed to establish MongoDB connection: %v", err)
	}

	return client, nil
}

func ConnectWithRetry(uri string, maxRetries int, baseDelay time.Duration) (*mongo.Client, error) {
	var client *mongo.Client
	var err error

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Printf("Trying to connect to MongoDB (attempt %d/%d)...", i+1, maxRetries)
		client, err = mongo.Connect(options.Client().ApplyURI(uri))

		if err == nil {
			slog.Info("Successfully connected DB")
			return client, nil
		}

		// Verify the connection
		err := client.Ping(ctx, nil)
		if err != nil {
			slog.Error("Failed to ping MongoDB", "error", err)
			return nil, err
		}

		log.Printf("Connection failed: %v", err)

		// Exponential backoff
		sleep := baseDelay * time.Duration(1<<i)
		log.Printf("Retrying in %v...", sleep)
		time.Sleep(sleep)
	}

	return nil, fmt.Errorf("could not connect to MongoDB after %d attempts: %v", maxRetries, err)
}
