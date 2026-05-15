package core

import (
	models "nxt_match_event_manager_api/internal/models/config"
	router "nxt_match_event_manager_api/internal/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type APISserver struct {
	clientDB *mongo.Database
	server   *gin.Engine
	dbConf   *models.DBconfig
	hostConf *models.HostConfig
	ctx      gin.Context
}

func NewAPIServer(clientDb *mongo.Database, s *gin.Engine, c *models.DBconfig, h *models.HostConfig) *APISserver {
	return &APISserver{
		clientDB: clientDb,
		server:   s,
		dbConf:   c,
		hostConf: h,
	}
}

func (s *APISserver) Run() {

	// Define Routes Config
	routes := router.NewRouterConfig(s.server, s.hostConf.AllowedOrigins)
	routes.InitConfig()

	//Routes
	// notifHandler.RegisterRoutes(routes.RoutesGroup())

	// Create multiplexing server
	NewServerGrpcInstance(s.hostConf, s.server, s.clientDB).EstablishedServer()

}
