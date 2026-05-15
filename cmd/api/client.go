package core

import (
	"fmt"
	"log"
	"log/slog"

	models "nxt_match_event_manager_api/internal/models/config"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGrpc struct {
	RouterGroup *gin.RouterGroup
	hostConf    *models.HostConfig
}

func NewClientGrpc(rg *gin.RouterGroup, h *models.HostConfig) *ClientGrpc {
	return &ClientGrpc{
		RouterGroup: rg,
		hostConf:    h,
	}
}

func (client *ClientGrpc) EstablishClient(host string, port string, serviceType string) *grpc.ClientConn {

	address := fmt.Sprintf("%s:%s", host, port)
	slog.Info("Establing gRPC Client Connection", "service", serviceType, "address", address)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("cannot dial server : %v", err)
	}

	return conn

}
