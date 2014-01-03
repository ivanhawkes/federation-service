package loottable

import (
	"appengine"
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"time"
)

const (
	rootPath = "/loottable"
)

// The various states for a loottable resource.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type LootShallow struct {
	Id   string `datastore:"-" json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	Link string `datastore:"-" json:"link" xml:"link"`
}

type Loottable struct {
	LootShallow
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`
}

type LoottableApi struct {
	Path string
}

func init() {
	log.Printf("Loottable: Register")
}

// Register the routes we require for this resource type.
//
func (api LoottableApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path(rootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.create).
		// Swagger documentation.
		Doc("create a new loottable").
		Param(ws.BodyParameter("Loottable", "representation of a loottable").DataType("loottable.Loottable")).
		Reads(Loottable{}))

	ws.Route(ws.GET("/{loottable-id}").To(api.read).
		// Swagger documentation.
		Doc("read a loottable").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")).
		Writes(Loottable{}))

	ws.Route(ws.PUT("/{loottable-id}").To(api.update).
		// Swagger documentation.
		Doc("update an existing loottable").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")).
		Param(ws.BodyParameter("Loottable", "representation of a loottable").DataType("loottable.Loottable")).
		Reads(Loottable{}))

	ws.Route(ws.DELETE("/{loottable-id}").To(api.delete).
		// Swagger documentation.
		Doc("delete a loottable").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")))

	restful.Add(ws)
}

// Create a new resource.
//
func (api *LoottableApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	loottable := new(Loottable)
	err := r.ReadEntity(&loottable)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Set some fields that need special handling.
	loottable.LastModified = time.Now()
	loottable.Status = StatusActive

	// The resource belongs to this loottable.
	//loottable.UserId = user.Current(c).ID

	// TODO: Should be ancestor to a federation.
	// Set a user as our ancestor...this is done by querying for the key for the current user.
	/*	var ancestor *datastore.Key
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
		}*/

	// Store the loottable.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "loottable", nil /*ancestor*/), loottable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The resource Id.
	loottable.Id = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", rootPath+"/"+k.Encode())

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	loottable.Link = rootPath + "/" + k.Encode()

	// Return the resultant entity.
	w.WriteHeader(http.StatusCreated)
	w.WriteEntity(loottable)
}

// Read the resource.
//
func (api LoottableApi) read(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("loottable-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the entity from the datastore.
	loottable := Loottable{}
	if err := datastore.Get(c, k, &loottable); err != nil {
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
	//if loottable.UserId != user.Current(c).ID {
	//	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	//	return
	//}

	// Set their Id.
	loottable.Id = k.Encode()

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	loottable.Link = rootPath + "/" + k.Encode()

	w.WriteEntity(loottable)
}

// Update the resource.
//
func (api *LoottableApi) update(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("loottable-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshall the entity from the request into a struct.
	loottable := new(Loottable)
	err = r.ReadEntity(&loottable)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Loottable{}
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
	// if old.UserId != user.Current(c).ID {
	// 	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	// 	return
	// }

	// Since the whole entity is re-written, we need to assign any invariant fields again
	// e.g. the owner of the entity.
	// loottable.UserId = user.Current(c).ID

	// Keep track of the last modification date.
	loottable.LastModified = time.Now()

	// Attempt to overwrite the old entity.
	_, err = datastore.Put(c, k, loottable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Let them know it succeeded.
	w.WriteHeader(http.StatusNoContent)
}

// Delete the resource.
//
func (api *LoottableApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("loottable-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Loottable{}
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
	// if old.UserId != user.Current(c).ID {
	// 	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	// 	return
	// }

	// Delete the entity.
	if err := datastore.Delete(c, k); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Success notification.
	w.WriteHeader(http.StatusNoContent)
}
