package login

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/andreweggleston/DeathByDagger/config"
	"github.com/andreweggleston/DeathByDagger/controllers/controllerhelpers"
	"github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"time"
)

type UserInfo struct {
	UserInfo	User	`json:"userinfo"`
}

type User struct {
	Sub           string `json:"sub"`
	Username       string `json:"preferred_username"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
}

var (
	conf = &oauth2.Config{
		ClientID:config.Constants.OpenIDClientID,
		ClientSecret:config.Constants.OpenIDClientSec,
		RedirectURL:"http://"+config.Constants.LoginRedirectPath+"/callback",
		Scopes:[]string{"openid, profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + config.Constants.OpenIDUrl + "/auth",
			TokenURL: "https://" + config.Constants.OpenIDUrl + "/token",
		},
	}
	state = ""
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	retrievedState := r.FormValue("state")
	if retrievedState != state {
		fmt.Printf("oauth state is invalid, expected '%s', got '%s' \n", state, retrievedState)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Printf("r.ParseForm() failed with %s\n", err)
	}

	code := r.Form.Get("code")

	tok, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := conf.Client(context.TODO(), tok)

	resp, err := client.Get("https://"+config.Constants.OpenIDUrl+"/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var user = new(User)
	err = json.Unmarshal(data, &user)
	if err != nil {
		logrus.Fatal(err)
	}


	p, err := player.GetPlayerByCSHUsername(user.Username)
	if err != nil {

		p, err = player.NewPlayer(user.Username)

		if err != nil {
			logrus.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		p.Name = user.GivenName + " " + user.FamilyName
		p.Sub = user.Sub
		p.Email = user.Email

		databaseDagger.DB.Create(p)
	}


	key := controllerhelpers.NewToken(p)
	cookie := &http.Cookie{
		Name:     "auth-jwt",
		Value:    key,
		Path:     "/",
		Domain:   config.Constants.CookieDomain,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   config.Constants.SecureCookies,
	}

	http.SetCookie(w, cookie)

	http.Redirect(w, r, "http://"+config.Constants.LoginRedirectPath, http.StatusFound)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	rand.Read(b)
	state = base64.StdEncoding.EncodeToString(b)
	url := conf.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth-jwt")
	if err != nil { //idiot wasnt even logged in LUL
		return
	}

	cookie.Domain = config.Constants.CookieDomain
	cookie.MaxAge = -1
	cookie.Expires = time.Time{}
	http.SetCookie(w, cookie)

	http.Redirect(w, r, "http://"+config.Constants.LoginRedirectPath, http.StatusFound)
}