package handler

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	slackbotConfig "kleos/config/slack_config"
	"kleos/cron"
	achiementDb "kleos/db/achievement_db"
	"kleos/db/kleosDb"
	managersDb "kleos/db/managerDb"
	"kleos/db/shedlock_db"
	userCountDb "kleos/db/user_count_db"
	"kleos/db/usersDb"
	"kleos/slack_service/processor"
	"kleos/slack_service/processor/action"
	"kleos/slack_service/resolver"
)

type SlackHandler struct {
	logger                *zap.Logger
	socketmodeClient      *socketmode.Client
	slashCommandProcessor *processor.SlashCommandProcessor
	slashCommandResolver  *resolver.SlashCommandStart
	blockActionProcessor  *action.BlockActionProcessor
	submissionProcessor   *action.ViewSubmissionProcessor
	homeHandler           *SlackHomeHandler
}

const PROD_CRON = "prod_cron"

func NewSlackHandler(logger *zap.Logger, socketModeClient *socketmode.Client, gormClient *gorm.DB, hNudgeSocketMode *socketmode.Client) *SlackHandler {
	slackbotClient := slackbotConfig.NewSlackClient(logger, socketModeClient, hNudgeSocketMode)
	kleosService := kleosDb.NewKleosRepository(logger, gormClient)
	achievementService := achiementDb.NewAchievementRepository(logger, gormClient)
	userService := usersDb.NewUserRepository(logger, gormClient)
	managerService := managersDb.NewManagerRepository(logger, gormClient)
	userCountService := userCountDb.NewUserCountRepository(logger, gormClient)
	shedlockService := shedlock_db.NewShedlockRepository(gormClient, logger)
	slashCommandProcessor := processor.NewSlashCommandProcessor(logger, socketModeClient, slackbotClient)
	slashCommandResolver := resolver.NewSlashCommandResolver(logger, slashCommandProcessor)
	blockActionProcessor := action.NewBlockActionProcessor(logger, socketModeClient, slackbotClient, achievementService)
	submissionProcessor := action.NewViewSubmissionProcessor(logger, socketModeClient, slackbotClient, kleosService,
		achievementService, userService, managerService, userCountService)
	homeHandler := NewSlackHomeHandler(logger, socketModeClient, kleosService, userService, achievementService)

	go func() {
		if viper.GetString("env.check") == PROD_CRON {
			cron.RunJob(hNudgeSocketMode, gormClient, logger, kleosService, userService, managerService, achievementService, userCountService, shedlockService)
		}
		return
	}()

	return &SlackHandler{
		logger:                logger,
		socketmodeClient:      socketModeClient,
		slashCommandProcessor: slashCommandProcessor,
		slashCommandResolver:  slashCommandResolver,
		blockActionProcessor:  blockActionProcessor,
		submissionProcessor:   submissionProcessor,
		homeHandler:           homeHandler,
	}
}

func (sh *SlackHandler) KleosConnect() {
	sh.logger.Info("Going to Start Kleos Socket Mode")

	go func() {
		for evt := range sh.socketmodeClient.Events {
			switch evt.Type {

			case socketmode.EventTypeConnecting:
				{
					sh.logger.Info("connecting to slackbot with socket mode")
				}
			case socketmode.EventTypeConnectionError:
				{
					sh.logger.Error("appToken : " + viper.GetString("app.token"))
					sh.logger.Error("botToken : " + viper.GetString("bot.token"))
					sh.logger.Error("Kleos connection failed.")
				}
			case socketmode.EventTypeConnected:
				{
					sh.logger.Info("connected to slackbot with socket mode")
				}
			case socketmode.EventTypeEventsAPI:
				{
					Event := evt.Data.(slackevents.EventsAPIEvent)
					//need to find a better implementation
					innerEvent, ok := Event.InnerEvent.Data.(*slackevents.AppHomeOpenedEvent)
					if !ok {
						//sh.logger.Info("Error while getting the inner event")
						continue
					}
					eventType := innerEvent.Type
					tab := innerEvent.Tab
					userId := innerEvent.User

					if "app_home_opened" == eventType && "home" == tab {
						sh.homeHandler.KleosHomeConnect(userId)
					}
				}
			case socketmode.EventTypeInteractive:
				{
					callback, _ := evt.Data.(slack.InteractionCallback)
					switch callback.Type {
					case slack.InteractionTypeBlockActions:
						{
							sh.logger.Info("received interaction type block action", zap.String("action_id",
								callback.ActionID), zap.String("block_id", callback.BlockID))
							sh.blockActionProcessor.ProcessCommand(callback, evt.Request)
						}
					case slack.InteractionTypeViewSubmission:
						{
							sh.logger.Info("received interaction type view submission", zap.String("action_id",
								callback.ActionID), zap.String("block_id", callback.BlockID))
							sh.submissionProcessor.ProcessCommand(callback, evt.Request)
						}
					}

				}
			case socketmode.EventTypeSlashCommand:
				{
					sh.slashCommandResolver.GiveKleos(&evt)
				}
			case socketmode.EventTypeHello:
				{
					continue
				}
			default:
				{
					sh.logger.Error("UnExpected Event Type Received", zap.Any("event_type", evt.Type))
				}

			}
		}
	}()

	_ = sh.socketmodeClient.Run()
}
