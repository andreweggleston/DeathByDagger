package routes

import (
	"github.com/andreweggleston/DeathByDagger/controllers"
	"net/http"
)

type route struct {
	pattern string
	handler http.HandlerFunc
}

var httpRoutes = []route{
	{"/", controllers.MainHandler},
}