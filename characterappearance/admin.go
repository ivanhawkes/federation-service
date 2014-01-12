package characterappearance

import (
	"appengine"
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	//	"log"
	"net/http"
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

	ws.Route(ws.DELETE("/{resource-id}").To(res.delete).
		// Swagger documentation.
		Doc("Delete an existing resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Unmodified-Since", "Conditional modifier").DataType("RFC3339Nano Date")))

	restful.Add(ws)
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
