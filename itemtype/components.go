package ItemTypes

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
	FoodTypeMeat = iota
	FoodTypeVegetable
	FoodTypeFruit
	FoodTypeGrain
	FoodTypeLegume
	FoodTypeNuts
	FoodTypeSeeds
	FoodTypeBerry
	FoodTypeRoot
	FoodTypeSalt
	FoodTypeSpice
	FoodTypeMeal // separate ingredients and results or not since results can also be ingredients?
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
	// moveable, etc...
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

// Consumables

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
