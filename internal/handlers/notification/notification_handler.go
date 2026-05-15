package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	cc "nxt_match_event_manager_api/internal/constants"
	common "nxt_match_event_manager_api/internal/models/common"
	models "nxt_match_event_manager_api/internal/models/notification"
	notifModels "nxt_match_event_manager_api/internal/models/notification"
	"strconv"
	"time"

	"nxt_match_event_manager_api/internal/utils"
	loggers "nxt_match_event_manager_api/internal/utils/loggers"
	"nxt_match_event_manager_api/pb/buff"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *NotificationHandler) registerToken(ctx *gin.Context) {
	var notif notifModels.RegisterTokenRequest
	var actionType string
	ttl := time.Until(time.Now().AddDate(0, 1, 0)) // max ttl for redis hash token (1 month)

	// validate valid payload
	utils.ValidatePayload(ctx, &notif)

	//insert to persistence
	if err := h.notifService.CreateNotificationTokenService(ctx, &notif); err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	// insert or update to in-memory
	fieldExist, err := h.redisClient.HashFieldExists(ctx, cc.FCM_NOTIFICATION_KEY, notif.UserId)

	if err != nil {
		loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	if fieldExist {

		v, err := h.redisClient.GetHashSet(ctx, cc.FCM_NOTIFICATION_KEY, notif.UserId)

		if err != nil {
			loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
			return
		}

		// cast to array
		var data []notifModels.RegisterTokenRequest
		if err := json.Unmarshal([]byte(v), &data); err != nil {
			loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
			return
		}

		/*
		** Validate some users might have multiple register tokens
		 */
		found := false
		for _, t := range data {
			if string(t.Token) == notif.Token {
				found = true
				break
			}
		}

		if found {
			actionType = fmt.Sprintf("%s, %s", "[redis]", cc.TOKEN_ALREADY_STORED)
		} else {

			arr := append(data, notif)
			b, err := json.Marshal(arr)

			if err != nil {
				loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
				return
			}

			// update the new token
			_, err = h.redisClient.UpsertHashSet(ctx, cc.FCM_NOTIFICATION_KEY, notif.UserId, string(b), ttl)

			if err != nil {
				loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
				return
			}
			actionType = fmt.Sprintf("%s, %s", "[redis]", cc.TOKEN_ALREADY_STORED)
		}

	} else {
		/**
		* 	if no field found create a new one
		*	with array of Register Token Request
		**/
		arr := []notifModels.RegisterTokenRequest{}
		arr = append(arr, notif)

		b, err := json.Marshal(arr)

		if err != nil {
			loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
			return
		}

		// insert token
		_, err = h.redisClient.UpsertHashSet(ctx, cc.FCM_NOTIFICATION_KEY, notif.UserId, string(b), ttl)

		if err != nil {
			loggers.GetCommonError(ctx, err.Error(), http.StatusInternalServerError)
			return
		}

		actionType = fmt.Sprintf("%s, %s", "[redis]", cc.SUCCESFULLY_CREATED_TOKEN)
	}

	loggers.SuccessResponse(ctx, actionType)
}

func (h *NotificationHandler) handleGetNotificationByLimit(ctx *gin.Context) {

	p := &models.PayloadQueryNotificationRef{
		NotificationTokenId: ctx.Param(cc.ID),
		CursorId:            ctx.Query(cc.CURSOR),
		Limit:               ctx.Query(cc.LIMIT),
		Platform:            ctx.Query(cc.PLATFORM),
		FCMToken:            ctx.Query(cc.FCM_TOKEN),
	}

	cursor_id, cErr := utils.MustHex(p.CursorId)

	// //validate cursor id if not empty
	// skip no cursor id to next iterration
	if len(p.CursorId) > 0 {
		if cErr != nil {
			loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid hex id"))
			return
		}
	}

	conv, err := strconv.Atoi(p.Limit)

	if err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid limit type :%v", err))
		return
	}

	if len(p.NotificationTokenId) == 0 {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid notif token id:%v", err))
		return
	}

	res, err := h.notifService.GetAllNotificationService(ctx, p, cursor_id, conv)

	if err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("service error :%v", err))
		return
	}

	var nextCursor string
	l, pErr := strconv.Atoi(string(p.Limit))

	if pErr != nil {
		loggers.StatusBadRequestError(ctx, pErr)
		return
	}

	if len(*res) == l {
		// dereference group
		var id = (*res)[len(*res)-1].Id

		if !bson.ObjectID.IsZero(id) {
			nextCursor = id.Hex()

		} else {
			nextCursor = ""
		}
	}

	loggers.StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS_GET_NOTIFICATIONS,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body: &common.NotificationPaginationResponse{
			Notifications: &res,
			NextCursor:    nextCursor,
			HasMore:       nextCursor != "",
		},
	})
}

func (h *NotificationHandler) handleGroupJoinStatus(ctx *gin.Context) {

	var notif notifModels.GroupInfoUpdateStatus
	var groupReq = buff.GroupNotificationStatusRequest{}
	utils.ValidatePayload(ctx, &notif)

	ids := []string{
		notif.ID,
		notif.NotificationId,
		notif.GroupJoiners[0].UserId,
	}

	_, err := utils.ValidateMultipleHexId(ids)
	if err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid hex id"))
		return
	}

	// deep copy
	copier.Copy(&groupReq, notif)
	groupStatusReq := &buff.GroupStatusRequest{GroupNotifStatusRequest: &groupReq}

	if err := h.grpcHandler.UpdateGroupNotifStatus(ctx, groupStatusReq); err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid update status :%v", err))
		return
	}

	loggers.StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS_UPDATE_NOTIFICATION_STATUS,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body:       "notification status updated successfully",
	})

}

func (h *NotificationHandler) handleEventJoinStatus(ctx *gin.Context) {

	var notif notifModels.EventInfoUpdateStatus
	var eventReq = buff.EventNotificationStatusRequest{}
	utils.ValidatePayload(ctx, &notif)

	ids := []string{
		notif.ID,
		notif.NotificationId,
		notif.EventJoiner[0].UserId,
	}

	_, err := utils.ValidateMultipleHexId(ids)

	if err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid hex id"))
		return
	}

	// deep copy
	copier.Copy(&eventReq, notif)
	eventStatusReq := &buff.EventStatusRequest{EventNotifStatusRequest: &eventReq}

	if err := h.grpcHandler.UpdateEventNotifStatus(ctx, eventStatusReq); err != nil {
		loggers.StatusBadRequestError(ctx, fmt.Errorf("invalid update status :%v", err))
		return
	}

	loggers.StatusOK(ctx, &common.SuccessResponse{
		SuccessID:  bson.NewObjectID(),
		Status:     cc.SUCCESS_UPDATE_NOTIFICATION_STATUS,
		HttpCode:   http.StatusOK,
		ResponseAt: time.Now(),
		Body:       "notification status updated successfully",
	})
}
