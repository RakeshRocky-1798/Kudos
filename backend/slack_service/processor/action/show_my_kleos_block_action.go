package action

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
	"kleos/slack_service/processor/view"
)

type ShowMyKleos struct {
	socketModeClient *socketmode.Client
	logger           *zap.Logger
	//kleosService     *entity.KleosEntity
}

func ShowMyKleosBlockAction(client *socketmode.Client, logger *zap.Logger) *ShowMyKleos {
	return &ShowMyKleos{
		socketModeClient: client,
		logger:           logger,
	}
}

func (sk *ShowMyKleos) ProcessAction(request *socketmode.Request, callback slack.InteractionCallback) {

	modal := view.GenerateShowMyKleosModal(callback.Channel.ID)
	_, err := sk.socketModeClient.OpenView(callback.TriggerID, modal)
	if err != nil {
		sk.logger.Error("[Show My Kleos] kleos view command failed.",
			zap.String("trigger_id", callback.TriggerID),
			zap.String("channel_id", callback.Channel.ID),
			zap.Error(err))
		return
	}
	var payload interface{}
	sk.socketModeClient.Ack(*request, payload)
}
