package socket

import (
	"github.com/andreweggleston/DeathByDagger/controllers/controllerhelpers/hooks"
	"github.com/andreweggleston/DeathByDagger/controllers/socket/handler"
	"github.com/andreweggleston/DeathByDagger/inside/pprof"
	"github.com/andreweggleston/DeathByDagger/routes/socket"
	"github.com/dgrijalva/jwt-go"
)

func RegisterHandlers() {
	socket.AuthServer.OnDisconnect = hooks.OnDisconnect
	socket.UnauthServer.OnDisconnect = func(string, *jwt.Token) {pprof.Clients.Add(-1)}

	socket.AuthServer.Register(handler.Global{})
	socket.AuthServer.Register(handler.Player{})

	socket.UnauthServer.Register(handler.Unauth{})

}