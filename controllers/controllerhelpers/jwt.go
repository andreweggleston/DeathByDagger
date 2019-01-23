package controllerhelpers

import (
	"encoding/base64"
	"fmt"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	signingKey []byte
)

func init() {
	if config.Constants.CookieStoreSecret == "secret" {
		logrus.Warning("Using an insecure encryption key")
		signingKey = []byte("secret")
		return
	}

	var err error
	signingKey, err = base64.StdEncoding.DecodeString(config.Constants.CookieStoreSecret)
	if err != nil {
		logrus.Fatal(err)
	}
}

func NewToken() { //todo implement player

}

func verifyToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return signingKey, nil
}

func GetToken(r *http.Request) (*jwt.Token, error) {
	cookie, err := r.Cookie("auth-jwt")
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &DaggerClaims{},verifyToken)
	return token, err
}

func GetPlayer(token *jwt.Token) /* *player.Player */ {
	//TODO
}