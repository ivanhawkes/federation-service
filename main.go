package main

import (
	"characterappearance"
	"characterlocation"
	"characters"
	"factions"
	"loottable"
	//	"profiles"
	"realms"
	"storagecontainer"
	"storageitem"
	"users"
	"zones"
)

func init() {
	// Register all the routes we need and document the interface.
	characters.CharacterApi{}.Register()
	factions.FactionApi{}.Register()
	//	profiles.ProfileApi{}.Register()
	realms.RealmApi{}.Register()
	users.UserApi{}.Register()
	zones.ZoneApi{}.Register()
	loottable.Resource{}.Register()
	characterappearance.Resource{}.Register()
	characterlocation.Resource{}.Register()
	storagecontainer.Resource{}.Register()
	storageitem.Resource{}.Register()
}
