package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	cc "nxt_match_event_manager_api/internal/constants"
	"time"

	"github.com/redis/go-redis/v9"
)

// InitConsumerGroup creates the consumer group if it doesn't exist yet.
// Passing "$" as the start ID means only new messages from this point onward.
// Pass "0" to process all existing messages on first boot.

func (c *Client) InitGroupStream(ctx context.Context) error {

	// Create a consumer group (only once)
	err := c.rdb.XGroupCreateMkStream(
		ctx,
		StreamKey,
		consumerGroup, "$").Err()

	if err != nil && err.Error() != cc.BUSY_GROUP_ALREADY_EXISTS {
		return fmt.Errorf("[redis] err: %v", err)
	}

	return nil
}

func (c *Client) InitConsumerStream(ctx context.Context) error {

	// Create a consumer group (only once)
	err := c.rdb.XGroupCreateConsumer(
		ctx,
		StreamKey,
		consumerGroup,
		ConsumerName).Err()

	if err != nil && err.Error() != cc.BUSY_GROUP_ALREADY_EXISTS {
		return fmt.Errorf("[redis] err: %v", err)
	}

	return nil
}

// StartConsumeNotification loops forever, reading batches from the stream and forwarding
// messages to the right WebSocket connections via the hub.
func (c *Client) StartConsumeNotification(ctx context.Context, deviceToken string) error {

	total := 0

	for {
		res, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: ConsumerName,
			Streams:  []string{StreamKey, ">"}, // consume new notifications incoming
			Count:    10,
			Block:    5 * time.Second,
		}).Result()

		// res, err := c.ReadPending(ctx,
		// 	c.GetStreamKey(),
		// 	c.GetConsumerGroup(),
		// 	c.GetConsumerName(),
		// 	">", // args
		// 	10,  // counts
		// )

		if err != nil {
			slog.Warn("[redis] unable to read message", slog.String("err", err.Error()))
			return nil
		}

		// If no messages, break avoid infinite loop
		if len(res[0].Messages) == 0 || total == len(res[0].Messages) {
			break
		}

		for _, stream := range res {
			for _, msgs := range stream.Messages {
				total++

				if msgs.Values[cc.DEVICE_TOKEN_ID] == deviceToken {

					slog.Info("msg [redis]", msgs)
					c.HandleMessage(ctx, c.rdb, msgs)
				}
			}
		}
	}

	return nil
}

// handleMessage delivers a single stream entry to the target user's WebSocket
// connections, then ACKs the message so it won't be redelivered.
// If the user is offline, we intentionally do NOT ACK — the message stays
// in the PEL (pending entries list) and can be reclaimed later via XAUTOCLAIM.
func (c *Client) HandleMessage(ctx context.Context,
	client *redis.Client,
	msg redis.XMessage) {

	notifId, _ := msg.Values[cc.NOTIFICATION_TOKEN_ID].(string)
	if notifId == "" {
		// Malformed entry — ACK and skip so it doesn't block the PEL.
		client.XAck(ctx, StreamKey, consumerGroup, msg.ID)
		return
	}

	payload, err := json.Marshal(map[string]interface{}{
		cc.NOTIFICATION_TOKEN_ID: msg.Values[cc.NOTIFICATION_TOKEN_ID],
		cc.PAYLOAD:               msg.Values[cc.PAYLOAD],
	})

	if err != nil {
		slog.Error("marshal error:", err)
		return
	}

	delivered := c.hub.SendToUser(notifId, payload)

	slog.Info("is delivered", delivered)

	if delivered {
		// User is online and received the message — safe to ACK.
		if err := client.XAck(ctx, StreamKey, consumerGroup, msg.ID).Err(); err != nil {
			slog.Error("xack error:", err)
		}
	}
	// If not delivered (offline), message stays unACK'd.
	// Run a separate XAUTOCLAIM job to reclaim and retry old PEL entries.
}
