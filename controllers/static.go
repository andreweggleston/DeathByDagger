package controllers

import (
	"fmt"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/andreweggleston/DeathByDagger/inside/version"
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

	errtempl := mainTempl.Execute(w, map[string]interface{}{
		"LoggedIn": 	err == nil,
		"MockLogin": 	config.Constants.MockupAuth,
		"BuildDate": 	version.BuildDate,
		"GitCommit":	version.GitCommit,
		"GitBranch":	version.GitBranch,
		"BuildInfo": fmt.Sprintf("Build using %s on %s (%s %s)", runtime.Version(), version.Hostname, runtime.GOOS, runtime.GOARCH),
	})
	if errtempl != nil {
		logrus.Error(err)
	}
}