package action

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
	achiementDb "kleos/db/achievement_db"
	"kleos/db/kleosDb"
	"kleos/db/usersDb"
	"kleos/slack_service/processor/view"
)

type ShowMyKleosSubmitAction struct {
	socketModeClient      *socketmode.Client
	logger                *zap.Logger
	kleosRepository       *kleosDb.KleosRepository
	achievementRepository *achiementDb.AchievementRepository
	userRepository        *usersDb.UserRepository
}

func NewShowMyKleosSubmitAction(client *socketmode.Client,
	logger *zap.Logger,
	kleosRepository *kleosDb.KleosRepository, achievementRepository *achiementDb.AchievementRepository,
	userRepository *usersDb.UserRepository) *ShowMyKleosSubmitAction {

	return &ShowMyKleosSubmitAction{
		socketModeClient:      client,
		logger:                logger,
		kleosRepository:       kleosRepository,
		achievementRepository: achievementRepository,
		userRepository:        userRepository,
	}
}

func (ssa ShowMyKleosSubmitAction) ShowMyKleosCommandProcessing(callback slack.InteractionCallback,
	request *socketmode.Request) {

	currentUser := callback.User.ID
	blocks := view.GenerateShowMyKleosResultModal(currentUser, ssa.kleosRepository, ssa.achievementRepository, ssa.userRepository)
	message := slack.MsgOptionBlocks(blocks...)
	_, err := ssa.socketModeClient.PostEphemeral(callback.View.PrivateMetadata, callback.User.ID, message)

	if err != nil {
		ssa.logger.Error("[Show kleos submit processor] error posting message to channel",
			zap.String("channel_id", callback.Channel.ID),
			zap.Error(err))
	}

	var payload interface{}
	ssa.socketModeClient.Ack(*request, payload)
}
