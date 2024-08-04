package action

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
	"kleos/slack_service/processor/view"
)

type ShowLeaderBoard struct {
	socketModeClient *socketmode.Client
	logger           *zap.Logger
	//kleosService     *entity.KleosEntity
}

func ShowLeaderBoardBlockAction(client *socketmode.Client, logger *zap.Logger) *ShowLeaderBoard {
	return &ShowLeaderBoard{
		socketModeClient: client,
		logger:           logger,
	}
}

func (sl *ShowLeaderBoard) ProcessAction(request *socketmode.Request, callback slack.InteractionCallback) {
	modal := view.GenerateShowLeaderBoardModal(callback.Channel.ID)
	_, err := sl.socketModeClient.OpenView(callback.TriggerID, modal)
	if err != nil {
		sl.logger.Error("[Show Leader Board] Leader board command failed.",
			zap.String("trigger_id", callback.TriggerID),
			zap.String("channel_id", callback.Channel.ID),
			zap.Error(err))
		return
	}
	var payload interface{}
	sl.socketModeClient.Ack(*request, payload)
}
