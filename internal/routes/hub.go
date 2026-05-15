package router

import (
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection for a user.
type HubClient struct {
	NotificationTokenID string
	TokenDevice         string
	Conn                *websocket.Conn
	Send                chan []byte
}

// Hub tracks all active WebSocket connections, grouped by userID.
// A single user can have multiple connections (e.g. two devices).
type Hub struct {
	clients      map[string]map[*HubClient]bool
	registered   []string // user_id + sha_id
	tokenDevices []string // fcm token per device
	mu           sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*HubClient]bool),
	}
}

func (h *Hub) Register(client *HubClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[client.NotificationTokenID] == nil {
		h.clients[client.NotificationTokenID] = make(map[*HubClient]bool)
	}
	h.clients[client.NotificationTokenID][client] = true

	// all registered tokens
	h.registered = append(h.registered, client.NotificationTokenID)
	h.tokenDevices = append(h.tokenDevices, client.TokenDevice)
}

func (h *Hub) Unregister(client *HubClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[client.NotificationTokenID]; ok {
		delete(conns, client)
		if len(conns) == 0 {
			delete(h.clients, client.NotificationTokenID)

			// remove from registered slice
			_ = slices.DeleteFunc(h.registered, func(e string) bool {
				return e == client.NotificationTokenID
			})
		}
	}

	close(client.Send)
}

// SendToUser delivers a message to all active connections for the given user.
// Returns true if at least one connection received the message (user is online).
func (h *Hub) SendToUser(notifId string, message []byte) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, ok := h.clients[notifId]
	if !ok || len(conns) == 0 {
		return false // user is offline
	}

	for client := range conns {
		select {
		case client.Send <- message:
		default:
			// Channel full — drop and let the write pump clean up
		}
	}
	return true
}

func (h *Hub) RegisterClients() []string {
	return h.registered
}

func (h *Hub) GetActiveDeviceTokens() []string {
	return h.tokenDevices
}
