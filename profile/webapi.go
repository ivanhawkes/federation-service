package profile

import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"net/http"
	//"fmt"
	"time"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	p := Profile{
		createDate: time.Now(),
		firstName:     "Ivan",
		nickName:     "Socks",
		lastName:     "Hawkes",
		Account:  user.Current(c).String(),
	}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "employee", nil), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
