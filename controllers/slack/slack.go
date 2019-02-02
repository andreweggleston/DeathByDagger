package slack

import (
	"github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/helpers"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/davecgh/go-spew/spew"
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

	spew.Dump(ev)

	if ev.SubType == "bot_message" {
		return nil
	}

	logrus.Infof("Incoming message from userid=%s (username=%s):\t%s", ev.User, ev.Username, ev.Msg.Text)

	if _, err := player.GetPlayerBySlackUserID(ev.User); err != nil {
		logrus.Infof("user not found in db, querying ldap for slackuid=%s", ev.User)
		if usernameEntries, err := s.L.SearchForSlackUID(ev.User); err != nil{
			if s.sendMessage("Couldn't find you in the ldap DB. Have you linked your slack to ldap with http://eac.csh.rit.edu ?", ev.Channel) != nil {
				logrus.Error("Failed to send eac message")
			}
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
	return nil
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
