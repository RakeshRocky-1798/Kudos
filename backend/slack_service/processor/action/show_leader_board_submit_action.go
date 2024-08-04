package action

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
	"kleos/db/kleosDb"
	"kleos/db/usersDb"
	"kleos/slack_service/processor/view"
)

type ShowLeaderBoardSubmitAction struct {
	socketModeClient *socketmode.Client
	logger           *zap.Logger
	kleosRepository  *kleosDb.KleosRepository
	userRepository   *usersDb.UserRepository
}

func NewShowLeaderBoardSubmitAction(client *socketmode.Client, logger *zap.Logger,
	kleosRepository *kleosDb.KleosRepository, userRepository *usersDb.UserRepository) *ShowLeaderBoardSubmitAction {
	return &ShowLeaderBoardSubmitAction{
		socketModeClient: client,
		logger:           logger,
		kleosRepository:  kleosRepository,
		userRepository:   userRepository,
	}
}

func (slbsa ShowLeaderBoardSubmitAction) ShowLeaderBoardCommandProcessing(callback slack.InteractionCallback,
	request *socketmode.Request) {

	//currentUser := callback.User.ID
	currentUser := 1
	blocks := view.GenerateShowLeaderBoardResultModal(currentUser, slbsa.kleosRepository, slbsa.userRepository)
	message := slack.MsgOptionBlocks(blocks...)
	_, err := slbsa.socketModeClient.PostEphemeral(callback.View.PrivateMetadata, callback.User.ID, message)
	if err != nil {
		slbsa.logger.Error("[Show Leader Board submit processor] error posting message to channel",
			zap.String("channel_id", callback.Channel.ID), zap.Error(err))
	}

	var payload interface{}
	slbsa.socketModeClient.Ack(*request, payload)
}
