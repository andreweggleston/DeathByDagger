package player

import "github.com/andreweggleston/DeathByDagger/helpers"

func (p *Player) DecoratePlayerTags() []string {
	tags := []string{helpers.RoleNames[p.Role]}
	return tags
}

func (p *Player) SetPlayerProfile() {
	p.Name = p.Alias()
}
