package ws

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	cc "nxt_match_event_manager_api/internal/constants"
	hub "nxt_match_event_manager_api/internal/routes"
	"nxt_match_event_manager_api/internal/utils/loggers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (ws *WebSocketConfig) HandleWsNotification(c *gin.Context) {

	// --- Auth: extract notificationId and deviceToken from query ---
	notifTokenId := c.Query(cc.NOTIFICATION_TOKEN_ID)
	deviceTokenId := c.Query(cc.DEVICE_TOKEN_ID)

	if len(notifTokenId) == 0 {
		// http.Error(w, "missing userId", http.StatusBadRequest)
		loggers.GetCommonError(c, "[WS] notification token id required", http.StatusBadRequest)
		return
	}

	if len(deviceTokenId) == 0 {
		loggers.GetCommonError(c, "[WS] device token required", http.StatusBadRequest)
		return
	}

	conn, err := ws.handleUpgrader(c).Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		loggers.GetCommonError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	client := &hub.HubClient{
		NotificationTokenID: notifTokenId,
		TokenDevice:         deviceTokenId,
		Conn:                conn,
		Send:                make(chan []byte, 256),
	}

	ws.hub.Register(client)
	log.Printf("user %s connected, w/ device token %s", notifTokenId, deviceTokenId)

	// Replay any messages the client missed while offline.
	go ws.ReplayMissed(context.Background(), client)
	go writePump(client)
	ws.readPump(ws.hub, client) // blocks until connection closes
}

// readPump keeps the connection alive by handling pings/pongs.
// When it returns the connection is considered dead.
func (ws *WebSocketConfig) readPump(hub *hub.Hub, client *hub.HubClient) {
	defer func() {
		hub.Unregister(client)
		client.Conn.Close()
		log.Printf("user %s disconnected", client.NotificationTokenID)
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// We don't expect the client to send data, but we must read
		// to receive control frames (pong, close).
		if _, _, err := client.Conn.ReadMessage(); err != nil {
			break
		}
	}
}

// writePump drains client.Send and writes messages to the WebSocket.
// It also sends periodic pings to detect dead connections.
func writePump(client *hub.HubClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// replayMissed reads stream entries newer than lastID for this user and
// sends them down the WebSocket so the client catches up after reconnect.
func (ws *WebSocketConfig) ReplayMissed(ctx context.Context, client *hub.HubClient) {

	// "0" means read all PENDING messages for this consumer (unACKed)
	msgs, err := ws.redisClient.ReadPending(ctx,
		ws.redisClient.GetStreamKey(),
		ws.redisClient.GetConsumerGroup(),
		ws.redisClient.GetConsumerName(),
		"0", // args
		100, // counts
	)

	if err != nil {
		slog.Error("XReadGroup failed", "err", err)
		return
	}

	payload := make([]map[string]interface{}, 0, len(msgs))

	for _, stream := range msgs {
		for _, msg := range stream.Messages {

			if msg.Values[cc.DEVICE_TOKEN_ID] != client.TokenDevice {
				continue
			}

			m := map[string]interface{}{
				cc.NOTIFICATION_TOKEN_ID: msg.Values[cc.NOTIFICATION_TOKEN_ID],
				cc.PAYLOAD:               msg.Values[cc.PAYLOAD],
			}

			payload = append(payload, m)

			// ACK to Redis — removes from pending
			err := ws.redisClient.AckMessage(ctx,
				ws.redisClient.GetStreamKey(),
				ws.redisClient.GetConsumerGroup(),
				msg.ID)

			if err != nil {
				slog.Error("XAck failed", "msgId", msg.ID, "err", err)
				continue
			}
		}
	}

	slog.Info("payload:", len(payload))
	// send payload
	marshall, err := json.Marshal(payload)

	if err != nil {
		slog.Warn("[ws] unable to marshall")
	}

	client.Send <- marshall
}
