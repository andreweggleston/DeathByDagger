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
	Client *slack.Client
	BotID  string
}

func (s *SlackListener) ListenAndResponse() {
	rtm := s.Client.NewRTM()

	go rtm.ManageConnection()

	logrus.Info("Listening for Slack messages")


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

	logrus.Debug("Incoming message: ", ev)


	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")
	if len(m) == 0 || m[0] != "kill" {
		return fmt.Errorf("invalid message")
	}



	if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("you want to kill %s", m[1]), false)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

