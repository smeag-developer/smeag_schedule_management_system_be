package services

import (
	"nxt_match_event_manager_api/internal/interfaces"
	models "nxt_match_event_manager_api/internal/models/notification"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationService struct {
	mu   sync.Mutex
	repo interfaces.NotificationRespositoryInterface
}

func NewNotificationService(
	repo interfaces.NotificationRespositoryInterface) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) GetAllNotificationService(
	ctx *gin.Context,
	m *models.PayloadQueryNotificationRef,
	cursor bson.ObjectID, limit int) (*[]models.NotificationListRequest, error) {

	return s.repo.GetAllNotificationsByLimit(ctx, m, cursor, limit)
}

func (s *NotificationService) BroadCastToUserService(ctx *gin.Context) error {
	return s.repo.BroadCastToUser(ctx)
}

func (s *NotificationService) RetrieveNotificationKeyService(ctx *gin.Context, id bson.ObjectID) (*[]interface{}, error) {
	return s.repo.FindFCMTokenNotificationRepository(ctx, id)
}

func (s *NotificationService) CreateNotificationTokenService(ctx *gin.Context, model *models.RegisterTokenRequest) error {
	return s.repo.CreateNotificationTokenRepository(ctx, model)
}
