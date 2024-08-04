package processor

import (
	"fmt"
	slackbotConfig "kleos/config/slack_config"

	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type SlashCommandProcessor struct {
	logger             *zap.Logger
	socketModeClient   *socketmode.Client
	slashCommandAction *SlashCommandAction
}

func NewSlashCommandProcessor(logger *zap.Logger,
	socketModeClient *socketmode.Client, slackBot *slackbotConfig.Client) *SlashCommandProcessor {
	return &SlashCommandProcessor{
		logger:             logger,
		socketModeClient:   socketModeClient,
		slashCommandAction: NewSlasCommandAction(logger, socketModeClient, slackBot),
	}
}

func (scp *SlashCommandProcessor) ProcessSlashCommand(evt *socketmode.Event) {

	defer func() {
		if r := recover(); r != nil {
			scp.logger.Error(fmt.Sprintf("[SCP] Exception occurred: %v", r.(error)))
		}
	}()

	scp.slashCommandAction.PerformAction(evt)
}
