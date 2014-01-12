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
func (res *Resource) RegisterShard() {
	ws := new(restful.WebService)

	ws.
		Path(res.ShardRootPath()).
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

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
			result.Entry[i].Link.Href = result.Entry[i].PreferredLink(k)
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
			result.Entry[i].Link.Href = result.Entry[i].PreferredLink(k)
		}
	}

	// Cache Control: By allowing a short cache time here we can reduce database calls and cost.
	w.AddHeader("Cache-Control", "max-age=900,must-revalidate")
	w.AddHeader(restful.HEADER_LastModified, time.Now().Format(time.RFC3339Nano))

	w.WriteEntity(result)
}
