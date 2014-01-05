package loottable

import (
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"time"
)

const (
	kind           = "loottable"
	adminRootPath  = "/admin/" + kind
	serverRootPath = "/server/" + kind
	clientRootPath = "/client/" + kind
)

// Status values for records of this resource type.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type Link struct {
	Rel  string `datastore:"-" json:"rel" xml:"rel"`
	Href string `datastore:"-" json:"href" xml:"href"`
}

type Shallow struct {
	Key          string    `datastore:"-" json:"key" xml:"key"`
	LastModified time.Time `json:"last_modified" xml:"last-modified"`
	Version      int       `json:"version" xml:"version"`
	Status       int       `json:"status" xml:"status"`
	Link         Link      `datastore:"-" json:"link" xml:"link"`
	Name         string    `json:"name" xml:"name"`
}

type ProbabilityEntry struct {
	ItemId      int64   `datastore:"ItemId" json:"item_id" xml:"item-id"`
	Probability float32 `datastore:"Probability" json:"probability" xml:"probability"`
	Quantity    int16   `datastore:"Quantity" json:"quantity" xml:"quantity"`
}

type LootTable struct {
	Shallow
	AllowPreload  bool               `json:"allow_preload" xml:"allow-preload"`
	Probabilities []ProbabilityEntry `json:"probabilities" xml:"probabilities"`
}

type ListSummary struct {
	LootTables []Shallow `json:"loot_tables" xml:"loot-tables"`
}

type LootQuery struct {
	LootTables []LootTable `json:"loot_tables" xml:"loot-tables"`
}

type ResourceApi struct {
	Path string
}

func init() {
	log.Printf("Registering " + kind)
}

// Register the routes we require for this resource type.
//
func (api ResourceApi) Register() {
	api.RegisterAdmin()
	api.RegisterServer()
}

// Attempts to create a valid key for a resource.
//
func (api ResourceApi) getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("resource-id"))
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid.\n")
		return nil, err
	}

	// Check for shenanigans with the key.
	if k.Kind() != kind {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid for this type of resource.\n")
		return nil, err
	}

	return k, nil
}
