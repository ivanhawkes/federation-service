package Items

import (
	"time"
)

const (
	rootPath = "/item"
)

// Think over the states needed to control an item's lifespan from creation, to placing in the world,
// controlled by a faction, picked up by a player, placed into inventory, in-use, destroyed.
const (
	StatusCreated = iota
	StatusWorldOwned
	StatusFactionOwned
	StatusCharacterOwned
	StatusDestroyPending
	StatusDestroyed
)

type ItemShallow struct {
	Id   string `datastore:"-" json:"id" xml:"id"`
	Link string `datastore:"-" json:"link" xml:"link"`
	Name string `json:"name" xml:"name"`
}

type Item struct {
	ItemShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`

	// Components that define this item.
	Components []interface{}

	// The current owner of this item. This could be a player, a faction, a realm / zone, even the system. This needs thought
	// on how to implement it for all these cases.
	OwnerId int64 `json:"owner-id" xml:"owner-id"`
}

type ItemApi struct {
	Path string
}
