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
	log.Printf("Registering " + kind + " server services")
}

// Register the routes we require for this resource type.
//
func (api LootTableApi) RegisterServer() {
	ws := new(restful.WebService)

	ws.
		Path(serverRootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/{loottable-id}").To(api.get).
		// Swagger documentation.
		Doc("read a loot table").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")).
		Writes(LootTable{}))

	ws.Route(ws.HEAD("/{loottable-id}").To(api.head).
		// Swagger documentation.
		Doc("return the document headers").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")))

	ws.Route(ws.GET("/summary").To(api.summary).
		// Swagger documentation.
		Doc("returns a summary of all the loot tables").
		//Writes(LootSummary{}))
		Writes([]LootShallow{}))

	ws.Route(ws.GET("/all").To(api.all).
		// Swagger documentation.
		Doc("returns a complete listing of all the loot tables").
		Writes(LootQuery{}))

	restful.Add(ws)
}

// Get a representation of the resource from our datastore.
//
func (api LootTableApi) get(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
		return
	} else {
		// Retrieve the entity from the datastore.
		loottable := LootTable{}
		if err := datastore.Get(c, k, &loottable); err != nil {
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
		loottable.Key = k.Encode()

		// Provide a link for ease of API usage.
		loottable.Link.Rel = "self"
		loottable.Link.Href = serverRootPath + "/" + k.Encode()

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, loottable.LastModified.String())
		w.AddHeader("ETag", strconv.Itoa(loottable.Version))

		// Output the response body.
		w.WriteEntity(loottable)
	}
}

// Return the headers for a resource.
//
func (api LootTableApi) head(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
		return
	} else {
		// Retrieve the entity from the datastore.
		loottable := LootTable{}
		if err := datastore.Get(c, k, &loottable); err != nil {
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
		loottable.Key = k.Encode()

		// Provide a link for ease of API usage.
		loottable.Link.Rel = "self"
		loottable.Link.Href = serverRootPath + "/" + k.Encode()

		// Only return the headers.
		w.AddHeader(restful.HEADER_LastModified, loottable.LastModified.String())
		w.AddHeader("ETag", strconv.Itoa(loottable.Version))

		// No response body required for this verb.
		w.WriteHeader(http.StatusNoContent)
	}
}

// Retrieve a summary of all the loot tables.
//
func (api LootTableApi) summary(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var q *datastore.Query
	var result LootSummary

	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
		q = datastore.NewQuery(kind).
			Project("LastModified", "Version", "Status", "Name")
	} else {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		} else {
			q = datastore.NewQuery(kind).
				Project("LastModified", "Version", "Status", "Name").
				Filter("LastModified >=", t)
		}
	}

	if keys, err := q.GetAll(c, &result.LootTables); err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	} else {
		for i, k := range keys {
			result.LootTables[i].Key = k.Encode()
			result.LootTables[i].Link.Rel = "self"
			result.LootTables[i].Link.Href = serverRootPath + "/" + k.Encode()
		}
	}

	w.WriteEntity(result)

}

// Retrieve a summary of all the loot tables.
//
func (api LootTableApi) all(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var result LootQuery
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
				Filter("LastModified >=", t)
		}
	}

	if keys, err := q.GetAll(c, &result.LootTables); err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	} else {
		for i, k := range keys {
			result.LootTables[i].Key = k.Encode()
			result.LootTables[i].Link.Rel = "self"
			result.LootTables[i].Link.Href = serverRootPath + "/" + k.Encode()
		}
	}

	w.WriteEntity(result)
}
