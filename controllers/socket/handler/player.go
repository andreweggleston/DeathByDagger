package handler

import (
	"errors"
	chelpers "github.com/andreweggleston/DeathByDagger/controllers/controllerhelpers"
	"github.com/andreweggleston/DeathByDagger/helpers"
	"github.com/andreweggleston/DeathByDagger/helpers/wsevent"
	"github.com/andreweggleston/DeathByDagger/models/player"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Player struct{}

func (Player) Name(s string) string {
	return string((s[0])+32) + s[1:]
}

func (Player) PlayerSettingsGet(so *wsevent.Client, args struct {
	Key *string `json:"key"`
}) interface{} {

	player := chelpers.GetPlayer(so.Token)
	if *args.Key == "*" {
		return newResponse(player.Settings)
	}

	setting := player.GetSetting(*args.Key)
	return newResponse(setting)
}

func (Player) PlayerSettingsSet(so *wsevent.Client, args struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}) interface{} {

	player := chelpers.GetPlayer(so.Token)

	switch *args.Key {
	case "siteAlias":
		if len(*args.Value) > 32 {
			return errors.New("Site alias must be under 32 characters long.")
		}
		player.SetSetting(*args.Key, *args.Value)

		player.SetPlayerProfile()
		so.EmitJSON(helpers.NewRequest("playerProfile", player))

	default:
		player.SetSetting(*args.Key, *args.Value)
	}

	return emptySuccess
}

func (Player) PlayerProfile(so *wsevent.Client, args struct {
	Cshusername *string `json:"cshusername"`
}) interface{} {
	cshusername := *args.Cshusername
	if cshusername == "" {
		cshusername = so.Token.Claims.(*chelpers.DaggerClaims).CSHUsername
	}

	player, err := player.GetPlayerByCSHUsername(cshusername)
	if err != nil {
		return err
	}

	player.SetPlayerProfile()

	logrus.Debug("Recieved request for: ", cshusername)

	return newResponse(player)
}

var (
	changeMu = new(sync.RWMutex)
	//stores the last time the player made a change to
	//the twitch bot's status (leave/join their channel)
	lastTextBotChange = make(map[uint]time.Time)
)


func (Player) PlayerEnableTextBot(so *wsevent.Client, _ struct{}) interface{} {
	player := chelpers.GetPlayer(so.Token)
	if player.PhoneNumber == "" {
		return errors.New("Please add a phone number first.")
	}

	logrus.Debug("Player ", player.Name, " enabled text bot")

	changeMu.RLock()
	last := lastTextBotChange[player.ID]
	changeMu.RUnlock()
	if time.Since(last) < time.Minute {
		return errors.New("Please wait for a minute before changing the bot's status")
	}


	changeMu.Lock()
	lastTextBotChange[player.ID] = time.Now()
	changeMu.Unlock()

	time.AfterFunc(1*time.Minute, func() {
		changeMu.Lock()
		delete(lastTextBotChange, player.ID)
		changeMu.Unlock()
	})

	return emptySuccess
}

func (Player) PlayerDisableTextBot(so *wsevent.Client, _ struct{}) interface{} {
	player := chelpers.GetPlayer(so.Token)
	if player.PhoneNumber == "" {
		return errors.New("Please add a phone number first.")
	}

	logrus.Debug("Player", player.Name, "disabled text bot")

	changeMu.RLock()
	last := lastTextBotChange[player.ID]
	changeMu.RUnlock()
	if time.Since(last) < time.Minute {
		return errors.New("Please wait for a minute before changing the bot's status")
	}


	changeMu.Lock()
	lastTextBotChange[player.ID] = time.Now()
	changeMu.Unlock()

	time.AfterFunc(1*time.Minute, func() {
		changeMu.Lock()
		delete(lastTextBotChange, player.ID)
		changeMu.Unlock()
	})
	return emptySuccess
}

func (Player) MarkTarget(so *wsevent.Client, _ struct{}) interface{}  {
	player := chelpers.GetPlayer(so.Token)
	player.MarkTarget()
	err := player.Save()


	if err != nil{
		return errors.New("Something broke. Sorry! Contact an admin!")
	}
	return emptySuccess
}

func (Player) ConfirmOwnMark(so *wsevent.Client, _ struct{}) interface{}  {
	player := chelpers.GetPlayer(so.Token)
	player.ConfirmOwnMark()
	err := player.Save()

	if err != nil{
		return errors.New("Something broke. Sorry! Contact an admin!")
	}
	return emptySuccess
}

func (Player) DenyOwnMark(so *wsevent.Client, _ struct{}) interface{}  {
	player := chelpers.GetPlayer(so.Token)
	player.DenyOwnMark()
	err := player.Save()

	if err != nil{
		return errors.New("Something broke. Sorry! Contact an admin!")
	}
	return emptySuccess
}