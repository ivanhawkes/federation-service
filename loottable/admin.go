package loottable

import (
	"appengine"
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {
	log.Printf("Registering " + kind + " admin services")
}

// Register the routes we require for this resource type.
//
func (api ResourceApi) RegisterAdmin() {
	ws := new(restful.WebService)

	ws.
		Path(adminRootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.post).
		// Swagger documentation.
		Doc("Create a new resource").
		Param(ws.BodyParameter("loottable.Resource", "representation of a resource").DataType("loottable.Resource")).
		Reads(Resource{}).
		Writes(Resource{}))

	ws.Route(ws.PUT("/{resource-id}").To(api.put).
		// Swagger documentation.
		Doc("Update an existing resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.BodyParameter("loottable.Resource", "representation of a resource").DataType("loottable.Resource")).
		Reads(Resource{}))

	ws.Route(ws.DELETE("/{resource-id}").To(api.delete).
		// Swagger documentation.
		Doc("Delete an existing resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")))

	restful.Add(ws)
}

// Create a new resource.
//
func (api *ResourceApi) post(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	resource := new(Resource)
	err := r.ReadEntity(&resource)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusNotAcceptable, err.Error())
		return
	}

	// Set some fields that need special handling.
	resource.LastModified = time.Now()
	resource.Status = StatusActive
	resource.Version = 1

	// Store the resource.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, kind, nil), resource)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	// The resource Key.
	resource.Key = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", serverRootPath+"/"+k.Encode())

	// Provide a link for ease of API usage.
	resource.Link.Rel = "self"
	resource.Link.Href = serverRootPath + "/" + k.Encode()

	// Set the headers.
	w.WriteHeader(http.StatusCreated)
	w.AddHeader(restful.HEADER_LastModified, resource.LastModified.String())
	w.AddHeader("ETag", strconv.Itoa(resource.Version))

	// Output the response body.
	w.WriteEntity(resource)
}

// Update the resource.
//
func (api *ResourceApi) put(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
		return
	} else {

		// Marshall the entity from the request into a struct.
		resource := new(Resource)
		err = r.ReadEntity(&resource)
		if err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		}

		// Retrieve the old entity from the datastore.
		old := new (Resource)
		if err := datastore.Get(c, k, old); err != nil {
			if err.Error() == "datastore: no such entity" {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusNotFound, err.Error())
				return
			} else {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusInternalServerError, err.Error())
				return
			}
			return
		}

		// Keep track of the last modification date.
		resource.LastModified = time.Now()
		resource.Version = old.Version + 1

		// Attempt to overwrite the old entity.
		_, err = datastore.Put(c, k, resource)
		if err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, resource.LastModified.String())
		w.AddHeader("ETag", strconv.Itoa(resource.Version))

		// Let them know it succeeded.
		w.WriteHeader(http.StatusNoContent)
	}
}

// Delete the resource.
//
func (api *ResourceApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
		return
	} else {

		// Retrieve the old entity from the datastore.
		old := new (Resource)
		if err := datastore.Get(c, k, old); err != nil {
			if err.Error() == "datastore: no such entity" {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusNotFound, err.Error())
				return
			} else {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusInternalServerError, err.Error())
				return
			}
			return
		}

		// Delete the entity.
		if err := datastore.Delete(c, k); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		// Success notification.
		w.WriteHeader(http.StatusNoContent)
	}
}
