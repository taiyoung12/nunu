package handler

import (
	"context"
	"strings"

	"nunu/internal/service"
	"nunu/pkg/log"
	"nunu/pkg/slackclient"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type SlackHandler struct {
	client   *slackclient.Client
	agentSvc service.AgentService
	slackSvc service.SlackService
	logger   *log.Logger
}

func NewSlackHandler(
	client *slackclient.Client,
	agentSvc service.AgentService,
	slackSvc service.SlackService,
	logger *log.Logger,
) *SlackHandler {
	return &SlackHandler{
		client:   client,
		agentSvc: agentSvc,
		slackSvc: slackSvc,
		logger:   logger,
	}
}

// Listen starts listening for Socket Mode events.
func (h *SlackHandler) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-h.client.SocketMode.Events:
			switch evt.Type {
			case socketmode.EventTypeEventsAPI:
				h.handleEventsAPI(ctx, evt)
			case socketmode.EventTypeInteractive:
				h.handleInteractive(ctx, evt)
			}
		}
	}
}

func (h *SlackHandler) handleEventsAPI(ctx context.Context, evt socketmode.Event) {
	eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
	if !ok {
		return
	}

	h.client.SocketMode.Ack(*evt.Request)

	switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		h.handleMention(ctx, ev)
	case *slackevents.MessageEvent:
		// Handle DMs
		if ev.ChannelType == "im" && ev.BotID == "" {
			h.handleDirectMessage(ctx, ev)
		}
	}
}

func (h *SlackHandler) handleInteractive(ctx context.Context, evt socketmode.Event) {
	callback, ok := evt.Data.(slack.InteractionCallback)
	if !ok {
		return
	}

	h.client.SocketMode.Ack(*evt.Request)

	// Handle reaction-based feedback
	for _, action := range callback.ActionCallback.BlockActions {
		switch action.ActionID {
		case "feedback_thumbsup":
			go h.processFeedback(ctx, callback.Channel.ID, callback.Message.Timestamp, "thumbsup")
		case "feedback_thumbsdown":
			go h.processFeedback(ctx, callback.Channel.ID, callback.Message.Timestamp, "thumbsdown")
		}
	}
}

func (h *SlackHandler) handleMention(ctx context.Context, ev *slackevents.AppMentionEvent) {
	// Strip the bot mention from the text
	question := strings.TrimSpace(ev.Text)
	// Remove <@BOT_ID> mention
	if idx := strings.Index(question, ">"); idx != -1 {
		question = strings.TrimSpace(question[idx+1:])
	}

	if question == "" {
		return
	}

	threadTS := ev.ThreadTimeStamp
	if threadTS == "" {
		threadTS = ev.TimeStamp
	}

	go h.processQuestion(ctx, ev.Channel, ev.User, threadTS, question)
}

func (h *SlackHandler) handleDirectMessage(ctx context.Context, ev *slackevents.MessageEvent) {
	if ev.Text == "" || ev.SubType != "" {
		return
	}

	threadTS := ev.ThreadTimeStamp
	if threadTS == "" {
		threadTS = ev.TimeStamp
	}

	go h.processQuestion(ctx, ev.Channel, ev.User, threadTS, ev.Text)
}

func (h *SlackHandler) processQuestion(ctx context.Context, channelID, userID, threadTS, question string) {
	h.logger.Info("processing question",
		zap.String("channel", channelID),
		zap.String("user", userID),
		zap.String("question", question),
	)

	// Send "thinking" message
	thinkingTS, err := h.slackSvc.SendThinking(channelID, threadTS)
	if err != nil {
		h.logger.Error("failed to send thinking message", zap.Error(err))
		return
	}

	// Run the agent
	answer, err := h.agentSvc.HandleQuestion(ctx, channelID, userID, threadTS, question)
	if err != nil {
		h.logger.Error("agent error", zap.Error(err))
		h.slackSvc.UpdateWithAnswer(channelID, thinkingTS, ":x: 오류가 발생했습니다. 잠시 후 다시 시도해주세요.")
		return
	}

	// Update thinking message with the answer
	if err := h.slackSvc.UpdateWithAnswer(channelID, thinkingTS, answer); err != nil {
		h.logger.Error("failed to update message", zap.Error(err))
	}
}

func (h *SlackHandler) processFeedback(ctx context.Context, channelID, threadTS, feedback string) {
	if err := h.agentSvc.HandleFeedback(ctx, channelID, threadTS, feedback); err != nil {
		h.logger.Error("failed to save feedback", zap.Error(err))
	}
}
