package view

import (
	"github.com/slack-go/slack"
	"kleos/web_service/utils"
)

func GenerateShowLeaderBoardModal(channel string) slack.ModalViewRequest {

	titleText := slack.NewTextBlockObject(slack.PlainTextType, "Kudos", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Send", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn", "Appreciation Champions Board\n", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	StringToShow := []string{"Kudos Received"}
	optionLeaderBoard := CreateCommonBlock(StringToShow)

	leaderBoardText := slack.NewTextBlockObject(slack.PlainTextType, "Top 10 Kudos", false, false)
	leaderBoardPlaceholder := slack.NewTextBlockObject(slack.PlainTextType, "Choose from below", false, false)
	leaderBoardTypeOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic,
		leaderBoardPlaceholder, "type", optionLeaderBoard...)
	leaderBoardBlock := slack.NewInputBlock("leader_board_type",
		leaderBoardText, nil, leaderBoardTypeOption)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			leaderBoardBlock,
		},
	}

	return slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           titleText,
		Close:           closeText,
		Submit:          submitText,
		Blocks:          blocks,
		PrivateMetadata: channel,
		CallbackID:      utils.StartshowLeaderBoard,
	}
}
