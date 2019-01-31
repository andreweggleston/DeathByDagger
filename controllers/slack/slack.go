package slack

import (
	"fmt"
	"github.com/andreweggleston/DeathByDagger/helpers"
	"github.com/andreweggleston/DeathByDagger/models/player"
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
	L	   *helpers.LDAP
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

	if _, err := player.GetPlayerBySlackUserID(ev.User); err != nil {
		//todo: initiate ldap search
	}

	if m[0] == "setusername" && len(m) == 2 {
		p, err := player.GetPlayerBySlackUserID(ev.User)
		if p != nil {
			if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("You've already set your username. If you typed it wrong, contact an admin."), false)); err != nil {
				return fmt.Errorf("failed to post username message: %s", err)
			}
			return nil
		}

		user, err := player.GetPlayerByCSHUsername(m[1])

		if err != nil {
			if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("That CSH username doesn't exist in our database. Make sure you log in first."), false)); err != nil {
				return fmt.Errorf("failed to post username message: %s", err)
			}
			return err
		}

		if user.SlackUserID != "" {
			if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("That user has already set their slack username. Stop trying to impersonate people!"), false)); err != nil {
				return fmt.Errorf("failed to post username message: %s", err)
			}
			return nil
		}
		if user.SlackUserID == "" {
			user.SetSlackUserID(ev.User)
			if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Set your slack username successfully!"), false)); err != nil {
				return fmt.Errorf("failed to post username message: %s", err)
			}
		}
		if err := user.Save(); err != nil {
			return fmt.Errorf("failed to save user's slackusername")
		}

		return nil
	}

	if len(m) != 1 || m[0] != "killtarget" {
		return fmt.Errorf("invalid message")
	}

	user, err := player.GetPlayerBySlackUserID(ev.User)

	if err != nil {
		return err
	}
	target, err := player.GetPlayerByCSHUsername(user.Target)
	if err != nil {
		return err
	}
	target.MarkForDeath()
	channel, _, _, err := s.Client.OpenConversation(&slack.OpenConversationParameters{Users:[]string{target.SlackUserID}})
	if err == nil {
		attachment := slack.Attachment{
			Text:       "Were you killed?",
			CallbackID: "killConfirm",
			Actions: []slack.AttachmentAction{
				{
					Name:  "confirm",
					Type:  "button",
					Text:  "Yes",
					Value: "confirm",
				},
				{

					Name:  "deny",
					Type:  "button",
					Text:  "Nope",
					Value: "deny",
				},
			},
		}

		if _, _, err := s.Client.PostMessage(channel.ID, slack.MsgOptionText("You've been marked for death!", false), slack.MsgOptionAttachments(attachment)); err != nil {
			return fmt.Errorf("failed to post interactive message: %s", err)
		}
	} else {
		logrus.Error("Couldn't open conversation? ", err)
	}
	if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Marked <@%s> as dead. When they confirm that they've been killed, you will recieve your next target.", target.SlackUserID), false)); err != nil {
		return fmt.Errorf("failed to post kill message: %s", err)
	}

	return nil
}

func (s *SlackListener) sendMessage(message string, ch *slack.Channel) error {
	_, _, err := s.Client.PostMessage(ch.ID, slack.MsgOptionText(message, false))
	return err
}
