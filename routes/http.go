package routes

import (
	"github.com/andreweggleston/DeathByDagger/controllers"
	"github.com/andreweggleston/DeathByDagger/controllers/login"
	slackhelper "github.com/andreweggleston/DeathByDagger/controllers/slack"
	"net/http"
)

type route struct {
	pattern string
	handler http.HandlerFunc
}

var InteractionHandler slackhelper.InteractionHandler

var httpRoutes = []route{
	{"/", controllers.MainHandler},
	{"/websocket/", controllers.SocketHandler},

	{"/login", login.LoginHandler},
	{"/callback", login.CallbackHandler},
	{"/logout", login.LogoutHandler},

	{"/slackinteraction", InteractionHandler.InteractionHandler},
}

func SetupHTTP(mux *http.ServeMux, listener *slackhelper.SlackListener) {
	InteractionHandler.S = listener
	for _, httpRoute := range httpRoutes {
		mux.HandleFunc(httpRoute.pattern, httpRoute.handler)
	}
}
