package action

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
	slackbotConfig "kleos/config/slack_config"
	achiementDb "kleos/db/achievement_db"
	"kleos/web_service/utils"
)

type BlockActionProcessor struct {
	logger               *zap.Logger
	socketModeClient     *socketmode.Client
	giveKleosBlockAction *StartKleosBlockAction
	showMyKleos          *ShowMyKleos
	showLeaderBoard      *ShowLeaderBoard
	showHelp             *HelpCommandsAction
}

func NewBlockActionProcessor(logger *zap.Logger, socketModeClient *socketmode.Client,
	slackbotClient *slackbotConfig.Client,
	achievementRepository *achiementDb.AchievementRepository) *BlockActionProcessor {
	return &BlockActionProcessor{
		logger:               logger,
		socketModeClient:     socketModeClient,
		giveKleosBlockAction: GiveKleosBlockAction(socketModeClient, logger, achievementRepository),
		showMyKleos:          ShowMyKleosBlockAction(socketModeClient, logger),
		showLeaderBoard:      ShowLeaderBoardBlockAction(socketModeClient, logger),
		showHelp:             NewHelpCommandsAction(socketModeClient),
	}
}

func (bap *BlockActionProcessor) ProcessCommand(callback slack.InteractionCallback,
	request *socketmode.Request) {
	defer func() {
		if r := recover(); r != nil {
			bap.logger.Error(fmt.Sprintf("[BAP] Exception occurred: %v", r.(error)))
		}
	}()

	actionId := utils.BlockActionType(callback.ActionCallback.BlockActions[0].ActionID)
	bap.logger.Info("process button callback event",
		zap.Any("action_id", actionId), zap.String("channel", callback.Channel.Name),
		zap.String("user_id", callback.User.ID),
		zap.String("user_name", callback.User.Name))
	switch actionId {
	case utils.GiveKleos:
		{
			bap.logger.Info("Inside Give Kleos Block")
			bap.giveKleosBlockAction.ProcessAction(request, callback)
		}
	case utils.ShowMyKleos:
		{
			bap.logger.Info("Inside Show Kleos Block")
			bap.showMyKleos.ProcessAction(request, callback)
		}
	case utils.ShowLeaderBoard:
		{
			bap.logger.Info("Inside Leader Board Kleos Block")
			bap.showLeaderBoard.ProcessAction(request, callback)
		}
	case utils.HelpKleos:
		{
			bap.logger.Info("Inside Help Kleos Block")
			bap.showHelp.ProcessAction(callback.User.ID, callback.Channel.ID, request)
		}
	}
}
