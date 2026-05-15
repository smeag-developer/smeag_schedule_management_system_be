package main

import (
	"context"
	"log"
	"log/slog"
	conf "smeag_sms_be/config"
	cc "smeag_sms_be/internal/constants"
	"smeag_sms_be/pkg/database"
	"smeag_sms_be/pkg/schema"
)

var BuildEnv string

func main() {
	/*** Initialize Mongo DB Connector*/
	dbConf, dErr := conf.GetDBConfig(BuildEnv)
	hostConf, sErr := conf.GetHostConfig(BuildEnv)

	if dErr != nil {
		slog.Error("Failed to initialize DB config", "error", dErr)
		panic(cc.UNINITIALIZED_PANIC_STATE)
	}

	if sErr != nil {
		slog.Error("Failed to initialize host config", "error", sErr)
		panic(cc.UNINITIALIZED_PANIC_STATE)
	}

	// Connect to MongoDB
	client, err := database.MongoInitConnector(dbConf)

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error on disconnect: %v", err)
		}
	}()

	if err != nil {
		slog.Error("Error connecting to MongoDB %v", err)
	}

	// auto create collections
	schemaBuilder := schema.NewSchemaBuilder()
	schemaBuilder.PrepareInitDB(client, dbConf.DBName)

	if err := schemaBuilder.PrepareCollections(client.Database(dbConf.DBName)); err != nil {
		slog.Error("Failed to prepare collections", "error", err)
		panic(cc.UNINITIALIZED_PANIC_STATE)
	}

	// Connect Server
	conf.ServerConfig(dbConf, hostConf, *client)
}
