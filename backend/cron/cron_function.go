package cron

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	common_func "kleos/cron/common"
	achiementDb "kleos/db/achievement_db"
	"kleos/db/kleosDb"
	managersDb "kleos/db/managerDb"
	"kleos/db/shedlock_db"
	userCountDb "kleos/db/user_count_db"
	"kleos/db/usersDb"
	"math/rand"
	"time"
)

func RunJob(socketModeClient *socketmode.Client, db *gorm.DB, logger *zap.Logger,
	kleosRepository *kleosDb.KleosRepository,
	userRepository *usersDb.UserRepository, managerRepository *managersDb.ManagerRepository,
	achievementRepository *achiementDb.AchievementRepository,
	userCountRepo *userCountDb.UserCountRepository, shedlockRepository *shedlock_db.ShedlockRepository) {

	shedlockConfig := NewLockerDbWithLockTime(100)

	err := shedlockConfig.AddFunc(viper.GetString("cron.job.name"), viper.GetString("cron.spec"), shedlockRepository, func() {
		SendCommsInSlackWeekly(socketModeClient, logger, userRepository)
	})
	if err != nil {
		logger.Error("Error while adding cron job", zap.Error(err))
	}
	shedlockConfig.Start()
}

func SendCommsInSlackWeekly(socketModeClient *socketmode.Client, logger *zap.Logger, userRepository *usersDb.UserRepository) {

	logger.Info("Running cron job to send comms in slack weekly")
	users, err := userRepository.GetUsersWithSlackId()
	if err != nil {
		logger.Error("Error while fetching users with slack id", zap.Error(err))
		return
	}
	rand.Seed(time.Now().UnixNano())
	randomColor := GenerateRandomColor()

	for _, user := range *users {
		slackId := user.SlackUserId
		if slackId == "" {
			continue
		}
		text := ":star2: Friendly Reminder: Don't forget to appreciate your peers today! A little recognition goes a long way in boosting morale and fostering positivity. Take a moment to acknowledge their hard work and support :star2:\n\n" +
			"To send a kudos, type `/kudos` in the message box of private chat or any public channel\n" +
			"To send kudos from the dashboard, click on: https://kudos.navi.com/"
		attachment := slack.Attachment{
			Text:  text,
			Color: randomColor,
		}
		_, _, err := socketModeClient.Client.PostMessage(slackId, slack.MsgOptionAttachments(attachment))
		if err != nil {
			logger.Error("Error while sending message to user", zap.Error(err))
		}
	}

}

func SendLeaderBoardData(socketModeClient *socketmode.Client, logger *zap.Logger, kleosService *kleosDb.KleosRepository, userService *usersDb.UserRepository) {

	slackChannelToPost := viper.GetString("slack.channel")

	logger.Info("Running cron job to send leader board data")
	rand.Seed(time.Now().UnixNano())

	blocks := common_func.LeaderBoardDataPublish(logger, kleosService, userService)
	message := slack.MsgOptionBlocks(blocks...)
	_, _, err := socketModeClient.PostMessage(slackChannelToPost, message)
	if err != nil {
		return
	}
}

func GenerateRandomColor() string {
	const hexCharset = "0123456789ABCDEF"
	color := "#"
	for i := 0; i < 6; i++ {
		color += string(hexCharset[rand.Intn(len(hexCharset))])
	}
	return color
}
