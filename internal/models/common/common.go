package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SuccessResponse struct {
	SuccessID  bson.ObjectID `json:"successId"`
	Status     string        `json:"status"`
	HttpCode   int           `json:"httpCode"`
	ResponseAt time.Time     `json:"requestedAt"`
	Body       interface{}   `json:"body"`
}

type ErrorResponse struct {
	ErrorCode int       `json:"errorCode"`
	TimeStamp time.Time `json:"timeStamp"`
	ErrorMsg  string    `json:"errorMsg"`
}

type NotificationPaginationResponse struct {
	Notifications interface{} `json:"notifications"`
	NextCursor    string      `json:"nextCursor"`
	HasMore       bool        `json:"hasMore"`
}
