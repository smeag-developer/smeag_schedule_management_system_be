package services

import (
	fcmModel "nxt_match_event_manager_api/internal/models/fcm"

	"github.com/gin-gonic/gin"
)

func (s *NotificationService) CreateGroupNotificationService(ctx *gin.Context, data []fcmModel.StatePushPayload) error {
	return s.repo.CreateNotificationRepository(ctx, data)
}
