package app

import (
	"kleos/cmd/app/clients"
	"kleos/cmd/app/handler"
	"kleos/cmd/app/middleware"
	slackbotConfig "kleos/config/slack_config"
	"kleos/metrics/metric"
	"kleos/web_service"
	"net/http"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.elastic.co/apm/module/apmgin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	gin           *gin.Engine
	db            *gorm.DB
	logger        *zap.Logger
	mjolnirClient *clients.MjolnirClient
}

func NewServer(gin *gin.Engine, db *gorm.DB, logger *zap.Logger, mjolnirClient *clients.MjolnirClient) *Server {
	return &Server{
		gin:           gin,
		db:            db,
		logger:        logger,
		mjolnirClient: mjolnirClient,
	}
}

func (s *Server) Handler(kleosGroup *gin.RouterGroup) {
	s.logger.Info("Inside Kleos Handler")
	go s.kleosHandler()
	s.Start(kleosGroup)
}

func (s *Server) createMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers",
			"Origin, Content-Type, Content-Length, Accept-Encoding, Authorization, X-Session-Token")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("X-Session-Token", "*")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
		ctx.Next()
	}
}

func (s *Server) kleosHandler() {
	s.logger.Info("Going to Start Kleos Handler Socket Mode")
	kleosClient := NewKleosClient(s.logger)
	hNudgeClient := NewHNudgeClient(s.logger)
	slackClient := handler.NewSlackHandler(s.logger, kleosClient.socketModeClient, s.db, hNudgeClient.socketModeClient)
	slackClient.KleosConnect()
}

func (s *Server) Start(kleosGroup *gin.RouterGroup) {

	s.gin.Use(apmgin.Middleware(s.gin))
	s.gin.Use(ginzap.RecoveryWithZap(s.logger, true))
	s.gin.Use(s.createMiddleWare())
	kleosGroup.Use(s.createMiddleWare())
	s.kleosClientHandler(kleosGroup)
	s.logger.Info("starting kleos server", zap.String("port", viper.GetString("port")))
	err := s.gin.Run(":8080")
	if err != nil {
		s.logger.Error("Failed to start server", zap.Error(err))
		return
	}
}

func (s *Server) kleosClientHandler(kleosGroup *gin.RouterGroup) {
	kleosClient := NewKleosClient(s.logger)
	hNudgeClient := NewHNudgeClient(s.logger)
	slackbotClient := slackbotConfig.NewSlackClient(s.logger, kleosClient.socketModeClient, hNudgeClient.socketModeClient)
	kleosHandler := web_service.NewKleosService(s.gin, s.logger, s.db, kleosClient.socketModeClient, slackbotClient, hNudgeClient.socketModeClient)

	kleosGroup.Use(middleware.MetricMiddleware())
	kleosGroup.Use(middleware.AuthMiddleware(s.mjolnirClient))

	kleosGroup.POST("/giveKleos", kleosHandler.GiveKleosFromWeb)
	kleosGroup.POST("/createUser", kleosHandler.CreateUser)
	kleosGroup.POST("/createUserBulk", kleosHandler.CreateUserBulk)
	kleosGroup.GET("/getAllUsers", kleosHandler.GetAllUsers)
	kleosGroup.GET("/leaderboard", kleosHandler.GetLeaderBoard)
	kleosGroup.GET("/achievement", kleosHandler.GetAllAchievement)
	kleosGroup.GET("/getKleos/:userid", kleosHandler.GetKleosDashboard)
	kleosGroup.GET("/getPaginatedInfo/:userid", kleosHandler.GetPaginatedInfo)
	kleosGroup.GET("/getKleosReceived/:userid", kleosHandler.GetKleosReceived)

	//admin routes
	kleosGroup.GET("/getAdminData", kleosHandler.NewAdminData)
	kleosGroup.GET("/getAdminData/csv/:userid", kleosHandler.GetAdminDataCsv)
	kleosGroup.GET("/admin/triggerComms", kleosHandler.TriggerLeaderboardData)
	kleosGroup.GET("/admin/triggerCommsWeekly", kleosHandler.TriggerWeeklyComms)
	kleosGroup.GET("/getAdminData/xls/:userid", kleosHandler.GetAdminXlsxFile)

	metric.AdminHandler()
}

func HealthCheckAPI(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})
}
