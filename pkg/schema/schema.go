package schema

import (
	"context"
	"fmt"
	"log"
	cc "nxt_match_event_manager_api/internal/constants"
	"slices"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SchemaBuilder struct {
	collections []string
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		collections: slices.Compact([]string{ // avoid duplicate collection names
			cc.NOTIFICATIONS,
			cc.MATCH_NOTIFICATIONS,
		}),
	}
}

// Verify if the database exists, if not create one
func (s *SchemaBuilder) PrepareInitDB(client *mongo.Client, db_name string) {

	c, err := client.ListDatabaseNames(context.Background(), bson.M{})

	if err != nil {
		log.Printf("un-initialize the database: %v", db_name)
		panic(fmt.Sprintf("error %v", err))
	}

	if !slices.Contains(c, db_name) {
		client.Database(db_name)
	}
}

// Verify if the collections exists, if not create one
func (s *SchemaBuilder) PrepareCollections(db *mongo.Database) error {

	for _, v := range s.collections {
		collects, err := db.ListCollectionNames(context.Background(), bson.M{})

		if err != nil {
			return err
		}

		if !slices.Contains(collects, v) {
			err := db.CreateCollection(context.Background(), v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
