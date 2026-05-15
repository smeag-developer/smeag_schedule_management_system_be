package repository

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type StudentRepository struct {
	StudentsCol *mongo.Collection
}

func NewStudentRepository(db *mongo.Database) *StudentRepository {
	return &StudentRepository{
		StudentsCol: db.Collection("student"),
	}
}

func insertDoc[T any](ctx context.Context, col *mongo.Collection, doc T) (bson.ObjectID, error) {
	res, err := col.InsertOne(ctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			slog.Error("duplicate key error", "error", err.Error())
			return bson.NilObjectID, fmt.Errorf("student already exists")
		}
		return bson.NilObjectID, err
	}

	oid, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.NilObjectID, fmt.Errorf("invalid object ID type")
	}

	return oid, nil
}

func updateDoc(ctx context.Context, col *mongo.Collection, id bson.ObjectID, fields bson.M) error {
	res, err := col.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.D{{Key: "$set", Value: fields}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("student not found")
	}
	return nil
}

func deleteDoc(ctx context.Context, col *mongo.Collection, id bson.ObjectID) error {
	res, err := col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("student not found")
	}
	return nil
}

func findWithPagination[T any](ctx context.Context, col *mongo.Collection, filter bson.D, limit int64) (*[]T, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := col.Find(ctx, filter, opts)
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
