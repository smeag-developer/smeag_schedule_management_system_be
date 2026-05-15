package interfaces

import (
	fcmModel "nxt_match_event_manager_api/internal/models/fcm"
	models "nxt_match_event_manager_api/internal/models/notification"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationRespositoryInterface interface {
	RequestGroupJoin(ctx *gin.Context, model *models.NotificationRequest) error
	RequestCreateEvent(ctx *gin.Context, model *models.EventNotificationRequest) error
	GetAllNotificationsByLimit(ctx *gin.Context, m *models.PayloadQueryNotificationRef, cursor bson.ObjectID, limit int) (*[]models.NotificationListRequest, error)

	UpsertNotificationStatus(ctx *gin.Context, status string, notifId bson.ObjectID) error
	BroadCastToUser(ctx *gin.Context) error

	GetAllEventNotificationByUser(ctx *gin.Context, token_id string, cursor_id bson.ObjectID, limit int) (*[]models.NotificationRequest, error)
	GetAllMatchNotificationByUser(ctx *gin.Context, token_id string, cursor_id bson.ObjectID, limit int) (*[]models.NotificationRequest, error)

	// Match Repositories
	CreateMatchJoinNotification(ctx *gin.Context, model *models.NotificationRequest) error

	CreateNotificationRepository(ctx *gin.Context, model []fcmModel.StatePushPayload) error
	CreateNotificationTokenRepository(ctx *gin.Context, model *models.RegisterTokenRequest) error
	FindFCMTokenNotificationRepository(ctx *gin.Context, userId bson.ObjectID) (*[]interface{}, error)
}
