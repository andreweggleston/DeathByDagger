package helpers

import "github.com/andreweggleston/DeathByDagger/helpers/authority"

const (
	RolePlayer authority.AuthRole = iota
	RoleMod
	RoleAdmin
	RoleDeveloper
)

var RoleNames = map[authority.AuthRole]string{
	RoleDeveloper: "developer",
	RolePlayer:    "player",
	RoleMod:       "moderator",
	RoleAdmin:     "administrator",
}

var RoleMap = map[string]authority.AuthRole{
	"player":        RolePlayer,
	"moderator":     RoleMod,
	"administrator": RoleAdmin,
	"developer":     RoleDeveloper,
}

// You can't change the order of these
const (
	ActionBanJoin authority.AuthAction = iota
	ActionBanCreate
	ActionBanChat
	ActionChangeRole
	ActionViewLogs
	ActionViewPage //view admin pages
	ActionDeleteChat
	ModifyServers //add/remove servers
)

var ActionNames = map[authority.AuthAction]string{
	ActionBanCreate: "ActionBanCreate",
	ActionBanJoin:   "ActionBanJoin",
	ActionBanChat:   "ActionBanChat",

	ActionChangeRole: "ActionChangeRole",
}

func init() {
	RoleDeveloper.Allow(ActionViewPage)

	RoleMod.Inherit(RolePlayer)
	RoleMod.Allow(ActionBanChat)
	RoleMod.Allow(ActionBanJoin)
	RoleMod.Allow(ActionBanCreate)
	RoleMod.Allow(ActionViewLogs)
	RoleMod.Allow(ActionViewPage)
	RoleMod.Allow(ActionDeleteChat)
	RoleMod.Allow(ModifyServers)

	RoleAdmin.Inherit(RoleMod)
	RoleAdmin.Allow(ActionChangeRole)
}
