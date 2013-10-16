package main

import (
	"appengine"
	"appengine/memcache"
	"github.com/emicklei/go-restful"
	"net/http"
)

// This example is functionally the same as ../restful-user-service.go
// but it`s supposed to run on Goole App Engine (GAE)

// Our simple example struct.
type User struct {
	Id, Name string
}

type UserWebApi struct {
}

func init() {
	u := UserWebApi{}
	u.register()
}

func (u UserWebApi) register() {
	ws := new(restful.WebService)

	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("").To(u.readAll).
		// docs
		Doc("get all users").
		//Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.GET("/{user-id}").To(u.read).
		// docs
		Doc("get a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(User{})) // on the response

	ws.Route(ws.POST("").To(u.update).
		// docs
		Doc("create a user").
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{})) // from the request

	ws.Route(ws.PUT("/{user-id}").To(u.create).
		// docs
		Doc("update a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Param(ws.BodyParameter("User", "representation of a user").DataType("main.User")).
		Reads(User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.delete).
		// docs
		Doc("delete a user").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	restful.Add(ws)
}

// GET http://localhost:8080/users
// TODO: broken until I switch from memcache over to datastore
func (u UserWebApi) readAll(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	id := request.PathParameter("user-id")
	usr := new(User)
	_, err := memcache.Gob.Get(c, id, &usr)
	if err != nil || len(usr.Id) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// GET http://localhost:8080/users/1
//
func (u UserWebApi) read(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	id := request.PathParameter("user-id")
	usr := new(User)
	_, err := memcache.Gob.Get(c, id, &usr)
	if err != nil || len(usr.Id) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
//
func (u *UserWebApi) update(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err == nil {
		item := &memcache.Item{
			Key:    usr.Id,
			Object: &usr,
		}
		err = memcache.Gob.Set(c, item)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// POST http://localhost:8080/users
// <User><Id>1</Id><Name>Melissa</Name></User>
//
func (u *UserWebApi) create(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	usr := User{Id: request.PathParameter("user-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		item := &memcache.Item{
			Key:    usr.Id,
			Object: &usr,
		}
		err = memcache.Gob.Add(c, item)
		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
//
func (u *UserWebApi) delete(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	id := request.PathParameter("user-id")
	err := memcache.Delete(c, id)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	}
}
