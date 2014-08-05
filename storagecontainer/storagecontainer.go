package storagecontainer

import (
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"resource"
)

// Status values for records of this resource type.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
)

type Shallow struct {
	resource.BaseResource
}

type Resource struct {
	Shallow

	// The owner key is a key belonging to one of these kinds - Character, Faction, Profile, Account
	OwnerKey datastore.Key `json:"owner_key" xml:"owner-key"`

	// The general class of storage for this container e.g.
	// "Character", "Character Bank", "Mail", "Faction Bank", "Quest", "Resource"
	Class string `json:"class" xml:"class"`

	// The specific name of the container e.g.
	// "Character 01", "Character Bank 01", Mail", "Resource.Ore", "Resource.Herb"
	Name string `json:"name" xml:"name"`

	// The maximum number of slots available within this container. The client is responsible for managing this figure and
	// ensuring they don't add more items than there are slots.
	SlotMax int32 `json:"slot_max" xml:"slot_max"`
}

func (s Shallow) Kind() string {
	return "storagecontainer"
}

func (s Shallow) AdminRootPath() string {
	return "/api/admin/" + s.Kind()
}

func (s Shallow) ShardRootPath() string {
	return "/api/shard/" + s.Kind()
}

func (s Shallow) AccountRootPath() string {
	return "/api/account/" + s.Kind()
}

func (s Shallow) ClientRootPath() string {
	return "/api/client/" + s.Kind()
}

func (s Shallow) PreferredLink(k *datastore.Key) string {
	return s.ShardRootPath() + "/" + k.Encode()
}

func init() {
}

// Register the routes we require for this resource type.
//
func (res Resource) Register() {
	log.Printf(Shallow{}.Kind() + " Register")
	//	res.RegisterAdmin()
	res.RegisterShard()
	// res.RegisterAccount()
	// res.RegisterClient()
}

// Attempts to create a valid key for a resource.
//
func (res *Resource) getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("resource-id"))
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid.\n")
		return nil, err
	}

	// Check for shenanigans with the key.
	if k.Kind() != res.Kind() {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid for this type of resource.\n")
		return nil, err
	}

	return k, nil
}
