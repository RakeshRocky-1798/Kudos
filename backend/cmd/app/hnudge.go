package app

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type HNudgeSlack struct {
	socketModeClient *socketmode.Client
	logger           *zap.Logger
}

func NewHNudgeClient(logger *zap.Logger) *HNudgeSlack {

	socketModeClient := HNudgeConnect()
	return &HNudgeSlack{
		socketModeClient: socketModeClient,
		logger:           logger,
	}
}

func HNudgeConnect() *socketmode.Client {

	hAppToken := viper.GetString("happ.token")
	hBotToken := viper.GetString("hbot.token")

	api := slack.New(
		hBotToken,
		slack.OptionDebug(false),
		slack.OptionAppLevelToken(hAppToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
	)

	return client
}
