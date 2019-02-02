package handler

import (
	"github.com/andreweggleston/DeathByDagger/helpers/wsevent"
	"github.com/andreweggleston/DeathByDagger/models/player"
)

type Unauth struct{}

func (Unauth) Name(s string) string {
	return string((s[0])+32) + s[1:]
}

func (Unauth) PlayerProfile(so *wsevent.Client, args struct {
	Studentid *string `json:"studentid"`
}) interface{} {

	player, err := player.GetPlayerByCSHUsername(*args.Studentid)
	if err != nil {
		return err
	}

	player.SetPlayerProfile()
	return newResponse(player)
}
