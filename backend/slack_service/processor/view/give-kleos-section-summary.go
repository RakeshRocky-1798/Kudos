package view

import (
	"fmt"
	"github.com/slack-go/slack"
	"math/rand"
	"time"
)

func GiveKleosSummarySection(achievement string, message string, senderId string, receiverId string) slack.Blocks {

	messageGiven := message
	levelGiven := achievement

	return slack.Blocks{
		BlockSet: []slack.Block{
			buildUserTagSection(receiverId, senderId),
			buildGiveKleosSection(levelGiven, messageGiven),
		},
	}
}

func ReceiverSlackBotSection(receiverId string, senderId string,
	levelGiven string, messageGiven string) slack.Attachment {

	rand.Seed(time.Now().UnixNano())
	randomColor := GenerateRandomColor()
	text := fmt.Sprintf("Yay <@%s>, You have revieved a kudos from <@%s>\n\n*Appreciation:* %s\n*Message:* %s", receiverId, senderId, levelGiven, messageGiven)

	attachment := slack.Attachment{
		Text:  text,
		Color: randomColor,
	}
	return attachment
}

func buildUserTagSection(receiverId string, senderId string) *slack.SectionBlock {

	textBlock := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf(":tada: WOW! <@%s>, you have received an appreciation from <@%s>!", receiverId, senderId),
		false, false)
	sectionBlock := slack.NewSectionBlock(textBlock, nil, nil)
	return sectionBlock
}

func buildGiveKleosSection(levelGiven string, messageGiven string) *slack.SectionBlock {

	textBlock := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("\n*Appreciation:* %s\n*Message:* %s", levelGiven, messageGiven),
		false, false)
	sectionBlock := slack.NewSectionBlock(textBlock, nil, nil)
	return sectionBlock
}

func GenerateRandomColor() string {
	const hexCharset = "0123456789ABCDEF"
	color := "#"
	for i := 0; i < 6; i++ {
		color += string(hexCharset[rand.Intn(len(hexCharset))])
	}
	return color
}
