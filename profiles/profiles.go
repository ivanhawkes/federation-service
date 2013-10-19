package profiles

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/emicklei/go-restful"
	"net/http"
	"time"
)

type Profile struct {
	Id           string    `datastore:"-" json:"id" xml:"id"`
	UserId       string    `datastore:"UserId" json:"-" xml:"-"`
	LastModified time.Time `json:"-" xml:"-"`
	FirstName    string    `json:"first_name" xml:"first-name"`
	NickName     string    `json:"nick_name" xml:"nick-name"`
	LastName     string    `json:"last_name" xml:"last-name"`
	FactionId    int64     `datastore:"FactionId" json:"faction_id" xml:"faction-id"` // The faction you belong to, if any.
}

type ProfileApi struct {
	Path string
}

func init() {
	api := ProfileApi{Path: "/profiles"}
	api.register()
}

// Register the routes we require for this resource type.
//
func (api ProfileApi) register() {
	ws := new(restful.WebService)

	ws.
		Path(api.Path).
		// You can specify consumes and produces per route as well.
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.create).
		// Swagger documentation.
		Doc("create a new profile").
		Param(ws.BodyParameter("Profile", "representation of a profile").DataType("profiles.Profile")).
		Reads(Profile{}))

	ws.Route(ws.GET("/{profile-id}").To(api.read).
		// Swagger documentation.
		Doc("read a profile").
		Param(ws.PathParameter("profile-id", "identifier for a profile").DataType("string")).
		Writes(Profile{}))

	ws.Route(ws.PUT("/{profile-id}").To(api.update).
		// Swagger documentation.
		Doc("update an existing profile").
		Param(ws.PathParameter("profile-id", "identifier for a profile").DataType("string")).
		Param(ws.BodyParameter("Profile", "representation of a profile").DataType("profiles.Profile")).
		Reads(Profile{}))

	ws.Route(ws.DELETE("/{profile-id}").To(api.delete).
		// Swagger documentation.
		Doc("delete a profile").
		Param(ws.PathParameter("profile-id", "identifier for a profile").DataType("string")))

	restful.Add(ws)
}

// Create a new resource.
//
func (api *ProfileApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	p := new(Profile)
	err := r.ReadEntity(&p)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Ensure we start with a sensible value for this field.
	p.LastModified = time.Now()

	// The resource belongs to this user.
	p.UserId = user.Current(c).ID

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

	// Store the profile.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "profiles", ancestor), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The resource Id.
	p.Id = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", api.Path+"/"+k.Encode())

	// Return the resultant entity.
	w.WriteHeader(http.StatusCreated)
	w.WriteEntity(p)
}

// Read the resource.
//
func (api ProfileApi) read(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("profile-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the entity from the datastore.
	p := Profile{}
	if err := datastore.Get(c, k, &p); err != nil {
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
	if p.UserId != user.Current(c).ID {
		http.Error(w, "You do not have access to this resource", http.StatusForbidden)
		return
	}

	// Set their Id.
	p.Id = k.Encode()

	w.WriteEntity(p)
}

// Update the resource.
//
func (api *ProfileApi) update(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("profile-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshall the entity from the request into a struct.
	p := new(Profile)
	err = r.ReadEntity(&p)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Profile{}
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
	p.UserId = user.Current(c).ID

	// Keep track of the last modification date.
	p.LastModified = time.Now()

	// Attempt to overwrite the old entity.
	_, err = datastore.Put(c, k, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Let them know it succeeded.
	w.WriteHeader(http.StatusNoContent)
}

// Delete the resource.
//
func (api *ProfileApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("profile-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the old entity from the datastore.
	old := Profile{}
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
