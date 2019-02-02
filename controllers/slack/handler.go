package slack

import (
	"encoding/json"
	"fmt"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

var VerificationToken = config.Constants.SlackVerificationToken

type InteractionHandler struct {
	S *SlackListener
}

func (h *InteractionHandler) SlashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logrus.Warnf("Invalid Method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s, err := slack.SlashCommandParse(r)

	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if s.ValidateToken(config.Constants.SlackVerificationToken) {
		logrus.Warn("Verification failed during slash command handling")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	switch s.Command {
	case "/marktarget":
		p, err := player.GetPlayerBySlackUserID(s.UserID)
		if err != nil {
			if h.S.sendMessage("Either you arent in the DB or something went wrong on our end. If you think this is a mistake, contact an admin.", s.ChannelID) != nil {
				logrus.Error("Failed to send marktarget message")
			}
			return
		}
		target, err := player.GetPlayerByCSHUsername(p.Target)
		if err != nil {
			if h.S.sendMessage("Failed to mark your target for death... Something went wrong on our end.", s.ChannelID) != nil {
				logrus.Error("Failed to send marktarget message")
			}
			return
		}
		target.MarkForDeath()
		channel, _, _, err := h.S.Client.OpenConversation(&slack.OpenConversationParameters{Users: []string{target.SlackUserID}})
		if err != nil {
			logrus.Error("Couldn't open target conversation: ", err)
			if h.S.sendMessage("Something went wrong when marking your target for death.", s.ChannelID) != nil {
				logrus.Error("Couldn't send message warning of failed target convo opening.")
			}
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
		if _, _, err := h.S.Client.PostMessage(channel.ID, slack.MsgOptionText("You've been marked for death!", false), slack.MsgOptionAttachments(attachment)); err != nil {
			_ = h.S.sendMessage("Something went wrong when marking your target for death.", s.ChannelID) //if an error exists here we're fucked
			logrus.Errorf("failed to post interactive message: %s", err)
		}
		if _, _, err := h.S.Client.PostMessage(s.ChannelID, slack.MsgOptionText(fmt.Sprintf("Marked <@%s> as dead. When they confirm that they've been killed, you will recieve your next target.", target.SlackUserID), false)); err != nil {
			logrus.Errorf("failed to post kill message: %s", err)
		}
	default:
		if h.S.sendMessage("You didn't send an actual message... idiot.", s.ChannelID) != nil {
			logrus.Error("failed to send message informing of bad slash command")
		}
	}
}

func (h *InteractionHandler) InteractionHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		logrus.Warnf("Invalid Method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		logrus.Errorf("Failed to unescape request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		logrus.Errorf("Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	if message.Token != VerificationToken {
		logrus.Errorf("Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if message.CallbackID == "killConfirm" {

		_, _, err := h.S.Client.DeleteMessage(message.Channel.ID, message.OriginalMessage.Timestamp)
		if err != nil {
			logrus.Errorf("Error while deleting interactive message: %s", err)
		}

		msg := ""
		msg2 := ""
		user, err := player.GetPlayerBySlackUserID(message.User.ID)
		assassin, err2 := player.GetPlayerByTarget(user.CSHUsername)
		if err != nil {
			logrus.Error(err)
			return
		}
		if err2 != nil {
			logrus.Error(err)
			return
		}

		switch message.ActionCallback.Actions[0].Value {
		case "confirm":
			msg = "You're dead! Sorry!"
			user.ConfirmOwnMark()
			if assassin.UpdatePlayerData() != nil {
				logrus.Error("Player's evaporated from the db...")
			}
			target, _ := player.GetPlayerByCSHUsername(assassin.Target)
			msg2 = fmt.Sprintf("Your new target is <@%s>, and you now have %d kills", target.SlackUserID, assassin.Kills)

		case "deny":
			msg = "You've denied your mark. Keep playing!"
			user.DenyOwnMark()
			msg2 = "Your target claims they weren't killed. If you believe this is'nt true, contact an admin."
		default:
			logrus.Warn("Recieved a response value from callback that wasn't expected")
		}
		_, _, err = h.S.Client.PostMessage(message.Channel.ID, slack.MsgOptionText(msg, false))
		if err != nil {
			logrus.Errorf("Error while posting response to interactive message: %s", err)
		}

		channel, _, _, err := h.S.Client.OpenConversation(&slack.OpenConversationParameters{Users:[]string{assassin.SlackUserID}})
		if err != nil {
			logrus.Errorf("Error while posting response to interactive message: %s", err)
			return
		}
		_, _, err = h.S.Client.PostMessage(channel.ID, slack.MsgOptionText(msg2, false))
		if err != nil {
			logrus.Errorf("Error while posting response to interactive message: %s", err)
		}


		return
	}
}
