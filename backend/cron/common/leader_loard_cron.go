package common_func

import (
	"fmt"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"kleos/db/kleosDb"
	"kleos/db/usersDb"
	"kleos/slack_service/processor/view"
	"strconv"
)

func LeaderBoardDataPublish(logger *zap.Logger,
	kleosService *kleosDb.KleosRepository,
	userService *usersDb.UserRepository) []slack.Block {

	groupValue, error1 := kleosService.GetLeaderBoardBySlackIdUser(0)
	if error1 != nil {
		logger.Error("Error while getting your leader", zap.Error(error1))
		return []slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "Error while getting your leader", false, false), nil, nil),
		}
	}

	var leaderboardEntries []view.LeaderboardEntry

	for _, value := range *groupValue {
		userId, _ := strconv.Atoi(value.ReceiverID)
		userInfo, err := userService.GetSlackIdFromId(userId)
		if err != nil {
			logger.Error("Error while getting your leader", zap.Error(err))
			return []slack.Block{
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "Error while getting your leader", false, false), nil, nil),
			}
		}

		leaderboardEntries = append(leaderboardEntries, view.LeaderboardEntry{
			RealName: userInfo.RealName,
			Count:    value.Count,
		})
	}

	leaderboardText := CreateShowLeaderBoard(leaderboardEntries)

	block := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		leaderboardText, false, false), nil, nil)

	return []slack.Block{
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", ">*Leaderboard data for the current week*", false, false), nil, nil),
		block,
	}
}

func CreateShowLeaderBoard(leaderBoardDict []view.LeaderboardEntry) string {
	return buildDescriptionText(leaderBoardDict)
}

func buildDescriptionText(leaderBoardDict []view.LeaderboardEntry) string {
	var result string
	for _, value := range leaderBoardDict {
		result += fmt.Sprintf("*%s* - `%v`\n", value.RealName, value.Count)
	}

	return result
}
