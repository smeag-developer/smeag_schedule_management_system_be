package repository

import (
	"context"
	"fmt"
	"log/slog"
	cc "nxt_match_event_manager_api/internal/constants"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type NotificationRepository struct {
	NotificationsCol     *mongo.Collection
	NotificationTokenCol *mongo.Collection
	MatchNotificationCol *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) *NotificationRepository {
	return &NotificationRepository{
		NotificationsCol:     db.Collection(cc.NOTIFICATIONS),
		NotificationTokenCol: db.Collection(cc.NOTIFICATION_TOKENS),
	}
}

func insertDocument[T any](ctx context.Context, col *mongo.Collection, doc T) (bson.ObjectID, error) {
	res, err := col.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			slog.Error("Duplicate key error: already exists", err.Error(), "%s")
			return bson.NilObjectID, nil
		}
		slog.Error("Error inserting document : %v", err.Error(), "%s")
		return bson.NilObjectID, err
	}

	oid, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		slog.Error("Unable to retrieve object id")
		return bson.NilObjectID, fmt.Errorf("invalid object ID type")
	}

	return oid, nil
}

func insertManyDocument[T any](ctx context.Context, col *mongo.Collection, docs T) ([]bson.ObjectID, error) {
	res, err := col.InsertMany(ctx, docs)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			slog.Error("Duplicate key error: already exists", "error ", err.Error())
			return nil, err
		}
		slog.Error("Error inserting document: %v", "error", err.Error())
		return nil, err
	}

	// Convert []interface{} → []bson.ObjectID
	oids := make([]bson.ObjectID, 0, len(res.InsertedIDs))
	for _, id := range res.InsertedIDs {
		oid, ok := id.(bson.ObjectID)
		if !ok {
			slog.Error("Unable to cast inserted ID to ObjectID")
			return nil, fmt.Errorf("invalid object ID type")
		}
		oids = append(oids, oid)
	}

	return oids, nil
}

func updateOne[T any](ctx context.Context, col *mongo.Collection, id bson.ObjectID, update bson.M) error {

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedDoc T

	err := col.FindOneAndUpdate(ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$set", Value: update}},
		opts,
	).Decode(&updatedDoc)

	if err != nil {
		return err
	}

	return nil
}

func MultiAggregateQuery[T any](ctx context.Context, col *mongo.Collection, pipeline mongo.Pipeline) (*[]T, error) {

	cursor, err := col.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []T

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func findOne[T any](ctx context.Context, col *mongo.Collection, id bson.ObjectID) (*T, error) {
	var raw T

	filter := bson.M{"_id": bson.M{"$eq": id}}
	err := col.FindOne(ctx, filter).Decode(&raw)

	if err != nil {
		return &raw, err
	}

	return &raw, err
}

func findMany[T any](ctx context.Context, col *mongo.Collection, filter interface{}) (*[]T, error) {
	var results []T

	cursor, err := col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func flatUpdate(ctx context.Context, col *mongo.Collection, keyId bson.M, docs bson.M) error {
	res, err := col.UpdateOne(ctx, keyId, docs, options.UpdateOne().SetUpsert(true))

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 && res.UpsertedCount == 0 {
		return fmt.Errorf("no document matched filter")
	}
	return nil
}
