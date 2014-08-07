package accounts

import (
	"appengine"
	"appengine/datastore"
	// "appengine/user"
	"errors"
	"github.com/emicklei/go-restful"
	"net/http"
	"resource"
	"strconv"
	"time"
)

// Status values for records of this resource type.
const (
	StatusActivationPending = iota
	StatusActive
	StatusDeactivated
	StatusPermanentlyBanned
	StatusDeletionPending
	StatusDeleted
)

// Access level control. Access is ranked from lowest to highest to enable simple arithmetic checks
// when appropriate.
const (
	AccessLevelGeneral = iota
	AccessLevelFederationManager
	AccessLevelCommunityManager
	AccessLevelGameMaster
	AccessLevelAdmin
	AccessLevelSiteAdmin
)

const (
	Kind     = "accounts"
	RootPath = "/nexus8/" + Kind
)

func PreferredLink(k *datastore.Key) string {
	return RootPath + "/" + k.Encode()
}

type Api struct {
}

type Resource struct {
	// Primary email address for this account holder.
	Email string `json:"email" xml:"email"`

	// First name for this account holder.
	FirstName string `json:"first_name" xml:"first-name"`

	// Last name for this account holder.
	LastName string `json:"last_name" xml:"last-name"`

	// Gravatar for this account holder.
	AvatarUrl string `json:"avatar_url" xml:"avatar-url"`

	// Access level control.
	AccessLevel int `json:"access_level" xml:"access-level"`
}

type OpenIdInfo struct {
	// Primary email address for this account holder.
	OpenId string `json:"open_id" xml:"open-id"`

	// First name for this account holder.
	AuthDomain string `json:"auth_domain" xml:"auth-domain"`

	// Last name for this account holder.
	FederatedIdentity string `json:"federated_identity" xml:"federated-identity"`

	// Gravatar for this account holder.
	FederatedProvider string `json:"federated_provider" xml:"federated-provider"`
}

type ResourceMeta struct {
	Api
	resource.Meta
	Resource
	OpenIdInfo
}

type ResourceRequest struct {
	Api
	Resource
}

type ResourceResponse struct {
	Api
	resource.Meta
	Resource
	OpenIdInfo
}

// Register the routes we require for this resource type.
//
func Register() {
	ws := new(restful.WebService)

	ws.
		Path(RootPath).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		Doc("User account management.")

	ws.Route(ws.PUT("/{resource-id}").To(put).
		Doc("Update an existing resource").
		Operation("putAccount").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.BodyParameter("accounts.Resource", "representation of a resource").DataType("accounts.Resource")).
		Param(ws.HeaderParameter("If-Unmodified-Since", "Conditional modifier").DataType("RFC3339Nano Date")).
		Reads(ResourceRequest{}))

	ws.Route(ws.GET("/{resource-id}").To(get).
		Doc("Read a resource").
		Operation("getAccounts").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")).
		Writes(ResourceResponse{}))

	ws.Route(ws.HEAD("/{resource-id}").To(head).
		Doc("Returns the headers for a resource").
		Operation("headAccount").
		Param(ws.PathParameter("resource-id", "key for an existing resource").DataType("string")).
		Param(ws.HeaderParameter("If-Modified-Since", "Optional conditional modifier").DataType("RFC3339Nano Date")))

	restful.Add(ws)
}

// Attempts to create a valid key for a resource.
//
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

// Tests to see if the current user has an account with access greater than or equal to the accessLevel requested.
func AccessLevelGE(r *restful.Request, w *restful.Response, accessLevel int) error {
	// c := appengine.NewContext(r.Request)

	// var accs []ResourceMeta
	// q := datastore.NewQuery(Kind).
	// 	Filter("OpenId =", user.Current(c).ID).
	// 	Filter("AccessLevel >=", accessLevel).
	// 	Limit(1).
	// 	KeysOnly()

	// if keys, err := q.GetAll(c, accs); err != nil {
	// 	resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
	// 	return err
	// } else {

	// 	if len(keys) == 0 {
	// 		// No matching account found.
	// 		resource.WriteError(w, resource.NewError(http.StatusForbidden, "/html/error/statusforbidden", "User account does not appear to have the required permissions."))
	// 		return errors.New("User account does not appear to have the required permissions.")
	// 	}
	// }

	return nil
}

// Tests to see if the current user has an account with the requested accessLevel only.
func AccessLevelEQ(r *restful.Request, w *restful.Response, accessLevel int) error {
	// c := appengine.NewContext(r.Request)

	// var accs []ResourceMeta
	// q := datastore.NewQuery(Kind).
	// 	Filter("OpenId =", user.Current(c).ID).
	// 	Filter("AccessLevel =", accessLevel).
	// 	Limit(1).
	// 	KeysOnly()

	// if keys, err := q.GetAll(c, accs); err != nil {
	// 	resource.WriteError(w, resource.NewError(http.StatusInternalServerError, "/html/error/statusinternalservererror", err.Error()))
	// 	return err
	// } else {

	// 	if len(keys) == 0 {
	// 		// No matching account found.
	// 		resource.WriteError(w, resource.NewError(http.StatusForbidden, "/html/error/statusforbidden", "User account does not appear to have the required permissions."))
	// 		return errors.New("User account does not appear to have the required permissions.")
	// 	}
	// }

	return nil
}

// Update the resource.
//
func put(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

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

		// Fill in any old values not found in the current request.
		res.OpenId = old.OpenId

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

/// Read a resource
//
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
//
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
