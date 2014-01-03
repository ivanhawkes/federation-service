package main

import (
	"characters"
	"factions"
	"loottable"
	"realms"
	"users"
	"zones"
)

func init() {
	// Register all the routes we need and document the interface.
	characters.CharacterApi{}.Register()
	factions.FactionApi{}.Register()
	realms.RealmApi{}.Register()
	users.UserApi{}.Register()
	zones.ZoneApi{}.Register()
	loottable.LootTableApi{}.Register()
}
