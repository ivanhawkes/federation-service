package main

import (
	"appengine"
	"appengine/datastore"
	//	"appengine/memcache"
	"github.com/emicklei/go-restful"
	"net/http"
	"time"

//	"appengine/profile"
)

// This example is functionally the same as ../restful-profile-service.go
// but it`s supposed to run on Goole App Engine (GAE)

// Our simple example struct.
type Profile struct {
	LastModified time.Time `json:"last_modified"`
	//ApplicationId datastore.Key `json:"application_id"`
	//AccountId     datastore.Key `json:"account_id"`
	FirstName string `json:"first_name"`
	NickName  string `json:"nick_name"`
	LastName  string `json:"last_name"`
	// We might need an int to append to nickname to make it unique like in GW2 e.g. socks.451
}

type ProfileWebApi struct {
}

func init() {
	u := ProfileWebApi{}
	u.register()
}

func (u ProfileWebApi) register() {
	ws := new(restful.WebService)

	ws.
		Path("/profiles").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.POST("").To(u.create).
		// docs
		Doc("create a profile").
		Param(ws.BodyParameter("Profile", "representation of a profile").DataType("main.Profile")).
		//		Reads(Profile{})) // from the request
		Reads(datastore.Key{})) // from the request

	ws.Route(ws.GET("/{profile-id}").To(u.read).
		// docs
		Doc("get a profile").
		Param(ws.PathParameter("profile-id", "identifier of the profile").DataType("string")).
		Writes(Profile{})) // on the response

	//	ws.Route(ws.GET("").To(u.readAll).
	// docs
	//		Doc("get all profiles").
	//		Writes(Profile{})) // on the response

	ws.Route(ws.PUT("/{profile-id}").To(u.update).
		// docs
		Doc("update a profile").
		Param(ws.PathParameter("profile-id", "identifier of the profile").DataType("string")).
		Param(ws.BodyParameter("Profile", "representation of a profile").DataType("main.Profile")).
		Reads(Profile{})) // from the request

	ws.Route(ws.DELETE("/{profile-id}").To(u.delete).
		// docs
		Doc("delete a profile").
		Param(ws.PathParameter("profile-id", "identifier of the profile").DataType("string")))

	restful.Add(ws)
}

// POST http://localhost:8080/profiles
// {"first_name": "Ivan", "nick_name": "Socks", "last_name": "Hawkes"}
//
func (u *ProfileWebApi) create(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	p := new(Profile)
	err := r.ReadEntity(&p)
	if err == nil {
		// Tag the modified datetime.
		p.LastModified = time.Now()

		k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "profiles", nil), p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return value is the Id we stored the data under.
		w.WriteEntity(k)
	} else {
		w.WriteError(http.StatusNotAcceptable, err)
	}
}

// GET http://localhost:8080/profiles
// TODO: broken until I switch from memcache over to datastore
func (u ProfileWebApi) readAll(r *restful.Request, w *restful.Response) {
	//	c := appengine.NewContext(r.Request)
	//	id := r.PathParameter("profile-id")
	//	prof := new(Profile)
	//	_, err := memcache.Gob.Get(c, id, &prof)
	//	if err != nil || len(prof.Id) == 0 {
	//		w.WriteErrorString(http.StatusNotFound, "Profile could not be found.")
	//	} else {
	//		w.WriteEntity(prof)
	//	}
}

// GET http://localhost:8080/profiles/ahdkZXZ-ZmVkZXJhdGlvbi1zZXJ2aWNlc3IVCxIIcHJvZmlsZXMYgICAgICAiAsM
//
func (u ProfileWebApi) read(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)
	id := r.PathParameter("profile-id")
	k := datastore.NewKey(c, "profiles", id, 0, nil)
	//w.WriteEntity(k.StringID())
	//w.WriteEntity(id)
//	return

	p := new(Profile)
	if err := datastore.Get(c, k, &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteEntity("buut")
	return
	//		w.WriteEntity(id)

	//	prof := new(Profile)
	//	_, err := memcache.Gob.Get(c, id, &prof)
	//	if err != nil || len(prof.Id) == 0 {
	//		w.WriteErrorString(http.StatusNotFound, "Profile could not be found.")
	//	} else {
	//		w.WriteEntity(prof)
	//	}
}

// PUT http://localhost:8080/profiles/1
// <Profile><Id>1</Id><Name>Melissa Raspberry</Name></Profile>
//
func (u *ProfileWebApi) update(r *restful.Request, w *restful.Response) {
	//	c := appengine.NewContext(r.Request)
	//	prof := Profile{Id: r.PathParameter("profile-id")}
	//	err := r.ReadEntity(&prof)
	//	if err == nil {
	//		item := &memcache.Item{
	//			Key:    prof.Id,
	//			Object: &prof,
	//		}
	//		err = memcache.Gob.Add(c, item)
	//		if err != nil {
	//			w.WriteError(http.StatusInternalServerError, err)
	//			return
	//		}
	//		w.WriteHeader(http.StatusCreated)
	//		w.WriteEntity(prof)
	//	} else {
	//		w.WriteError(http.StatusInternalServerError, err)
	//	}
}

// DELETE http://localhost:8080/profiles/1
//
func (u *ProfileWebApi) delete(r *restful.Request, w *restful.Response) {
	//	c := appengine.NewContext(r.Request)
	//	id := r.PathParameter("profile-id")
	//	err := memcache.Delete(c, id)
	//	if err != nil {
	//		w.WriteError(http.StatusInternalServerError, err)
	//	}
}
