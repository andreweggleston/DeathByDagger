package slack

import (
	"fmt"
	"github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/helpers"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
	L      *helpers.LDAP
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

	logrus.Infof("Incoming message from userid=%s (username=%s):\n\t%s", ev.User, ev.Username, ev.Msg.Text)

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")

	if _, err := player.GetPlayerBySlackUserID(ev.User); err != nil {
		logrus.Infof("user not found in db, querying ldap for slackuid=%s", ev.User)
		if usernameEntries, err := s.L.SearchForSlackUID(ev.User); err != nil{
			return err
		} else {
			if result, err := checkWhitelist(usernameEntries[0].Attributes[0].Values[0]); result && err==nil {
				p := player.NewPlayer(usernameEntries[0].Attributes[0].Values[0])
				p.SlackUserID = ev.User
				p.Name = usernameEntries[0].Attributes[1].Values[0]
				databaseDagger.DB.Create(p)
				if s.sendMessage("You've been added to the game!", ev.Channel) != nil {
					logrus.Error("Failed to send success message")
				}
				return nil
			} else if !result {
				if s.sendMessage("You aren't on the whitelist. If you just paid the entry fee, go bug an admin to add you to the whitelist.", ev.Channel) != nil {
					logrus.Error("Failed to send whitelist message")
				}
				return nil
			} else {
				return err
			}
		}
	}

	switch m[0] {
	case "killtarget":
		user, err := player.GetPlayerBySlackUserID(ev.User)

		if err != nil {
			return err
		}
		target, err := player.GetPlayerByCSHUsername(user.Target)
		if err != nil {
			return err
		}
		target.MarkForDeath()
		channel, _, _, err := s.Client.OpenConversation(&slack.OpenConversationParameters{Users: []string{target.SlackUserID}})
		if err != nil {
			logrus.Error("Couldn't open target conversation: ", err)
			return s.sendMessage("Something went wrong when marking your target for death.", ev.Channel)
		}
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
			_ = s.sendMessage("Something went wrong when marking your target for death.", ev.Channel) //if an error exists here we're fucked
			return fmt.Errorf("failed to post interactive message: %s", err)
		}
		if _, _, err := s.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Marked <@%s> as dead. When they confirm that they've been killed, you will recieve your next target.", target.SlackUserID), false)); err != nil {
			return fmt.Errorf("failed to post kill message: %s", err)
		}

		return nil
	default:
		return s.sendMessage("You didn't send an actual command. Idiot.", ev.Channel)
	}

}

func checkWhitelist(username string) (bool, error) {
	b, err := ioutil.ReadFile("whitelist.txt")

	if err != nil {
		return false, err
	}

	file := string(b)
	return strings.Contains(file, username), nil
}

func (s *SlackListener) sendMessage(message string, ch string) error {
	_, _, err := s.Client.PostMessage(ch, slack.MsgOptionText(message, false))
	return err
}
