package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/andreweggleston/DeathByDagger/inside/version"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	flagGen = flag.Bool("genkey", false, "write a 32bit key for encrypting cookies, then exit")
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

	//initialize controller templates

	//initialize db and set max conns
	//do migrations

	httpMux := http.NewServeMux()
	//setup routes
	//do handlers

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		<-sig
		shutdown()
		os.Exit(0)
	}()

}

func shutdown() {
	logrus.Info("RECIEVED SIGINT/SIGTERM")
}