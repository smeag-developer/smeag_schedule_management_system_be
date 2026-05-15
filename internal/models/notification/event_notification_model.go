package models

import (
	"time"
)

type EventNotificationRequest struct {
	NotificationReq   NotificationReqInfo `bson:"notification_req" json:"notification_req"`
	Events            EventNotifInfo      `bson:"events" json:"events"`
	Body              NotificationBody    `bson:"body" json:"body"`
	Requestor         RequestorInfo       `bson:"requestor" json:"requestor"`
	NotificationMsgId string              `bson:"notification_msg_id" json:"notification_msg_id,omitempty"`
}

// Event
type EventNotifInfo struct {
	EventId        string `json:"event_id" bson:"event_id"`
	EventName      string `json:"event_name" bson:"event_name"`
	EventCreatorId string `json:"event_creator_id" bson:"event_creator_id"`
	Role           string `bson:"role" json:"role"`
	Status         string `bson:"status" json:"status"`
}

type EventInfo struct {
	EventId          string     `json:"event_id" bson:"event_id"`
	EventName        string     `json:"event_name" bson:"event_name"`
	EventDescription string     `json:"event_description" bson:"event_description"`
	EventVenue       EventVenue `json:"event_venue" bson:"event_venue"`
	EventStartTime   string     `json:"event_start_time" bson:"event_start_time"`
	EventEndTime     string     `json:"event_end_time" bson:"event_end_time"`
	EventTimezone    string     `json:"event_timezone" bson:"event_timezone"`
	EventDate        time.Time  `json:"event_timedate" bson:"event_timedate"`
	EventLevel       string     `json:"event_level" bson:"event_level"`
	EventType        string     `json:"event_type" bson:"event_type"`
	IsFree           bool       `json:"is_free" bson:"is_free"`
	Price            float32    `json:"price" bson:"price"`
	Currency         string     `json:"currency" bson:"currency"`
	GenderCategory   string     `json:"gender_category" bson:"gender_category"`
}

type EventJoiner struct {
	UserId     string `json:"user_id" bson:"user_id"`
	Role       string `json:"role" bson:"role"`
	Status     string `json:"status" bson:"status"`
	DateJoined string `json:"date_joined" bson:"date_joined"`
}

type EventVenue struct {
	Name               string  `json:"name" bson:"name"` // Specific venue name (e.g., "Cebu Sports Complex")
	Location           string  `json:"location" bson:"location"`
	Type               string  `json:"type" bson:"type"`
	Latitude           float64 `json:"latitude" bson:"latitude"`                         // For map integration
	Longitude          float64 `json:"longitude" bson:"longitude"`                       // For map integration
	BookingEventStatus string  `json:"booking_event_status" bson:"booking_event_status"` // Reserved, Waiting for Booking , Pending
}

type EventInfoUpdateStatus struct {
	ID                      string        `bson:"_id,omitempty" json:"id,omitempty"`
	NotificationId          string        `bson:"notification_id" json:"notification_id"`
	UserRequestorActivityId string        `bson:"user_requestor_activity_id" json:"user_requestor_activity_id"`
	EventJoiner             []EventJoiner `bson:"event_joiner" json:"event_joiner"`
}

type EventOrganizer struct {
	OrganizerId      string `json:"organizer_id" bson:"organizer_id"`
	OrganizerName    string `json:"organizer_name" bson:"organizer_name"`
	OrganizerContact string `json:"organizer_contact" bson:"organizer_contact"`
}

type RegisterTokenRequest struct {
	UserId   string `json:"user_id" bson:"user_id"`
	Token    string `json:"token" bson:"token"`
	Platform string `json:"platform" bson:"platform"`
}
