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
func (res *Resource) RegisterShard() {
	ws := new(restful.WebService)

	ws.
		Path(res.ShardRootPath()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{resource-id}").To(res.get).
		// Swagger documentation.
		Doc("Read a resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(*res))

	ws.Route(ws.HEAD("/{resource-id}").To(res.head).
		// Swagger documentation.
		Doc("Returns the headers for a resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")))

	ws.Route(ws.GET("/summary").To(res.listSummary).
		// Swagger documentation.
		Doc("Summary list of all resources").
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ListSummary{}))

	ws.Route(ws.GET("/all").To(res.listAll).
		// Swagger documentation.
		Doc("Comprehensive list of all resources").
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ListComprehensive{}))

	restful.Add(ws)
}

// Read a resource
//
func (res *Resource) get(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := res.getKey(r, w); err != nil {
		return
	} else {

		// Check if they want to limit the query using a modified since date.
		if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince != "" {
			if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
				return
			} else {
				// Check the versioning information and see if we can tell them the
				// resource is unmodified.
				q := datastore.NewQuery(res.Kind()).
					Filter("__key__ =", k).
					Filter("LastModified >", t).
					KeysOnly()

				if keys, err := q.GetAll(c, nil); err != nil {
					w.AddHeader("Content-Type", "text/plain")
					w.WriteErrorString(http.StatusInternalServerError, err.Error())
					return
				} else {

					if len(keys) == 0 {
						// Not modified.
						w.WriteHeader(http.StatusNotModified)
						w.WriteEntity(nil)
						return
					}
				}
			}
		}

		// Retrieve the entity from the datastore.
		if err := datastore.Get(c, k, res); err != nil {
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

		// Set their Key.
		res.Key = k.Encode()

		// Provide a link for ease of API usage.
		res.Link.Rel = "self"
		res.Link.Href = res.ShardRootPath() + k.Encode()

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(res.Revision))

		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
		w.AddHeader("Cache-Control", "max-age=14400,must-revalidate")

		// Output the response body.
		w.WriteEntity(res)
	}
}

// Returns the headers for a resource
//
func (res *Resource) head(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := res.getKey(r, w); err != nil {
		return
	} else {

		// Check if they want to limit the query using a modified since date.
		if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince != "" {
			if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
				w.AddHeader("Content-Type", "text/plain")
				w.WriteErrorString(http.StatusNotAcceptable, err.Error())
				return
			} else {
				// Check the versioning information and see if we can tell them the
				// resource is unmodified.
				q := datastore.NewQuery(res.Kind()).
					Filter("__key__ =", k).
					Filter("LastModified >", t).
					KeysOnly()

				if keys, err := q.GetAll(c, nil); err != nil {
					w.AddHeader("Content-Type", "text/plain")
					w.WriteErrorString(http.StatusInternalServerError, err.Error())
					return
				} else {

					if len(keys) == 0 {
						// Not modified.
						w.WriteHeader(http.StatusNotModified)
						w.WriteEntity(nil)
						return
					}
				}
			}
		}

		// Retrieve the entity from the datastore.
		if err := datastore.Get(c, k, res); err != nil {
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

		// Set their Key.
		res.Key = k.Encode()

		// Provide a link for ease of API usage.
		res.Link.Rel = "self"
		res.Link.Href = res.ShardRootPath() + k.Encode()

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(res.Revision))

		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
		w.AddHeader("Cache-Control", "max-age=14400,must-revalidate")

		// No response body required for this verb.
		w.WriteHeader(http.StatusNoContent)
		w.WriteEntity(nil)
	}
}

// Summary list of all resources
//
func (res *Resource) listSummary(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var q *datastore.Query
	var result ListSummary

	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
		q = datastore.NewQuery(res.Kind()).
			Project("LastModified", "Revision", "Status", "Name")
	} else {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		} else {
			q = datastore.NewQuery(res.Kind()).
				Project("LastModified", "Revision", "Status", "Name").
				Filter("LastModified >", t)
		}
	}

	if keys, err := q.GetAll(c, &result.Entry); err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	} else {
		for i, k := range keys {
			result.Entry[i].Key = k.Encode()
			result.Entry[i].Link.Rel = "self"
			result.Entry[i].Link.Href = result.Entry[i].ShardRootPath() + "/" + k.Encode()
		}
	}

	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
	w.AddHeader("Cache-Control", "max-age=900,must-revalidate")
	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

	w.WriteEntity(result)

}

// Comprehensive list of all resources
//
func (res *Resource) listAll(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var result ListComprehensive
	var q *datastore.Query

	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
		q = datastore.NewQuery(res.Kind())
	} else {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		} else {
			q = datastore.NewQuery(res.Kind()).
				Filter("LastModified >", t)
		}
	}

	if keys, err := q.GetAll(c, &result.Entry); err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	} else {
		for i, k := range keys {
			result.Entry[i].Key = k.Encode()
			result.Entry[i].Link.Rel = "self"
			result.Entry[i].Link.Href = result.Entry[i].ShardRootPath() + "/" + k.Encode()
		}
	}

	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
	w.AddHeader("Cache-Control", "max-age=900,must-revalidate")
	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

	w.WriteEntity(result)
}
