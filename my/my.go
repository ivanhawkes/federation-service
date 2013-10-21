package my

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"github.com/emicklei/go-restful"
	"net/http"
	"profiles"
	"log"
)

type MyApi struct {
	Path string
}

func init() {
    log.Printf("My: Register")
	api := MyApi{Path: "/my"}
	api.Register()
}

// Register the routes we require for this resource type.
//
func (api MyApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path(api.Path).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/profiles").To(api.profiles).
		// Swagger documentation.
		Doc("return a list of all my profiles").
		Writes(profiles.Profile{}))

	restful.Add(ws)
}

// GET http://localhost:8080/my/profiles
//
func (api MyApi) profiles(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Find a list of all the profiles for the current user.
	// TODO: Make this an ancestor query.
	q := datastore.NewQuery("profiles").Filter("UserId =", user.Current(c).ID)

	//
	var profs []profiles.Profile
	var keys []*datastore.Key
	var err error
	if keys, err = q.GetAll(c, &profs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: we need to make sure the results are placed in a container to prevent XML / JSON errors.
	// Provide an Id for every result returned.
	for i, _ := range keys {
		profs[i].Id = keys[i].Encode()
	}

	// Return the results.
	w.WriteEntity(profs)
}
