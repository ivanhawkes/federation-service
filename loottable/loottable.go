package loottable

import (
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"resource"
)

const (
	kind           = "loottable"
	adminRootPath  = "/api/client/admin/" + kind
	shardRootPath  = "/api/shard/" + kind
	clientRootPath = "/api/client/" + kind
)

// Status values for records of this resource type.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type ProbabilityEntry struct {
	ItemId      int64   `datastore:"ItemId" json:"item_id" xml:"item-id"`
	Probability float32 `datastore:"Probability" json:"probability" xml:"probability"`
	Quantity    int16   `datastore:"Quantity" json:"quantity" xml:"quantity"`
}

type Shallow struct {
	resource.BaseResource
	Name string `json:"name" xml:"name"`
}

type Resource struct {
	Shallow
	AllowPreload  bool               `json:"allow_preload" xml:"allow-preload"`
	Probabilities []ProbabilityEntry `json:"probabilities" xml:"probabilities"`
}

type ListSummary struct {
	Entry []Shallow `json:"entry" xml:"entry"`
}

type ListComprehensive struct {
	Entry []Resource `json:"entry" xml:"entry"`
}

func (s Shallow) Kind () string {
	return kind
}

func (s Shallow) AdminRootPath () string {
	return "/api/client/admin/" + s.Kind ()
}

func (s Shallow) ShardRootPath () string {
	return "/api/shard/" + s.Kind ()
}

func (s Shallow) ClientRootPath () string {
	return "/api/client/" + s.Kind ()
}

func init() {
	log.Printf("Registering " + kind)
}

// Register the routes we require for this resource type.
//
func (res Resource) Register() {
	res.RegisterAdmin()
	res.RegisterShard()
}

// Attempts to create a valid key for a resource.
//
func (*Resource) getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {

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
