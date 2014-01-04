package realms

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"time"
)

const (
	rootPath = "/realms"
)

// The various states for a realm resource.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type RealmShallow struct {
	Id   string `datastore:"-" json:"id" xml:"id"`
	Link string `datastore:"-" json:"link" xml:"link"`
	Name string `json:"name" xml:"name"`
}

type Realm struct {
	RealmShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`
}

type RealmApi struct {
	Path string
}

func init() {
	log.Printf("Realms: Register")
}

// Register the routes we require for this resource type.
//
func (api RealmApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path(rootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.create).
		// Swagger documentation.
		Doc("create a new realm").
		Param(ws.BodyParameter("Realm", "representation of a realm").DataType("realms.Realm")).
		Reads(Realm{}))

	ws.Route(ws.GET("/{realm-id}").To(api.read).
		// Swagger documentation.
		Doc("read a realm").
		Param(ws.PathParameter("realm-id", "identifier for a realm").DataType("string")).
		Writes(Realm{}))

	ws.Route(ws.PUT("/{realm-id}").To(api.update).
		// Swagger documentation.
		Doc("update an existing realm").
		Param(ws.PathParameter("realm-id", "identifier for a realm").DataType("string")).
		Param(ws.BodyParameter("Realm", "representation of a realm").DataType("realms.Realm")).
		Reads(Realm{}))

	ws.Route(ws.DELETE("/{realm-id}").To(api.delete).
		// Swagger documentation.
		Doc("delete a realm").
		Param(ws.PathParameter("realm-id", "identifier for a realm").DataType("string")))

	restful.Add(ws)
}

// Create a new resource.
//
func (api *RealmApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	realm := new(Realm)
	err := r.ReadEntity(&realm)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Set some fields that need special handling.
	realm.LastModified = time.Now()
	realm.Status = StatusActive

	// The resource belongs to this realm.
	realm.UserId = user.Current(c).ID

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

	// Store the realm.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "realms", nil /*ancestor*/), realm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The resource Id.
	realm.Id = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", rootPath+"/"+k.Encode())

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	realm.Link = rootPath + "/" + k.Encode()

	// Return the resultant entity.
	w.WriteHeader(http.StatusCreated)
	w.WriteEntity(realm)
}

// Read the resource.
//
func (api RealmApi) read(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("realm-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the entity from the datastore.
	realm := Realm{}
	if err := datastore.Get(c, k, &realm); err != nil {
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
	//if realm.UserId != user.Current(c).ID {
	//	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	//	return
	//}

	// Set their Id.
	realm.Id = k.Encode()

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	realm.Link = rootPath + "/" + k.Encode()

	w.WriteEntity(realm)
}

// Update the resource.
//
func (api *RealmApi) update(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("realm-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshall the entity from the request into a struct.
	realm := new(Realm)
	err = r.ReadEntity(&realm)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Realm{}
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
	realm.UserId = user.Current(c).ID

	// Keep track of the last modification date.
	realm.LastModified = time.Now()

	// Attempt to overwrite the old entity.
	_, err = datastore.Put(c, k, realm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Let them know it succeeded.
	w.WriteHeader(http.StatusNoContent)
}

// Delete the resource.
//
func (api *RealmApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("realm-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Realm{}
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
