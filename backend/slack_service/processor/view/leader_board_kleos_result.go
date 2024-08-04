package view

import (
	"fmt"
	"github.com/slack-go/slack"
	"kleos/db/kleosDb"
	"kleos/db/usersDb"
	"strconv"
)

type LeaderboardEntry struct {
	RealName string
	Count    int
}

func GenerateShowLeaderBoardResultModal(user int,
	kleosRepository *kleosDb.KleosRepository,
	userRepo *usersDb.UserRepository) []slack.Block {

	groupValue, error1 := kleosRepository.GetLeaderBoardBySlackIdUser(user)
	if error1 != nil {
		return []slack.Block{
			slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
				"Error while getting your leader", false, false), nil, nil),
		}
	}

	var leaderboardEntries []LeaderboardEntry

	for _, value := range *groupValue {
		userId, _ := strconv.Atoi(value.ReceiverID)
		userInfo, err := userRepo.GetSlackIdFromId(userId)
		if err != nil {
			return []slack.Block{
				slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
					"Error while getting your leader", false, false), nil, nil),
			}
		}
		leaderboardEntries = append(leaderboardEntries, LeaderboardEntry{
			RealName: userInfo.RealName,
			Count:    value.Count,
		})
	}

	leaderboardText := CreateShowLeaderBoard(leaderboardEntries)

	block := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn",
		leaderboardText, false, false), nil, nil)

	return []slack.Block{block}
}

func CreateShowLeaderBoard(leaderBoardDict []LeaderboardEntry) string {
	return buildDescriptionText(leaderBoardDict)
}

func buildDescriptionText(leaderBoardDict []LeaderboardEntry) string {
	var result string
	for _, value := range leaderBoardDict {
		result += fmt.Sprintf("*%s* - `%v`\n", value.RealName, value.Count)
	}

	return result
}
