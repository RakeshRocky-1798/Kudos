package action

import (
	achiementDb "kleos/db/achievement_db"
	"kleos/slack_service/processor/view"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type StartKleosBlockAction struct {
	socketModeClient   *socketmode.Client
	logger             *zap.Logger
	achievementService *achiementDb.AchievementRepository
}

func GiveKleosBlockAction(client *socketmode.Client, logger *zap.Logger,
	achievementService *achiementDb.AchievementRepository) *StartKleosBlockAction {
	return &StartKleosBlockAction{
		socketModeClient:   client,
		logger:             logger,
		achievementService: achievementService,
	}
}

func (s *StartKleosBlockAction) ProcessAction(request *socketmode.Request,
	callback slack.InteractionCallback) {

	achievements, err := s.achievementService.GetAllAchievementName()
	if err != nil || achievements == nil {
		s.logger.Error("[SIP] failed while getting all achievements")
		return
	}

	modal := view.GiveKleosModel(callback.Channel.ID, achievements)

	_, err1 := s.socketModeClient.OpenView(callback.TriggerID, modal)
	if err1 != nil {
		s.logger.Error("[SIP] houston slackbot open view command failed.",
			zap.String("trigger_id", callback.TriggerID),
			zap.String("channel_id", callback.Channel.ID), zap.Error(err1))
	}

	s.logger.Info("[SIP] kleos successfully send model to slackbot",
		zap.String("trigger_id", callback.TriggerID))
	var payload interface{}
	s.socketModeClient.Ack(*request, payload)
}
