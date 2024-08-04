package view

import (
	"github.com/slack-go/slack"
	achiementDb "kleos/db/achievement_db"
	"kleos/web_service/utils"
	"strconv"
)

func GiveKleosModel(channel string, achievement []achiementDb.AchievementEntity) slack.ModalViewRequest {

	titleText := slack.NewTextBlockObject(slack.PlainTextType, "Kudos", false, false)
	closeText := slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false)
	submitText := slack.NewTextBlockObject(slack.PlainTextType, "Give Kudos", false, false)

	allMembersPlaceholder := slack.NewTextBlockObject(slack.PlainTextType,
		"Enter Name to search", false, false)
	allMemberTypeOption := slack.NewOptionsSelectBlockElement(slack.MultiOptTypeUser,
		allMembersPlaceholder, "member_slack_select")
	allMemberTypeText := slack.NewTextBlockObject(slack.PlainTextType,
		"Give Kudos to", false, false)
	allMemberBlock := slack.NewInputBlock("user_type",
		allMemberTypeText, nil, allMemberTypeOption)

	messageTitleText := slack.NewTextBlockObject(slack.PlainTextType,
		"Send a nice note", false, false)
	messageTitlePlaceholder := slack.NewTextBlockObject(slack.PlainTextType,
		"Describe their actions and how they impacted you or the business.", false, false)
	messageTitleElement := slack.NewPlainTextInputBlockElement(messageTitlePlaceholder, "title")
	messageTitleElement.Multiline = true
	messageTitle := slack.NewInputBlock("Description", messageTitleText, nil, messageTitleElement)

	achievementSelectionOptions := CreateBlockForAchievement(achievement)
	achievementSelectionText := slack.NewTextBlockObject(slack.PlainTextType,
		"Appreciation", false, false)
	achievementSelectionPlaceHolder := slack.NewTextBlockObject(slack.PlainTextType,
		"Select an appreciation type", false, false)
	achievementSelectionOption := slack.NewOptionsSelectBlockElement(slack.OptTypeStatic,
		achievementSelectionPlaceHolder, "level_give_kleos", achievementSelectionOptions...)
	achievementSelectionBlock := slack.NewInputBlock("level_give_kleos",
		achievementSelectionText, nil, achievementSelectionOption)

	//channelPlaceholder := slack.NewTextBlockObject(slack.PlainTextType,
	//	"Select Channel", false, false)
	//channelTypeOption := slack.NewOptionsSelectBlockElement(slack.OptTypeChannels,
	//	channelPlaceholder, "channel_slack_select")
	//channelTypeText := slack.NewTextBlockObject(slack.PlainTextType,
	//	"Channel To Post", false, false)
	//channelBlock := slack.NewInputBlock("channel_type",
	//	channelTypeText, nil, channelTypeOption)

	yesRadioOption := slack.NewOptionBlockObject("yes", slack.NewTextBlockObject(slack.MarkdownType, "Yes", false, false), nil)
	noRadioOption := slack.NewOptionBlockObject("no", slack.NewTextBlockObject(slack.MarkdownType, "No", false, false), nil)
	checkboxSelectionOption := slack.NewRadioButtonsBlockElement("radio_buttons", yesRadioOption, noRadioOption)
	checkboxSelectionOption.InitialOption = yesRadioOption
	checkboxSelectionBlock := slack.NewInputBlock(string(slack.METRadioButtons), slack.NewTextBlockObject(slack.PlainTextType, "Post in #kudos slack channel", false, false), nil, checkboxSelectionOption)

	emptyBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", " ", false, false),
		nil,
		nil,
	)

	infoSectionBlock := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", "> Kudos will also be shared to the recipient's manager", false, false),
		nil,
		nil,
	)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			allMemberBlock,
			achievementSelectionBlock,
			messageTitle,
			//channelBlock,
			emptyBlock,
			checkboxSelectionBlock,
			infoSectionBlock,
		},
	}

	return slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           titleText,
		Close:           closeText,
		Submit:          submitText,
		Blocks:          blocks,
		PrivateMetadata: channel,
		CallbackID:      string(utils.StartGiveKleos),
	}
}

func CreateCommonBlock(st []string) []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(st))
	for i, o := range st {
		localObj := o
		optionText := slack.NewTextBlockObject(slack.PlainTextType,
			localObj, false, false)
		optionBlockObjects = append(optionBlockObjects,
			slack.NewOptionBlockObject(strconv.Itoa(int(i)), optionText, nil))
	}
	return optionBlockObjects
}

func CreateBlockForAchievement(options []achiementDb.AchievementEntity) []*slack.OptionBlockObject {
	optionBlockObjects := make([]*slack.OptionBlockObject, 0, len(options))
	for _, o := range options {
		localObj := o
		optionText := slack.NewTextBlockObject(slack.PlainTextType,
			localObj.AchievementName, false, false)
		optionBlockObjects = append(optionBlockObjects,
			slack.NewOptionBlockObject(strconv.Itoa(int(localObj.ID)), optionText, nil))
	}
	return optionBlockObjects
}
