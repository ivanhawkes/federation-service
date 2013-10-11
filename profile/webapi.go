package profile

import (
//	"appengine"
//	"appengine/datastore"
//	"appengine/user"
//	"encoding/json"
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/profiles/", handleProfiles)
}

func serveError(w http.ResponseWriter, r *http.Request, status int) {
	// Let them know the error code.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	switch status {
	// 4xx
	case http.StatusBadRequest:
		io.WriteString(w, "Bad request\n")
	case http.StatusUnauthorized:
		io.WriteString(w, "Unauthorised\n")
	case http.StatusPaymentRequired:
		io.WriteString(w, "Payment required\n")
	case http.StatusForbidden:
		io.WriteString(w, "Forbidden\n")
	case http.StatusNotFound:
		io.WriteString(w, "Not found\n")
	case http.StatusMethodNotAllowed:
		io.WriteString(w, "Method not allowed\n")
	case http.StatusNotAcceptable:
		io.WriteString(w, "Not Acceptable\n")
	case http.StatusProxyAuthRequired:
		io.WriteString(w, "Proxy auth required\n")
	case http.StatusRequestTimeout:
		io.WriteString(w, "Request timeout\n")
	case http.StatusConflict:
		io.WriteString(w, "conflict\n")
	case http.StatusGone:
		io.WriteString(w, "Gone\n")
	case http.StatusLengthRequired:
		io.WriteString(w, "Length required\n")
	case http.StatusPreconditionFailed:
		io.WriteString(w, "Pre-condition failed\n")
	case http.StatusRequestEntityTooLarge:
		io.WriteString(w, "Request entity too large\n")
	case http.StatusRequestURITooLong:
		io.WriteString(w, "Request URI too long\n")
	case http.StatusUnsupportedMediaType:
		io.WriteString(w, "Unsupported media type\n")
	case http.StatusRequestedRangeNotSatisfiable:
		io.WriteString(w, "Rquested range not satisfiable\n")
	case http.StatusExpectationFailed:
		io.WriteString(w, "Expectation failed\n")
	case http.StatusTeapot:
		io.WriteString(w, "Teapot\n")

	// 5xx
	case http.StatusInternalServerError:
		io.WriteString(w, "Internal server error\n")
	case http.StatusNotImplemented:
		io.WriteString(w, "Not implemented\n")
	case http.StatusBadGateway:
		io.WriteString(w, "Bad gateway\n")
	case http.StatusServiceUnavailable:
		io.WriteString(w, "Service unavailable\n")
	case http.StatusGatewayTimeout:
		io.WriteString(w, "Gateway timeout\n")
	case http.StatusHTTPVersionNotSupported:
		io.WriteString(w, "Version not supported\n")

		// ???
	default:
		io.WriteString(w, "Unknown error\n")
	}
}

func handleProfiles(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "POST":
		handlePost(w, r)
	case "PUT":
		handlePut(w, r)
	case "DELETE":
		handleDelete(w, r)
	default:
		serveError(w, r, http.StatusNotFound)
		return
	}

	return
}


func handleGet(w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)
//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	io.WriteString(w, r.URL.Path) // <- the URL path

	io.WriteString(w, r.URL.Path)
	// Check if there is an Id
	io.WriteString(w, r.FormValue("id"))
	// NOTE!
	// We need to pull the Id from the URL immediately after the /profiles/ part.



	//io.WriteString(w, "get\n")
	return
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "muff\n")
	//c := appengine.NewContext(r)
	//err := PostJson(c, r.url, )

/*	if len(p.firstName) < 2 {
		serveError(w, r, http.StatusBadRequest)
		io.WriteString(w, "firstName is too short\n")
		return
	}
	if len(p.nickName) < 5 {
		serveError(w, r, http.StatusBadRequest)
		io.WriteString(w, "nickName is too short\n")
		return
	}
	if len(p.lastName) < 2 {
		serveError(w, r, http.StatusBadRequest)
		io.WriteString(w, "lastName is too short\n")
		return
	}*/

	//p := NewProfile(fn, nn, ln, user.Current(c).String())

	//k := datastore.NewIncompleteKey(c, "Profile", nil)
	//if _, err := datastore.Put(c, k, p); err != nil {
//		http.Error(w, err.Error(), 500)
//		return
	//}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

//	b, err := json.Marshal(p)
//	if err != nil {
//		io.WriteString(w, "error man!")
//	} else {
//		io.WriteString(w, string(b))
//	}

io.WriteString(w, "what up?")

	return
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, "put\n")
	return
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, "delete\n")
	return
}
