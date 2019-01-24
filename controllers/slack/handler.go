package slack

import (
	"encoding/json"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

var VerificationToken = config.Constants.SlackVerificatoinToken

type InteractionHandler struct {
	S *SlackListener
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

		_, _, err := h.S.Client.DeleteMessage(message.OriginalMessage.Channel, message.OriginalMessage.Timestamp)
		if err != nil {
			logrus.Errorf("Error while deleting interactive message: %s", err)
		}

		msg := ""

		switch message.ActionCallback.Actions[0].Value {
		case "confirm":
			msg = "You're dead! Sorry!"
		case "deny":
			msg = "You've denied your mark. Keep playing!"
		default:
			logrus.Warn("Recieved a response value from callback that wasn't expected")
		}

		_, _, err = h.S.Client.PostMessage(message.Channel.ID, slack.MsgOptionText(msg, false))
		if err != nil{
			logrus.Errorf("Error while posting response to interactive message: %s", err)
		}

		return
	}
}

