package view

import (
	"github.com/slack-go/slack"
	"kleos/web_service/utils"
)

func NewKleosBlock() map[string]interface{} {
	payload := map[string]interface{}{
		"blocks": []slack.Block{
			slack.NewActionBlock("give_kleos_button",
				slack.NewButtonBlockElement(
					string(utils.GiveKleos),
					"give_kleos_button_value",
					&slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "Give Kudos",
					},
				),
				slack.NewButtonBlockElement(
					string(utils.ShowMyKleos),
					"show_my_kleos_button_value",
					&slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "Show My Kudos",
					},
				),
				slack.NewButtonBlockElement(
					string(utils.ShowLeaderBoard),
					"show_leader_board_button_value",
					&slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "Appreciation Champions Board",
					},
				),
				slack.NewButtonBlockElement(
					string(utils.HelpKleos),
					"help_kleos_button_value",
					&slack.TextBlockObject{
						Type: slack.PlainTextType,
						Text: "Help",
					},
				)),
		},
	}
	return payload
}
