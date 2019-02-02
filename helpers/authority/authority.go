package authority

import "encoding/gob"

type AuthAction int

type AuthRole int

var permissions = make(map[AuthRole]map[AuthAction]bool)

func init() {
	gob.Register(AuthAction(0))
	gob.Register(AuthRole(0))
}

func (role AuthRole) Allow(action AuthAction) AuthRole {
	amap, ok := permissions[role]
	if !ok {
		amap = make(map[AuthAction]bool)
		permissions[role] = amap
	}

	amap[action] = true
	return role
}

func (role AuthRole) Disallow(action AuthAction) AuthRole {
	amap, ok := permissions[role]
	if !ok {
		amap = make(map[AuthAction]bool)
		permissions[role] = amap
	}

	amap[action] = false
	return role
}

func (role AuthRole) Inherit(otherrole AuthRole) AuthRole {
	mymap, ok := permissions[role]
	if !ok {
		mymap = make(map[AuthAction]bool)
		permissions[role] = mymap
	}

	othermap, otherok := permissions[otherrole]
	if !otherok {
		return role
	}

	for entry, val := range othermap {
		mymap[entry] = val
	}

	return role
}

func (role AuthRole) Can(action AuthAction) bool {
	mymap, ok := permissions[role]
	return ok && mymap[action]
}

func Can(roleInt int, action AuthAction) bool {
	var role = AuthRole(roleInt)
	return role.Can(action)
}

func Reset() {
	permissions = make(map[AuthRole]map[AuthAction]bool)
}
