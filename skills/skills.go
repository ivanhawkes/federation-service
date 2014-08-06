package Skills

import (
	"time"
)

const (
	rootPath = "/skills"
)

// The various states for a resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
)

type SkillsShallow struct {
	Id   string `datastore:"-" json:"id" xml:"id"`
	Link string `datastore:"-" json:"link" xml:"link"`
	Name string `json:"name" xml:"name"`
}

type Skills struct {
	SkillsShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`

	// Components that define this itemtype.
	Components []interface{}

	// Icon to display in the UI
	Icon string

	// Duration required to cast.
	CastTime float32

	// Time before this can be used again.
	Cooldown float32
}

type SkillsApi struct {
	Path string
}
