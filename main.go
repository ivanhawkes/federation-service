package main

import (
	"characters"
	"factions"
	"federations"
	"my"
	"profiles"
	"realms"
	"users"
	"zones"
)

func init() {
	// Register all the routes we need and document the interface.
	characters.CharacterApi{}.Register()
	factions.FactionApi{}.Register()
	federations.FederationApi{}.Register()
	my.MyApi{}.Register()
	profiles.ProfileApi{}.Register()
	realms.RealmApi{}.Register()
	users.UserApi{}.Register()
	zones.ZoneApi{}.Register()
}
