package core

import (
	"fmt"
	"log"
	"net"
	"net/http"
	models "nxt_match_event_manager_api/internal/models/config"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	host     *models.HostConfig
	router   *gin.Engine
	clientDb mongo.Database
}

func NewServerGrpcInstance(h *models.HostConfig, r *gin.Engine, cdb *mongo.Database) *Server {
	return &Server{
		host:     h,
		router:   r,
		clientDb: *cdb,
	}
}

func (s *Server) EstablishedServer() {

	// convert port
	port, err := strconv.Atoi(s.host.Port)
	if err != nil {
		log.Fatalf("port err:%v", err)
	}

	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("cannot start server :%v", err)
	}

	m := cmux.New(listener)
	httpL := m.Match(cmux.HTTP1Fast())

	// Serve Http/1 and handle gin router engine
	httpServer := &http.Server{Handler: s.router}
	go httpServer.Serve(httpL)

	// Serve mutliplexing
	if err := m.Serve(); err != nil {
		log.Fatalf("unable to serve :%v", err)
	}

}
