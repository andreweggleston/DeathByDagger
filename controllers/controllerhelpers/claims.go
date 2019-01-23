package controllerhelpers

import (
	"errors"
	"github.com/andreweggleston/DeathByDagger/helpers/authority"
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
	//TODO: find in DB where id is id and csh_username is CSHUsername
	return count != 0
}

func (c DaggerClaims) Valid() error{
	if !playerExists(c.PlayerID, c.CSHUsername){
		return errors.New("player not found")
	}
	return nil
}
