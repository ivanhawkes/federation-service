package items

// import (
// 	"time"
// )

const (
	QualityTrash = iota
	QualityAverage
	QualityAboveAverage
	QualityRare
	QualityEpic
	QualityLengendary
	QualityArtifact
)

// This component allows an item to have it's own unique name within the world so players
// are allowed to give their sword / etc a name. It can also allow GMs to create unique and
// powerful legendary / artifact level items with unique names.

type UniqueName struct {
	UniqueName string
}

// Quality is different to the stock itemtype.

type Improvement struct {
	NewQuality uint8
}

// Allow one of more of these as a way to enhance a weapon.

type WeaponEnhancement struct {
	EnchantId int64
}

// Note: not sure if we want armour enchants but maybe stub it out for armour improvements.
// Can this be made generic enough to support gems / enchants / atatchments (looks improvement) or should we
// break that out into separate paths?

type ArmourEnhancement struct {
	HelmId      int64
	ChestId     int64
	LegsId      int64
	LeftHandId  int64
	RightHandId int64
	FeetId      int64
	ForearmsId  int64
}

// Dependancies: Model, Physicalise
type Container struct {
	StorageId   int64 // reference to the storage area that holds items in this container.
	Locked      bool
	RequiredKey int64 // might use a more generic method than explicit keys, this needs to work with loot chests, bank storage, boss chests, etc
	OwnerId     int64
}

type Model struct { // overide the default in itemtype - should use one from itemtype rather than redefine - here to stub out my thoughts.
	Mesh      string
	Material  string
	SpecialFx string
}
