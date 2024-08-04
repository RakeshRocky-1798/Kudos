package app

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type KleosSlack struct {
	socketModeClient *socketmode.Client
	logger           *zap.Logger
}

func NewKleosClient(logger *zap.Logger) *KleosSlack {
	socketModeClient := SlackConnect(logger)

	return &KleosSlack{
		socketModeClient: socketModeClient,
		logger:           logger,
	}
}

func SlackConnect(logger *zap.Logger) *socketmode.Client {

	appToken := viper.GetString("app.token")
	botToken := viper.GetString("bot.token")
	
	api := slack.New(
		botToken,
		slack.OptionDebug(false),
		slack.OptionAppLevelToken(appToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
	)

	return client
}
