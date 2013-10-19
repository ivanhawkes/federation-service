package factions

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/emicklei/go-restful"
	"net/http"
	"time"
)

// The various states for a faction resource.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted)

type FactionShallow struct {
	Id string `datastore:"-" json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	Link string `datastore:"-" json:"link" xml:"link"`
}

type Faction struct {
	FactionShallow
	UserId string `datastore:"UserId" json:"-" xml:"-"` // Owner and by extension, leader. TODO: perhaps not require the leader to be this User.
	LastModified time.Time `json:"-" xml:"-"`
	Status int `json:"status" xml:"status"`
}

type FactionApi struct {
	Path string
}

func init() {
	api := FactionApi{Path: "/factions"}
	api.register()
}

// Register the routes we require for this resource type.
//
func (api FactionApi) register() {
	ws := new(restful.WebService)

	ws.
		Path(api.Path).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.create).
		// Swagger documentation.
		Doc("create a new faction").
		Param(ws.BodyParameter("Faction", "representation of a faction").DataType("factions.Faction")).
		Reads(Faction{}))

	ws.Route(ws.GET("/{faction-id}").To(api.read).
		// Swagger documentation.
		Doc("read a faction").
		Param(ws.PathParameter("faction-id", "identifier for a faction").DataType("string")).
		Writes(Faction{}))

	ws.Route(ws.PUT("/{faction-id}").To(api.update).
		// Swagger documentation.
		Doc("update an existing faction").
		Param(ws.PathParameter("faction-id", "identifier for a faction").DataType("string")).
		Param(ws.BodyParameter("Faction", "representation of a faction").DataType("factions.Faction")).
		Reads(Faction{}))

	ws.Route(ws.DELETE("/{faction-id}").To(api.delete).
		// Swagger documentation.
		Doc("delete a faction").
		Param(ws.PathParameter("faction-id", "identifier for a faction").DataType("string")))

	restful.Add(ws)
}

// Create a new resource.
//
func (api *FactionApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	f := new(Faction)
	err := r.ReadEntity(&f)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Set some fields that need special handling.
	f.LastModified = time.Now()
	f.Status = StatusActive

	// The resource belongs to this faction.
	f.UserId = user.Current(c).ID

	// Set a user as our ancestor...this is done by querying for the key for the current user.
	var ancestor *datastore.Key
	q := datastore.NewQuery("users").
		Filter("UserId =", user.Current(c).ID).
		KeysOnly()
	if keys, err := q.GetAll(c, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if keys == nil {
			http.Error(w, "There is no user resource for this login account", http.StatusNotAcceptable)
			return
		}
		ancestor = keys[0]
	}

	// Store the faction.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "factions", ancestor), f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The resource Id.
	f.Id = k.Encode ()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", api.Path+"/"+k.Encode())

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	f.Link = api.Path+"/"+k.Encode()

	// Return the resultant entity.
	w.WriteHeader(http.StatusCreated)
	w.WriteEntity(f)
}

// Read the resource.
//
func (api FactionApi) read(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("faction-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the entity from the datastore.
	f := Faction{}
	if err := datastore.Get(c, k, &f); err != nil {
		if err.Error() == "datastore: no such entity" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check we own the resource before allowing them to view it.
	// Optionally, return a 404 instead to help prevent guessing ids.
	// TODO: Allow admins access.
	//if f.UserId != user.Current(c).ID {
	//	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	//	return
	//}

	// Set their Id.
	f.Id = k.Encode ()

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	f.Link = api.Path+"/"+k.Encode()

	w.WriteEntity(f)
}

// Update the resource.
//
func (api *FactionApi) update(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("faction-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshall the entity from the request into a struct.
	f := new(Faction)
	err = r.ReadEntity(&f)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Faction{}
	if err := datastore.Get(c, k, &old); err != nil {
		if err.Error() == "datastore: no such entity" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check we own the resource before allowing them to update it.
	// Optionally, return a 404 instead to help prevent guessing ids.
	// TODO: Allow admins access.
	if old.UserId != user.Current(c).ID {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	// Since the whole entity is re-written, we need to assign any invariant fields again
	// e.g. the owner of the entity.
	f.UserId = user.Current(c).ID

	// Keep track of the last modification date.
	f.LastModified = time.Now()

	// Attempt to overwrite the old entity.
	_, err = datastore.Put(c, k, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Let them know it succeeded.
	w.WriteHeader(http.StatusNoContent)
}

// Delete the resource.
//
func (api *FactionApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("faction-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Faction{}
	if err := datastore.Get(c, k, &old); err != nil {
		if err.Error() == "datastore: no such entity" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check we own the resource before allowing them to delete it.
	// Optionally, return a 404 instead to help prevent guessing ids.
	// TODO: Allow admins access.
	if old.UserId != user.Current(c).ID {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	// Delete the entity.
	if err := datastore.Delete(c, k); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Success notification.
	w.WriteHeader(http.StatusNoContent)
}
