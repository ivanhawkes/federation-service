package characterappearance

import (
	"appengine"
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	//	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {
}

// Register the routes we require for this resource type.
//
func (res *Resource) RegisterAdmin() {
	ws := new(restful.WebService)

	ws.
		Path(res.AdminRootPath()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(res.post).
		// Swagger documentation.
		Doc("Create a new resource").
		Param(ws.BodyParameter("loottable.Resource", "representation of a resource").DataType("loottable.Resource")).
		Reads(*res).
		Writes(*res))

	ws.Route(ws.PUT("/{resource-id}").To(res.put).
		// Swagger documentation.
		Doc("Update an existing resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.BodyParameter("loottable.Resource", "representation of a resource").DataType("loottable.Resource")).
		Param(ws.HeaderParameter("If-Unmodified-Since", "Conditional modifier").DataType("RFC3339Nano Date")).
		Reads(*res))

	ws.Route(ws.DELETE("/{resource-id}").To(res.delete).
		// Swagger documentation.
		Doc("Delete an existing resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Unmodified-Since", "Conditional modifier").DataType("RFC3339Nano Date")))

	restful.Add(ws)
}

// Create a new resource.
//
func (res *Resource) post(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	err := r.ReadEntity(res)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusNotAcceptable, err.Error())
		return
	}

	// Set some fields that need special handling.
	res.LastModified = time.Now()
	res.Status = StatusActive
	res.Revision = 1

	// Store the resource.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, res.Kind(), nil), res)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	// The resource Key.
	res.Key = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", res.ShardRootPath()+k.Encode())

	// Provide a link for ease of API usage.
	res.Link.Rel = "self"
	res.Link.Href = res.ShardRootPath() + k.Encode()

	// Set the headers.
	w.WriteHeader(http.StatusCreated)
	w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
	w.AddHeader("ETag", strconv.Itoa(res.Revision))

	// Output the response body.
	w.WriteEntity(res)
}

// Update the resource.
//
func (res *Resource) put(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := res.getKey(r, w); err != nil {
		return
	} else {

		// Marshall the entity from the request into a struct.
		err = r.ReadEntity(res)
		if err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		}

		// Retrieve the old entity from the datastore.
		old := new(Resource)
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

		// Use conditional put - check LastModified before doing anything.
		if ifUnmodifiedSince := r.HeaderParameter("If-Unmodified-Since"); ifUnmodifiedSince == "" {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusForbidden, "Unconditional updates are not supported. Please provide 'If-Unmodified-Since' headers.")
			return
		} else {
			if t, err := time.Parse(time.RFC3339Nano, ifUnmodifiedSince); err != nil {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
				return
			} else {
				if t.Before(old.LastModified) {
					w.AddHeader("Content-Type", "text/plain")
					w.WriteErrorString(http.StatusPreconditionFailed, "The resource has been modified recently. Refresh your copy and try again if updating is still desireable.")
					return
				}
			}
		}

		// Keep track of the last modification date.
		res.LastModified = time.Now()
		res.Revision = old.Revision + 1

		// Attempt to overwrite the old entity.
		_, err = datastore.Put(c, k, res)
		if err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(res.Revision))

		// Let them know it succeeded.
		w.WriteHeader(http.StatusNoContent)
		w.WriteEntity(nil)
	}
}

// Delete the resource.
//
func (res *Resource) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := res.getKey(r, w); err != nil {
		return
	} else {

		// Retrieve the old entity from the datastore.
		old := new(Resource)
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

		// Use conditional delete - check LastModified before doing anything.
		if ifUnmodifiedSince := r.HeaderParameter("If-Unmodified-Since"); ifUnmodifiedSince == "" {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusForbidden, "Unconditional deletes are not supported. Please provide 'If-Unmodified-Since' headers.")
			return
		} else {
			if t, err := time.Parse(time.RFC3339Nano, ifUnmodifiedSince); err != nil {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
				return
			} else {
				if t.Before(old.LastModified) {
					w.AddHeader("Content-Type", "text/plain")
					w.WriteErrorString(http.StatusPreconditionFailed, "The resource has been modified recently. Refresh your copy and try again if deletion is still desireable.")
					return
				}
			}
		}

		// Delete the entity.
		if err := datastore.Delete(c, k); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		// Success notification.
		w.WriteHeader(http.StatusNoContent)
		w.WriteEntity(nil)
	}
}
