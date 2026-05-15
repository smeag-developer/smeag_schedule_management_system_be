package services

import (
	fcmModel "nxt_match_event_manager_api/internal/models/fcm"

	"github.com/gin-gonic/gin"
)

// func (s *NotificationService) FindFCMTokenUserServices(ctx *gin.Context, userId bson.ObjectID) error {
// 	return s.repo.FindFCMTokenNotificationRepository(ctx, userId)
// }

func (s *NotificationService) CreateMatchNotificationService(ctx *gin.Context, notif []fcmModel.StatePushPayload) error {
	return s.repo.CreateNotificationRepository(ctx, notif)
}
