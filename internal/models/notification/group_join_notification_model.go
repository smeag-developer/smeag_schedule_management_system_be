package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GroupInfo struct {
	ID                bson.ObjectID         `bson:"_id,omitempty" json:"id,omitempty"`
	GroupId           string                `bson:"group_id" json:"group_id"`
	GroupCreatorId    string                `bson:"group_creator_id" json:"group_creator_id"`
	GroupName         string                `bson:"group_name" json:"group_name"`
	Status            string                `bson:"status" json:"status"`
	Role              string                `bson:"role" json:"role"`
	ResponseTimeStamp timestamppb.Timestamp `bson:"response_time_stamp" json:"response_time_stamp"`
}

type GroupJoiner struct {
	UserId     string `bson:"user_id" json:"user_id"`
	Role       string `bson:"role" json:"role"`
	Status     string `bson:"status" json:"status"`
	DateJoined string `bson:"date_joined" json:"date_joined"`
}

type GroupInfoUpdateStatus struct {
	ID                      string        `bson:"_id,omitempty" json:"id,omitempty"`
	NotificationId          string        `bson:"notification_id" json:"notification_id"`
	UserRequestorActivityId string        `bson:"user_requestor_activity_id" json:"user_requestor_activity_id"`
	GroupJoiners            []GroupJoiner `bson:"group_joiners" json:"group_joiners"`
}
