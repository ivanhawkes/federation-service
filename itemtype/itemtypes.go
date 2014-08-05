package ItemTypes

const (
	rootPath = "/itemtypes"
)

// The various states for a resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
)

type ItemTypeShallow struct {
	Id   string `datastore:"-" json:"id" xml:"id"`
	Link string `datastore:"-" json:"link" xml:"link"`
	Name string `json:"name" xml:"name"`
}

type ItemType struct {
	ItemTypeShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`

	// Components that define this itemtype.
	Components []interface{}
}

type ItemTypeApi struct {
	Path string
}
