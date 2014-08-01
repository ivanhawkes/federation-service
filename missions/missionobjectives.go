package missions

import (
	"appengine"
	"appengine/datastore"
)

type MissionObjectiveX struct {
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`
}

type KillXofY struct {
	Required int           `json:"required" xml:"Required"`
	Creature datastore.Key `json:"creature" xml:"Creature"`
}

// This might be able to be handled by KillXofY with a Required of 1.
type KillNamedY struct {
	Creature datastore.Key `json:"creature" xml:"Creature"`
}

// Not efficient, make it use a list for the items and kills.
type KillXofYCollectIofJ struct {
	RequiredKills int           `json:"status" xml:"status"`
	Creature      datastore.Key `json:"creature" xml:"Creature"`
	RequiredItem1 int           `json:"required-items-1" xml:"RequiredItem1"`
	Item1         datastore.Key `json:"item-1" xml:"Item1"`
	RequiredItem2 int           `json:"required-item-2" xml:"RequiredItem2"`
	Item2         datastore.Key `json:"item-2" xml:"Item2"`
	RequiredItem3 int           `json:"required-item-3" xml:"RequiredItem3"`
	Item3         datastore.Key `json:"item-3" xml:"Item3"`
	RequiredItem4 int           `json:"required-item-4" xml:"RequiredItem4"`
	Item4         datastore.Key `json:"item-4" xml:"Item4"`
}

// Not efficient, make it use a list for the items and kills.
type CollectXofY struct {
	RequiredItem1 int           `json:"required-items-1" xml:"RequiredItem1"`
	Item1         datastore.Key `json:"item-1" xml:"Item1"`
	RequiredItem2 int           `json:"required-item-2" xml:"RequiredItem2"`
	Item2         datastore.Key `json:"item-2" xml:"Item2"`
	RequiredItem3 int           `json:"required-item-3" xml:"RequiredItem3"`
	Item3         datastore.Key `json:"item-3" xml:"Item3"`
	RequiredItem4 int           `json:"required-item-4" xml:"RequiredItem4"`
	Item4         datastore.Key `json:"item-4" xml:"Item4"`
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
type RaidBoss struct {
	Boss string
}

//
// Dungeon
//
type DungeonBoss struct {
	Boss string
}

// Get to a location
type DungeonLocation struct {
	Location string
}

//
// World boss
//
type WorldBoss struct {
	Boss string
}
