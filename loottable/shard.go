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
	log.Printf("Registering " + kind + " shard services")
}

// Register the routes we require for this resource type.
//
func (api ResourceApi) RegisterShard() {
	ws := new(restful.WebService)

	ws.
		Path(shardRootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{resource-id}").To(api.get).
		// Swagger documentation.
		Doc("Read a resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(Resource{}))

	ws.Route(ws.HEAD("/{resource-id}").To(api.head).
		// Swagger documentation.
		Doc("Returns the headers for a resource").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")))

	ws.Route(ws.GET("/summary").To(api.listSummary).
		// Swagger documentation.
		Doc("Summary list of all resources").
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ListSummary{}))

	ws.Route(ws.GET("/all").To(api.listAll).
		// Swagger documentation.
		Doc("Comprehensive list of all resources").
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ListComprehensive{}))

	restful.Add(ws)
}

// Read a resource
//
func (api ResourceApi) get(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
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
				q := datastore.NewQuery(kind).
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
		resource := new(Resource)
		if err := datastore.Get(c, k, resource); err != nil {
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
		resource.Key = k.Encode()

		// Provide a link for ease of API usage.
		resource.Link.Rel = "self"
		resource.Link.Href = shardRootPath + "/" + k.Encode()

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, resource.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(resource.Revision))

		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
		w.AddHeader("Cache-Control", "max-age=3600,must-revalidate")

		// Output the response body.
		w.WriteEntity(resource)
	}
}

// Returns the headers for a resource
//
func (api ResourceApi) head(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
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
				q := datastore.NewQuery(kind).
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
		resource := new(Resource)
		if err := datastore.Get(c, k, resource); err != nil {
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
		resource.Key = k.Encode()

		// Provide a link for ease of API usage.
		resource.Link.Rel = "self"
		resource.Link.Href = shardRootPath + "/" + k.Encode()

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, resource.LastModified.Format(time.RFC3339Nano))
		w.AddHeader("ETag", strconv.Itoa(resource.Revision))

		// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
		w.AddHeader("Cache-Control", "max-age=3600,must-revalidate")

		// No response body required for this verb.
		w.WriteHeader(http.StatusNoContent)
		w.WriteEntity(nil)
	}
}

// Summary list of all resources
//
func (api ResourceApi) listSummary(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var q *datastore.Query
	var result ListSummary

	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
		q = datastore.NewQuery(kind).
			Project("LastModified", "Revision", "Status", "Name")
	} else {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		} else {
			q = datastore.NewQuery(kind).
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
			result.Entry[i].Link.Href = shardRootPath + "/" + k.Encode()
		}
	}

	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
	w.AddHeader("Cache-Control", "max-age=3600,must-revalidate")
	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

	w.WriteEntity(result)

}

// Comprehensive list of all resources
//
func (api ResourceApi) listAll(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var result ListComprehensive
	var q *datastore.Query

	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
		q = datastore.NewQuery(kind)
	} else {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		} else {
			q = datastore.NewQuery(kind).
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
			result.Entry[i].Link.Href = shardRootPath + "/" + k.Encode()
		}
	}

	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
	w.AddHeader("Cache-Control", "max-age=3600,must-revalidate")
	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

	w.WriteEntity(result)
}
