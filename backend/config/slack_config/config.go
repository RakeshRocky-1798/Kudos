package slackbotConfig

import (
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type Client struct {
	SocketModeClient *socketmode.Client
	HNudgeSMode      *socketmode.Client
	logger           *zap.Logger
}

func NewSlackClient(logger *zap.Logger, socketModeClient *socketmode.Client, hNudge *socketmode.Client) *Client {
	return &Client{
		SocketModeClient: socketModeClient,
		HNudgeSMode:      hNudge,
		logger:           logger,
	}
}
