package processor

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
	slackbotConfig "kleos/config/slack_config"
	view "kleos/slack_service/processor/view"
)

type SlashCommandAction struct {
	logger           *zap.Logger
	socketModeClient *socketmode.Client
	slackBot         *slackbotConfig.Client
}

func NewSlasCommandAction(logger *zap.Logger,
	socketModeClient *socketmode.Client, slackBot *slackbotConfig.Client) *SlashCommandAction {

	return &SlashCommandAction{
		logger:           logger,
		socketModeClient: socketModeClient,
		slackBot:         slackBot,
	}
}

func (sca *SlashCommandAction) PerformAction(evt *socketmode.Event) {

	cmd, ok := evt.Data.(slack.SlashCommand)
	sca.logger.Info("processing kleos command", zap.Any("payload", cmd))
	if !ok {
		sca.logger.Error("event data to kleos command conversion failed", zap.Any("data", evt))
		return
	}
	var payload = view.NewKleosBlock()
	sca.socketModeClient.Ack(*evt.Request, payload)
}
