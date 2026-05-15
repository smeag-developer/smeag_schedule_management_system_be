package redis

import (
	"context"
	"encoding/json"
	"log/slog"
	fcmModel "nxt_match_event_manager_api/internal/models/fcm"

	"github.com/redis/go-redis/v9"
)

// EnqueueNotification adds a message to the Redis Stream
func (c *Client) EnqueueNotification(ctx context.Context, streamKey string, notif []fcmModel.StatePushPayload) (string, error) {
	//

	mId := make([]string, len(notif))

	for _, v := range notif {
		// serialize payload first
		data, err := json.Marshal(v.Payload)
		if err != nil {
			return "", err
		}

		m, err := c.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamKey,
			Values: map[string]interface{}{
				"notification_token_id": v.NotificationTokenId, // reference of notification receiver
				"payload":               data,
				"device_token_id":       v.DeviceTokenId,
				"attempt":               0,
			},
		}).Result()

		if err != nil {
			slog.Error("[redis] notification producer:", err)
			break
		}
		mId = append(mId, m)
	}

	return "", nil
}
