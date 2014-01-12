package storageitem

import (
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"resource"
	"time"
)

// Status values for records of this resource type.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type Shallow struct {
	resource.BaseResource
}

type Manufacture struct {
	RealmID	int16 `json:"realm_id" xml:"realm-id"`
	ShardID	int16 `json:"shard_id" xml:"shard-id"`	
	Created time.Time `json:"created" xml:"created"`
}

type Resource struct {
	Shallow
	
	// The owner key is a key belonging to one of these kinds - Character, Faction, Profile, Account
	OwnerKey datastore.Key `json:"owner_key" xml:"owner-key"`

	// The storage container in which this item is current stored.
	StorageContainer datastore.Key `json:"storage_container_key" xml:"storage-container-key"`

	// This needs to be unique. Make sure to update it as required when moving from one container to another.
	SlotID	int32 `json:"slot_id" xml:"slot-id"`

	// A count of how many of the given item are stored in this slot.
	Count int32  `json:"count" xml:"count"`

	// If the item is bound, then this key will belong to one of these kinds - Character, Faction, Profile, Account
	BoundKey datastore.Key `json:"bound_key" xml:"bound-key"`

	// Each item should carry a little information about it's manufacture.
	Manufacture Manufacture `json:"manufacture" xml:"manufacture"`
}

func (s Shallow) Kind() string {
	return "storageitem"
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
