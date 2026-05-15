package repository

import (
	fcmModel "nxt_match_event_manager_api/internal/models/fcm"
	models "nxt_match_event_manager_api/internal/models/notification"
	pl "nxt_match_event_manager_api/internal/models/pipeline"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *NotificationRepository) RequestGroupJoin(ctx *gin.Context, model *models.NotificationRequest) error {
	_, err := insertDocument(ctx, r.NotificationsCol, model)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) CreateNotificationRepository(ctx *gin.Context, notif []fcmModel.StatePushPayload) error {
	_, err := insertManyDocument(ctx, r.NotificationsCol, notif)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) GetAllNotificationsByLimit(
	ctx *gin.Context, m *models.PayloadQueryNotificationRef, cursor bson.ObjectID, limit int) (*[]models.NotificationListRequest, error) {

	return MultiAggregateQuery[models.NotificationListRequest](
		ctx, r.NotificationsCol,
		pl.FETCH_ALL_NOTIFICATIONS_BY_LIMIT(m.NotificationTokenId, m.Platform, m.FCMToken, cursor, limit))
}

func (r *NotificationRepository) UpsertNotificationStatus(ctx *gin.Context, status string, notifId bson.ObjectID) error {

	update := bson.M{"notification_req.status": status}

	return updateOne[models.NotificationRequest](ctx, r.NotificationsCol, notifId, update)
}

func (r *NotificationRepository) BroadCastToUser(ctx *gin.Context) error {
	return nil
}

func (r *NotificationRepository) CreateNotificationTokenRepository(ctx *gin.Context, model *models.RegisterTokenRequest) error {

	keyId := bson.M{"token": model.Token}
	docs := bson.M{
		"$set": bson.M{
			"user_id":  model.UserId,
			"token":    model.Token,
			"platform": model.Platform,
		},
	}

	return flatUpdate(ctx, r.NotificationTokenCol, keyId, docs)
}

func (r *NotificationRepository) FindFCMTokenNotificationRepository(
	ctx *gin.Context, userId bson.ObjectID) (*[]interface{}, error) {
	return findMany[interface{}](ctx, r.NotificationTokenCol, userId)
}
