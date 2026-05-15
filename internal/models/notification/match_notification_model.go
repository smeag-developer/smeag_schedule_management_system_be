package models

import "time"

type MatchNotificationRequest struct {
	NotificationReq   NotificationReqInfo   `bson:"notification_req" json:"notification_req"`
	MatchInfo         MatchInfo             `bson:"match_info" json:"match_info"`
	Body              MatchNotificationBody `bson:"body" json:"body"`
	Requestor         RequestorInfo         `bson:"requestor" json:"requestor"`
	NotificationMsgId string                `bson:"notification_msg_id" json:"notification_msg_id,omitempty"`
}

type MatchNotificationBody struct {
	UserId     string `bson:"user_id" json:"user_id"`
	Title      string `bson:"title" json:"title"`
	ImageURL   string `bson:"image_url" json:"image_url"` // image or and icon of notification
	Body       string `bson:"body" json:"body"`
	Platform   string `bson:"platform" json:"platform"`
	TTLSeconds int    `bson:"ttl_seconds" json:"ttl_seconds"`
}

type MatchJoinRequest struct {
	UserReceiverId    string    `json:"user_receiver_id"`
	UserRequestorInfo User      `json:"user_requestor_info"`
	MatchInfo         MatchInfo `json:"match_info"`
	UserActivityId    string    `json:"user_activity_id"`
	RequestTime       time.Time `json:"request_time"`
}

type MatchInfo struct {
	Id         string  `json:"id" bson:"_id,omitempty"`
	MatchName  string  `json:"match_name" bson:"match_name"`
	MatchType  string  `json:"match_type" bson:"match_type"`
	LandMark   string  `json:"land_mark" bson:"land_mark"`
	IsPaid     bool    `json:"is_paid" bson:"is_paid"`
	MatchPrice float64 `json:"match_price" bson:"match_price"`
	StartHour  string  `json:"start_hour" bson:"start_hour"`
	EndHour    string  `json:"end_hour" bson:"end_hour"`
	Status     string  `json:"status" bson:"status"` //
}

type User struct {
	Id             string `json:"id" bson:"id,omitempty"`
	FirstName      string `json:"firstName" bson:"firstName"`
	LastName       string `json:"lastName" bson:"lastName"`
	Email          string `json:"email" bson:"email"`
	ProfilePicture string `json:"profilePicture" bson:"profilePicture"`
}
