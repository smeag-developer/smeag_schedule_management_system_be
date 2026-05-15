package ws

import (
	"log"
	"log/slog"
	"net/http"
	cc "nxt_match_event_manager_api/internal/constants"
	models "nxt_match_event_manager_api/internal/models/config"
	redisClient "nxt_match_event_manager_api/internal/redis"
	hub "nxt_match_event_manager_api/internal/routes"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Remove invalid const block and use struct fields for these variables instead.
// These variables are already properly declared as fields in WebSocketConfig.

type WebSocketConfig struct {
	upgrader         *websocket.Upgrader
	userConn         map[string]*websocket.Conn
	multipleUserConn map[*websocket.Conn]bool
	redisClient      *redisClient.Client
	hostConf         *models.HostConfig
	hub              *hub.Hub
	mu               sync.RWMutex
}

func NewWebSocketConfig(redisClient *redisClient.Client, hostConf *models.HostConfig, h *hub.Hub) *WebSocketConfig {
	return &WebSocketConfig{
		hostConf: hostConf,
		//default value
		upgrader:         &websocket.Upgrader{},
		userConn:         make(map[string]*websocket.Conn),
		multipleUserConn: make(map[*websocket.Conn]bool),
		hub:              h,
		redisClient:      redisClient,
	}
}

func (ws *WebSocketConfig) handleUpgrader(c *gin.Context) *websocket.Upgrader {
	ws.upgrader.CheckOrigin = func(r *http.Request) bool {
		// return true
		// origin := r.Header.Get("Origin")
		// allowedOrigins := map[string]bool{
		// 	"http://localhost:8081": true,
		// 	"http://127.0.0.1:8081": true,
		// }
		// return ws.validateOrigin(c, allowedOrigins, origin)

		// for testing only allow
		return true
	}

	return ws.upgrader
}

func (ws *WebSocketConfig) validateOrigin(c *gin.Context, allowedOrigins map[string]bool, origin string) bool {

	if !allowedOrigins[origin] {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Origin not allowed"})
		return false
	}

	// Set CORS headers before upgrading
	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	return allowedOrigins[origin]
}

func (ws *WebSocketConfig) HandleBroadCastToUsers(c *gin.Context) {

	conn, err := ws.handleUpgrader(c).Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "websocket ugprade error"})
		return
	}

	defer conn.Close()

	ws.multipleUserConn[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			delete(ws.multipleUserConn, conn)
			break
		}
	}
}

func (ws *WebSocketConfig) HandleWS(c *gin.Context) {

	id := c.Query(cc.NOTIFICATION_TOKEN_ID)
	if len(id) == 0 {
		slog.Error("[ws] no id found", slog.String("notification_token_id", id))
		return
	}

	conn, err := ws.handleUpgrader(c).Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		slog.Error("[ws] upgrader", slog.String("err", err.Error()))
		return
	}

	ws.mu.Lock()
	ws.userConn[id] = conn
	ws.mu.Unlock()

	log.Printf("User %s connected via WebSocket\n", id)

	defer func() {
		ws.mu.Lock()
		delete(ws.userConn, id)
		ws.mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("User %s disconnected\n", id)
			break
		}
	}
}
