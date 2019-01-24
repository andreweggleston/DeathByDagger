package controllerhelpers

import (
	"encoding/base64"
	"fmt"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
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

func NewToken(player *player.Player) string {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = DaggerClaims{
		PlayerID:  		player.ID,
		CSHUsername:	player.CSHUsername,
		Role: 			player.Role,
		IssuedAt: 		time.Now().Unix(),
		Issuer: 		config.Constants.PublicAddress,

	}

	str, err := token.SignedString([]byte(signingKey))
	if err != nil {
		logrus.Error(err)
	}

	return str
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

	token, err := jwt.ParseWithClaims(cookie.Value, &DaggerClaims{}, verifyToken)
	return token, err
}

func GetPlayer(token *jwt.Token) *player.Player {
	player, _ := player.GetPlayerByID(token.Claims.(*DaggerClaims).PlayerID) //error is ignored because we know player exists
	return player
}
