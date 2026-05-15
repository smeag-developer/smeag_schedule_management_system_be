package models

import (
	"time"
)

type PushPayload struct {
	Token      string            `bson:"token" json:"token"`
	Title      string            `bson:"title" json:"title"`
	Body       string            `bson:"body" json:"body"`
	Data       map[string]string `bson:"data" json:"data"`
	Platform   string            `bson:"platform" json:"platform"`
	TTLSeconds int               `bson:"ttl" json:"ttl"` // 0 = use FCM default (4 weeks)
}

type StatePushPayload struct {
	Payload               PushPayload           `bson:"payload" json:"payload"`
	Status                string                `bson:"status"`
	Error                 string                `bson:"error"`
	NotificationTokenId   string                `bson:"notification_token_id"` // user_id + hash256 -> used to as notification identifier
	NotificationType      string                `bson:"notification_type"`
	DeviceTokenId         string                `bson:"device_token_id"`
	Requestor             Requestor             `bson:"requestor"`
	NotificationReference NotificationReference `bson:"notification_reference" json:"notification_reference"`
	Timestamp             time.Time             `bson:"time_stamp"`
}
type Requestor struct {
	UserId         string `bson:"user_id"`
	UserActivityId string `bson:"user_activity_id"`
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
