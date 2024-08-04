package logger_config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func LoadKleosConfig(logger *zap.Logger) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("kleos-config")
	viper.SetConfigType("properties")
	viper.SetConfigFile("./config/logger_config/application.properties")
	//viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("[Logger]Error while loading kleos configuration", zap.Error(err))
		os.Exit(1)
	}
}
