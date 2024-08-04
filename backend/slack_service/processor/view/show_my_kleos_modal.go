package view

import (
	"github.com/slack-go/slack"
	"kleos/web_service/utils"
)

func GenerateShowMyKleosModal(channel string) slack.ModalViewRequest {

	titleText := slack.NewTextBlockObject(slack.PlainTextType, "Kudos", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Send", false, false)

	headerText := slack.NewTextBlockObject("mrkdwn",
		"Show My Kudos", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	StringToShow := []string{"Appreciation received till now"}
	optionMyKleos := CreateCommonBlock(StringToShow)

	myKleosText := slack.NewTextBlockObject(slack.PlainTextType,
		"My Kudos", false, false)
	myKleosPlaceholder := slack.NewTextBlockObject(slack.PlainTextType,
		"Choose from below", false, false)
	myKleosTypeOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic,
		myKleosPlaceholder, "type", optionMyKleos...)
	myKleosBlock := slack.NewInputBlock("my_kleos_type", myKleosText, nil,
		myKleosTypeOption)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			headerSection,
			myKleosBlock,
		},
	}

	return slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           titleText,
		Close:           closeText,
		Submit:          submitText,
		Blocks:          blocks,
		PrivateMetadata: channel,
		CallbackID:      utils.StartshowMyKleos,
	}
}
