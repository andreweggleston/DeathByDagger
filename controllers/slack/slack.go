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
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

	//for msg := range rtm.IncomingEvents {
	//	switch ev := msg.Data.(type) {
	//	case *slack.MessageEvent:
	//		if err := s.handleMessageEvent(ev); err != nil {
	//			logrus.Error("Failed to handle message: %s", err)
	//		}
	//
	//	default:
	//		logrus.Debugf("Message event: %s", ev)
	//	}
	//}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {

	logrus.Debug("Incoming message: %v", ev)

	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.BotID)) {
		return nil
	}



	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 || m[0] != "kill" {
		return fmt.Errorf("invalid message")
	}


	if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText("Whats up doc", false)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

