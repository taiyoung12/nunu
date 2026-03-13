package service

import (
	"fmt"
	"strings"

	"nunu/pkg/log"
	"nunu/pkg/slackclient"
)

type SlackService interface {
	SendThinking(channelID, threadTS string) (string, error)
	UpdateWithAnswer(channelID, messageTS, answer string) error
	SendError(channelID, threadTS string, err error) error
}

type slackService struct {
	client *slackclient.Client
	logger *log.Logger
}

func NewSlackService(
	client *slackclient.Client,
	logger *log.Logger,
) SlackService {
	return &slackService{
		client: client,
		logger: logger,
	}
}

func (s *slackService) SendThinking(channelID, threadTS string) (string, error) {
	_, ts, err := s.client.PostThreadReply(channelID, threadTS, ":hourglass_flowing_sand: 생각 중...")
	if err != nil {
		return "", fmt.Errorf("send thinking message: %w", err)
	}
	return ts, nil
}

func (s *slackService) UpdateWithAnswer(channelID, messageTS, answer string) error {
	// Format the answer for Slack
	formatted := formatForSlack(answer)
	return s.client.UpdateMessage(channelID, messageTS, formatted)
}

func (s *slackService) SendError(channelID, threadTS string, err error) error {
	msg := fmt.Sprintf(":x: 오류가 발생했습니다: %v", err)
	_, _, sendErr := s.client.PostThreadReply(channelID, threadTS, msg)
	return sendErr
}

func formatForSlack(answer string) string {
	// The answer from the LLM is already markdown-formatted.
	// Slack uses mrkdwn which is similar but not identical.
	// Basic conversions:
	answer = strings.ReplaceAll(answer, "```sql", "```")
	return answer
}
