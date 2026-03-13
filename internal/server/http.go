package server

import (
	"github.com/gin-gonic/gin"
	"nunu/internal/middleware"
	"nunu/pkg/log"
	"nunu/pkg/server/http"
	"github.com/spf13/viper"
)

func NewHTTPServer(
	conf *viper.Viper,
	logger *log.Logger,
) *http.Server {
	if conf.GetString("env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	s := http.NewServer(
		gin.Default(),
		logger,
		http.WithServerHost(conf.GetString("http.host")),
		http.WithServerPort(conf.GetInt("http.port")),
	)

	s.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
	)

	// Health check
	s.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	// Serve CSV files
	csvPath := conf.GetString("csv.storage_path")
	if csvPath == "" {
		csvPath = "./storage/csv"
	}
	s.Static("/files/csv", csvPath)

	return s
}
