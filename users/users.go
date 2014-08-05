package users

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
	rootPath = "/users"
)

// The various states for a user resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type UserShallow struct {
	Id        string `datastore:"-" json:"id" xml:"id"`
	Link      string `datastore:"-" json:"link" xml:"link"`
	FirstName string `json:"first_name" xml:"first-name"`
	LastName  string `json:"last_name" xml:"last-name"`
	AvatarUrl string `json:"avatar_url" xml:"avatar-url"`
}

type User struct {
	UserShallow
	UserId       string    `datastore:"UserId" json:"-" xml:"-"` // The external Id for the user who this represents. Comes from user authentication library.
	LastModified time.Time `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`
}

type UserApi struct {
	Path string
}

func init() {
	log.Printf("Users: Register")
}

// Register the routes we require for this resource type.
//
func (api UserApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path(rootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.create).
		// Swagger documentation.
		Doc("create a new user").
		Param(ws.BodyParameter("User", "representation of a user").DataType("users.User")).
		Reads(User{}))

	ws.Route(ws.GET("/{user-id}").To(api.read).
		// Swagger documentation.
		Doc("read a user").
		Param(ws.PathParameter("user-id", "identifier for a user").DataType("string")).
		Writes(User{}))

	ws.Route(ws.PUT("/{user-id}").To(api.update).
		// Swagger documentation.
		Doc("update an existing user").
		Param(ws.PathParameter("user-id", "identifier for a user").DataType("string")).
		Param(ws.BodyParameter("User", "representation of a user").DataType("users.User")).
		Reads(User{}))

	ws.Route(ws.DELETE("/{user-id}").To(api.delete).
		// Swagger documentation.
		Doc("delete a user").
		Param(ws.PathParameter("user-id", "identifier for a user").DataType("string")))

	restful.Add(ws)
}

// Create a new resource.
//
func (api *UserApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	u := new(User)
	err := r.ReadEntity(&u)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Set some fields that need special handling.
	u.LastModified = time.Now()
	u.Status = StatusActive

	// The resource belongs to this user.
	u.UserId = user.Current(c).ID

	// Store the user.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "users", nil), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// The resource Id.
	u.Id = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", rootPath+"/"+k.Encode())

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	u.Link = rootPath + "/" + k.Encode()

	// Return the resultant entity.
	w.WriteHeader(http.StatusCreated)
	w.WriteEntity(u)
}

// Read the resource.
//
func (api UserApi) read(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("user-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the entity from the datastore.
	u := User{}
	if err := datastore.Get(c, k, &u); err != nil {
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
	//if u.UserId != user.Current(c).ID {
	//	http.Error(w, "You do not have access to this resource", http.StatusForbidden)
	//	return
	//}

	// Set their Id.
	u.Id = k.Encode()

	// Provide a link for ease of API usage.
	// TODO: This should be a fully qualified path.
	u.Link = rootPath + "/" + k.Encode()

	w.WriteEntity(u)
}

// Update the resource.
//
func (api *UserApi) update(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("user-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Marshall the entity from the request into a struct.
	u := new(User)
	err = r.ReadEntity(&u)
	if err != nil {
		w.WriteError(http.StatusNotAcceptable, err)
		return
	}

	// Retrieve the old entity from the datastore.
	old := User{}
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
	u.UserId = user.Current(c).ID

	// Keep track of the last modification date.
	u.LastModified = time.Now()

	// Attempt to overwrite the old entity.
	_, err = datastore.Put(c, k, u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Let them know it succeeded.
	w.WriteHeader(http.StatusNoContent)
}

// Delete the resource.
//
func (api *UserApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("user-id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the old entity from the datastore.
	old := User{}
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
