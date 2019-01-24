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
				logrus.Errorf("Failed to handle message: %s", err)
			}

		}
	}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {

	if ev.SubType == "bot_message" {
		return nil
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")

	if m[0] == "test" {
		attachment := slack.Attachment{
			Text:"Were you killed?",
			CallbackID: "killConfirm",
			Actions: []slack.AttachmentAction{
				{
					Name: "confirm",
					Type: "button",
					Text: "Yes",
					Value: "confirm",
				},
				{

					Name: "deny",
					Type: "button",
					Text: "Nope",
					Value: "deny",
				},

			},
		}

		if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText("You've been marked for death!", false), slack.MsgOptionAttachments(attachment)); err != nil {
			return fmt.Errorf("failed to post message: %s", err)
		}
	}

	if len(m) != 2 || m[0] != "kill"{
		return fmt.Errorf("invalid message")
	}



	if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("you want to kill %s", m[1]), false)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	//TODO: do something with the selected users data

	return nil
}

