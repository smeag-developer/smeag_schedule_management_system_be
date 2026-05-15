package repository

import (
	notifModel "nxt_match_event_manager_api/internal/models/notification"
	pl "nxt_match_event_manager_api/internal/models/pipeline"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *NotificationRepository) GetAllEventNotificationByUser(
	ctx *gin.Context,
	token_id string,
	cursor bson.ObjectID,
	limit int) (*[]notifModel.NotificationRequest, error) {

	return MultiAggregateQuery[notifModel.NotificationRequest](ctx, r.NotificationsCol, pl.GET_ALL_NOTIFICATIONS_BY_LIMIT(token_id, cursor, limit))
}

func (r *NotificationRepository) RequestCreateEvent(ctx *gin.Context, model *notifModel.EventNotificationRequest) error {
	_, err := insertDocument(ctx, r.NotificationsCol, model)
	if err != nil {
		return err
	}

	return nil
}
