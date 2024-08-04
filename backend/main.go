package main

import (
	"kleos/cmd/app"
	"kleos/cmd/app/clients"
	"kleos/config/logger_config"
	"kleos/config/postgres_config"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	logger_config.LoadKleosConfig(logger)
	_ = godotenv.Load()

	logger.Info("[Main] Logger initiated for Kleos")
	logger.Info(viper.GetString("postgres.dsn"))

	command := &cobra.Command{
		Use:   "Kleos",
		Short: "Reward and Recognition tool",
		Long:  "Reward and Recognition tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := gin.New()
			r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
			r.Use(ginzap.RecoveryWithZap(logger, true))

			//err := db.RunDatabaseMigrations(logger)
			//if err != nil {
			//	return err
			//}

			app.HealthCheckAPI(r)

			kleos := r.Group("/kleos")
			dbClient := postgres_config.NewGormClient(viper.GetString("postgres.dsn"), logger)

			httpClient := clients.NewHttpClient()
			mjolnirClient := clients.NewMjolnirClient(
				httpClient.HttpClient,
				viper.GetString("mjolnir.service.url"),
				viper.GetString("mjolnir.realm.id"), logger,
			)

			sv := app.NewServer(r, dbClient, logger, mjolnirClient)
			sv.Handler(kleos)
			return nil
		},
	}

	if err := command.Execute(); err != nil {
		logger.Error("[Kleos_Main] Kleos core command execution failed", zap.Error(err))
		os.Exit(1)
	}
}
