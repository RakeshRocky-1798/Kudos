package action

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	comms "kleos/cmd/app/comms"
	slackbotConfig "kleos/config/slack_config"
	achiementDb "kleos/db/achievement_db"
	"kleos/db/kleosDb"
	managersDb "kleos/db/managerDb"
	userCountDb "kleos/db/user_count_db"
	"kleos/db/usersDb"
	"kleos/slack_service/processor/view"
	"regexp"
	"strconv"
	"time"
)

type GiveKleosAction struct {
	client             *socketmode.Client
	logger             *zap.Logger
	slackbotClient     *slackbotConfig.Client
	kleosService       *kleosDb.KleosRepository
	achievementService *achiementDb.AchievementRepository
	userService        *usersDb.UserRepository
	managerService     *managersDb.ManagerRepository
	userCountService   *userCountDb.UserCountRepository
}

func NewGiveKleosAction(client *socketmode.Client, logger *zap.Logger,
	slackbotClient *slackbotConfig.Client, kleosService *kleosDb.KleosRepository,
	achievementService *achiementDb.AchievementRepository,
	userService *usersDb.UserRepository, managerService *managersDb.ManagerRepository,
	userCountService *userCountDb.UserCountRepository) *GiveKleosAction {
	return &GiveKleosAction{
		client:             client,
		logger:             logger,
		slackbotClient:     slackbotClient,
		kleosService:       kleosService,
		achievementService: achievementService,
		userService:        userService,
		managerService:     managerService,
		userCountService:   userCountService,
	}
}

func (gk *GiveKleosAction) GiveKleosModalCommandProcessing(callback slack.InteractionCallback,
	request *socketmode.Request) {

	limitPerWeekGiven := viper.GetInt("limit.per.week.g")

	userGivenInfo, err := gk.userService.GetUserIdFromSlackId(callback.User.ID)
	if err != nil {
		gk.logger.Error("Error while getting the user info", zap.Error(err))
		return
	}
	_, currentWeek := time.Now().ISOWeek()

	givenDetails, err := gk.userCountService.GetUserGiveCountCurrentWeek(strconv.Itoa(userGivenInfo.Id), fmt.Sprintf("%d", currentWeek))
	if err != nil {
		gk.logger.Error("Error while getting the given details", zap.Error(err))
		return
	}

	if givenDetails.GivenCount >= limitPerWeekGiven {
		gk.logger.Info("User has reached the limit for this week")
		msgOption := slack.MsgOptionText(fmt.Sprintf("User has reached the limit for this week"),
			false)
		_, err := gk.client.PostEphemeral(callback.View.PrivateMetadata,
			callback.User.ID, msgOption)
		if err != nil {
			gk.logger.Error("[Limit_Per_Week] Error", zap.Error(err))
			return
		}
		var payload interface{}
		gk.client.Ack(*request, payload)
		return
	}

	channelID := callback.View.PrivateMetadata
	gk.logger.Info("Sender channel id is ", zap.String("channelID", channelID))

	receiverIdList := callback.View.State.Values["user_type"]["member_slack_select"].SelectedUsers

	for _, receiverSlackId := range receiverIdList {
		gk.buildAndGiveKleos(callback, receiverSlackId)
	}
	var payload interface{}
	gk.client.Ack(*request, payload)
}

func (gk *GiveKleosAction) buildAndGiveKleos(callback slack.InteractionCallback, receiverSlackId string) {

	giveKleosRequest, requestMap := gk.buildGiveKleos(callback, receiverSlackId)

	if giveKleosRequest == nil {
		gk.logger.Error("One on One Limit Exceeded/ Same user")
		return
	}

	_, err11 := gk.kleosService.GiveKleos(giveKleosRequest)

	if err11 != nil {
		gk.logger.Error("[CIP] Error while giving Kleos", zap.Error(err11))
		return
	}

	requestMap["level_give_kleos"] = callback.View.State.Values["level_give_kleos"]["level_give_kleos"].SelectedOption.Text.Text

	senderSlackId := callback.User.ID

	givenInfo, err := gk.userService.GetUserIdFromSlackId(receiverSlackId)
	if err != nil {
		gk.logger.Error("Error while getting the given user info", zap.Error(err))
	}

	myInfo, err := gk.userService.GetUserIdFromSlackId(senderSlackId)
	if err != nil {
		gk.logger.Error("Error while getting the user info", zap.Error(err))
	}

	managerInfo, err := gk.managerService.GetManagerDataFromId(givenInfo.ManagerId)
	if err != nil {
		gk.logger.Error("Error while getting the manager info", zap.Error(err))
	}
	gk.logger.Info("Manager Info is", zap.Any("ManagerInfo", managerInfo))

	postToSlack := requestMap["radio_buttons"] == "yes"
	achievement := requestMap["level_give_kleos"]
	message := requestMap["title"]
	receiverEmail := givenInfo.Email
	receiverName := givenInfo.RealName

	senderEmail := myInfo.Email
	senderSlackImage := myInfo.SlackImageUrl
	senderName := myInfo.RealName

	managerEmail := managerInfo.Email

	gk.PostSlackCommunicationsAndTriggerEmail(senderSlackId, receiverSlackId, achievement, message, postToSlack,
		receiverEmail, receiverName, managerEmail, senderName, senderSlackImage, senderEmail)

	//go func() {
	//
	//	err := gk.postSlackBotMessage(senderId, individualUserId, requestMap["level_give_kleos"], message)
	//
	//	if err != nil {
	//		gk.logger.Error("Error will posting the give kleos summary")
	//	}
	//
	//	comms.TriggerComms(gk.logger, givenInfo.Email, managerInfo.Email, givenInfo.RealName, message, myInfo.SlackImageUrl, myInfo.RealName, myInfo.Email)
	//
	//	if postToSlack == "yes" {
	//		_, err = gk.postKleosCardToSlack(senderId, individualUserId, requestMap)
	//
	//		if err != nil {
	//			gk.logger.Error("Error will posting the give kleos summary")
	//		}
	//	}
	//}()

	return
}

func (gk *GiveKleosAction) PostSlackCommunicationsAndTriggerEmail(senderSlackId string, receiverSlackId string,
	achievement string, message string, slackFlag bool, receiverEmail string, receiverName string, managerEmail string, senderName string,
	senderSlackImage string, senderEmail string) {

	go func() {
		err := gk.postSlackBotMessage(senderSlackId, receiverSlackId, achievement, message)
		if err != nil {
			gk.logger.Error("Error while posting slack message", zap.Error(err))
		}

		comms.TriggerComms(gk.logger, receiverEmail, managerEmail, receiverName, message, senderSlackImage, senderName, senderEmail)

		if slackFlag {
			_, err := gk.postKleosCardToSlack(senderSlackId, receiverSlackId, achievement, message)
			if err != nil {
				gk.logger.Error("[Web] [Slack Channel] Error while posting slack message", zap.Error(err))
			}
		}
		return
	}()
}

//func (gk *GiveKleosAction) PostCardFromWeb(receiverId string,
//	message string,
//	achievementName string, email string, senderId string) (*string, error) {
//
//	var requestMap = make(map[string]string)
//	requestMap["level_give_kleos"] = achievementName
//	requestMap["title"] = message
//
//	color := "#CBC3E3"
//	fileUrl := selectImageForAchievement(requestMap["level_give_kleos"])
//
//	blocks := view.GiveKleosSummarySection(requestMap, senderId, receiverId)
//
//	imgAttach := slack.Attachment{
//		ImageURL:  fileUrl,
//		Title:     "Hurray",
//		TitleLink: fileUrl,
//		Color:     color,
//	}
//	_, timestamp, err := gk.client.Client.PostMessage(requestMap["channel_slack_select"],
//		slack.MsgOptionBlocks(blocks.BlockSet...),
//		slack.MsgOptionAttachments(imgAttach),
//	)
//	if err != nil {
//		gk.logger.Error("Error will posting message to the channel", zap.Error(err))
//	}
//
//	return &timestamp, nil
//}

func (gk *GiveKleosAction) postSlackBotMessage(senderSlackId string, receiverSlackId string, achievementName string, message string) error {

	senderBlock := slack.MsgOptionCompose(
		slack.MsgOptionText(
			fmt.Sprintf("Yay <@%s>, Kudos has been sent to <@%s>", senderSlackId, receiverSlackId),
			false),
		slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
			AsUser: true,
		}),
	)

	receiverBlock := view.ReceiverSlackBotSection(receiverSlackId, senderSlackId,
		achievementName, message)

	isSuccess := gk.postBotForSenderAndReceiver(senderSlackId, receiverSlackId, senderBlock, receiverBlock)

	if isSuccess {
		gk.logger.Info("Successfully posted bot messages", zap.String("senderId", senderSlackId), zap.String("receiverId", receiverSlackId), zap.Any("receiverId", time.Now()))
	}

	return nil
}

func (gk *GiveKleosAction) postBotForSenderAndReceiver(senderId string, receiverId string, senderBlock slack.MsgOption, receiverBlock slack.Attachment) bool {

	isSuccess := true

	//Sending message to the sender
	_, _, err := gk.client.Client.PostMessage(senderId, senderBlock)
	if err != nil {
		gk.logger.Error("[Slack] [Kleos Bot] Error will posting message to bot for sender", zap.Error(err))
		gk.logger.Info("Posting message to bot failed, continuing with the channel & email communication")
		isSuccess = false
	}

	//Sending message to the receiver
	_, _, err = gk.client.Client.PostMessage(receiverId, slack.MsgOptionAttachments(receiverBlock))
	if err != nil {
		gk.logger.Error("[Slack] [Kleos Bot] Error will posting message to bot receiver", zap.Error(err))
		gk.logger.Info("Posting message to bot failed, continuing with the channel & email communication")
		isSuccess = false
	}

	return isSuccess
}

func (gk *GiveKleosAction) postKleosCardToSlack(senderSlackId string, receiverSlackId string, achievement string, message string) (*string, error) {

	slackChannelToPost := viper.GetString("slack.channel")

	color := "#CBC3E3"

	fileUrl := selectImageForAchievement(achievement)

	blocks := view.GiveKleosSummarySection(achievement, message, senderSlackId, receiverSlackId)

	imgAttach := slack.Attachment{
		ImageURL:  fileUrl,
		Title:     "Hurray!",
		TitleLink: fileUrl,
		Color:     color,
	}

	gk.logger.Info("posting to slack channel: ", zap.String("slackChannelToPost", slackChannelToPost))

	_, timestamp, err := gk.client.Client.PostMessage(slackChannelToPost,
		slack.MsgOptionBlocks(blocks.BlockSet...),
		slack.MsgOptionAttachments(imgAttach),
	)

	if err != nil {
		gk.logger.Error("[Slack] [Kleos Channel] Error will posting message to the kudos channel", zap.Error(err))
	}

	return &timestamp, nil
}

func (gk *GiveKleosAction) buildGiveKleos(callback slack.InteractionCallback, individualUserId string) (*kleosDb.CreateKleosRequest, map[string]string) {

	blockActions := callback.View.State.Values
	var giveKleosRequest kleosDb.CreateKleosRequest
	var requestMap = make(map[string]string)

	for _, actions := range blockActions {
		for actionID, a := range actions {
			if string(a.Type) == string(slack.METPlainTextInput) {
				requestMap[actionID] = a.Value
			}
			if string(a.Type) == slack.OptTypeStatic {
				requestMap[actionID] = a.SelectedOption.Value
			}
			if string(a.Type) == string(slack.METRadioButtons) {
				requestMap[actionID] = a.SelectedOption.Value
			}
		}
	}

	requestMap["member_slack_select"] = individualUserId

	currentTime := time.Now()
	requestMap["year"] = fmt.Sprintf("%d", currentTime.Year())
	requestMap["month"] = fmt.Sprintf("%d", currentTime.Month())
	_, currentWeek := currentTime.ISOWeek()
	requestMap["week"] = fmt.Sprintf("%d", currentWeek)
	requestMap["day"] = fmt.Sprintf("%d", currentTime.Day())

	//userId
	givenId, err11 := gk.userService.GetUserIdFromSlackId(requestMap["member_slack_select"])
	if err11 != nil {
		gk.logger.Error("Error while getting the given user info", zap.Error(err11))
		return nil, nil
	}
	requestMap["member_slack_select"] = strconv.Itoa(givenId.Id)

	userInfo, err := gk.userService.GetUserIdFromSlackId(requestMap["member_slack_select"])
	if err != nil {
		gk.logger.Error("Error while getting the user info", zap.Error(err))
		return nil, nil
	}

	//my slackID
	userInfo, err1 := gk.userService.GetUserIdFromSlackId(callback.User.ID)
	if err1 != nil {
		gk.logger.Error("Error while getting the user info", zap.Error(err1))
		return nil, nil
	}
	requestMap["from_user"] = strconv.Itoa(userInfo.Id)

	if userInfo.Id == givenId.Id {
		gk.logger.Error("User cannot give kudos to himself")
		msgOption := slack.MsgOptionText(fmt.Sprintf("User cannot give kleos to himself"),
			false)
		_, err := gk.client.PostEphemeral(callback.View.PrivateMetadata,
			callback.User.ID, msgOption)
		if err != nil {
			gk.logger.Error("User cannot give kudos to himself", zap.Error(err))
			return nil, nil
		}
		return nil, nil
	}

	limitOneOnOne := viper.GetInt("limit.one.on.one")

	oneOnOneDetails, err := gk.kleosService.GetOneOnOneCount(requestMap["from_user"],
		requestMap["member_slack_select"], requestMap["week"])
	if err != nil {
		gk.logger.Error("Error while getting the one on one details", zap.Error(err))
		return nil, nil
	}

	if oneOnOneDetails.Count >= limitOneOnOne {
		gk.logger.Info("User one on one limit exceeded")
		msgOption := slack.MsgOptionText(fmt.Sprintf("User one on one limit exceeded for this week"),
			false)
		_, err := gk.client.PostEphemeral(callback.View.PrivateMetadata,
			callback.User.ID, msgOption)
		if err != nil {
			gk.logger.Error("[Limit_Per_Week] Error", zap.Error(err))
			return nil, nil
		}
		return nil, nil
	}

	newMap, _ := json.Marshal(requestMap)
	err = json.Unmarshal(newMap, &giveKleosRequest)
	if err != nil {
		gk.logger.Error("Error while unmarshal", zap.Error(err1))
		return nil, nil
	}

	giveToCount, err := gk.userCountService.CurrentWeekCount(requestMap["member_slack_select"], fmt.Sprintf("%d", currentWeek))
	if err != nil {
		gk.logger.Error("Error while getting the given details", zap.Error(err))
		return nil, nil
	}

	if giveToCount.UserId == "" {
		var userCountData = &userCountDb.CreateUserCountRequest{
			UserId:        requestMap["member_slack_select"],
			GivenCount:    0,
			ReceivedCount: 1,
			Month:         fmt.Sprintf("%d", currentTime.Month()),
			Week:          fmt.Sprintf("%d", currentWeek),
		}
		_, err := gk.userCountService.AddUserCountEntity(userCountData)
		if err != nil {
			gk.logger.Error("Error while adding the user count entity", zap.Error(err))
			return nil, nil
		}
	} else {
		gk.userCountService.UpdateReceivedCount(requestMap["member_slack_select"], fmt.Sprintf("%d", currentWeek))
	}

	err = gk.userService.UpdateUserReceivedCount(requestMap["member_slack_select"])
	if err != nil {
		gk.logger.Error("Error while updating the user received count", zap.Error(err))
		return nil, nil
	}

	givenFromCount, err := gk.userCountService.CurrentWeekCount(requestMap["from_user"], fmt.Sprintf("%d", currentWeek))
	if err != nil {
		gk.logger.Error("Error while getting the given details", zap.Error(err))
		return nil, nil
	}

	if givenFromCount.UserId == "" {
		var userCountData = &userCountDb.CreateUserCountRequest{
			UserId:        requestMap["from_user"],
			GivenCount:    1,
			ReceivedCount: 0,
			Month:         fmt.Sprintf("%d", currentTime.Month()),
			Week:          fmt.Sprintf("%d", currentWeek),
		}
		_, err := gk.userCountService.AddUserCountEntity(userCountData)
		if err != nil {
			gk.logger.Error("Error while adding the user count entity", zap.Error(err))
			return nil, nil
		}
	} else {
		gk.userCountService.UpdateGivenCount(requestMap["from_user"], fmt.Sprintf("%d", currentWeek))
	}

	err = gk.userService.UpdateUserGivenCount(requestMap["from_user"])
	if err != nil {
		gk.logger.Error("Error while updating the user given count", zap.Error(err))
		return nil, nil
	}

	return &giveKleosRequest, requestMap
}

func selectImageForAchievement(emojiText string) string {

	emojiPattern := `:[a-zA-Z0-9_]+:`
	emojiRegex := regexp.MustCompile(emojiPattern)
	emojis := emojiRegex.FindString(emojiText)

	switch emojis {
	case ":thumbsup:":
		return viper.GetString("thumbsup.image")
	case ":sports_medal:":
		return viper.GetString("sports.image")
	case ":rocket:":
		return viper.GetString("rocket.image")
	case ":fire:":
		return viper.GetString("fire.image")
	}
	return viper.GetString("common.image")
}
