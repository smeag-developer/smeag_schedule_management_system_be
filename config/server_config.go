package config

import (
	"io"
	"log/slog"
	core "nxt_match_event_manager_api/cmd/api"
	models "nxt_match_event_manager_api/internal/models/config"
	"os"

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
