package server

import (
	"context"

	"nunu/internal/handler"
	"nunu/pkg/log"
	"nunu/pkg/slackclient"

	"go.uber.org/zap"
)

// SlackServer manages the Slack Socket Mode connection.
type SlackServer struct {
	client       *slackclient.Client
	slackHandler *handler.SlackHandler
	logger       *log.Logger
}

func NewSlackServer(
	client *slackclient.Client,
	slackHandler *handler.SlackHandler,
	logger *log.Logger,
) *SlackServer {
	return &SlackServer{
		client:       client,
		slackHandler: slackHandler,
		logger:       logger,
	}
}

func (s *SlackServer) Start(ctx context.Context) error {
	s.logger.Info("starting Slack Socket Mode server")

	// Start event listener in a goroutine
	go s.slackHandler.Listen(ctx)

	// Run Socket Mode (blocking)
	err := s.client.SocketMode.RunContext(ctx)
	if err != nil {
		s.logger.Error("slack socket mode error", zap.Error(err))
	}
	return err
}

func (s *SlackServer) Stop(ctx context.Context) error {
	s.logger.Info("stopping Slack Socket Mode server")
	return nil
}
