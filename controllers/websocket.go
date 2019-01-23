package controllers

import (
	"errors"
	"fmt"
	chelpers "github.com/andreweggleston/DeathByDagger/controllers/controllerhelpers"
	"github.com/andreweggleston/DeathByDagger/controllers/controllerhelpers/hooks"
	"github.com/andreweggleston/DeathByDagger/helpers"
	"github.com/andreweggleston/DeathByDagger/helpers/wsevent"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/andreweggleston/DeathByDagger/routes/socket"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}

func SocketHandler(w http.ResponseWriter, r *http.Request) {
	token, err := chelpers.GetToken(r)
	if err != nil && err != http.ErrNoCookie { //invalid jwt token
		token = nil
	}

	//check if player is in the whitelist
	if token == nil {
		// player isn't logged in,
		// and access is restricted to logged in people
		http.Error(w, "Not logged in", http.StatusForbidden)
		return
	}

	var so *wsevent.Client

	if token != nil { //received valid jwt
		so, err = socket.AuthServer.NewClient(upgrader, w, r)
	} else {
		so, err = socket.UnauthServer.NewClient(upgrader, w, r)
	}

	if err != nil {
		return
	}

	so.Token = token

	//logrus.Debug("Connected to Socket")
	err = SocketInit(so)
	if err != nil {
		logrus.Error(err)
		so.Close()
		return
	}
}

var ErrRecordNotFound = errors.New("Player record for found.")

//SocketInit initializes the websocket connection for the provided socket
func SocketInit(so *wsevent.Client) error {
	loggedIn := so.Token != nil

	if loggedIn {
		hooks.AfterConnect(socket.AuthServer, so)
		cshusername := so.Token.Claims.(*chelpers.DaggerClaims).CSHUsername

		player, err := player.GetPlayerByCSHUsername(cshusername)
		if err != nil {
			return fmt.Errorf("Couldn't find player record for %s", cshusername)
		}

		hooks.AfterConnectLoggedIn(so, player)
	} else {
		hooks.AfterConnect(socket.UnauthServer, so)
		so.EmitJSON(helpers.NewRequest("playerSettings", "{}"))
		so.EmitJSON(helpers.NewRequest("playerProfile", "{}"))
	}

	so.EmitJSON(helpers.NewRequest("socketInitialized", "{}"))

	return nil
}
