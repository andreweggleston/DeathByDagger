package controllers

import (
	"fmt"
	"github.com/andreweggleston/DeathByDagger/controllers/controllerhelpers"
	"github.com/andreweggleston/DeathByDagger/inside/version"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"runtime"
)

var (
	mainTempl *template.Template
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var p *player.Player
	tok, err := controllerhelpers.GetToken(r)

	if err == nil {
		p = controllerhelpers.GetPlayer(tok)
	}

	errtempl := mainTempl.Execute(w, map[string]interface{}{
		"LoggedIn":  err == nil,
		"Player":    p,
		"BuildDate": version.BuildDate,
		"GitCommit": version.GitCommit,
		"GitBranch": version.GitBranch,
		"BuildInfo": fmt.Sprintf("Build using %s on %s (%s %s)", runtime.Version(), version.Hostname, runtime.GOOS, runtime.GOARCH),
	})
	if errtempl != nil {
		logrus.Error("Error!")
	}
}
