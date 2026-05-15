package core

import (
	studentHandler "smeag_sms_be/internal/handlers/student"
	models "smeag_sms_be/internal/models/config"
	studentRepository "smeag_sms_be/internal/repositories/student"
	router "smeag_sms_be/internal/routes"
	studentService "smeag_sms_be/internal/services/student"

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

	// Student routes
	studentRepo := studentRepository.NewStudentRepository(s.clientDB)
	studentSvc := studentService.NewStudentService(studentRepo)
	studentHandler.NewStudentHandlerInit(studentSvc).RegisterRoutes(routes.RoutesGroup())

	// Serve server
	NewServerInstance(s.hostConf, s.server, s.clientDB).EstablishedServer()

}
