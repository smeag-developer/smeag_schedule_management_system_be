package handlers

import (
	notif_group_handler "nxt_match_event_manager_api/internal/handlers/grpc"
	redisclient "nxt_match_event_manager_api/internal/redis"
	hub "nxt_match_event_manager_api/internal/routes"
	ws "nxt_match_event_manager_api/internal/routes/socket"
	fcmClient "nxt_match_event_manager_api/internal/services/fcm"
	service "nxt_match_event_manager_api/internal/services/notification"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	grpcHandler  *notif_group_handler.GrpcNotification
	notifService *service.NotificationService
	redisClient  *redisclient.Client
	fcmClient    *fcmClient.Client
	wsConf       *ws.WebSocketConfig
	hub          *hub.Hub
}

func NewHandler(
	notif_service *service.NotificationService,
	group_handler *notif_group_handler.GrpcNotification,
	client *redisclient.Client,
	fcmClient *fcmClient.Client,
	wsConf *ws.WebSocketConfig,
	h *hub.Hub,
) *NotificationHandler {
	return &NotificationHandler{
		notifService: notif_service,
		grpcHandler:  group_handler,
		redisClient:  client,
		fcmClient:    fcmClient,
		wsConf:       wsConf,
		hub:          h,
	}
}
func (h *NotificationHandler) RegisterRoutes(router *gin.RouterGroup) {

	router.GET("/notification/all/:id", h.handleGetNotificationByLimit)

	//Websocket Handler
	router.GET("/ws/notifications", h.wsConf.HandleWsNotification)

	router.POST("/notification/group/join", h.handleGroupJoinNotification)
	router.POST("/notification/event/join", h.handleJoinEventNotification)
	router.POST("/notification/match/join", h.handleRequestJoinMatchHandler)
	router.PUT("/notification/update/group/join/status", h.handleGroupJoinStatus)
	router.PUT("/notification/update/event/join/status", h.handleEventJoinStatus)
	router.POST("/notification/register-token", h.registerToken)
}
