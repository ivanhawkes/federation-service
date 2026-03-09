package storagecontainers

import (
	"errors"
	"federation-services/accounts"
	"federation-services/resource"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"

	"google.golang.org/appengine/datastore"

	"google.golang.org/appengine"
)

// The various states for a federation resource.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusDeletionPending
	StatusDeleted
)

const (
	Kind     = "storagecontainers"
	RootPath = "/federation/" + Kind
)

func PreferredLink(k *datastore.Key) string {
	return RootPath + "/" + k.Encode()
}

type Api struct {
}

type Resource struct {
	// The owner key is a key belonging to one of these kinds - Character, Faction, Profile, Account
	OwnerKey datastore.Key `json:"owner_key" xml:"owner-key"`

	// The general class of storage for this container e.g.
	// "Character", "Character Bank", "Mail", "Faction Bank", "Quest", "Resource"
	Class string `json:"class" xml:"class"`

	// The specific name of the container e.g.
	// "Character 01", "Character Bank 01", Mail", "Resource.Ore", "Resource.Herb"
	Name string `json:"name" xml:"name"`

	// The maximum number of slots available within this container. The client is responsible for managing this figure and
	// ensuring they don't add more items than there are slots.
	SlotMax int32 `json:"slot_max" xml:"slot_max"`
}

type ResourceMeta struct {
	Api
	resource.Meta
	Resource
}

type ResourceKey struct {
	Api
	Key datastore.Key `datastore:"-" json:"key" xml:"key"`
	Resource
}

type ResourceRequest struct {
	Api
	Resource
}

type ResourceResponse struct {
	ResourceMeta
}

type ListResource struct {
	Entry []Resource `json:"entry" xml:"entry"`
}

type ListResourceKey struct {
	Entry []ResourceKey `json:"entry" xml:"entry"`
}

type ListResourceMeta struct {
	Entry []ResourceMeta `json:"entry" xml:"entry"`
}

// Register the routes we require for this resource type.
func Register() {
	ws := new(restful.WebService)

	ws.
		Path(RootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("StorageContainer management.")

	ws.Route(ws.POST("").To(post).
		Doc("Create a new resource").
		Operation("postStorageContainer").
		Param(ws.BodyParameter("StorageContainer.Resource", "representation of a resource").DataType("StorageContainer.Resource")).
		Reads(ResourceRequest{}).
		Writes(ResourceResponse{}))

	ws.Route(ws.PUT("/{resource-id}").To(put).
		Doc("Update an existing resource").
		Operation("putStorageContainer").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.BodyParameter("StorageContainer.Resource", "representation of a resource").DataType("StorageContainer.Resource")).
		Param(ws.HeaderParameter("If-Unmodified-Since", "Conditional modifier").DataType("RFC3339Nano Date")).
		Reads(ResourceRequest{}))

	ws.Route(ws.GET("/{resource-id}").To(get).
		Doc("Read a resource").
		Operation("getStorageContainer").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ResourceResponse{}))

	ws.Route(ws.HEAD("/{resource-id}").To(head).
		Doc("Returns the headers for a resource").
		Operation("headStorageContainer").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")))

	ws.Route(ws.GET("/list").To(listAll).
		Doc("Get a list of resources").
		Operation("listStorageContainer").
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ResourceResponse{}))

	restful.Add(ws)
}

// Attempts to create a valid key for a resource.
func getKey(r *restful.Request, w *restful.Response) (*datastore.Key, error) {
	// Decode the request parameter to determine the key for the entity.
	k, err := datastore.DecodeKey(r.PathParameter("resource-id"))
	if err != nil {
		resource.WriteError(w, resource.NewError(http.StatusBadRequest, "/html/error/invalidkey", "The key is not valid."))
		return nil, err
	}

	// Check for shenanigans with the key.
	if k.Kind() != Kind {
		resource.WriteError(w, resource.NewError(http.StatusBadRequest, "/html/error/invalidkeymalicious", "The key is not valid for this type of resource."))
		return nil, errors.New("The key is not valid for this type of resource.")
	}

	return k, nil
}

// Create a new resource.
func post(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Auth check.
	if err := accounts.AccessLevelGE(r, w, accounts.AccessLevelAdmin); err != nil {
		return
	}

	// Marshall the entity from the request into a struct.
	res := new(ResourceMeta)
	err := r.ReadEntity(res)
	if err != nil {
		resource.WriteError(w, resource.NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
		return
	}

	// Set the meta data for the resource.
	res.LastModified = time.Now()
	res.Status = StatusActive
	res.Revision = 1

	// TODO: The OwnerKey must be set at this point!

	// Store the resource.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, Kind, nil), res)
	if err != nil {
		resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
		return
	}

	// The resource Key.
	res.Key = *k

	// Let them know the location of the newly created resource.
	w.AddHeader("Location", PreferredLink(k))

	// Provide a link for ease of API usage.
	res.Link.Rel = "self"
	res.Link.Href = PreferredLink(k)

	// Set the headers.
	w.WriteHeader(http.StatusCreated)
	w.AddHeader(restful.HEADER_LastModified, res.LastModified.Format(time.RFC3339Nano))
	w.AddHeader("ETag", strconv.Itoa(res.Revision))

	// Output the response body.
	w.WriteEntity(res)
}

// Update the resource.
func put(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Auth check.
	if err := accounts.AccessLevelGE(r, w, accounts.AccessLevelAdmin); err != nil {
		return
	}

	// Grab the key and validate it.
	if k, err := getKey(r, w); err != nil {
		return
	} else {

		// Marshall the entity from the request into a struct.
		res := new(ResourceResponse)
		err = r.ReadEntity(res)
		if err != nil {
			resource.WriteError(w, resource.NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
			return
		}

		// Retrieve the old entity from the datastore.
		old := new(ResourceResponse)
		if err := datastore.Get(c, k, old); err != nil {
			if err.Error() == "datastore: no such entity" {
				resource.WriteError(w, resource.NewError(http.StatusNotFound, "/html/error/statusnotfound", err.Error()))
				return
			} else {
				resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
				return
			}
			return
		}

		// Use conditional put - check LastModified before doing anything.
		if ifUnmodifiedSince := r.HeaderParameter("If-Unmodified-Since"); ifUnmodifiedSince == "" {
			resource.WriteError(w, resource.NewError(http.StatusForbidden, "/html/error/statusforbidden", "Unconditional updates are not supported. Please provide 'If-Unmodified-Since' headers."))
			return
		} else {
			if t, err := time.Parse(time.RFC3339Nano, ifUnmodifiedSince); err != nil {
				resource.WriteError(w, resource.NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
				return
			} else {
				if t.Before(old.LastModified) {
					resource.WriteError(w, resource.NewError(http.StatusPreconditionFailed, "/html/error/statuspreconditionfailed", "The resource has been modified recently. Refresh your copy and try again if updating is still desireable."))
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
			resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
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

// Read a resource
func get(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := getKey(r, w); err != nil {
		return
	} else {
		// Handle cases where they want to check if the resource has been modified since...
		if ok, err := resource.IfModifiedSince(r, w, c, Kind, k); err != nil || ok != true {
			return
		}

		// Retrieve the entity from the datastore.
		res := new(ResourceResponse)
		if err := datastore.Get(c, k, res); err != nil {
			if err.Error() == "datastore: no such entity" {
				resource.WriteError(w, resource.NewError(http.StatusNotFound, "/html/error/statusnotfound", err.Error()))
				return
			} else {
				resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
				return
			}
			return
		}

		// Set their Key.
		res.Key = *k

		// Provide a link for ease of API usage.
		res.Link.Rel = "self"
		res.Link.Href = PreferredLink(k)

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
func head(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Grab the key and validate it.
	if k, err := getKey(r, w); err != nil {
		return
	} else {
		// Handle cases where they want to check if the resource has been modified since...
		if ok, err := resource.IfModifiedSince(r, w, c, Kind, k); err != nil || ok != true {
			return
		}

		// Retrieve the entity from the datastore.
		res := new(ResourceResponse)
		if err := datastore.Get(c, k, res); err != nil {
			if err.Error() == "datastore: no such entity" {
				resource.WriteError(w, resource.NewError(http.StatusNotFound, "/html/error/statusnotfound", err.Error()))
				return
			} else {
				resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
				return
			}
			return
		}

		// Set their Key.
		res.Key = *k

		// Provide a link for ease of API usage.
		res.Link.Rel = "self"
		res.Link.Href = PreferredLink(k)

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

// Read a resource
func listAll(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	var result ListResourceKey
	var q *datastore.Query

	// Check if they want to limit the query using a modified since date.
	if ifModifiedSince := r.HeaderParameter("If-Modified-Since"); ifModifiedSince == "" {
		q = datastore.NewQuery(Kind).
			Project("Category", "Subcategory", "Description").
			Filter("Status =", StatusActive)
	} else {
		if t, err := time.Parse(time.RFC3339Nano, ifModifiedSince); err != nil {
			w.AddHeader("Content-Type", "text/plain")
			w.WriteErrorString(http.StatusNotAcceptable, err.Error())
			return
		} else {
			q = datastore.NewQuery(Kind).
				Project("Category", "Subcategory", "Description").
				Filter("Status =", StatusActive).
				Filter("LastModified >", t)
		}
	}

	if keys, err := q.GetAll(c, &result.Entry); err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	} else {
		for i, k := range keys {
			// TODO: why does XML not emit the key correctly (JSON does)?
			result.Entry[i].Key = *k
		}
	}

	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
	//	w.AddHeader("Cache-Control", "max-age=900,must-revalidate")
	//	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

	w.WriteEntity(result)
}
