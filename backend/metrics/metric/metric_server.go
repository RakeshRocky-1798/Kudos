package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"kleos/config/logger_config"
)

func AdminHandler() {
	logger, _ := zap.NewProduction()
	logger_config.LoadKleosConfig(logger)
	_ = godotenv.Load()

	ginServer := gin.New()
	ginServer.GET("/actuator/metrics", gin.WrapH(promhttp.Handler()))
	go func() {
		err := ginServer.Run(":8081")
		if err != nil {
			logger.Error("Error while starting metrics", zap.Error(err))
		}
	}()
}
