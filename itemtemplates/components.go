package itemtemplates

import (
	"time"
)

// Non instance related components for each type of item that may appear in the game world.
// Instanced components should be made as item components.

const (
	QualityTrash = iota
	QualityAverage
	QualityAboveAverage
	QualityRare
	QualityEpic
	QualityLengendary
	QualityArtifact
)

const (
	ResourceCookingGeneric = iota
	ResourceCookingBerry
	ResourceCookingFruit
	ResourceCookingGrain
	ResourceCookingLegume
	ResourceCookingMeat
	ResourceCookingNuts
	ResourceCookingRoot
	ResourceCookingSalt
	ResourceCookingSeeds
	ResourceCookingSpice
	ResourceCookingVegetable
)

const (
	NotBound = iota
	CharacterBound
	AccountBound
	FactionBound
	RealmBound
	FederationBound
)

const (
	UniqueCharacter = iota
	UniqueAccount
)

const (
	ResourceClothingGeneric = iota
)

type Name struct { // required
	HiddenName  string
	DisplayName string
}

type Quality struct { // required - but may be overidden on the actual item due to improvements.
	Quality uint8
}

type Model struct {
	Mesh      string
	Material  string
	SpecialFx string
}

type Physicalise struct {
	Mass float32
}

// Anything that can be destroyed in the world.

type Destructable struct {
	HitPoints int32
}

type Throwable struct {
	Distance float32 // metres
}

type Tradeable struct {
	BindRule uint8
}

type Unique struct {
	Unique uint8
}

// Can this item respawn?

type Respawnable struct {
	Timer uint32 // milliseconds
}

// Can be interacted with in the world, useable in their inventory.

type Usable struct {
	// Message to display when the player is aiming at this item.
	Message string
}

// Can be picked up out of the world into the character's inventory.

type Pickable struct {
	// Message to display when the player is aiming at this item.
	Message string
}

// Perhaps this stays in the weapon template?
type Pose struct {
	Pose    string
	AimPose string
}

type Book struct {
	Title           string
	AuthorId        int64
	PublicationDate time.Time // needs to refer to our time system in game
	Content         int64     // key to contents in database / filesystem
}

type EquipBuff struct {
	BuffId int64
}

// Consumables

type ConsumeBuff struct {
	BuffId   int64
	Duration int32 // miliseconds
}

type Food struct {
}

type Drink struct {
}

type Potion struct {
}

type Scroll struct {
}

type Recipe struct {
}

// Crafting. Will some materials be useful for multiple crafts?

type CraftMaterial struct {
}

type Gem struct {
}

type Jewellry struct {
}

type Clothing struct {
}

type Tool struct {
}

// Wearables

type WearableClothing struct { // ? wearable and clothing?
}

// Weapons

type Weapon struct {
	// e.g. mace 1h, sword 2h
	Category int64

	// Class definition inside CryEngine - could be derived from the Category
	WeaponClass string

	// Base speed for simple attacks.
	Speed int

	// Base DPS for simple attacks.
	Dps float32

	// A template to apply which will create the needed file for CryEngine to use this as a weapon. There will need
	// to be generic ones for each class, and some specific ones for edge cases.
	XmlTemplate string
}

// Armour

type ArmourHead struct {
}

type ArmourChest struct {
}

type ArmourLegs struct {
}

type ArmourArms struct {
}

type ArmourHands struct {
}

type ArmourFeet struct {
}

// Currency

type CurrencyCents struct {
	Amount uint64
}

type CurrencyPVPBadge struct {
	Amount uint8
}

type CurrencyDungeonBadge struct {
	Amount uint8
}

type CurrencyRaidBadge struct {
	Amount uint8
}

// Keys to things
type Key struct {
}
