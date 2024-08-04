package action

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
	slackbotConfig "kleos/config/slack_config"
	achiementDb "kleos/db/achievement_db"
	"kleos/db/kleosDb"
	managersDb "kleos/db/managerDb"
	userCountDb "kleos/db/user_count_db"
	"kleos/db/usersDb"
	"kleos/web_service/utils"
)

type ViewSubmissionProcessor struct {
	logger           *zap.Logger
	socketModeClient *socketmode.Client
	giveKleosAction  *GiveKleosAction
	ShowMyKleos      *ShowMyKleosSubmitAction
	ShowLeaderBoard  *ShowLeaderBoardSubmitAction
	kleosService     *kleosDb.KleosRepository
}

func NewViewSubmissionProcessor(logger *zap.Logger, socketModeClient *socketmode.Client,
	slackbotClient *slackbotConfig.Client,
	kleosService *kleosDb.KleosRepository,
	achievementService *achiementDb.AchievementRepository,
	userService *usersDb.UserRepository, managerService *managersDb.ManagerRepository, userCountService *userCountDb.UserCountRepository) *ViewSubmissionProcessor {
	return &ViewSubmissionProcessor{
		logger:           logger,
		socketModeClient: socketModeClient,
		giveKleosAction: NewGiveKleosAction(socketModeClient, logger, slackbotClient,
			kleosService, achievementService, userService, managerService, userCountService),
		ShowMyKleos:     NewShowMyKleosSubmitAction(socketModeClient, logger, kleosService, achievementService, userService),
		ShowLeaderBoard: NewShowLeaderBoardSubmitAction(socketModeClient, logger, kleosService, userService),
		kleosService:    kleosService,
	}
}

func (v *ViewSubmissionProcessor) ProcessCommand(callback slack.InteractionCallback, request *socketmode.Request) {
	//defer func() {
	//	r := recover()
	//	if r != nil {
	//		v.logger.Error(fmt.Sprintf("[VSP] Exception occurred: %v", r.(error)))
	//	}
	//}()

	var callBackID = utils.ViewSubmissionType(callback.View.CallbackID)

	switch callBackID {
	case utils.StartGiveKleos:
		{
			v.giveKleosAction.GiveKleosModalCommandProcessing(callback, request)
		}
	case utils.StartshowMyKleos:
		{
			v.ShowMyKleos.ShowMyKleosCommandProcessing(callback, request)
		}
	case utils.StartshowLeaderBoard:
		{
			v.ShowLeaderBoard.ShowLeaderBoardCommandProcessing(callback, request)
		}
	default:
		{
			return
		}
	}
}
