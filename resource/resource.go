package resource

import (
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"google.golang.org/appengine/datastore"

	"google.golang.org/appengine"
)

type Link struct {
	Rel  string `datastore:"-" json:"rel" xml:"rel"`
	Href string `datastore:"-" json:"href" xml:"href"`
}

type Versioning struct {
	LastModified time.Time `json:"last_modified" xml:"last-modified"`
	Revision     int       `json:"revision" xml:"revision"`
}

type Meta struct {
	Versioning
	Key    datastore.Key `datastore:"-" json:"key" xml:"key"`
	Status int           `json:"status" xml:"status"`
	Link   Link          `datastore:"-" json:"link" xml:"link"`
}

// Custom error reporting struct to streamline error reporting to the client.
type ResourceError struct {
	HttpStatus  int    `datastore:"-" json:"http_status" xml:"http-status"` // The HTTP status of the error.
	Uri         string `datastore:"-" json:"uri" xml:"uri"`                 // A Uri which can explain the nature of the error.
	Description string `datastore:"-" json:"description" xml:"description"` // A brief description of the error.
}

// Implement the Error interface.
func (resErr ResourceError) Error() string {
	return resErr.Description
}

// Write the error out to the response steam.
func WriteError(w *restful.Response, resErr *ResourceError) {
	w.WriteHeader(resErr.HttpStatus)
	w.WriteEntity(resErr)
}

// Convenience constructor for resource based errors.
func NewError(httpStatus int, uri string, description string) *ResourceError {
	resErr := new(ResourceError)
	resErr.HttpStatus = httpStatus
	resErr.Uri = uri
	resErr.Description = description

	return resErr
}

// Performs a check to see if the resource has been modified since a given datetime.
// Return value indicates if you should keep processing after calling this routine.
func IfModifiedSince(r *restful.Request, w *restful.Response, c appengine.Context, kind string, k *datastore.Key) (bool, error) {
	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince != "" {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			WriteError(w, NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
			return false, err
		} else {
			// Check the versioning information and see if we can tell them the
			// resource is unmodified.
			q := datastore.NewQuery(kind).
				Filter("__key__ =", k).
				Filter("LastModified >", t).
				KeysOnly()

			if keys, err := q.GetAll(c, nil); err != nil {
				WriteError(w, NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
				return false, err
			} else {

				if len(keys) == 0 {
					// Not modified.
					w.WriteHeader(http.StatusNotModified)
					w.WriteEntity(nil)
					return false, nil
				}
			}
		}
	}

	return true, nil
}
