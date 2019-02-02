package databaseDagger

import (
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/url"
	"sync"
)

var (
	DB          *gorm.DB
	dbMutex     sync.Mutex
	initialized = false
	DBUrl       url.URL
)

func Init() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if initialized {
		return
	}

	DBUrl = url.URL{
		Scheme:   "postgres",
		Host:     config.Constants.DbAddr,
		RawQuery: "sslmode=disable",
	}

	logrus.Info("Connecting to DB on ", DBUrl.String())

	DBUrl.User = url.UserPassword(config.Constants.DbUsername, config.Constants.DbPassword)

	var err error
	DB, err = gorm.Open("postgres", DBUrl.String())
	if err != nil {
		logrus.Fatal(err.Error())
	}

	DB.SetLogger(logrus.StandardLogger())

	logrus.Info("Connected!")
	initialized = true
}
