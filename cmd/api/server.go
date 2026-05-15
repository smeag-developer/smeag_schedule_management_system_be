package core

import (
	"fmt"
	"log"
	models "smeag_sms_be/internal/models/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServerInstance struct {
	host     *models.HostConfig
	router   *gin.Engine
	clientDb mongo.Database
}

func NewServerInstance(h *models.HostConfig, r *gin.Engine, cdb *mongo.Database) *ServerInstance {
	return &ServerInstance{
		host:     h,
		router:   r,
		clientDb: *cdb,
	}
}

func (s *ServerInstance) EstablishedServer() {
	address := fmt.Sprintf(":%s", s.host.Port)

	if err := s.router.Run(address); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
