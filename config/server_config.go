package config

import (
	"io"
	"log/slog"
	"os"
	core "smeag_sms_be/cmd/api"
	models "smeag_sms_be/internal/models/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ServerConfig(dbConf *models.DBconfig, hostConf *models.HostConfig, clientDb mongo.Client) {

	gin.ForceConsoleColor() // add additonal log colorS
	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	ginServer := gin.Default()
	db := clientDb.Database(dbConf.DBName)

	// Middlewares
	ShowMascot()
	slog.Info("App Running", "addr", hostConf.Host)
	slog.Info("Running in --" + dbConf.BuildEnv + " mode & setup config")
	slog.Info("MongoDB connection established", "addr", dbConf.DBHost)
	slog.Info("Starting Server", "port", hostConf.Port)
	slog.Info("Connected to DB", "db", dbConf.DBName)

	// Run Server
	core.NewAPIServer(db, ginServer, dbConf, hostConf).Run()

}
