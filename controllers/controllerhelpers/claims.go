package controllerhelpers

import (
	"errors"
	"github.com/andreweggleston/DeathByDagger/helpers/authority"
	db "github.com/andreweggleston/DeathByDagger/databaseDagger"
	"github.com/andreweggleston/DeathByDagger/models/player"
)

type DaggerClaims struct {
	PlayerID    uint               `json:"player_id"`
	CSHUsername string             `json:"csh_username"`
	Role        authority.AuthRole `json:"role"`
	IssuedAt    int64              `json:"iat"`
	Issuer      string             `json:"iss"`
}

func playerExists(id uint, CSHUsername string) bool {
	var count int
	db.DB.Model(&player.Player{}).Where("id = ? AND csh_username = ?", id, CSHUsername).Count(&count)
	return count != 0
}

func (c DaggerClaims) Valid() error{
	if !playerExists(c.PlayerID, c.CSHUsername){
		return errors.New("player not found")
	}
	return nil
}
