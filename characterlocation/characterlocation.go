package characterlocation

import (
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"resource"
)

// Status values for records of this resource type.
const (
	StatusThisFederation = iota
	StatusAnotherFederation
)

type Shallow struct {
	resource.BaseResource
}

type Position struct {
	X float32 `json:"x" xml:"x"`
	Y float32 `json:"y" xml:"y"`
	Z float32 `json:"z" xml:"z"`
}

type Direction struct {
	Yaw float32 `json:"yaw" xml:"yaw"`
	Pitch float32 `json:"pitch" xml:"pitch"`
	Roll float32 `json:"roll" xml:"roll"`
}

type Resource struct {
	Shallow
	CharacterKey datastore.Key `json:"character_key" xml:"character-key"`
	RealmKey     datastore.Key `json:"realm_key" xml:"realm-key"`
	ZoneKey      datastore.Key `json:"zone_key" xml:"zone-key"`
	ShardKey     datastore.Key `json:"shard_key" xml:"shard-key"`
	Position     Position      `json:"position" xml:"position"`
	Direction    Direction     `json:"direction" xml:"direction"`
}

func (s Shallow) Kind() string {
	return "characterlocation"
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
	return s.ClientRootPath() + "/" + k.Encode()
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
