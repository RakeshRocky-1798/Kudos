package resolver

import (
	"kleos/slack_service/processor"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type SlashCommandStart struct {
	logger                *zap.Logger
	slashCommandProcessor *processor.SlashCommandProcessor
}

func NewSlashCommandResolver(logger *zap.Logger,
	slashCommandProcessor *processor.SlashCommandProcessor) *SlashCommandStart {
	return &SlashCommandStart{
		logger:                logger,
		slashCommandProcessor: slashCommandProcessor,
	}
}

func (scs *SlashCommandStart) GiveKleos(evt *socketmode.Event) {

	scs.logger.Info("Inside EventTypeSlashCommand : GiveKleos")
	command, _ := evt.Data.(slack.SlashCommand)
	switch command.Command {
	default:
		scs.logger.Info("Kleos processing slash command",
			zap.String("command", command.Text),
			zap.String("channel_name", command.ChannelName),
			zap.String("user_name", command.UserName))

		scs.slashCommandProcessor.ProcessSlashCommand(evt)
	}
}
