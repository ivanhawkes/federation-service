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

// The various states for a loottable resource.
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

type LootShallow struct {
	Key          string    `datastore:"-" json:"key" xml:"key"`
	LastModified time.Time `json:"last_modified" xml:"last-modified"`
	Version      int       `json:"version" xml:"version"`
	Status       int       `json:"status" xml:"status"`
	Name         string    `json:"name" xml:"name"`
	Link         Link      `datastore:"-" json:"link" xml:"link"`
}

type LootEntry struct {
	ItemId      int64   `datastore:"ItemId" json:"item_id" xml:"item-id"`
	Probability float32 `datastore:"Probability" json:"probability" xml:"probability"`
	Quantity    int16   `datastore:"Quantity" json:"quantity" xml:"quantity"`
}

type LootTable struct {
	LootShallow
	AllowPreload  bool        `json:"allow_preload" xml:"allow-preload"`
	Probabilities []LootEntry `json:"probabilities" xml:"probabilities"`
}

type LootSummary struct {
	LootTables []LootShallow `json:"loot_tables" xml:"loot-tables"`
}

type LootQuery struct {
	LootTables []LootTable `json:"loot_tables" xml:"loot-tables"`
}

type LootTableApi struct {
	Path string
}

func init() {
	log.Printf("Registering " + kind)
}

// Register the routes we require for this resource type.
//
func (api LootTableApi) Register() {
	api.RegisterAdmin()
	api.RegisterServer()
}

// Attempts to create a valid key for a resource.
//
func (api LootTableApi) getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("loottable-id"))
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
