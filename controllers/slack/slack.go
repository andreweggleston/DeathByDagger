package slack

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

type SlackListener struct {
	client    *slack.Client
	botID     string
}

func (s *SlackListener) ListenAndResponse() {
	rtm := s.client.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				logrus.Error("Failed to handle message: %s", err)
			}
		}
	}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 || m[0] != "kill" {
		return fmt.Errorf("invalid message")
	}


	if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText("Whats up doc", false)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

