package hooks

import (
	"github.com/andreweggleston/DeathByDagger/controllers/socket/sessions"
	"github.com/dgrijalva/jwt-go"
)
import	chelpers "github.com/andreweggleston/DeathByDagger/controllers/controllerhelpers"

func OnDisconnect(socketID string, token *jwt.Token){
	if token != nil{
		player := chelpers.GetPlayer(token)//getplayer w token
		if player == nil {
			return
		}
		sessions.RemoveSocket(socketID, player.CSHUsername)
	}
}