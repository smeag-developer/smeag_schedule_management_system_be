package handlers

import (
	"log/slog"
	"net/http"
	cc "nxt_match_event_manager_api/internal/constants"
	models "nxt_match_event_manager_api/internal/models/notification"
	notifModels "nxt_match_event_manager_api/internal/models/notification"
	redisClient "nxt_match_event_manager_api/internal/redis"
	"nxt_match_event_manager_api/internal/utils"
	loggers "nxt_match_event_manager_api/internal/utils/loggers"
	"nxt_match_event_manager_api/pb/buff"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *NotificationHandler) handleJoinEventNotification(ctx *gin.Context) {
	start := time.Now()

	var notif notifModels.NotificationRequest
	utils.ValidatePayload(ctx, &notif)

	var buffEventReq buff.EventJoinNotificationRequest
	copier.Copy(&buffEventReq, &notif)
	req := &buff.CreateEventJoinRequest{EventInfoRequest: &buffEventReq}

	// ✅ Run gRPC + Redis cache check concurrently
	type cacheResult struct {
		exists    bool
		tokenInfo []interface{}
		err       error
	}

	grpcErrCh := make(chan error, 1)
	cacheCh := make(chan cacheResult, 1)

	go func() {
		grpcErrCh <- h.grpcHandler.RequestJoinEventService(ctx, req)
	}()

	go func() {
		exists, err := h.redisClient.HashFieldExists(ctx, cc.FCM_NOTIFICATION_KEY, notif.ReceiverInfo.Id)
		if err != nil {
			cacheCh <- cacheResult{err: err}
			return
		}
		if exists {
			tokens, err := h.redisClient.GetHMSETtoken(ctx, cc.FCM_NOTIFICATION_KEY, notif.ReceiverInfo.Id)
			cacheCh <- cacheResult{exists: true, tokenInfo: tokens, err: err}
		} else {
			cacheCh <- cacheResult{exists: false}
		}
	}()

	// Collect gRPC result
	if err := <-grpcErrCh; err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// Collect cache result
	cr := <-cacheCh
	if cr.err != nil {
		loggers.GetCommonError(ctx, cr.err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ DB fallback only when needed
	var tokenInfo []interface{}
	if !cr.exists {
		bsonId, _ := bson.ObjectIDFromHex(notif.ReceiverInfo.Id)
		ptrTokenInfo, err := h.notifService.RetrieveNotificationKeyService(ctx, bsonId)
		if err != nil {
			loggers.GetCommonError(ctx, cc.UNABLE_TO_FIND_NOTIF_TOKENS, http.StatusNotFound)
			return
		}
		if ptrTokenInfo != nil {
			tokenInfo = *ptrTokenInfo
		}
	} else {
		tokenInfo = cr.tokenInfo
	}

	ts, err := utils.SliceInterfaceToStruct[models.RegisterTokenRequest](tokenInfo)
	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	dv := map[string]string{
		"type":                  "join_request",
		"notification_event_id": notif.GroupInfo.GroupId,
	}

	shaId := utils.MustStringSHA256(notif.ReceiverInfo.Id)
	fcmPayload := h.fcmClient.FCMProccesor(ctx, ts, &notif, dv, shaId)

	var wg sync.WaitGroup
	var dbErr, enqueueErr error

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

	if dbErr != nil {
		loggers.GetCommonError(ctx, dbErr.Error(), http.StatusBadRequest)
		return
	}
	if enqueueErr != nil {
		loggers.GetCommonError(ctx, enqueueErr.Error(), http.StatusBadRequest) // ✅ fixed
		return
	}

	// ✅ Build a set once — O(n) lookup instead of O(n*m)
	activeTokenSet := make(map[string]struct{})
	for _, t := range h.hub.GetActiveDeviceTokens() {
		activeTokenSet[t] = struct{}{}
	}

	for _, p := range fcmPayload {
		if _, ok := activeTokenSet[p.DeviceTokenId]; ok {
			h.redisClient.StartConsumeNotification(ctx, p.DeviceTokenId)
			break
		}
	}

	slog.Info("total handleJoinEventNotification", "elapsed", time.Since(start))
}
