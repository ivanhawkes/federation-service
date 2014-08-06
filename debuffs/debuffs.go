package Debuffs

import (
	"time"
)

const (
	rootPath = "/debuffs"
)

// The various states for a resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
)

type DebuffsShallow struct {
	Id   string `datastore:"-" json:"id" xml:"id"`
	Link string `datastore:"-" json:"link" xml:"link"`
	Name string `json:"name" xml:"name"`
}

type Debuffs struct {
	DebuffsShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`

	// Components that define this itemtype.
	Components []interface{}

	// Icon to display in the UI
	Icon string
}

type DebuffsApi struct {
	Path string
}
