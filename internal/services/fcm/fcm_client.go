package fcm

import (
	"context"
	"fmt"
	"log"
	cc "nxt_match_event_manager_api/internal/constants"
	"time"

	models "nxt_match_event_manager_api/internal/models/fcm"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	msgClient *messaging.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	opt := option.WithCredentialsFile("./" + cc.FIREBASE_ACCOUNT_JSON_FILE)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	msgClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{msgClient: msgClient}, nil
}

func (c *Client) PrepareToSend(ctx context.Context, p models.PushPayload) error {
	//prepare ttl
	ttl := time.Duration(p.TTLSeconds) * time.Second
	if p.TTLSeconds == 0 {
		ttl = 7 * 24 * time.Hour // sensible default: 7 days
	}

	switch p.Platform {
	case cc.ANDROID:
		return c.sendToAndroidUser(ctx, p, ttl)
	case cc.IOS:
		return c.sendToIosUser(ctx, p, ttl)
	default:
		return nil
	}
}

func (c *Client) sendToAndroidUser(ctx context.Context, p models.PushPayload, ttl time.Duration) error {

	msg := &messaging.Message{
		Token: p.Token,
		Notification: &messaging.Notification{
			Title: p.Title,
			Body:  p.Body,
		},
		Data: p.Data,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			TTL:      &ttl, // FCM will hold & deliver when device reconnects
			Notification: &messaging.AndroidNotification{
				Sound: "default",
				// ClickAction: "btn-match-action",
			},
		},
	}

	return c.Send(ctx, msg)
}

func (c *Client) sendToIosUser(ctx context.Context, p models.PushPayload, ttl time.Duration) error {

	msg := &messaging.Message{
		Token: p.Token,
		Notification: &messaging.Notification{
			Title: p.Title,
			Body:  p.Body,
		},
		Data: p.Data,
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-expiration": fmt.Sprintf("%d", time.Now().Add(ttl).Unix()),
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Sound: "default",
					// Badge: messaging.BadgeCount(1),
				},
			},
		},
	}

	return c.Send(ctx, msg)
}

func (c *Client) Send(ctx context.Context, msg *messaging.Message) error {
	response, err := c.msgClient.Send(ctx, msg)

	if err != nil {
		log.Printf("[FCM] Error sending message: %v", err)
		return err
	} else {
		log.Printf("[FCM] Successfully sent message: %s", response)
	}

	return nil
}
