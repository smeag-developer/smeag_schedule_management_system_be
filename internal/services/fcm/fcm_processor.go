package fcm

import (
	"context"
	"fmt"
	"log/slog"
	cc "nxt_match_event_manager_api/internal/constants"
	fcmModel "nxt_match_event_manager_api/internal/models/fcm"
	models "nxt_match_event_manager_api/internal/models/notification"
	"sync"
	"time"
)

/*
* Process FCM notification with go routine workers
* @Params:
* @Context- context
* @RegisterToken - Slice of RegisterTokens (notification receiver)
* @NotificationModel - Notification information
* @Data - FCM data params
 */
func (c *Client) FCMProccesor(ctx context.Context,
	rtr []models.RegisterTokenRequest,
	notif *models.NotificationRequest,
	data map[string]string,
	shaId string,
) []fcmModel.StatePushPayload {

	var wg sync.WaitGroup
	notifState := make([]fcmModel.StatePushPayload, len(rtr))

	for i, t := range rtr {
		wg.Add(1)
		go func(idx int, token string, platform string) {
			defer wg.Done()
			c.ProcessFcmSend(ctx, notif, &notifState[idx], data, token, platform, shaId)
			slog.Info(token)
		}(i, t.Token, t.Platform) // ← Capture loop variables explicitly
	}

	wg.Wait()

	return notifState
}

// Send Notification
func (c *Client) ProcessFcmSend(ctx context.Context,
	notif *models.NotificationRequest,
	rc *fcmModel.StatePushPayload,
	data map[string]string,
	token, platform string,
	shaId string,
) error {

	var mu sync.Mutex

	payload := fcmModel.PushPayload{
		Token:      token,
		Title:      notif.Body.Title,
		Body:       notif.Body.Body,
		Data:       data,
		Platform:   platform,
		TTLSeconds: notif.Body.Ttl,
	}

	mu.Lock()
	err := c.PrepareToSend(ctx, payload)
	mu.Unlock()

	if err != nil {
		*rc = fcmModel.StatePushPayload{
			Payload:             payload,
			Status:              cc.STATUS_NOTIF_ERROR,
			Error:               err.Error(),
			Timestamp:           time.Now(),
			NotificationTokenId: "",
			Requestor:           fcmModel.Requestor{},
		}

		return fmt.Errorf("[FCM]: %v", err)
	}

	// convert sha256 the requestor id
	// identifier for receiver of notifications

	fcmState := fcmModel.StatePushPayload{
		Payload:             payload,
		Status:              cc.STATUS_NOTIF_DELIVERED,
		Error:               "-",
		NotificationTokenId: shaId,
		NotificationType:    notif.NotificationReq.NotificationType,
		DeviceTokenId:       token,
		Requestor: fcmModel.Requestor{
			UserId:         notif.Requestor.Id,
			UserActivityId: notif.Requestor.UserActivityId,
		},
		Timestamp: time.Now(),
	}

	// referecen type on which notification belongs
	switch notif.NotificationReq.NotificationType {
	case cc.EVENT_JOIN_REQUEST_TYPE:
		fcmState.NotificationReference.EventRef.EventId = notif.EventInfo.EventId
	case cc.GROUP_JOIN_REQUEST_TYPE:
		fcmState.NotificationReference.GroupRef.GroupId = notif.GroupInfo.GroupId
	case cc.MATCH_JOIN_REQUEST_TYPE:
		fcmState.NotificationReference.MatchRef.MatchId = notif.MatchInfo.Id
	}

	// is success delivery
	*rc = fcmState

	return nil
}
