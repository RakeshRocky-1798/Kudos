package action

import (
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"strings"
)

type HelpCommandsAction struct {
	socketModeClient *socketmode.Client
}

func NewHelpCommandsAction(client *socketmode.Client) *HelpCommandsAction {
	return &HelpCommandsAction{
		socketModeClient: client,
	}
}

func (action *HelpCommandsAction) ProcessAction(userID, channelID string, request *socketmode.Request) {

	webpageUrl := viper.GetString("dashboard.url")
	guideLineUrl := viper.GetString("guideline.url")

	webpageText := "For more details, please visit our webpage: " + string(webpageUrl)
	guideLineText := "For Kudos guidelines, please visit: " + string(guideLineUrl)

	var formattedHelpMessage string

	formattedHelpMessage += fmt.Sprintf("%s\n%s", strings.TrimSpace(webpageText), strings.TrimSpace(guideLineText))

	_, err := action.socketModeClient.Client.PostEphemeral(channelID, userID, slack.MsgOptionText(fmt.Sprintf("```%s```", formattedHelpMessage), false))
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to post help command response. %+v", err))
	}
	var payload interface{}
	action.socketModeClient.Ack(*request, payload)
}
