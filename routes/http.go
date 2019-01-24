package routes

import (
	"github.com/andreweggleston/DeathByDagger/controllers"
	"github.com/andreweggleston/DeathByDagger/controllers/login"
	"net/http"
)

type route struct {
	pattern string
	handler http.HandlerFunc
}

var httpRoutes = []route{
	{"/", controllers.MainHandler},
	{"/websocket/", controllers.SocketHandler},

	{"/login", login.LoginHandler},
	{"/callback", login.CallbackHandler},
	{"/logout", login.LogoutHandler},
}

func SetupHTTP(mux *http.ServeMux) {
	for _, httpRoute := range httpRoutes {
		mux.HandleFunc(httpRoute.pattern, httpRoute.handler)
	}
}