package models

import (
	"time"

	fcmModel "nxt_match_event_manager_api/internal/models/fcm"

	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/* Common Notification Model */

type NotificationRequest struct {
	Id              bson.ObjectID       `bson:"_id,omitempty" json:"id,omitempty"`
	NotificationReq NotificationReqInfo `bson:"notification_req" json:"notification_req"`
	Body            NotificationBody    `bson:"body" json:"body"`
	GroupInfo       GroupInfo           `bson:"group_info" json:"group_info"`
	EventInfo       EventNotifInfo      `bson:"event_info" json:"event_info"`
	MatchInfo       MatchInfo           `bson:"match_info" json:"match_info"`
	Requestor       RequestorInfo       `bson:"requestor" json:"requestor"`
	ReceiverInfo    ReceiverInfo        `bson:"receiver_info" json:"receiver_info"`
}

type NotificationToken struct {
	Id        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserId    bson.ObjectID `bson:"user_id" json:"user_id"`
	Token     string        `bson:"token" json:"token"`
	Platform  string        `bson:"platform" json:"platform"`
	Timestamp time.Time     `bson:"time_stamp" json:"time_stamp"`
}

type PayloadQueryNotificationRef struct {
	NotificationTokenId string
	CursorId            string
	Limit               string
	Platform            string
	FCMToken            string
}

type RequestorInfo struct {
	Id             string `bson:"_id,omitempty" json:"id,omitempty"`
	UserActivityId string `bson:"user_activity_id" json:"user_activity_id"`
	FirstName      string `bson:"first_name" json:"first_name"`
	LastName       string `bson:"last_name" json:"last_name"`
	ProfilePic     string `bson:"profile_pic" json:"profile_pic"`
}

type ReceiverInfo struct {
	Id             string `bson:"_id,omitempty" json:"id,omitempty"`
	UserActivityId string `bson:"user_activity_id" json:"user_activity_id"`
	FirstName      string `bson:"first_name" json:"first_name"`
	LastName       string `bson:"last_name" json:"last_name"`
	ProfilePic     string `bson:"profile_pic" json:"profile_pic"`
}

type NotificationReqInfo struct {
	Id               bson.ObjectID         `bson:"_id,omitempty" json:"id,omitempty"`
	Status           string                `bson:"status" json:"status"`
	NotificationType string                `bson:"notification_type" json:"notification_type"`   // e.g., "group_join_request"
	RequestTimeStamp timestamppb.Timestamp `bson:"request_time_stamp" json:"request_time_stamp"` // google grpc timestamp
}

type NotificationBody struct {
	Title            string `bson:"title" json:"title"`
	ImageURL         string `bson:"image_url" json:"image_url"` // image or and icon of notification
	Body             string `bson:"body" json:"body"`
	IsClicked        bool   `bson:"is_clicked" json:"is_clicked"`
	TimeStampRequest string `bson:"time_stamp_request" json:"time_stamp_request"`
	Ttl              int    `bson:"ttl" json:"ttl"`
}

type ProtoTimestamp struct {
	Seconds int64 `bson:"seconds" json:"seconds,omitempty"`
	Nanos   int32 `bson:"nanos" json:"nanos,omitempty"`
}

type NotificationEvent struct {
	UserID     string            `json:"user_id"`
	Platform   string            `json:"platform"`
	Title      string            `json:"title"`
	Body       string            `json:"body"`
	Data       map[string]string `json:"data,omitempty"`
	TTLSeconds int               `json:"ttl_seconds,omitempty"` // passed to FCM
}

type NotificationListRequest struct {
	Id                    bson.ObjectID         `bson:"_id,omitempty" json:"id,omitempty"`
	Payload               fcmModel.PushPayload  `bson:"payload" json:"payload"`
	Status                string                `bson:"status" json:"status"`
	Error                 string                `bson:"error" json:"error"`
	NotificationTokenId   string                `bson:"notification_token_id" json:"notification_token_id"` // user_id + hash256 -> used to as notification identifier
	NotificationType      string                `bson:"notification_type" json:"notification_type"`
	Requestor             Requestor             `bson:"requestor" json:"requestor"`
	NotificationReference NotificationReference `bson:"notification_reference" json:"notification_reference"`
	Timestamp             time.Time             `bson:"time_stamp" json:"time_stamp"`
}

type Requestor struct {
	UserId         string `bson:"user_id" json:"user_id"`
	UserActivityId string `bson:"user_activity_id" json:"user_activity_id"`
}

type NotificationReference struct {
	GroupRef GroupNotificationReference    `bson:"group_reference" json:"group_reference"`
	EventRef EventNotificationReference    `bson:"event_reference" json:"event_reference"`
	MatchRef MatchRefNotificationReference `bson:"match_reference" json:"match_reference"`
}

type GroupNotificationReference struct {
	GroupId string `bson:"group_id" json:"group_id"`
}

type EventNotificationReference struct {
	EventId string `bson:"event_id" json:"event_id"`
}

type MatchRefNotificationReference struct {
	MatchId string `bson:"match_id" json:"match_id"`
}
