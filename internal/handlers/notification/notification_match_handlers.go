package handlers

import (
	"net/http"
	cc "nxt_match_event_manager_api/internal/constants"
	models "nxt_match_event_manager_api/internal/models/notification"
	redisClient "nxt_match_event_manager_api/internal/redis"
	"nxt_match_event_manager_api/internal/utils"
	loggers "nxt_match_event_manager_api/internal/utils/loggers"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
)

func (h *NotificationHandler) handleRequestJoinMatchHandler(ctx *gin.Context) {
	var notif models.NotificationRequest

	if err := ctx.ShouldBindJSON(&notif); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": cc.NO_PAYLOAD_FOUND})
		return
	}

	tokenInfo, err := h.redisClient.GetHMSETtoken(ctx, cc.FCM_NOTIFICATION_KEY, notif.ReceiverInfo.Id)
	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	ts, err := utils.SliceInterfaceToStruct[models.RegisterTokenRequest](tokenInfo)
	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	dv := map[string]string{
		"type":                  "join_request",
		"notification_match_id": notif.GroupInfo.GroupId, // ← see note below
	}

	shaId := utils.MustStringSHA256(notif.ReceiverInfo.Id)
	fcmPayload := h.fcmClient.FCMProccesor(ctx, ts, &notif, dv, shaId)

	var wg sync.WaitGroup
	var enqueueErr, dbErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		_, enqueueErr = h.redisClient.EnqueueNotification(ctx, redisClient.StreamKey, fcmPayload)
	}()

	go func() {
		defer wg.Done()
		dbErr = h.notifService.CreateMatchNotificationService(ctx, fcmPayload)
	}()

	wg.Wait()

	if enqueueErr != nil {
		loggers.GetCommonError(ctx, enqueueErr.Error(), http.StatusInternalServerError)
		return
	}

	if dbErr != nil {
		loggers.GetCommonError(ctx, dbErr.Error(), http.StatusBadRequest)
		return
	}

	// try to send (active token)
	for _, p := range fcmPayload {
		if slices.Contains(h.hub.GetActiveDeviceTokens(), p.DeviceTokenId) {
			h.redisClient.StartConsumeNotification(ctx, p.DeviceTokenId)
			break
		}
	}
}
