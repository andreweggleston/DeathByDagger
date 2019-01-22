package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type constants struct {
	ListenAddress     string   `envconfig:"SA_SERVER_ADDR" default:"0.0.0.0:8081" doc:"Address to serve on"`
	AllowedOrigins    []string `envconfig:"ALLOWED_ORIGINS" default:"*"`
	DbAddr            string   `envconfig:"DATABASE_ADDR" doc:"Database Address"`
	DbDatabase        string   `envconfig:"DATABASE_NAME" doc:"Database Name"`
	DbUsername        string   `envconfig:"DATABASE_USERNAME" doc:"Database Username"`
	DbPassword        string   `envconfig:"DATABASE_PASSWORD" doc:"Database password"`
	CookieStoreSecret string   `envconfig:"COOKIE_STORE_SECRET" default:"secret" doc:"base64 encoded key to use for encrypting cookies"`
	SecureCookies     bool     `envconfig:"SECURE_COOKIE" doc:"Enable 'secure' flag on cookies" default:"false"`
}

var Constants = constants{}

func init() {
	err := envconfig.Process("PARM", &Constants)
	if err != nil {
		logrus.Fatal(err)
	}
}
