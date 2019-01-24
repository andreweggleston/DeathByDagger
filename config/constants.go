package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type constants struct {
	ListenAddress     string   `envconfig:"SERVER_ADDR" default:"0.0.0.0:8081" doc:"Address to serve on"`
	PublicAddress     string   `envconfig:"PUBLIC_ADDR" doc:"Publicly accessible address for the server, requires schema"`
	AllowedOrigins    []string `envconfig:"ALLOWED_ORIGINS" default:"*"`
	DbAddr	          string   `envconfig:"DATABASE_ADDR" doc:"Database Address"`
	DbDatabase 	      string   `envconfig:"DATABASE_NAME" doc:"Database Name"`
	DbUsername        string   `envconfig:"DATABASE_USERNAME" doc:"Database Username"`
	DbPassword        string   `envconfig:"DATABASE_PASSWORD" doc:"Database password"`
	CookieStoreSecret string   `envconfig:"COOKIE_STORE_SECRET" default:"secret" doc:"base64 encoded key to use for encrypting cookies"`
	SecureCookies     bool     `envconfig:"SECURE_COOKIE" doc:"Enable 'secure' flag on cookies" default:"false"`
	OpenIDUrl		  string   `envconfig:"OPENID_AUTH_DOMAIN" doc:"Address of the openid Auth0 domain"`
	OpenIDClientID	  string   `envconfig:"OPENID_CLIENT_ID" doc:"Auth0 Client ID"`
	OpenIDClientSec   string   `envconfig:"OPENID_CLIENT_SECRET" doc:"Auth0 Client Secret"`
	CookieDomain      string   `envconfig:"SERVER_COOKIE_DOMAIN" default:"" doc:"Cookie URL domain"`
	LoginRedirectPath string   `envconfig:"SERVER_REDIRECT_PATH" doc:"URL to redirect user to after a successful login"`
	SlackToken string   `envconfig:"SLACK_TOKEN" doc:"Slack token for bot for interactive messages"`
}

var Constants = constants{}

func init() {
	err := envconfig.Process("DBD", &Constants)
	if err != nil {
		logrus.Fatal(err)
	}
}
