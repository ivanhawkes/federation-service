package verbs

import (
	"federation-services/accounts"
	"federation-services/resource"
	"log"
	"net/http"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful/v3"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type RegisterResource struct {
	Email string `json:"email" xml:"email"`
}

type RegisterRequest struct {
	RegisterResource
}

type RegisterResponse struct {
	resource.Meta
	RegisterResource
}

func init() {
}

// Register the verb.
func RegisterVerbRegister() {
	log.Printf("VERB: /register")

	ws := new(restful.WebService)

	ws.
		Path("/register").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.POST("").To(postRegister).
		Doc("Register a new account").
		Operation("verbRegister").
		Param(ws.BodyParameter("Accounts.Resource", "representation of a resource").DataType("Accounts.Resource")).
		Reads(RegisterRequest{}).
		Writes(RegisterResponse{}))

	restful.Add(ws)
}

// Create a new resource.
func postRegister(r *restful.Request, w *restful.Response) {
	c := appengine.NewContext(r.Request)

	// Marshall the entity from the request into a struct.
	res := new(RegisterRequest)
	err := r.ReadEntity(res)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusNotAcceptable, err.Error())
		return
	}

	// Create a new account based on the limited information from registration.
	acc := new(accounts.ResourceMeta)

	// Marshall the entity from the request into a struct.
	resp := new(RegisterResponse)
	err = r.ReadEntity(resp)
	if err != nil {
		resource.WriteError(w, resource.NewError(http.StatusNotAcceptable, "/html/error/statusnotacceptable", err.Error()))
		return
	}

	// TODO: don't allow them to register the same email or ID twice

	// Set the meta data for the resource.
	resp.LastModified = time.Now()
	resp.Status = accounts.StatusActive
	resp.Revision = 1

	// TODO: we want the email address and ID to come from OAuth

	// Grab their email address.
	//	resp.Email = user.Current(c).Email

	// Copy over the information we have.
	acc.LastModified = resp.LastModified
	acc.Status = resp.Status
	acc.Revision = resp.Revision
	acc.AccessLevel = accounts.AccessLevelGeneral

	// Snag the information provided by OpenID  OAuth.
	acc.Email = resp.Email
	// acc.OpenId = user.Current(c).ID
	// acc.AuthDomain = user.Current(c).AuthDomain
	// acc.FederatedIdentity = user.Current(c).FederatedIdentity
	// acc.FederatedProvider = user.Current(c).FederatedProvider

	// If we don't have a valid email address for them, mark the account as not yet active.
	// TODO: need to send them an email validation - work out where in the process to add this.
	if len(acc.Email) == 0 {
		acc.Status = accounts.StatusActivationPending
	}

	// Store the resource.
	k, err := datastore.Put(c, datastore.NewIncompleteKey(c, "accounts", nil), acc)
	if err != nil {
		w.AddHeader("Content-Type", "text/plain")
		w.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	// The resource Key.
	resp.Key = *k

	// Let them know the location of the newly created resource.
	// TODO: Use a safe Url path append function.
	w.AddHeader("Location", accounts.PreferredLink(k))

	// Provide a link for ease of API usage.
	resp.Link.Rel = "self"
	resp.Link.Href = accounts.PreferredLink(k)

	// Set the headers.
	w.WriteHeader(http.StatusCreated)
	w.AddHeader(restful.HEADER_LastModified, resp.LastModified.Format(time.RFC3339Nano))
	w.AddHeader("ETag", strconv.Itoa(resp.Revision))

	// Output the response body.
	w.WriteEntity(resp)
}
