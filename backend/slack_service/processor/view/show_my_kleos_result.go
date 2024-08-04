package view

import (
	"fmt"
	"github.com/slack-go/slack"
	achiementDb "kleos/db/achievement_db"
	"kleos/db/kleosDb"
	"kleos/db/usersDb"
	"strconv"
)

func GenerateShowMyKleosResultModal(user string, kleosRepository *kleosDb.KleosRepository,
	achievementRepository *achiementDb.AchievementRepository,
	userRepository *usersDb.UserRepository) []slack.Block {

	userInfo, err := userRepository.GetUserIdFromSlackId(user)
	if err != nil {
		return []slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
				"Error while getting your kleos", false, false),
				nil, nil),
		}
	}

	groupValue, error1 := kleosRepository.GetMyKleosPerAchievement(strconv.Itoa(userInfo.Id))
	if error1 != nil {
		return []slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
				"Error while getting your kleos", false, false),
				nil, nil),
		}
	}

	achievementData, err := achievementRepository.GetAllAchievementName()
	if err != nil {
		return []slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
				"Error while getting your kleos", false, false),
				nil, nil),
		}
	}

	var kleosDict = make(map[string]int)
	for _, value := range *groupValue {
		for _, achievement := range achievementData {
			achId := strconv.Itoa(int(achievement.ID))
			if value.Achievement == achId {
				kleosDict[achievement.AchievementName] = value.Count
			}
		}
	}

	blocks := CreateShowMyKleos(kleosDict)

	return blocks
}

func CreateShowMyKleos(kleosDict map[string]int) []slack.Block {

	return []slack.Block{
		buildDescriptionBlock(kleosDict),
	}
}

func buildDescriptionBlock(kleosDict map[string]int) *slack.SectionBlock {
	var fields []*slack.TextBlockObject

	for key, value := range kleosDict {
		fields = append(fields, slack.NewTextBlockObject("mrkdwn",
			fmt.Sprintf("*%s* - `%v`\n", key, value), false, false))
	}

	block := slack.NewSectionBlock(nil, fields, nil)
	return block
}
