package handler

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	achiementDb "kleos/db/achievement_db"
	"kleos/db/kleosDb"
	"kleos/db/usersDb"
	"strconv"
)

type SlackHomeHandler struct {
	logger            *zap.Logger
	socketmodeClient  *socketmode.Client
	kleosEntity       *kleosDb.KleosRepository
	userEntity        *usersDb.UserRepository
	achievementEntity *achiementDb.AchievementRepository
}

func NewSlackHomeHandler(logger *zap.Logger, socketmodeClient *socketmode.Client, kleosEntity *kleosDb.KleosRepository, userEntity *usersDb.UserRepository, achievementEntity *achiementDb.AchievementRepository) *SlackHomeHandler {
	return &SlackHomeHandler{
		logger:            logger,
		socketmodeClient:  socketmodeClient,
		kleosEntity:       kleosEntity,
		userEntity:        userEntity,
		achievementEntity: achievementEntity,
	}
}

func (shh *SlackHomeHandler) KleosHomeConnect(userId string) {

	newUser, err := shh.userEntity.GetUserIdFromSlackId(userId)
	if err != nil {
		shh.logger.Error("[Home_Handler]Error in getting data from db", zap.Error(err))
		return
	}
	newUserId := fmt.Sprintf("%d", newUser.Id)

	homeView := slack.HomeTabViewRequest{
		Type: slack.VTHomeTab,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", "*Welcome to Kudos!*", false, false),
					nil,
					nil,
				),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", "Not appreciated yet, go ahead and 'Give Kudos'!", false, false),
					nil,
					nil,
				),
				slack.NewActionBlock(
					"button",
					slack.NewButtonBlockElement(
						"give_kleos",
						"giveKleosAction",
						slack.NewTextBlockObject("plain_text", "Give Kudos", true, false),
					),
				),
				//slack.NewDividerBlock(),
				//slack.NewDividerBlock(),
				//GetKleosLimitPerWeek(shh.logger),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", "*Appreciation received till now*", false, false),
					nil,
					nil,
				),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				GetDataForSection(newUserId, shh.kleosEntity, shh.logger, shh.achievementEntity),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				slack.NewSectionBlock(
					slack.NewTextBlockObject("mrkdwn", "*Here's a quick snapshot of your journey through Kudos so far...*", false, false),
					nil,
					nil,
				),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				ActivitySection(newUserId, shh.kleosEntity, shh.logger),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				GetMyWebpage(),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
				GetDashboard(),
				slack.NewDividerBlock(),
				slack.NewDividerBlock(),
			},
		},
	}

	_, err = shh.socketmodeClient.Client.PublishView(userId, homeView, "")
	if err != nil {
		shh.logger.Error("Error in publishing view", zap.Error(err))

	}
}

func GetDataForSection(userId string, kleosEntity *kleosDb.KleosRepository, logger *zap.Logger, achievementEntity *achiementDb.AchievementRepository) *slack.SectionBlock {

	userData, err := kleosEntity.GetMyKleosPerAchievement(userId)
	if err != nil {
		logger.Error("Error in getting data from db", zap.Error(err))
		return slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Error while getting your kleos", false, false),
			nil,
			nil,
		)
	}

	achievementData, err := achievementEntity.GetAllAchievementName()
	if err != nil {
		logger.Error("Error in getting data from db", zap.Error(err))
		return slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Error while getting your kleos", false, false),
			nil,
			nil,
		)
	}

	var kleosDict = make(map[string]int)
	for _, value := range *userData {
		for _, value1 := range achievementData {
			if value.Achievement == strconv.Itoa(int(value1.ID)) {
				kleosDict[value1.DisplayName] = value.Count
			}
		}
	}

	cV := ""

	for key, value := range kleosDict {
		cV = cV + fmt.Sprintf("%s - `%v`\n", key, value)
	}

	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", cV, false, false),
		nil,
		nil,
	)
}

func ActivitySection(userId string, kleosEntity *kleosDb.KleosRepository, logger *zap.Logger) *slack.SectionBlock {

	givenKleosCount, err := kleosEntity.KleosGivenCountPerUser(userId)
	if err != nil {
		logger.Error("[Home_Handler]Error in getting data from db", zap.Error(err))
		return slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "[Home_Handler]Error while getting your kleos", false, false),
			nil,
			nil,
		)
	}

	receivedKleosCount, err := kleosEntity.KleosReceivedPerUser(userId)

	if err != nil {
		logger.Error("[Home_Handler]Error in getting data from db", zap.Error(err))
		return slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "[Home_Handler]Error while getting your kleos", false, false),
			nil,
			nil,
		)
	}

	var newDict = make(map[string]int)

	newDict["No_of_kudos_given"] = givenKleosCount.Count
	newDict["No_of_kudos_received"] = receivedKleosCount.Count

	dataCm := ""

	for key, value := range newDict {
		if key == "No_of_kudos_given" {
			dataCm = dataCm + fmt.Sprintf("*%s* - `%v`\n", "Appreciation Received", value)
		} else {
			dataCm = dataCm + fmt.Sprintf("*%s* - `%v`\n", "Appreciation Given", value)
		}
	}

	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", dataCm, false, false),
		nil,
		nil,
	)
}

func GetKleosLimitPerWeek(logger *zap.Logger) *slack.SectionBlock {
	limit := viper.GetInt("limit.per.week")
	logger.Info("Max Kleos per week is : ", zap.Int("limit", limit))
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("You can give `%d` kleos per week", limit), false, false),
		nil,
		nil,
	)
}

func GetMyWebpage() *slack.SectionBlock {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "For detailed Kudos guidelines.... : https://qa-kudos-ui.np.navi-sa.in/about-kudos", false, false),
		nil,
		nil,
	)
}

func GetDashboard() *slack.SectionBlock {
	return slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "Kudos dashboard.... : https://qa-kudos-ui.np.navi-sa.in", false, false),
		nil,
		nil,
	)
}
