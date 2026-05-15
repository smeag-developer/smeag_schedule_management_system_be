package handlers

import (
	"net/http"
	cc "nxt_match_event_manager_api/internal/constants"
	models "nxt_match_event_manager_api/internal/models/notification"
	notifModels "nxt_match_event_manager_api/internal/models/notification"
	redisClient "nxt_match_event_manager_api/internal/redis"
	"nxt_match_event_manager_api/internal/utils"
	loggers "nxt_match_event_manager_api/internal/utils/loggers"
	"nxt_match_event_manager_api/pb/buff"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func (h *NotificationHandler) handleGroupJoinNotification(ctx *gin.Context) {

	var notif notifModels.NotificationRequest
	var buffGroupReq buff.GroupJoinNotificationRequest

	utils.ValidatePayload(ctx, &notif)
	copier.Copy(&buffGroupReq, &notif)

	req := &buff.CreateGroupJoinRequest{GroupInfoRequest: &buffGroupReq}

	err := h.grpcHandler.RequestGroupJoinService(ctx, req)
	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	exist, err := h.redisClient.HashFieldExists(ctx, cc.FCM_NOTIFICATION_KEY, notif.ReceiverInfo.Id)
	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exist {
		loggers.GetCommonError(ctx, "field does not exist", http.StatusNotFound)
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
		"notification_group_id": notif.GroupInfo.GroupId,
	}

	shaId := utils.MustStringSHA256(notif.ReceiverInfo.Id)
	fcmPayload := h.fcmClient.FCMProccesor(ctx, ts, &notif, dv, shaId)

	var wg sync.WaitGroup
	var dbErr error
	var enqueueErr error

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

	// try to send if (device token)
	for _, p := range fcmPayload {
		if slices.Contains(h.hub.GetActiveDeviceTokens(), p.DeviceTokenId) {
			h.redisClient.StartConsumeNotification(ctx, p.DeviceTokenId)
			break
		}
	}
}
