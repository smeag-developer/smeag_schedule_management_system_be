package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	model "nxt_match_event_manager_api/internal/models/notification"

	hub "nxt_match_event_manager_api/internal/routes"

	"github.com/redis/go-redis/v9"
)

const (
	StreamKey      = "notifications:stream"
	GroupName      = "notification-workers"
	tokenKeyPrefix = "fcm:token:"
	streamQueueKey = "notifications:queue"
	deadLetterKey  = "notifications:dlq"
	consumerGroup  = "notification-consumer-group-workers"
	tokenTTL       = 30 * 24 * time.Hour
	MaxRetries     = 5
	ConsumerName   = "nxt_match_event_notification_consumer" // make this unique per server instance (e.g. hostname)
	StreamMaxLen   = 50_000                                  // keep last 50k entries; tune to your traffic
)

type Client struct {
	rdb *redis.Client
	hub *hub.Hub
}

func NewClient(addr, password string, h *hub.Hub) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &Client{
		rdb: rdb,
		hub: h,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// EnsureConsumerGroup creates the stream group idempotently
func (c *Client) EnsureConsumerGroup(ctx context.Context) error {
	err := c.rdb.XGroupCreateMkStream(ctx, streamQueueKey, consumerGroup, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

func (c *Client) StoreToken(ctx context.Context, userID, platform, token string) error {
	key := fmt.Sprintf("%s%s:%s", tokenKeyPrefix, userID, platform)
	return c.rdb.Set(ctx, key, token, tokenTTL).Err()
}

func (c *Client) GetToken(ctx context.Context, userID, platform string) (string, error) {
	key := fmt.Sprintf("%s%s:%s", tokenKeyPrefix, userID, platform)
	return c.rdb.Get(ctx, key).Result()
}

func (c *Client) DeleteToken(ctx context.Context, userID, platform string) error {
	key := fmt.Sprintf("%s%s:%s", tokenKeyPrefix, userID, platform)
	return c.rdb.Del(ctx, key).Err()
}

// ReadPending reads unprocessed messages from the consumer group
func (c *Client) ReadPending(ctx context.Context,
	streamKey, consumerGroup, consumerName string,
	args string, count int64) ([]redis.XStream, error) {

	streams, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{streamKey, args},
		Count:    count,
		Block:    2 * time.Second,
	}).Result()

	if err == redis.Nil {
		return nil, nil
	}

	return streams, nil
}

// ReclaimStalled re-claims messages idle longer than idleTime (for retry after backoff)
func (c *Client) ReclaimStalled(ctx context.Context, consumerName string, idleTime time.Duration) ([]redis.XMessage, error) {

	result, _, err := c.rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   streamQueueKey,
		Group:    consumerGroup,
		Consumer: consumerName,
		MinIdle:  idleTime,
		Start:    "0-0",
		Count:    10,
	}).Result()

	if err != nil {
		return nil, err
	}

	return result, nil
}

// AckMessage acknowledges successful delivery
func (c *Client) AckMessage(ctx context.Context, streamKey, consumerGroup string, msgID string) error {
	return c.rdb.XAck(ctx, streamKey, consumerGroup, msgID).Err()
}

// MoveToDeadLetter removes from stream and saves to DLQ with reason
func (c *Client) MoveToDeadLetter(ctx context.Context, msgID string, event model.NotificationEvent, reason string) error {
	data, _ := json.Marshal(event)
	pipe := c.rdb.Pipeline()
	pipe.XAdd(ctx, &redis.XAddArgs{
		Stream: deadLetterKey,
		Values: map[string]interface{}{
			"payload":    data,
			"reason":     reason,
			"failed_at":  time.Now().Unix(),
			"message_id": msgID,
		},
	})
	pipe.XAck(ctx, streamQueueKey, consumerGroup, msgID)
	_, err := pipe.Exec(ctx)
	return err
}

// IncrAttempt stores retry count per message in Redis with expiry
func (c *Client) IncrAttempt(ctx context.Context, msgID string) (int64, error) {
	key := "fcm:retry:" + msgID
	count, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	c.rdb.Expire(ctx, key, 48*time.Hour)
	return count, nil
}

// Redis HASH SET
func (c *Client) GetHashSet(ctx context.Context, key string, field string) (string, error) {

	val, err := c.rdb.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		// Field does not exist
		return "", nil
	}
	if err != nil {
		slog.Error("unable to get hash field", "error", err)
		return "", err
	}
	return val, nil
}

// Users redis hash to access field and value
func (c *Client) GetHMSETtoken(ctx context.Context, k string, field string) ([]interface{}, error) {

	d, err := c.rdb.HMGet(ctx, k, field).Result()
	if err != nil {
		slog.Error("[redis] unable to to get hash", "error", err)
		return nil, err
	}

	return d, nil
}

func (c *Client) HashFieldExists(ctx context.Context, key string, field string) (bool, error) {
	exists, err := c.rdb.HExists(ctx, key, field).Result()
	if err != nil {
		slog.Error("unable to check hash field", "error", err)
		return false, err
	}
	// Temporary: confirm what Redis actually sees
	slog.Info("HExists result", "key", key, "field", field, "exists", exists)
	return exists, nil
}

func (c *Client) UpsertHashSet(ctx context.Context, key string, field string, value interface{}, time time.Duration) (int64, error) {
	id, err := c.rdb.HSet(ctx, key, field, value).Result()
	if err != nil {
		slog.Error("unable to store hash", "error", err)
		return 0, err
	}

	// ✅ Pass `field` instead of `value`
	_, err = c.SetHExpiration(ctx, key, field, time)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Client) SetHExpiration(ctx context.Context, key string, field string, time time.Duration) ([]int64, error) {
	ex, err := c.rdb.HExpire(ctx, key, time, field).Result()
	if err != nil {
		slog.Error("unable set expiration", "error", err)
		return nil, err
	}

	return ex, nil
}

func (c *Client) XRange(ctx context.Context, streamKey string, lastId string, specialId string) ([]redis.XMessage, error) {
	msgs, err := c.rdb.XRange(ctx, StreamKey, lastId, specialId).Result()
	if err != nil {
		log.Println("xrange error:", err)
		return []redis.XMessage{}, err
	}

	return msgs, nil
}

// Uses redis set insertion and update
func (c *Client) AddSet(ctx context.Context, key string, d interface{}) (int64, error) {

	s, err := c.rdb.SAdd(ctx, key, d, 0).Result()

	if err != nil {
		slog.Error("unable to store hash", "error", err)
		return 0, err
	}

	return s, err
}

// Uses to retrieve user tokens
// format : [fcm:notification:{user_id}] -> tokens
func (c *Client) GetSetMemberToken(ctx context.Context, key string) ([]string, error) {

	ids, err := c.rdb.SMembers(ctx, key).Result()

	if err != nil {
		slog.Error("unable to store hash", "error", err)
		return nil, nil
	}

	return ids, nil
}

func (c *Client) GetStreamKey() string {
	return StreamKey
}

func (c *Client) GetConsumerName() string {
	return ConsumerName
}

func (c *Client) GetConsumerGroup() string {
	return consumerGroup
}

func (c *Client) GetGroupName() string {
	return GroupName
}
