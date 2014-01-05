package loottable

import (
	"appengine"
	"appengine/datastore"
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
	"strconv"
	"time"
	//	"io"
	//	"bytes"
)

const (
	kind     = "loottable"
	rootPath = "/server/" + kind
)

// The various states for a loottable resource.
const (
	StatusActive = iota
	StatusDeactivated
	StatusPendingActivation
	StatusDeletionPending
	StatusDeleted
)

type Link struct {
	Rel  string `datastore:"-" json:"rel" xml:"rel"`
	Href string `datastore:"-" json:"href" xml:"href"`
}

type LootShallow struct {
	Key          string    `datastore:"-" json:"key" xml:"key"`
	LastModified time.Time `json:"last_modified" xml:"last-modified"`
	Version      int       `json:"version" xml:"version"`
	Status       int       `json:"status" xml:"status"`
	Name         string    `json:"name" xml:"name"`
	Link         Link      `datastore:"-" json:"link" xml:"link"`
}

type LootEntry struct {
	ItemId      int64   `datastore:"ItemId" json:"item_id" xml:"item-id"`
	Probability float32 `datastore:"Probability" json:"probability" xml:"probability"`
	Quantity    int16   `datastore:"Quantity" json:"quantity" xml:"quantity"`
}

type LootTable struct {
	LootShallow
	AllowPreload  bool        `json:"allow_preload" xml:"allow-preload"`
	Probabilities []LootEntry `json:"probabilities" xml:"probabilities"`
}

type LootSummary struct {
	LootTables []LootShallow `json:"loot_tables" xml:"loot-tables"`
}

type LootQuery struct {
	LootTables []LootTable `json:"loot_tables" xml:"loot-tables"`
}

type LootTableApi struct {
	Path string
}

func init() {
	log.Printf("LootTable: Register")
}

// Register the routes we require for this resource type.
//
func (api LootTableApi) Register() {
	ws := new(restful.WebService)

	ws.
		Path(rootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(api.post).
		// Swagger documentation.
		Doc("create a new loot table").
		Param(ws.BodyParameter("LootTable", "representation of a loottable").DataType("loottable.LootTable")).
		Reads(LootTable{}))

	ws.Route(ws.GET("/{loottable-id}").To(api.get).
		// Swagger documentation.
		Doc("read a loot table").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")).
		Writes(LootTable{}))

	ws.Route(ws.HEAD("/{loottable-id}").To(api.head).
		// Swagger documentation.
		Doc("return the document headers").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")))

	ws.Route(ws.PUT("/{loottable-id}").To(api.put).
		// Swagger documentation.
		Doc("update an existing loot table").
		Param(ws.PathParameter("loottable-id", "identifier for a loottable").DataType("string")).
		Param(ws.BodyParameter("LootTable", "representation of a loottable").DataType("loottable.LootTable")).
		Reads(LootTable{}))

	ws.Route(ws.DELETE("/{loottable-id}").To(api.delete).
		// Swagger documentation.
		Doc("delete a loot table").
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

// Attempts to create a valid key for a resource.
//
func (api LootTableApi) getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {

	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("loottable-id"))
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid.\n")
		return nil, err
	}

	// Check for shenanigans with the key.
	if k.Kind() != kind {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusBadRequest, "The key is not valid for this type of resource.\n")
		return nil, err
	}

	return k, nil
}

// Create a new resource.
//
func (api *LootTableApi) post(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	loottable := new(LootTable)
	err := r.ReadEntity(&loottable)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusNotAcceptable, err.Error())
		return
	}

	// Set some fields that need special handling.
	loottable.LastModified = time.Now()
	loottable.Status = StatusActive
	loottable.Version = 1

	// Store the loottable.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, kind, nil), loottable)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	// The resource Key.
	loottable.Key = k.Encode()

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", rootPath+"/"+k.Encode())

	// Provide a link for ease of API usage.
	loottable.Link.Rel = "self"
	loottable.Link.Href = rootPath + "/" + k.Encode()

	// Set the headers.
	w.WriteHeader(http.StatusCreated)
	w.AddHeader(restful.HEADER_LastModified, loottable.LastModified.String())
	w.AddHeader("ETag", strconv.Itoa(loottable.Version))

	// Output the response body.
	w.WriteEntity(loottable)
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
		loottable.Link.Href = rootPath + "/" + k.Encode()

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
		loottable.Link.Href = rootPath + "/" + k.Encode()

		// Only return the headers.
		w.AddHeader(restful.HEADER_LastModified, loottable.LastModified.String())
		w.AddHeader("ETag", strconv.Itoa(loottable.Version))

		// No response body required for this verb.
		w.WriteHeader(http.StatusNoContent)
	}
}

// Update the resource.
//
func (api *LootTableApi) put(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
		return
	} else {

		// Marshall the entity from the request into a struct.
		loottable := new(LootTable)
		err = r.ReadEntity(&loottable)
		if err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		}

		// Retrieve the old entity from the datastore.
		old := LootTable{}
		if err := datastore.Get(c, k, &old); err != nil {
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
		loottable.LastModified = time.Now()
		loottable.Version = old.Version + 1

		// Attempt to overwrite the old entity.
		_, err = datastore.Put(c, k, loottable)
		if err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}

		// Set the headers.
		w.AddHeader(restful.HEADER_LastModified, loottable.LastModified.String())
		w.AddHeader("ETag", strconv.Itoa(loottable.Version))

		// Let them know it succeeded.
		w.WriteHeader(http.StatusNoContent)
	}
}

// Delete the resource.
//
func (api *LootTableApi) delete(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := api.getKey(r, w); err != nil {
		return
	} else {

		// Retrieve the old entity from the datastore.
		old := LootTable{}
		if err := datastore.Get(c, k, &old); err != nil {
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
			result.LootTables[i].Link.Href = rootPath + "/" + k.Encode()
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
			result.LootTables[i].Link.Href = rootPath + "/" + k.Encode()
		}
	}

	w.WriteEntity(result)
}
