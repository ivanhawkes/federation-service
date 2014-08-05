package Items

import (
	"time"
)

const (
	rootPath = "/item"
)

// The various states for a resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
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
