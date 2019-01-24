package slack

import (
	"encoding/json"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

var VerificationToken = config.Constants.SlackVerificatoinToken

func InteractionHandler(w http.ResponseWriter, r *http.Request) {
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

	spew.Dump(message)
}
