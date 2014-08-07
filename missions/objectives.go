package missions

import (
	"appengine/datastore"
)

type KillXofY struct {
	Required   int           `json:"required" xml:"Required"`
	CreatureId datastore.Key `json:"creature_id" xml:"creature-id"`
}

// This might be able to be handled by KillXofY with a Required of 1.
type KillNamedY struct {
	Creature datastore.Key `json:"creature_id" xml:"creature-id"`
}

type OneOrMoreItems struct {
	Amount int           `json:"amount" xml:"amount"`
	ItemId datastore.Key `json:"item_id" xml:"item-id"`
}

// Not efficient, make it use a list for the items and kills.
type CollectXofY struct {
	ItemList []OneOrMoreItems
}

// TODO: Escort mission.
type Escort struct {
	Stuff string
}

// TODO: Defend a location or a person.
type DefendX struct {
	Stuff string
}

//
// Crafting quests.
//

// TODO: Craft skill mission.
type CraftSkill struct {
	Stuff string
}

// TODO: Create a craft item.
type CraftManufacture struct {
	Stuff string
}

// TODO: Upgrade an item.
type CraftUpgrade struct {
	Stuff string
}

// TODO: Repair an item
type CraftRepair struct {
	Stuff string
}

//
// PVP quests.
//

// TODO: PVP kill
type PVPKillXofY struct {
	Stuff string
}

// TODO: PVP - victory
type PVPVictory struct {
	Zone string
}

//
// Vehicle quests.
//

// TODO: a bombing run
type VehicleBombingRun struct {
	Zone string
}

// TODO: turret - kill things with a turrent weapon
type VehicleTurrent struct {
	Zone string
}

// TODO: win a race
type VehicleRace struct {
	Zone string
}

//
// Class quest
//

// TODO:
type ClassMission struct {
	Zone string
}

//
// Training
//

// TODO:
type TrainingWeapon struct {
	Weapon string
}

//
// Professions
//

// TODO:
type ProfessionGatherHerbs struct {
	Weapon string
}

type ProfessionCooking struct {
	Weapon string
}

//
// Lore
//
type LoreLearn struct {
	Stuff string
}

//
// Raid
//
type RaidCompletion struct {
	ZoneId datastore.Key `json:"zone_id" xml:"zone-id"`
}

type RaidBoss struct {
	CharacterId datastore.Key `json:"character_id" xml:"character-id"`
}

//
// Dungeon
//
type DungeonCompletion struct {
	ZoneId datastore.Key `json:"zone_id" xml:"zone-id"`
}

type DungeonBoss struct {
	CharacterId datastore.Key `json:"character_id" xml:"character-id"`
}

// Get to a location
type DungeonLocation struct {
	Location string
}

//
// World boss
//
type WorldBoss struct {
	CharacterId datastore.Key `json:"character_id" xml:"-id"`
}
