package slackclient

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Client wraps the Slack API and Socket Mode clients.
type Client struct {
	API        *slack.Client
	SocketMode *socketmode.Client
	Logger     *zap.Logger
}

func NewClient(conf *viper.Viper, logger *zap.Logger) *Client {
	botToken := conf.GetString("slack.bot_token")
	appToken := conf.GetString("slack.app_token")

	api := slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
	)

	sm := socketmode.New(
		api,
		socketmode.OptionLog(newZapSlackLogger(logger)),
	)

	return &Client{
		API:        api,
		SocketMode: sm,
		Logger:     logger,
	}
}

// PostMessage sends a message to a channel.
func (c *Client) PostMessage(channelID, text string, options ...slack.MsgOption) (string, string, error) {
	opts := append([]slack.MsgOption{slack.MsgOptionText(text, false)}, options...)
	ch, ts, err := c.API.PostMessage(channelID, opts...)
	if err != nil {
		return "", "", fmt.Errorf("slack post message error: %w", err)
	}
	return ch, ts, nil
}

// UpdateMessage updates an existing message.
func (c *Client) UpdateMessage(channelID, timestamp, text string, options ...slack.MsgOption) error {
	opts := append([]slack.MsgOption{slack.MsgOptionText(text, false)}, options...)
	_, _, _, err := c.API.UpdateMessage(channelID, timestamp, opts...)
	if err != nil {
		return fmt.Errorf("slack update message error: %w", err)
	}
	return nil
}

// PostThreadReply sends a threaded reply.
func (c *Client) PostThreadReply(channelID, threadTS, text string) (string, string, error) {
	return c.PostMessage(channelID, text, slack.MsgOptionTS(threadTS))
}

// zapSlackLogger adapts zap.Logger to slack's logger interface.
type zapSlackLogger struct {
	logger *zap.Logger
}

func newZapSlackLogger(logger *zap.Logger) *zapSlackLogger {
	return &zapSlackLogger{logger: logger.Named("slack")}
}

func (l *zapSlackLogger) Output(calldepth int, s string) error {
	l.logger.Debug(s)
	return nil
}
