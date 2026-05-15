package pipeline

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var GET_ALL_NOTIFICATIONS_BY_LIMIT = func(token string, cursorId bson.ObjectID, limit int) mongo.Pipeline {
	notifPipeline := mongo.Pipeline{
		// Stage 1: Match by token
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "notification_token_id", Value: token},
		}}},
	}

	// Stage 2: Apply pagination cursor if provided
	if cursorId != bson.NilObjectID {
		notifPipeline = append(notifPipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$lt", Value: cursorId},
				}},
			}},
		})
	}
	// Stage 4: Sort and Limit
	notifPipeline = append(notifPipeline,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "body.time_stamp_request", Value: -1}}}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	return notifPipeline
}

var FETCH_ALL_NOTIFICATIONS_BY_LIMIT = func(token string, platform string, payloadToken string, cursorId bson.ObjectID, limit int) mongo.Pipeline {
	notifPipeline := mongo.Pipeline{
		// Stage 1: Match by token, platform, and payload token
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "notification_token_id", Value: token},
			{Key: "payload.platform", Value: platform},
			{Key: "payload.token", Value: payloadToken},
		}}},
	}

	// Stage 2: Apply pagination cursor if provided
	if cursorId != bson.NilObjectID {
		notifPipeline = append(notifPipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$lt", Value: cursorId},
				}},
			}},
		})
	}

	// Stage 3: Sort and Limit
	notifPipeline = append(notifPipeline,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "time_stamp", Value: -1}}}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	return notifPipeline
}
