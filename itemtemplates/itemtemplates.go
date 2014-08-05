package ItemTemplates

import (
	"time"
)

const (
	rootPath = "/itemtemplates"
)

// The various states for a resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
)

type ItemTemplatesShallow struct {
	Id   string `datastore:"-" json:"id" xml:"id"`
	Link string `datastore:"-" json:"link" xml:"link"`
	Name string `json:"name" xml:"name"`
}

type ItemTemplates struct {
	ItemTemplatesShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`

	// Components that define this itemtype.
	Components []interface{}

	// Can be equipped to an inventory slot.
	IsEquipable bool

	// Can be dropped onto the ground in the world.
	IsDropable bool

	// Can be destroyed from inventory.
	IsDestroyable bool

	// Icon to display in the UI
	Icon string
}

type ItemTemplatesApi struct {
	Path string
}
