package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/andreweggleston/DeathByDagger/controllers"
	"github.com/nlopes/slack"
	slackhelper "github.com/andreweggleston/DeathByDagger/controllers/slack"
	"github.com/andreweggleston/DeathByDagger/controllers/socket"
	"github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/databaseDagger/migrations"
	"github.com/andreweggleston/DeathByDagger/helpers"
	"github.com/andreweggleston/DeathByDagger/inside/version"
	"github.com/andreweggleston/DeathByDagger/routes"
	socketServer "github.com/andreweggleston/DeathByDagger/routes/socket"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	flagGen = flag.Bool("genkey", false, "write a 32bit key for encrypting cookies, then exit")
	dbMaxopen = flag.Int("db-maxopen", 80, "maximum number of open database connections allowed.")
)

func main() {
	flag.Parse()

	if *flagGen {
		key := make([]byte, 64)
		_, err := rand.Read(key)
		if err != nil {
			logrus.Fatal(err)
		}

		base64Key := base64.StdEncoding.EncodeToString(key)
		fmt.Println(base64Key)
		return
	}

	logrus.Debug("Commit: ", version.GitCommit)
	logrus.Debug("Branch: ", version.GitBranch)
	logrus.Debug("Build date: ", version.BuildDate)

	controllers.InitTemplates()

	databaseDagger.Init()
	databaseDagger.DB.DB().SetMaxOpenConns(*dbMaxopen)
	migrations.Do()

	httpMux := http.NewServeMux()
	routes.SetupHTTP(httpMux)
	socket.RegisterHandlers()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:		config.Constants.AllowedOrigins,
		AllowedMethods:		[]string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowCredentials:	true,
	}).Handler(httpMux)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		<-sig
		shutdown()
		os.Exit(0)
	}()

	logrus.Info("Serving on ", config.Constants.ListenAddress)
	logrus.Info("Hosting on ", config.Constants.PublicAddress)

	client := slack.New(config.Constants.SlackToken)

	slackListener := &slackhelper.SlackListener{
		client,
		config.Constants.SlackBotID,
	}

	go slackListener.ListenAndResponse()

	logrus.Fatal(http.ListenAndServe(config.Constants.ListenAddress, corsHandler))

}

func shutdown() {
	logrus.Info("RECIEVED SIGINT/SIGTERM")
	logrus.Info("Waiting for GlobalWait")
	helpers.GlobalWait.Wait()
	logrus.Info("waiting for socket requests to complete.")
	socketServer.Wait()
	logrus.Info("closing all active websocket connections")
	socketServer.AuthServer.Close()
}